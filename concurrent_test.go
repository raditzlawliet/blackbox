package blackbox

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestConcurrentWrapper_NoDataLoss verifies that wrapping a concrete box with
func TestConcurrentWrapperNoDataLoss(t *testing.T) {
	producers := 4
	itemsPerProducer := 500
	total := producers * itemsPerProducer

	// Use a concrete FIFO and wrap it for concurrency.
	fifo := NewFIFO[int](0, 64)
	box := NewConcurrent[int](fifo)

	var pwg sync.WaitGroup
	var cwg sync.WaitGroup

	errs := make(chan error)

	// map to record consumed items
	seen := make(map[int]int)
	var seenMu sync.Mutex

	// producers
	pwg.Add(producers)
	for p := 0; p < producers; p++ {
		pid := p
		go func() {
			defer pwg.Done()
			base := pid * itemsPerProducer
			for i := 0; i < itemsPerProducer; i++ {
				val := base + i
				if err := box.Put(val); err != nil {
					errs <- fmt.Errorf("Put returned error: %v", err)
				}
				// a tiny sleep to increase interleaving
				time.Sleep(time.Microsecond)
			}
		}()
	}

	// consumers
	var consumed int64
	consumers := 6
	cwg.Add(consumers)
	for i := 0; i < consumers; i++ {
		go func() {
			defer cwg.Done()
			for {
				item, err := box.Get()
				if err == ErrEmptyBlackBox {
					// If producers are done and we've consumed all, exit.
					if atomic.LoadInt64(&consumed) >= int64(total) {
						return
					}
					// Otherwise wait briefly and retry.
					time.Sleep(time.Microsecond * 50)
					continue
				}
				atomic.AddInt64(&consumed, 1)
				seenMu.Lock()
				seen[item]++
				seenMu.Unlock()

				// Quick exit if done
				if atomic.LoadInt64(&consumed) >= int64(total) {
					return
				}
			}
		}()
	}

	pwg.Wait()
	cwg.Wait()

	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err) // fail immediately if any goroutine failed
		}
	}

	// Validate
	if int(atomic.LoadInt64(&consumed)) != total {
		t.Fatalf("expected %d consumed items, got %d", total, consumed)
	}

	// Check seen map has exactly the produced items and each exactly once
	if len(seen) != total {
		t.Fatalf("expected seen map size %d, got %d", total, len(seen))
	}
	for p := 0; p < producers; p++ {
		base := p * itemsPerProducer
		for i := 0; i < itemsPerProducer; i++ {
			v := base + i
			if cnt, ok := seen[v]; !ok {
				t.Fatalf("missing item %d", v)
			} else if cnt != 1 {
				t.Fatalf("item %d seen %d times (want 1)", v, cnt)
			}
		}
	}
}

// TestConcurrentWrapper_StrategiesBasic ensures that wrapping boxes is working with different strategies
func TestConcurrentWrapper_StrategiesBasic(t *testing.T) {
	const producers = 3
	const itemsPerProducer = 200
	const total = producers * itemsPerProducer

	strategies := []Strategy{StrategyFIFO, StrategyLIFO, StrategyRandom}

	for _, strat := range strategies {
		// Create a box via New with the given strategy and wrap it
		box := New[int](WithStrategy(strat), WithInitialCapacity(64))
		cbox := NewConcurrent[int](box)

		var pwg sync.WaitGroup
		var cwg sync.WaitGroup

		errs := make(chan error)

		var consumed int64
		pwg.Add(producers)

		for p := 0; p < producers; p++ {
			pid := p
			go func() {
				defer pwg.Done()
				base := pid * itemsPerProducer
				for i := 0; i < itemsPerProducer; i++ {
					val := base + i
					if err := cbox.Put(val); err != nil {
						errs <- fmt.Errorf("Put error for strategy %v: %v", strat, err)
					}
				}
			}()
		}

		// Consumers
		consumers := 4
		cwg.Add(consumers)
		for i := 0; i < consumers; i++ {
			go func() {
				defer cwg.Done()
				for {
					_, err := cbox.Get()
					if err == ErrEmptyBlackBox {
						if atomic.LoadInt64(&consumed) >= int64(total) {
							return
						}
						time.Sleep(time.Microsecond * 10)
						continue
					}
					atomic.AddInt64(&consumed, 1)
					if atomic.LoadInt64(&consumed) >= int64(total) {
						return
					}
				}
			}()
		}

		pwg.Wait()
		cwg.Wait()

		close(errs)
		for err := range errs {
			if err != nil {
				t.Fatal(err) // fail immediately if any goroutine failed
			}
		}

		if int(atomic.LoadInt64(&consumed)) != total {
			t.Fatalf("strategy %v: expected %d consumed, got %d", strat, total, consumed)
		}
	}
}

func TestConcurrentWrapper_AccessorsAndClean(t *testing.T) {
	fifo := NewFIFO[int](3, 3)
	box := NewConcurrent[int](fifo)

	// Initially empty
	if !box.IsEmpty() {
		t.Fatalf("expected empty box")
	}
	if got := box.Size(); got != 0 {
		t.Fatalf("expected size 0, got %d", got)
	}
	if got := box.MaxSize(); got != 3 {
		t.Fatalf("expected max size 3, got %d", got)
	}

	if _, err := box.Peek(); err != ErrEmptyBlackBox {
		t.Fatalf("expected ErrEmptyBlackBox on Peek, got %v", err)
	}

	for i := 0; i < 3; i++ {
		if err := box.Put(i); err != nil {
			t.Fatalf("Put returned unexpected error: %v", err)
		}
	}

	if !box.IsFull() {
		t.Fatalf("expected box to be full")
	}
	if got := box.Size(); got != 3 {
		t.Fatalf("expected size 3, got %d", got)
	}

	v, err := box.Peek()
	if err != nil {
		t.Fatalf("Peek returned unexpected error: %v", err)
	}
	if v != 0 {
		t.Fatalf("Peek returned %v, want 0", v)
	}
	if got := box.Size(); got != 3 {
		t.Fatalf("Peek should not remove item; size expected 3, got %d", got)
	}

	if err := box.Put(99); err != ErrBlackBoxFull {
		t.Fatalf("expected ErrBlackBoxFull on Put, got %v", err)
	}

	box.Clean()
	if !box.IsEmpty() {
		t.Fatalf("expected empty after Clean")
	}
	if got := box.Size(); got != 0 {
		t.Fatalf("expected size 0 after Clean, got %d", got)
	}
	if _, err := box.Get(); err != ErrEmptyBlackBox {
		t.Fatalf("expected ErrEmptyBlackBox after Clean on Get, got %v", err)
	}

	if got := box.MaxSize(); got != 3 {
		t.Fatalf("expected max size to remain 3 after Clean, got %d", got)
	}
}

func benchmarkConcurrentPut(b *testing.B, box BlackBox[int]) {
	cb := NewConcurrent(box)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = cb.Put(1)
		}
	})
}

func benchmarkConcurrentGet(b *testing.B, box BlackBox[int]) {
	cb := NewConcurrent(box)

	for i := 0; i < b.N; i++ {
		_ = cb.Put(i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = cb.Get()
		}
	})
}

func BenchmarkConcurrentFIFO_Put(b *testing.B) {
	box := New[int](
		WithStrategy(StrategyFIFO),
		WithInitialCapacity(b.N),
	)
	benchmarkConcurrentPut(b, box)
}

func BenchmarkConcurrentFIFO_Get(b *testing.B) {
	box := New[int](
		WithStrategy(StrategyFIFO),
		WithInitialCapacity(b.N),
	)
	benchmarkConcurrentGet(b, box)
}

func BenchmarkConcurrentConcreteFIFO_Put(b *testing.B) {
	box := NewFIFO[int](0, b.N)
	benchmarkConcurrentPut(b, box)
}

func BenchmarkConcurrentConcreteFIFO_Get(b *testing.B) {
	box := NewFIFO[int](0, b.N)
	benchmarkConcurrentGet(b, box)
}

func BenchmarkConcurrentLIFO_Put(b *testing.B) {
	box := New[int](
		WithStrategy(StrategyLIFO),
		WithInitialCapacity(b.N),
	)
	benchmarkConcurrentPut(b, box)
}

func BenchmarkConcurrentLIFO_Get(b *testing.B) {
	box := New[int](
		WithStrategy(StrategyLIFO),
		WithInitialCapacity(b.N),
	)
	benchmarkConcurrentGet(b, box)
}

func BenchmarkConcurrentConcreteLIFO_Put(b *testing.B) {
	box := NewLIFO[int](0, b.N)
	benchmarkConcurrentPut(b, box)
}

func BenchmarkConcurrentConcreteLIFO_Get(b *testing.B) {
	box := NewLIFO[int](0, b.N)
	benchmarkConcurrentGet(b, box)
}

func BenchmarkConcurrentRandom_Put(b *testing.B) {
	box := New[int](
		WithStrategy(StrategyRandom),
		WithInitialCapacity(b.N),
		WithSeed(42),
	)
	benchmarkConcurrentPut(b, box)
}

func BenchmarkConcurrentRandom_Get(b *testing.B) {
	box := New[int](
		WithStrategy(StrategyRandom),
		WithInitialCapacity(b.N),
		WithSeed(42),
	)
	benchmarkConcurrentGet(b, box)
}

func BenchmarkConcurrentConcreteRandom_Put(b *testing.B) {
	rng := rand.New(rand.NewSource(42))
	box := NewRandom[int](0, b.N, rng)
	benchmarkConcurrentPut(b, box)
}

func BenchmarkConcurrentConcreteRandom_Get(b *testing.B) {
	rng := rand.New(rand.NewSource(42))
	box := NewRandom[int](0, b.N, rng)
	benchmarkConcurrentGet(b, box)
}

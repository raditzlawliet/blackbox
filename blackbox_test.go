package blackbox

import (
	"math/rand"
	"testing"
	"time"
)

func TestFIFOStrategy(t *testing.T) {
	box := New[int](WithStrategy(StrategyFIFO))

	// Test Put
	for i := 1; i <= 5; i++ {
		err := box.Put(i)
		if err != nil {
			t.Fatalf("Failed to put item %d: %v", i, err)
		}
	}

	if box.Size() != 5 {
		t.Errorf("Expected size 5, got %d", box.Size())
	}

	// Test FIFO order (First In First Out)
	for i := 1; i <= 5; i++ {
		item, err := box.Get()
		if err != nil {
			t.Fatalf("Failed to get item: %v", err)
		}
		if item != i {
			t.Errorf("Expected item %d, got %d", i, item)
		}
	}

	if !box.IsEmpty() {
		t.Error("Box should be empty")
	}
}

func TestLIFOStrategy(t *testing.T) {
	box := New[int](WithStrategy(StrategyLIFO))

	// Test Put
	for i := 1; i <= 5; i++ {
		err := box.Put(i)
		if err != nil {
			t.Fatalf("Failed to put item %d: %v", i, err)
		}
	}

	if box.Size() != 5 {
		t.Errorf("Expected size 5, got %d", box.Size())
	}

	// Test LIFO order (Last In First Out)
	for i := 5; i >= 1; i-- {
		item, err := box.Get()
		if err != nil {
			t.Fatalf("Failed to get item: %v", err)
		}
		if item != i {
			t.Errorf("Expected item %d, got %d", i, item)
		}
	}

	if !box.IsEmpty() {
		t.Error("Box should be empty")
	}
}

func TestRandomStrategy(t *testing.T) {
	box := New[int](WithStrategy(StrategyRandom))

	// Test Put
	for i := 1; i <= 5; i++ {
		err := box.Put(i)
		if err != nil {
			t.Fatalf("Failed to put item %d: %v", i, err)
		}
	}

	if box.Size() != 5 {
		t.Errorf("Expected size 5, got %d", box.Size())
	}

	// Test Random retrieval (just verify all items are retrieved)
	retrieved := make(map[int]bool)
	for range 5 {
		item, err := box.Get()
		if err != nil {
			t.Fatalf("Failed to get item: %v", err)
		}
		if item < 1 || item > 5 {
			t.Errorf("Got unexpected item: %d", item)
		}
		retrieved[item] = true
	}

	if len(retrieved) != 5 {
		t.Errorf("Expected 5 unique items, got %d", len(retrieved))
	}

	if !box.IsEmpty() {
		t.Error("Box should be empty")
	}
}

func TestFIFOWithGrowth(t *testing.T) {
	// Test FIFO ring buffer growth
	box := New[int](
		WithStrategy(StrategyFIFO),
		WithInitialCapacity(4),
	)

	// Add items to trigger growth
	for i := 1; i <= 10; i++ {
		err := box.Put(i)
		if err != nil {
			t.Fatalf("Failed to put item %d: %v", i, err)
		}
	}

	if box.Size() != 10 {
		t.Errorf("Expected size 10, got %d", box.Size())
	}

	// Verify FIFO order is maintained after growth
	for i := 1; i <= 10; i++ {
		item, err := box.Get()
		if err != nil {
			t.Fatalf("Failed to get item: %v", err)
		}
		if item != i {
			t.Errorf("Expected item %d, got %d", i, item)
		}
	}
}

func TestMaxSize(t *testing.T) {
	box := New[int](
		WithStrategy(StrategyLIFO),
		WithMaxSize(3),
	)

	// Add up to max size
	for i := 1; i <= 3; i++ {
		err := box.Put(i)
		if err != nil {
			t.Fatalf("Failed to put item %d: %v", i, err)
		}
	}

	// Try to add beyond max size
	err := box.Put(4)
	if err != ErrBlackBoxFull {
		t.Errorf("Expected ErrBlackBoxFull, got %v", err)
	}

	if !box.IsFull() {
		t.Error("Box should be full")
	}
}

func TestClean(t *testing.T) {
	box := New[int](WithStrategy(StrategyLIFO))

	for i := 1; i <= 5; i++ {
		box.Put(i)
	}

	box.Clean()

	if !box.IsEmpty() {
		t.Error("Box should be empty after Clean()")
	}

	if box.Size() != 0 {
		t.Errorf("Expected size 0, got %d", box.Size())
	}
}

func TestPeek(t *testing.T) {
	box := New[int](WithStrategy(StrategyLIFO))

	box.Put(1)
	box.Put(2)
	box.Put(3)

	// Peek should return last item without removing
	item, err := box.Peek()
	if err != nil {
		t.Fatalf("Failed to peek: %v", err)
	}
	if item != 3 {
		t.Errorf("Expected peek to return 3, got %d", item)
	}

	// Size should remain unchanged
	if box.Size() != 3 {
		t.Errorf("Expected size 3 after peek, got %d", box.Size())
	}
}

func BenchmarkLIFOPut(b *testing.B) {
	box := New[int](WithStrategy(StrategyLIFO), WithInitialCapacity(b.N))
	i := 0
	b.ResetTimer()
	for b.Loop() {
		box.Put(i)
		i++
	}
}

func BenchmarkLIFOGet(b *testing.B) {
	box := New[int](WithStrategy(StrategyLIFO), WithInitialCapacity(b.N))
	for i := 0; i < b.N; i++ {
		box.Put(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = box.Get()
	}
}

func BenchmarkConcreteLIFOPut(b *testing.B) {
	box := NewLIFO[int](0, b.N)
	i := 0
	b.ResetTimer()
	for b.Loop() {
		box.Put(i)
		i++
	}
}

func BenchmarkConcreteLIFOGet(b *testing.B) {
	box := NewLIFO[int](0, b.N)
	for i := 0; i < b.N; i++ {
		box.Put(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = box.Get()
	}
}

func BenchmarkFIFOPut(b *testing.B) {
	box := New[int](WithStrategy(StrategyFIFO), WithInitialCapacity(b.N))
	i := 0
	b.ResetTimer()
	for b.Loop() {
		box.Put(i)
		i++
	}
}

func BenchmarkFIFOGet(b *testing.B) {
	box := New[int](WithStrategy(StrategyFIFO), WithInitialCapacity(b.N))
	for i := 0; i < b.N; i++ {
		box.Put(i)
	}
	b.ResetTimer()
	for b.Loop() {
		_, _ = box.Get()
	}
}

func BenchmarkConcreteFIFOPut(b *testing.B) {
	box := NewFIFO[int](0, b.N)
	i := 0
	b.ResetTimer()
	for b.Loop() {
		box.Put(i)
		i++
	}
}

func BenchmarkConcreteFIFOGet(b *testing.B) {
	box := NewFIFO[int](0, b.N)
	for i := 0; i < b.N; i++ {
		box.Put(i)
	}
	b.ResetTimer()
	for b.Loop() {
		_, _ = box.Get()
	}
}

func BenchmarkRandomPut(b *testing.B) {
	box := New[int](WithStrategy(StrategyRandom), WithInitialCapacity(b.N))
	i := 0
	b.ResetTimer()
	for b.Loop() {
		box.Put(i)
		i++
	}
}

func BenchmarkRandomGet(b *testing.B) {
	box := New[int](WithStrategy(StrategyRandom), WithInitialCapacity(b.N))
	for i := 0; i < b.N; i++ {
		box.Put(i)
	}
	b.ResetTimer()
	for b.Loop() {
		_, _ = box.Get()
	}
}

func BenchmarkConcreteRandomPut(b *testing.B) {
	box := NewRandom[int](0, b.N, rand.New(rand.NewSource(time.Now().UnixNano())))
	i := 0
	b.ResetTimer()
	for b.Loop() {
		box.Put(i)
		i++
	}
}
func BenchmarkConcreteRandomGet(b *testing.B) {
	box := NewRandom[int](0, b.N, rand.New(rand.NewSource(time.Now().UnixNano())))
	for i := 0; i < b.N; i++ {
		box.Put(i)
	}
	b.ResetTimer()
	for b.Loop() {
		_, _ = box.Get()
	}
}

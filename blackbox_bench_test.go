package blackbox

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkLIFOPut(b *testing.B) {
	box := New[int](WithStrategy(StrategyLIFO), WithInitialCapacity(b.N))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		box.Put(i)
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		box.Put(i)
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		box.Put(i)
	}
}

func BenchmarkFIFOGet(b *testing.B) {
	box := New[int](WithStrategy(StrategyFIFO), WithInitialCapacity(b.N))
	for i := 0; i < b.N; i++ {
		box.Put(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = box.Get()
	}
}

func BenchmarkConcreteFIFOPut(b *testing.B) {
	box := NewFIFO[int](0, b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		box.Put(i)
	}
}

func BenchmarkConcreteFIFOGet(b *testing.B) {
	box := NewFIFO[int](0, b.N)
	for i := 0; i < b.N; i++ {
		box.Put(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = box.Get()
	}
}

func BenchmarkRandomPut(b *testing.B) {
	box := New[int](WithStrategy(StrategyRandom), WithInitialCapacity(b.N))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		box.Put(i)
	}
}

func BenchmarkRandomGet(b *testing.B) {
	box := New[int](WithStrategy(StrategyRandom), WithInitialCapacity(b.N))
	for i := 0; i < b.N; i++ {
		box.Put(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = box.Get()
	}
}

func BenchmarkConcreteRandomPut(b *testing.B) {
	box := NewRandom[int](0, b.N, rand.New(rand.NewSource(time.Now().UnixNano())))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		box.Put(i)
	}
}
func BenchmarkConcreteRandomGet(b *testing.B) {
	box := NewRandom[int](0, b.N, rand.New(rand.NewSource(time.Now().UnixNano())))
	for i := 0; i < b.N; i++ {
		box.Put(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = box.Get()
	}
}

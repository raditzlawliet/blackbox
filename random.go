package blackbox

import (
	"math/rand"
)

type randomBox[T any] struct {
	items   []T
	rng     *rand.Rand
	maxSize int
}

// NewRandom creates a new Random blackbox with the specified maximum size, capacity and rng.
// Returns a concrete instance of lifo blackbox without interface.
func NewRandom[T any](maxSize, capacity int, rng *rand.Rand) *randomBox[T] {
	return &randomBox[T]{
		items:   make([]T, 0, capacity),
		maxSize: maxSize,
		rng:     rng,
	}
}

func (b *randomBox[T]) Put(item T) error {
	if b.maxSize > 0 && len(b.items) >= b.maxSize {
		return ErrBlackBoxFull
	}
	b.items = append(b.items, item)
	return nil
}

func (b *randomBox[T]) Get() (T, error) {
	if len(b.items) == 0 {
		var zero T
		return zero, ErrEmptyBlackBox
	}

	idx := b.rng.Intn(len(b.items))
	item := b.items[idx]
	lastIdx := len(b.items) - 1
	b.items[idx] = b.items[lastIdx]
	b.items = b.items[:lastIdx]
	return item, nil
}

// Peek returns a random item from the blackbox without removing it.
// In Random Strategy, Peek() behaviour will return different items when called multiple times,
// and not guaranteed to be the same item when Get() called as the last call to Peek().
func (b *randomBox[T]) Peek() (T, error) {
	if len(b.items) == 0 {
		var zero T
		return zero, ErrEmptyBlackBox
	}
	idx := b.rng.Intn(len(b.items))
	return b.items[idx], nil
}

func (b *randomBox[T]) Size() int {
	return len(b.items)
}

func (b *randomBox[T]) MaxSize() int {
	return b.maxSize
}

func (b *randomBox[T]) IsFull() bool {
	return b.maxSize > 0 && len(b.items) >= b.maxSize
}

func (b *randomBox[T]) IsEmpty() bool {
	return len(b.items) == 0
}

func (b *randomBox[T]) Clean() {
	b.items = b.items[:0]
}

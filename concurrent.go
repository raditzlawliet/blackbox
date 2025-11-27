package blackbox

import "sync"

// concurrentBox is a simple goroutine-safe wrapper around any BlackBox[T].
// It serializes all method calls with a mutex.
type concurrentBox[T any] struct {
	inner BlackBox[T]
	mu    sync.Mutex
}

// NewConcurrent wraps any BlackBox[T] and returns a goroutine-safe BlackBox[T].
// This is an opt-in wrapper; use the plain boxes directly for maximum
// performance when you don't need concurrency.
func NewConcurrent[T any](inner BlackBox[T]) BlackBox[T] {
	return &concurrentBox[T]{inner: inner}
}

func (c *concurrentBox[T]) Put(item T) error {
	c.mu.Lock()
	err := c.inner.Put(item)
	c.mu.Unlock()
	return err
}

func (c *concurrentBox[T]) Get() (T, error) {
	c.mu.Lock()
	item, err := c.inner.Get()
	c.mu.Unlock()
	return item, err
}

func (c *concurrentBox[T]) Peek() (T, error) {
	c.mu.Lock()
	item, err := c.inner.Peek()
	c.mu.Unlock()
	return item, err
}

func (c *concurrentBox[T]) Size() int {
	c.mu.Lock()
	size := c.inner.Size()
	c.mu.Unlock()
	return size
}

func (c *concurrentBox[T]) MaxSize() int {
	c.mu.Lock()
	size := c.inner.MaxSize()
	c.mu.Unlock()
	return size
}

func (c *concurrentBox[T]) IsFull() bool {
	c.mu.Lock()
	isFull := c.inner.IsFull()
	c.mu.Unlock()
	return isFull
}

func (c *concurrentBox[T]) IsEmpty() bool {
	c.mu.Lock()
	isEmpty := c.inner.IsEmpty()
	c.mu.Unlock()
	return isEmpty
}

func (c *concurrentBox[T]) Clean() {
	c.mu.Lock()
	c.inner.Clean()
	c.mu.Unlock()
}

// Compile-time assertion that concurrentBox implements BlackBox[T].
var _ BlackBox[any] = (*concurrentBox[any])(nil)

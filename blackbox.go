package blackbox

import (
	"errors"
	"math/rand"
	"time"
)

// Strategy defines how items are retrieved from the blackbox
type Strategy int

const (
	StrategyRandom Strategy = iota // Default: random retrieval
	StrategyFIFO                   // First In First Out
	StrategyLIFO                   // Last In First Out
)

// blackBoxConfig holds configuration that doesn't depend on the generic type T
type blackBoxConfig struct {
	strategy        Strategy
	maxSize         int // 0 means unlimited
	initialCapacity int // 0 means use default
	rng             *rand.Rand
}

// BlackBox is a generic container that stores items of type T
type BlackBox[T any] struct {
	blackBoxConfig
	items []T

	// Ring buffer fields (only used by FIFO strategy)
	fifoHead int // FIFO only: read position
	fifoTail int // FIFO only: write position
	fifoSize int // FIFO only: current number of items
}

var (
	ErrEmptyBlackBox = errors.New("blackbox is empty")
	ErrBlackBoxFull  = errors.New("blackbox is full")
)

const (
	defaultInitialCapacity = 8
	growthFactor           = 2
)

// Option is a function that configures the non-generic blackbox config
type Option func(*blackBoxConfig)

// WithStrategy sets the retrieval strategy for the blackbox
func WithStrategy(strategy Strategy) Option {
	return func(c *blackBoxConfig) {
		c.strategy = strategy
	}
}

// WithMaxSize sets the maximum capacity of the blackbox (0 = unlimited)
func WithMaxSize(size int) Option {
	return func(c *blackBoxConfig) {
		c.maxSize = size
	}
}

// WithSeed sets a custom random seed for reproducible random behavior
func WithSeed(seed int64) Option {
	return func(c *blackBoxConfig) {
		c.rng = rand.New(rand.NewSource(seed))
	}
}

// WithInitialCapacity sets the initial capacity to avoid early reallocations
// - FIFO: Pre-allocates ring buffer with fixed size
// - LIFO/Random: Pre-allocates slice capacity (grows automatically via append)
func WithInitialCapacity(capacity int) Option {
	return func(c *blackBoxConfig) {
		if capacity > 0 {
			c.initialCapacity = capacity
		}
	}
}

// New creates a new BlackBox with the specified options
func New[T any](opts ...Option) *BlackBox[T] {
	// Initialize config with defaults
	config := blackBoxConfig{
		strategy:        StrategyRandom,
		maxSize:         0, // unlimited by default
		initialCapacity: 0, // will use defaultInitialCapacity if 0
		rng:             rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	// Apply all options
	for _, opt := range opts {
		opt(&config)
	}

	// Create blackbox with config
	b := &BlackBox[T]{
		blackBoxConfig: config,
		fifoHead:       0,
		fifoTail:       0,
	}

	// Initialize items based on strategy
	capacity := config.initialCapacity
	if capacity == 0 {
		capacity = defaultInitialCapacity
	}

	if config.strategy == StrategyFIFO {
		// FIFO uses ring buffer - pre-allocate fixed size
		b.items = make([]T, capacity)
		b.fifoSize = 0
	} else {
		// LIFO/Random use native slice growth - start with zero length
		b.items = make([]T, 0, capacity)
	}

	return b
}

// fifoGrow expands the ring buffer for FIFO strategy only
func (b *BlackBox[T]) fifoGrow() {
	newCapacity := len(b.items) * growthFactor
	if b.maxSize > 0 && newCapacity > b.maxSize {
		newCapacity = b.maxSize
	}

	newItems := make([]T, newCapacity)

	// Handle ring buffer wrap-around for FIFO
	if b.fifoHead < b.fifoTail {
		copy(newItems, b.items[b.fifoHead:b.fifoTail])
	} else {
		n := copy(newItems, b.items[b.fifoHead:])
		copy(newItems[n:], b.items[:b.fifoTail])
	}
	b.fifoHead = 0
	b.fifoTail = b.fifoSize

	b.items = newItems
}

// Put adds an item to the box
func (b *BlackBox[T]) Put(item T) error {
	switch b.strategy {
	case StrategyFIFO:
		// Check max capacity
		if b.maxSize > 0 && b.fifoSize >= b.maxSize {
			return ErrBlackBoxFull
		}

		// Grow ring buffer if full
		if b.fifoSize >= len(b.items) {
			b.fifoGrow()
		}

		// Ring buffer: add at tail
		b.items[b.fifoTail] = item
		b.fifoTail = (b.fifoTail + 1) % len(b.items)
		b.fifoSize++

	case StrategyLIFO, StrategyRandom:
		// Check capacity for LIFO/Random
		if b.maxSize > 0 && len(b.items) >= b.maxSize {
			return ErrBlackBoxFull
		}
		// Use native Go slice append
		b.items = append(b.items, item)
	}

	return nil
}

// Get retrieves and removes an item from the box based on the strategy
func (b *BlackBox[T]) Get() (T, error) {
	if b.IsEmpty() {
		var zero T
		return zero, ErrEmptyBlackBox
	}

	var item T

	switch b.strategy {
	case StrategyFIFO:
		//ring buffer pop from head
		item = b.items[b.fifoHead]
		var zeroVal T
		b.items[b.fifoHead] = zeroVal // Clear reference for GC
		b.fifoHead = (b.fifoHead + 1) % len(b.items)
		b.fifoSize--

	case StrategyLIFO:
		// stack pop from end
		lastIdx := len(b.items) - 1
		item = b.items[lastIdx]
		b.items = b.items[:lastIdx] // Shrink slice

	case StrategyRandom:
		// swap with last element and pop
		idx := b.rng.Intn(len(b.items))
		item = b.items[idx]
		lastIdx := len(b.items) - 1
		b.items[idx] = b.items[lastIdx] // Swap with last element
		b.items = b.items[:lastIdx]     // Shrink slice
	}

	return item, nil
}

// Peek returns an item without removing it
func (b *BlackBox[T]) Peek() (T, error) {
	var zero T

	if b.IsEmpty() {
		return zero, ErrEmptyBlackBox
	}

	switch b.strategy {
	case StrategyFIFO:
		return b.items[b.fifoHead], nil
	case StrategyLIFO:
		return b.items[len(b.items)-1], nil
	case StrategyRandom:
		idx := b.rng.Intn(len(b.items))
		return b.items[idx], nil
	}

	return zero, nil
}

// Size returns the number of items in the box
func (b *BlackBox[T]) Size() int {
	if b.strategy == StrategyFIFO {
		return b.fifoSize
	}
	return len(b.items)
}

// Capacity returns the current internal capacity
func (b *BlackBox[T]) Capacity() int {
	return cap(b.items)
}

// MaxSize returns the maximum capacity (0 = unlimited)
func (b *BlackBox[T]) MaxSize() int {
	return b.maxSize
}

// IsFull returns true if the blackbox has reached maximum capacity
func (b *BlackBox[T]) IsFull() bool {
	if b.maxSize == 0 {
		return false
	}
	if b.strategy == StrategyFIFO {
		return b.fifoSize >= b.maxSize
	}
	return len(b.items) >= b.maxSize
}

// IsEmpty returns true if the blackbox has no items
func (b *BlackBox[T]) IsEmpty() bool {
	if b.strategy == StrategyFIFO {
		return b.fifoSize == 0
	}
	return len(b.items) == 0
}

// Clean removes all items from the blackbox
func (b *BlackBox[T]) Clean() {
	if b.strategy == StrategyFIFO {
		// Clear ring buffer references for GC
		var zero T
		for i := 0; i < b.fifoSize; i++ {
			idx := (b.fifoHead + i) % len(b.items)
			b.items[idx] = zero
		}
		b.fifoHead = 0
		b.fifoTail = 0
		b.fifoSize = 0
	} else {
		// LIFO/Random: simply reset slice
		b.items = b.items[:0]
	}
}

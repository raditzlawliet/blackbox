package blackbox

import (
	"errors"
	"math/rand"
	"time"
)

var (
	ErrEmptyBlackBox = errors.New("blackbox is empty")
	ErrBlackBoxFull  = errors.New("blackbox is full")
)

const (
	defaultInitialCapacity = 8
	growthFactor           = 2
)

// BlackBox is a generic container with a configurable retrieval strategy.
//
// Implementations of BlackBox[T] provide a way to store items of type T and
// retrieve them according to a strategy (FIFO, LIFO, Random).
//
// Method behavior (common across implementations):
//   - Put(item T) error
//     Insert an item into the blackbox. If the blackbox has a configured
//     maximum capacity and is already full, Put returns ErrBlackBoxFull.
//   - Get() (T, error)
//     Remove and return an item according to the configured retrieval strategy.
//     If the blackbox is empty, Get returns a zero value of T and ErrEmptyBlackBox.
//   - Peek() (T, error)
//     Return an item according to the configured retrieval strategy without
//     removing it. If the blackbox is empty, Peek returns a zero value of T
//     and ErrEmptyBlackBox. Note that for strategies like Random, repeated
//     calls to Peek may return different items.
//   - Size() int
//     Return the current number of stored items.
//   - MaxSize() int
//     Return the configured maximum capacity (0 means unlimited).
//   - IsFull() bool
//     Return true when the blackbox has a non-zero MaxSize and Size() >= MaxSize().
//   - IsEmpty() bool
//     Return true when Size() == 0.
//   - Clean()
//     Remove all items from the blackbox, resetting its size to zero.
//   - Items() []T
//     Return a copy of all items in the blackbox.
//
// Implementations returned by New[T] honor these semantics but differ in
// selection behavior (StrategyFIFO, StrategyLIFO, StrategyRandom).
type BlackBox[T any] interface {
	Put(item T) error
	Get() (T, error)
	Peek() (T, error)
	Size() int
	MaxSize() int
	IsFull() bool
	IsEmpty() bool
	Clean()
	Items() []T
}

// Strategy defines how items are retrieved from the blackbox
type Strategy int

const (
	StrategyRandom Strategy = iota // Default: random retrieval
	StrategyFIFO                   // First In First Out
	StrategyLIFO                   // Last In First Out
)

// config holds common configuration
type config struct {
	strategy        Strategy
	maxSize         int
	initialCapacity int
	seed            int64
	useSeed         bool
	useMaxSize      bool
}

// Option is a function that configures the blackbox
type Option func(*config)

// WithStrategy sets the retrieval strategy for the blackbox
func WithStrategy(strategy Strategy) Option {
	return func(c *config) {
		c.strategy = strategy
	}
}

// WithMaxSize sets the maximum capacity of the blackbox (0 = unlimited)
func WithMaxSize(size int) Option {
	return func(c *config) {
		c.maxSize = size
		c.useMaxSize = true
	}
}

// WithSeed sets a custom random seed for reproducible random behavior (Random Strategy)
func WithSeed(seed int64) Option {
	return func(c *config) {
		c.seed = seed
		c.useSeed = true
	}
}

// WithInitialCapacity sets the initial capacity to avoid early reallocations
func WithInitialCapacity(capacity int) Option {
	return func(c *config) {
		if capacity > 0 {
			c.initialCapacity = capacity
		}
	}
}

// parseOptions parses options into config
func parseOptions(opts []Option) config {
	cfg := config{
		maxSize:         0,
		initialCapacity: 0,
		useSeed:         false,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.initialCapacity == 0 {
		cfg.initialCapacity = defaultInitialCapacity
	}
	return cfg
}

// New creates a new BlackBox with the specified options.
//
// The returned implementation depends on the configured Strategy:
//   - StrategyFIFO -> FIFO behavior (first inserted is first returned)
//   - StrategyLIFO -> LIFO behavior (last inserted is first returned)
//   - StrategyRandom -> Random selection behavior (requires an RNG)
//
// For the Random strategy, if WithSeed was used the RNG will be seeded with
// the provided seed for reproducible behavior; otherwise a time-based seed is used.
func New[T any](opts ...Option) BlackBox[T] {
	cfg := parseOptions(opts)
	switch cfg.strategy {
	case StrategyFIFO:
		return NewFIFO[T](cfg.maxSize, cfg.initialCapacity)
	case StrategyLIFO:
		return NewLIFO[T](cfg.maxSize, cfg.initialCapacity)
	case StrategyRandom:
		fallthrough
	default:
		var rng *rand.Rand
		if cfg.useSeed {
			rng = rand.New(rand.NewSource(cfg.seed))
		} else {
			rng = rand.New(rand.NewSource(time.Now().UnixNano()))
		}
		return NewRandom[T](cfg.maxSize, cfg.initialCapacity, rng)
	}
}

// NewFrom creates a new BlackBox with existing data and the specified options
// items are copied so it safe to use the original slice after the blackbox is created.
// InitialCapacity will use items length
func NewFrom[T any](data []T, opts ...Option) BlackBox[T] {
	cfg := parseOptions(opts)
	if cfg.maxSize > 0 && cfg.maxSize < len(data) {
		cfg.maxSize = len(data)
	}
	switch cfg.strategy {
	case StrategyFIFO:
		return NewFIFOFrom[T](data, cfg.maxSize)
	case StrategyLIFO:
		return NewLIFOFrom[T](data, cfg.maxSize)
	case StrategyRandom:
		fallthrough
	default:
		var rng *rand.Rand
		if cfg.useSeed {
			rng = rand.New(rand.NewSource(cfg.seed))
		} else {
			rng = rand.New(rand.NewSource(time.Now().UnixNano()))
		}
		return NewRandomFrom[T](data, cfg.maxSize, rng)
	}
}

// NewFromBlackBox creates a new BlackBox with existing data and the specified options
// items are copied so it safe to use the original slice after the blackbox is created.
// InitialCapacity will use items length.
// MaxSize always has minimum box.MaxSize() or 0.
func NewFromBlackBox[T any](box BlackBox[T], opts ...Option) BlackBox[T] {
	cfg := parseOptions(opts)
	if cfg.useMaxSize {
		if cfg.maxSize > 0 && cfg.maxSize < box.Size() {
			cfg.maxSize = box.Size()
		}
	} else {
		cfg.maxSize = box.MaxSize()
	}
	switch cfg.strategy {
	case StrategyFIFO:
		return NewFIFOFromBlackBox[T](box, cfg.maxSize)
	case StrategyLIFO:
		return NewLIFOFromBlackBox[T](box, cfg.maxSize)
	case StrategyRandom:
		fallthrough
	default:
		var rng *rand.Rand
		if cfg.useSeed {
			rng = rand.New(rand.NewSource(cfg.seed))
		} else {
			rng = rand.New(rand.NewSource(time.Now().UnixNano()))
		}
		return NewRandomFromBlackBox[T](box, cfg.maxSize, rng)
	}
}

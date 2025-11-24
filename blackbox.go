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

// BlackBox is the interface that will be implemented by all strategy
type BlackBox[T any] interface {
	Put(item T) error
	Get() (T, error)
	Peek() (T, error)
	Size() int
	MaxSize() int
	IsFull() bool
	IsEmpty() bool
	Clean()
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

// New creates a new BlackBox with the specified options
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

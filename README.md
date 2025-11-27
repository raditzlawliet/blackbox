# BlackBox

[![Go Reference](https://pkg.go.dev/badge/github.com/raditzlawliet/blackbox.svg)](https://pkg.go.dev/github.com/raditzlawliet/blackbox)
[![codecov](https://codecov.io/gh/raditzlawliet/blackbox/graph/badge.svg?token=VJH1C5BFLN)](https://codecov.io/gh/raditzlawliet/blackbox)
[![Go Report Card](https://goreportcard.com/badge/github.com/raditzlawliet/blackbox)](https://goreportcard.com/report/github.com/raditzlawliet/blackbox)

A generic Go library that creates a literal "black box" - throw anything in, and see what comes out! Perfect for when you need unpredictability, or just want to manage collections with different retrieval strategies.

BlackBox is a type-safe, generic container where you can:

- **Put** anything in
- **Peek** at what might come out next (without removing it)
- **Get** item out using different strategies

The mystery? You can't see what's inside - you can only peek at one item at a time or get them out. Hehe...

## Features

- **Random Strategy**: Get a random item each time (default)
- **LIFO Strategy**: Last In, First Out (stack behavior)
- **FIFO Strategy**: First In, First Out (queue behavior)
- **Type-Safe**: Using Go generics
- **Zero Allocation**: No heap allocations during Get/Peek operations

## Requirements

- **Go 1.18 or higher** (requires generics support)

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/raditzlawliet/blackbox"
)

func main() {
    // Create a new black box with default random strategy
    box := blackbox.New[string]()

    // Put items in
    _ = box.Put("apple")
    _ = box.Put("banana")
    _ = box.Put("cherry")

    // Peek at what might come out (doesn't remove it)
    if item, err := box.Peek(); err == nil {
        fmt.Println("Peeked:", item)
    }

    // Get items out (removes them)
    for !box.IsEmpty() {
        item, _ := box.Get()
        fmt.Println("Got:", item)
    }
}
```

## Strategies

### Random (Default)

Get a random item each time — useful for draws and lotteries.

```go
box := blackbox.New[int](blackbox.WithStrategy(blackbox.StrategyRandom))
_ = box.Put(1)
_ = box.Put(2)
_ = box.Put(3)
```

### LIFO (Stack)

Last In, First Out — like a stack of plates.

```go
box := blackbox.New[string](blackbox.WithStrategy(blackbox.StrategyLIFO))
_ = box.Put("first")
_ = box.Put("second")
_ = box.Put("third")
```

### FIFO (Queue)

First In, First Out — a traditional queue.

```go
box := blackbox.New[string](blackbox.WithStrategy(blackbox.StrategyFIFO))
_ = box.Put("first")
_ = box.Put("second")
_ = box.Put("third")
```

## Configuration Options

- `WithMaxSize(int)`: set logical maximum number of items (0 = unlimited)
- `WithInitialCapacity(int)`: pre-allocate underlying storage to avoid early reallocations
- `WithSeed(int64)`: seed the RNG for the Random strategy (reproducible behavior)

## API Reference

Methods common to all boxes:

- `Put(item T) error` — insert an item (returns `ErrBlackBoxFull` if max size reached)
- `Get() (T, error)` — remove and return an item (returns `ErrEmptyBlackBox` if empty)
- `Peek() (T, error)` — view next item without removing
- `Size() int` — current number of items
- `MaxSize() int` — configured maximum size (0 = unlimited). You can't have zero-max-size box, for what?
- `IsFull() bool`, `IsEmpty() bool`
- `Clean()` — remove all items

Concrete constructors available for performance-sensitive use:

- `NewFIFO[T] (maxSize, capacity int) *fifoBox[T]`
- `NewLIFO[T] (maxSize, capacity int) *lifoBox[T]`
- `NewRandom[T] (maxSize, capacity int, rng *rand.Rand) *randomBox[T]`

Use the generic `New[T](opts...) BlackBox[T]` factory for convenience and option-based configuration.

## Concurrency

The core implementations (`fifo`, `lifo`, `random`) are intentionally lightweight and are **not** goroutine-safe by default to preserve single-threaded performance.

If you need safe concurrent access, we provide a simple, opt-in wrapper: `NewConcurrent`.

- `NewConcurrent(box)` returns a `BlackBox[T]` that serializes all calls with a mutex.
- This approach keeps the fast, lock-free implementations unchanged while offering an easy way to share a box across goroutines.

Example (concurrent wrapper):

```go
package main

import (
    "fmt"
    "sync"
    "github.com/raditzlawliet/blackbox"
)

func main() {
    // Create a concrete FIFO and wrap it for concurrent access
    fifo := blackbox.NewFIFO[int](0, 16)
    cbox := blackbox.NewConcurrent[int](fifo)

    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        defer wg.Done()
        _ = cbox.Put(42)
    }()

    go func() {
        defer wg.Done()
        if v, err := cbox.Get(); err == nil {
            fmt.Println("Got:", v)
        }
    }()

    wg.Wait()
}
```

Notes:

- See [`examples/concurrent`](example/concurrent/main.go) for a small runnable demo that shows producers and consumers using `NewConcurrent`.
- The concurrent wrapper serializes operations with a single `sync.Mutex`.

## Examples

- [`examples/basic`](example/basic/main.go) — basic usage for Random / LIFO / FIFO
- [`examples/concurrent`](example/concurrent/main.go) — simple concurrent usage demonstrating `NewConcurrent`
- [`examples/task_queue`](example/task_queue/main.go) — FIFO task queue
- [`examples/lucky_draw`](example/lucky_draw/main.go) — Random strategy example
- [`examples/undo_stack`](example/undo_stack/main.go) — LIFO undo/redo sample
- [`examples/concrete_types`](example/concrete_types/main.go) — direct constructor usage (concrete types)

## Performance

- FIFO uses a ring buffer for efficient Put/Get operations.
- LIFO uses append/slice operations.
- Random uses swap-with-last removal to keep operations efficient.
- For single-threaded hot paths, prefer the concrete constructors (`NewFIFO`, `NewLIFO`, `NewRandom`) when possible.

## Contributing

Feel free to:

- Report bugs
- Suggest new strategies
- Improve performance
- Add more features
- Share your creative use cases

---

Made with 🥰 and a bit of chaos...

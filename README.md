# BlackBox

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
    "raditzlawliet/blackbox"
)

func main() {
    // Create a new black box with default random strategy
    box := blackbox.New[string]()

    // Put items in
    box.Put("apple")
    box.Put("banana")
    box.Put("cherry")

    // Peek at what might come out (doesn't remove it)
    item, _ := box.Peek()
    fmt.Println("Peeked:", item)

    // Get items out (removes them)
    for !box.IsEmpty() {
        item, _ := box.Get()
        fmt.Println("Got:", item)
    }
}
```

## Strategies

### Random (Default)

Get a random item each time - perfect for lucky draws!

```go
box := blackbox.New[int](blackbox.WithStrategy(blackbox.StrategyRandom))
box.Put(1)
box.Put(2)
box.Put(3)

// Might get 2, then 1, then 3... who knows! 🎲
```

### LIFO (Stack)

Last In, First Out - like a stack of plates!

```go
box := blackbox.New[string](blackbox.WithStrategy(blackbox.StrategyLIFO))
box.Put("first")
box.Put("second")
box.Put("third")

box.Get() // Returns "third"
box.Get() // Returns "second"
box.Get() // Returns "first"
```

### FIFO (Queue)

First In, First Out - like a fair queue!

```go
box := blackbox.New[string](blackbox.WithStrategy(blackbox.StrategyFIFO))
box.Put("first")
box.Put("second")
box.Put("third")

box.Get() // Returns "first"
box.Get() // Returns "second"
box.Get() // Returns "third"
```

## Configuration Options

### Maximum Size

Limit how many items the box can hold:

```go
box := blackbox.New[int](
    blackbox.WithMaxSize(100), // Max 100 items
)

// When full:
err := box.Put(item) // Returns ErrBlackBoxFull
```

### Initial Capacity

Pre-allocate memory for better performance:

```go
box := blackbox.New[int](
    blackbox.WithInitialCapacity(1000), // Pre-allocate for 1000 items, useful on FIFO
)
```

### Custom Random Seed (Only Available for: StrategyRandom)

Make random behavior reproducible:

```go
box := blackbox.New[int](
    blackbox.WithStrategy(blackbox.StrategyRandom),
    blackbox.WithSeed(42), // Same seed = same random sequence
)
```

## API Reference

### Creating a BlackBox

```go
// Create with default options (random strategy)
box := blackbox.New[T]()

// Create with custom options
box := blackbox.New[T](options...)
```

### Adding Items

```go
err := box.Put(item)
// Returns ErrBlackBoxFull if max capacity reached
```

### Retrieving Items

```go
// Get and remove an item
item, err := box.Get()
// Returns ErrEmptyBlackBox if empty

// Peek at an item without removing
item, err := box.Peek()
// Returns ErrEmptyBlackBox if empty
```

### Status Checks

```go
// Check how many items are inside (but not what they are!)
size := box.Size()

// Check if empty
isEmpty := box.IsEmpty()

// Check if full (when max size is set)
isFull := box.IsFull()

// Get current internal capacity
capacity := box.Capacity()

// Get maximum size (0 = unlimited)
maxSize := box.MaxSize()
```

### Maintenance

```go
// Remove all items
box.Clean()
```

## Examples

Example implementations are available in the `examples/` directory:

- **[examples/basic](examples/basic/main.go)** - Basic usage of all three strategies (Random, LIFO, FIFO)
- **[examples/lucky_draw](examples/lucky_draw/main.go)** - Lucky draw system with random winner selection
- **[examples/task_queue](examples/task_queue/main.go)** - FIFO task queue with capacity management
- **[examples/undo_stack](examples/undo_stack/main.go)** - LIFO undo/redo system for command history
- **[examples/nested_blackbox](examples/nested_blackbox/main.go)** - Nested blackbox patterns
- **[examples/concrete_types](examples/concrete_types/main.go)** - Direct box creation with concrete types

### Lucky Draw System

```go
participants := blackbox.New[string](blackbox.WithStrategy(blackbox.StrategyRandom))
participants.Put("Alice")
participants.Put("Bob")
participants.Put("Charlie")

winner, _ := participants.Get()
fmt.Println("Winner:", winner)
```

#### Task Queue with Capacity

```go
taskQueue := blackbox.New[Task](
    blackbox.WithStrategy(blackbox.StrategyFIFO),
    blackbox.WithMaxSize(100),
)

for _, task := range tasks {
    if err := taskQueue.Put(task); err != nil {
        log.Printf("Queue full: %v", err)
    }
}
```

### Undo Stack

```go
undoStack := blackbox.New[Command](blackbox.WithStrategy(blackbox.StrategyLIFO))

cmd.Execute()
undoStack.Put(cmd)

// Undo
lastCmd, _ := undoStack.Get()
lastCmd.Undo()
```

## Type Support

BlackBox works with any type thanks to Go generics:

```go
// Basic types
intBox := blackbox.New[int]()
stringBox := blackbox.New[string]()
floatBox := blackbox.New[float64]()

// Structs
type Person struct {
    Name string
    Age  int
}
peopleBox := blackbox.New[Person]()

// Pointers
ptrBox := blackbox.New[*MyStruct]()

// Interfaces
interfaceBox := blackbox.New[io.Reader]()

// Custom types
type UserID string
userBox := blackbox.New[UserID]()
```

## Performance

BlackBox is highly optimized:

- **Random**: Removal using swap with last
- **LIFO**: Using slice operations
- **FIFO**: Using ring buffer

### Concrete Type

BlackBox also provide direct constructor each strategy so it can enable compiler optimizations

```go
lifoBox := blackbox.NewLIFO[string](5, 10)
fifoBox := blackbox.NewFIFO[string](5, 10)
randomBox := blackbox.NewRandom[int](0, b.N, rand.New(rand.NewSource(time.Now().UnixNano())))
```

### Benchmark Interface vs Direct Concrete Type

_This benchmark shows the performance difference between using the interface and the direct concrete type for each strategy. Only for reference, may different in realworld use case. Choose your self._

```sh
goos: windows
goarch: amd64
pkg: raditzlawliet/blackbox
cpu: 12th Gen Intel(R) Core(TM) i7-12700H
BenchmarkLIFOPut-20              	10000000	         9.007 ns/op	      49 B/op	       0 allocs/op
BenchmarkLIFOGet-20              	10000000	         1.727 ns/op	       0 B/op	       0 allocs/op
BenchmarkConcreteLIFOPut-20      	10000000	         5.923 ns/op	      49 B/op	       0 allocs/op
BenchmarkConcreteLIFOGet-20      	10000000	         0.2586 ns/op	       0 B/op	       0 allocs/op
BenchmarkFIFOPut-20              	10000000	         8.294 ns/op	      26 B/op	       0 allocs/op
BenchmarkFIFOGet-20              	10000000	         1.266 ns/op	       0 B/op	       0 allocs/op
BenchmarkConcreteFIFOPut-20      	10000000	         8.216 ns/op	      26 B/op	       0 allocs/op
BenchmarkConcreteFIFOGet-20      	10000000	         0.8256 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandomPut-20            	10000000	         6.864 ns/op	      49 B/op	       0 allocs/op
BenchmarkRandomGet-20            	10000000	         1.600 ns/op	       0 B/op	       0 allocs/op
BenchmarkConcreteRandomPut-20    	10000000	         5.747 ns/op	      49 B/op	       0 allocs/op
BenchmarkConcreteRandomGet-20    	10000000	         0.8947 ns/op	       0 B/op	       0 allocs/op
```

## Philosophy

Sometimes you need predictability. Sometimes you need chaos. BlackBox gives you both, wrapped in a type-safe, efficient package. Put something in. See what comes out. Enjoy the mystery!

## Contributing

Feel free to:

- Report bugs
- Suggest new strategies
- Improve performance
- Add more features
- Share your creative use cases

---

Made with 🥰 and a bit of chaos...

package main

import (
	"fmt"
	"math/rand"
	"raditzlawliet/blackbox"
)

func main() {
	concreteLifoExample()
	fmt.Println()

	concreteFifoExample()
	fmt.Println()

	concreteRandomExample()
	fmt.Println()
}

// concreteLifoExample demonstrates NewLIFO with concrete type
func concreteLifoExample() {
	// Create LIFO box directly - returns *lifoBox[T] (concrete type)
	// NewLIFO(maxSize, capacity)
	stack := blackbox.NewLIFO[string](5, 10)

	fmt.Println("✅ Created concrete LIFO stack")
	fmt.Printf("   Type: *lifoBox[string]\n")
	fmt.Printf("   Max size: %d\n", stack.MaxSize())

	// Push items
	items := []string{"Bottom", "Middle", "Top"}
	for _, item := range items {
		stack.Put(item)
		fmt.Printf("Push: %s (size: %d)\n", item, stack.Size())
	}

	// Pop items (LIFO order)
	fmt.Println("\nPopping from stack:")
	for !stack.IsEmpty() {
		item, _ := stack.Get()
		fmt.Printf("Pop: %s (remaining: %d)\n", item, stack.Size())
	}
}

// concreteFifoExample demonstrates NewFIFO with concrete type
func concreteFifoExample() {
	// Create FIFO box directly - returns *fifoBox[T] (concrete type)
	// NewFIFO(maxSize, capacity)
	queue := blackbox.NewFIFO[int](0, 100)

	fmt.Println("✅ Created concrete FIFO queue")
	fmt.Printf("   Type: *fifoBox[int]\n")

	// Enqueue items
	for i := 1; i <= 5; i++ {
		queue.Put(i * 10)
		fmt.Printf("Enqueue: %d (size: %d)\n", i*10, queue.Size())
	}

	// Dequeue items (FIFO order)
	fmt.Println("\nDequeuing from queue:")
	for !queue.IsEmpty() {
		item, _ := queue.Get()
		fmt.Printf("Dequeue: %d (remaining: %d)\n", item, queue.Size())
	}
}

// concreteRandomExample demonstrates NewRandom with concrete type
func concreteRandomExample() {
	// Create Random box directly with reproducible seed
	// NewRandom(maxSize, capacity, rng)
	rng := rand.New(rand.NewSource(42))
	randomBox := blackbox.NewRandom[string](0, 10, rng)

	fmt.Println("✅ Created concrete Random box")
	fmt.Printf("   Type: *randomBox[string]\n")
	fmt.Println("   Seed: 42 (reproducible)")

	// Add participants
	participants := []string{"Alice", "Bob", "Charlie", "Diana", "Eve"}
	for _, name := range participants {
		randomBox.Put(name)
		fmt.Printf("Add: %s (size: %d)\n", name, randomBox.Size())
	}

	// Draw winners randomly
	fmt.Println("\nRandom drawing:")
	for i := 1; !randomBox.IsEmpty(); i++ {
		winner, _ := randomBox.Get()
		fmt.Printf("Winner #%d: %s (remaining: %d)\n", i, winner, randomBox.Size())
	}
}

package main

import (
	"fmt"

	"github.com/raditzlawliet/blackbox"
)

func main() {
	randomExample()
	fmt.Println()

	lifoExample()
	fmt.Println()

	fifoExample()
}

// randomExample demonstrates the random strategy
func randomExample() {
	fmt.Println("📦 Random Strategy Example")
	fmt.Println("---------------------------")

	// Create a blackbox with random strategy (default)
	box := blackbox.New[string](
		blackbox.WithStrategy(blackbox.StrategyRandom),
		blackbox.WithSeed(42), // Use seed for reproducible results
	)

	// Put items in
	fruits := []string{"🍎 Apple", "🍌 Banana", "🍒 Cherry", "🍇 Grape", "🍊 Orange"}
	for _, fruit := range fruits {
		box.Put(fruit)
		fmt.Printf("Put: %s\n", fruit)
	}

	fmt.Printf("Box size: %d items\n", box.Size())

	// Peek at what might come out
	item, _ := box.Peek()
	fmt.Printf("\nPeeking (without removing): %s\n", item)
	fmt.Printf("Box size after peek: %d items (unchanged!)\n\n", box.Size())

	// Get items out randomly
	for !box.IsEmpty() {
		item, _ := box.Get()
		fmt.Printf("Got: %s (remaining: %d)\n", item, box.Size())
	}

	fmt.Printf("Box is now empty: %v\n", box.IsEmpty())
}

// lifoExample demonstrates the LIFO (Last In, First Out) strategy
func lifoExample() {
	fmt.Println("📚 LIFO Strategy Example (Stack)")
	fmt.Println("---------------------------------")

	// Create a blackbox with LIFO strategy
	box := blackbox.New[string](
		blackbox.WithStrategy(blackbox.StrategyLIFO),
	)

	// Put items in
	actions := []string{"1. First action", "2. Second action", "3. Third action", "4. Fourth action"}
	for _, action := range actions {
		box.Put(action)
		fmt.Printf("Put: %s\n", action)
	}

	fmt.Printf("Box size: %d items\n", box.Size())

	// Peek at what comes out next
	item, _ := box.Peek()
	fmt.Printf("\nPeeking at top (without removing): %s\n", item)

	// Get items out in LIFO order (last in, first out)
	fmt.Println("\nGetting items out (Last In, First Out):")
	for !box.IsEmpty() {
		item, _ := box.Get()
		fmt.Printf("Got: %s (remaining: %d)\n", item, box.Size())
	}
}

// fifoExample demonstrates the FIFO (First In, First Out) strategy
func fifoExample() {
	fmt.Println("📋 FIFO Strategy Example (Queue)")
	fmt.Println("---------------------------------")

	// Create a blackbox with FIFO strategy
	box := blackbox.New[string](
		blackbox.WithStrategy(blackbox.StrategyFIFO),
		blackbox.WithMaxSize(5), // Limit to 5 items
	)

	// Put items in
	customers := []string{"👤 Alice", "👤 Bob", "👤 Charlie", "👤 Diana", "👤 Eve", "👤 Frank"}
	for _, customer := range customers {
		err := box.Put(customer)
		if err != nil {
			fmt.Printf("X Queue full! %s has to wait\n", customer)
		} else {
			fmt.Printf("Put: %s\n", customer)
		}
	}

	fmt.Printf("\nQueue size: %d customers\n", box.Size())
	fmt.Printf("Queue is full: %v\n", box.IsFull())

	// Peek at who's next
	item, _ := box.Peek()
	fmt.Printf("\nNext customer to be served: %s\n", item)

	// Try to add one more (should fail)
	err := box.Put(customers[5])
	if err != nil {
		fmt.Printf("\nX Cannot add Frank: %v\n\n", err)
	}

	// Serve customers in FIFO order (first in, first out)
	for !box.IsEmpty() {
		customer, _ := box.Get()
		fmt.Printf("Serving: %s (waiting: %d)\n", customer, box.Size())
	}

	// Now Frank can join
	fmt.Println("\nNow there's space for Frank:")
	box.Put(customers[5])
	fmt.Printf("Put: Frank (queue size: %d)\n", box.Size())
	customer, _ := box.Get()
	fmt.Printf("Serving: %s (waiting: %d)\n", customer, box.Size())
}

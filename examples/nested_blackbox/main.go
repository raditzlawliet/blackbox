package main

import (
	"fmt"
	"raditzlawliet/blackbox"
)

// PrizeItem represents an item with a name and content
type PrizeItem struct {
	Name    string
	Content string
}

func main() {
	// Create inner blackboxes that contain references
	fmt.Println("\nCreating inner blackbox (Level 1)...")

	// We'll store box descriptions instead of actual boxes
	innerBlackBox1 := blackbox.New[PrizeItem](blackbox.WithStrategy(blackbox.StrategyRandom))
	innerBlackBox1.Put(PrizeItem{Name: "Red Mystery", Content: "Precious item"})
	innerBlackBox1.Put(PrizeItem{Name: "Blue Surprise", Content: "Artistic item"})
	fmt.Println("✓ Inner Box #1 created with 2 items descriptions")

	innerBlackBox2 := blackbox.New[PrizeItem](blackbox.WithStrategy(blackbox.StrategyFIFO))
	innerBlackBox2.Put(PrizeItem{Name: "Green Treat", Content: "Delicious treat"})
	innerBlackBox2.Put(PrizeItem{Name: "Silver Gift", Content: "Special gift"})
	fmt.Println("✓ Inner Box #2 created with 2 items descriptions")

	// Level 3: Create outer blackbox using pointer to blackboxes
	fmt.Println("\nCreating outer blackbox (Level 2)...")

	outerBlackBox := blackbox.New[*blackbox.BlackBox[PrizeItem]](blackbox.WithStrategy(blackbox.StrategyRandom))
	outerBlackBox.Put(innerBlackBox1)
	outerBlackBox.Put(innerBlackBox2)
	fmt.Println("✓ Outer box created containing 2 middle boxes")

	// Now let's unpack everything!
	fmt.Println("\n" + "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Starting to unpack the mystery!")
	fmt.Printf("Outer box contains: %d boxes\n\n", outerBlackBox.Size())

	// Unpack outer box
	unpackCount := 1
	for !outerBlackBox.IsEmpty() {
		fmt.Printf("Opening package #%d from outer box...\n", unpackCount)

		middleBox, _ := outerBlackBox.Get()
		fmt.Printf("   Found a middle box with %d items inside!\n", middleBox.Size())

		// Unpack middle box
		itemCount := 1
		for !middleBox.IsEmpty() {
			box, _ := middleBox.Get()
			fmt.Printf("   📦 Item #%d: %s\n", itemCount, box.Name)
			fmt.Printf("      Description: %s\n\n", box.Content)
			itemCount++
		}

		unpackCount++
	}
}

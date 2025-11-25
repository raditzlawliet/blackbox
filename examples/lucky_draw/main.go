package main

import (
	"fmt"

	"github.com/raditzlawliet/blackbox"
)

func main() {
	// Create a blackbox with random strategy for fair random selection
	participants := blackbox.New[string](
		blackbox.WithStrategy(blackbox.StrategyRandom),
		blackbox.WithSeed(42), // Remove this for truly random results
	)

	// Register participants
	names := []string{
		"Alice Johnson",
		"Bob Smith",
		"Charlie Brown",
		"Diana Prince",
		"Eve Anderson",
		"Frank Miller",
		"Grace Lee",
		"Henry Wilson",
		"Ivy Chen",
		"Jack Taylor",
	}

	fmt.Println("Registering participants...")
	for _, name := range names {
		participants.Put(name)
		fmt.Printf("✓ %s registered\n", name)
	}

	fmt.Printf("\nTotal participants: %d\n\n", participants.Size())

	// Draw winners
	prizes := []string{"🥇 First Prize - $1000", "🥈 Second Prize - $500", "🥉 Third Prize - $250"}

	fmt.Println("Drawing winners...")

	for _, prize := range prizes {
		if participants.IsEmpty() {
			fmt.Println("No more participants!")
			break
		}

		// Dramatic pause effect
		fmt.Printf("Drawing for %s...\n", prize)

		winner, _ := participants.Get()
		fmt.Printf("🎉 Winner: %s\n\n", winner)
	}

	fmt.Printf("Remaining participants: %d\n", participants.Size())

	// Consolation prizes for everyone else
	if !participants.IsEmpty() {
		fmt.Println("Consolation Prize Winners:")
		count := 1
		for !participants.IsEmpty() {
			person, _ := participants.Get()
			fmt.Printf("%d. %s - Gift Card $50\n", count, person)
			count++
		}
	}

	fmt.Println("\nLucky draw completed!")
}

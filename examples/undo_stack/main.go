package main

import (
	"fmt"
	"raditzlawliet/blackbox"
)

// Command represents an action that can be executed and undone
type Command struct {
	ID          int
	Action      string
	Description string
	Data        string
}

// Execute simulates executing a command
func (c Command) Execute() {
	fmt.Printf("▶️  Executing: %s (%s)\n", c.Action, c.Description)
}

// Undo simulates undoing a command
func (c Command) Undo() {
	fmt.Printf("◀️  Undoing: %s (%s)\n", c.Action, c.Description)
}

func main() {
	// Create an undo stack with LIFO strategy
	undoStack := blackbox.New[Command](
		blackbox.WithStrategy(blackbox.StrategyLIFO),
		blackbox.WithMaxSize(10), // Keep last 10 commands
	)

	// Create a redo stack
	redoStack := blackbox.New[Command](
		blackbox.WithStrategy(blackbox.StrategyLIFO),
	)

	// Define some commands
	commands := []Command{
		{ID: 1, Action: "CREATE", Description: "Create new document", Data: "document.txt"},
		{ID: 2, Action: "WRITE", Description: "Write 'Hello World'", Data: "Hello World"},
		{ID: 3, Action: "WRITE", Description: "Write 'This is a test'", Data: "This is a test"},
		{ID: 4, Action: "FORMAT", Description: "Bold text", Data: "bold"},
		{ID: 5, Action: "WRITE", Description: "Write 'More content'", Data: "More content"},
	}

	// Execute commands and add to undo stack
	for _, cmd := range commands {
		cmd.Execute()
		undoStack.Put(cmd)
	}

	fmt.Printf("\nUndo stack size: %d commands\n", undoStack.Size())

	// Peek at what command would be undone next
	if !undoStack.IsEmpty() {
		nextCmd, _ := undoStack.Peek()
		fmt.Printf("Next command to undo: #%d - %s\n", nextCmd.ID, nextCmd.Action)
	}

	// Undo last 3 commands
	fmt.Println("\nUndoing last 3 commands...")
	for i := 0; i < 3 && !undoStack.IsEmpty(); i++ {
		cmd, _ := undoStack.Get()
		cmd.Undo()
		redoStack.Put(cmd) // Move to redo stack
		fmt.Printf("Undo stack: %d | Redo stack: %d\n\n", undoStack.Size(), redoStack.Size())
	}

	// Show current state
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Current state:\n")
	fmt.Printf("- Undo stack: %d commands\n", undoStack.Size())
	fmt.Printf("- Redo stack: %d commands\n", redoStack.Size())

	// Redo 2 commands
	fmt.Println("\nRedoing 2 commands...")
	for i := 0; i < 2 && !redoStack.IsEmpty(); i++ {
		cmd, _ := redoStack.Get()
		cmd.Execute()
		undoStack.Put(cmd) // Move back to undo stack
		fmt.Printf("Undo stack: %d | Redo stack: %d\n\n", undoStack.Size(), redoStack.Size())
	}

	// Execute a new command (this clears redo stack in real apps)
	fmt.Println("Executing new command...")
	newCmd := Command{ID: 6, Action: "SAVE", Description: "Save document", Data: "document.txt"}
	newCmd.Execute()
	undoStack.Put(newCmd)

	// Clear redo stack (new action invalidates redo history)
	redoStack.Clean()
	fmt.Println("\n(Redo stack cleared - new action invalidates redo history)")

	// Final state
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Final state:\n")
	fmt.Printf("- Undo stack: %d commands\n", undoStack.Size())
	fmt.Printf("- Redo stack: %d commands\n", redoStack.Size())

	// Demonstrate undo all
	fmt.Println("\nUndoing all remaining commands...")
	count := 1
	for !undoStack.IsEmpty() {
		cmd, _ := undoStack.Get()
		fmt.Printf("%d. ", count)
		cmd.Undo()
		count++
	}

	fmt.Printf("Undo stack is empty: %v\n", undoStack.IsEmpty())
}

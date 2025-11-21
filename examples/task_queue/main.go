package main

import (
	"fmt"
	"raditzlawliet/blackbox"
	"time"
)

// Task represents a work item to be processed
type Task struct {
	ID          int
	Name        string
	Description string
}

func main() {
	// Create a task queue with FIFO strategy and max capacity
	taskQueue := blackbox.New[Task](
		blackbox.WithStrategy(blackbox.StrategyFIFO),
		blackbox.WithMaxSize(5), // Limit to 5 concurrent tasks
		blackbox.WithInitialCapacity(5),
	)

	// Prepare tasks
	tasks := []Task{
		{ID: 1, Name: "Process Payment", Description: "Process customer payment #1234"},
		{ID: 2, Name: "Send Email", Description: "Send welcome email to new user"},
		{ID: 3, Name: "Generate Report", Description: "Generate monthly sales report"},
		{ID: 4, Name: "Backup Database", Description: "Perform daily database backup"},
		{ID: 5, Name: "Update Cache", Description: "Refresh application cache"},
		{ID: 6, Name: "Cleanup Logs", Description: "Remove old log files"},
		{ID: 7, Name: "Sync Data", Description: "Synchronize with external API"},
	}

	// Try to add all tasks
	fmt.Println("Adding tasks to queue...")
	rejectedTasks := []Task{}

	for _, task := range tasks {
		err := taskQueue.Put(task)
		if err != nil {
			fmt.Printf("X Queue full! Task #%d (%s) rejected\n", task.ID, task.Name)
			rejectedTasks = append(rejectedTasks, task)
		} else {
			fmt.Printf("✓ Task #%d queued: %s\n", task.ID, task.Name)
		}
	}

	fmt.Printf("\nQueue Status:\n")
	fmt.Printf("- Current size: %d\n", taskQueue.Size())
	fmt.Printf("- Max capacity: %d\n", taskQueue.MaxSize())
	fmt.Printf("- Is full: %v\n", taskQueue.IsFull())
	fmt.Printf("- Rejected tasks: %d\n", len(rejectedTasks))

	// Peek at next task
	if !taskQueue.IsEmpty() {
		nextTask, _ := taskQueue.Peek()
		fmt.Printf("\nNext task to process: #%d - %s\n", nextTask.ID, nextTask.Name)
	}

	// Process tasks
	fmt.Println("\nProcessing tasks (FIFO order)...")
	processedCount := 0

	for !taskQueue.IsEmpty() {
		task, _ := taskQueue.Get()
		fmt.Printf("  Processing Task #%d: %s\n", task.ID, task.Name)
		fmt.Printf("    Description: %s\n", task.Description)

		// Simulate task processing
		time.Sleep(100 * time.Millisecond)

		fmt.Printf("      Completed! (Remaining in queue: %d)\n\n", taskQueue.Size())
		processedCount++

		// Add rejected tasks back when there's space
		if len(rejectedTasks) > 0 && !taskQueue.IsFull() {
			rejectedTask := rejectedTasks[0]
			rejectedTasks = rejectedTasks[1:]

			err := taskQueue.Put(rejectedTask)
			if err == nil {
				fmt.Printf("  Previously rejected Task #%d added to queue\n\n", rejectedTask.ID)
			}
		}
	}

	fmt.Printf("Total tasks processed: %d\n", processedCount)
	fmt.Printf("Queue is empty: %v\n", taskQueue.IsEmpty())
}

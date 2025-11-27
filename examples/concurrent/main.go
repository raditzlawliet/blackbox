package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/raditzlawliet/blackbox"
)

func main() {
	// Create a concrete FIFO and wrap it with the concurrent wrapper.
	// Use a modest capacity for the example.
	fifo := blackbox.NewFIFO[int](0, 16)
	cbox := blackbox.NewConcurrent[int](fifo)

	producers := 3
	itemsPerProducer := 10
	consumers := 2

	var wgProducers sync.WaitGroup
	var wgConsumers sync.WaitGroup

	// Channel to signal that all producers finished producing.
	producersDone := make(chan struct{})

	// Start producers.
	wgProducers.Add(producers)
	for p := 0; p < producers; p++ {
		id := p + 1
		go func(pid int) {
			defer wgProducers.Done()
			for i := 0; i < itemsPerProducer; i++ {
				item := pid*100 + i // produce a distinguishable item
				if err := cbox.Put(item); err != nil {
					// Should not happen here since maxSize is 0 (unlimited)
					fmt.Printf("producer %d: failed to put %d: %v\n", pid, item, err)
				} else {
					fmt.Printf("producer %d: put %d\n", pid, item)
				}
				// Sleep a bit to simulate work and interleave producers/consumers
				time.Sleep(10 * time.Millisecond)
			}
		}(id)
	}

	// Start consumers.
	totalItems := producers * itemsPerProducer
	wgConsumers.Add(consumers)
	for c := 0; c < consumers; c++ {
		id := c + 1
		go func(cid int) {
			defer wgConsumers.Done()
			consumed := 0
			for {
				item, err := cbox.Get()
				if err == blackbox.ErrEmptyBlackBox {
					// If producers are done and box is empty, we're finished.
					select {
					case <-producersDone:
						if cbox.IsEmpty() {
							// nothing more to consume
							return
						}
						// else continue trying
					default:
						// producers still running, wait a bit and retry
						time.Sleep(15 * time.Millisecond)
					}
					continue
				}
				// Successfully got an item
				fmt.Printf("consumer %d: got %d\n", cid, item)
				consumed++
				// Optional small delay to simulate work
				time.Sleep(20 * time.Millisecond)

				// Quick exit if we've consumed everything (best-effort)
				if consumed >= totalItems {
					return
				}
			}
		}(id)
	}

	// Wait for producers to finish, then close the done channel.
	go func() {
		wgProducers.Wait()
		close(producersDone)
	}()

	// Wait for consumers to finish.
	wgConsumers.Wait()

	// Drain any remaining items (should be none).
	for !cbox.IsEmpty() {
		it, err := cbox.Get()
		if err != nil {
			break
		}
		fmt.Printf("drain: got %v\n", it)
	}

	fmt.Println("All done.")
}

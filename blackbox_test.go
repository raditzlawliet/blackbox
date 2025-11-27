package blackbox

import (
	"slices"
	"testing"
)

func TestFIFOStrategy(t *testing.T) {
	box := New[int](WithStrategy(StrategyFIFO))

	// Test Put
	for i := 1; i <= 5; i++ {
		err := box.Put(i)
		if err != nil {
			t.Fatalf("Failed to put item %d: %v", i, err)
		}
	}

	if box.Size() != 5 {
		t.Errorf("Expected size 5, got %d", box.Size())
	}

	// Test FIFO order (First In First Out)
	for i := 1; i <= 5; i++ {
		item, err := box.Get()
		if err != nil {
			t.Fatalf("Failed to get item: %v", err)
		}
		if item != i {
			t.Errorf("Expected item %d, got %d", i, item)
		}
	}

	if !box.IsEmpty() {
		t.Error("Box should be empty")
	}

	if _, err := box.Get(); err != ErrEmptyBlackBox {
		t.Error("Should be error Box is empty")
	}
}

func TestLIFOStrategy(t *testing.T) {
	box := New[int](WithStrategy(StrategyLIFO))

	// Test Put
	for i := 1; i <= 5; i++ {
		err := box.Put(i)
		if err != nil {
			t.Fatalf("Failed to put item %d: %v", i, err)
		}
	}

	if box.Size() != 5 {
		t.Errorf("Expected size 5, got %d", box.Size())
	}

	// Test LIFO order (Last In First Out)
	for i := 5; i >= 1; i-- {
		item, err := box.Get()
		if err != nil {
			t.Fatalf("Failed to get item: %v", err)
		}
		if item != i {
			t.Errorf("Expected item %d, got %d", i, item)
		}
	}

	if !box.IsEmpty() {
		t.Error("Box should be empty")
	}

	if _, err := box.Get(); err != ErrEmptyBlackBox {
		t.Error("Should be error Box is empty")
	}
}

func TestRandomStrategy(t *testing.T) {
	box := New[int](WithStrategy(StrategyRandom))

	// Test Put
	for i := 1; i <= 5; i++ {
		err := box.Put(i)
		if err != nil {
			t.Fatalf("Failed to put item %d: %v", i, err)
		}
	}

	if box.Size() != 5 {
		t.Errorf("Expected size 5, got %d", box.Size())
	}

	// Test Random retrieval (just verify all items are retrieved)
	retrieved := make(map[int]bool)
	for i := 0; i < 5; i++ {
		item, err := box.Get()
		if err != nil {
			t.Fatalf("Failed to get item: %v", err)
		}
		if item < 1 || item > 5 {
			t.Errorf("Got unexpected item: %d", item)
		}
		retrieved[item] = true
	}

	if len(retrieved) != 5 {
		t.Errorf("Expected 5 unique items, got %d", len(retrieved))
	}

	if !box.IsEmpty() {
		t.Error("Box should be empty")
	}

	if _, err := box.Get(); err != ErrEmptyBlackBox {
		t.Error("Should be error Box is empty")
	}
}

func TestRandomStrategyWithSeed(t *testing.T) {
	seed := int64(42)

	box1 := New[int](
		WithStrategy(StrategyRandom),
		WithSeed(seed),
	)
	box2 := New[int](
		WithStrategy(StrategyRandom),
		WithSeed(seed),
	)

	// Test Put for both boxes
	for i := 1; i <= 5; i++ {
		if err := box1.Put(i); err != nil {
			t.Fatalf("Failed to put item %d into box1: %v", i, err)
		}
		if err := box2.Put(i); err != nil {
			t.Fatalf("Failed to put item %d into box2: %v", i, err)
		}
	}

	// Retrieve sequences from both boxes and ensure they are identical
	seq1 := make([]int, 0, 5)
	seq2 := make([]int, 0, 5)
	for i := 0; i < 5; i++ {
		a, err := box1.Get()
		if err != nil {
			t.Fatalf("Failed to get item from box1: %v", err)
		}
		b, err := box2.Get()
		if err != nil {
			t.Fatalf("Failed to get item from box2: %v", err)
		}
		seq1 = append(seq1, a)
		seq2 = append(seq2, b)
	}

	for i := 0; i < 5; i++ {
		if seq1[i] != seq2[i] {
			t.Fatalf("Expected sequences to be identical for same seed, but differ at index %d: %d vs %d", i, seq1[i], seq2[i])
		}
	}

	// With a different seed, expect (with high probability) a different sequence
	box3 := New[int](
		WithStrategy(StrategyRandom),
		WithSeed(7),
	)
	for i := 1; i <= 5; i++ {
		if err := box3.Put(i); err != nil {
			t.Fatalf("Failed to put item %d into box3: %v", i, err)
		}
	}
	seq3 := make([]int, 0, 5)
	for i := 0; i < 5; i++ {
		v, err := box3.Get()
		if err != nil {
			t.Fatalf("Failed to get item from box3: %v", err)
		}
		seq3 = append(seq3, v)
	}

	same := true
	for i := 0; i < 5; i++ {
		if seq1[i] != seq3[i] {
			same = false
			break
		}
	}
	if same {
		t.Error("Expected different sequence for a different seed, but sequences were identical (very unlikely)")
	}
}

func TestFIFOWithGrowth(t *testing.T) {
	// Test FIFO ring buffer growth
	box := New[int](
		WithStrategy(StrategyFIFO),
		WithInitialCapacity(4),
	)

	// Add items to trigger growth
	for i := 1; i <= 10; i++ {
		err := box.Put(i)
		if err != nil {
			t.Fatalf("Failed to put item %d: %v", i, err)
		}
	}

	if box.Size() != 10 {
		t.Errorf("Expected size 10, got %d", box.Size())
	}

	// Verify FIFO order is maintained after growth
	for i := 1; i <= 10; i++ {
		item, err := box.Get()
		if err != nil {
			t.Fatalf("Failed to get item: %v", err)
		}
		if item != i {
			t.Errorf("Expected item %d, got %d", i, item)
		}
	}
}

func TestMaxSize(t *testing.T) {
	strategies := []Strategy{StrategyFIFO, StrategyLIFO, StrategyRandom}
	for _, strategy := range strategies {
		box := New[int](
			WithStrategy(strategy),
			WithMaxSize(3),
		)

		// Add up to max size
		for i := 1; i <= 3; i++ {
			err := box.Put(i)
			if err != nil {
				t.Fatalf("Failed to put item %d: %v", i, err)
			}
		}

		// Try to add beyond max size
		err := box.Put(4)
		if err != ErrBlackBoxFull {
			t.Errorf("Expected ErrBlackBoxFull, got %v", err)
		}

		if !box.IsFull() {
			t.Error("Box should be full")
		}

		if box.MaxSize() != 3 {
			t.Errorf("Expected max size 3, got %d", box.MaxSize())
		}
	}
}

func TestClean(t *testing.T) {
	strategies := []Strategy{StrategyFIFO, StrategyLIFO, StrategyRandom}
	for _, strategy := range strategies {
		box := New[int](WithStrategy(strategy))

		for i := 1; i <= 5; i++ {
			box.Put(i)
		}

		box.Clean()

		if !box.IsEmpty() {
			t.Error("Box should be empty after Clean()")
		}

		if box.Size() != 0 {
			t.Errorf("Expected size 0, got %d", box.Size())
		}
	}
}

func TestPeek(t *testing.T) {
	strategies := []Strategy{StrategyFIFO, StrategyLIFO, StrategyRandom}
	for _, strategy := range strategies {
		box := New[int](WithStrategy(strategy))

		box.Put(1)
		box.Put(2)
		box.Put(3)

		// Peek should return last item without removing
		item, err := box.Peek()
		if err != nil {
			t.Fatalf("Failed to peek: %v", err)
		}
		switch strategy {
		case StrategyFIFO:
			if item != 1 {
				t.Errorf("Expected peek to return 1, got %d", item)
			}
		case StrategyLIFO:
			if item != 3 {
				t.Errorf("Expected peek to return 3, got %d", item)
			}
		case StrategyRandom:
			if !slices.Contains([]int{1, 2, 3}, item) {
				t.Errorf("Expected peek to return 1 to 3, got %d", item)
			}
		}

		// Size should remain unchanged
		if box.Size() != 3 {
			t.Errorf("Expected size 3 after peek, got %d", box.Size())
		}

		// Get all items
		for i := 0; i < 3; i++ {
			box.Get()
		}

		// Should be error
		if _, err := box.Peek(); err != ErrEmptyBlackBox {
			t.Errorf("Expected ErrEmptyBlackBox, got %v", err)
		}

		// Should be empty
		if box.Size() != 0 {
			t.Errorf("Expected size 0 after get all and peek, got %d", box.Size())
		}
	}
}

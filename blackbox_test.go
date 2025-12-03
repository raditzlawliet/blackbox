package blackbox

import (
	"math/rand"
	"testing"
	"time"
)

func ContainsInt(s []int, v int) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}

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

	originalItemsBox3 := box3.Items() // for box 4
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

	// With a different seed, expect (with high probability) a different sequence but using NewFrom
	box4 := NewFrom[int](
		originalItemsBox3,
		WithStrategy(StrategyRandom),
		WithSeed(7),
	)
	seq4 := make([]int, 0, 5)
	for i := 0; i < 5; i++ {
		v, err := box4.Get()
		if err != nil {
			t.Fatalf("Failed to get item from box4: %v", err)
		}
		seq4 = append(seq4, v)
	}

	same4 := true
	for i := 0; i < 5; i++ {
		if seq1[i] != seq4[i] {
			same4 = false
			break
		}
	}
	if same4 {
		t.Error("Expected different sequence for a different seed, but sequences were identical (very unlikely)")
	}

	// With a different seed, expect (with high probability) a different sequence but using NewFrom
	box5shadow := NewFrom[int](
		originalItemsBox3,
		WithStrategy(StrategyRandom),
	)
	box5 := NewFromBox[int](
		box5shadow,
		WithStrategy(StrategyRandom),
		WithSeed(7),
	)
	seq5 := make([]int, 0, 5)
	for i := 0; i < 5; i++ {
		v, err := box5.Get()
		if err != nil {
			t.Fatalf("Failed to get item from box4: %v", err)
		}
		seq5 = append(seq5, v)
	}

	same5 := true
	for i := 0; i < 5; i++ {
		if seq1[i] != seq5[i] {
			same5 = false
			break
		}
	}
	if same5 {
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
			if !ContainsInt([]int{1, 2, 3}, item) {
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

func TestItems(t *testing.T) {
	strategies := []Strategy{StrategyFIFO, StrategyLIFO, StrategyRandom}
	for _, strategy := range strategies {
		box := New[int](WithStrategy(strategy))

		box.Put(1)
		box.Put(2)
		box.Put(3)
		box.Get()
		box.Get()

		items := box.Items()
		if len(items) != 1 {
			t.Errorf("Strategy=%v Expected items length only 1, got %d", strategy, len(items))
		}

		switch strategy {
		case StrategyFIFO:
			if items[0] != 3 {
				t.Errorf("Expected item is 3, got %d", items[0])
			}
		case StrategyLIFO:
			if items[0] != 1 {
				t.Errorf("Expected item is 1, got %d", items[0])
			}
		case StrategyRandom:
			if !ContainsInt([]int{1, 2, 3}, items[0]) {
				t.Errorf("Expected item is 1 to 3, got %d", items[0])
			}
		}

		box.Put(4)
		box.Put(5)
		box.Put(6)
		box.Get()
		box.Get()

		items2 := box.Items()
		if len(items2) != 2 {
			t.Errorf("Strategy=%v Expected items length only 2, got %d", strategy, len(items2))
		}

		switch strategy {
		case StrategyFIFO:
			if items[0] != 3 {
				t.Errorf("Expected item is 3, got %d", items[0])
			}
		case StrategyLIFO:
			if items[0] != 1 {
				t.Errorf("Expected item is 1, got %d", items[0])
			}
		case StrategyRandom:
			if !ContainsInt([]int{1, 2, 3, 4, 5, 6}, items[0]) {
				t.Errorf("Expected item is 1 to 6, got %d", items[0])
			}
		}

		box.Clean()
		items3 := box.Items()
		if len(items3) != 0 {
			t.Errorf("Strategy=%v Expected items length only 0, got %d", strategy, len(items3))
		}
	}
}

func TestItemsFIFO(t *testing.T) {
	box := NewFIFO[int](0, defaultInitialCapacity)
	for i := 0; i < 40; i++ {
		box.Put(i)
	}
	for i := 0; i < 10; i++ {
		box.Get()
	}
	if box.head != 10 {
		t.Errorf("Expected head is 10, got %d", box.head)
	}
	items := box.Items()
	if len(items) != 30 {
		t.Errorf("Expected items length only 30, got %d", len(items))
	}

	// should be head not reset until tail >= head
	for i := 0; i < 30; i++ {
		box.Put(i)
	}
	if box.head != 10 {
		t.Errorf("Expected head is 0, got %d", box.head)
	}
	items2 := box.Items()
	if len(items2) != 60 {
		t.Errorf("Expected items length only 60, got %d", len(items))
	}

	// should be reset the head
	for i := 0; i < 10; i++ {
		box.Put(i)
	}
	if box.head != 0 {
		t.Errorf("Expected head is 0, got %d", box.head)
	}
	items3 := box.Items()
	if len(items3) != 70 {
		t.Errorf("Expected items length only 70, got %d", len(items))
	}
}

func TestNewFrom(t *testing.T) {
	strategies := []Strategy{StrategyFIFO, StrategyLIFO, StrategyRandom}
	for _, strategy := range strategies {
		data := []int{1, 2, 3}
		box := NewFrom[int](data, WithStrategy(strategy))
		if box.Size() != 3 {
			t.Errorf("Expected size is 3, got %d", box.Size())
		}
		if box.MaxSize() != 0 {
			t.Errorf("Expected max size is 0, got %d", box.MaxSize())
		}

		// MaxSize() will be set, because MaxSize() > Size()
		box2 := NewFromBox[int](box, WithStrategy(strategy), WithMaxSize(20))
		if box2.Size() != 3 {
			t.Errorf("Expected size is 3, got %d", box2.Size())
		}
		if box2.MaxSize() != 20 {
			t.Errorf("Expected max size is 20, got %d", box2.MaxSize())
		}

		// set MaxSize will be ignore because MaxSize() < Size()
		box3 := NewFromBox[int](box2, WithStrategy(strategy), WithMaxSize(1))
		if box3.Size() != 3 {
			t.Errorf("Expected size is 3, got %d", box3.Size())
		}
		if box3.MaxSize() != 3 {
			t.Errorf("Expected max size is 3, got %d", box3.MaxSize())
		}

		// set MaxSize will be ignore because MaxSize() < Size()
		box4 := NewFrom[int](box3.Items(), WithStrategy(strategy), WithMaxSize(1))
		if box4.Size() != 3 {
			t.Errorf("Expected size is 3, got %d", box4.Size())
		}
		if box4.MaxSize() != 3 {
			t.Errorf("Expected max size is 3, got %d", box4.MaxSize())
		}

		// inheritance MaxSize from previous box
		box5 := NewFromBox[int](box4, WithStrategy(strategy))
		if box5.Size() != 3 {
			t.Errorf("Expected size is 3, got %d", box5.Size())
		}
		if box5.MaxSize() != 3 {
			t.Errorf("Expected max size is 3, got %d", box5.MaxSize())
		}
	}
}
func TestNewFromConcrete(t *testing.T) {
	someItems := []int{1, 2, 3}
	lifoBox := NewLIFOFrom[int](someItems, 1)
	fifoBox := NewFIFOFrom[int](someItems, 1)
	randomBox := NewRandomFrom[int](someItems, 1, rand.New(rand.NewSource(time.Now().UnixNano())))

	lifoItem, _ := lifoBox.Get()
	fifoItem, _ := fifoBox.Get()
	randomItem, _ := randomBox.Get()

	if lifoItem != 3 {
		t.Errorf("Expected lifoItem is 3, got %d", lifoItem)
	}
	if fifoItem != 1 {
		t.Errorf("Expected fifoItem is 1, got %d", fifoItem)
	}
	if !ContainsInt([]int{1, 2, 3}, randomItem) {
		t.Errorf("Expected randomItem is 1 to 3, got %d", randomItem)
	}

	if lifoBox.Size() != 2 {
		t.Errorf("Expected items length only 2, got %d", lifoBox.Size())
	}

	if fifoBox.Size() != 2 {
		t.Errorf("Expected items length only 2, got %d", fifoBox.Size())
	}

	if randomBox.Size() != 2 {
		t.Errorf("Expected items length only 2, got %d", randomBox.Size())
	}

	// lifoBox should be 1, 2
	// fifoBox should be 2, 3
	// randomBox should be between 1, 2, 3 but only contains 2 items
	newFifoBox := NewFIFOFromBox[int](lifoBox, 1)
	newFifoItem, _ := newFifoBox.Get()
	if newFifoItem != 1 {
		t.Errorf("Expected newFifoItem is 1, got %d", newFifoItem)
	}
	if newFifoBox.Size() != 1 {
		t.Errorf("Expected newFifoBox should be 1, got %d", newFifoBox.Size())
	}
	if lifoBox.Size() != 2 {
		t.Errorf("Expected lifoBox should be 2, got %d", lifoBox.Size())
	}

	newLifoBox := NewLIFOFromBox[int](fifoBox, 1)
	newLifoItem, _ := newLifoBox.Get()
	if newLifoItem != 3 {
		t.Errorf("Expected newLifoItem is 3, got %d", newLifoItem)
	}
	if newLifoBox.Size() != 1 {
		t.Errorf("Expected newLifoBox should be 1, got %d", newLifoBox.Size())
	}
	if fifoBox.Size() != 2 {
		t.Errorf("Expected fifoBox should be 2, got %d", fifoBox.Size())
	}

	newRandomBox := NewRandomFromBox[int](fifoBox, 1, rand.New(rand.NewSource(time.Now().UnixNano())))
	newRandomItem, _ := newRandomBox.Get()
	if !ContainsInt([]int{2, 3}, newRandomItem) {
		t.Errorf("Expected newRandomItem either 2 or 3, got %d", newRandomItem)
	}
	if newRandomBox.Size() != 1 {
		t.Errorf("Expected newRandomBox should be 1, got %d", newRandomBox.Size())
	}
	if fifoBox.Size() != 2 {
		t.Errorf("Expected fifoBox should be 2, got %d", fifoBox.Size())
	}
}

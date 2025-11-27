package blackbox

import (
	"slices"
	"testing"
)

func TestFIFOGrowWithZero(t *testing.T) {
	b := NewFIFO[int](0, 0)
	for i := 0; i < 8; i++ {
		b.Put(i)
	}
	if b.Size() != 8 {
		t.Fatalf("expected to be 8, got %d", b.Size())
	}
}

func TestFIFOGrowCopiesContiguousRangeWhenHeadLessThanTail(t *testing.T) {
	// Create a fifo with capacity 8 and populate it with distinct values.
	b := NewFIFO[int](0, 8)
	for i := 0; i < 8; i++ {
		b.items[i] = i
	}

	// Arrange a contiguous region: head < tail
	b.head = 1
	b.tail = 5
	b.size = 4

	oldLen := len(b.items)
	// Call grow directly to exercise the head < tail branch.
	b.grow()

	// After grow, head should be reset to 0 and tail should equal size.
	if b.head != 0 {
		t.Fatalf("expected head to be 0 after grow, got %d", b.head)
	}
	if b.tail != b.size {
		t.Fatalf("expected tail to be %d after grow, got %d", b.size, b.tail)
	}

	// New capacity should be oldLen * growthFactor (no maxSize restriction in this case).
	expectedCap := oldLen * growthFactor
	if len(b.items) != expectedCap {
		t.Fatalf("expected capacity %d, got %d", expectedCap, len(b.items))
	}

	// Verify the items were copied in order from b.items[1:5].
	want := []int{1, 2, 3, 4}
	got := b.items[:b.size]
	if !slices.Equal(got, want) {
		t.Fatalf("items mismatch: want %v got %v", want, got)
	}
}

func TestFIFOGrowRespectsMaxSizeAndCopiesWrapAround(t *testing.T) {
	// Create a fifo with initial capacity 8 and a maxSize that will limit growth.
	max := 5
	b := NewFIFO[int](max, 8)
	for i := 0; i < 8; i++ {
		b.items[i] = i
	}

	// Arrange a wrapped buffer: head > tail
	// For head=6, tail=3 and len=8, size should be (8-6)+3 = 5.
	b.head = 6
	b.tail = 3
	b.size = (len(b.items) - b.head) + b.tail // 2 + 3 = 5

	// Call grow directly to exercise both the wrap-copy branch and the maxSize cap branch.
	b.grow()

	// After grow, capacity should be equal to max (since newCapacity would exceed max).
	if len(b.items) != max {
		t.Fatalf("expected capacity to be capped at maxSize %d, got %d", max, len(b.items))
	}

	// head should be reset to 0 and tail set to size
	if b.head != 0 {
		t.Fatalf("expected head to be 0 after grow, got %d", b.head)
	}
	if b.tail != b.size {
		t.Fatalf("expected tail to be %d after grow, got %d", b.size, b.tail)
	}

	// Verify the wrapped items were copied in correct order:
	// original sequence (from head=6, tail=3): [6,7,0,1,2]
	want := []int{6, 7, 0, 1, 2}
	got := b.items[:b.size]
	if !slices.Equal(got, want) {
		t.Fatalf("wrapped copy mismatch: want %v got %v", want, got)
	}
}

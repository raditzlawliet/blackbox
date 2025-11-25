package blackbox

type fifoBox[T any] struct {
	items   []T
	head    int
	tail    int
	size    int
	maxSize int
}

// NewFIFO creates a new FIFO blackbox with the specified maximum size and capacity.
// Returns a concrete instance of lifo blackbox without interface.
func NewFIFO[T any](maxSize, capacity int) *fifoBox[T] {
	return &fifoBox[T]{
		items:   make([]T, capacity),
		head:    0,
		tail:    0,
		size:    0,
		maxSize: maxSize,
	}
}

func (b *fifoBox[T]) grow() {
	newCapacity := len(b.items) * growthFactor
	if b.maxSize > 0 && newCapacity > b.maxSize {
		newCapacity = b.maxSize
	}

	newItems := make([]T, newCapacity)

	if b.head < b.tail {
		copy(newItems, b.items[b.head:b.tail])
	} else {
		n := copy(newItems, b.items[b.head:])
		copy(newItems[n:], b.items[:b.tail])
	}
	b.head = 0
	b.tail = b.size
	b.items = newItems
}

func (b *fifoBox[T]) Put(item T) error {
	if b.maxSize > 0 && b.size >= b.maxSize {
		return ErrBlackBoxFull
	}

	if b.size >= len(b.items) {
		b.grow()
	}

	b.items[b.tail] = item
	b.tail = (b.tail + 1) % len(b.items)
	b.size++
	return nil
}

func (b *fifoBox[T]) Get() (T, error) {
	if b.size == 0 {
		var zero T
		return zero, ErrEmptyBlackBox
	}

	item := b.items[b.head]
	var zero T
	b.items[b.head] = zero
	b.head = (b.head + 1) % len(b.items)
	b.size--
	return item, nil
}

func (b *fifoBox[T]) Peek() (T, error) {
	if b.size == 0 {
		var zero T
		return zero, ErrEmptyBlackBox
	}
	return b.items[b.head], nil
}

func (b *fifoBox[T]) Size() int {
	return b.size
}

func (b *fifoBox[T]) MaxSize() int {
	return b.maxSize
}

func (b *fifoBox[T]) IsFull() bool {
	return b.maxSize > 0 && b.size >= b.maxSize
}

func (b *fifoBox[T]) IsEmpty() bool {
	return b.size == 0
}

func (b *fifoBox[T]) Clean() {
	var zero T
	for i := 0; i < b.size; i++ {
		idx := (b.head + i) % len(b.items)
		b.items[idx] = zero
	}
	b.head = 0
	b.tail = 0
	b.size = 0
}

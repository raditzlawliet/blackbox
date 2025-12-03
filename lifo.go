package blackbox

type lifoBox[T any] struct {
	items   []T
	maxSize int
}

// NewLIFO creates a new LIFO blackbox with the specified maximum size and capacity.
// Returns a concrete instance of lifo blackbox without interface.
func NewLIFO[T any](maxSize, capacity int) *lifoBox[T] {
	return &lifoBox[T]{
		items:   make([]T, 0, capacity),
		maxSize: maxSize,
	}
}

// NewLIFOFrom creates a new LIFO blackbox from a slice of items with the specified maximum size.
// items are copied so it safe to use the original slice after the blackbox is created.
func NewLIFOFrom[T any](items []T, maxSize int) *lifoBox[T] {
	if maxSize != 0 && maxSize < len(items) {
		maxSize = len(items)
	}
	newItems := make([]T, len(items))
	copy(newItems, items)
	return &lifoBox[T]{
		items:   newItems,
		maxSize: maxSize,
	}
}

// NewLIFOFromBox creates a new LIFO blackbox from a blackbox with the specified maximum size.
// items are copied so it safe to use the original blackbox after the blackbox is created.
func NewLIFOFromBox[T any](box BlackBox[T]) *lifoBox[T] {
	return &lifoBox[T]{
		items:   box.Items(),
		maxSize: box.MaxSize(),
	}
}

func (b *lifoBox[T]) Put(item T) error {
	if b.maxSize > 0 && len(b.items) >= b.maxSize {
		return ErrBlackBoxFull
	}
	b.items = append(b.items, item)
	return nil
}

func (b *lifoBox[T]) Get() (T, error) {
	if len(b.items) == 0 {
		var zero T
		return zero, ErrEmptyBlackBox
	}
	lastIdx := len(b.items) - 1
	item := b.items[lastIdx]
	b.items = b.items[:lastIdx]
	return item, nil
}

func (b *lifoBox[T]) Peek() (T, error) {
	if len(b.items) == 0 {
		var zero T
		return zero, ErrEmptyBlackBox
	}
	return b.items[len(b.items)-1], nil
}

func (b *lifoBox[T]) Size() int {
	return len(b.items)
}

func (b *lifoBox[T]) MaxSize() int {
	return b.maxSize
}

func (b *lifoBox[T]) IsFull() bool {
	return b.maxSize > 0 && len(b.items) >= b.maxSize
}

func (b *lifoBox[T]) IsEmpty() bool {
	return len(b.items) == 0
}

func (b *lifoBox[T]) Clean() {
	b.items = b.items[:0]
}

func (b *lifoBox[T]) Items() []T {
	items := make([]T, len(b.items))
	copy(items, b.items)
	return items
}

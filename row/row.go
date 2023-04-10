package row

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Row is a generic item container with thread-safe item rotation.
type Row[T any] struct {
	items []T
	pivot int

	limiter rate.Sometimes
	lock    sync.Mutex
}

type OptFunc[T any] func(row *Row[T])

// WithInterval set the liminal interval between two effective rotation.
//
// Default to 100 millisecond.
func WithInterval[T any](value time.Duration) OptFunc[T] {
	return func(row *Row[T]) {
		row.limiter = rate.Sometimes{
			Interval: value,
		}
	}
}

// FromSlice create a new Row container with the given items.
//
// The items slice is copied into internal state.
func FromSlice[T any](items []T, opts ...OptFunc[T]) *Row[T] {
	row := &Row[T]{
		pivot: 0,

		limiter: rate.Sometimes{
			Interval: 100 * time.Millisecond,
		},
		lock:  sync.Mutex{},
		items: make([]T, len(items)),
	}
	// make a copy
	copy(row.items, items)

	for _, opt := range opts {
		opt(row)
	}

	return row
}

// Items returns a channel of items in the Row container.
func (r *Row[T]) Items() <-chan T {
	size := len(r.items)
	out := make(chan T, len(r.items))
	for idx, count := r.pivot, 0; count < size; idx, count = ((idx + 1) % size), count+1 {
		out <- r.items[idx]
	}
	close(out)
	return out
}

// Rotate makes the internal state of the Row to rotate left by one.
//
// At most one rotate can take effect in pre-defined interval.
// For rest attempting of rotate, nothing will happen.
func (r *Row[T]) Rotate() {
	r.limiter.Do(func() {
		r.lock.Lock()
		defer r.lock.Unlock()

		r.pivot = (r.pivot + 1) % len(r.items)
	})
}
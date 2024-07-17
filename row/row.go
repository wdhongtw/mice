package row

import (
	"context"
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
// Zero interval means there is no limit.
// Panics on negative value.
func WithInterval[T any](value time.Duration) OptFunc[T] {
	return func(row *Row[T]) {
		row.limiter = rate.Sometimes{Interval: value}
		if value < 0 {
			panic("got negative interval value")
		}
		if value == 0 {
			row.limiter = rate.Sometimes{Every: 1}
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

// Takes returns a channel of items in the Row container.
//
// Caller MUST cancel the context to avoid resource leak.
// This interface require additional routine when used, but unlike Items and ToSlice,
// the size of allocated memory won't grows linearly with the total item count,
// and can run faster.
//
// Consider use this interface when there are lots of items.
func (r *Row[T]) Takes(ctx context.Context) <-chan T {
	size := len(r.items)
	out := make(chan T)
	go func() {
		defer close(out)
		for idx, count := r.pivot, 0; count < size; idx, count = ((idx + 1) % size), count+1 {
			select {
			case <-ctx.Done():
				return
			default:
			}
			out <- r.items[idx]
		}
	}()
	return out
}

// ToSlice just like Items but returns items in slice.
//
// This would be the fastest interface if the length of items is expected to be short.
func (r *Row[T]) ToSlice() []T {
	size := len(r.items)
	out := make([]T, size)
	for idx, count := r.pivot, 0; count < size; idx, count = ((idx + 1) % size), count+1 {
		out[count] = r.items[idx]
	}
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

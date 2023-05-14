package async

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type FutureFunc[T any] func(ctx context.Context) (T, error)

// FutureGroup is a container that collect execution result from multiple jobs.
type FutureGroup[T any] struct {
	jobs []FutureFunc[T]
}

// NewGroup create a FutureGroup
func NewGroup[T any]() *FutureGroup[T] {
	return &FutureGroup[T]{}
}

// Go registers a job that produce a T on success.
//
// The job is called with a context, which is cancelled automatically if some other job return error.
func (rg *FutureGroup[T]) Go(job FutureFunc[T]) {
	rg.jobs = append(rg.jobs, job)
}

// Wait blocks wait until all jobs complete, returns results and the first error (if any).
//
// Just like sync.WaitGroup and errgroup.Group,
// caller should Add/Launch job and Wait in the same goroutine.
func (rg *FutureGroup[T]) Wait(ctx context.Context) ([]T, error) {
	result := make([]T, len(rg.jobs))
	group, ctx := errgroup.WithContext(ctx)
	for idx, job := range rg.jobs {
		idx, job := idx, job
		group.Go(func() error {
			t, err := job(ctx)

			// safe concurrent assess to allocated slice, tested with go test -race
			result[idx] = t
			return err
		})
	}

	err := group.Wait()
	return result, err
}

package async

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestFutureGroup(t *testing.T) {
	t.Run("GetResultFromAllJobs", func(t *testing.T) {
		rg := NewGroup[int]()
		rg.Go(func(_ context.Context) (int, error) {
			return 1, nil
		})
		rg.Go(func(_ context.Context) (int, error) {
			time.Sleep(1 * time.Millisecond)
			return 2, nil
		})

		result, err := rg.Wait(context.Background())
		if err != nil {
			t.Errorf("wait return error: %v", err)
		}
		if !sliceElemEqual([]int{1, 2}, result) {
			t.Errorf("expected: %v, got %v", []int{1, 2}, result)
		}
	})
	t.Run("PartialErrorCauseFinalError", func(t *testing.T) {
		rg := NewGroup[any]()
		rg.Go(func(_ context.Context) (any, error) {
			return nil, nil
		})
		rg.Go(func(_ context.Context) (any, error) {
			return nil, fmt.Errorf("some-error")
		})

		_, err := rg.Wait(context.Background())
		if err == nil {
			t.Errorf("wait return no error")
		}
	})
	t.Run("EmptyJobsListWorksFine", func(t *testing.T) {
		rg := NewGroup[string]()

		result, err := rg.Wait(context.Background())
		if err != nil {
			t.Errorf("wait return error: %v", err)
		}
		if !sliceElemEqual([]string{}, result) {
			t.Errorf("expected: %v, got %v", []string{}, result)
		}
	})
}

func BenchmarkFutureGroup(b *testing.B) {
	b.Run("FourStringJobSuccess", func(b *testing.B) {
		expected := []string{"", "", "", ""}
		ctx := context.Background()
		for i := 0; i < b.N; i++ {
			rg := NewGroup[string]()
			for i := 0; i < 4; i++ {
				rg.Go(func(ctx context.Context) (string, error) {
					return "", nil
				})
			}

			result, err := rg.Wait(ctx)
			if err != nil {
				b.Errorf("wait return error: %v", err)
			}

			if !sliceElemEqual(expected, result) {
				b.Errorf("expected: %v, got %v", expected, result)
			}
		}
	})
	b.Run("FourAnyJobError", func(b *testing.B) {
		ctx := context.Background()
		for i := 0; i < b.N; i++ {
			rg := NewGroup[any]()
			rg.Go(func(ctx context.Context) (any, error) {
				return nil, fmt.Errorf("some-error")
			})
			for i := 0; i < 4; i++ {
				rg.Go(func(ctx context.Context) (any, error) {
					return nil, nil
				})
			}

			_, err := rg.Wait(ctx)
			if err == nil {
				b.Errorf("wait return no error")
			}
		}
	})
}

func sliceElemEqual[T comparable](left, right []T) bool {
	leftMap := map[T]struct{}{}
	rightMap := map[T]struct{}{}

	for _, item := range left {
		leftMap[item] = struct{}{}
	}
	for _, item := range right {
		rightMap[item] = struct{}{}
	}
	return reflect.DeepEqual(leftMap, rightMap)
}

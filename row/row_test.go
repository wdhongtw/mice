package row

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestRow(t *testing.T) {
	t.Run("ZeroSliceCanUse", func(t *testing.T) {
		row := FromSlice([]int{})

		result := asSlice(row.Items())

		if !reflect.DeepEqual([]int{}, result) {
			t.Fatalf("result is not empty slice")
		}
	})
	t.Run("NilSliceCanUse", func(t *testing.T) {
		row := FromSlice[int](nil)

		result := asSlice(row.Items())

		if !reflect.DeepEqual([]int{}, result) {
			t.Fatalf("result is not empty slice")
		}
	})
	t.Run("RotateCanWork", func(t *testing.T) {
		row := FromSlice([]int{1, 2, 3}, WithInterval[int](1*time.Second))

		row.Rotate()

		result := asSlice(row.Items())
		if !reflect.DeepEqual([]int{2, 3, 1}, result) {
			t.Fatalf("result is not expected, got %+v", result)
		}
	})
	t.Run("MultipleRotateCompress", func(t *testing.T) {
		row := FromSlice([]int{1, 2, 3}, WithInterval[int](1*time.Second))

		wg := sync.WaitGroup{}
		for idx := 0; idx < 100; idx++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				row.Rotate()
			}()
		}
		wg.Wait()

		result := asSlice(row.Items())
		if !reflect.DeepEqual([]int{2, 3, 1}, result) {
			t.Fatalf("result is not expected, got %+v", result)
		}
	})
	t.Run("PointerTypeIsFine", func(t *testing.T) {
		name := "alice"
		row := FromSlice([]*string{&name})

		result := asSlice(row.Items())
		if result[0] != &name {
			t.Fatalf("result is not expected, got %+v", result)
		}
	})
}

func asSlice[T any](items <-chan T) []T {
	result := []T{}
	for item := range items {
		result = append(result, item)
	}

	return result
}

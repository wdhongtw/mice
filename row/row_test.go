package row

import (
	"context"
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
	t.Run("RotateWorksFine", func(t *testing.T) {
		row := FromSlice([]int{1, 2, 3}, WithInterval[int](0*time.Second))

		row.Rotate()
		row.Rotate()
		row.Rotate()

		result := asSlice(row.Items())
		if !reflect.DeepEqual([]int{1, 2, 3}, result) {
			t.Fatalf("result is not expected, got %+v", result)
		}
	})
	t.Run("ToSliceWorksFine", func(t *testing.T) {
		row := FromSlice([]int{1, 2, 3}, WithInterval[int](0*time.Second))

		row.Rotate()
		row.Rotate()

		result := row.ToSlice()
		if !reflect.DeepEqual([]int{3, 1, 2}, result) {
			t.Fatalf("result is not expected, got %+v", result)
		}
	})
	t.Run("TakesAlsoWorksFine", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		row := FromSlice([]int{1, 2, 3}, WithInterval[int](0*time.Second))

		result := asSlice(row.Takes(ctx))
		if !reflect.DeepEqual([]int{1, 2, 3}, result) {
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

func BenchmarkRowShort(b *testing.B) {
	type empty struct{}

	items := []*empty{&(empty{}), &(empty{})}
	row := FromSlice(items)
	b.Run("TestItemsPointerFirst", func(b *testing.B) {
		for idx := 0; idx < b.N; idx++ {
			<-row.Items()
		}
	})
	b.Run("TestTakesPointerFirst", func(b *testing.B) {
		for idx := 0; idx < b.N; idx++ {
			go func() {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				<-row.Takes(ctx)
			}()
		}
	})
	b.Run("TestToSlicePointerFirst", func(b *testing.B) {
		for idx := 0; idx < b.N; idx++ {
			_ = row.ToSlice()[0]
		}
	})
	b.Run("TestRawPointerFirst", func(b *testing.B) {
		for idx := 0; idx < b.N; idx++ {
			_ = items[0]
		}
	})
}

func BenchmarkRowLong(b *testing.B) {
	type empty struct{}

	items := []*empty{}
	for idx := 0; idx < 0x400; idx++ {
		items = append(items, &empty{})
	}
	row := FromSlice(items)
	b.Run("TestItemsPointerFirst", func(b *testing.B) {
		for idx := 0; idx < b.N; idx++ {
			<-row.Items()
		}
	})
	b.Run("TestTakesPointerFirst", func(b *testing.B) {
		for idx := 0; idx < b.N; idx++ {
			go func() {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				<-row.Takes(ctx)
			}()
		}
	})
	b.Run("TestToSlicePointerFirst", func(b *testing.B) {
		for idx := 0; idx < b.N; idx++ {
			_ = row.ToSlice()[0]
		}
	})
	b.Run("TestRawPointerFirst", func(b *testing.B) {
		for idx := 0; idx < b.N; idx++ {
			_ = items[0]
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

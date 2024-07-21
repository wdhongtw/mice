package flow

import "iter"

// Empty yield no value.
func Empty[V any]() iter.Seq[V] {
	return func(yield func(V) bool) {
	}
}

func Empty2[K, V any]() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
	}
}

// Pack yield one value.
func Pack[V any](v V) iter.Seq[V] {
	return func(yield func(V) bool) {
		yield(v)
	}
}

func Pack2[K, V any](k K, v V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		yield(k, v)
	}
}

// Chain chains multiple sequences together.
//
// Only for demonstration, consider "x/exp/xiter.Concat" once the official one is available.
func Chain[V any](sequences ...iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, sequence := range sequences {
			for val := range sequence {
				if !yield(val) {
					return
				}
			}
		}
	}
}

func Chain2[K, V any](sequences ...iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, sequence := range sequences {
			for key, val := range sequence {
				if !yield(key, val) {
					return
				}
			}
		}
	}
}

// Take takes first N items from the sequence, yield all items if exhausted.
func Take[V any](sequence iter.Seq[V], count int) iter.Seq[V] {
	return func(yield func(V) bool) {
		c := 0
		for val := range sequence {
			c += 1
			if c > count {
				return
			}
			if !yield(val) {
				return
			}
		}
	}
}

func Take2[K, V any](sequence iter.Seq2[K, V], count int) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		c := 0
		for key, val := range sequence {
			c += 1
			if c > count {
				return
			}
			if !yield(key, val) {
				return
			}
		}
	}
}

// Drop drops first N items from the sequence, yield no item if exhausted.
func Drop[V any](sequence iter.Seq[V], count int) iter.Seq[V] {
	return func(yield func(V) bool) {
		c := 0
		for val := range sequence {
			c += 1
			if c <= count {
				continue
			}
			if !yield(val) {
				return
			}
		}
	}
}

func Drop2[K, V any](sequence iter.Seq2[K, V], count int) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		c := 0
		for key, val := range sequence {
			c += 1
			if c <= count {
				continue
			}
			if !yield(key, val) {
				return
			}
		}
	}
}

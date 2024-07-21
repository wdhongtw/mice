package flow

import "iter"

// Since yields infinite sequence of numbers from "start".
func Since(start int) iter.Seq[int] {
	return func(yield func(int) bool) {
		val := start
		for {
			if !yield(val) {
				return
			}
			val += 1
		}
	}
}

// Forward yields [start, end) in normal order, yields no values if start >= end.
func Forward(begin int, end int) iter.Seq[int] {
	return func(yield func(int) bool) {
		idx := begin
		for {
			if idx >= end {
				return
			}
			if !yield(idx) {
				return
			}
			idx += 1
		}
	}
}

// Backward yields [start, end) in reversed order, yields no values if start >= end.
func Backward(begin int, end int) iter.Seq[int] {
	return func(yield func(int) bool) {
		idx := end
		for {
			if idx <= begin {
				return
			}
			if !yield(idx - 1) {
				return
			}
			idx -= 1
		}
	}
}

package flow

import "iter"

// Any return c1 || c2 || ... || cn, return false if input is empty.
func Any(conditions iter.Seq[bool]) bool {
	for condition := range conditions {
		if condition {
			return true
		}
	}
	return false
}

// Any return c1 && c2 && ... && cn, return true if input is empty.
func All(conditions iter.Seq[bool]) bool {
	for condition := range conditions {
		if !condition {
			return false
		}
	}
	return true
}

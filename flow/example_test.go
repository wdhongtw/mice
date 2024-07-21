package flow

import (
	"fmt"
	"iter"
)

func Example() {
	type Node struct {
		val   int
		left  *Node
		right *Node
	}

	var traverse func(node *Node) iter.Seq[int]
	traverse = func(node *Node) iter.Seq[int] {
		// Empty is useful as base case during recursive generator chaining.
		if node == nil {
			return Empty[int]()
		}
		// Pack is useful to promote a single value into a iterable for chaining.
		return Chain(
			traverse(node.left),
			Pack(node.val),
			traverse(node.right),
		)
	}

	root := &Node{
		val: 3,
		left: &Node{
			val:  2,
			left: &Node{val: 1},
		},
		right: &Node{
			val:  5,
			left: &Node{val: 4},
		},
	}

	var results []int
	for val := range traverse(root) {
		results = append(results, val)
		if val == 4 {
			break
		}
	}
	fmt.Printf("%v\n", results)
	// Output: [1 2 3 4]
}

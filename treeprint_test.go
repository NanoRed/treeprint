package treeprint

import (
	"math/rand"
	"testing"
	"time"
)

func TestPrint(t *testing.T) {

	getNumber := func() []int {
		number := make([]int, 0)
		for i := 0; i < 20; i++ {
			rand.Seed(time.Now().UnixNano())
			number = append(number, rand.Intn(100000))
		}
		t.Logf("\n%v", number)
		return number
	}

	// binary search tree test
	tree := &BinarySearchTree{}
	for key, val := range getNumber() {
		tree.Append(val, key)
	}
	t.Logf("\n%s", Sprint(tree.root))

	// TODO() multi leaf node tree test
}

// BinarySearchTree binary search tree for testing
type BinarySearchTree struct {
	root *BinarySearchNode
}

// BinarySearchNode binary search tree node
type BinarySearchNode struct {
	Key   int
	Value interface{}
	Left  *BinarySearchNode
	Right *BinarySearchNode
}

// Append append a new node to the tree
func (t *BinarySearchTree) Append(key int, val interface{}) {
	update(&t.root, key, val)
}

func update(node **BinarySearchNode, key int, val interface{}) {
	if *node == nil {
		*node = &BinarySearchNode{
			Key:   key,
			Value: val,
		}
		return
	} else if key == (*node).Key {
		(*node).Value = val
		return
	}
	if key < (*node).Key {
		update(&(*node).Left, key, val)
	} else {
		update(&(*node).Right, key, val)
	}
}

// GetKey implement treeprint
func (n *BinarySearchNode) GetKey() interface{} {
	return n.Key
}

// GetValue implement treeprint
func (n *BinarySearchNode) GetValue() interface{} {
	return ""
}

// RangeNode implement treeprint
func (n *BinarySearchNode) RangeNode() chan TreeNode {
	c := make(chan TreeNode, 2)
	c <- n.Left
	c <- n.Right
	close(c)
	return c
}

package btreeprint

import (
	"math/rand"
	"testing"
	"time"
)

func TestPrint(t *testing.T) {
	number := make([]int, 0)
	for i := 0; i < 20; i++ {
		rand.Seed(time.Now().UnixNano())
		number = append(number, rand.Intn(100000))
	}

	tree := &Tree{}
	for key, val := range number {
		tree.Append(val, key)
	}

	t.Logf("\n%v", number)
	t.Logf("\n%s", Sprint(tree.root))
}

// Tree a simple binary search tree for testing
type Tree struct {
	root *Node
}

// Node tree node
type Node struct {
	Key   int
	Value interface{}
	Left  *Node
	Right *Node
}

// Append append a new node to the tree
func (t *Tree) Append(key int, val interface{}) {
	update(&t.root, key, val)
}

func update(node **Node, key int, val interface{}) {
	if *node == nil {
		*node = &Node{
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

// GetKey implement interface
func (n *Node) GetKey() int {
	return n.Key
}

// GetValue implement interface
func (n *Node) GetValue() interface{} {
	return n.Value
}

// GetLeftNode implement interface
func (n *Node) GetLeftNode() BtreeNode {
	return n.Left
}

// GetRightNode implement interface
func (n *Node) GetRightNode() BtreeNode {
	return n.Right
}

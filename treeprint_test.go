package treeprint

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestPrint(t *testing.T) {

	// binary tree test
	number := make([]int, 0)
	for i := 0; i < 20; i++ {
		rand.Seed(time.Now().UnixNano())
		number = append(number, rand.Intn(100000))
	}
	tree := &BinaryTree{}
	for _, val := range number {
		tree.Append(val, nil)
	}
	t.Logf("\n%s", Sprint(tree.root))

	// multi leaf node tree test
	tree2 := &MultiLeafTree{}
	rand.Seed(time.Now().UnixNano())
	_ = tree2.Append([]int{}, rand.Intn(100000))
	for i := 0; i < 5; i++ {
		_ = tree2.Append([]int{i}, rand.Intn(100000))
	}
	layer2Index := rand.Intn(5)
	for i := 0; i < 5; i++ {
		_ = tree2.Append([]int{layer2Index, i}, rand.Intn(100000))
	}
	layer3Index := rand.Intn(5)
	for i := 0; i < 5; i++ {
		_ = tree2.Append([]int{layer2Index, layer3Index, i}, rand.Intn(100000))
	}
	t.Logf("\n%s", Sprint(tree2.root))
}


// ------------------------------
// binary tree struct for testing
// ------------------------------

// BinaryTree tree structure
type BinaryTree struct {
	root  *BinaryTreeNode
	count uint
}

// BinaryTreeNode tree node structure
type BinaryTreeNode struct {
	Key    int
	Value  interface{}
	Height uint
	Left   *BinaryTreeNode
	Right  *BinaryTreeNode
}

// Append append a new node to the tree
func (t *BinaryTree) Append(key int, val interface{}) {

	// nodes that need to be increased in height
	increaseHeight := make(map[*BinaryTreeNode]struct{})
	increaseTracking := func(parent *BinaryTreeNode, selected *BinaryTreeNode, another *BinaryTreeNode) {
		// if selected leaf node height equal or higher than the another one, than it means has to increase
		if another == nil || (selected != nil && selected.Height >= another.Height) {
			increaseHeight[parent] = struct{}{}
		} else {
			increaseHeight = make(map[*BinaryTreeNode]struct{})
		}
	}

	// search node
	current := &t.root
	for *current != nil {
		if key < (*current).Key {
			increaseTracking((*current), (*current).Left, (*current).Right)
			current = &(*current).Left
		} else if key > (*current).Key  {
			increaseTracking((*current), (*current).Right, (*current).Left)
			current = &(*current).Right
		} else {
			break
		}
	}
	if *current == nil {
		*current = &BinaryTreeNode{
			Key:    key,
			Value:  val,
			Height: 1,
		}
		t.count++
		for node := range increaseHeight {
			node.Height++
		}
	} else {
		(*current).Value = val
	}
}

// GetKey implement treeprint
func (n *BinaryTreeNode) GetKey() interface{} {
	return n.Key
}

// GetValue implement treeprint
func (n *BinaryTreeNode) GetValue() interface{} {
	return fmt.Sprintf("(h:%d)", n.Height) // n.Value
}

// RangeNode implement treeprint
func (n *BinaryTreeNode) RangeNode() chan TreeNode {
	c := make(chan TreeNode, 2)
	c <- n.Left
	c <- n.Right
	close(c)
	return c
}


// ----------------------------------
// multi-node tree struct for testing
// ----------------------------------

// MultiLeafTree multi-node tree structure
type MultiLeafTree struct {
	root  *MultiLeafTreeNode
	count uint
}

// MultiLeafTreeNode multi-node tree node structure
type MultiLeafTreeNode struct {
	Value interface{}
	Leaf []*MultiLeafTreeNode
}

// Append append a new node to the tree by index
// pass []int{} to index means root node
// pass []int{0} means the 1st leaf node of root
// pass []int{0, 1} means the 2nd leaf node of the 1st leaf node of root
// etc.
func (t *MultiLeafTree) Append(index []int, val interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("index out of range")
		}
	}()
	entry := &t.root
	layer := len(index)
	for i := 0; i < layer; i++ {
		if i == layer-1 {
			if index[i] == len((*entry).Leaf) {
				(*entry).Leaf = append((*entry).Leaf, nil)
			}
		}
		entry = &(*entry).Leaf[index[i]]
	}
	if *entry == nil {
		*entry = &MultiLeafTreeNode{
			Value: val,
		}
		t.count++
	} else {
		(*entry).Value = val
	}
	return
}

// GetKey implement treeprint
func (n *MultiLeafTreeNode) GetKey() interface{} {
	return n.Value
}

// GetValue implement treeprint
func (n *MultiLeafTreeNode) GetValue() interface{} {
	return ""
}

// RangeNode implement treeprint
func (n *MultiLeafTreeNode) RangeNode() chan TreeNode {
	count := len(n.Leaf)
	c := make(chan TreeNode, count)
	for i := 0; i < count; i++ {
		c <- n.Leaf[i]
	}
	close(c)
	return c
}
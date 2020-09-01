package btreeprint

import (
	"fmt"
	"reflect"
	"strings"
)

// BtreeNode binary tree interface
type BtreeNode interface {
	GetKey() int
	GetValue() interface{}
	GetLeftNode() BtreeNode
	GetRightNode() BtreeNode
}

// Sprint get print text
func Sprint(entryNode BtreeNode) string {
	rt := reflect.TypeOf(entryNode)
	rv := reflect.ValueOf(entryNode)
	if rt.Kind() != reflect.Ptr {
		return "<error: entryNode has to be a node pointer>"
	} else if rv.IsNil() {
		return "<nil>"
	}
	type nodeInfo struct {
		node       BtreeNode
		layer      int
		count      int
		index      int
		len        int
		str        string
		leftNode   *nodeInfo
		rightNode  *nodeInfo
		parentNode *nodeInfo
	}
	layer := make([][]*nodeInfo, 0)
	first := &nodeInfo{
		node:  entryNode,
		layer: 1,
		count: 1,
		str:   fmt.Sprintf("%d(%v)", entryNode.GetKey(), entryNode.GetValue()),
	}
	first.len = len(first.str)
	queue := []*nodeInfo{first}
	currentLayer := 2
	currentIndex := 0
	currentCount := 1
	i := 0
	for {
		if i > len(queue)-1 {
			break
		}
		current := queue[i]
		leftLeafNode := current.node.GetLeftNode()
		if !reflect.ValueOf(leftLeafNode).IsNil() {
			leftNodeInfo := &nodeInfo{
				node:       leftLeafNode,
				layer:      current.layer + 1,
				str:        fmt.Sprintf("%d(%v)", leftLeafNode.GetKey(), leftLeafNode.GetValue()),
				parentNode: current,
			}
			leftNodeInfo.len = len(leftNodeInfo.str)
			if leftNodeInfo.layer != currentLayer {
				currentLayer = leftNodeInfo.layer
				currentIndex = 0
				currentCount = 1
			}
			leftNodeInfo.count = currentCount
			currentCount++
			leftNodeInfo.index = currentIndex
			currentIndex += leftNodeInfo.len + 1
			queue = append(queue, leftNodeInfo)
			current.leftNode = leftNodeInfo
		}
		rightLeafNode := current.node.GetRightNode()
		if !reflect.ValueOf(rightLeafNode).IsNil() {
			rightNodeInfo := &nodeInfo{
				node:       rightLeafNode,
				layer:      current.layer + 1,
				str:        fmt.Sprintf("%d(%v)", rightLeafNode.GetKey(), rightLeafNode.GetValue()),
				parentNode: current,
			}
			rightNodeInfo.len = len(rightNodeInfo.str)
			if rightNodeInfo.layer != currentLayer {
				currentLayer = rightNodeInfo.layer
				currentIndex = 0
				currentCount = 1
			}
			rightNodeInfo.count = currentCount
			currentCount++
			rightNodeInfo.index = currentIndex
			currentIndex += rightNodeInfo.len + 1
			queue = append(queue, rightNodeInfo)
			current.rightNode = rightNodeInfo
		}
		if current.layer > len(layer) {
			layer = append(layer, make([]*nodeInfo, 0))
		}
		layer[current.layer-1] = append(layer[current.layer-1], current)
		i++
	}
	var alignLeft func(*nodeInfo)
	var alignRight func(*nodeInfo)
	alignLeft = func(current *nodeInfo) {
		if current.leftNode == nil {
			return
		} else if val := current.index - current.leftNode.index; val > 0 { // 下一层移位
			for k := current.leftNode.count - 1; k < len(layer[current.layer]); k++ {
				layer[current.layer][k].index += val
			}
		} else if val < 0 { // 本层移位
			for k := current.count - 1; k < len(layer[current.layer-1]); k++ {
				layer[current.layer-1][k].index += -val
			}
		}
		for j := current.leftNode.count - 1; j < len(layer[current.layer]); j++ {
			alignLeft(layer[current.layer][j])
			alignRight(layer[current.layer][j])
		}
	}
	alignRight = func(current *nodeInfo) {
		if current.rightNode == nil {
			return
		} else if !strings.Contains(current.str, "-+") {
			if val := current.rightNode.index - current.index - current.len; val > 0 {
				tmp := ""
				for k := 0; k < val; k++ {
					tmp += "-"
				}
				tmp += "+"
				current.str += tmp
				newLen := len(current.str)
				offset := newLen - current.len
				current.len = newLen
				for k := current.count; k < len(layer[current.layer-1]); k++ {
					layer[current.layer-1][k].index += offset
				}
			} else {
				tmp := "-+"
				current.str += tmp
				newLen := len(current.str)
				offset := newLen - current.len
				current.len = newLen
				for k := current.count; k < len(layer[current.layer-1]); k++ {
					layer[current.layer-1][k].index += offset
				}
				offset2 := -val + 1
				for k := current.rightNode.count - 1; k < len(layer[current.layer]); k++ {
					layer[current.layer][k].index += offset2
				}
			}
		} else {
			for k := current.rightNode.count - 1; k < len(layer[current.layer]); k++ {
				layer[current.layer][k].index += current.index + current.len - 1 - current.rightNode.index
			}
		}
		for j := current.rightNode.count - 1; j < len(layer[current.layer]); j++ {
			alignLeft(layer[current.layer][j])
			alignRight(layer[current.layer][j])
		}
	}
	for i := len(layer) - 2; i >= 0; i-- {
		for j := 0; j < len(layer[i]); j++ {
			alignLeft(layer[i][j])
			alignRight(layer[i][j])
		}
	}
	text := ""
	for i := 0; i < len(layer); i++ {
		line1 := ""
		line2 := ""
		for j := 0; j < len(layer[i]); j++ {
			if val := layer[i][j].index - len(line1); val > 0 {
				for k := 0; k < val; k++ {
					line1 += " "
					line2 += " "
				}
			}
			line1 += layer[i][j].str
			if layer[i][j].leftNode != nil {
				line2 += "|"
			} else {
				line2 += " "
			}
			for k := 0; k < layer[i][j].len-2; k++ {
				line2 += " "
			}
			if layer[i][j].rightNode != nil {
				line2 += "|"
			} else {
				line2 += " "
			}
		}
		line1 = strings.ReplaceAll(line1, "-", "─")
		line1 = strings.ReplaceAll(line1, "+", "┐")
		text += line1 + "\n"
		text += line2 + "\n"
	}
	return text
}

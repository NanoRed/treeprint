package treeprint

import (
	"fmt"
	"reflect"
	"strings"
)

// TreeNode binary tree interface
type TreeNode interface {
	GetKey() interface{}
	GetValue() interface{}
	RangeNode() chan TreeNode
}

// Sprint get print text
func Sprint(entryNode TreeNode) string {
	rt := reflect.TypeOf(entryNode)
	rv := reflect.ValueOf(entryNode)
	if rt.Kind() != reflect.Ptr {
		return "<error: entryNode has to be a node pointer>"
	} else if rv.IsNil() {
		return "<nil>"
	}
	type nodeInfo struct {
		node       TreeNode
		layer      int
		count      int
		index      int
		len        int
		str        string
		str2       string
		hasBar     bool
		segLen     map[int]int
		parentNode *nodeInfo
		leafNode   []*nodeInfo
	}
	layer := make([][]*nodeInfo, 0)
	entryNodeInfo := &nodeInfo{
		node:   entryNode,
		layer:  1,
		count:  1,
		str:    fmt.Sprintf("%v(%v)", entryNode.GetKey(), entryNode.GetValue()),
		segLen: make(map[int]int),
	}
	entryNodeInfo.len = len(entryNodeInfo.str)
	for i := 0; i < entryNodeInfo.len; i++ {
		entryNodeInfo.str2 += " "
	}
	queue := []*nodeInfo{entryNodeInfo}
	currentLayer := 2
	currentLayerIndex := 0
	currentLayerCount := 1
	i := 0
	for {
		if i > len(queue)-1 {
			break
		}
		current := queue[i]
		leafIndex := 0
		for currentLeafNode := range current.node.RangeNode() {
			if reflect.ValueOf(currentLeafNode).IsNil() {
				current.leafNode = append(current.leafNode, nil)
			} else {
				currentLeafNodeInfo := &nodeInfo{
					node:       currentLeafNode,
					layer:      current.layer + 1,
					str:        fmt.Sprintf("%v(%v)", currentLeafNode.GetKey(), currentLeafNode.GetValue()),
					segLen:     make(map[int]int),
					parentNode: current,
				}
				currentLeafNodeInfo.len = len(currentLeafNodeInfo.str)
				currentLeafNodeInfo.segLen = map[int]int{ 1: currentLeafNodeInfo.len }
				for j := 0; j < currentLeafNodeInfo.len; j++ {
					currentLeafNodeInfo.str2 += " "
				}
				if currentLeafNodeInfo.layer != currentLayer {
					currentLayer = currentLeafNodeInfo.layer
					currentLayerIndex = 0
					currentLayerCount = 1
				}
				currentLeafNodeInfo.count = currentLayerCount
				currentLayerCount++
				currentLeafNodeInfo.index = currentLayerIndex
				currentLayerIndex += currentLeafNodeInfo.len + 1
				queue = append(queue, currentLeafNodeInfo)
				current.leafNode = append(current.leafNode, currentLeafNodeInfo)
				current.segLen[currentLeafNodeInfo.count] = current.len
			}
			leafIndex++
		}
		for ; leafIndex < 2; leafIndex++ {
			current.leafNode = append(current.leafNode, nil)
		}
		if current.layer > len(layer) {
			layer = append(layer, make([]*nodeInfo, 0))
		}
		layer[current.layer-1] = append(layer[current.layer-1], current)
		i++
	}
	var alignFirst func(*nodeInfo)
	var alignOther func(*nodeInfo)
	alignFirst = func(current *nodeInfo) {
		if current.leafNode[0] == nil {
			return
		} else if val := current.index - current.leafNode[0].index; val > 0 { // 下一层移位
			for k := current.leafNode[0].count - 1; k < len(layer[current.layer]); k++ {
				layer[current.layer][k].index += val
			}
		} else if val < 0 { // 本层移位
			for k := current.count - 1; k < len(layer[current.layer-1]); k++ {
				layer[current.layer-1][k].index += -val
			}
		}
		for j := current.leafNode[0].count - 1; j < len(layer[current.layer]); j++ {
			alignFirst(layer[current.layer][j])
			alignOther(layer[current.layer][j])
		}
	}
	alignOther = func(current *nodeInfo) {
		lastLeafIndex := len(current.leafNode) - 1
		for i := 1; i <= lastLeafIndex; i++ {
			if current.leafNode[i] == nil {
				continue
			} else if !current.hasBar {
				if val := current.leafNode[i].index - current.index - current.len; val > 0 {
					for k := 0; k < val; k++ {
						current.str += "-"
						current.str2 += " "
					}
					current.str += "+"
					current.str2 += "|"
					newLen := len(current.str)
					offset := newLen - current.len
					current.len = newLen
					current.segLen[current.leafNode[i].count] = current.len
					for k := current.count; k < len(layer[current.layer-1]); k++ {
						layer[current.layer-1][k].index += offset
					}
				} else {
					current.str += "-+"
					current.str2 += " |"
					newLen := len(current.str)
					offset := newLen - current.len
					current.len = newLen
					current.segLen[current.leafNode[i].count] = current.len
					for k := current.count; k < len(layer[current.layer-1]); k++ {
						layer[current.layer-1][k].index += offset
					}
					offset2 := -val + 1
					for k := current.leafNode[i].count - 1; k < len(layer[current.layer]); k++ {
						layer[current.layer][k].index += offset2
					}
				}
				if i == lastLeafIndex {
					current.hasBar = true
				}
			} else {
				parentIndex := current.index + current.segLen[current.leafNode[i].count] - 1
				selfIndex := current.leafNode[i].index
				gap := parentIndex - selfIndex
				for k := current.leafNode[i].count - 1; k < len(layer[current.layer]); k++ {
					layer[current.layer][k].index += gap
				}
			}
			for j := current.leafNode[i].count - 1; j < len(layer[current.layer]); j++ {
				alignFirst(layer[current.layer][j])
				alignOther(layer[current.layer][j])
			}
		}
	}
	for i := len(layer) - 2; i >= 0; i-- {
		for j := 0; j < len(layer[i]); j++ {
			alignFirst(layer[i][j])
			alignOther(layer[i][j])
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
			if layer[i][j].leafNode[0] != nil {
				line2 += "|"
			} else {
				line2 += " "
			}
			line2 += layer[i][j].str2[1:]
		}
		line1 = strings.ReplaceAll(line1, "-+-", "─┬─")
		line1 = strings.ReplaceAll(line1, "-+", "─┐")
		line1 = strings.ReplaceAll(line1, "-", "─")
		text += line1 + "\n"
		text += line2 + "\n"
	}
	return text
}

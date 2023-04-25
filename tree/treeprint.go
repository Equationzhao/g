/*
	The MIT License (MIT)
	Copyright © 2016 Maxim Kupriianov <max@kc.vc>

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the “Software”), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in
	all copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
	THE SOFTWARE.
*/

package tree

import (
	"fmt"
	"github.com/valyala/bytebufferpool"
	"io"
	"strings"
	"sync"
)

// value defines any value
type value = string

// tree represents a tree structure with leaf-nodes and branch-nodes.
type tree interface {
	// AddNode adds a new node to a branch.
	AddNode(v value) tree
	// AddBranch adds a new branch node (a level deeper).
	AddBranch(v value) tree

	// String renders the tree or subtree as a string.
	String() string
}

type node struct {
	Root  *node
	Nodes []*node
	Value value
	m     sync.RWMutex
}

func (n *node) AddNode(v value) tree {
	newNode := &node{
		Root:  n,
		Value: v,
		Nodes: make([]*node, 0, nodeSize),
	}
	n.m.Lock()
	defer n.m.Unlock()
	n.Nodes = append(n.Nodes, newNode)
	return n
}

func (n *node) AddBranch(v value) tree {
	branch := &node{
		Root:  n,
		Value: v,
		Nodes: make([]*node, 0, nodeSize),
	}
	n.m.Lock()
	defer n.m.Unlock()
	n.Nodes = append(n.Nodes, branch)
	return branch
}

func (n *node) String() string {
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)
	level := 0
	var levelsEnded = make(map[int]bool, len(n.Nodes))
	if n.Root == nil {
		_, _ = buf.WriteString(fmt.Sprintf("%v", n.Value))
		_ = buf.WriteByte('\n')
	} else {
		edge := EdgeTypeMid
		if len(n.Nodes) == 0 {
			edge = EdgeTypeEnd
			levelsEnded[level] = true
		}
		printValues(buf, 0, levelsEnded, edge, n)
	}
	if len(n.Nodes) > 0 {
		printNodes(buf, level, levelsEnded, n.Nodes)
	}
	return buf.String()
}

func printNodes(wr io.Writer, level int, levelsEnded map[int]bool, nodes []*node) {

	for i, node := range nodes {
		edge := EdgeTypeMid
		if i == len(nodes)-1 {
			levelsEnded[level] = true // set the value in the map
			edge = EdgeTypeEnd
		}
		printValues(wr, level, levelsEnded, edge, node)
		if len(node.Nodes) > 0 {
			printNodes(wr, level+1, levelsEnded, node.Nodes)
		}
	}
}

func printValues(wr io.Writer, level int, levelsEnded map[int]bool, edge EdgeType, node *node) {

	for i := 0; i < level; i++ {
		if levelsEnded[i] {
			_, _ = wr.Write([]byte(strings.Repeat(" ", IndentSize+1)))
			continue
		}
		_, _ = fmt.Fprintf(wr, "%s%s", EdgeTypeLink, strings.Repeat(" ", IndentSize))
	}

	val := node.Value

	_, _ = fmt.Fprintf(wr, "%s %s\n", edge, val)
}

type EdgeType string

const (
	EdgeTypeLink EdgeType = "│"
	EdgeTypeMid  EdgeType = "├──"
	EdgeTypeEnd  EdgeType = "└──"
)

// IndentSize is the number of spaces per tree level.
const IndentSize = 3

const nodeSize = 20

// New Generates new tree
func New() tree {
	return &node{Value: ".", Nodes: make([]*node, 0, nodeSize)}
}

// NewWithRoot Generates new tree with the given root value
func NewWithRoot(root value) tree {
	return &node{Value: root, Nodes: make([]*node, 0, nodeSize)}
}

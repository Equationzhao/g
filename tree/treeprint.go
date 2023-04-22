package tree

import (
	"fmt"
	"io"
	"strings"
	"sync"
)

// value defines any value
type value = string

// metaValue defines any meta value
type metaValue any

// tree represents a tree structure with leaf-nodes and branch-nodes.
type tree interface {
	// AddNode adds a new node to a branch.
	AddNode(v value) tree
	// AddMetaNode adds a new node with meta value provided to a branch.
	AddMetaNode(meta metaValue, v value) tree
	// AddBranch adds a new branch node (a level deeper).
	AddBranch(v value) tree
	// AddMetaBranch adds a new branch node (a level deeper) with meta value provided.
	AddMetaBranch(meta metaValue, v value) tree
	// Branch converts a leaf-node to a branch-node,
	// applying this on a branch-node does no effect.
	Branch() tree

	// String renders the tree or subtree as a string.
	String() string
}

type node struct {
	Root  *node
	Meta  metaValue
	Value value
	m     sync.Mutex
	Nodes []*node
}

func (n *node) AddNode(v value) tree {
	n.m.Lock()
	defer n.m.Unlock()
	n.Nodes = append(n.Nodes, &node{
		Root:  n,
		Value: v,
	})
	return n
}

func (n *node) AddMetaNode(meta metaValue, v value) tree {
	n.m.Lock()
	defer n.m.Unlock()
	n.Nodes = append(n.Nodes, &node{
		Root:  n,
		Meta:  meta,
		Value: v,
		Nodes: make([]*node, 0, nodeSize),
	})
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

func (n *node) AddMetaBranch(meta metaValue, v value) tree {
	branch := &node{
		Root:  n,
		Meta:  meta,
		Value: v,
		Nodes: make([]*node, 0, nodeSize),
	}
	n.m.Lock()
	defer n.m.Unlock()
	n.Nodes = append(n.Nodes, branch)
	return branch
}

func (n *node) Branch() tree {
	n.Root = nil
	return n
}

func (n *node) String() string {
	buf := new(strings.Builder)
	level := 0
	var levelsEnded = make(map[int]bool, len(n.Nodes))
	if n.Root == nil {
		if n.Meta != nil {
			buf.WriteString(fmt.Sprintf("[%v]  %v", n.Meta, n.Value))
		} else {
			buf.WriteString(fmt.Sprintf("%v", n.Value))
		}
		buf.WriteByte('\n')
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
			_, _ = fmt.Fprint(wr, strings.Repeat(" ", IndentSize+1))
			continue
		}
		_, _ = fmt.Fprintf(wr, "%s%s", EdgeTypeLink, strings.Repeat(" ", IndentSize))
	}

	val := node.Value
	meta := node.Meta

	if meta != nil {
		_, _ = fmt.Fprintf(wr, "%s [%v]  %v\n", edge, meta, val)
		return
	}
	_, _ = fmt.Fprintf(wr, "%s %v\n", edge, val)
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

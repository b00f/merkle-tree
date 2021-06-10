package merkletree

import (
	"math"
)

type Hasher func([]byte) []byte

type Tree struct {
	nodes     map[int]*node
	maxWidth  int
	maxHeight int

	hasher Hasher
}

type node struct {
	width  int
	height int
	hash   []byte
}

// nodeID return the node ID (four bytes):
// +-+---+
// |h| w |
// +-+---+
// h: height
// w: width
//
func nodeID(width, height int) int {
	return ((height & 0xff) << 24) | (width & 0xffffff)
}

func New(hasher Hasher) *Tree {
	return &Tree{
		nodes:  make(map[int]*node),
		hasher: hasher,
	}
}

func (t *Tree) createNode(width, height int) *node {
	return &node{
		width:  width,
		height: height,
	}
}

func (t *Tree) getNode(width, height int) *node {
	id := nodeID(width, height)
	return t.nodes[id]
}

func (t *Tree) getOrCreateNode(width, height int) *node {
	id := nodeID(width, height)
	node, ok := t.nodes[id]
	if !ok {
		node = t.createNode(width, height)
		t.nodes[id] = node
	}
	return node
}

func (t *Tree) invalidateNode(width, height int) {
	n := t.getOrCreateNode(width, height)
	n.hash = nil
}

func (t *Tree) recalculateHeight(maxWidth int) {
	if maxWidth > t.maxWidth {
		t.maxWidth = maxWidth

		maxHeight := math.Log2(float64(maxWidth))
		if math.Remainder(maxHeight, 1.0) != 0 {
			t.maxHeight = int(math.Trunc(maxHeight)) + 2
		} else {
			t.maxHeight = int(math.Trunc(maxHeight)) + 1
		}
	}
}

func (t *Tree) SetBlockData(no int, data []byte) {
	t.recalculateHeight(no + 1)

	h := t.hasher(data)
	node := t.getOrCreateNode(no, 0)
	node.hash = h

	w := no / 2
	for h := 1; h < t.maxHeight; h++ {
		t.invalidateNode(w, h)
		w = w / 2
	}
}

func (t *Tree) Root() []byte {
	return t.nodeHash(0, t.maxHeight-1)
}

func (t *Tree) nodeHash(width, height int) []byte {
	n := t.getNode(width, height)
	if n == nil {
		n = t.getNode(width-1, height)
		if n == nil {
			panic("invalid merkle tree")
		}
	}
	if n.hash != nil {
		return n.hash
	}

	left := t.nodeHash(width*2, height-1)
	right := t.nodeHash(width*2+1, height-1)

	data := make([]byte, len(left)+len(right))
	copy(data[:len(left)], left)
	copy(data[len(left):], right)

	n.hash = t.hasher(data[:])
	return n.hash
}

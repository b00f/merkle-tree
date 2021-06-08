package merkletree

import "math"

type Hasher func([]byte) []byte

type Tree struct {
	nodes     map[int]*node
	maxWidth  int
	maxHeight int

	hasher Hasher
}

type node struct {
	height int
	width  int
	hash   []byte
}

// nodeID return the node ID (four bytes):
// +-+---+
// |h| w |
// +-+---+
//
func nodeID(height, width int) int {
	return ((height & 0xff) << 24) | (width & 0xffffff)
}

func New(hasher Hasher) *Tree {
	return &Tree{
		nodes:  make(map[int]*node),
		hasher: hasher,
	}
}
func (t *Tree) createNode(height, width int) *node {
	return &node{
		height: height,
		width:  width,
	}
}

func (t *Tree) getNode(height, width int) *node {
	id := nodeID(height, width)
	return t.nodes[id]
}

func (t *Tree) getOrCreateNode(height, width int) *node {
	id := nodeID(height, width)
	node, ok := t.nodes[id]
	if !ok {
		node = t.createNode(height, width)
		t.nodes[id] = node
	}
	return node
}

func (t *Tree) invalidateNode(height, width int) {
	n := t.getOrCreateNode(height, width)
	n.hash = nil
}

func (t *Tree) SetBlockData(no int, data []byte) {
	h := t.hasher(data)
	node := t.getOrCreateNode(0, no)
	node.hash = h

	if no+1 > t.maxWidth {
		maxHeight := math.Log2(float64(no + 2))

		t.maxWidth = no + 1
		t.maxHeight = int(math.Trunc(maxHeight))
	}

	w := no / 2
	for h := 1; h < t.maxHeight+1; h++ {
		t.invalidateNode(h, w)
		w = w / 2
	}
}

func (t *Tree) Root() []byte {
	return t.nodeHash(t.maxHeight, 0)
}

func (t *Tree) nodeHash(height, width int) []byte {
	n := t.getNode(height, width)
	if n == nil {
		n = t.getNode(height, width-1)
		if n == nil {
			panic("invalid merkle tree")
		}
	}
	if n.hash != nil {
		return n.hash
	}

	left := t.nodeHash(height-1, width*2)
	right := t.nodeHash(height-1, width*2+1)

	data := make([]byte, len(left)+len(right))
	copy(data[:len(left)], left)
	copy(data[len(left):], right)

	n.hash = t.hasher(data[:])
	return n.hash
}

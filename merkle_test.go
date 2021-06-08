package merkletree

import (
	"encoding/hex"
	"testing"

	"github.com/dchest/blake2b"
	"github.com/magiconair/properties/assert"
)

func HashBlake2b256(data []byte) []byte {
	h := blake2b.Sum256(data)
	return h[:]
}

func TestNodeID(t *testing.T) {
	assert.Equal(t, nodeID(0, 0), 0x00000000)
	assert.Equal(t, nodeID(1, 0), 0x01000000)
	assert.Equal(t, nodeID(0, 1), 0x00000001)
	assert.Equal(t, nodeID(1, 1), 0x01000001)
	assert.Equal(t, nodeID(0xff, 0xffffff), 0xffffffff)
	assert.Equal(t, nodeID(0x00, 0xffffff), 0x00ffffff)
	assert.Equal(t, nodeID(0x77, 0xff00ff), 0x77ff00ff)
}

func Hash256(data []byte) []byte {
	h := blake2b.Sum256(data)
	return h[:]
}

func TestMerkleTree(t *testing.T) {
	tree := New(HashBlake2b256)
	tree.SetBlockData(0, []byte("a"))
	tree.SetBlockData(1, []byte("b"))
	tree.SetBlockData(2, []byte("c"))

	expected, _ := hex.DecodeString("e6061997a9011668bcf216020aaad9cc7f5f34d5b6f78f1e63ef6257c1aa1f37")
	assert.Equal(t, tree.Root(), expected)
}

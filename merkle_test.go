package merkletree

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func HashSha256(data []byte) []byte {
	h := sha256.Sum256(data)
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

func TestCalculateHeight(t *testing.T) {
	tree := New(HashSha256)

	tree.recalculateHeight(0)
	assert.Equal(t, tree.maxHeight, 0)

	tree.recalculateHeight(1)
	assert.Equal(t, tree.maxHeight, 1)

	tree.recalculateHeight(2)
	assert.Equal(t, tree.maxHeight, 2)

	tree.recalculateHeight(4)
	assert.Equal(t, tree.maxHeight, 3)

	tree.recalculateHeight(5)
	assert.Equal(t, tree.maxHeight, 4)

	tree.recalculateHeight(8)
	assert.Equal(t, tree.maxHeight, 4)

	tree.recalculateHeight(9)
	assert.Equal(t, tree.maxHeight, 5)
}

func TestMerkleTree(t *testing.T) {
	tree := New(HashSha256)

	data := []string{
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
		"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	root := []string{
		"559aead08264d5795d3909718cdd05abd49572e84fe55590eef31a88a08fdffd",
		"63956f0ce48edc48a0d528cb0b5d58e4d625afb14d63ca1bb9950eb657d61f40",
		"420940ee1c7a73de80cfa2554efb4e6cec7ea745fed73108ccb06886054df8c6",
		"1b3faa3fcc5ed50cd8592f805c6f8fce976b8582c739b26a6f3613b7f9b13617",
		"cc7461bf32d59f9249796e6e7691f304a3221825732284b9c049b1dd2c689b56",
		"80da0f7dd1db3cd8cfeb4a054bc3f096a97b2392030c45fb76aee7663206bc5d",
		"9fa660bf1b4f58ff50e1823ac5c0350fbfb4281e7d0fdb1fc2c6d42373486d80",
		"b7d3319e67bfc42bc2013a6ca62152fcf5f43a6696adfb70e8bed24ed91e6af6",
		"d2dd5753257c547cf788887711001b6f459f47ccfbaa759928e568b78febe4c9",
		"2585391c79b18af75704d60d2441eb42fde911c12197b7e4f67f628b60eb51ae",
		"a9e6d67397975c4ea9e5c43d2822c6a4e4bbaa18d97bd6c7434b0fcf06992a0d",
		"d7453274f321407edf5f18c36aefcbfb8f177ad92ab65350deb5f8d913c8dec9",
		"344af865adf38b73fb5f7207d64825619571cadf09ae2985392649617a54d98b",
		"3fcb8edf4b06cc5caf1badcb7599830efd43f4f276c87ea51a8324d2da8739c4",
		"7974cfdd80de34c36f665c3726f6c6e569ef5311ac0ffbd1a7e6e67119d54521",
		"6d7256dcb0989bded454695e2c2ddab5cd9d390b4d2af13597b733758def5c73",
		"cf9c549a1390b46bf5808970914149118beb5721af83e6b441b2cac32d424306",
		"c1bffda62df3ecbac683bebde877a20a95b704fe99ba504023e7a593309e9302",
		"4c3b7395e45a2afc86a7978b4665c181ed34a93cd8b334c7af8e32671514b843",
		"8b3fb8638a84f2cc7883b223e4f34bd46119ea7bfde482eb9c93f3c87cf176d8",
		"c2bbcd611b2684b33a9df3d7f799cd83fb84ec6c5e72df18419320cc16685fd9",
		"5fc7f51fc8049472f5e810c8006776a679bb797d39d63ae272314dc430eccacf",
		"038db10eb4d4ecebcc06efa6e9ebe5300e0417a16f499355142ce00e78354fdf",
		"b72b7ac9f0a971ba6f51c6821577f43b8bc9339e006b728c306c6a345a39b6cd",
		"aca0a42b6da93c26f628aef660e575c49f1b36276fa31ec3eb09b943d4794efd",
		"b40b1080634cb680fa9099ab5af25a7a3a3333ef9c3b3ba1143489b4d90e56e8",
	}

	for i, d := range data {
		tree.SetBlockData(i, []byte(d))
		expected, _ := hex.DecodeString(root[i])
		assert.Equal(t, tree.Root(), expected, "Root %d not matched", i)
	}

	// Modifying some data blocks
	tree.SetBlockData(0, []byte("a"))
	tree.SetBlockData(21, []byte("v"))
	expected, _ := hex.DecodeString("5d179ab8afb696b57f433d820406530dae4952a02ba04317c97c86cdff207795")
	assert.Equal(t, tree.Root(), expected)

}

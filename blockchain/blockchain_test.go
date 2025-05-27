package blockchain

import (
	"bytes"
	"math"
	"math/big"
	"reflect"
	"testing"
)

func BenchmarkNewBlock(b *testing.B) {
	data := "test data"
	prevHash := []byte{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewBlock(data, prevHash)
	}
}

func BenchmarkAddBlock(b *testing.B) {
	chain := InItBlockChain()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		chain.AddBlock("test data")
	}
}

func TestNewBlock(t *testing.T) {
	data := "Test data"
	prevHash := []byte("previous hash")
	block := NewBlock(data, prevHash)

	if block == nil {
		t.Fatal("NewBlock returned nil")
	}

	if !bytes.Equal(block.Data, []byte(data)) {
		t.Errorf("Expected data %s, got %s", data, block.Data)
	}

	if !bytes.Equal(block.PrevHash, prevHash) {
		t.Errorf("Expected previous hash %x, got %x", prevHash, block.PrevHash)
	}

	if len(block.Hash) == 0 {
		t.Error("New block hash is empty")
	}

	if block.Nonce == 0 {
		t.Error("New block nonce is 0")
	}
}

func TestSerializeDeserializeBlock(t *testing.T) {
	data := "Test data"
	prevHash := []byte("previous hash")
	block := NewBlock(data, prevHash)

	serializedBlock := block.Serialize()
	deserializedBlock := new(Block).Deserialize(serializedBlock)

	if !reflect.DeepEqual(block, deserializedBlock) {
		t.Errorf("Original and deserialized blocks are not equal. Original: %+v, Deserialized: %+v", block, deserializedBlock)
	}
}

func TestInItBlockChain(t *testing.T) {
	chain := InItBlockChain()

	if chain == nil {
		t.Fatal("InItBlockChain returned nil")
	}

	if len(chain.Blocks) != 1 {
		t.Errorf("Expected 1 block in new chain, got %d", len(chain.Blocks))
	}

	if !bytes.Equal(chain.Blocks[0].Data, []byte("Genesis")) {
		t.Errorf("Expected Genesis block data 'Genesis', got %s", chain.Blocks[0].Data)
	}
}

func TestAddBlock(t *testing.T) {
	chain := InItBlockChain()
	initialBlocksCount := len(chain.Blocks)
	data := "Test block data"

	chain.AddBlock(data)

	if len(chain.Blocks) != initialBlocksCount+1 {
		t.Errorf("Expected %d blocks after adding, got %d", initialBlocksCount+1, len(chain.Blocks))
	}

	newBlock := chain.Blocks[len(chain.Blocks)-1]
	if !bytes.Equal(newBlock.Data, []byte(data)) {
		t.Errorf("Expected new block data %s, got %s", data, newBlock.Data)
	}

	prevBlock := chain.Blocks[len(chain.Blocks)-2]
	if !bytes.Equal(newBlock.PrevHash, prevBlock.Hash) {
		t.Errorf("Expected new block prevHash %x, got %x", prevBlock.Hash, newBlock.PrevHash)
	}
}

func TestNewProofOfWork(t *testing.T) {
	block := &Block{Data: []byte("test"), PrevHash: []byte("prev")}
	pow := NewProofOfWork(block)

	if pow == nil {
		t.Fatal("NewProofOfWork returned nil")
	}

	expectedTarget := big.NewInt(1)
	expectedTarget.Lsh(expectedTarget, uint(256-Difficulty))

	if pow.Target.Cmp(expectedTarget) != 0 {
		t.Errorf("Expected target %s, got %s", expectedTarget.Text(16), pow.Target.Text(16))
	}
}

func TestRun(t *testing.T) {
	block := &Block{Data: []byte("test"), PrevHash: []byte("prev")}
	pow := NewProofOfWork(block)

	nonce, hash := pow.Run()

	if nonce < 0 {
		t.Error("Run returned invalid nonce")
	}

	if len(hash) == 0 {
		t.Error("Run returned empty hash")
	}

	block.Nonce = nonce
	block.Hash = hash

	if !pow.Validate() {
		t.Error("Validated failed after Run")
	}
}

func TestValidate(t *testing.T) {
	// Test case with a valid block
	data := []byte("valid test")
	prevHash := []byte("prev")

	// Create a block and run Proof of Work to find a valid nonce and hash
	block := &Block{Data: data, PrevHash: prevHash}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Nonce = nonce
	block.Hash = hash

	validPoW := NewProofOfWork(block)
	if !validPoW.Validate() {
		t.Error("Validate failed for valid block")
	}

	// Test case with an invalid block (nonce changed)
	invalidBlock := &Block{
		Hash:     hash,
		Data:     data,
		PrevHash: prevHash,
		Nonce:    nonce + 1, // Incorrect nonce
	}
	invalidPoW := NewProofOfWork(invalidBlock)
	if invalidPoW.Validate() {
		t.Error("Validate succeeded for invalid block")
	}
}

func TestToHex(t *testing.T) {
	tests := []struct {
		input    int64
		expected []byte
	}{
		{0, []byte{0, 0, 0, 0, 0, 0, 0, 0}},
		{1, []byte{0, 0, 0, 0, 0, 0, 0, 1}},
		{255, []byte{0, 0, 0, 0, 0, 0, 0, 255}},
		{256, []byte{0, 0, 0, 0, 0, 0, 1, 0}},
		{math.MaxInt64, []byte{127, 255, 255, 255, 255, 255, 255, 255}},
	}

	for _, tt := range tests {
		result := ToHex(tt.input)
		if !bytes.Equal(result, tt.expected) {
			t.Errorf("ToHex(%d): Expected %x, got %x", tt.input, tt.expected, result)
		}
	}
}

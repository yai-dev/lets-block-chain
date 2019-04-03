package lets_block_chain

import "testing"

func TestNewBlock(t *testing.T) {
	block := NewBlock("TestBlock", []byte{})
	t.Logf("\nPrev Hash: `%x`\nData: `%s`\nHash: `%x`\n", block.Prev, block.Data, block.Hash)
}

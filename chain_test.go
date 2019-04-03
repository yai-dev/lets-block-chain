package lbc

import (
	"strconv"
	"testing"
)

func TestNewBlockChain(t *testing.T) {
	chain := NewBlockChain()
	defer func() {
		if err := chain.db.Close(); err != nil {
			t.Fatalf("Close chain failure with error: %s", err)
		}
	}()
	t.Logf("Block Chain: `%+v`", chain)
}

func TestBlockChain_AddBlock(t *testing.T) {
	chain := NewBlockChain()
	defer func() {
		if err := chain.db.Close(); err != nil {
			t.Fatalf("Close chain failure with error: %s", err)
		}
	}()
	chain.AddBlock("World of Warcraft")
	chain.AddBlock("Dota 2")
	t.Logf("Block Chain: `%+v`", chain)

	iterator := chain.iterator()
	for {
		block := iterator.Next()

		pow := newProofWork(block)
		t.Logf("\nPrev block hash: %x\nData: %s\nHash: %x\nPoW: %s\n",
			block.Prev,
			block.Data,
			block.Hash,
			strconv.FormatBool(pow.Validate()))

		if len(block.Prev) == 0 {
			break
		}
	}
}

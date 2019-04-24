package lbc

import "testing"

func TestNewBlock(t *testing.T) {
	transactions := []*Transaction{NewCoinBaseTX("Candy Ye", "Gavin Sun")}
	block := NewBlock(transactions, []byte{})
	t.Logf("\nPrev Hash: `%x`\nTransaction Num: `%d`\nHash: `%x`\n", block.Prev, len(block.Transactions), block.Hash)
}

package lbc

import (
	"strconv"
	"testing"
)

func TestNewBlockChain(t *testing.T) {
	chain := NewBlockChain("Gavin Sun")
	defer chain.Close()
	defer func() {
		if err := chain.db.Close(); err != nil {
			t.Fatalf("Close chain failure with error: %s", err)
		}
	}()
	t.Logf("Block Chain: `%+v`", chain)
}

func TestBlockChain_AddBlock(t *testing.T) {
	chain := NewBlockChain("Gavin Sun")
	defer chain.Close()
	defer func() {
		if err := chain.db.Close(); err != nil {
			t.Fatalf("Close chain failure with error: %s", err)
		}
	}()
	transactions1 := []*Transaction{NewCoinBaseTX("Candy Ye", "Gavin Sun")}
	transactions2 := []*Transaction{
		NewCoinBaseTX("Candy Ye", "Gavin Sun"),
		NewCoinBaseTX("Gavin Sun", "Candy Ye"),
	}
	chain.AddBlock(transactions1)
	chain.AddBlock(transactions2)
	t.Logf("\nBlock Chain: `%+v`", chain)

	iterator := chain.Iterator()
	for {
		block := iterator.Next()

		pow := NewProofWork(block)
		t.Logf("\nPrev block hash: %x\nTransaction Num: %d\nHash: %x\nPoW: %s\n",
			block.Prev,
			len(block.Transactions),
			block.Hash,
			strconv.FormatBool(pow.Validate()))

		if len(block.Prev) == 0 {
			break
		}
	}
}

func TestBlockChain_FindSpendableOutputs(t *testing.T) {
	chain := NewBlockChain("Gavin Sun")
	defer chain.Close()
	accumulated, unspentOutputs := chain.FindSpendableOutputs("Gavin Sun", 1)
	t.Logf("\n`Gavin Sun`'s Accumulated is `%d`, unspent outputs num is `%d`", accumulated, len(unspentOutputs))
}

func TestBlockChain_FindUnspentTransactions(t *testing.T) {
	chain := NewBlockChain("Gavin Sun")
	defer chain.Close()
	transactions := chain.FindUnspentTransactions("Gavin Sun")
	for _, tx := range transactions {
		t.Logf("\nTransaction ID: %x\nTransaction Input Num: %d\nTransaction Output Num: %d", tx.ID, len(tx.VIn), len(tx.VOut))
	}
}

func TestBlockChain_FindUTXO(t *testing.T) {
	chain := NewBlockChain("Gavin Sun")
	defer chain.Close()
	outputs := chain.FindUTXO("Gavin Sun")
	t.Logf("\n`Gavin Sun` has `%d` UTXO", len(outputs))
}

func TestBlockChain_Close(t *testing.T) {
	chain := NewBlockChain("Gavin Sun")
	defer chain.Close()
}

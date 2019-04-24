package lbc

import (
	"fmt"
	"testing"
)

func TestNewProofWork(t *testing.T) {
	transactions := []*Transaction{NewCoinBaseTX("Candy Ye", "Gavin Sun")}
	b := NewBlock(transactions, []byte{})
	pow := NewProofWork(b)
	t.Logf("ProofOfWork: `%+v`", pow)
}

func TestProofOfWork_Run(t *testing.T) {
	transactions := []*Transaction{NewCoinBaseTX("Candy Ye", "Gavin Sun")}
	b := NewBlock(transactions, []byte{})
	pow := NewProofWork(b)
	t.Logf("ProofOfWork: `%+v`", pow)
	count, hash := pow.Run()
	fmt.Printf("Block`%d`'{\npow count:`%d`, \nhash:`%x`\n}\n", len(b.Transactions), count, hash)
}

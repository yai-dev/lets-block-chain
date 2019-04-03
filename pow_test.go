package lets_block_chain

import (
	"fmt"
	"testing"
)

func TestNewProofWork(t *testing.T) {
	b := NewBlock("test pow", []byte{})
	pow := newProofWork(b)
	t.Logf("ProofOfWork: `%+v`", pow)
}

func TestProofOfWork_Run(t *testing.T) {
	b := NewBlock("test pow", []byte{})
	pow := newProofWork(b)
	t.Logf("ProofOfWork: `%+v`", pow)
	count, hash := pow.Run()
	fmt.Printf("Block`%s`'{\npow count:`%d`, \nhash:`%x`\n}\n", b.Data, count, hash)
}

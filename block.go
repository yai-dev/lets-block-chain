package lbc

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

// built-in block version
const (
	blockVersion     = 1
	genesisBlockData = "Genesis Block"
)

// Block represents a block in block-chain
type Block struct {
	// payload stored in this block
	Data []byte
	// previous block 32bit hash
	Prev []byte
	// 32bit hash for this block
	Hash []byte
	// created time for this block
	Timestamp int64
	// block version number
	Version int32
	// block nonce
	Nonce int
}

// newGenesisBlock will generate a genesis block,
// block data will be constant data,
// genesis block's prev is nil
func newGenesisBlock() *Block {
	return NewBlock(genesisBlockData, []byte{})
}

// NewBlock will create a new block with spec data and prev block hash
func NewBlock(data string, prev []byte) *Block {
	_block := &Block{
		Timestamp: time.Now().Unix(),
		Data:      []byte(data),
		Prev:      prev,
		Hash:      []byte{},
		Version:   blockVersion,
		Nonce:     0,
	}

	_pow := newProofWork(_block)
	nonce, hash := _pow.Run()

	_block.Nonce = nonce
	_block.Hash = hash
	return _block
}

// serialize the block struct to bytes array
func (b *Block) serialize() []byte {
	var blobBuf bytes.Buffer
	enc := gob.NewEncoder(&blobBuf)

	if err := enc.Encode(b); err != nil {
		panic(fmt.Sprintf("Serialized block failed, error with:`%s`\n", err))
	}

	return blobBuf.Bytes()
}

// deserializeBlock to block struct with spec bytes array
func deserializeBlock(blob []byte) *Block {
	var b Block
	dec := gob.NewDecoder(bytes.NewReader(blob))

	if err := dec.Decode(&b); err != nil {
		panic(fmt.Sprintf("Deserialized block failed, error with:`%s`", err))
	}

	return &b
}

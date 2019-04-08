package lbc

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

const (
	// targetBits represents difficulty of generate block
	targetBits = 24
	// maxNonce represents Hashcash max count
	maxNonce = math.MaxInt64
)

// ProofOfWork represents had the auth to generate block,
// using `Hashcash` algorithm to limit the efficiency of
// block generation
type ProofOfWork struct {
	// block pointer
	block *Block
	// pow target, init with 1, left shift 256 - targetBits
	target *big.Int
}

// NewProofWork will create pow for spec block
func NewProofWork(b *Block) *ProofOfWork {
	_target := big.NewInt(1)
	_target.Lsh(_target, uint(256-targetBits))
	return &ProofOfWork{
		block:  b,
		target: _target,
	}
}

// prepare hash data
func (pow *ProofOfWork) prepare(nonce int) []byte {
	return bytes.Join(
		[][]byte{
			pow.block.Prev,
			pow.block.Data,
			intToHex(pow.block.Timestamp),
			intToHex(int64(targetBits)),
			intToHex(int64(nonce)),
		}, []byte{})
}

// Run Hashcash algorithm for generating block,
// only hash reach the targetBits will return
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < maxNonce {
		data := pow.prepare(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			fmt.Printf("\r%x", hash)
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

// Validate Proof of Work
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepare(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}

// intToHex converts an int64 to a byte array
func intToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

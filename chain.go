package lbc

import (
	"fmt"
	"github.com/boltdb/bolt"
)

const (
	dbFilePath       = "./lbc.db"
	dbFilePerm       = 0666
	blocksBucketName = "__blocks__"
	lastBlockKey     = "__l__"
)

type (
	// BlockChain represents a chain of ordered blocks
	BlockChain struct {
		tip []byte
		db  *bolt.DB
	}

	// blockChainIterator will read block in db one by one
	blockChainIterator struct {
		currentHash []byte
		db          *bolt.DB
	}
)

// NewBlockChain will return a new block-chain,
// the first block will be Genesis Block
func NewBlockChain(dbPath string) *BlockChain {
	var _tip []byte
	_db, err := bolt.Open(dbPath, dbFilePerm, nil)

	err = _db.Update(func(tx *bolt.Tx) error {
		_bkt := tx.Bucket([]byte(blocksBucketName))

		if _bkt == nil {
			genesis := newGenesisBlock()
			_bkt, err = tx.CreateBucket([]byte(blocksBucketName))

			err = _bkt.Put(genesis.Hash, genesis.serialize())
			if err != nil {
				return fmt.Errorf("create block chain failed, error with:`%s`", err)
			}
			err = _bkt.Put([]byte(lastBlockKey), genesis.Hash)
			if err != nil {
				return fmt.Errorf("create block chain failed, error with:`%s`", err)
			}

			_tip = genesis.Hash
		} else {
			_tip = _bkt.Get([]byte(lastBlockKey))
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	_bc := &BlockChain{
		tip: _tip,
		db:  _db,
	}
	return _bc
}

// AddBlock to this chain
func (bc *BlockChain) AddBlock(data string) {
	var lHash []byte

	// need get last block's hash
	err := bc.db.View(func(tx *bolt.Tx) error {
		_bkt := tx.Bucket([]byte(blocksBucketName))
		lHash = _bkt.Get([]byte(lastBlockKey))
		return nil
	})

	_newBlock := NewBlock(data, lHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		_bkt := tx.Bucket([]byte(blocksBucketName))

		if err := _bkt.Put(_newBlock.Hash, _newBlock.serialize()); err != nil {
			return err
		}

		if err = _bkt.Put([]byte(lastBlockKey), _newBlock.Hash); err != nil {
			return nil
		}

		bc.tip = _newBlock.Hash

		return nil
	})

	if err != nil {
		panic(err)
	}
}

// Iterator used for iteration block chain
func (bc *BlockChain) Iterator() *blockChainIterator {
	return &blockChainIterator{
		currentHash: bc.tip,
		db:          bc.db,
	}
}

// Next block in spec Iterator's block chain
func (i *blockChainIterator) Next() *Block {
	var _b *Block

	_ = i.db.View(func(tx *bolt.Tx) error {
		_bkt := tx.Bucket([]byte(blocksBucketName))
		encodedBlock := _bkt.Get(i.currentHash)
		_b = deserializeBlock(encodedBlock)
		return nil
	})

	i.currentHash = _b.Prev

	return _b
}

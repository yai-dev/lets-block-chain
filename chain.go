package lbc

import (
	"encoding/hex"
	"fmt"

	"github.com/boltdb/bolt"
)

const (
	dbFilePath          = "./db/lbc.db"
	dbFilePerm          = 0666
	blocksBucketName    = "__blocks__"
	lastBlockKey        = "__l__"
	genesisCoinBaseData = `Alibaba's Jack Ma Again Endorses China's '996' Overtime Culture as Testament to Professional Passion`
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
func NewBlockChain(address string) *BlockChain {
	var _tip []byte
	_db, err := bolt.Open(dbFilePath, dbFilePerm, nil)

	err = _db.Update(func(tx *bolt.Tx) error {
		_bkt := tx.Bucket([]byte(blocksBucketName))

		if _bkt == nil {
			genesis := newGenesisBlock(NewCoinBaseTX(address, genesisCoinBaseData))
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
func (bc *BlockChain) AddBlock(transactions []*Transaction) {
	var lHash []byte

	// need get last block's hash
	err := bc.db.View(func(tx *bolt.Tx) error {
		_bkt := tx.Bucket([]byte(blocksBucketName))
		lHash = _bkt.Get([]byte(lastBlockKey))
		return nil
	})

	_newBlock := NewBlock(transactions, lHash)

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

// FindUnspentTransactions will find the transactions include the unspent
// transaction output, unspent transaction output represent that these
// transaction output has not been referenced in any transaction include yet
func (bc *BlockChain) FindUnspentTransactions(addr string) []Transaction {
	var unspentTXs []Transaction
	spentTXOutputs := make(map[string][]int)
	iterator := bc.Iterator()

	for {
		_block := iterator.Next()

		for _, tx := range _block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.VOut {
				if spentTXOutputs[txID] != nil {
					for _, spentOut := range spentTXOutputs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				if out.CanBeUnlockedWith(addr) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			if !tx.isCoinBaseTx() {
				for _, in := range tx.VIn {
					if in.CanUnlockOutputWith(addr) {
						inTxID := hex.EncodeToString(in.TxID)
						spentTXOutputs[inTxID] = append(spentTXOutputs[inTxID], in.VOut)
					}
				}
			}
		}

		if len(_block.Prev) == 0 {
			break
		}
	}

	return unspentTXs
}

// FindSpendableOutputs will find the all spendable outputs in current chain, and ensure that
// they store enough value
func (bc *BlockChain) FindSpendableOutputs(addr string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(addr)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.VOut {
			if out.CanBeUnlockedWith(addr) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

// FindUTXO will find the specified address's unspent transaction outputs
func (bc *BlockChain) FindUTXO(addr string) []TXOutput {
	var UTXOs []TXOutput
	unspentTransactions := bc.FindUnspentTransactions(addr)

	for _, tx := range unspentTransactions {
		for _, out := range tx.VOut {
			if out.CanBeUnlockedWith(addr) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
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

// CLose block chain
func (bc *BlockChain) Close() {
	if err := bc.db.Close(); err != nil {
		panic("Close block chain error!")
	}
}

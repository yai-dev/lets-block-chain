package lbc

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
)

// In the real block-chain, ScriptPubKey and ScriptSig will
// use `Script` programing language to implement, in our case,
// we will use public key and signature to implement it.
type (
	// TXOutput represents a simple transaction output,
	// contain a virtual coins value and script to lock
	// these virtual coins, only unlock script can be use
	// these virtual coins
	TXOutput struct {
		Value        int
		ScriptPubKey string
	}

	// TXInput represents a simple transaction input,
	// contain output id of the last transaction and
	// the index of this output in that transaction
	// (because a transaction may have multiple outputs),
	// ScriptSig is a script that provides data that unlocks the
	// ScriptPubKey field in the TXOutput
	TXInput struct {
		TxID      []byte
		VOut      int
		ScriptSig string
	}

	// Transaction use UTXO model, a transaction consists of a
	// combination of inputs(TXInput) and outputs(TXOutput)
	Transaction struct {
		ID   []byte
		VIn  []TXInput
		VOut []TXOutput
	}
)

const (
	subSidy = 10
)

func (tx *Transaction) setID() {
	var _byteBuf bytes.Buffer
	_enc := gob.NewEncoder(&_byteBuf)

	if err := _enc.Encode(tx); err != nil {
		panic("Cannot set id for transaction!")
	}
	_hash := sha256.Sum256(_byteBuf.Bytes())
	tx.ID = _hash[:]
}

// isCoinBaseTx return True if the current transaction is a CoinBase transaction.
func (tx *Transaction) isCoinBaseTx() bool {
	return len(tx.VIn) == 1 && len(tx.VIn[0].TxID) == 0 && tx.VIn[0].VOut == -1
}

// NewCoinBaseTX will publish new virtual coins
func NewCoinBaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	tx := Transaction{
		ID: nil,
		VIn: []TXInput{
			{
				TxID:      []byte{},
				VOut:      -1,
				ScriptSig: data,
			},
		},
		VOut: []TXOutput{
			{
				Value:        subSidy,
				ScriptPubKey: to,
			},
		},
	}

	tx.setID()
	return &tx
}

// NewUTXOTransaction will create new general transaction output, in this step,
// we need finds all unspent outputs to reference in inputs
func NewUTXOTransaction(from, to string, amount int, bc *BlockChain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	if acc < amount {
		panic("Not enough funds!")
	}

	// Build a list of inputs
	for txId, outs := range validOutputs {
		ID, err := hex.DecodeString(txId)
		if err != nil {
			panic(fmt.Sprintf("Can't create new UTXO transaction with error: %s", err))
		}

		for _, out := range outs {
			input := TXInput{ID, out, from}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	outputs = append(outputs, TXOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.setID()

	return &tx
}

// CanUnlockOutputWith specified data
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

// CanBeUnlockedWith specified data
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

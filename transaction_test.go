package lbc

import "testing"

func TestNewCoinBaseTX(t *testing.T) {
	tx := NewCoinBaseTX("Candy Ye", "Just for handout!")
	t.Logf("\nTransaction ID: %x\nTransaction VIn: %+v\nTransaction VOut: %+v", tx.ID, tx.VIn, tx.VOut)
}

func TestNewUTXOTransaction(t *testing.T) {
	_bc := NewBlockChain("Gavin Sun")
	defer _bc.Close()
	tx := NewUTXOTransaction("Gavin Sun", "Candy Ye", 1, _bc)
	t.Logf("\nUTXO Transaction ID: %x\nUTXO Transaction VIn: %+v\nUTXO Transaction VOut: %+v", tx.ID, tx.VIn, tx.VOut)
}

func TestTXInput_CanUnlockOutputWith(t *testing.T) {
	input := &TXInput{[]byte{}, -1, "Gavin Sun"}
	input.CanUnlockOutputWith("Gavin Sun")
}

func TestTXOutput_CanBeUnlockedWith(t *testing.T) {
	output := &TXOutput{1, "Candy Ye"}
	output.CanBeUnlockedWith("Candy Ye")
}

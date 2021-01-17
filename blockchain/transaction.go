package blockchain

import (
	"encoding/hex"
	"log"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"bytes"
)

type Transaction struct {
	ID []byte
	Inputs []TxInput
	Outputs []TxOutput
}

func CoinBaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coin to %s", to)
	}
	input := TxInput{[]byte{}, -1, to}
	output := TxOutput{100, to}
	tx := Transaction{nil, []TxInput{input}, []TxOutput{output}}
	tx.SetID()
	return &tx
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	Handle(err)
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func NewTransaction(from, to string, amount int, chain *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	acc, validOutputs := chain.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("Error! Not enough fund!")
	}

	outputs = append(outputs, TxOutput{amount, to})

	if acc > amount {
		outputs = append(outputs, TxOutput{acc-amount, from})
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		Handle(err)
		
		for _, out := range outs {
			inputs = append(inputs, TxInput{txID, out, from})
		}
	}

	tx := Transaction{[]byte{}, inputs, outputs}
	tx.SetID()

	return &tx
}

func (tx *Transaction) IsCoinBase() bool {
	return len(tx.Inputs)==1 && len(tx.Inputs[0].ID)==0 && tx.Inputs[0].Out==-1
}

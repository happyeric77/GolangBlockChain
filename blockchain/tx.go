package blockchain

import (
	"bytes"
	"project4/wallet"
)


type TxOutput struct {
	Value int
	PubKeyHash []byte
}

type TxInput struct {
	ID []byte
	Out int
	Sigature []byte
	Pubkey []byte
}

func NewOutput(value int, address string) *TxOutput {
	txo := &TxOutput{value, nil}
	txo.Lock([]byte(address))
	return txo
}

// UsesKey method is to hash the pubkey insde TxInput and see if the result the sample as the PubKeyHash in output
func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.Pubkey)
	return bytes.Compare(pubKeyHash, lockingHash) == 0
}

// Lock method is used when new output is generated to certain address in a new Transaction (reverse the address to its Pubkey Hash and put it into the new output)
func (out *TxOutput) Lock(address []byte) {
	pubKeyHash := wallet.Base58Decode(address)
	pubKeyHash = pubKeyHash[1:len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}
package blockchain

import (
	"project4/wallet"
	"strings"
	"math/big"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/ecdsa"
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
		randData := make([]byte, 24)
		_, err := rand.Read(randData)
		Handle(err)
		data = fmt.Sprintf("%x", randData)
	}

	input := TxInput{[]byte{}, -1, nil, []byte(data)}
	output := NewOutput(20, to)
	tx := Transaction{nil, []TxInput{input}, []TxOutput{*output}}
	tx.ID = tx.Hash()
	return &tx
}

func (tx *Transaction) Serialize() []byte {
	var encoded bytes.Buffer
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	Handle(err)
	return encoded.Bytes()
}

func (tx *Transaction) Hash() []byte {
	var hash [32]byte
	txCopy := *tx
	txCopy.ID = []byte{}
	hash = sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

func NewTransaction(from, to string, amount int, UTXO *UTXOSet) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	wallets, err := wallet.CreateWallets()
	Handle(err)
	w := wallets.GetWallet(from)
	pubKeyHash := wallet.PublicKeyHash(w.PublicKey)

	acc, validOutputs := UTXO.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		log.Panic("Error! Not enough fund!")
	}

	outputs = append(outputs, *NewOutput(amount, to))


	if acc > amount {
		outputs = append(outputs, *NewOutput(acc-amount, from))
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		Handle(err)
		
		for _, out := range outs {
			inputs = append(inputs, TxInput{txID, out, nil, w.PublicKey})
		}
	}

	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	UTXO.Blockchain.SignTransaction(&tx, w.PrivateKey)

	return &tx
}

func (tx *Transaction) IsCoinBase() bool {
	return len(tx.Inputs)==1 && len(tx.Inputs[0].ID)==0 && tx.Inputs[0].Out==-1
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinBase() {
		return
	}

	for _, in := range tx.Inputs {
		if prevTXs[hex.EncodeToString(in.ID)].ID == nil {
			log.Panic("ERROR: previous transaction is not correct")
		}
	}
	
	txCopy := tx.TrimmedCopy() 
	for inId, in := range txCopy.Inputs {
		prevTX := prevTXs[hex.EncodeToString(in.ID)]
		txCopy.Inputs[inId].Sigature = nil
		txCopy.Inputs[inId].Pubkey = prevTX.Outputs[in.Out].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[inId].Pubkey = nil
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		Handle(err)
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Inputs[inId].Sigature = signature
	}
}

func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, in := range tx.Inputs {
		inputs = append(inputs, TxInput{in.ID, in.Out, nil, nil})
	}

	for _, out := range tx.Outputs {
		outputs = append(outputs, TxOutput{out.Value, out.PubKeyHash})
	}
	return Transaction{tx.ID, inputs, outputs}
}

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinBase() {
		return true
	}
	for _, in := range tx.Inputs {
		if prevTXs[hex.EncodeToString(in.ID)].ID == nil {
			return false
		}
	}
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inId, in := range txCopy.Inputs {
		prevTX := prevTXs[hex.EncodeToString(in.ID)]
		txCopy.Inputs[inId].Sigature = nil
		txCopy.Inputs[inId].Pubkey = prevTX.Outputs[in.Out].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[inId].Pubkey = nil

		r := big.Int{}
		s := big.Int{}

		sigLen := len(in.Sigature)
		r.SetBytes(in.Sigature[:(sigLen/2)])
		s.SetBytes(in.Sigature[(sigLen/2): ])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(in.Pubkey)
		x.SetBytes(in.Pubkey[:(keyLen/2)])
		y.SetBytes(in.Pubkey[(keyLen/2):])

		rawPubkey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubkey, txCopy.ID, &r, &s) == false {
			return false
		}
	}
	return true
}

func (tx *Transaction) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("---Transaction: %x", tx.ID))
	for i, input := range tx.Inputs {
		lines = append(lines, fmt.Sprintf("Input#: %d", i))
		lines = append(lines, fmt.Sprintf("TXID: %x", input.ID))
		lines = append(lines, fmt.Sprintf("Out: %d", input.Out))
		lines = append(lines, fmt.Sprintf("Public key: %x", input.Pubkey))
		lines = append(lines, fmt.Sprintf("Signature: %x", input.Sigature))
	}

	for i, output := range tx.Outputs {
		lines = append(lines, fmt.Sprintf("Output#: %d", i))
		lines = append(lines, fmt.Sprintf("Value: %d", output.Value))
		lines = append(lines, fmt.Sprintf("Script: %x", output.PubKeyHash))
	}
	return strings.Join(lines, "\n")
}
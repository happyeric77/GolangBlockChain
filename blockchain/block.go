package blockchain

import (
	"log"
	"encoding/gob"
	"bytes"
)


// Block struct 
type Block struct {
	Hash []byte
	Transactions []*Transaction
	PrevHash []byte
	Nonce int
}

func (b *Block) HashTransactions() []byte {	
	var txDatas [][]byte
	// var txHash [32]byte
	for _, tx := range b.Transactions{
		txDatas = append(txDatas, tx.Serialize() )
	}
	mkTree := NewMerkleTree(txDatas)
	
	// txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return mkTree.RootNode.Data[:]
}

// CreateBlock function allows to create a new block
func CreateBlock(txs []*Transaction, prevHash []byte) *Block {
	block := &Block{Hash: []byte{}, Transactions: txs, PrevHash: prevHash, Nonce: 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return block
}


// Genesis function allows to create Genesis block without previous hash
func Genesis(coinBase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinBase}, []byte{})
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)
	Handle(err)
	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	Handle(err)
	return &block
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

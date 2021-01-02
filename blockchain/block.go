package blockchain

// import (
// 	"bytes"
// 	"crypto/sha256"
// )

// BlockChain struct
// type BlockChain struct {
// 	Blocks []*Block
// 	// The effient method for demo. The real blockchain type, the more complicated struct will be implemented later.
// }

// Block struct 
type Block struct {
	Hash []byte
	Data []byte
	PrevHash []byte
	Nonce int
}

// // DeriveHash Block method allows current block to derive hash based on data and prevHash
// func (b *Block) DeriveHash(){
// 	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
// 	hash := sha256.Sum256(info) 
// 	// Sha256 is farely simple compare to the real way to calc hash. (demo only)
// 	// More secured hashing function will be implement later.
// 	b.Hash = hash[:]
// }

// CreateBlock function allows to create a new block
func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{Hash: []byte{}, Data: []byte(data), PrevHash: prevHash, Nonce: 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return block
}

// AddBlock method takes data string without output value to add an block into current blockchain.
// func (chain *BlockChain) AddBlock(data string) {
// 	prevBlock := chain.Blocks[len(chain.Blocks)-1]
// 	new := CreateBlock(data, prevBlock.Hash)
// 	chain.Blocks = append(chain.Blocks, new)
// }

// Genesis function allows to create Genesis block without previous hash
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

// InitBlockChain allows to inital a new blockchain
// func InitBlockChain() *BlockChain {
// 	return &BlockChain{[]*Block{Genesis()}}
// }
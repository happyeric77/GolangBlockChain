package blockchain

// BlockChain struct
type BlockChain struct {
	Blocks []*Block
	// The effient method for demo. The real blockchain type, the more complicated struct will be implemented later.
}

// AddBlock method takes data string without output value to add an block into current blockchain.
func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	new := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, new)
}

// InitBlockChain allows to inital a new blockchain
func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}
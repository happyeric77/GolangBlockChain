package main
import (
	"project4/blockchain"
	"fmt"
)

func main(){
	chain := blockchain.InitBlockChain()
	chain.AddBlock("First block after genesis")
	chain.AddBlock("Second block after genesis")
	chain.AddBlock("Third block after genesis")

	for i, block := range chain.Blocks {
		fmt.Println("Block No.", i)
		fmt.Printf("Previous Hash is '%x'\n", block.PrevHash)
		fmt.Printf("Data in block is '%s'\n", block.Data)
		fmt.Printf("Hash is '%x'\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Println(pow.Validate())

		fmt.Println("========================")
	}
}
package main
import (
	"flag"
	"os"
	"project4/blockchain"
	"fmt"
	"runtime"
)

type CommandLine struct {
	blockchain *blockchain.BlockChain
}

func (cli *CommandLine) PrintUsage(){
	fmt.Println("Usage: ")
	fmt.Println("add -block BLOCK_DATA - Add a block to the chain")
	fmt.Println("print - Print the blocks in the chain")
}

func (cli *CommandLine) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) AddBlock(data string){
	cli.blockchain.AddBlock(data)
	fmt.Println("Block added")
}

func (cli *CommandLine) PrintChain() {
	
	iter := cli.blockchain.Iterator()
	count := 1
	
	for {
		block := iter.Next()
		fmt.Println("Block No.", count)
		fmt.Printf("Previous Hash is '%x'\n", block.PrevHash)
		fmt.Printf("Data in block is '%s'\n", block.Data)
		fmt.Printf("Hash is '%x'\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Println(pow.Validate())

		fmt.Println("========================")
		if len(block.PrevHash) == 0{
			break
		}
		count ++
	}
}

func (cli *CommandLine) Run(){
	cli.ValidateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	default:
		cli.PrintUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed(){
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.AddBlock(*addBlockData)
	}

	if printChainCmd.Parsed(){
		cli.PrintChain()
	}
}

func main(){
	defer os.Exit(0)
	chain := blockchain.InitBlockChain()
	defer chain.Database.Close()
	cli := CommandLine{chain}
	cli.Run()
	// chain.AddBlock("First block after genesis")
	// chain.AddBlock("Second block after genesis")
	// chain.AddBlock("Third block after genesis")

	// for i, block := range chain.Blocks {
	// 	fmt.Println("Block No.", i)
	// 	fmt.Printf("Previous Hash is '%x'\n", block.PrevHash)
	// 	fmt.Printf("Data in block is '%s'\n", block.Data)
	// 	fmt.Printf("Hash is '%x'\n", block.Hash)

	// 	pow := blockchain.NewProof(block)
	// 	fmt.Println(pow.Validate())

	// 	fmt.Println("========================")
	// }

}
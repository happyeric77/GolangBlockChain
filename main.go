package main
import (
	"flag"
	"os"
	"project4/blockchain"
	"fmt"
	"runtime"
)

type CommandLine struct {}

func (cli *CommandLine) PrintUsage(){
	fmt.Println("Usage: ")
	fmt.Println("getbalance -address ADDRESS: Get the balance for the address")
	fmt.Println("createblockchain -address ADDRESS: Creates a blockchain and sends genesis reward to address")
	fmt.Println("printchain: Print the blocks in the chain")
	fmt.Println("send -from FROM -to TO -amount AMOUNT: send amount of coins")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		runtime.Goexit()
	}
}

// func (cli *CommandLine) AddBlock(data string){
// 	cli.blockchain.AddBlock(data)
// 	fmt.Println("Block added")
// }

func (cli *CommandLine) printChain() {
	chain := blockchain.ContinueBlockChain("")
	defer chain.Database.Close()
	iter := chain.Iterator()
	count := 1
	
	for {
		block := iter.Next()
		fmt.Println("Block No.", count)
		fmt.Printf("Previous Hash is '%x'\n", block.PrevHash)
		fmt.Printf("Transactions in block is '%v'\n", block.Transactions)
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

func (cli *CommandLine) createBlockChain(address string) {
	chain := blockchain.InitBlockChain(address)
	defer chain.Database.Close()
	fmt.Println("New block chain created")
}

func (cli *CommandLine) getBalance(address string) {
	chain := blockchain.ContinueBlockChain("")
	UTXOs := chain.FindUTXOs(address)
	balance := 0
	for _, output := range UTXOs {
		balance += output.Value
	}
	fmt.Printf("%s's balance: $%v\n", address, balance)
}

func (cli *CommandLine) send (from, to string, amount int) {
	chain := blockchain.ContinueBlockChain("")
	txn := blockchain.NewTransaction("Eric", "Yuko", amount, chain)

	chain.AddBlock([]*blockchain.Transaction{txn,})
	fmt.Println("Success!")
}

func (cli *CommandLine) run(){
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "Get the balance for the address" )
	createBlockChainAddress := createBlockChainCmd.String("address", "", "Creates a blockchain and sends genesis reward to address")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")
	
	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "createblockchain":
		err := createBlockChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "send":
		err := sendCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	default:
		cli.PrintUsage()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed(){
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockChainCmd.Parsed(){
		if *createBlockChainAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockChain(*createBlockChainAddress)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount == 0 {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

	if printChainCmd.Parsed(){
		cli.printChain()
	}
}

func main(){
	defer os.Exit(0)
	// chain := blockchain.InitBlockChain()
	// defer chain.Database.Close()
	cli := CommandLine{}
	cli.run()
	
}
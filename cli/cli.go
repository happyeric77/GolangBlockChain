package cli

import (
	"log"
	
	"flag"
	"os"
	"project4/blockchain"
	"project4/wallet"
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
	fmt.Println("createwallet: Create a new wallet")
	fmt.Println("listaddresses: List all addresses in the wallet file")
	fmt.Println("reindexutxo: Rebuilds the UTXOSet")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) reindexUTXO() {
	chain := blockchain.ContinueBlockChain("")
	defer chain.Database.Close()
	UTXOSet := blockchain.UTXOSet{chain}
	UTXOSet.Reindex()
	count := UTXOSet.CountTransactions()
	fmt.Printf("Done! There are %d transactions in UTXOSet\n", count)
}

func (cli *CommandLine) listAddresses() {
	wallets, _ := wallet.CreateWallets()
	addresses := wallets.GetAllWalletAddress()
	for _, address := range addresses {
		fmt.Println(address)
	}
}

func (cli *CommandLine) createWallet(){
	wallets, _ := wallet.CreateWallets()
	address := wallets.AddWallet()
	wallets.SaveFile()
	fmt.Printf("New wallet's address: %s\n", address)
}

func (cli *CommandLine) printChain() {
	chain := blockchain.ContinueBlockChain("")
	defer chain.Database.Close()
	iter := chain.Iterator()
	count := 1
	
	for {
		block := iter.Next()
		fmt.Println("Block No.", count)
		fmt.Printf("Hash is '%x'\n", block.Hash)
		fmt.Printf("Previous Hash is '%x'\n", block.PrevHash)
		fmt.Printf("Transactions in block is '%v'\n", block.Transactions)

		for _, tx := range block.Transactions {
			println(tx)
		}		

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
	if !wallet.ValidateAddress(address) {
		log.Panic("Address is not valid")
	}
	chain := blockchain.InitBlockChain(address)
	defer chain.Database.Close()

	UTXOSet := blockchain.UTXOSet{chain}
	UTXOSet.Reindex()

	fmt.Println("New block chain created")
}

func (cli *CommandLine) getBalance(address string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("Address is not valid")
	}
	chain := blockchain.ContinueBlockChain(address)
	defer chain.Database.Close()
	UTXOSet := blockchain.UTXOSet{chain}

	balance := 0	
	pubKeyHash := wallet.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1:len(pubKeyHash)-4]
	UTXOs := UTXOSet.FindUTXO(pubKeyHash)
	for _, output := range UTXOs {
		balance += output.Value
	}
	fmt.Printf("%s's balance: $%v\n", address, balance)
}

func (cli *CommandLine) send (from, to string, amount int) {
	if !wallet.ValidateAddress(from) {
		log.Panic("Address is not valid")
	}
	if !wallet.ValidateAddress(to) {
		log.Panic("Address is not valid")
	}
	chain := blockchain.ContinueBlockChain("")
	defer chain.Database.Close()
	UTXOSet := blockchain.UTXOSet{chain}
	txn := blockchain.NewTransaction(from, to, amount, &UTXOSet)
	coinBaseTx := blockchain.CoinBaseTx(from, "")

	block := chain.AddBlock([]*blockchain.Transaction{coinBaseTx, txn})
	UTXOSet.Update(block)
	fmt.Println("Success!")
}

func (cli *CommandLine) Run(){
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
	reindexUTXOCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "Get the balance for the address" )
	createBlockChainAddress := createBlockChainCmd.String("address", "", "Creates a blockchain and sends genesis reward to address")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")
	
	switch os.Args[1] {
	case "reindexutxo":
		err := reindexUTXOCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

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

	if reindexUTXOCmd.Parsed(){
		cli.reindexUTXO()
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

	if listAddressesCmd.Parsed(){
		cli.listAddresses()
	}

	if createWalletCmd.Parsed(){
		cli.createWallet()
	}
}
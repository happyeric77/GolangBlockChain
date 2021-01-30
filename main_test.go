package main

import (
	"testing"
	"project4/blockchain"
	"project4/wallet"
)


func TestWallets(t *testing.T){
	t.Log("========== Test Wallet Section ============")

	/*=============
		Step 0: 
		Create a wallets file
	=============*/
	wallets, err := wallet.CreateWallets()
	if err != nil {
		t.Error(err)
	}

	/*=============
		Step 1: 
		Get all wallet addresses in the wallets file
	=============*/
	addresses := wallets.GetAllWalletAddress()
	for i, add := range addresses {
		t.Logf("wallet #%d: %s", i, add)
	}

	/*=============
		Step 2: 
		If needed, create a new wallet and save it into wallets file
	=============*/
	// wallet := wallet.MakeWallet()
	// address := string(wallet.Address())
	// t.Log(address)	
	// wallets.Wallets[address] = wallet
	// wallets.SaveFile()

	t.Log("========== Test Wallet Section Finish Line ============")
}




func TestBlockChain(t *testing.T){
	t.Log("========== Test Block Chain Section ============")

	/*=============
		Step 0: 
		get the wallets   
	=============*/
	wallets, _ := wallet.CreateWallets()
	w := wallets.GetWallet("1Cz47VHmAhyNCuQuJKFQ4MXJETLUk1j4JU")
	w2 := wallets.GetWallet("1DZjBsAcuWwDpfCLcAm15jgsSh9Amm33gQ")


	/*============= 
		Step 1.0: 
		If no blockchain, init one  
	=============*/
	// address := "1Cz47VHmAhyNCuQuJKFQ4MXJETLUk1j4JU"
	// chain := blockchain.InitBlockChain(address)

	/*=============
		Step 1.1:
		If blockchain exists, retrieve it 
	=============*/
	chain := blockchain.ContinueBlockChain("")

	/*=============
		Step 1.2:
		Create a persistence layer UTXOSet
	=============*/
	UTXO := blockchain.UTXOSet{chain}

	/* =============
		Step 1.2: 
		If no existing transaction in blockchain, create a coin base genesis transaction
	=============*/
	// coinBaseTx := blockchain.CoinBaseTx("1Cz47VHmAhyNCuQuJKFQ4MXJETLUk1j4JU", "Mined by Eric")
	// chain.AddBlock([]*blockchain.Transaction{coinBaseTx,})

	/*============= 
		Step 2: 
		Iterate through blockchain and print out the elements 
	=============*/
	iter := chain.Iterator()
	count := 1

	for {
		block := iter.Next()
		t.Logf("---------block%d start line-----------", count)
		t.Logf("Previous Hash is '%x'\n", block.PrevHash)
		t.Logf("Transactions in block is '%v'\n", block.Transactions)
		t.Logf("Hash is '%x'\n", block.Hash)

		pow := blockchain.NewProof(block)
		t.Log(pow.Validate())

		t.Logf("---------block%d finish line-----------", count)
		if len(block.PrevHash) == 0{
			break
		}
		count ++
	}	


	/*============= 
		Step 3.1:
		Reindex the UTXOSet persistence layer
	=============*/
	UTXO.Reindex()

	/*============= 
		Step 3.2:
		Get the balance of wallet #1 
	=============*/
	pubKeyHash := wallet.PublicKeyHash(w.PublicKey)
	UTXOs := UTXO.FindUTXO(pubKeyHash)

	balance := 0
	for _, output := range UTXOs {
		balance += output.Value
	}
	t.Logf("Eric's wallet (%s) balance TXN: $%d", w.Address(), balance)

	/*============= 
		Step 4:
		Get the balance of wallet #2 
	=============*/
	pubKeyHash2 := wallet.PublicKeyHash(w2.PublicKey)
	UTXOsYuko := UTXO.FindUTXO(pubKeyHash2)
	balanceYuko := 0
	for _, output := range UTXOsYuko {
		balanceYuko += output.Value
	}
	t.Logf("Yuko's wallet (%s) balance TXN: $%d",w2.Address(), balanceYuko)

	/*============= 
		Step 5:
		If needed, make a transaction from wallet#1 to wallet#2 
	=============*/
	// txn := blockchain.NewTransaction(string(w.Address()), string(w2.Address()), 30, chain)
	// chain.AddBlock([]*blockchain.Transaction{txn,})
	t.Log("========== Test Block Chain Section Finsh Line ============")


}


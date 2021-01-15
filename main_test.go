package main

import (
	"testing"
	"project4/blockchain"
)

func TestBlcokChain(t *testing.T){
	chain := blockchain.ContinueBlockChain("")

	// coinBaseTx := blockchain.CoinBaseTx("Eric", "Mined by Eric")

	// chain.AddBlock([]*blockchain.Transaction{coinBaseTx,})

	iter := chain.Iterator()
	count := 1

	for {
		block := iter.Next()
		t.Log("Block No.", count)
		t.Logf("Previous Hash is '%x'\n", block.PrevHash)
		t.Logf("Transactions in block is '%v'\n", block.Transactions)
		t.Logf("Hash is '%x'\n", block.Hash)

		pow := blockchain.NewProof(block)
		t.Log(pow.Validate())

		t.Log("========================")
		if len(block.PrevHash) == 0{
			break
		}
		count ++
	}	

	UTXOs := chain.FindUTXOs("Eric")
	balance := 0
	for _, output := range UTXOs {
		balance += output.Value
	}
	t.Log("Eric's balance Before TXN: $", balance)

	UTXOsYuko := chain.FindUTXOs("Yuko")
	balanceYuko := 0
	for _, output := range UTXOsYuko {
		balanceYuko += output.Value
	}
	t.Log("Yuko's balance Before TXN: $", balanceYuko)
	

	// txn := blockchain.NewTransaction("Eric", "Yuko", 30, chain)

	// chain.AddBlock([]*blockchain.Transaction{txn,})

	// UTXOsEric := chain.FindUTXOs("Eric")
	// balanceEric := 0
	// for _, output := range UTXOsEric {
	// 	balanceEric += output.Value
	// }
	// t.Log("Eric's balance After TXN: $", balanceEric)

	// UTXOsYuko1 := chain.FindUTXOs("Yuko")
	// balanceYuko1 := 0
	// for _, output := range UTXOsYuko1 {
	// 	balanceYuko1 += output.Value
	// }
	// t.Log("Yuko's balance Before TXN: $", balanceYuko1)
}
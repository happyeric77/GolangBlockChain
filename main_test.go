package main

import (
	"testing"
	"project4/blockchain"
)

func TestBlcokChain(t *testing.T){
	chain := blockchain.ContinueBlockChain("Eric")
	// t.Log(chain.LastHash)
	// iter := chain.Iterator()
	// for {
	// 	block := iter.Next()

	// 	t.Log(block.Nonce)
	// 	if len(block.PrevHash) == 0 {
	// 		break
	// 	}
	// }

	coinBaseTx := blockchain.CoinBaseTx("Eric", "Mined by Eric")

	chain.AddBlock([]*blockchain.Transaction{coinBaseTx,})

	txn := blockchain.NewTransaction("Eric", "Yuko", 30, chain)

	t.Log(txn.Outputs)

}
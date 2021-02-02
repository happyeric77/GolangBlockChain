package blockchain

import (
	"bytes"
	"encoding/hex"
	"log"
	"github.com/dgraph-io/badger/v2"
)

var (
	utxoprefix = []byte("utxo-")
	prefixLength = len(utxoprefix)
)

type UTXOSet struct {
	Blockchain *BlockChain
}

func (u UTXOSet) FindSpendableOutputs (pubKeyHash []byte, amount int) (int, map[string][]int){
	unspentOuts := make(map[string][]int)
	accumulated := 0
	db := u.Blockchain.Database

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek(utxoprefix); it.ValidForPrefix(utxoprefix); it.Next(){
			k := it.Item().Key()
			v, err := it.Item().ValueCopy(nil)
			Handle(err)
			k = bytes.TrimPrefix(k, utxoprefix)
			txID := hex.EncodeToString(k)
			outs := DeserializeOutputs(v)
			for outIdx, out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
					accumulated += out.Value
					unspentOuts[txID] = append(unspentOuts[txID], outIdx)
				}
			}
		}
		return nil
	})
	Handle(err)
	return accumulated, unspentOuts
}


func (u UTXOSet) FindUTXO(pubKeyHash []byte) []TxOutput {
	var UTXOs []TxOutput
	db := u.Blockchain.Database
	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek(utxoprefix); it.ValidForPrefix(utxoprefix); it.Next(){
			v, err := it.Item().ValueCopy(nil)
			Handle(err)
			outs := DeserializeOutputs(v)
			for _, out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
		}
		return nil
	})
	Handle(err)
	return UTXOs
}

func (u UTXOSet) CountTransactions() int {
	db := u.Blockchain.Database
	count := 0
	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek(utxoprefix); it.ValidForPrefix(utxoprefix); it.Next() {
			count ++
		}
		return nil
	})
	Handle(err)
	return count
}


func (u UTXOSet) Reindex() {
	db := u.Blockchain.Database
	u.DeleteByPrefix(utxoprefix)
	UTXO := u.Blockchain.FindUTXO()
	
	err := db.Update(func(txn *badger.Txn) error {
		for txId, outs := range UTXO {
			key, err := hex.DecodeString(txId)
			if err != nil {
				return err
			}
			key = append(utxoprefix, key...)
			err = txn.Set(key, outs.Serialize())
			Handle(err)
		}
		return nil
	})
	Handle(err)
}


func (u *UTXOSet) Update(block *Block) {
	db := u.Blockchain.Database
	err := db.Update(func(txn *badger.Txn)error {
		for _, tx := range block.Transactions {
			if tx.IsCoinBase() == false {
				for _, in := range tx.Inputs {
					updatedOutputs := TxOutputs{}
					inID := append(utxoprefix, in.ID...)
					item, err := txn.Get(inID)
					Handle(err)
					err = item.Value(func(val []byte) error {
						outs := DeserializeOutputs(val)
						for outIdx, out := range outs.Outputs {
							if outIdx != in.Out {
								updatedOutputs.Outputs = append(updatedOutputs.Outputs, out)
							}
						}
						if len(updatedOutputs.Outputs) == 0 {
							if err := txn.Delete(inID); err != nil {
								log.Panic(err)
							}
						} else {
							if err := txn.Set(inID, updatedOutputs.Serialize()); err != nil {
								log.Panic(err)
							}
						}
						return nil
					})
				}
			}
			newOutputs := TxOutputs{}
			for _, out := range tx.Outputs {
				newOutputs.Outputs = append(newOutputs.Outputs, out)
			}

			txID := append(utxoprefix, tx.ID...)
			if err := txn.Set(txID, newOutputs.Serialize()); err != nil {
				log.Panic(err)
			}
			return nil
			
		}
		return nil
	})
	Handle(err)
}


func (u *UTXOSet) DeleteByPrefix(prefix []byte) {
	deleteKeys := func (keysForDelete [][]byte) error {
		if err := u.Blockchain.Database.Update(func(txn *badger.Txn) error {
			for _, key := range keysForDelete {
				if err := txn.Delete(key); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
		return nil
	}

	collectSize := 100000
	err := u.Blockchain.Database.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		keysToDelete := make([][]byte, 0, collectSize)
		keysCollected := 0
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			key := it.Item().KeyCopy(nil)
			keysToDelete = append(keysToDelete, key)
			keysCollected ++
			if keysCollected == collectSize {
				if err := deleteKeys(keysToDelete); err != nil {
					log.Panic(err)
				}
				keysCollected = 0
				keysToDelete = make([][]byte, 0, collectSize)
			}
		}

		if keysCollected > 0 {
			if err := deleteKeys(keysToDelete); err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
	Handle(err)
}

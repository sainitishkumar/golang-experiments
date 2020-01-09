package main

import (
	"encoding/hex"

	"github.com/boltdb/bolt"
)

var chainBucket = []byte("chain")

// UTXOSet for the unspentTX outputs
type UTXOSet struct {
	BlockChain *BlockChain
}

// Reindex the chain state from the blockchain.db file
func (utxoset UTXOSet) Reindex() {
	db := utxoset.BlockChain.db
	_ = db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket(chainBucket)
		tx.CreateBucket(chainBucket)
		return nil
	})
	unspentTXO := utxoset.BlockChain.FindUTXO()
	_ = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(chainBucket)
		for id, outs := range unspentTXO {
			key, _ := hex.DecodeString(id)
			bucket.Put(key, outs.Serialise())
		}
		return nil
	})
}

// FindSpendableOutputs finds spendable outputs for an address
func (utxoset UTXOSet) FindSpendableOutputs(pubkeyhash []byte, amount int) (int, map[string][]int) {
	unspentOut := make(map[string][]int)
	accumulated := 0
	db := utxoset.BlockChain.db

	_ = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(chainBucket)
		cursor := bucket.Cursor()
		for id, outputs := cursor.First(); id != nil; id, outputs = cursor.Next() {
			txid := hex.EncodeToString(id)
			outs := DeSerializeOutputs(outputs)
			for outid, out := range outs.Outputs {
				if out.IsLockedWith(pubkeyhash) && accumulated < amount {
					accumulated += out.Value
					unspentOut[txid] = append(unspentOut[txid], outid)
				}
			}
		}
		return nil
	})

	return accumulated, unspentOut
}

// FindUTXO from chain bucket
func (utxoset UTXOSet) FindUTXO(pubkeyhash []byte) []TXOutput {
	var utxo []TXOutput
	db := utxoset.BlockChain.db
	_ = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(chainBucket)
		cursor := bucket.Cursor()
		for id, out := cursor.First(); id != nil; id, out = cursor.Next() {
			outputs := DeSerializeOutputs(out)
			for _, outs := range outputs.Outputs {
				if outs.IsLockedWith(pubkeyhash) {
					utxo = append(utxo, outs)
				}
			}
		}
		return nil
	})
	return utxo
}

// Update the chain bucket with new block's transactions
func (utxoset UTXOSet) Update(block *Block) {
	db := utxoset.BlockChain.db
	_ = db.Update(func(bolttx *bolt.Tx) error {
		bucket := bolttx.Bucket(chainBucket)
		for _, tx := range block.Transactions {
			if !tx.IsCoinBaseTX() {
				for _, vin := range tx.Vin {
					updatedOuts := TXOutputs{}
					outBytes := bucket.Get(tx.TXid)
					outs := DeSerializeOutputs(outBytes)
					for outid, out := range outs.Outputs {
						if outid != vin.Vout {
							updatedOuts.Outputs = append(updatedOuts.Outputs, out)
						}
					}
					// if no updated outs remove them from the bucket
					if len(updatedOuts.Outputs) == 0 {
						bucket.Delete(vin.TXid)
					} else {
						bucket.Put(vin.TXid, updatedOuts.Serialise())
					}
				}
			}
			// Put the newly generated outs into the bucket
			newOuts := TXOutputs{}
			for _, out := range tx.Vout {
				newOuts.Outputs = append(newOuts.Outputs, out)
			}
			bucket.Put(tx.TXid, newOuts.Serialise())
		}
		return nil
	})
}

// CountTransactions no of tx in UTXOSet
func (utxoset UTXOSet) CountTransactions() int {
	count := 0
	db := utxoset.BlockChain.db
	_ = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(chainBucket)
		cursor := bucket.Cursor()
		for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
			count++
		}
		return nil
	})
	return count
}

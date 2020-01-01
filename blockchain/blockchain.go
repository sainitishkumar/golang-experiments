package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const dbfile = "blockchain.db"

var blocksbucket = []byte("blocks")
var chainbucket = []byte("chain")

const genesisCoinbaseData = "Coinbase tx data for genesis block"

// BlockChain is structure containing slice of blocks
// type BlockChain struct {
// 	Blocks []*Block
// }

// BlockChain contains tip and database pointer
type BlockChain struct {
	tip []byte // last block's hash
	db  *bolt.DB
}

// BlockChainIterator for printing out the BC
type BlockChainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

/* CreateBlockChain checks if blockchain is present in the db or not
if present reads it into an object
otherwise creates a new genesis block and add to BC
DB structure:
Structure of bitcoin core: 1 bucket for blocks, 1 bucket for chain state
Bitcoin:
In blocks, the key -> value pairs are:

'b' + 32-byte block hash -> block index record
'f' + 4-byte file number -> file information record
'l' -> 4-byte file number: the last block file number used
'R' -> 1-byte boolean: whether we're in the process of reindexing
'F' + 1-byte flag name length + flag name string -> 1 byte boolean: various flags that can be on or off
't' + 32-byte transaction hash -> transaction index record
In chainstate, the key -> value pairs are:

'c' + 32-byte transaction hash -> unspent transaction output record for that transaction
'B' -> 32-byte block hash: the block hash up to which the database represents the unspent transaction outputs

Sample:
32-byte block-hash -> Serialized Block struct
'l'                -> The hash of the last block in a chain
*/

func dbExists() bool {
	if _, err := os.Stat(dbfile); os.IsNotExist(err) {
		return false
	}

	return true
}

// CreateBlockChain method
func CreateBlockChain(pubkeyhash string) *BlockChain {
	if dbExists() {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}
	var tip []byte
	db, err := bolt.Open(dbfile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		cbtx := NewCoionbaseTX(pubkeyhash, genesisCoinbaseData)
		genesis := CreateGenesisBlock(cbtx)
		b, err := tx.CreateBucket([]byte(blocksbucket))
		if err != nil {
			log.Panic(err)
		}
		err = b.Put(genesis.BlockHash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte("l"), genesis.BlockHash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.BlockHash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	bc := BlockChain{tip, db}
	return &bc
}

// GetBlockChain Return the created blockchain
func GetBlockChain(address string) *BlockChain {
	if dbExists() == false {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbfile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksbucket))
		tip = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := BlockChain{tip, db}

	return &bc
}

// MineBlock func adds a new block to the main BlockChain struct
func (bc *BlockChain) MineBlock(transactions []*Transaction) {
	// prevBlockHash := bc.Blocks[len(bc.Blocks)-1].BlockHash
	// b := CreateBlock(blockData, prevBlockHash)
	// bc.Blocks = append(bc.Blocks, b)

	// verify TX
	for _, tx := range transactions {
		if bc.VerifyTX(tx) != true {
			fmt.Println("Not correct tx")
			os.Exit(2)
		}
	}

	var tip []byte
	_ = bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(blocksbucket)
		tip = bucket.Get([]byte("l"))
		return nil
	})

	prevBlockHash := tip
	b := CreateBlock(transactions, prevBlockHash)

	// add to DB
	_ = bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(blocksbucket)
		bucket.Put(b.BlockHash, b.Serialize())
		bucket.Put([]byte("l"), b.BlockHash)
		bc.tip = b.BlockHash
		return nil
	})

}

/***   Printing the Blockchain    ***/

// GetIterator returns iterator at the tip of BC
func (bc *BlockChain) GetIterator() *BlockChainIterator {
	temp := BlockChainIterator{bc.tip, bc.db}
	return &temp
	// return &BlockChainIterator{nil, nil}
}

// GetBlock returns the block the iterator is pointing to
func (bcIter *BlockChainIterator) GetBlock() *Block {
	var b *Block
	_ = bcIter.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(blocksbucket)
		blockBytes := bucket.Get(bcIter.currentHash)
		b = DeSerialize(blockBytes)
		return nil
	})
	// point iter to prev block
	bcIter.currentHash = b.PrevBlockHash
	return b
}

// PrintBlockChain prints the BC in the order of blocks
// func (bc *BlockChain) PrintBlockChain() {

// }

// FindUnspentTX finds all unspent transactions, adds them up and returns them
func (bc *BlockChain) FindUnspentTX(pubKeyHash []byte) []Transaction {
	var unspentTX []Transaction
	spentTXO := make(map[string][]int) // vals is an array of ints
	bcIter := bc.GetIterator()
	// fmt.Println(bc)
	// return unspentTX
	// bcIter := BlockChainIterator{bc.tip, bc.db}
	for {
		b := bcIter.GetBlock()
		for _, tx := range b.Transactions {
			txid := hex.EncodeToString(tx.TXid)
		Label:
			for outidx, out := range tx.Vout {
				if spentTXO[txid] != nil {
					// output was spent
					for _, spentOut := range spentTXO[txid] {
						if spentOut == outidx {
							continue Label
						}
					}
				}
				if out.IsLockedWith(pubKeyHash) {
					unspentTX = append(unspentTX, *tx)
				}
			}
			if !tx.IsCoinBaseTX() {
				for _, in := range tx.Vin {
					if in.UsesKey(pubKeyHash) {
						intxid := hex.EncodeToString(in.TXid)
						spentTXO[intxid] = append(spentTXO[intxid], in.Vout)
					}
				}
			}
		}
		// if prev block is genesis block stop
		if len(b.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTX
}

// FindUTXO finds all unspent transaction outputs
func (bc *BlockChain) FindUTXO(pubkeyhash []byte) []TXOutput {
	var UTXOs []TXOutput
	unspentTX := bc.FindUnspentTX(pubkeyhash)

	for _, tx := range unspentTX {
		for _, out := range tx.Vout {
			if out.IsLockedWith(pubkeyhash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs
}

// FindSpendableOutputs finds spendable outputs for an address
func (bc *BlockChain) FindSpendableOutputs(pubkeyhash []byte, amount int) (int, map[string][]int) {
	unspentOut := make(map[string][]int)
	unspentTx := bc.FindUnspentTX(pubkeyhash)
	accumulated := 0
Label:
	for _, tx := range unspentTx {
		txid := hex.EncodeToString(tx.TXid)
		for outid, out := range tx.Vout {
			if out.IsLockedWith(pubkeyhash) && accumulated < amount {
				accumulated += out.Value
				unspentOut[txid] = append(unspentOut[txid], outid)
				if accumulated >= amount {
					break Label
				}
			}
		}
	}
	return accumulated, unspentOut
}

// FindTransaction return tx with id
func (bc *BlockChain) FindTransaction(TXid []byte) (Transaction, error) {
	bci := bc.GetIterator()
	for {
		b := bci.GetBlock()
		for _, tx := range b.Transactions {
			if bytes.Compare(tx.TXid, TXid) == 0 {
				return *tx, nil
			}
		}
		if len(b.PrevBlockHash) == 0 {
			break
		}
	}
	return Transaction{}, errors.New("No tx found with id: " + string(TXid))
}

// SignTX signs a TX
func (bc *BlockChain) SignTX(tx *Transaction, prikey ecdsa.PrivateKey) {
	PrevTX := make(map[string]Transaction)
	for _, vin := range tx.Vin {
		prevtx, _ := bc.FindTransaction(vin.TXid)
		PrevTX[hex.EncodeToString(vin.TXid)] = prevtx
	}
	tx.Sign(prikey, PrevTX)
}

// VerifyTX verifies for authenticity
func (bc *BlockChain) VerifyTX(tx *Transaction) bool {
	PrevTX := make(map[string]Transaction)
	for _, vin := range tx.Vin {
		prevtx, _ := bc.FindTransaction(vin.TXid)
		PrevTX[hex.EncodeToString(vin.TXid)] = prevtx
	}
	return tx.Verify(PrevTX)
}

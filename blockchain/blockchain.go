package main

import (
	"github.com/boltdb/bolt"
)

const dbfile = "blockchain.db"

var blocksbucket = []byte("blocks")
var chainbucket = []byte("chain")

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

// CreateBlockChain method
func CreateBlockChain() *BlockChain {

	var tip []byte
	db, _ := bolt.Open(dbfile, 0600, nil)

	_ = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(blocksbucket)
		// check if BC is already present or not
		if bucket == nil {
			bucket, _ := tx.CreateBucket(blocksbucket)
			genBlock := CreateGenesisBlock()
			bucket.Put(genBlock.BlockHash, genBlock.Serialize())
			bucket.Put([]byte("l"), genBlock.BlockHash)
			tip = genBlock.BlockHash
		} else {
			tip = bucket.Get([]byte("l"))
		}
		return nil
	})

	// b := CreateGenesisBlock()
	// blocks := []*Block{b}
	// bc := &BlockChain{blocks}

	bc := &BlockChain{tip, db}

	return bc
}

// AddBlock func adds a new block to the main BlockChain struct
func (bc *BlockChain) AddBlock(blockData string) {
	// prevBlockHash := bc.Blocks[len(bc.Blocks)-1].BlockHash
	// b := CreateBlock(blockData, prevBlockHash)
	// bc.Blocks = append(bc.Blocks, b)
	var tip []byte
	_ = bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(blocksbucket)
		tip = bucket.Get([]byte("l"))
		return nil
	})

	prevBlockHash := tip
	b := CreateBlock(blockData, prevBlockHash)

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
	return &BlockChainIterator{bc.tip, bc.db}
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

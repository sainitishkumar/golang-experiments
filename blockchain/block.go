package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
)

// Block contains info regarding the block
type Block struct {
	// Index         int64
	Timestamp int64
	// BlockData     []byte
	Transactions  []*Transaction
	PrevBlockHash []byte
	BlockHash     []byte
	Nonce         int64
}

// SetBlockHash sets a block b's hash by calculating with its params
// func (b *Block) SetBlockHash() {
// 	var temp []byte
// 	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
// 	// index := []byte(strconv.FormatInt(b.Index, 10))
// 	nonce := []byte(strconv.FormatInt(b.Nonce, 10))
// 	// temp = bytes.Join([][]byte{index, b.PrevBlockHash, timestamp, b.Data, nonce}, []byte{})
// 	temp = bytes.Join([][]byte{b.PrevBlockHash, b.BlockData, timestamp, nonce}, []byte{})
// 	hash := sha256.Sum256(temp)
// 	b.BlockHash = hash[:]
// }

// CreateGenesisBlock creates a genesis block and returns it
func CreateGenesisBlock(coinbase *Transaction) *Block {
	b := CreateBlock([]*Transaction{coinbase}, []byte{})
	return b
}

//CreateBlock creates and returns a new block with the given data
func CreateBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	b := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0}
	// b.SetBlockHash()
	pow := NewProofofWork(b)
	nonce, hash := pow.Run()
	b.Nonce = nonce
	b.BlockHash = hash
	return b
}

// Serialize converts the block struct into byte array
// this will be helpful as BoltBD(key -> value) stores in byte array format
func (b *Block) Serialize() []byte {
	var temp bytes.Buffer
	encoder := gob.NewEncoder(&temp)
	err := encoder.Encode(b)
	_ = err
	return temp.Bytes()
}

// DeSerialize converts the bytes read from BoltDB into blocks
func DeSerialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	_ = err
	return &block
}

// HashTransactions hashes all the tx in a block and returns the hash value
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txhash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.TXid)
	}
	txhash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txhash[:]
}

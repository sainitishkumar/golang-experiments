package main

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

// Block contains info regarding the block
type Block struct {
	// Index         int64
	Timestamp     int64
	BlockData     []byte
	PrevBlockHash []byte
	BlockHash     []byte
	Nonce         int64
}

// SetBlockHash sets a block b's hash by calculating with its params
func (b *Block) SetBlockHash() {
	var temp []byte
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	// index := []byte(strconv.FormatInt(b.Index, 10))
	nonce := []byte(strconv.FormatInt(b.Nonce, 10))
	// temp = bytes.Join([][]byte{index, b.PrevBlockHash, timestamp, b.Data, nonce}, []byte{})
	temp = bytes.Join([][]byte{b.PrevBlockHash, b.BlockData, timestamp, nonce}, []byte{})
	hash := sha256.Sum256(temp)
	b.BlockHash = hash[:]
}

// CreateGenesisBlock creates a genesis block and returns it
func CreateGenesisBlock() *Block {
	b := CreateBlock("Genesis block rocks", []byte{})
	return b
}

//CreateBlock creates and returns a new block with the given data
func CreateBlock(blockData string, prevBlockHash []byte) *Block {
	b := &Block{time.Now().Unix(), []byte(blockData), prevBlockHash, []byte{}, 0}
	// b.SetBlockHash()
	pow := NewProofofWork(b)
	nonce, hash := pow.Run()
	b.Nonce = nonce
	b.BlockHash = hash
	return b
}

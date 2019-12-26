package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

const targetbits = 20

var maxNonce = math.MaxInt64

// ProofofWork struct contains the req info for the miner
// contains the block that needs to be hashed and target need to be acheived
type ProofofWork struct {
	block  *Block
	target *big.Int
}

// NewProofofWork generates and returns a new POW object
func NewProofofWork(b *Block) *ProofofWork {
	target := big.NewInt(1)
	target = target.Lsh(target, uint(256-targetbits)) // so that we can compare the targets first bits
	pow := ProofofWork{b, target}
	return &pow
}

// GetData creates the data from the blockchain
// and nonce value and returns an array of bytes
// the data is combination of following
/*
PrevBlockHash,
BlockData,
Timestamp,
targetbits,
nonce
*/
func (pow *ProofofWork) GetData(nonce int64) []byte {
	var data []byte
	data = bytes.Join([][]byte{
		pow.block.PrevBlockHash,
		pow.block.HashTransactions(),
		IntToHex(pow.block.Timestamp),
		IntToHex(targetbits),
		IntToHex(nonce)}, []byte{})
	return data
}

// Run function generates the nonce value for POW and hash to validate with target bits
func (pow *ProofofWork) Run() (int64, []byte) {
	nonce := 0
	var hash [32]byte
	var hashint big.Int
	for nonce < maxNonce {
		data := pow.GetData(int64(nonce))
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashint.SetBytes(hash[:])
		// check if hash is less than the target value
		if hashint.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Printf("\n")
	return int64(nonce), hash[:]
}

// ValidatePow validates a pow by checking its hash with target
func (pow *ProofofWork) ValidatePow() bool {
	var hashint big.Int
	nonce := pow.block.Nonce
	data := pow.GetData(nonce)
	hash := sha256.Sum256(data)
	hashint.SetBytes(hash[:])

	// if hash value is less than target then it is POW
	if hashint.Cmp(pow.target) == -1 {
		return true
	}
	return false
}

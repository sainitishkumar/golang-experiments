package main

import (
	"encoding/binary"
	"fmt"
)

// IntToHex changes val into its hexa-decimal format
func IntToHex(val int64) []byte {
	b := make([]byte, binary.MaxVarintLen64)
	binary.LittleEndian.PutUint64(b, uint64(val))
	return b
}

func main() {
	bc := CreateBlockChain()
	bc.AddBlock("First block")
	bc.AddBlock("Second block")

	for _, b := range bc.Blocks {
		pow := NewProofofWork(b)
		valid := pow.ValidatePow()
		fmt.Printf("Prev. hash: %x\n", b.PrevBlockHash)
		fmt.Printf("Data: %s\n", b.BlockData)
		fmt.Printf("Hash: %x\n", b.BlockHash)
		fmt.Println("Validity by POW: ", valid)
		fmt.Println()
	}
}

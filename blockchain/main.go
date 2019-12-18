package main

import (
	"fmt"
)

func main() {
	bc := CreateBlockChain()
	bc.AddBlock("First block")
	bc.AddBlock("Second block")

	for _, b := range bc.Blocks {
		fmt.Printf("Prev. hash: %x\n", b.PrevBlockHash)
		fmt.Printf("Data: %s\n", b.BlockData)
		fmt.Printf("Hash: %x\n", b.BlockHash)
		fmt.Println()
	}
}

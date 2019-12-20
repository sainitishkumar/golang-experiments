package main

import (
	"flag"
	"fmt"
	"os"
)

type cli struct {
	bc *BlockChain
}

func (Cli *cli) PrintUsage() {
	fmt.Println("Use either print or addblock functionality")
	fmt.Println("blockchain addblock <data>")
	fmt.Println("blockchain printchain")
}

func (Cli *cli) PrintChain() {
	bcIter := Cli.bc.GetIterator()
	block := bcIter.GetBlock()
	for string(block.BlockData) != "Genesis block rocks" {
		b := block
		pow := NewProofofWork(b)
		valid := pow.ValidatePow()
		fmt.Printf("Prev. hash: %x\n", b.PrevBlockHash)
		fmt.Printf("Data: %s\n", b.BlockData)
		fmt.Printf("Hash: %x\n", b.BlockHash)
		fmt.Println("Validity by POW: ", valid)
		fmt.Println()
		block = bcIter.GetBlock()
	}
}

func (Cli *cli) Run() {
	if len(os.Args) < 2 {
		Cli.PrintUsage()
		os.Exit(2)
	}
	addblock := flag.NewFlagSet("addblock", flag.ExitOnError)
	printchain := flag.NewFlagSet("printchain", flag.ExitOnError)

	addblockdata := addblock.String("data", "", "data block_data")

	switch os.Args[1] {
	case "addblock":
		addblock.Parse(os.Args[2:])
	case "printchain":
		printchain.Parse(os.Args[2:])
	default:
		Cli.PrintUsage()
		os.Exit(2)
	}

	if addblock.Parsed() {
		if *addblockdata == "" {
			addblock.Usage()
			os.Exit(2)
		}
		Cli.bc.AddBlock(*addblockdata)
	} else if printchain.Parsed() {
		Cli.PrintChain()
	}

}

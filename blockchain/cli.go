package main

import (
	"flag"
	"fmt"
	"os"
)

// Cli struct
type Cli struct {
	bc *BlockChain
}

// PrintUsage prints the usage of the CLI
func (cli *Cli) PrintUsage() {
	fmt.Println("Use either print or addblock functionality")
	fmt.Println("blockchain addblock <data>")
	fmt.Println("blockchain printchain")
}

// PrintChain prints the BC from newest to oldest
func (cli *Cli) PrintChain() {
	bcIter := cli.bc.GetIterator()
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

// Run is for parsing CL args and execute them
func (cli *Cli) Run() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
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
		cli.PrintUsage()
		os.Exit(2)
	}

	if addblock.Parsed() {
		if *addblockdata == "" {
			addblock.Usage()
			os.Exit(2)
		}
		cli.bc.AddBlock(*addblockdata)
	} else if printchain.Parsed() {
		cli.PrintChain()
	}
}

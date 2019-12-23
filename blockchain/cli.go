package main

import (
	"flag"
	"fmt"
	"log"
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
	// bcIter := cli.bc.GetIterator()
	// block := bcIter.GetBlock()
	// for {
	// 	b := block
	// 	pow := NewProofofWork(b)
	// 	valid := pow.ValidatePow()
	// 	fmt.Printf("Prev. hash: %x\n", b.PrevBlockHash)
	// 	fmt.Printf("Hash: %x\n", b.BlockHash)
	// 	fmt.Println("Validity by POW: ", valid)
	// 	fmt.Println()
	// 	block = bcIter.GetBlock()
	// }
	bc := CreateBlockChain("")
	defer bc.db.Close()
	bci := bc.GetIterator()
	for {
		block := bci.GetBlock()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.BlockHash)
		// pow := NewProofofWork(block)

		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *Cli) createBlockchain(address string) {
	bc := CreateBlockChain(address)
	bc.db.Close()
	fmt.Println("Done!")
}

// Run is for parsing CL args and execute them
func (cli *Cli) Run() {
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	// getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	// sendFrom := sendCmd.String("from", "", "Source wallet address")
	// sendTo := sendCmd.String("to", "", "Destination wallet address")
	// sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.PrintUsage()
		os.Exit(1)
	}

	// if getBalanceCmd.Parsed() {
	// 	if *getBalanceAddress == "" {
	// 		getBalanceCmd.Usage()
	// 		os.Exit(1)
	// 	}
	// 	cli.getBalance(*getBalanceAddress)
	// }

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.PrintChain()
	}

	// if sendCmd.Parsed() {
	// 	if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
	// 		sendCmd.Usage()
	// 		os.Exit(1)
	// 	}

	// 	cli.send(*sendFrom, *sendTo, *sendAmount)
	// }
}

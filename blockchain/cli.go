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
	bc := GetBlockChain("")
	defer bc.db.Close()
	bci := bc.GetIterator()
	for {
		block := bci.GetBlock()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.BlockHash)
		pow := NewProofofWork(block)
		fmt.Println("POW: ", pow.ValidatePow())
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

// getBalance returns the number of coins address holds
func (cli *Cli) getBalance(address string) int {
	bc := GetBlockChain(address)
	defer bc.db.Close()

	// fmt.Println(bc.db.GoString())

	balance := 0
	unspentTX := bc.FindUTXO(address)
	for _, out := range unspentTX {
		balance += out.Value
	}
	fmt.Printf("Balance of %s is %d", address, balance)
	return balance
}

// send func to send coins
func (cli *Cli) send(from, to string, amount int) {
	bc := GetBlockChain(from)
	defer bc.db.Close()

	tx := NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*Transaction{tx})
	fmt.Println("Sent ", amount, "from", from, "to", to)
}

// Run is for parsing CL args and execute them
func (cli *Cli) Run() {
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
		}
	default:
		cli.PrintUsage()
		os.Exit(2)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}

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

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}

package main

// BlockChain is structure containing slice of blocks
type BlockChain struct {
	Blocks []*Block
}

// CreateBlockChain creates an slice containing the genesis block
func CreateBlockChain() *BlockChain {
	b := CreateGenesisBlock()
	blocks := []*Block{b}
	bc := &BlockChain{blocks}
	return bc
}

// AddBlock func adds a new block to the main BlockChain struct
func (bc *BlockChain) AddBlock(blockData string) {
	prevBlockHash := bc.Blocks[len(bc.Blocks)-1].BlockHash
	b := CreateBlock(blockData, prevBlockHash)
	bc.Blocks = append(bc.Blocks, b)
}

package main

import (
	"encoding/binary"
)

// IntToHex changes val into its hexa-decimal format
func IntToHex(val int64) []byte {
	b := make([]byte, binary.MaxVarintLen64)
	binary.LittleEndian.PutUint64(b, uint64(val))
	return b
}

func main() {
	bc := CreateBlockChain()
	defer bc.db.Close()

	Cli := cli{bc}
	Cli.Run()
}

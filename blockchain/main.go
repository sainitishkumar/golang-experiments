// utils and other helper functions
package main

import (
	"bytes"
	"encoding/binary"
	"math/big"

	"github.com/btcsuite/btcutil/base58"
)

var (
	b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")
)

// Base58Encode base58 encoding of the payload
func Base58Encode(payload []byte) []byte {
	enc := []byte(base58.Encode(payload))
	return enc
}

// Base58Decode decode from base58 to original
func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	zeroBytes := 0

	for b := range input {
		if b == 0x00 {
			zeroBytes++
		}
	}

	payload := input[zeroBytes:]
	for _, b := range payload {
		charIndex := bytes.IndexByte(b58Alphabet, b)
		result.Mul(result, big.NewInt(58))
		result.Add(result, big.NewInt(int64(charIndex)))
	}

	decoded := result.Bytes()
	decoded = append(bytes.Repeat([]byte{byte(0x00)}, zeroBytes), decoded...)

	return decoded
}

// IntToHex changes val into its hexa-decimal format
func IntToHex(val int64) []byte {
	b := make([]byte, binary.MaxVarintLen64)
	binary.LittleEndian.PutUint64(b, uint64(val))
	return b
}

func main() {
	cli := CLI{}
	cli.Run()
}

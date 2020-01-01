package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"golang.org/x/crypto/ripemd160"
)

const (
	version          = byte(0x00)
	addressChkSumLen = 4
)

// Wallet struct
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// NewKeyPair creates a private key using elliptic curve cryptography
func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	ellipticCurve := elliptic.P256()
	prikey, _ := ecdsa.GenerateKey(ellipticCurve, rand.Reader)
	pubkey := append(prikey.X.Bytes(), prikey.Y.Bytes()...)
	return *prikey, pubkey
}

// NewWallet returns a new wallet
func NewWallet() *Wallet {
	prikey, pubkey := NewKeyPair()
	wallet := Wallet{prikey, pubkey}
	return &wallet
}

// HashPubKey returns RIMEMD160(SHA256(PubKey))
func HashPubKey(pubkey []byte) []byte {
	shapub := sha256.Sum256(pubkey)
	ripemd := ripemd160.New()
	ripemd.Write(shapub[:])
	hash := ripemd.Sum(nil)
	return hash
}

// PayloadCheckSum return payload checksum
func PayloadCheckSum(payload []byte) []byte {
	temp := sha256.Sum256(payload)
	sha := sha256.Sum256(temp[:])
	chksum := sha[:addressChkSumLen]
	return chksum
}

// GetAddress func ref: resources/address-generation-scheme.png
// address = base58encode(version+pubhash+checksum)
func (w Wallet) GetAddress() []byte {
	pubhash := HashPubKey(w.PublicKey)
	versionPayload := append([]byte{version}, pubhash...)
	chkSum := PayloadCheckSum(versionPayload)
	payload := append(versionPayload, chkSum...)
	address := Base58Encode(payload)
	return address
}

// ValidateAddress Validates given Blockchain address
func ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChkSumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChkSumLen]
	targetChecksum := PayloadCheckSum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

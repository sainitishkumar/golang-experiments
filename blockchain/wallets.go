package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
)

const (
	walletFile = "wallet.dat"
)

// Wallets struct
type Wallets struct {
	Wallets map[string]*Wallet
}

// CreateWallet creates a wallet and returns its address
func (wallets *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := fmt.Sprintf("%s", wallet.GetAddress())
	wallets.Wallets[address] = wallet
	return address
}

// LoadfromFile loads from the walletFile
func (wallets *Wallets) LoadfromFile() {
	filecontent, _ := ioutil.ReadFile(walletFile)
	var wallet Wallets
	gob.Register(elliptic.P256())
	dec := gob.NewDecoder(bytes.NewReader(filecontent))
	dec.Decode(&wallet)
	wallets.Wallets = wallet.Wallets
}

// NewWallets load from wallet file
func NewWallets() *Wallets {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)
	wallets.LoadfromFile()
	return &wallets
}

// SavetoFile func
func (wallets *Wallets) SavetoFile() {
	var content bytes.Buffer
	gob.Register(elliptic.P256())
	enc := gob.NewEncoder(&content)
	enc.Encode(wallets)
	ioutil.WriteFile(walletFile, content.Bytes(), 0644)
}

// GetWallet func
func (wallets *Wallets) GetWallet(address string) Wallet {
	return *wallets.Wallets[address]
}

// GetAddresses func
func (wallets *Wallets) GetAddresses() []string {
	var addresses []string
	for address := range wallets.Wallets {
		addresses = append(addresses, address)
	}
	return addresses
}

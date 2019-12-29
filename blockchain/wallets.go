package main

// Wallets struct
type Wallets struct {
	Wallets map[string]*Wallet
}

// LoadfromFile loads from the walletFile
func (wallets *Wallets) LoadfromFile() {

}

// NewWallets create new wallets
func NewWallets() *Wallets {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)
	wallets.LoadfromFile()
	return &wallets
}

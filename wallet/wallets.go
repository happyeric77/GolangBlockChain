package wallet

import (
	"crypto/elliptic"
	"fmt"
	"os"
	"io/ioutil"
	"log"
	"encoding/gob"
	"bytes"
)

const (
	walletFile = "./tmp/wallets.data"
)

type Wallets struct {
	Wallets map[string]*Wallet
}

func (ws *Wallets) LoadFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	var wallets Wallets

	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		return err
	}

	gob.Register(elliptic.P256())
	err = gob.NewDecoder(bytes.NewReader(fileContent)).Decode(&wallets)
	if err != nil {
		return err
	}

	ws.Wallets = wallets.Wallets
	return nil
}

func CreateWallets() (*Wallets, error){
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)
	err := wallets.LoadFile()
	return &wallets, err
}

func (ws *Wallets) AddWallet() string {
	wallet := MakeWallet()
	address := fmt.Sprintf("%s", wallet.Address())
	// address := string(wallet.Address())
	ws.Wallets[address] = wallet
	return address
}

func (ws Wallets) GetAllWalletAddress() []string {
	var addresses []string
	for add := range ws.Wallets{
		addresses = append(addresses, add)
	}
	return addresses
}

func (ws Wallets) GetWallet (address string) Wallet {
	return *ws.Wallets[address]
}


func (ws *Wallets) SaveFile() {
	content := bytes.NewBuffer(nil)

	gob.Register(elliptic.P256())

	err := gob.NewEncoder(content).Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile, []byte(content.Bytes()), 0644)
	if err != nil {
		log.Panic(err)
	}
}
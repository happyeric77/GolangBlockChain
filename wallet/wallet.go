package wallet

import (
	"crypto/sha256"
	"log"
	"crypto/rand"
	"crypto/elliptic"
	"crypto/ecdsa"
	"golang.org/x/crypto/ripemd160"
	"fmt"
)

const (
	checksumLength = 4
	version = byte(0x00)
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey []byte
}

func (w Wallet) Address() []byte {
	pubkeyHash := PublicKeyHash(w.PublicKey)
	versionedHash := append([]byte{version}, pubkeyHash...)
	checksum := Checksum(versionedHash)
	fullHash := append(versionedHash, checksum...)
	address := Base58Encode(fullHash)
	fmt.Printf("pubkey: %x\n", w.PublicKey)
	fmt.Printf("pubkey Hash: %x\n", pubkeyHash)
	fmt.Printf("address: %s\n", address)
	return address
}

func NewKeyPair() (ecdsa.PrivateKey, []byte){
	curve := elliptic.P256() // Output of the elliptic curve (interface type) will be 256 bytes
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pub := append(private.X.Bytes(), private.Y.Bytes()...)
	return *private, pub
}

func MakeWallet() *Wallet {
	private, public := NewKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}

func PublicKeyHash(pubkey []byte) []byte {
	pubHash := sha256.Sum256(pubkey)
	hasher := ripemd160.New()
	_, err := hasher.Write(pubHash[:])
	if err != nil {
		log.Panic(err)
	}
	publicRipMD := hasher.Sum(nil)
	return publicRipMD
}

func Checksum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])
	return secondHash[:checksumLength]
}
package blockchain

import (
	"crypto/sha256"
	"log"
	"encoding/binary"
	"bytes"
	"math/big"
	"math"
	"fmt"
)

const Difficulty uint = 13

type ProofOfWork struct {
	Block *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256 - Difficulty))
	pow := &ProofOfWork{ b, target}
	return pow
}

func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			ToByte(int64(nonce)),
			ToByte(int64(Difficulty)),
		}, 
		[]byte{},
	)
	return data
}

func ToByte (num int64) []byte {
	buff := new (bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var byteHash [32]byte
	nonce := 0

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		byteHash = sha256.Sum256(data)
		fmt.Printf("%x\n", byteHash)
		intHash.SetBytes(byteHash[:])

		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce ++
		}
	}
	fmt.Println("=============")
	return nonce, byteHash[:]
}
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int
	data := pow.InitData(pow.Block.Nonce)
	byteHash := sha256.Sum256(data)
	intHash.SetBytes(byteHash[:])
	return intHash.Cmp(pow.Target) == -1
}
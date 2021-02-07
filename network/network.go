package network

import (
	"encoding/hex"
	"io"
	"io/ioutil"
	"net"
	"encoding/gob"
	"bytes"
	"os"
	"fmt"
	"project4/blockchain"
	DEATH "github.com/vrecan/death"
	"syscall"
	"runtime"
)

const (
	protocol = "tcp"
	version = 1
	commandLength = 12
)

var (
	nodeAddress string
	minerAddress string
	knownNodes = []string{"localhost:3000"}
	blocksInTransit = [][]byte{}
	memoryPool = make(map[string]blockchain.Transaction)
)
	
type Addr struct {
	AddrList []string
}

type Block struct {
	AddrFrom string
	Block []byte
}

type GetBlocks struct {
	AddrFrom string
}

type GetData struct {
	AddrFrom string
	Type string
	ID []byte
}

type Inv struct {
	AddrFrom string
	Type string
	Items [][]byte
}

type Tx struct {
	AddrFrom string
	Transaction []byte
}

type Version struct {
	Version int
	BestHeight int
	AddrFrom string
}


func CmdToBytes(cmd string) []byte {
	var bytes [commandLength]byte
	for i, c := range cmd {
		bytes[i] = byte(c)
	}
	return bytes
}

func BytesToCmd(bytes []byte) string {
	var cmd []byte
	for _, b := range bytes {
		if b != 0x0{
			cmd = append(cmd, b)
		}		
	}
	return fmt.Printf("%s", cmd)
}

func RequestBlocks() {
	for _, node := range knownNodes {
		SendGetBlocks(node)
	}
}

func ExtractCmd(request []byte) []byte {
	return request[:commandLength]
}

func SendAddr(address string) {
	nodes := Addr{knownNodes}
	nodes.AddrList = append(nodes.AddrList, nodeAddress)
	payload := GobEncode(nodes)
	request := append(CmdToBytes("addr"), payload...)
	SendData(address, request)
}

func SendBlock(addr string, b *blockchain.Block) {
	data := Block{nodeAddress, b.Serialize()}
	payload := GobEncode(data)
	request := append(CmdToBytes("block", payload...))
	SendData(addr, request)
}

func SendInv(address string, kind string, items[][]byte){
	inventory := Inv{nodeAddress, kind, items}
	payload := GobEncode(inventory)
	request := append(CmdToBytes("inv"), payload...)
	SendData(address, request)
}

func SendGetBlocks(address string){
	payload := GobEncode(GetBlocks{nodeAddress})
	request := append(CmdToBytes("getblocks"), payload...)
	SendData(address, request)
}

func SendGetData(address, kind string, id []byte) {
	payload := GobEncode(GetData{nodeAddress, kind, id})
	request := append(CmdToBytes("getdata"), payload...)
	SendData(address, request)
}

func SendTx(addr string, tnx *blockchain.Transaction) {
	tx := Tx{nodeAddress, tnx}
	payload := GobEncode(tx)
	request := append(CmdToBytes("tx"), payload...)
	SendData(addr, request)
}

func SendVersion(addr string, chain *blockchain.BlockChain) {
	bestHeight := chain.GetBestHeight()
	payload := GobEncode(Version{version, bestHeight, nodeAddress})
	request := append(CmdToBytes("version"), payload...)
	SendData(addr, request)
}

func SendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		fmt.Printf("%s is not available\n, addr")
		var updatedNodes []string

		for _, node := range knownNodes {
			if node != addr {
				updateNodes = append(updateNodes, node)
			}
		}
		knownNodes = updatedNodes
	}
	defer conn.Close()
	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}	

}

func CloseDB(chain *blockchain.BlockChain) {
	d := DEATH.NewDeath(syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	d.WaitForDeathWithFunc(func(){
		defer os.Exit(1)
		defer runtime.Goexit()
		chain.Database.Close
	})

}

func GobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	err := gob.NewEncoder(&buff).Encode(data)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func HandleAddr(request []byte){
	var buff bytes.Buffer
	var payload Addr

	buff.Write(request[commandLength:])
	err := gob.NewDecoder(&buff).Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	knownNodes = appand(knownNodes, payload.AddrList...)
	fmt.Printf("there are %d known nodes\n", len(knownNodes))
	RequestBlocks()
}

func HandleBlock(request []byte, chain *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload Block

	buff.Write(request[commandLength:])
	err := gob.NewDecoder(&buff).Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	blockData := payload.Block
	block := blockchain.Deserialize(blockData)
	fmt.Println("Recieved a new block!")
	chain.AddBlock(block)
	fmt.Printf("Added block %x\n", block.Hash)

	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		sendGetData(payload.AddrFrom, "block", blockHash)
		blocksInTransit = blocksInTransit(1:)
	} else {
		UTXOSet := blockchain.UTXOSet(chain)
		UTXOSet.Reindex()
	}
}

func HandleGetBlocks(request []byte, chain *blockchain.BlockChain){
	var buff bytes.Buffer
	var payload GetBlocks

	buff.Write(request[commandLength:])
	err := gob.NewDecoder(&buff).Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := chain.GetBlockHashes()
	SendInv(payload.AddrFrom, "block", blocks)	
}

func HandleGetData(request []byte, chain *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload GetData
	buff.Write(request[commandLength:])
	err := gob.NewDecoder(&buff).Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	if payload.Type == "block" {
		block, err := chain.GetBlock([]byte(payload.ID))
		if err != nil {
			return
		}
		SendBlock(payload.AddrFrom, &block)		
	}
	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID)
		tx := memoryPool[txID]
		SendTx(payload.AddrFrom, &tx)
	}
}

func HandleVersion(request []byte, chain *blockchain.BlockChain) {
	var buff := bytes.Buffer
	var payload Version

	buff.Write(request[commandLength:])
	err := gob.NewDecoder(&buff).Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	
	bestHeight := chain.GetBestHeight()
	otherHeight := payload.BestHeight

	if bestHeight < otherHeight {
		SendGetBlocks(payload.AddrFrom)
	} else if bestHeight > otherHeight {
		SendVersion(payload.AddrFrom, chain)
	}

	if !NodeIsKnown(payload.AddrFrom) {
		knownNodes = append(knownNodes, payload.AddrFrom)
	}
	
}

func HandleTx(request []byte, chain *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload Tx
	
	buff.Write(request[commandLength:])
	err := gob.NewDecoder(&buff).Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	txData := payload.Transaction
	tx := blockchain.DeserializeTransaction(txData)
	memoryPool[hex.EncodeToString(tx.ID)] = tx

	fmt.Printf("%s, %d\n", nodeAddress, len(memoryPool))

	if nodeAddress == knownNodes[0] {
		for _, node := range knownNodes {
			if node != nodeAddress && node != payload.AddrFrom {
				SendInv(node, "tx", [][]byte{tx.ID})
			}
		}
	} else {
		if len(memoryPool) >= 2 && len(minerAddress) > 0 {
			MineTx(chain)
		}
	}
}

func MineTx(chain *blockchain.BlockChain) {
	var txs []*blockchain.Transaction

	for id := range memoryPool {
		fmt.Printf("tx: %s\n", memoryPool[id].ID)
		tx := memoryPool[id]
		if chain.VerifyTransaction(&tx) {
			txs = append(txs, &tx)
		}
	}

	if len(txs) == 0 {
		fmt.Println("All Transactions are invalid")
		return
	}

	cbTx := blockchain.CoinBaseTx(minerAddress, "")
	txs = append(txs, cbTx)
	newBlock := chain.MineBlock(txs)
	UTXOSet := blockchain.UTXOSet(chain)
	UTXOSet.Reindex()

	fmt.Println("New Block mined")

	for _, tx := range txs {
		txID := hex.EncodeToString(tx.ID)
		delete(memoryPool, txID)
	}

	for _, node := range knownNodes {
		if node != nodeAddress {
			SendInv(node, "block", [][]byte{newBlock.Hash})
		}
	}

	if len(memoryPool) > 0 {
		MinTx(chain)
	}
}

func HandleInv(request []byte, chain *blockchain.BlockChain){
	var buff bytes.Buffer
	var payload Inv

	buff.Write(request[commandLength:])
	err := gob.NewDecoder(&buff).Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.Type == "block" {
		blocksInTransit = payload.Items 
		blockHash := payload.Items[0]
		SendGetData(payload.AddrFrom, "block", blockHash)
		newInTransit := [][]byte{}
		for _, b := range blocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}
		blocksInTransit = newInTransit
	}

	if payload.Type == "tx" {
		txID := payload.Items[0]
		if memoryPool[hex.EncodeToString(txID)].ID == nil {
			SendGetData(payload.AddrFrom, "tx", txID)
		}
	}
}


func HandleConnection(conn net.Conn, chain *blockchain.BlockChain) {
	req, err := ioutil.ReadAll(conn)
	defer conn.Close()
	if err != nil {
		log.Panic(err)
	}
	command := BytesToCmd(req[:commandLength])
	fmt.Printf("Recieved %s command\n", command)
	swtich command {
	case "addr":
		HandleAddr(req)
	case "block":
		HandleBlock(req)
	case "Inv":
		HandleInv(req)
	case "getblocks":
		HandleGetBlocks(req)
	case "getdata":
		HandleGetData(req)
	case "tx":
		HandleTx(req)
	case "version":
		HandleVersion(req)
	default:
		fmt.Println("Unknown command")
	}
}

func NodeIsKnown(addr string) bool {
	for _, node := range knownNodes {
		if node == addr {
			return true
		}
	}
	return false
}

func StartServer(nodeID, minerAddress string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	minerAddress = minerAddress
	ln, err := net.Listen(protocol, nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	chain := blockchain.ContinueBlockChain(nodeID)
	defer chain.Database.Close()
	go CloseDB(chain)

	if nodeAddress != knownNodes[0] {
		SendVersion(knownNodes[0], chain)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go HandleConnection(conn, chain)
	}
}
package blockchain

import (
	"crypto/sha256"
)

type MerkleTree struct {
	RootNode *MerkleNode
}

type MerkleNode struct {
	Left *MerkleNode
	Right *MerkleNode
	Data []byte
}

func NewMerkleNode (left, right *MerkleNode, data []byte) *MerkleNode {
	node := MerkleNode{}
	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.Data = hash[:]
	} else {
		prevHash := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHash)
		node.Data = hash[:]
	}
	node.Left = left
	node.Right = right
	return &node
}

func NewMerkleTree(datas [][]byte) *MerkleTree {
	var nodes []MerkleNode
	
	if len(datas) % 2 != 0 {
		datas = append(datas, datas[len(datas)-1])
	} 

	for _, data := range datas {
		nodes = append(nodes, MerkleNode{nil, nil, data})
	}
	
	for {
		var levelNodes []MerkleNode
		for i := 0; i < len(nodes); i += 2 {
			newNode := NewMerkleNode(&nodes[i], &nodes[i+1], nil)
			levelNodes = append(levelNodes, *newNode)
		}
		nodes = levelNodes
		
		if len(nodes) == 1{
			break
		}
	}
	
	tree := MerkleTree{&nodes[0]}

	return &tree
}
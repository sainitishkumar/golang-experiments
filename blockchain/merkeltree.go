package main

import "crypto/sha256"

// MerkelNode stores left and right pointers
type MerkelNode struct {
	left  *MerkelNode
	right *MerkelNode
	Data  []byte
}

// MerkelTree struct
type MerkelTree struct {
	RootNode *MerkelNode
}

// NewMerkelNode func
func NewMerkelNode(left, right *MerkelNode, data []byte) *MerkelNode {
	var node MerkelNode
	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.Data = hash[:]
		node.left = nil
		node.right = nil
	} else {
		hash := append(left.Data, right.Data...)
		node.Data = hash[:]
		node.left = left
		node.right = right
	}
	return &node
}

// NewMerkelTree create a new mkl tree and return the rootnode
func NewMerkelTree(data [][]byte) *MerkelTree {
	var nodes []MerkelNode
	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}
	for _, dat := range data {
		temp := NewMerkelNode(nil, nil, dat)
		nodes = append(nodes, *temp)
	}
	for i := 0; i < len(data)/2; i++ {
		var lvl []MerkelNode
		for j := 0; j < len(nodes); j += 2 {
			temp := NewMerkelNode(&nodes[j], &nodes[j+1], nil)
			lvl = append(lvl, *temp)
		}
		nodes = lvl
	}
	mkltree := MerkelTree{&nodes[0]}
	return &mkltree
}

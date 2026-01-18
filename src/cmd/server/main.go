package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"secure-logging/src/pkg"

	"github.com/gin-gonic/gin"
)

type Node struct {
	Blockchain *pkg.Blockchain
	Nodes      map[string]bool
	Mu         sync.Mutex
	Port       string
}

func NewNode(port string) *Node {
	return &Node{
		Blockchain: pkg.NewBlockchain(4),
		Nodes:      make(map[string]bool),
		Port:       port,
	}
}

func (n *Node) AddLog(c *gin.Context) {
	var input struct {
		Event     string `json:"event"`
		PublicKey string `json:"public_key"`
		Signature string `json:"signature"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing values"})
		return
	}

	// Verify Signature
	if err := pkg.VerifyEvent(input.PublicKey, input.Event, input.Signature); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	n.Mu.Lock()
	defer n.Mu.Unlock()

	lastBlock := n.Blockchain.GetLastBlock()
	data := pkg.LogData{
		Event:     input.Event,
		PublicKey: input.PublicKey,
		Signature: input.Signature,
	}

	proof, ts := n.Blockchain.ProofOfWork(lastBlock.Index+1, data, lastBlock.Hash)

	newBlock := pkg.Block{
		Index:        lastBlock.Index + 1,
		Timestamp:    ts,
		Data:         data,
		PreviousHash: lastBlock.Hash,
		Proof:        proof,
	}
	newBlock.Hash = newBlock.CalculateHash()

	if n.Blockchain.AddBlock(newBlock) {
		go n.BroadcastBlock(newBlock)
		c.JSON(http.StatusCreated, gin.H{"message": "Log added successfully", "block": newBlock})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add block"})
	}
}

func (n *Node) GetChain(c *gin.Context) {
	n.Mu.Lock()
	defer n.Mu.Unlock()
	c.JSON(http.StatusOK, gin.H{
		"chain":  n.Blockchain.Chain,
		"length": len(n.Blockchain.Chain),
	})
}

func (n *Node) RegisterNodes(c *gin.Context) {
	var input struct {
		Nodes []string `json:"nodes"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	n.Mu.Lock()
	for _, node := range input.Nodes {
		n.Nodes[node] = true
	}
	totalNodes := []string{}
	for k := range n.Nodes {
		totalNodes = append(totalNodes, k)
	}
	n.Mu.Unlock()

	c.JSON(http.StatusCreated, gin.H{"message": "New nodes added", "total_nodes": totalNodes})
}

func (n *Node) ReceiveBlock(c *gin.Context) {
	var block pkg.Block
	if err := c.ShouldBindJSON(&block); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid block data"})
		return
	}

	n.Mu.Lock()
	defer n.Mu.Unlock()

	if n.Blockchain.AddBlock(block) {
		c.String(http.StatusCreated, "Block accepted")
	} else {
		c.String(http.StatusBadRequest, "Block rejected")
	}
}

func (n *Node) BroadcastBlock(block pkg.Block) {
	n.Mu.Lock()
	nodes := []string{}
	for node := range n.Nodes {
		nodes = append(nodes, node)
	}
	n.Mu.Unlock()

	blockData, _ := json.Marshal(block)
	for _, node := range nodes {
		url := fmt.Sprintf("http://%s/receive-block", node)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(blockData))
		if err == nil {
			resp.Body.Close()
		}
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	node := NewNode(port)

	r := gin.Default()
	r.POST("/logs", node.AddLog)
	r.GET("/chain", node.GetChain)
	r.POST("/nodes/register", node.RegisterNodes)
	r.POST("/receive-block", node.ReceiveBlock)

	log.Printf("Node starting on port %s", port)
	r.Run(":" + port)
}

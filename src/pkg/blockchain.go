package pkg

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type LogData struct {
	Event     string `json:"event"`
	PublicKey string `json:"public_key"`
	Signature string `json:"signature"`
}

type Block struct {
	Index        int       `json:"index"`
	Timestamp    int64     `json:"timestamp"`
	Data         interface{} `json:"data"`
	PreviousHash string    `json:"previous_hash"`
	Proof        int       `json:"proof"`
	Hash         string    `json:"hash"`
}

func (b *Block) CalculateHash() string {
	record := fmt.Sprintf("%d%d%v%s%d", b.Index, b.Timestamp, b.Data, b.PreviousHash, b.Proof)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

type Blockchain struct {
	Chain      []Block
	Difficulty int
}

func NewBlockchain(difficulty int) *Blockchain {
	bc := &Blockchain{
		Difficulty: difficulty,
	}
	bc.createGenesisBlock()
	return bc
}

func (bc *Blockchain) createGenesisBlock() {
	genesis := Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		Data:         "Genesis Block",
		PreviousHash: "0",
		Proof:        0,
	}
	genesis.Hash = genesis.CalculateHash()
	bc.Chain = append(bc.Chain, genesis)
}

func (bc *Blockchain) GetLastBlock() Block {
	return bc.Chain[len(bc.Chain)-1]
}

func (bc *Blockchain) IsValidNewBlock(newBlock Block, previousBlock Block) bool {
	if previousBlock.Index+1 != newBlock.Index {
		return false
	}
	if previousBlock.Hash != newBlock.PreviousHash {
		return false
	}
	if newBlock.Hash != newBlock.CalculateHash() {
		return false
	}
	if !strings.HasPrefix(newBlock.Hash, strings.Repeat("0", bc.Difficulty)) {
		return false
	}
	return true
}

func (bc *Blockchain) ProofOfWork(index int, data interface{}, prevHash string) (int, int64) {
	proof := 0
	timestamp := time.Now().Unix()
	for {
		record := fmt.Sprintf("%d%d%v%s%d", index, timestamp, data, prevHash, proof)
		h := sha256.New()
		h.Write([]byte(record))
		hash := hex.EncodeToString(h.Sum(nil))
		if strings.HasPrefix(hash, strings.Repeat("0", bc.Difficulty)) {
			return proof, timestamp
		}
		proof++
	}
}

func (bc *Blockchain) AddBlock(newBlock Block) bool {
	if bc.IsValidNewBlock(newBlock, bc.GetLastBlock()) {
		bc.Chain = append(bc.Chain, newBlock)
		return true
	}
	return false
}

func (bc *Blockchain) IsChainValid(chain []Block) bool {
	if len(chain) == 0 {
		return false
	}
	// Simplified validation for this version
	for i := 1; i < len(chain); i++ {
		if !bc.IsValidNewBlock(chain[i], chain[i-1]) {
			return false
		}
	}
	return true
}

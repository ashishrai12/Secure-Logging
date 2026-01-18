package main

import (
	"testing"
)

func TestBlockchainGenesis(t *testing.T) {
	bc := NewBlockchain(1)
	if len(bc.Chain) != 1 {
		t.Errorf("Expected 1 block, got %d", len(bc.Chain))
	}
	if bc.Chain[0].Data != "Genesis Block" {
		t.Errorf("Expected Genesis Block, got %v", bc.Chain[0].Data)
	}
}

func TestAddBlock(t *testing.T) {
	bc := NewBlockchain(1)
	lastBlock := bc.GetLastBlock()
	data := "Test Log"
	proof, _ := bc.ProofOfWork(lastBlock.Index+1, data, lastBlock.Hash)

	newBlock := Block{
		Index:        lastBlock.Index + 1,
		Timestamp:    0, // ProofOfWork returns the TS it used, but we'll manually test
		Data:         data,
		PreviousHash: lastBlock.Hash,
		Proof:        proof,
	}
	newBlock.Hash = newBlock.CalculateHash()

	if !bc.AddBlock(newBlock) {
		t.Errorf("Failed to add valid block")
	}
}

func TestIdentitySigning(t *testing.T) {
	id, _ := NewIdentity()
	data := "Important System Event"
	sig, _ := id.SignEvent(data)
	pubPEM, _ := id.GetPublicKeyPEM()

	if err := VerifyEvent(pubPEM, data, sig); err != nil {
		t.Errorf("Signature verification failed: %v", err)
	}

	if err := VerifyEvent(pubPEM, "Tampered Data", sig); err == nil {
		t.Errorf("Tampered data should not pass verification")
	}
}

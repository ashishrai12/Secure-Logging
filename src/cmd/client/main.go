package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"secure-logging/src/pkg"
)

func main() {
	nodeURL := os.Getenv("NODE_URL")
	if nodeURL == "" {
		nodeURL = "http://localhost:5001"
	}

	fmt.Println("Initializing Client Identity...")
	id, err := pkg.NewIdentity()
	if err != nil {
		log.Fatalf("Failed to create identity: %v", err)
	}

	pubKeyPEM, _ := id.GetPublicKeyPEM()
	
	event := "SECURITY_ALERT: Unauthorized login attempt detected on production server"
	fmt.Printf("Signing event: %s\n", event)

	signature, err := id.SignEvent(event)
	if err != nil {
		log.Fatalf("Failed to sign event: %v", err)
	}

	payload := map[string]string{
		"event":      event,
		"public_key": pubKeyPEM,
		"signature":  signature,
	}

	jsonPayload, _ := json.Marshal(payload)

	fmt.Printf("Submitting log to node at %s...\n", nodeURL)
	resp, err := http.Post(nodeURL+"/logs", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Fatalf("Failed to submit log: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Node Response: %s\n", string(body))

	// Fetch chain
	fmt.Println("\nFetching Blockchain...")
	resp, err = http.Get(nodeURL + "/chain")
	if err != nil {
		log.Fatalf("Failed to fetch chain: %v", err)
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Current Blockchain: %s\n", string(body))
}

package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

type Identity struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func NewIdentity() (*Identity, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return &Identity{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
	}, nil
}

func (id *Identity) GetPublicKeyPEM() (string, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(id.PublicKey)
	if err != nil {
		return "", err
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})
	return string(pubPEM), nil
}

func (id *Identity) SignEvent(eventData string) (string, error) {
	hashed := sha256.Sum256([]byte(eventData))
	signature, err := rsa.SignPSS(rand.Reader, id.PrivateKey, crypto.SHA256, hashed[:], &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
	})
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

func VerifyEvent(publicKeyPEM string, eventData string, signatureB64 string) error {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return errors.New("failed to parse PEM block containing the public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return errors.New("not an RSA public key")
	}

	signature, err := base64.StdEncoding.DecodeString(signatureB64)
	if err != nil {
		return err
	}

	hashed := sha256.Sum256([]byte(eventData))
	return rsa.VerifyPSS(rsaPub, crypto.SHA256, hashed[:], signature, &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
	})
}

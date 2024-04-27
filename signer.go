package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

func setSigner(fpath string) (ssh.Signer, error) {
	var err error

	keyData, err := os.ReadFile(fpath)
	if err != nil {
		return nil, fmt.Errorf("failed reading private key file: %v", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("invalid private key format")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %v", err)
	}

	signer, err := ssh.NewSignerFromKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("error creating signer: %v", err)
	}

	return signer, nil
}

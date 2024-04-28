package main // import "keygen"

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func savePrivateKey() (privateKey *rsa.PrivateKey, err error) {
	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %s", err)
	}

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	privateKeyFile, err := os.Create("private_key.pem")
	if err != nil {
		return nil, fmt.Errorf("failed to create private key file: %s", err)
	}
	defer privateKeyFile.Close()

	if err = pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return nil, fmt.Errorf("failed to write private key to file: %s", err)
	}

	return privateKey, nil
}

func savePublicKey(privateKey *rsa.PrivateKey) (err error) {
	publicKey := privateKey.PublicKey

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %s", err)
	}

	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	publicKeyFile, err := os.Create("public_key.pem")
	if err != nil {
		return fmt.Errorf("failed to create public key file: %s", err)
	}
	defer publicKeyFile.Close()

	if err := pem.Encode(publicKeyFile, publicKeyPEM); err != nil {
		return fmt.Errorf("failed to write public key to file: %s", err)
	}

	return nil
}

func main() {
	privateKey, err := savePrivateKey()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Println("Private key generated and saved to private_key.pem")

	err = savePublicKey(privateKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	fmt.Println("Public key generated and saved to public_key.pem")
}

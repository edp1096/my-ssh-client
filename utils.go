package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"io"
	"os"
)

func ContainsMapKey[K comparable, V any](m map[K]V, key K) bool {
	_, ok := m[key]
	return ok
}

func saveHostData(fileName string, data interface{}, key []byte) error {
	var buf bytes.Buffer
	iv := make([]byte, aes.BlockSize)

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(&buf)
	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return err
	}

	_, err = file.Write(iv)
	if err != nil {
		return err
	}

	stream := cipher.NewCFBEncrypter(block, iv)

	writer := &cipher.StreamWriter{S: stream, W: file}
	_, err = writer.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func loadHostData(fileName string, key []byte, decryptedData interface{}) error {
	encryptedData := make([]byte, 4096)

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Read(encryptedData)
	if err != nil && err != io.EOF {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	iv := encryptedData[:aes.BlockSize]
	stream := cipher.NewCFBDecrypter(block, iv)

	encryptedData = encryptedData[aes.BlockSize:]

	reader := &cipher.StreamReader{S: stream, R: bytes.NewReader(encryptedData)}
	decoder := gob.NewDecoder(reader)
	err = decoder.Decode(decryptedData)
	if err != nil {
		return err
	}

	return nil
}

// Not use
func CreateSampleHostData() {
	hosts := []HostInfo{
		{Name: "Local", Address: "localhost", Port: 10122, Username: "user", Password: "12345"},
		{Name: "Local using key", Address: "localhost", Port: 10222, Username: "user", PrivateKeyFile: "my_private_key.pem"},
	}

	fileName := "hosts.dat"
	// key := []byte("0123456789!#$%^&*()abcdefghijklm") // AES key (32byte = 256bit)

	err := saveHostData(fileName, hosts, key)
	if err != nil {
		fmt.Println("error saving data:", err)
		return
	}

	fmt.Println("data is saved to", fileName)

	// var decodedHosts []HostInfo
	// err = loadHostData(fileName, key, &decodedHosts)
	// if err != nil {
	// 	fmt.Println("Error loading data:", err)
	// 	return
	// }

	// fmt.Println("data is loaded:", decodedHosts)
}

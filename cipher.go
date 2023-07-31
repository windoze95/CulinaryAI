package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"os"
)

type CipherConfig struct {
	SecretKey []byte
}

func NewCipherConfig() *CipherConfig {
	secretKey := os.Getenv(gc.Env.SecretKey)
	if secretKey == "" {
		log.Fatalf("%s must be set", gc.Env.SecretKey)
	}
	return &CipherConfig{
		SecretKey: []byte(secretKey),
	}
}

// Encrypt encrypts the plaintext with the secret key
func Encrypt(config *CipherConfig, plaintext string) (string, error) {
	block, err := aes.NewCipher(config.SecretKey)
	if err != nil {
		return "Error creating block cipher", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "Error reading cipher size", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts the ciphertext with the secret key
func Decrypt(config *CipherConfig, ciphertext string) (string, error) {
	block, err := aes.NewCipher(config.SecretKey)
	if err != nil {
		return "", err
	}

	decodedCiphertext, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	if len(decodedCiphertext) < aes.BlockSize {
		return "", errors.New("ciphertext is too short")
	}

	iv := decodedCiphertext[:aes.BlockSize]
	decodedCiphertext = decodedCiphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(decodedCiphertext, decodedCiphertext)

	return string(decodedCiphertext), nil
}

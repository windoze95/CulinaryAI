package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

type CipherConfig struct {
	EncryptionKey []byte
}

func GetOpenAIKeyCipherConfig() *CipherConfig {
	encryptionKey := os.Getenv(gc.Env.OpenAIKeyEncryptionKey)
	if encryptionKey == "" {
		log.Fatalf("%s must be set", gc.Env.OpenAIKeyEncryptionKey)
	}
	return &CipherConfig{
		EncryptionKey: []byte(encryptionKey),
	}
}

func encryptOpenAIKey(plaintext string) (string, error) {
	return encrypt(GetOpenAIKeyCipherConfig(), plaintext)
}

func decryptOpenAIKey(ciphertext string) (string, error) {
	return decrypt(GetOpenAIKeyCipherConfig(), ciphertext)
}

// Encrypt encrypts the plaintext with the secret key
func encrypt(config *CipherConfig, plaintext string) (string, error) {
	block, err := aes.NewCipher(config.EncryptionKey)
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
func decrypt(config *CipherConfig, ciphertext string) (string, error) {
	block, err := aes.NewCipher(config.EncryptionKey)
	if err != nil {
		return "", err
	}
	fmt.Println("openai decrypt 1")
	decodedCiphertext, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		fmt.Println("openai decrypt err 1")
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

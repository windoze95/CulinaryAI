package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"os"
)

type CipherConfig struct {
	EncryptionKey []byte
}

func GetOpenAIKeyCipherConfig() *CipherConfig {
	encryptionKeyHex := os.Getenv(gc.Env.OpenAIKeyEncryptionKey)
	if encryptionKeyHex == "" {
		log.Fatal("Openai key encryption key must be set")
	}
	encryptionKey, err := hex.DecodeString(encryptionKeyHex)
	if err != nil {
		log.Fatalf("Unable to decode openai key encryption key hex")
	}
	return &CipherConfig{
		EncryptionKey: encryptionKey,
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

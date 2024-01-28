package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

// CipherConfig is the configuration for the cipher.
type CipherConfig struct {
	EncryptionKey []byte
}

// GetOpenAIKeyCipherConfig gets the cipher config for the OpenAI key.
func GetOpenAIKeyCipherConfig(encryptionKeyHex string) (*CipherConfig, error) {
	if encryptionKeyHex == "" {
		return nil, errors.New("openai key encryption key must be set")
	}
	encryptionKey, err := hex.DecodeString(encryptionKeyHex)
	if err != nil {
		return nil, fmt.Errorf("unable to decode openai key encryption key hex: %v", err)
	}
	return &CipherConfig{
		EncryptionKey: encryptionKey,
	}, nil
}

// EncryptOpenAIKey encrypts the OpenAI key with the secret key.
func EncryptOpenAIKey(encryptionKeyHex string, plaintext string) (string, error) {
	config, err := GetOpenAIKeyCipherConfig(encryptionKeyHex)
	if err != nil {
		return "", err
	}
	return encrypt(config, plaintext)
}

// DecryptOpenAIKey decrypts the OpenAI key with the secret key.
func DecryptOpenAIKey(encryptionKeyHex string, ciphertext string) (string, error) {
	config, err := GetOpenAIKeyCipherConfig(encryptionKeyHex)
	if err != nil {
		return "", err
	}
	return decrypt(config, ciphertext)
}

// encrypt encrypts the plaintext with the secret key.
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

// decrypt decrypts the ciphertext with the secret key.
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

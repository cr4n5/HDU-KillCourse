package util

import (
	"crypto/aes"
	"encoding/base64"
)

// AesEncrypt
func AesEncrypt(key, plainText string) (string, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	paddedText := PKCS7Padding([]byte(plainText), block.BlockSize())
	cipherText := make([]byte, len(paddedText))

	for bs, be := 0, block.BlockSize(); bs < len(paddedText); bs, be = bs+block.BlockSize(), be+block.BlockSize() {
		block.Encrypt(cipherText[bs:be], paddedText[bs:be])
	}

	encodedCipherText := base64.StdEncoding.EncodeToString(cipherText)
	return encodedCipherText, nil
}

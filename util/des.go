package util

import (
	"bytes"
	"crypto/des"
	"encoding/base64"
)

// PKCS7Padding pads the plaintext to be a multiple of the block size
func PKCS7Padding(plainText []byte, blockSize int) []byte {
	padding := blockSize - len(plainText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plainText, padText...)
}

// Encrypt encrypts the plaintext using DES algorithm with the given key
func DesEncrypt(key, plainText string) (string, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}

	block, err := des.NewCipher(keyBytes)
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

package util

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"math/big"
)

func RsaEncrypt(publicKey string, data string) (string, error) {
	// 解码公钥
	pubKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return "", err
	}
	// 转化成十六进制
	pubKeyHex := hex.EncodeToString(pubKey)

	// 创建公钥
	pub := new(rsa.PublicKey)
	pub.N = new(big.Int)
	pub.N.SetString(pubKeyHex, 16)
	pub.E = 65537

	// 加密
	cipher, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(data))
	if err != nil {
		return "", err
	}

	// 转换为16进制
	cipherHex := hex.EncodeToString(cipher)
	cipherByte, err := hex.DecodeString(cipherHex)
	if err != nil {
		return "", err
	}
	// 转换为base64
	return base64.StdEncoding.EncodeToString(cipherByte), nil
}

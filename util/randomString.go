package util

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"time"
)

func GenerateRandomString(n int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// GenerateCsrfValue 生成 Csrf-Value
func GenerateCsrfValue(n string) string {
	// Base64 编码
	t := base64.StdEncoding.EncodeToString([]byte(n))
	// 插入 t 本身
	o := t[:len(t)/2] + t + t[len(t)/2:]
	// 计算 MD5 哈希值
	hash := md5.Sum([]byte(o))
	return hex.EncodeToString(hash[:])
}

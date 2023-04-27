package util

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func HmacSha1WithBase64(secret, value string) string {
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(value))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)) //进行base64编码
}

func HmacSha256WithBase64(secret, value string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(value))
	return hex.EncodeToString(mac.Sum(nil))
}

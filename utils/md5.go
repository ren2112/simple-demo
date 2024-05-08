package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	tempStr := h.Sum(nil)
	return hex.EncodeToString(tempStr)
}

func MakePassword(plainpwd, salt string) string {
	return Md5Encode(plainpwd + salt)
}

func ValidPassword(plainpwd, salt, password string) bool {
	return Md5Encode(plainpwd+salt) == password
}

package utils

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword bcrypt 加密密码
func HashPassword(pwd string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(b), err
}

// ParsePassword 校验明文密码与 bcrypt 密文是否一致
func ParsePassword(password, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}

// GenerateSession 生成随机 session id（32 字节，URL 安全 Base64）
func GenerateSession() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic("cannot generate session: " + err.Error())
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

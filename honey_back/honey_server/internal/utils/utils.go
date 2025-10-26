package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/crypto/bcrypt"
)

// GenerateRSA 生成RSA密钥对（2048）
// 程序初次启用时生成（生产部署）
func GenerateRSA() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// HashPassword 密码加密
func HashPassword(pwd string) (string, error) {
	crypted, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(crypted), err
}

// ParsePassword 密码对比
func ParsePassword(password, hashed string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		return false, err
	}
	return true, nil
}

// GenerateJWT Generate JwtToken
func GenerateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": username,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := token.SignedString([]byte("secret"))

	return tokenString, err
}

// ParseJWT 解析JWT
func ParseJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, ok := claims["username"].(string)
		if !ok {
			return "", errors.New("username is not a string")
		}
		return username, nil
	}

	return "", err
}

// GenerateSession 生成随机的session_id（32字节，Base64编码，确保随机性）
func GenerateSession() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic("Cannot generate session: " + err.Error())
	}
	return base64.URLEncoding.EncodeToString(b)
}

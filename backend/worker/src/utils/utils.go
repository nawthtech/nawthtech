package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/cloudflare/cloudflare-go/d1"
)

type D1Database struct {
	Client *d1.Client
	DB     *d1.DB
}

func GetD1Database(dsn string) (*D1Database, error) {
	client, err := d1.NewClient(dsn)
	if err != nil {
		return nil, err
	}
	db := client.DB("")
	return &D1Database{Client: client, DB: db}, nil
}

// JWT
func GenerateJWT(payload map[string]interface{}, secret string, expiration time.Duration) (string, error) {
	claims := jwt.MapClaims{}
	for k, v := range payload {
		claims[k] = v
	}
	claims["exp"] = time.Now().Add(expiration).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenStr, secret string) (bool, map[string]interface{}) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return false, nil
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return true, claims
	}
	return false, nil
}
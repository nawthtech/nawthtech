package utils

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nawthtech/nawthtech/backend/internal/config"
)

// Claims
type CustomClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT generates access and refresh tokens
func GenerateJWT(cfg *config.Config, userID, email, role string) (string, string, error) {
	if cfg == nil {
		return "", "", errors.New("config required")
	}
	secret := cfg.Auth.JWTSecret
	if secret == "" {
		return "", "", errors.New("JWT_SECRET not set")
	}

	now := time.Now()
	accessClaims := CustomClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(cfg.Auth.JWTExpiration)),
			Issuer:    "nawthtech",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	// refresh token (longer expiry)
	refreshClaims := jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(cfg.Auth.RefreshExpiration)),
		Issuer:    "nawthtech",
		Subject:   userID,
	}
	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refreshTokenObj.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// VerifyJWT verifies token and returns claims
func VerifyJWT(cfg *config.Config, tokenStr string) (*CustomClaims, error) {
	if cfg == nil {
		return nil, errors.New("config required")
	}
	secret := cfg.Auth.JWTSecret
	if secret == "" {
		return nil, errors.New("JWT_SECRET not set")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func GetUserIDFromContext(ctx context.Context) string {
	// helper to extract user id from context; used by handlers (you'll implement middleware to set it)
	if ctx == nil {
		return ""
	}
	if v := ctx.Value("user_id"); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

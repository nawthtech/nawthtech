package utils

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JSONResponse تسهل إرسال استجابة JSON
func JSONResponse(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// GenerateJWT ينشئ JWT صالح لمدة محددة (مثال: 24 ساعة)
func GenerateJWT(userID string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", ErrMissingJWTSecret
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateJWT يتحقق من صحة JWT ويستخرج user_id
func ValidateJWT(tokenString string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", ErrMissingJWTSecret
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// التأكد من طريقة التوقيع
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidJWT
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["user_id"].(string); ok {
			return userID, nil
		}
	}

	return "", ErrInvalidJWT
}

// الأخطاء المخصصة
var (
	ErrMissingJWTSecret = &JWTError{"JWT_SECRET is not set in environment"}
	ErrInvalidJWT       = &JWTError{"Invalid JWT token"}
)

// JWTError خطأ مخصص للتعامل مع JWT
type JWTError struct {
	Message string
}

func (e *JWTError) Error() string {
	return e.Message
}
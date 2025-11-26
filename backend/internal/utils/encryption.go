package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

var encryptionKey []byte

func InitEncryption(key string) error {
	if key == "" {
		return nil // لا يوجد مفتاح تشفير، ربما بيئة تطوير
	}

	// تحويل المفتاح إلى 32 بايت (AES-256)
	keyBytes := []byte(key)
	if len(keyBytes) < 32 {
		// تمديد المفتاح إذا كان أقصر من 32 بايت
		paddedKey := make([]byte, 32)
		copy(paddedKey, keyBytes)
		encryptionKey = paddedKey
	} else {
		encryptionKey = keyBytes[:32]
	}
	return nil
}

func Encrypt(text string) (string, error) {
	if encryptionKey == nil {
		return text, nil // لا يوجد تشفير، إرجاع النص كما هو
	}

	plaintext := []byte(text)

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(cryptoText string) (string, error) {
	if encryptionKey == nil {
		return cryptoText, nil // لا يوجد تشفير، إرجاع النص كما هو
	}

	ciphertext, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

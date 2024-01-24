package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword 返回一个密码的哈希值
func HashPassword(password string) (string, error) {
	// bcrypt.DefaultCost = 10
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate hash from password: %w", err)
	}

	return string(bytes), nil
}

func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

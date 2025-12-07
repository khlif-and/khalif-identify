package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

)

func HashPassword(p string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(p), 14)
	return string(bytes), err
}

func CheckPasswordHash(p, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(p))
	return err == nil
}

func GenerateToken(userID string, role, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
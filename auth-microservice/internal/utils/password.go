package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ValidatePasswordResetToken(token string) (string, error) {
	userID, err := ExtractUserIDFromJWT(token)
	if err != nil {
		return "", errors.New("invalid or expired password reset token")
	}
	return userID, nil
}

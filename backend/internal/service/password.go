package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	// DefaultBcryptCost niveau de complexité recommandé pour la production
	DefaultBcryptCost = 12
)

// HashPassword génère un hash sécurisé bcrypt à partir du mot de passe en clair.
func HashPassword(password string) (string, error) {
	if len(password) < 6 {
		return "", errors.New("le mot de passe doit comporter au moins 6 caractères")
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultBcryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash vérifie si le mot de passe en clair correspond au hash stocké.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"autoschool/backend/internal/models"
)

// GenerateAccessToken émet un JWT signé contenant l'identité utilisateur et les identifiants de multi-tenancy.
func GenerateAccessToken(
	user *models.User,
	schoolID string,
	agenceID *string,
	secret string,
	expirationHours int,
) (string, int64, error) {
	expDuration := time.Duration(expirationHours) * time.Hour
	expiresAt := time.Now().Add(expDuration)

	claims := &models.JWTClaims{
		UserID:      user.ID,
		AutoEcoleID: schoolID,
		AgenceID:    agenceID,
		Role:        user.Role,
		Telephone:   user.Telephone,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			Issuer:    "autoschool-backend",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", 0, fmt.Errorf("erreur signature JWT: %w", err)
	}

	return tokenString, int64(expDuration.Seconds()), nil
}

// ValidateAccessToken parse et valide la signature et l'expiration d'un JWT.
func ValidateAccessToken(tokenString string, secret string) (*models.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("méthode de signature inattendue: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token invalide ou corrompu")
	}

	return claims, nil
}

// GenerateRefreshToken génère un token aléatoire cryptographique de 32 octets (hex 64 chars)
// et retourne la valeur en clair (pour le client) ainsi que son hash SHA-256 (pour la base de données).
func GenerateRefreshToken() (plainToken string, tokenHash string, err error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", fmt.Errorf("impossible de générer un token aléatoire: %w", err)
	}

	plainToken = hex.EncodeToString(bytes)
	tokenHash = HashRefreshToken(plainToken)
	return plainToken, tokenHash, nil
}

// HashRefreshToken applique SHA-256 sur le token en clair pour stockage sécurisé en base.
func HashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthSession représente une session active en base de données PostgreSQL.
type AuthSession struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	AutoEcoleID string     `json:"auto_ecole_id"`
	TokenHash   string     `json:"token_hash"`
	IPAddress   *string    `json:"ip_address,omitempty"`
	UserAgent   *string    `json:"user_agent,omitempty"`
	ExpiresAt   time.Time  `json:"expires_at"`
	RevokedAt   *time.Time `json:"revoked_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// JWTClaims porte les données d'identité et de multi-tenancy encodées dans le token JWT.
type JWTClaims struct {
	UserID      string  `json:"user_id"`
	AutoEcoleID string  `json:"auto_ecole_id"`
	AgenceID    *string `json:"agence_id,omitempty"`
	Role        string  `json:"role"`
	Telephone   string  `json:"telephone"`
	jwt.RegisteredClaims
}

// AuthResponse payload renvoyé après une connexion ou un renouvellement réussi.
type AuthResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    int64         `json:"expires_in"` // Durée en secondes
	User         *UserResponse `json:"user"`
	School       *AutoEcole    `json:"school"`
	Agence       *Agence       `json:"agence,omitempty"`
}

package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"autoschool/backend/internal/database"
	"autoschool/backend/internal/models"
)

// SessionRepository gère la persistance des sessions et refresh tokens.
type SessionRepository struct {
	db *database.DB
}

func NewSessionRepository(db *database.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Create enregistre une nouvelle session en base.
func (r *SessionRepository) Create(ctx context.Context, session *models.AuthSession) error {
	query := `
		INSERT INTO auth_sessions (
			id, user_id, auto_ecole_id, token_hash, ip_address, user_agent, expires_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		) RETURNING created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query,
		session.ID, session.UserID, session.AutoEcoleID, session.TokenHash,
		session.IPAddress, session.UserAgent, session.ExpiresAt,
	).Scan(&session.CreatedAt, &session.UpdatedAt)

	if err != nil {
		return fmt.Errorf("erreur création session auth: %w", err)
	}
	return nil
}

// GetValidByTokenHash recherche une session active non expirée et non révoquée par le hash de son token.
func (r *SessionRepository) GetValidByTokenHash(ctx context.Context, tokenHash string) (*models.AuthSession, error) {
	query := `
		SELECT id, user_id, auto_ecole_id, token_hash, ip_address, user_agent,
		       expires_at, revoked_at, created_at, updated_at
		FROM auth_sessions
		WHERE token_hash = $1 
		  AND revoked_at IS NULL 
		  AND expires_at > NOW()
		LIMIT 1
	`
	var s models.AuthSession
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&s.ID, &s.UserID, &s.AutoEcoleID, &s.TokenHash, &s.IPAddress, &s.UserAgent,
		&s.ExpiresAt, &s.RevokedAt, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erreur recherche session valide: %w", err)
	}
	return &s, nil
}

// Revoke révoque une session spécifique (par exemple lors d'un logout).
func (r *SessionRepository) Revoke(ctx context.Context, tokenHash string) error {
	query := `
		UPDATE auth_sessions
		SET revoked_at = NOW(), updated_at = NOW()
		WHERE token_hash = $1 AND revoked_at IS NULL
	`
	_, err := r.db.ExecContext(ctx, query, tokenHash)
	if err != nil {
		return fmt.Errorf("erreur révocation session: %w", err)
	}
	return nil
}

// RevokeAllForUser révoque toutes les sessions actives d'un utilisateur (déconnexion globale).
func (r *SessionRepository) RevokeAllForUser(ctx context.Context, userID string) error {
	query := `
		UPDATE auth_sessions
		SET revoked_at = NOW(), updated_at = NOW()
		WHERE user_id = $1 AND revoked_at IS NULL
	`
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("erreur révocation de toutes les sessions: %w", err)
	}
	return nil
}

// CleanExpiredSessions supprime les sessions expirées de plus de 7 jours.
func (r *SessionRepository) CleanExpiredSessions(ctx context.Context) (int64, error) {
	query := `
		DELETE FROM auth_sessions
		WHERE expires_at < $1
	`
	res, err := r.db.ExecContext(ctx, query, time.Now().AddDate(0, 0, -7))
	if err != nil {
		return 0, fmt.Errorf("erreur nettoyage sessions expirées: %w", err)
	}
	return res.RowsAffected()
}

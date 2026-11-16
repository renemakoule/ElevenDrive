package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// DB enveloppe la connexion PostgreSQL.
type DB struct {
	*sql.DB
}

// Connect établit et configure le pool de connexions vers PostgreSQL (Neon ou local).
func Connect(databaseURL string) (*DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'initialisation du driver postgres: %w", err)
	}

	// Configuration du pool de connexions
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	// Vérification de la disponibilité de la base avec un timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("impossible d'établir la connexion avec la base de données: %w", err)
	}

	return &DB{db}, nil
}

// WithTenantTx exécute une fonction transactionnelle en injectant l'isolation Row-Level Security (RLS)
func (db *DB) WithTenantTx(ctx context.Context, autoEcoleID string, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("erreur démarrage transaction: %w", err)
	}
	defer tx.Rollback()

	if autoEcoleID != "" {
		if _, err := uuid.Parse(autoEcoleID); err != nil {
			return fmt.Errorf("identifiant auto_ecole_id invalide pour RLS: %w", err)
		}
		// Injection paramétrée sécurisée (anti-injection SQL) pour les politiques RLS PostgreSQL
		if _, err := tx.ExecContext(ctx, "SELECT set_config('app.current_auto_ecole_id', $1, true)", autoEcoleID); err != nil {
			return fmt.Errorf("erreur configuration RLS: %w", err)
		}
	}

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}

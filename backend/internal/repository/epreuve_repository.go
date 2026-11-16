package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"autoschool/backend/internal/database"
	"autoschool/backend/internal/models"
	"github.com/google/uuid"
)

type EpreuveRepository struct {
	db *database.DB
}

func NewEpreuveRepository(db *database.DB) *EpreuveRepository {
	return &EpreuveRepository{db: db}
}

// EnsureTableExists s'assure que la table epreuves_pdf existe dans PostgreSQL sans nécessiter d'opération manuelle
func (r *EpreuveRepository) EnsureTableExists(ctx context.Context) error {
	ddl := `
	CREATE TABLE IF NOT EXISTS epreuves_pdf (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		auto_ecole_id UUID NOT NULL REFERENCES auto_ecoles(id) ON DELETE CASCADE,
		titre VARCHAR(150) NOT NULL,
		categorie_permis VARCHAR(10) NOT NULL,
		description TEXT,
		annee INT NOT NULL DEFAULT 2024,
		session VARCHAR(100),
		duree_minutes INT DEFAULT 45,
		nombre_questions INT DEFAULT 40,
		avec_corrige BOOLEAN DEFAULT TRUE,
		fichier_nom VARCHAR(150) NOT NULL,
		fichier_taille VARCHAR(30),
		fichier_data TEXT NOT NULL,
		ajoute_par VARCHAR(120),
		telechargements_count INT DEFAULT 0,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_epreuves_school ON epreuves_pdf(auto_ecole_id);
	CREATE INDEX IF NOT EXISTS idx_epreuves_categorie ON epreuves_pdf(auto_ecole_id, categorie_permis);
	`
	_, err := r.db.ExecContext(ctx, ddl)
	return err
}

// GetEpreuves récupère les épreuves réelles en base de données pour une auto-école, filtrables par catégorie de permis
func (r *EpreuveRepository) GetEpreuves(ctx context.Context, autoEcoleID string, categorie string) ([]models.EpreuvePdf, error) {
	_ = r.EnsureTableExists(ctx)

	baseQuery := `
		SELECT id, auto_ecole_id, titre, categorie_permis, description,
		       annee, session, duree_minutes, nombre_questions, avec_corrige,
		       fichier_nom, fichier_taille, fichier_data, ajoute_par,
		       telechargements_count, created_at, updated_at
		FROM epreuves_pdf
		WHERE auto_ecole_id = $1
	`
	args := []interface{}{autoEcoleID}

	if categorie != "" && strings.ToUpper(categorie) != "ALL" {
		baseQuery += ` AND (UPPER(categorie_permis) = UPPER($2) OR UPPER(categorie_permis) = 'TOUTES')`
		args = append(args, categorie)
	}

	baseQuery += ` ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération des épreuves en base: %w", err)
	}
	defer rows.Close()

	var list []models.EpreuvePdf
	for rows.Next() {
		var ep models.EpreuvePdf
		var desc, sess, taille, ajoutePar sql.NullString

		if err := rows.Scan(
			&ep.ID, &ep.AutoEcoleID, &ep.Titre, &ep.CategoriePermis, &desc,
			&ep.Annee, &sess, &ep.DureeMinutes, &ep.NombreQuestions, &ep.AvecCorrige,
			&ep.FichierNom, &taille, &ep.FichierData, &ajoutePar,
			&ep.TelechargementsCount, &ep.CreatedAt, &ep.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if desc.Valid {
			ep.Description = &desc.String
		}
		if sess.Valid {
			ep.Session = &sess.String
		}
		if taille.Valid {
			ep.FichierTaille = &taille.String
		}
		if ajoutePar.Valid {
			ep.AjoutePar = &ajoutePar.String
		}

		list = append(list, ep)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if list == nil {
		list = []models.EpreuvePdf{}
	}
	return list, nil
}

// CreerEpreuve insère une épreuve PDF authentique en base de données PostgreSQL
func (r *EpreuveRepository) CreerEpreuve(ctx context.Context, autoEcoleID string, input *models.CreerEpreuveInput) (*models.EpreuvePdf, error) {
	_ = r.EnsureTableExists(ctx)

	epID := uuid.New().String()
	now := time.Now().UTC()

	query := `
		INSERT INTO epreuves_pdf (
			id, auto_ecole_id, titre, categorie_permis, description,
			annee, session, duree_minutes, nombre_questions, avec_corrige,
			fichier_nom, fichier_taille, fichier_data, ajoute_par,
			telechargements_count, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13, $14,
			0, $15, $15
		)
		RETURNING id, created_at, updated_at
	`

	ep := &models.EpreuvePdf{
		ID:                   epID,
		AutoEcoleID:          autoEcoleID,
		Titre:                input.Titre,
		CategoriePermis:      input.CategoriePermis,
		Description:          input.Description,
		Annee:                input.Annee,
		Session:              input.Session,
		DureeMinutes:         input.DureeMinutes,
		NombreQuestions:      input.NombreQuestions,
		AvecCorrige:          input.AvecCorrige,
		FichierNom:           input.FichierNom,
		FichierTaille:        input.FichierTaille,
		FichierData:          input.FichierData,
		AjoutePar:            input.AjoutePar,
		TelechargementsCount: 0,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	err := r.db.QueryRowContext(
		ctx, query,
		ep.ID, ep.AutoEcoleID, ep.Titre, ep.CategoriePermis, ep.Description,
		ep.Annee, ep.Session, ep.DureeMinutes, ep.NombreQuestions, ep.AvecCorrige,
		ep.FichierNom, ep.FichierTaille, ep.FichierData, ep.AjoutePar,
		now,
	).Scan(&ep.ID, &ep.CreatedAt, &ep.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("erreur insertion épreuve en base: %w", err)
	}

	return ep, nil
}

// SupprimerEpreuve retire définitivement une épreuve de la base de données
func (r *EpreuveRepository) SupprimerEpreuve(ctx context.Context, autoEcoleID, epreuveID string) error {
	query := `DELETE FROM epreuves_pdf WHERE id = $1 AND auto_ecole_id = $2`
	res, err := r.db.ExecContext(ctx, query, epreuveID, autoEcoleID)
	if err != nil {
		return fmt.Errorf("erreur suppression épreuve: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("épreuve introuvable ou non autorisée")
	}
	return nil
}

// IncrementerTelechargement met à jour le compteur d'accès de l'épreuve
func (r *EpreuveRepository) IncrementerTelechargement(ctx context.Context, autoEcoleID, epreuveID string) error {
	query := `
		UPDATE epreuves_pdf
		SET telechargements_count = telechargements_count + 1, updated_at = NOW()
		WHERE id = $1 AND auto_ecole_id = $2
	`
	_, err := r.db.ExecContext(ctx, query, epreuveID, autoEcoleID)
	return err
}

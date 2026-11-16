package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"autoschool/backend/internal/database"
	"autoschool/backend/internal/models"

	"github.com/google/uuid"
)

type FlotteRepository struct {
	db *database.DB
}

func NewFlotteRepository(db *database.DB) *FlotteRepository {
	return &FlotteRepository{db: db}
}

// GetVehicules liste l'ensemble des véhicules appartenant aux agences de l'auto-école.
func (r *FlotteRepository) GetVehicules(ctx context.Context, autoEcoleID string) ([]models.Vehicule, error) {
	var res []models.Vehicule
	err := database.Retry(ctx, func() error {
		return database.Safe(func() error {
			var e error
			res, e = r.getVehiculesOnce(ctx, autoEcoleID)
			return e
		})
	})
	return res, err
}

func (r *FlotteRepository) getVehiculesOnce(ctx context.Context, autoEcoleID string) ([]models.Vehicule, error) {
	query := `
		SELECT v.id, v.agence_id, v.immatriculation, v.marque_modele, v.type_boite,
		       v.kilometrage_actuel,
		       to_char(v.date_expiration_assurance, 'YYYY-MM-DD'),
		       to_char(v.date_visite_technique, 'YYYY-MM-DD'),
		       v.statut, v.created_at, v.updated_at
		FROM vehicules v
		JOIN agences a ON a.id = v.agence_id
		WHERE a.auto_ecole_id = $1
		ORDER BY v.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, autoEcoleID)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération des véhicules: %w", err)
	}
	defer rows.Close()

	var list []models.Vehicule
	for rows.Next() {
		var v models.Vehicule
		if err := rows.Scan(
			&v.ID, &v.AgenceID, &v.Immatriculation, &v.MarqueModele, &v.TypeBoite,
			&v.KilometrageActuel, &v.DateExpirationAssurance, &v.DateVisiteTechnique,
			&v.Statut, &v.CreatedAt, &v.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, v)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if list == nil {
		list = []models.Vehicule{}
	}
	return list, nil
}

// CreerVehicule enregistre un véhicule dans la base de données.
func (r *FlotteRepository) CreerVehicule(ctx context.Context, autoEcoleID, agenceID string, input *models.CreerVehiculeInput) (*models.Vehicule, error) {
	targetAgenceID := agenceID
	if input.AgenceID != nil && strings.TrimSpace(*input.AgenceID) != "" {
		candidateID := strings.TrimSpace(*input.AgenceID)
		var validID string
		err := r.db.QueryRowContext(ctx, `SELECT id FROM agences WHERE id = $1 AND auto_ecole_id = $2`, candidateID, autoEcoleID).Scan(&validID)
		if err != nil {
			return nil, errors.New("l'agence spécifiée est introuvable ou n'appartient pas à votre auto-école")
		}
		targetAgenceID = validID
	}
	if targetAgenceID == "" {
		err := r.db.QueryRowContext(ctx, `SELECT id FROM agences WHERE auto_ecole_id = $1 ORDER BY est_siege DESC, created_at ASC LIMIT 1`, autoEcoleID).Scan(&targetAgenceID)
		if err != nil {
			return nil, errors.New("aucune agence trouvée pour associer le véhicule")
		}
	}

	typeBoite := "MANUELLE"
	if strings.ToUpper(input.TypeBoite) == "AUTOMATIQUE" {
		typeBoite = "AUTOMATIQUE"
	}

	statut := "DISPONIBLE"
	if input.Statut != "" {
		statut = strings.ToUpper(input.Statut)
	}

	vehID := uuid.New().String()
	query := `
		INSERT INTO vehicules (
			id, agence_id, immatriculation, marque_modele, type_boite,
			kilometrage_actuel, date_expiration_assurance, date_visite_technique,
			statut, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7::date, $8::date, $9, NOW(), NOW())
		RETURNING id, agence_id, immatriculation, marque_modele, type_boite,
		          kilometrage_actuel,
		          to_char(date_expiration_assurance, 'YYYY-MM-DD'),
		          to_char(date_visite_technique, 'YYYY-MM-DD'),
		          statut, created_at, updated_at
	`
	var v models.Vehicule
	err := r.db.QueryRowContext(
		ctx, query,
		vehID, targetAgenceID, strings.ToUpper(strings.TrimSpace(input.Immatriculation)),
		strings.TrimSpace(input.MarqueModele), typeBoite, input.KilometrageActuel,
		input.DateExpirationAssurance, input.DateVisiteTechnique, statut,
	).Scan(
		&v.ID, &v.AgenceID, &v.Immatriculation, &v.MarqueModele, &v.TypeBoite,
		&v.KilometrageActuel, &v.DateExpirationAssurance, &v.DateVisiteTechnique,
		&v.Statut, &v.CreatedAt, &v.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("erreur création véhicule: %w", err)
	}
	return &v, nil
}

// UpdateVehiculeStatut modifie le statut d'un véhicule (ex: DISPONIBLE, EN_ENTRETIEN, EN_PANNE).
func (r *FlotteRepository) UpdateVehiculeStatut(ctx context.Context, autoEcoleID, vehiculeID, statut string) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE vehicules v
		SET statut = $1, updated_at = NOW()
		FROM agences a
		WHERE v.id = $2 AND v.agence_id = a.id AND a.auto_ecole_id = $3
	`, statut, vehiculeID, autoEcoleID)
	if err != nil {
		return fmt.Errorf("erreur mise à jour statut véhicule: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("véhicule introuvable")
	}
	return nil
}

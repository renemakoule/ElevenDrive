package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"autoschool/backend/internal/database"
	"autoschool/backend/internal/models"
)

type CommunicationRepository struct {
	db *database.DB
}

func NewCommunicationRepository(db *database.DB) *CommunicationRepository {
	return &CommunicationRepository{db: db}
}

// -------------------------------------------------------------
// ALERTES
// -------------------------------------------------------------

// GetAllAlerts récupère toutes les alertes (Super-Admin)
func (r *CommunicationRepository) GetAllAlerts(ctx context.Context) ([]models.NetworkAlert, error) {
	query := `
		SELECT 
			a.id, a.titre, a.message, a.severite, a.cible_type, a.cible_auto_ecole_id,
			COALESCE(e.nom_commercial, 'Toutes les auto-écoles') as cible_school_name,
			a.bloquante, a.est_active, a.date_envoi, a.created_at
		FROM reseau_alertes a
		LEFT JOIN auto_ecoles e ON e.id = a.cible_auto_ecole_id
		ORDER BY a.date_envoi DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération alertes: %w", err)
	}
	defer rows.Close()

	var list []models.NetworkAlert
	for rows.Next() {
		var a models.NetworkAlert
		var cibleSchoolID sql.NullString
		if err := rows.Scan(
			&a.ID, &a.Titre, &a.Message, &a.Severite, &a.CibleType, &cibleSchoolID,
			&a.CibleSchoolName, &a.Bloquante, &a.EstActive, &a.DateEnvoi, &a.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("erreur scan alerte: %w", err)
		}
		if cibleSchoolID.Valid {
			val := cibleSchoolID.String
			a.CibleSchoolID = &val
		}
		list = append(list, a)
	}
	if list == nil {
		list = []models.NetworkAlert{}
	}
	return list, nil
}

// GetActiveAlertsForSchool récupère les alertes actives destinées à une école donnée ou globales
func (r *CommunicationRepository) GetActiveAlertsForSchool(ctx context.Context, schoolID string) ([]models.NetworkAlert, error) {
	query := `
		SELECT 
			a.id, a.titre, a.message, a.severite, a.cible_type, a.cible_auto_ecole_id,
			COALESCE(e.nom_commercial, 'Toutes les auto-écoles') as cible_school_name,
			a.bloquante, a.est_active, a.date_envoi, a.created_at
		FROM reseau_alertes a
		LEFT JOIN auto_ecoles e ON e.id = a.cible_auto_ecole_id
		WHERE a.est_active = TRUE AND (a.cible_type = 'all' OR a.cible_auto_ecole_id = $1)
		ORDER BY a.date_envoi DESC
	`

	rows, err := r.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération alertes école: %w", err)
	}
	defer rows.Close()

	var list []models.NetworkAlert
	for rows.Next() {
		var a models.NetworkAlert
		var cibleSchoolID sql.NullString
		if err := rows.Scan(
			&a.ID, &a.Titre, &a.Message, &a.Severite, &a.CibleType, &cibleSchoolID,
			&a.CibleSchoolName, &a.Bloquante, &a.EstActive, &a.DateEnvoi, &a.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("erreur scan alerte école: %w", err)
		}
		if cibleSchoolID.Valid {
			val := cibleSchoolID.String
			a.CibleSchoolID = &val
		}
		list = append(list, a)
	}
	if list == nil {
		list = []models.NetworkAlert{}
	}
	return list, nil
}

// CreateAlert diffuse une nouvelle alerte
func (r *CommunicationRepository) CreateAlert(ctx context.Context, input *models.CreateAlertInput) (*models.NetworkAlert, error) {
	var cibleID *string
	if input.CibleType == "single" && input.CibleSchoolID != nil && strings.TrimSpace(*input.CibleSchoolID) != "" {
		cibleID = input.CibleSchoolID
	}

	query := `
		INSERT INTO reseau_alertes (
			titre, message, severite, cible_type, cible_auto_ecole_id, bloquante, est_active, date_envoi
		) VALUES ($1, $2, $3, $4, $5, $6, TRUE, NOW())
		RETURNING id, titre, message, severite, cible_type, cible_auto_ecole_id, bloquante, est_active, date_envoi, created_at
	`

	var a models.NetworkAlert
	var scannedCibleID sql.NullString
	err := r.db.QueryRowContext(
		ctx, query,
		input.Titre, input.Message, input.Severite, input.CibleType, cibleID, input.Bloquante,
	).Scan(
		&a.ID, &a.Titre, &a.Message, &a.Severite, &a.CibleType, &scannedCibleID,
		&a.Bloquante, &a.EstActive, &a.DateEnvoi, &a.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("erreur insertion alerte: %w", err)
	}
	if scannedCibleID.Valid {
		val := scannedCibleID.String
		a.CibleSchoolID = &val
	}

	// Résoudre le nom d'école
	if a.CibleSchoolID != nil {
		_ = r.db.QueryRowContext(ctx, `SELECT nom_commercial FROM auto_ecoles WHERE id = $1`, *a.CibleSchoolID).Scan(&a.CibleSchoolName)
	} else {
		a.CibleSchoolName = "Toutes les auto-écoles"
	}

	return &a, nil
}

// ToggleAlertStatus active ou révoque une alerte
func (r *CommunicationRepository) ToggleAlertStatus(ctx context.Context, id string) (bool, error) {
	var nouvStatut bool
	query := `
		UPDATE reseau_alertes
		SET est_active = NOT est_active
		WHERE id = $1
		RETURNING est_active
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&nouvStatut)
	if err != nil {
		return false, fmt.Errorf("erreur bascule statut alerte: %w", err)
	}
	return nouvStatut, nil
}

// DeleteAlert supprime définitivement une alerte
func (r *CommunicationRepository) DeleteAlert(ctx context.Context, id string) error {
	query := `DELETE FROM reseau_alertes WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("erreur suppression alerte: %w", err)
	}
	return nil
}

// -------------------------------------------------------------
// ANNONCES & ÉVÉNEMENTS
// -------------------------------------------------------------

// GetAllAnnouncements récupère toutes les annonces (Super-Admin)
func (r *CommunicationRepository) GetAllAnnouncements(ctx context.Context) ([]models.NetworkAnnouncement, error) {
	query := `
		SELECT 
			an.id, an.titre, an.contenu, an.type_format, an.image_url, an.categorie,
			an.cible_type, an.cible_auto_ecole_id,
			COALESCE(e.nom_commercial, 'Toutes les auto-écoles') as cible_school_name,
			an.date_evenement, an.lien_action, an.date_publication, an.created_at
		FROM reseau_annonces an
		LEFT JOIN auto_ecoles e ON e.id = an.cible_auto_ecole_id
		ORDER BY an.date_publication DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération annonces: %w", err)
	}
	defer rows.Close()

	var list []models.NetworkAnnouncement
	for rows.Next() {
		var a models.NetworkAnnouncement
		var imageURL, cibleSchoolID, dateEvenement, lienAction sql.NullString
		if err := rows.Scan(
			&a.ID, &a.Titre, &a.Contenu, &a.TypeFormat, &imageURL, &a.Categorie,
			&a.CibleType, &cibleSchoolID, &a.CibleSchoolName, &dateEvenement, &lienAction,
			&a.DatePublication, &a.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("erreur scan annonce: %w", err)
		}
		if imageURL.Valid {
			val := imageURL.String
			a.ImageURL = &val
		}
		if cibleSchoolID.Valid {
			val := cibleSchoolID.String
			a.CibleSchoolID = &val
		}
		if dateEvenement.Valid {
			val := dateEvenement.String
			a.DateEvenement = &val
		}
		if lienAction.Valid {
			val := lienAction.String
			a.LienAction = &val
		}
		list = append(list, a)
	}
	if list == nil {
		list = []models.NetworkAnnouncement{}
	}
	return list, nil
}

// GetAnnouncementsForSchool récupère les annonces destinées à une auto-école ou globales
func (r *CommunicationRepository) GetAnnouncementsForSchool(ctx context.Context, schoolID string) ([]models.NetworkAnnouncement, error) {
	query := `
		SELECT 
			an.id, an.titre, an.contenu, an.type_format, an.image_url, an.categorie,
			an.cible_type, an.cible_auto_ecole_id,
			COALESCE(e.nom_commercial, 'Toutes les auto-écoles') as cible_school_name,
			an.date_evenement, an.lien_action, an.date_publication, an.created_at
		FROM reseau_annonces an
		LEFT JOIN auto_ecoles e ON e.id = an.cible_auto_ecole_id
		WHERE an.cible_type = 'all' OR an.cible_auto_ecole_id = $1
		ORDER BY an.date_publication DESC
	`

	rows, err := r.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération annonces école: %w", err)
	}
	defer rows.Close()

	var list []models.NetworkAnnouncement
	for rows.Next() {
		var a models.NetworkAnnouncement
		var imageURL, cibleSchoolID, dateEvenement, lienAction sql.NullString
		if err := rows.Scan(
			&a.ID, &a.Titre, &a.Contenu, &a.TypeFormat, &imageURL, &a.Categorie,
			&a.CibleType, &cibleSchoolID, &a.CibleSchoolName, &dateEvenement, &lienAction,
			&a.DatePublication, &a.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("erreur scan annonce école: %w", err)
		}
		if imageURL.Valid {
			val := imageURL.String
			a.ImageURL = &val
		}
		if cibleSchoolID.Valid {
			val := cibleSchoolID.String
			a.CibleSchoolID = &val
		}
		if dateEvenement.Valid {
			val := dateEvenement.String
			a.DateEvenement = &val
		}
		if lienAction.Valid {
			val := lienAction.String
			a.LienAction = &val
		}
		list = append(list, a)
	}
	if list == nil {
		list = []models.NetworkAnnouncement{}
	}
	return list, nil
}

// CreateAnnouncement publie une nouvelle annonce
func (r *CommunicationRepository) CreateAnnouncement(ctx context.Context, input *models.CreateAnnouncementInput) (*models.NetworkAnnouncement, error) {
	var cibleID *string
	if input.CibleType == "single" && input.CibleSchoolID != nil && strings.TrimSpace(*input.CibleSchoolID) != "" {
		cibleID = input.CibleSchoolID
	}

	query := `
		INSERT INTO reseau_annonces (
			titre, contenu, type_format, image_url, categorie, cible_type, cible_auto_ecole_id,
			date_evenement, lien_action, date_publication
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		RETURNING id, titre, contenu, type_format, image_url, categorie, cible_type, cible_auto_ecole_id,
		          date_evenement, lien_action, date_publication, created_at
	`

	var a models.NetworkAnnouncement
	var scannedImage, scannedCibleID, scannedDateEvt, scannedLien sql.NullString
	err := r.db.QueryRowContext(
		ctx, query,
		input.Titre, input.Contenu, input.TypeFormat, input.ImageURL, input.Categorie,
		input.CibleType, cibleID, input.DateEvenement, input.LienAction,
	).Scan(
		&a.ID, &a.Titre, &a.Contenu, &a.TypeFormat, &scannedImage, &a.Categorie,
		&a.CibleType, &scannedCibleID, &scannedDateEvt, &scannedLien, &a.DatePublication, &a.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("erreur insertion annonce: %w", err)
	}

	if scannedImage.Valid {
		v := scannedImage.String
		a.ImageURL = &v
	}
	if scannedCibleID.Valid {
		v := scannedCibleID.String
		a.CibleSchoolID = &v
	}
	if scannedDateEvt.Valid {
		v := scannedDateEvt.String
		a.DateEvenement = &v
	}
	if scannedLien.Valid {
		v := scannedLien.String
		a.LienAction = &v
	}

	if a.CibleSchoolID != nil {
		_ = r.db.QueryRowContext(ctx, `SELECT nom_commercial FROM auto_ecoles WHERE id = $1`, *a.CibleSchoolID).Scan(&a.CibleSchoolName)
	} else {
		a.CibleSchoolName = "Toutes les auto-écoles"
	}

	return &a, nil
}

// DeleteAnnouncement supprime une annonce
func (r *CommunicationRepository) DeleteAnnouncement(ctx context.Context, id string) error {
	query := `DELETE FROM reseau_annonces WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("erreur suppression annonce: %w", err)
	}
	return nil
}

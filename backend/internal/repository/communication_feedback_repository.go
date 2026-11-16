package repository

import (
	"context"
	"database/sql"
	"fmt"

	"autoschool/backend/internal/models"
)

// -------------------------------------------------------------
// AVIS & RETOURS D'EXPÉRIENCE
// -------------------------------------------------------------

// GetAllFeedbacks récupère tous les avis avec filtres optionnels
func (r *CommunicationRepository) GetAllFeedbacks(ctx context.Context, starsFilter int) ([]models.NetworkFeedback, error) {
	query := `
		SELECT 
			fb.id, fb.auto_ecole_id, COALESCE(e.nom_commercial, 'Auto-École Partenaire') as auto_ecole_nom,
			fb.auteur, fb.role, fb.etoiles, fb.categorie, fb.commentaire, fb.statut,
			fb.reponse_admin, fb.created_at
		FROM reseau_avis fb
		LEFT JOIN auto_ecoles e ON e.id = fb.auto_ecole_id
	`
	var args []interface{}
	if starsFilter >= 1 && starsFilter <= 5 {
		query += ` WHERE fb.etoiles = $1`
		args = append(args, starsFilter)
	}
	query += ` ORDER BY fb.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération avis: %w", err)
	}
	defer rows.Close()

	var list []models.NetworkFeedback
	for rows.Next() {
		var f models.NetworkFeedback
		var autoEcoleID, reponseAdmin sql.NullString
		if err := rows.Scan(
			&f.ID, &autoEcoleID, &f.AutoEcoleNom, &f.Auteur, &f.Role, &f.Etoiles,
			&f.Categorie, &f.Commentaire, &f.Statut, &reponseAdmin, &f.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("erreur scan avis: %w", err)
		}
		if autoEcoleID.Valid {
			v := autoEcoleID.String
			f.AutoEcoleID = &v
		}
		if reponseAdmin.Valid {
			v := reponseAdmin.String
			f.ReponseAdmin = &v
		}
		list = append(list, f)
	}
	if list == nil {
		list = []models.NetworkFeedback{}
	}
	return list, nil
}

// CreateFeedback enregistre un avis
func (r *CommunicationRepository) CreateFeedback(ctx context.Context, input *models.CreateFeedbackInput) (*models.NetworkFeedback, error) {
	query := `
		INSERT INTO reseau_avis (
			auto_ecole_id, auteur, role, etoiles, categorie, commentaire, statut, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, 'nouveau', NOW())
		RETURNING id, auto_ecole_id, auteur, role, etoiles, categorie, commentaire, statut, reponse_admin, created_at
	`

	var f models.NetworkFeedback
	var autoEcoleID, reponseAdmin sql.NullString
	err := r.db.QueryRowContext(
		ctx, query,
		input.AutoEcoleID, input.Auteur, input.Role, input.Etoiles, input.Categorie, input.Commentaire,
	).Scan(
		&f.ID, &autoEcoleID, &f.Auteur, &f.Role, &f.Etoiles,
		&f.Categorie, &f.Commentaire, &f.Statut, &reponseAdmin, &f.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("erreur insertion avis: %w", err)
	}
	if autoEcoleID.Valid {
		v := autoEcoleID.String
		f.AutoEcoleID = &v
		_ = r.db.QueryRowContext(ctx, `SELECT nom_commercial FROM auto_ecoles WHERE id = $1`, v).Scan(&f.AutoEcoleNom)
	}
	if f.AutoEcoleNom == "" {
		f.AutoEcoleNom = "Auto-École Pilote"
	}

	return &f, nil
}

// ReplyFeedback enregistre la réponse officielle du Super-Admin
func (r *CommunicationRepository) ReplyFeedback(ctx context.Context, id, reponse, statut string) error {
	if statut == "" {
		statut = "resolu"
	}
	query := `
		UPDATE reseau_avis
		SET reponse_admin = $1, statut = $2
		WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, query, reponse, statut, id)
	if err != nil {
		return fmt.Errorf("erreur réponse avis: %w", err)
	}
	return nil
}

// UpdateFeedbackStatus met à jour le statut d'un avis
func (r *CommunicationRepository) UpdateFeedbackStatus(ctx context.Context, id, statut string) error {
	query := `
		UPDATE reseau_avis
		SET statut = $1
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, statut, id)
	if err != nil {
		return fmt.Errorf("erreur statut avis: %w", err)
	}
	return nil
}

// DeleteFeedback supprime un avis
func (r *CommunicationRepository) DeleteFeedback(ctx context.Context, id string) error {
	query := `DELETE FROM reseau_avis WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("erreur suppression avis: %w", err)
	}
	return nil
}

// -------------------------------------------------------------
// SUPPORT TICKETS
// -------------------------------------------------------------

// CreateSupportTicket soumet une demande de support
func (r *CommunicationRepository) CreateSupportTicket(ctx context.Context, input *models.CreateSupportTicketInput) (*models.SupportTicket, error) {
	query := `
		INSERT INTO reseau_support_tickets (
			auto_ecole_id, nom, email, sujet, priorite, message, statut, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, 'ouvert', NOW())
		RETURNING id, auto_ecole_id, nom, email, sujet, priorite, message, statut, created_at
	`

	var t models.SupportTicket
	var autoEcoleID sql.NullString
	err := r.db.QueryRowContext(
		ctx, query,
		input.AutoEcoleID, input.Nom, input.Email, input.Sujet, input.Priorite, input.Message,
	).Scan(
		&t.ID, &autoEcoleID, &t.Nom, &t.Email, &t.Sujet, &t.Priorite, &t.Message, &t.Statut, &t.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("erreur création ticket support: %w", err)
	}

	if autoEcoleID.Valid {
		v := autoEcoleID.String
		t.AutoEcoleID = &v
		_ = r.db.QueryRowContext(ctx, `SELECT nom_commercial FROM auto_ecoles WHERE id = $1`, v).Scan(&t.AutoEcoleNom)
	}

	return &t, nil
}

// GetAllSupportTickets liste tous les tickets pour le Super-Admin
func (r *CommunicationRepository) GetAllSupportTickets(ctx context.Context) ([]models.SupportTicket, error) {
	query := `
		SELECT 
			st.id, st.auto_ecole_id, COALESCE(e.nom_commercial, 'Général') as auto_ecole_nom,
			st.nom, st.email, st.sujet, st.priorite, st.message, st.statut, st.created_at
		FROM reseau_support_tickets st
		LEFT JOIN auto_ecoles e ON e.id = st.auto_ecole_id
		ORDER BY st.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération tickets: %w", err)
	}
	defer rows.Close()

	var list []models.SupportTicket
	for rows.Next() {
		var t models.SupportTicket
		var autoEcoleID sql.NullString
		if err := rows.Scan(
			&t.ID, &autoEcoleID, &t.AutoEcoleNom, &t.Nom, &t.Email, &t.Sujet,
			&t.Priorite, &t.Message, &t.Statut, &t.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("erreur scan ticket: %w", err)
		}
		if autoEcoleID.Valid {
			v := autoEcoleID.String
			t.AutoEcoleID = &v
		}
		list = append(list, t)
	}
	if list == nil {
		list = []models.SupportTicket{}
	}
	return list, nil
}


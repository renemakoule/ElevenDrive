package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"autoschool/backend/internal/models"
)

// EnregistrerSeanceCode : Émargement en salle et crédit +2h par présence
func (r *GuichetRepository) EnregistrerSeanceCode(
	ctx context.Context,
	autoEcoleID, agenceID string,
	input *models.EnregistrerSeanceCodeInput,
) (*models.EnregistrerSeanceCodeResponse, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	seanceID := uuid.New().String()
	seanceQuery := `
		INSERT INTO seances_code (
			id, agence_id, theme_traite, formateur_nom, heure_debut, heure_fin
		) VALUES ($1, $2, $3, $4, $5, $6)
	`
	if _, err := tx.ExecContext(ctx, seanceQuery, seanceID, agenceID, input.ThemeTraite, input.FormateurNom, input.HeureDebut, input.HeureFin); err != nil {
		return nil, fmt.Errorf("erreur création séance code: %w", err)
	}

	presentsCount := 0
	for _, p := range input.Presences {
		if p.EstPresent {
			presentsCount++
			_, _ = tx.ExecContext(ctx, `
				INSERT INTO presences_code (id, seance_id, candidat_id, est_present, score_test)
				VALUES ($1, $2, $3, $4, $5)
			`, uuid.New().String(), seanceID, p.CandidatID, p.EstPresent, p.ScoreTest)

			_, _ = tx.ExecContext(ctx, `
				UPDATE dossiers_inscription
				SET heures_code_effectuees = heures_code_effectuees + 2, updated_at = NOW()
				WHERE candidat_id = $1 AND auto_ecole_id = $2
			`, p.CandidatID, autoEcoleID)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.EnregistrerSeanceCodeResponse{
		SeanceID:      seanceID,
		TotalPresents: presentsCount,
		DateSeance:    time.Now().Format("2006-01-02"),
		Statut:        "SUCCES",
	}, nil
}

// GetSeancesCode : Liste les séances de code programmées (Présentiel / En ligne)
func (r *GuichetRepository) GetSeancesCode(ctx context.Context, autoEcoleID, filtreType string) ([]models.SeanceCode, error) {
	// Clôture automatique des séances dont la date et l'heure de fin sont expirées
	_, _ = r.db.ExecContext(ctx, `
		UPDATE seances_code sc
		SET statut_seance = 'TERMINEE'
		FROM agences a
		WHERE a.id = sc.agence_id
		  AND a.auto_ecole_id = $1
		  AND COALESCE(sc.statut_seance, 'PLANIFIEE') NOT IN ('TERMINEE', 'ANNULEE')
		  AND (
			sc.date_seance < CURRENT_DATE
			OR (sc.date_seance = CURRENT_DATE AND COALESCE(NULLIF(sc.heure_fin, ''), '23:59') <= to_char(NOW(), 'HH24:MI'))
		  )
	`, autoEcoleID)

	query := `
		SELECT sc.id, a.auto_ecole_id, sc.agence_id, to_char(sc.date_seance, 'YYYY-MM-DD'),
		       sc.heure_debut, sc.heure_fin, sc.theme_traite, sc.formateur_nom, sc.moniteur_id,
		       COALESCE(sc.type_seance, 'PRESENTIEL'), sc.lien_visio, sc.mot_de_passe_visio,
		       COALESCE(sc.statut_seance, 'PLANIFIEE')
		FROM seances_code sc
		JOIN agences a ON a.id = sc.agence_id
		WHERE a.auto_ecole_id = $1
		  AND COALESCE(sc.heure_debut, '') <> '00:00'
	`
	args := []interface{}{autoEcoleID}
	if filtreType != "" && strings.ToUpper(filtreType) != "ALL" {
		query += ` AND UPPER(COALESCE(sc.type_seance, 'PRESENTIEL')) = UPPER($2)`
		args = append(args, filtreType)
	}
	query += ` ORDER BY sc.date_seance DESC, sc.heure_debut ASC LIMIT 50`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération séances code: %w", err)
	}
	defer rows.Close()

	var list []models.SeanceCode
	for rows.Next() {
		var s models.SeanceCode
		var monID, lien, mdp sql.NullString
		if err := rows.Scan(
			&s.ID, &s.AutoEcoleID, &s.AgenceID, &s.DateSeance,
			&s.HeureDebut, &s.HeureFin, &s.ThemeTraite, &s.FormateurNom, &monID,
			&s.TypeSeance, &lien, &mdp, &s.StatutSeance,
		); err != nil {
			return nil, err
		}
		if monID.Valid {
			s.MoniteurID = &monID.String
		}
		if lien.Valid {
			s.LienVisio = &lien.String
		}
		if mdp.Valid {
			s.MotDePasseVisio = &mdp.String
		}
		list = append(list, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if list == nil {
		list = []models.SeanceCode{}
	}
	return list, nil
}

// CreerSeanceCode : Planifie une nouvelle séance théorique de code
func (r *GuichetRepository) CreerSeanceCode(ctx context.Context, autoEcoleID, agenceID string, input models.CreerSeanceCodeInput) (*models.SeanceCode, error) {
	if agenceID == "" {
		_ = r.db.QueryRowContext(ctx, `SELECT id FROM agences WHERE auto_ecole_id = $1 ORDER BY est_siege DESC LIMIT 1`, autoEcoleID).Scan(&agenceID)
	}
	seanceID := uuid.New().String()
	typeSeance := "PRESENTIEL"
	if strings.ToUpper(input.TypeSeance) == "EN_LIGNE" {
		typeSeance = "EN_LIGNE"
	}
	dateSeance := input.DateSeance
	if strings.TrimSpace(dateSeance) == "" {
		dateSeance = time.Now().Format("2006-01-02")
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO seances_code (
			id, agence_id, theme_traite, formateur_nom, moniteur_id,
			date_seance, heure_debut, heure_fin, type_seance,
			lien_visio, mot_de_passe_visio, statut_seance
		) VALUES ($1, $2, $3, $4, $5, $6::date, $7, $8, $9, $10, $11, 'PLANIFIEE')
	`, seanceID, agenceID, input.ThemeTraite, input.FormateurNom, input.MoniteurID,
		dateSeance, input.HeureDebut, input.HeureFin, typeSeance,
		input.LienVisio, input.MotDePasseVisio,
	)
	if err != nil {
		return nil, fmt.Errorf("erreur planification séance code: %w", err)
	}

	return &models.SeanceCode{
		ID:              seanceID,
		AutoEcoleID:     autoEcoleID,
		AgenceID:        agenceID,
		DateSeance:      dateSeance,
		HeureDebut:      input.HeureDebut,
		HeureFin:        input.HeureFin,
		ThemeTraite:     input.ThemeTraite,
		FormateurNom:    input.FormateurNom,
		MoniteurID:      input.MoniteurID,
		TypeSeance:      typeSeance,
		LienVisio:       input.LienVisio,
		MotDePasseVisio: input.MotDePasseVisio,
		StatutSeance:    "PLANIFIEE",
	}, nil
}

// ReprogrammerSeanceCode : Met à jour la date/horaires d'un cours expiré et le repasse en PLANIFIEE
func (r *GuichetRepository) ReprogrammerSeanceCode(ctx context.Context, autoEcoleID, id, dateSeance, heureDebut, heureFin string) (*models.SeanceCode, error) {
	var s models.SeanceCode
	var monID, lien, mdp sql.NullString
	err := r.db.QueryRowContext(ctx, `
		UPDATE seances_code sc
		SET date_seance = $1::date,
		    heure_debut = $2,
		    heure_fin = $3,
		    statut_seance = 'PLANIFIEE'
		FROM agences a
		WHERE sc.id = $4 AND a.id = sc.agence_id AND a.auto_ecole_id = $5
		RETURNING sc.id, a.auto_ecole_id, sc.agence_id, to_char(sc.date_seance, 'YYYY-MM-DD'),
		          sc.heure_debut, sc.heure_fin, sc.theme_traite, sc.formateur_nom, sc.moniteur_id,
		          COALESCE(sc.type_seance, 'PRESENTIEL'), sc.lien_visio, sc.mot_de_passe_visio, sc.statut_seance
	`, dateSeance, heureDebut, heureFin, id, autoEcoleID).Scan(
		&s.ID, &s.AutoEcoleID, &s.AgenceID, &s.DateSeance,
		&s.HeureDebut, &s.HeureFin, &s.ThemeTraite, &s.FormateurNom, &monID,
		&s.TypeSeance, &lien, &mdp, &s.StatutSeance,
	)
	if err != nil {
		return nil, fmt.Errorf("séance introuvable ou reprogrammation impossible: %w", err)
	}
	if monID.Valid {
		s.MoniteurID = &monID.String
	}
	if lien.Valid {
		s.LienVisio = &lien.String
	}
	if mdp.Valid {
		s.MotDePasseVisio = &mdp.String
	}
	return &s, nil
}

// Tutoriels Vidéos Pédagogiques (CRUD + compteur de vues)
func (r *GuichetRepository) GetTutorielsVideos(ctx context.Context, autoEcoleID, categorie string) ([]models.TutorielVideo, error) {
	query := `
		SELECT id, auto_ecole_id, titre, categorie_permis, description, url_video,
		       COALESCE(plateforme, 'YOUTUBE'), COALESCE(duree_minutes, 5), miniature_url,
		       COALESCE(vues_count, 0), ajoute_par, to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
		FROM tutoriels_video
		WHERE auto_ecole_id = $1
	`
	args := []interface{}{autoEcoleID}
	if categorie != "" && strings.ToUpper(categorie) != "ALL" {
		query += ` AND (UPPER(categorie_permis) = UPPER($2) OR UPPER(categorie_permis) = 'TOUTES')`
		args = append(args, categorie)
	}
	query += ` ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return []models.TutorielVideo{}, nil
	}
	defer rows.Close()

	var list []models.TutorielVideo
	for rows.Next() {
		var t models.TutorielVideo
		var desc, mini, ajoute sql.NullString
		if err := rows.Scan(
			&t.ID, &t.AutoEcoleID, &t.Titre, &t.CategoriePermis, &desc, &t.UrlVideo,
			&t.Plateforme, &t.DureeMinutes, &mini, &t.VuesCount, &ajoute, &t.CreatedAt,
		); err == nil {
			if desc.Valid {
				t.Description = &desc.String
			}
			if mini.Valid {
				t.MiniatureURL = &mini.String
			}
			if ajoute.Valid {
				t.AjoutePar = &ajoute.String
			}
			list = append(list, t)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if list == nil {
		list = []models.TutorielVideo{}
	}
	return list, nil
}

func (r *GuichetRepository) CreerTutorielVideo(ctx context.Context, autoEcoleID string, input models.CreerTutorielVideoInput) (*models.TutorielVideo, error) {
	id := uuid.New().String()
	plateforme := "YOUTUBE"
	if input.Plateforme != nil && *input.Plateforme != "" {
		plateforme = *input.Plateforme
	}
	duree := 5
	if input.DureeMinutes != nil && *input.DureeMinutes > 0 {
		duree = *input.DureeMinutes
	}
	cat := input.CategoriePermis
	if strings.TrimSpace(cat) == "" {
		cat = "TOUTES"
	}

	var createdAt string
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO tutoriels_video (
			id, auto_ecole_id, titre, categorie_permis, description, url_video,
			plateforme, duree_minutes, miniature_url, vues_count, ajoute_par
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 0, $10)
		RETURNING to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
	`, id, autoEcoleID, input.Titre, cat, input.Description, input.UrlVideo,
		plateforme, duree, input.MiniatureURL, input.AjoutePar,
	).Scan(&createdAt)
	if err != nil {
		return nil, fmt.Errorf("erreur création tutoriel vidéo: %w", err)
	}

	return &models.TutorielVideo{
		ID:              id,
		AutoEcoleID:     autoEcoleID,
		Titre:           input.Titre,
		CategoriePermis: cat,
		Description:     input.Description,
		UrlVideo:        input.UrlVideo,
		Plateforme:      plateforme,
		DureeMinutes:    duree,
		MiniatureURL:    input.MiniatureURL,
		VuesCount:       0,
		AjoutePar:       input.AjoutePar,
		CreatedAt:       createdAt,
	}, nil
}

func (r *GuichetRepository) SupprimerTutorielVideo(ctx context.Context, autoEcoleID, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM tutoriels_video WHERE id = $1 AND auto_ecole_id = $2`, id, autoEcoleID)
	return err
}

func (r *GuichetRepository) IncrementerVueTutoriel(ctx context.Context, autoEcoleID, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE tutoriels_video SET vues_count = COALESCE(vues_count, 0) + 1, updated_at = NOW() WHERE id = $1 AND auto_ecole_id = $2`, id, autoEcoleID)
	return err
}

// Supports de Cours Théoriques (CRUD + compteur de téléchargements)
func (r *GuichetRepository) GetSupportsCours(ctx context.Context, autoEcoleID, categorie string) ([]models.SupportCours, error) {
	query := `
		SELECT id, auto_ecole_id, titre, categorie_permis, description,
		       COALESCE(type_support, 'PDF'), fichier_nom, fichier_taille, fichier_data,
		       url_externe, ajoute_par, COALESCE(telechargements_count, 0),
		       to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
		FROM supports_cours
		WHERE auto_ecole_id = $1
	`
	args := []interface{}{autoEcoleID}
	if categorie != "" && strings.ToUpper(categorie) != "ALL" {
		query += ` AND (UPPER(categorie_permis) = UPPER($2) OR UPPER(categorie_permis) = 'TOUTES')`
		args = append(args, categorie)
	}
	query += ` ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return []models.SupportCours{}, nil
	}
	defer rows.Close()

	var list []models.SupportCours
	for rows.Next() {
		var s models.SupportCours
		var desc, fNom, fTaille, fData, urlExt, ajoute sql.NullString
		if err := rows.Scan(
			&s.ID, &s.AutoEcoleID, &s.Titre, &s.CategoriePermis, &desc,
			&s.TypeSupport, &fNom, &fTaille, &fData, &urlExt, &ajoute,
			&s.TelechargementsCount, &s.CreatedAt,
		); err == nil {
			if desc.Valid {
				s.Description = &desc.String
			}
			if fNom.Valid {
				s.FichierNom = &fNom.String
			}
			if fTaille.Valid {
				s.FichierTaille = &fTaille.String
			}
			if fData.Valid {
				s.FichierData = &fData.String
			}
			if urlExt.Valid {
				s.UrlExterne = &urlExt.String
			}
			if ajoute.Valid {
				s.AjoutePar = &ajoute.String
			}
			list = append(list, s)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if list == nil {
		list = []models.SupportCours{}
	}
	return list, nil
}

func (r *GuichetRepository) CreerSupportCours(ctx context.Context, autoEcoleID string, input models.CreerSupportCoursInput) (*models.SupportCours, error) {
	id := uuid.New().String()
	cat := input.CategoriePermis
	if strings.TrimSpace(cat) == "" {
		cat = "TOUTES"
	}
	typeSup := input.TypeSupport
	if strings.TrimSpace(typeSup) == "" {
		typeSup = "PDF"
	}

	var createdAt string
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO supports_cours (
			id, auto_ecole_id, titre, categorie_permis, description, type_support,
			fichier_nom, fichier_taille, fichier_data, url_externe, ajoute_par, telechargements_count
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, 0)
		RETURNING to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
	`, id, autoEcoleID, input.Titre, cat, input.Description, typeSup,
		input.FichierNom, input.FichierTaille, input.FichierData, input.UrlExterne, input.AjoutePar,
	).Scan(&createdAt)
	if err != nil {
		return nil, fmt.Errorf("erreur création support de cours: %w", err)
	}

	return &models.SupportCours{
		ID:                   id,
		AutoEcoleID:          autoEcoleID,
		Titre:                input.Titre,
		CategoriePermis:      cat,
		Description:          input.Description,
		TypeSupport:          typeSup,
		FichierNom:           input.FichierNom,
		FichierTaille:        input.FichierTaille,
		FichierData:          input.FichierData,
		UrlExterne:           input.UrlExterne,
		AjoutePar:            input.AjoutePar,
		TelechargementsCount: 0,
		CreatedAt:            createdAt,
	}, nil
}

func (r *GuichetRepository) SupprimerSupportCours(ctx context.Context, autoEcoleID, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM supports_cours WHERE id = $1 AND auto_ecole_id = $2`, id, autoEcoleID)
	return err
}

func (r *GuichetRepository) IncrementerTelechargementSupport(ctx context.Context, autoEcoleID, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE supports_cours SET telechargements_count = COALESCE(telechargements_count, 0) + 1, updated_at = NOW() WHERE id = $1 AND auto_ecole_id = $2`, id, autoEcoleID)
	return err
}

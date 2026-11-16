package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/google/uuid"

	"autoschool/backend/internal/models"
)

// EnregistrerQuizResultat : Enregistre le score obtenu à un test d'entraînement théorique
func (r *GuichetRepository) EnregistrerQuizResultat(
	ctx context.Context,
	autoEcoleID string,
	input models.EnregistrerQuizResultatInput,
) error {
	var agenceID string
	err := r.db.QueryRowContext(ctx, `
		SELECT agence_id FROM candidats WHERE id = $1 AND auto_ecole_id = $2
	`, input.CandidatID, autoEcoleID).Scan(&agenceID)
	if err != nil {
		return fmt.Errorf("candidat introuvable: %w", err)
	}

	seanceID := uuid.New().String()
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO seances_code (
			id, agence_id, theme_traite, formateur_nom, date_seance, heure_debut, heure_fin
		) VALUES (
			$1, $2, 'Quiz d''entraînement Code PWA', 'Auto-Évaluation Apprenant', CURRENT_DATE, '00:00', '00:00'
		)
	`, seanceID, agenceID)
	if err != nil {
		return fmt.Errorf("erreur création séance quiz: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO presences_code (
			id, seance_id, candidat_id, est_present, score_test, heure_pointage, created_at
		) VALUES (
			$1, $2, $3, true, $4, NOW(), NOW()
		)
	`, uuid.New().String(), seanceID, input.CandidatID, input.Score)
	if err != nil {
		return fmt.Errorf("erreur enregistrement score quiz: %w", err)
	}

	return nil
}

// ==============================================================================
// FORMULAIRES DE QUIZ INTERACTIFS (STYLE GOOGLE FORMS) & RÉPONSES ÉLÈVES
// ==============================================================================

func (r *GuichetRepository) EnsureQuizTablesExist(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS quiz_formulaires (
			id UUID PRIMARY KEY,
			auto_ecole_id UUID NOT NULL,
			moniteur_id TEXT,
			moniteur_nom TEXT NOT NULL DEFAULT 'Moniteur',
			titre TEXT NOT NULL,
			description TEXT,
			categorie_permis TEXT NOT NULL DEFAULT 'TOUTES',
			duree_minutes INT NOT NULL DEFAULT 15,
			seuil_reussite_pct INT NOT NULL DEFAULT 80,
			est_publie BOOLEAN NOT NULL DEFAULT TRUE,
			questions_json JSONB NOT NULL DEFAULT '[]'::jsonb,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS quiz_reponses_eleves (
			id UUID PRIMARY KEY,
			quiz_id UUID NOT NULL,
			quiz_titre TEXT NOT NULL DEFAULT '',
			auto_ecole_id UUID NOT NULL,
			candidat_id TEXT NOT NULL,
			candidat_nom TEXT NOT NULL,
			numero_dossier TEXT NOT NULL DEFAULT '',
			categorie_permis TEXT NOT NULL DEFAULT 'B',
			score_obtenu INT NOT NULL DEFAULT 0,
			total_points INT NOT NULL DEFAULT 1,
			pourcentage INT NOT NULL DEFAULT 0,
			est_reussi BOOLEAN NOT NULL DEFAULT FALSE,
			duree_secondes INT NOT NULL DEFAULT 0,
			reponses_detail_json JSONB NOT NULL DEFAULT '[]'::jsonb,
			submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *GuichetRepository) GetQuizFormulaires(ctx context.Context, autoEcoleID, categorie string) ([]models.QuizFormulaire, error) {
	_ = r.EnsureQuizTablesExist(ctx)

	query := `
		SELECT q.id, q.auto_ecole_id, q.moniteur_id, q.moniteur_nom, q.titre, q.description,
		       q.categorie_permis, q.duree_minutes, q.seuil_reussite_pct, q.est_publie,
		       COALESCE(q.questions_json, '[]'::jsonb),
		       (SELECT COUNT(*) FROM quiz_reponses_eleves qr WHERE qr.quiz_id = q.id) AS reponses_count,
		       to_char(q.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
		FROM quiz_formulaires q
		WHERE q.auto_ecole_id = $1
	`
	args := []interface{}{autoEcoleID}
	if categorie != "" && strings.ToUpper(categorie) != "ALL" {
		query += ` AND (UPPER(q.categorie_permis) = UPPER($2) OR UPPER(q.categorie_permis) = 'TOUTES')`
		args = append(args, categorie)
	}
	query += ` ORDER BY q.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return []models.QuizFormulaire{}, nil
	}
	defer rows.Close()

	var list []models.QuizFormulaire
	for rows.Next() {
		var q models.QuizFormulaire
		var monID, desc sql.NullString
		var qJSON []byte
		if err := rows.Scan(
			&q.ID, &q.AutoEcoleID, &monID, &q.MoniteurNom, &q.Titre, &desc,
			&q.CategoriePermis, &q.DureeMinutes, &q.SeuilReussitePct, &q.EstPublie,
			&qJSON, &q.ReponsesCount, &q.CreatedAt,
		); err == nil {
			if monID.Valid {
				q.MoniteurID = &monID.String
			}
			if desc.Valid {
				q.Description = &desc.String
			}
			if len(qJSON) > 0 {
				_ = json.Unmarshal(qJSON, &q.Questions)
			}
			if q.Questions == nil {
				q.Questions = []models.QuizQuestionItem{}
			}
			list = append(list, q)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if list == nil {
		list = []models.QuizFormulaire{}
	}
	return list, nil
}

func (r *GuichetRepository) CreerQuizFormulaire(ctx context.Context, autoEcoleID string, input models.CreerQuizFormulaireInput) (*models.QuizFormulaire, error) {
	_ = r.EnsureQuizTablesExist(ctx)

	id := uuid.New().String()
	cat := strings.TrimSpace(input.CategoriePermis)
	if cat == "" {
		cat = "TOUTES"
	}
	duree := input.DureeMinutes
	if duree <= 0 {
		duree = 15
	}
	seuil := input.SeuilReussitePct
	if seuil <= 0 || seuil > 100 {
		seuil = 80
	}
	monNom := strings.TrimSpace(input.MoniteurNom)
	if monNom == "" {
		monNom = "Moniteur"
	}
	if input.Questions == nil {
		input.Questions = []models.QuizQuestionItem{}
	}
	qJSON, _ := json.Marshal(input.Questions)

	var createdAt string
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO quiz_formulaires (
			id, auto_ecole_id, moniteur_id, moniteur_nom, titre, description,
			categorie_permis, duree_minutes, seuil_reussite_pct, est_publie, questions_json
		) VALUES ($1::uuid, $2::uuid, $3, $4, $5, $6, $7, $8, $9, $10, $11::jsonb)
		RETURNING to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
	`, id, autoEcoleID, input.MoniteurID, monNom, input.Titre, input.Description,
		cat, duree, seuil, input.EstPublie, string(qJSON),
	).Scan(&createdAt)
	if err != nil {
		return nil, fmt.Errorf("erreur création formulaire quiz: %w", err)
	}

	return &models.QuizFormulaire{
		ID:               id,
		AutoEcoleID:      autoEcoleID,
		MoniteurID:       input.MoniteurID,
		MoniteurNom:      monNom,
		Titre:            input.Titre,
		Description:      input.Description,
		CategoriePermis:  cat,
		DureeMinutes:     duree,
		SeuilReussitePct: seuil,
		EstPublie:        input.EstPublie,
		Questions:        input.Questions,
		ReponsesCount:    0,
		CreatedAt:        createdAt,
	}, nil
}

func (r *GuichetRepository) SupprimerQuizFormulaire(ctx context.Context, autoEcoleID, id string) error {
	_ = r.EnsureQuizTablesExist(ctx)
	_, _ = r.db.ExecContext(ctx, `DELETE FROM quiz_reponses_eleves WHERE quiz_id = $1::uuid AND auto_ecole_id = $2::uuid`, id, autoEcoleID)
	_, err := r.db.ExecContext(ctx, `DELETE FROM quiz_formulaires WHERE id = $1::uuid AND auto_ecole_id = $2::uuid`, id, autoEcoleID)
	return err
}

func (r *GuichetRepository) GetQuizSoumissions(ctx context.Context, autoEcoleID, quizID, candidatID string) ([]models.QuizSoumissionEleve, error) {
	_ = r.EnsureQuizTablesExist(ctx)

	query := `
		SELECT id, quiz_id, quiz_titre, auto_ecole_id, candidat_id, candidat_nom,
		       numero_dossier, categorie_permis, score_obtenu, total_points, pourcentage,
		       est_reussi, duree_secondes, COALESCE(reponses_detail_json, '[]'::jsonb),
		       to_char(submitted_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
		FROM quiz_reponses_eleves
		WHERE auto_ecole_id = $1::uuid
	`
	args := []interface{}{autoEcoleID}
	argIdx := 2
	if quizID != "" {
		query += fmt.Sprintf(` AND quiz_id = $%d::uuid`, argIdx)
		args = append(args, quizID)
		argIdx++
	}
	if candidatID != "" {
		query += fmt.Sprintf(` AND candidat_id = $%d`, argIdx)
		args = append(args, candidatID)
	}
	query += ` ORDER BY submitted_at DESC LIMIT 200`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return []models.QuizSoumissionEleve{}, nil
	}
	defer rows.Close()

	var list []models.QuizSoumissionEleve
	for rows.Next() {
		var s models.QuizSoumissionEleve
		var dJSON []byte
		if err := rows.Scan(
			&s.ID, &s.QuizID, &s.QuizTitre, &s.AutoEcoleID, &s.CandidatID, &s.CandidatNom,
			&s.NumeroDossier, &s.CategoriePermis, &s.ScoreObtenu, &s.TotalPoints, &s.Pourcentage,
			&s.EstReussi, &s.DureeSecondes, &dJSON, &s.SubmittedAt,
		); err == nil {
			if len(dJSON) > 0 {
				_ = json.Unmarshal(dJSON, &s.Details)
			}
			if s.Details == nil {
				s.Details = []models.QuizReponseDetailItem{}
			}
			list = append(list, s)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if list == nil {
		list = []models.QuizSoumissionEleve{}
	}
	return list, nil
}

func (r *GuichetRepository) SoumettreQuizReponse(ctx context.Context, autoEcoleID string, input models.SoumettreQuizReponseInput) (*models.QuizSoumissionEleve, error) {
	_ = r.EnsureQuizTablesExist(ctx)

	id := uuid.New().String()
	if input.Details == nil {
		input.Details = []models.QuizReponseDetailItem{}
	}
	dJSON, _ := json.Marshal(input.Details)

	var submittedAt string
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO quiz_reponses_eleves (
			id, quiz_id, quiz_titre, auto_ecole_id, candidat_id, candidat_nom,
			numero_dossier, categorie_permis, score_obtenu, total_points, pourcentage,
			est_reussi, duree_secondes, reponses_detail_json
		) VALUES ($1::uuid, $2::uuid, $3, $4::uuid, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14::jsonb)
		RETURNING to_char(submitted_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
	`, id, input.QuizID, input.QuizTitre, autoEcoleID, input.CandidatID, input.CandidatNom,
		input.NumeroDossier, input.CategoriePermis, input.ScoreObtenu, input.TotalPoints,
		input.Pourcentage, input.EstReussi, input.DureeSecondes, string(dJSON),
	).Scan(&submittedAt)
	if err != nil {
		return nil, fmt.Errorf("erreur enregistrement réponse quiz élève: %w", err)
	}

	// Mettre également à jour le dernier score sur 40 dans presences_code pour le suivi du dossier candidat
	scoreSur40 := int(math.Round(float64(input.Pourcentage) * 40.0 / 100.0))
	if input.CandidatID != "" {
		_ = r.EnregistrerQuizResultat(ctx, autoEcoleID, models.EnregistrerQuizResultatInput{
			CandidatID: input.CandidatID,
			Score:      scoreSur40,
		})
	}

	return &models.QuizSoumissionEleve{
		ID:              id,
		QuizID:          input.QuizID,
		QuizTitre:       input.QuizTitre,
		AutoEcoleID:     autoEcoleID,
		CandidatID:      input.CandidatID,
		CandidatNom:     input.CandidatNom,
		NumeroDossier:   input.NumeroDossier,
		CategoriePermis: input.CategoriePermis,
		ScoreObtenu:     input.ScoreObtenu,
		TotalPoints:     input.TotalPoints,
		Pourcentage:     input.Pourcentage,
		EstReussi:       input.EstReussi,
		DureeSecondes:   input.DureeSecondes,
		Details:         input.Details,
		SubmittedAt:     submittedAt,
	}, nil
}


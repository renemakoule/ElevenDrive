package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"autoschool/backend/internal/models"
)

func defaultCompetencesREMC(c1Niveau int) []models.CompetenceItem {
	return []models.CompetenceItem{
		{ID: "c1", Code: "A1", Categorie: "Maniement", Libelle: "Installation au poste de conduite & réglages", Niveau: c1Niveau},
		{ID: "c2", Code: "A2", Categorie: "Maniement", Libelle: "Démarrage en côte & point de patinage", Niveau: 2},
		{ID: "c3", Code: "B1", Categorie: "Circulation", Libelle: "Franchissement des ronds-points & giratoires", Niveau: 2},
		{ID: "c4", Code: "B2", Categorie: "Circulation", Libelle: "Insertion et circulation sur voie rapide", Niveau: 1},
		{ID: "c5", Code: "C1", Categorie: "Manœuvres", Libelle: "Stationnement en créneau & en bataille", Niveau: 2},
		{ID: "c6", Code: "D1", Categorie: "Autonomie", Libelle: "Anticipation des piétons & distances de sécurité", Niveau: 2},
	}
}

// GetPlanningConduite : Récupère les séances de conduite du jour (ou toutes les dates si dateSeance == "ALL")
func (r *GuichetRepository) GetPlanningConduite(ctx context.Context, autoEcoleID, moniteurID, dateSeance string) ([]models.SeanceConduite, error) {
	trimmedDate := strings.TrimSpace(dateSeance)
	isAllDates := strings.EqualFold(trimmedDate, "ALL")
	if trimmedDate == "" {
		dateSeance = time.Now().Format("2006-01-02")
	} else {
		dateSeance = trimmedDate
	}

	defaultComps := defaultCompetencesREMC(3)

	query := `
		SELECT sc.id, sc.auto_ecole_id, sc.agence_id, sc.dossier_id,
		       c.id AS candidat_id, c.nom || ' ' || c.prenoms AS candidat_nom,
		       c.telephone_whatsapp, c.photo_url,
		       f.intitule AS formule_nom,
		       sc.moniteur_id,
		       COALESCE(u.nom || ' ' || u.prenoms, 'Moniteur Assigné') AS moniteur_nom,
		       sc.vehicule_id,
		       COALESCE(v.immatriculation, '') AS vehicule_immatriculation,
		       COALESCE(v.marque_modele, '') AS vehicule_modele,
		       COALESCE(v.kilometrage_actuel, 0) AS vehicule_km,
		       to_char(sc.date_seance, 'YYYY-MM-DD') AS date_seance,
		       sc.heure_debut, sc.heure_fin, sc.duree_minutes, sc.statut,
		       COALESCE(d.heures_conduite_effectuees, 0) AS heures_deja_effectuees,
		       COALESCE(f.heures_pratiques, 20) AS heures_forfait_total,
		       COALESCE(sc.competences_travaillees, '[]'::jsonb),
		       sc.commentaire_moniteur,
		       sc.signature_eleve_url
		FROM seances_conduite sc
		JOIN dossiers_inscription d ON d.id = sc.dossier_id
		JOIN candidats c ON c.id = d.candidat_id
		JOIN formules_permis f ON f.id = d.formule_id
		LEFT JOIN utilisateurs u ON u.id = sc.moniteur_id
		LEFT JOIN vehicules v ON v.id = sc.vehicule_id
		WHERE sc.auto_ecole_id = $1
	`
	var rows *sql.Rows
	var err error
	if isAllDates {
		query += ` ORDER BY sc.date_seance DESC, sc.heure_debut ASC LIMIT 100`
		rows, err = r.db.QueryContext(ctx, query, autoEcoleID)
	} else {
		query += ` AND sc.date_seance = $2::date ORDER BY sc.heure_debut ASC`
		rows, err = r.db.QueryContext(ctx, query, autoEcoleID, dateSeance)
	}
	if err == nil {
		defer rows.Close()
		var list []models.SeanceConduite
		for rows.Next() {
			var s models.SeanceConduite
			var compJSON []byte
			var monID, vehID *string
			var vehKm int
			if err := rows.Scan(
				&s.ID, &s.AutoEcoleID, &s.AgenceID, &s.DossierID,
				&s.CandidatID, &s.CandidatNom,
				&s.TelephoneWhatsApp, &s.PhotoURL,
				&s.FormuleNom,
				&monID,
				&s.MoniteurNom,
				&vehID,
				&s.VehiculeImmatriculation,
				&s.VehiculeModele,
				&vehKm,
				&s.DateSeance,
				&s.HeureDebut, &s.HeureFin, &s.DureeMinutes, &s.Statut,
				&s.HeuresDejaEffectuees,
				&s.HeuresForfaitTotal,
				&compJSON,
				&s.CommentaireMoniteur,
				&s.SignatureEleveURL,
			); err == nil {
				s.MoniteurID = monID
				s.VehiculeID = vehID
				if len(compJSON) > 0 {
					_ = json.Unmarshal(compJSON, &s.Competences)
				}
				if len(s.Competences) == 0 {
					s.Competences = defaultComps
				}
				if vehKm > 0 {
					startKm := vehKm
					endKm := vehKm + 15
					if s.Statut == "EFFECTUE" && vehKm >= 15 {
						startKm = vehKm - 15
						endKm = vehKm
					}
					s.KmDebut = &startKm
					s.KmFin = &endKm
				}
				list = append(list, s)
			}
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
		if list == nil {
			list = []models.SeanceConduite{}
		}
		return list, nil
	}

	return []models.SeanceConduite{}, nil
}

// ValiderSeanceConduite : Émargement numérique, mise à jour des compétences, crédit d'heures et km
func (r *GuichetRepository) ValiderSeanceConduite(
	ctx context.Context,
	autoEcoleID string,
	moniteurID string,
	input *models.ValiderSeanceConduiteInput,
) (*models.SeanceConduite, error) {
	if input.DossierID == "" {
		return nil, errors.New("dossier_id obligatoire")
	}
	if input.DureeHeures <= 0 {
		input.DureeHeures = 2
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("erreur démarrage transaction: %w", err)
	}
	defer tx.Rollback()

	// Récupérer agence_id du dossier
	var agenceID string
	err = tx.QueryRowContext(ctx, `SELECT agence_id FROM dossiers_inscription WHERE id = $1 AND auto_ecole_id = $2`, input.DossierID, autoEcoleID).Scan(&agenceID)
	if err != nil {
		return nil, fmt.Errorf("dossier d'inscription introuvable: %w", err)
	}

	// Mettre à jour les heures de conduite du dossier
	_, err = tx.ExecContext(ctx, `
		UPDATE dossiers_inscription
		SET heures_conduite_effectuees = heures_conduite_effectuees + $1, updated_at = NOW()
		WHERE id = $2 AND auto_ecole_id = $3
	`, input.DureeHeures, input.DossierID, autoEcoleID)
	if err != nil {
		return nil, fmt.Errorf("erreur mise à jour heures dossier: %w", err)
	}

	// Mettre à jour le kilométrage du véhicule si fourni
	if input.VehiculeID != nil && *input.VehiculeID != "" && input.KmFin > 0 {
		_, _ = tx.ExecContext(ctx, `
			UPDATE vehicules
			SET kilometrage_actuel = GREATEST(kilometrage_actuel, $1), updated_at = NOW()
			WHERE id = $2
		`, input.KmFin, *input.VehiculeID)
	}

	// Mettre à jour ou insérer la séance de conduite
	compJSON, _ := json.Marshal(input.Competences)
	seanceID := strings.TrimSpace(input.SeanceID)
	var monIDVal *string
	if moniteurID != "" {
		monIDVal = &moniteurID
	}

	updatedExisting := false
	if seanceID != "" && !strings.HasPrefix(seanceID, "plan-") {
		res, errUpd := tx.ExecContext(ctx, `
			UPDATE seances_conduite
			SET statut = 'EFFECTUE',
			    duree_minutes = $1,
			    competences_travaillees = $2::jsonb,
			    commentaire_moniteur = $3,
			    signature_eleve_url = $4,
			    moniteur_id = COALESCE($5, moniteur_id),
			    vehicule_id = COALESCE($6, vehicule_id),
			    updated_at = NOW()
			WHERE id = $7 AND auto_ecole_id = $8
		`, input.DureeHeures*60, string(compJSON), input.Commentaire, input.SignatureEleve, monIDVal, input.VehiculeID, seanceID, autoEcoleID)
		if errUpd == nil {
			if rowsAff, _ := res.RowsAffected(); rowsAff > 0 {
				updatedExisting = true
			}
		}
	}

	if !updatedExisting {
		seanceID = uuid.New().String()
		insertQuery := `
			INSERT INTO seances_conduite (
				id, auto_ecole_id, agence_id, dossier_id, moniteur_id, vehicule_id,
				date_seance, heure_debut, heure_fin, duree_minutes, statut,
				competences_travaillees, commentaire_moniteur, signature_eleve_url,
				created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6,
				CURRENT_DATE, '08:00', '10:00', $7, 'EFFECTUE',
				$8::jsonb, $9, $10,
				NOW(), NOW()
			)
		`
		if _, err := tx.ExecContext(
			ctx, insertQuery,
			seanceID, autoEcoleID, agenceID, input.DossierID, monIDVal, input.VehiculeID,
			input.DureeHeures*60, string(compJSON), input.Commentaire, input.SignatureEleve,
		); err != nil {
			return nil, fmt.Errorf("erreur enregistrement séance conduite: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("erreur validation finale: %w", err)
	}

	return &models.SeanceConduite{
		ID:                  seanceID,
		AutoEcoleID:         autoEcoleID,
		DossierID:           input.DossierID,
		VehiculeID:          input.VehiculeID,
		DureeMinutes:        input.DureeHeures * 60,
		Statut:              "EFFECTUE",
		DateSeance:          time.Now().Format("2006-01-02"),
		Competences:         input.Competences,
		CommentaireMoniteur: &input.Commentaire,
		KmDebut:             &input.KmDebut,
		KmFin:               &input.KmFin,
		SignatureEleveURL:   &input.SignatureEleve,
	}, nil
}

// ReserverSeanceConduite : Réservation d'un créneau de conduite par l'apprenant
func (r *GuichetRepository) ReserverSeanceConduite(
	ctx context.Context,
	autoEcoleID string,
	input models.ReserverSeanceConduiteInput,
) (*models.SeanceConduite, error) {
	// 1. Vérifier le dossier et son éligibilité (règle 2e tranche ou dérogation)
	var d models.DossierInscription
	err := r.db.QueryRowContext(ctx, `
		SELECT id, agence_id, candidat_id, montant_total, solde_restant, derogation_conduite_accordee
		FROM dossiers_inscription
		WHERE id = $1 AND auto_ecole_id = $2
	`, input.DossierID, autoEcoleID).Scan(
		&d.ID, &d.AgenceID, &d.CandidatID, &d.MontantTotal, &d.SoldeRestant, &d.DerogationConduiteAccordee,
	)
	if err != nil {
		return nil, fmt.Errorf("dossier d'inscription introuvable: %w", err)
	}

	// Seuil 40% (2e tranche) ou dérogation accordée (ou forcée par la direction/secrétariat)
	estEligible := (d.SoldeRestant <= d.MontantTotal*0.4) || d.DerogationConduiteAccordee || input.DerogationForcee
	if !estEligible {
		return nil, fmt.Errorf("leçon bloquée : le solde restant (%.0f FCFA) dépasse le seuil autorisé (60%% réglés requis)", d.SoldeRestant)
	}

	// 2. Si aucun véhicule spécifié, prendre un véhicule disponible de l'agence
	var vehIDVal *string = input.VehiculeID
	if vehIDVal == nil || *vehIDVal == "" {
		var firstVehID string
		errVeh := r.db.QueryRowContext(ctx, `
			SELECT id FROM vehicules
			WHERE agence_id = $1 AND statut = 'DISPONIBLE'
			ORDER BY created_at ASC LIMIT 1
		`, d.AgenceID).Scan(&firstVehID)
		if errVeh == nil && firstVehID != "" {
			vehIDVal = &firstVehID
		}
	}

	// 3. Vérifier le moniteur spécifié ou affecter le premier moniteur actif de l'auto-école
	var monIDVal *string = input.MoniteurID
	if monIDVal != nil && *monIDVal != "" {
		var existsID, existingRole string
		errCheck := r.db.QueryRowContext(ctx, `
			SELECT id, role FROM utilisateurs WHERE id = $1 AND auto_ecole_id = $2
		`, *monIDVal, autoEcoleID).Scan(&existsID, &existingRole)
		if errCheck != nil || existsID == "" {
			monIDVal = nil
		} else if existingRole != "MONITEUR" {
			var activeMonID string
			if errM := r.db.QueryRowContext(ctx, `
				SELECT id FROM utilisateurs
				WHERE auto_ecole_id = $1 AND role = 'MONITEUR' AND statut = 'ACTIF'
				ORDER BY created_at ASC LIMIT 1
			`, autoEcoleID).Scan(&activeMonID); errM == nil && activeMonID != "" {
				monIDVal = &activeMonID
			}
		}
	}
	if monIDVal == nil || *monIDVal == "" {
		var firstMonID string
		errMon := r.db.QueryRowContext(ctx, `
			SELECT id FROM utilisateurs
			WHERE auto_ecole_id = $1 AND role = 'MONITEUR' AND statut = 'ACTIF'
			ORDER BY created_at ASC LIMIT 1
		`, autoEcoleID).Scan(&firstMonID)
		if errMon == nil && firstMonID != "" {
			monIDVal = &firstMonID
		}
	}

	// 4. Insérer dans seances_conduite avec la grille officielle des compétences REMC
	defaultComps := defaultCompetencesREMC(2)
	defaultCompsJSON, _ := json.Marshal(defaultComps)

	seanceID := uuid.New().String()
	insertQuery := `
		INSERT INTO seances_conduite (
			id, auto_ecole_id, agence_id, dossier_id, moniteur_id, vehicule_id,
			date_seance, heure_debut, heure_fin, duree_minutes, statut,
			competences_travaillees, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7::date, $8, $9, 120, 'PROGRAMME',
			$10::jsonb, NOW(), NOW()
		)
	`
	_, err = r.db.ExecContext(
		ctx, insertQuery,
		seanceID, autoEcoleID, d.AgenceID, d.ID, monIDVal, vehIDVal,
		input.DateSeance, input.HeureDebut, input.HeureFin, string(defaultCompsJSON),
	)
	if err != nil {
		return nil, fmt.Errorf("erreur réservation leçon de conduite: %w", err)
	}

	return &models.SeanceConduite{
		ID:           seanceID,
		AutoEcoleID:  autoEcoleID,
		AgenceID:     d.AgenceID,
		DossierID:    d.ID,
		MoniteurID:   monIDVal,
		VehiculeID:   vehIDVal,
		DateSeance:   input.DateSeance,
		HeureDebut:   input.HeureDebut,
		HeureFin:     input.HeureFin,
		DureeMinutes: 120,
		Statut:       "PROGRAMME",
		Competences:  defaultComps,
	}, nil
}

// AnnulerSeanceConduite : Supprime ou annule une leçon de conduite programmée
func (r *GuichetRepository) AnnulerSeanceConduite(ctx context.Context, autoEcoleID, seanceID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM seances_conduite
		WHERE id = $1 AND auto_ecole_id = $2
	`, seanceID, autoEcoleID)
	return err
}

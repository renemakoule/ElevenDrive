package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"

	"autoschool/backend/internal/models"
)

// SearchCandidats : Recherche multi-critères avec dossier et formule associés
func (r *GuichetRepository) SearchCandidats(ctx context.Context, autoEcoleID string, term string) ([]models.Candidat, error) {
	term = strings.TrimSpace(term)
	query := `
		SELECT c.id, c.auto_ecole_id, c.agence_id, c.nom, c.prenoms, c.genre,
		       to_char(c.date_naissance, 'YYYY-MM-DD'), c.lieu_naissance, c.nationalite,
		       c.profession, c.quartier_residence, c.telephone_whatsapp, c.telephone_urgence,
		       c.groupe_sanguin, c.photo_url, c.type_piece_id, c.numero_piece_id,
		       c.visite_medicale_statut, c.created_at, c.updated_at,
		       d.numero_dossier, f.intitule, d.solde_restant, d.statut_dossier,
		       COALESCE(f.categorie_permis, 'B') AS categorie_permis
		FROM candidats c
		LEFT JOIN dossiers_inscription d ON d.candidat_id = c.id AND d.auto_ecole_id = c.auto_ecole_id
		LEFT JOIN formules_permis f ON f.id = d.formule_id
		WHERE c.auto_ecole_id = $1
	`
	var rows *sql.Rows
	var err error

	if term != "" {
		likeTerm := "%" + strings.ToLower(term) + "%"
		query += ` AND (LOWER(c.nom) LIKE $2 OR LOWER(c.prenoms) LIKE $2 OR c.telephone_whatsapp LIKE $2 OR LOWER(c.numero_piece_id) LIKE $2 OR LOWER(COALESCE(d.numero_dossier, '')) LIKE $2)`
		query += ` ORDER BY c.created_at DESC LIMIT 50`
		rows, err = r.db.QueryContext(ctx, query, autoEcoleID, likeTerm)
	} else {
		query += ` ORDER BY c.created_at DESC LIMIT 50`
		rows, err = r.db.QueryContext(ctx, query, autoEcoleID)
	}

	if err != nil {
		return nil, fmt.Errorf("erreur recherche candidats: %w", err)
	}
	defer rows.Close()

	var list []models.Candidat
	for rows.Next() {
		var c models.Candidat
		var numDos, formNom, statDos, catPermis sql.NullString
		var soldeRestant sql.NullFloat64

		if err := rows.Scan(
			&c.ID, &c.AutoEcoleID, &c.AgenceID, &c.Nom, &c.Prenoms, &c.Genre,
			&c.DateNaissance, &c.LieuNaissance, &c.Nationalite, &c.Profession,
			&c.QuartierResidence, &c.TelephoneWhatsApp, &c.TelephoneUrgence,
			&c.GroupeSanguin, &c.PhotoURL, &c.TypePieceID, &c.NumeroPieceID,
			&c.VisiteMedicaleStatut, &c.CreatedAt, &c.UpdatedAt,
			&numDos, &formNom, &soldeRestant, &statDos, &catPermis,
		); err != nil {
			return nil, err
		}

		if numDos.Valid {
			c.NumeroDossier = &numDos.String
		}
		if formNom.Valid {
			c.FormuleNom = &formNom.String
		}
		if catPermis.Valid {
			c.CategoriePermis = &catPermis.String
		}
		if soldeRestant.Valid {
			c.SoldeRestant = &soldeRestant.Float64
		}
		if statDos.Valid {
			c.StatutDossier = &statDos.String
		}

		list = append(list, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

// GetCandidatDossier : Dossier complet avec formule, solde et paiements
func (r *GuichetRepository) GetCandidatDossier(ctx context.Context, autoEcoleID, candidatID string) (*models.DossierDetailComplet, error) {
	// Candidat
	var c models.Candidat
	queryCand := `
		SELECT id, auto_ecole_id, agence_id, nom, prenoms, genre,
		       to_char(date_naissance, 'YYYY-MM-DD'), lieu_naissance, nationalite,
		       profession, quartier_residence, telephone_whatsapp, telephone_urgence,
		       groupe_sanguin, photo_url, type_piece_id, numero_piece_id,
		       visite_medicale_statut, created_at, updated_at
		FROM candidats
		WHERE auto_ecole_id = $1 AND id = $2
	`
	err := r.db.QueryRowContext(ctx, queryCand, autoEcoleID, candidatID).Scan(
		&c.ID, &c.AutoEcoleID, &c.AgenceID, &c.Nom, &c.Prenoms, &c.Genre,
		&c.DateNaissance, &c.LieuNaissance, &c.Nationalite, &c.Profession,
		&c.QuartierResidence, &c.TelephoneWhatsApp, &c.TelephoneUrgence,
		&c.GroupeSanguin, &c.PhotoURL, &c.TypePieceID, &c.NumeroPieceID,
		&c.VisiteMedicaleStatut, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("candidat introuvable: %w", err)
	}

	// Dossier
	var d models.DossierInscription
	queryDos := `
		SELECT id, auto_ecole_id, agence_id, candidat_id, formule_id, numero_dossier,
		       to_char(date_inscription, 'YYYY-MM-DD'), montant_total, remise_accordee,
		       solde_restant, statut_dossier, heures_code_effectuees, heures_conduite_effectuees,
		       derogation_conduite_accordee, created_at, updated_at
		FROM dossiers_inscription
		WHERE auto_ecole_id = $1 AND candidat_id = $2
		ORDER BY created_at DESC
		LIMIT 1
	`
	err = r.db.QueryRowContext(ctx, queryDos, autoEcoleID, candidatID).Scan(
		&d.ID, &d.AutoEcoleID, &d.AgenceID, &d.CandidatID, &d.FormuleID, &d.NumeroDossier,
		&d.DateInscription, &d.MontantTotal, &d.RemiseAccordee, &d.SoldeRestant,
		&d.StatutDossier, &d.HeuresCodeEffectuees, &d.HeuresConduiteEffectuees,
		&d.DerogationConduiteAccordee, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("dossier d'inscription introuvable: %w", err)
	}

	// Formule
	var f models.FormulePermis
	queryFormule := `
		SELECT id, auto_ecole_id, categorie_permis, intitule, description,
		       heures_theoriques, heures_pratiques, prix_total, acompte_minimum,
		       est_actif, created_at, updated_at
		FROM formules_permis
		WHERE id = $1
	`
	_ = r.db.QueryRowContext(ctx, queryFormule, d.FormuleID).Scan(
		&f.ID, &f.AutoEcoleID, &f.CategoriePermis, &f.Intitule, &f.Description,
		&f.HeuresTheoriques, &f.HeuresPratiques, &f.PrixTotal, &f.AcompteMinimum,
		&f.EstActif, &f.CreatedAt, &f.UpdatedAt,
	)

	// Paiements
	queryPay := `
		SELECT p.id, p.auto_ecole_id, p.agence_id, p.dossier_id, p.session_caisse_id,
		       p.montant, p.mode_reglement, p.reference_transaction, p.reference_recu,
		       p.solde_apres_paiement, to_char(p.date_paiement, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
		       p.encaisse_par_id, p.est_annule, p.created_at
		FROM paiements p
		WHERE p.dossier_id = $1
		ORDER BY p.date_paiement DESC
	`
	payRows, err := r.db.QueryContext(ctx, queryPay, d.ID)
	paiements := []models.Paiement{}
	if err == nil {
		defer payRows.Close()
		for payRows.Next() {
			var p models.Paiement
			if err := payRows.Scan(
				&p.ID, &p.AutoEcoleID, &p.AgenceID, &p.DossierID, &p.SessionCaisseID,
				&p.Montant, &p.ModeReglement, &p.ReferenceTransaction, &p.ReferenceRecu,
				&p.SoldeApresPaiement, &p.DatePaiement, &p.EncaisseParID,
				&p.EstAnnule, &p.CreatedAt,
			); err == nil {
				p.CandidatNomComplet = fmt.Sprintf("%s %s", c.Nom, c.Prenoms)
				paiements = append(paiements, p)
			}
		}
		if err := payRows.Err(); err != nil {
			return nil, err
		}
	}

	// Calcul des pourcentages de progression
	progCode := 0
	if f.HeuresTheoriques > 0 {
		progCode = int(math.Min(100, float64(d.HeuresCodeEffectuees)/float64(f.HeuresTheoriques)*100))
	}
	progConduite := 0
	if f.HeuresPratiques > 0 {
		progConduite = int(math.Min(100, float64(d.HeuresConduiteEffectuees)/float64(f.HeuresPratiques)*100))
	}

	// Dernier score de test code réel enregistré
	var dernierScore *int
	var scVal int
	errScore := r.db.QueryRowContext(ctx, `
		SELECT score_test 
		FROM presences_code 
		WHERE candidat_id = $1 AND score_test > 0 
		ORDER BY created_at DESC 
		LIMIT 1
	`, candidatID).Scan(&scVal)
	if errScore == nil {
		dernierScore = &scVal
	}

	return &models.DossierDetailComplet{
		Candidat:               c,
		Dossier:                d,
		Formule:                f,
		Paiements:              paiements,
		ProgressionCodePct:     progCode,
		ProgressionConduitePct: progConduite,
		DernierScoreCode:       dernierScore,
	}, nil
}

// CreerCandidatAvecDossier : Enrôlement atomique en transaction
func (r *GuichetRepository) CreerCandidatAvecDossier(
	ctx context.Context,
	autoEcoleID, agenceID, caissierID, caissierNom string,
	input *models.CreerCandidatInput,
) (*models.Candidat, *models.DossierDetailComplet, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()

	// Récupérer la formule choisie
	var f models.FormulePermis
	err = tx.QueryRowContext(ctx, `
		SELECT id, auto_ecole_id, categorie_permis, intitule, prix_total, acompte_minimum, heures_theoriques, heures_pratiques
		FROM formules_permis WHERE id = $1 AND auto_ecole_id = $2
	`, input.FormuleID, autoEcoleID).Scan(
		&f.ID, &f.AutoEcoleID, &f.CategoriePermis, &f.Intitule, &f.PrixTotal,
		&f.AcompteMinimum, &f.HeuresTheoriques, &f.HeuresPratiques,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("formule de permis invalide ou introuvable: %w", err)
	}

	candID := uuid.New().String()
	dosID := uuid.New().String()

	// Déterminer le numéro séquentiel du dossier
	var countDossiers int
	_ = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM dossiers_inscription WHERE auto_ecole_id = $1`, autoEcoleID).Scan(&countDossiers)
	numDossier := fmt.Sprintf("DOS-%d-%04d", time.Now().Year(), countDossiers+1)

	medStatut := "VALIDE_SANS_RESERVE"
	if input.VisiteMedicaleStatut != nil && *input.VisiteMedicaleStatut != "" {
		medStatut = *input.VisiteMedicaleStatut
	}

	// 1. Insertion candidat
	candQuery := `
		INSERT INTO candidats (
			id, auto_ecole_id, agence_id, nom, prenoms, genre, date_naissance,
			lieu_naissance, telephone_whatsapp, type_piece_id, numero_piece_id,
			visite_medicale_statut, photo_url
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		) RETURNING created_at, updated_at
	`
	var cand models.Candidat
	cand.ID = candID
	cand.AutoEcoleID = autoEcoleID
	cand.AgenceID = agenceID
	cand.Nom = strings.ToUpper(strings.TrimSpace(input.Nom))
	cand.Prenoms = strings.TrimSpace(input.Prenoms)
	cand.Genre = input.Genre
	cand.DateNaissance = input.DateNaissance
	cand.LieuNaissance = input.LieuNaissance
	cand.TelephoneWhatsApp = strings.TrimSpace(input.TelephoneWhatsApp)
	cand.TypePieceID = input.TypePieceID
	cand.NumeroPieceID = strings.TrimSpace(input.NumeroPieceID)
	cand.VisiteMedicaleStatut = medStatut
	cand.PhotoURL = input.PhotoURL

	err = tx.QueryRowContext(ctx, candQuery,
		cand.ID, cand.AutoEcoleID, cand.AgenceID, cand.Nom, cand.Prenoms, cand.Genre,
		cand.DateNaissance, cand.LieuNaissance, cand.TelephoneWhatsApp, cand.TypePieceID,
		cand.NumeroPieceID, cand.VisiteMedicaleStatut, cand.PhotoURL,
	).Scan(&cand.CreatedAt, &cand.UpdatedAt)
	if err != nil {
		return nil, nil, fmt.Errorf("erreur création candidat: %w", err)
	}

	// 2. Insertion dossier
	soldeRestant := math.Max(0, f.PrixTotal-input.AcompteInitial)
	modeFormation := "PRESENTIEL"
	if input.ModeFormationTheorique != nil && *input.ModeFormationTheorique != "" {
		modeFormation = *input.ModeFormationTheorique
	}
	dosQuery := `
		INSERT INTO dossiers_inscription (
			id, auto_ecole_id, agence_id, candidat_id, formule_id, numero_dossier,
			montant_total, remise_accordee, solde_restant, statut_dossier, mode_formation_theorique
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, 0, $8, 'EN_COURS', $9
		) RETURNING to_char(date_inscription, 'YYYY-MM-DD'), created_at, updated_at
	`
	var dos models.DossierInscription
	dos.ID = dosID
	dos.AutoEcoleID = autoEcoleID
	dos.AgenceID = agenceID
	dos.CandidatID = candID
	dos.FormuleID = f.ID
	dos.NumeroDossier = numDossier
	dos.MontantTotal = f.PrixTotal
	dos.SoldeRestant = soldeRestant
	dos.StatutDossier = "EN_COURS"

	err = tx.QueryRowContext(ctx, dosQuery,
		dos.ID, dos.AutoEcoleID, dos.AgenceID, dos.CandidatID, dos.FormuleID,
		dos.NumeroDossier, dos.MontantTotal, dos.SoldeRestant, modeFormation,
	).Scan(&dos.DateInscription, &dos.CreatedAt, &dos.UpdatedAt)
	if err != nil {
		return nil, nil, fmt.Errorf("erreur création dossier: %w", err)
	}

	// 3. Paiement de l'acompte initial si supérieur à 0
	paiements := []models.Paiement{}
	if input.AcompteInitial > 0 {
		payID := uuid.New().String()
		recuNum := fmt.Sprintf("REC-%d-%04d", time.Now().Year(), time.Now().Unix()%10000)

		payQuery := `
			INSERT INTO paiements (
				id, auto_ecole_id, agence_id, dossier_id, montant, mode_reglement,
				reference_recu, solde_apres_paiement, encaisse_par_id
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9
			) RETURNING to_char(date_paiement, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), created_at
		`
		var p models.Paiement
		p.ID = payID
		p.AutoEcoleID = autoEcoleID
		p.AgenceID = agenceID
		p.DossierID = dosID
		p.Montant = input.AcompteInitial
		p.ModeReglement = input.ModeReglement
		p.ReferenceRecu = recuNum
		p.SoldeApresPaiement = soldeRestant
		p.EncaisseParID = &caissierID
		p.CandidatNomComplet = fmt.Sprintf("%s %s", cand.Nom, cand.Prenoms)

		err = tx.QueryRowContext(ctx, payQuery,
			p.ID, p.AutoEcoleID, p.AgenceID, p.DossierID, p.Montant, p.ModeReglement,
			p.ReferenceRecu, p.SoldeApresPaiement, p.EncaisseParID,
		).Scan(&p.DatePaiement, &p.CreatedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("erreur enregistrement paiement acompte: %w", err)
		}

		recordReceiptAndCashTx(ctx, tx, payID, recuNum, agenceID, p.ModeReglement, p.Montant)
		paiements = append(paiements, p)
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}

	complet := &models.DossierDetailComplet{
		Candidat:               cand,
		Dossier:                dos,
		Formule:                f,
		Paiements:              paiements,
		ProgressionCodePct:     0,
		ProgressionConduitePct: 0,
	}

	return &cand, complet, nil
}

// UpdateCandidatPhoto : Met à jour la photo d'identité 4x4 d'un candidat
func (r *GuichetRepository) UpdateCandidatPhoto(ctx context.Context, autoEcoleID, candidatID, photoURL string) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE candidats SET photo_url = $1, updated_at = NOW()
		WHERE id = $2 AND auto_ecole_id = $3
	`, photoURL, candidatID, autoEcoleID)
	if err != nil {
		return fmt.Errorf("erreur mise à jour photo candidat: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("candidat introuvable")
	}
	return nil
}

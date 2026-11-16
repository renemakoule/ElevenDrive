package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/google/uuid"

	"autoschool/backend/internal/models"
)

func recordReceiptAndCashTx(
	ctx context.Context,
	tx *sql.Tx,
	payID, recuNum, agenceID, modeReglement string,
	montant float64,
) string {
	secHash := fmt.Sprintf("SECHASH-%s-%s", payID[:8], recuNum)
	_, _ = tx.ExecContext(ctx, `
		INSERT INTO recus_caisse (id, paiement_id, numero_recu, code_securite_hash)
		VALUES ($1, $2, $3, $4)
	`, uuid.New().String(), payID, recuNum, secHash)

	if modeReglement == "ESPECES" {
		_, _ = tx.ExecContext(ctx, `
			UPDATE sessions_caisse
			SET total_especes_theorique = total_especes_theorique + $1, updated_at = NOW()
			WHERE agence_id = $2 AND statut = 'OUVERTE'
		`, montant, agenceID)
	}
	return secHash
}

// EnregistrerPaiement : Enregistrement tranche de paiement avec trigger SQL de solde
func (r *GuichetRepository) EnregistrerPaiement(
	ctx context.Context,
	autoEcoleID, agenceID, caissierID string,
	input *models.EnregistrerPaiementInput,
) (*models.EnregistrerPaiementResponse, error) {
	if input.Montant <= 0 {
		return nil, fmt.Errorf("le montant doit être supérieur à zéro")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Récupérer le dossier
	var d models.DossierInscription
	err = tx.QueryRowContext(ctx, `
		SELECT id, auto_ecole_id, agence_id, numero_dossier, solde_restant
		FROM dossiers_inscription
		WHERE id = $1 AND auto_ecole_id = $2
	`, input.DossierID, autoEcoleID).Scan(&d.ID, &d.AutoEcoleID, &d.AgenceID, &d.NumeroDossier, &d.SoldeRestant)
	if err != nil {
		return nil, fmt.Errorf("dossier d'inscription introuvable: %w", err)
	}

	if input.Montant > d.SoldeRestant {
		return nil, fmt.Errorf("le montant versé (%.0f FCFA) dépasse le solde restant (%.0f FCFA)", input.Montant, d.SoldeRestant)
	}

	if agenceID == "" {
		agenceID = d.AgenceID
	}

	nouveauSolde := math.Max(0, d.SoldeRestant-input.Montant)
	payID := uuid.New().String()
	recuNum := fmt.Sprintf("REC-%d-%04d", time.Now().Year(), time.Now().Unix()%10000)

	payQuery := `
		INSERT INTO paiements (
			id, auto_ecole_id, agence_id, dossier_id, montant, mode_reglement,
			reference_transaction, reference_recu, solde_apres_paiement, encaisse_par_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		) RETURNING to_char(date_paiement, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
	`
	var datePaiement string
	err = tx.QueryRowContext(ctx, payQuery,
		payID, autoEcoleID, agenceID, d.ID, input.Montant, input.ModeReglement,
		input.ReferenceTransaction, recuNum, nouveauSolde, caissierID,
	).Scan(&datePaiement)
	if err != nil {
		return nil, fmt.Errorf("erreur enregistrement paiement: %w", err)
	}

	// Mettre à jour le solde restant sur le dossier
	_, err = tx.ExecContext(ctx, `
		UPDATE dossiers_inscription
		SET solde_restant = $1, updated_at = NOW()
		WHERE id = $2
	`, nouveauSolde, d.ID)
	if err != nil {
		return nil, fmt.Errorf("erreur mise à jour solde dossier: %w", err)
	}

	secHash := recordReceiptAndCashTx(ctx, tx, payID, recuNum, agenceID, input.ModeReglement, input.Montant)

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.EnregistrerPaiementResponse{
		PaiementID:          payID,
		NumeroRecu:          recuNum,
		CodeSecuriteHash:    secHash,
		MontantVerse:        input.Montant,
		NouveauSoldeRestant: nouveauSolde,
		DatePaiement:        datePaiement,
		NumeroDossier:       d.NumeroDossier,
	}, nil
}

// GetSessionCaisseDetail : Détail session caisse avec ventilation
func (r *GuichetRepository) GetSessionCaisseDetail(ctx context.Context, autoEcoleID, agenceID string) (*models.SessionCaisseDetail, error) {
	// Récupérer la dernière session réelle de l'auto-école (ou de l'agence)
	var sess models.SessionCaisse
	querySess := `
		SELECT s.id, s.agence_id, s.caissier_id, s.date_ouverture, s.date_fermeture, s.fond_de_caisse,
		       s.total_especes_theorique, s.statut, s.observations
		FROM sessions_caisse s
		JOIN agences a ON a.id = s.agence_id
		WHERE a.auto_ecole_id = $1 AND ($2 = '' OR s.agence_id::text = $2)
		ORDER BY s.date_ouverture DESC
		LIMIT 1
	`
	err := r.db.QueryRowContext(ctx, querySess, autoEcoleID, agenceID).Scan(
		&sess.ID, &sess.AgenceID, &sess.CaissierID, &sess.DateOuverture, &sess.DateFermeture,
		&sess.FondDeCaisse, &sess.TotalEspecesTheorique, &sess.Statut, &sess.Observations,
	)

	// Ventilation des encaissements (depuis l'ouverture de session, ou du jour UTC si aucune session ouverte)
	var vent models.VentilationCaisse
	var nbTransactions int
	nowUTC := time.Now().UTC()
	startOfDayUTC := time.Date(nowUTC.Year(), nowUTC.Month(), nowUTC.Day(), 0, 0, 0, 0, time.UTC)
	ventSince := startOfDayUTC
	if err == nil && sess.Statut == "OUVERTE" {
		ventSince = sess.DateOuverture
	} else if err == nil && sess.DateFermeture != nil && sess.DateFermeture.After(startOfDayUTC) {
		ventSince = *sess.DateFermeture
	}

	queryVent := `
		SELECT
			COALESCE(SUM(CASE WHEN mode_reglement = 'ESPECES' THEN montant ELSE 0 END), 0) AS esp,
			COALESCE(SUM(CASE WHEN mode_reglement = 'WAVE' THEN montant ELSE 0 END), 0) AS wave,
			COALESCE(SUM(CASE WHEN mode_reglement = 'ORANGE_MONEY' THEN montant ELSE 0 END), 0) AS om,
			COALESCE(SUM(CASE WHEN mode_reglement = 'MTN_MOMO' THEN montant ELSE 0 END), 0) AS mtn,
			COALESCE(SUM(CASE WHEN mode_reglement = 'MOOV' THEN montant ELSE 0 END), 0) AS moov,
			COALESCE(SUM(CASE WHEN mode_reglement = 'VIREMENT' THEN montant ELSE 0 END), 0) AS vir,
			COALESCE(SUM(montant), 0) AS total,
			COUNT(*) AS nb
		FROM paiements
		WHERE auto_ecole_id = $1
		  AND ($2 = '' OR agence_id::text = $2)
		  AND COALESCE(est_annule, false) = false
		  AND COALESCE(reference_transaction, '') <> 'REPRISE_HISTORIQUE'
		  AND date_paiement >= $3
	`
	if vErr := r.db.QueryRowContext(ctx, queryVent, autoEcoleID, agenceID, ventSince).Scan(
		&vent.Especes, &vent.Wave, &vent.OrangeMoney, &vent.MtnMomo, &vent.Moov,
		&vent.Virement, &vent.TotalEncaisse, &nbTransactions,
	); vErr != nil {
		log.Printf("[GetSessionCaisseDetail] queryVent error: %v", vErr)
	}

	paiementsList := []models.Paiement{}
	if pRows, pErr := r.db.QueryContext(ctx, `
		SELECT p.id, p.auto_ecole_id, p.agence_id, p.dossier_id, p.montant, p.mode_reglement,
		       p.reference_transaction, p.reference_recu, p.solde_apres_paiement,
		       to_char(p.date_paiement, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
		       COALESCE(c.nom || ' ' || c.prenoms, ''), COALESCE(p.est_annule, false)
		FROM paiements p
		LEFT JOIN dossiers_inscription d ON d.id = p.dossier_id
		LEFT JOIN candidats c ON c.id = d.candidat_id
		WHERE p.auto_ecole_id = $1
		  AND ($2 = '' OR p.agence_id::text = $2)
		  AND COALESCE(p.est_annule, false) = false
		  AND COALESCE(p.reference_transaction, '') <> 'REPRISE_HISTORIQUE'
		  AND p.date_paiement >= $3
		ORDER BY p.date_paiement DESC
	`, autoEcoleID, agenceID, ventSince); pErr == nil {
		for pRows.Next() {
			var p models.Paiement
			if pRows.Scan(
				&p.ID, &p.AutoEcoleID, &p.AgenceID, &p.DossierID, &p.Montant, &p.ModeReglement,
				&p.ReferenceTransaction, &p.ReferenceRecu, &p.SoldeApresPaiement,
				&p.DatePaiement, &p.CandidatNomComplet, &p.EstAnnule,
			) == nil {
				paiementsList = append(paiementsList, p)
			}
		}
		pRows.Close()
	} else {
		log.Printf("[GetSessionCaisseDetail] paiementsList error: %v", pErr)
	}

	// Si aucune session n'est encore enregistrée dans sessions_caisse, renvoyer les recettes réelles du jour
	if err != nil {
		return &models.SessionCaisseDetail{
			SessionOuverte:        false,
			ID:                    "",
			AgenceID:              agenceID,
			CaissierID:            "",
			CaissierNom:           "Session du jour",
			DateOuverture:         ventSince.Format(time.RFC3339),
			FondDeCaisse:          0,
			TotalEspecesTheorique: vent.Especes,
			TotalTiroirAttendu:    vent.Especes,
			Statut:                "OUVERTE",
			Ventilation:           vent,
			NombreTransactions:    nbTransactions,
			Paiements:             paiementsList,
		}, nil
	}

	// Nom du caissier
	caissierNom := "Agent Guichet"
	if sess.CaissierID != "" {
		_ = r.db.QueryRowContext(ctx, `SELECT nom || ' ' || prenoms FROM utilisateurs WHERE id = $1`, sess.CaissierID).Scan(&caissierNom)
	}

	return &models.SessionCaisseDetail{
		SessionOuverte:        sess.Statut == "OUVERTE",
		ID:                    sess.ID,
		AgenceID:              sess.AgenceID,
		CaissierID:            sess.CaissierID,
		CaissierNom:           caissierNom,
		DateOuverture:         sess.DateOuverture.Format(time.RFC3339),
		FondDeCaisse:          sess.FondDeCaisse,
		TotalEspecesTheorique: vent.Especes,
		TotalTiroirAttendu:    sess.FondDeCaisse + vent.Especes,
		Statut:                sess.Statut,
		Ventilation:           vent,
		NombreTransactions:    nbTransactions,
		Paiements:             paiementsList,
	}, nil
}

// CloturerSessionCaisse : Clôture définitive et verrouillage du tiroir
func (r *GuichetRepository) CloturerSessionCaisse(
	ctx context.Context,
	autoEcoleID, agenceID, caissierID string,
	input *models.CloturerSessionCaisseInput,
) error {
	query := `
		UPDATE sessions_caisse
		SET total_especes_compte = $1,
		    ecart_caisse = $1 - (fond_de_caisse + total_especes_theorique),
		    statut = 'CLOTUREE',
		    date_fermeture = NOW(),
		    observations = $2,
		    updated_at = NOW()
		WHERE agence_id = $3 AND statut = 'OUVERTE'
	`
	res, err := r.db.ExecContext(ctx, query, input.TotalEspecesCompte, input.Observations, agenceID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		// Insérer directement une clôture enregistrée
		insertQuery := `
			INSERT INTO sessions_caisse (
				id, agence_id, caissier_id, date_fermeture, fond_de_caisse,
				total_especes_theorique, total_especes_compte, ecart_caisse, statut, observations
			) VALUES (
				$1, $2, $3, NOW(), 0, 0, $4, 0, 'CLOTUREE', $5
			)
		`
		_, err = r.db.ExecContext(ctx, insertQuery, uuid.New().String(), agenceID, caissierID, input.TotalEspecesCompte, input.Observations)
		return err
	}
	return nil
}

// AccorderDerogationConduite : Lève le verrou sur la conduite avant le seuil réglementaire (60%)
func (r *GuichetRepository) AccorderDerogationConduite(ctx context.Context, autoEcoleID, dossierID string) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE dossiers_inscription
		SET derogation_conduite_accordee = TRUE, updated_at = NOW()
		WHERE id = $1 AND auto_ecole_id = $2
	`, dossierID, autoEcoleID)
	if err != nil {
		return fmt.Errorf("erreur dérogation conduite: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("dossier d'inscription introuvable pour cette auto-école")
	}
	return nil
}

// ReimprimerRecu : Incrémente le compteur de duplicata et retourne le détail du reçu
func (r *GuichetRepository) ReimprimerRecu(ctx context.Context, autoEcoleID, paiementID string) (*models.RecuImprimable, error) {
	_, _ = r.db.ExecContext(ctx, `
		UPDATE recus_caisse rc
		SET nombre_impressions = nombre_impressions + 1, derniere_impression_at = NOW()
		FROM paiements p
		WHERE rc.paiement_id = p.id AND p.id = $1 AND p.auto_ecole_id = $2
	`, paiementID, autoEcoleID)

	var recu models.RecuImprimable
	err := r.db.QueryRowContext(ctx, `
		SELECT COALESCE(rc.numero_recu, p.reference_recu),
		       COALESCE(rc.code_securite_hash, 'SECHASH-' || SUBSTRING(p.id::text, 1, 8)),
		       COALESCE(rc.nombre_impressions, 1),
		       p.montant, d.solde_restant, p.mode_reglement,
		       to_char(p.date_paiement, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
		       c.nom || ' ' || c.prenoms, d.numero_dossier, f.intitule,
		       COALESCE(u.nom || ' ' || u.prenoms, 'Agent Guichet')
		FROM paiements p
		LEFT JOIN recus_caisse rc ON rc.paiement_id = p.id
		JOIN dossiers_inscription d ON d.id = p.dossier_id
		JOIN candidats c ON c.id = d.candidat_id
		JOIN formules_permis f ON f.id = d.formule_id
		LEFT JOIN utilisateurs u ON u.id = p.encaisse_par_id
		WHERE p.id = $1 AND p.auto_ecole_id = $2
	`, paiementID, autoEcoleID).Scan(
		&recu.NumeroRecu, &recu.CodeSecuriteHash, &recu.NombreImpressions,
		&recu.MontantVerse, &recu.SoldeRestant, &recu.ModeReglement,
		&recu.DatePaiement, &recu.CandidatNom, &recu.NumeroDossier,
		&recu.FormuleNom, &recu.CaissierNom,
	)
	if err != nil {
		return nil, fmt.Errorf("reçu introuvable: %w", err)
	}
	recu.EstDuplicata = recu.NombreImpressions > 1
	return &recu, nil
}

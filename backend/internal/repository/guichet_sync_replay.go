package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"

	"autoschool/backend/internal/models"
)

// ReplaySyncEvents rejoue les opérations hors-ligne dans les vraies tables PostgreSQL de façon idempotente.
func (r *GuichetRepository) ReplaySyncEvents(
	ctx context.Context,
	autoEcoleID, agenceID string,
	events []models.SyncEventPayload,
) (*models.SyncBatchResponse, error) {
	nowISO := time.Now().UTC().Format(time.RFC3339)
	processed, duplicates, errCount := 0, 0, 0

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	for _, ev := range events {
		evID := ev.EventID
		if _, err := uuid.Parse(evID); err != nil {
			evID = uuid.New().String()
		}
		entID := ev.EntityID
		if _, err := uuid.Parse(entID); err != nil {
			entID = uuid.New().String()
		}
		payloadBytes, _ := json.Marshal(ev.Payload)
		if len(payloadBytes) == 0 {
			payloadBytes = []byte("{}")
		}
		localTS := strings.TrimSpace(ev.LocalTimestampUTC)
		if localTS == "" {
			localTS = nowISO
		}

		_, _ = tx.ExecContext(ctx, "SAVEPOINT sp_sync_ev")
		res, errIns := tx.ExecContext(ctx, `
			INSERT INTO journal_sync_events (
				event_id, auto_ecole_id, agence_id, nom_table, type_operation,
				entite_id, donnees_payload, horodatage_local, statut_sync
			) VALUES ($1::uuid, $2::uuid, NULLIF($3, '')::uuid, $4, $5, $6::uuid, $7::jsonb, $8::timestamptz, 'SYNCHRONISE')
			ON CONFLICT (event_id) DO NOTHING
		`, evID, autoEcoleID, agenceID, ev.EntityName, ev.Operation, entID, string(payloadBytes), localTS)
		if errIns != nil {
			_, _ = tx.ExecContext(ctx, "ROLLBACK TO SAVEPOINT sp_sync_ev")
			errCount++
			continue
		}
		if rowsAff, _ := res.RowsAffected(); rowsAff == 0 {
			_, _ = tx.ExecContext(ctx, "RELEASE SAVEPOINT sp_sync_ev")
			duplicates++
			continue
		}

		if errReplay := replayEntityMutation(ctx, tx, autoEcoleID, agenceID, ev); errReplay != nil {
			_, _ = tx.ExecContext(ctx, "ROLLBACK TO SAVEPOINT sp_sync_ev")
			errCount++
			continue
		}
		_, _ = tx.ExecContext(ctx, "RELEASE SAVEPOINT sp_sync_ev")
		processed++
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.SyncBatchResponse{
		Statut:            "SUCCES",
		EvenementsTraites: processed,
		DoublonsIgnores:   duplicates,
		Erreurs:           errCount,
		SyncedAt:          nowISO,
		Message:           fmt.Sprintf("%d transaction(s) synchronisée(s) et appliquée(s) sur Neon PostgreSQL", processed),
	}, nil
}

func replayEntityMutation(ctx context.Context, tx *sql.Tx, autoEcoleID, agenceID string, ev models.SyncEventPayload) error {
	p := ev.Payload
	op := strings.ToUpper(strings.TrimSpace(ev.Operation))
	getStr := func(k string) string {
		if v, ok := p[k]; ok && v != nil {
			return strings.TrimSpace(fmt.Sprintf("%v", v))
		}
		return ""
	}
	getFloat := func(k string) float64 {
		if v, ok := p[k]; ok && v != nil {
			if f, ok := v.(float64); ok {
				return f
			}
		}
		return 0
	}
	getBool := func(k string) bool {
		if v, ok := p[k]; ok && v != nil {
			if b, ok := v.(bool); ok {
				return b
			}
		}
		return false
	}

	switch ev.EntityName {
	case "candidats":
		candID := getStr("candidat_id")
		dosID := getStr("dossier_id")
		formuleID := getStr("formule_id")
		if candID == "" || dosID == "" || formuleID == "" {
			return nil
		}
		nom := strings.ToUpper(getStr("nom"))
		prenoms := getStr("prenoms")
		genre := getStr("genre")
		dateNaiss := getStr("date_naissance")
		lieuNaiss := getStr("lieu_naissance")
		typePiece := normalizeTypePieceID(getStr("type_piece_id"))
		numPiece := getStr("numero_piece_id")
		visite := getStr("visite_medicale_statut")
		if nom == "" || prenoms == "" || genre == "" || dateNaiss == "" || lieuNaiss == "" {
			return fmt.Errorf("données candidat incomplètes dans l'événement %s", candID)
		}
		_, err := tx.ExecContext(ctx, `
			INSERT INTO candidats (
				id, auto_ecole_id, agence_id, nom, prenoms, genre, date_naissance,
				lieu_naissance, telephone_whatsapp, type_piece_id, numero_piece_id, visite_medicale_statut
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NULLIF($12, ''))
			ON CONFLICT (id) DO NOTHING
		`, candID, autoEcoleID, agenceID, nom, prenoms, genre, dateNaiss, lieuNaiss, getStr("telephone_whatsapp"), typePiece, numPiece, visite)
		if err != nil {
			return err
		}
		numDos := getStr("numero_dossier")
		montantTotal := getFloat("montant_total")
		soldeRestant := getFloat("solde_restant")
		modeForm := getStr("mode_formation_theorique")
		_, err = tx.ExecContext(ctx, `
			INSERT INTO dossiers_inscription (
				id, auto_ecole_id, agence_id, candidat_id, formule_id, numero_dossier,
				montant_total, remise_accordee, solde_restant, statut_dossier, mode_formation_theorique
			) VALUES ($1, $2, $3, $4, $5, $6, $7, 0, $8, 'EN_COURS', COALESCE(NULLIF($9, ''), 'PRESENTIEL'))
			ON CONFLICT (id) DO NOTHING
		`, dosID, autoEcoleID, agenceID, candID, formuleID, numDos, montantTotal, soldeRestant, modeForm)
		return err

	case "dossiers_inscription":
		dosID := getStr("dossier_id")
		if dosID == "" {
			dosID = ev.EntityID
		}
		if dosID == "" {
			return nil
		}
		if getBool("derogation_conduite_accordee") {
			_, err := tx.ExecContext(ctx, `
				UPDATE dossiers_inscription SET derogation_conduite_accordee = true, updated_at = NOW()
				WHERE id = $1 AND auto_ecole_id = $2
			`, dosID, autoEcoleID)
			return err
		}
		return nil

	case "paiements":
		payID := getStr("paiement_id")
		dosID := getStr("dossier_id")
		montant := getFloat("montant")
		if payID == "" || dosID == "" || montant <= 0 {
			return nil
		}
		recu := getStr("reference_recu")
		modeReg := getStr("mode_reglement")
		var currentSolde float64
		_ = tx.QueryRowContext(ctx, `SELECT solde_restant FROM dossiers_inscription WHERE id = $1 AND auto_ecole_id = $2`, dosID, autoEcoleID).Scan(&currentSolde)
		newSolde := math.Max(0, currentSolde-montant)

		res, err := tx.ExecContext(ctx, `
			INSERT INTO paiements (
				id, auto_ecole_id, agence_id, dossier_id, montant, mode_reglement,
				reference_recu, solde_apres_paiement
			) VALUES ($1, $2, $3, $4, $5, NULLIF($6, ''), $7, $8)
			ON CONFLICT (id) DO NOTHING
		`, payID, autoEcoleID, agenceID, dosID, montant, modeReg, recu, newSolde)
		if err != nil {
			return err
		}
		if aff, _ := res.RowsAffected(); aff > 0 {
			_, _ = tx.ExecContext(ctx, `UPDATE dossiers_inscription SET solde_restant = $1, updated_at = NOW() WHERE id = $2`, newSolde, dosID)
		}
		return nil

	case "formules_permis":
		id := getStr("id")
		if id == "" {
			id = ev.EntityID
		}
		if id == "" {
			return nil
		}
		if op == "DELETE" {
			_, err := tx.ExecContext(ctx, `DELETE FROM formules_permis WHERE id = $1 AND auto_ecole_id = $2`, id, autoEcoleID)
			return err
		}
		if op == "UPDATE" {
			if _, ok := p["est_actif"]; ok && !getBool("est_actif") {
				_, err := tx.ExecContext(ctx, `UPDATE formules_permis SET est_actif = false, updated_at = NOW() WHERE id = $1 AND auto_ecole_id = $2`, id, autoEcoleID)
				return err
			}
			nouveauPrix := getFloat("prix_total")
			if nouveauPrix > 0 {
				if _, hasAc := p["acompte_minimum"]; hasAc {
					_, err := tx.ExecContext(ctx, `UPDATE formules_permis SET prix_total = $1, acompte_minimum = $2, updated_at = NOW() WHERE id = $3 AND auto_ecole_id = $4`,
						nouveauPrix, getFloat("acompte_minimum"), id, autoEcoleID)
					return err
				}
				_, err := tx.ExecContext(ctx, `UPDATE formules_permis SET prix_total = $1, updated_at = NOW() WHERE id = $2 AND auto_ecole_id = $3`,
					nouveauPrix, id, autoEcoleID)
				return err
			}
			return nil
		}
		_, err := tx.ExecContext(ctx, `
			INSERT INTO formules_permis (
				id, auto_ecole_id, categorie_permis, intitule, description,
				heures_theoriques, heures_pratiques, prix_total, acompte_minimum, est_actif
			) VALUES ($1, $2, $3, $4, NULLIF($5, ''), $6, $7, $8, $9, true)
			ON CONFLICT (id) DO NOTHING
		`, id, autoEcoleID, getStr("categorie_permis"), getStr("intitule"), getStr("description"),
			int(getFloat("heures_theoriques")), int(getFloat("heures_pratiques")), getFloat("prix_total"), getFloat("acompte_minimum"))
		return err

	case "vehicules":
		id := getStr("id")
		if id == "" {
			id = ev.EntityID
		}
		if id == "" {
			return nil
		}
		if op == "UPDATE" {
			statut := getStr("statut")
			if statut == "" {
				return nil
			}
			_, err := tx.ExecContext(ctx, `UPDATE vehicules SET statut = $1, updated_at = NOW() WHERE id = $2`, statut, id)
			return err
		}
		_, err := tx.ExecContext(ctx, `
			INSERT INTO vehicules (
				id, agence_id, immatriculation, marque_modele, type_boite,
				kilometrage_actuel, date_expiration_assurance, date_visite_technique, statut
			) VALUES ($1, $2, $3, $4, $5, $6, $7::date, $8::date, $9)
			ON CONFLICT (id) DO NOTHING
		`, id, agenceID, getStr("immatriculation"), getStr("marque_modele"), getStr("type_boite"),
			int(getFloat("kilometrage_actuel")), getStr("date_expiration_assurance"), getStr("date_visite_technique"), getStr("statut"))
		return err

	case "sessions_caisse":
		id := getStr("id")
		if id == "" {
			id = ev.EntityID
		}
		if id == "" {
			return nil
		}
		if op == "UPDATE" {
			_, err := tx.ExecContext(ctx, `
				UPDATE sessions_caisse
				SET total_especes_compte = $1,
				    ecart_caisse = $1 - (fond_de_caisse + total_especes_theorique),
				    statut = 'CLOTUREE',
				    date_fermeture = COALESCE(NULLIF($2, '')::timestamptz, NOW()),
				    observations = NULLIF($3, '')
				WHERE id = $4 AND agence_id = $5
			`, getFloat("total_especes_compte"), getStr("date_fermeture"), getStr("observations"), id, agenceID)
			return err
		}
		caissierID := getStr("caissier_id")
		if caissierID == "" {
			return nil
		}
		_, err := tx.ExecContext(ctx, `
			INSERT INTO sessions_caisse (id, agence_id, caissier_id, date_ouverture, fond_de_caisse, total_especes_theorique, statut)
			VALUES ($1, $2, $3, COALESCE(NULLIF($4, '')::timestamptz, NOW()), $5, 0, 'OUVERTE')
			ON CONFLICT (id) DO NOTHING
		`, id, agenceID, caissierID, getStr("date_ouverture"), getFloat("fond_de_caisse"))
		return err

	case "seances_code":
		id := getStr("seance_id")
		if id == "" {
			id = ev.EntityID
		}
		if id == "" {
			return nil
		}
		if op == "UPDATE" {
			_, err := tx.ExecContext(ctx, `
				UPDATE seances_code
				SET date_seance = $1::date, heure_debut = $2, heure_fin = $3, statut_seance = 'PLANIFIEE'
				WHERE id = $4
			`, getStr("date_seance"), getStr("heure_debut"), getStr("heure_fin"), id)
			return err
		}
		theme := getStr("theme_traite")
		if theme == "" {
			theme = getStr("theme")
		}
		formateur := getStr("formateur_nom")
		if formateur == "" {
			formateur = getStr("formateur")
		}
		if theme == "" || formateur == "" {
			return nil
		}
		_, err := tx.ExecContext(ctx, `
			INSERT INTO seances_code (
				id, auto_ecole_id, agence_id, date_seance, heure_debut, heure_fin,
				theme_traite, formateur_nom, moniteur_id, type_seance, lien_visio, mot_de_passe_visio, statut_seance, statut
			) VALUES (
				$1, $2, $3, COALESCE(NULLIF($4, '')::date, CURRENT_DATE), $5, $6,
				$7, $8, NULLIF($9, '')::uuid, COALESCE(NULLIF($10, ''), 'PRESENTIEL'),
				NULLIF($11, ''), NULLIF($12, ''), COALESCE(NULLIF($13, ''), 'PLANIFIEE'), COALESCE(NULLIF($13, ''), 'PLANIFIEE')
			)
			ON CONFLICT (id) DO NOTHING
		`, id, autoEcoleID, agenceID, getStr("date_seance"), getStr("heure_debut"), getStr("heure_fin"),
			theme, formateur, getStr("moniteur_id"), getStr("type_seance"), getStr("lien_visio"), getStr("mot_de_passe_visio"), getStr("statut"))
		return err

	case "candidats_examen":
		id := getStr("id")
		if id == "" {
			id = ev.EntityID
		}
		sessID := getStr("session_examen_id")
		dosID := getStr("dossier_id")
		candID := getStr("candidat_id")
		if id == "" || sessID == "" || dosID == "" || candID == "" {
			return nil
		}
		_, err := tx.ExecContext(ctx, `
			INSERT INTO candidats_examen (
				id, session_examen_id, dossier_id, candidat_id, numero_table,
				resultat_code, resultat_conduite, resultat_final
			) VALUES ($1, $2, $3, $4, $5, 'EN_ATTENTE', 'EN_ATTENTE', 'EN_ATTENTE')
			ON CONFLICT (id) DO NOTHING
		`, id, sessID, dosID, candID, getStr("numero_table"))
		return err
	}
	return nil
}

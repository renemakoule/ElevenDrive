package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"autoschool/backend/internal/models"
)

// GetSessionExamenActive : Récupère la session d'examen active avec les candidats éligibles et non-éligibles
func (r *GuichetRepository) GetSessionExamenActive(ctx context.Context, autoEcoleID string) (*models.SessionExamenActiveResponse, error) {
	var sess models.SessionExamen
	var inspecteurNom, observations sql.NullString
	querySess := `
		SELECT id, auto_ecole_id, agence_id, code_session, to_char(date_examen, 'YYYY-MM-DD'),
		       centre_examen, categorie_permis, statut_session, inspecteur_nom, observations,
		       created_at, updated_at
		FROM sessions_examen
		WHERE auto_ecole_id = $1 AND statut_session IN ('EN_PREPARATION', 'BORDEREAU_DEPOSE', 'RESULTATS_PUBLIES')
		ORDER BY date_examen DESC
		LIMIT 1
	`
	var hasSession = true
	err := r.db.QueryRowContext(ctx, querySess, autoEcoleID).Scan(
		&sess.ID, &sess.AutoEcoleID, &sess.AgenceID, &sess.CodeSession, &sess.DateExamen,
		&sess.CentreExamen, &sess.CategoriePermis, &sess.StatutSession,
		&inspecteurNom, &observations, &sess.CreatedAt, &sess.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			hasSession = false
		} else {
			return nil, fmt.Errorf("erreur récupération session examen: %w", err)
		}
	}
	if inspecteurNom.Valid {
		sess.InspecteurNom = &inspecteurNom.String
	}
	if observations.Valid {
		sess.Observations = &observations.String
	}

	queryCands := `
		SELECT 
			d.id AS dossier_id, d.numero_dossier, d.solde_restant,
			d.heures_code_effectuees, d.heures_conduite_effectuees,
			d.derogation_conduite_accordee,
			f.heures_theoriques, f.heures_pratiques, f.categorie_permis,
			c.id AS candidat_id, c.nom, c.prenoms, c.genre,
			to_char(c.date_naissance, 'YYYY-MM-DD'), c.lieu_naissance,
			c.telephone_whatsapp, c.type_piece_id, c.numero_piece_id,
			c.photo_url, c.visite_medicale_statut,
			ce.id AS candidat_examen_id,
			COALESCE(ce.numero_table, '') AS numero_table,
			COALESCE(ce.resultat_code, 'EN_ATTENTE') AS resultat_code,
			COALESCE(ce.resultat_conduite, 'EN_ATTENTE') AS resultat_conduite,
			COALESCE(ce.resultat_final, 'EN_ATTENTE') AS resultat_final,
			to_char(ce.date_notification_envoyee, 'YYYY-MM-DD"T"HH24:MI:SS') AS date_notification_envoyee
		FROM dossiers_inscription d
		JOIN candidats c ON c.id = d.candidat_id
		JOIN formules_permis f ON f.id = d.formule_id
		LEFT JOIN candidats_examen ce ON ce.dossier_id = d.id AND ce.session_examen_id = $2
		WHERE d.auto_ecole_id = $1
		ORDER BY d.created_at ASC
	`
	candsSessID := sess.ID
	if !hasSession {
		candsSessID = "00000000-0000-0000-0000-000000000000"
	}

	rows, err := r.db.QueryContext(ctx, queryCands, autoEcoleID, candsSessID)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération candidats examen: %w", err)
	}
	defer rows.Close()

	var eligibles []models.CandidatExamenItem
	var ineligibles []models.DossierNonEligible
	tableIndex := 1

	for rows.Next() {
		var dossierID, numeroDossier string
		var soldeRestant float64
		var hCodeEff, hCondEff int
		var derogConduite bool
		var hCodeReq, hCondReq int
		var catPermis string
		var candidatID, nom, prenoms, genre, dateNais, lieuNais, tel, typePiece, numPiece, visiteStat string
		var photoURL sql.NullString
		var ceID, ceTable, resCode, resCond, resFinal, notifDate sql.NullString

		if err := rows.Scan(
			&dossierID, &numeroDossier, &soldeRestant,
			&hCodeEff, &hCondEff, &derogConduite,
			&hCodeReq, &hCondReq, &catPermis,
			&candidatID, &nom, &prenoms, &genre,
			&dateNais, &lieuNais, &tel, &typePiece, &numPiece,
			&photoURL, &visiteStat,
			&ceID, &ceTable, &resCode, &resCond, &resFinal, &notifDate,
		); err != nil {
			return nil, err
		}

		soldeOk := soldeRestant <= 0
		codeOk := hCodeEff >= hCodeReq
		conduiteOk := hCondEff >= hCondReq || derogConduite
		visiteOk := visiteStat == "VALIDE_SANS_RESERVE" || visiteStat == "VALIDE_AVEC_VERRES"

		isForced := ceID.Valid && ceID.String != ""
		isEligible := isForced || (soldeOk && codeOk && conduiteOk && visiteOk)

		if isEligible {
			item := models.CandidatExamenItem{
				ID:                       dossierID,
				SessionExamenID:          sess.ID,
				DossierID:                dossierID,
				NumeroDossier:            numeroDossier,
				CandidatID:               candidatID,
				Nom:                      nom,
				Prenoms:                  prenoms,
				Genre:                    genre,
				DateNaissance:            dateNais,
				LieuNaissance:            lieuNais,
				TelephoneWhatsApp:        tel,
				TypePieceID:              typePiece,
				NumeroPieceID:            numPiece,
				NumeroTable:              fmt.Sprintf("TABLE-%03d", tableIndex),
				ResultatCode:             "EN_ATTENTE",
				ResultatConduite:         "EN_ATTENTE",
				ResultatFinal:            "EN_ATTENTE",
				SoldeRestant:             soldeRestant,
				VisiteMedicaleStatut:     visiteStat,
				HeuresCodeEffectuees:     hCodeEff,
				HeuresCodeRequises:       hCodeReq,
				HeuresConduiteEffectuees: hCondEff,
				HeuresConduiteRequises:   hCondReq,
				EstEligible:              true,
				EstDerogation:            isForced && !(soldeOk && codeOk && conduiteOk && visiteOk),
			}
			tableIndex++
			if photoURL.Valid {
				item.PhotoURL = &photoURL.String
			}
			if ceID.Valid && ceID.String != "" {
				item.ID = ceID.String
			}
			if ceTable.Valid && ceTable.String != "" {
				item.NumeroTable = ceTable.String
			}
			if resCode.Valid && resCode.String != "" {
				item.ResultatCode = resCode.String
			}
			if resCond.Valid && resCond.String != "" {
				item.ResultatConduite = resCond.String
			}
			if resFinal.Valid && resFinal.String != "" {
				item.ResultatFinal = resFinal.String
			}
			if notifDate.Valid && notifDate.String != "" {
				item.DateNotificationEnvoyee = &notifDate.String
			}
			eligibles = append(eligibles, item)
		} else {
			var motifs []string
			if !soldeOk {
				motifs = append(motifs, fmt.Sprintf("Solde impayé : %.0f F CFA", soldeRestant))
			}
			if !codeOk {
				motifs = append(motifs, fmt.Sprintf("Code : %d/%d h", hCodeEff, hCodeReq))
			}
			if !conduiteOk {
				motifs = append(motifs, fmt.Sprintf("Conduite : %d/%d h", hCondEff, hCondReq))
			}
			if !visiteOk {
				motifs = append(motifs, "Visite médicale non validée")
			}
			ineligibles = append(ineligibles, models.DossierNonEligible{
				NumeroDossier: numeroDossier,
				Nom:           nom + " " + prenoms,
				Motifs:        motifs,
				DossierID:     &dossierID,
				CandidatID:    &candidatID,
			})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if eligibles == nil {
		eligibles = []models.CandidatExamenItem{}
	}
	if ineligibles == nil {
		ineligibles = []models.DossierNonEligible{}
	}

	var respSession *models.SessionExamen
	if hasSession {
		sess.TotalCandidats = len(eligibles)
		for _, c := range eligibles {
			switch c.ResultatFinal {
			case "ADMIS":
				sess.TotalAdmis++
			case "AJOURNE":
				sess.TotalAjournes++
			}
		}
		respSession = &sess
	}

	return &models.SessionExamenActiveResponse{
		Session:     respSession,
		Candidats:   eligibles,
		Ineligibles: ineligibles,
	}, nil
}

// CreerOuModifierSessionExamen : Planification ou modification explicite d'une session (aucune valeur par défaut)
func (r *GuichetRepository) CreerOuModifierSessionExamen(
	ctx context.Context,
	autoEcoleID string,
	input models.CreerOuModifierSessionExamenInput,
) (*models.SessionExamen, error) {
	var sess models.SessionExamen
	var inspecteurNom, observations sql.NullString

	if input.ID != nil && *input.ID != "" {
		queryUpdate := `
			UPDATE sessions_examen
			SET code_session = $3,
			    date_examen = $4::date,
			    centre_examen = $5,
			    categorie_permis = $6,
			    inspecteur_nom = $7,
			    observations = $8,
			    statut_session = COALESCE($9, statut_session),
			    updated_at = NOW()
			WHERE id = $1 AND auto_ecole_id = $2
			RETURNING id, auto_ecole_id, agence_id, code_session, to_char(date_examen, 'YYYY-MM-DD'),
			          centre_examen, categorie_permis, statut_session, inspecteur_nom, observations,
			          created_at, updated_at
		`
		err := r.db.QueryRowContext(
			ctx, queryUpdate,
			*input.ID, autoEcoleID, input.CodeSession, input.DateExamen,
			input.CentreExamen, input.CategoriePermis, input.InspecteurNom, input.Observations, input.StatutSession,
		).Scan(
			&sess.ID, &sess.AutoEcoleID, &sess.AgenceID, &sess.CodeSession, &sess.DateExamen,
			&sess.CentreExamen, &sess.CategoriePermis, &sess.StatutSession,
			&inspecteurNom, &observations, &sess.CreatedAt, &sess.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erreur modification session examen: %w", err)
		}
	} else {
		var agenceID string
		if input.AgenceID != nil && *input.AgenceID != "" {
			agenceID = *input.AgenceID
		} else {
			_ = r.db.QueryRowContext(ctx, `SELECT id FROM agences WHERE auto_ecole_id = $1 ORDER BY est_siege DESC LIMIT 1`, autoEcoleID).Scan(&agenceID)
		}

		sessID := uuid.New().String()
		statut := "EN_PREPARATION"
		if input.StatutSession != nil && *input.StatutSession != "" {
			statut = *input.StatutSession
		}

		queryInsert := `
			INSERT INTO sessions_examen (
				id, auto_ecole_id, agence_id, code_session, date_examen,
				centre_examen, categorie_permis, statut_session, inspecteur_nom, observations
			) VALUES (
				$1, $2, $3, $4, $5::date,
				$6, $7, $8, $9, $10
			)
			RETURNING id, auto_ecole_id, agence_id, code_session, to_char(date_examen, 'YYYY-MM-DD'),
			          centre_examen, categorie_permis, statut_session, inspecteur_nom, observations,
			          created_at, updated_at
		`
		err := r.db.QueryRowContext(
			ctx, queryInsert,
			sessID, autoEcoleID, agenceID, input.CodeSession, input.DateExamen,
			input.CentreExamen, input.CategoriePermis, statut, input.InspecteurNom, input.Observations,
		).Scan(
			&sess.ID, &sess.AutoEcoleID, &sess.AgenceID, &sess.CodeSession, &sess.DateExamen,
			&sess.CentreExamen, &sess.CategoriePermis, &sess.StatutSession,
			&inspecteurNom, &observations, &sess.CreatedAt, &sess.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erreur création session examen: %w", err)
		}
	}

	if inspecteurNom.Valid {
		sess.InspecteurNom = &inspecteurNom.String
	}
	if observations.Valid {
		sess.Observations = &observations.String
	}

	return &sess, nil
}

// SupprimerSessionExamen : Annule et supprime une session d'examen
func (r *GuichetRepository) SupprimerSessionExamen(ctx context.Context, autoEcoleID string, sessionID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM sessions_examen WHERE id = $1 AND auto_ecole_id = $2`, sessionID, autoEcoleID)
	if err != nil {
		return fmt.Errorf("erreur suppression session examen: %w", err)
	}
	return nil
}

// EnregistrerResultatsExamen : Enregistre les PV et met à jour le statut des dossiers
func (r *GuichetRepository) EnregistrerResultatsExamen(
	ctx context.Context,
	autoEcoleID string,
	input *models.EnregistrerResultatsExamenInput,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("erreur démarrage transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		UPDATE sessions_examen
		SET statut_session = 'RESULTATS_PUBLIES', updated_at = NOW()
		WHERE id = $1 AND auto_ecole_id = $2
	`, input.SessionExamenID, autoEcoleID)
	if err != nil {
		return fmt.Errorf("erreur mise à jour session: %w", err)
	}

	for idx, res := range input.Resultats {
		defaultTableNum := fmt.Sprintf("TABLE-%03d", idx+1)
		upsertQuery := `
			INSERT INTO candidats_examen (
				id, session_examen_id, dossier_id, candidat_id,
				numero_table, resultat_code, resultat_conduite, resultat_final
			)
			SELECT 
				COALESCE(ce.id, gen_random_uuid()), $1, d.id, d.candidat_id,
				COALESCE(NULLIF(ce.numero_table, ''), $2), $3, $4, $5
			FROM dossiers_inscription d
			LEFT JOIN candidats_examen ce ON ce.dossier_id = d.id AND ce.session_examen_id = $1
			WHERE d.auto_ecole_id = $6 AND (d.id::text = $7 OR ce.id::text = $7)
			ON CONFLICT (session_examen_id, candidat_id) DO UPDATE SET
				resultat_code = EXCLUDED.resultat_code,
				resultat_conduite = EXCLUDED.resultat_conduite,
				resultat_final = EXCLUDED.resultat_final
			RETURNING dossier_id
		`
		var matchedDossierID string
		errUpsert := tx.QueryRowContext(
			ctx, upsertQuery,
			input.SessionExamenID, defaultTableNum,
			res.ResultatCode, res.ResultatConduite, res.ResultatFinal,
			autoEcoleID, res.CandidatExamenID,
		).Scan(&matchedDossierID)
		if errUpsert != nil {
			continue
		}

		switch res.ResultatFinal {
		case "ADMIS":
			_, _ = tx.ExecContext(ctx, `UPDATE dossiers_inscription SET statut_dossier = 'ADMIS', updated_at = NOW() WHERE id = $1 AND auto_ecole_id = $2`, matchedDossierID, autoEcoleID)
		case "AJOURNE":
			_, _ = tx.ExecContext(ctx, `UPDATE dossiers_inscription SET statut_dossier = 'AJOURNE', updated_at = NOW() WHERE id = $1 AND auto_ecole_id = $2`, matchedDossierID, autoEcoleID)
		}
	}

	return tx.Commit()
}

// ForcerCandidatExamen : Dérogation Direction manuelle pour inscrire un candidat inéligible au bordereau d'État
func (r *GuichetRepository) ForcerCandidatExamen(
	ctx context.Context,
	autoEcoleID string,
	input models.ForcerCandidatExamenInput,
) (*models.CandidatExamenItem, error) {
	var candID string
	err := r.db.QueryRowContext(ctx, `
		SELECT candidat_id FROM dossiers_inscription WHERE id = $1 AND auto_ecole_id = $2
	`, input.DossierID, autoEcoleID).Scan(&candID)
	if err != nil {
		return nil, fmt.Errorf("dossier d'inscription introuvable : %w", err)
	}

	var count int
	_ = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM candidats_examen WHERE session_examen_id = $1
	`, input.SessionExamenID).Scan(&count)
	numeroTable := fmt.Sprintf("TABLE-%03d", count+1)

	candidatExamenID := uuid.New().String()
	query := `
		INSERT INTO candidats_examen (
			id, session_examen_id, dossier_id, candidat_id,
			numero_table, resultat_code, resultat_conduite, resultat_final
		) VALUES (
			$1, $2, $3, $4, $5, 'EN_ATTENTE', 'EN_ATTENTE', 'EN_ATTENTE'
		)
		ON CONFLICT (session_examen_id, candidat_id) DO UPDATE SET
			dossier_id = EXCLUDED.dossier_id,
			numero_table = COALESCE(candidats_examen.numero_table, EXCLUDED.numero_table)
		RETURNING id, numero_table
	`
	err = r.db.QueryRowContext(
		ctx, query,
		candidatExamenID, input.SessionExamenID, input.DossierID, candID, numeroTable,
	).Scan(&candidatExamenID, &numeroTable)
	if err != nil {
		return nil, fmt.Errorf("impossible de forcer l'inscription à l'examen : %w", err)
	}

	return &models.CandidatExamenItem{
		ID:               candidatExamenID,
		SessionExamenID:  input.SessionExamenID,
		DossierID:        input.DossierID,
		CandidatID:       candID,
		NumeroTable:      numeroTable,
		ResultatCode:     "EN_ATTENTE",
		ResultatConduite: "EN_ATTENTE",
		ResultatFinal:    "EN_ATTENTE",
		EstEligible:      true,
		EstDerogation:    true,
	}, nil
}

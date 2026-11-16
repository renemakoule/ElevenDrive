package repository

import (
	"context"

	"autoschool/backend/internal/database"
	"autoschool/backend/internal/models"
)

// GetTableauGerant : Créances, impayés et assidus au code
func (r *GuichetRepository) GetTableauGerant(ctx context.Context, autoEcoleID string) (*models.TableauGerantResponse, error) {
	var res *models.TableauGerantResponse
	err := database.Retry(ctx, func() error {
		return database.Safe(func() error {
			var e error
			res, e = r.getTableauGerantOnce(ctx, autoEcoleID)
			return e
		})
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *GuichetRepository) getTableauGerantOnce(ctx context.Context, autoEcoleID string) (*models.TableauGerantResponse, error) {
	queryImpayes := `
		SELECT d.id, d.numero_dossier, c.nom || ' ' || c.prenoms, c.telephone_whatsapp,
		       f.intitule, d.montant_total, d.solde_restant,
		       (d.montant_total - d.solde_restant) AS deja_paye,
		       CASE WHEN d.montant_total > 0 THEN ROUND(((d.montant_total - d.solde_restant) / d.montant_total) * 100) ELSE 0 END AS taux,
		       d.heures_code_effectuees, f.heures_theoriques, to_char(d.date_inscription, 'YYYY-MM-DD')
		FROM dossiers_inscription d
		JOIN candidats c ON c.id = d.candidat_id
		JOIN formules_permis f ON f.id = d.formule_id
		WHERE d.auto_ecole_id = $1 AND d.solde_restant > 0
		ORDER BY d.solde_restant DESC
	`
	rows, err := r.db.QueryContext(ctx, queryImpayes, autoEcoleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var impayes []models.ImpayeItem
	totalCreances := 0.0

	for rows.Next() {
		var item models.ImpayeItem
		if err := rows.Scan(
			&item.DossierID, &item.NumeroDossier, &item.CandidatNom, &item.TelephoneWhatsApp,
			&item.Formule, &item.MontantTotal, &item.SoldeRestant, &item.DejaPaye,
			&item.TauxReglementPct, &item.HeuresCode, &item.HeuresCodeRequises, &item.DateInscription,
		); err == nil {
			totalCreances += item.SoldeRestant
			impayes = append(impayes, item)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var assidus []models.AssiduItem
	queryAssidus := `
		SELECT d.id, d.numero_dossier, c.nom || ' ' || c.prenoms, c.telephone_whatsapp,
		       f.intitule, d.heures_code_effectuees, f.heures_theoriques,
		       (d.solde_restant <= 0 OR d.derogation_conduite_accordee) AS conduite_debloquee_financierement,
		       CASE WHEN d.montant_total > 0 THEN ROUND(((d.montant_total - d.solde_restant) / d.montant_total) * 100) ELSE 0 END AS taux,
		       d.solde_restant
		FROM dossiers_inscription d
		JOIN candidats c ON c.id = d.candidat_id
		JOIN formules_permis f ON f.id = d.formule_id
		WHERE d.auto_ecole_id = $1 AND d.statut_dossier NOT IN ('ARCHIVE', 'ABANDONNE')
		  AND d.heures_code_effectuees >= (f.heures_theoriques * 0.7)
		ORDER BY d.heures_code_effectuees DESC
		LIMIT 15
	`
	rowsAssidus, err := r.db.QueryContext(ctx, queryAssidus, autoEcoleID)
	if err == nil {
		defer rowsAssidus.Close()
		for rowsAssidus.Next() {
			var a models.AssiduItem
			if err := rowsAssidus.Scan(
				&a.DossierID, &a.NumeroDossier, &a.CandidatNom, &a.TelephoneWhatsApp,
				&a.Formule, &a.HeuresCodeEffectuees, &a.HeuresTheoriques,
				&a.ConduiteDebloqueeFinancierement, &a.TauxFinancierPct, &a.SoldeRestant,
			); err == nil {
				assidus = append(assidus, a)
			}
		}
		if err := rowsAssidus.Err(); err != nil {
			return nil, err
		}
	}
	if assidus == nil {
		assidus = []models.AssiduItem{}
	}

	// Ventilation réelle par mode de paiement sur l'ensemble de l'auto-école
	var vent models.VentilationCaisse
	queryVent := `
		SELECT
			COALESCE(SUM(CASE WHEN mode_reglement = 'ESPECES' THEN montant ELSE 0 END), 0) AS esp,
			COALESCE(SUM(CASE WHEN mode_reglement = 'WAVE' THEN montant ELSE 0 END), 0) AS wave,
			COALESCE(SUM(CASE WHEN mode_reglement = 'ORANGE_MONEY' THEN montant ELSE 0 END), 0) AS om,
			COALESCE(SUM(CASE WHEN mode_reglement = 'MTN_MOMO' THEN montant ELSE 0 END), 0) AS mtn,
			COALESCE(SUM(CASE WHEN mode_reglement = 'MOOV' THEN montant ELSE 0 END), 0) AS moov,
			COALESCE(SUM(CASE WHEN mode_reglement = 'VIREMENT' THEN montant ELSE 0 END), 0) AS vir,
			COALESCE(SUM(montant), 0) AS total
		FROM paiements
		WHERE auto_ecole_id = $1
	`
	_ = r.db.QueryRowContext(ctx, queryVent, autoEcoleID).Scan(
		&vent.Especes, &vent.Wave, &vent.OrangeMoney, &vent.MtnMomo, &vent.Moov,
		&vent.Virement, &vent.TotalEncaisse,
	)

	// Évolution mensuelle réelle des encaissements (6 derniers mois)
	var recettesMensuelles []models.RecetteMensuelle
	queryMois := `
		SELECT to_char(date_trunc('month', date_paiement), 'YYYY-MM') AS mois,
		       COALESCE(SUM(montant), 0) AS total
		FROM paiements
		WHERE auto_ecole_id = $1 AND date_paiement >= (CURRENT_DATE - INTERVAL '6 months')
		GROUP BY date_trunc('month', date_paiement)
		ORDER BY date_trunc('month', date_paiement) ASC
	`
	rowsMois, errMois := r.db.QueryContext(ctx, queryMois, autoEcoleID)
	if errMois == nil {
		defer rowsMois.Close()
		for rowsMois.Next() {
			var rm models.RecetteMensuelle
			if err := rowsMois.Scan(&rm.Mois, &rm.Montant); err == nil {
				recettesMensuelles = append(recettesMensuelles, rm)
			}
		}
		if err := rowsMois.Err(); err != nil {
			return nil, err
		}
	}
	if recettesMensuelles == nil {
		recettesMensuelles = []models.RecetteMensuelle{}
	}

	// Inscriptions réelles par formule
	var inscriptionsFormules []models.FormuleStat
	queryFormules := `
		SELECT f.id, f.intitule, COUNT(d.id) AS nb
		FROM formules_permis f
		LEFT JOIN dossiers_inscription d ON d.formule_id = f.id AND d.auto_ecole_id = $1
		WHERE f.auto_ecole_id = $1
		GROUP BY f.id, f.intitule
		ORDER BY nb DESC, f.intitule ASC
	`
	rowsFormules, errFormules := r.db.QueryContext(ctx, queryFormules, autoEcoleID)
	if errFormules == nil {
		defer rowsFormules.Close()
		for rowsFormules.Next() {
			var fs models.FormuleStat
			if err := rowsFormules.Scan(&fs.FormuleID, &fs.Intitule, &fs.Inscrits); err == nil {
				inscriptionsFormules = append(inscriptionsFormules, fs)
			}
		}
		if err := rowsFormules.Err(); err != nil {
			return nil, err
		}
	}
	if inscriptionsFormules == nil {
		inscriptionsFormules = []models.FormuleStat{}
	}

	// Derniers paiements réels encaissés
	var derniersPaiements []models.PaiementRecent
	queryPaiements := `
		SELECT p.id, COALESCE(p.reference_recu, ''), d.numero_dossier,
		       c.nom || ' ' || c.prenoms, f.intitule, p.montant,
		       p.mode_reglement, to_char(p.date_paiement, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
		FROM paiements p
		JOIN dossiers_inscription d ON d.id = p.dossier_id
		JOIN candidats c ON c.id = d.candidat_id
		JOIN formules_permis f ON f.id = d.formule_id
		WHERE p.auto_ecole_id = $1
		ORDER BY p.date_paiement DESC
		LIMIT 10
	`
	rowsPaiements, errPaiements := r.db.QueryContext(ctx, queryPaiements, autoEcoleID)
	if errPaiements == nil {
		defer rowsPaiements.Close()
		for rowsPaiements.Next() {
			var pr models.PaiementRecent
			if err := rowsPaiements.Scan(
				&pr.ID, &pr.NumeroRecu, &pr.NumeroDossier,
				&pr.CandidatNom, &pr.FormuleNom, &pr.Montant,
				&pr.ModeReglement, &pr.DatePaiement,
			); err == nil {
				derniersPaiements = append(derniersPaiements, pr)
			}
		}
		if err := rowsPaiements.Err(); err != nil {
			return nil, err
		}
	}
	if derniersPaiements == nil {
		derniersPaiements = []models.PaiementRecent{}
	}

	// Taux de réussite réel aux examens
	var tauxReussite float64
	queryTaux := `
		SELECT
			COALESCE(
				ROUND(
					(COUNT(CASE WHEN resultat_final = 'ADMIS' THEN 1 END)::numeric /
					 NULLIF(COUNT(CASE WHEN resultat_final IN ('ADMIS', 'AJOURNE') THEN 1 END), 0)::numeric) * 100
				, 1),
			0.0)
		FROM candidats_examen ce
		JOIN sessions_examen se ON se.id = ce.session_examen_id
		WHERE se.auto_ecole_id = $1
	`
	_ = r.db.QueryRowContext(ctx, queryTaux, autoEcoleID).Scan(&tauxReussite)

	// Total candidats actifs
	var totalCandidats int
	_ = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM candidats WHERE auto_ecole_id = $1`, autoEcoleID).Scan(&totalCandidats)

	return &models.TableauGerantResponse{
		TotalCreancesImpayes:  totalCreances,
		NombreDossiersImpayes: len(impayes),
		Impayes:               impayes,
		ElevesAssidusCode:     assidus,
		Ventilation:           vent,
		RecettesMensuelles:    recettesMensuelles,
		InscriptionsFormules:  inscriptionsFormules,
		DerniersPaiements:     derniersPaiements,
		TauxReussiteExamen:    tauxReussite,
		TotalCandidatsActifs:  totalCandidats,
	}, nil
}


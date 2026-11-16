package repository

import (
	"context"
	"time"

	"autoschool/backend/internal/models"
)

type SyncUserSnapshot struct {
	ID             string  `json:"id"`
	AutoEcoleID    string  `json:"auto_ecole_id"`
	AgenceID       string  `json:"agence_id"`
	Nom            string  `json:"nom"`
	Prenoms        string  `json:"prenoms"`
	Telephone      string  `json:"telephone"`
	Role           string  `json:"role"`
	CodePinManager *string `json:"code_pin_manager,omitempty"`
	Statut         string  `json:"statut"`
}

type SyncSessionCaisseSnapshot struct {
	ID                    string   `json:"id"`
	AgenceID              string   `json:"agence_id"`
	CaissierID            string   `json:"caissier_id"`
	DateOuverture         string   `json:"date_ouverture"`
	DateFermeture         *string  `json:"date_fermeture,omitempty"`
	FondDeCaisse          float64  `json:"fond_de_caisse"`
	TotalEspecesTheorique float64  `json:"total_especes_theorique"`
	TotalEspecesCompte    *float64 `json:"total_especes_compte,omitempty"`
	EcartCaisse           *float64 `json:"ecart_caisse,omitempty"`
	Statut                string   `json:"statut"`
	Observations          *string  `json:"observations,omitempty"`
}

type SyncSeanceCodeSnapshot struct {
	ID              string  `json:"id"`
	AutoEcoleID     string  `json:"auto_ecole_id"`
	AgenceID        string  `json:"agence_id"`
	DateSeance      string  `json:"date_seance"`
	HeureDebut      string  `json:"heure_debut"`
	HeureFin        string  `json:"heure_fin"`
	ThemeTraite     string  `json:"theme_traite"`
	FormateurNom    string  `json:"formateur_nom"`
	MoniteurID      *string `json:"moniteur_id,omitempty"`
	TypeSeance      string  `json:"type_seance"`
	LienVisio       *string `json:"lien_visio,omitempty"`
	MotDePasseVisio *string `json:"mot_de_passe_visio,omitempty"`
	StatutSeance    string  `json:"statut_seance"`
}

type SyncCandidatExamenSnapshot struct {
	ID               string  `json:"id"`
	SessionExamenID  string  `json:"session_examen_id"`
	DossierID        string  `json:"dossier_id"`
	CandidatID       string  `json:"candidat_id"`
	NumeroTable      *string `json:"numero_table,omitempty"`
	ResultatCode     string  `json:"resultat_code"`
	ResultatConduite string  `json:"resultat_conduite"`
	ResultatFinal    string  `json:"resultat_final"`
}

type SyncPullSnapshot struct {
	AutoEcole       *models.AutoEcole            `json:"auto_ecole,omitempty"`
	Agences         []models.Agence              `json:"agences"`
	Utilisateurs    []SyncUserSnapshot           `json:"utilisateurs"`
	Formules        []models.FormulePermis       `json:"formules"`
	Candidats       []models.Candidat            `json:"candidats"`
	Dossiers        []models.DossierInscription  `json:"dossiers"`
	SessionsCaisse  []SyncSessionCaisseSnapshot  `json:"sessions_caisse"`
	Paiements       []models.Paiement            `json:"paiements"`
	Vehicules       []models.Vehicule            `json:"vehicules"`
	SeancesCode     []SyncSeanceCodeSnapshot     `json:"seances_code"`
	SessionsExamen  []models.SessionExamen       `json:"sessions_examen"`
	CandidatsExamen []SyncCandidatExamenSnapshot `json:"candidats_examen"`
	SyncedAt        string                       `json:"synced_at"`
}

// GetSyncPullSnapshot retourne un instantané complet de l'auto-école pour hydrater SQLite local.
func (r *GuichetRepository) GetSyncPullSnapshot(ctx context.Context, autoEcoleID string) (*SyncPullSnapshot, error) {
	snap := &SyncPullSnapshot{
		Agences:         []models.Agence{},
		Utilisateurs:    []SyncUserSnapshot{},
		Formules:        []models.FormulePermis{},
		Candidats:       []models.Candidat{},
		Dossiers:        []models.DossierInscription{},
		SessionsCaisse:  []SyncSessionCaisseSnapshot{},
		Paiements:       []models.Paiement{},
		Vehicules:       []models.Vehicule{},
		SeancesCode:     []SyncSeanceCodeSnapshot{},
		SessionsExamen:  []models.SessionExamen{},
		CandidatsExamen: []SyncCandidatExamenSnapshot{},
		SyncedAt:        time.Now().UTC().Format(time.RFC3339),
	}

	var ae models.AutoEcole
	err := r.db.QueryRowContext(ctx, `
		SELECT id, code_unique, nom_commercial, COALESCE(numero_agrement, ''), COALESCE(telephone, ''), email, COALESCE(devise, 'XOF')
		FROM auto_ecoles WHERE id = $1
	`, autoEcoleID).Scan(&ae.ID, &ae.CodeUnique, &ae.NomCommercial, &ae.NumeroAgrement, &ae.Telephone, &ae.Email, &ae.Devise)
	if err == nil {
		snap.AutoEcole = &ae
	}

	if agRows, err := r.db.QueryContext(ctx, `
		SELECT id, auto_ecole_id, nom, COALESCE(adresse, ''), COALESCE(telephone, '')
		FROM agences WHERE auto_ecole_id = $1
	`, autoEcoleID); err == nil {
		for agRows.Next() {
			var ag models.Agence
			if agRows.Scan(&ag.ID, &ag.AutoEcoleID, &ag.Nom, &ag.Adresse, &ag.Telephone) == nil {
				snap.Agences = append(snap.Agences, ag)
			}
		}
		agRows.Close()
	}

	if uRows, err := r.db.QueryContext(ctx, `
		SELECT id, auto_ecole_id, agence_id, nom, prenoms, telephone, role, code_pin_manager, statut
		FROM utilisateurs WHERE auto_ecole_id = $1
	`, autoEcoleID); err == nil {
		for uRows.Next() {
			var u SyncUserSnapshot
			if uRows.Scan(&u.ID, &u.AutoEcoleID, &u.AgenceID, &u.Nom, &u.Prenoms, &u.Telephone, &u.Role, &u.CodePinManager, &u.Statut) == nil {
				snap.Utilisateurs = append(snap.Utilisateurs, u)
			}
		}
		uRows.Close()
	}

	if formules, err := r.GetFormules(ctx, autoEcoleID); err == nil && formules != nil {
		snap.Formules = formules
	}

	if cRows, err := r.db.QueryContext(ctx, `
		SELECT id, auto_ecole_id, agence_id, nom, prenoms, genre,
		       to_char(date_naissance, 'YYYY-MM-DD'), COALESCE(lieu_naissance, ''),
		       COALESCE(nationalite, 'Ivoirienne'), telephone_whatsapp, photo_url,
		       type_piece_id, numero_piece_id, visite_medicale_statut, created_at
		FROM candidats WHERE auto_ecole_id = $1
		ORDER BY created_at DESC LIMIT 1000
	`, autoEcoleID); err == nil {
		for cRows.Next() {
			var c models.Candidat
			if cRows.Scan(
				&c.ID, &c.AutoEcoleID, &c.AgenceID, &c.Nom, &c.Prenoms, &c.Genre,
				&c.DateNaissance, &c.LieuNaissance, &c.Nationalite, &c.TelephoneWhatsApp,
				&c.PhotoURL, &c.TypePieceID, &c.NumeroPieceID, &c.VisiteMedicaleStatut, &c.CreatedAt,
			) == nil {
				snap.Candidats = append(snap.Candidats, c)
			}
		}
		cRows.Close()
	}

	if dRows, err := r.db.QueryContext(ctx, `
		SELECT id, auto_ecole_id, agence_id, candidat_id, formule_id, numero_dossier,
		       to_char(date_inscription, 'YYYY-MM-DD'), montant_total, remise_accordee, solde_restant,
		       statut_dossier, heures_code_effectuees, heures_conduite_effectuees,
		       COALESCE(derogation_conduite_accordee, false)
		FROM dossiers_inscription WHERE auto_ecole_id = $1
		ORDER BY created_at DESC LIMIT 1000
	`, autoEcoleID); err == nil {
		for dRows.Next() {
			var d models.DossierInscription
			if dRows.Scan(
				&d.ID, &d.AutoEcoleID, &d.AgenceID, &d.CandidatID, &d.FormuleID, &d.NumeroDossier,
				&d.DateInscription, &d.MontantTotal, &d.RemiseAccordee, &d.SoldeRestant,
				&d.StatutDossier, &d.HeuresCodeEffectuees, &d.HeuresConduiteEffectuees,
				&d.DerogationConduiteAccordee,
			) == nil {
				snap.Dossiers = append(snap.Dossiers, d)
			}
		}
		dRows.Close()
	}

	if sRows, err := r.db.QueryContext(ctx, `
		SELECT s.id, s.agence_id, s.caissier_id,
		       to_char(s.date_ouverture, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
		       CASE WHEN s.date_fermeture IS NOT NULL THEN to_char(s.date_fermeture, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') ELSE NULL END,
		       s.fond_de_caisse, s.total_especes_theorique, s.total_especes_compte, s.ecart_caisse, s.statut, s.observations
		FROM sessions_caisse s
		JOIN agences a ON a.id = s.agence_id
		WHERE a.auto_ecole_id = $1
		ORDER BY s.date_ouverture DESC LIMIT 100
	`, autoEcoleID); err == nil {
		for sRows.Next() {
			var s SyncSessionCaisseSnapshot
			if sRows.Scan(
				&s.ID, &s.AgenceID, &s.CaissierID, &s.DateOuverture, &s.DateFermeture,
				&s.FondDeCaisse, &s.TotalEspecesTheorique, &s.TotalEspecesCompte, &s.EcartCaisse, &s.Statut, &s.Observations,
			) == nil {
				snap.SessionsCaisse = append(snap.SessionsCaisse, s)
			}
		}
		sRows.Close()
	}

	if pRows, err := r.db.QueryContext(ctx, `
		SELECT p.id, p.auto_ecole_id, p.agence_id, p.dossier_id, p.session_caisse_id,
		       p.montant, p.mode_reglement, p.reference_transaction, p.reference_recu,
		       p.solde_apres_paiement, to_char(p.date_paiement, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
		       p.encaisse_par_id, COALESCE(p.est_annule, false)
		FROM paiements p
		WHERE p.auto_ecole_id = $1
		ORDER BY p.date_paiement DESC LIMIT 1000
	`, autoEcoleID); err == nil {
		for pRows.Next() {
			var p models.Paiement
			if pRows.Scan(
				&p.ID, &p.AutoEcoleID, &p.AgenceID, &p.DossierID, &p.SessionCaisseID,
				&p.Montant, &p.ModeReglement, &p.ReferenceTransaction, &p.ReferenceRecu,
				&p.SoldeApresPaiement, &p.DatePaiement, &p.EncaisseParID, &p.EstAnnule,
			) == nil {
				snap.Paiements = append(snap.Paiements, p)
			}
		}
		pRows.Close()
	}

	if vRows, err := r.db.QueryContext(ctx, `
		SELECT v.id, v.agence_id, v.immatriculation, v.marque_modele, v.type_boite,
		       v.kilometrage_actuel, to_char(v.date_expiration_assurance, 'YYYY-MM-DD'),
		       to_char(v.date_visite_technique, 'YYYY-MM-DD'), v.statut, v.created_at, v.updated_at
		FROM vehicules v
		JOIN agences a ON a.id = v.agence_id
		WHERE a.auto_ecole_id = $1
	`, autoEcoleID); err == nil {
		for vRows.Next() {
			var v models.Vehicule
			if vRows.Scan(
				&v.ID, &v.AgenceID, &v.Immatriculation, &v.MarqueModele, &v.TypeBoite,
				&v.KilometrageActuel, &v.DateExpirationAssurance, &v.DateVisiteTechnique,
				&v.Statut, &v.CreatedAt, &v.UpdatedAt,
			) == nil {
				snap.Vehicules = append(snap.Vehicules, v)
			}
		}
		vRows.Close()
	}

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

	if scRows, err := r.db.QueryContext(ctx, `
		SELECT id, auto_ecole_id, agence_id, to_char(date_seance, 'YYYY-MM-DD'), heure_debut, heure_fin,
		       theme_traite, formateur_nom, moniteur_id, COALESCE(type_seance, 'PRESENTIEL'),
		       lien_visio, mot_de_passe_visio, COALESCE(statut_seance, 'PLANIFIEE')
		FROM seances_code WHERE auto_ecole_id = $1
		ORDER BY date_seance DESC LIMIT 200
	`, autoEcoleID); err == nil {
		for scRows.Next() {
			var sc SyncSeanceCodeSnapshot
			if scRows.Scan(
				&sc.ID, &sc.AutoEcoleID, &sc.AgenceID, &sc.DateSeance, &sc.HeureDebut, &sc.HeureFin,
				&sc.ThemeTraite, &sc.FormateurNom, &sc.MoniteurID, &sc.TypeSeance,
				&sc.LienVisio, &sc.MotDePasseVisio, &sc.StatutSeance,
			) == nil {
				snap.SeancesCode = append(snap.SeancesCode, sc)
			}
		}
		scRows.Close()
	}

	if seRows, err := r.db.QueryContext(ctx, `
		SELECT id, auto_ecole_id, agence_id, code_session, to_char(date_examen, 'YYYY-MM-DD'),
		       centre_examen, categorie_permis, statut_session, inspecteur_nom
		FROM sessions_examen WHERE auto_ecole_id = $1
		ORDER BY date_examen DESC LIMIT 50
	`, autoEcoleID); err == nil {
		for seRows.Next() {
			var se models.SessionExamen
			if seRows.Scan(
				&se.ID, &se.AutoEcoleID, &se.AgenceID, &se.CodeSession, &se.DateExamen,
				&se.CentreExamen, &se.CategoriePermis, &se.StatutSession, &se.InspecteurNom,
			) == nil {
				snap.SessionsExamen = append(snap.SessionsExamen, se)
			}
		}
		seRows.Close()
	}

	if ceRows, err := r.db.QueryContext(ctx, `
		SELECT ce.id, ce.session_examen_id, ce.dossier_id, ce.candidat_id, ce.numero_table,
		       ce.resultat_code, ce.resultat_conduite, ce.resultat_final
		FROM candidats_examen ce
		JOIN sessions_examen se ON se.id = ce.session_examen_id
		WHERE se.auto_ecole_id = $1
		LIMIT 500
	`, autoEcoleID); err == nil {
		for ceRows.Next() {
			var ce SyncCandidatExamenSnapshot
			if ceRows.Scan(
				&ce.ID, &ce.SessionExamenID, &ce.DossierID, &ce.CandidatID, &ce.NumeroTable,
				&ce.ResultatCode, &ce.ResultatConduite, &ce.ResultatFinal,
			) == nil {
				snap.CandidatsExamen = append(snap.CandidatsExamen, ce)
			}
		}
		ceRows.Close()
	}

	return snap, nil
}

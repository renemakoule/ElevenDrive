package repository

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"
	"time"

	"autoschool/backend/internal/database"
	"autoschool/backend/internal/models"
)

type SuperAdminRepository struct {
	db *database.DB
}

func NewSuperAdminRepository(db *database.DB) *SuperAdminRepository {
	return &SuperAdminRepository{db: db}
}

// GetGlobalStats calcule les métriques d'ensemble pour le Super-Admin
func (r *SuperAdminRepository) GetGlobalStats(ctx context.Context) (*models.SuperAdminStats, error) {
	start := time.Now()
	stats := &models.SuperAdminStats{
		CloudStatus: "ONLINE",
	}

	err := database.Retry(ctx, func() error {
		return database.Safe(func() error {
			// 1. Écoles
			_ = r.db.QueryRowContext(ctx, `SELECT COUNT(*), COUNT(CASE WHEN est_actif THEN 1 END) FROM auto_ecoles`).Scan(&stats.TotalSchools, &stats.ActiveSchools)

			// 2. Utilisateurs
			_ = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM utilisateurs`).Scan(&stats.TotalUsers)

			// 3. Candidats
			_ = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM candidats`).Scan(&stats.TotalCandidats)

			// 4. Volume financier total
			_ = r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(montant), 0) FROM paiements`).Scan(&stats.TotalVolumePaiement)

			return nil
		})
	})
	if err != nil {
		stats.CloudStatus = "DEGRADED"
	}

	stats.ClusterLatencyMs = time.Since(start).Milliseconds()
	return stats, nil
}

// GetAllSchools liste l'ensemble des établissements enregistrés avec statistiques agrégées, quotas et licences
func (r *SuperAdminRepository) GetAllSchools(ctx context.Context, search string) ([]models.SchoolItemSummary, error) {
	query := `
		SELECT 
			e.id, e.code_unique, e.nom_commercial, e.numero_agrement, e.telephone, e.email,
			e.devise, e.pays, COALESCE(e.ville, ''), e.est_actif,
			COALESCE(e.max_machines, 3),
			COALESCE((SELECT COUNT(*) FROM machines_auto_ecoles m WHERE m.auto_ecole_id = e.id AND m.est_active), 0),
			e.date_expiration_licence,
			COALESCE(e.type_licence, 'ANNUEL'),
			COALESCE(e.modules_actifs, '["guichet","caisse","conduite","examens","flotte","archivage_cloud"]'),
			e.created_at, e.updated_at,
			(SELECT COUNT(*) FROM agences a WHERE a.auto_ecole_id = e.id) as agences_count,
			(SELECT COUNT(*) FROM utilisateurs u WHERE u.auto_ecole_id = e.id) as users_count,
			(SELECT COUNT(*) FROM candidats c WHERE c.auto_ecole_id = e.id) as candidats_count,
			COALESCE((SELECT SUM(p.montant) FROM paiements p WHERE p.auto_ecole_id = e.id), 0) as total_volume
		FROM auto_ecoles e
	`
	var args []interface{}
	if search != "" {
		query += ` WHERE LOWER(e.nom_commercial) LIKE LOWER($1) OR LOWER(e.code_unique) LIKE LOWER($1) OR LOWER(e.numero_agrement) LIKE LOWER($1)`
		args = append(args, "%"+search+"%")
	}
	query += ` ORDER BY e.created_at DESC`

	var rows []models.SchoolItemSummary
	err := database.Retry(ctx, func() error {
		return database.Safe(func() error {
			dbRows, err := r.db.QueryContext(ctx, query, args...)
			if err != nil {
				return fmt.Errorf("erreur récupération liste écoles super-admin: %w", err)
			}
			defer dbRows.Close()

			var tempRows []models.SchoolItemSummary
			for dbRows.Next() {
				var s models.SchoolItemSummary
				if err := dbRows.Scan(
					&s.ID, &s.CodeUnique, &s.NomCommercial, &s.NumeroAgrement, &s.Telephone, &s.Email,
					&s.Devise, &s.Pays, &s.Ville, &s.EstActif,
					&s.MaxMachines, &s.MachinesActives,
					&s.DateExpirationLicence, &s.TypeLicence, &s.ModulesActifs,
					&s.CreatedAt, &s.UpdatedAt,
					&s.AgencesCount, &s.UsersCount, &s.CandidatsCount, &s.TotalVolume,
				); err != nil {
					return err
				}

				// Calcul dynamique précis du statut de licence et des jours restants
				if s.DateExpirationLicence != nil {
					diff := time.Until(*s.DateExpirationLicence)
					days := int(math.Ceil(diff.Hours() / 24.0))
					s.JoursRestants = days
					if days <= 0 {
						s.JoursRestants = 0
						s.StatutLicence = "EXPIREE"
					} else if !s.EstActif {
						s.StatutLicence = "SUSPENDUE"
					} else if days <= 5 {
						s.StatutLicence = "EXPIRATION_PROCHE"
					} else {
						s.StatutLicence = "ACTIVE"
					}
				} else {
					s.JoursRestants = 9999
					s.StatutLicence = "ACTIVE"
				}

				tempRows = append(tempRows, s)
			}
			if err := dbRows.Err(); err != nil {
				return err
			}

			// Charger les agences rattachées avec leur localisation exacte (adresses, villes, communes)
			agencesRows, err := r.db.QueryContext(ctx, `
				SELECT id, auto_ecole_id, code_agence, nom, COALESCE(adresse, ''), COALESCE(telephone, ''),
				       COALESCE(ville, ''), COALESCE(commune, ''), est_siege, created_at, updated_at
				FROM agences
				ORDER BY est_siege DESC, nom ASC
			`)
			if err == nil {
				defer agencesRows.Close()
				agencesMap := make(map[string][]*models.Agence)
				for agencesRows.Next() {
					var a models.Agence
					if errScan := agencesRows.Scan(
						&a.ID, &a.AutoEcoleID, &a.CodeAgence, &a.Nom, &a.Adresse, &a.Telephone,
						&a.Ville, &a.Commune, &a.EstSiege, &a.CreatedAt, &a.UpdatedAt,
					); errScan == nil {
						agencesMap[a.AutoEcoleID] = append(agencesMap[a.AutoEcoleID], &a)
					}
				}
				for i := range tempRows {
					if agList, ok := agencesMap[tempRows[i].ID]; ok {
						tempRows[i].Agences = agList
					} else {
						tempRows[i].Agences = []*models.Agence{}
					}
				}
			}

			rows = tempRows
			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// UpdateSchoolAdmin modifie complètement une auto-école
func (r *SuperAdminRepository) UpdateSchoolAdmin(ctx context.Context, input *models.UpdateSchoolAdminInput) error {
	query := `
		UPDATE auto_ecoles SET
			nom_commercial = $1,
			code_unique = $2,
			numero_agrement = $3,
			telephone = $4,
			email = $5,
			ville = $6,
			pays = $7,
			devise = $8,
			est_actif = $9,
			max_machines = $10,
			date_expiration_licence = $11,
			type_licence = $12,
			modules_actifs = $13,
			updated_at = NOW()
		WHERE id = $14
	`
	_, err := r.db.ExecContext(ctx, query,
		input.NomCommercial, input.CodeUnique, input.NumeroAgrement,
		input.Telephone, input.Email, input.Ville, input.Pays, input.Devise,
		input.EstActif, input.MaxMachines, input.DateExpirationLicence,
		input.TypeLicence, input.ModulesActifs, input.SchoolID,
	)
	return err
}

// UpdateSchoolLicence définit spécifiquement le quota de machines et la durée de validité
func (r *SuperAdminRepository) UpdateSchoolLicence(ctx context.Context, input *models.UpdateSchoolLicenceInput) error {
	var expDate time.Time
	mode := strings.ToUpper(strings.TrimSpace(input.Mode))

	if input.DateExpirationLicence != nil {
		expDate = *input.DateExpirationLicence
	} else if mode == "KEEP_DATE" {
		var currentExp sql.NullTime
		_ = r.db.QueryRowContext(ctx, `SELECT date_expiration_licence FROM auto_ecoles WHERE id = $1`, input.SchoolID).Scan(&currentExp)
		if currentExp.Valid {
			expDate = currentExp.Time
		} else {
			expDate = time.Now().AddDate(0, 0, 30)
		}
	} else {
		baseDate := time.Now()
		if mode == "EXTEND" {
			// Prolonger à partir de la date d'expiration actuelle si encore active
			var currentExp sql.NullTime
			_ = r.db.QueryRowContext(ctx, `SELECT date_expiration_licence FROM auto_ecoles WHERE id = $1`, input.SchoolID).Scan(&currentExp)
			if currentExp.Valid && currentExp.Time.After(baseDate) {
				baseDate = currentExp.Time
			}
		}

		switch strings.ToUpper(input.DureeUnite) {
		case "ANNEES", "ANNEE":
			expDate = baseDate.AddDate(input.DureeValeur, 0, 0)
		case "MOIS":
			expDate = baseDate.AddDate(0, input.DureeValeur, 0)
		default: // "JOURS"
			expDate = baseDate.AddDate(0, 0, input.DureeValeur)
		}
	}

	query := `
		UPDATE auto_ecoles SET
			max_machines = $1,
			type_licence = $2,
			date_expiration_licence = $3,
			updated_at = NOW()
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query, input.MaxMachines, input.TypeLicence, expDate, input.SchoolID)
	return err
}

// UpdateSchoolModules définit la liste des modules autorisés en JSON
func (r *SuperAdminRepository) UpdateSchoolModules(ctx context.Context, schoolID string, modulesJSON string) error {
	query := `UPDATE auto_ecoles SET modules_actifs = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, modulesJSON, schoolID)
	return err
}

// DeleteSchool effectue la suppression sécurisée d'une auto-école
func (r *SuperAdminRepository) DeleteSchool(ctx context.Context, schoolID string) error {
	// Suppression atomique
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Purge des sessions, caisses, candidats et collaborateurs rattachés
	_, _ = tx.ExecContext(ctx, `DELETE FROM machines_auto_ecoles WHERE auto_ecole_id = $1`, schoolID)
	_, _ = tx.ExecContext(ctx, `DELETE FROM utilisateurs WHERE auto_ecole_id = $1`, schoolID)
	_, _ = tx.ExecContext(ctx, `DELETE FROM agences WHERE auto_ecole_id = $1`, schoolID)

	// 2. Suppression de l'auto-école elle-même
	res, err := tx.ExecContext(ctx, `DELETE FROM auto_ecoles WHERE id = $1`, schoolID)
	if err != nil {
		return fmt.Errorf("erreur suppression auto-école: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("auto-école introuvable avec l'ID %s", schoolID)
	}

	return tx.Commit()
}

// LogAudit enregistre une action d'administration dans le journal système
func (r *SuperAdminRepository) LogAudit(ctx context.Context, action, cibleType, cibleID, details, ip string) error {
	query := `INSERT INTO audit_logs_super_admin (action, cible_type, cible_id, details, ip_address) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, action, cibleType, cibleID, details, ip)
	return err
}

// GetAuditLogs récupère les dernières entrées du journal d'audit
func (r *SuperAdminRepository) GetAuditLogs(ctx context.Context, limit int) ([]models.AuditLogEntry, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}
	query := `SELECT id, action, cible_type, cible_id, COALESCE(details, ''), COALESCE(ip_address, ''), created_at FROM audit_logs_super_admin ORDER BY created_at DESC LIMIT $1`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.AuditLogEntry
	for rows.Next() {
		var e models.AuditLogEntry
		if err := rows.Scan(&e.ID, &e.Action, &e.CibleType, &e.CibleID, &e.Details, &e.IPAddress, &e.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

// GetAllUsers liste tous les utilisateurs enregistrés avec leur auto-école
func (r *SuperAdminRepository) GetAllUsers(ctx context.Context, search string) ([]models.GlobalUserSummary, error) {
	query := `
		SELECT 
			u.id, u.auto_ecole_id, COALESCE(e.nom_commercial, 'Système Root') as auto_ecole_nom,
			u.nom, u.prenoms, u.telephone, u.email, u.role, u.statut, u.created_at
		FROM utilisateurs u
		LEFT JOIN auto_ecoles e ON e.id = u.auto_ecole_id
	`
	var args []interface{}
	if search != "" {
		query += ` WHERE LOWER(u.nom) LIKE LOWER($1) OR LOWER(u.prenoms) LIKE LOWER($1) OR LOWER(u.telephone) LIKE LOWER($1) OR LOWER(COALESCE(u.email, '')) LIKE LOWER($1) OR LOWER(COALESCE(e.nom_commercial, '')) LIKE LOWER($1)`
		args = append(args, "%"+search+"%")
	}
	query += ` ORDER BY u.created_at DESC LIMIT 200`

	var rows []models.GlobalUserSummary
	err := database.Retry(ctx, func() error {
		return database.Safe(func() error {
			dbRows, err := r.db.QueryContext(ctx, query, args...)
			if err != nil {
				return fmt.Errorf("erreur récupération liste utilisateurs super-admin: %w", err)
			}
			defer dbRows.Close()

			var temp []models.GlobalUserSummary
			for dbRows.Next() {
				var u models.GlobalUserSummary
				if err := dbRows.Scan(
					&u.ID, &u.AutoEcoleID, &u.AutoEcoleNom,
					&u.Nom, &u.Prenoms, &u.Telephone, &u.Email,
					&u.Role, &u.Statut, &u.CreatedAt,
				); err != nil {
					return err
				}
				temp = append(temp, u)
			}
			if err := dbRows.Err(); err != nil {
				return err
			}
			rows = temp
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// GetAllCandidats liste tous les candidats inscrits sur la plateforme
func (r *SuperAdminRepository) GetAllCandidats(ctx context.Context, search string) ([]models.GlobalCandidatSummary, error) {
	query := `
		SELECT 
			c.id, c.auto_ecole_id, COALESCE(e.nom_commercial, 'Auto-école') as auto_ecole_nom,
			c.nom, c.prenoms, c.telephone_whatsapp, c.genre,
			COALESCE(d.numero_dossier, ''), COALESCE(d.statut_dossier, 'EN_COURS'),
			COALESCE(d.montant_total, 0), COALESCE(d.solde_restant, 0),
			c.created_at
		FROM candidats c
		LEFT JOIN auto_ecoles e ON e.id = c.auto_ecole_id
		LEFT JOIN LATERAL (
			SELECT numero_dossier, statut_dossier, montant_total, solde_restant
			FROM dossiers d
			WHERE d.candidat_id = c.id
			ORDER BY d.created_at DESC
			LIMIT 1
		) d ON true
	`
	var args []interface{}
	if search != "" {
		query += ` WHERE LOWER(c.nom) LIKE LOWER($1) OR LOWER(c.prenoms) LIKE LOWER($1) OR LOWER(c.telephone_whatsapp) LIKE LOWER($1) OR LOWER(COALESCE(d.numero_dossier, '')) LIKE LOWER($1) OR LOWER(COALESCE(e.nom_commercial, '')) LIKE LOWER($1)`
		args = append(args, "%"+search+"%")
	}
	query += ` ORDER BY c.created_at DESC LIMIT 200`

	var rows []models.GlobalCandidatSummary
	err := database.Retry(ctx, func() error {
		return database.Safe(func() error {
			dbRows, err := r.db.QueryContext(ctx, query, args...)
			if err != nil {
				return fmt.Errorf("erreur récupération liste candidats super-admin: %w", err)
			}
			defer dbRows.Close()

			var temp []models.GlobalCandidatSummary
			for dbRows.Next() {
				var c models.GlobalCandidatSummary
				if err := dbRows.Scan(
					&c.ID, &c.AutoEcoleID, &c.AutoEcoleNom,
					&c.Nom, &c.Prenoms, &c.TelephoneWhatsApp, &c.Genre,
					&c.NumeroDossier, &c.StatutDossier,
					&c.MontantTotal, &c.SoldeRestant,
					&c.CreatedAt,
				); err != nil {
					return err
				}
				temp = append(temp, c)
			}
			if err := dbRows.Err(); err != nil {
				return err
			}
			rows = temp
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// ToggleSchoolStatus active ou suspend un établissement
func (r *SuperAdminRepository) ToggleSchoolStatus(ctx context.Context, schoolID string, estActif bool) error {
	_, err := r.db.ExecContext(ctx, `UPDATE auto_ecoles SET est_actif = $1, updated_at = NOW() WHERE id = $2`, estActif, schoolID)
	return err
}

// GetSuperAdminUser recherche un utilisateur avec le rôle SUPER_ADMIN
func (r *SuperAdminRepository) GetSuperAdminByIdentifier(ctx context.Context, identifiant string) (*models.User, error) {
	query := `
		SELECT id, auto_ecole_id, agence_id, nom, prenoms, telephone, email,
		       mot_de_passe_hash, code_pin_manager, role, statut, created_at, updated_at
		FROM utilisateurs
		WHERE (telephone = $1 OR LOWER(COALESCE(email, '')) = LOWER($1))
		  AND UPPER(role) = 'SUPER_ADMIN'
		LIMIT 1
	`
	var u models.User
	err := r.db.QueryRowContext(ctx, query, identifiant).Scan(
		&u.ID, &u.AutoEcoleID, &u.AgenceID, &u.Nom, &u.Prenoms, &u.Telephone, &u.Email,
		&u.MotDePasseHash, &u.CodePinManager, &u.Role, &u.Statut, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

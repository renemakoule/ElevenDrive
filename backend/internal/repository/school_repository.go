package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"autoschool/backend/internal/database"
	"autoschool/backend/internal/models"
)

// SchoolRepository gère les opérations en base sur les auto-écoles et leurs agences.
type SchoolRepository struct {
	db *database.DB
}

func NewSchoolRepository(db *database.DB) *SchoolRepository {
	return &SchoolRepository{db: db}
}

// CreateSchoolWithAgencyAndAdmin crée atomiquement dans une transaction :
// 1. L'Auto-École
// 2. L'Agence Siège
// 3. Le compte utilisateur Directeur
func (r *SchoolRepository) CreateSchoolWithAgencyAndAdmin(
	ctx context.Context,
	school *models.AutoEcole,
	agence *models.Agence,
	admin *models.User,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("impossible de démarrer la transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Insérer l'auto-école avec sa licence d'essai Super-Admin (30 jours)
	maxMachines := school.MaxMachines
	if maxMachines <= 0 {
		maxMachines = 3
	}
	typeLicence := school.TypeLicence
	if typeLicence == "" {
		typeLicence = "TRIAL"
	}
	modules := school.ModulesActifs
	if modules == "" {
		modules = `["guichet","caisse","conduite","examens","flotte","archivage_cloud"]`
	}

	querySchool := `
		INSERT INTO auto_ecoles (
			id, code_unique, nom_commercial, numero_agrement, telephone, email, devise, pays, est_actif,
			max_machines, type_licence, date_expiration_licence, modules_actifs
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9,
			$10, $11, $12, $13
		) RETURNING created_at, updated_at
	`
	err = tx.QueryRowContext(ctx, querySchool,
		school.ID, school.CodeUnique, school.NomCommercial, school.NumeroAgrement,
		school.Telephone, school.Email, school.Devise, school.Pays, school.EstActif,
		maxMachines, typeLicence, school.DateExpirationLicence, modules,
	).Scan(&school.CreatedAt, &school.UpdatedAt)
	if err != nil {
		return fmt.Errorf("erreur insertion auto-école: %w", err)
	}

	// 2. Insérer l'agence siège
	queryAgence := `
		INSERT INTO agences (
			id, auto_ecole_id, code_agence, nom, adresse, telephone, ville, commune, est_siege
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		) RETURNING created_at, updated_at
	`
	err = tx.QueryRowContext(ctx, queryAgence,
		agence.ID, school.ID, agence.CodeAgence, agence.Nom, agence.Adresse,
		agence.Telephone, agence.Ville, agence.Commune, agence.EstSiege,
	).Scan(&agence.CreatedAt, &agence.UpdatedAt)
	if err != nil {
		return fmt.Errorf("erreur insertion agence siège: %w", err)
	}

	// 3. Insérer le compte Directeur
	queryAdmin := `
		INSERT INTO utilisateurs (
			id, auto_ecole_id, agence_id, nom, prenoms, telephone, email,
			mot_de_passe_hash, code_pin_manager, role, statut
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		) RETURNING created_at, updated_at
	`
	err = tx.QueryRowContext(ctx, queryAdmin,
		admin.ID, school.ID, agence.ID, admin.Nom, admin.Prenoms, admin.Telephone, admin.Email,
		admin.MotDePasseHash, admin.CodePinManager, admin.Role, admin.Statut,
	).Scan(&admin.CreatedAt, &admin.UpdatedAt)
	if err != nil {
		return fmt.Errorf("erreur insertion compte directeur: %w", err)
	}

	// 4. Catalogue vierge par défaut : l'auto-école définit elle-même ses tarifs
	return tx.Commit()
}

// GetByID recherche une auto-école par son identifiant unique UUID.
func (r *SchoolRepository) GetByID(ctx context.Context, id string) (*models.AutoEcole, error) {
	query := `
		SELECT e.id, e.code_unique, e.nom_commercial, e.numero_agrement, e.telephone, e.email,
		       e.devise, e.pays, COALESCE(e.ville, ''), e.logo_url, e.est_actif,
		       COALESCE(e.max_machines, 3),
		       COALESCE((SELECT COUNT(*) FROM machines_auto_ecoles m WHERE m.auto_ecole_id = e.id AND m.est_active), 0),
		       COALESCE(e.type_licence, 'TRIAL'), e.date_expiration_licence, COALESCE(e.modules_actifs, ''),
		       e.created_at, e.updated_at
		FROM auto_ecoles e
		WHERE e.id = $1
	`
	var school models.AutoEcole
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&school.ID, &school.CodeUnique, &school.NomCommercial, &school.NumeroAgrement,
		&school.Telephone, &school.Email, &school.Devise, &school.Pays, &school.Ville,
		&school.LogoURL, &school.EstActif,
		&school.MaxMachines, &school.MachinesActives, &school.TypeLicence, &school.DateExpirationLicence, &school.ModulesActifs,
		&school.CreatedAt, &school.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erreur recherche auto-école par id: %w", err)
	}
	return &school, nil
}

// GetByCodeUnique recherche une auto-école par son slug ou code unique.
func (r *SchoolRepository) GetByCodeUnique(ctx context.Context, slug string) (*models.AutoEcole, error) {
	query := `
		SELECT e.id, e.code_unique, e.nom_commercial, e.numero_agrement, e.telephone, e.email,
		       e.devise, e.pays, COALESCE(e.ville, ''), e.logo_url, e.est_actif,
		       COALESCE(e.max_machines, 3),
		       COALESCE((SELECT COUNT(*) FROM machines_auto_ecoles m WHERE m.auto_ecole_id = e.id AND m.est_active), 0),
		       COALESCE(e.type_licence, 'TRIAL'), e.date_expiration_licence, COALESCE(e.modules_actifs, ''),
		       e.created_at, e.updated_at
		FROM auto_ecoles e
		WHERE LOWER(e.code_unique) = LOWER($1)
	`
	var school models.AutoEcole
	err := database.Retry(ctx, func() error {
		return r.db.QueryRowContext(ctx, query, slug).Scan(
			&school.ID, &school.CodeUnique, &school.NomCommercial, &school.NumeroAgrement,
			&school.Telephone, &school.Email, &school.Devise, &school.Pays, &school.Ville,
			&school.LogoURL, &school.EstActif,
			&school.MaxMachines, &school.MachinesActives, &school.TypeLicence, &school.DateExpirationLicence, &school.ModulesActifs,
			&school.CreatedAt, &school.UpdatedAt,
		)
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erreur recherche auto-école par code: %w", err)
	}
	return &school, nil
}

// UpdateSchool met à jour les informations officielles d'une auto-école.
func (r *SchoolRepository) UpdateSchool(ctx context.Context, s *models.AutoEcole) error {
	query := `
		UPDATE auto_ecoles
		SET nom_commercial = $1,
		    code_unique = $2,
		    numero_agrement = $3,
		    telephone = $4,
		    email = $5,
		    devise = $6,
		    pays = $7,
		    updated_at = NOW()
		WHERE id = $8
		RETURNING updated_at
	`
	err := r.db.QueryRowContext(ctx, query,
		s.NomCommercial, s.CodeUnique, s.NumeroAgrement,
		s.Telephone, s.Email, s.Devise, s.Pays, s.ID,
	).Scan(&s.UpdatedAt)
	if err != nil {
		return fmt.Errorf("erreur mise à jour auto-école: %w", err)
	}
	return nil
}

// GetAgenceByID récupère les données d'une agence spécifique.
func (r *SchoolRepository) GetAgenceByID(ctx context.Context, agenceID string) (*models.Agence, error) {
	query := `
		SELECT id, auto_ecole_id, code_agence, nom, adresse, telephone,
		       ville, commune, est_siege, created_at, updated_at
		FROM agences
		WHERE id = $1
	`
	var agence models.Agence
	err := r.db.QueryRowContext(ctx, query, agenceID).Scan(
		&agence.ID, &agence.AutoEcoleID, &agence.CodeAgence, &agence.Nom, &agence.Adresse,
		&agence.Telephone, &agence.Ville, &agence.Commune, &agence.EstSiege,
		&agence.CreatedAt, &agence.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erreur recherche agence par id: %w", err)
	}
	return &agence, nil
}

// GetAgencesBySchoolID liste toutes les agences d'une auto-école.
func (r *SchoolRepository) GetAgencesBySchoolID(ctx context.Context, schoolID string) ([]*models.Agence, error) {
	query := `
		SELECT id, auto_ecole_id, code_agence, nom, adresse, telephone,
		       ville, commune, est_siege, created_at, updated_at
		FROM agences
		WHERE auto_ecole_id = $1
		ORDER BY est_siege DESC, nom ASC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return nil, fmt.Errorf("erreur listage agences: %w", err)
	}
	defer rows.Close()

	var agences []*models.Agence
	for rows.Next() {
		var a models.Agence
		if err := rows.Scan(
			&a.ID, &a.AutoEcoleID, &a.CodeAgence, &a.Nom, &a.Adresse,
			&a.Telephone, &a.Ville, &a.Commune, &a.EstSiege,
			&a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("erreur lecture ligne agence: %w", err)
		}
		agences = append(agences, &a)
	}

	return agences, rows.Err()
}

// GetSiegeAgence récupère l'agence siège principale d'une auto-école.
func (r *SchoolRepository) GetSiegeAgence(ctx context.Context, schoolID string) (*models.Agence, error) {
	query := `
		SELECT id, auto_ecole_id, code_agence, nom, adresse, telephone,
		       ville, commune, est_siege, created_at, updated_at
		FROM agences
		WHERE auto_ecole_id = $1
		ORDER BY est_siege DESC, created_at ASC
		LIMIT 1
	`
	var a models.Agence
	err := r.db.QueryRowContext(ctx, query, schoolID).Scan(
		&a.ID, &a.AutoEcoleID, &a.CodeAgence, &a.Nom, &a.Adresse,
		&a.Telephone, &a.Ville, &a.Commune, &a.EstSiege,
		&a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erreur recherche agence siege: %w", err)
	}
	return &a, nil
}

// CreateAgence insère une nouvelle succursale dans la base de données.
func (r *SchoolRepository) CreateAgence(ctx context.Context, a *models.Agence) error {
	if a.EstSiege {
		_, _ = r.db.ExecContext(ctx, `UPDATE agences SET est_siege = false WHERE auto_ecole_id = $1`, a.AutoEcoleID)
	}

	query := `
		INSERT INTO agences (
			id, auto_ecole_id, code_agence, nom, adresse, telephone,
			ville, commune, est_siege
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		) RETURNING created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query,
		a.ID, a.AutoEcoleID, a.CodeAgence, a.Nom, a.Adresse,
		a.Telephone, a.Ville, a.Commune, a.EstSiege,
	).Scan(&a.CreatedAt, &a.UpdatedAt)

	if err != nil {
		return fmt.Errorf("erreur creation agence: %w", err)
	}
	return nil
}

// UpdateAgence met à jour les informations d'une succursale.
func (r *SchoolRepository) UpdateAgence(ctx context.Context, a *models.Agence) error {
	if a.EstSiege {
		_, _ = r.db.ExecContext(ctx, `UPDATE agences SET est_siege = false WHERE auto_ecole_id = $1 AND id <> $2`, a.AutoEcoleID, a.ID)
	}

	query := `
		UPDATE agences
		SET code_agence = $1, nom = $2, adresse = $3, telephone = $4,
		    ville = $5, commune = $6, est_siege = $7, updated_at = NOW()
		WHERE id = $8 AND auto_ecole_id = $9
		RETURNING updated_at
	`
	err := r.db.QueryRowContext(ctx, query,
		a.CodeAgence, a.Nom, a.Adresse, a.Telephone,
		a.Ville, a.Commune, a.EstSiege, a.ID, a.AutoEcoleID,
	).Scan(&a.UpdatedAt)

	if err != nil {
		return fmt.Errorf("erreur modification agence: %w", err)
	}
	return nil
}

// DeleteAgence supprime une succursale (interdit si agence siège).
func (r *SchoolRepository) DeleteAgence(ctx context.Context, schoolID, agenceID string) error {
	var estSiege bool
	err := r.db.QueryRowContext(ctx, `SELECT est_siege FROM agences WHERE id = $1 AND auto_ecole_id = $2`, agenceID, schoolID).Scan(&estSiege)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("agence introuvable")
		}
		return fmt.Errorf("erreur verification agence: %w", err)
	}
	if estSiege {
		return errors.New("impossible de supprimer l'agence siege principale")
	}

	res, err := r.db.ExecContext(ctx, `DELETE FROM agences WHERE id = $1 AND auto_ecole_id = $2`, agenceID, schoolID)
	if err != nil {
		return fmt.Errorf("erreur suppression agence: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("aucune agence supprimee")
	}
	return nil
}

// CheckAndRegisterMachine vérifie et enregistre l'empreinte matérielle d'un poste de travail
func (r *SchoolRepository) CheckAndRegisterMachine(ctx context.Context, schoolID, fingerprint, machineName, osName string, maxMachines int) error {
	if fingerprint == "" {
		return nil // Non applicable (ex: accès mobile PWA)
	}
	if maxMachines <= 0 {
		maxMachines = 3
	}

	// 1. Vérifier si la machine est déjà enregistrée pour cette école
	var isActive bool
	err := r.db.QueryRowContext(ctx, `SELECT est_active FROM machines_auto_ecoles WHERE auto_ecole_id = $1 AND machine_fingerprint = $2`, schoolID, fingerprint).Scan(&isActive)
	if err == nil {
		if !isActive {
			return errors.New("ce poste de travail a été révoqué par le Super-Administrateur")
		}
		// Mettre à jour la dernière activité
		_, _ = r.db.ExecContext(ctx, `UPDATE machines_auto_ecoles SET derniere_activite = NOW() WHERE auto_ecole_id = $1 AND machine_fingerprint = $2`, schoolID, fingerprint)
		return nil
	}

	// 2. Si nouvelle machine : vérifier le quota de postes autorisés
	var count int
	_ = r.db.QueryRowContext(ctx, `SELECT count(*) FROM machines_auto_ecoles WHERE auto_ecole_id = $1 AND est_active = true`, schoolID).Scan(&count)
	if count >= maxMachines {
		return fmt.Errorf("quota de postes atteint (%d/%d postes autorisés). Veuillez contacter le Super-Administrateur pour ajouter un poste", count, maxMachines)
	}

	// 3. Enregistrer le nouveau poste
	if machineName == "" {
		machineName = fmt.Sprintf("Poste Desktop %d", count+1)
	}
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO machines_auto_ecoles (auto_ecole_id, machine_fingerprint, nom_machine, os, est_active, derniere_activite)
		VALUES ($1, $2, $3, $4, true, NOW())
		ON CONFLICT (auto_ecole_id, machine_fingerprint) DO UPDATE SET derniere_activite = NOW(), est_active = true
	`, schoolID, fingerprint, machineName, osName)
	return err
}

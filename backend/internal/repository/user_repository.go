package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"autoschool/backend/internal/database"
	"autoschool/backend/internal/models"
)

// UserRepository gère les opérations en base sur les utilisateurs.
type UserRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create insère un nouvel utilisateur dans la base de données.
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO utilisateurs (
			id, auto_ecole_id, agence_id, nom, prenoms, telephone, email,
			mot_de_passe_hash, code_pin_manager, role, statut
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		) RETURNING created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query,
		user.ID, user.AutoEcoleID, user.AgenceID, user.Nom, user.Prenoms,
		user.Telephone, user.Email, user.MotDePasseHash, user.CodePinManager,
		user.Role, user.Statut,
	).Scan(&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("erreur création utilisateur: %w", err)
	}
	return nil
}

// GetByID recherche un utilisateur par son UUID.
func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT id, auto_ecole_id, agence_id, nom, prenoms, telephone, email,
		       mot_de_passe_hash, code_pin_manager, role, statut, created_at, updated_at
		FROM utilisateurs
		WHERE id = $1
	`
	var u models.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.AutoEcoleID, &u.AgenceID, &u.Nom, &u.Prenoms, &u.Telephone, &u.Email,
		&u.MotDePasseHash, &u.CodePinManager, &u.Role, &u.Statut, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erreur recherche utilisateur par id: %w", err)
	}
	return &u, nil
}

// GetByIdentifierInSchool recherche un utilisateur par téléphone OU email au sein d'une auto-école donnée.
func (r *UserRepository) GetByIdentifierInSchool(ctx context.Context, schoolID string, identifier string) (*models.User, error) {
	query := `
		SELECT id, auto_ecole_id, agence_id, nom, prenoms, telephone, email,
		       mot_de_passe_hash, code_pin_manager, role, statut, created_at, updated_at
		FROM utilisateurs
		WHERE auto_ecole_id = $1 
		  AND (telephone = $2 OR LOWER(COALESCE(email, '')) = LOWER($2))
		LIMIT 1
	`
	var u models.User
	err := r.db.QueryRowContext(ctx, query, schoolID, identifier).Scan(
		&u.ID, &u.AutoEcoleID, &u.AgenceID, &u.Nom, &u.Prenoms, &u.Telephone, &u.Email,
		&u.MotDePasseHash, &u.CodePinManager, &u.Role, &u.Statut, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erreur recherche utilisateur par identifiant: %w", err)
	}
	return &u, nil
}

// GetByIdentifierGlobal recherche les correspondances pour un identifiant donné (téléphone ou email)
// à travers l'ensemble des écoles pour résoudre le tenant lors du login sans slug fourni.
func (r *UserRepository) GetByIdentifierGlobal(ctx context.Context, identifier string) ([]*models.User, error) {
	query := `
		SELECT id, auto_ecole_id, agence_id, nom, prenoms, telephone, email,
		       mot_de_passe_hash, code_pin_manager, role, statut, created_at, updated_at
		FROM utilisateurs
		WHERE telephone = $1 OR LOWER(COALESCE(email, '')) = LOWER($1)
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, identifier)
	if err != nil {
		return nil, fmt.Errorf("erreur recherche globale utilisateur: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.ID, &u.AutoEcoleID, &u.AgenceID, &u.Nom, &u.Prenoms, &u.Telephone, &u.Email,
			&u.MotDePasseHash, &u.CodePinManager, &u.Role, &u.Statut, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("erreur lecture ligne utilisateur: %w", err)
		}
		users = append(users, &u)
	}

	return users, rows.Err()
}

// UpdatePasswordByID met à jour le hash du mot de passe d'un utilisateur sans filtre d'auto-école.
func (r *UserRepository) UpdatePasswordByID(ctx context.Context, userID string, passwordHash string) error {
	query := `
		UPDATE utilisateurs
		SET mot_de_passe_hash = $1, updated_at = NOW()
		WHERE id = $2
	`
	res, err := r.db.ExecContext(ctx, query, passwordHash, userID)
	if err != nil {
		return fmt.Errorf("erreur mise à jour mot de passe: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("utilisateur introuvable")
	}
	return nil
}

// VerifyManagerPin vérifie si un code PIN manager correspond à un gérant/directeur actif de l'auto-école.
func (r *UserRepository) VerifyManagerPin(ctx context.Context, schoolID string, pin string) (*models.User, error) {
	query := `
		SELECT id, auto_ecole_id, agence_id, nom, prenoms, telephone, email,
		       mot_de_passe_hash, code_pin_manager, role, statut, created_at, updated_at
		FROM utilisateurs
		WHERE auto_ecole_id = $1 
		  AND role = 'DIRECTEUR' 
		  AND code_pin_manager = $2
		  AND statut = 'ACTIF'
		LIMIT 1
	`
	var u models.User
	err := r.db.QueryRowContext(ctx, query, schoolID, pin).Scan(
		&u.ID, &u.AutoEcoleID, &u.AgenceID, &u.Nom, &u.Prenoms, &u.Telephone, &u.Email,
		&u.MotDePasseHash, &u.CodePinManager, &u.Role, &u.Statut, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erreur vérification code pin: %w", err)
	}
	return &u, nil
}

// GetUsersBySchool liste l'ensemble des collaborateurs d'une auto-école avec le nom de leur agence.
func (r *UserRepository) GetUsersBySchool(ctx context.Context, schoolID string) ([]*models.UserTeamMember, error) {
	query := `
		SELECT u.id, u.auto_ecole_id, u.agence_id,
		       COALESCE(a.nom, 'Siège Social') AS agence_nom,
		       u.nom, u.prenoms, u.telephone, u.email, u.role, u.statut,
		       to_char(u.created_at, 'YYYY-MM-DD HH24:MI') AS created_at
		FROM utilisateurs u
		LEFT JOIN agences a ON a.id = u.agence_id
		WHERE u.auto_ecole_id = $1
		ORDER BY 
			CASE u.role 
				WHEN 'DIRECTEUR' THEN 1 
				WHEN 'SECRETAIRE' THEN 2 
				WHEN 'CAISSIER' THEN 3 
				WHEN 'MONITEUR' THEN 4 
				ELSE 5 
			END,
			u.nom ASC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return nil, fmt.Errorf("erreur listage collaborateurs: %w", err)
	}
	defer rows.Close()

	var list []*models.UserTeamMember
	for rows.Next() {
		var m models.UserTeamMember
		if err := rows.Scan(
			&m.ID, &m.AutoEcoleID, &m.AgenceID, &m.AgenceNom,
			&m.Nom, &m.Prenoms, &m.Telephone, &m.Email,
			&m.Role, &m.Statut, &m.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("erreur scan collaborateur: %w", err)
		}
		list = append(list, &m)
	}
	return list, rows.Err()
}

// Update met à jour les informations d'un utilisateur existant.
func (r *UserRepository) Update(ctx context.Context, u *models.User) error {
	query := `
		UPDATE utilisateurs
		SET nom = $1, prenoms = $2, telephone = $3, email = $4,
		    role = $5, agence_id = $6, statut = $7,
		    code_pin_manager = COALESCE($8, code_pin_manager),
		    updated_at = NOW()
		WHERE id = $9 AND auto_ecole_id = $10
	`
	res, err := r.db.ExecContext(ctx, query,
		u.Nom, u.Prenoms, u.Telephone, u.Email,
		u.Role, u.AgenceID, u.Statut, u.CodePinManager,
		u.ID, u.AutoEcoleID,
	)
	if err != nil {
		return fmt.Errorf("erreur mise à jour utilisateur: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("utilisateur introuvable ou déjà supprimé")
	}
	return nil
}

// UpdatePassword met à jour le mot de passe d'un utilisateur.
func (r *UserRepository) UpdatePassword(ctx context.Context, userID, schoolID, newHash string) error {
	query := `
		UPDATE utilisateurs
		SET mot_de_passe_hash = $1, updated_at = NOW()
		WHERE id = $2 AND auto_ecole_id = $3
	`
	_, err := r.db.ExecContext(ctx, query, newHash, userID, schoolID)
	return err
}

// Delete supprime un utilisateur de la base pour une auto-école donnée.
func (r *UserRepository) Delete(ctx context.Context, userID string, schoolID string) error {
	query := `DELETE FROM utilisateurs WHERE id = $1 AND auto_ecole_id = $2`
	res, err := r.db.ExecContext(ctx, query, userID, schoolID)
	if err != nil {
		return fmt.Errorf("erreur suppression utilisateur: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("utilisateur introuvable ou déjà supprimé")
	}
	return nil
}

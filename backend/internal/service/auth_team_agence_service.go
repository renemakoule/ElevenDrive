package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"autoschool/backend/internal/models"
)

// RegisterUser ajoute un nouveau collaborateur au sein d'une auto-école.
func (s *AuthService) RegisterUser(
	ctx context.Context,
	actor *models.User,
	input *models.RegisterUserInput,
) (*models.UserResponse, error) {
	actorRole := models.NormalizeRole(actor.Role)
	inputRole := models.NormalizeRole(input.Role)

	if actorRole != models.RoleDirecteur {
		if (actorRole == models.RoleSecretaire || actorRole == models.RoleCaissier) && inputRole != models.RoleEleve {
			return nil, errors.New("les secrétaires ne sont autorisées qu'à enregistrer des comptes élèves")
		}
		if actorRole != models.RoleSecretaire && actorRole != models.RoleCaissier {
			return nil, errors.New("droit insuffisant pour créer un utilisateur")
		}
	}

	if !models.IsValidRole(input.Role) {
		return nil, fmt.Errorf("rôle utilisateur invalide: %s", input.Role)
	}
	if strings.TrimSpace(input.Nom) == "" || strings.TrimSpace(input.Prenoms) == "" {
		return nil, errors.New("le nom et prénom sont obligatoires")
	}
	if strings.TrimSpace(input.Telephone) == "" {
		return nil, errors.New("le numéro de téléphone est obligatoire")
	}
	if len(input.MotDePasse) < 6 {
		return nil, errors.New("le mot de passe doit comporter au moins 6 caractères")
	}

	existing, err := s.userRepo.GetByIdentifierInSchool(ctx, actor.AutoEcoleID, input.Telephone)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("un utilisateur avec le téléphone %s existe déjà dans cette auto-école", input.Telephone)
	}

	pwdHash, err := HashPassword(input.MotDePasse)
	if err != nil {
		return nil, fmt.Errorf("erreur hachage mot de passe: %w", err)
	}

	agenceID := input.AgenceID
	if agenceID != nil {
		cleanID := strings.TrimSpace(*agenceID)
		if cleanID == "" || strings.HasPrefix(cleanID, "ag-") || len(cleanID) < 32 {
			agenceID = nil
		} else {
			ag, err := s.schoolRepo.GetAgenceByID(ctx, cleanID)
			if err != nil || ag == nil || ag.AutoEcoleID != actor.AutoEcoleID {
				agenceID = nil
			}
		}
	}

	if agenceID == nil || *agenceID == "" {
		if siege, err := s.schoolRepo.GetSiegeAgence(ctx, actor.AutoEcoleID); err == nil && siege != nil {
			agenceID = &siege.ID
		} else if actor.AgenceID != nil && !strings.HasPrefix(*actor.AgenceID, "ag-") && len(*actor.AgenceID) >= 32 {
			agenceID = actor.AgenceID
		}
	}

	newUser := &models.User{
		ID:             uuid.New().String(),
		AutoEcoleID:    actor.AutoEcoleID,
		AgenceID:       agenceID,
		Nom:            strings.ToUpper(strings.TrimSpace(input.Nom)),
		Prenoms:        strings.TrimSpace(input.Prenoms),
		Telephone:      strings.TrimSpace(input.Telephone),
		Email:          input.Email,
		MotDePasseHash: &pwdHash,
		CodePinManager: input.CodePinManager,
		Role:           models.NormalizeRole(input.Role),
		Statut:         models.StatusActif,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("erreur enregistrement utilisateur: %w", err)
	}

	if s.notifService != nil {
		roleDir := "directeur"
		schoolIDVal := actor.AutoEcoleID
		userName := fmt.Sprintf("%s %s", newUser.Prenoms, newUser.Nom)
		_, _ = s.notifService.NotifyAndBroadcast(ctx, &models.CreateNotificationInput{
			AutoEcoleID: &schoolIDVal,
			RoleCible:   &roleDir,
			Titre:       "👤 Nouveau compte collaborateur créé",
			Message:     fmt.Sprintf("Le collaborateur %s a été ajouté avec le rôle %s.", userName, newUser.Role),
			Type:        "info",
			Module:      "systeme",
		})

		_ = s.notifService.LogAudit(ctx, &models.SystemAuditLog{
			Action:      "CREATION_UTILISATEUR",
			CibleType:   "utilisateur",
			CibleID:     newUser.ID,
			Details:     fmt.Sprintf("Création du compte %s (rôle: %s, tél: %s)", userName, newUser.Role, newUser.Telephone),
			AutoEcoleID: &schoolIDVal,
			UserID:      &actor.ID,
			Role:        &actor.Role,
		})
	}

	return newUser.ToResponse(), nil
}

// GetSchoolUsers renvoie tous les collaborateurs d'une auto-école.
func (s *AuthService) GetSchoolUsers(ctx context.Context, schoolID string) ([]*models.UserTeamMember, error) {
	return s.userRepo.GetUsersBySchool(ctx, schoolID)
}

// UpdateUser modifie les informations d'un collaborateur existant.
func (s *AuthService) UpdateUser(ctx context.Context, claims *models.JWTClaims, userID string, input *models.UpdateUserInput) (*models.UserResponse, error) {
	actorRole := models.NormalizeRole(claims.Role)
	if actorRole != models.RoleDirecteur && actorRole != models.RoleSecretaire {
		return nil, errors.New("droit insuffisant pour modifier un collaborateur")
	}

	targetUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if targetUser == nil || targetUser.AutoEcoleID != claims.AutoEcoleID {
		return nil, errors.New("collaborateur introuvable au sein de cette auto-école")
	}

	if !models.IsValidRole(input.Role) {
		return nil, fmt.Errorf("rôle utilisateur invalide: %s", input.Role)
	}
	if strings.TrimSpace(input.Nom) == "" || strings.TrimSpace(input.Prenoms) == "" {
		return nil, errors.New("le nom et prénom sont obligatoires")
	}
	if strings.TrimSpace(input.Telephone) == "" {
		return nil, errors.New("le numéro de téléphone est obligatoire")
	}

	if targetUser.Telephone != strings.TrimSpace(input.Telephone) {
		existing, err := s.userRepo.GetByIdentifierInSchool(ctx, claims.AutoEcoleID, input.Telephone)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != userID {
			return nil, fmt.Errorf("un utilisateur avec le téléphone %s existe déjà", input.Telephone)
		}
	}

	targetUser.Nom = strings.ToUpper(strings.TrimSpace(input.Nom))
	targetUser.Prenoms = strings.TrimSpace(input.Prenoms)
	targetUser.Telephone = strings.TrimSpace(input.Telephone)
	targetUser.Email = input.Email
	targetUser.Role = models.NormalizeRole(input.Role)
	if input.AgenceID != nil && *input.AgenceID != "" {
		targetUser.AgenceID = input.AgenceID
	}
	if input.Statut != "" {
		targetUser.Statut = strings.ToUpper(input.Statut)
	}

	if err := s.userRepo.Update(ctx, targetUser); err != nil {
		return nil, err
	}

	if strings.TrimSpace(input.MotDePasse) != "" {
		if len(input.MotDePasse) < 6 {
			return nil, errors.New("le nouveau mot de passe doit comporter au moins 6 caractères")
		}
		newHash, err := HashPassword(input.MotDePasse)
		if err != nil {
			return nil, fmt.Errorf("erreur hachage mot de passe: %w", err)
		}
		if err := s.userRepo.UpdatePassword(ctx, userID, claims.AutoEcoleID, newHash); err != nil {
			return nil, err
		}
	}

	return targetUser.ToResponse(), nil
}

// DeleteUser supprime un collaborateur (sauf le directeur lui-même).
func (s *AuthService) DeleteUser(ctx context.Context, claims *models.JWTClaims, userID string) error {
	if models.NormalizeRole(claims.Role) != models.RoleDirecteur {
		return errors.New("seul le directeur est habilité à supprimer un collaborateur")
	}
	if claims.UserID == userID {
		return errors.New("impossible de supprimer votre propre compte directeur")
	}

	targetUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if targetUser == nil || targetUser.AutoEcoleID != claims.AutoEcoleID {
		return errors.New("collaborateur introuvable")
	}

	return s.userRepo.Delete(ctx, userID, claims.AutoEcoleID)
}

// GetSchoolAgences renvoie toutes les agences d'une auto-école.
func (s *AuthService) GetSchoolAgences(ctx context.Context, schoolID string) ([]*models.Agence, error) {
	return s.schoolRepo.GetAgencesBySchoolID(ctx, schoolID)
}

// CreateAgence crée une nouvelle agence pour l'auto-école en base PostgreSQL.
func (s *AuthService) CreateAgence(ctx context.Context, claims *models.JWTClaims, input *models.CreateAgenceInput) (*models.Agence, error) {
	actorRole := models.NormalizeRole(claims.Role)
	if actorRole != models.RoleDirecteur && actorRole != models.RoleSecretaire {
		return nil, errors.New("droit insuffisant pour créer une agence")
	}

	nom := strings.TrimSpace(input.Nom)
	if nom == "" {
		return nil, errors.New("le nom de l'agence est obligatoire")
	}

	codeAgence := strings.TrimSpace(input.CodeAgence)
	if codeAgence == "" {
		codeAgence = fmt.Sprintf("AG-%d", time.Now().Unix()%10000)
	}

	newAgence := &models.Agence{
		ID:          uuid.New().String(),
		AutoEcoleID: claims.AutoEcoleID,
		CodeAgence:  codeAgence,
		Nom:         nom,
		Adresse:     strings.TrimSpace(input.Adresse),
		Telephone:   strings.TrimSpace(input.Telephone),
		Ville:       strings.TrimSpace(input.Ville),
		Commune:     strings.TrimSpace(input.Commune),
		EstSiege:    input.EstSiege,
	}

	if err := s.schoolRepo.CreateAgence(ctx, newAgence); err != nil {
		return nil, err
	}

	return newAgence, nil
}

// UpdateAgence modifie une succursale existante.
func (s *AuthService) UpdateAgence(ctx context.Context, claims *models.JWTClaims, agenceID string, input *models.UpdateAgenceInput) (*models.Agence, error) {
	actorRole := models.NormalizeRole(claims.Role)
	if actorRole != models.RoleDirecteur && actorRole != models.RoleSecretaire {
		return nil, errors.New("droit insuffisant pour modifier une agence")
	}

	existing, err := s.schoolRepo.GetAgenceByID(ctx, agenceID)
	if err != nil {
		return nil, err
	}
	if existing == nil || existing.AutoEcoleID != claims.AutoEcoleID {
		return nil, errors.New("agence introuvable au sein de cette auto-école")
	}

	nom := strings.TrimSpace(input.Nom)
	if nom == "" {
		return nil, errors.New("le nom de l'agence est obligatoire")
	}

	existing.Nom = nom
	if strings.TrimSpace(input.CodeAgence) != "" {
		existing.CodeAgence = strings.TrimSpace(input.CodeAgence)
	}
	existing.Adresse = strings.TrimSpace(input.Adresse)
	existing.Telephone = strings.TrimSpace(input.Telephone)
	existing.Ville = strings.TrimSpace(input.Ville)
	existing.Commune = strings.TrimSpace(input.Commune)
	existing.EstSiege = input.EstSiege

	if err := s.schoolRepo.UpdateAgence(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// DeleteAgence supprime une agence non-siège.
func (s *AuthService) DeleteAgence(ctx context.Context, claims *models.JWTClaims, agenceID string) error {
	if models.NormalizeRole(claims.Role) != models.RoleDirecteur {
		return errors.New("seul le directeur peut supprimer une agence")
	}

	return s.schoolRepo.DeleteAgence(ctx, claims.AutoEcoleID, agenceID)
}


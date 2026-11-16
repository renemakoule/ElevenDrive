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

// RefreshToken valide un refresh token, révoque l'ancienne session et émet de nouveaux jetons.
func (s *AuthService) RefreshToken(
	ctx context.Context,
	refreshToken string,
	clientIP, userAgent string,
) (*models.AuthResponse, error) {
	tokenHash := HashRefreshToken(refreshToken)

	session, err := s.sessRepo.GetValidByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, errors.New("session invalide ou expirée, veuillez vous reconnecter")
	}

	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil || user == nil {
		return nil, errors.New("utilisateur introuvable")
	}
	if user.Statut != models.StatusActif {
		return nil, errors.New("compte suspendu ou inactif")
	}

	school, err := s.schoolRepo.GetByID(ctx, session.AutoEcoleID)
	if err != nil || school == nil {
		return nil, errors.New("auto-école introuvable")
	}

	var agence *models.Agence
	if user.AgenceID != nil {
		agence, _ = s.schoolRepo.GetAgenceByID(ctx, *user.AgenceID)
	}

	_ = s.sessRepo.Revoke(ctx, tokenHash)

	return s.createSessionAndResponse(ctx, user, school, agence, clientIP, userAgent)
}

// Logout révoque la session active associée au refresh token.
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := HashRefreshToken(refreshToken)
	return s.sessRepo.Revoke(ctx, tokenHash)
}

// VerifyPin vérifie la validité d'un code PIN manager pour autoriser une action critique au guichet.
func (s *AuthService) VerifyPin(ctx context.Context, schoolID string, pin string) (*models.UserResponse, error) {
	if strings.TrimSpace(pin) == "" {
		return nil, errors.New("code pin requis")
	}
	manager, err := s.userRepo.VerifyManagerPin(ctx, schoolID, pin)
	if err != nil {
		return nil, err
	}
	if manager == nil {
		return nil, errors.New("code pin manager invalide ou non autorisé")
	}
	return manager.ToResponse(), nil
}

// GetProfile renvoie le profil utilisateur avec son auto-école et son agence.
func (s *AuthService) GetProfile(ctx context.Context, userID string) (*models.UserResponse, *models.AutoEcole, *models.Agence, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, nil, nil, errors.New("utilisateur introuvable")
	}

	school, err := s.schoolRepo.GetByID(ctx, user.AutoEcoleID)
	if err != nil {
		return nil, nil, nil, err
	}

	var agence *models.Agence
	if user.AgenceID != nil {
		agence, _ = s.schoolRepo.GetAgenceByID(ctx, *user.AgenceID)
	}

	return user.ToResponse(), school, agence, nil
}

// UpdateMyProfile met à jour le profil de l'utilisateur authentifié selon son rôle et ses permissions.
func (s *AuthService) UpdateMyProfile(ctx context.Context, claims *models.JWTClaims, input *models.UpdateUserInput) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil || user == nil {
		return nil, errors.New("utilisateur introuvable")
	}

	if strings.TrimSpace(input.Nom) != "" {
		user.Nom = strings.ToUpper(strings.TrimSpace(input.Nom))
	}
	if strings.TrimSpace(input.Prenoms) != "" {
		user.Prenoms = strings.TrimSpace(input.Prenoms)
	}
	if strings.TrimSpace(input.Telephone) != "" {
		newTel := strings.TrimSpace(input.Telephone)
		if newTel != user.Telephone {
			existing, err := s.userRepo.GetByIdentifierInSchool(ctx, user.AutoEcoleID, newTel)
			if err != nil {
				return nil, err
			}
			if existing != nil && existing.ID != user.ID {
				return nil, fmt.Errorf("ce numéro de téléphone (%s) est déjà utilisé par un autre compte", newTel)
			}
			user.Telephone = newTel
		}
	}
	if input.Email != nil {
		cleanEmail := strings.TrimSpace(*input.Email)
		if cleanEmail == "" {
			user.Email = nil
		} else {
			user.Email = &cleanEmail
		}
	}

	userRole := models.NormalizeRole(user.Role)
	if input.CodePinManager != nil && strings.TrimSpace(*input.CodePinManager) != "" {
		if userRole != models.RoleDirecteur {
			return nil, errors.New("seul le Directeur est habilité à définir un code PIN Manager")
		}
		pin := strings.TrimSpace(*input.CodePinManager)
		user.CodePinManager = &pin
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
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
		if err := s.userRepo.UpdatePassword(ctx, user.ID, user.AutoEcoleID, newHash); err != nil {
			return nil, err
		}
	}

	return user.ToResponse(), nil
}

// UpdateSchoolInfo met à jour les informations de l'établissement (strictement réservé au Directeur).
func (s *AuthService) UpdateSchoolInfo(ctx context.Context, claims *models.JWTClaims, input *models.AutoEcole) (*models.AutoEcole, error) {
	if models.NormalizeRole(claims.Role) != models.RoleDirecteur {
		return nil, errors.New("seul le Directeur est habilité à modifier les informations officielles de l'auto-école")
	}

	school, err := s.schoolRepo.GetByID(ctx, claims.AutoEcoleID)
	if err != nil || school == nil {
		return nil, errors.New("auto-école introuvable")
	}

	if strings.TrimSpace(input.NomCommercial) != "" {
		school.NomCommercial = strings.TrimSpace(input.NomCommercial)
	}
	if strings.TrimSpace(input.CodeUnique) != "" {
		newSlug := Slugify(input.CodeUnique)
		if !strings.EqualFold(newSlug, school.CodeUnique) {
			existing, err := s.schoolRepo.GetByCodeUnique(ctx, newSlug)
			if err != nil {
				return nil, err
			}
			if existing != nil && existing.ID != school.ID {
				return nil, fmt.Errorf("l'identifiant web '%s' est déjà pris par une autre auto-école", newSlug)
			}
			school.CodeUnique = newSlug
		}
	}
	if strings.TrimSpace(input.NumeroAgrement) != "" {
		school.NumeroAgrement = strings.TrimSpace(input.NumeroAgrement)
	}
	if strings.TrimSpace(input.Telephone) != "" {
		school.Telephone = strings.TrimSpace(input.Telephone)
	}
	if input.Email != nil {
		school.Email = input.Email
	}
	if strings.TrimSpace(input.Devise) != "" {
		school.Devise = strings.TrimSpace(input.Devise)
	}
	if strings.TrimSpace(input.Pays) != "" {
		school.Pays = strings.TrimSpace(input.Pays)
	}

	if err := s.schoolRepo.UpdateSchool(ctx, school); err != nil {
		return nil, err
	}
	return school, nil
}

// GetSchoolInfoBySlug récupère les métadonnées publiques d'une auto-école pour les liens d'accès.
func (s *AuthService) GetSchoolInfoBySlug(ctx context.Context, slug string) (*models.SchoolPublicInfo, error) {
	school, err := s.schoolRepo.GetByCodeUnique(ctx, slug)
	if err != nil {
		return nil, err
	}
	if school == nil {
		return nil, errors.New("auto-école introuvable")
	}

	agences, err := s.schoolRepo.GetAgencesBySchoolID(ctx, school.ID)
	if err != nil {
		return nil, err
	}

	return &models.SchoolPublicInfo{
		ID:             school.ID,
		CodeUnique:     school.CodeUnique,
		NomCommercial:  school.NomCommercial,
		NumeroAgrement: school.NumeroAgrement,
		Telephone:      school.Telephone,
		Email:          school.Email,
		Devise:         school.Devise,
		Pays:           school.Pays,
		LogoURL:        school.LogoURL,
		Agences:        agences,
	}, nil
}

// createSessionAndResponse génère les tokens et insère la session PostgreSQL.
func (s *AuthService) createSessionAndResponse(
	ctx context.Context,
	user *models.User,
	school *models.AutoEcole,
	agence *models.Agence,
	clientIP, userAgent string,
) (*models.AuthResponse, error) {
	accessToken, expiresIn, err := GenerateAccessToken(
		user,
		school.ID,
		user.AgenceID,
		s.cfg.JWTSecret,
		s.cfg.JWTExpirationHours,
	)
	if err != nil {
		return nil, err
	}

	plainRefreshToken, tokenHash, err := GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().AddDate(0, 0, s.cfg.RefreshExpirationDay)

	var ipPtr, uaPtr *string
	if clientIP != "" {
		ipPtr = &clientIP
	}
	if userAgent != "" {
		uaPtr = &userAgent
	}

	session := &models.AuthSession{
		ID:          uuid.New().String(),
		UserID:      user.ID,
		AutoEcoleID: school.ID,
		TokenHash:   tokenHash,
		IPAddress:   ipPtr,
		UserAgent:   uaPtr,
		ExpiresAt:   expiresAt,
	}

	if err := s.sessRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("erreur persistance session: %w", err)
	}

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: plainRefreshToken,
		ExpiresIn:    expiresIn,
		User:         user.ToResponse(),
		School:       school,
		Agence:       agence,
	}, nil
}


package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"autoschool/backend/internal/models"
)

// OAuthLoginInput représente la requête d'authentification OAuth d'un utilisateur.
type OAuthLoginInput struct {
	Provider string `json:"provider"`
	Email    string `json:"email"`
	Nom      string `json:"nom"`
	Prenoms  string `json:"prenoms"`
	IDToken  string `json:"id_token"`
	SchoolID string `json:"school_id,omitempty"`
}

// OAuthResult expose le statut de la requête OAuth : compte existant connecté ou inscription requise.
type OAuthResult struct {
	Status       string               `json:"status"` // "authenticated" | "needs_registration" | "select_school"
	AuthResponse *models.AuthResponse `json:"auth_response,omitempty"`
	Provider     string               `json:"provider"`
	Email        string               `json:"email"`
	Nom          string               `json:"nom"`
	Prenoms      string               `json:"prenoms"`
	Message      string               `json:"message,omitempty"`
	Ticket       string               `json:"ticket,omitempty"`
	Schools      []SchoolChoice       `json:"schools,omitempty"`
}

// OAuthUserInfo contient les informations officielles retournées par le fournisseur OAuth.
type OAuthUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

func verifyOAuthToken(ctx context.Context, provider, idToken, expectedEmail string) error {
	if strings.TrimSpace(idToken) == "" {
		if os.Getenv("ENV") != "production" {
			return nil
		}
		return errors.New("jeton d'identité OAuth (id_token) obligatoire pour toute authentification sociale")
	}

	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "google":
		client := &http.Client{Timeout: 10 * time.Second}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://oauth2.googleapis.com/tokeninfo?id_token="+url.QueryEscape(idToken), nil)
		if err != nil {
			return fmt.Errorf("erreur initialisation requête validation OAuth: %w", err)
		}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("impossible de joindre le service de vérification Google: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return errors.New("jeton Google OAuth invalide, corrompu ou expiré")
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
			return errors.New("réponse du fournisseur Google OAuth illisible")
		}

		verifiedEmail, _ := payload["email"].(string)
		if !strings.EqualFold(strings.TrimSpace(verifiedEmail), strings.TrimSpace(expectedEmail)) {
			return errors.New("l'adresse email certifiée par Google ne correspond pas à l'email réclamé")
		}

		isVerified := false
		if v, ok := payload["email_verified"].(bool); ok {
			isVerified = v
		} else if s, ok := payload["email_verified"].(string); ok {
			isVerified = strings.EqualFold(s, "true")
		}
		if !isVerified {
			return errors.New("l'adresse email n'a pas été confirmée par votre compte Google")
		}
		return nil

	default:
		parts := strings.Split(idToken, ".")
		if len(parts) != 3 {
			return errors.New("format de jeton OAuth invalide (structure JWT attendue)")
		}
		return nil
	}
}

// LoginWithOAuth vérifie si l'utilisateur existe déjà via son email social.
func (s *AuthService) LoginWithOAuth(
	ctx context.Context,
	input *OAuthLoginInput,
	clientIP, userAgent string,
) (*OAuthResult, error) {
	provider := strings.ToLower(strings.TrimSpace(input.Provider))
	email := strings.ToLower(strings.TrimSpace(input.Email))
	idToken := strings.TrimSpace(input.IDToken)

	if email != "" {
		users, err := s.userRepo.GetByIdentifierGlobal(ctx, email)
		if err == nil && len(users) > 0 {
			if errToken := verifyOAuthToken(ctx, provider, idToken, email); errToken != nil {
				return nil, fmt.Errorf("échec de validation OAuth : %w", errToken)
			}

			candidates := filterOAuthCandidates(users, input.SchoolID)
			if len(candidates) > 1 {
				return s.buildSchoolChoice(ctx, provider, email, candidates)
			}
			targetUser := users[0]
			if len(candidates) == 1 {
				targetUser = candidates[0]
			}
			if targetUser.Statut == models.StatusActif {
				school, errSchool := s.schoolRepo.GetByID(ctx, targetUser.AutoEcoleID)
				if errSchool == nil && school != nil {
					if !school.EstActif {
						return nil, fmt.Errorf("l'établissement '%s' est actuellement suspendu par le Super-Administrateur", school.NomCommercial)
					}
					if school.DateExpirationLicence != nil && time.Now().After(*school.DateExpirationLicence) {
						return nil, fmt.Errorf("la licence de l'établissement '%s' a expiré le %s. Veuillez contacter le Super-Administrateur", school.NomCommercial, school.DateExpirationLicence.Format("02/01/2006"))
					}

					var targetAgence *models.Agence
					if targetUser.AgenceID != nil {
						targetAgence, _ = s.schoolRepo.GetAgenceByID(ctx, *targetUser.AgenceID)
					}
					authResp, errResp := s.createSessionAndResponse(ctx, targetUser, school, targetAgence, clientIP, userAgent)
					if errResp == nil {
						return &OAuthResult{
							Status:       "authenticated",
							AuthResponse: authResp,
							Provider:     provider,
							Email:        email,
							Nom:          targetUser.Nom,
							Prenoms:      targetUser.Prenoms,
							Message:      "Connexion réussie via votre compte social authentifié.",
						}, nil
					}
				}
			}
		}
	}

	nom := strings.TrimSpace(input.Nom)
	prenoms := strings.TrimSpace(input.Prenoms)
	if nom == "" && prenoms == "" && email != "" {
		parts := strings.Split(email, "@")
		prenoms = strings.Title(parts[0])
	}

	return &OAuthResult{
		Status:   "needs_registration",
		Provider: provider,
		Email:    email,
		Nom:      nom,
		Prenoms:  prenoms,
		Message:  "Aucun compte auto-école n'est associé à cette identité. Veuillez renseigner le nom de votre établissement pour finaliser la création.",
	}, nil
}

// GetOAuthAuthURL génère l'URL d'autorisation officielle du fournisseur (Google, Facebook, LinkedIn, Apple)
func (s *AuthService) resolveOAuthRedirectURI(configured, provider string) (string, error) {
	if strings.TrimSpace(configured) != "" {
		return strings.TrimSpace(configured), nil
	}
	if strings.TrimSpace(s.cfg.PublicAPIURL) != "" {
		return fmt.Sprintf("%s/api/auth/%s/callback", s.cfg.PublicAPIURL, provider), nil
	}
	return "", fmt.Errorf("PUBLIC_API_URL ou %s_REDIRECT_URI doit être défini dans le fichier .env", strings.ToUpper(provider))
}

func (s *AuthService) GetOAuthAuthURL(provider, state string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "google":
		if s.cfg.GoogleClientID == "" || strings.HasPrefix(s.cfg.GoogleClientID, "votre_") {
			return "", errors.New("GOOGLE_CLIENT_ID non configuré dans le fichier .env")
		}
		redirectURI, err := s.resolveOAuthRedirectURI(s.cfg.GoogleRedirectURI, "google")
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(
			"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=openid%%20profile%%20email&access_type=offline&prompt=select_account&state=%s",
			url.QueryEscape(s.cfg.GoogleClientID),
			url.QueryEscape(redirectURI),
			url.QueryEscape(state),
		), nil

	case "facebook":
		if s.cfg.FacebookAppID == "" || strings.HasPrefix(s.cfg.FacebookAppID, "votre_") {
			return "", errors.New("FACEBOOK_APP_ID non configuré dans le fichier .env")
		}
		redirectURI, err := s.resolveOAuthRedirectURI(s.cfg.FacebookRedirectURI, "facebook")
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(
			"https://www.facebook.com/v19.0/dialog/oauth?client_id=%s&redirect_uri=%s&scope=email,public_profile&state=%s",
			url.QueryEscape(s.cfg.FacebookAppID),
			url.QueryEscape(redirectURI),
			url.QueryEscape(state),
		), nil

	case "linkedin":
		if s.cfg.LinkedInClientID == "" || strings.HasPrefix(s.cfg.LinkedInClientID, "votre_") {
			return "", errors.New("LINKEDIN_CLIENT_ID non configuré dans le fichier .env")
		}
		redirectURI, err := s.resolveOAuthRedirectURI(s.cfg.LinkedInRedirectURI, "linkedin")
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(
			"https://www.linkedin.com/oauth/v2/authorization?response_type=code&client_id=%s&redirect_uri=%s&scope=openid%%20profile%%20email&state=%s",
			url.QueryEscape(s.cfg.LinkedInClientID),
			url.QueryEscape(redirectURI),
			url.QueryEscape(state),
		), nil

	case "apple":
		if s.cfg.AppleClientID == "" || strings.HasPrefix(s.cfg.AppleClientID, "com.votreautoecole") {
			return "", errors.New("APPLE_CLIENT_ID non configuré dans le fichier .env")
		}
		redirectURI, err := s.resolveOAuthRedirectURI(s.cfg.AppleRedirectURI, "apple")
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(
			"https://appleid.apple.com/auth/authorize?client_id=%s&redirect_uri=%s&response_type=code&response_mode=query&scope=name%%20email&state=%s",
			url.QueryEscape(s.cfg.AppleClientID),
			url.QueryEscape(redirectURI),
			url.QueryEscape(state),
		), nil

	default:
		return "", fmt.Errorf("fournisseur OAuth '%s' non pris en charge", provider)
	}
}

// HandleOAuthCallback traite la réponse du fournisseur OAuth et connecte ou prépare l'inscription de l'utilisateur
func (s *AuthService) HandleOAuthCallback(
	ctx context.Context,
	provider, code string,
	clientIP, userAgent string,
) (*OAuthResult, error) {
	userInfo, idToken, err := s.ExchangeOAuthCode(ctx, provider, code)
	if err != nil {
		return nil, err
	}

	loginInput := &OAuthLoginInput{
		Provider: provider,
		Email:    userInfo.Email,
		Nom:      userInfo.FamilyName,
		Prenoms:  userInfo.GivenName,
		IDToken:  idToken,
	}

	return s.LoginWithOAuth(ctx, loginInput, clientIP, userAgent)
}


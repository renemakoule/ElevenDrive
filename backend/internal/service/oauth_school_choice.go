package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"autoschool/backend/internal/models"
)

// SchoolChoice décrit un compte (école + agence) proposé à l'utilisateur
// lorsque la même identité sociale correspond à plusieurs établissements.
type SchoolChoice struct {
	UserID        string `json:"user_id"`
	SchoolID      string `json:"school_id"`
	NomCommercial string `json:"nom_commercial"`
	CodeUnique    string `json:"code_unique"`
	AgenceID      string `json:"agence_id,omitempty"`
	AgenceNom     string `json:"agence_nom,omitempty"`
	Role          string `json:"role"`
}

type pendingSchoolChoice struct {
	Email     string
	Provider  string
	UserIDs   map[string]bool
	ExpiresAt time.Time
}

var oauthChoiceTickets sync.Map // ticket -> *pendingSchoolChoice

const oauthChoiceTTL = 5 * time.Minute

func newChoiceTicket() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// filterOAuthCandidates ne garde que les comptes actifs (et ceux de l'école demandée si schoolID est fourni).
func filterOAuthCandidates(users []*models.User, schoolID string) []*models.User {
	var out []*models.User
	for _, u := range users {
		if u.Statut != models.StatusActif {
			continue
		}
		if schoolID != "" && u.AutoEcoleID != schoolID {
			continue
		}
		out = append(out, u)
	}
	return out
}

// buildSchoolChoice crée un ticket à usage unique (5 min) et la liste des écoles/agences à proposer.
// La vérification cryptographique du jeton du provider a déjà été faite avant l'appel.
func (s *AuthService) buildSchoolChoice(ctx context.Context, provider, email string, users []*models.User) (*OAuthResult, error) {
	ticket, err := newChoiceTicket()
	if err != nil {
		return nil, err
	}

	pending := &pendingSchoolChoice{
		Email:     email,
		Provider:  provider,
		UserIDs:   map[string]bool{},
		ExpiresAt: time.Now().Add(oauthChoiceTTL),
	}

	choices := make([]SchoolChoice, 0, len(users))
	for _, u := range users {
		school, err := s.schoolRepo.GetByID(ctx, u.AutoEcoleID)
		if err != nil || school == nil {
			continue
		}
		c := SchoolChoice{
			UserID:        u.ID,
			SchoolID:      school.ID,
			NomCommercial: school.NomCommercial,
			CodeUnique:    school.CodeUnique,
			Role:          string(u.Role),
		}
		if u.AgenceID != nil {
			if ag, _ := s.schoolRepo.GetAgenceByID(ctx, *u.AgenceID); ag != nil {
				c.AgenceID = ag.ID
				c.AgenceNom = ag.Nom
			}
		}
		pending.UserIDs[u.ID] = true
		choices = append(choices, c)
	}

	if len(choices) == 0 {
		return nil, errors.New("aucune auto-école active associée à ce compte")
	}

	oauthChoiceTickets.Store(ticket, pending)

	return &OAuthResult{
		Status:   "select_school",
		Provider: provider,
		Email:    email,
		Ticket:   ticket,
		Schools:  choices,
		Message:  "Plusieurs auto-écoles sont associées à ce compte. Choisissez l'établissement (et l'agence) à ouvrir.",
	}, nil
}

// CompleteOAuthSchoolChoice finalise la connexion après le choix de l'école par l'utilisateur.
func (s *AuthService) CompleteOAuthSchoolChoice(
	ctx context.Context,
	ticket, userID string,
	clientIP, userAgent string,
) (*models.AuthResponse, error) {
	v, ok := oauthChoiceTickets.LoadAndDelete(ticket) // usage unique
	if !ok {
		return nil, errors.New("choix expiré ou invalide, veuillez vous reconnecter")
	}
	pending := v.(*pendingSchoolChoice)
	if time.Now().After(pending.ExpiresAt) {
		return nil, errors.New("choix expiré, veuillez vous reconnecter")
	}
	if !pending.UserIDs[userID] {
		return nil, errors.New("établissement non autorisé pour ce compte")
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil || user.Statut != models.StatusActif {
		return nil, errors.New("compte introuvable ou inactif")
	}
	school, err := s.schoolRepo.GetByID(ctx, user.AutoEcoleID)
	if err != nil {
		return nil, err
	}
	if school == nil {
		return nil, errors.New("auto-école introuvable")
	}
	var agence *models.Agence
	if user.AgenceID != nil {
		agence, _ = s.schoolRepo.GetAgenceByID(ctx, *user.AgenceID)
	}
	return s.createSessionAndResponse(ctx, user, school, agence, clientIP, userAgent)
}


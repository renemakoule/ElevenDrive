package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"autoschool/backend/internal/service"
)

// OAuthDesktopSessionResult stocke l'état d'un flux OAuth initié par l'application Desktop
type OAuthDesktopSessionResult struct {
	Status       string                 `json:"status"` // "pending" | "authenticated" | "needs_registration" | "error"
	Token        string                 `json:"token,omitempty"`
	RefreshToken string                 `json:"refresh_token,omitempty"`
	Role         string                 `json:"role,omitempty"`
	Provider     string                 `json:"provider,omitempty"`
	Email        string                 `json:"email,omitempty"`
	Nom          string                 `json:"nom,omitempty"`
	Prenoms      string                 `json:"prenoms,omitempty"`
	Error        string                 `json:"error,omitempty"`
	Ticket       string                 `json:"ticket,omitempty"`
	Schools      []service.SchoolChoice `json:"schools,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
}

var oauthDesktopSessions sync.Map

// OAuth gère l'authentification via un fournisseur social (Google, Facebook, LinkedIn, Apple).
// POST /api/auth/oauth
func (h *AuthHandler) OAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	var input service.OAuthLoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps de requête JSON invalide")
		return
	}

	res, err := h.authService.LoginWithOAuth(r.Context(), &input, getClientIP(r), r.UserAgent())
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	Success(w, http.StatusOK, res)
}

// OAuthSelectSchool finalise une connexion sociale après choix de l'école/agence.
// POST /api/auth/oauth/select-school  {ticket, user_id}
func (h *AuthHandler) OAuthSelectSchool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}
	var input struct {
		Ticket string `json:"ticket"`
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps de requête JSON invalide")
		return
	}
	resp, err := h.authService.CompleteOAuthSchoolChoice(r.Context(), input.Ticket, input.UserID, getClientIP(r), r.UserAgent())
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}
	Success(w, http.StatusOK, resp)
}

func (h *AuthHandler) GoogleAuthRedirect(w http.ResponseWriter, r *http.Request) {
	h.handleOAuthRedirect(w, r, "google")
}

func (h *AuthHandler) GoogleAuthCallback(w http.ResponseWriter, r *http.Request) {
	h.handleOAuthCallback(w, r, "google")
}

func (h *AuthHandler) FacebookAuthRedirect(w http.ResponseWriter, r *http.Request) {
	h.handleOAuthRedirect(w, r, "facebook")
}

func (h *AuthHandler) FacebookAuthCallback(w http.ResponseWriter, r *http.Request) {
	h.handleOAuthCallback(w, r, "facebook")
}

func (h *AuthHandler) LinkedInAuthRedirect(w http.ResponseWriter, r *http.Request) {
	h.handleOAuthRedirect(w, r, "linkedin")
}

func (h *AuthHandler) LinkedInAuthCallback(w http.ResponseWriter, r *http.Request) {
	h.handleOAuthCallback(w, r, "linkedin")
}

func (h *AuthHandler) AppleAuthRedirect(w http.ResponseWriter, r *http.Request) {
	h.handleOAuthRedirect(w, r, "apple")
}

func (h *AuthHandler) AppleAuthCallback(w http.ResponseWriter, r *http.Request) {
	h.handleOAuthCallback(w, r, "apple")
}

func (h *AuthHandler) handleOAuthRedirect(w http.ResponseWriter, r *http.Request, provider string) {
	sessionID := r.URL.Query().Get("session_id")
	redirectOrigin := r.URL.Query().Get("origin")
	if redirectOrigin == "" {
		redirectOrigin = r.URL.Query().Get("redirect_origin")
	}
	if redirectOrigin == "" {
		if referer := r.Header.Get("Referer"); referer != "" {
			if u, err := url.Parse(referer); err == nil {
				redirectOrigin = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
			}
		}
	}
	if redirectOrigin == "" {
		redirectOrigin = h.authService.GetFrontendURL()
	}

	stateData := url.Values{}
	stateData.Set("origin", redirectOrigin)
	if sessionID != "" {
		stateData.Set("session_id", sessionID)
		oauthDesktopSessions.Store(sessionID, &OAuthDesktopSessionResult{
			Status:    "pending",
			CreatedAt: time.Now(),
		})
	}

	authURL, err := h.authService.GetOAuthAuthURL(provider, stateData.Encode())
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// GetOAuthSession vérifie l'état d'une session OAuth initiée depuis l'application Desktop
func (h *AuthHandler) GetOAuthSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		Error(w, http.StatusBadRequest, "Le paramètre 'session_id' est obligatoire")
		return
	}

	val, ok := oauthDesktopSessions.Load(sessionID)
	if !ok {
		Success(w, http.StatusOK, map[string]interface{}{"status": "pending"})
		return
	}

	res := val.(*OAuthDesktopSessionResult)
	if res.Status == "authenticated" || res.Status == "needs_registration" || res.Status == "select_school" || res.Status == "error" {
		oauthDesktopSessions.Delete(sessionID)
	}

	Success(w, http.StatusOK, res)
}

func (h *AuthHandler) resolveOAuthCallbackState(r *http.Request) (string, string) {
	state := r.URL.Query().Get("state")
	redirectOrigin := h.authService.GetFrontendURL()
	sessionID := r.URL.Query().Get("session_id")

	if state == "" {
		return redirectOrigin, sessionID
	}

	if vals, err := url.ParseQuery(state); err == nil && len(vals) > 0 {
		if sID := vals.Get("session_id"); sID != "" {
			sessionID = sID
		}
		if o := vals.Get("origin"); o != "" {
			redirectOrigin = strings.TrimRight(o, "/")
		}
	}
	if sessionID == "" {
		if unescaped, err := url.QueryUnescape(state); err == nil && unescaped != "" {
			if vals, err2 := url.ParseQuery(unescaped); err2 == nil && len(vals) > 0 {
				if sID := vals.Get("session_id"); sID != "" {
					sessionID = sID
				}
				if o := vals.Get("origin"); o != "" {
					redirectOrigin = strings.TrimRight(o, "/")
				}
			} else if strings.HasPrefix(unescaped, "http") {
				redirectOrigin = strings.TrimRight(unescaped, "/")
			}
		}
	}
	return redirectOrigin, sessionID
}

func (h *AuthHandler) handleOAuthCallback(w http.ResponseWriter, r *http.Request, provider string) {
	code := r.URL.Query().Get("code")
	oauthErr := r.URL.Query().Get("error")
	redirectOrigin, sessionID := h.resolveOAuthCallbackState(r)
	isDesktop := sessionID != ""

	renderOAuthErr := func(msg string) {
		if sessionID != "" {
			oauthDesktopSessions.Store(sessionID, &OAuthDesktopSessionResult{
				Status:    "error",
				Error:     msg,
				CreatedAt: time.Now(),
			})
		}
		target := fmt.Sprintf("%s/?oauth_status=error&oauth_error=%s", redirectOrigin, url.QueryEscape(msg))
		renderOAuthResponse(w, target, map[string]interface{}{"status": "error", "error": msg}, isDesktop)
	}

	if oauthErr != "" {
		renderOAuthErr(oauthErr)
		return
	}
	if code == "" {
		renderOAuthErr("Code d'autorisation manquant")
		return
	}

	result, err := h.authService.HandleOAuthCallback(r.Context(), provider, code, getClientIP(r), r.UserAgent())
	if err != nil {
		renderOAuthErr(err.Error())
		return
	}

	if result.Status == "authenticated" && result.AuthResponse != nil {
		role := string(result.AuthResponse.User.Role)
		if sessionID != "" {
			oauthDesktopSessions.Store(sessionID, &OAuthDesktopSessionResult{
				Status:       "authenticated",
				Token:        result.AuthResponse.AccessToken,
				RefreshToken: result.AuthResponse.RefreshToken,
				Role:         role,
				Provider:     provider,
				CreatedAt:    time.Now(),
			})
		}
		target := fmt.Sprintf("%s/?oauth_status=authenticated&token=%s&refresh_token=%s&role=%s&provider=%s",
			redirectOrigin,
			url.QueryEscape(result.AuthResponse.AccessToken),
			url.QueryEscape(result.AuthResponse.RefreshToken),
			url.QueryEscape(role),
			url.QueryEscape(provider),
		)
		renderOAuthResponse(w, target, map[string]interface{}{
			"status":        "authenticated",
			"token":         result.AuthResponse.AccessToken,
			"refresh_token": result.AuthResponse.RefreshToken,
			"role":          role,
			"user":          result.AuthResponse.User,
			"school":        result.AuthResponse.School,
		}, isDesktop)
		return
	}

	if result.Status == "select_school" {
		if sessionID != "" {
			oauthDesktopSessions.Store(sessionID, &OAuthDesktopSessionResult{
				Status:    "select_school",
				Provider:  provider,
				Email:     result.Email,
				Ticket:    result.Ticket,
				Schools:   result.Schools,
				CreatedAt: time.Now(),
			})
		}
		schoolsJSON, _ := json.Marshal(result.Schools)
		target := fmt.Sprintf("%s/?oauth_status=select_school&provider=%s&ticket=%s&schools=%s",
			redirectOrigin,
			url.QueryEscape(provider),
			url.QueryEscape(result.Ticket),
			url.QueryEscape(string(schoolsJSON)),
		)
		renderOAuthResponse(w, target, map[string]interface{}{
			"status":   "select_school",
			"provider": provider,
			"ticket":   result.Ticket,
			"schools":  result.Schools,
		}, isDesktop)
		return
	}

	if sessionID != "" {
		oauthDesktopSessions.Store(sessionID, &OAuthDesktopSessionResult{
			Status:    "needs_registration",
			Provider:  provider,
			Email:     result.Email,
			Nom:       result.Nom,
			Prenoms:   result.Prenoms,
			CreatedAt: time.Now(),
		})
	}
	target := fmt.Sprintf("%s/?oauth_status=needs_registration&provider=%s&email=%s&nom=%s&prenoms=%s",
		redirectOrigin,
		url.QueryEscape(provider),
		url.QueryEscape(result.Email),
		url.QueryEscape(result.Nom),
		url.QueryEscape(result.Prenoms),
	)
	renderOAuthResponse(w, target, map[string]interface{}{
		"status":   "needs_registration",
		"provider": provider,
		"email":    result.Email,
		"nom":      result.Nom,
		"prenoms":  result.Prenoms,
	}, isDesktop)
}

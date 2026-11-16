package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"autoschool/backend/internal/middleware"
	"autoschool/backend/internal/models"
	"autoschool/backend/internal/service"
)

// AuthHandler expose les contrôleurs HTTP liés à l'authentification.
type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// getClientIP extrait l'adresse IP cliente réelle depuis les en-têtes de proxy ou la connexion.
func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}
	ipPort := r.RemoteAddr
	if idx := strings.LastIndex(ipPort, ":"); idx != -1 {
		return ipPort[:idx]
	}
	return ipPort
}

// RegisterSchool gère l'onboarding complet d'une nouvelle auto-école.
// POST /api/auth/register-school
func (h *AuthHandler) RegisterSchool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	var input models.RegisterSchoolInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		fmt.Printf("❌ [RegisterSchool] Erreur décodage JSON: %v\n", err)
		Error(w, http.StatusBadRequest, "Corps de requête JSON invalide: "+err.Error())
		return
	}

	ip := getClientIP(r)
	ua := r.UserAgent()

	fmt.Printf("📥 [RegisterSchool] Début inscription: école=%q, agrément=%q, tel=%q, pays=%q, ville=%q, dir=%s %s\n",
		input.NomCommercial, input.NumeroAgrement, input.TelephoneDirecteur, input.Pays, input.VilleAgence, input.PrenomsDirecteur, input.NomDirecteur)

	resp, err := h.authService.RegisterSchool(r.Context(), &input, ip, ua)
	if err != nil {
		fmt.Printf("❌ [RegisterSchool] Échec: %v\n", err)
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Printf("✅ [RegisterSchool] Inscription réussie: école=%s (id=%s, slug=%s, admin=%s)\n",
		resp.School.NomCommercial, resp.School.ID, resp.School.CodeUnique, resp.User.ID)
	Success(w, http.StatusCreated, resp)
}

// Login connecte un utilisateur.
// POST /api/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	var input models.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps de requête JSON invalide")
		return
	}

	resp, err := h.authService.Login(r.Context(), &input, getClientIP(r), r.UserAgent())
	if err != nil {
		Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	Success(w, http.StatusOK, resp)
}

// Refresh renouvelle un Access Token à partir d'un Refresh Token valide.
// POST /api/auth/refresh
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || strings.TrimSpace(body.RefreshToken) == "" {
		Error(w, http.StatusBadRequest, "Le champ 'refresh_token' est obligatoire")
		return
	}

	resp, err := h.authService.RefreshToken(r.Context(), body.RefreshToken, getClientIP(r), r.UserAgent())
	if err != nil {
		Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	Success(w, http.StatusOK, resp)
}

// Logout révoque la session courante.
// POST /api/auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || strings.TrimSpace(body.RefreshToken) == "" {
		Error(w, http.StatusBadRequest, "Le champ 'refresh_token' est obligatoire")
		return
	}

	if err := h.authService.Logout(r.Context(), body.RefreshToken); err != nil {
		Error(w, http.StatusInternalServerError, "Erreur lors de la déconnexion")
		return
	}

	Success(w, http.StatusOK, map[string]string{
		"message": "Déconnexion réussie",
	})
}

// VerifyPin valide le code PIN manager pour une dérogation guichet.
// POST /api/auth/verify-pin
func (h *AuthHandler) VerifyPin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	var input models.VerifyPinInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps de requête JSON invalide")
		return
	}

	resp, err := h.authService.VerifyPin(r.Context(), input.AutoEcoleID, input.CodePin)
	if err != nil {
		Error(w, http.StatusForbidden, err.Error())
		return
	}

	Success(w, http.StatusOK, resp)
}

// Me renvoie ou met à jour les informations de l'utilisateur actuellement connecté selon son rôle et ses permissions.
// GET / PUT / PATCH /api/auth/me
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		Error(w, http.StatusUnauthorized, "Authentification requise")
		return
	}

	switch r.Method {
	case http.MethodGet:
		user, school, agence, err := h.authService.GetProfile(r.Context(), claims.UserID)
		if err != nil {
			Error(w, http.StatusNotFound, err.Error())
			return
		}

		role := models.NormalizeRole(user.Role)
		permissions := map[string]bool{
			"can_manage_school":   role == models.RoleDirecteur,
			"can_manage_agences":  role == models.RoleDirecteur,
			"can_view_agences":    role == models.RoleDirecteur || role == models.RoleSecretaire || role == models.RoleCaissier,
			"can_configure_pin":   role == models.RoleDirecteur,
			"can_manage_users":    role == models.RoleDirecteur || role == models.RoleSecretaire,
			"can_change_password": role != models.RoleEleve,
		}

		Success(w, http.StatusOK, map[string]interface{}{
			"user":        user,
			"school":      school,
			"agence":      agence,
			"role":        role,
			"permissions": permissions,
		})

	case http.MethodPut, http.MethodPatch:
		var input models.UpdateUserInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			Error(w, http.StatusBadRequest, "Corps de requête JSON invalide")
			return
		}
		updated, err := h.authService.UpdateMyProfile(r.Context(), claims, &input)
		if err != nil {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}
		Success(w, http.StatusOK, updated)

	default:
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
	}
}

// UpdateSchool met à jour les informations officielles de l'auto-école (Directeur uniquement).
// PUT / PATCH /api/auth/school
func (h *AuthHandler) UpdateSchool(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		Error(w, http.StatusUnauthorized, "Authentification requise")
		return
	}

	var input models.AutoEcole
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps de requête JSON invalide")
		return
	}

	updated, err := h.authService.UpdateSchoolInfo(r.Context(), claims, &input)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	Success(w, http.StatusOK, updated)
}

// GetSchoolInfo renvoie les informations publiques d'une auto-école par son slug ou code unique.
// GET /api/auth/school?slug=...
func (h *AuthHandler) GetSchoolInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	slug := r.URL.Query().Get("slug")
	if strings.TrimSpace(slug) == "" {
		Error(w, http.StatusBadRequest, "Paramètre 'slug' manquant dans l'URL")
		return
	}

	info, err := h.authService.GetSchoolInfoBySlug(r.Context(), slug)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	Success(w, http.StatusOK, info)
}

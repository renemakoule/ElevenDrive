package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"autoschool/backend/internal/middleware"
	"autoschool/backend/internal/models"
	"autoschool/backend/internal/realtime"
)

// RegisterUser gère l'ajout d'un collaborateur par le Directeur.
// POST /api/auth/register-user
func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		Error(w, http.StatusUnauthorized, "Authentification requise")
		return
	}

	actor := &models.User{
		ID:          claims.UserID,
		AutoEcoleID: claims.AutoEcoleID,
		AgenceID:    claims.AgenceID,
		Role:        claims.Role,
	}

	var input models.RegisterUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps de requête JSON invalide")
		return
	}

	resp, err := h.authService.RegisterUser(r.Context(), actor, &input)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	realtime.Broadcast(claims.AutoEcoleID, "EQUIPE_USERS_UPDATED", resp)
	Success(w, http.StatusCreated, resp)
}

// GetUsers renvoie tous les collaborateurs de l'auto-école connectée.
// GET /api/auth/users
func (h *AuthHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		Error(w, http.StatusUnauthorized, "Authentification requise")
		return
	}

	list, err := h.authService.GetSchoolUsers(r.Context(), claims.AutoEcoleID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusOK, list)
}

func extractUserIDFromRequest(r *http.Request, ignoredLastSegment string) string {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		parts := strings.Split(strings.TrimRight(r.URL.Path, "/"), "/")
		if len(parts) > 0 {
			last := parts[len(parts)-1]
			if last != "users" && last != ignoredLastSegment {
				userID = last
			}
		}
	}
	return userID
}

// UpdateUser met à jour les informations d'un collaborateur.
// PUT /api/auth/users?id=... ou PUT /api/auth/users/{id}
func (h *AuthHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch && r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée, utilisez PUT ou PATCH")
		return
	}

	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		Error(w, http.StatusUnauthorized, "Authentification requise")
		return
	}

	userID := extractUserIDFromRequest(r, "update")

	var input models.UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps de requête JSON invalide")
		return
	}

	if userID == "" {
		Error(w, http.StatusBadRequest, "L'identifiant du collaborateur est requis")
		return
	}

	resp, err := h.authService.UpdateUser(r.Context(), claims, userID, &input)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	realtime.Broadcast(claims.AutoEcoleID, "EQUIPE_USERS_UPDATED", resp)
	Success(w, http.StatusOK, resp)
}

// DeleteUser supprime un collaborateur (réservé au Directeur).
// DELETE /api/auth/users?id=... ou DELETE /api/auth/users/{id}
func (h *AuthHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée, utilisez DELETE")
		return
	}

	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		Error(w, http.StatusUnauthorized, "Authentification requise")
		return
	}

	userID := extractUserIDFromRequest(r, "delete")
	if userID == "" {
		Error(w, http.StatusBadRequest, "L'identifiant du collaborateur à supprimer est requis")
		return
	}

	if err := h.authService.DeleteUser(r.Context(), claims, userID); err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	realtime.Broadcast(claims.AutoEcoleID, "EQUIPE_USERS_UPDATED", map[string]interface{}{
		"deleted": true,
		"id":      userID,
	})
	Success(w, http.StatusOK, map[string]interface{}{
		"message": "Collaborateur supprimé avec succès",
		"id":      userID,
	})
}

// GetAgences renvoie toutes les agences de l'auto-école connectée.
// GET /api/auth/agences
func (h *AuthHandler) GetAgences(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		Error(w, http.StatusUnauthorized, "Authentification requise")
		return
	}

	list, err := h.authService.GetSchoolAgences(r.Context(), claims.AutoEcoleID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusOK, list)
}

// CreateAgence crée une nouvelle agence.
// POST /api/auth/agences
func (h *AuthHandler) CreateAgence(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		Error(w, http.StatusUnauthorized, "Authentification requise")
		return
	}

	var input models.CreateAgenceInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps de requête JSON invalide")
		return
	}

	agence, err := h.authService.CreateAgence(r.Context(), claims, &input)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	realtime.Broadcast(claims.AutoEcoleID, "EQUIPE_USERS_UPDATED", agence)
	Success(w, http.StatusCreated, agence)
}

func extractAgenceIDFromRequest(r *http.Request) string {
	agenceID := strings.TrimPrefix(r.URL.Path, "/api/auth/agences/")
	agenceID = strings.TrimPrefix(agenceID, "/api/agences/")
	if agenceID == "" || agenceID == r.URL.Path {
		agenceID = r.URL.Query().Get("id")
	}
	return agenceID
}

// UpdateAgence met à jour une agence existante.
// PUT /api/auth/agences/{id} ou PUT /api/auth/agences?id={id}
func (h *AuthHandler) UpdateAgence(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		Error(w, http.StatusUnauthorized, "Authentification requise")
		return
	}

	agenceID := extractAgenceIDFromRequest(r)
	if agenceID == "" {
		Error(w, http.StatusBadRequest, "Identifiant d'agence obligatoire")
		return
	}

	var input models.UpdateAgenceInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps de requête JSON invalide")
		return
	}

	agence, err := h.authService.UpdateAgence(r.Context(), claims, agenceID, &input)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	realtime.Broadcast(claims.AutoEcoleID, "EQUIPE_USERS_UPDATED", agence)
	Success(w, http.StatusOK, agence)
}

// DeleteAgence supprime une succursale non-siège.
// DELETE /api/auth/agences/{id} ou DELETE /api/auth/agences?id={id}
func (h *AuthHandler) DeleteAgence(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		Error(w, http.StatusUnauthorized, "Authentification requise")
		return
	}

	agenceID := extractAgenceIDFromRequest(r)
	if agenceID == "" {
		Error(w, http.StatusBadRequest, "Identifiant d'agence obligatoire")
		return
	}

	if err := h.authService.DeleteAgence(r.Context(), claims, agenceID); err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	realtime.Broadcast(claims.AutoEcoleID, "EQUIPE_USERS_UPDATED", map[string]interface{}{
		"deleted":   true,
		"agence_id": agenceID,
	})
	Success(w, http.StatusOK, map[string]interface{}{
		"message": "Agence supprimée avec succès",
		"id":      agenceID,
	})
}


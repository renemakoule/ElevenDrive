package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"autoschool/backend/internal/middleware"
	"autoschool/backend/internal/models"
	"autoschool/backend/internal/realtime"
	"autoschool/backend/internal/service"
)

type EpreuveHandler struct {
	epreuveService *service.EpreuveService
}

func NewEpreuveHandler(epreuveService *service.EpreuveService) *EpreuveHandler {
	return &EpreuveHandler{epreuveService: epreuveService}
}

// HandleEpreuves gère GET, POST et DELETE sur /api/epreuves
func (h *EpreuveHandler) HandleEpreuves(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		Error(w, http.StatusUnauthorized, "Authentification requise")
		return
	}

	switch r.Method {
	case http.MethodGet:
		categorie := r.URL.Query().Get("categorie")
		list, err := h.epreuveService.GetEpreuves(r.Context(), claims.AutoEcoleID, categorie)
		if err != nil {
			Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		Success(w, http.StatusOK, list)

	case http.MethodPost:
		role := models.NormalizeRole(claims.Role)
		if role != models.RoleDirecteur && role != models.RoleSecretaire && role != models.RoleMoniteur {
			Error(w, http.StatusForbidden, "Privilèges insuffisants : seul le personnel pédagogique peut ajouter une épreuve")
			return
		}

		// Protection DoS mémoire : limiter la taille maximale du corps JSON à 25 Mo
		r.Body = http.MaxBytesReader(w, r.Body, 25*1024*1024)

		var input models.CreerEpreuveInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			Error(w, http.StatusBadRequest, "Corps JSON invalide: "+err.Error())
			return
		}
		ep, err := h.epreuveService.CreerEpreuve(r.Context(), claims.AutoEcoleID, &input)
		if err != nil {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}
		realtime.Broadcast(claims.AutoEcoleID, "EPREUVES_UPDATED", ep)
		Success(w, http.StatusCreated, ep)

	case http.MethodDelete:
		role := models.NormalizeRole(claims.Role)
		if role != models.RoleDirecteur && role != models.RoleSecretaire {
			Error(w, http.StatusForbidden, "Privilèges insuffisants : seul le Directeur ou la Secrétaire peut supprimer une épreuve")
			return
		}

		epreuveID := r.URL.Query().Get("id")
		if strings.TrimSpace(epreuveID) == "" {
			Error(w, http.StatusBadRequest, "Paramètre 'id' manquant")
			return
		}
		if err := h.epreuveService.SupprimerEpreuve(r.Context(), claims.AutoEcoleID, epreuveID); err != nil {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}
		realtime.Broadcast(claims.AutoEcoleID, "EPREUVES_UPDATED", map[string]interface{}{
			"deleted": true,
			"id":      epreuveID,
		})
		Success(w, http.StatusOK, map[string]interface{}{
			"message": "Épreuve supprimée avec succès",
			"id":      epreuveID,
		})

	default:
		Error(w, http.StatusMethodNotAllowed, "Méthode HTTP non autorisée")
	}
}

// HandleTelechargement gère POST /api/epreuves/telecharger?id=...
func (h *EpreuveHandler) HandleTelechargement(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode HTTP non autorisée")
		return
	}
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		Error(w, http.StatusUnauthorized, "Authentification requise")
		return
	}
	epreuveID := r.URL.Query().Get("id")
	if strings.TrimSpace(epreuveID) == "" {
		Error(w, http.StatusBadRequest, "Paramètre 'id' manquant")
		return
	}
	if err := h.epreuveService.IncrementerTelechargement(r.Context(), claims.AutoEcoleID, epreuveID); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	realtime.Broadcast(claims.AutoEcoleID, "EPREUVES_UPDATED", map[string]interface{}{
		"incremented": true,
		"id":          epreuveID,
	})
	Success(w, http.StatusOK, map[string]interface{}{"success": true})
}

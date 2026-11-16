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

type FlotteHandler struct {
	flotteService *service.FlotteService
}

func NewFlotteHandler(flotteService *service.FlotteService) *FlotteHandler {
	return &FlotteHandler{flotteService: flotteService}
}

// HandleVehicules gère GET /api/flotte/vehicules et POST /api/flotte/vehicules
func (h *FlotteHandler) HandleVehicules(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		Error(w, http.StatusUnauthorized, "Authentification requise")
		return
	}

	switch r.Method {
	case http.MethodGet:
		list, err := h.flotteService.GetVehicules(r.Context(), claims.AutoEcoleID)
		if err != nil {
			Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		Success(w, http.StatusOK, list)

	case http.MethodPost:
		role := models.NormalizeRole(claims.Role)
		if role != models.RoleDirecteur && role != models.RoleSecretaire {
			Error(w, http.StatusForbidden, "Privilèges insuffisants : seul le Directeur ou la Secrétaire peut enregistrer un véhicule")
			return
		}

		var input models.CreerVehiculeInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			Error(w, http.StatusBadRequest, "Corps JSON invalide")
			return
		}
		agenceID := ""
		if claims.AgenceID != nil {
			agenceID = *claims.AgenceID
		}
		veh, err := h.flotteService.CreerVehicule(r.Context(), claims.AutoEcoleID, agenceID, &input)
		if err != nil {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}
		realtime.Broadcast(claims.AutoEcoleID, "FLOTTE_VEHICULES_UPDATED", veh)
		Success(w, http.StatusCreated, veh)

	default:
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
	}
}

// UpdateStatut gère PATCH /api/flotte/vehicules/statut
func (h *FlotteHandler) UpdateStatut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		Error(w, http.StatusUnauthorized, "Authentification requise")
		return
	}

	var input struct {
		VehiculeID string `json:"vehicule_id"`
		Statut     string `json:"statut"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps JSON invalide")
		return
	}
	if strings.TrimSpace(input.VehiculeID) == "" {
		Error(w, http.StatusBadRequest, "vehicule_id requis")
		return
	}

	err := h.flotteService.UpdateStatut(r.Context(), claims.AutoEcoleID, input.VehiculeID, input.Statut)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}
	realtime.Broadcast(claims.AutoEcoleID, "FLOTTE_VEHICULES_UPDATED", map[string]interface{}{
		"vehicule_id": input.VehiculeID,
		"statut":      input.Statut,
	})
	Success(w, http.StatusOK, map[string]string{"message": "Statut du véhicule mis à jour avec succès"})
}

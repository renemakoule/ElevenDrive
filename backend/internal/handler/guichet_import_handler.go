package handler

import (
	"encoding/json"
	"net/http"

	"autoschool/backend/internal/models"
	"autoschool/backend/internal/realtime"
)

// POST /api/candidats/import
func (h *GuichetHandler) ImporterCandidatsBatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims, agenceID, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	var input models.ImporterCandidatsBatchInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps de requête JSON invalide")
		return
	}

	result, err := h.guichetService.ImporterCandidatsBatch(
		r.Context(), claims.AutoEcoleID, agenceID, claims.UserID, &input,
	)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if result.TotalCrees > 0 {
		realtime.Broadcast(claims.AutoEcoleID, "CANDIDAT_DOSSIER_UPDATED", map[string]interface{}{
			"action":      "BATCH_IMPORTED",
			"total_crees": result.TotalCrees,
		})
	}

	Success(w, http.StatusCreated, result)
}

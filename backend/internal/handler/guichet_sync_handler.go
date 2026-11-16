package handler

import (
	"net/http"

	"autoschool/backend/internal/realtime"
)

// GET /api/sync/pull
func (h *GuichetHandler) PullSyncSnapshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims, _, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	snapshot, err := h.guichetService.GetSyncPullSnapshot(r.Context(), claims.AutoEcoleID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusOK, snapshot)
}

func (h *GuichetHandler) broadcastSyncApplied(autoEcoleID string, count int) {
	if count > 0 {
		realtime.Broadcast(autoEcoleID, "CANDIDAT_DOSSIER_UPDATED", map[string]interface{}{
			"action": "SYNC_BATCH_APPLIED",
			"count":  count,
		})
	}
}


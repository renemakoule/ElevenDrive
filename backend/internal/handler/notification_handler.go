package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"autoschool/backend/internal/middleware"
	"autoschool/backend/internal/models"
	"autoschool/backend/internal/service"
)

type NotificationHandler struct {
	notifService *service.NotificationService
}

func NewNotificationHandler(notifService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notifService: notifService}
}

// GetNotifications renvoie les notifications du compte connecté avec le décompte des non lues
// GET /api/notifications?unread=true&limit=50
func (h *NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		Error(w, http.StatusUnauthorized, "Authentification requise")
		return
	}

	unreadOnly := r.URL.Query().Get("unread") == "true"
	limit := 50
	if lStr := r.URL.Query().Get("limit"); lStr != "" {
		if val, err := strconv.Atoi(lStr); err == nil && val > 0 {
			limit = val
		}
	}

	list, unreadCount, err := h.notifService.GetUserNotifications(
		r.Context(),
		claims.AutoEcoleID,
		claims.UserID,
		claims.Role,
		limit,
		unreadOnly,
	)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusOK, map[string]interface{}{
		"notifications": list,
		"unread_count":  unreadCount,
	})
}

// MarkAsRead passe une notification spécifique à l'état lu
// POST /api/notifications/read
func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		Error(w, http.StatusUnauthorized, "Authentification requise")
		return
	}

	var body struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.ID == "" {
		body.ID = r.URL.Query().Get("id")
	}
	if body.ID == "" {
		Error(w, http.StatusBadRequest, "Identifiant de notification requis")
		return
	}

	if err := h.notifService.MarkAsRead(r.Context(), body.ID); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusOK, map[string]interface{}{
		"id":     body.ID,
		"est_lu": true,
	})
}

// MarkAllAsRead passe toutes les notifications du compte courant à l'état lu
// POST /api/notifications/read-all
func (h *NotificationHandler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		Error(w, http.StatusUnauthorized, "Authentification requise")
		return
	}

	if err := h.notifService.MarkAllAsRead(r.Context(), claims.AutoEcoleID, claims.UserID, claims.Role); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

// GetAuditLogs renvoie les logs d'audit filtrés
// GET /api/notifications/audit
func (h *NotificationHandler) GetAuditLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		Error(w, http.StatusUnauthorized, "Authentification requise")
		return
	}

	roleNorm := models.NormalizeRole(claims.Role)
	schoolFilter := claims.AutoEcoleID
	if roleNorm == models.RoleSuperAdmin {
		schoolFilter = r.URL.Query().Get("auto_ecole_id")
	}

	limit := 100
	if lStr := r.URL.Query().Get("limit"); lStr != "" {
		if val, err := strconv.Atoi(lStr); err == nil && val > 0 {
			limit = val
		}
	}

	logs, err := h.notifService.GetAuditLogs(r.Context(), schoolFilter, limit)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusOK, map[string]interface{}{
		"logs": logs,
	})
}

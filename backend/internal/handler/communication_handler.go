package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"autoschool/backend/internal/middleware"
	"autoschool/backend/internal/models"
	"autoschool/backend/internal/service"
)

type CommunicationHandler struct {
	commService *service.CommunicationService
}

func NewCommunicationHandler(commService *service.CommunicationService) *CommunicationHandler {
	return &CommunicationHandler{commService: commService}
}

// -------------------------------------------------------------
// SUPER-ADMIN : ALERTES
// -------------------------------------------------------------

// HandleSuperAdminAlerts dispatche GET, POST, DELETE sur /api/super-admin/alerts
func (h *CommunicationHandler) HandleSuperAdminAlerts(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil || models.NormalizeRole(claims.Role) != models.RoleSuperAdmin {
		Error(w, http.StatusForbidden, "Accès strictement réservé au Super-Administrateur")
		return
	}

	switch r.Method {
	case http.MethodGet:
		list, err := h.commService.GetAllAlerts(r.Context())
		if err != nil {
			Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		Success(w, http.StatusOK, list)

	case http.MethodPost:
		var input models.CreateAlertInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			Error(w, http.StatusBadRequest, "Corps de requête invalide: "+err.Error())
			return
		}
		alert, err := h.commService.CreateAlert(r.Context(), &input)
		if err != nil {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}
		Success(w, http.StatusCreated, alert)

	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" {
			var body struct {
				ID string `json:"id"`
			}
			_ = json.NewDecoder(r.Body).Decode(&body)
			id = body.ID
		}
		if id == "" {
			Error(w, http.StatusBadRequest, "Identifiant d'alerte requis")
			return
		}
		if err := h.commService.DeleteAlert(r.Context(), id); err != nil {
			Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		Success(w, http.StatusOK, map[string]interface{}{"deleted": true, "id": id})

	default:
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
	}
}

// ToggleAlert bascule l'état actif/inactif d'une alerte
// POST ou PATCH /api/super-admin/alerts/toggle
func (h *CommunicationHandler) ToggleAlert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPatch {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}
	claims := middleware.GetClaims(r.Context())
	if claims == nil || models.NormalizeRole(claims.Role) != models.RoleSuperAdmin {
		Error(w, http.StatusForbidden, "Accès réservé au Super-Administrateur")
		return
	}

	var body struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.ID == "" {
		body.ID = r.URL.Query().Get("id")
	}
	if body.ID == "" {
		Error(w, http.StatusBadRequest, "Identifiant d'alerte manquant")
		return
	}

	active, err := h.commService.ToggleAlertStatus(r.Context(), body.ID)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}
	Success(w, http.StatusOK, map[string]interface{}{"id": body.ID, "est_active": active})
}

// -------------------------------------------------------------
// SUPER-ADMIN : ANNONCES
// -------------------------------------------------------------

// HandleSuperAdminAnnouncements dispatche GET, POST, DELETE sur /api/super-admin/announcements
func (h *CommunicationHandler) HandleSuperAdminAnnouncements(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil || models.NormalizeRole(claims.Role) != models.RoleSuperAdmin {
		Error(w, http.StatusForbidden, "Accès strictement réservé au Super-Administrateur")
		return
	}

	switch r.Method {
	case http.MethodGet:
		list, err := h.commService.GetAllAnnouncements(r.Context())
		if err != nil {
			Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		Success(w, http.StatusOK, list)

	case http.MethodPost:
		var input models.CreateAnnouncementInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			Error(w, http.StatusBadRequest, "Corps de requête invalide: "+err.Error())
			return
		}
		ann, err := h.commService.CreateAnnouncement(r.Context(), &input)
		if err != nil {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}
		Success(w, http.StatusCreated, ann)

	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" {
			var body struct {
				ID string `json:"id"`
			}
			_ = json.NewDecoder(r.Body).Decode(&body)
			id = body.ID
		}
		if id == "" {
			Error(w, http.StatusBadRequest, "Identifiant d'annonce requis")
			return
		}
		if err := h.commService.DeleteAnnouncement(r.Context(), id); err != nil {
			Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		Success(w, http.StatusOK, map[string]interface{}{"deleted": true, "id": id})

	default:
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
	}
}

// -------------------------------------------------------------
// SUPER-ADMIN : AVIS & RETOURS
// -------------------------------------------------------------

// HandleSuperAdminFeedbacks dispatche GET, POST, DELETE sur /api/super-admin/feedbacks
func (h *CommunicationHandler) HandleSuperAdminFeedbacks(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil || models.NormalizeRole(claims.Role) != models.RoleSuperAdmin {
		Error(w, http.StatusForbidden, "Accès strictement réservé au Super-Administrateur")
		return
	}

	switch r.Method {
	case http.MethodGet:
		stars := 0
		if sStr := r.URL.Query().Get("etoiles"); sStr != "" {
			stars, _ = strconv.Atoi(sStr)
		}
		list, err := h.commService.GetAllFeedbacks(r.Context(), stars)
		if err != nil {
			Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		Success(w, http.StatusOK, list)

	case http.MethodPost:
		var input models.CreateFeedbackInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			Error(w, http.StatusBadRequest, "Corps de requête invalide: "+err.Error())
			return
		}
		fb, err := h.commService.CreateFeedback(r.Context(), &input)
		if err != nil {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}
		Success(w, http.StatusCreated, fb)

	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" {
			var body struct {
				ID string `json:"id"`
			}
			_ = json.NewDecoder(r.Body).Decode(&body)
			id = body.ID
		}
		if id == "" {
			Error(w, http.StatusBadRequest, "Identifiant de l'avis requis")
			return
		}
		if err := h.commService.DeleteFeedback(r.Context(), id); err != nil {
			Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		Success(w, http.StatusOK, map[string]interface{}{"deleted": true, "id": id})

	default:
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
	}
}

// ReplyFeedback répond officiellement à un avis
// POST /api/super-admin/feedbacks/reply
func (h *CommunicationHandler) ReplyFeedback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}
	claims := middleware.GetClaims(r.Context())
	if claims == nil || models.NormalizeRole(claims.Role) != models.RoleSuperAdmin {
		Error(w, http.StatusForbidden, "Accès réservé au Super-Administrateur")
		return
	}

	var input models.ReplyFeedbackInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps de requête invalide: "+err.Error())
		return
	}

	if err := h.commService.ReplyFeedback(r.Context(), &input); err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}
	Success(w, http.StatusOK, map[string]interface{}{"success": true, "id": input.FeedbackID})
}

// UpdateFeedbackStatus met à jour le statut d'un avis
// PATCH /api/super-admin/feedbacks/status
func (h *CommunicationHandler) UpdateFeedbackStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch && r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}
	claims := middleware.GetClaims(r.Context())
	if claims == nil || models.NormalizeRole(claims.Role) != models.RoleSuperAdmin {
		Error(w, http.StatusForbidden, "Accès réservé au Super-Administrateur")
		return
	}

	var input models.UpdateFeedbackStatusInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps de requête invalide: "+err.Error())
		return
	}

	if err := h.commService.UpdateFeedbackStatus(r.Context(), input.FeedbackID, input.Statut); err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}
	Success(w, http.StatusOK, map[string]interface{}{"success": true, "id": input.FeedbackID, "statut": input.Statut})
}

// GetSupportTickets liste les tickets de support (Super-Admin)
// GET /api/super-admin/support-tickets
func (h *CommunicationHandler) GetSupportTickets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}
	claims := middleware.GetClaims(r.Context())
	if claims == nil || models.NormalizeRole(claims.Role) != models.RoleSuperAdmin {
		Error(w, http.StatusForbidden, "Accès réservé au Super-Administrateur")
		return
	}

	tickets, err := h.commService.GetAllSupportTickets(r.Context())
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, http.StatusOK, tickets)
}

// -------------------------------------------------------------
// ENDPOINTS CLIENTS / AUTO-ÉCOLES / ASSISTANCE
// -------------------------------------------------------------

// SubmitSupportTicket soumet un ticket de support depuis le bouton d'aide
// POST /api/support/ticket
func (h *CommunicationHandler) SubmitSupportTicket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	var input models.CreateSupportTicketInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps de requête invalide: "+err.Error())
		return
	}

	// Si un token JWT est présent, attacher automatiquement l'auto-école
	claims := middleware.GetClaims(r.Context())
	if claims != nil && claims.AutoEcoleID != "" && input.AutoEcoleID == nil {
		schoolID := claims.AutoEcoleID
		input.AutoEcoleID = &schoolID
	}

	ticket, err := h.commService.CreateSupportTicket(r.Context(), &input)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	Success(w, http.StatusCreated, ticket)
}

// GetClientAlerts renvoie les alertes actives pour l'auto-école connectée
// GET /api/communication/alerts
func (h *CommunicationHandler) GetClientAlerts(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	schoolID := ""
	if claims != nil {
		schoolID = claims.AutoEcoleID
	}
	list, err := h.commService.GetActiveAlertsForSchool(r.Context(), schoolID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, http.StatusOK, list)
}

// GetClientAnnouncements renvoie les annonces destinées à l'auto-école connectée
// GET /api/communication/announcements
func (h *CommunicationHandler) GetClientAnnouncements(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	schoolID := ""
	if claims != nil {
		schoolID = claims.AutoEcoleID
	}
	list, err := h.commService.GetAnnouncementsForSchool(r.Context(), schoolID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, http.StatusOK, list)
}

// SubmitClientFeedback enregistre un avis d'auto-école
// POST /api/communication/feedbacks
func (h *CommunicationHandler) SubmitClientFeedback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	var input models.CreateFeedbackInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps de requête invalide: "+err.Error())
		return
	}

	claims := middleware.GetClaims(r.Context())
	if claims != nil && claims.AutoEcoleID != "" && input.AutoEcoleID == nil {
		schoolID := claims.AutoEcoleID
		input.AutoEcoleID = &schoolID
	}

	fb, err := h.commService.CreateFeedback(r.Context(), &input)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	Success(w, http.StatusCreated, fb)
}

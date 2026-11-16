package handler

import (
	"encoding/json"
	"net/http"

	"autoschool/backend/internal/middleware"
	"autoschool/backend/internal/models"
	"autoschool/backend/internal/service"
)

type SuperAdminHandler struct {
	superAdminService *service.SuperAdminService
}

func NewSuperAdminHandler(superAdminService *service.SuperAdminService) *SuperAdminHandler {
	return &SuperAdminHandler{superAdminService: superAdminService}
}

// LoginSuperAdmin gère l'authentification console du super-administrateur
// POST /api/super-admin/login
func (h *SuperAdminHandler) LoginSuperAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	var input models.SuperAdminLoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps de requête JSON invalide")
		return
	}

	ip := getClientIP(r)
	ua := r.UserAgent()

	resp, err := h.superAdminService.LoginSuperAdmin(r.Context(), &input, ip, ua)
	if err != nil {
		Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	Success(w, http.StatusOK, resp)
}

// GetDashboard récupère les métriques globales et la liste des auto-écoles
// GET /api/super-admin/dashboard
func (h *SuperAdminHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims := middleware.GetClaims(r.Context())
	if claims == nil || models.NormalizeRole(claims.Role) != models.RoleSuperAdmin {
		Error(w, http.StatusForbidden, "Accès strictement réservé au Super-Administrateur")
		return
	}

	search := r.URL.Query().Get("q")
	stats, schools, err := h.superAdminService.GetDashboardData(r.Context(), search)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusOK, map[string]interface{}{
		"stats":   stats,
		"schools": schools,
	})
}

// GetUsers liste tous les utilisateurs enregistrés pour la vue Utilisateurs
// GET /api/super-admin/users
func (h *SuperAdminHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims := middleware.GetClaims(r.Context())
	if claims == nil || models.NormalizeRole(claims.Role) != models.RoleSuperAdmin {
		Error(w, http.StatusForbidden, "Accès strictement réservé au Super-Administrateur")
		return
	}

	search := r.URL.Query().Get("q")
	users, err := h.superAdminService.GetUsersData(r.Context(), search)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusOK, map[string]interface{}{
		"users": users,
	})
}

// GetCandidats liste tous les candidats enregistrés pour la vue Candidats
// GET /api/super-admin/candidats
func (h *SuperAdminHandler) GetCandidats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims := middleware.GetClaims(r.Context())
	if claims == nil || models.NormalizeRole(claims.Role) != models.RoleSuperAdmin {
		Error(w, http.StatusForbidden, "Accès strictement réservé au Super-Administrateur")
		return
	}

	search := r.URL.Query().Get("q")
	candidats, err := h.superAdminService.GetCandidatsData(r.Context(), search)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusOK, map[string]interface{}{
		"candidats": candidats,
	})
}

// ToggleSchoolStatus active ou suspend un établissement
// POST /api/super-admin/schools/status
func (h *SuperAdminHandler) ToggleSchoolStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims := middleware.GetClaims(r.Context())
	if claims == nil || models.NormalizeRole(claims.Role) != models.RoleSuperAdmin {
		Error(w, http.StatusForbidden, "Accès strictement réservé au Super-Administrateur")
		return
	}

	var input models.ToggleSchoolStatusInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps de requête JSON invalide")
		return
	}

	if err := h.superAdminService.ToggleSchoolStatus(r.Context(), input.SchoolID, input.EstActif, getClientIP(r)); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusOK, map[string]interface{}{
		"success":   true,
		"school_id": input.SchoolID,
		"est_actif": input.EstActif,
	})
}

// UpdateSchool met à jour les informations d'une auto-école
// PUT /api/super-admin/schools
func (h *SuperAdminHandler) UpdateSchool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims := middleware.GetClaims(r.Context())
	if claims == nil || models.NormalizeRole(claims.Role) != models.RoleSuperAdmin {
		Error(w, http.StatusForbidden, "Accès strictement réservé au Super-Administrateur")
		return
	}

	var input models.UpdateSchoolAdminInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps de requête JSON invalide")
		return
	}

	if err := h.superAdminService.UpdateSchoolAdmin(r.Context(), &input, getClientIP(r)); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"school":  input,
	})
}

// UpdateLicence met à jour le nombre de machines et la durée de validité
// POST /api/super-admin/schools/licence
func (h *SuperAdminHandler) UpdateLicence(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims := middleware.GetClaims(r.Context())
	if claims == nil || models.NormalizeRole(claims.Role) != models.RoleSuperAdmin {
		Error(w, http.StatusForbidden, "Accès strictement réservé au Super-Administrateur")
		return
	}

	var input models.UpdateSchoolLicenceInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps de requête JSON invalide")
		return
	}

	if err := h.superAdminService.UpdateSchoolLicence(r.Context(), &input, getClientIP(r)); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"licence": input,
	})
}

// UpdateModules configure les modules activés pour une auto-école
// POST /api/super-admin/schools/modules
func (h *SuperAdminHandler) UpdateModules(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims := middleware.GetClaims(r.Context())
	if claims == nil || models.NormalizeRole(claims.Role) != models.RoleSuperAdmin {
		Error(w, http.StatusForbidden, "Accès strictement réservé au Super-Administrateur")
		return
	}

	var input models.UpdateSchoolModulesInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps de requête JSON invalide")
		return
	}

	if err := h.superAdminService.UpdateSchoolModules(r.Context(), &input, getClientIP(r)); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"modules": input.ModulesActifs,
	})
}

// DeleteSchool supprime définitivement une auto-école
// POST /api/super-admin/schools/delete
func (h *SuperAdminHandler) DeleteSchool(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims := middleware.GetClaims(r.Context())
	if claims == nil || models.NormalizeRole(claims.Role) != models.RoleSuperAdmin {
		Error(w, http.StatusForbidden, "Accès strictement réservé au Super-Administrateur")
		return
	}

	var input models.DeleteSchoolInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps de requête JSON invalide")
		return
	}

	if err := h.superAdminService.DeleteSchool(r.Context(), input.SchoolID, input.ConfirmName, getClientIP(r)); err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	Success(w, http.StatusOK, map[string]interface{}{
		"success":   true,
		"school_id": input.SchoolID,
	})
}

// GetAuditLogs renvoie les logs d'administration système
// GET /api/super-admin/audit
func (h *SuperAdminHandler) GetAuditLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims := middleware.GetClaims(r.Context())
	if claims == nil || models.NormalizeRole(claims.Role) != models.RoleSuperAdmin {
		Error(w, http.StatusForbidden, "Accès strictement réservé au Super-Administrateur")
		return
	}

	logs, err := h.superAdminService.GetAuditLogs(r.Context(), 100)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusOK, map[string]interface{}{
		"logs": logs,
	})
}

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

type GuichetHandler struct {
	guichetService *service.GuichetService
}

func NewGuichetHandler(guichetService *service.GuichetService) *GuichetHandler {
	return &GuichetHandler{guichetService: guichetService}
}

func requireClaimsWithAgence(w http.ResponseWriter, r *http.Request) (*models.JWTClaims, string, bool) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		Error(w, http.StatusUnauthorized, "Authentification requise")
		return nil, "", false
	}
	agenceID := ""
	if claims.AgenceID != nil {
		agenceID = *claims.AgenceID
	}
	return claims, agenceID, true
}

// GET /api/guichet/stats
func (h *GuichetHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	claims, agenceID, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	stats, err := h.guichetService.GetGuichetStats(r.Context(), claims.AutoEcoleID, agenceID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, http.StatusOK, stats)
}

// GET /api/formules
func (h *GuichetHandler) GetFormules(w http.ResponseWriter, r *http.Request) {
	claims, _, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	formules, err := h.guichetService.GetFormules(r.Context(), claims.AutoEcoleID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, http.StatusOK, formules)
}

// POST /api/formules
func (h *GuichetHandler) CreerFormule(w http.ResponseWriter, r *http.Request) {
	claims, _, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	roleUpper := models.NormalizeRole(claims.Role)
	if roleUpper != models.RoleDirecteur && roleUpper != models.RoleSecretaire && roleUpper != models.RoleCaissier {
		Error(w, http.StatusForbidden, "Accès refusé : Seul le Directeur ou la Secrétaire peut créer une formule et fixer son prix initial")
		return
	}

	var input models.CreerFormuleInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Format de données JSON invalide")
		return
	}

	formule, err := h.guichetService.CreerFormule(r.Context(), claims.AutoEcoleID, claims.Role, input)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	realtime.Broadcast(claims.AutoEcoleID, "FORMULES_TARIFS_UPDATED", formule)
	Success(w, http.StatusCreated, formule)
}

// PATCH / PUT /api/formules/prix (Directeur uniquement)
func (h *GuichetHandler) ModifierPrixFormule(w http.ResponseWriter, r *http.Request) {
	claims, _, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	if models.NormalizeRole(claims.Role) != models.RoleDirecteur {
		Error(w, http.StatusForbidden, "Accès formellement refusé : Seul le Directeur est habilité à modifier le prix d'un permis déjà fixé")
		return
	}

	var input models.ModifierPrixFormuleInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Format de données JSON invalide")
		return
	}

	formule, err := h.guichetService.ModifierPrixFormule(r.Context(), claims.AutoEcoleID, claims.Role, input)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	realtime.Broadcast(claims.AutoEcoleID, "FORMULES_TARIFS_UPDATED", formule)
	Success(w, http.StatusOK, formule)
}

// GET /api/candidats?q=...
func (h *GuichetHandler) SearchCandidats(w http.ResponseWriter, r *http.Request) {
	claims, _, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	candidats, err := h.guichetService.SearchCandidats(r.Context(), claims.AutoEcoleID, r.URL.Query().Get("q"))
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Success(w, http.StatusOK, candidats)
}

// GET /api/candidats/:id/dossier
func (h *GuichetHandler) GetCandidatDossier(w http.ResponseWriter, r *http.Request) {
	claims, _, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	candidatID := r.URL.Query().Get("candidat_id")
	if candidatID == "" {
		Error(w, http.StatusBadRequest, "Paramètre 'candidat_id' obligatoire")
		return
	}

	dossier, err := h.guichetService.GetCandidatDossier(r.Context(), claims.AutoEcoleID, candidatID)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}
	Success(w, http.StatusOK, dossier)
}

// POST /api/candidats
func (h *GuichetHandler) CreerCandidat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims, agenceID, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	var input models.CreerCandidatInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps de requête JSON invalide")
		return
	}

	candidat, dossierComplet, err := h.guichetService.CreerCandidatAvecDossier(
		r.Context(), claims.AutoEcoleID, agenceID, claims.UserID, claims.Role, &input,
	)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	realtime.Broadcast(claims.AutoEcoleID, "CANDIDAT_DOSSIER_UPDATED", map[string]interface{}{
		"action":         "CREATED",
		"candidat":       candidat,
		"dossier_detail": dossierComplet,
	})
	Success(w, http.StatusCreated, map[string]interface{}{
		"candidat":       candidat,
		"dossier_detail": dossierComplet,
	})
}

// PATCH /api/candidats/{id}/photo
func (h *GuichetHandler) UpdateCandidatPhoto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch && r.Method != http.MethodPut && r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims, _, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	candidatID := r.URL.Query().Get("candidat_id")
	if candidatID == "" {
		trimmed := strings.TrimPrefix(r.URL.Path, "/api/candidats/")
		parts := strings.Split(strings.Trim(trimmed, "/"), "/")
		if len(parts) >= 1 && parts[0] != "" && parts[0] != "dossier" {
			candidatID = parts[0]
		}
	}
	if candidatID == "" {
		Error(w, http.StatusBadRequest, "Identifiant du candidat manquant")
		return
	}

	var input models.UpdatePhotoInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps JSON invalide")
		return
	}

	if err := h.guichetService.UpdateCandidatPhoto(r.Context(), claims.AutoEcoleID, candidatID, input.PhotoURL); err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	realtime.Broadcast(claims.AutoEcoleID, "CANDIDAT_DOSSIER_UPDATED", map[string]interface{}{
		"action":      "PHOTO",
		"candidat_id": candidatID,
	})
	Success(w, http.StatusOK, map[string]interface{}{
		"success":     true,
		"candidat_id": candidatID,
	})
}

// POST /api/paiements
func (h *GuichetHandler) EnregistrerPaiement(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims, agenceID, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	var input models.EnregistrerPaiementInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps de requête JSON invalide")
		return
	}

	resp, err := h.guichetService.EnregistrerPaiement(r.Context(), claims.AutoEcoleID, agenceID, claims.UserID, &input)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	realtime.Broadcast(claims.AutoEcoleID, "PAIEMENT_ENREGISTRE", resp)
	realtime.Broadcast(claims.AutoEcoleID, "CANDIDAT_DOSSIER_UPDATED", map[string]interface{}{
		"action":     "PAIEMENT",
		"dossier_id": input.DossierID,
	})
	Success(w, http.StatusCreated, resp)
}

// GET /api/paiements/recu?paiement_id=...
func (h *GuichetHandler) ReimprimerRecu(w http.ResponseWriter, r *http.Request) {
	claims, _, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	paiementID := r.URL.Query().Get("paiement_id")
	if paiementID == "" {
		Error(w, http.StatusBadRequest, "Paramètre 'paiement_id' obligatoire")
		return
	}

	recu, err := h.guichetService.ReimprimerRecu(r.Context(), claims.AutoEcoleID, paiementID)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	Success(w, http.StatusOK, recu)
}

// GET /api/caisse/session-detail
func (h *GuichetHandler) GetSessionCaisse(w http.ResponseWriter, r *http.Request) {
	claims, agenceID, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	detail, err := h.guichetService.GetSessionCaisseDetail(r.Context(), claims.AutoEcoleID, agenceID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusOK, detail)
}

// POST /api/caisse/cloturer
func (h *GuichetHandler) CloturerSessionCaisse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims, agenceID, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	var input models.CloturerSessionCaisseInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps de requête JSON invalide")
		return
	}

	if err := h.guichetService.CloturerSessionCaisse(r.Context(), claims.AutoEcoleID, agenceID, claims.UserID, &input); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	realtime.Broadcast(claims.AutoEcoleID, "CAISSE_SESSION_UPDATED", map[string]interface{}{
		"action":    "CLOTURE",
		"agence_id": agenceID,
	})
	Success(w, http.StatusOK, map[string]string{
		"message": "Session de caisse clôturée avec succès",
	})
}

// GET /api/manager/dashboard
func (h *GuichetHandler) GetTableauGerant(w http.ResponseWriter, r *http.Request) {
	claims, _, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	data, err := h.guichetService.GetTableauGerant(r.Context(), claims.AutoEcoleID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusOK, data)
}

// POST /api/manager/derogation
func (h *GuichetHandler) AccorderDerogation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims, _, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	var input models.AccorderDerogationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps JSON invalide")
		return
	}

	err := h.guichetService.AccorderDerogationConduite(r.Context(), claims.AutoEcoleID, input.DossierID, input.CodePIN, input.Motif)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	realtime.Broadcast(claims.AutoEcoleID, "CANDIDAT_DOSSIER_UPDATED", map[string]interface{}{
		"action":     "DEROGATION",
		"dossier_id": input.DossierID,
	})
	Success(w, http.StatusOK, map[string]interface{}{
		"success":    true,
		"message":    "Dérogation accordée avec succès. Accès à la conduite déverrouillé.",
		"dossier_id": input.DossierID,
	})
}

// POST /api/sync/batch
func (h *GuichetHandler) SynchroniserLotEvenements(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims, agenceID, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	var input models.SyncBatchInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps JSON invalide")
		return
	}

	resp, err := h.guichetService.SynchroniserLotEvenements(r.Context(), claims.AutoEcoleID, agenceID, input.Events)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusOK, resp)
}

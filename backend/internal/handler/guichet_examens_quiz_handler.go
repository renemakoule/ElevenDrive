package handler

import (
	"encoding/json"
	"net/http"

	"autoschool/backend/internal/models"
	"autoschool/backend/internal/realtime"
)

// GET /api/examens/session-active
func (h *GuichetHandler) GetSessionExamenActive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims, _, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	data, err := h.guichetService.GetSessionExamenActive(r.Context(), claims.AutoEcoleID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusOK, data)
}

// POST /api/examens/sessions
func (h *GuichetHandler) CreerOuModifierSessionExamen(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims, _, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	var input models.CreerOuModifierSessionExamenInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps JSON invalide: "+err.Error())
		return
	}

	sess, err := h.guichetService.CreerOuModifierSessionExamen(r.Context(), claims.AutoEcoleID, input)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	realtime.Broadcast(claims.AutoEcoleID, "EXAMENS_SESSION_UPDATED", sess)
	Success(w, http.StatusOK, sess)
}

// POST / DELETE /api/examens/sessions/supprimer
func (h *GuichetHandler) SupprimerSessionExamen(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims, _, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	var input models.SupprimerSessionExamenInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps JSON invalide: "+err.Error())
		return
	}

	if err := h.guichetService.SupprimerSessionExamen(r.Context(), claims.AutoEcoleID, input.SessionID); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	realtime.Broadcast(claims.AutoEcoleID, "EXAMENS_SESSION_UPDATED", map[string]interface{}{
		"deleted":    true,
		"session_id": input.SessionID,
	})
	Success(w, http.StatusOK, map[string]string{"message": "Session d'examen supprimée avec succès"})
}

// POST /api/examens/resultats
func (h *GuichetHandler) EnregistrerResultatsExamen(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims, _, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	var input models.EnregistrerResultatsExamenInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps JSON invalide")
		return
	}

	if err := h.guichetService.EnregistrerResultatsExamen(r.Context(), claims.AutoEcoleID, &input); err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	realtime.Broadcast(claims.AutoEcoleID, "EXAMENS_SESSION_UPDATED", map[string]interface{}{
		"action":     "RESULTATS_PUBLIES",
		"session_id": input.SessionExamenID,
	})
	Success(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Résultats d'examen et procès-verbal enregistrés avec succès",
	})
}

// POST /api/examens/forcer-candidat
func (h *GuichetHandler) ForcerCandidatExamen(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims, _, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	var input models.ForcerCandidatExamenInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps JSON invalide: "+err.Error())
		return
	}

	cand, err := h.guichetService.ForcerCandidatExamen(r.Context(), claims.AutoEcoleID, input)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	realtime.Broadcast(claims.AutoEcoleID, "EXAMENS_SESSION_UPDATED", map[string]interface{}{
		"action":     "CANDIDAT_FORCE",
		"session_id": input.SessionExamenID,
		"dossier_id": input.DossierID,
	})
	Success(w, http.StatusOK, cand)
}

// POST /api/pedagogie/quiz-resultat
func (h *GuichetHandler) EnregistrerQuizResultat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims, _, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	var input models.EnregistrerQuizResultatInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps JSON invalide")
		return
	}

	if err := h.guichetService.EnregistrerQuizResultat(r.Context(), claims.AutoEcoleID, input); err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	realtime.Broadcast(claims.AutoEcoleID, "QUIZ_RESULTAT_UPDATED", map[string]interface{}{
		"candidat_id": input.CandidatID,
		"score":       input.Score,
	})
	realtime.Broadcast(claims.AutoEcoleID, "CANDIDAT_DOSSIER_UPDATED", map[string]interface{}{
		"action":      "QUIZ",
		"candidat_id": input.CandidatID,
		"score":       input.Score,
	})
	Success(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Score du quiz enregistré",
	})
}

// GET, POST, DELETE /api/pedagogie/quiz-forms
func (h *GuichetHandler) HandleQuizFormulaires(w http.ResponseWriter, r *http.Request) {
	claims, _, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	switch r.Method {
	case http.MethodGet:
		list, err := h.guichetService.GetQuizFormulaires(r.Context(), claims.AutoEcoleID, r.URL.Query().Get("categorie"))
		if err != nil {
			Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		Success(w, http.StatusOK, list)
	case http.MethodPost:
		role := models.NormalizeRole(claims.Role)
		if role != models.RoleDirecteur && role != models.RoleSecretaire && role != models.RoleMoniteur {
			Error(w, http.StatusForbidden, "Privilèges insuffisants : seul le personnel pédagogique peut créer un formulaire de quiz")
			return
		}
		var input models.CreerQuizFormulaireInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			Error(w, http.StatusBadRequest, "Corps JSON invalide")
			return
		}
		created, err := h.guichetService.CreerQuizFormulaire(r.Context(), claims.AutoEcoleID, input)
		if err != nil {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}
		realtime.Broadcast(claims.AutoEcoleID, "QUIZ_FORM_UPDATED", created)
		Success(w, http.StatusCreated, created)
	case http.MethodDelete:
		role := models.NormalizeRole(claims.Role)
		if role != models.RoleDirecteur && role != models.RoleSecretaire {
			Error(w, http.StatusForbidden, "Privilèges insuffisants : seul le Directeur ou la Secrétaire peut supprimer un quiz")
			return
		}
		id := r.URL.Query().Get("id")
		if id == "" {
			Error(w, http.StatusBadRequest, "Paramètre 'id' requis")
			return
		}
		if err := h.guichetService.SupprimerQuizFormulaire(r.Context(), claims.AutoEcoleID, id); err != nil {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}
		realtime.Broadcast(claims.AutoEcoleID, "QUIZ_FORM_UPDATED", map[string]interface{}{
			"deleted": true,
			"id":      id,
		})
		Success(w, http.StatusOK, map[string]bool{"deleted": true})
	default:
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
	}
}

// GET, POST /api/pedagogie/quiz-submissions
func (h *GuichetHandler) HandleQuizSoumissions(w http.ResponseWriter, r *http.Request) {
	claims, _, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	switch r.Method {
	case http.MethodGet:
		list, err := h.guichetService.GetQuizSoumissions(
			r.Context(), claims.AutoEcoleID, r.URL.Query().Get("quiz_id"), r.URL.Query().Get("candidat_id"),
		)
		if err != nil {
			Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		Success(w, http.StatusOK, list)
	case http.MethodPost:
		var input models.SoumettreQuizReponseInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			Error(w, http.StatusBadRequest, "Corps JSON invalide")
			return
		}
		sub, err := h.guichetService.SoumettreQuizReponse(r.Context(), claims.AutoEcoleID, input)
		if err != nil {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}
		realtime.Broadcast(claims.AutoEcoleID, "QUIZ_SUBMISSION_RECEIVED", sub)
		realtime.Broadcast(claims.AutoEcoleID, "QUIZ_RESULTAT_UPDATED", map[string]interface{}{
			"candidat_id": sub.CandidatID,
			"score":       sub.Pourcentage,
		})
		Success(w, http.StatusCreated, sub)
	default:
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
	}
}


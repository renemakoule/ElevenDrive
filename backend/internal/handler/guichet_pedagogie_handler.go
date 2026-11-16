package handler

import (
	"encoding/json"
	"net/http"

	"autoschool/backend/internal/models"
	"autoschool/backend/internal/realtime"
)

// POST /api/pedagogie/seance-code
func (h *GuichetHandler) EnregistrerSeanceCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims, agenceID, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	var input models.EnregistrerSeanceCodeInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps de requête JSON invalide")
		return
	}

	resp, err := h.guichetService.EnregistrerSeanceCode(r.Context(), claims.AutoEcoleID, agenceID, &input)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	realtime.Broadcast(claims.AutoEcoleID, "COURS_CODE_UPDATED", resp)
	Success(w, http.StatusCreated, resp)
}

// GET / DELETE /api/pedagogie/planning-conduite
func (h *GuichetHandler) GetPlanningConduite(w http.ResponseWriter, r *http.Request) {
	claims, _, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	if r.Method == http.MethodDelete {
		seanceID := r.URL.Query().Get("id")
		if seanceID == "" {
			Error(w, http.StatusBadRequest, "Paramètre 'id' de la séance obligatoire")
			return
		}
		if err := h.guichetService.AnnulerSeanceConduite(r.Context(), claims.AutoEcoleID, seanceID); err != nil {
			Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		realtime.Broadcast(claims.AutoEcoleID, "PLANNING_CONDUITE_UPDATED", map[string]interface{}{
			"action":    "DELETED",
			"seance_id": seanceID,
		})
		Success(w, http.StatusOK, map[string]string{"message": "Séance de conduite supprimée avec succès"})
		return
	}

	if r.Method != http.MethodGet {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	dateSeance := r.URL.Query().Get("date")
	moniteurID := r.URL.Query().Get("moniteur_id")

	if moniteurID == "" && claims.Role == "MONITEUR" {
		moniteurID = claims.UserID
	}

	list, err := h.guichetService.GetPlanningConduite(r.Context(), claims.AutoEcoleID, moniteurID, dateSeance)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if moniteurID != "" && moniteurID != "ALL" {
		filtered := make([]models.SeanceConduite, 0)
		for _, s := range list {
			if s.MoniteurID == nil || *s.MoniteurID == "" || *s.MoniteurID == moniteurID {
				filtered = append(filtered, s)
			}
		}
		if len(filtered) == 0 && len(list) > 0 && claims.Role == "MONITEUR" {
			Success(w, http.StatusOK, list)
			return
		}
		Success(w, http.StatusOK, filtered)
		return
	}

	Success(w, http.StatusOK, list)
}

// POST /api/pedagogie/seance-conduite/valider
func (h *GuichetHandler) ValiderSeanceConduite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims, _, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	var input models.ValiderSeanceConduiteInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps JSON invalide")
		return
	}

	res, err := h.guichetService.ValiderSeanceConduite(r.Context(), claims.AutoEcoleID, claims.UserID, &input)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	realtime.Broadcast(claims.AutoEcoleID, "LIVRET_EMARGEMENT_UPDATED", res)
	realtime.Broadcast(claims.AutoEcoleID, "PLANNING_CONDUITE_UPDATED", res)
	realtime.Broadcast(claims.AutoEcoleID, "FLOTTE_VEHICULES_UPDATED", map[string]interface{}{
		"vehicule_id": input.VehiculeID,
		"km_fin":      input.KmFin,
	})
	Success(w, http.StatusOK, res)
}

// POST /api/pedagogie/reserver-conduite
func (h *GuichetHandler) ReserverSeanceConduite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	claims, _, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	var input models.ReserverSeanceConduiteInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Corps JSON invalide")
		return
	}

	seance, err := h.guichetService.ReserverSeanceConduite(r.Context(), claims.AutoEcoleID, input)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	realtime.Broadcast(claims.AutoEcoleID, "PLANNING_CONDUITE_UPDATED", seance)
	Success(w, http.StatusCreated, seance)
}

// GET & POST /api/pedagogie/seances-code
func (h *GuichetHandler) HandleSeancesCode(w http.ResponseWriter, r *http.Request) {
	claims, agenceID, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	switch r.Method {
	case http.MethodGet:
		list, err := h.guichetService.GetSeancesCode(r.Context(), claims.AutoEcoleID, r.URL.Query().Get("type"))
		if err != nil {
			Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		Success(w, http.StatusOK, list)
	case http.MethodPost:
		role := models.NormalizeRole(claims.Role)
		if role != models.RoleDirecteur && role != models.RoleSecretaire && role != models.RoleMoniteur {
			Error(w, http.StatusForbidden, "Privilèges insuffisants : seul le personnel pédagogique peut planifier une séance")
			return
		}
		var input models.CreerSeanceCodeInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			Error(w, http.StatusBadRequest, "Corps JSON invalide")
			return
		}
		seance, err := h.guichetService.CreerSeanceCode(r.Context(), claims.AutoEcoleID, agenceID, input)
		if err != nil {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}
		realtime.Broadcast(claims.AutoEcoleID, "COURS_CODE_UPDATED", seance)
		Success(w, http.StatusCreated, seance)
	case http.MethodPut:
		role := models.NormalizeRole(claims.Role)
		if role != models.RoleDirecteur && role != models.RoleSecretaire {
			Error(w, http.StatusForbidden, "Privilèges insuffisants : seul le directeur ou la secrétaire peut reprogrammer un cours bloqué")
			return
		}
		var body struct {
			ID         string `json:"id"`
			DateSeance string `json:"date_seance"`
			HeureDebut string `json:"heure_debut"`
			HeureFin   string `json:"heure_fin"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			Error(w, http.StatusBadRequest, "Corps JSON invalide")
			return
		}
		seance, err := h.guichetService.ReprogrammerSeanceCode(r.Context(), claims.AutoEcoleID, body.ID, body.DateSeance, body.HeureDebut, body.HeureFin)
		if err != nil {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}
		realtime.Broadcast(claims.AutoEcoleID, "COURS_CODE_UPDATED", seance)
		Success(w, http.StatusOK, seance)
	default:
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
	}
}

// GET, POST, DELETE /api/pedagogie/tutoriels
func (h *GuichetHandler) HandleTutorielsVideos(w http.ResponseWriter, r *http.Request) {
	claims, _, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	switch r.Method {
	case http.MethodGet:
		list, err := h.guichetService.GetTutorielsVideos(r.Context(), claims.AutoEcoleID, r.URL.Query().Get("categorie"))
		if err != nil {
			Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		Success(w, http.StatusOK, list)
	case http.MethodPost:
		role := models.NormalizeRole(claims.Role)
		if role != models.RoleDirecteur && role != models.RoleSecretaire && role != models.RoleMoniteur {
			Error(w, http.StatusForbidden, "Privilèges insuffisants : seul le personnel pédagogique peut ajouter un tutoriel")
			return
		}
		var input models.CreerTutorielVideoInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			Error(w, http.StatusBadRequest, "Corps JSON invalide")
			return
		}
		tut, err := h.guichetService.CreerTutorielVideo(r.Context(), claims.AutoEcoleID, input)
		if err != nil {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}
		realtime.Broadcast(claims.AutoEcoleID, "TUTORIELS_UPDATED", tut)
		Success(w, http.StatusCreated, tut)
	case http.MethodDelete:
		role := models.NormalizeRole(claims.Role)
		if role != models.RoleDirecteur && role != models.RoleSecretaire {
			Error(w, http.StatusForbidden, "Privilèges insuffisants : seul le Directeur ou la Secrétaire peut supprimer un tutoriel")
			return
		}
		id := r.URL.Query().Get("id")
		if id == "" {
			Error(w, http.StatusBadRequest, "Paramètre 'id' requis")
			return
		}
		if err := h.guichetService.SupprimerTutorielVideo(r.Context(), claims.AutoEcoleID, id); err != nil {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}
		realtime.Broadcast(claims.AutoEcoleID, "TUTORIELS_UPDATED", map[string]interface{}{
			"deleted": true,
			"id":      id,
		})
		Success(w, http.StatusOK, map[string]bool{"deleted": true})
	default:
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
	}
}

// POST /api/pedagogie/tutoriels/vue?id=...
func (h *GuichetHandler) IncrementerVueTutoriel(w http.ResponseWriter, r *http.Request) {
	claims, _, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		Error(w, http.StatusBadRequest, "Paramètre 'id' requis")
		return
	}
	_ = h.guichetService.IncrementerVueTutoriel(r.Context(), claims.AutoEcoleID, id)
	realtime.Broadcast(claims.AutoEcoleID, "TUTORIELS_UPDATED", map[string]interface{}{
		"incremented": true,
		"id":          id,
	})
	Success(w, http.StatusOK, map[string]bool{"incremented": true})
}

// GET, POST, DELETE /api/pedagogie/supports-cours
func (h *GuichetHandler) HandleSupportsCours(w http.ResponseWriter, r *http.Request) {
	claims, _, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}

	switch r.Method {
	case http.MethodGet:
		list, err := h.guichetService.GetSupportsCours(r.Context(), claims.AutoEcoleID, r.URL.Query().Get("categorie"))
		if err != nil {
			Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		Success(w, http.StatusOK, list)
	case http.MethodPost:
		role := models.NormalizeRole(claims.Role)
		if role != models.RoleDirecteur && role != models.RoleSecretaire && role != models.RoleMoniteur {
			Error(w, http.StatusForbidden, "Privilèges insuffisants : seul le personnel pédagogique peut ajouter un support de cours")
			return
		}
		var input models.CreerSupportCoursInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			Error(w, http.StatusBadRequest, "Corps JSON invalide")
			return
		}
		sup, err := h.guichetService.CreerSupportCours(r.Context(), claims.AutoEcoleID, input)
		if err != nil {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}
		realtime.Broadcast(claims.AutoEcoleID, "SUPPORTS_COURS_UPDATED", sup)
		Success(w, http.StatusCreated, sup)
	case http.MethodDelete:
		role := models.NormalizeRole(claims.Role)
		if role != models.RoleDirecteur && role != models.RoleSecretaire {
			Error(w, http.StatusForbidden, "Privilèges insuffisants : seul le Directeur ou la Secrétaire peut supprimer un support")
			return
		}
		id := r.URL.Query().Get("id")
		if id == "" {
			Error(w, http.StatusBadRequest, "Paramètre 'id' requis")
			return
		}
		if err := h.guichetService.SupprimerSupportCours(r.Context(), claims.AutoEcoleID, id); err != nil {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}
		realtime.Broadcast(claims.AutoEcoleID, "SUPPORTS_COURS_UPDATED", map[string]interface{}{
			"deleted": true,
			"id":      id,
		})
		Success(w, http.StatusOK, map[string]bool{"deleted": true})
	default:
		Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
	}
}

// POST /api/pedagogie/supports-cours/telecharger?id=...
func (h *GuichetHandler) IncrementerTelechargementSupport(w http.ResponseWriter, r *http.Request) {
	claims, _, ok := requireClaimsWithAgence(w, r)
	if !ok {
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		Error(w, http.StatusBadRequest, "Paramètre 'id' requis")
		return
	}
	_ = h.guichetService.IncrementerTelechargementSupport(r.Context(), claims.AutoEcoleID, id)
	realtime.Broadcast(claims.AutoEcoleID, "SUPPORTS_COURS_UPDATED", map[string]interface{}{
		"incremented": true,
		"id":          id,
	})
	Success(w, http.StatusOK, map[string]bool{"incremented": true})
}


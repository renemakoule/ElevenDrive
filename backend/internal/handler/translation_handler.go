package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"autoschool/backend/internal/service"
)

// TranslationHandler gère les requêtes HTTP de traduction de l'application.
type TranslationHandler struct {
	service service.TranslationService
}

// NewTranslationHandler instancie le contrôleur HTTP de traduction.
func NewTranslationHandler(s service.TranslationService) *TranslationHandler {
	return &TranslationHandler{service: s}
}

// TranslateRequest structure de la charge utile de traduction reçue du frontend.
type TranslateRequest struct {
	Text       string   `json:"text,omitempty"`
	Texts      []string `json:"texts,omitempty"`
	TargetLang string   `json:"target_lang"`
	SourceLang string   `json:"source_lang,omitempty"`
}

// Translate traite les demandes de traduction envoyées par l'application Desktop ou Web.
// POST /api/traduire
func (h *TranslationHandler) Translate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "Méthode HTTP non autorisée, utilisez POST")
		return
	}

	var req TranslateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "Corps de requête JSON invalide")
		return
	}

	targetLang := strings.TrimSpace(req.TargetLang)
	if targetLang == "" {
		Error(w, http.StatusBadRequest, "Le champ 'target_lang' est obligatoire (ex: 'en', 'es', 'ar')")
		return
	}

	// 1. Traduction par lot (Batch)
	if len(req.Texts) > 0 {
		translatedList, err := h.service.TranslateBatch(r.Context(), req.Texts, targetLang, req.SourceLang)
		if err != nil {
			Error(w, http.StatusBadGateway, "Erreur lors de la traduction : "+err.Error())
			return
		}
		Success(w, http.StatusOK, map[string]interface{}{
			"target_lang":  targetLang,
			"translations": translatedList,
		})
		return
	}

	// 2. Traduction d'un texte simple
	if strings.TrimSpace(req.Text) == "" {
		Success(w, http.StatusOK, map[string]interface{}{
			"original":    "",
			"translation": "",
			"target_lang": targetLang,
		})
		return
	}

	translated, err := h.service.Translate(r.Context(), req.Text, targetLang, req.SourceLang)
	if err != nil {
		// En cas d'erreur de clé ou de réseau, on renvoie le texte avec code 200 en fallback
		Success(w, http.StatusOK, map[string]interface{}{
			"original":    req.Text,
			"translation": req.Text,
			"target_lang": targetLang,
			"warning":     err.Error(),
		})
		return
	}

	Success(w, http.StatusOK, map[string]interface{}{
		"original":    req.Text,
		"translation": translated,
		"target_lang": targetLang,
	})
}

// Statut renvoie l'état de configuration du moteur DeepL.
// GET /api/traduire/statut
func (h *TranslationHandler) Statut(w http.ResponseWriter, r *http.Request) {
	Success(w, http.StatusOK, map[string]interface{}{
		"provider":      "DeepL Official API",
		"is_configured": h.service.IsConfigured(),
	})
}

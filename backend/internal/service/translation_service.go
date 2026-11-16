package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// TranslationService interface abstraite pour le service de traduction automatique.
type TranslationService interface {
	Translate(ctx context.Context, text string, targetLang string, sourceLang string) (string, error)
	TranslateBatch(ctx context.Context, texts []string, targetLang string, sourceLang string) ([]string, error)
	IsConfigured() bool
}

type deepLTranslationService struct {
	apiKey     string
	isPro      bool
	endpoint   string
	httpClient *http.Client
	cache      map[string]string
	mu         sync.RWMutex
}

// NewDeepLTranslationService initialise le service DeepL officiel côté serveur.
func NewDeepLTranslationService(apiKey string) TranslationService {
	cleanKey := strings.TrimSpace(apiKey)
	isPro := len(cleanKey) > 0 && !strings.HasSuffix(cleanKey, ":fx")

	endpoint := "https://api-free.deepl.com/v2/translate"
	if isPro {
		endpoint = "https://api.deepl.com/v2/translate"
	}

	return &deepLTranslationService{
		apiKey:     cleanKey,
		isPro:      isPro,
		endpoint:   endpoint,
		httpClient: &http.Client{Timeout: 12 * time.Second},
		cache:      make(map[string]string),
	}
}

// IsConfigured indique si une clé API DeepL est configurée.
func (s *deepLTranslationService) IsConfigured() bool {
	return len(s.apiKey) > 0
}

// normalizeTarget adapte le code langue aux spécifications exactes de DeepL API.
func normalizeTarget(lang string) string {
	l := strings.TrimSpace(strings.ToUpper(lang))
	switch l {
	case "EN", "EN-US", "EN-GB", "ENGLISH":
		return "EN-US"
	case "ES", "SPANISH", "ESPAÑOL":
		return "ES"
	case "AR", "ARABIC", "العربية":
		return "AR"
	case "FR", "FRENCH", "FRANÇAIS":
		return "FR"
	default:
		return l
	}
}

// normalizeSource adapte le code langue source.
func normalizeSource(lang string) string {
	l := strings.TrimSpace(strings.ToUpper(lang))
	if l == "" || l == "FR" || l == "FRANÇAIS" {
		return "FR"
	}
	return l
}

type deepLRequestPayload struct {
	Text       []string `json:"text"`
	TargetLang string   `json:"target_lang"`
	SourceLang string   `json:"source_lang,omitempty"`
}

type deepLResponsePayload struct {
	Translations []struct {
		DetectedSourceLanguage string `json:"detected_source_language"`
		Text                   string `json:"text"`
	} `json:"translations"`
	Message string `json:"message,omitempty"`
}

const maxTranslationCacheEntries = 2000

func (s *deepLTranslationService) setCache(key, val string) {
	if len(s.cache) >= maxTranslationCacheEntries {
		counter := 0
		for k := range s.cache {
			delete(s.cache, k)
			counter++
			if counter >= maxTranslationCacheEntries/5 {
				break
			}
		}
	}
	s.cache[key] = val
}

// Translate traduit un texte unique avec mise en cache mémoire thread-safe et bornée.
func (s *deepLTranslationService) Translate(ctx context.Context, text string, targetLang string, sourceLang string) (string, error) {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return "", nil
	}

	target := normalizeTarget(targetLang)
	source := normalizeSource(sourceLang)

	// Pas de traduction si la cible est identique à la source
	if target == source || (target == "FR" && source == "FR") {
		return text, nil
	}

	cacheKey := target + "::" + trimmed

	// Vérification dans le cache mémoire
	s.mu.RLock()
	cached, found := s.cache[cacheKey]
	s.mu.RUnlock()
	if found {
		return cached, nil
	}

	// Si aucune clé configurée, repli gracieux sur le texte original
	if !s.IsConfigured() {
		log.Printf("⚠️ DeepL: aucune clé DEEPL_API_KEY configurée. Repli sur le texte original.")
		return text, nil
	}

	results, err := s.TranslateBatch(ctx, []string{trimmed}, target, source)
	if err != nil {
		return text, err
	}

	if len(results) > 0 && results[0] != "" {
		s.mu.Lock()
		s.setCache(cacheKey, results[0])
		s.mu.Unlock()
		return results[0], nil
	}

	return text, nil
}

// TranslateBatch traduit un ensemble de textes en une seule requête HTTP DeepL
// avec mise en cache mémoire thread-safe pour économiser le quota.
func (s *deepLTranslationService) TranslateBatch(ctx context.Context, texts []string, targetLang string, sourceLang string) ([]string, error) {
	if len(texts) == 0 {
		return []string{}, nil
	}

	target := normalizeTarget(targetLang)
	source := strings.TrimSpace(strings.ToUpper(sourceLang))

	// Pas de traduction si la cible est le français (langue native de l'application)
	if target == "FR" {
		return texts, nil
	}

	if !s.IsConfigured() {
		return texts, fmt.Errorf("clé API DeepL non configurée (définir DEEPL_API_KEY dans backend/.env)")
	}

	results := make([]string, len(texts))
	var uncachedIndices []int
	var uncachedTexts []string

	// 1. Consulter le cache mémoire thread-safe
	s.mu.RLock()
	for i, t := range texts {
		trimmed := strings.TrimSpace(t)
		if trimmed == "" {
			results[i] = t
			continue
		}
		cacheKey := target + "::" + trimmed
		if cachedVal, ok := s.cache[cacheKey]; ok {
			results[i] = cachedVal
		} else {
			uncachedIndices = append(uncachedIndices, i)
			uncachedTexts = append(uncachedTexts, trimmed)
		}
	}
	s.mu.RUnlock()

	// Si tous les textes étaient déjà en cache, renvoyer immédiatement (0 ms, 0 caractère consommé)
	if len(uncachedTexts) == 0 {
		return results, nil
	}

	// 2. Préparer la charge utile pour DeepL API
	payload := deepLRequestPayload{
		Text:       uncachedTexts,
		TargetLang: target,
	}
	if source != "" && source != "AUTO" {
		payload.SourceLang = source
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return texts, fmt.Errorf("erreur de sérialisation JSON DeepL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return texts, fmt.Errorf("erreur de création de la requête HTTP DeepL: %w", err)
	}

	req.Header.Set("Authorization", "DeepL-Auth-Key "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		log.Printf("❌ DeepL HTTP Error: %v", err)
		return texts, fmt.Errorf("échec de l'appel HTTP DeepL: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return texts, fmt.Errorf("erreur de lecture du corps DeepL: %w", err)
	}

	if resp.StatusCode == http.StatusForbidden {
		log.Printf("❌ DeepL API: Clé invalide ou refusée (HTTP 403)")
		return texts, fmt.Errorf("clé API DeepL non valide ou accès interdit")
	}

	if resp.StatusCode == 456 { // Quota Exceeded
		log.Printf("⚠️ DeepL API: Quota mensuel de caractères épuisé (HTTP 456)")
		return texts, fmt.Errorf("quota mensuel DeepL épuisé")
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("❌ DeepL API HTTP %d: %s", resp.StatusCode, string(respBody))
		return texts, fmt.Errorf("DeepL API a renvoyé le statut HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var parsed deepLResponsePayload
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return texts, fmt.Errorf("erreur de décodage JSON DeepL: %w", err)
	}

	// 3. Réinjecter les traductions reçues dans le tableau de résultats et dans le cache
	s.mu.Lock()
	for idx, item := range parsed.Translations {
		if idx < len(uncachedIndices) {
			origIndex := uncachedIndices[idx]
			origText := uncachedTexts[idx]
			results[origIndex] = item.Text
			s.setCache(target+"::"+origText, item.Text)
		}
	}
	s.mu.Unlock()

	return results, nil
}

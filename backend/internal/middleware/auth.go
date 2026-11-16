package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"autoschool/backend/internal/models"
	"autoschool/backend/internal/service"
)

type contextKey string

const (
	claimsContextKey contextKey = "user_claims"
)

// RequireAuth intercepte la requête, valide le jeton Bearer JWT et stocke les claims dans le contexte.
func RequireAuth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondUnauthorized(w, "En-tête d'autorisation manquant")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				respondUnauthorized(w, "Format du jeton Bearer invalide")
				return
			}

			claims, err := service.ValidateAccessToken(parts[1], jwtSecret)
			if err != nil {
				respondUnauthorized(w, "Jeton d'authentification invalide ou expiré")
				return
			}

			// Injection des claims dans le contexte HTTP
			ctx := context.WithValue(r.Context(), claimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole vérifie que l'utilisateur connecté possède l'un des rôles requis.
func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaims(r.Context())
			if claims == nil {
				respondUnauthorized(w, "Authentification requise")
				return
			}

			userRole := models.NormalizeRole(claims.Role)
			isAuthorized := false
			for _, role := range allowedRoles {
				normalizedAllowed := models.NormalizeRole(role)
				if strings.EqualFold(userRole, normalizedAllowed) {
					isAuthorized = true
					break
				}
			}

			if !isAuthorized {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"success": false,
					"error":   "Accès interdit : privilèges insuffisants",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetClaims extrait les claims JWT du contexte de la requête.
func GetClaims(ctx context.Context) *models.JWTClaims {
	claims, ok := ctx.Value(claimsContextKey).(*models.JWTClaims)
	if !ok {
		return nil
	}
	return claims
}

// GetUser renvoie les claims de l'utilisateur connecté depuis le contexte.
func GetUser(ctx context.Context) *models.JWTClaims {
	return GetClaims(ctx)
}

func respondUnauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   msg,
	})
}

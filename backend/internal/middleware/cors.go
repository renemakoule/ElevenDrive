package middleware

import (
	"net/http"
	"strings"
)

// CORS configure les en-têtes Cross-Origin Resource Sharing.
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Vérifier si l'origine est reconnue et autorisée
			allowWildcard := len(allowedOrigins) == 1 && allowedOrigins[0] == "*"
			isTrustedOrigin := false

			if origin != "" {
				if strings.HasPrefix(origin, "http://localhost:") ||
					strings.HasPrefix(origin, "http://127.0.0.1:") ||
					strings.HasPrefix(origin, "http://tauri.localhost") ||
					strings.HasPrefix(origin, "https://tauri.localhost") ||
					strings.HasPrefix(origin, "tauri://") {
					isTrustedOrigin = true
				} else {
					for _, allowed := range allowedOrigins {
						if allowed != "*" && strings.EqualFold(allowed, origin) {
							isTrustedOrigin = true
							break
						}
					}
				}
			}

			if isTrustedOrigin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			} else if allowWildcard {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}

			if isTrustedOrigin || allowWildcard {
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-Tenant-ID")
				w.Header().Set("Access-Control-Max-Age", "86400")
			}

			// Traiter la pré-requête OPTIONS immédiatement
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

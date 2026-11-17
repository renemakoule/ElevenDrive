package middleware

import (
	"net/http"
	"strings"
)

// CORS configure les en-têtes Cross-Origin Resource Sharing de façon robuste
// pour le client Desktop Tauri, le Web PWA, le développement local et la production.
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Déterminer si l'accès est autorisé
			allowWildcard := len(allowedOrigins) == 0 || (len(allowedOrigins) == 1 && allowedOrigins[0] == "*")
			isTrustedOrigin := allowWildcard

			if origin != "" && !isTrustedOrigin {
				// Autoriser systématiquement les environnements de développement et desktop natifs
				if strings.HasPrefix(origin, "http://localhost") ||
					strings.HasPrefix(origin, "https://localhost") ||
					strings.HasPrefix(origin, "http://127.0.0.1") ||
					strings.HasPrefix(origin, "https://127.0.0.1") ||
					strings.Contains(origin, "tauri.localhost") ||
					strings.HasPrefix(origin, "tauri://") ||
					strings.Contains(origin, "ipc.localhost") ||
					strings.HasSuffix(origin, ".onrender.com") ||
					strings.HasSuffix(origin, ".netlify.app") ||
					strings.HasSuffix(origin, ".vercel.app") ||
					strings.HasSuffix(origin, ".fly.dev") {
					isTrustedOrigin = true
				} else {
					for _, allowed := range allowedOrigins {
						if allowed == "*" || strings.EqualFold(allowed, origin) {
							isTrustedOrigin = true
							break
						}
					}
				}
			}

			// Injection systématique des en-têtes CORS
			if origin != "" && (isTrustedOrigin || allowWildcard) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			} else if allowWildcard {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}

			if isTrustedOrigin || allowWildcard {
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD")
				
				// Renvoyer les en-têtes demandés par le navigateur pour garantir le passage du pré-vol (OPTIONS)
				if reqHeaders := r.Header.Get("Access-Control-Request-Headers"); reqHeaders != "" {
					w.Header().Set("Access-Control-Allow-Headers", reqHeaders)
				} else {
					w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-Tenant-ID, Accept, Origin")
				}
				
				w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Range, Content-Disposition, X-Tenant-ID")
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

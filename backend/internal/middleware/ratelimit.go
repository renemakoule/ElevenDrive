package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type tokenBucket struct {
	tokens   float64
	lastSeen time.Time
}

type ipLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*tokenBucket
	rate     float64 // jetons ajoutés par seconde
	burst    float64 // capacité maximale de rafale
	maxIps   int
	lastPurge time.Time
}

func newIPLimiter(ratePerSec, burst float64) *ipLimiter {
	return &ipLimiter{
		buckets:   make(map[string]*tokenBucket),
		rate:      ratePerSec,
		burst:     burst,
		maxIps:    50000,
		lastPurge: time.Now(),
	}
}

func (l *ipLimiter) allow(ip string) bool {
	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()

	if now.Sub(l.lastPurge) > 2*time.Minute || len(l.buckets) > l.maxIps {
		for k, b := range l.buckets {
			if now.Sub(b.lastSeen) > 5*time.Minute {
				delete(l.buckets, k)
			}
		}
		if len(l.buckets) > l.maxIps {
			l.buckets = make(map[string]*tokenBucket)
		}
		l.lastPurge = now
	}

	b, ok := l.buckets[ip]
	if !ok {
		l.buckets[ip] = &tokenBucket{
			tokens:   l.burst - 1,
			lastSeen: now,
		}
		return true
	}

	elapsed := now.Sub(b.lastSeen).Seconds()
	b.tokens += elapsed * l.rate
	if b.tokens > l.burst {
		b.tokens = l.burst
	}
	b.lastSeen = now

	if b.tokens < 1 {
		return false
	}
	b.tokens -= 1
	return true
}

func extractClientIP(r *http.Request) string {
	if cf := strings.TrimSpace(r.Header.Get("CF-Connecting-IP")); cf != "" {
		return cf
	}
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if ip := strings.TrimSpace(parts[0]); ip != "" {
			return ip
		}
	}
	if xri := strings.TrimSpace(r.Header.Get("X-Real-IP")); xri != "" {
		return xri
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}
	return r.RemoteAddr
}

// RateLimit protège le serveur contre les rafales DoS et le brute-force par IP.
func RateLimit(next http.Handler) http.Handler {
	strictLimiter := newIPLimiter(5, 25)   // Auth, SuperAdmin login, Traduction publique
	generalLimiter := newIPLimiter(30, 150) // Routes API standards

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		ip := extractClientIP(r)
		path := r.URL.Path

		isStrict := strings.HasPrefix(path, "/api/auth/") ||
			strings.HasPrefix(path, "/api/traduire") ||
			path == "/api/superadmin/login"

		allowed := true
		if isStrict {
			allowed = strictLimiter.allow(ip)
		} else {
			allowed = generalLimiter.allow(ip)
		}

		if !allowed {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "2")
			w.WriteHeader(http.StatusTooManyRequests)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Trop de requêtes détectées. Veuillez patienter quelques secondes.",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

// SecurityAndBodyLimit applique les en-têtes de sécurité et limite la taille des corps HTTP.
func SecurityAndBodyLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		if r.Body != nil && r.Method != http.MethodGet && r.Method != http.MethodHead {
			var maxBytes int64 = 10 << 20 // 10 Mo par défaut
			path := r.URL.Path
			if strings.Contains(path, "/upload") ||
				strings.HasPrefix(path, "/api/epreuves") ||
				strings.HasPrefix(path, "/api/pedagogie/supports-cours") ||
				strings.HasPrefix(path, "/api/archivage") {
				maxBytes = 25 << 20 // 25 Mo pour les documents PDF / archives
			}
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
		}

		next.ServeHTTP(w, r)
	})
}


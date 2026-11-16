package middleware

import (
	"bufio"
	"errors"
	"log"
	"net"
	"net/http"
	"time"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Hijack délègue le hijacking TCP sous-jacent pour les connexions WebSocket temps réel.
func (rw *responseWriterWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := rw.ResponseWriter.(http.Hijacker); ok {
		rw.statusCode = http.StatusSwitchingProtocols
		return hijacker.Hijack()
	}
	return nil, nil, errors.New("http.Hijacker non supporté par le ResponseWriter")
}

// Flush délègue le flush HTTP pour le streaming temps réel.
func (rw *responseWriterWrapper) Flush() {
	if flusher, ok := rw.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Logger enregistre chaque requête entrante avec sa méthode, son URL, son code HTTP et sa durée.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapper := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapper, r)

		duration := time.Since(start)
		log.Printf("[%s] %s %d - %v", r.Method, r.URL.Path, wrapper.statusCode, duration)
	})
}

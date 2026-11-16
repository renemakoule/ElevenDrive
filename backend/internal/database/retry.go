package database

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// IsTransientError détecte les erreurs réseau / pooler (Neon PgBouncer) qui disparaissent au nouvel essai.
func IsTransientError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	for _, s := range []string{
		"unnamed prepared statement does not exist",
		"prepared statement",
		"bad connection",
		"connection reset",
		"broken pipe",
		"unexpected eof",
		"driver: bad connection",
	} {
		if strings.Contains(msg, s) {
			return true
		}
	}
	return false
}

// Safe exécute fn et convertit un panic du driver (ex: lib/pq "index out of range" sur flux TCP corrompu)
// en erreur transitoire, afin de pouvoir réessayer au lieu de couper la connexion HTTP.
func Safe(fn func() error) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("driver: bad connection (panic récupéré: %v)", rec)
		}
	}()
	return fn()
}

// Retry exécute fn jusqu'à 3 fois en cas d'erreur transitoire (lecture idempotente uniquement).
func Retry(ctx context.Context, fn func() error) error {
	var err error
	for attempt := 0; attempt < 3; attempt++ {
		if err = fn(); err == nil || !IsTransientError(err) {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(attempt+1) * 300 * time.Millisecond):
		}
	}
	return err
}


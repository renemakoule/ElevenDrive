package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config regroupe les paramètres d'exécution de l'application.
// Aucun mot de passe ni identifiant factice n'est défini par défaut.
type Config struct {
	Port                 string
	AppEnv               string
	FrontendURL          string
	PublicAPIURL         string
	DatabaseURL          string
	JWTSecret            string
	JWTExpirationHours   int
	RefreshExpirationDay int
	AllowedOrigins       []string
	DeepLAPIKey          string
	R2AccountID          string
	R2AccessKeyID        string
	R2SecretAccessKey    string
	R2BucketName         string
	R2PublicURL          string
	GoogleClientID       string
	GoogleClientSecret   string
	GoogleRedirectURI    string
	FacebookAppID        string
	FacebookAppSecret    string
	FacebookRedirectURI  string
	LinkedInClientID     string
	LinkedInClientSecret string
	LinkedInRedirectURI  string
	AppleClientID        string
	AppleTeamID          string
	AppleKeyID           string
	AppleRedirectURI     string
}

// Load charge la configuration depuis l'environnement ou un fichier .env.
// Renvoie une erreur si les variables obligatoires (DATABASE_URL, JWT_SECRET) sont absentes.
func Load() (*Config, error) {
	_ = godotenv.Load()
	_ = godotenv.Load("backend/.env")
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	appEnv := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	if appEnv == "" {
		appEnv = strings.ToLower(strings.TrimSpace(os.Getenv("GO_ENV")))
	}

	frontendURL := strings.TrimRight(strings.TrimSpace(os.Getenv("FRONTEND_URL")), "/")
	if frontendURL == "" && appEnv != "production" {
		frontendURL = "http://localhost:1420"
	}

	publicAPIURL := strings.TrimRight(strings.TrimSpace(os.Getenv("PUBLIC_API_URL")), "/")
	if publicAPIURL == "" {
		publicAPIURL = strings.TrimRight(strings.TrimSpace(os.Getenv("BACKEND_URL")), "/")
	}
	if publicAPIURL == "" && appEnv != "production" {
		publicAPIURL = "http://localhost:" + port
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if strings.TrimSpace(databaseURL) == "" {
		return nil, errors.New("la variable d'environnement DATABASE_URL est requise et ne peut être vide")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if strings.TrimSpace(jwtSecret) == "" {
		return nil, errors.New("la variable d'environnement JWT_SECRET est requise et ne peut être vide")
	}
	if len(jwtSecret) < 16 {
		return nil, fmt.Errorf("la clé JWT_SECRET doit contenir au moins 16 caractères pour des raisons de sécurité")
	}

	jwtExpHours := 24
	if expStr := os.Getenv("JWT_EXPIRATION_HOURS"); expStr != "" {
		if val, err := strconv.Atoi(expStr); err == nil && val > 0 {
			jwtExpHours = val
		}
	}

	refreshExpDays := 30
	if refStr := os.Getenv("REFRESH_EXPIRATION_DAYS"); refStr != "" {
		if val, err := strconv.Atoi(refStr); err == nil && val > 0 {
			refreshExpDays = val
		}
	}

	allowedOrigins := []string{"*"}
	if originsStr := os.Getenv("ALLOWED_ORIGINS"); originsStr != "" {
		parts := strings.Split(originsStr, ",")
		cleaned := make([]string, 0, len(parts))
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				cleaned = append(cleaned, trimmed)
			}
		}
		if len(cleaned) > 0 {
			allowedOrigins = cleaned
		}
	}

	deeplKey := strings.TrimSpace(os.Getenv("DEEPL_API_KEY"))
	if deeplKey == "" {
		deeplKey = strings.TrimSpace(os.Getenv("VITE_DEEPL_API_KEY"))
	}

	r2AccountID := strings.TrimSpace(os.Getenv("R2_ACCOUNT_ID"))
	r2AccessKeyID := strings.TrimSpace(os.Getenv("R2_ACCESS_KEY_ID"))
	r2SecretAccessKey := strings.TrimSpace(os.Getenv("R2_SECRET_ACCESS_KEY"))
	r2BucketName := strings.TrimSpace(os.Getenv("R2_BUCKET_NAME"))
	r2PublicURL := strings.TrimSpace(os.Getenv("R2_PUBLIC_URL"))

	googleClientID := strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID"))
	googleClientSecret := strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_SECRET"))
	googleRedirectURI := strings.TrimSpace(os.Getenv("GOOGLE_REDIRECT_URI"))

	facebookAppID := strings.TrimSpace(os.Getenv("FACEBOOK_APP_ID"))
	facebookAppSecret := strings.TrimSpace(os.Getenv("FACEBOOK_APP_SECRET"))
	facebookRedirectURI := strings.TrimSpace(os.Getenv("FACEBOOK_REDIRECT_URI"))

	linkedInClientID := strings.TrimSpace(os.Getenv("LINKEDIN_CLIENT_ID"))
	linkedInClientSecret := strings.TrimSpace(os.Getenv("LINKEDIN_CLIENT_SECRET"))
	linkedInRedirectURI := strings.TrimSpace(os.Getenv("LINKEDIN_REDIRECT_URI"))

	appleClientID := strings.TrimSpace(os.Getenv("APPLE_CLIENT_ID"))
	appleTeamID := strings.TrimSpace(os.Getenv("APPLE_TEAM_ID"))
	appleKeyID := strings.TrimSpace(os.Getenv("APPLE_KEY_ID"))
	appleRedirectURI := strings.TrimSpace(os.Getenv("APPLE_REDIRECT_URI"))

	return &Config{
		Port:                 port,
		AppEnv:               appEnv,
		FrontendURL:          frontendURL,
		PublicAPIURL:         publicAPIURL,
		DatabaseURL:          databaseURL,
		JWTSecret:            jwtSecret,
		JWTExpirationHours:   jwtExpHours,
		RefreshExpirationDay: refreshExpDays,
		AllowedOrigins:       allowedOrigins,
		DeepLAPIKey:          deeplKey,
		R2AccountID:          r2AccountID,
		R2AccessKeyID:        r2AccessKeyID,
		R2SecretAccessKey:    r2SecretAccessKey,
		R2BucketName:         r2BucketName,
		R2PublicURL:          r2PublicURL,
		GoogleClientID:       googleClientID,
		GoogleClientSecret:   googleClientSecret,
		GoogleRedirectURI:    googleRedirectURI,
		FacebookAppID:        facebookAppID,
		FacebookAppSecret:    facebookAppSecret,
		FacebookRedirectURI:  facebookRedirectURI,
		LinkedInClientID:     linkedInClientID,
		LinkedInClientSecret: linkedInClientSecret,
		LinkedInRedirectURI:  linkedInRedirectURI,
		AppleClientID:        appleClientID,
		AppleTeamID:          appleTeamID,
		AppleKeyID:           appleKeyID,
		AppleRedirectURI:     appleRedirectURI,
	}, nil
}

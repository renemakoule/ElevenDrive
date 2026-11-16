package service

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"autoschool/backend/internal/config"
)

type R2StorageService struct {
	accountID       string
	accessKeyID     string
	secretAccessKey string
	bucketName      string
	publicURL       string
	httpClient      *http.Client
}

func NewR2StorageService(cfg *config.Config) *R2StorageService {
	return &R2StorageService{
		accountID:       cfg.R2AccountID,
		accessKeyID:     cfg.R2AccessKeyID,
		secretAccessKey: cfg.R2SecretAccessKey,
		bucketName:      cfg.R2BucketName,
		publicURL:       strings.TrimSuffix(cfg.R2PublicURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// IsConfigured vérifie si les informations Cloudflare R2 sont renseignées
func (s *R2StorageService) IsConfigured() bool {
	return s.accountID != "" && s.accessKeyID != "" && s.secretAccessKey != "" && s.bucketName != ""
}

func (s *R2StorageService) GetBucketName() string {
	return s.bucketName
}

func (s *R2StorageService) GetPublicURL() string {
	return s.publicURL
}

func hmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}

func getSigningKey(secretKey, dateStamp, region, service string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+secretKey), dateStamp)
	kRegion := hmacSHA256(kDate, region)
	kService := hmacSHA256(kRegion, service)
	return hmacSHA256(kService, "aws4_request")
}

// UploadObject téléverse un objet binaire sur Cloudflare R2 avec signature AWS v4
func (s *R2StorageService) UploadObject(ctx context.Context, objectKey string, data []byte, contentType string) (string, error) {
	if !s.IsConfigured() {
		return "", errors.New("les identifiants Cloudflare R2 ne sont pas configurés dans .env")
	}

	objectKey = strings.TrimPrefix(objectKey, "/")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	endpointHost := fmt.Sprintf("%s.r2.cloudflarestorage.com", s.accountID)
	requestURI := fmt.Sprintf("/%s/%s", s.bucketName, objectKey)
	fullURL := fmt.Sprintf("https://%s%s", endpointHost, requestURI)

	now := time.Now().UTC()
	amzDate := now.Format("20060102T150405Z")
	dateStamp := now.Format("20060102")
	region := "auto"
	serviceName := "s3"

	// 1. Hash SHA-256 du contenu
	contentHash := sha256.Sum256(data)
	contentHashHex := hex.EncodeToString(contentHash[:])

	// 2. Création de la requête HTTP
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, fullURL, bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("création requête R2: %w", err)
	}

	req.Header.Set("Host", endpointHost)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-amz-date", amzDate)
	req.Header.Set("x-amz-content-sha256", contentHashHex)

	// 3. Construction de la requête canonique
	canonicalHeaders := fmt.Sprintf("content-type:%s\nhost:%s\nx-amz-content-sha256:%s\nx-amz-date:%s\n",
		contentType, endpointHost, contentHashHex, amzDate)
	signedHeaders := "content-type;host;x-amz-content-sha256;x-amz-date"

	canonicalRequest := fmt.Sprintf("PUT\n%s\n\n%s\n%s\n%s",
		requestURI, canonicalHeaders, signedHeaders, contentHashHex)

	canonicalRequestHash := sha256.Sum256([]byte(canonicalRequest))
	canonicalRequestHashHex := hex.EncodeToString(canonicalRequestHash[:])

	// 4. Chaîne à signer
	credentialScope := fmt.Sprintf("%s/%s/%s/aws4_request", dateStamp, region, serviceName)
	stringToSign := fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s\n%s",
		amzDate, credentialScope, canonicalRequestHashHex)

	// 5. Calcul de la signature
	signingKey := getSigningKey(s.secretAccessKey, dateStamp, region, serviceName)
	signature := hex.EncodeToString(hmacSHA256(signingKey, stringToSign))

	// 6. Header Authorization
	authHeader := fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		s.accessKeyID, credentialScope, signedHeaders, signature)
	req.Header.Set("Authorization", authHeader)

	// 7. Envoi de la requête
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("erreur transport Cloudflare R2: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("erreur R2 status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// 8. Retour de l'URL publique Cloudflare R2
	if s.publicURL != "" {
		return fmt.Sprintf("%s/%s", s.publicURL, objectKey), nil
	}
	return fullURL, nil
}

// DeleteObject supprime un objet dans Cloudflare R2
func (s *R2StorageService) DeleteObject(ctx context.Context, objectKey string) error {
	if !s.IsConfigured() {
		return errors.New("les identifiants Cloudflare R2 ne sont pas configurés")
	}

	objectKey = strings.TrimPrefix(objectKey, "/")
	endpointHost := fmt.Sprintf("%s.r2.cloudflarestorage.com", s.accountID)
	requestURI := fmt.Sprintf("/%s/%s", s.bucketName, objectKey)
	fullURL := fmt.Sprintf("https://%s%s", endpointHost, requestURI)

	now := time.Now().UTC()
	amzDate := now.Format("20060102T150405Z")
	dateStamp := now.Format("20060102")
	region := "auto"
	serviceName := "s3"
	contentHashHex := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" // Hash vide

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fullURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Host", endpointHost)
	req.Header.Set("x-amz-date", amzDate)
	req.Header.Set("x-amz-content-sha256", contentHashHex)

	canonicalHeaders := fmt.Sprintf("host:%s\nx-amz-content-sha256:%s\nx-amz-date:%s\n",
		endpointHost, contentHashHex, amzDate)
	signedHeaders := "host;x-amz-content-sha256;x-amz-date"

	canonicalRequest := fmt.Sprintf("DELETE\n%s\n\n%s\n%s\n%s",
		requestURI, canonicalHeaders, signedHeaders, contentHashHex)

	canonicalRequestHash := sha256.Sum256([]byte(canonicalRequest))
	canonicalRequestHashHex := hex.EncodeToString(canonicalRequestHash[:])

	credentialScope := fmt.Sprintf("%s/%s/%s/aws4_request", dateStamp, region, serviceName)
	stringToSign := fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s\n%s",
		amzDate, credentialScope, canonicalRequestHashHex)

	signingKey := getSigningKey(s.secretAccessKey, dateStamp, region, serviceName)
	signature := hex.EncodeToString(hmacSHA256(signingKey, stringToSign))

	authHeader := fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		s.accessKeyID, credentialScope, signedHeaders, signature)
	req.Header.Set("Authorization", authHeader)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("erreur suppression R2 status %d", resp.StatusCode)
	}
	return nil
}

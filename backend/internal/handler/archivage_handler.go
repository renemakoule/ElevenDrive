package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"autoschool/backend/internal/middleware"
	"autoschool/backend/internal/service"
)

type ArchivageHandler struct {
	r2Service *service.R2StorageService
}

func NewArchivageHandler(r2Service *service.R2StorageService) *ArchivageHandler {
	return &ArchivageHandler{r2Service: r2Service}
}

// GET /api/archivage/cloud/statut
func (h *ArchivageHandler) GetCloudStatut(w http.ResponseWriter, r *http.Request) {
	configured := h.r2Service.IsConfigured()
	Success(w, http.StatusOK, map[string]interface{}{
		"provider":      "Cloudflare R2",
		"is_configured": configured,
		"bucket_name":   h.r2Service.GetBucketName(),
		"public_url":    h.r2Service.GetPublicURL(),
		"timestamp":     time.Now().UTC().Format(time.RFC3339),
	})
}

type CloudUploadInput struct {
	CandidatID   string `json:"candidat_id"`
	DossierID    string `json:"dossier_id"`
	NumeroDossier string `json:"numero_dossier"`
	TypeDocument string `json:"type_document"`
	Titre        string `json:"titre"`
	NomFichier   string `json:"nom_fichier"`
	FichierBase64 string `json:"fichier_base64"`
	MimeType     string `json:"mime_type"`
}

// POST /api/archivage/cloud/upload
func (h *ArchivageHandler) UploadDocumentToCloud(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		Error(w, http.StatusUnauthorized, "Authentification requise")
		return
	}

	var input CloudUploadInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Payload JSON invalide")
		return
	}

	if input.NomFichier == "" || input.FichierBase64 == "" {
		Error(w, http.StatusBadRequest, "Nom de fichier et données base64 requis")
		return
	}

	// Nettoyage base64 (au cas où il contient data:image/png;base64,...)
	b64Data := input.FichierBase64
	mimeType := input.MimeType
	if idx := strings.Index(b64Data, ";base64,"); idx != -1 {
		prefix := b64Data[:idx]
		if mimeType == "" && strings.HasPrefix(prefix, "data:") {
			mimeType = strings.TrimPrefix(prefix, "data:")
		}
		b64Data = b64Data[idx+8:]
	}

	dataBytes, err := base64.StdEncoding.DecodeString(b64Data)
	if err != nil {
		Error(w, http.StatusBadRequest, "Impossible de décoder les données base64: "+err.Error())
		return
	}

	if mimeType == "" {
		mimeType = http.DetectContentType(dataBytes)
	}

	// Construction de la clé de stockage Cloudflare R2
	annee := time.Now().Format("2006")
	dosRef := input.NumeroDossier
	if dosRef == "" {
		dosRef = input.DossierID
	}
	if dosRef == "" {
		dosRef = "DOS-ARCHIVE"
	}

	cleanFileName := strings.ReplaceAll(filepath.Base(input.NomFichier), " ", "_")
	objectKey := fmt.Sprintf("archives/%s/%s/%s/%d_%s",
		claims.AutoEcoleID, annee, dosRef, time.Now().Unix(), cleanFileName)

	cloudURL, err := h.r2Service.UploadObject(r.Context(), objectKey, dataBytes, mimeType)
	if err != nil {
		Error(w, http.StatusInternalServerError, "Échec téléversement Cloudflare R2: "+err.Error())
		return
	}

	Success(w, http.StatusOK, map[string]interface{}{
		"statut":        "SCELLE_R2",
		"provider":      "Cloudflare R2",
		"bucket":        h.r2Service.GetBucketName(),
		"object_key":    objectKey,
		"public_url":    cloudURL,
		"nom_fichier":   input.NomFichier,
		"taille_octets": len(dataBytes),
		"mime_type":     mimeType,
		"date_scellement": time.Now().UTC().Format(time.RFC3339),
	})
}

// POST /api/archivage/cloud/dossier
func (h *ArchivageHandler) ArchiveDossierToCloud(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		Error(w, http.StatusUnauthorized, "Authentification requise")
		return
	}

	var dossierPayload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&dossierPayload); err != nil {
		Error(w, http.StatusBadRequest, "Payload JSON invalide")
		return
	}

	dosNum, _ := dossierPayload["numero_dossier"].(string)
	if dosNum == "" {
		dosNum = fmt.Sprintf("DOS-%d", time.Now().Unix())
	}

	jsonBytes, err := json.MarshalIndent(dossierPayload, "", "  ")
	if err != nil {
		Error(w, http.StatusInternalServerError, "Erreur sérialisation dossier: "+err.Error())
		return
	}

	annee := time.Now().Format("2006")
	objectKey := fmt.Sprintf("archives/%s/%s/%s/dossier_archive_%d.json",
		claims.AutoEcoleID, annee, dosNum, time.Now().Unix())

	cloudURL, err := h.r2Service.UploadObject(r.Context(), objectKey, jsonBytes, "application/json")
	if err != nil {
		Error(w, http.StatusInternalServerError, "Échec archivage dossier Cloudflare R2: "+err.Error())
		return
	}

	Success(w, http.StatusOK, map[string]interface{}{
		"statut":          "DOSSIER_SCELLE_R2",
		"provider":        "Cloudflare R2",
		"bucket":          h.r2Service.GetBucketName(),
		"object_key":      objectKey,
		"public_url":      cloudURL,
		"taille_octets":   len(jsonBytes),
		"date_archivage":  time.Now().UTC().Format(time.RFC3339),
	})
}

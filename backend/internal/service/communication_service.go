package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"autoschool/backend/internal/models"
	"autoschool/backend/internal/repository"
)

type CommunicationService struct {
	repo         *repository.CommunicationRepository
	notifService *NotificationService
}

func NewCommunicationService(repo *repository.CommunicationRepository) *CommunicationService {
	return &CommunicationService{repo: repo}
}

// SetNotificationService injecte le service de notification temps réel
func (s *CommunicationService) SetNotificationService(notifService *NotificationService) {
	s.notifService = notifService
}

// -------------------------------------------------------------
// ALERTES
// -------------------------------------------------------------

func (s *CommunicationService) GetAllAlerts(ctx context.Context) ([]models.NetworkAlert, error) {
	return s.repo.GetAllAlerts(ctx)
}

func (s *CommunicationService) GetActiveAlertsForSchool(ctx context.Context, schoolID string) ([]models.NetworkAlert, error) {
	return s.repo.GetActiveAlertsForSchool(ctx, schoolID)
}

func (s *CommunicationService) CreateAlert(ctx context.Context, input *models.CreateAlertInput) (*models.NetworkAlert, error) {
	input.Titre = strings.TrimSpace(input.Titre)
	input.Message = strings.TrimSpace(input.Message)

	if input.Titre == "" {
		return nil, errors.New("le titre de l'alerte est requis")
	}
	if input.Message == "" {
		return nil, errors.New("le message de l'alerte est requis")
	}
	if input.Severite != "info" && input.Severite != "warning" && input.Severite != "critical" {
		input.Severite = "warning"
	}
	if input.CibleType != "all" && input.CibleType != "single" {
		input.CibleType = "all"
	}
	if input.CibleType == "single" && (input.CibleSchoolID == nil || strings.TrimSpace(*input.CibleSchoolID) == "") {
		return nil, errors.New("une auto-école cible doit être spécifiée pour une alerte individuelle")
	}

	alert, err := s.repo.CreateAlert(ctx, input)
	if err != nil {
		return nil, err
	}

	if s.notifService != nil && alert != nil {
		notifType := "warning"
		if alert.Severite == "critical" {
			notifType = "critical"
		} else if alert.Severite == "info" {
			notifType = "info"
		}
		var schoolID *string
		if alert.CibleType == "single" {
			schoolID = alert.CibleSchoolID
		}
		_, _ = s.notifService.NotifyAndBroadcast(ctx, &models.CreateNotificationInput{
			AutoEcoleID: schoolID,
			Titre:       "🚨 " + alert.Titre,
			Message:     alert.Message,
			Type:        notifType,
			Module:      "systeme",
		})
	}

	return alert, nil
}

func (s *CommunicationService) ToggleAlertStatus(ctx context.Context, id string) (bool, error) {
	if strings.TrimSpace(id) == "" {
		return false, errors.New("identifiant d'alerte manquant")
	}
	return s.repo.ToggleAlertStatus(ctx, id)
}

func (s *CommunicationService) DeleteAlert(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("identifiant d'alerte manquant")
	}
	return s.repo.DeleteAlert(ctx, id)
}

// -------------------------------------------------------------
// ANNONCES & ÉVÉNEMENTS
// -------------------------------------------------------------

func (s *CommunicationService) GetAllAnnouncements(ctx context.Context) ([]models.NetworkAnnouncement, error) {
	return s.repo.GetAllAnnouncements(ctx)
}

func (s *CommunicationService) GetAnnouncementsForSchool(ctx context.Context, schoolID string) ([]models.NetworkAnnouncement, error) {
	return s.repo.GetAnnouncementsForSchool(ctx, schoolID)
}

func (s *CommunicationService) CreateAnnouncement(ctx context.Context, input *models.CreateAnnouncementInput) (*models.NetworkAnnouncement, error) {
	input.Titre = strings.TrimSpace(input.Titre)
	input.Contenu = strings.TrimSpace(input.Contenu)

	if input.Titre == "" {
		return nil, errors.New("le titre de l'annonce est requis")
	}
	if input.TypeFormat != "image_only" && input.Contenu == "" {
		return nil, errors.New("le texte explicatif de l'annonce est requis")
	}
	if input.TypeFormat == "image_only" && (input.ImageURL == nil || strings.TrimSpace(*input.ImageURL) == "") {
		return nil, errors.New("une image est obligatoire pour le format visuel simple")
	}
	if input.CibleType == "single" && (input.CibleSchoolID == nil || strings.TrimSpace(*input.CibleSchoolID) == "") {
		return nil, errors.New("une auto-école cible doit être sélectionnée")
	}

	return s.repo.CreateAnnouncement(ctx, input)
}

func (s *CommunicationService) DeleteAnnouncement(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("identifiant d'annonce manquant")
	}
	return s.repo.DeleteAnnouncement(ctx, id)
}

// -------------------------------------------------------------
// AVIS & NOTATIONS
// -------------------------------------------------------------

func (s *CommunicationService) GetAllFeedbacks(ctx context.Context, starsFilter int) ([]models.NetworkFeedback, error) {
	return s.repo.GetAllFeedbacks(ctx, starsFilter)
}

func (s *CommunicationService) CreateFeedback(ctx context.Context, input *models.CreateFeedbackInput) (*models.NetworkFeedback, error) {
	input.Auteur = strings.TrimSpace(input.Auteur)
	input.Commentaire = strings.TrimSpace(input.Commentaire)

	if input.Auteur == "" {
		input.Auteur = "Équipe Direction"
	}
	if input.Etoiles < 1 || input.Etoiles > 5 {
		return nil, errors.New("la note doit être comprise entre 1 et 5 étoiles")
	}
	if input.Commentaire == "" {
		return nil, errors.New("votre appréciation ou commentaire est requis")
	}

	return s.repo.CreateFeedback(ctx, input)
}

func (s *CommunicationService) ReplyFeedback(ctx context.Context, input *models.ReplyFeedbackInput) error {
	input.FeedbackID = strings.TrimSpace(input.FeedbackID)
	input.ReponseAdmin = strings.TrimSpace(input.ReponseAdmin)

	if input.FeedbackID == "" {
		return errors.New("identifiant de l'avis requis")
	}
	if input.ReponseAdmin == "" {
		return errors.New("la réponse officielle de l'administration ne peut pas être vide")
	}

	err := s.repo.ReplyFeedback(ctx, input.FeedbackID, input.ReponseAdmin, input.Statut)
	if err != nil {
		return err
	}

	if s.notifService != nil {
		roleDir := "directeur"
		_, _ = s.notifService.NotifyAndBroadcast(ctx, &models.CreateNotificationInput{
			RoleCible: &roleDir,
			Titre:     "⭐ Réponse officielle à votre avis",
			Message:   fmt.Sprintf("L'administration ElevenDrive a répondu : « %s »", input.ReponseAdmin),
			Type:      "success",
			Module:    "systeme",
		})
	}

	return nil
}

func (s *CommunicationService) UpdateFeedbackStatus(ctx context.Context, id, statut string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("identifiant de l'avis requis")
	}
	if statut != "nouveau" && statut != "en_cours" && statut != "resolu" {
		statut = "en_cours"
	}
	return s.repo.UpdateFeedbackStatus(ctx, id, statut)
}

func (s *CommunicationService) DeleteFeedback(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("identifiant de l'avis requis")
	}
	return s.repo.DeleteFeedback(ctx, id)
}

// -------------------------------------------------------------
// SUPPORT TICKETS
// -------------------------------------------------------------

func (s *CommunicationService) CreateSupportTicket(ctx context.Context, input *models.CreateSupportTicketInput) (*models.SupportTicket, error) {
	input.Nom = strings.TrimSpace(input.Nom)
	input.Email = strings.TrimSpace(input.Email)
	input.Sujet = strings.TrimSpace(input.Sujet)
	input.Message = strings.TrimSpace(input.Message)

	if input.Nom == "" {
		return nil, errors.New("le nom est requis")
	}
	if input.Email == "" {
		return nil, errors.New("l'adresse email est requise")
	}
	if input.Message == "" {
		return nil, errors.New("le descriptif de votre demande est requis")
	}

	ticket, err := s.repo.CreateSupportTicket(ctx, input)
	if err != nil {
		return nil, err
	}

	if s.notifService != nil && ticket != nil {
		roleSuperAdmin := "super_admin"
		_, _ = s.notifService.NotifyAndBroadcast(ctx, &models.CreateNotificationInput{
			RoleCible: &roleSuperAdmin,
			Titre:     "🎫 Nouveau ticket support technique",
			Message:   fmt.Sprintf("Demande de %s (%s) : %s", ticket.Nom, ticket.Email, ticket.Sujet),
			Type:      "warning",
			Module:    "super_admin",
			ActionURL: &ticket.ID,
		})
	}

	return ticket, nil
}

func (s *CommunicationService) GetAllSupportTickets(ctx context.Context) ([]models.SupportTicket, error) {
	return s.repo.GetAllSupportTickets(ctx)
}

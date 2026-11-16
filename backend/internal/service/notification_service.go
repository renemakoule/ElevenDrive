package service

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"autoschool/backend/internal/models"
	"autoschool/backend/internal/realtime"
	"autoschool/backend/internal/repository"
)

type NotificationService struct {
	repo       *repository.NotificationRepository
	schoolRepo *repository.SchoolRepository
}

func NewNotificationService(repo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

// SetSchoolRepository injecte le dépôt des auto-écoles pour la surveillance des licences et abonnements
func (s *NotificationService) SetSchoolRepository(schoolRepo *repository.SchoolRepository) {
	s.schoolRepo = schoolRepo
}

// CheckSubscriptionAlert vérifie la période d'essai ou d'abonnement et émet une alerte automatique si expiration proche ou échue
func (s *NotificationService) CheckSubscriptionAlert(ctx context.Context, autoEcoleID string) {
	if s.schoolRepo == nil || strings.TrimSpace(autoEcoleID) == "" {
		return
	}

	school, err := s.schoolRepo.GetByID(ctx, autoEcoleID)
	if err != nil || school == nil || school.DateExpirationLicence == nil {
		return
	}

	hoursUntil := time.Until(*school.DateExpirationLicence).Hours()
	daysRemaining := int(math.Ceil(hoursUntil / 24.0))

	roleDirecteur := models.RoleDirecteur

	// 1. Période d'essai ou licence totalement expirée (jours <= 0)
	if daysRemaining <= 0 {
		hasRecent, _ := s.repo.HasRecentNotification(ctx, autoEcoleID, "abonnement", "critical", 24)
		if !hasRecent {
			titre := "⛔ Période d'essai expirée"
			msg := fmt.Sprintf("Votre période d'essai gratuit de 30 jours pour %s est arrivée à expiration. Votre guichet et vos agences sont temporairement restreints. Cliquez sur 'Bon Plan' pour choisir une formule et réactiver immédiatement votre établissement.", school.NomCommercial)
			_, _ = s.NotifyAndBroadcast(ctx, &models.CreateNotificationInput{
				AutoEcoleID: &autoEcoleID,
				RoleCible:   &roleDirecteur,
				Titre:       titre,
				Message:     msg,
				Type:        "critical",
				Module:      "abonnement",
				Metadata: map[string]interface{}{
					"days_remaining": 0,
					"status":         "EXPIRED",
					"school_id":      school.ID,
				},
			})
		}
		return
	}

	// 2. Alerte d'expiration imminente (J-5 à J-1)
	if daysRemaining <= 5 {
		hasRecent, _ := s.repo.HasRecentNotification(ctx, autoEcoleID, "abonnement", "warning", 24)
		if !hasRecent {
			titre := fmt.Sprintf("⚠️ Période d'essai : J-%d avant expiration", daysRemaining)
			if daysRemaining == 1 {
				titre = "⚠️ Dernier jour de votre période d'essai gratuit !"
			}
			msg := fmt.Sprintf("Votre période d'essai gratuit de 30 jours pour %s se termine dans %d jour(s). Pensez à renouveler votre licence via le bouton 'Bon Plan' pour continuer vos opérations guichet sans coupure.", school.NomCommercial, daysRemaining)
			_, _ = s.NotifyAndBroadcast(ctx, &models.CreateNotificationInput{
				AutoEcoleID: &autoEcoleID,
				RoleCible:   &roleDirecteur,
				Titre:       titre,
				Message:     msg,
				Type:        "warning",
				Module:      "abonnement",
				Metadata: map[string]interface{}{
					"days_remaining": daysRemaining,
					"status":         "EXPIRING_SOON",
					"school_id":      school.ID,
				},
			})
		}
	}
}

// NotifyAndBroadcast enregistre la notification en base PostgreSQL et la diffuse instantanément par WebSocket
func (s *NotificationService) NotifyAndBroadcast(ctx context.Context, input *models.CreateNotificationInput) (*models.SystemNotification, error) {
	input.Titre = strings.TrimSpace(input.Titre)
	input.Message = strings.TrimSpace(input.Message)
	if input.Type == "" {
		input.Type = "info"
	}
	if input.Module == "" {
		input.Module = "systeme"
	}

	// 1. Sauvegarde persistante dans PostgreSQL
	notif, err := s.repo.CreateNotification(ctx, input)
	if err != nil {
		return nil, err
	}

	// 2. Diffusion WebSocket temps réel ciblée
	roleCible := ""
	if input.RoleCible != nil {
		roleCible = strings.ToLower(strings.TrimSpace(*input.RoleCible))
	}

	if roleCible == "super_admin" {
		realtime.BroadcastToSuperAdmin("NOTIFICATION_NEW", notif)
	} else if input.UserID != nil && strings.TrimSpace(*input.UserID) != "" && input.AutoEcoleID != nil {
		realtime.BroadcastToUser(*input.AutoEcoleID, *input.UserID, "NOTIFICATION_NEW", notif)
	} else if input.AutoEcoleID != nil && strings.TrimSpace(*input.AutoEcoleID) != "" {
		if roleCible != "" && roleCible != "all" {
			realtime.BroadcastToRole(*input.AutoEcoleID, roleCible, "NOTIFICATION_NEW", notif)
		} else {
			realtime.Broadcast(*input.AutoEcoleID, "NOTIFICATION_NEW", notif)
		}
	} else {
		realtime.BroadcastGlobal("NOTIFICATION_NEW", notif)
	}

	return notif, nil
}

// LogAudit enregistre une action d'audit en base de données
func (s *NotificationService) LogAudit(ctx context.Context, log *models.SystemAuditLog) error {
	return s.repo.LogAudit(ctx, log)
}

// GetUserNotifications renvoie la liste et le compteur des notifications non lues
func (s *NotificationService) GetUserNotifications(
	ctx context.Context,
	autoEcoleID string,
	userID string,
	role string,
	limit int,
	unreadOnly bool,
) ([]models.SystemNotification, int, error) {
	// Vérification proactive et silencieuse de l'expiration d'abonnement pour l'auto-école
	if autoEcoleID != "" && models.NormalizeRole(role) != models.RoleSuperAdmin {
		s.CheckSubscriptionAlert(ctx, autoEcoleID)
	}

	return s.repo.GetNotificationsForUser(ctx, autoEcoleID, userID, role, limit, unreadOnly)
}

// MarkAsRead passe une notification à l'état lu
func (s *NotificationService) MarkAsRead(ctx context.Context, id string) error {
	return s.repo.MarkAsRead(ctx, id)
}

// MarkAllAsRead passe toutes les notifications de l'utilisateur à l'état lu
func (s *NotificationService) MarkAllAsRead(ctx context.Context, autoEcoleID, userID, role string) error {
	return s.repo.MarkAllAsRead(ctx, autoEcoleID, userID, role)
}

// GetAuditLogs récupère les logs d'audit filtrés
func (s *NotificationService) GetAuditLogs(ctx context.Context, autoEcoleID string, limit int) ([]models.SystemAuditLog, error) {
	return s.repo.GetAuditLogs(ctx, autoEcoleID, limit)
}

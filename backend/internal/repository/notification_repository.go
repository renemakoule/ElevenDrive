package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"autoschool/backend/internal/database"
	"autoschool/backend/internal/models"
)

type NotificationRepository struct {
	db *database.DB
}

func NewNotificationRepository(db *database.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// CreateNotification insère une notification système dans PostgreSQL
func (r *NotificationRepository) CreateNotification(ctx context.Context, input *models.CreateNotificationInput) (*models.SystemNotification, error) {
	var metaJSON []byte
	if input.Metadata != nil {
		metaJSON, _ = json.Marshal(input.Metadata)
	}

	query := `
		INSERT INTO notifications_systeme (
			auto_ecole_id, user_id, role_cible, titre, message, type, module, action_url, metadata, est_lu, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, FALSE, NOW())
		RETURNING id, auto_ecole_id, user_id, role_cible, titre, message, type, module, action_url, metadata, est_lu, created_at
	`

	var n models.SystemNotification
	var autoEcoleID, userID, roleCible, actionURL sql.NullString
	var rawMeta []byte

	err := r.db.QueryRowContext(
		ctx, query,
		input.AutoEcoleID, input.UserID, input.RoleCible, input.Titre, input.Message,
		input.Type, input.Module, input.ActionURL, metaJSON,
	).Scan(
		&n.ID, &autoEcoleID, &userID, &roleCible, &n.Titre, &n.Message,
		&n.Type, &n.Module, &actionURL, &rawMeta, &n.EstLu, &n.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("erreur insertion notification: %w", err)
	}

	if autoEcoleID.Valid {
		v := autoEcoleID.String
		n.AutoEcoleID = &v
	}
	if userID.Valid {
		v := userID.String
		n.UserID = &v
	}
	if roleCible.Valid {
		v := roleCible.String
		n.RoleCible = &v
	}
	if actionURL.Valid {
		v := actionURL.String
		n.ActionURL = &v
	}
	if len(rawMeta) > 0 {
		_ = json.Unmarshal(rawMeta, &n.Metadata)
	}

	return &n, nil
}

// GetNotificationsForUser récupère la liste des notifications destinées à un utilisateur
// en tenant compte de son auto-école, de son identifiant personnel ou de son rôle.
func (r *NotificationRepository) GetNotificationsForUser(
	ctx context.Context,
	autoEcoleID string,
	userID string,
	role string,
	limit int,
	unreadOnly bool,
) ([]models.SystemNotification, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	roleNorm := strings.ToLower(models.NormalizeRole(role))
	var whereClauses []string
	var args []interface{}
	argIdx := 1

	if roleNorm == "super_admin" {
		whereClauses = append(whereClauses, "(role_cible = 'super_admin' OR (auto_ecole_id IS NULL AND user_id IS NULL))")
	} else {
		// Utilisateur d'auto-école : notifications directes de son école ou adressées à son rôle ou globales à son école
		if autoEcoleID != "" {
			sub := fmt.Sprintf("(auto_ecole_id = $%d AND (user_id IS NULL OR user_id = $%d) AND (role_cible IS NULL OR role_cible = '' OR LOWER(role_cible) = 'all' OR LOWER(role_cible) = $%d OR (LOWER(role_cible) IN ('manager','directeur','admin') AND $%d IN ('manager','directeur','admin'))))",
				argIdx, argIdx+1, argIdx+2, argIdx+2)
			whereClauses = append(whereClauses, sub)
			args = append(args, autoEcoleID, userID, roleNorm)
			argIdx += 3
		} else if userID != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("user_id = $%d", argIdx))
			args = append(args, userID)
			argIdx++
		}
	}

	if unreadOnly {
		whereClauses = append(whereClauses, "est_lu = FALSE")
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// 1. Compter les non lus
	unreadCountQuery := `SELECT COUNT(*) FROM notifications_systeme`
	if len(whereClauses) > 0 {
		// reconstruire avec est_lu = FALSE
		var unreadClauses []string
		for _, c := range whereClauses {
			if c != "est_lu = FALSE" {
				unreadClauses = append(unreadClauses, c)
			}
		}
		unreadClauses = append(unreadClauses, "est_lu = FALSE")
		unreadCountQuery += " WHERE " + strings.Join(unreadClauses, " AND ")
	} else {
		unreadCountQuery += " WHERE est_lu = FALSE"
	}

	var unreadCount int
	_ = r.db.QueryRowContext(ctx, unreadCountQuery, args...).Scan(&unreadCount)

	// 2. Sélectionner la liste
	selectQuery := fmt.Sprintf(`
		SELECT id, auto_ecole_id, user_id, role_cible, titre, message, type, module, action_url, metadata, est_lu, created_at
		FROM notifications_systeme
		%s
		ORDER BY created_at DESC
		LIMIT $%d
	`, whereSQL, argIdx)
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("erreur récupération notifications: %w", err)
	}
	defer rows.Close()

	var list []models.SystemNotification
	for rows.Next() {
		var n models.SystemNotification
		var schoolID, uID, rCible, actionURL sql.NullString
		var rawMeta []byte

		if err := rows.Scan(
			&n.ID, &schoolID, &uID, &rCible, &n.Titre, &n.Message,
			&n.Type, &n.Module, &actionURL, &rawMeta, &n.EstLu, &n.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("erreur scan notification: %w", err)
		}

		if schoolID.Valid {
			v := schoolID.String
			n.AutoEcoleID = &v
		}
		if uID.Valid {
			v := uID.String
			n.UserID = &v
		}
		if rCible.Valid {
			v := rCible.String
			n.RoleCible = &v
		}
		if actionURL.Valid {
			v := actionURL.String
			n.ActionURL = &v
		}
		if len(rawMeta) > 0 {
			_ = json.Unmarshal(rawMeta, &n.Metadata)
		}

		list = append(list, n)
	}
	if list == nil {
		list = []models.SystemNotification{}
	}

	return list, unreadCount, nil
}

// MarkAsRead passe une notification à l'état lu
func (r *NotificationRepository) MarkAsRead(ctx context.Context, id string) error {
	query := `UPDATE notifications_systeme SET est_lu = TRUE WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("erreur mise à jour notification lue: %w", err)
	}
	return nil
}

// MarkAllAsRead passe toutes les notifications d'un contexte à l'état lu
func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, autoEcoleID, userID, role string) error {
	roleNorm := strings.ToLower(models.NormalizeRole(role))
	var query string
	var args []interface{}

	if roleNorm == "super_admin" {
		query = `UPDATE notifications_systeme SET est_lu = TRUE WHERE est_lu = FALSE AND (role_cible = 'super_admin' OR (auto_ecole_id IS NULL AND user_id IS NULL))`
	} else if autoEcoleID != "" {
		query = `UPDATE notifications_systeme SET est_lu = TRUE WHERE est_lu = FALSE AND auto_ecole_id = $1 AND (user_id IS NULL OR user_id = $2) AND (role_cible IS NULL OR role_cible = '' OR LOWER(role_cible) = 'all' OR LOWER(role_cible) = $3 OR (LOWER(role_cible) IN ('manager','directeur','admin') AND $3 IN ('manager','directeur','admin')))`
		args = append(args, autoEcoleID, userID, roleNorm)
	} else if userID != "" {
		query = `UPDATE notifications_systeme SET est_lu = TRUE WHERE est_lu = FALSE AND user_id = $1`
		args = append(args, userID)
	} else {
		return nil
	}

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("erreur marquage global lu: %w", err)
	}
	return nil
}

// HasRecentNotification vérifie si une notification de même type et module existe déjà récemment pour éviter les doublons
func (r *NotificationRepository) HasRecentNotification(ctx context.Context, autoEcoleID, module, notifType string, sinceHours int) (bool, error) {
	if sinceHours <= 0 {
		sinceHours = 24
	}
	query := `
		SELECT COUNT(*) FROM notifications_systeme 
		WHERE auto_ecole_id = $1 AND module = $2 AND type = $3 
		  AND created_at >= NOW() - ($4 || ' hours')::INTERVAL
	`
	var count int
	err := r.db.QueryRowContext(ctx, query, autoEcoleID, module, notifType, fmt.Sprintf("%d", sinceHours)).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// LogAudit enregistre une trace d'action dans audit_logs_super_admin
func (r *NotificationRepository) LogAudit(ctx context.Context, log *models.SystemAuditLog) error {
	query := `
		INSERT INTO audit_logs_super_admin (
			action, cible_type, cible_id, details, ip_address, auto_ecole_id, user_id, auteur_nom, role, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
	`
	_, err := r.db.ExecContext(
		ctx, query,
		log.Action, log.CibleType, log.CibleID, log.Details, log.IPAddress,
		log.AutoEcoleID, log.UserID, log.AuteurNom, log.Role,
	)
	if err != nil {
		return fmt.Errorf("erreur insertion journal audit: %w", err)
	}
	return nil
}

// GetAuditLogs récupère les logs d'audit récents
func (r *NotificationRepository) GetAuditLogs(ctx context.Context, autoEcoleID string, limit int) ([]models.SystemAuditLog, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}

	query := `
		SELECT id, action, cible_type, cible_id, details, COALESCE(ip_address, ''),
		       auto_ecole_id, user_id, auteur_nom, role, created_at
		FROM audit_logs_super_admin
	`
	var args []interface{}
	if autoEcoleID != "" {
		query += ` WHERE auto_ecole_id = $1`
		args = append(args, autoEcoleID)
		query += ` ORDER BY created_at DESC LIMIT $2`
		args = append(args, limit)
	} else {
		query += ` ORDER BY created_at DESC LIMIT $1`
		args = append(args, limit)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("erreur lecture logs audit: %w", err)
	}
	defer rows.Close()

	var logs []models.SystemAuditLog
	for rows.Next() {
		var l models.SystemAuditLog
		var schoolID, uID, autNom, rName sql.NullString
		if err := rows.Scan(
			&l.ID, &l.Action, &l.CibleType, &l.CibleID, &l.Details, &l.IPAddress,
			&schoolID, &uID, &autNom, &rName, &l.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("erreur scan log audit: %w", err)
		}
		if schoolID.Valid {
			v := schoolID.String
			l.AutoEcoleID = &v
		}
		if uID.Valid {
			v := uID.String
			l.UserID = &v
		}
		if autNom.Valid {
			v := autNom.String
			l.AuteurNom = &v
		}
		if rName.Valid {
			v := rName.String
			l.Role = &v
		}
		logs = append(logs, l)
	}
	if logs == nil {
		logs = []models.SystemAuditLog{}
	}
	return logs, nil
}

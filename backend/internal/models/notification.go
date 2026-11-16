package models

import "time"

// SystemNotification représente une notification système persistante dans PostgreSQL
type SystemNotification struct {
	ID          string                 `json:"id"`
	AutoEcoleID *string                `json:"auto_ecole_id,omitempty"`
	UserID      *string                `json:"user_id,omitempty"`
	RoleCible   *string                `json:"role_cible,omitempty"`
	Titre       string                 `json:"titre"`
	Message     string                 `json:"message"`
	Type        string                 `json:"type"`   // "info" | "warning" | "critical" | "success"
	Module      string                 `json:"module"` // "guichet" | "caisse" | "conduite" | "examens" | "flotte" | "super_admin" | "systeme"
	ActionURL   *string                `json:"action_url,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	EstLu       bool                   `json:"est_lu"`
	CreatedAt   time.Time              `json:"created_at"`
}

// CreateNotificationInput charge utile pour persister et diffuser une notification
type CreateNotificationInput struct {
	AutoEcoleID *string                `json:"auto_ecole_id,omitempty"`
	UserID      *string                `json:"user_id,omitempty"`
	RoleCible   *string                `json:"role_cible,omitempty"`
	Titre       string                 `json:"titre"`
	Message     string                 `json:"message"`
	Type        string                 `json:"type"`
	Module      string                 `json:"module"`
	ActionURL   *string                `json:"action_url,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// SystemAuditLog représente une entrée de traçabilité d'audit dans PostgreSQL
type SystemAuditLog struct {
	ID          string    `json:"id"`
	Action      string    `json:"action"`
	CibleType   string    `json:"cible_type"`
	CibleID     string    `json:"cible_id"`
	Details     string    `json:"details"`
	IPAddress   string    `json:"ip_address"`
	AutoEcoleID *string   `json:"auto_ecole_id,omitempty"`
	UserID      *string   `json:"user_id,omitempty"`
	AuteurNom   *string   `json:"auteur_nom,omitempty"`
	Role        *string   `json:"role,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

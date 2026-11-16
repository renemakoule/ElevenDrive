package models

import "time"

// NetworkAlert représente une alerte diffusée à l'ensemble du réseau ou à une auto-école spécifique.
type NetworkAlert struct {
	ID              string     `json:"id"`
	Titre           string     `json:"titre"`
	Message         string     `json:"message"`
	Severite        string     `json:"severite"` // "info" | "warning" | "critical"
	CibleType       string     `json:"cible_type"` // "all" | "single"
	CibleSchoolID   *string    `json:"cible_school_id,omitempty"`
	CibleSchoolName string     `json:"cible_school_name,omitempty"`
	Bloquante       bool       `json:"bloquante"`
	EstActive       bool       `json:"est_active"`
	DateEnvoi       time.Time  `json:"date_envoi"`
	CreatedAt       time.Time  `json:"created_at"`
}

// CreateAlertInput représente la charge utile pour diffuser une alerte.
type CreateAlertInput struct {
	Titre         string  `json:"titre"`
	Message       string  `json:"message"`
	Severite      string  `json:"severite"`
	CibleType     string  `json:"cible_type"`
	CibleSchoolID *string `json:"cible_school_id,omitempty"`
	Bloquante     bool    `json:"bloquante"`
}

// NetworkAnnouncement représente une annonce ou un message événementiel du réseau.
type NetworkAnnouncement struct {
	ID              string    `json:"id"`
	Titre           string    `json:"titre"`
	Contenu         string    `json:"contenu"`
	TypeFormat      string    `json:"type_format"` // "text_only" | "text_image" | "image_only"
	ImageURL        *string   `json:"image_url,omitempty"`
	Categorie       string    `json:"categorie"`   // "evenement" | "mise_a_jour" | "formation" | "general"
	CibleType       string    `json:"cible_type"`  // "all" | "single"
	CibleSchoolID   *string   `json:"cible_school_id,omitempty"`
	CibleSchoolName string    `json:"cible_school_name,omitempty"`
	DateEvenement   *string   `json:"date_evenement,omitempty"`
	LienAction      *string   `json:"lien_action,omitempty"`
	DatePublication time.Time `json:"date_publication"`
	CreatedAt       time.Time `json:"created_at"`
}

// CreateAnnouncementInput représente la charge utile pour publier une annonce réseau.
type CreateAnnouncementInput struct {
	Titre         string  `json:"titre"`
	Contenu       string  `json:"contenu"`
	TypeFormat    string  `json:"type_format"`
	ImageURL      *string `json:"image_url,omitempty"`
	Categorie     string  `json:"categorie"`
	CibleType     string  `json:"cible_type"`
	CibleSchoolID *string `json:"cible_school_id,omitempty"`
	DateEvenement *string `json:"date_evenement,omitempty"`
	LienAction    *string `json:"lien_action,omitempty"`
}

// NetworkFeedback représente un avis ou retour d'expérience d'une auto-école.
type NetworkFeedback struct {
	ID           string    `json:"id"`
	AutoEcoleID  *string   `json:"auto_ecole_id,omitempty"`
	AutoEcoleNom string    `json:"auto_ecole_nom"`
	Auteur       string    `json:"auteur"`
	Role         string    `json:"role"`
	Etoiles      int       `json:"etoiles"` // 1 à 5
	Categorie    string    `json:"categorie"` // "ergonomie" | "guichet" | "pwa_mobile" | "examens" | "support" | "general"
	Commentaire  string    `json:"commentaire"`
	Statut       string    `json:"statut"` // "nouveau" | "en_cours" | "resolu"
	ReponseAdmin *string   `json:"reponse_admin,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// CreateFeedbackInput pour enregistrer un nouvel avis.
type CreateFeedbackInput struct {
	AutoEcoleID *string `json:"auto_ecole_id,omitempty"`
	Auteur      string  `json:"auteur"`
	Role        string  `json:"role"`
	Etoiles     int     `json:"etoiles"`
	Categorie   string  `json:"categorie"`
	Commentaire string  `json:"commentaire"`
}

// ReplyFeedbackInput pour envoyer une réponse officielle du Super-Admin.
type ReplyFeedbackInput struct {
	FeedbackID   string `json:"feedback_id"`
	ReponseAdmin string `json:"reponse_admin"`
	Statut       string `json:"statut,omitempty"`
}

// UpdateFeedbackStatusInput pour modifier le statut d'un avis.
type UpdateFeedbackStatusInput struct {
	FeedbackID string `json:"feedback_id"`
	Statut     string `json:"statut"` // "nouveau" | "en_cours" | "resolu"
}

// SupportTicket représente une demande de support envoyée au Super-Admin.
type SupportTicket struct {
	ID           string    `json:"id"`
	AutoEcoleID  *string   `json:"auto_ecole_id,omitempty"`
	AutoEcoleNom string    `json:"auto_ecole_nom,omitempty"`
	Nom          string    `json:"nom"`
	Email        string    `json:"email"`
	Sujet        string    `json:"sujet"`
	Priorite     string    `json:"priorite"`
	Message      string    `json:"message"`
	Statut       string    `json:"statut"` // "ouvert" | "en_cours" | "resolu"
	CreatedAt    time.Time `json:"created_at"`
}

// CreateSupportTicketInput pour soumettre un ticket.
type CreateSupportTicketInput struct {
	AutoEcoleID *string `json:"auto_ecole_id,omitempty"`
	Nom         string  `json:"nom"`
	Email       string  `json:"email"`
	Sujet       string  `json:"sujet"`
	Priorite    string  `json:"priorite"`
	Message     string  `json:"message"`
}


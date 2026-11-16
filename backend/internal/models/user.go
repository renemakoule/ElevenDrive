package models

import (
	"strings"
	"time"
)

// Rôles utilisateurs du système (conformes à la contrainte CHECK de PostgreSQL)
const (
	RoleSuperAdmin = "SUPER_ADMIN"
	RoleDirecteur  = "DIRECTEUR"
	RoleSecretaire = "SECRETAIRE"
	RoleCaissier   = "CAISSIER"
	RoleMoniteur   = "MONITEUR"
	RoleEleve      = "ELEVE"
)

// Statuts des comptes
const (
	StatusActif    = "ACTIF"
	StatusSuspendu = "SUSPENDU"
	StatusArchive  = "ARCHIVE"
)

// User représente un compte utilisateur dans la table 'utilisateurs'.
type User struct {
	ID             string    `json:"id"`
	AutoEcoleID    string    `json:"auto_ecole_id"`
	AgenceID       *string   `json:"agence_id,omitempty"`
	Nom            string    `json:"nom"`
	Prenoms        string    `json:"prenoms"`
	Telephone      string    `json:"telephone"`
	Email          *string   `json:"email,omitempty"`
	MotDePasseHash *string   `json:"-"`
	CodePinManager *string   `json:"-"`
	Role           string    `json:"role"`
	Statut         string    `json:"statut"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// UserResponse représente les informations publiques d'un utilisateur sans informations sensibles.
type UserResponse struct {
	ID          string  `json:"id"`
	AutoEcoleID string  `json:"auto_ecole_id"`
	AgenceID    *string `json:"agence_id,omitempty"`
	Nom         string  `json:"nom"`
	Prenoms     string  `json:"prenoms"`
	Telephone   string  `json:"telephone"`
	Email       *string `json:"email,omitempty"`
	Role        string  `json:"role"`
	Statut      string  `json:"statut"`
}

// UserTeamMember représente un collaborateur avec le nom de son agence de rattachement.
type UserTeamMember struct {
	ID          string  `json:"id"`
	AutoEcoleID string  `json:"auto_ecole_id"`
	AgenceID    *string `json:"agence_id,omitempty"`
	AgenceNom   string  `json:"agence_nom"`
	Nom         string  `json:"nom"`
	Prenoms     string  `json:"prenoms"`
	Telephone   string  `json:"telephone"`
	Email       *string `json:"email,omitempty"`
	Role        string  `json:"role"`
	Statut      string  `json:"statut"`
	CreatedAt   string  `json:"created_at"`
}

// ToResponse convertit une entité User en DTO public sécurisé.
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:          u.ID,
		AutoEcoleID: u.AutoEcoleID,
		AgenceID:    u.AgenceID,
		Nom:         u.Nom,
		Prenoms:     u.Prenoms,
		Telephone:   u.Telephone,
		Email:       u.Email,
		Role:        u.Role,
		Statut:      u.Statut,
	}
}

// RegisterUserInput charge la création d'un collaborateur au sein d'une auto-école.
type RegisterUserInput struct {
	AutoEcoleID    string  `json:"auto_ecole_id"`
	AgenceID       *string `json:"agence_id"`
	Nom            string  `json:"nom"`
	Prenoms        string  `json:"prenoms"`
	Telephone      string  `json:"telephone"`
	Email          *string `json:"email"`
	MotDePasse     string  `json:"mot_de_passe"`
	CodePinManager *string `json:"code_pin_manager"`
	Role           string  `json:"role"`
}

// LoginInput données d'identification.
type LoginInput struct {
	// Identifiant : Numéro de téléphone ou Email
	Identifiant string `json:"identifiant"`
	MotDePasse  string `json:"mot_de_passe"`
	// Optionnel : Slug ou Code de l'auto-école si l'utilisateur est rattaché à une école spécifique
	SchoolSlug string `json:"school_slug,omitempty"`
	// Identification du poste pour le contrôle de quota de machines (Desktop Tauri)
	MachineFingerprint string `json:"machine_fingerprint,omitempty"`
	MachineName        string `json:"machine_name,omitempty"`
	OS                 string `json:"os,omitempty"`
}

// VerifyPinInput vérification rapide du code PIN manager pour autorisations de guichet.
type VerifyPinInput struct {
	AutoEcoleID string `json:"auto_ecole_id"`
	CodePin     string `json:"code_pin"`
}

// UpdatePasswordInput changement de mot de passe.
type UpdatePasswordInput struct {
	AncienMotDePasse  string `json:"ancien_mot_de_passe"`
	NouveauMotDePasse string `json:"nouveau_mot_de_passe"`
}

// UpdateUserInput données de modification d'un collaborateur.
type UpdateUserInput struct {
	Nom            string  `json:"nom"`
	Prenoms        string  `json:"prenoms"`
	Telephone      string  `json:"telephone"`
	Email          *string `json:"email"`
	Role           string  `json:"role"`
	AgenceID       *string `json:"agence_id"`
	Statut         string  `json:"statut"`
	MotDePasse     string  `json:"mot_de_passe,omitempty"`
	CodePinManager *string `json:"code_pin_manager,omitempty"`
}

// NormalizeRole normalise les alias de rôles (ex: "manager"/"gerant" -> "DIRECTEUR", "caissière" -> "CAISSIER").
func NormalizeRole(role string) string {
	switch strings.ToUpper(strings.TrimSpace(role)) {
	case RoleSuperAdmin, "SUPERADMIN", "SUPER-ADMIN", "ROOT":
		return RoleSuperAdmin
	case RoleDirecteur, "MANAGER", "GERANT", "GÉRANT", "ADMIN", "ADMINISTRATEUR":
		return RoleDirecteur
	case RoleSecretaire, "SECRÉTAIRE":
		return RoleSecretaire
	case RoleCaissier, "CAISSIERE", "CAISSIÈRE":
		return RoleCaissier
	case RoleMoniteur, "INSTRUCTEUR":
		return RoleMoniteur
	case RoleEleve, "ÉLÈVE", "CANDIDAT", "APPRENANT":
		return RoleEleve
	default:
		return strings.ToUpper(strings.TrimSpace(role))
	}
}

// IsValidRole vérifie si le rôle fait partie des valeurs permises.
func IsValidRole(role string) bool {
	switch NormalizeRole(role) {
	case RoleSuperAdmin, RoleDirecteur, RoleSecretaire, RoleCaissier, RoleMoniteur, RoleEleve:
		return true
	default:
		return false
	}
}

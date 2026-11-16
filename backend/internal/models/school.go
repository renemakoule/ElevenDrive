package models

import (
	"time"
)

// AutoEcole représente l'entreprise titulaire de l'agrément ministériel (Tenant principal).
type AutoEcole struct {
	ID                    string     `json:"id"`
	CodeUnique            string     `json:"code_unique"`
	NomCommercial         string     `json:"nom_commercial"`
	NumeroAgrement        string     `json:"numero_agrement"`
	Telephone             string     `json:"telephone"`
	Email                 *string    `json:"email,omitempty"`
	Devise                string     `json:"devise"`
	Pays                  string     `json:"pays"`
	Ville                 string     `json:"ville,omitempty"`
	LogoURL               *string    `json:"logo_url,omitempty"`
	EstActif              bool       `json:"est_actif"`
	MaxMachines           int        `json:"max_machines"`
	MachinesActives       int        `json:"machines_actives"`
	DateExpirationLicence *time.Time `json:"date_expiration_licence,omitempty"`
	TypeLicence           string     `json:"type_licence,omitempty"`
	ModulesActifs         string     `json:"modules_actifs,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

// Agence représente une succursale physique d'exploitation de l'auto-école.
type Agence struct {
	ID          string    `json:"id"`
	AutoEcoleID string    `json:"auto_ecole_id"`
	CodeAgence  string    `json:"code_agence"`
	Nom         string    `json:"nom"`
	Adresse     string    `json:"adresse"`
	Telephone   string    `json:"telephone"`
	Ville       string    `json:"ville"`
	Commune     string    `json:"commune"`
	EstSiege    bool      `json:"est_siege"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RegisterSchoolInput structure requise pour créer une nouvelle auto-école avec son agence siège et son administrateur.
type RegisterSchoolInput struct {
	// Auto-École
	NomCommercial  string  `json:"nom_commercial"`
	CodeUnique     string  `json:"code_unique"` // Optionnel : si vide, généré automatiquement depuis le nom
	NumeroAgrement string  `json:"numero_agrement"`
	TelephoneEcole string  `json:"telephone_ecole"`
	EmailEcole     *string `json:"email_ecole"`
	Pays           string  `json:"pays"`
	Devise         string  `json:"devise"`

	// Agence Siège
	NomAgence       string `json:"nom_agence"`
	AdresseAgence   string `json:"adresse_agence"`
	TelephoneAgence string `json:"telephone_agence"`
	VilleAgence     string `json:"ville_agence"`
	CommuneAgence   string `json:"commune_agence"`

	// Compte Directeur Initial
	NomDirecteur       string  `json:"nom_directeur"`
	PrenomsDirecteur   string  `json:"prenoms_directeur"`
	TelephoneDirecteur string  `json:"telephone_directeur"`
	EmailDirecteur     *string `json:"email_directeur"`
	MotDePasse         string  `json:"mot_de_passe"`
	CodePinManager     *string `json:"code_pin_manager"`
}

// SchoolPublicInfo informations publiques accessibles aux moniteurs et candidats.
type SchoolPublicInfo struct {
	ID             string    `json:"id"`
	CodeUnique     string    `json:"code_unique"`
	NomCommercial  string    `json:"nom_commercial"`
	NumeroAgrement string    `json:"numero_agrement"`
	Telephone      string    `json:"telephone"`
	Email          *string   `json:"email,omitempty"`
	Devise         string    `json:"devise"`
	Pays           string    `json:"pays"`
	LogoURL        *string   `json:"logo_url,omitempty"`
	Agences        []*Agence `json:"agences"`
}

// CreateAgenceInput données requises pour créer une succursale d'exploitation.
type CreateAgenceInput struct {
	Nom        string `json:"nom"`
	CodeAgence string `json:"code_agence"`
	Adresse    string `json:"adresse"`
	Ville      string `json:"ville"`
	Commune    string `json:"commune"`
	Telephone  string `json:"telephone"`
	EstSiege   bool   `json:"est_siege"`
}

// UpdateAgenceInput données modifiables d'une succursale.
type UpdateAgenceInput struct {
	Nom        string `json:"nom"`
	CodeAgence string `json:"code_agence"`
	Adresse    string `json:"adresse"`
	Ville      string `json:"ville"`
	Commune    string `json:"commune"`
	Telephone  string `json:"telephone"`
	EstSiege   bool   `json:"est_siege"`
}


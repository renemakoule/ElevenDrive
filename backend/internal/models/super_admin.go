package models

import "time"

// SuperAdminStats regroupe les métriques globales de l'infrastructure
type SuperAdminStats struct {
	TotalSchools        int     `json:"total_schools"`
	ActiveSchools       int     `json:"active_schools"`
	TotalUsers          int     `json:"total_users"`
	TotalCandidats      int     `json:"total_candidats"`
	TotalVolumePaiement float64 `json:"total_volume_paiement"`
	CloudStatus         string  `json:"cloud_status"`
	ClusterLatencyMs    int64   `json:"cluster_latency_ms"`
}

// SchoolItemSummary représente une ligne d'école dans la console d'administration
type SchoolItemSummary struct {
	ID                    string     `json:"id"`
	CodeUnique            string     `json:"code_unique"`
	NomCommercial         string     `json:"nom_commercial"`
	NumeroAgrement        string     `json:"numero_agrement"`
	Telephone             string     `json:"telephone"`
	Email                 *string    `json:"email"`
	Devise                string     `json:"devise"`
	Pays                  string     `json:"pays"`
	Ville                 string     `json:"ville"`
	EstActif              bool       `json:"est_actif"`
	MaxMachines           int        `json:"max_machines"`
	MachinesActives       int        `json:"machines_actives"`
	DateExpirationLicence *time.Time `json:"date_expiration_licence"`
	TypeLicence           string     `json:"type_licence"`
	ModulesActifs         string     `json:"modules_actifs"`
	JoursRestants         int        `json:"jours_restants"`
	StatutLicence         string     `json:"statut_licence"`
	AgencesCount          int        `json:"agences_count"`
	UsersCount            int        `json:"users_count"`
	CandidatsCount        int        `json:"candidats_count"`
	TotalVolume           float64    `json:"total_volume"`
	Adresse               string     `json:"adresse,omitempty"`
	Agences               []*Agence  `json:"agences,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

// SuperAdminLoginInput pour l'authentification console
type SuperAdminLoginInput struct {
	Identifiant string `json:"identifiant"`
	MotDePasse  string `json:"mot_de_passe"`
	SecurityKey string `json:"security_key,omitempty"`
}

// UpdateSchoolAdminInput pour la modification complète d'une auto-école par le Super-Admin
type UpdateSchoolAdminInput struct {
	SchoolID              string     `json:"school_id"`
	NomCommercial         string     `json:"nom_commercial"`
	CodeUnique            string     `json:"code_unique"`
	NumeroAgrement        string     `json:"numero_agrement"`
	Telephone             string     `json:"telephone"`
	Email                 *string    `json:"email"`
	Ville                 string     `json:"ville"`
	Pays                  string     `json:"pays"`
	Devise                string     `json:"devise"`
	EstActif              bool       `json:"est_actif"`
	MaxMachines           int        `json:"max_machines"`
	DateExpirationLicence *time.Time `json:"date_expiration_licence"`
	TypeLicence           string     `json:"type_licence"`
	ModulesActifs         string     `json:"modules_actifs"`
}

// UpdateSchoolLicenceInput pour définir le nombre de machines et la durée d'utilisation
type UpdateSchoolLicenceInput struct {
	SchoolID              string     `json:"school_id"`
	MaxMachines           int        `json:"max_machines"`
	TypeLicence           string     `json:"type_licence"`   // "TRIAL", "MENSUEL", "TRIMESTRIEL", "SEMESTRIEL", "ANNUEL", "3 ANS", "CUSTOM"
	DureeValeur           int        `json:"duree_valeur"`   // Nombre de jours / mois / années
	DureeUnite            string     `json:"duree_unite"`    // "JOURS", "MOIS", "ANNEES"
	Mode                  string     `json:"mode,omitempty"` // "REPLACE" (depuis aujourd'hui), "EXTEND" (cumuler sur l'échéance active), "KEEP_DATE" (conserver la date actuelle)
	DateExpirationLicence *time.Time `json:"date_expiration_licence,omitempty"`
}

// UpdateSchoolModulesInput pour activer/désactiver les modules de droits d'accès
type UpdateSchoolModulesInput struct {
	SchoolID      string   `json:"school_id"`
	ModulesActifs []string `json:"modules_actifs"`
}

// DeleteSchoolInput pour la suppression sécurisée d'une auto-école
type DeleteSchoolInput struct {
	SchoolID    string `json:"school_id"`
	ConfirmName string `json:"confirm_name"`
}

// AuditLogEntry pour la traçabilité des actions Super-Admin
type AuditLogEntry struct {
	ID        string    `json:"id"`
	Action    string    `json:"action"`
	CibleType string    `json:"cible_type"`
	CibleID   string    `json:"cible_id"`
	Details   string    `json:"details"`
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}

// GlobalUserSummary représente un utilisateur de la plateforme pour la supervision Super-Admin
type GlobalUserSummary struct {
	ID           string    `json:"id"`
	AutoEcoleID  string    `json:"auto_ecole_id"`
	AutoEcoleNom string    `json:"auto_ecole_nom"`
	Nom          string    `json:"nom"`
	Prenoms      string    `json:"prenoms"`
	Telephone    string    `json:"telephone"`
	Email        *string   `json:"email"`
	Role         string    `json:"role"`
	Statut       string    `json:"statut"`
	CreatedAt    time.Time `json:"created_at"`
}

// GlobalCandidatSummary représente un candidat inscrit pour la supervision Super-Admin
type GlobalCandidatSummary struct {
	ID                string    `json:"id"`
	AutoEcoleID       string    `json:"auto_ecole_id"`
	AutoEcoleNom      string    `json:"auto_ecole_nom"`
	Nom               string    `json:"nom"`
	Prenoms           string    `json:"prenoms"`
	TelephoneWhatsApp string    `json:"telephone_whatsapp"`
	Genre             string    `json:"genre"`
	NumeroDossier     string    `json:"numero_dossier"`
	StatutDossier     string    `json:"statut_dossier"`
	MontantTotal      float64   `json:"montant_total"`
	SoldeRestant      float64   `json:"solde_restant"`
	CreatedAt         time.Time `json:"created_at"`
}

// ToggleSchoolStatusInput pour activer/suspendre un établissement
type ToggleSchoolStatusInput struct {
	SchoolID string `json:"school_id"`
	EstActif bool   `json:"est_actif"`
}

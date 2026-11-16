package models

import (
	"time"
)

// ==============================================================================
// 1. CATALOGUE FORMULES DE PERMIS
// ==============================================================================

type FormulePermis struct {
	ID               string    `json:"id"`
	AutoEcoleID      string    `json:"auto_ecole_id"`
	CategoriePermis  string    `json:"categorie_permis"`
	Intitule         string    `json:"intitule"`
	Description      *string   `json:"description,omitempty"`
	HeuresTheoriques int       `json:"heures_theoriques"`
	HeuresPratiques  int       `json:"heures_pratiques"`
	PrixTotal        float64   `json:"prix_total"`
	AcompteMinimum   float64   `json:"acompte_minimum"`
	EstActif         bool      `json:"est_actif"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CreerFormuleInput struct {
	CategoriePermis  string  `json:"categorie_permis"`
	Intitule         string  `json:"intitule"`
	Description      *string `json:"description,omitempty"`
	HeuresTheoriques int     `json:"heures_theoriques"`
	HeuresPratiques  int     `json:"heures_pratiques"`
	PrixTotal        float64 `json:"prix_total"`
	AcompteMinimum   float64 `json:"acompte_minimum"`
}

type ModifierPrixFormuleInput struct {
	FormuleID      string   `json:"formule_id"`
	NouveauPrix    float64  `json:"nouveau_prix"`
	AcompteMinimum *float64 `json:"acompte_minimum,omitempty"`
	Motif          *string  `json:"motif,omitempty"`
}

// ==============================================================================
// 2. CANDIDATS & DOSSIERS D'INSCRIPTION
// ==============================================================================

type Candidat struct {
	ID                   string     `json:"id"`
	AutoEcoleID          string     `json:"auto_ecole_id"`
	AgenceID             string     `json:"agence_id"`
	Nom                  string     `json:"nom"`
	Prenoms              string     `json:"prenoms"`
	Genre                string     `json:"genre"`
	DateNaissance        string     `json:"date_naissance"`
	LieuNaissance        string     `json:"lieu_naissance"`
	Nationalite          string     `json:"nationalite"`
	Profession           *string    `json:"profession,omitempty"`
	QuartierResidence    *string    `json:"quartier_residence,omitempty"`
	TelephoneWhatsApp    string     `json:"telephone_whatsapp"`
	TelephoneUrgence     *string    `json:"telephone_urgence,omitempty"`
	GroupeSanguin        *string    `json:"groupe_sanguin,omitempty"`
	PhotoURL             *string    `json:"photo_url,omitempty"`
	TypePieceID          string     `json:"type_piece_id"`
	NumeroPieceID        string     `json:"numero_piece_id"`
	DateExpirationPiece  *string    `json:"date_expiration_piece,omitempty"`
	VisiteMedicaleStatut string     `json:"visite_medicale_statut"`
	DateVisiteMedicale   *string    `json:"date_visite_medicale,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	SyncedAt             *time.Time `json:"synced_at,omitempty"`

	// Données du dossier actif associées (injectées lors des recherches / listes)
	NumeroDossier   *string  `json:"numero_dossier,omitempty"`
	FormuleNom      *string  `json:"formule_nom,omitempty"`
	CategoriePermis *string  `json:"categorie_permis,omitempty"`
	SoldeRestant    *float64 `json:"solde_restant,omitempty"`
	StatutDossier   *string  `json:"statut_dossier,omitempty"`
}

type DossierInscription struct {
	ID                         string     `json:"id"`
	AutoEcoleID                string     `json:"auto_ecole_id"`
	AgenceID                   string     `json:"agence_id"`
	CandidatID                 string     `json:"candidat_id"`
	FormuleID                  string     `json:"formule_id"`
	NumeroDossier              string     `json:"numero_dossier"`
	DateInscription            string     `json:"date_inscription"`
	MontantTotal               float64    `json:"montant_total"`
	RemiseAccordee             float64    `json:"remise_accordee"`
	SoldeRestant               float64    `json:"solde_restant"`
	StatutDossier              string     `json:"statut_dossier"`
	HeuresCodeEffectuees       int        `json:"heures_code_effectuees"`
	HeuresConduiteEffectuees   int        `json:"heures_conduite_effectuees"`
	DerogationConduiteAccordee bool       `json:"derogation_conduite_accordee"`
	CreatedAt                  time.Time  `json:"created_at"`
	UpdatedAt                  time.Time  `json:"updated_at"`
	SyncedAt                   *time.Time `json:"synced_at,omitempty"`
}

type DossierDetailComplet struct {
	Candidat               Candidat           `json:"candidat"`
	Dossier                DossierInscription `json:"dossier"`
	Formule                FormulePermis      `json:"formule"`
	Paiements              []Paiement         `json:"paiements"`
	ProgressionCodePct     int                `json:"progression_code_pct"`
	ProgressionConduitePct int                `json:"progression_conduite_pct"`
	DernierScoreCode       *int               `json:"dernier_score_code,omitempty"`
}

type CreerCandidatInput struct {
	Nom                    string  `json:"nom"`
	Prenoms                string  `json:"prenoms"`
	Genre                  string  `json:"genre"`
	DateNaissance          string  `json:"date_naissance"`
	LieuNaissance          string  `json:"lieu_naissance"`
	TelephoneWhatsApp      string  `json:"telephone_whatsapp"`
	TypePieceID            string  `json:"type_piece_id"`
	NumeroPieceID          string  `json:"numero_piece_id"`
	FormuleID              string  `json:"formule_id"`
	AcompteInitial         float64 `json:"acompte_initial"`
	ModeReglement          string  `json:"mode_reglement"`
	PhotoURL               *string `json:"photo_url,omitempty"`
	VisiteMedicaleStatut   *string `json:"visite_medicale_statut,omitempty"`
	ModeFormationTheorique *string `json:"mode_formation_theorique,omitempty"`
}

// ==============================================================================
// 3. PAIEMENTS & REÇUS THERMIQUES ESC/POS
// ==============================================================================

type Paiement struct {
	ID                   string     `json:"id"`
	AutoEcoleID          string     `json:"auto_ecole_id"`
	AgenceID             string     `json:"agence_id"`
	DossierID            string     `json:"dossier_id"`
	SessionCaisseID      *string    `json:"session_caisse_id,omitempty"`
	Montant              float64    `json:"montant"`
	ModeReglement        string     `json:"mode_reglement"`
	ReferenceTransaction *string    `json:"reference_transaction,omitempty"`
	ReferenceRecu        string     `json:"reference_recu"`
	SoldeApresPaiement   float64    `json:"solde_apres_paiement"`
	DatePaiement         string     `json:"date_paiement"`
	EncaisseParID        *string    `json:"encaisse_par_id,omitempty"`
	CandidatNomComplet   string     `json:"candidat_nom_complet"`
	EstAnnule            bool       `json:"est_annule"`
	CreatedAt            time.Time  `json:"created_at"`
	SyncedAt             *time.Time `json:"synced_at,omitempty"`
}

type EnregistrerPaiementInput struct {
	DossierID            string  `json:"dossier_id"`
	Montant              float64 `json:"montant"`
	ModeReglement        string  `json:"mode_reglement"`
	ReferenceTransaction *string `json:"reference_transaction,omitempty"`
}

type EnregistrerPaiementResponse struct {
	PaiementID          string  `json:"paiement_id"`
	NumeroRecu          string  `json:"numero_recu"`
	CodeSecuriteHash    string  `json:"code_securite_hash"`
	MontantVerse        float64 `json:"montant_verse"`
	NouveauSoldeRestant float64 `json:"nouveau_solde_restant"`
	DatePaiement        string  `json:"date_paiement"`
	NumeroDossier       string  `json:"numero_dossier"`
}

type RecuImprimable struct {
	NumeroRecu        string  `json:"numero_recu"`
	CodeSecuriteHash  string  `json:"code_securite_hash"`
	NombreImpressions int     `json:"nombre_impressions"`
	MontantVerse      float64 `json:"montant_verse"`
	SoldeRestant      float64 `json:"solde_restant"`
	ModeReglement     string  `json:"mode_reglement"`
	DatePaiement      string  `json:"date_paiement"`
	CandidatNom       string  `json:"candidat_nom"`
	NumeroDossier     string  `json:"numero_dossier"`
	FormuleNom        string  `json:"formule_nom"`
	CaissierNom       string  `json:"caissier_nom"`
	EstDuplicata      bool    `json:"est_duplicata"`
}

// ==============================================================================
// 4. TIROIR-CAISSE, SESSIONS & RAPPORT Z
// ==============================================================================

type SessionCaisse struct {
	ID                    string     `json:"id"`
	AgenceID              string     `json:"agence_id"`
	CaissierID            string     `json:"caissier_id"`
	DateOuverture         time.Time  `json:"date_ouverture"`
	DateFermeture         *time.Time `json:"date_fermeture,omitempty"`
	FondDeCaisse          float64    `json:"fond_de_caisse"`
	TotalEspecesTheorique float64    `json:"total_especes_theorique"`
	TotalEspecesCompte    *float64   `json:"total_especes_compte,omitempty"`
	EcartCaisse           *float64   `json:"ecart_caisse,omitempty"`
	Statut                string     `json:"statut"`
	Observations          *string    `json:"observations,omitempty"`
}

type VentilationCaisse struct {
	Especes       float64 `json:"especes"`
	Wave          float64 `json:"wave"`
	OrangeMoney   float64 `json:"orange_money"`
	MtnMomo       float64 `json:"mtn_momo"`
	Moov          float64 `json:"moov"`
	Virement      float64 `json:"virement"`
	TotalEncaisse float64 `json:"total_encaisse"`
}

type SessionCaisseDetail struct {
	SessionOuverte        bool              `json:"session_ouverte"`
	ID                    string            `json:"id"`
	AgenceID              string            `json:"agence_id"`
	CaissierID            string            `json:"caissier_id"`
	CaissierNom           string            `json:"caissier_nom"`
	DateOuverture         string            `json:"date_ouverture"`
	FondDeCaisse          float64           `json:"fond_de_caisse"`
	TotalEspecesTheorique float64           `json:"total_especes_theorique"`
	TotalTiroirAttendu    float64           `json:"total_tiroir_attendu"`
	Statut                string            `json:"statut"`
	Ventilation           VentilationCaisse `json:"ventilation"`
	NombreTransactions    int               `json:"nombre_transactions"`
	Paiements             []Paiement        `json:"paiements"`
}

type CloturerSessionCaisseInput struct {
	TotalEspecesCompte float64 `json:"total_especes_compte"`
	Observations       *string `json:"observations,omitempty"`
}

// ==============================================================================
// 5. PÉDAGOGIE (CODE & CONDUITE)
// ==============================================================================

type PresenceItemInput struct {
	CandidatID string `json:"candidat_id"`
	EstPresent bool   `json:"est_present"`
	ScoreTest  int    `json:"score_test"`
}

type EnregistrerSeanceCodeInput struct {
	ThemeTraite  string              `json:"theme_traite"`
	FormateurNom string              `json:"formateur_nom"`
	HeureDebut   string              `json:"heure_debut"`
	HeureFin     string              `json:"heure_fin"`
	Presences    []PresenceItemInput `json:"presences"`
}

type EnregistrerSeanceCodeResponse struct {
	SeanceID      string `json:"seance_id"`
	TotalPresents int    `json:"total_presents"`
	DateSeance    string `json:"date_seance"`
	Statut        string `json:"statut"`
}

type CompetenceItem struct {
	ID        string `json:"id"`
	Code      string `json:"code"`
	Categorie string `json:"categorie"`
	Libelle   string `json:"libelle"`
	Niveau    int    `json:"niveau"`
}

type SeanceConduite struct {
	ID                      string           `json:"id"`
	AutoEcoleID             string           `json:"auto_ecole_id"`
	AgenceID                string           `json:"agence_id"`
	DossierID               string           `json:"dossier_id"`
	CandidatID              string           `json:"candidat_id"`
	CandidatNom             string           `json:"candidat_nom"`
	TelephoneWhatsApp       string           `json:"telephone_whatsapp"`
	PhotoURL                *string          `json:"photo_url,omitempty"`
	FormuleNom              string           `json:"formule_nom"`
	MoniteurID              *string          `json:"moniteur_id,omitempty"`
	MoniteurNom             string           `json:"moniteur_nom"`
	VehiculeID              *string          `json:"vehicule_id,omitempty"`
	VehiculeImmatriculation string           `json:"vehicule_immatriculation"`
	VehiculeModele          string           `json:"vehicule_modele"`
	DateSeance              string           `json:"date_seance"`
	HeureDebut              string           `json:"heure_debut"`
	HeureFin                string           `json:"heure_fin"`
	DureeMinutes            int              `json:"duree_minutes"`
	Statut                  string           `json:"statut"`
	HeuresDejaEffectuees    int              `json:"heures_deja_effectuees"`
	HeuresForfaitTotal      int              `json:"heures_forfait_total"`
	Competences             []CompetenceItem `json:"competences"`
	CommentaireMoniteur     *string          `json:"commentaire_moniteur,omitempty"`
	KmDebut                 *int             `json:"km_debut,omitempty"`
	KmFin                   *int             `json:"km_fin,omitempty"`
	SignatureEleveURL       *string          `json:"signature_eleve_url,omitempty"`
}

type ValiderSeanceConduiteInput struct {
	SeanceID       string           `json:"seance_id"`
	DossierID      string           `json:"dossier_id"`
	VehiculeID     *string          `json:"vehicule_id,omitempty"`
	DureeHeures    int              `json:"duree_heures"`
	KmDebut        int              `json:"km_debut"`
	KmFin          int              `json:"km_fin"`
	Commentaire    string           `json:"commentaire"`
	Competences    []CompetenceItem `json:"competences"`
	SignatureEleve string           `json:"signature_eleve"`
}

// Les types Examens, Séances Code, Tutoriels, Supports et Quiz ont été extraits dans pedagogie_examens.go

// ==============================================================================
// 6. STATISTIQUES & TABLEAU DE BORD GÉRANT
// ==============================================================================

type GuichetStats struct {
	TotalCandidats          int     `json:"total_candidats"`
	FondDeCaisse            float64 `json:"fond_de_caisse"`
	TotalEspecesTheorique   float64 `json:"total_especes_theorique"`
	EvenementsSyncEnAttente int     `json:"evenements_sync_en_attente"`
	ModeConnectivite        string  `json:"mode_connectivite"`
}

type ImpayeItem struct {
	DossierID          string  `json:"dossier_id"`
	NumeroDossier      string  `json:"numero_dossier"`
	CandidatNom        string  `json:"candidat_nom"`
	TelephoneWhatsApp  string  `json:"telephone_whatsapp"`
	Formule            string  `json:"formule"`
	MontantTotal       float64 `json:"montant_total"`
	SoldeRestant       float64 `json:"solde_restant"`
	DejaPaye           float64 `json:"deja_paye"`
	TauxReglementPct   int     `json:"taux_reglement_pct"`
	HeuresCode         int     `json:"heures_code"`
	HeuresCodeRequises int     `json:"heures_code_requises"`
	DateInscription    string  `json:"date_inscription"`
}

type AssiduItem struct {
	DossierID                       string  `json:"dossier_id"`
	NumeroDossier                   string  `json:"numero_dossier"`
	CandidatNom                     string  `json:"candidat_nom"`
	TelephoneWhatsApp               string  `json:"telephone_whatsapp"`
	Formule                         string  `json:"formule"`
	HeuresCodeEffectuees            int     `json:"heures_code_effectuees"`
	HeuresTheoriques                int     `json:"heures_theoriques"`
	ConduiteDebloqueeFinancierement bool    `json:"conduite_debloquee_financierement"`
	TauxFinancierPct                int     `json:"taux_financier_pct"`
	SoldeRestant                    float64 `json:"solde_restant"`
}

type RecetteMensuelle struct {
	Mois    string  `json:"mois"`
	Montant float64 `json:"montant"`
}

type FormuleStat struct {
	FormuleID string `json:"formule_id"`
	Intitule  string `json:"intitule"`
	Inscrits  int    `json:"inscrits"`
}

type PaiementRecent struct {
	ID            string  `json:"id"`
	NumeroRecu    string  `json:"numero_recu"`
	NumeroDossier string  `json:"numero_dossier"`
	CandidatNom   string  `json:"candidat_nom"`
	FormuleNom    string  `json:"formule_nom"`
	Montant       float64 `json:"montant"`
	ModeReglement string  `json:"mode_reglement"`
	DatePaiement  string  `json:"date_paiement"`
}

type TableauGerantResponse struct {
	TotalCreancesImpayes  float64            `json:"total_creances_impayes"`
	NombreDossiersImpayes int                `json:"nombre_dossiers_impayes"`
	Impayes               []ImpayeItem       `json:"impayes"`
	ElevesAssidusCode     []AssiduItem       `json:"eleves_assidus_code"`
	Ventilation           VentilationCaisse  `json:"ventilation"`
	RecettesMensuelles    []RecetteMensuelle `json:"recettes_mensuelles"`
	InscriptionsFormules  []FormuleStat      `json:"inscriptions_formules"`
	DerniersPaiements     []PaiementRecent   `json:"derniers_paiements"`
	TauxReussiteExamen    float64            `json:"taux_reussite_examen"`
	TotalCandidatsActifs  int                `json:"total_candidats_actifs"`
}

type AccorderDerogationInput struct {
	DossierID string `json:"dossier_id"`
	CodePIN   string `json:"code_pin"`
	Motif     string `json:"motif"`
}

type ReserverSeanceConduiteInput struct {
	DossierID        string  `json:"dossier_id"`
	DateSeance       string  `json:"date_seance"` // YYYY-MM-DD
	HeureDebut       string  `json:"heure_debut"` // ex "08:00"
	HeureFin         string  `json:"heure_fin"`   // ex "10:00"
	VehiculeID       *string `json:"vehicule_id,omitempty"`
	MoniteurID       *string `json:"moniteur_id,omitempty"`
	DerogationForcee bool    `json:"derogation_forcee,omitempty"`
}

type EnregistrerQuizResultatInput struct {
	CandidatID string `json:"candidat_id"`
	Score      int    `json:"score"` // sur 40
}

type UpdatePhotoInput struct {
	PhotoURL string `json:"photo_url"`
}

type SyncEventPayload struct {
	EventID           string                 `json:"event_id"`
	EntityName        string                 `json:"entity_name"`
	Operation         string                 `json:"operation"`
	EntityID          string                 `json:"entity_id"`
	Payload           map[string]interface{} `json:"payload"`
	LocalTimestampUTC string                 `json:"local_timestamp_utc"`
}

type SyncBatchInput struct {
	Events []SyncEventPayload `json:"events"`
}

type SyncBatchResponse struct {
	Statut            string `json:"statut"`
	EvenementsTraites int    `json:"evenements_traites"`
	DoublonsIgnores   int    `json:"doublons_ignores"`
	Erreurs           int    `json:"erreurs"`
	SyncedAt          string `json:"synced_at"`
	Message           string `json:"message,omitempty"`
}

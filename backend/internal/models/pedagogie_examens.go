package models

import (
	"time"
)

// ==============================================================================
// 1. EXAMENS D'ÉTAT & BORDEREAUX MINISTÉRIELS
// ==============================================================================

type SessionExamen struct {
	ID              string    `json:"id"`
	AutoEcoleID     string    `json:"auto_ecole_id"`
	AgenceID        string    `json:"agence_id"`
	CodeSession     string    `json:"code_session"`
	DateExamen      string    `json:"date_examen"`
	CentreExamen    string    `json:"centre_examen"`
	CategoriePermis string    `json:"categorie_permis"`
	StatutSession   string    `json:"statut_session"`
	InspecteurNom   *string   `json:"inspecteur_nom,omitempty"`
	Observations    *string   `json:"observations,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	TotalCandidats  int       `json:"total_candidats"`
	TotalAdmis      int       `json:"total_admis"`
	TotalAjournes   int       `json:"total_ajournes"`
}

type DossierNonEligible struct {
	NumeroDossier string   `json:"numero_dossier"`
	Nom           string   `json:"nom"`
	Motifs        []string `json:"motifs"`
	DossierID     *string  `json:"dossier_id,omitempty"`
	CandidatID    *string  `json:"candidat_id,omitempty"`
}

type ForcerCandidatExamenInput struct {
	SessionExamenID string `json:"session_examen_id"`
	DossierID       string `json:"dossier_id"`
	Motif           string `json:"motif,omitempty"`
	CodePin         string `json:"code_pin,omitempty"`
}

type CandidatExamenItem struct {
	ID                       string   `json:"id"`
	SessionExamenID          string   `json:"session_examen_id"`
	DossierID                string   `json:"dossier_id"`
	NumeroDossier            string   `json:"numero_dossier"`
	CandidatID               string   `json:"candidat_id"`
	Nom                      string   `json:"nom"`
	Prenoms                  string   `json:"prenoms"`
	Genre                    string   `json:"genre"`
	DateNaissance            string   `json:"date_naissance"`
	LieuNaissance            string   `json:"lieu_naissance"`
	TelephoneWhatsApp        string   `json:"telephone_whatsapp"`
	TypePieceID              string   `json:"type_piece_id"`
	NumeroPieceID            string   `json:"numero_piece_id"`
	PhotoURL                 *string  `json:"photo_url,omitempty"`
	NumeroTable              string   `json:"numero_table"`
	ResultatCode             string   `json:"resultat_code"`
	ResultatConduite         string   `json:"resultat_conduite"`
	ResultatFinal            string   `json:"resultat_final"`
	DateNotificationEnvoyee  *string  `json:"date_notification_envoyee,omitempty"`
	SoldeRestant             float64  `json:"solde_restant"`
	VisiteMedicaleStatut     string   `json:"visite_medicale_statut"`
	HeuresCodeEffectuees     int      `json:"heures_code_effectuees"`
	HeuresCodeRequises       int      `json:"heures_code_requises"`
	HeuresConduiteEffectuees int      `json:"heures_conduite_effectuees"`
	HeuresConduiteRequises   int      `json:"heures_conduite_requises"`
	EstEligible              bool     `json:"est_eligible"`
	EstDerogation            bool     `json:"est_derogation"`
	MotifsIneligibilite      []string `json:"motifs_ineligibilite,omitempty"`
}

type SessionExamenActiveResponse struct {
	Session     *SessionExamen       `json:"session"`
	Candidats   []CandidatExamenItem `json:"candidats"`
	Ineligibles []DossierNonEligible `json:"ineligibles"`
}

type CreerOuModifierSessionExamenInput struct {
	ID              *string `json:"id,omitempty"`
	AgenceID        *string `json:"agence_id,omitempty"`
	CodeSession     string  `json:"code_session"`
	DateExamen      string  `json:"date_examen"`
	CentreExamen    string  `json:"centre_examen"`
	CategoriePermis string  `json:"categorie_permis"`
	InspecteurNom   *string `json:"inspecteur_nom,omitempty"`
	Observations    *string `json:"observations,omitempty"`
	StatutSession   *string `json:"statut_session,omitempty"`
}

type SupprimerSessionExamenInput struct {
	SessionID string `json:"session_id"`
}

type CandidatExamenDetail struct {
	ID                      string   `json:"id"`
	SessionExamenID         string   `json:"session_examen_id"`
	DossierID               string   `json:"dossier_id"`
	CandidatID              string   `json:"candidat_id"`
	NumeroTable             string   `json:"numero_table"`
	ResultatCode            string   `json:"resultat_code"`
	ResultatConduite        string   `json:"resultat_conduite"`
	ResultatFinal           string   `json:"resultat_final"`
	DateNotificationEnvoyee *string  `json:"date_notification_envoyee,omitempty"`
	CandidatNom             string   `json:"candidat_nom"`
	CandidatPrenoms         string   `json:"candidat_prenoms"`
	TelephoneWhatsApp       string   `json:"telephone_whatsapp"`
	NumeroPieceID           string   `json:"numero_piece_id"`
	NumeroDossier           string   `json:"numero_dossier"`
	SoldeRestant            float64  `json:"solde_restant"`
	HeuresCode              int      `json:"heures_code"`
	HeuresConduite          int      `json:"heures_conduite"`
	EstEligible             bool     `json:"est_eligible"`
	MotifsIneligibilite     []string `json:"motifs_ineligibilite,omitempty"`
}

type CreerSessionExamenInput struct {
	DateExamen      string   `json:"date_examen"`
	CentreExamen    string   `json:"centre_examen"`
	CategoriePermis string   `json:"categorie_permis"`
	InspecteurNom   *string  `json:"inspecteur_nom,omitempty"`
	CandidatsIDs    []string `json:"candidats_ids"`
}

type ResultatCandidatInput struct {
	CandidatExamenID string `json:"candidat_examen_id"`
	ResultatCode     string `json:"resultat_code"`
	ResultatConduite string `json:"resultat_conduite"`
	ResultatFinal    string `json:"resultat_final"`
}

type EnregistrerResultatsExamenInput struct {
	SessionExamenID string                  `json:"session_examen_id"`
	Resultats       []ResultatCandidatInput `json:"resultats"`
}

// ==============================================================================
// 2. SÉANCES CODE, TUTORIELS VIDÉO, SUPPORTS COURS & QUIZ INTERACTIFS
// ==============================================================================

type SeanceCode struct {
	ID              string  `json:"id"`
	AutoEcoleID     string  `json:"auto_ecole_id"`
	AgenceID        string  `json:"agence_id"`
	DateSeance      string  `json:"date_seance"`
	HeureDebut      string  `json:"heure_debut"`
	HeureFin        string  `json:"heure_fin"`
	ThemeTraite     string  `json:"theme_traite"`
	FormateurNom    string  `json:"formateur_nom"`
	MoniteurID      *string `json:"moniteur_id,omitempty"`
	TypeSeance      string  `json:"type_seance"`
	LienVisio       *string `json:"lien_visio,omitempty"`
	MotDePasseVisio *string `json:"mot_de_passe_visio,omitempty"`
	StatutSeance    string  `json:"statut_seance"`
}

type CreerSeanceCodeInput struct {
	ThemeTraite     string  `json:"theme_traite"`
	FormateurNom    string  `json:"formateur_nom"`
	MoniteurID      *string `json:"moniteur_id,omitempty"`
	DateSeance      string  `json:"date_seance"`
	HeureDebut      string  `json:"heure_debut"`
	HeureFin        string  `json:"heure_fin"`
	TypeSeance      string  `json:"type_seance"`
	LienVisio       *string `json:"lien_visio,omitempty"`
	MotDePasseVisio *string `json:"mot_de_passe_visio,omitempty"`
}

type TutorielVideo struct {
	ID              string  `json:"id"`
	AutoEcoleID     string  `json:"auto_ecole_id"`
	Titre           string  `json:"titre"`
	CategoriePermis string  `json:"categorie_permis"`
	Description     *string `json:"description,omitempty"`
	UrlVideo        string  `json:"url_video"`
	Plateforme      string  `json:"plateforme"`
	DureeMinutes    int     `json:"duree_minutes"`
	MiniatureURL    *string `json:"miniature_url,omitempty"`
	VuesCount       int     `json:"vues_count"`
	AjoutePar       *string `json:"ajoute_par,omitempty"`
	CreatedAt       string  `json:"created_at"`
}

type CreerTutorielVideoInput struct {
	Titre           string  `json:"titre"`
	CategoriePermis string  `json:"categorie_permis"`
	Description     *string `json:"description,omitempty"`
	UrlVideo        string  `json:"url_video"`
	Plateforme      *string `json:"plateforme,omitempty"`
	DureeMinutes    *int    `json:"duree_minutes,omitempty"`
	MiniatureURL    *string `json:"miniature_url,omitempty"`
	AjoutePar       *string `json:"ajoute_par,omitempty"`
}

type SupportCours struct {
	ID                   string  `json:"id"`
	AutoEcoleID          string  `json:"auto_ecole_id"`
	Titre                string  `json:"titre"`
	CategoriePermis      string  `json:"categorie_permis"`
	Description          *string `json:"description,omitempty"`
	TypeSupport          string  `json:"type_support"`
	FichierNom           *string `json:"fichier_nom,omitempty"`
	FichierTaille        *string `json:"fichier_taille,omitempty"`
	FichierData          *string `json:"fichier_data,omitempty"`
	UrlExterne           *string `json:"url_externe,omitempty"`
	AjoutePar            *string `json:"ajoute_par,omitempty"`
	TelechargementsCount int     `json:"telechargements_count"`
	CreatedAt            string  `json:"created_at"`
}

type CreerSupportCoursInput struct {
	Titre           string  `json:"titre"`
	CategoriePermis string  `json:"categorie_permis"`
	Description     *string `json:"description,omitempty"`
	TypeSupport     string  `json:"type_support"`
	FichierNom      *string `json:"fichier_nom,omitempty"`
	FichierTaille   *string `json:"fichier_taille,omitempty"`
	FichierData     *string `json:"fichier_data,omitempty"`
	UrlExterne      *string `json:"url_externe,omitempty"`
	AjoutePar       *string `json:"ajoute_par,omitempty"`
}

type QuizOptionItem struct {
	ID          string `json:"id"`
	Texte       string `json:"texte"`
	EstCorrecte bool   `json:"est_correcte"`
}

type QuizQuestionItem struct {
	ID              string           `json:"id"`
	Intitule        string           `json:"intitule"`
	TypeQuestion    string           `json:"type_question"`
	Points          int              `json:"points"`
	Obligatoire     bool             `json:"obligatoire"`
	ImageURL        *string          `json:"image_url,omitempty"`
	Explication     *string          `json:"explication,omitempty"`
	ReponseAttendue *string          `json:"reponse_attendue,omitempty"`
	Options         []QuizOptionItem `json:"options"`
}

type QuizFormulaire struct {
	ID               string             `json:"id"`
	AutoEcoleID      string             `json:"auto_ecole_id"`
	MoniteurID       *string            `json:"moniteur_id,omitempty"`
	MoniteurNom      string             `json:"moniteur_nom"`
	Titre            string             `json:"titre"`
	Description      *string            `json:"description,omitempty"`
	CategoriePermis  string             `json:"categorie_permis"`
	DureeMinutes     int                `json:"duree_minutes"`
	SeuilReussitePct int                `json:"seuil_reussite_pct"`
	EstPublie        bool               `json:"est_publie"`
	Questions        []QuizQuestionItem `json:"questions"`
	ReponsesCount    int                `json:"reponses_count"`
	CreatedAt        string             `json:"created_at"`
}

type CreerQuizFormulaireInput struct {
	MoniteurID       *string            `json:"moniteur_id,omitempty"`
	MoniteurNom      string             `json:"moniteur_nom"`
	Titre            string             `json:"titre"`
	Description      *string            `json:"description,omitempty"`
	CategoriePermis  string             `json:"categorie_permis"`
	DureeMinutes     int                `json:"duree_minutes"`
	SeuilReussitePct int                `json:"seuil_reussite_pct"`
	EstPublie        bool               `json:"est_publie"`
	Questions        []QuizQuestionItem `json:"questions"`
}

type QuizReponseDetailItem struct {
	QuestionID        string   `json:"question_id"`
	IntituleQuestion  string   `json:"intitule_question"`
	TypeQuestion      string   `json:"type_question"`
	ChoixIDs          []string `json:"choix_ids"`
	ChoixLabels       []string `json:"choix_labels"`
	TexteReponse      string   `json:"texte_reponse"`
	EstCorrecte       bool     `json:"est_correcte"`
	PointsObtenus     int      `json:"points_obtenus"`
	PointsMax         int      `json:"points_max"`
	BonneReponseLabel string   `json:"bonne_reponse_label"`
	Explication       string   `json:"explication,omitempty"`
}

type QuizSoumissionEleve struct {
	ID              string                  `json:"id"`
	QuizID          string                  `json:"quiz_id"`
	QuizTitre       string                  `json:"quiz_titre"`
	AutoEcoleID     string                  `json:"auto_ecole_id"`
	CandidatID      string                  `json:"candidat_id"`
	CandidatNom     string                  `json:"candidat_nom"`
	NumeroDossier   string                  `json:"numero_dossier"`
	CategoriePermis string                  `json:"categorie_permis"`
	ScoreObtenu     int                     `json:"score_obtenu"`
	TotalPoints     int                     `json:"total_points"`
	Pourcentage     int                     `json:"pourcentage"`
	EstReussi       bool                    `json:"est_reussi"`
	DureeSecondes   int                     `json:"duree_secondes"`
	Details         []QuizReponseDetailItem `json:"details"`
	SubmittedAt     string                  `json:"submitted_at"`
}

type SoumettreQuizReponseInput struct {
	QuizID          string                  `json:"quiz_id"`
	QuizTitre       string                  `json:"quiz_titre"`
	CandidatID      string                  `json:"candidat_id"`
	CandidatNom     string                  `json:"candidat_nom"`
	NumeroDossier   string                  `json:"numero_dossier"`
	CategoriePermis string                  `json:"categorie_permis"`
	ScoreObtenu     int                     `json:"score_obtenu"`
	TotalPoints     int                     `json:"total_points"`
	Pourcentage     int                     `json:"pourcentage"`
	EstReussi       bool                    `json:"est_reussi"`
	DureeSecondes   int                     `json:"duree_secondes"`
	Details         []QuizReponseDetailItem `json:"details"`
}


package models

import "time"

// EpreuvePdf représente un sujet ou une annale officielle stockée en base de données PostgreSQL
type EpreuvePdf struct {
	ID                   string    `json:"id"`
	AutoEcoleID          string    `json:"auto_ecole_id"`
	Titre                string    `json:"titre"`
	CategoriePermis      string    `json:"categorie_permis"`
	Description          *string   `json:"description,omitempty"`
	Annee                int       `json:"annee"`
	Session              *string   `json:"session,omitempty"`
	DureeMinutes         int       `json:"duree_minutes"`
	NombreQuestions      int       `json:"nombre_questions"`
	AvecCorrige          bool      `json:"avec_corrige"`
	FichierNom           string    `json:"fichier_nom"`
	FichierTaille        *string   `json:"fichier_taille,omitempty"`
	FichierData          string    `json:"fichier_data"`
	AjoutePar            *string   `json:"ajoute_par,omitempty"`
	TelechargementsCount int       `json:"telechargements_count"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// CreerEpreuveInput structure le payload lors du dépôt d'une épreuve par le Directeur ou la Secrétaire
type CreerEpreuveInput struct {
	Titre           string  `json:"titre"`
	CategoriePermis string  `json:"categorie_permis"`
	Description     *string `json:"description,omitempty"`
	Annee           int     `json:"annee"`
	Session         *string `json:"session,omitempty"`
	DureeMinutes    int     `json:"duree_minutes"`
	NombreQuestions int     `json:"nombre_questions"`
	AvecCorrige     bool    `json:"avec_corrige"`
	FichierNom      string  `json:"fichier_nom"`
	FichierTaille   *string `json:"fichier_taille,omitempty"`
	FichierData     string  `json:"fichier_data"`
	AjoutePar       *string `json:"ajoute_par,omitempty"`
}

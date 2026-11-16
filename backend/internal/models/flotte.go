package models

import "time"

type Vehicule struct {
	ID                     string    `json:"id"`
	AgenceID               string    `json:"agence_id"`
	Immatriculation        string    `json:"immatriculation"`
	MarqueModele           string    `json:"marque_modele"`
	TypeBoite              string    `json:"type_boite"`
	KilometrageActuel      int       `json:"kilometrage_actuel"`
	DateExpirationAssurance string    `json:"date_expiration_assurance"`
	DateVisiteTechnique    string    `json:"date_visite_technique"`
	Statut                 string    `json:"statut"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

type CreerVehiculeInput struct {
	AgenceID               *string `json:"agence_id,omitempty"`
	Immatriculation        string  `json:"immatriculation"`
	MarqueModele           string  `json:"marque_modele"`
	TypeBoite              string  `json:"type_boite"`
	KilometrageActuel      int     `json:"kilometrage_actuel"`
	DateExpirationAssurance string  `json:"date_expiration_assurance"`
	DateVisiteTechnique    string  `json:"date_visite_technique"`
	Statut                 string  `json:"statut"`
}

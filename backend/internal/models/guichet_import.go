package models

// ImporterCandidatsBatchInput représente une requête d'enrôlement en masse depuis Excel.
type ImporterCandidatsBatchInput struct {
	Candidats          []CreerCandidatInput `json:"candidats"`
	IgnorerDoublons    bool                 `json:"ignorer_doublons"`
	ImpacterCaisseJour bool                 `json:"impacter_caisse_jour"`
}

// ImportErreurLigne décrit une ligne ignorée ou en erreur lors de l'import batch.
type ImportErreurLigne struct {
	Index         int    `json:"index"`
	Nom           string `json:"nom"`
	Prenoms       string `json:"prenoms"`
	NumeroPieceID string `json:"numero_piece_id"`
	Raison        string `json:"raison"`
}

// ImportDossierCree résume un dossier créé avec succès lors de l'import batch.
type ImportDossierCree struct {
	CandidatID    string  `json:"candidat_id"`
	DossierID     string  `json:"dossier_id"`
	NumeroDossier string  `json:"numero_dossier"`
	NomComplet    string  `json:"nom_complet"`
	SoldeRestant  float64 `json:"solde_restant"`
}

// ImporterCandidatsBatchResult résume le bilan complet d'un import Excel.
type ImporterCandidatsBatchResult struct {
	TotalRecus    int                 `json:"total_recus"`
	TotalCrees    int                 `json:"total_crees"`
	TotalIgnores  int                 `json:"total_ignores"`
	TotalErreurs  int                 `json:"total_erreurs"`
	DossiersCrees []ImportDossierCree `json:"dossiers_crees"`
	Erreurs       []ImportErreurLigne `json:"erreurs"`
}

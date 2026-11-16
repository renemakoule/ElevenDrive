package service

import (
	"context"
	"errors"

	"autoschool/backend/internal/models"
	"autoschool/backend/internal/repository"
)

// ReprogrammerSeanceCode met à jour la date et les horaires d'une séance expirée.
func (s *GuichetService) ReprogrammerSeanceCode(ctx context.Context, autoEcoleID, id, dateSeance, heureDebut, heureFin string) (*models.SeanceCode, error) {
	if autoEcoleID == "" || id == "" {
		return nil, errors.New("identifiant auto-école ou séance manquant")
	}
	if dateSeance == "" || heureDebut == "" || heureFin == "" {
		return nil, errors.New("la nouvelle date et les horaires sont obligatoires")
	}
	return s.guichetRepo.ReprogrammerSeanceCode(ctx, autoEcoleID, id, dateSeance, heureDebut, heureFin)
}

// GetSyncPullSnapshot retourne un instantané complet de l'auto-école pour la synchronisation descendante.
func (s *GuichetService) GetSyncPullSnapshot(ctx context.Context, autoEcoleID string) (*repository.SyncPullSnapshot, error) {
	if autoEcoleID == "" {
		return nil, errors.New("identifiant auto-école manquant")
	}
	return s.guichetRepo.GetSyncPullSnapshot(ctx, autoEcoleID)
}

// ImporterCandidatsBatch valide et orchestre l'importation en masse de candidats depuis un tableur.
func (s *GuichetService) ImporterCandidatsBatch(
	ctx context.Context,
	autoEcoleID, agenceID, caissierID string,
	input *models.ImporterCandidatsBatchInput,
) (*models.ImporterCandidatsBatchResult, error) {
	if autoEcoleID == "" || agenceID == "" {
		return nil, errors.New("identifiant auto-école ou agence manquant")
	}
	if input == nil || len(input.Candidats) == 0 {
		return nil, errors.New("aucun candidat à importer")
	}
	if len(input.Candidats) > 1000 {
		return nil, errors.New("limite dépassée : maximum 1000 candidats par import")
	}

	return s.guichetRepo.ImporterCandidatsBatch(ctx, autoEcoleID, agenceID, caissierID, input)
}

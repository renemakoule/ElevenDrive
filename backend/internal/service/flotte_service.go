package service

import (
	"context"
	"errors"
	"strings"

	"autoschool/backend/internal/models"
	"autoschool/backend/internal/repository"
)

type FlotteService struct {
	flotteRepo *repository.FlotteRepository
}

func NewFlotteService(flotteRepo *repository.FlotteRepository) *FlotteService {
	return &FlotteService{flotteRepo: flotteRepo}
}

func (s *FlotteService) GetVehicules(ctx context.Context, autoEcoleID string) ([]models.Vehicule, error) {
	if autoEcoleID == "" {
		return nil, errors.New("identifiant auto-école manquant")
	}
	return s.flotteRepo.GetVehicules(ctx, autoEcoleID)
}

func (s *FlotteService) CreerVehicule(ctx context.Context, autoEcoleID, agenceID string, input *models.CreerVehiculeInput) (*models.Vehicule, error) {
	if autoEcoleID == "" {
		return nil, errors.New("identifiant auto-école manquant")
	}
	if strings.TrimSpace(input.Immatriculation) == "" {
		return nil, errors.New("numéro d'immatriculation requis")
	}
	if strings.TrimSpace(input.MarqueModele) == "" {
		return nil, errors.New("marque et modèle requis")
	}
	if strings.TrimSpace(input.DateExpirationAssurance) == "" {
		return nil, errors.New("date d'expiration de l'assurance requise")
	}
	if strings.TrimSpace(input.DateVisiteTechnique) == "" {
		return nil, errors.New("date de visite technique requise")
	}
	return s.flotteRepo.CreerVehicule(ctx, autoEcoleID, agenceID, input)
}

func (s *FlotteService) UpdateStatut(ctx context.Context, autoEcoleID, vehiculeID, statut string) error {
	if autoEcoleID == "" || vehiculeID == "" {
		return errors.New("identifiants manquants")
	}
	statutUpper := strings.ToUpper(strings.TrimSpace(statut))
	if statutUpper != "DISPONIBLE" && statutUpper != "EN_LECON" && statutUpper != "EN_PANNE" && statutUpper != "EN_ENTRETIEN" {
		return errors.New("statut invalide")
	}
	return s.flotteRepo.UpdateVehiculeStatut(ctx, autoEcoleID, vehiculeID, statutUpper)
}

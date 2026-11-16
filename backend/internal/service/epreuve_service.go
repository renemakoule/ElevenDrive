package service

import (
	"context"
	"errors"
	"strings"

	"autoschool/backend/internal/models"
	"autoschool/backend/internal/repository"
)

type EpreuveService struct {
	epreuveRepo *repository.EpreuveRepository
}

func NewEpreuveService(epreuveRepo *repository.EpreuveRepository) *EpreuveService {
	return &EpreuveService{epreuveRepo: epreuveRepo}
}

func (s *EpreuveService) GetEpreuves(ctx context.Context, autoEcoleID, categorie string) ([]models.EpreuvePdf, error) {
	if autoEcoleID == "" {
		return nil, errors.New("identifiant auto-école manquant")
	}
	return s.epreuveRepo.GetEpreuves(ctx, autoEcoleID, categorie)
}

func (s *EpreuveService) CreerEpreuve(ctx context.Context, autoEcoleID string, input *models.CreerEpreuveInput) (*models.EpreuvePdf, error) {
	if autoEcoleID == "" {
		return nil, errors.New("identifiant auto-école manquant")
	}
	if strings.TrimSpace(input.Titre) == "" {
		return nil, errors.New("le titre de l'épreuve est requis")
	}
	if strings.TrimSpace(input.CategoriePermis) == "" {
		return nil, errors.New("la catégorie de permis est requise")
	}
	if strings.TrimSpace(input.FichierData) == "" {
		return nil, errors.New("le document PDF est requis")
	}
	if input.Annee == 0 {
		input.Annee = 2024
	}
	if input.DureeMinutes <= 0 {
		input.DureeMinutes = 45
	}
	if input.NombreQuestions <= 0 {
		input.NombreQuestions = 40
	}
	return s.epreuveRepo.CreerEpreuve(ctx, autoEcoleID, input)
}

func (s *EpreuveService) SupprimerEpreuve(ctx context.Context, autoEcoleID, epreuveID string) error {
	if autoEcoleID == "" || epreuveID == "" {
		return errors.New("identifiants manquants")
	}
	return s.epreuveRepo.SupprimerEpreuve(ctx, autoEcoleID, epreuveID)
}

func (s *EpreuveService) IncrementerTelechargement(ctx context.Context, autoEcoleID, epreuveID string) error {
	if autoEcoleID == "" || epreuveID == "" {
		return errors.New("identifiants manquants")
	}
	return s.epreuveRepo.IncrementerTelechargement(ctx, autoEcoleID, epreuveID)
}

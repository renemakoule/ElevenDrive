package service

import (
	"context"
	"errors"
	"strings"

	"autoschool/backend/internal/models"
	"autoschool/backend/internal/repository"
)

type GuichetService struct {
	guichetRepo *repository.GuichetRepository
	userRepo    *repository.UserRepository
}

func NewGuichetService(guichetRepo *repository.GuichetRepository, userRepo *repository.UserRepository) *GuichetService {
	return &GuichetService{guichetRepo: guichetRepo, userRepo: userRepo}
}

func (s *GuichetService) GetFormules(ctx context.Context, autoEcoleID string) ([]models.FormulePermis, error) {
	if autoEcoleID == "" {
		return nil, errors.New("identifiant auto-école manquant")
	}
	return s.guichetRepo.GetFormules(ctx, autoEcoleID)
}

// CreerFormule permet au Directeur ou à la Secrétaire de définir une formule et de fixer son prix initial.
func (s *GuichetService) CreerFormule(
	ctx context.Context,
	autoEcoleID, userRole string,
	input models.CreerFormuleInput,
) (*models.FormulePermis, error) {
	if autoEcoleID == "" {
		return nil, errors.New("identifiant auto-école manquant")
	}

	// Règle métier : Le Directeur ou la Secrétaire fixe les prix initiaux des permis
	roleUpper := models.NormalizeRole(userRole)
	if roleUpper != models.RoleDirecteur && roleUpper != models.RoleSecretaire && roleUpper != models.RoleCaissier {
		return nil, errors.New("accès refusé : seul le Directeur ou la Secrétaire peut créer une formule et fixer son prix")
	}

	if strings.TrimSpace(input.Intitule) == "" {
		return nil, errors.New("l'intitulé officiel de la formule est obligatoire")
	}
	if strings.TrimSpace(input.CategoriePermis) == "" {
		return nil, errors.New("la catégorie de permis (A, B, C, D, E) est obligatoire")
	}
	if input.PrixTotal <= 0 {
		return nil, errors.New("le prix total fixé doit être strictement supérieur à zéro")
	}
	if input.AcompteMinimum < 0 || input.AcompteMinimum > input.PrixTotal {
		return nil, errors.New("l'acompte minimum ne peut excéder le prix total de la formule")
	}

	return s.guichetRepo.CreerFormule(ctx, autoEcoleID, input)
}

// ModifierPrixFormule permet EXCLUSIVEMENT au Directeur de modifier un prix déjà fixé.
func (s *GuichetService) ModifierPrixFormule(
	ctx context.Context,
	autoEcoleID, userRole string,
	input models.ModifierPrixFormuleInput,
) (*models.FormulePermis, error) {
	if autoEcoleID == "" {
		return nil, errors.New("identifiant auto-école manquant")
	}

	// Règle métier stricte : SEUL le Directeur peut modifier un prix qui a été fixé !
	roleUpper := models.NormalizeRole(userRole)
	if roleUpper != models.RoleDirecteur {
		return nil, errors.New("accès formellement refusé : seul le Directeur est habilité à modifier le prix d'un permis déjà fixé")
	}

	if strings.TrimSpace(input.FormuleID) == "" {
		return nil, errors.New("l'identifiant de la formule est obligatoire")
	}
	if input.NouveauPrix <= 0 {
		return nil, errors.New("le nouveau prix doit être strictement supérieur à zéro")
	}
	if input.AcompteMinimum != nil && (*input.AcompteMinimum < 0 || *input.AcompteMinimum > input.NouveauPrix) {
		return nil, errors.New("le nouvel acompte minimum ne peut pas excéder le nouveau prix")
	}

	return s.guichetRepo.ModifierPrixFormule(ctx, autoEcoleID, input)
}

func (s *GuichetService) GetGuichetStats(ctx context.Context, autoEcoleID, agenceID string) (*models.GuichetStats, error) {
	if autoEcoleID == "" {
		return nil, errors.New("identifiant auto-école manquant")
	}
	return s.guichetRepo.GetGuichetStats(ctx, autoEcoleID, agenceID)
}

func (s *GuichetService) SearchCandidats(ctx context.Context, autoEcoleID, query string) ([]models.Candidat, error) {
	if autoEcoleID == "" {
		return nil, errors.New("identifiant auto-école manquant")
	}
	return s.guichetRepo.SearchCandidats(ctx, autoEcoleID, query)
}

func (s *GuichetService) GetCandidatDossier(ctx context.Context, autoEcoleID, candidatID string) (*models.DossierDetailComplet, error) {
	if autoEcoleID == "" || candidatID == "" {
		return nil, errors.New("paramètres manquants")
	}
	return s.guichetRepo.GetCandidatDossier(ctx, autoEcoleID, candidatID)
}

func (s *GuichetService) CreerCandidatAvecDossier(
	ctx context.Context,
	autoEcoleID, agenceID, caissierID, caissierNom string,
	input *models.CreerCandidatInput,
) (*models.Candidat, *models.DossierDetailComplet, error) {
	if strings.TrimSpace(input.Nom) == "" || strings.TrimSpace(input.Prenoms) == "" {
		return nil, nil, errors.New("le nom et prénom du candidat sont obligatoires")
	}
	if strings.TrimSpace(input.TelephoneWhatsApp) == "" {
		return nil, nil, errors.New("le numéro WhatsApp du candidat est obligatoire")
	}
	if strings.TrimSpace(input.NumeroPieceID) == "" {
		return nil, nil, errors.New("le numéro de pièce d'identité est obligatoire")
	}
	if strings.TrimSpace(input.FormuleID) == "" {
		return nil, nil, errors.New("le choix d'une formule de permis est obligatoire")
	}

	return s.guichetRepo.CreerCandidatAvecDossier(ctx, autoEcoleID, agenceID, caissierID, caissierNom, input)
}

func (s *GuichetService) EnregistrerPaiement(
	ctx context.Context,
	autoEcoleID, agenceID, caissierID string,
	input *models.EnregistrerPaiementInput,
) (*models.EnregistrerPaiementResponse, error) {
	if strings.TrimSpace(input.DossierID) == "" {
		return nil, errors.New("identifiant de dossier manquant")
	}
	if input.Montant <= 0 {
		return nil, errors.New("le montant doit être supérieur à zéro")
	}
	return s.guichetRepo.EnregistrerPaiement(ctx, autoEcoleID, agenceID, caissierID, input)
}

func (s *GuichetService) GetSessionCaisseDetail(ctx context.Context, autoEcoleID, agenceID string) (*models.SessionCaisseDetail, error) {
	return s.guichetRepo.GetSessionCaisseDetail(ctx, autoEcoleID, agenceID)
}

func (s *GuichetService) CloturerSessionCaisse(
	ctx context.Context,
	autoEcoleID, agenceID, caissierID string,
	input *models.CloturerSessionCaisseInput,
) error {
	return s.guichetRepo.CloturerSessionCaisse(ctx, autoEcoleID, agenceID, caissierID, input)
}

func (s *GuichetService) EnregistrerSeanceCode(
	ctx context.Context,
	autoEcoleID, agenceID string,
	input *models.EnregistrerSeanceCodeInput,
) (*models.EnregistrerSeanceCodeResponse, error) {
	if strings.TrimSpace(input.ThemeTraite) == "" {
		return nil, errors.New("le thème traité est obligatoire")
	}
	return s.guichetRepo.EnregistrerSeanceCode(ctx, autoEcoleID, agenceID, input)
}

func (s *GuichetService) GetTableauGerant(ctx context.Context, autoEcoleID string) (*models.TableauGerantResponse, error) {
	return s.guichetRepo.GetTableauGerant(ctx, autoEcoleID)
}

func (s *GuichetService) AccorderDerogationConduite(
	ctx context.Context,
	autoEcoleID, dossierID, codePIN, motif string,
) error {
	if autoEcoleID == "" || strings.TrimSpace(dossierID) == "" {
		return errors.New("identifiants requis")
	}
	if strings.TrimSpace(codePIN) == "" {
		return errors.New("le code PIN manager est obligatoire")
	}
	// Vérification des droits et du PIN manager
	manager, err := s.userRepo.VerifyManagerPin(ctx, autoEcoleID, codePIN)
	if err != nil || manager == nil {
		return errors.New("code PIN manager invalide ou droits insuffisants")
	}

	return s.guichetRepo.AccorderDerogationConduite(ctx, autoEcoleID, dossierID)
}

func (s *GuichetService) GetPlanningConduite(
	ctx context.Context,
	autoEcoleID, moniteurID, dateSeance string,
) ([]models.SeanceConduite, error) {
	return s.guichetRepo.GetPlanningConduite(ctx, autoEcoleID, moniteurID, dateSeance)
}

func (s *GuichetService) ValiderSeanceConduite(
	ctx context.Context,
	autoEcoleID, moniteurID string,
	input *models.ValiderSeanceConduiteInput,
) (*models.SeanceConduite, error) {
	return s.guichetRepo.ValiderSeanceConduite(ctx, autoEcoleID, moniteurID, input)
}

func (s *GuichetService) GetSessionExamenActive(
	ctx context.Context,
	autoEcoleID string,
) (*models.SessionExamenActiveResponse, error) {
	return s.guichetRepo.GetSessionExamenActive(ctx, autoEcoleID)
}

func (s *GuichetService) ForcerCandidatExamen(
	ctx context.Context,
	autoEcoleID string,
	input models.ForcerCandidatExamenInput,
) (*models.CandidatExamenItem, error) {
	if autoEcoleID == "" {
		return nil, errors.New("identifiant auto-école manquant")
	}
	if strings.TrimSpace(input.SessionExamenID) == "" {
		return nil, errors.New("identifiant session examen requis")
	}
	if strings.TrimSpace(input.DossierID) == "" {
		return nil, errors.New("identifiant dossier requis")
	}
	return s.guichetRepo.ForcerCandidatExamen(ctx, autoEcoleID, input)
}

func (s *GuichetService) EnregistrerResultatsExamen(
	ctx context.Context,
	autoEcoleID string,
	input *models.EnregistrerResultatsExamenInput,
) error {
	return s.guichetRepo.EnregistrerResultatsExamen(ctx, autoEcoleID, input)
}

func (s *GuichetService) ReserverSeanceConduite(
	ctx context.Context,
	autoEcoleID string,
	input models.ReserverSeanceConduiteInput,
) (*models.SeanceConduite, error) {
	return s.guichetRepo.ReserverSeanceConduite(ctx, autoEcoleID, input)
}

func (s *GuichetService) AnnulerSeanceConduite(
	ctx context.Context,
	autoEcoleID string,
	seanceID string,
) error {
	return s.guichetRepo.AnnulerSeanceConduite(ctx, autoEcoleID, seanceID)
}

func (s *GuichetService) EnregistrerQuizResultat(
	ctx context.Context,
	autoEcoleID string,
	input models.EnregistrerQuizResultatInput,
) error {
	return s.guichetRepo.EnregistrerQuizResultat(ctx, autoEcoleID, input)
}

func (s *GuichetService) CreerOuModifierSessionExamen(
	ctx context.Context,
	autoEcoleID string,
	input models.CreerOuModifierSessionExamenInput,
) (*models.SessionExamen, error) {
	return s.guichetRepo.CreerOuModifierSessionExamen(ctx, autoEcoleID, input)
}

func (s *GuichetService) SupprimerSessionExamen(
	ctx context.Context,
	autoEcoleID string,
	sessionID string,
) error {
	return s.guichetRepo.SupprimerSessionExamen(ctx, autoEcoleID, sessionID)
}

func (s *GuichetService) ReimprimerRecu(ctx context.Context, autoEcoleID, paiementID string) (*models.RecuImprimable, error) {
	if autoEcoleID == "" || strings.TrimSpace(paiementID) == "" {
		return nil, errors.New("identifiant du paiement requis")
	}
	return s.guichetRepo.ReimprimerRecu(ctx, autoEcoleID, paiementID)
}

func (s *GuichetService) UpdateCandidatPhoto(ctx context.Context, autoEcoleID, candidatID, photoURL string) error {
	if autoEcoleID == "" || strings.TrimSpace(candidatID) == "" {
		return errors.New("identifiant du candidat requis")
	}
	return s.guichetRepo.UpdateCandidatPhoto(ctx, autoEcoleID, candidatID, photoURL)
}

func (s *GuichetService) GetSeancesCode(ctx context.Context, autoEcoleID, filtreType string) ([]models.SeanceCode, error) {
	return s.guichetRepo.GetSeancesCode(ctx, autoEcoleID, filtreType)
}

func (s *GuichetService) CreerSeanceCode(ctx context.Context, autoEcoleID, agenceID string, input models.CreerSeanceCodeInput) (*models.SeanceCode, error) {
	if strings.TrimSpace(input.ThemeTraite) == "" {
		return nil, errors.New("le thème traité est obligatoire")
	}
	return s.guichetRepo.CreerSeanceCode(ctx, autoEcoleID, agenceID, input)
}

func (s *GuichetService) GetTutorielsVideos(ctx context.Context, autoEcoleID, categorie string) ([]models.TutorielVideo, error) {
	return s.guichetRepo.GetTutorielsVideos(ctx, autoEcoleID, categorie)
}

func (s *GuichetService) CreerTutorielVideo(ctx context.Context, autoEcoleID string, input models.CreerTutorielVideoInput) (*models.TutorielVideo, error) {
	if strings.TrimSpace(input.Titre) == "" || strings.TrimSpace(input.UrlVideo) == "" {
		return nil, errors.New("le titre et l'URL de la vidéo sont obligatoires")
	}
	return s.guichetRepo.CreerTutorielVideo(ctx, autoEcoleID, input)
}

func (s *GuichetService) SupprimerTutorielVideo(ctx context.Context, autoEcoleID, id string) error {
	return s.guichetRepo.SupprimerTutorielVideo(ctx, autoEcoleID, id)
}

func (s *GuichetService) IncrementerVueTutoriel(ctx context.Context, autoEcoleID, id string) error {
	return s.guichetRepo.IncrementerVueTutoriel(ctx, autoEcoleID, id)
}

func (s *GuichetService) GetSupportsCours(ctx context.Context, autoEcoleID, categorie string) ([]models.SupportCours, error) {
	return s.guichetRepo.GetSupportsCours(ctx, autoEcoleID, categorie)
}

func (s *GuichetService) CreerSupportCours(ctx context.Context, autoEcoleID string, input models.CreerSupportCoursInput) (*models.SupportCours, error) {
	if strings.TrimSpace(input.Titre) == "" {
		return nil, errors.New("le titre du support de cours est obligatoire")
	}
	return s.guichetRepo.CreerSupportCours(ctx, autoEcoleID, input)
}

func (s *GuichetService) SupprimerSupportCours(ctx context.Context, autoEcoleID, id string) error {
	return s.guichetRepo.SupprimerSupportCours(ctx, autoEcoleID, id)
}

func (s *GuichetService) IncrementerTelechargementSupport(ctx context.Context, autoEcoleID, id string) error {
	return s.guichetRepo.IncrementerTelechargementSupport(ctx, autoEcoleID, id)
}

func (s *GuichetService) SynchroniserLotEvenements(
	ctx context.Context,
	autoEcoleID, agenceID string,
	events []models.SyncEventPayload,
) (*models.SyncBatchResponse, error) {
	if autoEcoleID == "" {
		return nil, errors.New("identifiant auto-école manquant")
	}
	return s.guichetRepo.SynchroniserLotEvenements(ctx, autoEcoleID, agenceID, events)
}

func (s *GuichetService) GetQuizFormulaires(ctx context.Context, autoEcoleID, categorie string) ([]models.QuizFormulaire, error) {
	return s.guichetRepo.GetQuizFormulaires(ctx, autoEcoleID, categorie)
}

func (s *GuichetService) CreerQuizFormulaire(ctx context.Context, autoEcoleID string, input models.CreerQuizFormulaireInput) (*models.QuizFormulaire, error) {
	if strings.TrimSpace(input.Titre) == "" {
		return nil, errors.New("le titre du formulaire de quiz est obligatoire")
	}
	if len(input.Questions) == 0 {
		return nil, errors.New("le quiz doit comporter au moins une question")
	}
	return s.guichetRepo.CreerQuizFormulaire(ctx, autoEcoleID, input)
}

func (s *GuichetService) SupprimerQuizFormulaire(ctx context.Context, autoEcoleID, id string) error {
	return s.guichetRepo.SupprimerQuizFormulaire(ctx, autoEcoleID, id)
}

func (s *GuichetService) GetQuizSoumissions(ctx context.Context, autoEcoleID, quizID, candidatID string) ([]models.QuizSoumissionEleve, error) {
	return s.guichetRepo.GetQuizSoumissions(ctx, autoEcoleID, quizID, candidatID)
}

func (s *GuichetService) SoumettreQuizReponse(ctx context.Context, autoEcoleID string, input models.SoumettreQuizReponseInput) (*models.QuizSoumissionEleve, error) {
	if strings.TrimSpace(input.QuizID) == "" {
		return nil, errors.New("l'identifiant du quiz est requis")
	}
	return s.guichetRepo.SoumettreQuizReponse(ctx, autoEcoleID, input)
}

package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"autoschool/backend/internal/config"
	"autoschool/backend/internal/models"
	"autoschool/backend/internal/realtime"
	"autoschool/backend/internal/repository"
)

// AuthService orchestre l'authentification et l'isolation multi-tenant.
type AuthService struct {
	cfg          *config.Config
	schoolRepo   *repository.SchoolRepository
	userRepo     *repository.UserRepository
	sessRepo     *repository.SessionRepository
	notifService *NotificationService
}

func NewAuthService(
	cfg *config.Config,
	schoolRepo *repository.SchoolRepository,
	userRepo *repository.UserRepository,
	sessRepo *repository.SessionRepository,
) *AuthService {
	return &AuthService{
		cfg:        cfg,
		schoolRepo: schoolRepo,
		userRepo:   userRepo,
		sessRepo:   sessRepo,
	}
}

// SetNotificationService injecte le service de notification et d'audit
func (s *AuthService) SetNotificationService(notifService *NotificationService) {
	s.notifService = notifService
}

// GetFrontendURL retourne l'URL frontend configurée via FRONTEND_URL dans .env.
func (s *AuthService) GetFrontendURL() string {
	return s.cfg.FrontendURL
}

// Slugify nettoie une chaîne pour en faire un identifiant URL valide.
func Slugify(text string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	slug := reg.ReplaceAllString(strings.ToLower(text), "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = "ecole-" + uuid.New().String()[:8]
	}
	return slug
}

func normalizeSchoolDevise(raw string) string {
	upper := strings.ToUpper(strings.TrimSpace(raw))
	switch {
	case strings.Contains(upper, "XAF"):
		return "XAF"
	case strings.Contains(upper, "GNF"):
		return "GNF"
	case strings.Contains(upper, "EUR"):
		return "EUR"
	case strings.Contains(upper, "USD"):
		return "USD"
	case len(upper) == 3:
		return upper
	default:
		return "XOF"
	}
}

// RegisterSchool crée une nouvelle auto-école avec son agence siège et son premier compte Directeur.
func (s *AuthService) RegisterSchool(
	ctx context.Context,
	input *models.RegisterSchoolInput,
	clientIP, userAgent string,
) (*models.AuthResponse, error) {
	if strings.TrimSpace(input.NomCommercial) == "" {
		return nil, errors.New("le nom commercial de l'auto-école est requis")
	}

	telDirecteur := strings.TrimSpace(input.TelephoneDirecteur)
	telEcole := strings.TrimSpace(input.TelephoneEcole)
	if telEcole == "" && telDirecteur != "" {
		telEcole = telDirecteur
		input.TelephoneEcole = telDirecteur
	}
	if telDirecteur == "" && telEcole != "" {
		telDirecteur = telEcole
		input.TelephoneDirecteur = telEcole
	}

	if telEcole == "" {
		return nil, errors.New("le numéro de téléphone de l'auto-école est requis")
	}
	if strings.TrimSpace(input.NomDirecteur) == "" || strings.TrimSpace(input.PrenomsDirecteur) == "" {
		return nil, errors.New("le nom et prénom du directeur sont requis")
	}
	if telDirecteur == "" {
		return nil, errors.New("le téléphone du directeur est requis")
	}
	if len(input.MotDePasse) < 6 {
		return nil, errors.New("le mot de passe doit comporter au moins 6 caractères")
	}

	codeUnique := Slugify(strings.TrimSpace(input.CodeUnique))
	if strings.TrimSpace(input.CodeUnique) == "" {
		codeUnique = Slugify(input.NomCommercial)
	}

	existing, err := s.schoolRepo.GetByCodeUnique(ctx, codeUnique)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		codeUnique = fmt.Sprintf("%s-%d", codeUnique, time.Now().Unix()%10000)
		if ex2, _ := s.schoolRepo.GetByCodeUnique(ctx, codeUnique); ex2 != nil {
			codeUnique = fmt.Sprintf("%s-%s", codeUnique, uuid.New().String()[:4])
		}
	}

	pwdHash, err := HashPassword(input.MotDePasse)
	if err != nil {
		return nil, fmt.Errorf("erreur hachage mot de passe: %w", err)
	}

	pays := strings.TrimSpace(input.Pays)
	if pays == "" {
		pays = "Cameroun"
	}

	schoolID := uuid.New().String()
	agenceID := uuid.New().String()
	directorID := uuid.New().String()
	trialExpiration := time.Now().AddDate(0, 0, 30)

	school := &models.AutoEcole{
		ID:                    schoolID,
		CodeUnique:            codeUnique,
		NomCommercial:         input.NomCommercial,
		NumeroAgrement:        input.NumeroAgrement,
		Telephone:             input.TelephoneEcole,
		Email:                 input.EmailEcole,
		Devise:                normalizeSchoolDevise(input.Devise),
		Pays:                  pays,
		EstActif:              true,
		MaxMachines:           3,
		TypeLicence:           "TRIAL",
		DateExpirationLicence: &trialExpiration,
		ModulesActifs:         `["guichet","caisse","conduite","examens","flotte","archivage_cloud"]`,
	}

	nomAgence := input.NomAgence
	if strings.TrimSpace(nomAgence) == "" {
		nomAgence = "Agence Principale"
	}
	adresseAgence := input.AdresseAgence
	if strings.TrimSpace(adresseAgence) == "" {
		adresseAgence = "Siège social"
	}
	telAgence := input.TelephoneAgence
	if strings.TrimSpace(telAgence) == "" {
		telAgence = input.TelephoneEcole
	}

	cityPrefix := "HQ"
	cleanVille := strings.ToUpper(strings.TrimSpace(input.VilleAgence))
	if len(cleanVille) >= 3 {
		cityPrefix = cleanVille[:3]
	} else if len(school.CodeUnique) >= 3 {
		cityPrefix = strings.ToUpper(school.CodeUnique[:3])
	}
	cleanSlug := strings.ToUpper(strings.ReplaceAll(school.CodeUnique, "-", ""))
	if len(cleanSlug) > 4 {
		cleanSlug = cleanSlug[:4]
	}

	agence := &models.Agence{
		ID:          agenceID,
		AutoEcoleID: schoolID,
		CodeAgence:  fmt.Sprintf("AG-%s-%s-01", cityPrefix, cleanSlug),
		Nom:         nomAgence,
		Adresse:     adresseAgence,
		Telephone:   telAgence,
		Ville:       input.VilleAgence,
		Commune:     input.CommuneAgence,
		EstSiege:    true,
	}

	director := &models.User{
		ID:             directorID,
		AutoEcoleID:    schoolID,
		AgenceID:       &agenceID,
		Nom:            input.NomDirecteur,
		Prenoms:        input.PrenomsDirecteur,
		Telephone:      input.TelephoneDirecteur,
		Email:          input.EmailDirecteur,
		MotDePasseHash: &pwdHash,
		CodePinManager: input.CodePinManager,
		Role:           models.RoleDirecteur,
		Statut:         models.StatusActif,
	}

	if err := s.schoolRepo.CreateSchoolWithAgencyAndAdmin(ctx, school, agence, director); err != nil {
		return nil, fmt.Errorf("échec de l'enregistrement de l'auto-école: %w", err)
	}

	if s.notifService != nil {
		s.notifySchoolRegistration(ctx, school, agence, director, input.TelephoneDirecteur, clientIP)
	}

	return s.createSessionAndResponse(ctx, director, school, agence, clientIP, userAgent)
}

func (s *AuthService) notifySchoolRegistration(
	ctx context.Context,
	school *models.AutoEcole,
	agence *models.Agence,
	director *models.User,
	telDirecteur, clientIP string,
) {
	roleSuperAdmin := "super_admin"
	schoolIDVal := school.ID
	dirName := fmt.Sprintf("%s %s", director.Prenoms, director.Nom)

	_, _ = s.notifService.NotifyAndBroadcast(ctx, &models.CreateNotificationInput{
		RoleCible: &roleSuperAdmin,
		Titre:     "🏢 Nouvelle Auto-École enregistrée",
		Message:   fmt.Sprintf("L'établissement %s (%s, %s) vient de s'inscrire. Directeur : %s (Tél: %s).", school.NomCommercial, agence.Ville, school.Pays, dirName, telDirecteur),
		Type:      "success",
		Module:    "super_admin",
		ActionURL: &schoolIDVal,
		Metadata: map[string]interface{}{
			"school_id":       school.ID,
			"nom_commercial":  school.NomCommercial,
			"code_unique":     school.CodeUnique,
			"ville":           agence.Ville,
			"pays":            school.Pays,
			"directeur_nom":   dirName,
			"telephone":       telDirecteur,
			"numero_agrement": school.NumeroAgrement,
		},
	})

	_ = s.notifService.LogAudit(ctx, &models.SystemAuditLog{
		Action:      "INSCRIPTION_AUTO_ECOLE",
		CibleType:   "auto_ecole",
		CibleID:     school.ID,
		Details:     fmt.Sprintf("Inscription officielle de l'auto-école %s (%s). Directeur: %s, Tél: %s, Ville: %s", school.NomCommercial, school.CodeUnique, dirName, telDirecteur, agence.Ville),
		IPAddress:   clientIP,
		AutoEcoleID: &schoolIDVal,
		UserID:      &director.ID,
		AuteurNom:   &dirName,
		Role:        &director.Role,
	})

	realtime.BroadcastToSuperAdmin("SCHOOL_REGISTERED", map[string]interface{}{
		"id":             school.ID,
		"nom_commercial": school.NomCommercial,
		"code_unique":    school.CodeUnique,
		"ville":          agence.Ville,
		"pays":           school.Pays,
		"directeur":      dirName,
		"telephone":      telDirecteur,
		"created_at":     time.Now().UTC().Format(time.RFC3339),
	})

	roleDirecteur := models.RoleDirecteur
	_, _ = s.notifService.NotifyAndBroadcast(ctx, &models.CreateNotificationInput{
		AutoEcoleID: &schoolIDVal,
		RoleCible:   &roleDirecteur,
		Titre:       "🎉 Bienvenue sur ElevenDrive !",
		Message:     fmt.Sprintf("Félicitations pour l'ouverture de %s ! Votre période d'essai gratuit de 30 jours (Plan Eleven Pro) est active. Toutes les fonctionnalités sont disponibles sans limite.", school.NomCommercial),
		Type:        "success",
		Module:      "abonnement",
		ActionURL:   &schoolIDVal,
		Metadata: map[string]interface{}{
			"school_id":               school.ID,
			"nom_commercial":          school.NomCommercial,
			"jours_gratuits":          30,
			"type_licence":            "TRIAL",
			"date_expiration_licence": school.DateExpirationLicence.Format(time.RFC3339),
		},
	})
}

// Login authentifie un utilisateur et émet un couple AccessToken / RefreshToken.
func (s *AuthService) Login(
	ctx context.Context,
	input *models.LoginInput,
	clientIP, userAgent string,
) (*models.AuthResponse, error) {
	identifiant := strings.TrimSpace(input.Identifiant)
	motDePasse := input.MotDePasse
	if identifiant == "" || motDePasse == "" {
		return nil, errors.New("l'identifiant et le mot de passe sont obligatoires")
	}

	targetUser, targetSchool, err := s.resolveLoginTarget(ctx, identifiant, motDePasse, input.SchoolSlug)
	if err != nil {
		return nil, err
	}
	if targetUser == nil || targetUser.MotDePasseHash == nil || !CheckPasswordHash(motDePasse, *targetUser.MotDePasseHash) {
		return nil, errors.New("identifiants invalides")
	}
	if targetUser.Statut != models.StatusActif {
		return nil, fmt.Errorf("ce compte est actuellement %s, veuillez contacter l'administration", strings.ToLower(targetUser.Statut))
	}

	if targetSchool != nil && targetUser.Role != models.RoleSuperAdmin {
		updatedSchool, err := s.enforceSchoolLicenceAndMachine(ctx, targetUser, targetSchool, input)
		if err != nil {
			return nil, err
		}
		targetSchool = updatedSchool
	}

	var targetAgence *models.Agence
	if targetUser.AgenceID != nil {
		targetAgence, _ = s.schoolRepo.GetAgenceByID(ctx, *targetUser.AgenceID)
	}

	if s.notifService != nil && targetSchool != nil {
		s.notifService.CheckSubscriptionAlert(ctx, targetSchool.ID)
	}

	return s.createSessionAndResponse(ctx, targetUser, targetSchool, targetAgence, clientIP, userAgent)
}

func (s *AuthService) resolveLoginTarget(ctx context.Context, identifiant, motDePasse, schoolSlug string) (*models.User, *models.AutoEcole, error) {
	if schoolSlug != "" {
		school, err := s.schoolRepo.GetByCodeUnique(ctx, schoolSlug)
		if err != nil {
			return nil, nil, err
		}
		if school != nil {
			user, err := s.userRepo.GetByIdentifierInSchool(ctx, school.ID, identifiant)
			return user, school, err
		}
	}

	users, err := s.userRepo.GetByIdentifierGlobal(ctx, identifiant)
	if err != nil {
		return nil, nil, err
	}
	if len(users) == 0 {
		return nil, nil, errors.New("identifiants invalides")
	}

	var matchedUser *models.User
	if len(users) == 1 {
		matchedUser = users[0]
	} else {
		for _, u := range users {
			if u.MotDePasseHash != nil && CheckPasswordHash(motDePasse, *u.MotDePasseHash) {
				matchedUser = u
				break
			}
		}
		if matchedUser == nil {
			return nil, nil, errors.New("identifiants invalides")
		}
	}

	school, err := s.schoolRepo.GetByID(ctx, matchedUser.AutoEcoleID)
	if err != nil {
		return nil, nil, err
	}
	return matchedUser, school, nil
}

func (s *AuthService) enforceSchoolLicenceAndMachine(
	ctx context.Context,
	user *models.User,
	school *models.AutoEcole,
	input *models.LoginInput,
) (*models.AutoEcole, error) {
	if !school.EstActif {
		return nil, fmt.Errorf("l'établissement '%s' est actuellement suspendu par le Super-Administrateur", school.NomCommercial)
	}
	if school.DateExpirationLicence != nil && time.Now().After(*school.DateExpirationLicence) {
		if user.Role != models.RoleDirecteur && user.Role != models.RoleSecretaire {
			return nil, fmt.Errorf("la période d'essai de l'établissement '%s' a expiré le %s. Veuillez contacter votre direction pour son renouvellement", school.NomCommercial, school.DateExpirationLicence.Format("02/01/2006"))
		}
	}

	if input.MachineFingerprint != "" {
		if err := s.schoolRepo.CheckAndRegisterMachine(ctx, school.ID, input.MachineFingerprint, input.MachineName, input.OS, school.MaxMachines); err != nil {
			return nil, err
		}
		if refreshed, errRefresh := s.schoolRepo.GetByID(ctx, school.ID); errRefresh == nil && refreshed != nil {
			school = refreshed
		}
		realtime.Broadcast(school.ID, "SCHOOL_LICENCE_UPDATED", map[string]interface{}{
			"school_id":        school.ID,
			"machines_actives": school.MachinesActives,
			"max_machines":     school.MaxMachines,
		})
		realtime.BroadcastToSuperAdmin("SUPER_ADMIN_REFRESH", map[string]interface{}{
			"action":           "MACHINE_CONNECTED",
			"school_id":        school.ID,
			"machines_actives": school.MachinesActives,
		})
	}
	return school, nil
}

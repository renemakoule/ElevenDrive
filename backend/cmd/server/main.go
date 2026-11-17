package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"autoschool/backend/internal/config"
	"autoschool/backend/internal/database"
	"autoschool/backend/internal/handler"
	"autoschool/backend/internal/middleware"
	"autoschool/backend/internal/models"
	"autoschool/backend/internal/realtime"
	"autoschool/backend/internal/repository"
	"autoschool/backend/internal/service"
)

func main() {
	log.Println("==================================================================")
	log.Println(" AutoSchool Backend API - Serveur d'Authentification Multi-Tenant ")
	log.Println("==================================================================")

	// 1. Chargement de la configuration (strict sans données par défaut ni mock)
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Erreur de configuration: %v\nAssurez-vous de définir les variables requises (DATABASE_URL, JWT_SECRET).", err)
	}

	// 2. Connexion à la base de données PostgreSQL (Neon Serverless)
	log.Println("🔄 Connexion au pool PostgreSQL...")
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ Impossible de se connecter à PostgreSQL: %v", err)
	}
	defer db.Close()
	log.Println("✅ Connexion PostgreSQL établie avec succès.")

	// Auto-migration idempotente du schéma Super-Admin (licences, machines, quotas, audit)
	_ = db.AutoMigrateSuperAdminSchema(context.Background())

	// 3. Initialisation des dépôts (Repositories)
	schoolRepo := repository.NewSchoolRepository(db)
	userRepo := repository.NewUserRepository(db)
	sessRepo := repository.NewSessionRepository(db)
	guichetRepo := repository.NewGuichetRepository(db)
	flotteRepo := repository.NewFlotteRepository(db)
	epreuveRepo := repository.NewEpreuveRepository(db)
	superAdminRepo := repository.NewSuperAdminRepository(db)
	commRepo := repository.NewCommunicationRepository(db)
	notifRepo := repository.NewNotificationRepository(db)
	_ = epreuveRepo.EnsureTableExists(context.Background())
	_ = guichetRepo.EnsureQuizTablesExist(context.Background())
	guichetRepo.EnsureSchemaCompat(context.Background())

	// 4. Initialisation des services métiers
	notifService := service.NewNotificationService(notifRepo)
	notifService.SetSchoolRepository(schoolRepo)
	authService := service.NewAuthService(cfg, schoolRepo, userRepo, sessRepo)
	authService.SetNotificationService(notifService)
	superAdminService := service.NewSuperAdminService(cfg, superAdminRepo, schoolRepo, userRepo, sessRepo)
	superAdminService.SetNotificationService(notifService)
	commService := service.NewCommunicationService(commRepo)
	commService.SetNotificationService(notifService)
	guichetService := service.NewGuichetService(guichetRepo, userRepo)
	flotteService := service.NewFlotteService(flotteRepo)
	epreuveService := service.NewEpreuveService(epreuveRepo)
	translationService := service.NewDeepLTranslationService(cfg.DeepLAPIKey)
	r2Service := service.NewR2StorageService(cfg)

	// 5. Initialisation des contrôleurs HTTP
	authHandler := handler.NewAuthHandler(authService)
	superAdminHandler := handler.NewSuperAdminHandler(superAdminService)
	commHandler := handler.NewCommunicationHandler(commService)
	notifHandler := handler.NewNotificationHandler(notifService)
	guichetHandler := handler.NewGuichetHandler(guichetService)
	flotteHandler := handler.NewFlotteHandler(flotteService)
	epreuveHandler := handler.NewEpreuveHandler(epreuveService)
	translationHandler := handler.NewTranslationHandler(translationService)
	archivageHandler := handler.NewArchivageHandler(r2Service)

	// 6. Définition des routes HTTP & Hub WebSocket Temps Réel
	mux := http.NewServeMux()
	wsHub := realtime.NewHub(cfg.JWTSecret)
	mux.HandleFunc("/api/ws", wsHub.ServeWS)
	mux.HandleFunc("/ws", wsHub.ServeWS)

	// Route de santé et routes d'accueil
	rootHandler := func(w http.ResponseWriter, r *http.Request) {
		cleanPath := strings.TrimRight(r.URL.Path, "/")
		if cleanPath != "" && cleanPath != "/api" && cleanPath != "/api/auth" {
			http.NotFound(w, r)
			return
		}
		handler.Success(w, http.StatusOK, map[string]interface{}{
			"service":   "AutoSchool ElevenDrive API",
			"status":    "ONLINE",
			"version":   "1.1.0",
			"database":  "CONNECTED",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	}
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/api", rootHandler)
	mux.HandleFunc("/api/", rootHandler)
	mux.HandleFunc("/api/auth", rootHandler)
	mux.HandleFunc("/api/auth/", rootHandler)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		handler.Success(w, http.StatusOK, map[string]interface{}{
			"status":    "UP",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"database":  "CONNECTED",
		})
	})

	// Middlewares de protection par JWT
	requireAuth := middleware.RequireAuth(cfg.JWTSecret)
	requireManager := middleware.RequireRole(models.RoleDirecteur, models.RoleSecretaire, models.RoleCaissier)

	// Routes Super-Admin Console
	mux.HandleFunc("/api/super-admin/login", superAdminHandler.LoginSuperAdmin)
	mux.Handle("/api/super-admin/dashboard", requireAuth(http.HandlerFunc(superAdminHandler.GetDashboard)))
	mux.Handle("/api/super-admin/users", requireAuth(http.HandlerFunc(superAdminHandler.GetUsers)))
	mux.Handle("/api/super-admin/candidats", requireAuth(http.HandlerFunc(superAdminHandler.GetCandidats)))
	mux.Handle("/api/super-admin/schools/status", requireAuth(http.HandlerFunc(superAdminHandler.ToggleSchoolStatus)))
	mux.Handle("/api/super-admin/schools/update", requireAuth(http.HandlerFunc(superAdminHandler.UpdateSchool)))
	mux.Handle("/api/super-admin/schools/licence", requireAuth(http.HandlerFunc(superAdminHandler.UpdateLicence)))
	mux.Handle("/api/super-admin/schools/modules", requireAuth(http.HandlerFunc(superAdminHandler.UpdateModules)))
	mux.Handle("/api/super-admin/schools/delete", requireAuth(http.HandlerFunc(superAdminHandler.DeleteSchool)))
	mux.Handle("/api/super-admin/audit", requireAuth(http.HandlerFunc(superAdminHandler.GetAuditLogs)))

	// Routes Super-Admin Diffusion & Retours Réseau
	mux.Handle("/api/super-admin/alerts", requireAuth(http.HandlerFunc(commHandler.HandleSuperAdminAlerts)))
	mux.Handle("/api/super-admin/alerts/toggle", requireAuth(http.HandlerFunc(commHandler.ToggleAlert)))
	mux.Handle("/api/super-admin/announcements", requireAuth(http.HandlerFunc(commHandler.HandleSuperAdminAnnouncements)))
	mux.Handle("/api/super-admin/feedbacks", requireAuth(http.HandlerFunc(commHandler.HandleSuperAdminFeedbacks)))
	mux.Handle("/api/super-admin/feedbacks/reply", requireAuth(http.HandlerFunc(commHandler.ReplyFeedback)))
	mux.Handle("/api/super-admin/feedbacks/status", requireAuth(http.HandlerFunc(commHandler.UpdateFeedbackStatus)))
	mux.Handle("/api/super-admin/support-tickets", requireAuth(http.HandlerFunc(commHandler.GetSupportTickets)))

	// Routes Publiques / Auto-Écoles : Communication & Assistance
	mux.HandleFunc("/api/support/ticket", commHandler.SubmitSupportTicket)
	mux.Handle("/api/communication/alerts", requireAuth(http.HandlerFunc(commHandler.GetClientAlerts)))
	mux.Handle("/api/communication/announcements", requireAuth(http.HandlerFunc(commHandler.GetClientAnnouncements)))
	mux.Handle("/api/communication/feedbacks", requireAuth(http.HandlerFunc(commHandler.SubmitClientFeedback)))

	// Routes Système : Notifications Multi-Tenant & Audit
	mux.Handle("/api/notifications", requireAuth(http.HandlerFunc(notifHandler.GetNotifications)))
	mux.Handle("/api/notifications/", requireAuth(http.HandlerFunc(notifHandler.GetNotifications)))
	mux.Handle("/api/notifications/read", requireAuth(http.HandlerFunc(notifHandler.MarkAsRead)))
	mux.Handle("/api/notifications/read-all", requireAuth(http.HandlerFunc(notifHandler.MarkAllAsRead)))
	mux.Handle("/api/notifications/audit", requireAuth(http.HandlerFunc(notifHandler.GetAuditLogs)))

	// Routes publiques d'authentification (compatibilité /api/auth et /api direct)
	mux.HandleFunc("/api/auth/register-school", authHandler.RegisterSchool)
	mux.HandleFunc("/api/register-school", authHandler.RegisterSchool)
	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.HandleFunc("/api/login", authHandler.Login)
	mux.HandleFunc("/api/auth/oauth", authHandler.OAuth)
	mux.HandleFunc("/api/oauth", authHandler.OAuth)
	mux.HandleFunc("/api/auth/oauth/select-school", authHandler.OAuthSelectSchool)
	mux.HandleFunc("/api/oauth/select-school", authHandler.OAuthSelectSchool)
	mux.HandleFunc("/api/auth/oauth/session", authHandler.GetOAuthSession)
	mux.HandleFunc("/api/oauth/session", authHandler.GetOAuthSession)
	mux.HandleFunc("/api/auth/google", authHandler.GoogleAuthRedirect)
	mux.HandleFunc("/api/auth/google/callback", authHandler.GoogleAuthCallback)
	mux.HandleFunc("/api/api/auth/google", authHandler.GoogleAuthRedirect)
	mux.HandleFunc("/api/api/auth/google/callback", authHandler.GoogleAuthCallback)
	mux.HandleFunc("/api/auth/facebook", authHandler.FacebookAuthRedirect)
	mux.HandleFunc("/api/auth/facebook/callback", authHandler.FacebookAuthCallback)
	mux.HandleFunc("/api/api/auth/facebook", authHandler.FacebookAuthRedirect)
	mux.HandleFunc("/api/api/auth/facebook/callback", authHandler.FacebookAuthCallback)
	mux.HandleFunc("/api/auth/linkedin", authHandler.LinkedInAuthRedirect)
	mux.HandleFunc("/api/auth/linkedin/callback", authHandler.LinkedInAuthCallback)
	mux.HandleFunc("/api/api/auth/linkedin", authHandler.LinkedInAuthRedirect)
	mux.HandleFunc("/api/api/auth/linkedin/callback", authHandler.LinkedInAuthCallback)
	mux.HandleFunc("/api/auth/apple", authHandler.AppleAuthRedirect)
	mux.HandleFunc("/api/auth/apple/callback", authHandler.AppleAuthCallback)
	mux.HandleFunc("/api/api/auth/apple", authHandler.AppleAuthRedirect)
	mux.HandleFunc("/api/api/auth/apple/callback", authHandler.AppleAuthCallback)
	mux.HandleFunc("/api/auth/refresh", authHandler.Refresh)
	mux.HandleFunc("/api/refresh", authHandler.Refresh)
	mux.HandleFunc("/api/auth/logout", authHandler.Logout)
	mux.HandleFunc("/api/logout", authHandler.Logout)
	mux.HandleFunc("/api/auth/verify-pin", authHandler.VerifyPin)
	mux.HandleFunc("/api/verify-pin", authHandler.VerifyPin)

	schoolDispatchHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			authHandler.GetSchoolInfo(w, r)
			return
		}
		if r.Method == http.MethodPut || r.Method == http.MethodPatch {
			requireAuth(http.HandlerFunc(authHandler.UpdateSchool)).ServeHTTP(w, r)
			return
		}
		handler.Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
	})
	mux.Handle("/api/auth/school", schoolDispatchHandler)
	mux.Handle("/api/school", schoolDispatchHandler)

	mux.Handle("/api/auth/me", requireAuth(http.HandlerFunc(authHandler.Me)))
	mux.Handle("/api/me", requireAuth(http.HandlerFunc(authHandler.Me)))

	// Gestion des utilisateurs / collaborateurs (Lecture ouverte aux membres authentifiés de l'auto-école, écriture réservée au management)
	userDispatchHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			authHandler.GetUsers(w, r)
		case http.MethodPost:
			requireManager(http.HandlerFunc(authHandler.RegisterUser)).ServeHTTP(w, r)
		case http.MethodPut, http.MethodPatch:
			requireManager(http.HandlerFunc(authHandler.UpdateUser)).ServeHTTP(w, r)
		case http.MethodDelete:
			requireManager(http.HandlerFunc(authHandler.DeleteUser)).ServeHTTP(w, r)
		default:
			handler.Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		}
	})

	mux.Handle("/api/auth/users", requireAuth(userDispatchHandler))
	mux.Handle("/api/auth/users/", requireAuth(userDispatchHandler))
	mux.Handle("/api/users", requireAuth(userDispatchHandler))
	mux.Handle("/api/users/", requireAuth(userDispatchHandler))

	// Gestion des agences / succursales (Lecture ouverte à tout utilisateur authentifié de l'école, écriture réservée au management)
	agenceDispatchHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			authHandler.GetAgences(w, r)
		case http.MethodPost:
			requireManager(http.HandlerFunc(authHandler.CreateAgence)).ServeHTTP(w, r)
		case http.MethodPut, http.MethodPatch:
			requireManager(http.HandlerFunc(authHandler.UpdateAgence)).ServeHTTP(w, r)
		case http.MethodDelete:
			requireManager(http.HandlerFunc(authHandler.DeleteAgence)).ServeHTTP(w, r)
		default:
			handler.Error(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		}
	})

	mux.Handle("/api/auth/agences", requireAuth(agenceDispatchHandler))
	mux.Handle("/api/auth/agences/", requireAuth(agenceDispatchHandler))
	mux.Handle("/api/agences", requireAuth(agenceDispatchHandler))
	mux.Handle("/api/agences/", requireAuth(agenceDispatchHandler))

	mux.Handle("/api/auth/register-user", requireAuth(requireManager(http.HandlerFunc(authHandler.RegisterUser))))
	mux.Handle("/api/register-user", requireAuth(requireManager(http.HandlerFunc(authHandler.RegisterUser))))

	requireDirecteur := middleware.RequireRole(models.RoleDirecteur)

	// Routes Guichet, Caisse & Pédagogie (isolées par auto_ecole_id)
	mux.Handle("/api/guichet/stats", requireAuth(http.HandlerFunc(guichetHandler.GetStats)))

	// Catalogue des Formules & Gestion des Tarifs :
	// Fixation initiale : Directeur ou Secrétaire. Modification d'un prix existant : STRICTEMENT Directeur.
	mux.Handle("/api/formules", requireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			guichetHandler.GetFormules(w, r)
		case http.MethodPost:
			guichetHandler.CreerFormule(w, r)
		case http.MethodPatch, http.MethodPut:
			guichetHandler.ModifierPrixFormule(w, r)
		default:
			handler.Error(w, http.StatusMethodNotAllowed, "Méthode HTTP non autorisée")
		}
	})))
	mux.Handle("/api/formules/prix", requireAuth(requireDirecteur(http.HandlerFunc(guichetHandler.ModifierPrixFormule))))
	mux.Handle("/api/candidats", requireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			guichetHandler.CreerCandidat(w, r)
		} else {
			guichetHandler.SearchCandidats(w, r)
		}
	})))
	mux.Handle("/api/candidats/import", requireAuth(http.HandlerFunc(guichetHandler.ImporterCandidatsBatch)))
	mux.Handle("/api/candidats/dossier", requireAuth(http.HandlerFunc(guichetHandler.GetCandidatDossier)))
	mux.Handle("/api/candidats/", requireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/photo") {
			guichetHandler.UpdateCandidatPhoto(w, r)
			return
		}
		guichetHandler.GetCandidatDossier(w, r)
	})))
	mux.Handle("/api/paiements", requireAuth(http.HandlerFunc(guichetHandler.EnregistrerPaiement)))
	mux.Handle("/api/paiements/recu", requireAuth(http.HandlerFunc(guichetHandler.ReimprimerRecu)))
	mux.Handle("/api/caisse/session-detail", requireAuth(http.HandlerFunc(guichetHandler.GetSessionCaisse)))
	mux.Handle("/api/caisse/cloturer", requireAuth(requireManager(http.HandlerFunc(guichetHandler.CloturerSessionCaisse))))
	mux.Handle("/api/pedagogie/seance-code", requireAuth(http.HandlerFunc(guichetHandler.EnregistrerSeanceCode)))
	mux.Handle("/api/pedagogie/seances-code", requireAuth(http.HandlerFunc(guichetHandler.HandleSeancesCode)))
	mux.Handle("/api/pedagogie/tutoriels", requireAuth(http.HandlerFunc(guichetHandler.HandleTutorielsVideos)))
	mux.Handle("/api/pedagogie/tutoriels/vue", requireAuth(http.HandlerFunc(guichetHandler.IncrementerVueTutoriel)))
	mux.Handle("/api/pedagogie/supports-cours", requireAuth(http.HandlerFunc(guichetHandler.HandleSupportsCours)))
	mux.Handle("/api/pedagogie/supports-cours/telecharger", requireAuth(http.HandlerFunc(guichetHandler.IncrementerTelechargementSupport)))
	mux.Handle("/api/pedagogie/planning-conduite", requireAuth(http.HandlerFunc(guichetHandler.GetPlanningConduite)))
	mux.Handle("/api/pedagogie/seance-conduite/valider", requireAuth(http.HandlerFunc(guichetHandler.ValiderSeanceConduite)))
	mux.Handle("/api/pedagogie/reserver-conduite", requireAuth(http.HandlerFunc(guichetHandler.ReserverSeanceConduite)))
	mux.Handle("/api/pedagogie/quiz-resultat", requireAuth(http.HandlerFunc(guichetHandler.EnregistrerQuizResultat)))
	mux.Handle("/api/pedagogie/quiz-forms", requireAuth(http.HandlerFunc(guichetHandler.HandleQuizFormulaires)))
	mux.Handle("/api/pedagogie/quiz-submissions", requireAuth(http.HandlerFunc(guichetHandler.HandleQuizSoumissions)))
	mux.Handle("/api/manager/dashboard", requireAuth(requireManager(http.HandlerFunc(guichetHandler.GetTableauGerant))))
	mux.Handle("/api/manager/derogation", requireAuth(requireManager(http.HandlerFunc(guichetHandler.AccorderDerogation))))

	// Examens d'État, Bordereaux & Procès-Verbaux (Réservés au Management)
	mux.Handle("/api/examens/session-active", requireAuth(http.HandlerFunc(guichetHandler.GetSessionExamenActive)))
	mux.Handle("/api/examens/sessions", requireAuth(requireManager(http.HandlerFunc(guichetHandler.CreerOuModifierSessionExamen))))
	mux.Handle("/api/examens/sessions/supprimer", requireAuth(requireManager(http.HandlerFunc(guichetHandler.SupprimerSessionExamen))))
	mux.Handle("/api/examens/forcer-candidat", requireAuth(requireManager(http.HandlerFunc(guichetHandler.ForcerCandidatExamen))))
	mux.Handle("/api/examens/resultats", requireAuth(requireManager(http.HandlerFunc(guichetHandler.EnregistrerResultatsExamen))))

	// Gestion de la Flotte Automobile & Alertes Techniques
	requireStaff := middleware.RequireRole(models.RoleDirecteur, models.RoleSecretaire, models.RoleCaissier, models.RoleMoniteur)
	mux.Handle("/api/flotte/vehicules", requireAuth(http.HandlerFunc(flotteHandler.HandleVehicules)))
	mux.Handle("/api/flotte/vehicules/statut", requireAuth(requireStaff(http.HandlerFunc(flotteHandler.UpdateStatut))))

	// Banque d'Épreuves Pédagogiques & Annales PDF (Multi-Tenant & Filtrage Catégorie)
	mux.Handle("/api/epreuves", requireAuth(http.HandlerFunc(epreuveHandler.HandleEpreuves)))
	mux.Handle("/api/epreuves/telecharger", requireAuth(http.HandlerFunc(epreuveHandler.HandleTelechargement)))

	// Moteur de Synchronisation Offline-First (Ingestion par lots vers Neon PostgreSQL)
	mux.Handle("/api/sync/batch", requireAuth(http.HandlerFunc(guichetHandler.SynchroniserLotEvenements)))
	// Ping mis en cache 5s : évite que des dizaines de pings parallèles saturent le pool de connexions.
	var (
		pingMu      sync.Mutex
		pingLastOK  time.Time
		pingLastErr error
	)
	mux.Handle("/api/sync/pull", requireAuth(http.HandlerFunc(guichetHandler.PullSyncSnapshot)))
	mux.HandleFunc("/api/sync/ping", func(w http.ResponseWriter, r *http.Request) {
		pingMu.Lock()
		if time.Since(pingLastOK) > 5*time.Second || pingLastErr != nil {
			ctxPing, cancelPing := context.WithTimeout(r.Context(), 3*time.Second)
			pingLastErr = db.PingContext(ctxPing)
			cancelPing()
			if pingLastErr == nil {
				pingLastOK = time.Now()
			}
		}
		err := pingLastErr
		pingMu.Unlock()
		if err != nil {
			handler.Error(w, http.StatusServiceUnavailable, "Cluster PostgreSQL Neon non joignable: "+err.Error())
			return
		}
		handler.Success(w, http.StatusOK, map[string]interface{}{
			"status":    "CONNECTED",
			"cluster":   "Neon Serverless PostgreSQL",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Service Universel de Traduction DeepL API (Accessible pour la vitrine publique Landing et le dashboard)
	mux.HandleFunc("/api/traduire", translationHandler.Translate)
	mux.HandleFunc("/api/traduire/", translationHandler.Translate)
	mux.HandleFunc("/api/traduire/statut", translationHandler.Statut)

	// Centre d'Archivage Cloudflare R2 (Bucket : maryam-logistic-bucket)
	mux.Handle("/api/archivage/cloud/statut", requireAuth(http.HandlerFunc(archivageHandler.GetCloudStatut)))
	mux.Handle("/api/archivage/cloud/upload", requireAuth(http.HandlerFunc(archivageHandler.UploadDocumentToCloud)))
	mux.Handle("/api/archivage/cloud/dossier", requireAuth(http.HandlerFunc(archivageHandler.ArchiveDossierToCloud)))

	// 7. Enveloppe globale des middlewares : Logger + CORS + Sécurité/Taille + Rate Limiting anti-DoS
	finalHandler := middleware.Logger(
		middleware.CORS(cfg.AllowedOrigins)(
			middleware.SecurityAndBodyLimit(
				middleware.RateLimit(mux),
			),
		),
	)

	// 8. Démarrage du serveur HTTP avec timeouts anti-Slowloris et arrêt gracieux (Graceful Shutdown)
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           finalHandler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 Mo max pour les en-têtes HTTP
	}

	go func() {
		log.Printf("🚀 Serveur démarré et à l'écoute sur http://0.0.0.0:%s\n", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Erreur serveur HTTP: %v", err)
		}
	}()

	// Capture des signaux système d'interruption
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Signal d'arrêt reçu, fermeture gracieuse du serveur...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Erreur lors de l'arrêt du serveur: %v", err)
	}

	log.Println("✅ Serveur AutoSchool arrêté avec succès.")
}

# AutoSchool - Backend d'Authentification Multi-Tenant (Go)

Ce dossier contient l'API Cloud et le service d'authentification natif en **Go (Golang)** pour la plateforme de gestion d'auto-écoles **AutoSchool**.

---

## 🔒 Principes d'Architecture & Sécurité

1. **Zéro Fausse Donnée / Zéro Mock** :
   - Aucun compte par défaut, mot de passe en dur ou donnée fictive n'est stocké dans le backend.
   - Les auto-écoles et leurs comptes directeurs sont créés dynamiquement via l'endpoint d'onboarding (`POST /api/auth/register-school`).
   - Le serveur refuse de démarrer si `DATABASE_URL` ou `JWT_SECRET` sont absents.

2. **Multi-Tenancy avec Row Level Security (RLS)** :
   - Chaque requête authentifiée isole les données au niveau PostgreSQL (`auto_ecole_id`).
   - Support du cloisonnement strict entre écoles et agences.

3. **Sécurité Cryptographique Avancée** :
   - Hachage des mots de passe avec **Bcrypt** (coût de travail 12).
   - Jetons d'accès **JWT signés HMAC-SHA256** (courte durée, contenant l'identité et le tenant).
   - Jetons de rafraîchissement (**Refresh Tokens**) cryptographiques de 32 octets aléatoires, hachés en **SHA-256** en base de données avec **rotation automatique** à chaque renouvellement.
   - Traçabilité des sessions actives (IP, User-Agent, expiration, révocation).
   - Code PIN manager pour les dérogations rapides au guichet physique.

---

## 📁 Structure du Projet

```
backend/
├── cmd/
│   └── server/
│       └── main.go                  # Point d'entrée du serveur HTTP & arrêt gracieux
├── internal/
│   ├── config/
│   │   └── config.go                # Chargement et validation stricte de l'environnement
│   ├── database/
│   │   └── db.go                    # Pool de connexions PostgreSQL & support RLS
│   ├── handler/
│   │   ├── auth_handler.go          # Contrôleurs HTTP d'authentification et onboarding
│   │   └── response.go              # Formateur de réponses JSON standardisées
│   ├── middleware/
│   │   ├── auth.go                  # Validation Bearer JWT & contrôle d'accès RBAC
│   │   ├── cors.go                  # Gestion des origines autorisées & pré-requêtes OPTIONS
│   │   └── logger.go                # Journalisation HTTP structurée
│   ├── models/
│   │   ├── user.go                  # Modèles utilisateurs, rôles et DTOs
│   │   ├── school.go                # Modèles auto-écoles, agences et onboarding
│   │   └── session.go               # Modèles de sessions et claims JWT
│   ├── repository/
│   │   ├── user_repository.go       # Requêtes PostgreSQL utilisateurs
│   │   ├── school_repository.go     # Requêtes PostgreSQL auto-écoles et agences
│   │   └── session_repository.go    # Requêtes PostgreSQL sessions et révocations
│   └── service/
│       ├── auth_service.go          # Logique métier d'authentification
│       ├── password.go              # Hachage et vérification de mots de passe
│       └── token.go                 # Signature JWT et tokens cryptographiques
├── database/
│   ├── neon_postgres_schema.sql     # Schéma complet PostgreSQL (DDL + RLS)
│   ├── migrations/
│   │   └── 001_auth_sessions.sql    # Migration de la table auth_sessions
│   └── README_NEON.md               # Guide d'exécution sur Neon.tech
├── .env.example                     # Modèle de configuration sans identifiant factice
└── go.mod                           # Définition du module Go
```

---

## 🚀 Démarrage Rapide

### 1. Prérequis
- **Go 1.22+** installé
- Instance PostgreSQL (Neon.tech recommandé) avec le schéma exécuté :
  - `backend/database/neon_postgres_schema.sql`
  - `backend/database/migrations/001_auth_sessions.sql`

### 2. Configuration de l'Environnement
Copiez le modèle et renseignez vos véritables paramètres :
```bash
cp .env.example .env
```

Éditez `.env` :
```env
PORT=8080
DATABASE_URL=postgresql://user:password@endpoint.neon.tech/neondb?sslmode=require
JWT_SECRET=votre_cle_secrete_aleatoire_tres_longue_d_au_moins_32_caracteres
JWT_EXPIRATION_HOURS=24
REFRESH_EXPIRATION_DAYS=30
ALLOWED_ORIGINS=http://localhost:5173,tauri://localhost
```

### 3. Téléchargement des dépendances et Lancement
```bash
cd backend
go mod download
go run cmd/server/main.go
```

Le serveur démarre sur `http://localhost:8080`.

---

## 📡 Endpoints de l'API Authentification

### Routes Publiques

| Méthode | Route | Description |
| :--- | :--- | :--- |
| `GET` | `/health` | Vérification de l'état de santé du serveur |
| `POST` | `/api/auth/register-school` | Onboarding : Création de l'auto-école, de l'agence siège et du compte Directeur |
| `POST` | `/api/auth/login` | Connexion (téléphone/email + mot de passe, slug optionnel) |
| `POST` | `/api/auth/refresh` | Renouvellement de l'Access Token avec rotation du Refresh Token |
| `POST` | `/api/auth/logout` | Révocation de la session en base |
| `POST` | `/api/auth/verify-pin` | Vérification du code PIN manager pour autoriser un encaissement ou dérogation |
| `GET` | `/api/auth/school?slug=...` | Métadonnées publiques de l'école pour portail candidat / moniteur |

### Routes Protégées (Bearer JWT)

| Méthode | Route | Rôle Requis | Description |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/auth/me` | Tous rôles actifs | Profil de l'utilisateur connecté avec détails école/agence |
| `POST` | `/api/auth/register-user` | `DIRECTEUR`, `SECRETAIRE` | Création d'un collaborateur (Secrétaire, Moniteur, Élève) |

---

## 📝 Exemples d'Appels JSON

### 1. Créer une nouvelle auto-école (`POST /api/auth/register-school`)
```json
{
  "nom_commercial": "Auto-École Ivoire Pro",
  "numero_agrement": "AGR-2026/045",
  "telephone_ecole": "+225 27 22 11 22 33",
  "email_ecole": "contact@ivoirepro.ci",
  "pays": "Côte d'Ivoire",
  "devise": "XOF",
  "nom_agence": "Agence Plateau",
  "adresse_agence": "Rue du Commerce, Immeuble Horizon",
  "telephone_agence": "+225 07 01 02 03 04",
  "ville_agence": "Abidjan",
  "commune_agence": "Plateau",
  "nom_directeur": "KOUASSI",
  "prenoms_directeur": "Jean-Marc",
  "telephone_directeur": "+225 07 10 20 30 40",
  "email_directeur": "jm.kouassi@ivoirepro.ci",
  "mot_de_passe": "MonMotDePasseTresSecurise2026!",
  "code_pin_manager": "7890"
}
```

### 2. Se connecter (`POST /api/auth/login`)
```json
{
  "identifiant": "+225 07 10 20 30 40",
  "mot_de_passe": "MonMotDePasseTresSecurise2026!"
}
```

Réponse :
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOi...",
    "refresh_token": "6d9a1f...",
    "expires_in": 86400,
    "user": {
      "id": "7b8e...",
      "auto_ecole_id": "a1b2...",
      "role": "DIRECTEUR",
      "statut": "ACTIF"
    },
    "school": { ... }
  }
}
```

# AutoSchool - Logiciel Hybride de Gestion d'Auto-École (Guichet Desktop & Cloud)

## 📌 Phase 3 Démarrée : Sprint 7 (Examens d'État, Bordereaux Ministériels A4 & PV Clavier)

Ce projet implémente l'architecture complète définie dans les dossiers de cadrage et de conception technique :
- **Architecture Desktop Natif** : Tauri v2 (Rust) + Vue 3 + TypeScript + TailwindCSS + Pinia.
- **Persistance Locale Offline-First** : SQLite en mode WAL (Write-Ahead Logging), ultra-résilient aux coupures brutales d'électricité et tolérant aux pannes de connexion Internet.
- **Base Cloud PostgreSQL Serverless (Neontech)** : Multi-tenant strict (`auto_ecole_id`, `agence_id`), isolation Row-Level Security (RLS), et ingestion par lots avec `process_sync_batch()`.
- **Ergonomie 100% Clavier** : Raccourcis `F1` à `F9` permettant un traitement des opérations de guichet, d'examen et de comptabilité en quelques secondes.

---

### ✨ Synthèse des Réalisations par Sprint :

#### 🔹 Sprint 0 : Socle & Design System Guichet
- Initialisation du projet hybride Desktop Tauri v2 + Vue 3.
- Base de données SQLite avec migrations complètes et tables transactionnelles.
- Layout scindé 50/50 optimisé pour les écrans d'agence avec indicateur de connectivité locale.

#### 🔹 Sprint 1 : Enrôlement Candidats & Webcam 4x4
- Module photo webcam avec mire de cadrage 4x4 officielle et déclenchement par touche `Espace`.
- Formulaire d'enrôlement rapide (`F1`) avec contrôle d'âge légal (>= 18 ans) et aptitude médicale.
- Génération et impression du Contrat d'Enseignement officiel A4 avec intégration de la photo du candidat.

#### 🔹 Sprint 2 : Caisse, Tranches & Reçus Thermiques ESC/POS
- Encaissement rapide (`F2`) supportant Espèces (Cash) et Mobile Money (Wave, Orange, MTN, Moov).
- Calcul automatique du rendu de monnaie et de la progression de règlement.
- Reçus thermiques 80 mm ESC/POS avec code de sécurité SHA-256 et QR code d'authenticité.
- Montant converti en toutes lettres selon la norme comptable OHADA.
- Détection des fraudes avec marquage officiel `*** DUPLICATA N° X ***`.
- Reçu d'attestation administrative formelle au format A5.
- Sessions de caisse avec Point de contrôle X, Clôture Z et partage du bilan sur WhatsApp.
- Verrouillage pédagogique de la conduite (<60% payé) avec dérogation par code PIN Manager (Code démo : `1234`).

#### 🔹 Sprint 3 : Code en Salle, Pointage & Tableau Gérant (Milestone 1)
- Pointage des Présences au Code (`F5`) : thèmes officiels, notation sur 40, crédit automatique de +2h.
- Feuille d'émargement officielle de code au format A4 prête pour les contrôles ministériels.
- Tableau de Bord du Gérant (`F6`) : créances, relances WhatsApp en 1 clic et liste des assidus au code.
- Matrice de Recette & Audit de conformité Milestone 1 (Critères DoD validés).

#### 🔹 Sprint 4 : Schéma PostgreSQL Cloud Neontech & Moteur de Sync
- **Schéma PostgreSQL Cloud ([`backend/database/neon_postgres_schema.sql`](backend/database/neon_postgres_schema.sql))** :
  - 14 tables relationnelles avec clés primaires UUID v4.
  - Multi-tenancy strict avec Row-Level Security (RLS) basé sur `auto_ecole_id`.
  - Triggers automatiques de recalcul de solde comptable (`fn_recalculer_solde_dossier`).
  - Procédure stockée d'ingestion par lots `process_sync_batch(p_agence_id, p_events)` avec idempotence et arbitrage Append-Only sans conflit.
- **Données Initiales ([`backend/database/neon_seed_data.sql`](backend/database/neon_seed_data.sql))** :
  - Auto-École La Référence, Agence Cocody Angré, catalogue formules A/B/C, moniteurs et véhicules.
- **Guide de Déploiement ([`backend/database/README_NEON.md`](backend/database/README_NEON.md))** :
  - Tutoriel pour déployer sur la console Neon.tech en 3 étapes.
- **Moteur Sync & Panneau Guichet ([`SyncStatusModal.vue`](src/components/guichet/SyncStatusModal.vue))** :
  - Badge interactif dans l'en-tête affichant l'état temps réel (🟢 Connecté Neon, 🟡 X en attente, 🔵 Synchro active).
  - Test de latence en direct, configuration du `DATABASE_URL` Neon, forçage manuel de la synchro et historique des lots transmis.

#### 🔹 Sprint 5 : PWA Mobile Moniteur (Planning & Livret d'Apprentissage Tactile)
- **Application Web Mobile PWA ([`MoniteurPwaView.vue`](src/components/moniteur/MoniteurPwaView.vue))** :
  - Interface Mobile-First optimisée smartphone pour utilisation en voiture (en zone blanche / sans Internet).
  - Planning quotidien du moniteur (*Kouamé Séverin*) avec créneaux de conduite (08h, 10h30, 14h), élèves assignés et véhicules affectés (*Peugeot 208 - 1245-HJ-01*).
  - Appels et messages WhatsApp directs en 1 clic vers l'élève.
  - Manifeste PWA (`/manifest.json`) permettant l'installation directe sur l'écran d'accueil du smartphone ("Ajouter à l'écran d'accueil").
- **Livret d'Apprentissage Tactile & Émargement Doigt ([`ConduiteSessionModal.vue`](src/components/moniteur/ConduiteSessionModal.vue))** :
  - Grille des compétences officielles du permis (Maniement, Circulation, Manœuvres, Autonomie) avec boutons tactiles 1 à 4 (Non acquis à Maîtrisé).
  - Relevé des compteurs kilométriques (Km Début, Km Fin, distance calculée).
  - **Pad de signature tactile HTML5 Canvas** : signature numérique tracée au doigt par l'élève et le moniteur directement sur l'écran.
  - **Décompte automatique de +2h de conduite** imputé instantanément au forfait du candidat avec génération d'un événement offline-first pour la synchronisation Cloud Neon.

#### 🔹 Sprint 6 : PWA Mobile Candidat (Espace Apprenant, Solde en Direct & Quiz Code)
- **Application Web Mobile Candidat ([`CandidatPwaView.vue`](src/components/candidat/CandidatPwaView.vue))** :
  - Espace apprenant moderne, accessible par simple lien WhatsApp (`/?mode=candidat` ou `/?dossier=REF-2026-B-0042`).
  - **Transparence Financière Totale** : Solde restant dû en temps réel (*50 000 FCFA*), jauge de règlement (67%), liens de paiement direct par Mobile Money (Wave CI, Orange Money, MTN MoMo).
  - **Historique Immuable des Reçus de Caisse** : Numéro officiel, montant, date, hash SHA-256 et statut anti-duplicata.
  - Double jauge de progression pédagogique : Code de la Route (14/20h - 70%) et Conduite Pratique (6/20h - 30%).
- **Module Séries de Quiz Code 24/7 ([`CodeQuizModal.vue`](src/components/candidat/CodeQuizModal.vue))** :
  - Séries interactives de questions officielles (Priorité à droite, sens interdit, feux jaunes, limitations en agglomération 60 km/h, distances de sécurité 2 secondes, carrefours giratoires).
  - Corrigé instantané avec justification officielle du code de la route et notation finale sur 40.
- **Module Réservation de Créneau Conduite ([`ReservationConduiteModal.vue`](src/components/candidat/ReservationConduiteModal.vue))** #### 🔹 Sprint 7 : Examens d'État, Bordereaux Ministériels A4 & PV Clavier
- **Gestion Complète des Sessions d'Examen ([`ExamensView.vue`](src/components/examens/ExamensView.vue))** :
  - Sessions d'examen théorique (Code) et pratique (Conduite) avec localisation du centre d'examen (*Centre Principal d'Abidjan-Treichville*), date et heure.
  - **Filtre d'Éligibilité Ministériel Automatique (Élimination stricte en 3 critères)** :
    1. *Solde = 0 FCFA* : Tout candidat ayant le moindre impayé est exclu de la liste de présentation.
    2. *Quotas d'heures validés* : Minimum de 20 heures de code et 20 heures de conduite certifiées.
    3. *Certificat Médical Valide* : Vérification de la conformité de l'aptitude physique.
  - **Saisie Ultra-Rapide des Résultats du PV Officiel au Clavier** :
    - Saisie continue des notes de Code (/40) et de Conduite (/20) avec bascule instantanée d'état (*Admis*, *Ajourné*, *Absent*).
    - Calcul en temps réel du taux de réussite global et par catégorie.
  - **Diffusion 1 Clic des Résultats & Convocations WhatsApp** :
    - Génération du message formaté avec date d'examen, numéro de table, lieu et rappel des pièces CNI à fournir.
    - Annonce solennelle d'admission ou message d'encouragement personnalisé pour la réinscription.
- **Bordereau Officiel Ministériel A4 Imprimable ([`BordereauMinisterielModal.vue`](src/components/examens/BordereauMinisterielModal.vue))** :
  - Document administratif A4 prêt pour inspection et dépôt auprès de la Direction des Transports Terrestres / Ministère des Transports.
  - Intégration de la photo d'identité 4x4 officielle, numéro de table, numéro de dossier, nom, date et lieu de naissance, numéro CNI et formule.
  - Cadres d'émargement, cachet officiel de l'auto-école et cadre réservé à l'Inspecteur d'État examinateur.
  - Prise en charge native de l'impression physique A4 (`@media print`) sans en-tête ni pied de page parasite.

---

## 🚀 Commandes de Démarrage

### 1. Installation des dépendances Frontend
Dans le terminal PowerShell à la racine du projet :
```bash
npm install
```

### 2. Démarrage en Mode Aperçu Web (Recommandé pour tester immédiatement)
Permet de tester instantanément toute l'ergonomie, les inscriptions, la webcam, les encaissements, le rapport Z, le code, le tableau gérant, la synchro Neon, la PWA Moniteur, l'Espace Candidat et les Examens d'État :
```bash
npm run dev
```
- Accès Guichet Desktop : `http://localhost:1420`
- Accès PWA Mobile Moniteur : `http://localhost:1420/?mode=moniteur`
- Accès PWA Mobile Candidat : `http://localhost:1420/?mode=candidat`
- Accès Module Examens d'État : `http://localhost:1420/?mode=examens`

### 3. Démarrage en Mode Desktop Natif (Tauri v2 + SQLite local)
Compile et lance l'exécutable natif Windows avec la base SQLite locale intégrée :
```bash
npm run tauri dev
```

---

## ⌨️ Guide Complet des Raccourcis Clavier du Guichet

| Touche | Action | Description |
| :---: | :--- | :--- |
| **`F1`** | **Nouvelle Inscription** | Ouvre le formulaire d'enrôlement avec capture photo 4x4 et choix de formule. |
| **`Espace`** | **Obturateur Webcam** | Capture instantanément la photo 4x4 conforme sans quitter le clavier des mains. |
| **`F2`** | **Encaissement Rapide** | Focus direct sur le champ montant avec suggestions de tranches (+25k, +50k...). |
| **`F3`** | **Recherche Dossier** | Recherche instantanée par nom, téléphone WhatsApp, numéro de dossier ou CNI. |
| **`F4`** | **Point Caisse (X & Z)** | Contrôle intermédiaire (Point X), clôture définitive (Rapport Z) et impression. |
| **`F5`** | **Présences Code** | Feuille de pointage en salle, notation sur 40 et crédit de 2 heures de cours. |
| **`F6`** | **Tableau Gérant** | Pilotage des impayés, relances WhatsApp et matrice de recette du MVP. |
| **`F7`** | **PWA Moniteur Conduite** | Planning mobile, livret d'apprentissage tactile et émargement au doigt. |
| **`F8`** | **Espace Candidat** | Espace apprenant, solde en direct, reçus et séries de quiz de code 24/7. |
| **`F9`** | **Examens d'État** | Sessions d'examen, filtrage éligibilité, bordereau A4 et PV clavier. |
| **`Badge Sync`** | **Panneau Neontech** | Test de latence Neon, saisie DATABASE_URL et forçage de synchronisation. |
| **`Entrée`** | **Valider / Imprimer** | Valide le paiement ou la saisie et ouvre le ticket thermique 80 mm. |
| **`Échap`** | **Fermer / Retour** | Ferme les modales ouvertes ou revient à l'écran scindé du guichet. |

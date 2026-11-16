# Guide de Déploiement du Schéma PostgreSQL sur Neontech (Neon.tech)

Ce guide décrit la mise en œuvre pas-à-pas du schéma de base de données Cloud **AutoSchool** sur [Neon.tech](https://neon.tech), la plateforme PostgreSQL Serverless moderne.

---

## 🚀 Étape 1 : Créer votre Projet PostgreSQL sur Neon.tech

1. Rendez-vous sur **[https://console.neon.tech](https://console.neon.tech)** et connectez-vous (via GitHub ou Google).
2. Cliquez sur **"New Project"** :
   - **Name** : `autoschool-cloud`
   - **Postgres version** : `16` ou `17` (Par défaut)
   - **Region** : Choisissez la région la plus proche de votre zone d'exploitation (ex: `eu-central-1` Francfort ou `eu-west-1` Irlande pour une latence minimale vers l'Afrique de l'Ouest).
3. Cliquez sur **"Create Project"**.
4. Neon affiche vos identifiants de connexion :
   - Copiez l'URL de connexion : `postgresql://[user]:[password]@[endpoint].neon.tech/neondb?sslmode=require`

---

## ⚡ Étape 2 : Exécuter le Schéma SQL sur Neon

Dans la console Neon.tech de votre projet :
1. Dans le menu de gauche, cliquez sur **"SQL Editor"**.
2. Ouvrez le fichier local [`backend/database/neon_postgres_schema.sql`](./neon_postgres_schema.sql), copiez son contenu intégral et collez-le dans l'éditeur SQL de Neon.
3. Cliquez sur le bouton vert **"Run"**.
   - ✅ Toutes les tables (`auto_ecoles`, `agences`, `candidats`, `dossiers_inscription`, `paiements`, `sessions_caisse`, `journal_sync_events`, etc.) sont créées.
   - ✅ Les index de recherche, les triggers de calcul de solde et la procédure d'ingestion par lots `process_sync_batch()` sont installés.
   - ✅ La sécurité au niveau des lignes (Row-Level Security - RLS) est activée.

---

## 🔑 Étape 3 : Exécuter la Migration Authentification

Dans le même **SQL Editor** de Neon.tech :
1. Ouvrez le fichier local [`backend/database/migrations/001_auth_sessions.sql`](./migrations/001_auth_sessions.sql).
2. Copiez et collez son contenu dans le SQL Editor.
3. Cliquez sur **"Run"**.
   - ✅ La table `auth_sessions`, les index de performance et les politiques de sécurité RLS sont créés.
   - ✅ Aucune fausse donnée n'est injectée : les auto-écoles et utilisateurs sont créés dynamiquement via l'API Go (`POST /api/auth/register-school`).

---

## 🔄 Étape 4 : Vérification de la Procédure de Synchronisation

Vous pouvez tester l'ingestion d'un lot d'événements hors-ligne directement dans le SQL Editor de Neon avec la commande suivante :

```sql
SELECT process_sync_batch(
    'b0000000-0000-0000-0000-000000000001'::UUID,
    '[
      {
        "event_id": "11111111-2222-3333-4444-555555555555",
        "entity_name": "paiements",
        "operation": "INSERT",
        "entity_id": "99999999-8888-7777-6666-555555555555",
        "payload": {
          "dossier_id": "e0000000-0000-0000-0000-000000000001",
          "montant": 25000.00,
          "mode_reglement": "WAVE",
          "reference_transaction": "WAVE-CI-778899",
          "reference_recu": "REC-2026-10-TEST",
          "solde_apres_paiement": 25000.00
        },
        "local_timestamp_utc": "2026-10-20T12:00:00Z"
      }
    ]'::jsonb
);
```

Résultat attendu :
```json
{
  "statut": "SUCCES",
  "evenements_traites": 1,
  "doublons_ignores": 0,
  "erreurs": 0,
  "synced_at": "2026-10-20T..."
}
```
Le paiement est immédiatement imputé au dossier sans aucun conflit !

---

## 📱 Étape 5 : Lier l'Application Desktop à Neon

1. Ouvrez l'application AutoSchool.
2. Dans l'en-tête, cliquez sur le badge de synchronisation ou utilisez le menu de configuration.
3. Collez votre chaîne de connexion Neon (`postgresql://...`).
4. Le moteur **SyncEngine** synchronisera automatiquement les transactions en tâche de fond dès que la connexion Internet est active.

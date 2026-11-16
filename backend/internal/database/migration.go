package database

import (
	"context"
	"log"
	"time"
)

// AutoMigrateSuperAdminSchema applique de façon totalement idempotente
// les colonnes et tables nécessaires pour la gestion des quotas, licences et audit.
func (db *DB) AutoMigrateSuperAdminSchema(ctx context.Context) error {
	queries := []string{
		// 1. Colonnes sur auto_ecoles
		`ALTER TABLE auto_ecoles ADD COLUMN IF NOT EXISTS max_machines INT DEFAULT 3;`,
		`ALTER TABLE auto_ecoles ADD COLUMN IF NOT EXISTS date_expiration_licence TIMESTAMPTZ DEFAULT (NOW() + INTERVAL '1 year');`,
		`ALTER TABLE auto_ecoles ADD COLUMN IF NOT EXISTS type_licence VARCHAR(50) DEFAULT 'ANNUEL';`,
		`ALTER TABLE auto_ecoles ADD COLUMN IF NOT EXISTS modules_actifs TEXT DEFAULT '["guichet","caisse","conduite","examens","flotte","archivage_cloud"]';`,
		`ALTER TABLE auto_ecoles ADD COLUMN IF NOT EXISTS ville VARCHAR(100) DEFAULT 'Douala';`,

		// 2. Table pour tracer les machines / postes enregistrés par auto-école
		`CREATE TABLE IF NOT EXISTS machines_auto_ecoles (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			auto_ecole_id UUID NOT NULL REFERENCES auto_ecoles(id) ON DELETE CASCADE,
			machine_fingerprint VARCHAR(128) NOT NULL,
			nom_machine VARCHAR(100) NOT NULL,
			os VARCHAR(50),
			derniere_activite TIMESTAMPTZ DEFAULT NOW(),
			est_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(auto_ecole_id, machine_fingerprint)
		);`,

		// 3. Table de journalisation d'audit des actions Super-Admin et Système
		`CREATE TABLE IF NOT EXISTS audit_logs_super_admin (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			action VARCHAR(100) NOT NULL,
			cible_type VARCHAR(50) NOT NULL,
			cible_id VARCHAR(100) NOT NULL,
			details TEXT,
			ip_address VARCHAR(50),
			auto_ecole_id UUID NULL,
			user_id UUID NULL,
			auteur_nom VARCHAR(100) NULL,
			role VARCHAR(50) NULL,
			created_at TIMESTAMPTZ DEFAULT NOW()
		);`,
		`ALTER TABLE audit_logs_super_admin ADD COLUMN IF NOT EXISTS auto_ecole_id UUID NULL;`,
		`ALTER TABLE audit_logs_super_admin ADD COLUMN IF NOT EXISTS user_id UUID NULL;`,
		`ALTER TABLE audit_logs_super_admin ADD COLUMN IF NOT EXISTS auteur_nom VARCHAR(100) NULL;`,
		`ALTER TABLE audit_logs_super_admin ADD COLUMN IF NOT EXISTS role VARCHAR(50) NULL;`,

		// 4. Alertes réseau diffusées par le Super-Admin
		`CREATE TABLE IF NOT EXISTS reseau_alertes (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			titre TEXT NOT NULL,
			message TEXT NOT NULL,
			severite VARCHAR(20) NOT NULL DEFAULT 'warning',
			cible_type VARCHAR(20) NOT NULL DEFAULT 'all',
			cible_auto_ecole_id UUID NULL REFERENCES auto_ecoles(id) ON DELETE CASCADE,
			bloquante BOOLEAN DEFAULT FALSE,
			est_active BOOLEAN DEFAULT TRUE,
			date_envoi TIMESTAMPTZ DEFAULT NOW(),
			created_at TIMESTAMPTZ DEFAULT NOW()
		);`,

		// 5. Annonces réseau & événements
		`CREATE TABLE IF NOT EXISTS reseau_annonces (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			titre TEXT NOT NULL,
			contenu TEXT NOT NULL,
			type_format VARCHAR(20) NOT NULL DEFAULT 'text_only',
			image_url TEXT,
			categorie VARCHAR(30) NOT NULL DEFAULT 'general',
			cible_type VARCHAR(20) NOT NULL DEFAULT 'all',
			cible_auto_ecole_id UUID NULL REFERENCES auto_ecoles(id) ON DELETE CASCADE,
			date_evenement VARCHAR(50),
			lien_action TEXT,
			date_publication TIMESTAMPTZ DEFAULT NOW(),
			created_at TIMESTAMPTZ DEFAULT NOW()
		);`,

		// 6. Avis, notes et retours des auto-écoles
		`CREATE TABLE IF NOT EXISTS reseau_avis (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			auto_ecole_id UUID NULL REFERENCES auto_ecoles(id) ON DELETE CASCADE,
			auteur TEXT NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'Directeur',
			etoiles INT NOT NULL CHECK (etoiles >= 1 AND etoiles <= 5),
			categorie VARCHAR(30) NOT NULL DEFAULT 'general',
			commentaire TEXT NOT NULL,
			statut VARCHAR(20) NOT NULL DEFAULT 'nouveau',
			reponse_admin TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW()
		);`,

		// 7. Tickets d'assistance et support technique
		`CREATE TABLE IF NOT EXISTS reseau_support_tickets (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			auto_ecole_id UUID NULL REFERENCES auto_ecoles(id) ON DELETE SET NULL,
			nom TEXT NOT NULL,
			email TEXT NOT NULL,
			sujet VARCHAR(100) NOT NULL,
			priorite VARCHAR(20) NOT NULL DEFAULT 'normale',
			message TEXT NOT NULL,
			statut VARCHAR(20) NOT NULL DEFAULT 'ouvert',
			created_at TIMESTAMPTZ DEFAULT NOW()
		);`,

		// 8. Notifications multi-tenant persistantes pour tous les intervenants
		`CREATE TABLE IF NOT EXISTS notifications_systeme (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			auto_ecole_id UUID NULL REFERENCES auto_ecoles(id) ON DELETE CASCADE,
			user_id UUID NULL REFERENCES utilisateurs(id) ON DELETE CASCADE,
			role_cible VARCHAR(50) NULL,
			titre TEXT NOT NULL,
			message TEXT NOT NULL,
			type VARCHAR(30) NOT NULL DEFAULT 'info',
			module VARCHAR(50) NOT NULL DEFAULT 'systeme',
			action_url TEXT,
			metadata JSONB,
			est_lu BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMPTZ DEFAULT NOW()
		);`,
		`CREATE INDEX IF NOT EXISTS idx_notifications_school ON notifications_systeme(auto_ecole_id, est_lu);`,
		`CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications_systeme(user_id, est_lu);`,
		`CREATE INDEX IF NOT EXISTS idx_notifications_role ON notifications_systeme(role_cible, est_lu);`,
	}

	for _, q := range queries {
		ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
		_, err := db.ExecContext(ctxTimeout, q)
		cancel()
		if err != nil {
			log.Printf("⚠️ Avertissement migration douce schema super-admin: %v (requête continuée)", err)
		}
	}

	return nil
}

use rusqlite::{Connection, Result};
use std::path::PathBuf;
use std::sync::Mutex;

pub struct Database {
    pub conn: Mutex<Connection>,
}

impl Database {
    fn resolve_or_create_machine_key(db_path: &PathBuf) -> String {
        if let Ok(env_key) = std::env::var("SQLCIPHER_KEY") {
            if !env_key.trim().is_empty() {
                return env_key;
            }
        }
        if let Some(dir) = db_path.parent() {
            let key_file = dir.join(".local_vault.key");
            if let Ok(existing) = std::fs::read_to_string(&key_file) {
                let trimmed = existing.trim().to_string();
                if !trimmed.is_empty() {
                    return trimmed;
                }
            }
            let generated = format!("{}{}", uuid::Uuid::new_v4().simple(), uuid::Uuid::new_v4().simple());
            let _ = std::fs::write(&key_file, &generated);
            return generated;
        }
        uuid::Uuid::new_v4().simple().to_string()
    }

    pub fn init(db_path: PathBuf) -> Result<Self> {
        let cipher_key = Self::resolve_or_create_machine_key(&db_path);
        let conn = Connection::open(db_path)?;

        // Clé unique par machine (active avec SQLCipher ; ignorée par SQLite standard)
        let _ = conn.execute_batch(&format!("PRAGMA key = '{}';", cipher_key));

        // Optimisations critiques pour terrain instable (WAL mode = résistance coupures de courant)
        conn.execute_batch(
            "PRAGMA journal_mode = WAL;
             PRAGMA synchronous = NORMAL;
             PRAGMA foreign_keys = ON;
             PRAGMA busy_timeout = 5000;"
        )?;

        let db = Database {
            conn: Mutex::new(conn),
        };

        // Création stricte du schéma sans aucune donnée mock ou fausse donnée
        db.create_tables()?;

        Ok(db)
    }

    fn create_tables(&self) -> Result<()> {
        let conn = self.conn.lock().unwrap();
        conn.execute_batch(
            "
            CREATE TABLE IF NOT EXISTS auto_ecoles (
                id TEXT PRIMARY KEY,
                code_unique TEXT UNIQUE NOT NULL,
                nom_commercial TEXT NOT NULL,
                numero_agrement TEXT NOT NULL,
                telephone TEXT NOT NULL,
                email TEXT,
                devise TEXT DEFAULT 'XOF',
                logo_url TEXT,
                created_at TEXT DEFAULT (datetime('now'))
            );

            CREATE TABLE IF NOT EXISTS agences (
                id TEXT PRIMARY KEY,
                auto_ecole_id TEXT NOT NULL REFERENCES auto_ecoles(id),
                nom TEXT NOT NULL,
                adresse TEXT NOT NULL,
                telephone TEXT NOT NULL,
                created_at TEXT DEFAULT (datetime('now'))
            );

            CREATE TABLE IF NOT EXISTS utilisateurs (
                id TEXT PRIMARY KEY,
                auto_ecole_id TEXT NOT NULL REFERENCES auto_ecoles(id),
                agence_id TEXT NOT NULL REFERENCES agences(id),
                nom TEXT NOT NULL,
                prenoms TEXT NOT NULL,
                telephone TEXT UNIQUE NOT NULL,
                role TEXT NOT NULL,
                code_pin_manager TEXT,
                statut TEXT DEFAULT 'ACTIF',
                created_at TEXT DEFAULT (datetime('now'))
            );

            CREATE TABLE IF NOT EXISTS formules_permis (
                id TEXT PRIMARY KEY,
                auto_ecole_id TEXT NOT NULL REFERENCES auto_ecoles(id),
                categorie_permis TEXT NOT NULL,
                intitule TEXT NOT NULL,
                description TEXT,
                heures_theoriques INTEGER DEFAULT 20,
                heures_pratiques INTEGER DEFAULT 20,
                prix_total REAL NOT NULL,
                acompte_minimum REAL NOT NULL,
                est_actif INTEGER DEFAULT 1
            );

            CREATE TABLE IF NOT EXISTS candidats (
                id TEXT PRIMARY KEY,
                auto_ecole_id TEXT NOT NULL REFERENCES auto_ecoles(id),
                agence_id TEXT NOT NULL REFERENCES agences(id),
                nom TEXT NOT NULL,
                prenoms TEXT NOT NULL,
                genre TEXT NOT NULL,
                date_naissance TEXT NOT NULL,
                lieu_naissance TEXT NOT NULL,
                nationalite TEXT DEFAULT 'Ivoirienne',
                profession TEXT,
                quartier_residence TEXT,
                telephone_whatsapp TEXT NOT NULL,
                telephone_urgence TEXT,
                groupe_sanguin TEXT,
                photo_url TEXT,
                type_piece_id TEXT NOT NULL,
                numero_piece_id TEXT NOT NULL,
                date_expiration_piece TEXT,
                visite_medicale_statut TEXT DEFAULT 'EN_ATTENTE',
                date_visite_medicale TEXT,
                created_at TEXT DEFAULT (datetime('now'))
            );

            CREATE TABLE IF NOT EXISTS dossiers_inscription (
                id TEXT PRIMARY KEY,
                auto_ecole_id TEXT NOT NULL REFERENCES auto_ecoles(id),
                agence_id TEXT NOT NULL REFERENCES agences(id),
                candidat_id TEXT NOT NULL REFERENCES candidats(id),
                formule_id TEXT NOT NULL REFERENCES formules_permis(id),
                numero_dossier TEXT UNIQUE NOT NULL,
                date_inscription TEXT DEFAULT (date('now')),
                montant_total REAL NOT NULL,
                remise_accordee REAL DEFAULT 0.0,
                solde_restant REAL NOT NULL,
                statut_dossier TEXT DEFAULT 'EN_COURS',
                heures_code_effectuees INTEGER DEFAULT 0,
                heures_conduite_effectuees INTEGER DEFAULT 0,
                derogation_conduite_accordee INTEGER DEFAULT 0,
                mode_formation_theorique TEXT DEFAULT 'PRESENTIEL'
            );

            CREATE TABLE IF NOT EXISTS sessions_caisse (
                id TEXT PRIMARY KEY,
                agence_id TEXT NOT NULL REFERENCES agences(id),
                caissier_id TEXT NOT NULL REFERENCES utilisateurs(id),
                date_ouverture TEXT NOT NULL,
                date_fermeture TEXT,
                fond_de_caisse REAL DEFAULT 0.0,
                total_especes_theorique REAL DEFAULT 0.0,
                total_especes_compte REAL,
                ecart_caisse REAL,
                statut TEXT DEFAULT 'OUVERTE',
                observations TEXT
            );

            CREATE TABLE IF NOT EXISTS paiements (
                id TEXT PRIMARY KEY,
                auto_ecole_id TEXT NOT NULL REFERENCES auto_ecoles(id),
                agence_id TEXT NOT NULL REFERENCES agences(id),
                dossier_id TEXT NOT NULL REFERENCES dossiers_inscription(id),
                session_caisse_id TEXT NOT NULL REFERENCES sessions_caisse(id),
                montant REAL NOT NULL,
                mode_reglement TEXT NOT NULL,
                reference_transaction TEXT,
                date_paiement TEXT NOT NULL,
                encaisse_par_id TEXT NOT NULL REFERENCES utilisateurs(id),
                est_annule INTEGER DEFAULT 0,
                motif_annulation TEXT
            );

            CREATE TABLE IF NOT EXISTS recus_caisse (
                id TEXT PRIMARY KEY,
                paiement_id TEXT UNIQUE NOT NULL REFERENCES paiements(id),
                numero_recu TEXT UNIQUE NOT NULL,
                code_securite_hash TEXT NOT NULL,
                nombre_impressions INTEGER DEFAULT 1
            );

            CREATE TABLE IF NOT EXISTS journal_sync_events (
                id TEXT PRIMARY KEY,
                agence_id TEXT NOT NULL,
                nom_table TEXT NOT NULL,
                type_operation TEXT NOT NULL,
                entite_id TEXT NOT NULL,
                donnees_payload TEXT NOT NULL,
                horodatage_local TEXT NOT NULL,
                statut_sync TEXT DEFAULT 'EN_ATTENTE'
            );

            CREATE TABLE IF NOT EXISTS seances_code (
                id TEXT PRIMARY KEY,
                auto_ecole_id TEXT NOT NULL,
                agence_id TEXT NOT NULL,
                date_seance TEXT NOT NULL,
                heure_debut TEXT NOT NULL,
                heure_fin TEXT NOT NULL,
                theme_traite TEXT NOT NULL,
                formateur_nom TEXT NOT NULL,
                moniteur_id TEXT,
                type_seance TEXT DEFAULT 'PRESENTIEL',
                lien_visio TEXT,
                mot_de_passe_visio TEXT,
                statut_seance TEXT DEFAULT 'PLANIFIEE',
                statut TEXT DEFAULT 'TERMINEE'
            );

            CREATE TABLE IF NOT EXISTS presences_code (
                id TEXT PRIMARY KEY,
                seance_id TEXT NOT NULL REFERENCES seances_code(id),
                candidat_id TEXT NOT NULL REFERENCES candidats(id),
                dossier_id TEXT NOT NULL REFERENCES dossiers_inscription(id),
                heure_pointage TEXT NOT NULL,
                score_test INTEGER,
                est_present INTEGER DEFAULT 1
            );

            CREATE TABLE IF NOT EXISTS sessions_examen (
                id TEXT PRIMARY KEY,
                auto_ecole_id TEXT NOT NULL,
                agence_id TEXT NOT NULL,
                code_session TEXT NOT NULL,
                date_examen TEXT NOT NULL,
                centre_examen TEXT NOT NULL,
                categorie_permis TEXT NOT NULL,
                statut_session TEXT DEFAULT 'EN_PREPARATION',
                inspecteur_nom TEXT,
                created_at TEXT DEFAULT (datetime('now'))
            );

            CREATE TABLE IF NOT EXISTS candidats_examen (
                id TEXT PRIMARY KEY,
                session_examen_id TEXT NOT NULL REFERENCES sessions_examen(id),
                dossier_id TEXT NOT NULL REFERENCES dossiers_inscription(id),
                candidat_id TEXT NOT NULL REFERENCES candidats(id),
                numero_table TEXT,
                resultat_code TEXT DEFAULT 'EN_ATTENTE',
                resultat_conduite TEXT DEFAULT 'EN_ATTENTE',
                resultat_final TEXT DEFAULT 'EN_ATTENTE',
                date_notification_envoyee TEXT
            );

            CREATE TABLE IF NOT EXISTS epreuves_pdf (
                id TEXT PRIMARY KEY,
                titre TEXT NOT NULL,
                categorie_permis TEXT NOT NULL,
                description TEXT,
                annee INTEGER NOT NULL,
                session TEXT,
                duree_minutes INTEGER DEFAULT 45,
                nombre_questions INTEGER DEFAULT 40,
                avec_corrige INTEGER DEFAULT 1,
                fichier_nom TEXT NOT NULL,
                fichier_taille TEXT,
                fichier_data TEXT NOT NULL,
                ajoute_par TEXT,
                telechargements INTEGER DEFAULT 0,
                created_at TEXT DEFAULT (datetime('now'))
            );

            CREATE TABLE IF NOT EXISTS tutoriels_video (
                id TEXT PRIMARY KEY,
                titre TEXT NOT NULL,
                categorie_permis TEXT NOT NULL DEFAULT 'TOUTES',
                description TEXT,
                url_video TEXT NOT NULL,
                plateforme TEXT DEFAULT 'YOUTUBE',
                duree_minutes INTEGER DEFAULT 5,
                miniature_url TEXT,
                vues_count INTEGER DEFAULT 0,
                ajoute_par TEXT,
                created_at TEXT DEFAULT (datetime('now'))
            );

            CREATE TABLE IF NOT EXISTS supports_cours (
                id TEXT PRIMARY KEY,
                titre TEXT NOT NULL,
                categorie_permis TEXT NOT NULL DEFAULT 'TOUTES',
                description TEXT,
                type_support TEXT DEFAULT 'PDF',
                fichier_nom TEXT,
                fichier_taille TEXT,
                fichier_data TEXT,
                url_externe TEXT,
                ajoute_par TEXT,
                telechargements INTEGER DEFAULT 0,
                created_at TEXT DEFAULT (datetime('now'))
            );

            CREATE TABLE IF NOT EXISTS vehicules (
                id TEXT PRIMARY KEY,
                agence_id TEXT NOT NULL,
                immatriculation TEXT UNIQUE NOT NULL,
                marque_modele TEXT NOT NULL,
                type_boite TEXT DEFAULT 'MANUELLE',
                kilometrage_actuel INTEGER DEFAULT 0,
                date_expiration_assurance TEXT NOT NULL,
                date_visite_technique TEXT NOT NULL,
                statut TEXT DEFAULT 'DISPONIBLE',
                created_at TEXT DEFAULT (datetime('now')),
                updated_at TEXT DEFAULT (datetime('now'))
            );

            CREATE INDEX IF NOT EXISTS idx_sqlite_epreuves_cat ON epreuves_pdf(categorie_permis);
            CREATE INDEX IF NOT EXISTS idx_sqlite_tutoriels_cat ON tutoriels_video(categorie_permis);
            CREATE INDEX IF NOT EXISTS idx_sqlite_supports_cat ON supports_cours(categorie_permis);
            "
        )?;

        // Migrations non-destructives pour bases SQLite déjà existantes
        let _ = conn.execute("ALTER TABLE utilisateurs ADD COLUMN code_pin_manager TEXT;", []);
        let _ = conn.execute("ALTER TABLE formules_permis ADD COLUMN description TEXT;", []);
        let _ = conn.execute("ALTER TABLE formules_permis ADD COLUMN tranches_paiement TEXT;", []);
        let _ = conn.execute("ALTER TABLE dossiers_inscription ADD COLUMN derogation_conduite_accordee INTEGER DEFAULT 0;", []);
        let _ = conn.execute("ALTER TABLE dossiers_inscription ADD COLUMN mode_formation_theorique TEXT DEFAULT 'PRESENTIEL';", []);
        let _ = conn.execute("ALTER TABLE sessions_caisse ADD COLUMN observations TEXT;", []);
        let _ = conn.execute("ALTER TABLE seances_code ADD COLUMN type_seance TEXT DEFAULT 'PRESENTIEL';", []);
        let _ = conn.execute("ALTER TABLE seances_code ADD COLUMN lien_visio TEXT;", []);
        let _ = conn.execute("ALTER TABLE seances_code ADD COLUMN mot_de_passe_visio TEXT;", []);
        let _ = conn.execute("ALTER TABLE seances_code ADD COLUMN statut_seance TEXT DEFAULT 'PLANIFIEE';", []);
        let _ = conn.execute("ALTER TABLE seances_code ADD COLUMN moniteur_id TEXT;", []);

        Ok(())
    }
}


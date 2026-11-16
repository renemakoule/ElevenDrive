-- ==============================================================================
-- AutoSchool - Schéma de Base de Données PostgreSQL Cloud (Neontech Serverless)
-- Version : 1.0.0 (Sprint 4 - Architecture Multi-Tenancy & Offline-First)
-- Compatible : PostgreSQL 15 / 16 / 17 (Neon Serverless Postgres)
-- ==============================================================================

-- 1. EXTENSIONS REQUISES
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "citext";

-- 2. FONCTIONS SYSTÈMES & TRIGGERS UTILITAIRES
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ==============================================================================
-- 3. MODULE 1 : MULTI-TENANCY & UTILISATEURS
-- ==============================================================================

-- Table: auto_ecoles (Entreprise titulaire de l'agrément)
CREATE TABLE IF NOT EXISTS auto_ecoles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code_unique VARCHAR(20) UNIQUE NOT NULL,
    nom_commercial VARCHAR(150) NOT NULL,
    numero_agrement VARCHAR(80) NOT NULL,
    telephone VARCHAR(30) NOT NULL,
    email citext,
    devise VARCHAR(10) DEFAULT 'XOF',
    pays VARCHAR(60) DEFAULT 'Côte d''Ivoire',
    logo_url TEXT,
    est_actif BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TRIGGER trg_auto_ecoles_updated_at
BEFORE UPDATE ON auto_ecoles
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Table: agences (Succursales d'exploitation)
CREATE TABLE IF NOT EXISTS agences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    auto_ecole_id UUID NOT NULL REFERENCES auto_ecoles(id) ON DELETE CASCADE,
    code_agence VARCHAR(20) NOT NULL,
    nom VARCHAR(100) NOT NULL,
    adresse TEXT NOT NULL,
    telephone VARCHAR(30) NOT NULL,
    ville VARCHAR(80) DEFAULT 'Abidjan',
    commune VARCHAR(80) DEFAULT 'Cocody',
    est_siege BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT uq_agence_code_per_school UNIQUE (auto_ecole_id, code_agence)
);

CREATE TRIGGER trg_agences_updated_at
BEFORE UPDATE ON agences
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Table: utilisateurs (Directeurs, Secrétaires, Moniteurs, Caissières)
CREATE TABLE IF NOT EXISTS utilisateurs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    auto_ecole_id UUID NOT NULL REFERENCES auto_ecoles(id) ON DELETE CASCADE,
    agence_id UUID REFERENCES agences(id) ON DELETE SET NULL,
    nom VARCHAR(80) NOT NULL,
    prenoms VARCHAR(120) NOT NULL,
    telephone VARCHAR(30) NOT NULL,
    email citext,
    mot_de_passe_hash TEXT,
    code_pin_manager VARCHAR(10), -- Utilisé pour les dérogations rapides au guichet
    role VARCHAR(25) NOT NULL CHECK (role IN ('DIRECTEUR', 'SECRETAIRE', 'CAISSIER', 'MONITEUR', 'ELEVE')),
    statut VARCHAR(20) DEFAULT 'ACTIF' CHECK (statut IN ('ACTIF', 'SUSPENDU', 'ARCHIVE')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT uq_user_telephone_per_school UNIQUE (auto_ecole_id, telephone)
);

CREATE TRIGGER trg_utilisateurs_updated_at
BEFORE UPDATE ON utilisateurs
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Table: auth_sessions (Tokens de session, traçabilité et révocation)
CREATE TABLE IF NOT EXISTS auth_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES utilisateurs(id) ON DELETE CASCADE,
    auto_ecole_id UUID NOT NULL REFERENCES auto_ecoles(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    ip_address VARCHAR(45),
    user_agent TEXT,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TRIGGER trg_auth_sessions_updated_at
BEFORE UPDATE ON auth_sessions
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX IF NOT EXISTS idx_auth_sessions_user_id ON auth_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_auth_sessions_token_hash ON auth_sessions(token_hash);
CREATE INDEX IF NOT EXISTS idx_auth_sessions_expires_at ON auth_sessions(expires_at);


-- ==============================================================================
-- 4. MODULE 2 : CANDIDATS, FORMULES & DOSSIERS D'INSCRIPTION
-- ==============================================================================

-- Table: formules_permis (Catalogue tarifaire)
CREATE TABLE IF NOT EXISTS formules_permis (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    auto_ecole_id UUID NOT NULL REFERENCES auto_ecoles(id) ON DELETE CASCADE,
    categorie_permis VARCHAR(5) NOT NULL CHECK (categorie_permis IN ('A', 'A1', 'B', 'C', 'D', 'E')),
    intitule VARCHAR(100) NOT NULL,
    description TEXT,
    heures_theoriques INT DEFAULT 20 CHECK (heures_theoriques >= 0),
    heures_pratiques INT DEFAULT 20 CHECK (heures_pratiques >= 0),
    prix_total NUMERIC(12,2) NOT NULL CHECK (prix_total >= 0),
    acompte_minimum NUMERIC(12,2) NOT NULL CHECK (acompte_minimum >= 0),
    est_actif BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TRIGGER trg_formules_permis_updated_at
BEFORE UPDATE ON formules_permis
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Table: candidats (Fiche d'état civil officielle de l'apprenant)
CREATE TABLE IF NOT EXISTS candidats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    auto_ecole_id UUID NOT NULL REFERENCES auto_ecoles(id) ON DELETE CASCADE,
    agence_id UUID NOT NULL REFERENCES agences(id) ON DELETE CASCADE,
    nom VARCHAR(80) NOT NULL,
    prenoms VARCHAR(120) NOT NULL,
    genre VARCHAR(1) NOT NULL CHECK (genre IN ('M', 'F')),
    date_naissance DATE NOT NULL,
    lieu_naissance VARCHAR(100) NOT NULL,
    nationalite VARCHAR(60) DEFAULT 'Ivoirienne',
    profession VARCHAR(100),
    quartier_residence VARCHAR(120),
    telephone_whatsapp VARCHAR(30) NOT NULL,
    telephone_urgence VARCHAR(30),
    groupe_sanguin VARCHAR(5),
    photo_url TEXT,
    type_piece_id VARCHAR(30) NOT NULL CHECK (type_piece_id IN ('CNI', 'PASSEPORT', 'ATTESTATION', 'CARTE_CONSULAIRE')),
    numero_piece_id VARCHAR(60) NOT NULL,
    date_expiration_piece DATE,
    visite_medicale_statut VARCHAR(30) DEFAULT 'VALIDE_SANS_RESERVE' CHECK (visite_medicale_statut IN ('EN_ATTENTE', 'VALIDE_SANS_RESERVE', 'VALIDE_AVEC_VERRES', 'INAPTE')),
    date_visite_medicale DATE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    synced_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TRIGGER trg_candidats_updated_at
BEFORE UPDATE ON candidats
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Table: dossiers_inscription (Contrat de formation liant le candidat)
CREATE TABLE IF NOT EXISTS dossiers_inscription (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    auto_ecole_id UUID NOT NULL REFERENCES auto_ecoles(id) ON DELETE CASCADE,
    agence_id UUID NOT NULL REFERENCES agences(id) ON DELETE CASCADE,
    candidat_id UUID NOT NULL REFERENCES candidats(id) ON DELETE CASCADE,
    formule_id UUID NOT NULL REFERENCES formules_permis(id),
    numero_dossier VARCHAR(40) NOT NULL,
    date_inscription DATE DEFAULT CURRENT_DATE,
    montant_total NUMERIC(12,2) NOT NULL CHECK (montant_total >= 0),
    remise_accordee NUMERIC(12,2) DEFAULT 0.00 CHECK (remise_accordee >= 0),
    solde_restant NUMERIC(12,2) NOT NULL CHECK (solde_restant >= 0),
    statut_dossier VARCHAR(25) DEFAULT 'EN_COURS' CHECK (statut_dossier IN ('EN_COURS', 'PREPARE_EXAMEN', 'ADMIS', 'AJOURNE', 'ABANDONNE')),
    heures_code_effectuees INT DEFAULT 0,
    heures_conduite_effectuees INT DEFAULT 0,
    derogation_conduite_accordee BOOLEAN DEFAULT FALSE,
    mode_formation_theorique VARCHAR(20) DEFAULT 'PRESENTIEL' CHECK (mode_formation_theorique IN ('PRESENTIEL', 'EN_LIGNE', 'HYBRIDE')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    synced_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT uq_numero_dossier_per_school UNIQUE (auto_ecole_id, numero_dossier)
);

CREATE TRIGGER trg_dossiers_inscription_updated_at
BEFORE UPDATE ON dossiers_inscription
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ==============================================================================
-- 5. MODULE 3 : CAISSE, PAIEMENTS & REÇUS ESC/POS
-- ==============================================================================

-- Table: sessions_caisse (Session quotidienne caissière)
CREATE TABLE IF NOT EXISTS sessions_caisse (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agence_id UUID NOT NULL REFERENCES agences(id) ON DELETE CASCADE,
    caissier_id UUID NOT NULL REFERENCES utilisateurs(id),
    date_ouverture TIMESTAMPTZ DEFAULT NOW(),
    date_fermeture TIMESTAMPTZ,
    fond_de_caisse NUMERIC(12,2) DEFAULT 0.00,
    total_especes_theorique NUMERIC(12,2) DEFAULT 0.00,
    total_especes_compte NUMERIC(12,2),
    ecart_caisse NUMERIC(12,2),
    statut VARCHAR(15) DEFAULT 'OUVERTE' CHECK (statut IN ('OUVERTE', 'CLOTUREE')),
    observations TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    synced_at TIMESTAMPTZ DEFAULT NOW()
);

-- Table: paiements (Grand livre comptable immuable)
CREATE TABLE IF NOT EXISTS paiements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    auto_ecole_id UUID NOT NULL REFERENCES auto_ecoles(id) ON DELETE CASCADE,
    agence_id UUID NOT NULL REFERENCES agences(id) ON DELETE CASCADE,
    dossier_id UUID NOT NULL REFERENCES dossiers_inscription(id) ON DELETE CASCADE,
    session_caisse_id UUID REFERENCES sessions_caisse(id),
    montant NUMERIC(12,2) NOT NULL CHECK (montant > 0),
    mode_reglement VARCHAR(25) NOT NULL CHECK (mode_reglement IN ('ESPECES', 'WAVE', 'ORANGE_MONEY', 'MTN_MOMO', 'MOOV', 'VIREMENT', 'CHEQUE')),
    reference_transaction VARCHAR(80),
    reference_recu VARCHAR(40) NOT NULL,
    solde_apres_paiement NUMERIC(12,2) NOT NULL CHECK (solde_apres_paiement >= 0),
    date_paiement TIMESTAMPTZ DEFAULT NOW(),
    encaisse_par_id UUID REFERENCES utilisateurs(id),
    est_annule BOOLEAN DEFAULT FALSE,
    motif_annulation TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    synced_at TIMESTAMPTZ DEFAULT NOW()
);

-- Table: recus_caisse (Spécifications d'impression & traçabilité anti-fraude)
CREATE TABLE IF NOT EXISTS recus_caisse (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    paiement_id UUID NOT NULL UNIQUE REFERENCES paiements(id) ON DELETE CASCADE,
    numero_recu VARCHAR(40) NOT NULL,
    code_securite_hash VARCHAR(64) NOT NULL,
    nombre_impressions INT DEFAULT 1,
    derniere_impression_at TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Trigger de calcul et mise à jour automatique du solde sur dossier lors d'un paiement
CREATE OR REPLACE FUNCTION fn_recalculer_solde_dossier()
RETURNS TRIGGER AS $$
DECLARE
    v_total_regle NUMERIC(12,2);
    v_montant_dossier NUMERIC(12,2);
    v_remise NUMERIC(12,2);
BEGIN
    SELECT montant_total, remise_accordee 
    INTO v_montant_dossier, v_remise
    FROM dossiers_inscription
    WHERE id = NEW.dossier_id;

    SELECT COALESCE(SUM(montant), 0)
    INTO v_total_regle
    FROM paiements
    WHERE dossier_id = NEW.dossier_id AND est_annule = FALSE;

    UPDATE dossiers_inscription
    SET solde_restant = GREATEST(0, (v_montant_dossier - v_remise) - v_total_regle),
        updated_at = NOW()
    WHERE id = NEW.dossier_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_paiement_after_insert
AFTER INSERT OR UPDATE ON paiements
FOR EACH ROW EXECUTE FUNCTION fn_recalculer_solde_dossier();

-- ==============================================================================
-- 6. MODULE 4 : PÉDAGOGIE (CODE & CONDUITE)
-- ==============================================================================

-- Table: seances_code (Sessions théoriques en salle ou en ligne)
CREATE TABLE IF NOT EXISTS seances_code (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agence_id UUID NOT NULL REFERENCES agences(id) ON DELETE CASCADE,
    theme_traite VARCHAR(150) NOT NULL,
    formateur_nom VARCHAR(100) NOT NULL,
    moniteur_id UUID REFERENCES utilisateurs(id),
    date_seance DATE DEFAULT CURRENT_DATE,
    heure_debut VARCHAR(10) NOT NULL,
    heure_fin VARCHAR(10) NOT NULL,
    type_seance VARCHAR(20) DEFAULT 'PRESENTIEL' CHECK (type_seance IN ('PRESENTIEL', 'EN_LIGNE')),
    lien_visio TEXT,
    mot_de_passe_visio VARCHAR(50),
    statut_seance VARCHAR(20) DEFAULT 'PLANIFIEE' CHECK (statut_seance IN ('PLANIFIEE', 'EN_COURS', 'TERMINEE', 'ANNULEE')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    synced_at TIMESTAMPTZ DEFAULT NOW()
);

-- Table: presences_code (Émargement individuel & notes d'examens blancs)
CREATE TABLE IF NOT EXISTS presences_code (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seance_id UUID NOT NULL REFERENCES seances_code(id) ON DELETE CASCADE,
    candidat_id UUID NOT NULL REFERENCES candidats(id) ON DELETE CASCADE,
    est_present BOOLEAN DEFAULT TRUE,
    score_test INT DEFAULT 0 CHECK (score_test >= 0 AND score_test <= 40),
    heure_pointage TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT uq_candidat_seance_code UNIQUE (seance_id, candidat_id)
);

-- Table: vehicules (Flotte automobile pour les cours pratiques)
CREATE TABLE IF NOT EXISTS vehicules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agence_id UUID NOT NULL REFERENCES agences(id) ON DELETE CASCADE,
    immatriculation VARCHAR(30) NOT NULL UNIQUE,
    marque_modele VARCHAR(80) NOT NULL,
    type_boite VARCHAR(15) DEFAULT 'MANUELLE' CHECK (type_boite IN ('MANUELLE', 'AUTOMATIQUE')),
    kilometrage_actuel INT DEFAULT 0,
    date_expiration_assurance DATE NOT NULL,
    date_visite_technique DATE NOT NULL,
    statut VARCHAR(20) DEFAULT 'DISPONIBLE' CHECK (statut IN ('DISPONIBLE', 'EN_LECON', 'EN_PANNE', 'EN_ENTRETIEN')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Table: seances_conduite (Planning leçons au volant)
CREATE TABLE IF NOT EXISTS seances_conduite (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    auto_ecole_id UUID NOT NULL REFERENCES auto_ecoles(id) ON DELETE CASCADE,
    agence_id UUID NOT NULL REFERENCES agences(id) ON DELETE CASCADE,
    dossier_id UUID NOT NULL REFERENCES dossiers_inscription(id) ON DELETE CASCADE,
    moniteur_id UUID REFERENCES utilisateurs(id),
    vehicule_id UUID REFERENCES vehicules(id),
    date_seance DATE NOT NULL,
    heure_debut VARCHAR(10) NOT NULL,
    heure_fin VARCHAR(10) NOT NULL,
    duree_minutes INT DEFAULT 60,
    statut VARCHAR(20) DEFAULT 'PROGRAMME' CHECK (statut IN ('PROGRAMME', 'EFFECTUE', 'ANNULE_ELEVE', 'REPORTE')),
    competences_travaillees JSONB DEFAULT '[]'::jsonb,
    commentaire_moniteur TEXT,
    signature_eleve_url TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    synced_at TIMESTAMPTZ DEFAULT NOW()
);

-- ==============================================================================
-- 7. MODULE 5 : EXAMENS D'ÉTAT & BORDEREAUX MINISTÉRIELS
-- ==============================================================================

-- Table: sessions_examen (Sessions officielles de passage d'examen du permis)
CREATE TABLE IF NOT EXISTS sessions_examen (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    auto_ecole_id UUID NOT NULL REFERENCES auto_ecoles(id) ON DELETE CASCADE,
    agence_id UUID NOT NULL REFERENCES agences(id) ON DELETE CASCADE,
    code_session VARCHAR(30) NOT NULL,
    date_examen DATE NOT NULL,
    centre_examen VARCHAR(120) NOT NULL,
    categorie_permis VARCHAR(5) NOT NULL,
    statut_session VARCHAR(25) DEFAULT 'EN_PREPARATION' CHECK (statut_session IN ('EN_PREPARATION', 'BORDEREAU_DEPOSE', 'RESULTATS_PUBLIES', 'CLOTUREE')),
    inspecteur_nom VARCHAR(100),
    observations TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Table: candidats_examen (Lignes du bordereau ministériel & PV de résultats)
CREATE TABLE IF NOT EXISTS candidats_examen (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_examen_id UUID NOT NULL REFERENCES sessions_examen(id) ON DELETE CASCADE,
    dossier_id UUID NOT NULL REFERENCES dossiers_inscription(id) ON DELETE CASCADE,
    candidat_id UUID NOT NULL REFERENCES candidats(id) ON DELETE CASCADE,
    numero_table VARCHAR(30),
    resultat_code VARCHAR(15) DEFAULT 'EN_ATTENTE' CHECK (resultat_code IN ('EN_ATTENTE', 'ADMIS', 'AJOURNE', 'ABSENT')),
    resultat_conduite VARCHAR(15) DEFAULT 'EN_ATTENTE' CHECK (resultat_conduite IN ('EN_ATTENTE', 'ADMIS', 'AJOURNE', 'ABSENT')),
    resultat_final VARCHAR(20) DEFAULT 'EN_ATTENTE' CHECK (resultat_final IN ('EN_ATTENTE', 'ADMIS', 'AJOURNE', 'ABSENT')),
    date_notification_envoyee TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT uq_candidat_session_examen UNIQUE (session_examen_id, candidat_id)
);

-- ==============================================================================
-- 8. MODULE 6 : JOURNAL DE SYNCHRONISATION OFFLINE-FIRST (Ingestion Cloud)
-- ==============================================================================

CREATE TABLE IF NOT EXISTS journal_sync_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL UNIQUE, -- Généré sur le client (Idempotence)
    auto_ecole_id UUID REFERENCES auto_ecoles(id) ON DELETE CASCADE,
    agence_id UUID NOT NULL REFERENCES agences(id) ON DELETE CASCADE,
    nom_table VARCHAR(50) NOT NULL,
    type_operation VARCHAR(10) NOT NULL CHECK (type_operation IN ('INSERT', 'UPDATE', 'DELETE')),
    entite_id UUID NOT NULL,
    donnees_payload JSONB NOT NULL,
    horodatage_local TIMESTAMPTZ NOT NULL,
    statut_sync VARCHAR(20) DEFAULT 'SYNCHRONISE' CHECK (statut_sync IN ('EN_ATTENTE', 'SYNCHRONISE', 'CONFLIT', 'REJETE')),
    message_erreur TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ==============================================================================
-- 8. INDEX POUR PERFORMANCES EXTRÊMES (RECHERCHES & SYNC)
-- ==============================================================================

CREATE INDEX IF NOT EXISTS idx_candidats_recherche ON candidats(nom, prenoms);
CREATE INDEX IF NOT EXISTS idx_candidats_telephone ON candidats(telephone_whatsapp);
CREATE INDEX IF NOT EXISTS idx_candidats_piece ON candidats(numero_piece_id);
CREATE INDEX IF NOT EXISTS idx_dossiers_numero ON dossiers_inscription(numero_dossier);
CREATE INDEX IF NOT EXISTS idx_dossiers_candidat ON dossiers_inscription(candidat_id);
CREATE INDEX IF NOT EXISTS idx_paiements_dossier ON paiements(dossier_id);
CREATE INDEX IF NOT EXISTS idx_paiements_date ON paiements(date_paiement);
CREATE INDEX IF NOT EXISTS idx_sync_events_agence ON journal_sync_events(agence_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_sync_events_idempotence ON journal_sync_events(event_id);

-- ==============================================================================
-- 9. PROCÉDURE STOCKÉE D'INGESTION PAR LOTS (Sync RPC Function)
-- Traite atomiquement un lot d'événements hors-ligne transmis par le Desktop
-- ==============================================================================

CREATE OR REPLACE FUNCTION process_sync_batch(
    p_agence_id UUID,
    p_events JSONB
)
RETURNS JSONB AS $$
DECLARE
    v_event JSONB;
    v_event_id UUID;
    v_table VARCHAR(50);
    v_op VARCHAR(10);
    v_entity_id UUID;
    v_payload JSONB;
    v_local_ts TIMESTAMPTZ;
    v_processed_count INT := 0;
    v_duplicate_count INT := 0;
    v_error_count INT := 0;
    v_auto_ecole_id UUID;
BEGIN
    -- Récupérer l'auto_ecole_id liée à l'agence
    SELECT auto_ecole_id INTO v_auto_ecole_id FROM agences WHERE id = p_agence_id;
    IF v_auto_ecole_id IS NULL THEN
        RAISE EXCEPTION 'Agence non trouvée ou non rattachée à une auto-école: %', p_agence_id;
    END IF;

    -- Parcourir chaque événement du lot
    FOR v_event IN SELECT * FROM jsonb_array_elements(p_events)
    LOOP
        v_event_id  := (v_event->>'event_id')::UUID;
        v_table     := v_event->>'entity_name';
        v_op        := v_event->>'operation';
        v_entity_id := (v_event->>'entity_id')::UUID;
        v_payload   := v_event->'payload';
        v_local_ts  := (v_event->>'local_timestamp_utc')::TIMESTAMPTZ;

        -- 1. Idempotence : si l'événement a déjà été traité, ignorer
        IF EXISTS (SELECT 1 FROM journal_sync_events WHERE event_id = v_event_id) THEN
            v_duplicate_count := v_duplicate_count + 1;
            CONTINUE;
        END IF;

        -- 2. Traitement selon la table cible
        BEGIN
            IF v_table = 'candidats' AND v_op = 'INSERT' THEN
                INSERT INTO candidats (
                    id, auto_ecole_id, agence_id, nom, prenoms, genre, 
                    date_naissance, lieu_naissance, telephone_whatsapp, 
                    type_piece_id, numero_piece_id, photo_url, visite_medicale_statut, created_at
                ) VALUES (
                    v_entity_id, v_auto_ecole_id, p_agence_id,
                    v_payload->>'nom', v_payload->>'prenoms', v_payload->>'genre',
                    (v_payload->>'date_naissance')::DATE, v_payload->>'lieu_naissance',
                    v_payload->>'telephone_whatsapp', v_payload->>'type_piece_id',
                    v_payload->>'numero_piece_id', v_payload->>'photo_url',
                    COALESCE(v_payload->>'visite_medicale_statut', 'VALIDE_SANS_RESERVE'),
                    v_local_ts
                ) ON CONFLICT (id) DO UPDATE SET
                    photo_url = COALESCE(EXCLUDED.photo_url, candidats.photo_url),
                    telephone_whatsapp = EXCLUDED.telephone_whatsapp,
                    updated_at = NOW();

            ELSIF v_table = 'dossiers_inscription' AND v_op = 'INSERT' THEN
                INSERT INTO dossiers_inscription (
                    id, auto_ecole_id, agence_id, candidat_id, formule_id,
                    numero_dossier, montant_total, solde_restant, date_inscription, created_at
                ) VALUES (
                    v_entity_id, v_auto_ecole_id, p_agence_id,
                    (v_payload->>'candidat_id')::UUID, (v_payload->>'formule_id')::UUID,
                    v_payload->>'numero_dossier', (v_payload->>'montant_total')::NUMERIC,
                    (v_payload->>'solde_restant')::NUMERIC, (v_payload->>'date_inscription')::DATE,
                    v_local_ts
                ) ON CONFLICT (id) DO NOTHING;

            ELSIF v_table = 'paiements' AND v_op = 'INSERT' THEN
                -- Règle Append-Only financière : jamais de conflit sur un paiement reçu
                INSERT INTO paiements (
                    id, auto_ecole_id, agence_id, dossier_id, session_caisse_id,
                    montant, mode_reglement, reference_transaction, reference_recu,
                    solde_apres_paiement, date_paiement, created_at
                ) VALUES (
                    v_entity_id, v_auto_ecole_id, p_agence_id,
                    (v_payload->>'dossier_id')::UUID, 
                    NULLIF(v_payload->>'session_caisse_id', '')::UUID,
                    (v_payload->>'montant')::NUMERIC, v_payload->>'mode_reglement',
                    v_payload->>'reference_transaction', v_payload->>'reference_recu',
                    (v_payload->>'solde_apres_paiement')::NUMERIC, v_local_ts, v_local_ts
                ) ON CONFLICT (id) DO NOTHING;

            ELSIF v_table = 'seances_code' AND v_op = 'INSERT' THEN
                INSERT INTO seances_code (
                    id, agence_id, theme_traite, formateur_nom,
                    heure_debut, heure_fin, created_at
                ) VALUES (
                    v_entity_id, p_agence_id, v_payload->>'theme_traite',
                    v_payload->>'formateur_nom', v_payload->>'heure_debut',
                    v_payload->>'heure_fin', v_local_ts
                ) ON CONFLICT (id) DO NOTHING;
            END IF;

            -- 3. Enregistrer l'événement comme synchronisé
            INSERT INTO journal_sync_events (
                event_id, auto_ecole_id, agence_id, nom_table, type_operation,
                entite_id, donnees_payload, horodatage_local, statut_sync
            ) VALUES (
                v_event_id, v_auto_ecole_id, p_agence_id, v_table, v_op,
                v_entity_id, v_payload, v_local_ts, 'SYNCHRONISE'
            );

            v_processed_count := v_processed_count + 1;

        EXCEPTION WHEN OTHERS THEN
            v_error_count := v_error_count + 1;
            INSERT INTO journal_sync_events (
                event_id, auto_ecole_id, agence_id, nom_table, type_operation,
                entite_id, donnees_payload, horodatage_local, statut_sync, message_erreur
            ) VALUES (
                v_event_id, v_auto_ecole_id, p_agence_id, v_table, v_op,
                v_entity_id, v_payload, v_local_ts, 'REJETE', SQLERRM
            );
        END;
    END LOOP;

    RETURN jsonb_build_object(
        'statut', 'SUCCES',
        'evenements_traites', v_processed_count,
        'doublons_ignores', v_duplicate_count,
        'erreurs', v_error_count,
        'synced_at', NOW()
    );
END;
$$ LANGUAGE plpgsql;

-- ==============================================================================
-- 10. ROW LEVEL SECURITY (RLS) - MULTI-TENANT ISOLATION
-- ==============================================================================

ALTER TABLE agences ENABLE ROW LEVEL SECURITY;
ALTER TABLE utilisateurs ENABLE ROW LEVEL SECURITY;
ALTER TABLE candidats ENABLE ROW LEVEL SECURITY;
ALTER TABLE formules_permis ENABLE ROW LEVEL SECURITY;
ALTER TABLE dossiers_inscription ENABLE ROW LEVEL SECURITY;
ALTER TABLE paiements ENABLE ROW LEVEL SECURITY;
ALTER TABLE vehicules ENABLE ROW LEVEL SECURITY;
ALTER TABLE seances_conduite ENABLE ROW LEVEL SECURITY;
ALTER TABLE auth_sessions ENABLE ROW LEVEL SECURITY;

-- Politique d'isolation : les requêtes ne voient que les données de leur auto_ecole_id
CREATE POLICY agences_isolation_policy ON agences
    FOR ALL USING (auto_ecole_id = NULLIF(current_setting('app.current_auto_ecole_id', true), '')::UUID OR current_setting('app.is_admin', true) = 'true');

CREATE POLICY candidats_isolation_policy ON candidats
    FOR ALL USING (auto_ecole_id = NULLIF(current_setting('app.current_auto_ecole_id', true), '')::UUID OR current_setting('app.is_admin', true) = 'true');

CREATE POLICY dossiers_isolation_policy ON dossiers_inscription
    FOR ALL USING (auto_ecole_id = NULLIF(current_setting('app.current_auto_ecole_id', true), '')::UUID OR current_setting('app.is_admin', true) = 'true');

CREATE POLICY paiements_isolation_policy ON paiements
    FOR ALL USING (auto_ecole_id = NULLIF(current_setting('app.current_auto_ecole_id', true), '')::UUID OR current_setting('app.is_admin', true) = 'true');

CREATE POLICY auth_sessions_isolation_policy ON auth_sessions
    FOR ALL USING (auto_ecole_id = NULLIF(current_setting('app.current_auto_ecole_id', true), '')::UUID OR current_setting('app.is_admin', true) = 'true');

-- ==============================================================================
-- 11. MODULE 7 : BANQUE D'ÉPREUVES & ANNALES PDF (Pédagogie & Révision)
-- ==============================================================================

CREATE TABLE IF NOT EXISTS epreuves_pdf (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    auto_ecole_id UUID NOT NULL REFERENCES auto_ecoles(id) ON DELETE CASCADE,
    titre VARCHAR(150) NOT NULL,
    categorie_permis VARCHAR(10) NOT NULL CHECK (categorie_permis IN ('A', 'A1', 'B', 'C', 'D', 'E', 'TOUTES')),
    description TEXT,
    annee INT NOT NULL DEFAULT 2024,
    session VARCHAR(100),
    duree_minutes INT DEFAULT 45 CHECK (duree_minutes > 0),
    nombre_questions INT DEFAULT 40 CHECK (nombre_questions > 0),
    avec_corrige BOOLEAN DEFAULT TRUE,
    fichier_nom VARCHAR(150) NOT NULL,
    fichier_taille VARCHAR(30),
    fichier_data TEXT NOT NULL,
    ajoute_par VARCHAR(120),
    telechargements_count INT DEFAULT 0 CHECK (telechargements_count >= 0),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

DROP TRIGGER IF EXISTS trg_epreuves_pdf_updated_at ON epreuves_pdf;
CREATE TRIGGER trg_epreuves_pdf_updated_at
BEFORE UPDATE ON epreuves_pdf
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX IF NOT EXISTS idx_epreuves_school ON epreuves_pdf(auto_ecole_id);
CREATE INDEX IF NOT EXISTS idx_epreuves_categorie ON epreuves_pdf(auto_ecole_id, categorie_permis);

ALTER TABLE epreuves_pdf ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS epreuves_pdf_isolation_policy ON epreuves_pdf;
CREATE POLICY epreuves_pdf_isolation_policy ON epreuves_pdf
    FOR ALL USING (auto_ecole_id = NULLIF(current_setting('app.current_auto_ecole_id', true), '')::UUID OR current_setting('app.is_admin', true) = 'true');

-- ==============================================================================
-- 10. MIGRATIONS NON-DESTRUCTIVES : COURS DE CODE EN LIGNE (DISTANCIEL)
-- ==============================================================================
ALTER TABLE dossiers_inscription ADD COLUMN IF NOT EXISTS mode_formation_theorique VARCHAR(20) DEFAULT 'PRESENTIEL';
ALTER TABLE seances_code ADD COLUMN IF NOT EXISTS type_seance VARCHAR(20) DEFAULT 'PRESENTIEL';
ALTER TABLE seances_code ADD COLUMN IF NOT EXISTS lien_visio TEXT;
ALTER TABLE seances_code ADD COLUMN IF NOT EXISTS mot_de_passe_visio VARCHAR(50);
ALTER TABLE seances_code ADD COLUMN IF NOT EXISTS statut_seance VARCHAR(20) DEFAULT 'PLANIFIEE';
ALTER TABLE seances_code ADD COLUMN IF NOT EXISTS moniteur_id UUID REFERENCES utilisateurs(id);

-- ==============================================================================
-- 12. MODULE 8 : TUTORIELS VIDÉO PÉDAGOGIQUES (GUIDES PRATIQUES & CODE)
-- ==============================================================================
CREATE TABLE IF NOT EXISTS tutoriels_video (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    auto_ecole_id UUID NOT NULL REFERENCES auto_ecoles(id) ON DELETE CASCADE,
    titre VARCHAR(150) NOT NULL,
    categorie_permis VARCHAR(10) NOT NULL DEFAULT 'TOUTES' CHECK (categorie_permis IN ('A', 'A1', 'B', 'C', 'D', 'E', 'TOUTES')),
    description TEXT,
    url_video TEXT NOT NULL,
    plateforme VARCHAR(30) DEFAULT 'YOUTUBE',
    duree_minutes INT DEFAULT 5,
    miniature_url TEXT,
    vues_count INT DEFAULT 0 CHECK (vues_count >= 0),
    ajoute_par VARCHAR(120),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tutoriels_school ON tutoriels_video(auto_ecole_id);
CREATE INDEX IF NOT EXISTS idx_tutoriels_categorie ON tutoriels_video(auto_ecole_id, categorie_permis);

ALTER TABLE tutoriels_video ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS tutoriels_video_isolation_policy ON tutoriels_video;
CREATE POLICY tutoriels_video_isolation_policy ON tutoriels_video
    FOR ALL USING (auto_ecole_id = NULLIF(current_setting('app.current_auto_ecole_id', true), '')::UUID OR current_setting('app.is_admin', true) = 'true');

-- ==============================================================================
-- 13. MODULE 9 : SUPPORTS DE COURS THÉORIQUES (PDF & LIENS DE COURS)
-- ==============================================================================
CREATE TABLE IF NOT EXISTS supports_cours (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    auto_ecole_id UUID NOT NULL REFERENCES auto_ecoles(id) ON DELETE CASCADE,
    titre VARCHAR(150) NOT NULL,
    categorie_permis VARCHAR(10) NOT NULL DEFAULT 'TOUTES' CHECK (categorie_permis IN ('A', 'A1', 'B', 'C', 'D', 'E', 'TOUTES')),
    description TEXT,
    type_support VARCHAR(20) NOT NULL DEFAULT 'PDF' CHECK (type_support IN ('PDF', 'LIEN_WEB')),
    fichier_nom VARCHAR(150),
    fichier_taille VARCHAR(30),
    fichier_data TEXT,
    url_externe TEXT,
    ajoute_par VARCHAR(120),
    telechargements_count INT DEFAULT 0 CHECK (telechargements_count >= 0),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_supports_school ON supports_cours(auto_ecole_id);
CREATE INDEX IF NOT EXISTS idx_supports_categorie ON supports_cours(auto_ecole_id, categorie_permis);

ALTER TABLE supports_cours ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS supports_cours_isolation_policy ON supports_cours;
CREATE POLICY supports_cours_isolation_policy ON supports_cours
    FOR ALL USING (auto_ecole_id = NULLIF(current_setting('app.current_auto_ecole_id', true), '')::UUID OR current_setting('app.is_admin', true) = 'true');


use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct AutoEcole {
    pub id: String,
    pub code_unique: String,
    pub nom_commercial: String,
    pub numero_agrement: String,
    pub telephone: String,
    pub email: Option<String>,
    pub devise: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Agence {
    pub id: String,
    pub auto_ecole_id: String,
    pub nom: String,
    pub adresse: String,
    pub telephone: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Utilisateur {
    pub id: String,
    pub nom: String,
    pub prenoms: String,
    pub telephone: String,
    pub role: String,
    pub statut: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Candidat {
    pub id: String,
    pub auto_ecole_id: String,
    pub agence_id: String,
    pub nom: String,
    pub prenoms: String,
    pub genre: String,
    pub date_naissance: String,
    pub lieu_naissance: String,
    pub nationalite: String,
    pub profession: Option<String>,
    pub quartier_residence: Option<String>,
    pub telephone_whatsapp: String,
    pub telephone_urgence: Option<String>,
    pub groupe_sanguin: Option<String>,
    pub photo_url: Option<String>,
    pub type_piece_id: String,
    pub numero_piece_id: String,
    pub date_expiration_piece: Option<String>,
    pub visite_medicale_statut: String,
    pub date_visite_medicale: Option<String>,
    pub created_at: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct FormulePermis {
    pub id: String,
    pub categorie_permis: String,
    pub intitule: String,
    pub heures_theoriques: i32,
    pub heures_pratiques: i32,
    pub prix_total: f64,
    pub acompte_minimum: f64,
    pub est_actif: bool,
    pub tranches_paiement: Option<serde_json::Value>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct DossierInscription {
    pub id: String,
    pub candidat_id: String,
    pub formule_id: String,
    pub numero_dossier: String,
    pub date_inscription: String,
    pub montant_total: f64,
    pub remise_accordee: f64,
    pub solde_restant: f64,
    pub statut_dossier: String,
    pub heures_code_effectuees: i32,
    pub heures_conduite_effectuees: i32,
    pub mode_formation_theorique: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct SeanceCode {
    pub id: String,
    pub auto_ecole_id: String,
    pub agence_id: String,
    pub date_seance: String,
    pub heure_debut: String,
    pub heure_fin: String,
    pub theme_traite: String,
    pub formateur_nom: String,
    pub moniteur_id: Option<String>,
    pub type_seance: String,
    pub lien_visio: Option<String>,
    pub mot_de_passe_visio: Option<String>,
    pub statut_seance: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct CreerSeanceCodePayload {
    pub theme_traite: String,
    pub formateur_nom: String,
    pub moniteur_id: Option<String>,
    pub date_seance: String,
    pub heure_debut: String,
    pub heure_fin: String,
    pub type_seance: String,
    pub lien_visio: Option<String>,
    pub mot_de_passe_visio: Option<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct SessionCaisse {
    pub id: String,
    pub agence_id: String,
    pub caissier_id: String,
    pub caissier_nom: String,
    pub date_ouverture: String,
    pub date_fermeture: Option<String>,
    pub fond_de_caisse: f64,
    pub total_especes_theorique: f64,
    pub total_especes_compte: Option<f64>,
    pub ecart_caisse: Option<f64>,
    pub statut: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Paiement {
    pub id: String,
    pub dossier_id: String,
    pub session_caisse_id: String,
    pub montant: f64,
    pub mode_reglement: String,
    pub reference_transaction: Option<String>,
    pub date_paiement: String,
    pub encaisse_par_id: String,
    pub reference_recu: String,
    pub solde_apres_paiement: f64,
    pub candidat_nom_complet: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct DossierDetailComplet {
    pub candidat: Candidat,
    pub dossier: DossierInscription,
    pub formule: FormulePermis,
    pub paiements: Vec<Paiement>,
    pub progression_code_pct: f64,
    pub progression_conduite_pct: f64,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct SyncEvent {
    pub id: String,
    pub nom_table: String,
    pub type_operation: String,
    pub entite_id: String,
    pub donnees_payload: String,
    pub horodatage_local: String,
    pub statut_sync: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct EpreuvePdf {
    pub id: String,
    pub titre: String,
    pub categorie_permis: String,
    pub description: Option<String>,
    pub annee: i32,
    pub session: Option<String>,
    pub duree_minutes: i32,
    pub nombre_questions: i32,
    pub avec_corrige: bool,
    pub fichier_nom: String,
    pub fichier_taille: Option<String>,
    pub fichier_data: String,
    pub ajoute_par: Option<String>,
    pub telechargements_count: i32,
    pub created_at: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct CreerEpreuvePayload {
    pub titre: String,
    pub categorie_permis: String,
    pub description: Option<String>,
    pub annee: i32,
    pub session: Option<String>,
    pub duree_minutes: i32,
    pub nombre_questions: i32,
    pub avec_corrige: bool,
    pub fichier_nom: String,
    pub fichier_taille: Option<String>,
    pub fichier_data: String,
    pub ajoute_par: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct TutorielVideo {
    pub id: String,
    pub titre: String,
    pub categorie_permis: String,
    pub description: Option<String>,
    pub url_video: String,
    pub plateforme: Option<String>,
    pub duree_minutes: Option<i32>,
    pub miniature_url: Option<String>,
    pub vues_count: Option<i32>,
    pub ajoute_par: Option<String>,
    pub created_at: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct CreerTutorielVideoPayload {
    pub titre: String,
    pub categorie_permis: String,
    pub description: Option<String>,
    pub url_video: String,
    pub plateforme: Option<String>,
    pub duree_minutes: Option<i32>,
    pub miniature_url: Option<String>,
    pub ajoute_par: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct SupportCours {
    pub id: String,
    pub titre: String,
    pub categorie_permis: String,
    pub description: Option<String>,
    pub type_support: String,
    pub fichier_nom: Option<String>,
    pub fichier_taille: Option<String>,
    pub fichier_data: Option<String>,
    pub url_externe: Option<String>,
    pub ajoute_par: Option<String>,
    pub telechargements_count: Option<i32>,
    pub created_at: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct CreerSupportCoursPayload {
    pub titre: String,
    pub categorie_permis: String,
    pub description: Option<String>,
    pub type_support: String,
    pub fichier_nom: Option<String>,
    pub fichier_taille: Option<String>,
    pub fichier_data: Option<String>,
    pub url_externe: Option<String>,
    pub ajoute_par: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct CreerFormulePayload {
    pub categorie_permis: String,
    pub intitule: String,
    pub description: Option<String>,
    pub heures_theoriques: i32,
    pub heures_pratiques: i32,
    pub prix_total: f64,
    pub acompte_minimum: f64,
    pub tranches_paiement: Option<serde_json::Value>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ModifierPrixFormulePayload {
    pub formule_id: String,
    pub nouveau_prix: f64,
    pub acompte_minimum: Option<f64>,
    pub motif: Option<String>,
    pub tranches_paiement: Option<serde_json::Value>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct CreerCandidatPayload {
    pub nom: String,
    pub prenoms: String,
    pub genre: String,
    pub date_naissance: String,
    pub lieu_naissance: String,
    pub telephone_whatsapp: String,
    pub type_piece_id: String,
    pub numero_piece_id: String,
    pub formule_id: String,
    pub acompte_initial: f64,
    pub mode_reglement: String,
    pub photo_url: Option<String>,
    pub visite_medicale_statut: Option<String>,
    pub mode_formation_theorique: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct CloturerSessionCaissePayload {
    pub total_especes_compte: f64,
    pub observations: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct AccorderDerogationPayload {
    pub dossier_id: String,
    pub code_pin: String,
    pub motif: Option<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Vehicule {
    pub id: String,
    pub agence_id: String,
    pub immatriculation: String,
    pub marque_modele: String,
    pub type_boite: String,
    pub kilometrage_actuel: i32,
    pub date_expiration_assurance: String,
    pub date_visite_technique: String,
    pub statut: String,
    pub created_at: String,
    pub updated_at: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct CreerVehiculePayload {
    pub agence_id: Option<String>,
    pub immatriculation: String,
    pub marque_modele: String,
    pub type_boite: String,
    pub kilometrage_actuel: i32,
    pub date_expiration_assurance: String,
    pub date_visite_technique: String,
    pub statut: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ForcerCandidatExamenPayload {
    pub session_examen_id: String,
    pub dossier_id: String,
    pub motif: Option<String>,
    pub code_pin: Option<String>,
}


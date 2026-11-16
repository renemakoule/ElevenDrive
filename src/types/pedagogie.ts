export type StatutVehicule =
  | "DISPONIBLE"
  | "EN_LECON"
  | "EN_PANNE"
  | "EN_ENTRETIEN"
  | "EN_MAINTENANCE";

export interface Vehicule {
  id: string;
  agence_id: string;
  immatriculation: string;
  marque_modele: string;
  type_boite: "MANUELLE" | "AUTOMATIQUE";
  kilometrage_actuel: number;
  date_expiration_assurance: string;
  date_visite_technique: string;
  statut: StatutVehicule;
  created_at?: string;
  updated_at?: string;
}

export interface CreerVehiculeInput {
  agence_id?: string;
  immatriculation: string;
  marque_modele: string;
  type_boite?: "MANUELLE" | "AUTOMATIQUE";
  kilometrage_actuel?: number;
  date_expiration_assurance: string;
  date_visite_technique: string;
  statut?: StatutVehicule;
}

export interface CompetenceItem {
  id: string;
  code: string;
  categorie: string; // Ex: "Maniement", "Circulation", "Partage & Risques", "Autonomie"
  libelle: string;
  niveau: 1 | 2 | 3 | 4; // 1: Non acquis, 2: En cours, 3: Acquis, 4: Maîtrisé
}

export interface SeanceConduite {
  id: string;
  dossier_id: string;
  candidat_id: string;
  candidat_nom: string;
  telephone_whatsapp: string;
  photo_url?: string;
  formule_nom: string;
  moniteur_id: string;
  moniteur_nom: string;
  vehicule_id: string;
  vehicule_immatriculation: string;
  vehicule_modele: string;
  date_seance: string;
  heure_debut: string;
  heure_fin: string;
  duree_minutes: number;
  statut: "PROGRAMME" | "EN_COURS" | "EFFECTUE" | "ANNULE_ELEVE" | "REPORTE";
  heures_deja_effectuees: number;
  heures_forfait_total: number;
  competences: CompetenceItem[];
  commentaire_moniteur?: string;
  km_debut?: number;
  km_fin?: number;
  signature_eleve_url?: string;
  signature_moniteur_url?: string;
}

export interface QuizOption {
  id: string;
  texte: string;
  est_correcte: boolean;
}

export interface QuizQuestion {
  id: string;
  theme: string;
  enonce: string;
  panneau_code?: string;
  situation_desc?: string;
  options: QuizOption[];
  explication: string;
}

export interface CreneauReservation {
  id: string;
  date: string;
  heure_debut: string;
  heure_fin: string;
  moniteur_nom: string;
  vehicule_nom: string;
  disponible: boolean;
}

export interface EpreuvePdf {
  id: string;
  auto_ecole_id?: string;
  titre: string;
  categorie_permis: string; // 'A' | 'A1' | 'B' | 'C' | 'D' | 'E' | 'TOUTES'
  description?: string;
  annee: number;
  session?: string;
  duree_minutes: number;
  nombre_questions: number;
  avec_corrige: boolean;
  fichier_nom: string;
  fichier_taille?: string;
  fichier_data: string; // Data URL Base64
  date_ajout: string;
  ajoute_par: string;
  telechargements_count: number;
}

export interface TutorielVideo {
  id: string;
  auto_ecole_id?: string;
  titre: string;
  categorie_permis: string; // 'A' | 'A1' | 'B' | 'C' | 'D' | 'E' | 'TOUTES'
  description?: string;
  url_video: string;
  plateforme?: "YOUTUBE" | "VIMEO" | "MP4_DIRECT" | "AUTRE";
  duree_minutes?: number;
  miniature_url?: string;
  vues_count?: number;
  ajoute_par?: string;
  created_at?: string;
}

export interface SupportCours {
  id: string;
  auto_ecole_id?: string;
  titre: string;
  categorie_permis: string; // 'A' | 'A1' | 'B' | 'C' | 'D' | 'E' | 'TOUTES'
  description?: string;
  type_support: "PDF" | "LIEN_WEB";
  fichier_nom?: string;
  fichier_taille?: string;
  fichier_data?: string; // Data URL Base64 si PDF
  url_externe?: string; // URL si lien externe
  ajoute_par?: string;
  telechargements_count?: number;
  created_at?: string;
}

export type QuizQuestionType =
  | "CHOIX_UNIQUE"
  | "CHOIX_MULTIPLE"
  | "VRAI_FAUX"
  | "REPONSE_COURTE";

export interface QuizOptionItem {
  id: string;
  texte: string;
  est_correcte: boolean;
}

export interface QuizQuestionItem {
  id: string;
  intitule: string;
  type_question: QuizQuestionType;
  points: number;
  obligatoire: boolean;
  image_url?: string;
  explication?: string;
  reponse_attendue?: string;
  options: QuizOptionItem[];
}

export interface QuizFormulaire {
  id: string;
  auto_ecole_id?: string;
  moniteur_id?: string;
  moniteur_nom: string;
  titre: string;
  description?: string;
  categorie_permis: string;
  duree_minutes: number;
  seuil_reussite_pct: number;
  est_publie: boolean;
  questions: QuizQuestionItem[];
  reponses_count?: number;
  created_at: string;
}

export interface QuizReponseDetailItem {
  question_id: string;
  intitule_question: string;
  type_question: QuizQuestionType;
  choix_ids: string[];
  choix_labels: string[];
  texte_reponse: string;
  est_correcte: boolean;
  points_obtenus: number;
  points_max: number;
  bonne_reponse_label: string;
  explication?: string;
}

export interface QuizSoumissionEleve {
  id: string;
  quiz_id: string;
  quiz_titre: string;
  auto_ecole_id?: string;
  candidat_id: string;
  candidat_nom: string;
  numero_dossier: string;
  categorie_permis: string;
  score_obtenu: number;
  total_points: number;
  pourcentage: number;
  est_reussi: boolean;
  duree_secondes: number;
  details: QuizReponseDetailItem[];
  submitted_at: string;
}

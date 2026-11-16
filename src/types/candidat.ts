export interface Candidat {
  id: string;
  auto_ecole_id: string;
  agence_id: string;
  nom: string;
  prenoms: string;
  genre: "M" | "F";
  date_naissance: string;
  lieu_naissance: string;
  nationalite: string;
  profession?: string;
  quartier_residence?: string;
  telephone_whatsapp: string;
  telephone?: string;
  telephone_urgence?: string;
  groupe_sanguin?: string;
  photo_url?: string;
  type_piece_id: string;
  numero_piece_id: string;
  date_expiration_piece?: string;
  visite_medicale_statut:
    | "EN_ATTENTE"
    | "VALIDE_SANS_RESERVE"
    | "VALIDE_AVEC_VERRES"
    | "INAPTE";
  date_visite_medicale?: string;
  created_at: string;
  numero_dossier?: string;
  formule_nom?: string;
  categorie_permis?: string;
  solde_restant?: number;
  statut_dossier?: string;
}

export interface TranchePaiementConfig {
  numero: number;
  intitule: string;
  montant: number;
  echeance?: string;
}

export interface FormulePermis {
  id: string;
  categorie_permis: string;
  intitule: string;
  description?: string;
  heures_theoriques: number;
  heures_pratiques: number;
  prix_total: number;
  acompte_minimum: number;
  est_actif: boolean;
  tranches_paiement?: TranchePaiementConfig[] | string;
}

export interface DossierInscription {
  id: string;
  candidat_id: string;
  formule_id: string;
  numero_dossier: string;
  date_inscription: string;
  montant_total: number;
  remise_accordee: number;
  solde_restant: number;
  statut_dossier:
    | "EN_COURS"
    | "PREPARE_EXAMEN"
    | "ADMIS"
    | "AJOURNE"
    | "ABANDONNE"
    | "ARCHIVE";
  heures_code_effectuees: number;
  heures_conduite_effectuees: number;
  derogation_conduite_accordee?: boolean;
  mode_formation_theorique?: "PRESENTIEL" | "EN_LIGNE" | "HYBRIDE";
}

export type ModeFormationTheorique = "PRESENTIEL" | "EN_LIGNE" | "HYBRIDE";
export type TypeSeanceCode = "PRESENTIEL" | "EN_LIGNE";
export type StatutSeanceCode =
  | "PLANIFIEE"
  | "EN_COURS"
  | "TERMINEE"
  | "ANNULEE";

export interface SeanceCode {
  id: string;
  auto_ecole_id?: string;
  agence_id?: string;
  theme_traite: string;
  formateur_nom: string;
  moniteur_id?: string;
  date_seance: string;
  heure_debut: string;
  heure_fin: string;
  type_seance: TypeSeanceCode;
  lien_visio?: string;
  mot_de_passe_visio?: string;
  statut_seance: StatutSeanceCode;
  nombre_inscrits?: number;
  nombre_presents?: number;
}

export interface Paiement {
  id: string;
  dossier_id: string;
  session_caisse_id: string;
  montant: number;
  mode_reglement:
    | "ESPECES"
    | "WAVE"
    | "ORANGE_MONEY"
    | "MTN_MOMO"
    | "MOOV"
    | "VIREMENT";
  reference_transaction?: string;
  date_paiement: string;
  encaisse_par_id: string;
  reference_recu: string;
  solde_apres_paiement: number;
  candidat_nom_complet: string;
}

export interface DossierDetailComplet {
  candidat: Candidat;
  dossier: DossierInscription;
  formule: FormulePermis;
  paiements: Paiement[];
  progression_code_pct: number;
  progression_conduite_pct: number;
  dernier_score_code?: number;
}

export interface GuichetStats {
  total_candidats: number;
  fond_de_caisse: number;
  total_especes_theorique: number;
  evenements_sync_en_attente: number;
  mode_connectivite: string;
}

export interface SessionCaisse {
  id: string;
  agence_id: string;
  caissier_id: string;
  caissier_nom: string;
  date_ouverture: string;
  date_fermeture?: string;
  fond_de_caisse: number;
  total_especes_theorique: number;
  total_especes_compte?: number;
  ecart_caisse?: number;
  statut: "OUVERTE" | "CLOTUREE";
}

import type { Agence } from "./auth";

export interface SuperAdminStats {
  total_schools: number;
  active_schools: number;
  total_users: number;
  total_candidats: number;
  total_volume_paiement: number;
  cloud_status: string;
  cluster_latency_ms: number;
}

export interface SchoolItemSummary {
  id: string;
  code_unique: string;
  nom_commercial: string;
  numero_agrement: string;
  telephone: string;
  email?: string;
  devise: string;
  pays: string;
  ville?: string;
  est_actif: boolean;
  max_machines: number;
  machines_actives: number;
  date_expiration_licence?: string;
  type_licence?: string;
  modules_actifs?: string;
  jours_restants?: number;
  statut_licence?: "ACTIVE" | "EXPIREE" | "SUSPENDUE";
  agences_count: number;
  users_count: number;
  candidats_count: number;
  total_volume: number;
  agences?: Agence[];
  created_at: string;
  updated_at: string;
}

export interface UpdateSchoolAdminPayload {
  school_id: string;
  nom_commercial: string;
  code_unique: string;
  numero_agrement: string;
  telephone: string;
  email?: string;
  ville: string;
  pays: string;
  devise: string;
  est_actif: boolean;
  max_machines: number;
  date_expiration_licence?: string;
  type_licence: string;
  modules_actifs: string;
}

export interface UpdateSchoolLicencePayload {
  school_id: string;
  max_machines: number;
  type_licence: string;
  duree_valeur: number;
  duree_unite: "JOURS" | "MOIS" | "ANNEES";
  mode?: "REPLACE" | "EXTEND" | "KEEP_DATE";
  date_expiration_licence?: string;
}

export interface AuditLogEntry {
  id: string;
  action: string;
  cible_type: string;
  cible_id: string;
  details: string;
  ip_address: string;
  created_at: string;
}

export interface SuperAdminGlobalUser {
  id: string;
  auto_ecole_id: string;
  auto_ecole_nom: string;
  nom: string;
  prenoms: string;
  telephone: string;
  email?: string;
  role: string;
  statut: string;
  created_at: string;
}

export interface SuperAdminGlobalCandidat {
  id: string;
  auto_ecole_id: string;
  auto_ecole_nom: string;
  nom: string;
  prenoms: string;
  telephone_whatsapp: string;
  genre: string;
  numero_dossier: string;
  statut_dossier: string;
  montant_total: number;
  solde_restant: number;
  created_at: string;
}

export interface VentilationCaisse {
  especes: number;
  wave: number;
  orange_money: number;
  mtn_momo: number;
  moov: number;
  virement: number;
  total_encaisse: number;
}

export interface RecetteMensuelle {
  mois: string;
  montant: number;
}

export interface FormuleStat {
  formule_id: string;
  intitule: string;
  inscrits: number;
}

export interface PaiementRecent {
  id: string;
  numero_recu: string;
  numero_dossier: string;
  candidat_nom: string;
  formule_nom: string;
  montant: number;
  mode_reglement: string;
  date_paiement: string;
}

export interface ImpayeDetail {
  dossier_id: string;
  numero_dossier: string;
  candidat_nom: string;
  telephone_whatsapp: string;
  formule: string;
  montant_total: number;
  solde_restant: number;
  deja_paye: number;
  taux_reglement_pct: number;
  heures_code: number;
  heures_code_requises: number;
  date_inscription: string;
}

export interface TableauGerantData {
  total_creances_impayes: number;
  nombre_dossiers_impayes: number;
  impayes: ImpayeDetail[];
  eleves_assidus_code: any[];
  ventilation: VentilationCaisse;
  recettes_mensuelles: RecetteMensuelle[];
  inscriptions_formules: FormuleStat[];
  derniers_paiements: PaiementRecent[];
  taux_reussite_examen: number;
  total_candidats_actifs: number;
}

export interface NetworkAlert {
  id: string;
  titre: string;
  message: string;
  severite: "info" | "warning" | "critical";
  cible_type: "all" | "single";
  cible_school_id?: string;
  cible_school_name?: string;
  bloquante: boolean;
  est_active: boolean;
  date_envoi: string;
  created_at: string;
}

export interface CreateAlertPayload {
  titre: string;
  message: string;
  severite: "info" | "warning" | "critical";
  cible_type: "all" | "single";
  cible_school_id?: string;
  bloquante: boolean;
}

export interface NetworkAnnouncement {
  id: string;
  titre: string;
  contenu: string;
  type_format: "text_only" | "text_image" | "image_only";
  image_url?: string;
  categorie: "evenement" | "mise_a_jour" | "formation" | "general";
  cible_type: "all" | "single";
  cible_school_id?: string;
  cible_school_name?: string;
  date_evenement?: string;
  lien_action?: string;
  date_publication: string;
  created_at: string;
}

export interface CreateAnnouncementPayload {
  titre: string;
  contenu: string;
  type_format: "text_only" | "text_image" | "image_only";
  image_url?: string;
  categorie: "evenement" | "mise_a_jour" | "formation" | "general";
  cible_type: "all" | "single";
  cible_school_id?: string;
  date_evenement?: string;
  lien_action?: string;
}

export interface SchoolFeedback {
  id: string;
  auto_ecole_id?: string;
  auto_ecole_nom: string;
  auteur: string;
  role: string;
  etoiles: number;
  categorie:
    | "ergonomie"
    | "guichet"
    | "pwa_mobile"
    | "examens"
    | "support"
    | "general";
  commentaire: string;
  statut: "nouveau" | "en_cours" | "resolu";
  reponse_admin?: string;
  created_at: string;
}

export interface CreateFeedbackPayload {
  auto_ecole_id?: string;
  auteur: string;
  role: string;
  etoiles: number;
  categorie:
    | "ergonomie"
    | "guichet"
    | "pwa_mobile"
    | "examens"
    | "support"
    | "general";
  commentaire: string;
}

export interface SupportTicket {
  id: string;
  auto_ecole_id?: string;
  auto_ecole_nom?: string;
  nom: string;
  email: string;
  sujet: string;
  priorite: string;
  message: string;
  statut: string;
  created_at: string;
}

export interface CreateSupportTicketPayload {
  auto_ecole_id?: string;
  nom: string;
  email: string;
  sujet: string;
  priorite: string;
  message: string;
}

export interface SystemNotificationItem {
  id: string;
  auto_ecole_id?: string;
  user_id?: string;
  role_cible?: string;
  titre: string;
  message: string;
  type: "info" | "warning" | "critical" | "success";
  module:
    | "guichet"
    | "caisse"
    | "conduite"
    | "examens"
    | "flotte"
    | "super_admin"
    | "systeme";
  action_url?: string;
  metadata?: Record<string, any>;
  est_lu: boolean;
  created_at: string;
}

export interface SystemAuditLogItem {
  id: string;
  action: string;
  cible_type: string;
  cible_id: string;
  details: string;
  ip_address: string;
  auto_ecole_id?: string;
  user_id?: string;
  auteur_nom?: string;
  role?: string;
  created_at: string;
}


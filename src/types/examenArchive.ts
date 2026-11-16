import type { Paiement } from "./candidat";

export interface SessionExamen {
  id: string;
  code_session: string;
  date_examen: string;
  centre_examen: string;
  categorie_permis: string;
  statut_session:
    | "EN_PREPARATION"
    | "BORDEREAU_DEPOSE"
    | "RESULTATS_PUBLIES"
    | "CLOTUREE";
  inspecteur_nom?: string;
  total_candidats: number;
  total_admis?: number;
  total_ajournes?: number;
}

export interface CandidatExamenItem {
  id: string;
  session_examen_id: string;
  dossier_id: string;
  numero_dossier: string;
  candidat_id: string;
  nom: string;
  prenoms: string;
  genre: "M" | "F";
  date_naissance: string;
  lieu_naissance: string;
  telephone_whatsapp: string;
  type_piece_id: string;
  numero_piece_id: string;
  photo_url?: string;
  numero_table: string;
  resultat_code: "EN_ATTENTE" | "ADMIS" | "AJOURNE" | "ABSENT";
  resultat_conduite: "EN_ATTENTE" | "ADMIS" | "AJOURNE" | "ABSENT";
  resultat_final: "EN_ATTENTE" | "ADMIS" | "AJOURNE" | "ABSENT";
  date_notification_envoyee?: string;
  solde_restant: number;
  visite_medicale_statut: string;
  heures_code_effectuees: number;
  heures_code_requises: number;
  heures_conduite_effectuees: number;
  heures_conduite_requises: number;
  est_eligible: boolean;
  est_derogation?: boolean;
  motifs_ineligibilite?: string[];
}

export interface DossierNonEligible {
  numero_dossier: string;
  nom: string;
  motifs: string[];
  dossier_id?: string;
  candidat_id?: string;
}

export interface ForcerCandidatExamenInput {
  session_examen_id: string;
  dossier_id: string;
  motif?: string;
  code_pin?: string;
}

export interface CreerOuModifierSessionExamenInput {
  id?: string;
  code_session: string;
  date_examen: string;
  centre_examen: string;
  categorie_permis: string;
  inspecteur_nom?: string;
  observations?: string;
  statut_session?:
    | "EN_PREPARATION"
    | "BORDEREAU_DEPOSE"
    | "RESULTATS_PUBLIES"
    | "CLOTUREE";
}

export interface SessionExamenActiveResponse {
  session: SessionExamen | null;
  candidats: CandidatExamenItem[];
  ineligibles: DossierNonEligible[];
}

export type MotifArchivage =
  | "PERMIS_OBTENU"
  | "DOSSIER_SOLDE"
  | "ABANDON_DESISTEMENT"
  | "PRESCRIPTION_LEGALE"
  | "TRANSFERT_AUTRE_ECOLE"
  | "AUTRE";

export type TypeDocumentImportant =
  | "PIECE_IDENTITE"
  | "VISITE_MEDICALE"
  | "PHOTO_IDENTITE"
  | "CONTRAT_FORMATION"
  | "RECU_PAIEMENT"
  | "ATTESTATION_EXAMEN"
  | "AUTRE";

export interface DocumentArchive {
  id: string;
  dossier_id: string;
  candidat_id: string;
  type_document: TypeDocumentImportant;
  titre: string;
  nom_fichier: string;
  taille_octets?: number;
  fichier_url_ou_base64: string;
  mime_type?: string;
  date_televersement: string;
  televerse_par_nom: string;
  date_archivage?: string;
  est_archive: boolean;
  notes?: string;
}

export interface DossierArchiveItem {
  id: string;
  auto_ecole_id: string;
  dossier_id: string;
  candidat_id: string;
  numero_dossier: string;
  candidat_nom: string;
  candidat_prenoms: string;
  telephone: string;
  categorie_permis: string;
  formule_intitule: string;
  montant_total: number;
  total_paye: number;
  solde_restant: number;
  heures_code_effectuees: number;
  heures_conduite_effectuees: number;
  statut_dossier_origine: string;
  date_inscription: string;
  date_archivage: string;
  archive_par_id: string;
  archive_par_nom: string;
  archive_par_role: string;
  motif_archivage: MotifArchivage;
  reference_boite_physique?: string;
  observations?: string;
  documents_count: number;
  paiements_count: number;
  documents?: DocumentArchive[];
  paiements?: Paiement[];
}

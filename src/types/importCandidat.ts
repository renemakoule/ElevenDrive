export interface CreerCandidatBatchItem {
  nom: string;
  prenoms: string;
  genre: "M" | "F";
  date_naissance: string;
  lieu_naissance: string;
  telephone_whatsapp: string;
  type_piece_id: string;
  numero_piece_id: string;
  formule_id: string;
  acompte_initial: number;
  mode_reglement: string;
  visite_medicale_statut?:
    | "VALIDE_SANS_RESERVE"
    | "VALIDE_AVEC_VERRES"
    | "EN_ATTENTE";
  mode_formation_theorique?: "PRESENTIEL" | "EN_LIGNE" | "HYBRIDE";
}

export interface ImporterCandidatsBatchInput {
  candidats: CreerCandidatBatchItem[];
  ignorer_doublons: boolean;
  impacter_caisse_jour: boolean;
}

export interface ImportErreurLigne {
  index: number;
  nom: string;
  prenoms: string;
  numero_piece_id: string;
  raison: string;
}

export interface ImportDossierCree {
  candidat_id: string;
  dossier_id: string;
  numero_dossier: string;
  nom_complet: string;
  solde_restant: number;
}

export interface ImporterCandidatsBatchResult {
  total_recus: number;
  total_crees: number;
  total_ignores: number;
  total_erreurs: number;
  dossiers_crees: ImportDossierCree[];
  erreurs: ImportErreurLigne[];
}

export interface ParsedExcelCandidateRow {
  rowNumber: number;
  rawFormuleLabel: string;
  payload: CreerCandidatBatchItem;
  formuleLabel: string;
  errors: string[];
  isDuplicateInFile: boolean;
  isValid: boolean;
}

export type SystemImportField =
  | "nom"
  | "prenoms"
  | "nom_complet"
  | "genre"
  | "date_naissance"
  | "lieu_naissance"
  | "telephone_whatsapp"
  | "type_piece_id"
  | "numero_piece_id"
  | "formule"
  | "mode_formation_theorique"
  | "visite_medicale_statut"
  | "acompte_initial"
  | "mode_reglement";

export type ExcelColumnMapping = Record<SystemImportField, string>;

export interface ExcelImportOptions {
  defaultFormuleId: string;
  useDefaultIfUnknownFormule: boolean;
  generateMissingPiece: boolean;
}

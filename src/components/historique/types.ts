export interface HistoriqueItem {
  id: string;
  timestamp: number;
  dateStr: string;
  heureStr: string;
  type:
    | "PAIEMENT"
    | "INSCRIPTION"
    | "CAISSE"
    | "CONDUITE"
    | "CODE"
    | "QUIZ"
    | "FLOTTE"
    | "DEROGATION"
    | "EXAMEN";
  titre: string;
  description: string;
  acteur?: string;
  badgeLabel: string;
  badgeClass: string;
  montant?: number;
  score?: number;
  isRealtimeLive?: boolean;
}

export const fcfa = (n: number): string => (n || 0).toLocaleString("fr-FR");

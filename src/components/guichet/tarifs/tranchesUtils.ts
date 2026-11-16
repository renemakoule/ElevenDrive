import type { TranchePaiementConfig } from "../../../types";

export type TranchesMode = "1" | "2" | "3" | "4" | "custom";

export const formatFcfa = (val: number | null | undefined): string => {
  return (val || 0).toLocaleString("fr-FR") + " FCFA";
};

export const parseTranches = (tranches: any): TranchePaiementConfig[] => {
  if (!tranches) return [];
  if (Array.isArray(tranches)) return tranches;
  if (typeof tranches === "string") {
    try {
      const parsed = JSON.parse(tranches);
      return Array.isArray(parsed) ? parsed : [];
    } catch {
      return [];
    }
  }
  return [];
};

export const genererTranchesPredefinies = (
  total: number,
  nb: number,
): TranchePaiementConfig[] => {
  if (!total || total <= 0) return [];

  if (nb === 1) {
    return [
      {
        numero: 1,
        intitule: "Paiement Comptant (100%)",
        montant: total,
        echeance: "À l'inscription",
      },
    ];
  }

  if (nb === 2) {
    const t1 = Math.round((total * 0.5) / 1000) * 1000;
    const t2 = total - t1;
    return [
      {
        numero: 1,
        intitule: "1ère Tranche (Acompte 50%)",
        montant: t1,
        echeance: "À l'inscription",
      },
      {
        numero: 2,
        intitule: "2ème Tranche (Solde 50%)",
        montant: t2,
        echeance: "Avant le début des leçons de conduite",
      },
    ];
  }

  if (nb === 3) {
    const t1 = Math.round((total * 0.4) / 1000) * 1000;
    const t2 = Math.round((total * 0.3) / 1000) * 1000;
    const t3 = total - t1 - t2;
    return [
      {
        numero: 1,
        intitule: "1ère Tranche (Acompte 40%)",
        montant: t1,
        echeance: "À l'inscription",
      },
      {
        numero: 2,
        intitule: "2ème Tranche (Conduite 30%)",
        montant: t2,
        echeance: "Au 1er cours de conduite pratique",
      },
      {
        numero: 3,
        intitule: "3ème Tranche (Solde Examen 30%)",
        montant: t3,
        echeance: "Avant dépôt du dossier d'examen",
      },
    ];
  }

  if (nb === 4) {
    const part = Math.round(total / 4 / 1000) * 1000;
    const t4 = total - part * 3;
    return [
      {
        numero: 1,
        intitule: "1ère Tranche (Acompte 25%)",
        montant: part,
        echeance: "À l'inscription",
      },
      {
        numero: 2,
        intitule: "2ème Tranche (Code validé 25%)",
        montant: part,
        echeance: "Après 10h de théorie en salle",
      },
      {
        numero: 3,
        intitule: "3ème Tranche (Mi-parcours conduite 25%)",
        montant: part,
        echeance: "Après 10h de pratique",
      },
      {
        numero: 4,
        intitule: "4ème Tranche (Solde Examen 25%)",
        montant: t4,
        echeance: "Avant bordereau officiel d'examen",
      },
    ];
  }

  return [];
};


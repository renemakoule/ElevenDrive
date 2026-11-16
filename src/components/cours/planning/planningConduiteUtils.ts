export const CATEGORIES_PERMIS = [
  { code: "ALL", label: "Toutes catégories" },
  { code: "B", label: "Permis B (Auto)" },
  { code: "A", label: "Permis A / Moto" },
  { code: "C", label: "Permis C (Poids Lourd)" },
  { code: "D", label: "Permis D (Transport)" },
  { code: "E", label: "Permis E (Remorque)" },
];

export const extraireCategorie = (nom?: string): string => {
  if (!nom) return "B";
  const upper = nom.toUpperCase();
  if (
    upper.includes("PERMIS A") ||
    upper.includes("MOTO") ||
    upper.includes("A1") ||
    upper.includes("A2")
  )
    return "A";
  if (
    upper.includes("PERMIS C") ||
    upper.includes("POIDS LOURD") ||
    upper.includes("CAMION")
  )
    return "C";
  if (
    upper.includes("PERMIS D") ||
    upper.includes("TRANSPORT") ||
    upper.includes("BUS")
  )
    return "D";
  if (upper.includes("PERMIS E") || upper.includes("REMORQUE")) return "E";
  return "B";
};

export const getCreneauIcon = (debut?: string): string => {
  if (!debut) return "🚗";
  const h = parseInt(debut.slice(0, 2), 10);
  if (h < 10) return "🌅";
  if (h < 13) return "☀️";
  if (h < 16) return "🚗";
  return "🌆";
};

export const getDureeCreneauLabel = (debut: string, fin: string): string => {
  try {
    const [hD, mD] = debut.split(":").map(Number);
    const [hF, mF] = fin.split(":").map(Number);
    const diff = hF * 60 + mF - (hD * 60 + mD);
    if (diff <= 0) return "";
    const h = Math.floor(diff / 60);
    const m = diff % 60;
    if (m === 0) return `(${h}h)`;
    return `(${h}h${m.toString().padStart(2, "0")})`;
  } catch {
    return "(2h)";
  }
};

export const getStatutBadge = (
  statut: string,
): { label: string; class: string } => {
  switch (statut) {
    case "EFFECTUE":
      return {
        label: "Validée / Effectuée",
        class: "bg-emerald-500/10 text-emerald-600 border-emerald-500/20",
      };
    case "EN_COURS":
      return {
        label: "En cours de leçon",
        class:
          "bg-amber-500/10 text-amber-600 border-amber-500/20 animate-pulse",
      };
    case "ANNULE_ELEVE":
    case "REPORTE":
      return {
        label: "Annulée / Reportée",
        class: "bg-red-500/10 text-red-600 border-red-500/20",
      };
    case "PROGRAMME":
    default:
      return {
        label: "Confirmée & Programmée",
        class: "bg-sky-500/10 text-sky-600 border-sky-500/20",
      };
  }
};

export const getProgressPercent = (
  effectuees?: number,
  total?: number,
): number => {
  return Math.min(100, Math.round(((effectuees || 0) / (total || 20)) * 100));
};


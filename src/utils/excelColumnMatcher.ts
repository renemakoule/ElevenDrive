import type {
  ExcelColumnMapping,
  FormulePermis,
  SystemImportField,
} from "../types";

export const normalizeTokenText = (raw: string): string =>
  raw
    .normalize("NFD")
    .replace(/[\u0300-\u036f]/g, "")
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, " ")
    .trim();

export const normalizeHeaderKey = (raw: string): string =>
  normalizeTokenText(raw).replace(/\s+/g, "_");

export const SYSTEM_FIELD_LABELS: {
  field: SystemImportField;
  label: string;
  required?: boolean;
  aliases: string[];
}[] = [
  {
    field: "nom",
    label: "Nom",
    required: true,
    aliases: ["nom", "nom_famille", "last_name", "nom_eleve"],
  },
  {
    field: "prenoms",
    label: "Prénoms",
    required: true,
    aliases: ["prenoms", "prenom", "first_name", "prenoms_eleve"],
  },
  {
    field: "nom_complet",
    label: "Nom complet (si 1 seule colonne)",
    aliases: [
      "nom_complet",
      "nom_et_prenoms",
      "nom_prenoms",
      "eleve",
      "candidat",
      "full_name",
    ],
  },
  { field: "genre", label: "Genre (M/F)", aliases: ["genre", "sexe", "civ"] },
  {
    field: "date_naissance",
    label: "Date de naissance",
    required: true,
    aliases: [
      "date_naissance",
      "naissance",
      "date_de_naissance",
      "dob",
      "ne_le",
    ],
  },
  {
    field: "lieu_naissance",
    label: "Lieu de naissance",
    aliases: ["lieu_naissance", "lieu", "lieu_de_naissance"],
  },
  {
    field: "telephone_whatsapp",
    label: "Téléphone / WhatsApp",
    required: true,
    aliases: [
      "telephone_whatsapp",
      "whatsapp",
      "telephone",
      "tel",
      "contact",
      "numero_tel",
      "mobile",
    ],
  },
  {
    field: "type_piece_id",
    label: "Type de pièce",
    aliases: ["type_piece_id", "type_piece", "piece_type", "nature_piece"],
  },
  {
    field: "numero_piece_id",
    label: "N° de pièce",
    required: true,
    aliases: [
      "numero_piece_id",
      "numero_piece",
      "piece",
      "cni",
      "num_piece",
      "n_piece",
    ],
  },
  {
    field: "formule",
    label: "Formule / Catégorie",
    required: true,
    aliases: [
      "formule",
      "formule_id",
      "categorie",
      "permis",
      "categorie_permis",
      "offre",
    ],
  },
  {
    field: "mode_formation_theorique",
    label: "Mode code",
    aliases: ["mode_formation_theorique", "mode_formation", "mode_code"],
  },
  {
    field: "visite_medicale_statut",
    label: "Visite médicale",
    aliases: ["visite_medicale_statut", "visite_medicale", "statut_medical"],
  },
  {
    field: "acompte_initial",
    label: "Acompte / Montant payé",
    aliases: [
      "acompte_initial",
      "acompte",
      "montant_paye",
      "versement",
      "avance",
      "deja_paye",
    ],
  },
  {
    field: "mode_reglement",
    label: "Mode de règlement",
    aliases: ["mode_reglement", "mode_paiement", "reglement", "paiement"],
  },
];

/** Détecte automatiquement la correspondance entre les entêtes Excel et les champs système. */
export function buildAutoColumnMapping(headers: string[]): ExcelColumnMapping {
  const mapping = {} as ExcelColumnMapping;
  const usedHeaders = new Set<string>();

  for (const def of SYSTEM_FIELD_LABELS) {
    let matched = "";
    for (const h of headers) {
      if (usedHeaders.has(h)) continue;
      const norm = normalizeHeaderKey(h);
      if (def.aliases.includes(norm)) {
        matched = h;
        break;
      }
    }
    if (!matched) {
      for (const h of headers) {
        if (usedHeaders.has(h)) continue;
        const norm = normalizeHeaderKey(h);
        if (def.aliases.some((a) => norm.includes(a) || a.includes(norm))) {
          matched = h;
          break;
        }
      }
    }
    if (matched) usedHeaders.add(matched);
    mapping[def.field] = matched;
  }

  return mapping;
}

/**
 * Résout intelligemment un libellé de formule issu d'un Excel externe
 * (ex: "Permis B – Code + Conduite" ou "Permis A – Moto") vers une formule du système.
 */
export function suggestFormuleForLabel(
  rawValue: string,
  formules: FormulePermis[],
): FormulePermis | undefined {
  const cleanRaw = rawValue.trim();
  if (!cleanRaw || formules.length === 0) return undefined;

  const lower = cleanRaw.toLowerCase();
  const exact =
    formules.find((f) => f.id.toLowerCase() === lower) ||
    formules.find((f) => f.intitule.trim().toLowerCase() === lower) ||
    formules.find((f) => f.categorie_permis.trim().toLowerCase() === lower);
  if (exact) return exact;

  const normInput = normalizeTokenText(cleanRaw);
  const inputTokens = new Set(normInput.split(" ").filter((t) => t.length > 0));

  // Extraire une catégorie éventuelle ("Permis A1", "Permis B ...", "Catégorie A", "B")
  const catMatch =
    normInput.match(/\bpermis\s+([a-e])(?:[12e])?\b/) ||
    normInput.match(/\bcat(?:egorie)?\s+([a-e])(?:[12e])?\b/) ||
    normInput.match(/^([a-e])(?:[12e])?\b/);
  let detectedCat = catMatch ? catMatch[1].toUpperCase() : null;

  if (!detectedCat) {
    if (
      /(recyclage|perfectionnement|performance|remise a niveau)/.test(normInput)
    )
      detectedCat = "REC";
    else if (/(moto|cyclo|scooter|2 roues|deux roues)/.test(normInput))
      detectedCat = "A";
    else if (/(poids lourd|camion)/.test(normInput)) detectedCat = "C";
    else if (/(transport en commun|autocar|bus)/.test(normInput))
      detectedCat = "D";
    else if (/(remorque|articule)/.test(normInput)) detectedCat = "E";
  }

  let bestFormule: FormulePermis | undefined;
  let bestScore = 0;

  for (const f of formules) {
    let score = 0;
    const normTitle = normalizeTokenText(f.intitule);
    const titleTokens = normTitle.split(" ").filter((t) => t.length > 0);
    const catUpper = f.categorie_permis.toUpperCase();

    if (
      detectedCat &&
      (catUpper === detectedCat ||
        (detectedCat === "REC" &&
          (catUpper === "REC" || catUpper === "RECYCLAGE")))
    ) {
      score += 10;
    }

    for (const t of titleTokens) {
      if (inputTokens.has(t) && t !== "permis") {
        score += 2;
      }
    }

    if (normTitle.includes(normInput) || normInput.includes(normTitle)) {
      score += 5;
    }

    if (score > bestScore) {
      bestScore = score;
      bestFormule = f;
    }
  }

  return bestScore > 0 ? bestFormule : undefined;
}

/** Liste les valeurs distinctes présentes dans la colonne Formule et propose un mapping initial. */
export function buildInitialFormulaValueMap(
  rawRows: Record<string, unknown>[],
  formuleColumn: string,
  formules: FormulePermis[],
): { labels: string[]; valueMap: Record<string, string> } {
  const uniqueSet = new Set<string>();
  if (formuleColumn) {
    for (const row of rawRows) {
      const val = String(row[formuleColumn] ?? "").trim();
      if (val) uniqueSet.add(val);
    }
  }
  const labels = Array.from(uniqueSet);
  const valueMap: Record<string, string> = {};
  for (const label of labels) {
    valueMap[label] = suggestFormuleForLabel(label, formules)?.id || "";
  }
  return { labels, valueMap };
}

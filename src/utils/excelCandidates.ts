import * as XLSX from "xlsx";
import type {
  CreerCandidatBatchItem,
  ExcelColumnMapping,
  ExcelImportOptions,
  FormulePermis,
  ImportErreurLigne,
  ParsedExcelCandidateRow,
} from "../types";
import { buildAutoColumnMapping } from "./excelColumnMatcher";

const pad2 = (n: number) => String(n).padStart(2, "0");

/** Normalise une date Excel (Date JS, numéro de série, JJ/MM/AAAA ou AAAA-MM-JJ) vers YYYY-MM-DD. */
export function normalizeExcelDate(val: unknown): string {
  if (val instanceof Date && !Number.isNaN(val.getTime())) {
    return `${val.getFullYear()}-${pad2(val.getMonth() + 1)}-${pad2(val.getDate())}`;
  }
  if (typeof val === "number" && val > 1000) {
    const parsed = XLSX.SSF.parse_date_code(val);
    if (parsed && parsed.y && parsed.m && parsed.d) {
      return `${parsed.y}-${pad2(parsed.m)}-${pad2(parsed.d)}`;
    }
  }
  const str = String(val ?? "").trim();
  if (!str) return "";
  if (/^\d{4}-\d{2}-\d{2}$/.test(str)) return str;
  const frMatch = str.match(/^(\d{1,2})[\/\-](\d{1,2})[\/\-](\d{4})$/);
  if (frMatch) {
    return `${frMatch[3]}-${pad2(Number(frMatch[2]))}-${pad2(Number(frMatch[1]))}`;
  }
  const parsed = new Date(str);
  if (!Number.isNaN(parsed.getTime())) {
    return `${parsed.getFullYear()}-${pad2(parsed.getMonth() + 1)}-${pad2(parsed.getDate())}`;
  }
  return "";
}

function computeAge(dateIso: string): number {
  const birth = new Date(dateIso);
  if (Number.isNaN(birth.getTime())) return -1;
  const now = new Date();
  let age = now.getFullYear() - birth.getFullYear();
  const m = now.getMonth() - birth.getMonth();
  if (m < 0 || (m === 0 && now.getDate() < birth.getDate())) age--;
  return age;
}

/** Génère et télécharge un modèle Excel prêt à être rempli par la secrétaire. */
export function downloadCandidatesTemplate(formules: FormulePermis[]): void {
  const firstFormule = formules[0]?.intitule || "Permis B Standard";
  const sampleRows = [
    {
      nom: "KOUASSI",
      prenoms: "Aya Marie",
      genre: "F",
      date_naissance: "2002-05-14",
      lieu_naissance: "Abidjan",
      telephone_whatsapp: "+225 0700000001",
      type_piece_id: "CNI",
      numero_piece_id: "CI001234567",
      formule: firstFormule,
      mode_formation_theorique: "PRESENTIEL",
      visite_medicale_statut: "VALIDE_SANS_RESERVE",
      acompte_initial: formules[0]?.acompte_minimum ?? 0,
      mode_reglement: "ESPECES",
    },
  ];

  const refRows = formules.map((f) => ({
    Intitule_Formule: f.intitule,
    Categorie: f.categorie_permis,
    Prix_Total_FCFA: f.prix_total,
    Acompte_Minimum_FCFA: f.acompte_minimum,
  }));

  const wb = XLSX.utils.book_new();
  const wsEleves = XLSX.utils.json_to_sheet(sampleRows);
  wsEleves["!cols"] = [
    { wch: 16 },
    { wch: 22 },
    { wch: 8 },
    { wch: 15 },
    { wch: 16 },
    { wch: 20 },
    { wch: 14 },
    { wch: 18 },
    { wch: 24 },
    { wch: 22 },
    { wch: 22 },
    { wch: 16 },
    { wch: 16 },
  ];
  XLSX.utils.book_append_sheet(wb, wsEleves, "Eleves");

  if (refRows.length > 0) {
    const wsRef = XLSX.utils.json_to_sheet(refRows);
    wsRef["!cols"] = [{ wch: 28 }, { wch: 12 }, { wch: 18 }, { wch: 20 }];
    XLSX.utils.book_append_sheet(wb, wsRef, "Formules_Disponibles");
  }

  XLSX.writeFile(wb, "modele-import-eleves-elevendrive.xlsx");
}

/** Lit un fichier Excel/CSV et extrait les entêtes, lignes brutes et le mapping auto-détecté. */
export async function inspectExcelFile(file: File): Promise<{
  headers: string[];
  rawRows: Record<string, unknown>[];
  autoMapping: ExcelColumnMapping;
}> {
  const buffer = await file.arrayBuffer();
  const workbook = XLSX.read(buffer, { type: "array", cellDates: true });
  const firstSheetName = workbook.SheetNames[0];
  if (!firstSheetName) {
    throw new Error("Le classeur Excel ne contient aucune feuille.");
  }

  const sheet = workbook.Sheets[firstSheetName];
  const rawRows = XLSX.utils.sheet_to_json<Record<string, unknown>>(sheet, {
    defval: "",
  });
  if (rawRows.length === 0) {
    throw new Error("La première feuille du fichier est vide.");
  }

  const headerSet = new Set<string>();
  for (const row of rawRows) {
    for (const k of Object.keys(row)) {
      if (k && !k.startsWith("__EMPTY")) headerSet.add(k);
    }
  }
  const headers = Array.from(headerSet);
  return {
    headers,
    rawRows,
    autoMapping: buildAutoColumnMapping(headers),
  };
}

function normalizeTypePiece(raw: string): string {
  const val = raw.trim().toUpperCase();
  if (!val) return "CNI";
  if (val.includes("CONSUL")) return "CONSULAIRE";
  if (val.includes("PASS")) return "PASSEPORT";
  if (val.includes("ATTEST")) return "ATTESTATION";
  if (val === "CNI" || val === "PERMIS" || val === "AUTRE") return val;
  return "AUTRE";
}

/** Transforme et valide les lignes brutes selon le mapping de colonnes et de formules choisi. */
export function parseWithMappings(
  rawRows: Record<string, unknown>[],
  mapping: ExcelColumnMapping,
  formulaValueMap: Record<string, string>,
  formules: FormulePermis[],
  options: ExcelImportOptions,
): ParsedExcelCandidateRow[] {
  const seenPieces = new Set<string>();
  const getStr = (row: Record<string, unknown>, col: string) =>
    col ? String(row[col] ?? "").trim() : "";

  return rawRows.map((raw, idx) => {
    const rowNumber = idx + 2;
    const errors: string[] = [];

    let nom = getStr(raw, mapping.nom).toUpperCase();
    let prenoms = getStr(raw, mapping.prenoms);
    if ((!nom || !prenoms) && mapping.nom_complet) {
      const full = getStr(raw, mapping.nom_complet);
      const parts = full.split(/\s+/).filter(Boolean);
      if (!nom && parts.length > 0) nom = parts[0].toUpperCase();
      if (!prenoms && parts.length > 1) prenoms = parts.slice(1).join(" ");
    }

    const genreRaw = getStr(raw, mapping.genre).toUpperCase();
    const genre: "M" | "F" = genreRaw.startsWith("F") ? "F" : "M";
    const dateRaw = mapping.date_naissance ? raw[mapping.date_naissance] : "";
    const dateNaissance = normalizeExcelDate(dateRaw);
    const lieuNaissance = getStr(raw, mapping.lieu_naissance) || "Abidjan";
    const tel = getStr(raw, mapping.telephone_whatsapp);
    const typePiece = normalizeTypePiece(getStr(raw, mapping.type_piece_id));

    let numeroPiece = getStr(raw, mapping.numero_piece_id);
    if (!numeroPiece && options.generateMissingPiece) {
      numeroPiece = `PROV-L${rowNumber}`;
    }

    const formuleText = getStr(raw, mapping.formule);
    let mappedFormuleId = formuleText ? formulaValueMap[formuleText] || "" : "";
    if (
      !mappedFormuleId &&
      (!formuleText || options.useDefaultIfUnknownFormule)
    ) {
      mappedFormuleId = options.defaultFormuleId;
    }
    const formule = formules.find((f) => f.id === mappedFormuleId);

    const modeFormationRaw = getStr(
      raw,
      mapping.mode_formation_theorique,
    ).toUpperCase();
    const modeFormation: "PRESENTIEL" | "EN_LIGNE" | "HYBRIDE" =
      modeFormationRaw === "EN_LIGNE" || modeFormationRaw === "HYBRIDE"
        ? modeFormationRaw
        : "PRESENTIEL";

    const visiteRaw = getStr(raw, mapping.visite_medicale_statut).toUpperCase();
    const visiteStatut:
      | "VALIDE_SANS_RESERVE"
      | "VALIDE_AVEC_VERRES"
      | "EN_ATTENTE" =
      visiteRaw === "VALIDE_AVEC_VERRES" || visiteRaw === "EN_ATTENTE"
        ? visiteRaw
        : "VALIDE_SANS_RESERVE";

    const acompteStr = getStr(raw, mapping.acompte_initial);
    const acompteParsed = Number(
      acompteStr.replace(/\s+/g, "").replace(",", "."),
    );
    const acompteInitial =
      Number.isFinite(acompteParsed) && acompteParsed >= 0 ? acompteParsed : 0;
    const modeReglement =
      getStr(raw, mapping.mode_reglement).toUpperCase() || "ESPECES";

    if (!nom) errors.push("Nom manquant");
    if (!prenoms) errors.push("Prénoms manquants");
    if (!tel || tel.replace(/\s+/g, "").length < 8) {
      errors.push("WhatsApp invalide (< 8 car.)");
    }
    if (!numeroPiece) errors.push("N° de pièce manquant");
    if (!dateNaissance) {
      errors.push("Date de naissance invalide");
    } else if (formule) {
      const age = computeAge(dateNaissance);
      const minAge = formule.categorie_permis === "A" ? 16 : 18;
      if (age < minAge) {
        errors.push(`Âge insuffisant (${age} ans < ${minAge} ans)`);
      }
    }
    if (!formule) {
      errors.push(`Formule non associée (${formuleText || "vide"})`);
    }

    const pieceNorm = numeroPiece.toLowerCase();
    let isDuplicateInFile = false;
    if (pieceNorm) {
      if (seenPieces.has(pieceNorm)) {
        isDuplicateInFile = true;
        errors.push("Doublon de pièce dans le fichier");
      } else {
        seenPieces.add(pieceNorm);
      }
    }

    const payload: CreerCandidatBatchItem = {
      nom,
      prenoms,
      genre,
      date_naissance: dateNaissance || "2000-01-01",
      lieu_naissance: lieuNaissance,
      telephone_whatsapp: tel,
      type_piece_id: typePiece,
      numero_piece_id: numeroPiece,
      formule_id: formule?.id || "",
      acompte_initial: acompteInitial,
      mode_reglement: modeReglement,
      visite_medicale_statut: visiteStatut,
      mode_formation_theorique: modeFormation,
    };

    return {
      rowNumber,
      rawFormuleLabel: formuleText,
      payload,
      formuleLabel: formule
        ? `${formule.intitule} (${formule.categorie_permis})`
        : formuleText || "Non définie",
      errors,
      isDuplicateInFile,
      isValid: errors.length === 0,
    };
  });
}

/** Exporte les erreurs d'import en CSV UTF-8 pour correction facile par la secrétaire. */
export function exportImportErrorsCsv(erreurs: ImportErreurLigne[]): void {
  const rows = [
    ["Ligne", "Nom", "Prénoms", "N° Pièce", "Raison"],
    ...erreurs.map((e) => [
      String(e.index),
      e.nom,
      e.prenoms,
      e.numero_piece_id,
      e.raison,
    ]),
  ];
  const content = "\uFEFF" + rows.map((r) => r.join(";")).join("\n");
  const link = document.createElement("a");
  link.href = URL.createObjectURL(
    new Blob([content], { type: "text/csv;charset=utf-8" }),
  );
  link.download = `rejets-import-eleves-${new Date().toISOString().slice(0, 10)}.csv`;
  link.click();
  URL.revokeObjectURL(link.href);
}

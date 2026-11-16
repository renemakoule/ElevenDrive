import ExcelJS from "exceljs";
import { api } from "../db/api";
import type { Candidat, DossierDetailComplet } from "../types";
import { pdfStr, saveOrDownloadBytes, wrapText } from "./pdfExport";

export interface CandidatExportRow {
  numero_dossier: string;
  nom: string;
  prenoms: string;
  genre: string;
  date_naissance: string;
  lieu_naissance: string;
  telephone_whatsapp: string;
  type_piece_id: string;
  numero_piece_id: string;
  visite_medicale_statut: string;
  formule: string;
  categorie: string;
  date_inscription: string;
  statut_dossier: string;
  montant_total: number;
  total_paye: number;
  solde_restant: number;
  statut_paiement: string;
  paiements_detail: string;
  paiements: DossierDetailComplet["paiements"];
}

async function fetchCandidatsDetails(
  candidats: Candidat[],
): Promise<CandidatExportRow[]> {
  const details: (DossierDetailComplet | null)[] = new Array(
    candidats.length,
  ).fill(null);
  const batchSize = 6;
  for (let i = 0; i < candidats.length; i += batchSize) {
    const slice = candidats.slice(i, i + batchSize);
    const res = await Promise.all(
      slice.map((c) => api.getCandidatDossier(c.id).catch(() => null)),
    );
    res.forEach((d, idx) => {
      details[i + idx] = d;
    });
  }

  return candidats.map((c, idx) => {
    const det = details[idx];
    const cand = det?.candidat || c;
    const dos = det?.dossier;
    const form = det?.formule;
    const paiements = (det?.paiements || []).filter(
      (p) => !(p as any).est_annule,
    );

    const montantTotal = dos?.montant_total ?? 0;
    const soldeRestant = dos?.solde_restant ?? c.solde_restant ?? 0;
    const totalPaye =
      paiements.length > 0
        ? paiements.reduce((s, p) => s + (p.montant || 0), 0)
        : Math.max(0, montantTotal - soldeRestant);

    let statutPaiement = "Non payé";
    if (montantTotal > 0 && soldeRestant <= 0) statutPaiement = "Soldé";
    else if (totalPaye > 0) statutPaiement = "Partiel";

    const paiementsDetail =
      paiements.length > 0
        ? paiements
            .map(
              (p) =>
                `${(p.date_paiement || "").slice(0, 10)} : ${(p.montant || 0).toLocaleString("fr-FR")} F (${p.mode_reglement}${p.reference_recu ? ` · ${p.reference_recu}` : ""})`,
            )
            .join(" | ")
        : "Aucun paiement effectué";

    return {
      numero_dossier: dos?.numero_dossier || c.numero_dossier || "",
      nom: cand.nom || "",
      prenoms: cand.prenoms || "",
      genre: cand.genre || "",
      date_naissance: cand.date_naissance || "",
      lieu_naissance: cand.lieu_naissance || "",
      telephone_whatsapp: cand.telephone_whatsapp || "",
      type_piece_id: cand.type_piece_id || "",
      numero_piece_id: cand.numero_piece_id || "",
      visite_medicale_statut: cand.visite_medicale_statut || "",
      formule: form?.intitule || c.formule_nom || "",
      categorie: form?.categorie_permis || c.categorie_permis || "",
      date_inscription: dos?.date_inscription || "",
      statut_dossier: dos?.statut_dossier || c.statut_dossier || "",
      montant_total: montantTotal,
      total_paye: totalPaye,
      solde_restant: soldeRestant,
      statut_paiement: statutPaiement,
      paiements_detail: paiementsDetail,
      paiements,
    };
  });
}

export async function exportCandidatsToExcel(
  candidats: Candidat[],
  schoolName: string,
): Promise<{ filename: string; savedPath?: string }> {
  const rows = await fetchCandidatsDetails(candidats);
  const wb = new ExcelJS.Workbook();
  wb.creator = schoolName || "AutoSchool";
  wb.created = new Date();

  const ws = wb.addWorksheet("Candidats & Paiements");
  ws.columns = [
    { header: "N° Dossier", key: "numero_dossier", width: 16 },
    { header: "Nom", key: "nom", width: 18 },
    { header: "Prénoms", key: "prenoms", width: 22 },
    { header: "Genre", key: "genre", width: 8 },
    { header: "Date Naissance", key: "date_naissance", width: 14 },
    { header: "Lieu Naissance", key: "lieu_naissance", width: 16 },
    { header: "Téléphone WhatsApp", key: "telephone_whatsapp", width: 18 },
    { header: "Type Pièce", key: "type_piece_id", width: 14 },
    { header: "N° Pièce", key: "numero_piece_id", width: 18 },
    { header: "Visite Médicale", key: "visite_medicale_statut", width: 20 },
    { header: "Formule", key: "formule", width: 26 },
    { header: "Catégorie", key: "categorie", width: 10 },
    { header: "Date Inscription", key: "date_inscription", width: 15 },
    { header: "Statut Dossier", key: "statut_dossier", width: 14 },
    { header: "Montant Forfait (FCFA)", key: "montant_total", width: 20 },
    { header: "Total Payé (FCFA)", key: "total_paye", width: 18 },
    { header: "Reste à Payer (FCFA)", key: "solde_restant", width: 20 },
    { header: "Statut Paiement", key: "statut_paiement", width: 16 },
    {
      header: "Détail Paiements Effectués",
      key: "paiements_detail",
      width: 55,
    },
  ];

  const headerRow = ws.getRow(1);
  headerRow.font = { bold: true, color: { argb: "FFFFFFFF" } };
  headerRow.fill = {
    type: "pattern",
    pattern: "solid",
    fgColor: { argb: "FF047857" },
  };

  rows.forEach((r) => ws.addRow(r));

  const sumForfait = rows.reduce((s, r) => s + r.montant_total, 0);
  const sumPaye = rows.reduce((s, r) => s + r.total_paye, 0);
  const sumReste = rows.reduce((s, r) => s + r.solde_restant, 0);
  const totalRow = ws.addRow({
    numero_dossier: "TOTAUX",
    nom: `${rows.length} candidat(s)`,
    montant_total: sumForfait,
    total_paye: sumPaye,
    solde_restant: sumReste,
  });
  totalRow.font = { bold: true };

  const wsPay = wb.addWorksheet("Détail Paiements");
  wsPay.columns = [
    { header: "N° Dossier", key: "numero_dossier", width: 16 },
    { header: "Candidat", key: "candidat", width: 28 },
    { header: "Date Paiement", key: "date", width: 16 },
    { header: "Montant (FCFA)", key: "montant", width: 16 },
    { header: "Mode Règlement", key: "mode", width: 16 },
    { header: "N° Reçu", key: "recu", width: 18 },
    { header: "Solde Après (FCFA)", key: "solde", width: 18 },
  ];
  wsPay.getRow(1).font = { bold: true, color: { argb: "FFFFFFFF" } };
  wsPay.getRow(1).fill = {
    type: "pattern",
    pattern: "solid",
    fgColor: { argb: "FF047857" },
  };
  rows.forEach((r) => {
    r.paiements.forEach((p) => {
      wsPay.addRow({
        numero_dossier: r.numero_dossier,
        candidat: `${r.nom} ${r.prenoms}`.trim(),
        date: (p.date_paiement || "").slice(0, 10),
        montant: p.montant,
        mode: p.mode_reglement,
        recu: p.reference_recu || "",
        solde: p.solde_apres_paiement,
      });
    });
  });

  const buffer = await wb.xlsx.writeBuffer();
  const filename = `candidats-paiements-${new Date().toISOString().slice(0, 10)}.xlsx`;
  return saveOrDownloadBytes(
    filename,
    new Uint8Array(buffer),
    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
  );
}

export async function exportCandidatsToPdf(
  candidats: Candidat[],
  schoolName: string,
): Promise<{ filename: string; savedPath?: string }> {
  const rows = await fetchCandidatsDetails(candidats);
  const pageWidth = 595.28;
  const pageHeight = 841.89;
  const marginX = 32;
  const pagesContent: string[] = [];
  let ops: string[] = [];
  let y = pageHeight - 40;

  const sumForfait = rows.reduce((s, r) => s + r.montant_total, 0);
  const sumPaye = rows.reduce((s, r) => s + r.total_paye, 0);
  const sumReste = rows.reduce((s, r) => s + r.solde_restant, 0);

  const drawText = (
    x: number,
    yPos: number,
    size: number,
    bold: boolean,
    text: string,
    rgb = "0.1 0.1 0.12",
  ) => {
    ops.push(
      `BT /${bold ? "F2" : "F1"} ${size} Tf ${rgb} rg ${x.toFixed(1)} ${yPos.toFixed(1)} Td ${pdfStr(text)} Tj ET`,
    );
  };
  const drawLine = (yPos: number, rgb = "0.85 0.87 0.9") => {
    ops.push(
      `${rgb} RG 0.6 w ${marginX} ${yPos.toFixed(1)} m ${(pageWidth - marginX).toFixed(1)} ${yPos.toFixed(1)} l S`,
    );
  };

  const startHeader = (pageNum: number) => {
    y = pageHeight - 40;
    drawText(
      marginX,
      y,
      13,
      true,
      `${schoolName || "AutoSchool"} - Etat des Candidats & Paiements`,
      "0.05 0.35 0.25",
    );
    y -= 15;
    drawText(
      marginX,
      y,
      8.5,
      false,
      `Exporte le ${new Date().toLocaleString("fr-FR")}  |  ${rows.length} candidat(s)  |  Page ${pageNum}`,
      "0.4 0.42 0.46",
    );
    y -= 13;
    drawText(
      marginX,
      y,
      8.5,
      true,
      `Forfaits : ${sumForfait.toLocaleString("fr-FR")} F  |  Paye : ${sumPaye.toLocaleString("fr-FR")} F  |  Reste a payer : ${sumReste.toLocaleString("fr-FR")} F`,
      "0.12 0.14 0.18",
    );
    y -= 8;
    drawLine(y, "0.05 0.45 0.3");
    y -= 14;
  };

  let pageNum = 1;
  startHeader(pageNum);

  for (const r of rows) {
    const payLines = wrapText(`Paiements : ${r.paiements_detail}`, 96);
    const needed = 42 + payLines.length * 10.5;
    if (y - needed < 42) {
      pagesContent.push(ops.join("\n"));
      ops = [];
      pageNum += 1;
      startHeader(pageNum);
    }

    drawText(
      marginX,
      y,
      9.5,
      true,
      `${r.numero_dossier || "-"}  |  ${r.nom} ${r.prenoms} (${r.genre})`,
      "0.08 0.1 0.14",
    );
    const soldeColor =
      r.solde_restant <= 0 ? "0.05 0.48 0.28" : "0.75 0.35 0.05";
    const soldeTxt =
      r.solde_restant <= 0
        ? "SOLDE (0 F)"
        : `Reste : ${r.solde_restant.toLocaleString("fr-FR")} F`;
    drawText(pageWidth - marginX - 115, y, 8.5, true, soldeTxt, soldeColor);
    y -= 11.5;

    drawText(
      marginX + 6,
      y,
      8,
      false,
      `Ne(e) : ${r.date_naissance} (${r.lieu_naissance})  |  Tel : ${r.telephone_whatsapp}  |  ${r.type_piece_id} ${r.numero_piece_id}  |  Visite : ${r.visite_medicale_statut}`,
      "0.3 0.32 0.36",
    );
    y -= 11;

    drawText(
      marginX + 6,
      y,
      8,
      false,
      `Formule : ${r.formule || "-"} (${r.categorie || "-"})  |  Forfait : ${r.montant_total.toLocaleString("fr-FR")} F  |  Paye : ${r.total_paye.toLocaleString("fr-FR")} F  |  Statut : ${r.statut_paiement}`,
      "0.22 0.24 0.28",
    );
    y -= 11;

    for (const pl of payLines) {
      drawText(marginX + 6, y, 7.8, false, pl, "0.38 0.4 0.45");
      y -= 10.5;
    }
    drawLine(y + 2);
    y -= 8;
  }

  pagesContent.push(ops.join("\n"));

  const objects: string[] = [];
  const addObj = (body: string) => {
    objects.push(body);
    return objects.length;
  };
  const f1 = addObj(
    "<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding /WinAnsiEncoding >>",
  );
  const f2 = addObj(
    "<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica-Bold /Encoding /WinAnsiEncoding >>",
  );
  const pagesIdPlaceholder = objects.length + 1 + pagesContent.length * 2;
  const pageIds: number[] = [];
  for (const stream of pagesContent) {
    const cid = addObj(
      `<< /Length ${stream.length} >>\nstream\n${stream}\nendstream`,
    );
    pageIds.push(
      addObj(
        `<< /Type /Page /Parent ${pagesIdPlaceholder} 0 R /MediaBox [0 0 ${pageWidth} ${pageHeight}] /Resources << /Font << /F1 ${f1} 0 R /F2 ${f2} 0 R >> >> /Contents ${cid} 0 R >>`,
      ),
    );
  }
  const pagesId = addObj(
    `<< /Type /Pages /Kids [${pageIds.map((id) => `${id} 0 R`).join(" ")}] /Count ${pageIds.length} >>`,
  );
  const catalogId = addObj(`<< /Type /Catalog /Pages ${pagesId} 0 R >>`);

  let pdf = "%PDF-1.4\n";
  const offsets: number[] = [0];
  for (let i = 0; i < objects.length; i++) {
    offsets.push(pdf.length);
    pdf += `${i + 1} 0 obj\n${objects[i]}\nendobj\n`;
  }
  const xrefOffset = pdf.length;
  pdf += `xref\n0 ${objects.length + 1}\n0000000000 65535 f \n`;
  for (let i = 1; i <= objects.length; i++) {
    pdf += `${String(offsets[i]).padStart(10, "0")} 00000 n \n`;
  }
  pdf += `trailer\n<< /Size ${objects.length + 1} /Root ${catalogId} 0 R >>\nstartxref\n${xrefOffset}\n%%EOF`;

  const outBytes = new Uint8Array(pdf.length);
  for (let i = 0; i < pdf.length; i++) {
    outBytes[i] = pdf.charCodeAt(i) & 0xff;
  }

  const filename = `candidats-paiements-${new Date().toISOString().slice(0, 10)}.pdf`;
  return saveOrDownloadBytes(filename, outBytes, "application/pdf");
}

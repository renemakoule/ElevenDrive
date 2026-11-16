import type { HistoriqueItem } from "../components/historique/types";
import { isTauri, invokeTauri } from "../db/httpClient";

interface ExportHistoriquePdfOptions {
  schoolName: string;
  roleTitle: string;
  totalActions: number;
  totalMontant: number;
  items: HistoriqueItem[];
}

function toWinAnsiBytes(text: string): number[] {
  const normalized = (text || "")
    .replace(/\u202F|\u00A0/g, " ")
    .replace(/[’‘]/g, "'")
    .replace(/[“”]/g, '"')
    .replace(/—|–/g, "-")
    .replace(/•/g, "-");
  const bytes: number[] = [];
  for (let i = 0; i < normalized.length; i++) {
    const code = normalized.charCodeAt(i);
    if (code === 0x20ac) {
      bytes.push(0x80); // € in WinAnsi
    } else if (code >= 0x20 && code <= 0xff) {
      bytes.push(code);
    } else {
      bytes.push(0x3f); // '?'
    }
  }
  return bytes;
}

export function pdfStr(text: string): string {
  const bytes = toWinAnsiBytes(text);
  let out = "(";
  for (const b of bytes) {
    if (b === 0x28 || b === 0x29 || b === 0x5c) {
      out += "\\" + String.fromCharCode(b);
    } else if (b >= 0x20 && b <= 0x7e) {
      out += String.fromCharCode(b);
    } else {
      out += "\\" + b.toString(8).padStart(3, "0");
    }
  }
  return out + ")";
}

export function wrapText(text: string, maxChars: number): string[] {
  const clean = (text || "").replace(/\s+/g, " ").trim();
  if (!clean) return [""];
  const words = clean.split(" ");
  const lines: string[] = [];
  let current = "";
  for (const w of words) {
    if (!current) {
      current = w.length > maxChars ? w.slice(0, maxChars) : w;
    } else if (current.length + 1 + w.length <= maxChars) {
      current += " " + w;
    } else {
      lines.push(current);
      current = w.length > maxChars ? w.slice(0, maxChars) : w;
    }
  }
  if (current) lines.push(current);
  return lines;
}

export async function saveOrDownloadBytes(
  filename: string,
  bytes: Uint8Array,
  mimeType: string,
): Promise<{ filename: string; savedPath?: string }> {
  if (isTauri) {
    const savedPath = await invokeTauri<string>(
      "sauvegarder_pdf_telechargements",
      {
        filename,
        bytes: Array.from(bytes),
      },
    );
    return { filename, savedPath };
  }

  const blob = new Blob([new Uint8Array(bytes).buffer], { type: mimeType });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = filename;
  a.style.display = "none";
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  setTimeout(() => URL.revokeObjectURL(url), 1500);
  return { filename };
}

export async function exportHistoriqueToPdf(
  opts: ExportHistoriquePdfOptions,
): Promise<{ filename: string; savedPath?: string }> {
  const pageWidth = 595.28; // A4
  const pageHeight = 841.89;
  const marginX = 36;
  const pagesContent: string[] = [];

  let ops: string[] = [];
  let y = pageHeight - 44;

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
    y = pageHeight - 44;
    drawText(
      marginX,
      y,
      14,
      true,
      `${opts.schoolName || "AutoSchool"} - Historique des operations`,
      "0.05 0.35 0.25",
    );
    y -= 16;
    const dateExport = new Date().toLocaleString("fr-FR");
    drawText(
      marginX,
      y,
      9,
      false,
      `Role : ${opts.roleTitle}   |   Exporte le : ${dateExport}   |   Page ${pageNum}`,
      "0.4 0.42 0.46",
    );
    y -= 14;
    drawText(
      marginX,
      y,
      9,
      true,
      `Total actions : ${opts.totalActions}   |   Cumul encaissements : ${opts.totalMontant.toLocaleString("fr-FR")} FCFA   |   Lignes exportees : ${opts.items.length}`,
      "0.15 0.17 0.2",
    );
    y -= 10;
    drawLine(y, "0.05 0.45 0.3");
    y -= 16;
  };

  let pageNum = 1;
  startHeader(pageNum);

  if (opts.items.length === 0) {
    drawText(
      marginX,
      y,
      10,
      false,
      "Aucune operation ne correspond aux filtres selectionnes.",
      "0.45 0.45 0.5",
    );
  } else {
    for (const item of opts.items) {
      const descLines = wrapText(item.description, 88);
      const neededHeight = 28 + descLines.length * 12;
      if (y - neededHeight < 44) {
        pagesContent.push(ops.join("\n"));
        ops = [];
        pageNum += 1;
        startHeader(pageNum);
      }

      const dateStr =
        item.dateStr ||
        (item.timestamp
          ? new Date(item.timestamp).toLocaleString("fr-FR")
          : "");
      const montantStr =
        item.montant && item.montant > 0
          ? ` | ${item.montant.toLocaleString("fr-FR")} FCFA`
          : "";
      drawText(
        marginX,
        y,
        9.5,
        true,
        `[${item.badgeLabel}] ${item.titre}${montantStr}`,
        "0.08 0.1 0.14",
      );
      if (dateStr) {
        drawText(
          pageWidth - marginX - 105,
          y,
          8,
          false,
          dateStr.slice(0, 20),
          "0.45 0.47 0.52",
        );
      }
      y -= 12;

      for (const line of descLines) {
        drawText(marginX + 8, y, 8.5, false, line, "0.28 0.3 0.35");
        y -= 11.5;
      }
      if (item.acteur) {
        drawText(
          marginX + 8,
          y,
          7.5,
          false,
          `Acteur : ${item.acteur}`,
          "0.48 0.5 0.55",
        );
        y -= 10;
      }
      drawLine(y + 2);
      y -= 10;
    }
  }

  pagesContent.push(ops.join("\n"));

  // Assemblage PDF 1.4
  const objects: string[] = [];
  const addObj = (body: string) => {
    objects.push(body);
    return objects.length;
  };

  const fontRegularId = addObj(
    "<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding /WinAnsiEncoding >>",
  );
  const fontBoldId = addObj(
    "<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica-Bold /Encoding /WinAnsiEncoding >>",
  );

  const pageIds: number[] = [];
  const pagesRootIdPlaceholder = objects.length + 1 + pagesContent.length * 2;

  for (const stream of pagesContent) {
    const contentId = addObj(
      `<< /Length ${stream.length} >>\nstream\n${stream}\nendstream`,
    );
    const pageId = addObj(
      `<< /Type /Page /Parent ${pagesRootIdPlaceholder} 0 R /MediaBox [0 0 ${pageWidth} ${pageHeight}] /Resources << /Font << /F1 ${fontRegularId} 0 R /F2 ${fontBoldId} 0 R >> >> /Contents ${contentId} 0 R >>`,
    );
    pageIds.push(pageId);
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

  const filename = `historique-activites-${new Date().toISOString().slice(0, 10)}.pdf`;
  return saveOrDownloadBytes(filename, outBytes, "application/pdf");
}

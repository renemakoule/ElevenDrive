/**
 * Extracteur de texte PDF haute fidélité basé sur Mozilla PDF.js (pdfjs-dist).
 * Décode tous les flux PDF (y compris FlateDecode, ObjStm, CMaps ToUnicode)
 * et extrait le texte page par page pour la synthèse vocale et la traduction.
 */
import * as pdfjsLib from 'pdfjs-dist';
import pdfWorker from 'pdfjs-dist/build/pdf.worker.min.mjs?url';

// Configuration du Worker PDF.js pour environnement Vite / Tauri / PWA
if (typeof window !== 'undefined' && pdfjsLib.GlobalWorkerOptions) {
  pdfjsLib.GlobalWorkerOptions.workerSrc = pdfWorker;
}

export interface PdfExtractionResult {
  texte: string;
  pages: string[];
  succes: boolean;
  source: 'FLUX_TEXTUEL' | 'METADONNEES_REPLI';
  nombreCaracteres: number;
}

export class PdfTextExtractor {
  /**
   * Extrait le texte complet d'un PDF sous forme de Data URL Base64, URL ou Uint8Array
   */
  public static async extraireTexte(
    base64OrDataUrl: string,
    repliMetadonnees?: { titre?: string; description?: string; categorie?: string }
  ): Promise<PdfExtractionResult> {
    try {
      if (!base64OrDataUrl || typeof base64OrDataUrl !== 'string') {
        return this.creerRepli(repliMetadonnees);
      }

      // 1. Décodage Base64 vers Uint8Array
      let bytes: Uint8Array;
      if (base64OrDataUrl.startsWith('data:') || base64OrDataUrl.includes('base64,')) {
        const dataUrlIdx = base64OrDataUrl.indexOf('base64,');
        const b64 = (dataUrlIdx !== -1 ? base64OrDataUrl.slice(dataUrlIdx + 7) : base64OrDataUrl).replace(/\s+/g, '');
        const binaryString = atob(b64);
        bytes = new Uint8Array(binaryString.length);
        for (let i = 0; i < binaryString.length; i++) {
          bytes[i] = binaryString.charCodeAt(i);
        }
      } else if (base64OrDataUrl.startsWith('http://') || base64OrDataUrl.startsWith('https://') || base64OrDataUrl.startsWith('blob:')) {
        const resp = await fetch(base64OrDataUrl);
        const arrayBuf = await resp.arrayBuffer();
        bytes = new Uint8Array(arrayBuf);
      } else {
        // Tentative atob directe
        const b64 = base64OrDataUrl.replace(/\s+/g, '');
        const binaryString = atob(b64);
        bytes = new Uint8Array(binaryString.length);
        for (let i = 0; i < binaryString.length; i++) {
          bytes[i] = binaryString.charCodeAt(i);
        }
      }

      // 2. Chargement via PDF.js
      const loadingTask = pdfjsLib.getDocument({
        data: bytes,
        useSystemFonts: true,
      });

      const doc = await loadingTask.promise;
      const pagesText: string[] = [];

      for (let pageNum = 1; pageNum <= doc.numPages; pageNum++) {
        const page = await doc.getPage(pageNum);
        const textContent = await page.getTextContent();
        let pageStr = '';

        for (const item of textContent.items) {
          if ('str' in item && typeof item.str === 'string') {
            const strVal = item.str;
            pageStr += strVal + ((item as any).hasEOL ? '\n' : ' ');
          }
        }

        const cleanPage = this.nettoyerTexte(pageStr);
        if (cleanPage.length > 0) {
          pagesText.push(cleanPage);
        }
      }

      const texteComplet = pagesText.join('\n\n');

      if (texteComplet.length >= 25) {
        return {
          texte: texteComplet,
          pages: pagesText,
          succes: true,
          source: 'FLUX_TEXTUEL',
          nombreCaracteres: texteComplet.length,
        };
      }

      // Si le PDF est un scan d'images ou sans texte extractible
      return this.creerRepli(repliMetadonnees);
    } catch (err) {
      console.warn('Extraction PDF via PDF.js échouée, bascule vers repli:', err);
      return this.creerRepli(repliMetadonnees);
    }
  }

  /**
   * Nettoie et formate le texte pour la lecture vocale et la traduction
   */
  public static nettoyerTexte(texte: string): string {
    if (!texte) return '';
    return texte
      .replace(/[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]/g, '') // Contrôles binaires
      .replace(/[ \t]+/g, ' ')
      .replace(/\r\n/g, '\n')
      .replace(/\n{3,}/g, '\n\n')
      .replace(/Page\s+\d+\s*(?:\/\s*\d+)?/gi, '') // "Page 1/4"
      .trim();
  }

  public static creerRepli(metas?: {
    titre?: string;
    description?: string;
    categorie?: string;
  }): PdfExtractionResult {
    const morceaux: string[] = [];
    if (metas?.titre) morceaux.push(metas.titre);
    if (metas?.categorie) morceaux.push(`Catégorie Permis ${metas.categorie}.`);
    if (metas?.description) morceaux.push(metas.description);

    if (morceaux.length === 0) {
      morceaux.push("Document d'apprentissage de l'auto-école.");
    }

    const texte = morceaux.join('. ');
    return {
      texte,
      pages: [texte],
      succes: false,
      source: 'METADONNEES_REPLI',
      nombreCaracteres: texte.length,
    };
  }
}

// ==============================================================================
// ElevenDrive - Service de Traduction DeepL API (Zero i18n, 100% Natif & Hybride)
// ==============================================================================

import { DEEPL_LOCAL_CACHE, type SupportedLanguage } from "./deepLDictionary";
import { API_BASE_URL } from "../config/env";
export type { SupportedLanguage };

export interface DeepLConfig {
  apiKey: string;
  isPro: boolean;
}

export interface LanguageMeta {
  code: SupportedLanguage;
  name: string;
  flag: string;
  deeplTarget: string;
  ttsLocale: string;
  dir: "ltr" | "rtl";
}

const getAuthHeaders = (): Record<string, string> => {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };
  if (typeof localStorage !== "undefined") {
    const raw =
      localStorage.getItem("elevendrive_active_session") ||
      localStorage.getItem("autoschool_active_session");
    if (raw) {
      try {
        const session = JSON.parse(raw);
        const token = session.token || session.access_token;
        if (token) {
          headers["Authorization"] = `Bearer ${token}`;
        }
      } catch {
        // Ignorer
      }
    }
  }
  return headers;
};

export const LANGUAGES_LIST: LanguageMeta[] = [
  {
    code: "fr",
    name: "Français",
    flag: "🇫🇷",
    deeplTarget: "FR",
    ttsLocale: "fr-FR",
    dir: "ltr",
  },
  {
    code: "en",
    name: "English",
    flag: "🇬🇧",
    deeplTarget: "EN-US",
    ttsLocale: "en-US",
    dir: "ltr",
  },
  {
    code: "es",
    name: "Español",
    flag: "🇪🇸",
    deeplTarget: "ES",
    ttsLocale: "es-ES",
    dir: "ltr",
  },
  {
    code: "ar",
    name: "العربية",
    flag: "🇸🇦",
    deeplTarget: "AR",
    ttsLocale: "ar-SA",
    dir: "rtl",
  },
];

class DeepLTranslationService {
  private dynamicCache: Record<string, Record<string, string>> = {};
  private apiKey: string = "";
  private isPro: boolean = false;
  private pendingOrProcessedTexts = new Set<string>();

  constructor() {
    this.chargerCacheLocal();
  }

  private chargerCacheLocal() {
    // 1. Clé issue des variables d'environnement (.env : VITE_DEEPL_API_KEY)
    try {
      const envKey = (import.meta.env?.VITE_DEEPL_API_KEY as string)?.trim();
      if (envKey) {
        this.apiKey = envKey;
        this.isPro = !envKey.endsWith(":fx");
      }
    } catch {
      // Ignorer si indisponible
    }

    // 2. Clé ou surcharge manuelle sauvegardée dans localStorage
    if (typeof localStorage !== "undefined") {
      try {
        const savedApiKey = localStorage.getItem("elevendrive_deepl_api_key");
        if (savedApiKey) {
          this.apiKey = savedApiKey.trim();
          this.isPro = !savedApiKey.endsWith(":fx");
        }
        const savedCache = localStorage.getItem("elevendrive_deepl_cache");
        if (savedCache) {
          this.dynamicCache = JSON.parse(savedCache);
        }
      } catch (err) {
        console.warn("Erreur chargement cache DeepL:", err);
      }
    }
  }

  private sauvegarderCacheLocal() {
    if (typeof localStorage !== "undefined") {
      try {
        localStorage.setItem(
          "elevendrive_deepl_cache",
          JSON.stringify(this.dynamicCache),
        );
      } catch (err) {
        console.warn("Erreur sauvegarde cache DeepL:", err);
      }
    }
  }

  public setApiKey(key: string) {
    this.apiKey = key.trim();
    this.isPro = !key.endsWith(":fx");
    if (typeof localStorage !== "undefined") {
      localStorage.setItem("elevendrive_deepl_api_key", this.apiKey);
    }
  }

  public getApiKey(): string {
    return this.apiKey;
  }

  /**
   * Vérifie formellement si une traduction existe déjà dans les dictionnaires locaux ou dynamiques
   */
  public hasTranslation(texte: string, cible: SupportedLanguage): boolean {
    if (!texte) return true;
    if (cible === "fr") return true;
    const trimmed = texte.trim();
    if (DEEPL_LOCAL_CACHE[cible]?.[trimmed] !== undefined) return true;
    if (this.dynamicCache[cible]?.[trimmed] !== undefined) return true;
    return false;
  }

  /**
   * Vérifie si un texte a déjà été demandé ou traité pour cette langue
   */
  public isPendingOrProcessed(
    texte: string,
    cible: SupportedLanguage,
  ): boolean {
    if (this.hasTranslation(texte, cible)) return true;
    return this.pendingOrProcessedTexts.has(`${cible}::${texte.trim()}`);
  }

  /**
   * Marque un texte comme en cours de traitement
   */
  public markPending(texte: string, cible: SupportedLanguage) {
    this.pendingOrProcessedTexts.add(`${cible}::${texte.trim()}`);
  }

  /**
   * Traduit un texte de façon synchrone et immédiate (utilise le cache certifié DeepL)
   */
  public traduireSynchrone(texte: string, cible: SupportedLanguage): string {
    if (!texte) return "";
    if (cible === "fr") return texte;

    // 1. Chercher dans le dictionnaire certifié local DeepL
    const dict = DEEPL_LOCAL_CACHE[cible];
    if (dict && dict[texte]) {
      return dict[texte];
    }

    // 2. Chercher dans le cache dynamique de session
    if (this.dynamicCache[cible] && this.dynamicCache[cible][texte]) {
      return this.dynamicCache[cible][texte];
    }

    // 3. Repli : si texte non traduit, retourner le texte original
    return texte;
  }

  /**
   * Traduction asynchrone appelant l'API officielle DeepL
   */
  public async traduireViaDeepL(
    texte: string,
    cible: SupportedLanguage,
  ): Promise<string> {
    if (!texte || !texte.trim()) return "";
    if (cible === "fr") return texte;

    // Vérifier si déjà présent en cache
    const enCache = this.traduireSynchrone(texte, cible);
    if (enCache !== texte) return enCache;

    // 1. Appel sécurisé via le Backend Go (la clé secrète DEEPL_API_KEY reste protégée côté serveur)
    try {
      const response = await fetch(`${API_BASE_URL}/traduire`, {
        method: "POST",
        headers: getAuthHeaders(),
        body: JSON.stringify({
          text: texte,
          target_lang: cible,
          source_lang: "FR",
        }),
      });

      if (response.ok) {
        const resJson = await response.json();
        const traduction =
          resJson?.data?.translation || resJson?.data?.translations?.[0];
        if (traduction && typeof traduction === "string" && traduction.trim()) {
          if (!this.dynamicCache[cible]) {
            this.dynamicCache[cible] = {};
          }
          this.dynamicCache[cible][texte] = traduction;
          this.sauvegarderCacheLocal();
          return traduction;
        }
      }
    } catch (backendErr) {
      // Le backend est potentiellement arrêté ou inaccessible en mode déconnecté
      console.warn(
        "Backend de traduction non joignable, tentative directe de secours:",
        backendErr,
      );
    }

    // 2. Repli de secours : appel direct si une clé API est configurée localement
    if (!this.apiKey) {
      return texte;
    }

    const meta = LANGUAGES_LIST.find((l) => l.code === cible);
    const targetCode = meta ? meta.deeplTarget : cible.toUpperCase();
    const endpoint = this.isPro
      ? "https://api.deepl.com/v2/translate"
      : "https://api-free.deepl.com/v2/translate";

    try {
      const response = await fetch(endpoint, {
        method: "POST",
        headers: {
          Authorization: `DeepL-Auth-Key ${this.apiKey}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          text: [texte],
          target_lang: targetCode,
          source_lang: "FR",
        }),
      });

      if (!response.ok) {
        throw new Error(`DeepL API HTTP ${response.status}`);
      }

      const data = await response.json();
      const traduction = data?.translations?.[0]?.text || texte;

      // Mettre en cache
      if (!this.dynamicCache[cible]) {
        this.dynamicCache[cible] = {};
      }
      this.dynamicCache[cible][texte] = traduction;
      this.sauvegarderCacheLocal();

      return traduction;
    } catch (err) {
      console.warn(`Échec DeepL API direct vers ${cible}:`, err);
      return texte;
    }
  }

  /**
   * Alias anglais pour la traduction asynchrone d'un texte unique
   */
  public async translateAsync(
    texte: string,
    cible: SupportedLanguage,
  ): Promise<string> {
    return this.traduireViaDeepL(texte, cible);
  }

  /**
   * Traduit automatiquement un lot de textes via l'API DeepL (Backend Go ou direct)
   * sans aucune traduction manuelle dans le code.
   */
  public async traduireLotViaBackend(
    textes: string[],
    cible: SupportedLanguage,
  ): Promise<Record<string, string>> {
    if (!textes || textes.length === 0 || cible === "fr") {
      return {};
    }

    // Filtrer les textes non vides et non encore traduits
    const uniques = Array.from(
      new Set(textes.map((t) => t.trim()).filter((t) => t.length > 0)),
    );
    const aTraduire = uniques.filter((t) => !this.hasTranslation(t, cible));

    if (aTraduire.length === 0) {
      return {};
    }

    // Bloquer les réémissions immédiates pendant le traitement
    aTraduire.forEach((t) => this.markPending(t, cible));

    const meta = LANGUAGES_LIST.find((l) => l.code === cible);
    const targetCode = meta ? meta.deeplTarget : cible.toUpperCase();
    const resultat: Record<string, string> = {};

    // Découpage par paquets de 30 pour respecter les limites HTTP et DeepL
    const TAILLE_LOT = 30;
    for (let i = 0; i < aTraduire.length; i += TAILLE_LOT) {
      const lotCourant = aTraduire.slice(i, i + TAILLE_LOT);

      // 1. Tenter l'appel groupé via le backend Go
      let lotTraduit = false;
      try {
        const response = await fetch(`${API_BASE_URL}/traduire`, {
          method: "POST",
          headers: getAuthHeaders(),
          body: JSON.stringify({
            texts: lotCourant,
            target_lang: targetCode,
            source_lang: "FR",
          }),
        });

        if (response.ok) {
          const resJson = await response.json();
          const translations: string[] = resJson?.data?.translations || [];
          if (Array.isArray(translations) && translations.length > 0) {
            if (!this.dynamicCache[cible]) {
              this.dynamicCache[cible] = {};
            }
            lotCourant.forEach((orig, idx) => {
              const trad = translations[idx];
              const val = trad && trad.trim() ? trad.trim() : orig;
              this.dynamicCache[cible][orig] = val;
              resultat[orig] = val;
            });
            lotTraduit = true;
          }
        }
      } catch (backendErr) {
        console.warn(
          "Backend /api/traduire non joignable pour le sous-lot, repli direct:",
          backendErr,
        );
      }

      // 2. Repli direct vers DeepL si nécessaire
      if (!lotTraduit && this.apiKey) {
        const endpoint = this.isPro
          ? "https://api.deepl.com/v2/translate"
          : "https://api-free.deepl.com/v2/translate";

        try {
          const response = await fetch(endpoint, {
            method: "POST",
            headers: {
              Authorization: `DeepL-Auth-Key ${this.apiKey}`,
              "Content-Type": "application/json",
            },
            body: JSON.stringify({
              text: lotCourant,
              target_lang: targetCode,
              source_lang: "FR",
            }),
          });

          if (response.ok) {
            const data = await response.json();
            const translations = data?.translations || [];
            if (!this.dynamicCache[cible]) {
              this.dynamicCache[cible] = {};
            }
            lotCourant.forEach((orig, idx) => {
              const trad = translations[idx]?.text;
              const val = trad && trad.trim() ? trad.trim() : orig;
              this.dynamicCache[cible][orig] = val;
              resultat[orig] = val;
            });
            lotTraduit = true;
          }
        } catch (directErr) {
          console.warn("Échec DeepL API direct pour le sous-lot:", directErr);
        }
      }

      // 3. Fallback de sécurité si échec
      if (!lotTraduit) {
        if (!this.dynamicCache[cible]) {
          this.dynamicCache[cible] = {};
        }
        lotCourant.forEach((orig) => {
          this.dynamicCache[cible][orig] = orig;
          resultat[orig] = orig;
        });
      }
    }

    this.sauvegarderCacheLocal();
    return resultat;
  }

  /**
   * Traduit une liste ordonnée de segments textuels (ex: phrases d'un PDF)
   * et retourne la liste dans le même ordre exact.
   */
  public async traduireSegments(
    segments: string[],
    cible: SupportedLanguage,
  ): Promise<string[]> {
    if (!segments || segments.length === 0 || cible === "fr") {
      return segments;
    }

    const tradMap = await this.traduireLotViaBackend(segments, cible);

    return segments.map((s) => {
      const trimmed = s.trim();
      return (
        tradMap[trimmed] ||
        this.dynamicCache[cible]?.[trimmed] ||
        DEEPL_LOCAL_CACHE[cible]?.[trimmed] ||
        s
      );
    });
  }
}

export const deepLTranslationService = new DeepLTranslationService();

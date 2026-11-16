// ==============================================================================
// ElevenDrive - Store Pinia de Gestion Multilingue (100% DeepL API, Zéro i18n)
// ==============================================================================

import { defineStore } from "pinia";
import {
  deepLTranslationService,
  LANGUAGES_LIST,
  type SupportedLanguage,
  type LanguageMeta,
} from "../services/deepLTranslationService";
import { speechAudioService } from "../services/speechAudioService";
import { useGuichetStore } from "./guichetStore";

// File d'attente globale pour regrouper les textes à traduire en un seul lot
const pendingTextsQueue = new Set<string>();
let batchDebounceTimer: any = null;

export const useLanguageStore = defineStore("language", {
  state: () => ({
    currentLanguage: (typeof localStorage !== "undefined" &&
    (localStorage.getItem("elevendrive_lang") as SupportedLanguage)
      ? (localStorage.getItem("elevendrive_lang") as SupportedLanguage)
      : "fr") as SupportedLanguage,
    isTranslating: false,
    isDropdownOpen: false,
    translationsRevision: 0, // Compteur incrémenté pour déclencher la réactivité automatique dans les composants
  }),

  getters: {
    currentMeta(state): LanguageMeta {
      return (
        LANGUAGES_LIST.find((l) => l.code === state.currentLanguage) ||
        LANGUAGES_LIST[0]
      );
    },
    availableLanguages(): LanguageMeta[] {
      return LANGUAGES_LIST;
    },
    isRTL(state): boolean {
      return state.currentLanguage === "ar";
    },
  },

  actions: {
    /**
     * Initialise la langue active au démarrage de l'application
     */
    init() {
      if (typeof window === "undefined") return;

      const savedLang = localStorage.getItem(
        "elevendrive_lang",
      ) as SupportedLanguage | null;
      if (savedLang && LANGUAGES_LIST.some((l) => l.code === savedLang)) {
        this.currentLanguage = savedLang;
      } else {
        this.currentLanguage = "fr";
      }

      this.applyDomSettings();
      speechAudioService.selectVoiceForLanguage(this.currentLanguage);
    },

    /**
     * Change la langue de l'interface et applique les changements réactivement
     */
    setLanguage(lang: SupportedLanguage) {
      if (!LANGUAGES_LIST.some((l) => l.code === lang)) return;

      this.currentLanguage = lang;
      this.isDropdownOpen = false;
      this.translationsRevision++;

      if (typeof localStorage !== "undefined") {
        localStorage.setItem("elevendrive_lang", lang);
      }

      this.applyDomSettings();
      speechAudioService.selectVoiceForLanguage(lang);

      try {
        const guichetStore = useGuichetStore();
        const meta = this.currentMeta;
        guichetStore.showToast(
          `${meta.flag} ${this.t("Langue")} : ${meta.name}`,
          "info",
        );
      } catch {
        // Silencieux si hors contexte store
      }
    },

    /**
     * Bascule l'état d'ouverture du sélecteur de langue
     */
    toggleDropdown() {
      this.isDropdownOpen = !this.isDropdownOpen;
    },

    closeDropdown() {
      this.isDropdownOpen = false;
    },

    /**
     * Traduit un texte de façon réactive et immédiate.
     * Si le texte n'a pas encore été traduit, déclenche automatiquement la traduction
     * groupée via le backend DeepL et réactualise l'interface dès réception !
     */
    t(texte: string): string {
      if (!texte || typeof texte !== "string") return "";
      const trimmed = texte.trim();
      if (!trimmed || this.currentLanguage === "fr") return texte;

      // Dépendance réactive Vue 3
      const _ = this.translationsRevision;

      // 1. Chercher dans le cache certifié ou dynamique
      if (deepLTranslationService.hasTranslation(trimmed, this.currentLanguage)) {
        return deepLTranslationService.traduireSynchrone(
          trimmed,
          this.currentLanguage,
        );
      }

      // 2. Déclencher uniquement si non déjà demandé ou en cours
      if (!deepLTranslationService.isPendingOrProcessed(trimmed, this.currentLanguage)) {
        this.demanderTraduction(trimmed);
      }

      // 3. Renvoyer le texte original en attendant la réponse de DeepL (évite les écrans blancs)
      return texte;
    },

    /**
     * Enregistre un texte dans la file d'attente pour traduction automatique groupée (Batch)
     */
    demanderTraduction(texte: string) {
      if (this.currentLanguage === "fr") return;
      const trimmed = texte.trim();
      if (!trimmed) return;

      if (deepLTranslationService.hasTranslation(trimmed, this.currentLanguage)) {
        return;
      }
      if (deepLTranslationService.isPendingOrProcessed(trimmed, this.currentLanguage)) {
        return;
      }

      deepLTranslationService.markPending(trimmed, this.currentLanguage);
      pendingTextsQueue.add(trimmed);

      if (batchDebounceTimer) {
        clearTimeout(batchDebounceTimer);
      }

      batchDebounceTimer = setTimeout(async () => {
        const textes = Array.from(pendingTextsQueue);
        pendingTextsQueue.clear();
        if (textes.length === 0) return;

        this.isTranslating = true;
        try {
          const nouveaux = await deepLTranslationService.traduireLotViaBackend(
            textes,
            this.currentLanguage,
          );
          if (Object.keys(nouveaux).length > 0) {
            // Déclenche la mise à jour réactive de tous les composants qui ont appelé t(texte)
            this.translationsRevision++;
          }
        } finally {
          this.isTranslating = false;
        }
      }, 100);
    },

    /**
     * Traduction asynchrone appelant l'API officielle DeepL pour les textes dynamiques
     */
    async translateAsync(texte: string): Promise<string> {
      if (!texte || !texte.trim()) return "";
      this.isTranslating = true;
      try {
        const resultat = await deepLTranslationService.traduireViaDeepL(
          texte,
          this.currentLanguage,
        );
        this.translationsRevision++;
        return resultat;
      } finally {
        this.isTranslating = false;
      }
    },

    /**
     * Configure la clé API DeepL
     */
    setDeepLApiKey(key: string) {
      deepLTranslationService.setApiKey(key);
    },

    /**
     * Applique les attributs dir (RTL/LTR) et lang sur le document HTML
     */
    applyDomSettings() {
      if (typeof document === "undefined") return;
      const isRtl = this.currentLanguage === "ar";
      document.documentElement.setAttribute("dir", isRtl ? "rtl" : "ltr");
      document.documentElement.setAttribute("lang", this.currentLanguage);
    },
  },
});

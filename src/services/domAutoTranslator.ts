// ==============================================================================
// ElevenDrive - Moteur Universel de Traduction Automatique du DOM (DeepL API)
// Traduit TOUS les composants de l'application automatiquement sans code en dur.
// ==============================================================================

import { useLanguageStore } from "../stores/languageStore";
import { deepLTranslationService } from "./deepLTranslationService";

class DomAutoTranslator {
  private observer: MutationObserver | null = null;
  private isApplyingTranslations = false;
  private scanDebounceTimer: any = null;
  private originalTextMap = new WeakMap<Node, string>();
  private originalAttrMap = new WeakMap<Element, Record<string, string>>();
  private translatedTextMap = new WeakMap<Node, string>();
  private initialized = false;

  private ignoredTags = new Set([
    "SCRIPT",
    "STYLE",
    "NOSCRIPT",
    "CODE",
    "PRE",
    "SVG",
    "CANVAS",
    "AUDIO",
    "VIDEO",
  ]);

  /**
   * Vérifie si un élément ou ses parents doivent être ignorés (code, logos, montants, matricules)
   */
  private isIgnoredElement(el: Element | null): boolean {
    let curr: Element | null = el;
    while (curr && curr !== document.body) {
      if (this.ignoredTags.has(curr.tagName)) return true;
      if (curr.hasAttribute("data-no-translate")) return true;
      if (curr.classList && curr.classList.contains("no-translate")) return true;
      // Ne pas altérer les matricules, identifiants techniques et codes mono
      if (curr.classList && curr.classList.contains("font-mono")) return true;
      curr = curr.parentElement;
    }
    return false;
  }

  /**
   * Filtre les chaînes qui ne doivent pas être traduites (nombres, dates, devises, URLs)
   */
  private shouldTranslate(str: string): boolean {
    if (!str) return false;
    const trimmed = str.trim();
    if (trimmed.length < 2) return false;

    // Ignorer si uniquement des chiffres, dates, symboles ou ponctuation
    if (/^[\d\s.,;:!?%€$£+\-/*#@()\[\]{}]+$/.test(trimmed)) return false;

    // Ignorer les formats de prix purs "15 000 FCFA", "0 F", "+25000 XOF"
    if (/^\+?[\d\s.,]+(FCFA|F|XOF|EUR|USD)$/i.test(trimmed)) return false;

    // Ignorer les URLs, emails, extensions de fichiers
    if (/^(https?:\/\/|[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}|\/api\/|\.[a-z0-9]{2,4})/i.test(trimmed)) {
      return false;
    }

    // Doit contenir au moins deux lettres
    const lettres = trimmed.match(/[a-zA-Z\u00C0-\u024F\u0600-\u06FF]/g);
    return Boolean(lettres && lettres.length >= 2);
  }

  /**
   * Initialise l'observateur automatique sur l'arbre DOM de l'application
   */
  public init() {
    if (typeof window === "undefined" || this.initialized) return;
    this.initialized = true;

    // Écouter UNIQUEMENT les vrais changements de langue ou nouvelles traductions reçues
    const languageStore = useLanguageStore();
    languageStore.$subscribe((mutation) => {
      const events = Array.isArray((mutation as any).events)
        ? (mutation as any).events
        : [(mutation as any).events];
      const hasRelevantChange = events.some(
        (e: any) => e?.key === "currentLanguage" || e?.key === "translationsRevision"
      );
      if (hasRelevantChange) {
        this.scheduleScan();
      }
    });

    // Observer les mutations réelles du DOM (nouveaux éléments ou modales ouvertes)
    const root = document.getElementById("app") || document.body;
    this.observer = new MutationObserver((mutations) => {
      let shouldRescan = false;
      for (const m of mutations) {
        if (m.type === "childList" && (m.addedNodes.length > 0 || m.removedNodes.length > 0)) {
          shouldRescan = true;
          break;
        }
        if (m.type === "characterData" && m.target) {
          const applied = this.translatedTextMap.get(m.target);
          if (applied !== undefined && m.target.nodeValue === applied) {
            // Ignorer notre propre mise à jour du texte traduit
            continue;
          }
          const orig = this.originalTextMap.get(m.target);
          if (!orig || orig !== m.target.nodeValue) {
            shouldRescan = true;
            break;
          }
        }
      }

      if (shouldRescan) {
        this.scheduleScan();
      }
    });

    this.observer.observe(root, {
      childList: true,
      subtree: true,
      characterData: true,
      attributes: true,
      attributeFilter: ["placeholder", "title"],
    });

    // Premier scan au démarrage
    this.scheduleScan();
  }

  /**
   * Planifie un scan dérebondé pour éviter de surcharger les cycles de rendu
   */
  public scheduleScan() {
    if (this.scanDebounceTimer) {
      clearTimeout(this.scanDebounceTimer);
    }
    this.scanDebounceTimer = setTimeout(() => {
      this.scanAndTranslate();
    }, 40);
  }

  /**
   * Scanne tous les nœuds de texte et attributs traduisibles du DOM
   */
  public scanAndTranslate(rootNode?: Node) {
    if (typeof window === "undefined") return;

    const languageStore = useLanguageStore();
    const currentLang = languageStore.currentLanguage;
    const isFrench = currentLang === "fr";

    const targetRoot = rootNode || document.getElementById("app") || document.body;
    const walker = document.createTreeWalker(
      targetRoot,
      NodeFilter.SHOW_TEXT | NodeFilter.SHOW_ELEMENT,
      {
        acceptNode: (node) => {
          if (node.nodeType === Node.ELEMENT_NODE) {
            const el = node as Element;
            if (this.isIgnoredElement(el)) {
              return NodeFilter.FILTER_REJECT;
            }
            return NodeFilter.FILTER_ACCEPT;
          }
          if (node.nodeType === Node.TEXT_NODE) {
            const parent = node.parentElement;
            if (this.isIgnoredElement(parent)) {
              return NodeFilter.FILTER_REJECT;
            }
            return NodeFilter.FILTER_ACCEPT;
          }
          return NodeFilter.FILTER_SKIP;
        },
      }
    );

    const pendingBatchTexts = new Set<string>();
    this.isApplyingTranslations = true;

    try {
      let current: Node | null = walker.currentNode;
      while (current) {
        // 1. Traduction des nœuds de texte
        if (current.nodeType === Node.TEXT_NODE) {
          const raw = current.nodeValue || "";
          let original = this.originalTextMap.get(current);

          if (!original) {
            if (this.shouldTranslate(raw)) {
              original = raw;
              this.originalTextMap.set(current, raw);
            }
          }

          if (original) {
            const leading = original.match(/^\s*/)?.[0] || "";
            const trailing = original.match(/\s*$/)?.[0] || "";
            const trimmed = original.trim();

            if (isFrench) {
              if (current.nodeValue !== original) {
                current.nodeValue = original;
              }
            } else {
              if (deepLTranslationService.hasTranslation(trimmed, currentLang)) {
                const traduit = deepLTranslationService.traduireSynchrone(trimmed, currentLang);
                const finalStr = `${leading}${traduit}${trailing}`;
                if (current.nodeValue !== finalStr) {
                  this.translatedTextMap.set(current, finalStr);
                  current.nodeValue = finalStr;
                }
              } else if (!deepLTranslationService.isPendingOrProcessed(trimmed, currentLang)) {
                pendingBatchTexts.add(trimmed);
              }
            }
          }
        }

        // 2. Traduction des attributs d'accessibilité et de formulaire (placeholder, title)
        if (current.nodeType === Node.ELEMENT_NODE) {
          const el = current as Element;
          if (el.tagName === "INPUT" || el.tagName === "TEXTAREA" || el.tagName === "BUTTON") {
            const placeholder = el.getAttribute("placeholder");
            const title = el.getAttribute("title");

            let attrStorage = this.originalAttrMap.get(el);
            if (!attrStorage) {
              attrStorage = {};
              if (placeholder && this.shouldTranslate(placeholder)) {
                attrStorage["placeholder"] = placeholder;
              }
              if (title && this.shouldTranslate(title)) {
                attrStorage["title"] = title;
              }
              if (Object.keys(attrStorage).length > 0) {
                this.originalAttrMap.set(el, attrStorage);
              }
            }

            if (attrStorage) {
              for (const [attrName, origVal] of Object.entries(attrStorage)) {
                const trimmedAttr = origVal.trim();
                if (isFrench) {
                  if (el.getAttribute(attrName) !== origVal) {
                    el.setAttribute(attrName, origVal);
                  }
                } else {
                  if (deepLTranslationService.hasTranslation(trimmedAttr, currentLang)) {
                    const traduit = deepLTranslationService.traduireSynchrone(trimmedAttr, currentLang);
                    if (el.getAttribute(attrName) !== traduit) {
                      el.setAttribute(attrName, traduit);
                    }
                  } else if (!deepLTranslationService.isPendingOrProcessed(trimmedAttr, currentLang)) {
                    pendingBatchTexts.add(trimmedAttr);
                  }
                }
              }
            }
          }
        }

        current = walker.nextNode();
      }

      // Si des textes non encore traduits ont été découverts dans le DOM,
      // on les transmet directement au batch automatique du backend Go !
      if (!isFrench && pendingBatchTexts.size > 0) {
        pendingBatchTexts.forEach((txt) => {
          languageStore.demanderTraduction(txt);
        });
      }
    } finally {
      this.isApplyingTranslations = false;
    }
  }

  /**
   * Nettoie les écouteurs lors de la destruction
   */
  public destroy() {
    if (this.observer) {
      this.observer.disconnect();
      this.observer = null;
    }
    this.initialized = false;
  }
}

export const domAutoTranslator = new DomAutoTranslator();

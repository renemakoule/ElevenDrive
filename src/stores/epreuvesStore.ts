import { defineStore } from "pinia";
import type { EpreuvePdf } from "../types";
import { api } from "../db/api";
import { useAuthStore } from "./authStore";

const STORAGE_PREFIX = "autoschool_epreuves_cache_";

function normalizePermisCat(raw?: string): string {
  const s = (raw || "B").toUpperCase().trim();
  if (s === "TOUTES" || s === "ALL" || s === "TOUS" || s.includes("TRONC"))
    return "TOUTES";
  if (s.includes("A") || s.includes("MOTO")) return "A";
  if (s.includes("C") || s.includes("POIDS")) return "C";
  if (s.includes("D") || s.includes("TRANSPORT") || s.includes("BUS"))
    return "D";
  if (s.includes("E") || s.includes("REMORQUE")) return "E";
  return "B";
}

function safeSaveEpreuvesCache(key: string, list: EpreuvePdf[]) {
  if (typeof localStorage === "undefined") return;
  try {
    localStorage.setItem(key, JSON.stringify(list));
  } catch {
    try {
      // Si le PDF Base64 dépasse le quota localStorage (~5 Mo), sauvegarder les métadonnées légères sans lever d'exception
      const light = list.map((e) => ({
        ...e,
        fichier_data:
          e.fichier_data && e.fichier_data.length > 120000
            ? ""
            : e.fichier_data,
      }));
      localStorage.setItem(key, JSON.stringify(light));
    } catch {
      // Ignorer silencieusement le dépassement de quota
    }
  }
}

export const useEpreuvesStore = defineStore("epreuves", {
  state: () => ({
    epreuvesList: [] as EpreuvePdf[],
    selectedCategorieFilter: "ALL",
    searchQuery: "",
    onlyWithCorrection: false,
    activePdfToPreview: null as EpreuvePdf | null,
    isUploadModalOpen: false,
    isLoading: false,
    error: null as string | null,
  }),

  getters: {
    /**
     * Épreuves réelles filtrées pour l'administration (Directeur / Secrétaire)
     */
    filteredEpreuves(state): EpreuvePdf[] {
      return state.epreuvesList.filter((ep) => {
        // Filtre catégorie
        if (
          state.selectedCategorieFilter !== "ALL" &&
          ep.categorie_permis !== state.selectedCategorieFilter
        ) {
          return false;
        }
        // Filtre corrigé
        if (state.onlyWithCorrection && !ep.avec_corrige) {
          return false;
        }
        // Filtre recherche textuelle
        if (state.searchQuery.trim()) {
          const q = state.searchQuery.toLowerCase();
          const matchTitle = ep.titre.toLowerCase().includes(q);
          const matchDesc = ep.description?.toLowerCase().includes(q) || false;
          const matchSession = ep.session?.toLowerCase().includes(q) || false;
          const matchCat = ep.categorie_permis.toLowerCase().includes(q);
          if (!matchTitle && !matchDesc && !matchSession && !matchCat) {
            return false;
          }
        }
        return true;
      });
    },

    /**
     * Statistiques réelles calculées depuis les enregistrements de la base de données
     */
    stats(state) {
      const total = state.epreuvesList.length;
      const permisB = state.epreuvesList.filter(
        (e) => e.categorie_permis === "B",
      ).length;
      const permisA = state.epreuvesList.filter(
        (e) => e.categorie_permis === "A" || e.categorie_permis === "A1",
      ).length;
      const poidsLourd = state.epreuvesList.filter((e) =>
        ["C", "D", "E"].includes(e.categorie_permis),
      ).length;
      const avecCorrige = state.epreuvesList.filter(
        (e) => e.avec_corrige,
      ).length;
      const totalDownloads = state.epreuvesList.reduce(
        (acc, curr) => acc + (curr.telechargements_count || 0),
        0,
      );

      return {
        total,
        permisB,
        permisA,
        poidsLourd,
        avecCorrige,
        totalDownloads,
      };
    },
  },

  actions: {
    /**
     * Chargement réel depuis la base de données PostgreSQL (API Go) avec fallback cache local
     */
    async loadEpreuves() {
      this.isLoading = true;
      this.error = null;
      const authStore = useAuthStore();
      const schoolId = authStore.currentSchool?.id || "default";
      const cacheKey = `${STORAGE_PREFIX}${schoolId}`;

      try {
        // 1. Appel vers la base de données via l'API Go
        const data = await api.getEpreuves();
        this.epreuvesList = data || [];

        // 2. Mise à jour sécurisée du cache local (sans risque de QuotaExceededError)
        safeSaveEpreuvesCache(cacheKey, this.epreuvesList);
      } catch (err: any) {
        this.error = err.message || "Impossible de joindre la base de données.";
        // 3. Fallback de secours sur le cache local uniquement si la liste mémoire est vide
        if (
          this.epreuvesList.length === 0 &&
          typeof localStorage !== "undefined"
        ) {
          const cached = localStorage.getItem(cacheKey);
          if (cached) {
            try {
              const parsed = JSON.parse(cached);
              if (Array.isArray(parsed)) {
                this.epreuvesList = parsed;
              }
            } catch {
              this.epreuvesList = [];
            }
          }
        }
      } finally {
        this.isLoading = false;
      }
    },

    /**
     * Insertion réelle d'une épreuve en base de données PostgreSQL
     */
    async addEpreuve(
      payload: Omit<EpreuvePdf, "id" | "date_ajout" | "telechargements_count">,
    ): Promise<EpreuvePdf> {
      this.isLoading = true;
      this.error = null;
      const authStore = useAuthStore();
      const schoolId = authStore.currentSchool?.id || "default";
      const cacheKey = `${STORAGE_PREFIX}${schoolId}`;

      try {
        // Appel direct vers l'API Go & PostgreSQL
        const created = await api.creerEpreuve({
          titre: payload.titre,
          categorie_permis: payload.categorie_permis,
          description: payload.description,
          annee: payload.annee,
          session: payload.session,
          duree_minutes: payload.duree_minutes,
          nombre_questions: payload.nombre_questions,
          avec_corrige: payload.avec_corrige,
          fichier_nom: payload.fichier_nom,
          fichier_taille: payload.fichier_taille,
          fichier_data: payload.fichier_data,
          ajoute_par: payload.ajoute_par,
        });

        this.epreuvesList = [
          created,
          ...this.epreuvesList.filter((e) => e.id !== created.id),
        ];

        // Sauvegarde sécurisée dans le cache local
        safeSaveEpreuvesCache(cacheKey, this.epreuvesList);

        return created;
      } catch (err: any) {
        this.error =
          err.message || "Erreur lors de l'enregistrement en base de données.";
        throw err;
      } finally {
        this.isLoading = false;
      }
    },

    /**
     * Suppression réelle d'une épreuve en base de données PostgreSQL
     */
    async deleteEpreuve(id: string): Promise<void> {
      this.isLoading = true;
      const authStore = useAuthStore();
      const schoolId = authStore.currentSchool?.id || "default";
      const cacheKey = `${STORAGE_PREFIX}${schoolId}`;

      try {
        await api.supprimerEpreuve(id);
        this.epreuvesList = this.epreuvesList.filter((e) => e.id !== id);
        safeSaveEpreuvesCache(cacheKey, this.epreuvesList);
      } catch (err: any) {
        this.error = err.message || "Erreur lors de la suppression en base.";
        throw err;
      } finally {
        this.isLoading = false;
      }
    },

    /**
     * Incrémente le compteur de téléchargements en base de données
     */
    async incrementDownload(id: string): Promise<void> {
      const ep = this.epreuvesList.find((e) => e.id === id);
      if (ep) {
        ep.telechargements_count = (ep.telechargements_count || 0) + 1;
      }
      try {
        await api.incrementerTelechargementEpreuve(id);
      } catch {
        // Ignorer l'erreur d'incrémentation en cas de micro-coupure réseau
      }
    },

    /**
     * FILTRE POUR L'ÉLÈVE :
     * Renvoie les épreuves correspondant à la catégorie normalisée de l'élève ou 'TOUTES'
     * (et replie sur l'ensemble des épreuves de l'école si aucune ne correspond exactement)
     */
    getEpreuvesForStudent(categorie: string): EpreuvePdf[] {
      const targetCat = normalizePermisCat(categorie);
      if (targetCat === "TOUTES") return this.epreuvesList;
      const matched = this.epreuvesList.filter((ep) => {
        const epCat = normalizePermisCat(ep.categorie_permis);
        return epCat === targetCat || epCat === "TOUTES";
      });
      return matched.length > 0 ? matched : this.epreuvesList;
    },

    /**
     * Déclenche le téléchargement direct d'un PDF réel
     */
    downloadPdf(epreuve: EpreuvePdf) {
      this.incrementDownload(epreuve.id);
      const link = document.createElement("a");
      link.href = epreuve.fichier_data;
      link.download = epreuve.fichier_nom.endsWith(".pdf")
        ? epreuve.fichier_nom
        : `${epreuve.fichier_nom}.pdf`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    },

    /**
     * Ouvre la modale d'aperçu d'un PDF réel
     */
    previewPdf(epreuve: EpreuvePdf) {
      this.activePdfToPreview = epreuve;
    },

    /**
     * Ferme la modale d'aperçu
     */
    closePreview() {
      this.activePdfToPreview = null;
    },
  },
});

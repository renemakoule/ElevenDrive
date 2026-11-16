import { defineStore } from "pinia";
import type { SupportCours } from "../types";
import { api } from "../db/api";
import { useAuthStore } from "./authStore";
import { useGuichetStore } from "./guichetStore";
import { openExternalLink } from "../utils/urlHelper";

const STORAGE_PREFIX = "autoschool_supports_cours_cache_";

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

function safeSaveSupportsCache(key: string, list: SupportCours[]) {
  if (typeof localStorage === "undefined") return;
  try {
    localStorage.setItem(key, JSON.stringify(list));
  } catch {
    try {
      const light = list.map((s) => ({
        ...s,
        fichier_data:
          s.fichier_data && s.fichier_data.length > 120000
            ? ""
            : s.fichier_data,
      }));
      localStorage.setItem(key, JSON.stringify(light));
    } catch {
      // Ignorer silencieusement le dépassement de quota
    }
  }
}

export const useSupportsCoursStore = defineStore("supportsCours", {
  state: () => ({
    supportsList: [] as SupportCours[],
    selectedCategorieFilter: "ALL",
    selectedTypeFilter: "ALL" as "ALL" | "PDF" | "LIEN_WEB",
    searchQuery: "",
    activePdfToPreview: null as SupportCours | null,
    isUploadModalOpen: false,
    isLoading: false,
    error: null as string | null,
  }),

  getters: {
    filteredSupports(state): SupportCours[] {
      return state.supportsList.filter((item) => {
        // Filtre catégorie
        if (
          state.selectedCategorieFilter !== "ALL" &&
          item.categorie_permis !== state.selectedCategorieFilter &&
          item.categorie_permis !== "TOUTES"
        ) {
          return false;
        }
        // Filtre type (PDF vs Lien web)
        if (
          state.selectedTypeFilter !== "ALL" &&
          item.type_support !== state.selectedTypeFilter
        ) {
          return false;
        }
        // Filtre recherche
        if (state.searchQuery.trim()) {
          const q = state.searchQuery.toLowerCase();
          const matchTitle = item.titre.toLowerCase().includes(q);
          const matchDesc =
            item.description?.toLowerCase().includes(q) || false;
          const matchCat = item.categorie_permis.toLowerCase().includes(q);
          if (!matchTitle && !matchDesc && !matchCat) {
            return false;
          }
        }
        return true;
      });
    },

    stats(state) {
      const total = state.supportsList.length;
      const pdfCount = state.supportsList.filter(
        (s) => s.type_support === "PDF",
      ).length;
      const lienCount = state.supportsList.filter(
        (s) => s.type_support === "LIEN_WEB",
      ).length;
      const permisB = state.supportsList.filter(
        (s) => s.categorie_permis === "B",
      ).length;
      const permisA = state.supportsList.filter(
        (s) => s.categorie_permis === "A" || s.categorie_permis === "A1",
      ).length;
      const poidsLourd = state.supportsList.filter((s) =>
        ["C", "D", "E"].includes(s.categorie_permis),
      ).length;
      const totalDownloads = state.supportsList.reduce(
        (acc, curr) => acc + (curr.telechargements_count || 0),
        0,
      );

      return {
        total,
        pdfCount,
        lienCount,
        permisB,
        permisA,
        poidsLourd,
        totalDownloads,
      };
    },
  },

  actions: {
    async loadSupports() {
      this.isLoading = true;
      this.error = null;
      const authStore = useAuthStore();
      const schoolId = authStore.currentSchool?.id || "default";
      const cacheKey = `${STORAGE_PREFIX}${schoolId}`;

      try {
        const data = await api.getSupportsCours();
        this.supportsList = data || [];

        safeSaveSupportsCache(cacheKey, this.supportsList);
      } catch (err: any) {
        this.error = err.message || "Impossible de joindre la base de données.";
        if (
          this.supportsList.length === 0 &&
          typeof localStorage !== "undefined"
        ) {
          const cached = localStorage.getItem(cacheKey);
          if (cached) {
            try {
              const parsed = JSON.parse(cached);
              if (Array.isArray(parsed)) {
                this.supportsList = parsed;
              }
            } catch {
              // ignore
            }
          }
        }
      } finally {
        this.isLoading = false;
      }
    },

    async ajouterSupport(payload: {
      titre: string;
      categorie_permis: string;
      description?: string;
      type_support: "PDF" | "LIEN_WEB";
      fichier_nom?: string;
      fichier_taille?: string;
      fichier_data?: string;
      url_externe?: string;
    }) {
      this.isLoading = true;
      const authStore = useAuthStore();
      const guichetStore = useGuichetStore();

      const auteur = authStore.currentUser?.nom
        ? `${authStore.currentUser.prenom} ${authStore.currentUser.nom}`
        : "Direction Pédagogique";

      try {
        const nouveauSupport = await api.creerSupportCours({
          ...payload,
          ajoute_par: auteur,
        });

        this.supportsList = [
          nouveauSupport,
          ...this.supportsList.filter((s) => s.id !== nouveauSupport.id),
        ];

        const schoolId = authStore.currentSchool?.id || "default";
        safeSaveSupportsCache(
          `${STORAGE_PREFIX}${schoolId}`,
          this.supportsList,
        );

        guichetStore.showToast(
          "Support de cours ajouté avec succès !",
          "success",
        );
        this.isUploadModalOpen = false;
        return nouveauSupport;
      } catch (err: any) {
        guichetStore.showToast(
          `Erreur ajout support : ${err.message}`,
          "error",
        );
        throw err;
      } finally {
        this.isLoading = false;
      }
    },

    async supprimerSupport(id: string) {
      const guichetStore = useGuichetStore();
      try {
        await api.supprimerSupportCours(id);
        this.supportsList = this.supportsList.filter((s) => s.id !== id);

        const authStore = useAuthStore();
        const schoolId = authStore.currentSchool?.id || "default";
        safeSaveSupportsCache(
          `${STORAGE_PREFIX}${schoolId}`,
          this.supportsList,
        );

        guichetStore.showToast(
          "Support de cours supprimé avec succès.",
          "info",
        );
      } catch (err: any) {
        guichetStore.showToast(`Erreur suppression : ${err.message}`, "error");
        throw err;
      }
    },

    previewPdf(support: SupportCours) {
      if (support.type_support === "LIEN_WEB" && support.url_externe) {
        openExternalLink(support.url_externe);
        return;
      }
      this.activePdfToPreview = support;
    },

    closePreview() {
      this.activePdfToPreview = null;
    },

    async downloadPdf(support: SupportCours) {
      const guichetStore = useGuichetStore();
      if (support.type_support === "LIEN_WEB" && support.url_externe) {
        openExternalLink(support.url_externe);
        return;
      }

      if (!support.fichier_data) {
        guichetStore.showToast("Données du document introuvables.", "error");
        return;
      }

      try {
        const link = document.createElement("a");
        link.href = support.fichier_data;
        link.download = support.fichier_nom || `${support.titre}.pdf`;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);

        await api.incrementerTelechargementSupport(support.id);
        if (typeof support.telechargements_count === "number") {
          support.telechargements_count++;
        }
        guichetStore.showToast(
          `Téléchargement de "${support.titre}" démarré.`,
          "info",
        );
      } catch (err: any) {
        guichetStore.showToast(
          `Erreur téléchargement : ${err.message}`,
          "error",
        );
      }
    },

    getSupportsForStudent(studentCategory: string): SupportCours[] {
      const targetCat = normalizePermisCat(studentCategory);
      if (targetCat === "TOUTES") return this.supportsList;
      const matched = this.supportsList.filter((s) => {
        const sCat = normalizePermisCat(s.categorie_permis);
        return sCat === targetCat || sCat === "TOUTES";
      });
      return matched.length > 0 ? matched : this.supportsList;
    },
  },
});

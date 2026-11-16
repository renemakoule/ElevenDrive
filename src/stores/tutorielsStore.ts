import { defineStore } from "pinia";
import type { TutorielVideo } from "../types";
import { api } from "../db/api";
import { useAuthStore } from "./authStore";
import { useGuichetStore } from "./guichetStore";

const STORAGE_PREFIX = "autoschool_tutoriels_cache_";

export const useTutorielsStore = defineStore("tutoriels", {
  state: () => ({
    tutorielsList: [] as TutorielVideo[],
    selectedCategorieFilter: "ALL",
    searchQuery: "",
    activeVideoToPlay: null as TutorielVideo | null,
    isUploadModalOpen: false,
    isLoading: false,
    error: null as string | null,
  }),

  getters: {
    filteredTutoriels(state): TutorielVideo[] {
      return state.tutorielsList.filter((item) => {
        if (
          state.selectedCategorieFilter !== "ALL" &&
          item.categorie_permis !== state.selectedCategorieFilter &&
          item.categorie_permis !== "TOUTES"
        ) {
          return false;
        }
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
      const total = state.tutorielsList.length;
      const permisB = state.tutorielsList.filter(
        (v) => v.categorie_permis === "B",
      ).length;
      const permisA = state.tutorielsList.filter(
        (v) => v.categorie_permis === "A" || v.categorie_permis === "A1",
      ).length;
      const poidsLourd = state.tutorielsList.filter((v) =>
        ["C", "D", "E"].includes(v.categorie_permis),
      ).length;
      const tousPermis = state.tutorielsList.filter(
        (v) => v.categorie_permis === "TOUTES",
      ).length;
      const totalVues = state.tutorielsList.reduce(
        (acc, curr) => acc + (curr.vues_count || 0),
        0,
      );

      return {
        total,
        permisB,
        permisA,
        poidsLourd,
        tousPermis,
        totalVues,
      };
    },
  },

  actions: {
    async loadTutoriels() {
      this.isLoading = true;
      this.error = null;
      const authStore = useAuthStore();
      const schoolId = authStore.currentSchool?.id || "default";
      const cacheKey = `${STORAGE_PREFIX}${schoolId}`;

      try {
        const data = await api.getTutorielsVideos();
        this.tutorielsList = data || [];

        if (typeof localStorage !== "undefined") {
          localStorage.setItem(cacheKey, JSON.stringify(this.tutorielsList));
        }
      } catch (err: any) {
        this.error = err.message || "Impossible de joindre la base de données.";
        if (typeof localStorage !== "undefined") {
          const cached = localStorage.getItem(cacheKey);
          if (cached) {
            try {
              const parsed = JSON.parse(cached);
              if (Array.isArray(parsed)) {
                this.tutorielsList = parsed;
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

    async ajouterTutoriel(payload: {
      titre: string;
      categorie_permis: string;
      description?: string;
      url_video: string;
      duree_minutes?: number;
      miniature_url?: string;
    }) {
      this.isLoading = true;
      const authStore = useAuthStore();
      const guichetStore = useGuichetStore();

      // Détecter plateforme vidéo
      let plateforme: "YOUTUBE" | "VIMEO" | "MP4_DIRECT" | "AUTRE" = "YOUTUBE";
      const url = payload.url_video.toLowerCase();
      if (url.includes("youtube.com") || url.includes("youtu.be")) {
        plateforme = "YOUTUBE";
      } else if (url.includes("vimeo.com")) {
        plateforme = "VIMEO";
      } else if (url.endsWith(".mp4") || url.endsWith(".webm")) {
        plateforme = "MP4_DIRECT";
      } else {
        plateforme = "AUTRE";
      }

      const formateur = authStore.currentUser?.nom
        ? `${authStore.currentUser.prenom} ${authStore.currentUser.nom}`
        : "Direction Pédagogique";

      try {
        const nouveauTutoriel = await api.creerTutorielVideo({
          ...payload,
          plateforme,
          ajoute_par: formateur,
        });

        this.tutorielsList.unshift(nouveauTutoriel);

        const schoolId = authStore.currentSchool?.id || "default";
        if (typeof localStorage !== "undefined") {
          localStorage.setItem(
            `${STORAGE_PREFIX}${schoolId}`,
            JSON.stringify(this.tutorielsList),
          );
        }

        guichetStore.showToast(
          "Tutoriel vidéo ajouté avec succès !",
          "success",
        );
        this.isUploadModalOpen = false;
        return nouveauTutoriel;
      } catch (err: any) {
        guichetStore.showToast(
          `Erreur ajout tutoriel : ${err.message}`,
          "error",
        );
        throw err;
      } finally {
        this.isLoading = false;
      }
    },

    async supprimerTutoriel(id: string) {
      const guichetStore = useGuichetStore();
      try {
        await api.supprimerTutorielVideo(id);
        this.tutorielsList = this.tutorielsList.filter((t) => t.id !== id);

        const authStore = useAuthStore();
        const schoolId = authStore.currentSchool?.id || "default";
        if (typeof localStorage !== "undefined") {
          localStorage.setItem(
            `${STORAGE_PREFIX}${schoolId}`,
            JSON.stringify(this.tutorielsList),
          );
        }

        guichetStore.showToast("Tutoriel vidéo supprimé avec succès.", "info");
      } catch (err: any) {
        guichetStore.showToast(`Erreur suppression : ${err.message}`, "error");
        throw err;
      }
    },

    async playVideo(video: TutorielVideo) {
      this.activeVideoToPlay = video;
      // Incrémenter vue
      try {
        await api.incrementerVueTutoriel(video.id);
        if (typeof video.vues_count === "number") {
          video.vues_count++;
        }
      } catch {
        // silence
      }
    },

    closePlayer() {
      this.activeVideoToPlay = null;
    },

    getTutorielsForStudent(studentCategory: string): TutorielVideo[] {
      const s = (studentCategory || "B").toUpperCase().trim();
      const targetCat =
        s === "TOUTES" || s === "ALL" || s.includes("TRONC")
          ? "TOUTES"
          : s.includes("A") || s.includes("MOTO")
            ? "A"
            : s.includes("C")
              ? "C"
              : s.includes("D")
                ? "D"
                : s.includes("E")
                  ? "E"
                  : "B";
      if (targetCat === "TOUTES") return this.tutorielsList;
      const matched = this.tutorielsList.filter((t) => {
        const tc = (t.categorie_permis || "TOUTES").toUpperCase().trim();
        return tc === targetCat || tc === "TOUTES" || tc.includes(targetCat);
      });
      return matched.length > 0 ? matched : this.tutorielsList;
    },
  },
});

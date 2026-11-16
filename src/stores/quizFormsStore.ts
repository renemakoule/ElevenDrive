import { defineStore } from "pinia";
import { api } from "../db/api";
import { useGuichetStore } from "./guichetStore";
import { realtimeSocket } from "../services/realtimeSocketService";
import type {
  QuizFormulaire,
  QuizQuestionItem,
  QuizSoumissionEleve,
  QuizReponseDetailItem,
} from "../types";

const STORAGE_FORMS_KEY = "autoschool_quiz_forms_v1";
const STORAGE_SUBS_KEY = "autoschool_quiz_submissions_v1";

function safeWriteJson(key: string, data: any): void {
  if (typeof localStorage === "undefined") return;
  try {
    localStorage.setItem(key, JSON.stringify(data));
  } catch {
    // Ignore quota error
  }
}

export const useQuizFormsStore = defineStore("quizForms", {
  state: () => ({
    formulaires: [] as QuizFormulaire[],
    soumissions: [] as QuizSoumissionEleve[],
    loading: false,
    submitting: false,
    selectedCategoryFilter: "ALL" as string,
    selectedQuizFilterId: "ALL" as string,
    lastReceivedSubmission: null as QuizSoumissionEleve | null,
  }),

  getters: {
    filteredFormulaires: (state): QuizFormulaire[] => {
      if (state.selectedCategoryFilter === "ALL") {
        return state.formulaires;
      }
      const norm = state.selectedCategoryFilter.toUpperCase();
      return state.formulaires.filter(
        (q) =>
          q.categorie_permis.toUpperCase() === norm ||
          q.categorie_permis.toUpperCase() === "TOUTES",
      );
    },

    getQuizForStudent: (state) => {
      return (studentCat: string): QuizFormulaire[] => {
        const raw = (studentCat || "B").toUpperCase().trim();
        const cleaned = raw.replace(/^PERMIS\s+/i, "").trim();
        const letterMatch = cleaned.match(/^([A-E]1?)\b/);
        const normalizedCat = letterMatch ? letterMatch[1] : cleaned;

        return state.formulaires.filter((q) => {
          if (!q.est_publie) return false;
          const qCat = (q.categorie_permis || "TOUTES").toUpperCase().trim();
          return (
            qCat === normalizedCat ||
            qCat === raw ||
            qCat === "TOUTES" ||
            qCat === "ALL" ||
            normalizedCat.startsWith(qCat)
          );
        });
      };
    },

    getSoumissionsByQuiz: (state) => {
      return (quizId: string): QuizSoumissionEleve[] => {
        if (!quizId || quizId === "ALL") {
          return state.soumissions;
        }
        return state.soumissions.filter((s) => s.quiz_id === quizId);
      };
    },

    getSoumissionsByCandidat: (state) => {
      return (candidatId: string): QuizSoumissionEleve[] => {
        if (!candidatId) return [];
        return state.soumissions.filter((s) => s.candidat_id === candidatId);
      };
    },
  },

  actions: {
    async loadAll() {
      await Promise.all([this.loadFormulaires(), this.loadSoumissions()]);
    },

    async loadFormulaires(categorie = "ALL") {
      this.loading = true;
      try {
        const list = await api.getQuizFormulaires(categorie);
        if (Array.isArray(list)) {
          this.formulaires = list;
          safeWriteJson(STORAGE_FORMS_KEY, list);
          return;
        }
      } catch {
        // Fallback localStorage si hors-ligne
      } finally {
        this.loading = false;
      }

      if (typeof localStorage !== "undefined") {
        try {
          const raw = localStorage.getItem(STORAGE_FORMS_KEY);
          if (raw) {
            const parsed = JSON.parse(raw);
            if (Array.isArray(parsed)) {
              this.formulaires = parsed;
            }
          }
        } catch {
          this.formulaires = [];
        }
      }
    },

    async loadSoumissions(quizId = "", candidatId = "") {
      try {
        const list = await api.getQuizSoumissions(quizId, candidatId);
        if (Array.isArray(list)) {
          this.soumissions = list;
          safeWriteJson(STORAGE_SUBS_KEY, list);
          return;
        }
      } catch {
        // Fallback localStorage si hors-ligne
      }

      if (typeof localStorage !== "undefined") {
        try {
          const raw = localStorage.getItem(STORAGE_SUBS_KEY);
          if (raw) {
            const parsed = JSON.parse(raw);
            if (Array.isArray(parsed)) {
              this.soumissions = parsed;
            }
          }
        } catch {
          this.soumissions = [];
        }
      }
    },

    async creerFormulaire(payload: {
      moniteur_id?: string;
      moniteur_nom: string;
      titre: string;
      description?: string;
      categorie_permis: string;
      duree_minutes: number;
      seuil_reussite_pct: number;
      est_publie: boolean;
      questions: QuizQuestionItem[];
    }): Promise<QuizFormulaire> {
      const guichetStore = useGuichetStore();
      this.submitting = true;
      try {
        let created: QuizFormulaire;
        try {
          created = await api.creerQuizFormulaire(payload);
        } catch {
          created = {
            id: "quiz-" + Date.now(),
            moniteur_id: payload.moniteur_id,
            moniteur_nom: payload.moniteur_nom,
            titre: payload.titre,
            description: payload.description,
            categorie_permis: payload.categorie_permis,
            duree_minutes: payload.duree_minutes,
            seuil_reussite_pct: payload.seuil_reussite_pct,
            est_publie: payload.est_publie,
            questions: payload.questions,
            reponses_count: 0,
            created_at: new Date().toISOString(),
          };
          realtimeSocket.sendEvent("QUIZ_FORM_UPDATED", created);
        }

        this.formulaires = [
          created,
          ...this.formulaires.filter((f) => f.id !== created.id),
        ];
        safeWriteJson(STORAGE_FORMS_KEY, this.formulaires);
        guichetStore.showToast(
          `Quiz « ${created.titre} » publié et envoyé aux élèves en temps réel !`,
          "success",
        );
        return created;
      } finally {
        this.submitting = false;
      }
    },

    async supprimerFormulaire(id: string) {
      const guichetStore = useGuichetStore();
      try {
        await api.supprimerQuizFormulaire(id);
      } catch {
        realtimeSocket.sendEvent("QUIZ_FORM_UPDATED", {
          action: "DELETED",
          id,
        });
      }
      this.formulaires = this.formulaires.filter((f) => f.id !== id);
      this.soumissions = this.soumissions.filter((s) => s.quiz_id !== id);
      safeWriteJson(STORAGE_FORMS_KEY, this.formulaires);
      safeWriteJson(STORAGE_SUBS_KEY, this.soumissions);
      guichetStore.showToast("Formulaire de quiz supprimé.", "info");
    },

    async soumettreReponseEleve(payload: {
      quiz_id: string;
      quiz_titre: string;
      candidat_id: string;
      candidat_nom: string;
      numero_dossier: string;
      categorie_permis: string;
      score_obtenu: number;
      total_points: number;
      pourcentage: number;
      est_reussi: boolean;
      duree_secondes: number;
      details: QuizReponseDetailItem[];
    }): Promise<QuizSoumissionEleve> {
      const guichetStore = useGuichetStore();
      this.submitting = true;
      try {
        let created: QuizSoumissionEleve;
        try {
          created = await api.soumettreQuizReponse(payload);
        } catch {
          created = {
            id: "sub-" + Date.now(),
            ...payload,
            submitted_at: new Date().toISOString(),
          };
          realtimeSocket.sendEvent("QUIZ_SUBMISSION_RECEIVED", created);
        }

        this.soumissions = [
          created,
          ...this.soumissions.filter((s) => s.id !== created.id),
        ];
        const quizTarget = this.formulaires.find(
          (q) => q.id === payload.quiz_id,
        );
        if (quizTarget) {
          quizTarget.reponses_count = (quizTarget.reponses_count || 0) + 1;
        }
        safeWriteJson(STORAGE_SUBS_KEY, this.soumissions);
        safeWriteJson(STORAGE_FORMS_KEY, this.formulaires);

        guichetStore.showToast(
          `Réponses envoyées au moniteur ! Score : ${created.score_obtenu}/${created.total_points} (${created.pourcentage}%)`,
          created.est_reussi ? "success" : "info",
        );
        return created;
      } finally {
        this.submitting = false;
      }
    },

    registerIncomingSubmission(sub: QuizSoumissionEleve) {
      if (!sub || !sub.id) return;
      this.lastReceivedSubmission = sub;
      const exists = this.soumissions.some((s) => s.id === sub.id);
      if (!exists) {
        this.soumissions = [sub, ...this.soumissions];
        const q = this.formulaires.find((f) => f.id === sub.quiz_id);
        if (q) {
          q.reponses_count = (q.reponses_count || 0) + 1;
        }
        safeWriteJson(STORAGE_SUBS_KEY, this.soumissions);
      }
    },
  },
});

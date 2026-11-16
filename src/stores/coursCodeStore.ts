import { defineStore } from "pinia";
import { api } from "../db/api";
import type { SeanceCode, TypeSeanceCode } from "../types";
import { useGuichetStore } from "./guichetStore";
import { JITSI_BASE_URL } from "../config/env";

export function isSeanceCodeExpired(
  s:
    | Pick<SeanceCode, "date_seance" | "heure_fin" | "statut_seance">
    | null
    | undefined,
  nowMs: number = Date.now(),
): boolean {
  if (!s) return true;
  if (s.statut_seance === "TERMINEE" || s.statut_seance === "ANNULEE") {
    return true;
  }
  const datePart = (s.date_seance || "").trim().slice(0, 10);
  if (!/^\d{4}-\d{2}-\d{2}$/.test(datePart)) return false;
  const endTime = (s.heure_fin || "").trim();
  const timePart = /^\d{2}:\d{2}/.test(endTime) ? endTime.slice(0, 5) : "23:59";
  const endAt = new Date(`${datePart}T${timePart}:00`).getTime();
  if (Number.isNaN(endAt)) return false;
  return nowMs >= endAt;
}

export const useCoursCodeStore = defineStore("coursCode", {
  state: () => ({
    seances: [] as SeanceCode[],
    loading: false,
    filtreType: "ALL" as "ALL" | "PRESENTIEL" | "EN_LIGNE",
    seanceActivePourVisio: null as SeanceCode | null,
    isPlanificationModalOpen: false,
  }),

  getters: {
    seancesFiltrees: (state) => {
      if (state.filtreType === "ALL") return state.seances;
      return state.seances.filter((s) => s.type_seance === state.filtreType);
    },

    seancesEnLigne: (state) => {
      return state.seances.filter((s) => s.type_seance === "EN_LIGNE");
    },

    seancesAujourdhui: (state) => {
      const today = new Date().toISOString().slice(0, 10);
      return state.seances.filter((s) => s.date_seance === today);
    },

    prochaineSeanceEnLigne: (state) => {
      const nowMs = Date.now();
      return (
        state.seances
          .filter(
            (s) =>
              s.type_seance === "EN_LIGNE" && !isSeanceCodeExpired(s, nowMs),
          )
          .sort(
            (a, b) =>
              a.date_seance.localeCompare(b.date_seance) ||
              a.heure_debut.localeCompare(b.heure_debut),
          )[0] || null
      );
    },
  },

  actions: {
    async loadSeances(filtre?: string) {
      this.loading = true;
      try {
        const type =
          filtre || (this.filtreType !== "ALL" ? this.filtreType : undefined);
        const data = await api.getSeancesCode(type);
        this.seances = data || [];
      } catch (err: any) {
        console.warn("Erreur chargement séances code:", err);
      } finally {
        this.loading = false;
      }
    },

    async planifierSeance(payload: {
      theme_traite: string;
      formateur_nom: string;
      moniteur_id?: string;
      date_seance: string;
      heure_debut: string;
      heure_fin: string;
      type_seance: TypeSeanceCode;
      lien_visio?: string;
      mot_de_passe_visio?: string;
    }) {
      this.loading = true;
      const guichetStore = useGuichetStore();
      try {
        // Si visio et aucun lien renseigné, générer une salle Jitsi Meet sécurisée
        let lien = payload.lien_visio;
        if (payload.type_seance === "EN_LIGNE" && (!lien || !lien.trim())) {
          const roomCode =
            "AutoSchool-" +
            Math.random().toString(36).substring(2, 8).toUpperCase();
        }

        const nouvelleSeance = await api.creerSeanceCode({
          ...payload,
          lien_visio: lien,
        });

        this.seances.unshift(nouvelleSeance);
        guichetStore.showToast(
          payload.type_seance === "EN_LIGNE"
            ? "Cours en ligne planifié avec succès ! Salon de visio activé."
            : "Séance en salle planifiée avec succès !",
          "success",
        );
        this.isPlanificationModalOpen = false;
        return nouvelleSeance;
      } catch (err: any) {
        guichetStore.showToast(
          "Erreur planification séance : " + err.message,
          "error",
        );
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async reprogrammerSeance(payload: {
      id: string;
      date_seance: string;
      heure_debut: string;
      heure_fin: string;
    }) {
      this.loading = true;
      const guichetStore = useGuichetStore();
      try {
        const updated = await api.reprogrammerSeanceCode(payload);
        const idx = this.seances.findIndex((s) => s.id === payload.id);
        if (idx !== -1) {
          this.seances[idx] = updated;
        }
        guichetStore.showToast(
          "Cours reprogrammé et réactivé avec succès !",
          "success",
        );
        return updated;
      } catch (err: any) {
        guichetStore.showToast(
          "Erreur reprogrammation : " + err.message,
          "error",
        );
        throw err;
      } finally {
        this.loading = false;
      }
    },

    rejoindreVisio(seance: SeanceCode) {
      if (!seance.lien_visio) {
        useGuichetStore().showToast(
          "Aucun lien de visioconférence associé à cette séance.",
          "error",
        );
        return;
      }
      this.seanceActivePourVisio = seance;
    },

    fermerVisio() {
      this.seanceActivePourVisio = null;
    },
  },
});

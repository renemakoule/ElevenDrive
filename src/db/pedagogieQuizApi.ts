import type {
  QuizFormulaire,
  QuizQuestionItem,
  QuizReponseDetailItem,
  QuizSoumissionEleve,
} from "../types";
import { realtimeSocket } from "../services/realtimeSocketService";
import { apiFetch } from "./httpClient";

export const pedagogieQuizApi = {
  async enregistrerQuizScore(payload: {
    candidat_id: string;
    score: number;
  }): Promise<any> {
    const res = await apiFetch<any>("/pedagogie/quiz-resultat", {
      method: "POST",
      body: JSON.stringify(payload),
    });
    realtimeSocket.sendEvent("QUIZ_RESULTAT_UPDATED", payload);
    return res;
  },

  async getQuizFormulaires(categorie = "ALL"): Promise<QuizFormulaire[]> {
    return apiFetch<QuizFormulaire[]>(
      `/pedagogie/quiz-forms?categorie=${encodeURIComponent(categorie)}`,
    );
  },

  async creerQuizFormulaire(payload: {
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
    const created = await apiFetch<QuizFormulaire>("/pedagogie/quiz-forms", {
      method: "POST",
      body: JSON.stringify(payload),
    });
    realtimeSocket.sendEvent("QUIZ_FORM_UPDATED", {
      action: "CREATED",
      id: created.id,
      titre: created.titre,
      categorie_permis: created.categorie_permis,
    });
    return created;
  },

  async supprimerQuizFormulaire(id: string): Promise<any> {
    const res = await apiFetch<any>(
      `/pedagogie/quiz-forms?id=${encodeURIComponent(id)}`,
      { method: "DELETE" },
    );
    realtimeSocket.sendEvent("QUIZ_FORM_UPDATED", { action: "DELETED", id });
    return res;
  },

  async getQuizSoumissions(
    quizId = "",
    candidatId = "",
  ): Promise<QuizSoumissionEleve[]> {
    const params = new URLSearchParams();
    if (quizId) params.set("quiz_id", quizId);
    if (candidatId) params.set("candidat_id", candidatId);
    const qs = params.toString();
    return apiFetch<QuizSoumissionEleve[]>(
      `/pedagogie/quiz-submissions${qs ? `?${qs}` : ""}`,
    );
  },

  async soumettreQuizReponse(payload: {
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
    const created = await apiFetch<QuizSoumissionEleve>(
      "/pedagogie/quiz-submissions",
      {
        method: "POST",
        body: JSON.stringify(payload),
      },
    );
    realtimeSocket.sendEvent("QUIZ_SUBMISSION_RECEIVED", created);
    realtimeSocket.sendEvent("QUIZ_RESULTAT_UPDATED", {
      candidat_id: created.candidat_id,
      score: created.pourcentage,
    });
    return created;
  },
};


import type {
  SessionExamenActiveResponse,
  SessionExamen,
  CreerOuModifierSessionExamenInput,
} from "../types";
import { realtimeSocket } from "../services/realtimeSocketService";
import { isTauri, apiFetch, invokeTauri, mirrorToTauri } from "./httpClient";

export const examensApi = {
  async getSessionExamenActive(): Promise<SessionExamenActiveResponse> {
    try {
      return await apiFetch<SessionExamenActiveResponse>(
        "/examens/session-active",
      );
    } catch {
      return {
        session: null,
        candidats: [],
        ineligibles: [],
      };
    }
  },

  async enregistrerSessionExamen(
    payload: CreerOuModifierSessionExamenInput,
  ): Promise<SessionExamen> {
    const res = await apiFetch<SessionExamen>("/examens/sessions", {
      method: "POST",
      body: JSON.stringify(payload),
    });
    realtimeSocket.sendEvent("EXAMEN_SESSION_UPDATED", res);
    return res;
  },

  async supprimerSessionExamen(
    sessionId: string,
  ): Promise<{ message: string }> {
    const res = await apiFetch<{ message: string }>(
      "/examens/sessions/supprimer",
      {
        method: "POST",
        body: JSON.stringify({ session_id: sessionId }),
      },
    );
    realtimeSocket.sendEvent("EXAMEN_SESSION_UPDATED", {
      deleted: true,
      session_id: sessionId,
    });
    return res;
  },

  async enregistrerResultatsExamen(payload: {
    session_examen_id: string;
    resultats: {
      candidat_examen_id: string;
      resultat_code: string;
      resultat_conduite: string;
      resultat_final: string;
    }[];
  }): Promise<any> {
    const res = await apiFetch<any>("/examens/resultats", {
      method: "POST",
      body: JSON.stringify(payload),
    });
    realtimeSocket.sendEvent("EXAMEN_SESSION_UPDATED", res);
    return res;
  },

  async forcerCandidatExamen(payload: {
    session_examen_id: string;
    dossier_id: string;
    motif?: string;
    code_pin?: string;
  }): Promise<any> {
    try {
      const res = await apiFetch<any>("/examens/forcer-candidat", {
        method: "POST",
        body: JSON.stringify(payload),
      });
      mirrorToTauri("forcer_candidat_examen", { payload });
      realtimeSocket.sendEvent("EXAMENS_SESSION_UPDATED", res);
      return res;
    } catch (err) {
      if (isTauri) {
        const local = await invokeTauri("forcer_candidat_examen", {
          payload,
        });
        realtimeSocket.sendEvent("EXAMENS_SESSION_UPDATED", local);
        return local;
      }
      throw err;
    }
  },
};


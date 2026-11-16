import type {
  Vehicule,
  CreerVehiculeInput,
  SeanceConduite,
  CompetenceItem,
} from "../types";
import { realtimeSocket } from "../services/realtimeSocketService";
import { isTauri, apiFetch, invokeTauri, mirrorToTauri } from "./httpClient";

export const flotteConduiteApi = {
  async getVehicules(): Promise<Vehicule[]> {
    try {
      return await apiFetch<Vehicule[]>("/flotte/vehicules");
    } catch (err) {
      if (isTauri) return invokeTauri<Vehicule[]>("get_vehicules");
      throw err;
    }
  },

  async creerVehicule(input: CreerVehiculeInput): Promise<Vehicule> {
    try {
      const created = await apiFetch<Vehicule>("/flotte/vehicules", {
        method: "POST",
        body: JSON.stringify(input),
      });
      mirrorToTauri("creer_vehicule", { payload: input });
      realtimeSocket.sendEvent("FLOTTE_VEHICULES_UPDATED", created);
      return created;
    } catch (err) {
      if (isTauri) {
        const local = await invokeTauri<Vehicule>("creer_vehicule", {
          payload: input,
        });
        realtimeSocket.sendEvent("FLOTTE_VEHICULES_UPDATED", local);
        return local;
      }
      throw err;
    }
  },

  async updateVehiculeStatut(vehiculeId: string, statut: string): Promise<any> {
    try {
      const res = await apiFetch<any>("/flotte/vehicules/statut", {
        method: "PATCH",
        body: JSON.stringify({ vehicule_id: vehiculeId, statut }),
      });
      mirrorToTauri("update_vehicule_statut", { vehiculeId, statut });
      realtimeSocket.sendEvent("FLOTTE_VEHICULES_UPDATED", {
        vehicule_id: vehiculeId,
        statut,
      });
      return res;
    } catch (err) {
      if (isTauri) {
        const local = await invokeTauri("update_vehicule_statut", {
          vehiculeId,
          statut,
        });
        realtimeSocket.sendEvent("FLOTTE_VEHICULES_UPDATED", {
          vehicule_id: vehiculeId,
          statut,
        });
        return local;
      }
      throw err;
    }
  },

  async getPlanningConduite(
    date?: string,
    moniteurId?: string,
  ): Promise<SeanceConduite[]> {
    const params = new URLSearchParams();
    if (date) params.append("date", date);
    if (moniteurId) params.append("moniteur_id", moniteurId);
    const qs = params.toString() ? `?${params.toString()}` : "";
    return apiFetch<SeanceConduite[]>(`/pedagogie/planning-conduite${qs}`);
  },

  async annulerSeanceConduite(id: string): Promise<any> {
    const res = await apiFetch<any>(
      `/pedagogie/planning-conduite?id=${encodeURIComponent(id)}`,
      { method: "DELETE" },
    );
    realtimeSocket.sendEvent("PLANNING_CONDUITE_UPDATED", {
      action: "DELETED",
      seance_id: id,
    });
    return res;
  },

  async validerSeanceConduite(payload: {
    sessionId: string;
    dossierId: string;
    vehiculeId?: string;
    dureeHeures: number;
    kmDebut: number;
    kmFin: number;
    commentaire: string;
    competences: CompetenceItem[];
    signatureEleve: string;
  }): Promise<any> {
    const res = await apiFetch<any>("/pedagogie/seance-conduite/valider", {
      method: "POST",
      body: JSON.stringify({
        seance_id: payload.sessionId,
        dossier_id: payload.dossierId,
        vehicule_id: payload.vehiculeId,
        duree_heures: payload.dureeHeures,
        km_debut: payload.kmDebut,
        km_fin: payload.kmFin,
        commentaire: payload.commentaire,
        competences: payload.competences,
        signature_eleve: payload.signatureEleve,
      }),
    });
    realtimeSocket.sendEvent("LIVRET_EMARGEMENT_UPDATED", res);
    realtimeSocket.sendEvent("PLANNING_CONDUITE_UPDATED", res);
    return res;
  },

  async reserverSeanceConduite(payload: {
    dossier_id: string;
    date_seance: string;
    heure_debut: string;
    heure_fin: string;
    vehicule_id?: string;
    moniteur_id?: string;
  }): Promise<SeanceConduite> {
    const res = await apiFetch<SeanceConduite>("/pedagogie/reserver-conduite", {
      method: "POST",
      body: JSON.stringify(payload),
    });
    realtimeSocket.sendEvent("PLANNING_CONDUITE_UPDATED", res);
    return res;
  },

  async reserverConduite(payload: {
    dossier_id: string;
    date_seance: string;
    heure_debut: string;
    heure_fin: string;
    vehicule_id?: string;
    moniteur_id?: string;
  }): Promise<SeanceConduite> {
    return flotteConduiteApi.reserverSeanceConduite(payload);
  },

  async getMoniteurs(): Promise<
    { id: string; nomComplet: string; role: string; telephone?: string }[]
  > {
    try {
      const { authBackend } = await import("../services/authBackendService");
      const users = await authBackend.getSchoolUsers();
      const filtered = users
        .filter((u) => u.role === "MONITEUR" || u.role === "DIRECTEUR")
        .sort((a, b) => {
          if (a.role === "MONITEUR" && b.role !== "MONITEUR") return -1;
          if (a.role !== "MONITEUR" && b.role === "MONITEUR") return 1;
          return 0;
        })
        .map((u) => ({
          id: u.id,
          nomComplet: `${u.prenoms || ""} ${u.nom || ""}`.trim() || u.nom,
          role: u.role,
          telephone: u.telephone,
        }));
      if (filtered.length > 0) return filtered;
    } catch {
      // Fallback vers SQLite local si hors-ligne
    }

    if (isTauri) {
      try {
        const list = await invokeTauri<any[]>("get_moniteurs");
        if (Array.isArray(list) && list.length > 0) {
          return list.map((m) => ({
            id: m.id,
            nomComplet:
              `${m.prenoms || ""} ${m.nom || ""}`.trim() || m.nom || "Moniteur",
            role: m.role || "MONITEUR",
            telephone: m.telephone,
          }));
        }
      } catch (err) {
        console.warn("Fallback get_moniteurs:", err);
      }
    }

    return [];
  },
};

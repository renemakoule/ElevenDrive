import type {
  EpreuvePdf,
  SeanceCode,
  TutorielVideo,
  SupportCours,
} from "../types";
import { realtimeSocket } from "../services/realtimeSocketService";
import { isTauri, apiFetch, invokeTauri, mirrorToTauri } from "./httpClient";
import { pedagogieQuizApi } from "./pedagogieQuizApi";

export const pedagogieContentApi = {
  ...pedagogieQuizApi,
  async enregistrerSeanceCode(payload: {
    theme_traite: string;
    formateur_nom: string;
    heure_debut: string;
    heure_fin: string;
    presences: Array<{
      candidat_id: string;
      est_present: boolean;
      score_test: number;
    }>;
  }): Promise<any> {
    try {
      const res = await apiFetch<any>("/pedagogie/seance-code", {
        method: "POST",
        body: JSON.stringify(payload),
      });
      mirrorToTauri("enregistrer_seance_code_avec_presences", {
        themeTraite: payload.theme_traite,
        formateurNom: payload.formateur_nom,
        heureDebut: payload.heure_debut,
        heureFin: payload.heure_fin,
        presences: payload.presences,
      });
      realtimeSocket.sendEvent("COURS_CODE_UPDATED", res);
      realtimeSocket.sendEvent("CANDIDAT_DOSSIER_UPDATED", res);
      return res;
    } catch (err) {
      if (isTauri) {
        const local = await invokeTauri(
          "enregistrer_seance_code_avec_presences",
          {
            themeTraite: payload.theme_traite,
            formateurNom: payload.formateur_nom,
            heureDebut: payload.heure_debut,
            heureFin: payload.heure_fin,
            presences: payload.presences,
          },
        );
        realtimeSocket.sendEvent("COURS_CODE_UPDATED", local);
        realtimeSocket.sendEvent("CANDIDAT_DOSSIER_UPDATED", local);
        return local;
      }
      throw err;
    }
  },

  async getEpreuves(categorie?: string): Promise<EpreuvePdf[]> {
    const q =
      categorie && categorie !== "ALL"
        ? `?categorie=${encodeURIComponent(categorie)}`
        : "";
    try {
      return await apiFetch<EpreuvePdf[]>(`/epreuves${q}`);
    } catch (err) {
      if (isTauri) {
        return invokeTauri<EpreuvePdf[]>("get_epreuves", {
          categorie: categorie || null,
        });
      }
      throw err;
    }
  },

  async creerEpreuve(payload: {
    titre: string;
    categorie_permis: string;
    description?: string;
    annee: number;
    session?: string;
    duree_minutes: number;
    nombre_questions: number;
    avec_corrige: boolean;
    fichier_nom: string;
    fichier_taille?: string;
    fichier_data: string;
    ajoute_par?: string;
  }): Promise<EpreuvePdf> {
    try {
      const created = await apiFetch<EpreuvePdf>("/epreuves", {
        method: "POST",
        body: JSON.stringify(payload),
      });
      mirrorToTauri("creer_epreuve", { payload });
      realtimeSocket.sendEvent("EPREUVES_UPDATED", {
        action: "CREATED",
        id: created.id,
        titre: created.titre,
      });
      return created;
    } catch (err) {
      if (isTauri) {
        const local = await invokeTauri<EpreuvePdf>("creer_epreuve", {
          payload,
        });
        realtimeSocket.sendEvent("EPREUVES_UPDATED", {
          action: "CREATED",
          id: local.id,
          titre: local.titre,
        });
        return local;
      }
      throw err;
    }
  },

  async supprimerEpreuve(id: string): Promise<any> {
    try {
      const res = await apiFetch<any>(
        `/epreuves?id=${encodeURIComponent(id)}`,
        { method: "DELETE" },
      );
      mirrorToTauri("supprimer_epreuve", { id });
      realtimeSocket.sendEvent("EPREUVES_UPDATED", { action: "DELETED", id });
      return res;
    } catch (err) {
      if (isTauri) {
        const local = await invokeTauri<boolean>("supprimer_epreuve", { id });
        realtimeSocket.sendEvent("EPREUVES_UPDATED", { action: "DELETED", id });
        return local;
      }
      throw err;
    }
  },

  async incrementerTelechargementEpreuve(id: string): Promise<any> {
    try {
      const res = await apiFetch<any>(
        `/epreuves/telecharger?id=${encodeURIComponent(id)}`,
        { method: "POST" },
      );
      mirrorToTauri("incrementer_telechargement_epreuve", { id });
      realtimeSocket.sendEvent("EPREUVES_UPDATED", {
        action: "DOWNLOADED",
        id,
      });
      return res;
    } catch (err) {
      if (isTauri) {
        return invokeTauri<boolean>("incrementer_telechargement_epreuve", {
          id,
        });
      }
      throw err;
    }
  },

  async getSeancesCode(filtreType?: string): Promise<SeanceCode[]> {
    const q =
      filtreType && filtreType !== "ALL"
        ? `?type=${encodeURIComponent(filtreType)}`
        : "";
    try {
      return await apiFetch<SeanceCode[]>(`/pedagogie/seances-code${q}`);
    } catch (err) {
      if (isTauri) return invokeTauri<SeanceCode[]>("get_seances_code", {
        filtreType: filtreType || null,
      });
      throw err;
    }
  },

  async creerSeanceCode(payload: {
    theme_traite: string;
    formateur_nom: string;
    moniteur_id?: string;
    date_seance: string;
    heure_debut: string;
    heure_fin: string;
    type_seance: "PRESENTIEL" | "EN_LIGNE";
    lien_visio?: string;
    mot_de_passe_visio?: string;
  }): Promise<SeanceCode> {
    try {
      const created = await apiFetch<SeanceCode>("/pedagogie/seances-code", {
        method: "POST",
        body: JSON.stringify(payload),
      });
      mirrorToTauri("creer_seance_code", { payload });
      realtimeSocket.sendEvent("COURS_CODE_UPDATED", created);
      return created;
    } catch (err) {
      if (isTauri) {
        const local = await invokeTauri<SeanceCode>("creer_seance_code", {
          payload,
        });
        realtimeSocket.sendEvent("COURS_CODE_UPDATED", local);
        return local;
      }
      throw err;
    }
  },

  async reprogrammerSeanceCode(payload: {
    id: string;
    date_seance: string;
    heure_debut: string;
    heure_fin: string;
  }): Promise<SeanceCode> {
    try {
      const updated = await apiFetch<SeanceCode>("/pedagogie/seances-code", {
        method: "PUT",
        body: JSON.stringify(payload),
      });
      mirrorToTauri("reprogrammer_seance_code", {
        id: payload.id,
        dateSeance: payload.date_seance,
        heureDebut: payload.heure_debut,
        heureFin: payload.heure_fin,
      });
      realtimeSocket.sendEvent("COURS_CODE_UPDATED", updated);
      return updated;
    } catch (err) {
      if (isTauri) {
        const local = await invokeTauri<SeanceCode>(
          "reprogrammer_seance_code",
          {
            id: payload.id,
            dateSeance: payload.date_seance,
            heureDebut: payload.heure_debut,
            heureFin: payload.heure_fin,
          },
        );
        realtimeSocket.sendEvent("COURS_CODE_UPDATED", local);
        return local;
      }
      throw err;
    }
  },

  async enregistrerSeanceCodeAvecPresences(payload: {
    theme_traite: string;
    formateur_nom: string;
    heure_debut: string;
    heure_fin: string;
    presences: {
      candidat_id: string;
      est_present: boolean;
      score_test: number;
    }[];
  }): Promise<any> {
    try {
      const res = await apiFetch<any>("/pedagogie/seance-code", {
        method: "POST",
        body: JSON.stringify(payload),
      });
      mirrorToTauri("enregistrer_seance_code_avec_presences", payload);
      realtimeSocket.sendEvent("COURS_CODE_UPDATED", res);
      realtimeSocket.sendEvent("CANDIDAT_DOSSIER_UPDATED", res);
      return res;
    } catch (err) {
      if (isTauri) {
        const local = await invokeTauri(
          "enregistrer_seance_code_avec_presences",
          payload,
        );
        realtimeSocket.sendEvent("COURS_CODE_UPDATED", local);
        realtimeSocket.sendEvent("CANDIDAT_DOSSIER_UPDATED", local);
        return local;
      }
      throw err;
    }
  },

  async getTutorielsVideos(categorie?: string): Promise<TutorielVideo[]> {
    const q =
      categorie && categorie !== "ALL"
        ? `?categorie=${encodeURIComponent(categorie)}`
        : "";
    try {
      return await apiFetch<TutorielVideo[]>(`/pedagogie/tutoriels${q}`);
    } catch {
      if (isTauri) {
        try {
          return await invokeTauri<TutorielVideo[]>("get_tutoriels_videos", {
            categorie: categorie || null,
          });
        } catch {}
      }
      return [];
    }
  },

  async creerTutorielVideo(payload: {
    titre: string;
    categorie_permis: string;
    description?: string;
    url_video: string;
    plateforme?: "YOUTUBE" | "VIMEO" | "MP4_DIRECT" | "AUTRE";
    duree_minutes?: number;
    miniature_url?: string;
    ajoute_par?: string;
  }): Promise<TutorielVideo> {
    try {
      const created = await apiFetch<TutorielVideo>("/pedagogie/tutoriels", {
        method: "POST",
        body: JSON.stringify(payload),
      });
      mirrorToTauri("creer_tutoriel_video", { payload });
      realtimeSocket.sendEvent("TUTORIELS_UPDATED", created);
      return created;
    } catch (err) {
      if (isTauri) {
        const local = await invokeTauri<TutorielVideo>(
          "creer_tutoriel_video",
          { payload },
        );
        realtimeSocket.sendEvent("TUTORIELS_UPDATED", local);
        return local;
      }
      throw err;
    }
  },

  async supprimerTutorielVideo(id: string): Promise<any> {
    try {
      const res = await apiFetch<any>(
        `/pedagogie/tutoriels?id=${encodeURIComponent(id)}`,
        { method: "DELETE" },
      );
      mirrorToTauri("supprimer_tutoriel_video", { id });
      realtimeSocket.sendEvent("TUTORIELS_UPDATED", { action: "DELETED", id });
      return res;
    } catch (err) {
      if (isTauri) {
        const local = await invokeTauri<boolean>("supprimer_tutoriel_video", {
          id,
        });
        realtimeSocket.sendEvent("TUTORIELS_UPDATED", {
          action: "DELETED",
          id,
        });
        return local;
      }
      throw err;
    }
  },

  async incrementerVueTutoriel(id: string): Promise<any> {
    try {
      const res = await apiFetch<any>(
        `/pedagogie/tutoriels/vue?id=${encodeURIComponent(id)}`,
        { method: "POST" },
      );
      mirrorToTauri("incrementer_vue_tutoriel", { id });
      realtimeSocket.sendEvent("TUTORIELS_UPDATED", { action: "VIEWED", id });
      return res;
    } catch (err) {
      if (isTauri) return invokeTauri<boolean>("incrementer_vue_tutoriel", { id });
      throw err;
    }
  },

  async getSupportsCours(categorie?: string): Promise<SupportCours[]> {
    const q =
      categorie && categorie !== "ALL"
        ? `?categorie=${encodeURIComponent(categorie)}`
        : "";
    try {
      return await apiFetch<SupportCours[]>(`/pedagogie/supports-cours${q}`);
    } catch {
      if (isTauri) {
        try {
          return await invokeTauri<SupportCours[]>("get_supports_cours", {
            categorie: categorie || null,
          });
        } catch {}
      }
      return [];
    }
  },

  async creerSupportCours(payload: {
    titre: string;
    categorie_permis: string;
    description?: string;
    type_support: "PDF" | "LIEN_WEB";
    fichier_nom?: string;
    fichier_taille?: string;
    fichier_data?: string;
    url_externe?: string;
    ajoute_par?: string;
  }): Promise<SupportCours> {
    try {
      const created = await apiFetch<SupportCours>(
        "/pedagogie/supports-cours",
        {
          method: "POST",
          body: JSON.stringify(payload),
        },
      );
      mirrorToTauri("creer_support_cours", { payload });
      realtimeSocket.sendEvent("SUPPORTS_COURS_UPDATED", {
        action: "CREATED",
        id: created.id,
        titre: created.titre,
      });
      return created;
    } catch (err) {
      if (isTauri) {
        const local = await invokeTauri<SupportCours>("creer_support_cours", {
          payload,
        });
        realtimeSocket.sendEvent("SUPPORTS_COURS_UPDATED", {
          action: "CREATED",
          id: local.id,
          titre: local.titre,
        });
        return local;
      }
      throw err;
    }
  },

  async supprimerSupportCours(id: string): Promise<any> {
    try {
      const res = await apiFetch<any>(
        `/pedagogie/supports-cours?id=${encodeURIComponent(id)}`,
        { method: "DELETE" },
      );
      mirrorToTauri("supprimer_support_cours", { id });
      realtimeSocket.sendEvent("SUPPORTS_COURS_UPDATED", {
        action: "DELETED",
        id,
      });
      return res;
    } catch (err) {
      if (isTauri) {
        const local = await invokeTauri<boolean>("supprimer_support_cours", {
          id,
        });
        realtimeSocket.sendEvent("SUPPORTS_COURS_UPDATED", {
          action: "DELETED",
          id,
        });
        return local;
      }
      throw err;
    }
  },

  async incrementerTelechargementSupport(id: string): Promise<any> {
    try {
      const res = await apiFetch<any>(
        `/pedagogie/supports-cours/telecharger?id=${encodeURIComponent(id)}`,
        { method: "POST" },
      );
      mirrorToTauri("incrementer_telechargement_support", { id });
      realtimeSocket.sendEvent("SUPPORTS_COURS_UPDATED", {
        action: "DOWNLOADED",
        id,
      });
      return res;
    } catch (err) {
      if (isTauri) {
        return invokeTauri<boolean>("incrementer_telechargement_support", {
          id,
        });
      }
      throw err;
    }
  },
};

import type {
  Candidat,
  DossierDetailComplet,
  FormulePermis,
  TranchePaiementConfig,
  GuichetStats,
} from "../types";
import { realtimeSocket } from "../services/realtimeSocketService";
import {
  isTauri,
  API_BASE_URL,
  getAuthHeaders,
  apiFetch,
  invokeTauri,
  mirrorToTauri,
} from "./httpClient";

export const guichetApi = {
  async getGuichetStats(): Promise<GuichetStats> {
    try {
      return await apiFetch<GuichetStats>("/guichet/stats");
    } catch (err) {
      if (isTauri) return invokeTauri("get_guichet_stats");
      throw err;
    }
  },

  async searchCandidats(query: string): Promise<Candidat[]> {
    try {
      return await apiFetch<Candidat[]>(
        `/candidats?q=${encodeURIComponent(query)}`,
      );
    } catch (err) {
      if (isTauri) return invokeTauri("search_candidats", { query });
      throw err;
    }
  },

  async getCandidatDossier(
    candidatId: string,
  ): Promise<DossierDetailComplet | null> {
    try {
      return await apiFetch<DossierDetailComplet>(
        `/candidats/dossier?candidat_id=${encodeURIComponent(candidatId)}`,
      );
    } catch (err) {
      if (isTauri) return invokeTauri("get_candidat_dossier", { candidatId });
      throw err;
    }
  },

  async getFormules(): Promise<FormulePermis[]> {
    try {
      return await apiFetch<FormulePermis[]>("/formules");
    } catch (err) {
      if (isTauri) return invokeTauri("get_formules");
      throw err;
    }
  },

  async creerFormule(input: {
    categorie_permis: string;
    intitule: string;
    description?: string;
    heures_theoriques: number;
    heures_pratiques: number;
    prix_total: number;
    acompte_minimum: number;
    tranches_paiement?: TranchePaiementConfig[];
  }): Promise<FormulePermis> {
    try {
      const created = await apiFetch<FormulePermis>("/formules", {
        method: "POST",
        body: JSON.stringify(input),
      });
      mirrorToTauri("creer_formule", { payload: input });
      realtimeSocket.sendEvent("FORMULES_TARIFS_UPDATED", created);
      return created;
    } catch (err) {
      if (isTauri) {
        const local = await invokeTauri<FormulePermis>("creer_formule", {
          payload: input,
        });
        realtimeSocket.sendEvent("FORMULES_TARIFS_UPDATED", local);
        return local;
      }
      throw err;
    }
  },

  async supprimerFormule(formuleId: string): Promise<boolean> {
    try {
      await apiFetch(`/formules/${encodeURIComponent(formuleId)}`, {
        method: "DELETE",
      });
      mirrorToTauri("supprimer_formule", { formuleId });
      realtimeSocket.sendEvent("FORMULES_TARIFS_UPDATED", {
        deleted_id: formuleId,
      });
      return true;
    } catch (err) {
      if (isTauri) {
        const res = await invokeTauri<boolean>("supprimer_formule", {
          formuleId,
        });
        realtimeSocket.sendEvent("FORMULES_TARIFS_UPDATED", {
          deleted_id: formuleId,
        });
        return res;
      }
      throw err;
    }
  },

  async supprimerToutesFormules(): Promise<number> {
    try {
      const res = await apiFetch<{ count: number }>("/formules/clear-all", {
        method: "DELETE",
      });
      mirrorToTauri("supprimer_toutes_formules");
      realtimeSocket.sendEvent("FORMULES_TARIFS_UPDATED", { cleared: true });
      return res.count;
    } catch (err) {
      if (isTauri) {
        const count = await invokeTauri<number>("supprimer_toutes_formules");
        realtimeSocket.sendEvent("FORMULES_TARIFS_UPDATED", { cleared: true });
        return count;
      }
      throw err;
    }
  },

  async modifierPrixFormule(input: {
    formule_id: string;
    nouveau_prix: number;
    acompte_minimum?: number;
    motif?: string;
    tranches_paiement?: TranchePaiementConfig[];
  }): Promise<FormulePermis> {
    try {
      const updated = await apiFetch<FormulePermis>("/formules/prix", {
        method: "PATCH",
        body: JSON.stringify(input),
      });
      mirrorToTauri("modifier_prix_formule", { payload: input });
      realtimeSocket.sendEvent("FORMULES_TARIFS_UPDATED", updated);
      return updated;
    } catch (err) {
      if (isTauri) {
        const local = await invokeTauri<FormulePermis>(
          "modifier_prix_formule",
          { payload: input },
        );
        realtimeSocket.sendEvent("FORMULES_TARIFS_UPDATED", local);
        return local;
      }
      throw err;
    }
  },

  async creerCandidatAvecDossier(data: {
    nom: string;
    prenoms: string;
    genre: "M" | "F";
    date_naissance: string;
    lieu_naissance: string;
    telephone_whatsapp: string;
    type_piece_id: string;
    numero_piece_id: string;
    formule_id: string;
    acompte_initial: number;
    mode_reglement: string;
    photo_url?: string;
    visite_medicale_statut?:
      | "VALIDE_SANS_RESERVE"
      | "VALIDE_AVEC_VERRES"
      | "EN_ATTENTE";
    mode_formation_theorique?: "PRESENTIEL" | "EN_LIGNE" | "HYBRIDE";
  }): Promise<{ candidat: Candidat; dossier_detail: DossierDetailComplet }> {
    try {
      const created = await apiFetch<{
        candidat: Candidat;
        dossier_detail: DossierDetailComplet;
      }>("/candidats", {
        method: "POST",
        body: JSON.stringify(data),
      });
      mirrorToTauri("creer_candidat_avec_dossier", { payload: data });
      realtimeSocket.sendEvent("CANDIDAT_DOSSIER_UPDATED", created);
      return created;
    } catch (err) {
      if (isTauri) {
        const local = await invokeTauri<{
          candidat: Candidat;
          dossier_detail: DossierDetailComplet;
        }>("creer_candidat_avec_dossier", { payload: data });
        realtimeSocket.sendEvent("CANDIDAT_DOSSIER_UPDATED", local);
        return local;
      }
      throw err;
    }
  },

  async enregistrerPaiement(payload: {
    dossier_id: string;
    montant: number;
    mode_reglement: string;
    reference_transaction?: string;
  }): Promise<{
    paiement_id: string;
    numero_recu: string;
    code_securite_hash: string;
    montant_verse: number;
    nouveau_solde_restant: number;
    date_paiement: string;
    numero_dossier: string;
  }> {
    try {
      const res = await apiFetch<{
        paiement_id: string;
        numero_recu: string;
        code_securite_hash: string;
        montant_verse: number;
        nouveau_solde_restant: number;
        date_paiement: string;
        numero_dossier: string;
      }>("/paiements", {
        method: "POST",
        body: JSON.stringify(payload),
      });
      mirrorToTauri("enregistrer_paiement", {
        dossierId: payload.dossier_id,
        montant: payload.montant,
        modeReglement: payload.mode_reglement,
        referenceTransaction: payload.reference_transaction,
      });
      realtimeSocket.sendEvent("PAIEMENT_ENREGISTRE", res);
      return res;
    } catch (err) {
      if (isTauri) {
        const local = await invokeTauri<any>("enregistrer_paiement", {
          dossierId: payload.dossier_id,
          montant: payload.montant,
          modeReglement: payload.mode_reglement,
          referenceTransaction: payload.reference_transaction,
        });
        realtimeSocket.sendEvent("PAIEMENT_ENREGISTRE", local);
        return local;
      }
      throw err;
    }
  },

  async getSessionCaisseDetail(): Promise<any> {
    try {
      return await apiFetch<any>("/caisse/session-detail");
    } catch (err) {
      if (isTauri) return invokeTauri("get_session_caisse_detail");
      throw err;
    }
  },

  async cloturerCaisse(payload: {
    total_especes_compte: number;
    observations?: string;
  }): Promise<any> {
    try {
      const res = await apiFetch<any>("/caisse/cloturer", {
        method: "POST",
        body: JSON.stringify(payload),
      });
      mirrorToTauri("cloturer_session_caisse", { payload });
      realtimeSocket.sendEvent("CAISSE_SESSION_UPDATED", res);
      return res;
    } catch (err) {
      if (isTauri) {
        const local = await invokeTauri("cloturer_session_caisse", {
          payload,
        });
        realtimeSocket.sendEvent("CAISSE_SESSION_UPDATED", local);
        return local;
      }
      throw err;
    }
  },

  async reimprimerRecu(paiementId: string): Promise<any> {
    try {
      return await apiFetch<any>(
        `/paiements/recu?paiement_id=${encodeURIComponent(paiementId)}`,
      );
    } catch (err) {
      if (isTauri) return invokeTauri("reimprimer_recu", { paiementId });
      throw err;
    }
  },

  async verifierCodeManager(pin: string): Promise<boolean> {
    try {
      let schoolId = "";
      let sessionPin = "";
      if (typeof localStorage !== "undefined") {
        const raw =
          localStorage.getItem("autoschool_active_session") ||
          localStorage.getItem("elevendrive_active_session");
        if (raw) {
          const s = JSON.parse(raw);
          schoolId = s.auto_ecole?.id || "";
          sessionPin = s.user?.code_pin_manager || "";
        }
      }

      if (sessionPin && sessionPin === pin) return true;

      const response = await fetch(`${API_BASE_URL}/auth/verify-pin`, {
        method: "POST",
        headers: getAuthHeaders(),
        body: JSON.stringify({ auto_ecole_id: schoolId, code_pin: pin }),
      });
      const result = await response.json();
      if (response.ok && result.success) return true;
    } catch {
      // Secours hors-ligne
    }

    if (isTauri) {
      try {
        return await invokeTauri("verifier_code_manager_derogation", { pin });
      } catch {}
    }

    if (typeof localStorage !== "undefined") {
      const raw =
        localStorage.getItem("autoschool_active_session") ||
        localStorage.getItem("elevendrive_active_session");
      if (raw) {
        try {
          const s = JSON.parse(raw);
          return !!(
            s.user?.code_pin_manager && s.user.code_pin_manager === pin
          );
        } catch {}
      }
    }
    return false;
  },

  async getImpayesEtAlertesGerant(): Promise<any> {
    try {
      return await apiFetch<any>("/manager/dashboard");
    } catch (err) {
      if (isTauri) return invokeTauri("get_impayes_et_alertes_gerant");
      throw err;
    }
  },

  async synchroniserEvenementsEnAttente(): Promise<any> {
    if (isTauri) {
      try {
        const pendingEvents = await invokeTauri<
          Array<{
            id: string;
            nom_table: string;
            type_operation: string;
            entite_id: string;
            donnees_payload: string;
            horodatage_local: string;
            statut_sync: string;
          }>
        >("get_evenements_en_attente");

        if (!pendingEvents || pendingEvents.length === 0) {
          return {
            statut: "SUCCES",
            evenements_traites: 0,
            doublons_ignores: 0,
            erreurs: 0,
            synced_at: new Date().toISOString(),
            message: "Aucun événement local en attente de synchronisation.",
          };
        }

        const formattedEvents = pendingEvents.map((ev) => {
          let parsedPayload: Record<string, any> = {};
          try {
            parsedPayload = JSON.parse(ev.donnees_payload || "{}");
          } catch {
            parsedPayload = {};
          }
          return {
            event_id: ev.id,
            entity_name: ev.nom_table,
            operation: ev.type_operation,
            entity_id: ev.entite_id,
            payload: parsedPayload,
            local_timestamp_utc: ev.horodatage_local,
          };
        });

        const cloudResp = await apiFetch<any>("/sync/batch", {
          method: "POST",
          body: JSON.stringify({ events: formattedEvents }),
        });

        await invokeTauri("marquer_evenements_synchronises", {
          eventIds: pendingEvents.map((e) => e.id),
        });

        return cloudResp;
      } catch (err) {
        console.warn(
          "Serveur Go non joignable pour /sync/batch, conservation en file locale:",
          err,
        );
        throw err;
      }
    }
    return {
      statut: "SUCCES",
      evenements_traites: 0,
      doublons_ignores: 0,
      erreurs: 0,
      synced_at: new Date().toISOString(),
      message:
        "Mode Web connecté directement à la base de données Neon PostgreSQL.",
    };
  },

  async updateCandidatPhoto(
    candidatId: string,
    photoUrl: string,
  ): Promise<any> {
    try {
      const res = await apiFetch<any>(
        `/candidats/${encodeURIComponent(candidatId)}/photo`,
        {
          method: "PATCH",
          body: JSON.stringify({ photo_url: photoUrl }),
        },
      );
      mirrorToTauri("update_candidat_photo", { candidatId, photoUrl });
      realtimeSocket.sendEvent("CANDIDAT_DOSSIER_UPDATED", {
        candidat_id: candidatId,
      });
      return res;
    } catch (err) {
      if (isTauri) {
        const local = await invokeTauri("update_candidat_photo", {
          candidatId,
          photoUrl,
        });
        realtimeSocket.sendEvent("CANDIDAT_DOSSIER_UPDATED", {
          candidat_id: candidatId,
        });
        return local;
      }
      throw err;
    }
  },

  async accorderDerogation(input: {
    dossier_id: string;
    code_pin: string;
    motif?: string;
  }): Promise<any> {
    try {
      const res = await apiFetch<any>("/manager/derogation", {
        method: "POST",
        body: JSON.stringify(input),
      });
      mirrorToTauri("accorder_derogation", { payload: input });
      realtimeSocket.sendEvent("CANDIDAT_DOSSIER_UPDATED", res);
      return res;
    } catch (err) {
      if (isTauri) {
        const local = await invokeTauri("accorder_derogation", {
          payload: input,
        });
        realtimeSocket.sendEvent("CANDIDAT_DOSSIER_UPDATED", local);
        return local;
      }
      throw err;
    }
  },
};

// ==============================================================================
// AutoSchool - Actions Caisse, Dérogations PIN & Séances Code (Guichet Store)
// ==============================================================================

import { api } from "../../db/api";
import { neonSyncEngine } from "../../services/neonSyncService";

export async function execLoadSessionCaisse(store: any) {
  try {
    store.sessionCaisseDetail = await api.getSessionCaisseDetail();
  } catch (err: any) {
    store.showToast(
      "Erreur chargement session caisse: " + err.message,
      "error",
    );
  }
}

export async function execCloturerCaisse(
  store: any,
  payload: { total_especes_compte: number; observations?: string },
) {
  store.loading = true;
  try {
    const res = await api.cloturerCaisse(payload);
    store.sessionCaisseDetail = res;
    store.stats = await api.getGuichetStats();
    neonSyncEngine.setPendingCount(store.stats.evenements_sync_en_attente);
    store.showToast(
      "Caisse clôturée avec succès. Rapport Z généré.",
      "success",
    );
    return res;
  } catch (err: any) {
    store.showToast("Erreur clôture caisse: " + err.message, "error");
    throw err;
  } finally {
    store.loading = false;
  }
}

export async function execReprintReceipt(store: any, paiementId: string) {
  store.loading = true;
  try {
    const data = await api.reimprimerRecu(paiementId);
    store.lastReceipt = {
      numero_recu: data.numero_recu,
      code_securite_hash: data.code_securite_hash,
      candidat_nom: data.candidat_nom,
      numero_dossier: data.numero_dossier,
      formule_nom: data.formule_nom,
      montant_verse: data.montant_verse,
      solde_restant: data.solde_restant,
      mode_reglement: data.mode_reglement,
      date_paiement: data.date_paiement,
      caissier_nom: data.caissier_nom,
      est_duplicata: true,
      nombre_impressions: data.nombre_impressions,
    };
    store.isReceiptModalOpen = true;
    store.showToast(
      `Reçu ${data.numero_recu} réimprimé (Mention DUPLICATA N° ${data.nombre_impressions})`,
      "info",
    );
  } catch (err: any) {
    store.showToast("Erreur réimpression reçu: " + err.message, "error");
  } finally {
    store.loading = false;
  }
}

export async function execUnlockConduiteWithPin(
  store: any,
  pin: string,
  dossierId: string,
  motif?: string,
): Promise<boolean> {
  try {
    const resp = await api.accorderDerogation({
      dossier_id: dossierId,
      code_pin: pin,
      motif,
    });
    if (resp && (resp.success || resp.status === "success" || !resp.error)) {
      if (!store.derogationsAccorde.includes(dossierId)) {
        store.derogationsAccorde.push(dossierId);
      }
      if (
        store.activeDossierDetail &&
        store.activeDossierDetail.dossier.id === dossierId
      ) {
        store.activeDossierDetail.dossier.derogation_conduite_accordee = true;
      }
      store.isManagerOverrideModalOpen = false;
      store.showToast(
        "Dérogation Manager enregistrée ! Accès à la conduite déverrouillé.",
        "success",
      );
      return true;
    }
  } catch (err: any) {
    store.showToast(
      err.message || "Code Manager invalide ou refus de dérogation",
      "error",
    );
  }
  return false;
}

export async function execLoadGerantDashboard(store: any) {
  store.loading = true;
  try {
    const data = await api.getImpayesEtAlertesGerant();
    if (data && typeof data === "object") {
      store.gerantDashboardData = {
        ...data,
        impayes: Array.isArray(data.impayes) ? data.impayes : [],
        eleves_assidus_code: Array.isArray(data.eleves_assidus_code)
          ? data.eleves_assidus_code
          : [],
        inscriptions_formules: Array.isArray(data.inscriptions_formules)
          ? data.inscriptions_formules
          : [],
        evolution_recettes: Array.isArray(data.evolution_recettes)
          ? data.evolution_recettes
          : [],
        ventilation_modes: Array.isArray(data.ventilation_modes)
          ? data.ventilation_modes
          : [],
        paiements_recents: Array.isArray(data.paiements_recents)
          ? data.paiements_recents
          : [],
      };
    } else {
      store.gerantDashboardData = data;
    }
  } catch (err: any) {
    store.showToast(
      "Erreur chargement tableau de bord gérant: " + err.message,
      "error",
    );
  } finally {
    store.loading = false;
  }
}

export async function execEnregistrerSeanceCode(
  store: any,
  payload: {
    theme_traite: string;
    formateur_nom: string;
    heure_debut: string;
    heure_fin: string;
    presences: Array<{
      candidat_id: string;
      est_present: boolean;
      score_test: number;
    }>;
  },
) {
  store.loading = true;
  try {
    const res = await api.enregistrerSeanceCode(payload);
    store.showToast(
      `Séance enregistrée ! ${res.total_presents} élèves pointés (+2h de code créditées par présent).`,
      "success",
    );
    if (store.selectedCandidatId) {
      await store.selectCandidat(store.selectedCandidatId);
    }
    store.stats = await api.getGuichetStats();
    neonSyncEngine.setPendingCount(store.stats.evenements_sync_en_attente);
    return res;
  } catch (err: any) {
    store.showToast(
      "Erreur enregistrement séance code: " + err.message,
      "error",
    );
  } finally {
    store.loading = false;
  }
}

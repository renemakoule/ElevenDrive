// ==============================================================================
// AutoSchool - Actions Flotte Automobile & Catalogue Formules (Guichet Store)
// ==============================================================================

import { api } from "../../db/api";
import type {
  CreerVehiculeInput,
  StatutVehicule,
  TranchePaiementConfig,
  Vehicule,
} from "../../types";

const VEHICULES_CACHE_KEY = "autoschool_vehicules_cache";

function persistVehiculesCache(vehicules: Vehicule[]) {
  if (typeof localStorage !== "undefined") {
    localStorage.setItem(VEHICULES_CACHE_KEY, JSON.stringify(vehicules));
  }
}

function getErrorMessage(err: any): string {
  if (!err) return "Erreur inconnue";
  if (typeof err === "string") return err;
  if (typeof err.message === "string" && err.message.trim()) return err.message;
  if (typeof err.error === "string" && err.error.trim()) return err.error;
  return String(err);
}

export async function execLoadFormules(store: any) {
  try {
    const list = await api.getFormules();
    store.formules = Array.isArray(list) ? list : [];
  } catch (err: any) {
    if (!Array.isArray(store.formules)) store.formules = [];
    store.showToast(
      "Erreur chargement des formules: " + getErrorMessage(err),
      "error",
    );
  }
}

export async function execCreerFormule(
  store: any,
  input: {
    categorie_permis: string;
    intitule: string;
    description?: string;
    heures_theoriques: number;
    heures_pratiques: number;
    prix_total: number;
    acompte_minimum: number;
    tranches_paiement?: TranchePaiementConfig[];
  },
) {
  store.loading = true;
  try {
    const nouvelleFormule = await api.creerFormule(input);
    if (!Array.isArray(store.formules)) {
      store.formules = [];
    }
    store.formules.push(nouvelleFormule);
    store.showToast(
      `Formule ${nouvelleFormule.intitule} créée et prix fixé à ${nouvelleFormule.prix_total.toLocaleString("fr-FR")} FCFA.`,
      "success",
    );
    return nouvelleFormule;
  } catch (err: any) {
    store.showToast(
      `Erreur création formule: ${getErrorMessage(err)}`,
      "error",
    );
    throw err;
  } finally {
    store.loading = false;
  }
}

export async function execModifierPrixFormule(
  store: any,
  formuleId: string,
  nouveauPrix: number,
  acompteMinimum?: number,
  motif?: string,
  tranchesPaiement?: TranchePaiementConfig[],
) {
  store.loading = true;
  try {
    const formuleMaj = await api.modifierPrixFormule({
      formule_id: formuleId,
      nouveau_prix: nouveauPrix,
      acompte_minimum: acompteMinimum,
      motif,
      tranches_paiement: tranchesPaiement,
    });
    const index = store.formules.findIndex((f: any) => f.id === formuleId);
    if (index !== -1) {
      store.formules[index] = formuleMaj;
    }
    store.showToast(
      `Prix de la formule mis à jour à ${formuleMaj.prix_total.toLocaleString("fr-FR")} FCFA par le Directeur.`,
      "success",
    );
    return formuleMaj;
  } catch (err: any) {
    store.showToast(
      `Erreur modification du prix: ${getErrorMessage(err)}`,
      "error",
    );
    throw err;
  } finally {
    store.loading = false;
  }
}

export async function execSupprimerFormule(store: any, formuleId: string) {
  store.loading = true;
  try {
    await api.supprimerFormule(formuleId);
    store.formules = store.formules.filter((f: any) => f.id !== formuleId);
    store.showToast("Formule retirée du catalogue avec succès.", "success");
  } catch (err: any) {
    store.showToast(
      `Erreur suppression formule: ${getErrorMessage(err)}`,
      "error",
    );
    throw err;
  } finally {
    store.loading = false;
  }
}

export async function execViderCatalogueFormules(store: any) {
  store.loading = true;
  try {
    await api.supprimerToutesFormules();
    store.formules = [];
    store.showToast(
      "Le catalogue a été vidé avec succès. Vous pouvez maintenant définir vos propres tarifs.",
      "success",
    );
  } catch (err: any) {
    store.showToast(
      `Erreur lors du nettoyage du catalogue: ${getErrorMessage(err)}`,
      "error",
    );
    throw err;
  } finally {
    store.loading = false;
  }
}

export async function execLoadVehicules(store: any) {
  store.flotteLoading = true;
  try {
    const list = await api.getVehicules();
    if (Array.isArray(list)) {
      store.vehicules = list;
      persistVehiculesCache(list);
      store.flotteLoading = false;
      return list;
    }
  } catch {
    // Vérifier le cache local hors-ligne
  }

  if (typeof localStorage !== "undefined") {
    const cached = localStorage.getItem(VEHICULES_CACHE_KEY);
    if (cached) {
      try {
        const parsed = JSON.parse(cached);
        if (Array.isArray(parsed)) {
          store.vehicules = parsed.filter(
            (v: any) => v.id !== "veh-01" && v.id !== "veh-02",
          );
          store.flotteLoading = false;
          return store.vehicules;
        }
      } catch {}
    }
  }

  store.flotteLoading = false;
  return store.vehicules;
}

export async function execAjouterVehicule(
  store: any,
  input: CreerVehiculeInput,
) {
  store.flotteLoading = true;
  try {
    const nouveauVehicule = await api.creerVehicule(input);
    store.vehicules.push(nouveauVehicule);
    persistVehiculesCache(store.vehicules);
    store.showToast(
      `Véhicule ${input.immatriculation} ajouté à la flotte avec succès !`,
      "success",
    );
    return nouveauVehicule;
  } catch {
    const localVeh: Vehicule = {
      id: "veh_" + Date.now(),
      agence_id: input.agence_id || "ag-primary",
      immatriculation: input.immatriculation.toUpperCase(),
      marque_modele: input.marque_modele,
      type_boite: input.type_boite || "MANUELLE",
      kilometrage_actuel: input.kilometrage_actuel || 0,
      date_visite_technique: input.date_visite_technique,
      date_expiration_assurance: input.date_expiration_assurance,
      statut: input.statut || "DISPONIBLE",
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    };
    store.vehicules.push(localVeh);
    persistVehiculesCache(store.vehicules);
    store.showToast(
      `Véhicule ${input.immatriculation} enregistré dans la flotte !`,
      "success",
    );
    return localVeh;
  } finally {
    store.flotteLoading = false;
  }
}

export async function execModifierStatutVehicule(
  store: any,
  vehiculeId: string,
  statut: StatutVehicule,
) {
  const applyLocalStatus = () => {
    const v = store.vehicules.find((item: any) => item.id === vehiculeId);
    if (v) {
      v.statut = statut;
    }
    persistVehiculesCache(store.vehicules);
  };

  try {
    await api.updateVehiculeStatut(vehiculeId, statut);
    applyLocalStatus();
    store.showToast(`Statut du véhicule mis à jour : ${statut}`, "success");
  } catch {
    applyLocalStatus();
    store.showToast(`Statut mis à jour : ${statut}`, "info");
  }
}

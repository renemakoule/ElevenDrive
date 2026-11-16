// ==============================================================================
// AutoSchool - Actions Profil, École, Agences & Temps Réel Super-Admin
// ==============================================================================

import type { Agence, AutoEcole, UserProfile } from "../../types";
import {
  authBackend,
  buildSchoolAccessLinks,
} from "../../services/authBackendService";
import { USER_PROFILES } from "./authGetters";

export async function execLoadProfileByRole(store: any): Promise<void> {
  const remote = await authBackend.getMyProfile();
  if (remote && remote.user) {
    store.customUserProfile = {
      ...store.customUserProfile,
      ...remote.user,
      dossier: store.customUserProfile?.dossier || remote.user.dossier,
    };
    store.currentRole = remote.user.role;
    if (remote.school) {
      store.currentSchool = remote.school;
      store.schoolAccessLinks =
        remote.school.links ||
        buildSchoolAccessLinks(remote.school.code_unique);
    }
    if (remote.agence) {
      store.currentAgence = remote.agence;
    }
    const session = authBackend.getActiveSession();
    if (session) {
      session.user = store.customUserProfile;
      if (remote.school) session.auto_ecole = remote.school;
      if (remote.agence) session.agence = remote.agence;
      authBackend.saveSession(session);
    }
  } else if (typeof localStorage !== "undefined") {
    const savedWeb = localStorage.getItem("autoschool_web_session");
    if (savedWeb) {
      try {
        const parsed = JSON.parse(savedWeb);
        if (parsed.user) {
          store.customUserProfile = parsed.user;
          if (parsed.role) store.currentRole = parsed.role;
        }
      } catch {}
    }
  }

  if (store.rolePermissions.canViewAgences) {
    await store.loadAgences();
  } else {
    store.loadAgencesLocally();
  }
}

export async function execUpdateUserProfile(
  store: any,
  updates: Partial<UserProfile> & { nouveau_mot_de_passe?: string },
): Promise<void> {
  if (!store.rolePermissions.canConfigurePin) {
    delete updates.code_pin_manager;
  }

  try {
    await authBackend.updateMyProfile({
      nom: updates.nom || store.currentUser.nom,
      prenoms: updates.prenom || store.currentUser.prenom,
      telephone: updates.telephone || store.currentUser.telephone || "",
      email: updates.email,
      mot_de_passe: updates.nouveau_mot_de_passe || undefined,
      code_pin_manager: store.rolePermissions.canConfigurePin
        ? updates.code_pin_manager
        : undefined,
    });
  } catch (err) {
    console.warn("Mise à jour profil en mode local (API non jointe):", err);
  }

  const { nouveau_mot_de_passe: _pw, ...profileUpdates } = updates;
  if (store.customUserProfile) {
    store.customUserProfile = {
      ...store.customUserProfile,
      ...profileUpdates,
      initial: (
        profileUpdates.prenom?.[0] ||
        store.customUserProfile.prenom?.[0] ||
        "U"
      ).toUpperCase(),
    };
  } else {
    store.customUserProfile = {
      ...USER_PROFILES[store.currentRole as keyof typeof USER_PROFILES],
      ...profileUpdates,
      role: store.currentRole,
    };
  }

  const session = authBackend.getActiveSession();
  if (session) {
    session.user = { ...session.user, ...profileUpdates };
    authBackend.saveSession(session);
  }

  if (typeof localStorage !== "undefined" && store.platform === "pwa") {
    const savedWeb = localStorage.getItem("autoschool_web_session");
    if (savedWeb) {
      try {
        const parsed = JSON.parse(savedWeb);
        parsed.user = store.customUserProfile;
        localStorage.setItem("autoschool_web_session", JSON.stringify(parsed));
      } catch {}
    }
  }
}

export async function execUpdateSchool(
  store: any,
  updates: Partial<AutoEcole>,
): Promise<void> {
  if (!store.rolePermissions.canManageSchool) {
    throw new Error(
      "Accès refusé : seul le Directeur est habilité à modifier les paramètres de l'auto-école.",
    );
  }

  try {
    const updatedRemote = await authBackend.updateSchoolInfo(updates);
    if (updatedRemote) {
      updates = { ...updates, ...updatedRemote };
    }
  } catch (err) {
    console.warn("Mise à jour auto-école en mode local:", err);
  }

  if (store.currentSchool) {
    store.currentSchool = { ...store.currentSchool, ...updates };
  }
  const newSlug = updates.code_unique || store.currentSchool?.code_unique;
  if (newSlug) {
    store.schoolAccessLinks = buildSchoolAccessLinks(newSlug);
  }
  const session = authBackend.getActiveSession();
  if (session) {
    session.auto_ecole = { ...session.auto_ecole, ...updates };
    if (newSlug) {
      session.auto_ecole.links = buildSchoolAccessLinks(newSlug);
    }
    authBackend.saveSession(session);
  }
}

export async function execLoadAgences(store: any): Promise<void> {
  try {
    const agences = await authBackend.getSchoolAgences();
    if (Array.isArray(agences) && agences.length > 0) {
      store.agencesList = agences;
      store.saveAgencesLocally();
      return;
    }
  } catch (err) {
    console.warn("Échec chargement agences depuis l'API:", err);
  }
  store.loadAgencesLocally();
}

export async function execUpdateAgence(
  store: any,
  agenceId: string,
  updates: Partial<Agence>,
): Promise<void> {
  if (store.currentAgence && store.currentAgence.id === agenceId) {
    store.currentAgence = { ...store.currentAgence, ...updates };
  }
  const idx = store.agencesList.findIndex((a: Agence) => a.id === agenceId);
  if (idx !== -1) {
    store.agencesList[idx] = { ...store.agencesList[idx], ...updates };
  }
  store.saveAgencesLocally();
  const session = authBackend.getActiveSession();
  if (session && session.agence?.id === agenceId) {
    session.agence = { ...session.agence, ...updates };
    authBackend.saveSession(session);
  }

  try {
    await authBackend.updateAgence(agenceId, updates);
  } catch (err) {
    console.warn("Échec persistance modification agence en base:", err);
  }
}

export async function execAddAgence(
  store: any,
  newAgence: Omit<Agence, "id">,
): Promise<Agence> {
  try {
    const created = await authBackend.createAgence(newAgence);
    store.agencesList = store.agencesList.filter(
      (a: Agence) =>
        a.id !== created.id && a.code_agence !== created.code_agence,
    );
    store.agencesList.push(created);
    store.saveAgencesLocally();
    return created;
  } catch (err) {
    console.error("Erreur addAgence API:", err);
    throw err;
  }
}

export async function execDeleteAgence(
  store: any,
  agenceId: string,
): Promise<void> {
  try {
    await authBackend.deleteAgence(agenceId);
  } catch (err) {
    console.warn("Erreur suppression backend agence:", err);
  }
  store.agencesList = store.agencesList.filter(
    (a: Agence) => a.id !== agenceId,
  );
  store.saveAgencesLocally();
}

export function execSaveAgencesLocally(store: any): void {
  if (typeof localStorage === "undefined") return;
  const key = `autoschool_agences_${store.currentSchool?.id || "default"}`;
  localStorage.setItem(key, JSON.stringify(store.agencesList));
}

export function execLoadAgencesLocally(store: any): void {
  if (typeof localStorage === "undefined") return;
  const key = `autoschool_agences_${store.currentSchool?.id || "default"}`;
  const raw = localStorage.getItem(key);
  if (raw) {
    try {
      const parsed = JSON.parse(raw);
      if (Array.isArray(parsed) && parsed.length > 0) {
        store.agencesList = parsed;
        return;
      }
    } catch {}
  }
  if (store.currentAgence) {
    store.agencesList = [store.currentAgence];
  } else if (store.currentSchool) {
    const cityPrefix =
      (store.currentSchool.ville || "HQ")
        .normalize("NFD")
        .replace(/[\u0300-\u036f]/g, "")
        .replace(/[^a-zA-Z0-9]/g, "")
        .slice(0, 3)
        .toUpperCase() || "HQ";

    store.agencesList = [
      {
        id: "ag-primary",
        auto_ecole_id: store.currentSchool.id,
        code_agence: `AG-${cityPrefix}-01`,
        nom: `Siège · ${store.currentSchool.ville || "Principal"}`,
        adresse: "Adresse du siège",
        telephone: store.currentSchool.telephone,
        ville: store.currentSchool.ville || "Douala",
        commune: "Centre",
        est_siege: true,
      },
    ];
  }
}

export function execHandleRealtimeSchoolEvent(
  store: any,
  event: { type: string; data: any },
): void {
  if (!store.currentSchool) return;
  const data = event.data || {};
  const schoolId = data.school_id || data.id;
  if (schoolId && schoolId !== store.currentSchool.id) return;

  console.log(
    "[AuthStore] Réception événement temps réel SuperAdmin:",
    event.type,
    data,
  );

  const updatedSchool = { ...store.currentSchool };

  if (event.type === "SCHOOL_LICENCE_UPDATED") {
    if (data.type_licence !== undefined)
      updatedSchool.type_licence = data.type_licence;
    if (data.max_machines !== undefined)
      updatedSchool.max_machines = data.max_machines;
    if (data.machines_actives !== undefined)
      updatedSchool.machines_actives = data.machines_actives;
    if (data.date_expiration_licence !== undefined) {
      updatedSchool.date_expiration_licence = data.date_expiration_licence;
    }
    if (data.est_actif !== undefined) updatedSchool.est_actif = data.est_actif;
  } else if (event.type === "SCHOOL_STATUS_CHANGED") {
    if (data.est_actif !== undefined) updatedSchool.est_actif = data.est_actif;
  } else if (event.type === "SCHOOL_MODULES_UPDATED") {
    if (data.modules_actifs !== undefined) {
      updatedSchool.modules_actifs =
        typeof data.modules_actifs === "string"
          ? data.modules_actifs
          : JSON.stringify(data.modules_actifs);
    }
  } else if (event.type === "SCHOOL_UPDATED") {
    Object.assign(updatedSchool, data);
  } else if (event.type === "SCHOOL_DELETED") {
    store.logout();
    if (typeof alert !== "undefined") {
      alert(
        data.message ||
          "Cet établissement a été supprimé par l'administration centrale.",
      );
    }
    return;
  }

  store.currentSchool = updatedSchool;

  if (
    store.activeSection !== "overview" &&
    !store.isModuleEnabled(store.activeSection)
  ) {
    store.activeSection = "overview";
  }

  const session = authBackend.getActiveSession();
  if (session && session.auto_ecole) {
    session.auto_ecole = { ...session.auto_ecole, ...store.currentSchool };
    authBackend.saveSession(session);
  }

  store.loadProfileByRole().catch(() => {});
}

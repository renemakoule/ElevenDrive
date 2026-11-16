// ==============================================================================
// AutoSchool - Client CRUD Collaborateurs & Profil Établissement
// ==============================================================================

import type {
  AutoEcole,
  Agence,
  UserProfile,
  RegisterUserPayload,
  UserTeamMember,
} from "../../types";
import {
  API_BASE_URL,
  buildSchoolAccessLinks,
  fetchAuthWithFallback,
  getActiveSessionFromStorage,
  mapRoleToUI,
} from "./authSessionHelpers";
import { fetchSchoolAgences } from "./authAgenceClient";

export {
  fetchSchoolAgences,
  createSchoolAgence,
  updateSchoolAgence,
  deleteSchoolAgence,
} from "./authAgenceClient";

const inFlightRegisterUsers = new Map<string, Promise<any>>();

/** Récupérer la liste des collaborateurs de l'auto-école */
export async function fetchSchoolUsers(): Promise<UserTeamMember[]> {
  const session = getActiveSessionFromStorage();
  const token = session?.token;
  const localTeamKey = `elevendrive_team_${session?.auto_ecole?.id || "default"}`;

  if (token) {
    try {
      const response = await fetchAuthWithFallback("/users", {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        const result = await response.json();
        if (result.success && Array.isArray(result.data)) {
          if (typeof localStorage !== "undefined") {
            localStorage.setItem(localTeamKey, JSON.stringify(result.data));
          }
          return result.data;
        }
      }
    } catch {
      // Lire la copie locale hors-ligne
    }
  }

  if (typeof localStorage !== "undefined") {
    const cached = localStorage.getItem(localTeamKey);
    if (cached) {
      try {
        return JSON.parse(cached);
      } catch {
        // ignore
      }
    }
  }

  if (session) {
    return [
      {
        id: session.user.id || "usr_dir_default",
        auto_ecole_id: session.auto_ecole?.id || "default",
        agence_id: session.agence?.id,
        agence_nom: session.agence?.nom || "Agence Siège",
        nom: session.user.nom || "Directeur",
        prenoms: session.user.prenom || "Gérant",
        telephone: session.user.telephone || "",
        email: session.user.email,
        role: "DIRECTEUR",
        statut: "ACTIF",
        created_at: new Date().toISOString().slice(0, 16).replace("T", " "),
      },
    ];
  }
  return [];
}

/** Enregistrer un nouveau collaborateur au sein de l'auto-école */
export async function registerSchoolUser(
  payload: RegisterUserPayload,
): Promise<any> {
  const cleanTel = (payload.telephone || "").replace(/\s+/g, "");
  if (cleanTel && inFlightRegisterUsers.has(cleanTel)) {
    return inFlightRegisterUsers.get(cleanTel);
  }

  const task = (async () => {
    const session = getActiveSessionFromStorage();
    const token = session?.token;

    let response: Response | null = null;
    let networkError: any = null;

    try {
      response = await fetchAuthWithFallback("/register-user", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          ...(token ? { Authorization: `Bearer ${token}` } : {}),
        },
        body: JSON.stringify({
          ...payload,
          auto_ecole_id: session?.auto_ecole.id,
        }),
      });
    } catch (err) {
      networkError = err;
    }

    if (response) {
      const result = await response.json();
      if (!response.ok || !result.success) {
        throw new Error(
          result.error || "Erreur lors de la création du collaborateur.",
        );
      }
      return result.data;
    }

    if (networkError) {
      console.warn(
        "Backend non joignable lors de l'enregistrement:",
        networkError,
      );
    }

    if (typeof localStorage !== "undefined" && session) {
      const localTeamKey = `elevendrive_team_${session.auto_ecole.id || "default"}`;
      const existing = await fetchSchoolUsers();
      const agences = await fetchSchoolAgences();
      const targetAgence =
        agences.find((a) => a.id === payload.agence_id) || session.agence;

      const newMember: UserTeamMember = {
        id: "usr_" + Date.now(),
        auto_ecole_id: session.auto_ecole.id,
        agence_id: payload.agence_id || session.agence?.id,
        agence_nom: targetAgence?.nom || "Agence Principale",
        nom: payload.nom.toUpperCase(),
        prenoms: payload.prenoms,
        telephone: payload.telephone,
        email: payload.email,
        role: payload.role,
        statut: "ACTIF",
        created_at: new Date().toISOString().slice(0, 16).replace("T", " "),
      };
      existing.push(newMember);
      localStorage.setItem(localTeamKey, JSON.stringify(existing));
      return newMember;
    }

    throw new Error("Impossible d'enregistrer le collaborateur.");
  })();

  if (cleanTel) {
    inFlightRegisterUsers.set(cleanTel, task);
  }

  try {
    return await task;
  } finally {
    if (cleanTel) {
      inFlightRegisterUsers.delete(cleanTel);
    }
  }
}

/** Modifier un collaborateur existant */
export async function updateSchoolUser(
  userId: string,
  payload: Partial<RegisterUserPayload> & { statut?: string },
): Promise<any> {
  const session = getActiveSessionFromStorage();
  const token = session?.token;
  const localTeamKey = `elevendrive_team_${session?.auto_ecole.id || "default"}`;

  let response: Response | null = null;
  try {
    response = await fetch(
      `${API_BASE_URL}/users?id=${encodeURIComponent(userId)}`,
      {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          ...(token ? { Authorization: `Bearer ${token}` } : {}),
        },
        body: JSON.stringify(payload),
      },
    );
  } catch {
    // Ignorer si offline
  }

  if (response) {
    const result = await response.json();
    if (!response.ok || !result.success) {
      throw new Error(
        result.error || "Erreur lors de la modification du collaborateur.",
      );
    }
    return result.data;
  }

  if (typeof localStorage !== "undefined" && session) {
    const existing = await fetchSchoolUsers();
    const agences = await fetchSchoolAgences();
    const targetAgence = agences.find((a) => a.id === payload.agence_id);

    const updated = existing.map((m) => {
      if (m.id === userId) {
        return {
          ...m,
          nom: payload.nom ? payload.nom.toUpperCase() : m.nom,
          prenoms: payload.prenoms || m.prenoms,
          telephone: payload.telephone || m.telephone,
          email: payload.email !== undefined ? payload.email : m.email,
          role: payload.role || m.role,
          agence_id: payload.agence_id || m.agence_id,
          agence_nom: targetAgence ? targetAgence.nom : m.agence_nom,
          statut: payload.statut || m.statut,
        };
      }
      return m;
    });
    localStorage.setItem(localTeamKey, JSON.stringify(updated));
    return updated.find((m) => m.id === userId);
  }

  throw new Error("Impossible de modifier le collaborateur.");
}

/** Supprimer un collaborateur */
export async function deleteSchoolUser(userId: string): Promise<boolean> {
  const session = getActiveSessionFromStorage();
  const token = session?.token;
  const localTeamKey = `elevendrive_team_${session?.auto_ecole.id || "default"}`;

  let response: Response | null = null;
  try {
    response = await fetch(
      `${API_BASE_URL}/users?id=${encodeURIComponent(userId)}`,
      {
        method: "DELETE",
        headers: token ? { Authorization: `Bearer ${token}` } : {},
      },
    );
  } catch {
    // Ignorer si offline
  }

  if (response) {
    const result = await response.json();
    if (!response.ok || !result.success) {
      throw new Error(
        result.error || "Erreur lors de la suppression du collaborateur.",
      );
    }
    if (typeof localStorage !== "undefined") {
      const existing = await fetchSchoolUsers();
      const filtered = existing.filter((m) => m.id !== userId);
      localStorage.setItem(localTeamKey, JSON.stringify(filtered));
    }
    return true;
  }

  if (typeof localStorage !== "undefined" && session) {
    const existing = await fetchSchoolUsers();
    const filtered = existing.filter((m) => m.id !== userId);
    localStorage.setItem(localTeamKey, JSON.stringify(filtered));
    return true;
  }

  throw new Error("Impossible de supprimer le collaborateur.");
}

/** Récupérer le profil complet de l'utilisateur connecté */
export async function fetchMyProfile(): Promise<{
  user: UserProfile;
  school?: AutoEcole;
  agence?: Agence;
  permissions?: {
    can_manage_school: boolean;
    can_manage_agences: boolean;
    can_view_agences: boolean;
    can_configure_pin: boolean;
    can_manage_users: boolean;
    can_change_password: boolean;
  };
} | null> {
  const session = getActiveSessionFromStorage();
  const token = session?.token;
  if (!token) return null;

  try {
    const response = await fetchAuthWithFallback("/me", {
      headers: { Authorization: `Bearer ${token}` },
      cache: "no-store",
    });
    if (!response.ok) return null;

    const result = await response.json();
    if (!result.success || !result.data?.user) return null;

    const u = result.data.user;
    const s = result.data.school;
    const a = result.data.agence;
    const { role, roleLabel } = mapRoleToUI(
      u.role || result.data.role || "SECRETAIRE",
    );

    const userProfile: UserProfile = {
      id: u.id,
      nom: u.nom,
      prenom: u.prenoms,
      role,
      roleLabel,
      initial: (u.prenoms?.[0] || u.nom?.[0] || "U").toUpperCase(),
      agence: a?.nom || session?.agence?.nom || "Agence Principale",
      auto_ecole_id: u.auto_ecole_id || s?.id,
      agence_id: u.agence_id || a?.id || "",
      telephone: u.telephone,
      email: u.email || undefined,
      code_pin_manager: session?.user?.code_pin_manager,
    };

    const school: AutoEcole | undefined = s
      ? {
          ...s,
          ville: a?.ville || session?.auto_ecole?.ville || "Douala",
          links: buildSchoolAccessLinks(s.code_unique),
        }
      : undefined;

    return {
      user: userProfile,
      school,
      agence: a || undefined,
      permissions: result.data.permissions,
    };
  } catch {
    return null;
  }
}

/** Mettre à jour le profil de l'utilisateur connecté */
export async function updateMyUserProfile(payload: {
  nom: string;
  prenoms: string;
  telephone: string;
  email?: string;
  mot_de_passe?: string;
  code_pin_manager?: string;
}): Promise<any> {
  const session = getActiveSessionFromStorage();
  const token = session?.token;
  if (!token) return null;

  const response = await fetchAuthWithFallback("/me", {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(payload),
  });
  const result = await response.json();
  if (!response.ok || !result.success) {
    throw new Error(result.error || "Erreur lors de la mise à jour du profil.");
  }
  return result.data;
}

/** Mettre à jour les informations de l'Auto-École (Réservé au Directeur) */
export async function updateSchoolDetails(
  payload: Partial<AutoEcole>,
): Promise<AutoEcole | null> {
  const session = getActiveSessionFromStorage();
  const token = session?.token;
  if (!token) return null;

  const response = await fetchAuthWithFallback("/school", {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(payload),
  });
  const result = await response.json();
  if (!response.ok || !result.success) {
    throw new Error(
      result.error || "Erreur lors de la mise à jour de l'établissement.",
    );
  }
  return result.data as AutoEcole;
}

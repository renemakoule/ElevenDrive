// ==============================================================================
// AutoSchool - Helpers Session, Stockage Local & Empreinte Machine
// ==============================================================================

import type {
  AutoEcole,
  UserProfile,
  AuthSessionData,
  SchoolAccessLinks,
  UserRole,
} from "../../types";
import {
  AUTH_API_BASE_URL,
  BACKEND_ORIGIN,
  getPublicAppOrigin,
} from "../../config/env";

export const STORAGE_KEY_ACTIVE_SESSION = "elevendrive_active_session";
export const STORAGE_KEY_REFRESH_TOKEN = "elevendrive_refresh_token";

/** URL de base de l'API Auth (configurée via VITE_API_URL) */
export const API_BASE_URL = AUTH_API_BASE_URL;

/** Obtenir ou générer une empreinte de machine unique et persistante (Desktop) */
export const getMachineFingerprint = (): string => {
  if (typeof localStorage === "undefined") return "";
  let id = localStorage.getItem("elevendrive_machine_id");
  if (!id) {
    id =
      "mach-" +
      Math.random().toString(36).substring(2, 10) +
      "-" +
      Date.now().toString(36);
    localStorage.setItem("elevendrive_machine_id", id);
  }
  return id;
};

/** Formate l'URL de base du portail Web */
export const getBaseUrl = (): string => getPublicAppOrigin();

/** Extrait l'origine du serveur backend depuis VITE_API_URL */
export const getBackendOrigin = (): string => BACKEND_ORIGIN;

/** Générateur de Slug propre (Code Unique d'Auto-École) */
export const generateSchoolSlug = (name: string): string => {
  const clean = name
    .toLowerCase()
    .normalize("NFD")
    .replace(/[\u0300-\u036f]/g, "")
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/-+/g, "-")
    .replace(/^-|-$/g, "");

  return clean || "auto-ecole-" + Math.floor(1000 + Math.random() * 9000);
};

/** Construit les liens publics d'accès d'une auto-école */
export const buildSchoolAccessLinks = (slug: string): SchoolAccessLinks => {
  const base = getBaseUrl();
  const cleanSlug = slug || "auto-ecole";
  const uniqueUrl = `${base}/${cleanSlug}`;
  return {
    portal_url: uniqueUrl,
    moniteur_link: uniqueUrl,
    eleve_link: uniqueUrl,
    inscription_link: uniqueUrl,
  };
};

/** Mapper les rôles de l'API Go vers l'UI */
export const mapRoleToUI = (
  role: string,
): { role: UserRole; roleLabel: string } => {
  switch (role.toUpperCase()) {
    case "DIRECTEUR":
      return { role: "manager", roleLabel: "Manager" };
    case "SECRETAIRE":
      return { role: "secretaire", roleLabel: "Secrétaire" };
    case "CAISSIER":
      return { role: "secretaire", roleLabel: "Caissière" };
    case "MONITEUR":
      return { role: "moniteur", roleLabel: "Moniteur" };
    case "ELEVE":
      return { role: "eleve", roleLabel: "Candidat Permis" };
    default:
      return { role: "secretaire", roleLabel: "Agent" };
  }
};

/** Effectue un fetch sur API_BASE_URL avec fallback /auth/... sur 404 */
export async function fetchAuthWithFallback(
  path: string,
  init?: RequestInit,
): Promise<Response> {
  let response = await fetch(`${API_BASE_URL}${path}`, init);
  if (response.status === 404) {
    response = await fetch(`${API_BASE_URL}/auth${path}`, init);
  }
  return response;
}

/** Construit un objet AuthSessionData normalisé à partir d'un payload API standard */
export function buildSessionFromAuthData(
  data: any,
  defaultAgenceName = "Agence principale",
  overridePin?: string,
): AuthSessionData {
  const { role, roleLabel } = mapRoleToUI(data.user.role);

  const userProfile: UserProfile = {
    id: data.user.id,
    nom: data.user.nom,
    prenom: data.user.prenoms,
    role,
    roleLabel,
    initial: (data.user.prenoms?.[0] || "U").toUpperCase(),
    agence: data.agence?.nom || defaultAgenceName,
    auto_ecole_id: data.school.id,
    agence_id: data.agence?.id || "",
    telephone: data.user.telephone,
    email: data.user.email || undefined,
    code_pin_manager: overridePin || data.user.code_pin_manager || undefined,
  };

  const school: AutoEcole = {
    ...data.school,
    links: buildSchoolAccessLinks(data.school.code_unique),
  };

  if (data.refresh_token && typeof localStorage !== "undefined") {
    localStorage.setItem(STORAGE_KEY_REFRESH_TOKEN, data.refresh_token);
  }

  return {
    user: userProfile,
    auto_ecole: school,
    agence: data.agence,
    token: data.access_token,
  };
}

/** Sauvegarde la session active et mémorise le compte pour le mode hors-ligne */
export function saveSessionToStorage(session: AuthSessionData): void {
  if (typeof localStorage === "undefined") return;
  const serialized = JSON.stringify(session);
  localStorage.setItem(STORAGE_KEY_ACTIVE_SESSION, serialized);
  localStorage.setItem("autoschool_active_session", serialized);

  try {
    const raw = localStorage.getItem("autoschool_known_accounts") || "[]";
    const accounts: any[] = JSON.parse(raw);
    const tel = session.user?.telephone?.replace(/[^0-9]/g, "") || "";
    const email = session.user?.email?.toLowerCase() || "";

    const filtered = accounts.filter((a) => {
      const aTel = a.telephone?.replace(/[^0-9]/g, "") || "";
      const aEmail = a.email?.toLowerCase() || "";
      if (tel && aTel && tel === aTel) return false;
      if (email && aEmail && email === aEmail) return false;
      return true;
    });

    filtered.push({
      email: session.user.email ? session.user.email.toLowerCase() : null,
      telephone: session.user.telephone || null,
      nom: session.user.nom,
      prenom: session.user.prenom,
      role: session.user.role,
      session,
    });
    localStorage.setItem("autoschool_known_accounts", JSON.stringify(filtered));
  } catch {}
}

export function getKnownAccountsFromStorage(): any[] {
  if (typeof localStorage === "undefined") return [];
  try {
    const raw = localStorage.getItem("autoschool_known_accounts");
    return raw ? JSON.parse(raw) : [];
  } catch {
    return [];
  }
}

export function getActiveSessionFromStorage(): AuthSessionData | null {
  if (typeof localStorage === "undefined") return null;
  const raw =
    localStorage.getItem(STORAGE_KEY_ACTIVE_SESSION) ||
    localStorage.getItem("autoschool_active_session");
  if (!raw) return null;
  try {
    return JSON.parse(raw);
  } catch {
    return null;
  }
}

export function clearSessionFromStorage(): void {
  if (typeof localStorage === "undefined") return;
  localStorage.removeItem(STORAGE_KEY_ACTIVE_SESSION);
  localStorage.removeItem("autoschool_active_session");
}

/** Authentification locale hors-ligne si Internet ou Neon est coupé */
export function authenticateOfflineFromStorage(
  identifiant: string,
): AuthSessionData | null {
  const cleanId = identifiant.trim().toLowerCase();
  const cleanPhone = cleanId.replace(/[^0-9]/g, "");

  const active = getActiveSessionFromStorage();
  if (active && active.user) {
    const uTel = (active.user.telephone || "").replace(/[^0-9]/g, "");
    const uEmail = (active.user.email || "").toLowerCase();
    if ((cleanPhone && uTel === cleanPhone) || (uEmail && uEmail === cleanId)) {
      return active;
    }
  }

  const accounts = getKnownAccountsFromStorage();
  const found = accounts.find((a) => {
    const aTel = (a.telephone || "").replace(/[^0-9]/g, "");
    const aEmail = (a.email || "").toLowerCase();
    if (cleanPhone && aTel && aTel === cleanPhone) return true;
    if (aEmail && aEmail === cleanId) return true;
    return false;
  });

  if (found && found.session) {
    return found.session;
  }

  return null;
}

// ==============================================================================
// AutoSchool - Client OAuth & Navigation Externe Desktop
// ==============================================================================

import type {
  AutoEcole,
  UserProfile,
  AuthSessionData,
  OAuthResult,
} from "../../types";
import {
  API_BASE_URL,
  buildSchoolAccessLinks,
  buildSessionFromAuthData,
  fetchAuthWithFallback,
  getActiveSessionFromStorage,
  getBackendOrigin,
  getKnownAccountsFromStorage,
  mapRoleToUI,
  saveSessionToStorage,
} from "./authSessionHelpers";

/**
 * Connexion via un fournisseur OAuth (Google, Facebook, LinkedIn, Apple)
 */
export async function loginWithOAuthProvider(
  provider: "google" | "facebook" | "linkedin" | "apple",
  profileData?: {
    email?: string;
    nom?: string;
    prenoms?: string;
    id_token?: string;
  },
): Promise<OAuthResult> {
  const providerNames: Record<string, string> = {
    google: "Google",
    facebook: "Facebook",
    linkedin: "LinkedIn",
    apple: "Apple",
  };
  const providerName = providerNames[provider] || provider;

  let result: any = null;
  try {
    const response = await fetch(`${API_BASE_URL}/oauth`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        provider,
        email: profileData?.email || "",
        nom: profileData?.nom || "",
        prenoms: profileData?.prenoms || "",
        id_token: profileData?.id_token || "",
      }),
    });
    result = await response.json();
  } catch {
    // Backend temporairement non joignable ou mode déconnecté
  }

  // 1. Si l'API backend renvoie authenticated
  if (
    result &&
    result.success &&
    result.data?.status === "authenticated" &&
    result.data.auth_response
  ) {
    const data = result.data.auth_response;
    const session = buildSessionFromAuthData(data, "Agence principale");
    saveSessionToStorage(session);

    return {
      status: "authenticated",
      session,
      provider,
      profile: {
        email: profileData?.email || data.user.email,
        nom: data.user.nom,
        prenoms: data.user.prenoms,
      },
    };
  }

  // 2. Vérification dans le cache local (mode hors-ligne ou session mémorisée)
  const activeSession = getActiveSessionFromStorage();
  if (
    activeSession &&
    profileData?.email &&
    activeSession.user?.email &&
    activeSession.user.email.toLowerCase() === profileData.email.toLowerCase()
  ) {
    return {
      status: "authenticated",
      session: activeSession,
      provider,
      profile: {
        email: activeSession.user.email,
        nom: activeSession.user.nom,
        prenoms: activeSession.user.prenom,
      },
    };
  }

  // 3. Vérification dans l'annuaire local des comptes créés
  const knownAccounts = getKnownAccountsFromStorage();
  const matched = profileData?.email
    ? knownAccounts.find(
        (acc) => acc.email?.toLowerCase() === profileData.email!.toLowerCase(),
      )
    : null;

  if (matched) {
    saveSessionToStorage(matched.session);
    return {
      status: "authenticated",
      session: matched.session,
      provider,
      profile: {
        email: matched.session.user.email,
        nom: matched.session.user.nom,
        prenoms: matched.session.user.prenom,
      },
    };
  }

  // 4. Compte non trouvé -> Inscription requise
  return {
    status: "needs_registration",
    provider,
    profile: {
      email: profileData?.email || "",
      nom: profileData?.nom || "",
      prenoms: profileData?.prenoms || "",
    },
    message: `Aucun compte associé à ce profil ${providerName}. Veuillez renseigner le nom de votre établissement pour créer votre auto-école.`,
  };
}

/** Ouvre une URL dans le navigateur externe par défaut du système */
export async function openExternalBrowserUrl(url: string): Promise<void> {
  if (typeof window === "undefined") return;

  const isTauri = "__TAURI_INTERNALS__" in window || "__TAURI__" in window;
  if (isTauri) {
    const { invoke } = await import("@tauri-apps/api/core");
    await invoke("ouvrir_url_externe", { url });
    return;
  }

  window.open(url, "_blank");
}

/** Interroge l'état d'une session d'authentification Desktop en cours */
export async function pollOAuthDesktopSession(sessionId: string): Promise<{
  status: "pending" | "authenticated" | "needs_registration" | "error";
  token?: string;
  refresh_token?: string;
  role?: string;
  provider?: string;
  email?: string;
  nom?: string;
  prenoms?: string;
  error?: string;
}> {
  const backendOrigin = getBackendOrigin();

  try {
    const res = await fetch(
      `${backendOrigin}/api/auth/oauth/session?session_id=${encodeURIComponent(sessionId)}`,
      { cache: "no-store" },
    );
    if (!res.ok) return { status: "pending" };
    const data = await res.json();
    if (data.success && data.data) {
      return data.data;
    }
    return { status: "pending" };
  } catch {
    return { status: "pending" };
  }
}

/** Déclenche l'authentification OAuth Desktop en ouvrant le navigateur système */
export async function triggerOAuthDesktop(
  provider: "google" | "facebook" | "linkedin" | "apple",
  sessionId: string,
): Promise<void> {
  const backendOrigin = getBackendOrigin();
  const origin =
  typeof window !== "undefined" ? window.location.origin : backendOrigin;
  const oauthUrl = `${backendOrigin}/api/auth/${provider}?session_id=${encodeURIComponent(sessionId)}&origin=${encodeURIComponent(origin)}`;

  await openExternalBrowserUrl(oauthUrl);
}

/** Déclenche le flux d'authentification OAuth 2.0 dans le navigateur Web */
export function triggerOAuthWebFlow(
  provider: "google" | "facebook" | "linkedin" | "apple",
): void {
  if (typeof window === "undefined") return;
  const origin = window.location.origin;
  const backendOrigin = getBackendOrigin();
  const oauthUrl = `${backendOrigin}/api/auth/${provider}?origin=${encodeURIComponent(origin)}`;

  const width = 520;
  const height = 640;
  const left = window.screenX + (window.outerWidth - width) / 2;
  const top = window.screenY + (window.outerHeight - height) / 2;

  const popup = window.open(
    oauthUrl,
    `oauth_${provider}`,
    `width=${width},height=${height},left=${left},top=${top},status=no,resizable=yes,scrollbars=yes`,
  );

  if (!popup || popup.closed || typeof popup.closed === "undefined") {
    window.location.href = oauthUrl;
  }
}

/** Restaure et valide une session utilisateur à partir d'un jeton JWT OAuth reçu */
export async function restoreSessionFromJwt(
  token: string,
): Promise<AuthSessionData | null> {
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
      u.role || result.data.role || "DIRECTEUR",
    );

    const userProfile: UserProfile = {
      id: u.id,
      nom: u.nom,
      prenom: u.prenoms,
      role,
      roleLabel,
      initial: (u.prenoms?.[0] || u.nom?.[0] || "U").toUpperCase(),
      agence: a?.nom || "Agence Principale",
      auto_ecole_id: u.auto_ecole_id || s?.id,
      agence_id: u.agence_id || a?.id || "",
      telephone: u.telephone,
      email: u.email || undefined,
      code_pin_manager: u.code_pin_manager || undefined,
    };

    const school: AutoEcole = s
      ? {
          ...s,
          ville: a?.ville || "Douala",
          links: buildSchoolAccessLinks(s.code_unique),
        }
      : {
          id: u.auto_ecole_id || "school-default",
          code_unique: "elevendrive",
          nom_commercial: "Auto-École",
          devise: "FCFA",
          pays: "Cameroun",
          ville: "Douala",
          telephone: u.telephone || "",
          est_actif: true,
          created_at: new Date().toISOString(),
          links: buildSchoolAccessLinks("elevendrive"),
        };

    const session: AuthSessionData = {
      user: userProfile,
      auto_ecole: school,
      agence: a,
      token,
    };

    saveSessionToStorage(session);
    return session;
  } catch (e) {
    console.error("Erreur restoreSessionFromToken:", e);
    return null;
  }
}

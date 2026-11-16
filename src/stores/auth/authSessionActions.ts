// ==============================================================================
// AutoSchool - Actions Authentification, OAuth & Inscription (Auth Store)
// ==============================================================================

import type {
  LoginCredentials,
  OAuthResult,
  RegisterSchoolPayload,
  UserRole,
} from "../../types";
import { authBackend } from "../../services/authBackendService";
import { useNotificationStore } from "../notificationStore";
import {
  applySessionToStore,
  startRealtimeNotifications,
} from "./authEnvSessionActions";

export function execInvalidateLocalSession(store: any): void {
  if (typeof localStorage !== "undefined") {
    localStorage.removeItem("autoschool_web_session");
    localStorage.removeItem("autoschool_active_session");
    localStorage.removeItem("elevendrive_active_session");
    localStorage.removeItem("elevendrive_refresh_token");
  }
  store.isAuthenticated = false;
  store.customUserProfile = null;
  store.authError =
    "Votre session n'est plus valide. Veuillez vous reconnecter.";
}

export async function execLoginSuperAdmin(
  store: any,
  identifiant: string,
  motDePasse: string,
) {
  store.isAuthLoading = true;
  store.authError = null;
  try {
    const session = await authBackend.loginSuperAdmin(identifiant, motDePasse);
    store.isAuthenticated = true;
    store.currentRole = "super_admin";
    store.customUserProfile = session.user;
    store.currentSchool = session.auto_ecole;
    store.currentAgence = session.agence;
    store.webPortalView = "super_admin_dashboard";
    return session;
  } catch (err: any) {
    store.authError =
      err?.message || "Identifiants Super-Administrateur invalides.";
    throw err;
  } finally {
    store.isAuthLoading = false;
  }
}

export async function execLogin(
  store: any,
  credentials: LoginCredentials,
  schoolSlug?: string,
) {
  store.isAuthLoading = true;
  store.authError = null;
  try {
    const slug =
      schoolSlug ||
      (store.hasExplicitSchoolInUrl ? store.schoolSlug : undefined) ||
      undefined;
    const session = await authBackend.loginDesktop(credentials, slug);
    applySessionToStore(store, session, true);
    startRealtimeNotifications(session.token);
    return session;
  } catch (err: any) {
    const offlineSession = authBackend.authenticateOffline(
      credentials.identifiant,
    );
    if (offlineSession) {
      authBackend.saveSession(offlineSession);
      applySessionToStore(store, offlineSession);
      return offlineSession;
    }

    store.authError = err?.message || "Erreur lors de la connexion.";
    throw err;
  } finally {
    store.isAuthLoading = false;
  }
}

export async function execApplyOAuthSession(
  store: any,
  accessToken: string,
  refreshToken?: string,
  _preferredRole?: UserRole,
) {
  store.isAuthLoading = true;
  store.authError = null;
  try {
    if (refreshToken && typeof localStorage !== "undefined") {
      localStorage.setItem("elevendrive_refresh_token", refreshToken);
    }

    const session = await authBackend.restoreSessionFromToken(accessToken);
    if (session) {
      applySessionToStore(store, session, true);
      store.loadAgencesLocally();
      return session;
    }
  } catch (err: any) {
    store.authError =
      err?.message || "Échec de validation de la session OAuth.";
    throw err;
  } finally {
    store.isAuthLoading = false;
  }
}

export async function execLoginWithOAuth(
  store: any,
  provider: "google" | "facebook" | "linkedin" | "apple",
  profileData?: {
    email?: string;
    nom?: string;
    prenoms?: string;
    id_token?: string;
  },
): Promise<OAuthResult> {
  store.isAuthLoading = true;
  store.authError = null;
  try {
    const res = await authBackend.loginWithOAuth(provider, profileData);
    if (res.status === "authenticated" && res.session) {
      applySessionToStore(store, res.session, true);
    }
    return res;
  } catch (err: any) {
    store.authError =
      err?.message || `Erreur lors de la connexion avec ${provider}.`;
    throw err;
  } finally {
    store.isAuthLoading = false;
  }
}

export async function execRegisterSchool(
  store: any,
  payload: RegisterSchoolPayload,
) {
  store.isAuthLoading = true;
  store.authError = null;
  try {
    const session = await authBackend.registerSchool(payload);
    store.registrationSuccessData = session;
    applySessionToStore(store, session);
    startRealtimeNotifications(session.token);
    return session;
  } catch (err: any) {
    store.authError =
      err?.message || "Erreur lors de l'enregistrement de l'auto-école.";
    throw err;
  } finally {
    store.isAuthLoading = false;
  }
}

export function execLogout(store: any): void {
  authBackend.clearSession();
  if (typeof localStorage !== "undefined") {
    localStorage.removeItem("autoschool_web_session");
    localStorage.removeItem("autoschool_active_session");
  }
  store.isAuthenticated = false;
  store.customUserProfile = null;
  store.registrationSuccessData = null;
  store.authError = null;
  store.webPortalView = "home";
  if (store.platform === "desktop") {
    store.currentRole = "secretaire";
  }
  try {
    const notifStore = useNotificationStore();
    notifStore.cleanup();
  } catch {
    // ignore
  }
}


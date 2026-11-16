// ==============================================================================
// AutoSchool - Service Client d'Authentification (Façade Unifiée)
// ==============================================================================

import type {
  AutoEcole,
  Agence,
  RegisterSchoolPayload,
  RegisterUserPayload,
  UserTeamMember,
  LoginCredentials,
  AuthSessionData,
  OAuthResult,
} from "../types";
import {
  API_BASE_URL,
  STORAGE_KEY_REFRESH_TOKEN,
  authenticateOfflineFromStorage,
  buildSchoolAccessLinks,
  clearSessionFromStorage,
  fetchAuthWithFallback,
  generateSchoolSlug,
  getActiveSessionFromStorage,
  getKnownAccountsFromStorage,
  saveSessionToStorage,
} from "./auth/authSessionHelpers";
import {
  loginAsSuperAdmin,
  loginWithCredentials,
  registerNewSchool,
} from "./auth/authCredentialsClient";
import {
  loginWithOAuthProvider,
  openExternalBrowserUrl,
  pollOAuthDesktopSession,
  restoreSessionFromJwt,
  triggerOAuthDesktop,
  triggerOAuthWebFlow,
} from "./auth/authOAuthClient";
import {
  createSchoolAgence,
  deleteSchoolAgence,
  deleteSchoolUser,
  fetchMyProfile,
  fetchSchoolAgences,
  fetchSchoolUsers,
  registerSchoolUser,
  updateMyUserProfile,
  updateSchoolAgence,
  updateSchoolDetails,
  updateSchoolUser,
} from "./auth/authTeamAgenceClient";

export { generateSchoolSlug, buildSchoolAccessLinks };

class AuthBackendService {
  /** Vérifie auprès du serveur que le compte de la session locale existe toujours */
  async verifySession(token?: string): Promise<boolean> {
    if (!token) return true;
    try {
      const res = await fetch(`${API_BASE_URL.replace(/\/auth\/?$/, "")}/me`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      return !(res.status === 401 || res.status === 403 || res.status === 404);
    } catch {
      return true;
    }
  }

  authenticateOffline(identifiant: string): AuthSessionData | null {
    return authenticateOfflineFromStorage(identifiant);
  }

  loadStoredSession(): AuthSessionData | null {
    return this.getActiveSession();
  }

  registerSchool(payload: RegisterSchoolPayload): Promise<AuthSessionData> {
    return registerNewSchool(payload);
  }

  login(
    credentials: LoginCredentials,
    schoolSlug?: string,
  ): Promise<AuthSessionData> {
    return loginWithCredentials(credentials, schoolSlug);
  }

  loginSuperAdmin(
    identifiant: string,
    motDePasse: string,
  ): Promise<AuthSessionData> {
    return loginAsSuperAdmin(identifiant, motDePasse);
  }

  loginDesktop(
    credentials: LoginCredentials,
    schoolSlug?: string,
  ): Promise<AuthSessionData> {
    return this.login(credentials, schoolSlug);
  }

  loginWithOAuth(
    provider: "google" | "facebook" | "linkedin" | "apple",
    profileData?: {
      email?: string;
      nom?: string;
      prenoms?: string;
      id_token?: string;
    },
  ): Promise<OAuthResult> {
    return loginWithOAuthProvider(provider, profileData);
  }

  openExternalBrowser(url: string) {
    return openExternalBrowserUrl(url);
  }

  pollOAuthSession(sessionId: string) {
    return pollOAuthDesktopSession(sessionId);
  }

  startOAuthDesktop(
    provider: "google" | "facebook" | "linkedin" | "apple",
    sessionId: string,
  ) {
    return triggerOAuthDesktop(provider, sessionId);
  }

  startOAuthFlow(provider: "google" | "facebook" | "linkedin" | "apple") {
    return triggerOAuthWebFlow(provider);
  }

  restoreSessionFromToken(token: string): Promise<AuthSessionData | null> {
    return restoreSessionFromJwt(token);
  }

  /** Récupérer les informations publiques d'une Auto-École via son slug */
  async getSchoolByCode(code: string): Promise<AutoEcole | null> {
    try {
      const response = await fetch(
        `${API_BASE_URL}/school?slug=${encodeURIComponent(code)}`,
      );
      const result = await response.json();
      if (!response.ok || !result.success) return null;

      const data = result.data;
      return {
        id: data.id,
        code_unique: data.code_unique,
        nom_commercial: data.nom_commercial,
        numero_agrement: data.numero_agrement || "AGR-MIN-AUTO/2024",
        telephone: data.telephone,
        email: data.email,
        devise: data.devise || "FCFA",
        pays: data.pays || "Cameroun",
        ville: data.agences?.[0]?.ville || "Douala",
        logo_url: data.logo_url,
        est_actif: true,
        links: buildSchoolAccessLinks(data.code_unique),
        created_at: new Date().toISOString(),
      };
    } catch {
      return null;
    }
  }

  /** Vérification d'un code PIN manager */
  async verifyManagerPin(schoolId: string, pin: string): Promise<boolean> {
    try {
      const response = await fetch(`${API_BASE_URL}/verify-pin`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ auto_ecole_id: schoolId, code_pin: pin }),
      });
      const result = await response.json();
      return response.ok && result.success;
    } catch {
      return false;
    }
  }

  /** Déconnexion et révocation de session */
  async logout(): Promise<void> {
    if (typeof localStorage !== "undefined") {
      const refreshToken = localStorage.getItem(STORAGE_KEY_REFRESH_TOKEN);
      if (refreshToken) {
        try {
          await fetch(`${API_BASE_URL}/logout`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ refresh_token: refreshToken }),
          });
        } catch {
          // Ignorer l'erreur réseau lors de la déconnexion
        }
      }
      localStorage.removeItem(STORAGE_KEY_REFRESH_TOKEN);
    }
    this.clearSession();
  }

  /** Renouvelle l'Access Token JWT via le Refresh Token persistant */
  async refreshAccessToken(): Promise<string | null> {
    if (typeof localStorage === "undefined") return null;
    const refreshToken = localStorage.getItem(STORAGE_KEY_REFRESH_TOKEN);
    if (!refreshToken) return null;

    try {
      const response = await fetchAuthWithFallback("/refresh", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ refresh_token: refreshToken }),
      });
      if (!response.ok) return null;

      const result = await response.json();
      if (!result.success || !result.data?.access_token) return null;

      const newAccessToken = result.data.access_token as string;
      if (result.data.refresh_token) {
        localStorage.setItem(
          STORAGE_KEY_REFRESH_TOKEN,
          result.data.refresh_token,
        );
      }

      const activeSession = this.getActiveSession();
      if (activeSession) {
        activeSession.token = newAccessToken;
        this.saveSession(activeSession);
      }
      return newAccessToken;
    } catch {
      return null;
    }
  }

  saveSession(session: AuthSessionData) {
    saveSessionToStorage(session);
  }

  getKnownAccounts(): any[] {
    return getKnownAccountsFromStorage();
  }

  getActiveSession(): AuthSessionData | null {
    return getActiveSessionFromStorage();
  }

  clearSession() {
    clearSessionFromStorage();
  }

  getSchoolUsers(): Promise<UserTeamMember[]> {
    return fetchSchoolUsers();
  }

  getSchoolAgences(): Promise<Agence[]> {
    return fetchSchoolAgences();
  }

  createAgence(payload: {
    nom: string;
    code_agence?: string;
    adresse?: string;
    ville?: string;
    commune?: string;
    telephone?: string;
    est_siege?: boolean;
  }): Promise<Agence> {
    return createSchoolAgence(payload);
  }

  updateAgence(id: string, updates: Partial<Agence>): Promise<Agence> {
    return updateSchoolAgence(id, updates);
  }

  deleteAgence(id: string): Promise<void> {
    return deleteSchoolAgence(id);
  }

  registerUser(payload: RegisterUserPayload): Promise<any> {
    return registerSchoolUser(payload);
  }

  updateUser(
    userId: string,
    payload: Partial<RegisterUserPayload> & { statut?: string },
  ): Promise<any> {
    return updateSchoolUser(userId, payload);
  }

  deleteUser(userId: string): Promise<boolean> {
    return deleteSchoolUser(userId);
  }

  getMyProfile() {
    return fetchMyProfile();
  }

  updateMyProfile(payload: {
    nom: string;
    prenoms: string;
    telephone: string;
    email?: string;
    mot_de_passe?: string;
    code_pin_manager?: string;
  }): Promise<any> {
    return updateMyUserProfile(payload);
  }

  updateSchoolInfo(payload: Partial<AutoEcole>): Promise<AutoEcole | null> {
    return updateSchoolDetails(payload);
  }
}

export const authBackend = new AuthBackendService();

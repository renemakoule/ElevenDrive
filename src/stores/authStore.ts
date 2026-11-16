import { defineStore } from "pinia";
import type {
  UserRole,
  AppPlatform,
  UserProfile,
  AutoEcole,
  Agence,
  SchoolAccessLinks,
  LoginCredentials,
  RegisterSchoolPayload,
  AuthSessionData,
  OAuthResult,
} from "../types";
import { authBackend } from "../services/authBackendService";
import {
  USER_PROFILES,
  computeIsLicenceExpired,
  computeIsModuleEnabled,
  computeIsRoleAllowed,
  computeRolePermissions,
  computeTrialDaysRemaining,
  computeWelcomeSubtitle,
  parseActiveModules,
} from "./auth/authGetters";
import {
  applyThemeToDom,
  runInitFromEnvironment,
} from "./auth/authEnvSessionActions";
import {
  execApplyOAuthSession,
  execInvalidateLocalSession,
  execLogin,
  execLoginSuperAdmin,
  execLoginWithOAuth,
  execLogout,
  execRegisterSchool,
} from "./auth/authSessionActions";
import {
  execAddAgence,
  execDeleteAgence,
  execHandleRealtimeSchoolEvent,
  execLoadAgences,
  execLoadAgencesLocally,
  execLoadProfileByRole,
  execSaveAgencesLocally,
  execUpdateAgence,
  execUpdateSchool,
  execUpdateUserProfile,
} from "./auth/authProfileAgenceActions";

export { USER_PROFILES };

export const useAuthStore = defineStore("auth", {
  state: () => ({
    isAuthenticated: false,
    currentRole: "secretaire" as UserRole,
    platform: "desktop" as AppPlatform,
    activeSection: "overview" as string,
    isSidebarCollapsed: false,
    isHelpModalOpen: false,
    isDark: false,
    customUserProfile: null as UserProfile | null,
    currentSchool: null as AutoEcole | null,
    currentAgence: null as Agence | null,
    agencesList: [] as Agence[],
    schoolAccessLinks: null as SchoolAccessLinks | null,
    authError: null as string | null,
    isAuthLoading: false,
    registrationSuccessData: null as AuthSessionData | null,
    webPortalView: "landing" as
      | "landing"
      | "home"
      | "auth"
      | "register"
      | "desktop_auth"
      | "super_admin_login"
      | "super_admin_dashboard",
    hasExplicitSchoolInUrl: false,
    pendingOAuthRegistration: null as {
      provider: "google" | "facebook" | "linkedin" | "apple";
      email: string;
      nom: string;
      prenoms: string;
    } | null,
  }),

  getters: {
    currentUser(state): UserProfile {
      if (
        state.customUserProfile &&
        state.customUserProfile.role === state.currentRole
      ) {
        return state.customUserProfile;
      }
      return USER_PROFILES[state.currentRole];
    },

    schoolName(state): string {
      return state.currentSchool?.nom_commercial || "ElevenDrive";
    },

    schoolSlug(state): string {
      return state.currentSchool?.code_unique || "";
    },

    schoolAgrement(state): string {
      return state.currentSchool?.numero_agrement || "Agrément en cours";
    },

    schoolPhone(state): string {
      return (
        state.currentSchool?.telephone ||
        state.customUserProfile?.telephone ||
        ""
      );
    },

    schoolAddress(state): string {
      if (state.currentAgence) {
        return `${state.currentAgence.nom}, ${state.currentAgence.ville || ""}`;
      }
      return state.currentSchool
        ? `${state.currentSchool.ville || ""}, ${state.currentSchool.pays || ""}`.trim()
        : "";
    },

    isCurrentRoleAllowed(state): boolean {
      return computeIsRoleAllowed(state.currentRole, state.platform);
    },

    canSetInitialPrice(state): boolean {
      return (
        state.currentRole === "manager" || state.currentRole === "secretaire"
      );
    },

    canModifyFixedPrice(state): boolean {
      return state.currentRole === "manager";
    },

    rolePermissions(state) {
      return computeRolePermissions(state.currentRole);
    },

    isLicenceExpired(state): boolean {
      return computeIsLicenceExpired(state.currentRole, state.currentSchool);
    },

    trialDaysRemaining(state): number {
      return computeTrialDaysRemaining(state.currentSchool);
    },

    activeModulesList(state): string[] {
      return parseActiveModules(state.currentRole, state.currentSchool);
    },

    isModuleEnabled(state): (moduleKey: string) => boolean {
      return (moduleKey: string): boolean =>
        computeIsModuleEnabled(
          state.currentRole,
          state.currentSchool,
          moduleKey,
        );
    },

    greetingMessage(state): string {
      const user =
        state.customUserProfile &&
        state.customUserProfile.role === state.currentRole
          ? state.customUserProfile
          : USER_PROFILES[state.currentRole];
      return `Hi, ${user.prenom} 👋`;
    },

    welcomeSubtitle(state): string {
      return computeWelcomeSubtitle(
        state.currentRole,
        state.currentSchool?.nom_commercial,
      );
    },
  },

  actions: {
    invalidateLocalSession() {
      execInvalidateLocalSession(this);
    },

    initFromEnvironment() {
      return runInitFromEnvironment(this);
    },

    setWebPortalView(
      view:
        | "landing"
        | "home"
        | "auth"
        | "register"
        | "desktop_auth"
        | "super_admin_login"
        | "super_admin_dashboard",
    ) {
      this.webPortalView = view;
    },

    loginWebAsCandidate(candidateData: {
      id: string;
      nom: string;
      prenom: string;
      telephone: string;
      dossier: string;
    }) {
      this.isAuthenticated = true;
      this.currentRole = "eleve";
      this.platform = "pwa";
      this.customUserProfile = {
        id: candidateData.id,
        nom: candidateData.nom,
        prenom: candidateData.prenom,
        role: "eleve",
        roleLabel: "Candidat",
        initial: (
          candidateData.prenom?.[0] ||
          candidateData.nom?.[0] ||
          "C"
        ).toUpperCase(),
        agence: this.schoolName,
        telephone: candidateData.telephone,
        dossier: candidateData.dossier,
      };

      if (typeof localStorage !== "undefined") {
        localStorage.setItem(
          "autoschool_web_session",
          JSON.stringify({
            user: this.customUserProfile,
            role: "eleve",
            candidatId: candidateData.id,
            schoolSlug: this.schoolSlug,
          }),
        );
      }
    },

    loginWebAsMoniteur(moniteurData: {
      id: string;
      nom: string;
      prenom: string;
      telephone: string;
    }) {
      this.isAuthenticated = true;
      this.currentRole = "moniteur";
      this.platform = "pwa";
      this.customUserProfile = {
        id: moniteurData.id,
        nom: moniteurData.nom,
        prenom: moniteurData.prenom,
        role: "moniteur",
        roleLabel: "Moniteur",
        initial: (moniteurData.prenom?.[0] || "M").toUpperCase(),
        agence: this.schoolName,
        telephone: moniteurData.telephone,
      };

      if (typeof localStorage !== "undefined") {
        localStorage.setItem(
          "autoschool_web_session",
          JSON.stringify({
            user: this.customUserProfile,
            role: "moniteur",
            schoolSlug: this.schoolSlug,
          }),
        );
      }
    },

    logoutWeb() {
      if (typeof localStorage !== "undefined") {
        localStorage.removeItem("autoschool_web_session");
      }
      this.isAuthenticated = false;
      this.webPortalView = "home";
      this.customUserProfile = null;
    },

    loginSuperAdmin(identifiant: string, motDePasse: string) {
      return execLoginSuperAdmin(this, identifiant, motDePasse);
    },

    login(credentials: LoginCredentials, schoolSlug?: string) {
      return execLogin(this, credentials, schoolSlug);
    },

    startOAuthFlow(provider: "google" | "facebook" | "linkedin" | "apple") {
      authBackend.startOAuthFlow(provider);
    },

    applyOAuthSession(
      accessToken: string,
      refreshToken?: string,
      preferredRole?: UserRole,
    ) {
      return execApplyOAuthSession(
        this,
        accessToken,
        refreshToken,
        preferredRole,
      );
    },

    loginWithOAuth(
      provider: "google" | "facebook" | "linkedin" | "apple",
      profileData?: {
        email?: string;
        nom?: string;
        prenoms?: string;
        id_token?: string;
      },
    ): Promise<OAuthResult> {
      return execLoginWithOAuth(this, provider, profileData);
    },

    registerSchool(payload: RegisterSchoolPayload) {
      return execRegisterSchool(this, payload);
    },

    logout() {
      execLogout(this);
    },

    applyTheme(isDark: boolean) {
      applyThemeToDom(isDark);
    },

    toggleDarkMode() {
      this.isDark = !this.isDark;
      this.applyTheme(this.isDark);
      if (typeof localStorage !== "undefined") {
        localStorage.setItem(
          "autoschool_theme",
          this.isDark ? "dark" : "light",
        );
      }
    },

    setRole(role: UserRole) {
      if (this.isAuthenticated && this.customUserProfile) {
        console.warn(
          `[Sécurité] Tentative non autorisée de changement de rôle vers ${role}. Session scellée pour ${this.customUserProfile.role}.`,
        );
        return;
      }

      this.currentRole = role;
      this.activeSection = "overview";

      if (role === "secretaire") {
        this.platform = "desktop";
      } else if (role === "moniteur" || role === "eleve") {
        this.platform = "pwa";
      }
    },

    setPlatform(platform: AppPlatform) {
      this.platform = platform;
    },

    setActiveSection(section: string) {
      this.activeSection = section;
    },

    toggleSidebar() {
      this.isSidebarCollapsed = !this.isSidebarCollapsed;
    },

    toggleHelp() {
      this.isHelpModalOpen = !this.isHelpModalOpen;
    },

    loadProfileByRole() {
      return execLoadProfileByRole(this);
    },

    updateUserProfile(
      updates: Partial<UserProfile> & { nouveau_mot_de_passe?: string },
    ) {
      return execUpdateUserProfile(this, updates);
    },

    updateSchool(updates: Partial<AutoEcole>) {
      return execUpdateSchool(this, updates);
    },

    loadAgences() {
      return execLoadAgences(this);
    },

    updateAgence(agenceId: string, updates: Partial<Agence>) {
      return execUpdateAgence(this, agenceId, updates);
    },

    addAgence(newAgence: Omit<Agence, "id">) {
      return execAddAgence(this, newAgence);
    },

    deleteAgence(agenceId: string) {
      return execDeleteAgence(this, agenceId);
    },

    saveAgencesLocally() {
      execSaveAgencesLocally(this);
    },

    loadAgencesLocally() {
      execLoadAgencesLocally(this);
    },

    handleRealtimeSchoolEvent(event: { type: string; data: any }) {
      execHandleRealtimeSchoolEvent(this, event);
    },
  },
});

// ==============================================================================
// AutoSchool - Initialisation Environnement, Thème & Session (Auth Store)
// ==============================================================================

import type { AuthSessionData, UserRole } from "../../types";
import {
  authBackend,
  buildSchoolAccessLinks,
  generateSchoolSlug,
} from "../../services/authBackendService";
import { useNotificationStore } from "../notificationStore";
import { USER_PROFILES } from "./authGetters";

export function applyThemeToDom(isDark: boolean): void {
  if (typeof document === "undefined") return;
  const root = document.documentElement;
  if (isDark) {
    root.classList.add("dark");
    root.setAttribute("data-theme", "dark");
    if (document.body) {
      document.body.classList.add("dark");
      document.body.setAttribute("data-theme", "dark");
    }
  } else {
    root.classList.remove("dark");
    root.removeAttribute("data-theme");
    if (document.body) {
      document.body.classList.remove("dark");
      document.body.removeAttribute("data-theme");
    }
  }
}

export function applySessionToStore(
  store: any,
  session: AuthSessionData,
  adaptPlatform = false,
): void {
  store.isAuthenticated = true;
  store.currentRole = session.user.role;
  store.customUserProfile = session.user;
  store.currentSchool = session.auto_ecole;
  store.currentAgence = session.agence;
  store.schoolAccessLinks =
    session.auto_ecole.links ||
    buildSchoolAccessLinks(session.auto_ecole.code_unique);

  if (
    adaptPlatform &&
    (session.user.role === "moniteur" || session.user.role === "eleve")
  ) {
    store.platform = "pwa";
  }
}

export function startRealtimeNotifications(token?: string): void {
  if (!token) return;
  const notifStore = useNotificationStore();
  notifStore.initRealtime(token);
  notifStore.loadNotifications();
}

export async function runInitFromEnvironment(store: any): Promise<void> {
  if (typeof window === "undefined") return;

  const urlParams = new URLSearchParams(window.location.search);

  // 1. Thème sombre / clair initialisation
  const savedTheme = localStorage.getItem("autoschool_theme");
  const prefersDark =
    window.matchMedia &&
    window.matchMedia("(prefers-color-scheme: dark)").matches;
  store.isDark = savedTheme === "dark" || (!savedTheme && prefersDark);
  store.applyTheme(store.isDark);

  // Retours OAuth 2.0 direct
  const oauthStatus = urlParams.get("oauth_status");
  if (oauthStatus === "authenticated") {
    const token = urlParams.get("token");
    const refreshToken = urlParams.get("refresh_token") || undefined;
    const role = (urlParams.get("role") as UserRole) || undefined;
    if (token) {
      await store.applyOAuthSession(token, refreshToken, role);
      const cleanUrl = window.location.origin + window.location.pathname;
      window.history.replaceState({}, document.title, cleanUrl);
      return;
    }
  } else if (oauthStatus === "needs_registration") {
    const rawProvider = urlParams.get("provider");
    const provider: "google" | "facebook" | "linkedin" | "apple" =
      rawProvider === "facebook" ||
      rawProvider === "linkedin" ||
      rawProvider === "apple"
        ? rawProvider
        : "google";

    store.pendingOAuthRegistration = {
      provider,
      email: urlParams.get("email") || "",
      nom: urlParams.get("nom") || "",
      prenoms: urlParams.get("prenoms") || "",
    };
    store.webPortalView = "desktop_auth";
    const cleanUrl = window.location.origin + window.location.pathname;
    window.history.replaceState({}, document.title, cleanUrl);
  } else if (oauthStatus === "error") {
    store.authError =
      urlParams.get("oauth_error") ||
      "Erreur lors de l'authentification OAuth.";
    store.webPortalView = "desktop_auth";
    const cleanUrl = window.location.origin + window.location.pathname;
    window.history.replaceState({}, document.title, cleanUrl);
  }

  // 2. Détection Desktop Tauri vs Web standard
  const isTauri = "__TAURI_INTERNALS__" in window || "__TAURI__" in window;
  if (urlParams.get("platform") === "desktop") {
    store.platform = "desktop";
  } else if (
    urlParams.get("platform") === "pwa" ||
    urlParams.get("platform") === "web"
  ) {
    store.platform = "pwa";
  } else {
    store.platform = isTauri ? "desktop" : "pwa";
  }

  // Routes dédiées Super-Admin
  const currentPath = window.location.pathname
    .toLowerCase()
    .replace(/\/+$/, "");
  if (currentPath === "/super-admin/dashboard") {
    const activeSession = authBackend.getActiveSession();
    if (activeSession && activeSession.user.role === "super_admin") {
      store.isAuthenticated = true;
      store.currentRole = "super_admin";
      store.customUserProfile = activeSession.user;
      store.currentSchool = activeSession.auto_ecole;
      store.webPortalView = "super_admin_dashboard";
      return;
    }
    store.isAuthenticated = false;
    store.webPortalView = "super_admin_login";
    return;
  }
  if (currentPath === "/super-admin/login" || currentPath === "/super-admin") {
    store.isAuthenticated = false;
    store.webPortalView = "super_admin_login";
    return;
  }

  // 3. Détection du code auto-école dans l'URL
  const explicitSchoolParam = urlParams.get("school") || urlParams.get("slug");
  let pathSegment: string | null = null;
  if (
    window.location.pathname &&
    window.location.pathname !== "/" &&
    window.location.pathname !== "/index.html"
  ) {
    const seg = window.location.pathname
      .replace(/^\/+|\/+$/g, "")
      .split("/")[0];
    if (
      seg &&
      !["api", "login", "register", "auth", "admin", "pwa", "index.html", "index"].includes(
        seg.toLowerCase(),
      )
    ) {
      pathSegment = seg;
    }
  }

  // Nettoyage visuel si le serveur statique a redirigé vers /index.html
  if (window.location.pathname === "/index.html") {
    const cleanOrigin = window.location.origin + (window.location.search || "");
    window.history.replaceState({}, document.title, cleanOrigin);
  }
  store.hasExplicitSchoolInUrl = Boolean(explicitSchoolParam || pathSegment);

  let schoolCode = explicitSchoolParam || pathSegment;
  if (!schoolCode && typeof localStorage !== "undefined") {
    schoolCode = localStorage.getItem("autoschool_last_slug");
  }
  if (!schoolCode) {
    const activeSession = authBackend.getActiveSession();
    if (activeSession?.auto_ecole?.code_unique) {
      schoolCode = activeSession.auto_ecole.code_unique;
    }
  }

  if (schoolCode) {
    if (typeof localStorage !== "undefined") {
      localStorage.setItem("autoschool_last_slug", schoolCode);
    }
    const resolved = await authBackend.getSchoolByCode(schoolCode);
    if (resolved) {
      store.currentSchool = resolved;
      store.schoolAccessLinks =
        resolved.links || buildSchoolAccessLinks(resolved.code_unique);
    } else {
      const formattedName = schoolCode
        .split("-")
        .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
        .join(" ");
      store.currentSchool = {
        id: `school-${schoolCode}`,
        code_unique: schoolCode,
        nom_commercial: formattedName,
        numero_agrement: "AGR-MIN-AUTO/2024",
        telephone: "+237 651 727 932",
        email: `contact@${schoolCode}.cm`,
        devise: "FCFA",
        pays: "Cameroun",
        ville: "Douala",
        est_actif: true,
        links: buildSchoolAccessLinks(schoolCode),
        created_at: new Date().toISOString(),
      };
      store.schoolAccessLinks = buildSchoolAccessLinks(schoolCode);
    }
  } else {
    store.currentSchool = {
      id: "school-default",
      code_unique: "elevendrive",
      nom_commercial: "ElevenDrive",
      numero_agrement: "AGR-MIN-AUTO/2026",
      telephone: "+237 651 727 932",
      email: "contact@elevendrive.cm",
      devise: "FCFA",
      pays: "Cameroun",
      ville: "Douala",
      est_actif: true,
      links: buildSchoolAccessLinks("elevendrive"),
      created_at: new Date().toISOString(),
    };
    store.schoolAccessLinks = buildSchoolAccessLinks("elevendrive");
  }

  const mode = urlParams.get("mode");
  const role = urlParams.get("role") as UserRole;
  const hasDossier = urlParams.has("dossier");

  if (role && USER_PROFILES[role]) {
    store.currentRole = role;
  }

  // 4. Restauration de la session selon la plateforme
  if (store.platform === "desktop") {
    const activeSession = authBackend.getActiveSession();
    if (activeSession) {
      if (activeSession.auto_ecole && !activeSession.auto_ecole.code_unique) {
        const dynamicSlug = generateSchoolSlug(
          activeSession.auto_ecole.nom_commercial || "auto-ecole",
        );
        activeSession.auto_ecole.code_unique = dynamicSlug;
        activeSession.auto_ecole.links = buildSchoolAccessLinks(dynamicSlug);
        authBackend.saveSession(activeSession);
      }

      applySessionToStore(store, activeSession);
      store.loadAgencesLocally();
      authBackend.verifySession(activeSession.token).then((ok) => {
        if (!ok) {
          store.invalidateLocalSession();
        } else {
          store.loadProfileByRole();
          startRealtimeNotifications(activeSession.token);
        }
      });
    } else {
      store.isAuthenticated = false;
    }
    return;
  }

  const savedWebSession = localStorage.getItem("autoschool_web_session");
  if (savedWebSession) {
    try {
      const parsed = JSON.parse(savedWebSession);
      if (parsed.user && parsed.role) {
        store.customUserProfile = parsed.user;
        store.currentRole = parsed.role;
        store.isAuthenticated = true;
        if (parsed.token) {
          authBackend.verifySession(parsed.token).then((ok) => {
            if (!ok) {
              store.invalidateLocalSession();
            } else {
              store.loadProfileByRole();
              startRealtimeNotifications(parsed.token);
            }
          });
        }
        return;
      }
    } catch {}
  }

  if (mode === "candidat" && hasDossier) {
    store.currentRole = "eleve";
    store.isAuthenticated = true;
    return;
  }

  store.isAuthenticated = false;
  if (mode === "login" || urlParams.get("auth") === "login") {
    store.webPortalView = store.hasExplicitSchoolInUrl
      ? "auth"
      : "desktop_auth";
  } else if (mode === "register" || urlParams.get("action") === "nouveau") {
    store.webPortalView = "desktop_auth";
  } else if (store.hasExplicitSchoolInUrl) {
    store.webPortalView = "home";
  } else {
    store.webPortalView = "landing";
  }
}

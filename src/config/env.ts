// ==============================================================================
// ElevenDrive - Configuration centralisée des URLs depuis les variables .env
// ==============================================================================

const rawApiUrl = (
  (import.meta.env?.VITE_API_URL as string | undefined) || ""
).trim();

const isDevMode = Boolean(import.meta.env?.DEV);

function isTauriWebviewOrigin(origin: string): boolean {
  return (
    origin.includes("tauri.localhost") || origin.startsWith("tauri://")
  );
}

/**
 * Résout l'origine du serveur Go (ex: https://api.mondomaine.com).
 * Utilise VITE_API_URL en priorité ; n'utilise localhost qu'en mode développement.
 */
export function resolveBackendOrigin(): string {
  if (rawApiUrl) {
    try {
      const parsed = new URL(rawApiUrl);
      return `${parsed.protocol}//${parsed.host}`;
    } catch {
      return rawApiUrl.replace(/\/api(\/auth)?\/?$/i, "").replace(/\/+$/, "");
    }
  }

  if (isDevMode) {
    return "http://localhost:8080";
  }

  if (
    typeof window !== "undefined" &&
    window.location?.origin &&
    !isTauriWebviewOrigin(window.location.origin)
  ) {
    return window.location.origin;
  }

  if (typeof console !== "undefined") {
    console.error(
      "[ElevenDrive] Variable VITE_API_URL manquante dans votre fichier .env / .env.production.",
    );
  }
  return "";
}

export const BACKEND_ORIGIN = resolveBackendOrigin();
export const API_BASE_URL = BACKEND_ORIGIN ? `${BACKEND_ORIGIN}/api` : "/api";
export const AUTH_API_BASE_URL = BACKEND_ORIGIN
  ? `${BACKEND_ORIGIN}/api/auth`
  : "/api/auth";

/** URL WebSocket vers /api/ws dérivée de VITE_API_URL (ws:// ou wss://). */
export function getWsApiBaseUrl(): string {
  const origin =
    BACKEND_ORIGIN ||
    (typeof window !== "undefined" ? window.location.origin : "");
  if (!origin) return "";
  const wsOrigin = origin
    .replace(/^http:\/\//i, "ws://")
    .replace(/^https:\/\//i, "wss://");
  return `${wsOrigin}/api/ws`;
}

/** URL publique du portail Web / QR Code (VITE_PUBLIC_APP_URL ou origine courante). */
export function getPublicAppOrigin(): string {
  const envPublic = (
    (import.meta.env?.VITE_PUBLIC_APP_URL as string | undefined) || ""
  )
    .trim()
    .replace(/\/+$/, "");
  if (envPublic) return envPublic;

  if (
    typeof window !== "undefined" &&
    window.location?.origin &&
    !isTauriWebviewOrigin(window.location.origin)
  ) {
    return window.location.origin;
  }
  return BACKEND_ORIGIN;
}

export const R2_PUBLIC_URL = (
  (import.meta.env?.VITE_R2_PUBLIC_URL as string | undefined) || ""
)
  .trim()
  .replace(/\/+$/, "");

export const JITSI_BASE_URL =
  (
    (import.meta.env?.VITE_JITSI_BASE_URL as string | undefined) || ""
  )
    .trim()
    .replace(/\/+$/, "") || "https://meet.jit.si";


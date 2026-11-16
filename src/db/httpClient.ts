import { authBackend } from "../services/authBackendService";
import { API_BASE_URL, BACKEND_ORIGIN } from "../config/env";

export { API_BASE_URL };

export const isTauri =
  typeof window !== "undefined" && "__TAURI_INTERNALS__" in window;

export const getAuthHeaders = (): Record<string, string> => {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };
  if (typeof localStorage !== "undefined") {
    const raw =
      localStorage.getItem("autoschool_active_session") ||
      localStorage.getItem("elevendrive_active_session");
    if (raw) {
      try {
        const session = JSON.parse(raw);
        if (session.token) {
          headers["Authorization"] = `Bearer ${session.token}`;
        }
      } catch {
        // Ignorer erreur de parsing locale
      }
    }
  }
  return headers;
};

export async function apiFetch<T>(
  endpoint: string,
  options: RequestInit = {},
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;
  let headers: Record<string, string> = {
    ...getAuthHeaders(),
    ...(options.headers as Record<string, string> | undefined),
  };

  try {
    let response = await fetch(url, { ...options, headers });

    if (response.status === 401) {
      try {
        const newToken = await authBackend.refreshAccessToken();
        if (newToken) {
          headers = { ...headers, Authorization: `Bearer ${newToken}` };
          response = await fetch(url, { ...options, headers });
        }
      } catch {
        // Ignorer l'échec de rafraîchissement
      }
    }

    const result = await response.json();
    if (!response.ok || !result.success) {
      throw new Error(result.error || `Erreur API (${response.status})`);
    }
    return result.data as T;
  } catch (err: any) {
    if (err instanceof TypeError && err.message.includes("fetch")) {
      throw new Error(
        `Impossible de contacter le serveur API (${BACKEND_ORIGIN || url}). Vérifiez que le backend est démarré et joignable.`,
      );
    }
    throw err;
  }
}

export async function invokeTauri<T>(
  cmd: string,
  args?: Record<string, any>,
): Promise<T> {
  try {
    const { invoke } = await import("@tauri-apps/api/core");
    return await invoke<T>(cmd, args);
  } catch (err: any) {
    if (err instanceof Error) throw err;
    const msg =
      typeof err === "string"
        ? err
        : err?.message ||
          err?.error ||
          JSON.stringify(err) ||
          "Erreur locale Tauri";
    throw new Error(msg);
  }
}

export function mirrorToTauri(cmd: string, args?: Record<string, any>): void {
  if (!isTauri) return;
  import("@tauri-apps/api/core")
    .then(({ invoke }) => invoke(cmd, args))
    .catch(() => {});
}

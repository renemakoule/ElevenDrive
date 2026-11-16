// ==============================================================================
// AutoSchool - Service de Synchronisation Cloud Neontech (Offline-First Sync Engine)
// Version : 1.1.0 (Détection de connectivité réelle & Sonde active)
// ==============================================================================

import { api } from "../db/api";
import { API_BASE_URL } from "../config/env";

export interface SyncEventItem {
  event_id: string;
  agency_id: string;
  entity_name: string;
  operation: "INSERT" | "UPDATE" | "DELETE";
  entity_id: string;
  payload: Record<string, any>;
  local_timestamp_utc: string;
  sync_status: "EN_ATTENTE" | "SYNCHRONISE" | "CONFLIT" | "REJETE";
}

export interface SyncBatchResult {
  statut: "SUCCES" | "ERREUR" | "HORS_LIGNE";
  evenements_traites: number;
  doublons_ignores: number;
  erreurs: number;
  synced_at: string;
  message?: string;
}

export interface SyncState {
  isOnline: boolean;
  isSyncing: boolean;
  neonConnected: boolean;
  neonEndpoint: string;
  lastSyncTimestamp: string | null;
  pendingCount: number;
  lastBatchResult: SyncBatchResult | null;
  syncHistory: Array<{
    id: string;
    timestamp: string;
    count: number;
    statut: string;
    details: string;
  }>;
}

class NeonSyncEngine {
  private state: SyncState = {
    isOnline: typeof navigator !== "undefined" ? navigator.onLine : true,
    isSyncing: false,
    neonConnected: typeof navigator !== "undefined" ? navigator.onLine : true,
    neonEndpoint:
      (typeof localStorage !== "undefined" &&
        localStorage.getItem("autoschool_neon_url")) ||
      (import.meta.env?.VITE_NEON_DISPLAY_ENDPOINT as string) ||
      "",
    lastSyncTimestamp: null,
    pendingCount: 0,
    lastBatchResult: null,
    syncHistory: [],
  };

  private listeners: Array<(state: SyncState) => void> = [];
  private syncTimer: number | null = null;
  private probeTimer: number | null = null;

  constructor() {
    if (typeof window !== "undefined") {
      // 1. Événements natifs immédiats du système d'exploitation
      window.addEventListener("online", () => {
        this.verifierConnectiviteReelle();
      });

      window.addEventListener("offline", () => {
        this.updateOnlineStatus(false);
      });

      // 2. Vérification de connectivité réelle dès l'initialisation
      this.verifierConnectiviteReelle();

      // 3. Sonde active régulière (toutes les 6 secondes) pour détecter si Internet
      // a été coupé sur la box même si la carte Wi-Fi/Ethernet reste connectée
      this.probeTimer = window.setInterval(() => {
        this.verifierConnectiviteReelle();
      }, 6000);

      // 4. Timer périodique de synchronisation (toutes les 25 secondes si en ligne)
      this.syncTimer = window.setInterval(() => {
        if (this.state.isOnline && !this.state.isSyncing) {
          this.triggerAutoSync();
        }
      }, 25000);
    }
  }

  /**
   * Sonde de connectivité Internet réelle
   * (Contourne le fait que navigator.onLine reste 'true' sur Windows quand le Wi-Fi est actif mais sans Internet)
   */
  public async verifierConnectiviteReelle(): Promise<boolean> {
    if (typeof navigator !== "undefined" && !navigator.onLine) {
      this.updateOnlineStatus(false);
      return false;
    }

    try {
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), 2500);

      try {
        const localPing = await fetch(`${API_BASE_URL}/sync/ping`, {
          method: "GET",
          cache: "no-store",
          signal: controller.signal,
        });
        if (localPing.ok) {
          clearTimeout(timeoutId);
          this.state.neonConnected = true;
          this.updateOnlineStatus(true);
          return true;
        }
      } catch {
        // Fallback sonde externe si le backend local n'a pas répondu
      }

      // Test HTTP universel très léger (Android/Chromium captive portal 204)
      await fetch("https://www.google.com/generate_204", {
        method: "HEAD",
        mode: "no-cors",
        cache: "no-store",
        signal: controller.signal,
      });

      clearTimeout(timeoutId);
      this.updateOnlineStatus(true);
      return true;
    } catch {
      // Échec : pas de connectivité Internet réelle
      this.updateOnlineStatus(false);
      return false;
    }
  }

  private updateOnlineStatus(enLigne: boolean) {
    const change =
      this.state.isOnline !== enLigne ||
      (enLigne === false && this.state.neonConnected !== false);
    this.state.isOnline = enLigne;
    if (!enLigne) {
      this.state.neonConnected = false;
    }
    if (change) {
      this.notify();
    }
  }

  public subscribe(listener: (state: SyncState) => void) {
    this.listeners.push(listener);
    listener({ ...this.state });
    return () => {
      this.listeners = this.listeners.filter((l) => l !== listener);
    };
  }

  private notify() {
    const snapshot = { ...this.state };
    for (const listener of this.listeners) {
      listener(snapshot);
    }
  }

  public getState(): SyncState {
    return { ...this.state };
  }

  public setNeonEndpoint(url: string) {
    this.state.neonEndpoint = url;
    if (typeof localStorage !== "undefined") {
      localStorage.setItem("autoschool_neon_url", url);
    }
    this.notify();
  }

  public setPendingCount(count: number) {
    this.state.pendingCount = count;
    this.notify();
  }

  public async testNeonConnection(
    url?: string,
  ): Promise<{ success: boolean; latencyMs: number; message: string }> {
    const targetUrl = url || this.state.neonEndpoint;
    const start = Date.now();

    const isInternetReel = await this.verifierConnectiviteReelle();
    if (!isInternetReel) {
      return {
        success: false,
        latencyMs: 0,
        message: "Aucun accès Internet détecté (Mode hors-ligne local actif)",
      };
    }

    try {
      const response = await fetch(`${API_BASE_URL}/sync/ping`, {
        method: "GET",
      });
      const latency = Date.now() - start;

      if (response.ok) {
        this.state.neonConnected = true;
        this.notify();
        return {
          success: true,
          latencyMs: latency,
          message: `Connecté avec succès au cluster Neon Postgres (${latency} ms)`,
        };
      }

      if (!targetUrl || targetUrl.trim().length === 0) {
        return {
          success: false,
          latencyMs: 0,
          message: "URL de connexion Neon non renseignée",
        };
      }

      this.state.neonConnected = false;
      this.notify();
      return {
        success: false,
        latencyMs: latency,
        message: `Cluster PostgreSQL non joignable (HTTP ${response.status})`,
      };
    } catch (err: any) {
      this.state.neonConnected = false;
      this.notify();
      return {
        success: false,
        latencyMs: 0,
        message:
          err.message || "Impossible de joindre le serveur API Go / Neon Cloud",
      };
    }
  }

  public async triggerAutoSync(): Promise<SyncBatchResult> {
    if (this.state.isSyncing) {
      return {
        statut: "HORS_LIGNE",
        evenements_traites: 0,
        doublons_ignores: 0,
        erreurs: 0,
        synced_at: new Date().toISOString(),
        message: "Synchronisation déjà en cours",
      };
    }

    if (!this.state.isOnline) {
      return {
        statut: "HORS_LIGNE",
        evenements_traites: 0,
        doublons_ignores: 0,
        erreurs: 0,
        synced_at: new Date().toISOString(),
        message: "Réseau indisponible (Mode hors-ligne local actif)",
      };
    }

    this.state.isSyncing = true;
    this.notify();

    try {
      // 1. PUSH : SQLite local -> Cloud Neon (rejeu des événements EN_ATTENTE)
      const res = await api.synchroniserEvenementsEnAttente();

      // 2. PULL : Cloud Neon -> SQLite local (hydratation sans écraser EN_ATTENTE)
      let pulledTotal = 0;
      try {
        const pullRes = await api.pullCloudSnapshot();
        if (pullRes.pulled && pullRes.stats) {
          pulledTotal =
            (pullRes.stats.candidats_appliques || 0) +
            (pullRes.stats.dossiers_appliques || 0) +
            (pullRes.stats.paiements_appliques || 0) +
            (pullRes.stats.formules_appliquees || 0) +
            (pullRes.stats.vehicules_appliques || 0) +
            (pullRes.stats.utilisateurs_appliques || 0) +
            (pullRes.stats.sessions_caisse_appliquees || 0) +
            (pullRes.stats.seances_code_appliquees || 0) +
            (pullRes.stats.sessions_examen_appliquees || 0) +
            (pullRes.stats.candidats_examen_appliquees || 0);
        }
      } catch (pullErr) {
        console.warn("Pull Cloud -> SQLite ignoré:", pullErr);
      }

      const pushedCount = res.evenements_traites || 0;
      const nowIso = res.synced_at || new Date().toISOString();

      const result: SyncBatchResult = {
        statut: "SUCCES",
        evenements_traites: pushedCount,
        doublons_ignores: res.doublons_ignores || 0,
        erreurs: res.erreurs || 0,
        synced_at: nowIso,
        message:
          res.message ||
          `Push: ${pushedCount} envoyé(s) · Pull: ${pulledTotal} élément(s) synchronisé(s)`,
      };

      this.state.lastSyncTimestamp = nowIso;
      this.state.pendingCount = 0;
      this.state.lastBatchResult = result;
      this.state.neonConnected = true;

      this.state.syncHistory.unshift({
        id: `batch-${Date.now()}`,
        timestamp: new Date().toLocaleTimeString("fr-FR"),
        count: pushedCount + pulledTotal,
        statut: "VALIDÉ",
        details: `Bidirectionnel — Push: ${pushedCount}, Pull: ${pulledTotal}`,
      });

      if (this.state.syncHistory.length > 20) {
        this.state.syncHistory.pop();
      }

      this.notify();
      return result;
    } catch (err: any) {
      this.state.neonConnected = false;
      this.state.isOnline = false;

      const errorResult: SyncBatchResult = {
        statut: "ERREUR",
        evenements_traites: 0,
        doublons_ignores: 0,
        erreurs: 1,
        synced_at: new Date().toISOString(),
        message:
          err.message || "Erreur lors de la synchronisation vers Neontech",
      };
      this.state.lastBatchResult = errorResult;
      this.notify();
      return errorResult;
    } finally {
      this.state.isSyncing = false;
      this.notify();
    }
  }
}

export const neonSyncEngine = new NeonSyncEngine();

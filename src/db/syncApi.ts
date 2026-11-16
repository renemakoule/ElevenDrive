import { isTauri, apiFetch, invokeTauri } from "./httpClient";

export interface SyncPullStats {
  formules_appliquees: number;
  candidats_appliques: number;
  dossiers_appliques: number;
  paiements_appliques: number;
  vehicules_appliques: number;
  utilisateurs_appliques?: number;
  sessions_caisse_appliquees?: number;
  seances_code_appliquees?: number;
  sessions_examen_appliquees?: number;
  candidats_examen_appliquees?: number;
  ignores_en_attente: number;
}

export interface SyncPullResult {
  pulled: boolean;
  stats?: SyncPullStats;
}

export const syncApi = {
  /**
   * Télécharge l'instantané Cloud Neon (GET /api/sync/pull) et l'applique
   * dans la base SQLite locale sans écraser les mutations EN_ATTENTE.
   */
  async pullCloudSnapshot(): Promise<SyncPullResult> {
    if (!isTauri) {
      return { pulled: false };
    }
    const snapshot = await apiFetch<Record<string, unknown>>("/sync/pull", {
      method: "GET",
    });
    const stats = await invokeTauri<SyncPullStats>("appliquer_snapshot_cloud", {
      snapshot,
    });
    return { pulled: true, stats };
  },
};


import { api } from "../../db/api";
import { neonSyncEngine } from "../../services/neonSyncService";
import type {
  ImporterCandidatsBatchInput,
  ImporterCandidatsBatchResult,
} from "../../types";

export async function execImportCandidatesBatch(
  store: any,
  payload: ImporterCandidatsBatchInput,
): Promise<ImporterCandidatsBatchResult> {
  store.loading = true;
  try {
    const res = await api.importerCandidatsBatch(payload);
    await store.search("");
    store.stats = await api.getGuichetStats();
    neonSyncEngine.setPendingCount(store.stats.evenements_sync_en_attente);

    if (res.total_crees > 0) {
      store.showToast(
        `${res.total_crees} élève(s) inscrit(s) depuis le fichier Excel.`,
        "success",
      );
    } else {
      store.showToast(
        "Aucun nouvel élève inscrit. Consultez le bilan.",
        "info",
      );
    }
    return res;
  } catch (err: any) {
    store.showToast(err.message || "Erreur lors de l'import Excel", "error");
    throw err;
  } finally {
    store.loading = false;
  }
}

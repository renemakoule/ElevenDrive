import type {
  ImporterCandidatsBatchInput,
  ImporterCandidatsBatchResult,
} from "../types";
import { realtimeSocket } from "../services/realtimeSocketService";
import {
  isTauri,
  apiFetch,
  invokeTauri,
  mirrorToTauri,
} from "./httpClient";

export const candidatsImportApi = {
  async importerCandidatsBatch(
    payload: ImporterCandidatsBatchInput,
  ): Promise<ImporterCandidatsBatchResult> {
    try {
      const res = await apiFetch<ImporterCandidatsBatchResult>(
        "/candidats/import",
        {
          method: "POST",
          body: JSON.stringify(payload),
        },
      );
      mirrorToTauri("importer_candidats_batch", { payload });
      if (res.total_crees > 0) {
        realtimeSocket.sendEvent("CANDIDAT_DOSSIER_UPDATED", {
          action: "BATCH_IMPORTED",
          total_crees: res.total_crees,
        });
      }
      return res;
    } catch (err) {
      if (isTauri) {
        const local = await invokeTauri<ImporterCandidatsBatchResult>(
          "importer_candidats_batch",
          { payload },
        );
        if (local.total_crees > 0) {
          realtimeSocket.sendEvent("CANDIDAT_DOSSIER_UPDATED", {
            action: "BATCH_IMPORTED",
            total_crees: local.total_crees,
          });
        }
        return local;
      }
      throw err;
    }
  },
};

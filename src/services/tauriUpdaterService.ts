import { ref } from "vue";
import { useGuichetStore } from "../stores/guichetStore";

export const updateAvailable = ref(false);
export const updateVersion = ref("");
export const updateNotes = ref("");
export const isCheckingUpdate = ref(false);
export const isInstallingUpdate = ref(false);
export const downloadProgress = ref(0);
export const isUpdateDismissed = ref(false);

let pendingUpdateHandle: any = null;

function isDesktopTauriEnv(): boolean {
  return (
    typeof window !== "undefined" &&
    ("__TAURI_INTERNALS__" in window || "__TAURI__" in window)
  );
}

export async function checkForAppUpdate(silent = true): Promise<void> {
  if (!isDesktopTauriEnv() || isCheckingUpdate.value) return;
  isCheckingUpdate.value = true;
  try {
    const { check } = await import("@tauri-apps/plugin-updater");
    const update = await check();
    if (update) {
      pendingUpdateHandle = update;
      updateVersion.value = update.version;
      updateNotes.value = update.body || "";
      updateAvailable.value = true;
      isUpdateDismissed.value = false;
    } else if (!silent) {
      useGuichetStore().showToast(
        "Votre application desktop est déjà à jour.",
        "success",
      );
    }
  } catch (err: any) {
    if (!silent) {
      useGuichetStore().showToast(
        `Vérification de mise à jour impossible : ${err?.message || String(err)}`,
        "error",
      );
    }
  } finally {
    isCheckingUpdate.value = false;
  }
}

export async function installAndRelaunchUpdate(): Promise<void> {
  if (!pendingUpdateHandle || isInstallingUpdate.value) return;
  isInstallingUpdate.value = true;
  downloadProgress.value = 0;

  try {
    let totalBytes = 0;
    let downloadedBytes = 0;

    await pendingUpdateHandle.downloadAndInstall((event: any) => {
      if (event.event === "Started") {
        totalBytes = event.data?.contentLength || 0;
      } else if (event.event === "Progress") {
        downloadedBytes += event.data?.chunkLength || 0;
        if (totalBytes > 0) {
          downloadProgress.value = Math.min(
            100,
            Math.round((downloadedBytes / totalBytes) * 100),
          );
        }
      } else if (event.event === "Finished") {
        downloadProgress.value = 100;
      }
    });

    const { relaunch } = await import("@tauri-apps/plugin-process");
    await relaunch();
  } catch (err: any) {
    useGuichetStore().showToast(
      `Échec de la mise à jour : ${err?.message || String(err)}`,
      "error",
    );
  } finally {
    isInstallingUpdate.value = false;
  }
}

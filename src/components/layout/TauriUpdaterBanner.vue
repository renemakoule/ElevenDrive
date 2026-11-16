<script setup lang="ts">
import {
  downloadProgress,
  installAndRelaunchUpdate,
  isInstallingUpdate,
  isUpdateDismissed,
  updateAvailable,
  updateVersion,
} from "../../services/tauriUpdaterService";
import { DownloadCloud, Loader2, X } from "lucide-vue-next";
</script>

<template>
  <div
    v-if="updateAvailable && !isUpdateDismissed"
    class="flex items-center justify-between gap-3 border-b border-brand-border bg-brand-light px-4 py-1.5 text-xs text-brand-dark shrink-0"
    role="status"
  >
    <div class="flex items-center gap-2 min-w-0">
      <DownloadCloud class="h-4 w-4 text-brand shrink-0" />
      <span class="truncate font-medium">
        Nouvelle version disponible (<strong>v{{ updateVersion }}</strong>) — Vos données locales sont conservées.
      </span>
      <span
        v-if="isInstallingUpdate && downloadProgress > 0"
        class="font-mono text-[11px] font-semibold"
      >
        {{ downloadProgress }}%
      </span>
    </div>

    <div class="flex items-center gap-2 shrink-0">
      <button
        type="button"
        :disabled="isInstallingUpdate"
        @click="installAndRelaunchUpdate"
        class="inline-flex h-6 items-center gap-1.5 rounded-[3px] bg-brand-dark px-2.5 text-[11px] font-semibold text-white transition hover:bg-brand disabled:opacity-60"
      >
        <Loader2 v-if="isInstallingUpdate" class="h-3 w-3 animate-spin" />
        <span>{{
          isInstallingUpdate
            ? "Installation..."
            : "Mettre à jour et redémarrer"
        }}</span>
      </button>

      <button
        v-if="!isInstallingUpdate"
        type="button"
        @click="isUpdateDismissed = true"
        class="rounded p-0.5 text-brand-dark/70 hover:bg-brand/10 hover:text-brand-dark"
        aria-label="Plus tard"
      >
        <X class="h-3.5 w-3.5" />
      </button>
    </div>
  </div>
</template>

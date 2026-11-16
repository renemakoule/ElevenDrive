<script setup lang="ts">
import type { DossierArchiveItem } from "../../types";
import { Trash2 } from "lucide-vue-next";

interface Props {
  isOpen: boolean;
  archive: DossierArchiveItem | null;
}

defineProps<Props>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "confirm"): void;
}>();
</script>

<template>
  <div
    v-if="isOpen && archive"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 backdrop-blur-xs"
    role="dialog"
    aria-modal="true"
  >
    <div
      class="w-full max-w-sm rounded-[4px] border border-rose-500/30 bg-card p-5 space-y-3"
    >
      <div class="flex items-center gap-2 text-rose-600">
        <Trash2 class="h-5 w-5" />
        <h3 class="text-sm font-bold text-foreground">
          Purge Définitive de l'Archive
        </h3>
      </div>

      <p class="text-xs text-muted-foreground">
        Êtes-vous absolument certain de vouloir purger définitivement l'archive
        du dossier
        <strong class="text-foreground">{{ archive.numero_dossier }}</strong>
        ({{ archive.candidat_nom }} {{ archive.candidat_prenoms }}) ?
        <br /><br />
        <span class="text-rose-600 font-semibold">
          Cette action est irréversible et supprime définitivement les documents
          scellés.
        </span>
      </p>

      <div
        class="flex items-center justify-end gap-2 pt-2 border-t border-border"
      >
        <button
          type="button"
          @click="emit('close')"
          class="px-3 py-1.5 rounded-[4px] border border-border bg-card text-xs font-medium text-foreground hover:bg-muted"
        >
          Annuler
        </button>

        <button
          type="button"
          @click="emit('confirm')"
          class="px-3 py-1.5 rounded-[4px] bg-rose-600 text-white font-bold text-xs hover:bg-rose-700 transition"
        >
          Purger Définitivement
        </button>
      </div>
    </div>
  </div>
</template>

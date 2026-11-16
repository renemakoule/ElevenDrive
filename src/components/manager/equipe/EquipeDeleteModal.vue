<script setup lang="ts">
import type { UserTeamMember } from "../../../types";
import { Trash2, RefreshCw } from "lucide-vue-next";

defineProps<{
  isOpen: boolean;
  member: UserTeamMember | null;
  isDeleting: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "confirm"): void;
}>();
</script>

<template>
  <div
    v-if="isOpen && member"
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm animate-in fade-in duration-200"
    @click.self="emit('close')"
  >
    <div
      class="relative w-full max-w-md overflow-hidden rounded-[5px] border border-border bg-card shadow-2xl animate-in zoom-in-95 duration-200 p-6"
    >
      <div
        class="flex items-center gap-3 text-rose-600 dark:text-rose-400 mb-4"
      >
        <div
          class="flex h-10 w-10 items-center justify-center rounded-xl bg-rose-500/10 border border-rose-500/20"
        >
          <Trash2 class="h-5 w-5" />
        </div>
        <div>
          <h3 class="text-sm font-bold text-foreground">
            Supprimer le collaborateur
          </h3>
          <p class="text-xs text-muted-foreground">
            Révocation définitive des accès
          </p>
        </div>
      </div>

      <p class="text-xs text-muted-foreground leading-relaxed">
        Êtes-vous sûr de vouloir supprimer définitivement
        <strong class="text-foreground"
          >{{ member.prenoms }} {{ member.nom }}</strong
        >
        (<span class="capitalize">{{ member.role.toLowerCase() }}</span
        >) ?
      </p>
      <p
        class="mt-2 text-[11px] text-muted-foreground/80 bg-muted/40 p-2.5 rounded border border-border"
      >
        Tous les identifiants de ce collaborateur seront désactivés
        immédiatement.
      </p>

      <div class="mt-6 flex items-center justify-end gap-2">
        <button
          type="button"
          @click="emit('close')"
          class="h-9 rounded-lg border border-border px-4 text-xs font-medium text-foreground hover:bg-muted transition"
        >
          Annuler
        </button>
        <button
          type="button"
          :disabled="isDeleting"
          @click="emit('confirm')"
          class="inline-flex h-9 items-center justify-center gap-2 rounded-lg bg-rose-600 px-4 text-xs font-semibold text-white hover:bg-rose-700 transition disabled:opacity-50"
        >
          <RefreshCw v-if="isDeleting" class="h-3.5 w-3.5 animate-spin" />
          <span>{{
            isDeleting ? "Suppression..." : "Supprimer définitivement"
          }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

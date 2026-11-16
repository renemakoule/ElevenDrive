<script setup lang="ts">
import { ref, watch } from "vue";
import { ShieldAlert, Key, X } from "lucide-vue-next";

const props = defineProps<{
  isOpen: boolean;
  cand: any | null;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "submit", payload: { pin: string; motif: string }): void;
}>();

const overridePin = ref("");
const overrideMotif = ref("");

watch(
  () => props.isOpen,
  (open) => {
    if (open) {
      overridePin.value = "";
      overrideMotif.value = "";
    }
  },
);

const handleSubmit = () => {
  if (!overridePin.value.trim()) return;
  emit("submit", {
    pin: overridePin.value.trim(),
    motif: overrideMotif.value.trim(),
  });
};
</script>

<template>
  <div
    v-if="isOpen && cand"
    class="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/55 p-0 sm:p-4 animate-in fade-in"
    role="dialog"
    aria-modal="true"
    @click.self="emit('close')"
  >
    <div
      class="w-full sm:max-w-sm overflow-hidden rounded-t-[3px] sm:rounded-[3px] border-t sm:border border-border bg-card pb-safe"
    >
      <div
        class="sm:hidden mx-auto mt-2 h-1 w-9 rounded-full bg-border shrink-0"
      />
      <header
        class="flex items-center justify-between gap-3 border-b border-border px-4 py-3 bg-muted/30"
      >
        <div class="flex items-center space-x-2 min-w-0">
          <div
            class="w-8 h-8 rounded-full bg-amber-500/10 border border-amber-500/20 text-amber-600 flex items-center justify-center shrink-0"
          >
            <ShieldAlert class="w-4 h-4" />
          </div>
          <div>
            <h2 class="text-sm font-semibold text-foreground">
              Dérogation Manager
            </h2>
            <p class="text-[11px] text-muted-foreground">
              Déblocage conduite · Code PIN requis
            </p>
          </div>
        </div>
        <button
          type="button"
          @click="emit('close')"
          class="flex h-7 w-7 items-center justify-center rounded-md text-muted-foreground transition hover:bg-muted hover:text-foreground"
        >
          <X class="h-4 w-4" />
        </button>
      </header>

      <form @submit.prevent="handleSubmit" class="p-4 space-y-3.5">
        <div
          class="p-2.5 rounded bg-amber-500/10 border border-amber-500/20 text-amber-800 dark:text-amber-300 text-xs"
        >
          Vous autorisez exceptionnellement l'élève
          <strong>{{ cand.candidat_nom }}</strong> (dossier
          {{ cand.numero_dossier }}) à débuter ses heures pratiques de conduite
          malgré un solde restant.
        </div>

        <div>
          <label class="block text-xs font-semibold text-foreground mb-1"
            >Code PIN Manager (4 chiffres) *</label
          >
          <div class="relative">
            <input
              v-model="overridePin"
              type="password"
              maxlength="6"
              required
              placeholder="••••"
              class="w-full h-9 pl-9 pr-3 rounded border border-border bg-background text-foreground text-center font-mono text-base tracking-widest focus:border-brand-dark focus:outline-none"
              autofocus
            />
            <Key
              class="w-4 h-4 text-muted-foreground absolute left-3 top-2.5"
            />
          </div>
        </div>

        <div>
          <label class="block text-xs font-medium text-muted-foreground mb-1"
            >Motif de la dérogation (optionnel)</label
          >
          <input
            v-model="overrideMotif"
            type="text"
            placeholder="Ex : Accord de paiement échelonné fin de mois"
            class="w-full h-8 px-2.5 rounded border border-border bg-background text-foreground text-xs focus:border-brand-dark focus:outline-none"
          />
        </div>

        <div
          class="flex items-center justify-end gap-2 pt-2 border-t border-border"
        >
          <button
            type="button"
            @click="emit('close')"
            class="px-3 py-1.5 rounded border border-border hover:bg-muted text-xs font-medium text-foreground transition"
          >
            Annuler
          </button>
          <button
            type="submit"
            class="px-4 py-1.5 rounded bg-brand-dark hover:bg-brand text-xs font-semibold text-white transition shadow-xs"
          >
            Valider la dérogation
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

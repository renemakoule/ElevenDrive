<script setup lang="ts">
import { ref, watch } from "vue";
import type { DossierNonEligible } from "../../types";
import { ShieldAlert, X, KeyRound } from "lucide-vue-next";

const props = defineProps<{
  isOpen: boolean;
  dossier: DossierNonEligible | null;
  isSubmitting: boolean;
  errorMsg: string;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "confirm", payload: { motif: string; codePin: string }): void;
}>();

const codePinManager = ref("");
const motifForcage = ref("Dérogation exceptionnelle accordée par la Direction");

watch(
  () => props.isOpen,
  (open) => {
    if (open) {
      codePinManager.value = "";
      motifForcage.value = "Dérogation exceptionnelle accordée par la Direction";
    }
  },
);

const submit = () => {
  if (props.isSubmitting || !codePinManager.value) return;
  emit("confirm", {
    motif: motifForcage.value,
    codePin: codePinManager.value,
  });
};
</script>

<template>
  <div
    v-if="isOpen && dossier"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 select-none backdrop-blur-xs"
  >
    <div
      class="w-full max-w-md overflow-hidden rounded-xl border border-border bg-card shadow-2xl animate-in fade-in zoom-in-95 duration-150"
    >
      <div
        class="flex items-center justify-between border-b border-border bg-amber-500/10 px-4 py-3"
      >
        <div class="flex items-center gap-2 text-amber-900 dark:text-amber-200">
          <ShieldAlert class="h-4 w-4 text-amber-600 dark:text-amber-400" />
          <h3 class="text-xs font-semibold uppercase tracking-wide">
            Dérogation Direction · Forçage au bordereau
          </h3>
        </div>
        <button
          type="button"
          @click="emit('close')"
          class="rounded-md p-1 text-muted-foreground hover:bg-muted hover:text-foreground transition"
        >
          <X class="h-4 w-4" />
        </button>
      </div>

      <div class="space-y-3.5 p-4 text-xs">
        <div
          class="rounded-lg border border-amber-300/80 bg-amber-50/70 p-3 text-[11px] text-amber-900 dark:border-amber-800/60 dark:bg-amber-950/30 dark:text-amber-200"
        >
          <p class="font-semibold text-amber-950 dark:text-amber-100 mb-1">
            Responsabilité officielle du Directeur
          </p>
          <p>
            Cet élève ne remplit pas actuellement tous les critères
            automatiques (solde non soldé, heures incomplètes ou visite
            médicale). Le forçage manuel l'inscrit directement sur le bordereau
            officiel A4 ainsi que sur le procès-verbal d'examen.
          </p>
        </div>

        <div class="rounded-md border border-border bg-muted/40 p-2.5">
          <div class="text-[10px] text-muted-foreground uppercase font-semibold">
            Candidat à inscrire
          </div>
          <div class="font-bold text-foreground text-sm mt-0.5">
            {{ dossier.nom }}
          </div>
          <div class="font-mono text-[10px] text-muted-foreground">
            Dossier N° {{ dossier.numero_dossier }}
          </div>

          <div class="mt-2 text-[10px]">
            <span class="font-medium text-foreground"
              >Motifs d'inéligibilité :</span
            >
            <ul
              class="list-disc list-inside mt-0.5 text-amber-700 dark:text-amber-400 space-y-0.5"
            >
              <li v-for="(m, i) in dossier.motifs" :key="i">{{ m }}</li>
            </ul>
          </div>
        </div>

        <div>
          <label class="block text-[10px] font-semibold text-foreground mb-1">
            Motif de la dérogation
          </label>
          <input
            type="text"
            v-model="motifForcage"
            placeholder="Ex : Accord exceptionnel de la Direction, solde échelonné"
            class="w-full h-8 px-2.5 rounded-md border border-border bg-background text-foreground text-xs focus:outline-none focus:ring-1 focus:ring-brand"
          />
        </div>

        <div>
          <label class="block text-[10px] font-semibold text-foreground mb-1">
            Code PIN Directeur requis (4 chiffres)
          </label>
          <div class="relative">
            <KeyRound
              class="absolute left-2.5 top-2 h-4 w-4 text-muted-foreground"
            />
            <input
              type="password"
              maxlength="6"
              v-model="codePinManager"
              placeholder="••••"
              class="w-full h-8 pl-8 pr-2.5 rounded-md border border-border bg-background text-foreground font-mono text-xs tracking-widest focus:outline-none focus:ring-1 focus:ring-brand"
              @keyup.enter="submit"
              autofocus
            />
          </div>
          <p
            v-if="errorMsg"
            class="mt-1 text-[11px] font-medium text-destructive"
          >
            {{ errorMsg }}
          </p>
        </div>
      </div>

      <div
        class="flex items-center justify-end gap-2 border-t border-border bg-muted/20 px-4 py-2.5"
      >
        <button
          type="button"
          @click="emit('close')"
          :disabled="isSubmitting"
          class="rounded-md border border-border px-3 py-1.5 text-xs font-medium text-muted-foreground hover:bg-muted hover:text-foreground transition disabled:opacity-50"
        >
          Annuler
        </button>
        <button
          type="button"
          @click="submit"
          :disabled="isSubmitting || !codePinManager"
          class="inline-flex items-center gap-1.5 rounded-md bg-amber-600 hover:bg-amber-700 px-3.5 py-1.5 text-xs font-semibold text-white transition disabled:opacity-50 shadow-xs"
        >
          <ShieldAlert class="h-3.5 w-3.5" />
          <span>{{
            isSubmitting
              ? "Validation en cours..."
              : "Valider la dérogation et inscrire"
          }}</span>
        </button>
      </div>
    </div>
  </div>
</template>


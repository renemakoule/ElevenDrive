<script setup lang="ts">
import type { DossierDetailComplet } from "../../../types";
import { AlertTriangle, CheckCircle2 } from "lucide-vue-next";

defineProps<{
  detail: DossierDetailComplet;
  pctCode: number;
  pctConduite: number;
  conduiteVerrouillee: boolean;
  derogationActive: boolean;
}>();

const emit = defineEmits<{
  (e: "open-override"): void;
}>();
</script>

<template>
  <section class="space-y-4 p-4">
    <h3 class="text-sm font-semibold text-foreground">Progression</h3>

    <!-- Code -->
    <div>
      <div class="mb-1.5 flex items-center justify-between text-xs">
        <span class="text-foreground">Code de la route</span>
        <span class="tabular-nums text-muted-foreground">
          <span class="font-medium text-foreground">
            {{ detail.dossier.heures_code_effectuees }} /
            {{ detail.formule.heures_theoriques }} h
          </span>
          · {{ pctCode }} %
        </span>
      </div>
      <div
        class="h-1 overflow-hidden rounded-full bg-muted"
        role="progressbar"
        :aria-valuenow="pctCode"
        aria-valuemin="0"
        aria-valuemax="100"
        aria-label="Progression du code"
      >
        <div
          class="h-full rounded-full bg-brand transition-all duration-500"
          :style="{ width: pctCode + '%' }"
        />
      </div>
    </div>

    <!-- Conduite -->
    <div>
      <div class="mb-1.5 flex items-center justify-between text-xs">
        <span class="text-foreground">Conduite pratique</span>
        <span class="tabular-nums text-muted-foreground">
          <span class="font-medium text-foreground">
            {{ detail.dossier.heures_conduite_effectuees }} /
            {{ detail.formule.heures_pratiques }} h
          </span>
          · {{ pctConduite }} %
        </span>
      </div>
      <div
        class="h-1 overflow-hidden rounded-full bg-muted"
        role="progressbar"
        :aria-valuenow="pctConduite"
        aria-valuemin="0"
        aria-valuemax="100"
        aria-label="Progression de la conduite"
      >
        <div
          class="h-full rounded-full bg-brand transition-all duration-500"
          :style="{ width: pctConduite + '%' }"
        />
      </div>

      <!-- Verrou de démarrage de la conduite -->
      <template v-if="conduiteVerrouillee">
        <div
          v-if="derogationActive"
          class="mt-3 flex items-center gap-2 rounded-lg bg-brand-light px-3 py-2 text-xs text-brand-dark"
        >
          <CheckCircle2 class="h-3.5 w-3.5 shrink-0 text-brand" />
          <span class="flex-1"
            >Dérogation manager active : conduite déverrouillée.</span
          >
        </div>

        <div
          v-else
          class="mt-3 flex items-center gap-2 rounded-lg bg-amber-50 px-3 py-2 text-xs text-amber-800"
        >
          <AlertTriangle class="h-3.5 w-3.5 shrink-0 text-amber-600" />
          <span class="flex-1"
            >Tranche 2 requise (≥ 60 % payé) pour débuter la conduite.</span
          >
          <button
            type="button"
            @click="emit('open-override')"
            class="shrink-0 rounded-md bg-amber-600 px-2.5 py-1 text-[11px] font-medium text-white transition hover:bg-amber-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-500/50"
          >
            Dérogation PIN
          </button>
        </div>
      </template>
    </div>
  </section>
</template>

<script setup lang="ts">
import { Phone, Send } from "lucide-vue-next";

defineProps<{
  impayes: any[];
}>();

const emit = defineEmits<{
  (e: "relance", cand: any): void;
}>();

const fcfa = (n: number) => (n ?? 0).toLocaleString("fr-FR");
const pct = (n: number) => Math.max(0, Math.min(100, Math.round(n ?? 0)));
</script>

<template>
  <div class="flex flex-1 flex-col overflow-hidden">
    <div class="flex items-center justify-between gap-3 px-4 py-3">
      <h2 class="text-sm font-semibold text-foreground">
        Dossiers avec solde dû
      </h2>
      <span class="text-[11px] text-muted-foreground"
        >« Relancer » prépare un message WhatsApp prêt à envoyer</span
      >
    </div>

    <ul
      v-if="impayes.length"
      class="flex-1 divide-y divide-border overflow-y-auto border-t border-border"
    >
      <li
        v-for="cand in impayes"
        :key="cand.dossier_id"
        class="flex flex-col gap-3 px-4 py-3 transition hover:bg-muted/50 sm:flex-row sm:items-center sm:justify-between"
      >
        <div class="min-w-0 space-y-1">
          <div class="flex flex-wrap items-center gap-x-2">
            <span class="text-xs font-semibold text-foreground">{{
              cand.candidat_nom
            }}</span>
            <span class="font-mono text-[11px] font-medium text-brand-dark">{{
              cand.numero_dossier
            }}</span>
            <span class="text-[11px] text-muted-foreground">{{
              cand.formule
            }}</span>
          </div>
          <div
            class="flex flex-wrap items-center gap-x-3 gap-y-0.5 text-[11px] text-muted-foreground"
          >
            <span class="inline-flex items-center gap-1">
              <Phone class="h-3 w-3" />
              {{ cand.telephone_whatsapp }}
            </span>
            <span
              >Code {{ cand.heures_code }} /
              {{ cand.heures_code_requises }} h</span
            >
            <span>Inscrit le {{ cand.date_inscription }}</span>
          </div>
        </div>

        <div
          class="flex w-full sm:w-auto shrink-0 items-center justify-between sm:justify-end gap-3 sm:gap-4"
        >
          <div class="w-40 sm:w-44">
            <div
              class="sm:text-right text-sm font-semibold tabular-nums text-amber-600"
            >
              {{ fcfa(cand.solde_restant) }} FCFA
            </div>
            <div
              class="mt-1.5 h-1 overflow-hidden rounded-full bg-muted"
              role="progressbar"
              :aria-valuenow="pct(cand.taux_reglement_pct)"
              aria-valuemin="0"
              aria-valuemax="100"
              :aria-label="`Règlement de ${cand.candidat_nom}`"
            >
              <div
                class="h-full rounded-full bg-brand"
                :style="{ width: pct(cand.taux_reglement_pct) + '%' }"
              />
            </div>
            <div
              class="mt-1 sm:text-right text-[11px] tabular-nums text-muted-foreground"
            >
              Réglé {{ fcfa(cand.deja_paye) }} F ·
              {{ pct(cand.taux_reglement_pct) }} %
            </div>
          </div>

          <button
            type="button"
            @click="emit('relance', cand)"
            class="inline-flex h-9 sm:h-8 items-center gap-1.5 rounded-[3px] border border-border bg-card px-3 text-xs font-medium text-brand-dark transition hover:bg-brand-light focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-dark/40"
          >
            <Send class="h-3.5 w-3.5" />
            Relancer
          </button>
        </div>
      </li>
    </ul>
    <p
      v-else
      class="border-t border-border px-4 py-10 text-center text-xs text-muted-foreground"
    >
      Aucun impayé à relancer.
    </p>
  </div>
</template>

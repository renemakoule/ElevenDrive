<script setup lang="ts">
import type { ImporterCandidatsBatchResult } from "../../../types";
import { exportImportErrorsCsv } from "../../../utils/excelCandidates";
import { CheckCircle2, Download, AlertCircle } from "lucide-vue-next";

const props = defineProps<{
  result: ImporterCandidatsBatchResult;
}>();

const fcfa = (n: number) => (n || 0).toLocaleString("fr-FR");
</script>

<template>
  <div class="space-y-3">
    <div class="grid grid-cols-3 gap-2 text-center">
      <div class="rounded-[3px] border border-emerald-200 bg-emerald-50 p-2.5">
        <div class="text-lg font-bold text-emerald-700">
          {{ result.total_crees }}
        </div>
        <div class="text-[10px] font-medium uppercase text-emerald-800">
          Inscrits
        </div>
      </div>
      <div class="rounded-[3px] border border-amber-200 bg-amber-50 p-2.5">
        <div class="text-lg font-bold text-amber-700">
          {{ result.total_ignores }}
        </div>
        <div class="text-[10px] font-medium uppercase text-amber-800">
          Doublons ignorés
        </div>
      </div>
      <div class="rounded-[3px] border border-rose-200 bg-rose-50 p-2.5">
        <div class="text-lg font-bold text-rose-700">
          {{ result.total_erreurs }}
        </div>
        <div class="text-[10px] font-medium uppercase text-rose-800">
          Erreurs
        </div>
      </div>
    </div>

    <div
      v-if="result.dossiers_crees.length > 0"
      class="max-h-40 overflow-auto rounded-[3px] border border-border"
    >
      <div class="border-b border-border bg-muted/60 px-2.5 py-1.5 text-[10px] font-semibold uppercase text-muted-foreground">
        Derniers dossiers créés
      </div>
      <ul class="divide-y divide-border text-[11px]">
        <li
          v-for="dos in result.dossiers_crees.slice(0, 30)"
          :key="dos.dossier_id"
          class="flex items-center justify-between px-2.5 py-1.5"
        >
          <span class="flex items-center gap-1.5 font-medium text-foreground">
            <CheckCircle2 class="h-3.5 w-3.5 text-emerald-600" />
            {{ dos.nom_complet }}
          </span>
          <span class="font-mono text-muted-foreground">
            {{ dos.numero_dossier }} · Solde: {{ fcfa(dos.solde_restant) }} F
          </span>
        </li>
      </ul>
    </div>

    <div
      v-if="result.erreurs.length > 0"
      class="space-y-1.5 rounded-[3px] border border-rose-200 bg-rose-50/50 p-2.5"
    >
      <div class="flex items-center justify-between">
        <span class="flex items-center gap-1.5 text-xs font-semibold text-rose-800">
          <AlertCircle class="h-3.5 w-3.5" />
          Lignes non importées ({{ result.erreurs.length }})
        </span>
        <button
          type="button"
          @click="exportImportErrorsCsv(props.result.erreurs)"
          class="inline-flex items-center gap-1 rounded-[3px] border border-rose-300 bg-white px-2 py-1 text-[10px] font-medium text-rose-800 hover:bg-rose-50"
        >
          <Download class="h-3 w-3" /> Exporter rejets (.csv)
        </button>
      </div>
      <ul class="max-h-28 overflow-auto space-y-0.5 text-[11px] text-rose-700">
        <li v-for="(err, idx) in result.erreurs" :key="idx">
          #{{ err.index }} {{ err.nom }} {{ err.prenoms }} ({{ err.numero_piece_id || "—" }}) : {{ err.raison }}
        </li>
      </ul>
    </div>
  </div>
</template>

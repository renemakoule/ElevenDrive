<script setup lang="ts">
import { computed, ref } from "vue";
import type { ParsedExcelCandidateRow } from "../../../types";
import { AlertTriangle, CheckCircle2, Filter } from "lucide-vue-next";

const props = defineProps<{
  rows: ParsedExcelCandidateRow[];
  excludedCount?: number;
}>();

const filterMode = ref<"ALL" | "ERRORS">("ALL");

const validCount = computed(() => props.rows.filter((r) => r.isValid).length);
const errorCount = computed(() => props.rows.filter((r) => !r.isValid).length);

const displayedRows = computed(() =>
  filterMode.value === "ERRORS"
    ? props.rows.filter((r) => !r.isValid)
    : props.rows,
);

const fcfa = (n: number) => (n || 0).toLocaleString("fr-FR");
</script>

<template>
  <div class="space-y-2.5">
    <div class="flex flex-wrap items-center justify-between gap-2">
      <div class="flex flex-wrap items-center gap-2 text-[11px]">
        <span
          class="inline-flex items-center gap-1 rounded-[3px] border border-emerald-200 bg-emerald-50 px-2 py-0.5 font-semibold text-emerald-700"
        >
          <CheckCircle2 class="h-3.5 w-3.5" />
          {{ validCount }} valide(s)
        </span>
        <span
          v-if="errorCount > 0"
          class="inline-flex items-center gap-1 rounded-[3px] border border-rose-200 bg-rose-50 px-2 py-0.5 font-semibold text-rose-700"
        >
          <AlertTriangle class="h-3.5 w-3.5" />
          {{ errorCount }} en erreur
        </span>
        <span
          v-if="excludedCount && excludedCount > 0"
          class="inline-flex items-center gap-1 rounded-[3px] border border-border bg-muted px-2 py-0.5 font-medium text-muted-foreground"
        >
          <Filter class="h-3 w-3" />
          {{ excludedCount }} exclue(s) par filtre
        </span>
      </div>

      <div class="flex items-center gap-1 text-[11px]">
        <button
          type="button"
          @click="filterMode = 'ALL'"
          :class="[
            'rounded-[3px] border px-2 py-1 font-medium transition',
            filterMode === 'ALL'
              ? 'border-brand bg-brand-light text-brand-dark'
              : 'border-border bg-card text-muted-foreground hover:bg-muted',
          ]"
        >
          Retenues ({{ rows.length }})
        </button>
        <button
          v-if="errorCount > 0"
          type="button"
          @click="filterMode = 'ERRORS'"
          :class="[
            'rounded-[3px] border px-2 py-1 font-medium transition',
            filterMode === 'ERRORS'
              ? 'border-rose-300 bg-rose-50 text-rose-700'
              : 'border-border bg-card text-muted-foreground hover:bg-muted',
          ]"
        >
          Erreurs ({{ errorCount }})
        </button>
      </div>
    </div>

    <div
      class="max-h-64 overflow-auto rounded-[3px] border border-border bg-card"
    >
      <table class="w-full border-collapse text-left text-[11px]">
        <thead
          class="sticky top-0 z-10 border-b border-border bg-muted text-[10px] uppercase text-muted-foreground"
        >
          <tr>
            <th class="px-2.5 py-2 bg-muted">Ligne</th>
            <th class="px-2.5 py-2 bg-muted">Élève</th>
            <th class="px-2.5 py-2 bg-muted">WhatsApp</th>
            <th class="px-2.5 py-2 bg-muted">Pièce</th>
            <th class="px-2.5 py-2 bg-muted">Formule</th>
            <th class="px-2.5 py-2 bg-muted text-right">Acompte</th>
            <th class="px-2.5 py-2 bg-muted">État</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border bg-card">
          <tr v-if="displayedRows.length === 0">
            <td colspan="7" class="px-3 py-6 text-center text-muted-foreground">
              Aucune ligne ne correspond aux formules sélectionnées.
            </td>
          </tr>
          <tr
            v-for="row in displayedRows"
            :key="row.rowNumber"
            :class="row.isValid ? '' : 'bg-rose-50/40 dark:bg-rose-950/20'"
          >
            <td class="px-2.5 py-1.5 font-mono text-muted-foreground">
              #{{ row.rowNumber }}
            </td>
            <td class="px-2.5 py-1.5 font-medium text-foreground">
              {{ row.payload.nom || "—" }} {{ row.payload.prenoms || "" }}
              <div class="text-[10px] text-muted-foreground">
                {{ row.payload.genre }} · {{ row.payload.date_naissance }}
              </div>
            </td>
            <td class="px-2.5 py-1.5 font-mono">
              {{ row.payload.telephone_whatsapp || "—" }}
            </td>
            <td class="px-2.5 py-1.5">
              {{ row.payload.type_piece_id }} ·
              <span class="font-mono">{{
                row.payload.numero_piece_id || "—"
              }}</span>
            </td>
            <td class="px-2.5 py-1.5">{{ row.formuleLabel }}</td>
            <td class="px-2.5 py-1.5 text-right font-mono">
              {{ fcfa(row.payload.acompte_initial) }} F
            </td>
            <td class="px-2.5 py-1.5">
              <span
                v-if="row.isValid"
                class="inline-flex items-center gap-1 text-emerald-700 font-medium"
              >
                <CheckCircle2 class="h-3 w-3" /> Prêt
              </span>
              <span v-else class="text-rose-700 font-medium">
                {{ row.errors.join(" · ") }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

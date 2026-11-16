<script setup lang="ts">
import type { CandidatExamenItem } from "../../../types";
import { Calendar, Plus, CheckCircle2 } from "lucide-vue-next";

defineProps<{
  hasActiveSession: boolean;
  candidats: CandidatExamenItem[];
  publie: boolean;
  libelleResultat: Record<string, string>;
}>();

const emit = defineEmits<{
  (e: "planifier"): void;
  (
    e: "toggle-result",
    candidat: CandidatExamenItem,
    champ: "resultat_code" | "resultat_conduite",
  ): void;
  (e: "validate-pv"): void;
}>();

const styleResultat = (r: string) =>
  r === "ADMIS"
    ? "border-brand-border bg-brand-light text-brand-dark"
    : r === "AJOURNE"
      ? "border-rose-300 bg-rose-50 text-rose-700 dark:bg-rose-950/40 dark:text-rose-300 dark:border-rose-800"
      : "border-border bg-muted text-muted-foreground";
</script>

<template>
  <div class="flex-1 space-y-2 overflow-y-auto p-3 text-[10px]">
    <div
      v-if="!hasActiveSession"
      class="rounded-lg border border-dashed border-border bg-muted/30 p-8 text-center"
    >
      <Calendar class="mx-auto h-8 w-8 text-muted-foreground/60 mb-2" />
      <h4 class="text-xs font-semibold text-foreground">
        Aucune session d'examen active
      </h4>
      <p class="text-[11px] text-muted-foreground max-w-sm mx-auto mt-1 mb-3">
        Vous devez planifier une session d'examen d'État avant de pouvoir
        enregistrer les résultats officiels des épreuves.
      </p>
      <button
        type="button"
        @click="emit('planifier')"
        class="inline-flex items-center gap-1.5 rounded-md bg-brand-dark px-3 py-1.5 text-xs font-semibold text-white hover:bg-brand transition"
      >
        <Plus class="h-3.5 w-3.5" />
        Planifier une session
      </button>
    </div>

    <template v-else>
      <div class="text-muted-foreground">
        Cliquez sur un bouton d'épreuve pour basculer le résultat entre Admis,
        Ajourné et Absent.
      </div>

      <div class="overflow-hidden rounded-[5px] border border-border">
        <table class="w-full border-collapse text-left text-[10px]">
          <thead class="bg-muted/40 text-muted-foreground">
            <tr>
              <th class="w-20 px-3 py-1.5 font-medium">Table</th>
              <th class="px-3 py-1.5 font-medium">Candidat</th>
              <th class="w-24 px-3 py-1.5 text-center font-medium">Code</th>
              <th class="w-24 px-3 py-1.5 text-center font-medium">Conduite</th>
              <th class="w-24 px-3 py-1.5 text-center font-medium">
                Résultat officiel
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-border">
            <tr v-if="candidats.length === 0">
              <td
                colspan="5"
                class="px-3 py-6 text-center text-muted-foreground"
              >
                Aucun candidat sur ce procès-verbal.
              </td>
            </tr>
            <tr
              v-for="c in candidats"
              :key="c.id"
              class="hover:bg-muted/60"
            >
              <td class="px-3 py-1.5 font-mono font-medium text-foreground">
                {{ c.numero_table }}
              </td>
              <td class="px-3 py-1.5">
                <div class="font-medium text-foreground">
                  {{ c.nom }} {{ c.prenoms }}
                </div>
                <div class="font-mono text-[9px] text-muted-foreground">
                  {{ c.numero_dossier }}
                </div>
              </td>
              <td class="px-3 py-1.5 text-center">
                <button
                  type="button"
                  @click="emit('toggle-result', c, 'resultat_code')"
                  :class="[
                    'w-20 rounded border py-0.5 text-[10px] font-medium transition',
                    styleResultat(c.resultat_code),
                  ]"
                >
                  {{ libelleResultat[c.resultat_code] }}
                </button>
              </td>
              <td class="px-3 py-1.5 text-center">
                <button
                  type="button"
                  @click="emit('toggle-result', c, 'resultat_conduite')"
                  :class="[
                    'w-20 rounded border py-0.5 text-[10px] font-medium transition',
                    styleResultat(c.resultat_conduite),
                  ]"
                >
                  {{ libelleResultat[c.resultat_conduite] }}
                </button>
              </td>
              <td class="px-3 py-1.5 text-center">
                <span
                  :class="[
                    'inline-block w-20 rounded border py-0.5 text-[10px] font-medium',
                    styleResultat(c.resultat_final),
                  ]"
                >
                  {{ libelleResultat[c.resultat_final] }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="flex items-center justify-end gap-2 pt-1">
        <span v-if="publie" class="text-[10px] font-medium text-brand"
          >✓ Procès-verbal validé & publié</span
        >
        <button
          type="button"
          @click="emit('validate-pv')"
          :disabled="candidats.length === 0"
          class="flex items-center gap-1 rounded-md bg-brand-dark px-2.5 py-1 text-[10px] font-semibold text-white hover:bg-brand disabled:opacity-40 transition"
        >
          <CheckCircle2 class="h-3 w-3" />
          Enregistrer et valider le procès-verbal
        </button>
      </div>
    </template>
  </div>
</template>


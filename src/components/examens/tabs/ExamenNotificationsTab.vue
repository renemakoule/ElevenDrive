<script setup lang="ts">
import type { CandidatExamenItem } from "../../../types";
import { Send, Plus, Award } from "lucide-vue-next";

defineProps<{
  hasActiveSession: boolean;
  candidats: CandidatExamenItem[];
  publie: boolean;
  libelleResultat: Record<string, string>;
}>();

const emit = defineEmits<{
  (e: "planifier"): void;
  (e: "send-convocation", candidat: CandidatExamenItem): void;
  (e: "send-result", candidat: CandidatExamenItem): void;
}>();
</script>

<template>
  <div class="flex-1 space-y-2 overflow-y-auto p-3 text-[10px]">
    <div
      v-if="!hasActiveSession"
      class="rounded-lg border border-dashed border-border bg-muted/30 p-8 text-center"
    >
      <Send class="mx-auto h-8 w-8 text-muted-foreground/60 mb-2" />
      <h4 class="text-xs font-semibold text-foreground">
        Session non planifiée
      </h4>
      <p class="text-[11px] text-muted-foreground max-w-sm mx-auto mt-1 mb-3">
        Planifiez la session d'examen pour pouvoir envoyer automatiquement les
        convocations individuelles par WhatsApp.
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
        Les convocations peuvent être transmises aux élèves retenus. Les
        notifications de résultats sont débloquées dès la validation du
        procès-verbal.
      </div>

      <div
        class="divide-y divide-border overflow-hidden rounded-[5px] border border-border"
      >
        <div
          v-if="candidats.length === 0"
          class="py-6 px-3 text-center text-muted-foreground"
        >
          Aucun candidat à notifier.
        </div>
        <div
          v-for="c in candidats"
          :key="c.id"
          class="flex items-center justify-between gap-3 px-3 py-1.5 hover:bg-muted/60"
        >
          <div class="min-w-0">
            <div class="truncate font-medium text-foreground">
              {{ c.nom }} {{ c.prenoms }}
              <span
                class="ml-1 rounded bg-muted px-1 py-px font-mono text-[9px] font-normal text-muted-foreground"
                >{{ c.numero_table }}</span
              >
            </div>
            <div class="text-[9px] text-muted-foreground">
              {{ c.telephone_whatsapp }} · Résultat :
              <strong class="text-foreground">{{
                libelleResultat[c.resultat_final]
              }}</strong>
            </div>
          </div>

          <div class="flex shrink-0 items-center gap-1.5">
            <button
              type="button"
              @click="emit('send-convocation', c)"
              class="flex items-center gap-1 rounded border border-border px-2 py-0.5 text-[10px] font-medium text-foreground hover:bg-muted transition"
            >
              <Send class="h-2.5 w-2.5 text-sky-600" />
              Convocation
            </button>
            <button
              type="button"
              @click="emit('send-result', c)"
              :disabled="!publie"
              :title="
                publie
                  ? ''
                  : 'Validez d\'abord le procès-verbal dans l\'onglet précédent'
              "
              class="flex items-center gap-1 rounded border border-border px-2 py-0.5 text-[10px] font-medium text-foreground hover:bg-muted disabled:cursor-not-allowed disabled:opacity-40 transition"
            >
              <Award class="h-2.5 w-2.5 text-brand" />
              Résultat
            </button>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>


<script setup lang="ts">
import type { QuizSoumissionEleve } from "../../../types";
import { Award, RotateCcw } from "lucide-vue-next";

defineProps<{
  result: QuizSoumissionEleve;
}>();

const emit = defineEmits<{
  (e: "restart"): void;
  (e: "close"): void;
}>();
</script>

<template>
  <div
    class="rounded-lg border-2 p-5 space-y-4 bg-card shadow-sm"
    :class="
      result.est_reussi ? 'border-emerald-500/50' : 'border-amber-500/50'
    "
  >
    <div
      class="flex flex-col sm:flex-row sm:items-center justify-between gap-3"
    >
      <div class="flex items-center gap-3">
        <div
          class="h-12 w-12 rounded-full flex items-center justify-center shrink-0"
          :class="
            result.est_reussi
              ? 'bg-emerald-500/15 text-emerald-600'
              : 'bg-amber-500/15 text-amber-600'
          "
        >
          <Award class="h-6 w-6" />
        </div>
        <div>
          <span
            class="inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-[10px] font-bold"
            :class="
              result.est_reussi
                ? 'bg-emerald-500/15 text-emerald-700 dark:text-emerald-300'
                : 'bg-amber-500/15 text-amber-700'
            "
          >
            ✓ Réponses envoyées en direct à votre moniteur
          </span>
          <h3 class="text-base font-extrabold text-foreground mt-1">
            Score obtenu : {{ result.score_obtenu }} /
            {{ result.total_points }} points ({{ result.pourcentage }}%)
          </h3>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <button
          type="button"
          @click="emit('restart')"
          class="inline-flex items-center gap-1.5 rounded-md border border-border bg-background px-3 py-1.5 text-xs font-semibold text-foreground hover:bg-muted cursor-pointer"
        >
          <RotateCcw class="h-3.5 w-3.5" />
          <span>Recommencer</span>
        </button>
        <button
          type="button"
          @click="emit('close')"
          class="rounded-md bg-brand-dark px-4 py-1.5 text-xs font-bold text-white hover:bg-brand cursor-pointer"
        >
          Terminer
        </button>
      </div>
    </div>

    <div class="space-y-2.5 pt-3 border-t border-border">
      <div
        v-for="(det, idx) in result.details"
        :key="det.question_id || idx"
        class="rounded-md border p-3.5 text-xs space-y-1.5"
        :class="
          det.est_correcte
            ? 'border-emerald-500/30 bg-emerald-500/5'
            : 'border-rose-500/30 bg-rose-500/5'
        "
      >
        <div class="flex items-start justify-between gap-2">
          <span class="font-bold text-foreground">
            {{ idx + 1 }}. {{ det.intitule_question }}
          </span>
          <span
            class="shrink-0 font-bold"
            :class="det.est_correcte ? 'text-emerald-600' : 'text-rose-600'"
          >
            {{ det.points_obtenus }} / {{ det.points_max }} pt(s)
          </span>
        </div>
        <div>
          Votre réponse :
          <strong
            :class="det.est_correcte ? 'text-emerald-700' : 'text-rose-600'"
          >
            {{ det.texte_reponse || "Aucune réponse" }}
          </strong>
        </div>
        <div
          v-if="!det.est_correcte && det.bonne_reponse_label"
          class="text-emerald-700 dark:text-emerald-300"
        >
          Bonne réponse : <strong>{{ det.bonne_reponse_label }}</strong>
        </div>
        <p v-if="det.explication" class="italic text-muted-foreground">
          💡 {{ det.explication }}
        </p>
      </div>
    </div>
  </div>
</template>


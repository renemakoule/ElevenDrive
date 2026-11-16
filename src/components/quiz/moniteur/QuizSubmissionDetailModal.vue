<script setup lang="ts">
import type { QuizSoumissionEleve } from "../../../types";
import { X } from "lucide-vue-next";

defineProps<{
  submission: QuizSoumissionEleve | null;
}>();

const emit = defineEmits<{
  (e: "close"): void;
}>();

const formatDateHeure = (iso: string) => {
  if (!iso) return "À l'instant";
  const d = new Date(iso);
  if (isNaN(d.getTime())) return iso;
  return d.toLocaleDateString("fr-FR", {
    day: "2-digit",
    month: "short",
    hour: "2-digit",
    minute: "2-digit",
  });
};

const formatDuration = (sec: number) => {
  if (!sec || sec <= 0) return "< 1 min";
  const m = Math.floor(sec / 60);
  const s = sec % 60;
  return m > 0 ? `${m} min ${s}s` : `${s}s`;
};
</script>

<template>
  <div
    v-if="submission"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-xs"
    @click.self="emit('close')"
  >
    <div
      class="max-h-[90vh] w-full max-w-3xl overflow-hidden rounded-xl border border-border bg-card shadow-2xl flex flex-col"
    >
      <div
        class="flex items-center justify-between border-b border-border bg-muted/40 px-5 py-4"
      >
        <div>
          <div class="flex items-center gap-2">
            <span
              class="rounded bg-brand-light px-2 py-0.5 text-[10px] font-bold text-brand-dark"
            >
              Copie individuelle élève
            </span>
            <span class="font-mono text-xs text-muted-foreground">
              {{ submission.numero_dossier }} · Permis
              {{ submission.categorie_permis }}
            </span>
          </div>
          <h3 class="text-base font-bold text-foreground mt-1">
            {{ submission.candidat_nom }} —
            <span class="text-brand-dark">{{ submission.quiz_titre }}</span>
          </h3>
        </div>

        <div class="flex items-center gap-3">
          <div class="text-right">
            <div class="text-sm font-extrabold text-foreground">
              {{ submission.score_obtenu }} / {{ submission.total_points }} pts
              ({{ submission.pourcentage }}%)
            </div>
            <div
              class="text-[10px] font-bold"
              :class="submission.est_reussi ? 'text-emerald-600' : 'text-amber-600'"
            >
              {{
                submission.est_reussi
                  ? "Objectif atteint"
                  : "Sous le seuil de réussite"
              }}
            </div>
          </div>
          <button
            type="button"
            @click="emit('close')"
            class="rounded-lg p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground cursor-pointer"
          >
            <X class="h-4 w-4" />
          </button>
        </div>
      </div>

      <div class="flex-1 overflow-y-auto p-5 space-y-3">
        <div
          v-for="(det, idx) in submission.details"
          :key="det.question_id || idx"
          class="rounded-lg border p-4 space-y-2"
          :class="
            det.est_correcte
              ? 'border-emerald-500/30 bg-emerald-500/5'
              : 'border-rose-500/30 bg-rose-500/5'
          "
        >
          <div class="flex items-start justify-between gap-2">
            <div class="text-xs sm:text-sm font-bold text-foreground">
              {{ idx + 1 }}. {{ det.intitule_question }}
            </div>
            <span
              class="shrink-0 rounded-full px-2 py-0.5 text-[10px] font-bold"
              :class="
                det.est_correcte
                  ? 'bg-emerald-600 text-white'
                  : 'bg-rose-600 text-white'
              "
            >
              {{ det.points_obtenus }} / {{ det.points_max }} pt(s)
            </span>
          </div>

          <div class="text-xs space-y-1">
            <div class="flex items-center gap-1.5">
              <span class="text-muted-foreground">Réponse de l'élève :</span>
              <strong
                :class="
                  det.est_correcte
                    ? 'text-emerald-700 dark:text-emerald-300'
                    : 'text-rose-600 dark:text-rose-400'
                "
              >
                {{
                  det.choix_labels && det.choix_labels.length > 0
                    ? det.choix_labels.join(" · ")
                    : det.texte_reponse || "Aucune réponse"
                }}
              </strong>
            </div>

            <div
              v-if="!det.est_correcte && det.bonne_reponse_label"
              class="flex items-center gap-1.5 text-emerald-700 dark:text-emerald-300"
            >
              <span>Bonne réponse attendue :</span>
              <strong>{{ det.bonne_reponse_label }}</strong>
            </div>

            <p
              v-if="det.explication"
              class="text-[11px] italic text-muted-foreground pt-1"
            >
              💡 {{ det.explication }}
            </p>
          </div>
        </div>
      </div>

      <div
        class="border-t border-border bg-muted/30 px-5 py-3 flex items-center justify-between text-xs"
      >
        <span class="text-muted-foreground">
          Soumis le {{ formatDateHeure(submission.submitted_at) }} · Durée :
          {{ formatDuration(submission.duree_secondes) }}
        </span>
        <button
          type="button"
          @click="emit('close')"
          class="rounded-md bg-brand-dark px-4 py-1.5 font-bold text-white hover:bg-brand transition cursor-pointer"
        >
          Fermer la copie
        </button>
      </div>
    </div>
  </div>
</template>


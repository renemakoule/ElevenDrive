<script setup lang="ts">
import type { QuizFormulaire, QuizSoumissionEleve } from "../../../types";
import {
  Brain,
  Radio,
  Sparkles,
  CheckCircle2,
  FileCheck2,
  Clock,
  Send,
} from "lucide-vue-next";

const props = defineProps<{
  studentCategory: string;
  isWsConnected: boolean;
  studentQuizzes: QuizFormulaire[];
  mySubmissions: QuizSoumissionEleve[];
}>();

const emit = defineEmits<{
  (e: "start", quiz: QuizFormulaire): void;
  (e: "open-simulator"): void;
}>();

const getMyBestSubmission = (
  quizId: string,
): QuizSoumissionEleve | undefined => {
  const subs = props.mySubmissions.filter((s) => s.quiz_id === quizId);
  if (subs.length === 0) return undefined;
  return subs.reduce((best, cur) =>
    cur.pourcentage > best.pourcentage ? cur : best,
  );
};
</script>

<template>
  <div class="space-y-4">
    <!-- Bannière d'en-tête -->
    <div
      class="rounded-lg border border-brand-border bg-brand-light/30 p-4 flex flex-col sm:flex-row sm:items-center justify-between gap-3"
    >
      <div class="flex items-center gap-3">
        <div
          class="h-10 w-10 rounded-lg bg-brand-light border border-brand-border flex items-center justify-center text-brand-dark shrink-0"
        >
          <Brain class="h-5 w-5 text-brand" />
        </div>
        <div>
          <div class="flex items-center gap-2">
            <h2 class="text-sm font-bold text-foreground">
              Formulaires de Quiz du Moniteur · Permis {{ studentCategory }}
            </h2>
            <span
              class="inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-bold"
              :class="
                isWsConnected
                  ? 'bg-emerald-500/15 text-emerald-700 dark:text-emerald-300'
                  : 'bg-amber-500/15 text-amber-700'
              "
            >
              <Radio class="h-3 w-3 animate-pulse" />
              {{ isWsConnected ? "En direct" : "Sync" }}
            </span>
          </div>
          <p class="text-xs text-muted-foreground mt-0.5">
            Répondez directement aux formulaires de quiz préparés par votre
            moniteur. Vos réponses lui sont envoyées instantanément.
          </p>
        </div>
      </div>

      <button
        type="button"
        @click="emit('open-simulator')"
        class="inline-flex items-center gap-1.5 rounded-lg border border-border bg-card px-3 py-2 text-xs font-bold text-foreground hover:bg-muted transition shrink-0 cursor-pointer"
      >
        <Sparkles class="h-3.5 w-3.5 text-brand" />
        <span>Simulateur 40 questions aléatoires</span>
      </button>
    </div>

    <!-- Aucun quiz publié -->
    <div
      v-if="studentQuizzes.length === 0"
      class="rounded-lg border border-dashed border-border bg-card p-12 text-center space-y-3"
    >
      <Brain class="h-10 w-10 text-muted-foreground/40 mx-auto" />
      <h3 class="text-sm font-bold text-foreground">
        Aucun formulaire de quiz publié pour le Permis {{ studentCategory }}
      </h3>
      <p class="text-xs text-muted-foreground max-w-md mx-auto">
        Dès que votre moniteur publie un nouveau formulaire de quiz, il
        apparaîtra automatiquement ici en temps réel. En attendant, vous pouvez
        vous entraîner sur le simulateur officiel.
      </p>
      <button
        type="button"
        @click="emit('open-simulator')"
        class="inline-flex items-center gap-1.5 rounded-lg bg-brand-dark px-4 py-2 text-xs font-bold text-white hover:bg-brand transition cursor-pointer"
      >
        <Sparkles class="h-4 w-4" />
        <span>Lancer une série d'entraînement libre</span>
      </button>
    </div>

    <!-- Grille des formulaires -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div
        v-for="qf in studentQuizzes"
        :key="qf.id"
        class="rounded-lg border border-border bg-card overflow-hidden flex flex-col justify-between shadow-xs hover:border-brand/60 transition"
      >
        <div>
          <div class="h-2 w-full bg-brand"></div>
          <div class="p-4 space-y-2.5">
            <div class="flex items-center justify-between gap-2">
              <span
                class="rounded bg-brand-light px-2 py-0.5 text-[10px] font-bold uppercase text-brand-dark border border-brand-border"
              >
                Permis {{ qf.categorie_permis }}
              </span>

              <span
                v-if="getMyBestSubmission(qf.id)"
                class="inline-flex items-center gap-1 rounded-full bg-emerald-500/15 px-2.5 py-0.5 text-[10px] font-bold text-emerald-700 dark:text-emerald-300"
              >
                <CheckCircle2 class="h-3 w-3" />
                Déjà répondu ({{ getMyBestSubmission(qf.id)?.pourcentage }}%)
              </span>
              <span
                v-else
                class="inline-flex items-center gap-1 rounded-full bg-amber-500/15 px-2.5 py-0.5 text-[10px] font-bold text-amber-700 dark:text-amber-300"
              >
                Nouveau Quiz à compléter
              </span>
            </div>

            <h3 class="text-sm font-bold text-foreground">{{ qf.titre }}</h3>

            <p
              v-if="qf.description"
              class="text-xs text-muted-foreground line-clamp-2"
            >
              {{ qf.description }}
            </p>

            <div
              class="flex flex-wrap items-center gap-3 pt-2 text-[11px] text-muted-foreground"
            >
              <span class="flex items-center gap-1">
                <FileCheck2 class="h-3.5 w-3.5 text-brand" />
                {{ qf.questions?.length || 0 }} questions
              </span>
              <span class="flex items-center gap-1">
                <Clock class="h-3.5 w-3.5" />
                {{ qf.duree_minutes }} min
              </span>
              <span>Par {{ qf.moniteur_nom || "Votre Moniteur" }}</span>
            </div>
          </div>
        </div>

        <div
          class="border-t border-border bg-muted/20 px-4 py-3 flex items-center justify-between"
        >
          <span class="text-[11px] font-medium text-muted-foreground">
            Seuil requis :
            <strong class="text-foreground"
              >{{ qf.seuil_reussite_pct }}%</strong
            >
          </span>

          <button
            type="button"
            @click="emit('start', qf)"
            class="inline-flex items-center gap-1.5 rounded-md bg-brand-dark px-3.5 py-1.5 text-xs font-bold text-white hover:bg-brand transition cursor-pointer"
          >
            <Send class="h-3.5 w-3.5" />
            <span>
              {{
                getMyBestSubmission(qf.id)
                  ? "Répondre à nouveau"
                  : "Répondre au Quiz"
              }}
            </span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>


<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import { useGuichetStore } from "../../stores/guichetStore";
import { useQuizFormsStore } from "../../stores/quizFormsStore";
import { realtimeSocket } from "../../services/realtimeSocketService";
import type { QuizSoumissionEleve } from "../../types";
import { Brain, Plus, FileCheck2, Radio, ListChecks, X } from "lucide-vue-next";
import QuizBuilderTab from "./moniteur/QuizBuilderTab.vue";
import QuizPublishedTab from "./moniteur/QuizPublishedTab.vue";
import QuizResponsesTab from "./moniteur/QuizResponsesTab.vue";
import QuizSubmissionDetailModal from "./moniteur/QuizSubmissionDetailModal.vue";

const guichetStore = useGuichetStore();
const quizStore = useQuizFormsStore();

const isWsConnected = realtimeSocket.isConnected;
let unsubscribeWs: (() => void) | null = null;

const activeSubTab = ref<"builder" | "responses" | "published">("builder");
const selectedQuizFilter = ref("ALL");
const inspectedSubmission = ref<QuizSoumissionEleve | null>(null);

const viewResponsesForQuiz = (quizId: string) => {
  selectedQuizFilter.value = quizId;
  activeSubTab.value = "responses";
};

onMounted(async () => {
  realtimeSocket.connect();
  await quizStore.loadAll();

  unsubscribeWs = realtimeSocket.subscribe(async (msg) => {
    if (msg.type === "QUIZ_SUBMISSION_RECEIVED") {
      if (msg.data?.id) {
        quizStore.registerIncomingSubmission(msg.data);
        guichetStore.showToast(
          `Nouvelle réponse reçue de ${msg.data.candidat_nom} : ${msg.data.score_obtenu}/${msg.data.total_points} (${msg.data.pourcentage}%)`,
          "success",
        );
      } else {
        await quizStore.loadSoumissions();
      }
    } else if (msg.type === "QUIZ_FORM_UPDATED") {
      await quizStore.loadFormulaires();
    }
  });
});

onUnmounted(() => {
  if (unsubscribeWs) {
    unsubscribeWs();
    unsubscribeWs = null;
  }
});
</script>

<template>
  <div class="space-y-4">
    <!-- En-tête & Navigation -->
    <div
      class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between border-b border-border pb-3"
    >
      <div>
        <div class="flex items-center gap-2">
          <Brain class="h-4 w-4 text-brand" />
          <h1 class="text-base font-bold text-foreground">
            Formulaires de Quiz & Réponses Élèves (Style Google Forms)
          </h1>
          <span
            class="inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-bold"
            :class="
              isWsConnected
                ? 'bg-emerald-500/15 text-emerald-700 dark:text-emerald-300'
                : 'bg-amber-500/15 text-amber-700'
            "
          >
            <Radio class="h-3 w-3 animate-pulse" />
            {{ isWsConnected ? "Temps réel WS" : "Sync locale" }}
          </span>
        </div>
        <p class="text-xs text-muted-foreground mt-0.5">
          Créez vos formulaires de quiz interactifs, envoyez-les instantanément
          aux élèves et consultez les copies corrigées de chaque élève en
          direct.
        </p>
      </div>

      <div
        class="inline-flex items-center rounded-lg border border-border bg-muted/40 p-1 gap-1 self-start"
      >
        <button
          type="button"
          @click="activeSubTab = 'builder'"
          class="flex items-center gap-1.5 rounded-md px-3 py-1.5 text-xs font-semibold transition cursor-pointer"
          :class="
            activeSubTab === 'builder'
              ? 'bg-card text-brand-dark shadow-xs border border-border'
              : 'text-muted-foreground hover:text-foreground'
          "
        >
          <Plus class="h-3.5 w-3.5" />
          <span>Créer un Quiz</span>
        </button>

        <button
          type="button"
          @click="activeSubTab = 'published'"
          class="flex items-center gap-1.5 rounded-md px-3 py-1.5 text-xs font-semibold transition cursor-pointer"
          :class="
            activeSubTab === 'published'
              ? 'bg-card text-brand-dark shadow-xs border border-border'
              : 'text-muted-foreground hover:text-foreground'
          "
        >
          <ListChecks class="h-3.5 w-3.5" />
          <span>Quiz Publiés</span>
          <span
            class="rounded-full bg-brand-light px-1.5 py-0.2 text-[10px] font-bold text-brand-dark"
          >
            {{ quizStore.formulaires.length }}
          </span>
        </button>

        <button
          type="button"
          @click="activeSubTab = 'responses'"
          class="flex items-center gap-1.5 rounded-md px-3 py-1.5 text-xs font-semibold transition cursor-pointer"
          :class="
            activeSubTab === 'responses'
              ? 'bg-card text-brand-dark shadow-xs border border-border'
              : 'text-muted-foreground hover:text-foreground'
          "
        >
          <FileCheck2 class="h-3.5 w-3.5" />
          <span>Réponses des Élèves</span>
          <span
            class="rounded-full px-1.5 py-0.2 text-[10px] font-bold"
            :class="
              quizStore.soumissions.length > 0
                ? 'bg-emerald-500 text-white'
                : 'bg-muted text-muted-foreground'
            "
          >
            {{ quizStore.soumissions.length }}
          </span>
        </button>
      </div>
    </div>

    <!-- Bannière notification en direct -->
    <div
      v-if="quizStore.lastReceivedSubmission"
      class="rounded-[5px] border border-emerald-500/40 bg-emerald-500/10 px-4 py-2.5 flex flex-wrap items-center justify-between gap-2 text-xs"
    >
      <div class="flex items-center gap-2">
        <span
          class="inline-flex items-center gap-1 rounded-full bg-emerald-600 px-2 py-0.5 text-[10px] font-bold text-white"
        >
          <Radio class="h-3 w-3 animate-pulse" />
          Nouvelle copie reçue en direct
        </span>
        <span class="font-bold text-foreground">
          {{ quizStore.lastReceivedSubmission.candidat_nom }}
        </span>
        <span class="text-muted-foreground">
          a répondu à « {{ quizStore.lastReceivedSubmission.quiz_titre }} » —
          Score :
          <strong class="text-foreground">
            {{ quizStore.lastReceivedSubmission.score_obtenu }} /
            {{ quizStore.lastReceivedSubmission.total_points }} pts ({{
              quizStore.lastReceivedSubmission.pourcentage
            }}%)
          </strong>
        </span>
      </div>
      <div class="flex items-center gap-2">
        <button
          type="button"
          @click="
            inspectedSubmission = quizStore.lastReceivedSubmission;
            activeSubTab = 'responses';
          "
          class="rounded bg-emerald-600 px-2.5 py-1 text-[11px] font-bold text-white hover:bg-emerald-700 transition cursor-pointer"
        >
          Voir sa copie
        </button>
        <button
          type="button"
          @click="quizStore.lastReceivedSubmission = null"
          class="text-muted-foreground hover:text-foreground cursor-pointer"
        >
          <X class="h-3.5 w-3.5" />
        </button>
      </div>
    </div>

    <QuizBuilderTab
      v-if="activeSubTab === 'builder'"
      @published="activeSubTab = 'published'"
    />

    <QuizPublishedTab
      v-else-if="activeSubTab === 'published'"
      @create-new="activeSubTab = 'builder'"
      @view-responses="viewResponsesForQuiz"
    />

    <QuizResponsesTab
      v-else-if="activeSubTab === 'responses'"
      v-model:selected-quiz-filter="selectedQuizFilter"
      @inspect="(sub) => (inspectedSubmission = sub)"
    />

    <QuizSubmissionDetailModal
      :submission="inspectedSubmission"
      @close="inspectedSubmission = null"
    />
  </div>
</template>

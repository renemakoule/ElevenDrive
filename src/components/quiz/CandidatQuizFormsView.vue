<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted } from "vue";
import { useQuizFormsStore } from "../../stores/quizFormsStore";
import { useGuichetStore } from "../../stores/guichetStore";
import { realtimeSocket } from "../../services/realtimeSocketService";
import CandidatQuizList from "./candidat/CandidatQuizList.vue";
import CandidatQuizResultSheet from "./candidat/CandidatQuizResultSheet.vue";
import type {
  QuizFormulaire,
  QuizQuestionItem,
  QuizReponseDetailItem,
  QuizSoumissionEleve,
} from "../../types";
import { CheckCircle2, Send, Clock, ArrowLeft } from "lucide-vue-next";

const props = defineProps<{
  candidatId: string;
  candidatNom: string;
  numeroDossier: string;
  studentCategory: string;
}>();

const emit = defineEmits<{
  (e: "open-simulator"): void;
  (e: "score-saved", scoreSur40: number): void;
}>();

const quizStore = useQuizFormsStore();
const guichetStore = useGuichetStore();
const isWsConnected = realtimeSocket.isConnected;
let unsubscribeWs: (() => void) | null = null;

const activeQuiz = ref<QuizFormulaire | null>(null);
const lastResult = ref<QuizSoumissionEleve | null>(null);
const startTimeMs = ref<number>(0);
const elapsedSeconds = ref<number>(0);
let timerHandle: ReturnType<typeof setInterval> | null = null;

const selectedChoices = reactive<Record<string, string[]>>({});
const shortAnswers = reactive<Record<string, string>>({});

const studentQuizzes = computed(() =>
  quizStore.getQuizForStudent(props.studentCategory),
);

const mySubmissions = computed(() =>
  quizStore.getSoumissionsByCandidat(props.candidatId),
);

const answeredQuestionsCount = computed(() => {
  if (!activeQuiz.value) return 0;
  return activeQuiz.value.questions.filter((q) => {
    if (q.type_question === "REPONSE_COURTE") {
      return Boolean(shortAnswers[q.id]?.trim());
    }
    return (selectedChoices[q.id]?.length || 0) > 0;
  }).length;
});

const startQuiz = (quiz: QuizFormulaire) => {
  activeQuiz.value = quiz;
  lastResult.value = null;
  Object.keys(selectedChoices).forEach((k) => delete selectedChoices[k]);
  Object.keys(shortAnswers).forEach((k) => delete shortAnswers[k]);

  quiz.questions.forEach((q) => {
    selectedChoices[q.id] = [];
    shortAnswers[q.id] = "";
  });

  startTimeMs.value = Date.now();
  elapsedSeconds.value = 0;
  if (timerHandle) clearInterval(timerHandle);
  timerHandle = setInterval(() => {
    elapsedSeconds.value = Math.floor((Date.now() - startTimeMs.value) / 1000);
  }, 1000);
};

const closeQuiz = () => {
  activeQuiz.value = null;
  lastResult.value = null;
  if (timerHandle) {
    clearInterval(timerHandle);
    timerHandle = null;
  }
};

const toggleChoice = (q: QuizQuestionItem, optId: string) => {
  if (!selectedChoices[q.id]) {
    selectedChoices[q.id] = [];
  }
  if (q.type_question === "CHOIX_UNIQUE" || q.type_question === "VRAI_FAUX") {
    selectedChoices[q.id] = [optId];
  } else {
    const idx = selectedChoices[q.id].indexOf(optId);
    if (idx === -1) {
      selectedChoices[q.id].push(optId);
    } else {
      selectedChoices[q.id].splice(idx, 1);
    }
  }
};

const submitQuizForm = async () => {
  const quiz = activeQuiz.value;
  if (!quiz) return;

  for (let i = 0; i < quiz.questions.length; i++) {
    const q = quiz.questions[i];
    if (!q.obligatoire) continue;
    if (q.type_question === "REPONSE_COURTE") {
      if (!shortAnswers[q.id]?.trim()) {
        guichetStore.showToast(
          `Veuillez répondre à la question obligatoire N°${i + 1}.`,
          "error",
        );
        return;
      }
    } else if (!selectedChoices[q.id] || selectedChoices[q.id].length === 0) {
      guichetStore.showToast(
        `Veuillez sélectionner une réponse pour la question obligatoire N°${i + 1}.`,
        "error",
      );
      return;
    }
  }

  if (timerHandle) {
    clearInterval(timerHandle);
    timerHandle = null;
  }

  let scoreObtenu = 0;
  let totalPoints = 0;
  const details: QuizReponseDetailItem[] = [];

  for (const q of quiz.questions) {
    const ptsMax = Number(q.points) || 1;
    totalPoints += ptsMax;

    if (q.type_question === "REPONSE_COURTE") {
      const saisie = (shortAnswers[q.id] || "").trim();
      const attendue = (q.reponse_attendue || "").trim();
      const isOk =
        saisie.length > 0 &&
        attendue.length > 0 &&
        saisie.toLowerCase() === attendue.toLowerCase();
      const pts = isOk ? ptsMax : 0;
      scoreObtenu += pts;

      details.push({
        question_id: q.id,
        intitule_question: q.intitule,
        type_question: q.type_question,
        choix_ids: [],
        choix_labels: [],
        texte_reponse: saisie,
        est_correcte: isOk,
        points_obtenus: pts,
        points_max: ptsMax,
        bonne_reponse_label: attendue,
        explication: q.explication || "",
      });
    } else {
      const chosenIds = selectedChoices[q.id] || [];
      const correctOpts = q.options.filter((o) => o.est_correcte);
      const correctIds = correctOpts.map((o) => o.id);
      const chosenLabels = q.options
        .filter((o) => chosenIds.includes(o.id))
        .map((o) => o.texte);
      const bonneLabel = correctOpts.map((o) => o.texte).join(" · ");

      const isOk =
        chosenIds.length === correctIds.length &&
        correctIds.every((id) => chosenIds.includes(id));
      const pts = isOk ? ptsMax : 0;
      scoreObtenu += pts;

      details.push({
        question_id: q.id,
        intitule_question: q.intitule,
        type_question: q.type_question,
        choix_ids: chosenIds,
        choix_labels: chosenLabels,
        texte_reponse: chosenLabels.join(", "),
        est_correcte: isOk,
        points_obtenus: pts,
        points_max: ptsMax,
        bonne_reponse_label: bonneLabel,
        explication: q.explication || "",
      });
    }
  }

  const pourcentage =
    totalPoints > 0 ? Math.round((scoreObtenu / totalPoints) * 100) : 0;
  const estReussi = pourcentage >= (quiz.seuil_reussite_pct || 80);
  const dureeSec = Math.max(
    1,
    Math.floor((Date.now() - startTimeMs.value) / 1000),
  );

  const submitted = await quizStore.soumettreReponseEleve({
    quiz_id: quiz.id,
    quiz_titre: quiz.titre,
    candidat_id: props.candidatId,
    candidat_nom: props.candidatNom,
    numero_dossier: props.numeroDossier,
    categorie_permis: props.studentCategory,
    score_obtenu: scoreObtenu,
    total_points: totalPoints,
    pourcentage,
    est_reussi: estReussi,
    duree_secondes: dureeSec,
    details,
  });

  lastResult.value = submitted;
  emit("score-saved", Math.round((pourcentage * 40) / 100));
};

const formatTimer = (sec: number) => {
  const m = Math.floor(sec / 60)
    .toString()
    .padStart(2, "0");
  const s = (sec % 60).toString().padStart(2, "0");
  return `${m}:${s}`;
};

onMounted(async () => {
  realtimeSocket.connect();
  await quizStore.loadAll();

  unsubscribeWs = realtimeSocket.subscribe(async (msg) => {
    if (msg.type === "QUIZ_FORM_UPDATED") {
      await quizStore.loadFormulaires();
      guichetStore.showToast(
        "Nouveau formulaire de Quiz publié par votre moniteur !",
        "info",
      );
    }
  });
});

onUnmounted(() => {
  if (timerHandle) clearInterval(timerHandle);
  if (unsubscribeWs) {
    unsubscribeWs();
    unsubscribeWs = null;
  }
});
</script>

<template>
  <div class="space-y-4">
    <CandidatQuizList
      v-if="!activeQuiz"
      :student-category="studentCategory"
      :is-ws-connected="isWsConnected"
      :student-quizzes="studentQuizzes"
      :my-submissions="mySubmissions"
      @start="startQuiz"
      @open-simulator="emit('open-simulator')"
    />

    <template v-else>
      <div class="max-w-3xl mx-auto space-y-4">
        <div
          class="flex flex-wrap items-center justify-between gap-2 rounded-lg border border-border bg-card px-4 py-2.5 shadow-xs"
        >
          <button
            type="button"
            @click="closeQuiz"
            class="inline-flex items-center gap-1.5 text-xs font-semibold text-muted-foreground hover:text-foreground cursor-pointer"
          >
            <ArrowLeft class="h-4 w-4" />
            <span>Retour aux formulaires de quiz</span>
          </button>

          <div class="flex items-center gap-3 text-xs">
            <span class="text-muted-foreground">
              Répondu :
              <strong class="text-foreground">
                {{ answeredQuestionsCount }} / {{ activeQuiz.questions.length }}
              </strong>
            </span>
            <span
              class="inline-flex items-center gap-1 rounded-full bg-muted px-2.5 py-1 font-mono font-bold text-foreground"
            >
              <Clock class="h-3.5 w-3.5 text-brand" />
              {{ formatTimer(elapsedSeconds) }}
            </span>
          </div>
        </div>

        <div
          class="rounded-lg border border-border bg-card overflow-hidden shadow-xs"
        >
          <div class="h-2.5 w-full bg-brand"></div>
          <div class="p-5 space-y-2">
            <div class="flex items-center justify-between">
              <span
                class="rounded bg-brand-light px-2 py-0.5 text-[10px] font-bold uppercase text-brand-dark"
              >
                Permis {{ activeQuiz.categorie_permis }} · Envoyé par
                {{ activeQuiz.moniteur_nom }}
              </span>
              <span class="text-xs font-semibold text-muted-foreground">
                Candidat :
                <strong class="text-foreground">{{ candidatNom }}</strong>
              </span>
            </div>
            <h2 class="text-base sm:text-lg font-extrabold text-foreground">
              {{ activeQuiz.titre }}
            </h2>
            <p
              v-if="activeQuiz.description"
              class="text-xs text-muted-foreground"
            >
              {{ activeQuiz.description }}
            </p>
          </div>
        </div>

        <CandidatQuizResultSheet
          v-if="lastResult"
          :result="lastResult"
          @restart="startQuiz(activeQuiz)"
          @close="closeQuiz"
        />

        <template v-else>
          <div
            v-for="(q, qIdx) in activeQuiz.questions"
            :key="q.id"
            class="rounded-lg border border-border bg-card p-5 space-y-3 shadow-xs"
          >
            <div class="flex items-start justify-between gap-3">
              <div class="text-xs sm:text-sm font-bold text-foreground">
                <span class="text-brand-dark mr-1">{{ qIdx + 1 }}.</span>
                {{ q.intitule }}
                <span v-if="q.obligatoire" class="text-rose-500 ml-0.5">*</span>
              </div>
              <span
                class="shrink-0 rounded bg-muted px-2 py-0.5 text-[10px] font-semibold text-muted-foreground"
              >
                {{ q.points }} pt{{ q.points > 1 ? "s" : "" }}
              </span>
            </div>

            <div
              v-if="q.type_question !== 'REPONSE_COURTE'"
              class="space-y-2 pt-1"
            >
              <p
                v-if="q.type_question === 'CHOIX_MULTIPLE'"
                class="text-[11px] text-muted-foreground"
              >
                Plusieurs réponses possibles (cochez toutes les bonnes réponses)
                :
              </p>

              <button
                v-for="opt in q.options"
                :key="opt.id"
                type="button"
                @click="toggleChoice(q, opt.id)"
                class="w-full flex items-center gap-3 rounded-lg border px-3.5 py-2.5 text-left text-xs font-medium transition cursor-pointer"
                :class="
                  selectedChoices[q.id]?.includes(opt.id)
                    ? 'border-brand bg-brand-light/40 text-foreground font-bold'
                    : 'border-border bg-background text-foreground hover:bg-muted/40'
                "
              >
                <span
                  class="flex h-4 w-4 shrink-0 items-center justify-center border transition"
                  :class="[
                    q.type_question === 'CHOIX_MULTIPLE'
                      ? 'rounded-xs'
                      : 'rounded-full',
                    selectedChoices[q.id]?.includes(opt.id)
                      ? 'border-brand bg-brand text-white'
                      : 'border-muted-foreground/50',
                  ]"
                >
                  <CheckCircle2
                    v-if="selectedChoices[q.id]?.includes(opt.id)"
                    class="h-3 w-3"
                  />
                </span>
                <span>{{ opt.texte }}</span>
              </button>
            </div>

            <div v-else class="pt-1">
              <input
                v-model="shortAnswers[q.id]"
                type="text"
                placeholder="Votre réponse..."
                class="w-full border-b-2 border-border bg-transparent px-1 py-2 text-xs sm:text-sm font-semibold text-foreground focus:border-brand focus:outline-none"
              />
            </div>
          </div>

          <div class="flex items-center justify-between pt-2">
            <button
              type="button"
              @click="closeQuiz"
              class="rounded-lg border border-border bg-card px-4 py-2 text-xs font-semibold text-muted-foreground hover:text-foreground cursor-pointer"
            >
              Annuler
            </button>

            <button
              type="button"
              @click="submitQuizForm"
              :disabled="quizStore.submitting"
              class="inline-flex items-center gap-2 rounded-lg bg-brand-dark px-6 py-2.5 text-xs font-bold text-white hover:bg-brand transition shadow-sm disabled:opacity-50 cursor-pointer"
            >
              <Send class="h-4 w-4" />
              <span>
                {{
                  quizStore.submitting
                    ? "Envoi en cours..."
                    : "Envoyer mes réponses au Moniteur"
                }}
              </span>
            </button>
          </div>
        </template>
      </div>
    </template>
  </div>
</template>

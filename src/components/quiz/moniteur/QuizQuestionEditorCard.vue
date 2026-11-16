<script setup lang="ts">
import type { QuizQuestionItem, QuizQuestionType } from "../../../types";
import { useGuichetStore } from "../../../stores/guichetStore";
import {
  CheckCircle2,
  X,
  Plus,
  ArrowUp,
  ArrowDown,
  Copy,
  Trash2,
} from "lucide-vue-next";

const props = defineProps<{
  question: QuizQuestionItem;
  index: number;
  totalQuestions: number;
}>();

const emit = defineEmits<{
  (e: "move", dir: -1 | 1): void;
  (e: "duplicate"): void;
  (e: "remove"): void;
}>();

const guichetStore = useGuichetStore();

const onQuestionTypeChange = (newType: QuizQuestionType) => {
  const q = props.question;
  q.type_question = newType;
  if (newType === "VRAI_FAUX") {
    q.options = [
      { id: `vf-1-${Date.now()}`, texte: "Vrai", est_correcte: true },
      { id: `vf-2-${Date.now()}`, texte: "Faux", est_correcte: false },
    ];
  } else if (newType === "REPONSE_COURTE") {
    q.options = [];
    if (!q.reponse_attendue) q.reponse_attendue = "";
  } else if (!q.options || q.options.length === 0) {
    q.options = [
      { id: `opt-1-${Date.now()}`, texte: "Option 1", est_correcte: true },
      { id: `opt-2-${Date.now()}`, texte: "Option 2", est_correcte: false },
    ];
  } else if (newType === "CHOIX_UNIQUE") {
    let foundCorrect = false;
    q.options.forEach((opt) => {
      if (opt.est_correcte && !foundCorrect) foundCorrect = true;
      else opt.est_correcte = false;
    });
    if (!foundCorrect && q.options[0]) q.options[0].est_correcte = true;
  }
};

const addOption = () => {
  props.question.options.push({
    id: `opt-${Date.now()}-${props.question.options.length + 1}`,
    texte: `Option ${props.question.options.length + 1}`,
    est_correcte: false,
  });
};

const removeOption = (optIdx: number) => {
  const q = props.question;
  if (q.options.length <= 2) {
    guichetStore.showToast(
      "Une question à choix doit comporter au moins 2 propositions.",
      "error",
    );
    return;
  }
  const removedWasCorrect = q.options[optIdx]?.est_correcte;
  q.options.splice(optIdx, 1);
  if (
    removedWasCorrect &&
    !q.options.some((o) => o.est_correcte) &&
    q.options[0]
  ) {
    q.options[0].est_correcte = true;
  }
};

const toggleCorrectOption = (optIdx: number) => {
  const q = props.question;
  if (q.type_question === "CHOIX_UNIQUE" || q.type_question === "VRAI_FAUX") {
    q.options.forEach((o, i) => {
      o.est_correcte = i === optIdx;
    });
  } else {
    const target = q.options[optIdx];
    if (target) target.est_correcte = !target.est_correcte;
  }
};
</script>

<template>
  <div
    class="rounded-lg border border-border bg-card shadow-xs overflow-hidden transition hover:border-brand/50"
  >
    <div class="flex">
      <div class="w-1.5 bg-brand/80 shrink-0"></div>

      <div class="flex-1 p-5 space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-12 gap-3 items-start">
          <div class="md:col-span-8">
            <div class="flex items-center gap-2 mb-1">
              <span
                class="inline-flex items-center justify-center rounded bg-brand-light px-2 py-0.5 text-[10px] font-extrabold text-brand-dark"
              >
                Question {{ index + 1 }}
              </span>
              <span
                v-if="question.obligatoire"
                class="text-[10px] font-bold text-rose-500"
              >
                * Obligatoire
              </span>
            </div>
            <input
              v-model="question.intitule"
              type="text"
              placeholder="Saisissez l'énoncé de la question ici..."
              class="w-full rounded-md border border-border bg-muted/20 px-3 py-2 text-xs sm:text-sm font-semibold text-foreground focus:border-brand focus:bg-background focus:outline-none"
            />
          </div>

          <div class="md:col-span-4">
            <label
              class="block text-[10px] font-semibold uppercase text-muted-foreground mb-1"
            >
              Type de réponse
            </label>
            <select
              :value="question.type_question"
              @change="
                onQuestionTypeChange(
                  ($event.target as HTMLSelectElement).value as QuizQuestionType,
                )
              "
              class="w-full rounded-md border border-border bg-background px-2.5 py-2 text-xs font-semibold text-foreground focus:border-brand focus:outline-none"
            >
              <option value="CHOIX_UNIQUE">🔘 Choix unique (QCM Radio)</option>
              <option value="CHOIX_MULTIPLE">
                ☑️ Cases à cocher (Choix multiple)
              </option>
              <option value="VRAI_FAUX">⚖️ Vrai / Faux</option>
              <option value="REPONSE_COURTE">
                ✍️ Réponse courte (Texte/Nombre)
              </option>
            </select>
          </div>
        </div>

        <!-- Propositions selon le type -->
        <div
          v-if="question.type_question !== 'REPONSE_COURTE'"
          class="space-y-2 pl-1"
        >
          <div
            class="flex items-center justify-between text-[11px] text-muted-foreground"
          >
            <span>
              Propositions de réponse (cliquez sur le bouton à gauche pour
              définir la
              <strong class="text-emerald-600 dark:text-emerald-400"
                >bonne réponse</strong
              >) :
            </span>
          </div>

          <div
            v-for="(opt, optIdx) in question.options"
            :key="opt.id"
            class="flex items-center gap-2.5 rounded-md border px-3 py-2 transition"
            :class="
              opt.est_correcte
                ? 'border-emerald-500/50 bg-emerald-500/5'
                : 'border-border bg-background'
            "
          >
            <button
              type="button"
              @click="toggleCorrectOption(optIdx)"
              class="flex h-5 w-5 shrink-0 items-center justify-center rounded transition cursor-pointer"
              :class="[
                question.type_question === 'CHOIX_MULTIPLE'
                  ? 'rounded-sm'
                  : 'rounded-full',
                opt.est_correcte
                  ? 'bg-emerald-600 text-white'
                  : 'border-2 border-muted-foreground/40 hover:border-brand',
              ]"
              :title="
                opt.est_correcte
                  ? 'Bonne réponse sélectionnée'
                  : 'Cliquer pour définir comme bonne réponse'
              "
            >
              <CheckCircle2 v-if="opt.est_correcte" class="h-3.5 w-3.5" />
            </button>

            <input
              v-model="opt.texte"
              type="text"
              :readonly="question.type_question === 'VRAI_FAUX'"
              placeholder="Texte de la proposition..."
              class="flex-1 bg-transparent text-xs font-medium text-foreground focus:outline-none"
            />

            <span
              v-if="opt.est_correcte"
              class="rounded-full bg-emerald-500/15 px-2 py-0.5 text-[10px] font-bold text-emerald-700 dark:text-emerald-300 shrink-0"
            >
              Bonne réponse
            </span>

            <button
              v-if="
                question.type_question !== 'VRAI_FAUX' &&
                question.options.length > 2
              "
              type="button"
              @click="removeOption(optIdx)"
              class="p-1 text-muted-foreground hover:text-rose-500 transition cursor-pointer"
              title="Retirer cette option"
            >
              <X class="h-3.5 w-3.5" />
            </button>
          </div>

          <button
            v-if="question.type_question !== 'VRAI_FAUX'"
            type="button"
            @click="addOption"
            class="inline-flex items-center gap-1.5 text-xs font-semibold text-brand-dark hover:underline pt-1 cursor-pointer"
          >
            <Plus class="h-3.5 w-3.5" />
            Ajouter une option
          </button>
        </div>

        <!-- Réponse courte -->
        <div
          v-else
          class="rounded-md border border-emerald-500/40 bg-emerald-500/5 p-3 space-y-1.5"
        >
          <label
            class="block text-[11px] font-bold text-emerald-700 dark:text-emerald-300"
          >
            Réponse exacte attendue (Clé de correction automatique) *
          </label>
          <input
            v-model="question.reponse_attendue"
            type="text"
            placeholder="Ex : 80 ou priorité à droite"
            class="w-full rounded border border-border bg-background px-3 py-1.5 text-xs font-semibold text-foreground focus:border-brand focus:outline-none"
          />
          <p class="text-[10px] text-muted-foreground">
            La réponse de l'élève sera comparée (sans distinction de
            majuscules/minuscules).
          </p>
        </div>

        <div>
          <input
            v-model="question.explication"
            type="text"
            placeholder="💡 Explication pédagogique affichée à l'élève après correction (optionnel)..."
            class="w-full rounded-md border border-dashed border-border bg-muted/15 px-3 py-1.5 text-xs text-muted-foreground focus:border-brand focus:text-foreground focus:outline-none"
          />
        </div>

        <!-- Pied de carte Question -->
        <div
          class="flex flex-wrap items-center justify-between gap-2 pt-3 border-t border-border text-xs"
        >
          <div class="flex items-center gap-3">
            <div class="flex items-center gap-1.5">
              <span class="text-muted-foreground">Barème :</span>
              <select
                v-model.number="question.points"
                class="rounded border border-border bg-background px-2 py-1 text-xs font-bold text-foreground"
              >
                <option :value="1">1 point</option>
                <option :value="2">2 points</option>
                <option :value="3">3 points</option>
                <option :value="4">4 points</option>
                <option :value="5">5 points</option>
              </select>
            </div>

            <label class="inline-flex items-center gap-1.5 cursor-pointer select-none">
              <input
                v-model="question.obligatoire"
                type="checkbox"
                class="rounded border-border text-brand focus:ring-brand"
              />
              <span class="text-muted-foreground">Obligatoire</span>
            </label>
          </div>

          <div class="flex items-center gap-1">
            <button
              type="button"
              @click="emit('move', -1)"
              :disabled="index === 0"
              class="p-1.5 rounded hover:bg-muted text-muted-foreground hover:text-foreground disabled:opacity-30 cursor-pointer"
              title="Monter la question"
            >
              <ArrowUp class="h-3.5 w-3.5" />
            </button>
            <button
              type="button"
              @click="emit('move', 1)"
              :disabled="index === totalQuestions - 1"
              class="p-1.5 rounded hover:bg-muted text-muted-foreground hover:text-foreground disabled:opacity-30 cursor-pointer"
              title="Descendre la question"
            >
              <ArrowDown class="h-3.5 w-3.5" />
            </button>
            <button
              type="button"
              @click="emit('duplicate')"
              class="p-1.5 rounded hover:bg-muted text-muted-foreground hover:text-foreground cursor-pointer"
              title="Dupliquer cette question"
            >
              <Copy class="h-3.5 w-3.5" />
            </button>
            <button
              type="button"
              @click="emit('remove')"
              class="p-1.5 rounded hover:bg-rose-500/10 text-rose-500 cursor-pointer"
              title="Supprimer cette question"
            >
              <Trash2 class="h-3.5 w-3.5" />
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>


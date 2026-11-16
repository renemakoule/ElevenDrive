<script setup lang="ts">
import { ref, computed } from "vue";
import { useAuthStore } from "../../../stores/authStore";
import { useGuichetStore } from "../../../stores/guichetStore";
import { useQuizFormsStore } from "../../../stores/quizFormsStore";
import type { QuizQuestionItem } from "../../../types";
import QuizQuestionEditorCard from "./QuizQuestionEditorCard.vue";
import { Sparkles, Plus, Send } from "lucide-vue-next";

const emit = defineEmits<{
  (e: "published"): void;
}>();

const authStore = useAuthStore();
const guichetStore = useGuichetStore();
const quizStore = useQuizFormsStore();

const formTitre = ref("Quiz Code de la Route — Priorités & Signalisation");
const formDescription = ref(
  "Répondez attentivement à chaque question. Vos réponses sont transmises instantanément à votre moniteur.",
);
const formCategorie = ref("TOUTES");
const formDuree = ref(15);
const formSeuil = ref(80);

const createDefaultQuestion = (index: number): QuizQuestionItem => ({
  id: `q-${Date.now()}-${index}-${Math.random().toString(36).slice(2, 6)}`,
  intitule: "",
  type_question: "CHOIX_UNIQUE",
  points: 2,
  obligatoire: true,
  explication: "",
  reponse_attendue: "",
  options: [
    { id: `opt-1-${Date.now()}`, texte: "Option 1", est_correcte: true },
    { id: `opt-2-${Date.now()}`, texte: "Option 2", est_correcte: false },
    { id: `opt-3-${Date.now()}`, texte: "Option 3", est_correcte: false },
  ],
});

const questions = ref<QuizQuestionItem[]>([
  {
    id: "q-init-1",
    intitule:
      "À une intersection sans signalisation, quelle règle de priorité s'applique par défaut ?",
    type_question: "CHOIX_UNIQUE",
    points: 2,
    obligatoire: true,
    explication:
      "En l'absence de signalisation, la règle générale est la priorité à droite.",
    options: [
      { id: "o-1", texte: "La priorité à droite", est_correcte: true },
      {
        id: "o-2",
        texte: "La priorité au véhicule le plus rapide",
        est_correcte: false,
      },
      { id: "o-3", texte: "La priorité à gauche", est_correcte: false },
    ],
  },
]);

const totalPointsFormulaire = computed(() =>
  questions.value.reduce((acc, q) => acc + (Number(q.points) || 1), 0),
);

const addQuestion = () => {
  questions.value.push(createDefaultQuestion(questions.value.length + 1));
};

const duplicateQuestion = (idx: number) => {
  const src = questions.value[idx];
  if (!src) return;
  const copy: QuizQuestionItem = {
    ...src,
    id: `q-${Date.now()}-${Math.random().toString(36).slice(2, 6)}`,
    intitule: `${src.intitule} (Copie)`,
    options: src.options.map((o, i) => ({
      ...o,
      id: `opt-${Date.now()}-${i}-${Math.random().toString(36).slice(2, 5)}`,
    })),
  };
  questions.value.splice(idx + 1, 0, copy);
};

const removeQuestion = (idx: number) => {
  if (questions.value.length <= 1) {
    guichetStore.showToast(
      "Le formulaire doit contenir au moins une question.",
      "error",
    );
    return;
  }
  questions.value.splice(idx, 1);
};

const moveQuestion = (idx: number, dir: -1 | 1) => {
  const target = idx + dir;
  if (target < 0 || target >= questions.value.length) return;
  const temp = questions.value[idx];
  questions.value[idx] = questions.value[target];
  questions.value[target] = temp;
};

const loadTemplateQuiz = () => {
  formTitre.value =
    "Évaluation Code — Signalisation, Dépassement & Distances de Sécurité";
  formDescription.value =
    "Série officielle préparée par votre moniteur. Répondez à toutes les questions obligatoires puis validez votre copie.";
  formCategorie.value = "B";
  formDuree.value = 15;
  formSeuil.value = 75;
  questions.value = [
    {
      id: "tpl-1",
      intitule:
        "Quelle est la distance minimale latérale à laisser lors du dépassement d'un deux-roues hors agglomération ?",
      type_question: "CHOIX_UNIQUE",
      points: 2,
      obligatoire: true,
      explication:
        "Hors agglomération, l'écart latéral minimal est de 1,50 mètre (1 mètre en agglomération).",
      options: [
        { id: "t1-o1", texte: "1,00 mètre", est_correcte: false },
        { id: "t1-o2", texte: "1,50 mètre", est_correcte: true },
        { id: "t1-o3", texte: "0,50 mètre", est_correcte: false },
      ],
    },
    {
      id: "tpl-2",
      intitule:
        "Quels éléments le conducteur doit-il obligatoirement contrôler avant de s'engager dans un dépassement ?",
      type_question: "CHOIX_MULTIPLE",
      points: 3,
      obligatoire: true,
      explication:
        "Il faut vérifier la signalisation (ligne continue/panneau), l'absence de danger dans les rétroviseurs et la visibilité vers l'avant.",
      options: [
        {
          id: "t2-o1",
          texte: "Qu'aucun usager ne soit déjà en train de me dépasser",
          est_correcte: true,
        },
        {
          id: "t2-o2",
          texte: "Que la visibilité vers l'avant est suffisante et dégagée",
          est_correcte: true,
        },
        {
          id: "t2-o3",
          texte: "Que le véhicule devant moi accélère",
          est_correcte: false,
        },
      ],
    },
    {
      id: "tpl-3",
      intitule:
        "Un feu jaune fixe oblige le conducteur à s'arrêter s'il peut le faire en toute sécurité.",
      type_question: "VRAI_FAUX",
      points: 2,
      obligatoire: true,
      explication:
        "L'arrêt au feu jaune fixe est obligatoire, sauf si l'arrêt brusque présente un danger.",
      options: [
        { id: "t3-o1", texte: "Vrai", est_correcte: true },
        { id: "t3-o2", texte: "Faux", est_correcte: false },
      ],
    },
    {
      id: "tpl-4",
      intitule:
        "Sur route mouillée hors agglomération (limitée à 90 km/h par temps sec), quelle est la vitesse maximale autorisée (en km/h) ?",
      type_question: "REPONSE_COURTE",
      points: 3,
      obligatoire: true,
      reponse_attendue: "80",
      explication:
        "Par temps de pluie, la limitation de vitesse sur route passe de 90 km/h à 80 km/h.",
      options: [],
    },
  ];
  guichetStore.showToast(
    "Modèle de formulaire Google Forms chargé avec 4 types de questions.",
    "info",
  );
};

const publierQuiz = async () => {
  if (!formTitre.value.trim()) {
    guichetStore.showToast(
      "Veuillez saisir un titre pour votre formulaire de quiz.",
      "error",
    );
    return;
  }

  for (let i = 0; i < questions.value.length; i++) {
    const q = questions.value[i];
    if (!q.intitule.trim()) {
      guichetStore.showToast(
        `Veuillez renseigner l'intitulé de la question N°${i + 1}.`,
        "error",
      );
      return;
    }
    if (q.type_question === "REPONSE_COURTE") {
      if (!q.reponse_attendue?.trim()) {
        guichetStore.showToast(
          `Veuillez indiquer la réponse attendue pour la question N°${i + 1}.`,
          "error",
        );
        return;
      }
    } else if (!q.options.some((o) => o.est_correcte)) {
      guichetStore.showToast(
        `Veuillez cocher au moins une bonne réponse pour la question N°${i + 1}.`,
        "error",
      );
      return;
    }
  }

  const u = authStore.currentUser;
  const moniteurNom =
    `${u?.prenom || ""} ${u?.nom || ""}`.trim() || "Moniteur Principal";

  await quizStore.creerFormulaire({
    moniteur_id: u?.id || undefined,
    moniteur_nom: moniteurNom,
    titre: formTitre.value.trim(),
    description: formDescription.value.trim() || undefined,
    categorie_permis: formCategorie.value,
    duree_minutes: Number(formDuree.value) || 15,
    seuil_reussite_pct: Number(formSeuil.value) || 80,
    est_publie: true,
    questions: questions.value,
  });

  emit("published");
};
</script>

<template>
  <div class="max-w-4xl mx-auto space-y-4">
    <!-- Carte En-tête du Formulaire -->
    <div class="rounded-lg border border-border bg-card shadow-xs overflow-hidden">
      <div class="h-2.5 w-full bg-brand"></div>
      <div class="p-5 space-y-4">
        <div
          class="flex flex-col sm:flex-row sm:items-center justify-between gap-2"
        >
          <span
            class="inline-flex items-center gap-1.5 text-[11px] font-bold uppercase tracking-wider text-brand-dark"
          >
            <Sparkles class="h-3.5 w-3.5 text-brand" />
            Constructeur de Quiz Interactif · Diffusion Temps Réel
          </span>
          <button
            type="button"
            @click="loadTemplateQuiz"
            class="inline-flex items-center gap-1.5 rounded-md border border-brand-border bg-brand-light/50 px-2.5 py-1 text-[11px] font-semibold text-brand-dark hover:bg-brand-light transition cursor-pointer"
          >
            <Sparkles class="h-3 w-3 text-brand" />
            Pré-remplir avec un modèle Code de la Route
          </button>
        </div>

        <div>
          <label class="block text-[11px] font-semibold text-muted-foreground mb-1">
            Titre du Formulaire de Quiz *
          </label>
          <input
            v-model="formTitre"
            type="text"
            placeholder="Ex : Série N°4 — Priorités, Croisements et Dépassements"
            class="w-full border-b-2 border-border bg-transparent px-1 py-2 text-base font-bold text-foreground focus:border-brand focus:outline-none transition"
          />
        </div>

        <div>
          <label class="block text-[11px] font-semibold text-muted-foreground mb-1">
            Description & Consignes pour l'élève
          </label>
          <textarea
            v-model="formDescription"
            rows="2"
            placeholder="Instructions affichées en haut du formulaire côté élève..."
            class="w-full rounded-md border border-border bg-muted/20 px-3 py-2 text-xs text-foreground focus:border-brand focus:outline-none"
          ></textarea>
        </div>

        <div
          class="grid grid-cols-2 sm:grid-cols-4 gap-3 pt-2 border-t border-border"
        >
          <div>
            <label
              class="block text-[10px] font-semibold uppercase text-muted-foreground mb-1"
            >
              Catégorie Permis ciblée
            </label>
            <select
              v-model="formCategorie"
              class="w-full rounded-md border border-border bg-background px-2.5 py-1.5 text-xs font-semibold text-foreground focus:border-brand focus:outline-none"
            >
              <option value="TOUTES">Toutes catégories (A, B, C, D...)</option>
              <option value="B">Permis B (Véhicule léger)</option>
              <option value="A">Permis A / A1 (Moto)</option>
              <option value="C">Permis C (Poids lourd)</option>
              <option value="D">Permis D (Transport commun)</option>
              <option value="E">Permis E (Remorque)</option>
            </select>
          </div>

          <div>
            <label
              class="block text-[10px] font-semibold uppercase text-muted-foreground mb-1"
            >
              Durée impartie (min)
            </label>
            <input
              v-model.number="formDuree"
              type="number"
              min="2"
              max="120"
              class="w-full rounded-md border border-border bg-background px-2.5 py-1.5 text-xs font-semibold text-foreground focus:border-brand focus:outline-none"
            />
          </div>

          <div>
            <label
              class="block text-[10px] font-semibold uppercase text-muted-foreground mb-1"
            >
              Seuil de réussite (%)
            </label>
            <input
              v-model.number="formSeuil"
              type="number"
              min="50"
              max="100"
              class="w-full rounded-md border border-border bg-background px-2.5 py-1.5 text-xs font-semibold text-foreground focus:border-brand focus:outline-none"
            />
          </div>

          <div
            class="flex flex-col justify-center rounded-md bg-muted/40 px-3 py-1.5 border border-border"
          >
            <span class="text-[10px] text-muted-foreground">Barème total</span>
            <span class="text-xs font-extrabold text-brand-dark">
              {{ questions.length }} question(s) · {{ totalPointsFormulaire }} pts
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Cartes de Questions Dynamiques -->
    <QuizQuestionEditorCard
      v-for="(q, qIdx) in questions"
      :key="q.id"
      :question="q"
      :index="qIdx"
      :total-questions="questions.length"
      @move="(dir) => moveQuestion(qIdx, dir)"
      @duplicate="duplicateQuestion(qIdx)"
      @remove="removeQuestion(qIdx)"
    />

    <!-- Barre d'actions finale -->
    <div
      class="rounded-lg border border-border bg-card p-4 flex flex-col sm:flex-row items-center justify-between gap-3 shadow-xs"
    >
      <button
        type="button"
        @click="addQuestion"
        class="w-full sm:w-auto inline-flex items-center justify-center gap-2 rounded-lg border border-brand-border bg-brand-light px-4 py-2 text-xs font-bold text-brand-dark hover:bg-brand-light/80 transition cursor-pointer"
      >
        <Plus class="h-4 w-4 text-brand" />
        <span>Ajouter une nouvelle question</span>
      </button>

      <button
        type="button"
        @click="publierQuiz"
        :disabled="quizStore.submitting"
        class="w-full sm:w-auto inline-flex items-center justify-center gap-2 rounded-lg bg-brand-dark px-5 py-2.5 text-xs font-bold text-white hover:bg-brand transition shadow-sm disabled:opacity-50 cursor-pointer"
      >
        <Send class="h-4 w-4" />
        <span>
          {{
            quizStore.submitting
              ? "Publication en cours..."
              : "Publier & Envoyer aux Élèves en Temps Réel"
          }}
        </span>
      </button>
    </div>
  </div>
</template>


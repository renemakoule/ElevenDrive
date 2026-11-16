<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue';
import type { QuizQuestion } from '../../types';
import { speechAudioService } from '../../services/speechAudioService';
import { X, Check, ArrowRight, RotateCcw, Award, BookOpen, Volume2, VolumeX } from 'lucide-vue-next';

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'finish', score: number): void;
}>();

const QUESTIONS: QuizQuestion[] = [
  {
    id: 'q1',
    theme: 'Priorités et intersections',
    enonce: 'À une intersection sans aucune signalisation (ni panneaux ni feux), qui a la priorité ?',
    panneau_code: 'INTERSECTION',
    options: [
      { id: 'a', texte: 'Le véhicule qui arrive à ma gauche', est_correcte: false },
      { id: 'b', texte: 'Le véhicule qui arrive à ma droite (priorité à droite)', est_correcte: true },
      { id: 'c', texte: 'Le véhicule qui roule le plus vite', est_correcte: false },
    ],
    explication: "En l'absence de signalisation ou d'indication d'un agent, la règle est la priorité au véhicule venant de la droite.",
  },
  {
    id: 'q2',
    theme: 'Panneaux de signalisation',
    enonce: 'Que signifie un panneau rond à fond rouge traversé par une barre blanche horizontale ?',
    panneau_code: 'SENS_INTERDIT',
    options: [
      { id: 'a', texte: 'Accès interdit à tous les véhicules (sens interdit)', est_correcte: true },
      { id: 'b', texte: 'Arrêt obligatoire (stop)', est_correcte: false },
      { id: 'c', texte: 'Stationnement interdit', est_correcte: false },
    ],
    explication: "Ce panneau indique un sens interdit : l'accès est défendu à tout véhicule par cette extrémité de la chaussée.",
  },
  {
    id: 'q3',
    theme: 'Feux tricolores',
    enonce: 'Le feu passe au jaune fixe (orange). Que devez-vous faire ?',
    panneau_code: 'FEU_ORANGE',
    options: [
      { id: 'a', texte: 'Accélérer pour passer avant le feu rouge', est_correcte: false },
      { id: 'b', texte: "M'arrêter, sauf si je suis trop près et que le freinage présente un danger immédiat pour les usagers qui me suivent", est_correcte: true },
      { id: 'c', texte: 'Klaxonner et passer', est_correcte: false },
    ],
    explication: "Le feu jaune impose l'arrêt. La seule exception est lorsque le véhicule est tellement engagé qu'un arrêt brutal provoquerait une collision par l'arrière.",
  },
  {
    id: 'q4',
    theme: 'Vitesse et agglomération',
    enonce: 'En agglomération (zone urbaine / ville), quelle est la vitesse maximale légale autorisée ?',
    panneau_code: 'VITESSE_VILLE',
    options: [
      { id: 'a', texte: '80 km/h', est_correcte: false },
      { id: 'b', texte: '60 km/h', est_correcte: true },
      { id: 'c', texte: '100 km/h', est_correcte: false },
    ],
    explication: "En agglomération, la limitation générale de vitesse est fixée à 60 km/h pour garantir la sécurité des piétons.",
  },
  {
    id: 'q5',
    theme: 'Lignes et dépassement',
    enonce: 'Puis-je franchir ou chevaucher une ligne blanche continue pour dépasser un véhicule lent ?',
    panneau_code: 'LIGNE_CONTINUE',
    options: [
      { id: 'a', texte: 'Oui, si le véhicule devant moi roule à moins de 30 km/h', est_correcte: false },
      { id: 'b', texte: "Non, le franchissement ou le chevauchement d'une ligne continue est interdit", est_correcte: true },
      { id: 'c', texte: 'Oui, avec les feux de détresse allumés', est_correcte: false },
    ],
    explication: "La ligne continue ne se franchit pas. Tout chevauchement constitue une infraction grave.",
  },
  {
    id: 'q6',
    theme: 'Distances de sécurité',
    enonce: "Quelle règle simple permet d'évaluer la distance de sécurité minimale avec le véhicule qui précède ?",
    panneau_code: 'DISTANCE_SECURI',
    options: [
      { id: 'a', texte: "Laisser au moins 2 secondes d'intervalle entre les deux véhicules", est_correcte: true },
      { id: 'b', texte: "Laisser la longueur d'une roue", est_correcte: false },
      { id: 'c', texte: 'Coller au pare-chocs pour économiser le carburant', est_correcte: false },
    ],
    explication: "L'intervalle de sécurité correspond au temps de réaction (1 seconde) plus une marge de freinage (1 seconde), soit 2 secondes au minimum.",
  },
  {
    id: 'q7',
    theme: 'Alcool et vigilance',
    enonce: "Quel est l'effet de l'alcool ou d'une fatigue excessive sur le conducteur ?",
    panneau_code: 'SECURITE_ALCOOL',
    options: [
      { id: 'a', texte: 'Il améliore les réflexes et la vision de nuit', est_correcte: false },
      { id: 'b', texte: 'Il allonge le temps de réaction et rétrécit le champ visuel (effet tunnel)', est_correcte: true },
      { id: 'c', texte: "Il n'a aucun impact si la vitesse reste faible", est_correcte: false },
    ],
    explication: "L'alcool et la fatigue dégradent les capacités du conducteur : le temps de réaction passe de 1 seconde à plus de 2 secondes et la perception des distances est altérée.",
  },
  {
    id: 'q8',
    theme: 'Ronds-points et giratoires',
    enonce: 'Sur un carrefour à sens giratoire précédé d\'un panneau « Cédez le passage », qui a la priorité ?',
    panneau_code: 'ROND_POINT',
    options: [
      { id: 'a', texte: "Les véhicules qui s'engagent sur le rond-point", est_correcte: false },
      { id: 'b', texte: "Les véhicules déjà engagés sur l'anneau", est_correcte: true },
      { id: 'c', texte: 'Le véhicule le plus lourd', est_correcte: false },
    ],
    explication: "Sur un giratoire avec « Cédez le passage », la priorité appartient aux usagers déjà engagés sur l'anneau.",
  },
];

const LETTRES = ['A', 'B', 'C', 'D'];
const NOTE_MAX = 40;

function melanger<T>(liste: T[]): T[] {
  const copie = [...liste];
  for (let i = copie.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [copie[i], copie[j]] = [copie[j], copie[i]];
  }
  return copie;
}

const preparer = () => QUESTIONS.map((q) => ({ ...q, options: melanger(q.options) }));

const questions = ref<QuizQuestion[]>(preparer());
const index = ref(0);
const choix = ref<string | null>(null);
const corrige = ref(false);
const score = ref(0);
const termine = ref(false);

const total = computed(() => questions.value.length);
const question = computed(() => questions.value[index.value]);
const repondues = computed(() => index.value + (corrige.value ? 1 : 0));
const progression = computed(() => Math.round((repondues.value / total.value) * 100));
const note = computed(() => Math.round((score.value / total.value) * NOTE_MAX));
const pourcentage = computed(() => Math.round((score.value / total.value) * 100));
const dernier = computed(() => index.value === total.value - 1);

const appreciation = computed(() => {
  const ratio = score.value / total.value;
  if (ratio >= 0.85) return 'Très bon niveau : priorités et panneaux sont bien maîtrisés.';
  if (ratio >= 0.6) return "Bonne base. Refaites la série pour viser le sans-faute avant l'examen.";
  return 'Continuez à vous entraîner, en particulier sur la circulation et les feux.';
});

const modeAudio = ref(true);

const toggleModeAudio = () => {
  modeAudio.value = !modeAudio.value;
  if (!modeAudio.value) {
    speechAudioService.arreter();
  } else {
    lireQuestionCourante();
  }
};

const lireQuestionCourante = () => {
  if (!question.value || termine.value) return;
  const q = question.value;
  const optionsTexte = q.options
    .map((opt, i) => `Option ${LETTRES[i]} : ${opt.texte}`)
    .join('. ');
  const texte = `Question ${index.value + 1}. Thème : ${q.theme}. ${q.enonce}. ${optionsTexte}`;
  speechAudioService.parlerCourt(texte);
};

const fermerModal = () => {
  speechAudioService.arreter();
  emit('close');
};

onMounted(() => {
  if (modeAudio.value) {
    setTimeout(lireQuestionCourante, 400);
  }
});

onUnmounted(() => {
  speechAudioService.arreter();
});

const choisir = (id: string) => {
  if (!corrige.value) choix.value = id;
};

const valider = () => {
  if (!choix.value) return;
  const juste = question.value.options.find((o) => o.id === choix.value)?.est_correcte ?? false;
  if (juste) score.value += 1;
  corrige.value = true;

  if (modeAudio.value) {
    const feedback = juste
      ? `Bonne réponse ! ${question.value.explication}`
      : `Réponse incorrecte. Explication : ${question.value.explication}`;
    speechAudioService.parlerCourt(feedback);
  }
};

const suivante = () => {
  corrige.value = false;
  choix.value = null;
  if (!dernier.value) {
    index.value += 1;
    if (modeAudio.value) {
      setTimeout(lireQuestionCourante, 300);
    }
    return;
  }
  termine.value = true;
  if (modeAudio.value) {
    speechAudioService.parlerCourt(
      `Quiz terminé ! Vous avez obtenu ${score.value} bonnes réponses sur ${total.value}, soit une note de ${note.value} sur ${NOTE_MAX}. ${appreciation.value}`
    );
  }
  emit('finish', note.value);
};

const recommencer = () => {
  questions.value = preparer();
  index.value = 0;
  choix.value = null;
  corrige.value = false;
  score.value = 0;
  termine.value = false;
  if (modeAudio.value) {
    setTimeout(lireQuestionCourante, 300);
  }
};

const styleOption = (opt: QuizQuestion['options'][number]) => {
  if (corrige.value) {
    if (opt.est_correcte) return 'border-brand bg-brand/10 text-brand-dark font-medium';
    if (choix.value === opt.id) return 'border-destructive bg-destructive/10 text-destructive';
    return 'border-border bg-card text-muted-foreground opacity-60';
  }
  return choix.value === opt.id
    ? 'border-brand-dark bg-brand/10 text-foreground ring-1 ring-brand-dark font-medium'
    : 'border-border bg-card text-foreground hover:bg-muted transition';
};
</script>

<template>
  <div
    class="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/50 p-0 sm:p-3 select-none"
    @click.self="fermerModal"
  >
    <div class="flex max-h-[92dvh] w-full sm:max-w-md flex-col overflow-hidden rounded-t-[3px] sm:rounded-[3px] border-t sm:border border-border bg-card pb-safe">
      <div class="sm:hidden mx-auto mt-2 h-1 w-9 rounded-full bg-border shrink-0" />
      <!-- En-tête -->
      <header class="flex items-center justify-between border-b border-border px-4 py-3">
        <div>
          <div class="text-xs font-semibold leading-tight text-foreground">Quiz du code de la route</div>
          <div class="text-[11px] text-muted-foreground">
            <template v-if="!termine">Question {{ index + 1 }} sur {{ total }} · {{ score }} / {{ repondues }}</template>
            <template v-else>Résultat</template>
          </div>
        </div>

        <div class="flex items-center space-x-2">
          <!-- Commutateur Audio -->
          <button
            type="button"
            @click="toggleModeAudio"
            class="flex items-center gap-1.5 px-2 py-1 rounded-[4px] text-[11px] font-semibold border transition cursor-pointer"
            :class="
              modeAudio
                ? 'bg-brand/10 border-brand/30 text-brand'
                : 'bg-muted border-border text-muted-foreground'
            "
            :title="modeAudio ? 'Désactiver la voix audio' : 'Activer la voix audio'"
          >
            <Volume2 v-if="modeAudio" class="w-3.5 h-3.5 text-brand" />
            <VolumeX v-else class="w-3.5 h-3.5 text-muted-foreground" />
            <span class="hidden sm:inline">{{ modeAudio ? 'Voix active' : 'Voix muette' }}</span>
          </button>

          <button
            @click="fermerModal"
            class="rounded p-1 text-muted-foreground hover:bg-muted hover:text-foreground transition cursor-pointer"
            aria-label="Fermer"
          >
            <X class="h-4 w-4" />
          </button>
        </div>
      </header>

      <!-- Progression -->
      <div class="h-1 bg-muted">
        <div class="h-full bg-brand transition-all duration-300" :style="{ width: progression + '%' }"></div>
      </div>

      <!-- Question -->
      <div v-if="!termine" class="flex-1 space-y-3 overflow-y-auto p-4">
        <div class="flex items-start justify-between gap-2">
          <div>
            <div class="text-[11px] font-medium text-brand-dark">{{ question.theme }}</div>
            <p class="mt-1 text-xs font-semibold leading-snug text-foreground">{{ question.enonce }}</p>
          </div>
          <button
            v-if="modeAudio"
            type="button"
            @click="lireQuestionCourante"
            class="shrink-0 flex items-center gap-1 text-[10px] font-bold text-brand bg-brand/10 hover:bg-brand/20 border border-brand/20 px-2 py-1 rounded transition cursor-pointer"
            title="Réécouter la question et les choix"
          >
            <Volume2 class="w-3 h-3" />
            <span>Réécouter</span>
          </button>
        </div>

        <div class="space-y-2">
          <button
            v-for="(opt, i) in question.options"
            :key="opt.id"
            type="button"
            @click="choisir(opt.id)"
            :disabled="corrige"
            :class="['flex w-full min-h-[48px] items-center gap-2.5 rounded-[3px] border px-3 py-2.5 text-left text-xs leading-snug transition', styleOption(opt)]"
          >
            <span class="flex h-6 w-6 shrink-0 items-center justify-center rounded-[3px] bg-muted text-[10px] font-semibold text-muted-foreground">
              {{ LETTRES[i] }}
            </span>
            <span class="flex-1">{{ opt.texte }}</span>
            <Check v-if="corrige && opt.est_correcte" class="h-4 w-4 shrink-0 text-brand-dark" />
            <X v-else-if="corrige && choix === opt.id" class="h-4 w-4 shrink-0 text-destructive" />
          </button>
        </div>

        <div v-if="corrige" class="rounded-[5px] border border-brand-border bg-brand/5 p-3">
          <div class="flex items-center gap-1.5 text-xs font-semibold text-brand-dark">
            <BookOpen class="h-3.5 w-3.5" />
            Explication
          </div>
          <p class="mt-1 text-xs leading-relaxed text-muted-foreground">{{ question.explication }}</p>
        </div>
      </div>

      <!-- Résultat -->
      <div v-else class="flex flex-1 flex-col items-center justify-center gap-3 p-6 text-center">
        <Award class="h-8 w-8 text-brand-dark" />
        <div>
          <div class="font-mono text-3xl font-bold text-foreground">
            {{ note }}<span class="text-base text-muted-foreground"> / {{ NOTE_MAX }}</span>
          </div>
          <div class="mt-1 text-xs text-muted-foreground">
            {{ score }} bonnes réponses sur {{ total }} ({{ pourcentage }} %)
          </div>
        </div>
        <p class="max-w-xs text-xs leading-relaxed text-muted-foreground">{{ appreciation }}</p>
      </div>

      <!-- Actions -->
      <footer class="flex items-center justify-between gap-2 border-t border-border px-4 py-3">
        <button
          type="button"
          @click="fermerModal"
          class="inline-flex h-10 sm:h-8 items-center rounded-[3px] border border-border bg-card px-3 text-xs font-medium text-foreground hover:bg-muted transition"
        >
          Fermer
        </button>

        <template v-if="!termine">
          <button
            v-if="!corrige"
            type="button"
            @click="valider"
            :disabled="!choix"
            class="flex-1 sm:flex-initial inline-flex h-10 sm:h-8 items-center justify-center gap-1.5 rounded-[3px] bg-brand-dark px-3 text-xs font-semibold text-white hover:bg-brand disabled:cursor-not-allowed disabled:opacity-40 transition"
          >
            <Check class="h-3.5 w-3.5" />
            Valider
          </button>
          <button
            v-else
            type="button"
            @click="suivante"
            class="flex-1 sm:flex-initial inline-flex h-10 sm:h-8 items-center justify-center gap-1.5 rounded-[3px] bg-brand-dark px-3 text-xs font-semibold text-white hover:bg-brand transition"
          >
            {{ dernier ? 'Voir le résultat' : 'Suivante' }}
            <ArrowRight class="h-3.5 w-3.5" />
          </button>
        </template>

        <button
          v-else
          type="button"
          @click="recommencer"
          class="flex-1 sm:flex-initial inline-flex h-10 sm:h-8 items-center justify-center gap-1.5 rounded-[3px] bg-brand-dark px-3 text-xs font-semibold text-white hover:bg-brand transition"
        >
          <RotateCcw class="h-3.5 w-3.5" />
          Recommencer
        </button>
      </footer>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import appLogo from "@root/app-icon.png";
import {
  WifiOff,
  Receipt,
  ShieldCheck,
  Sparkles,
  ArrowRight,
  ArrowLeft,
  X,
} from "lucide-vue-next";

const emit = defineEmits<{
  (e: "complete"): void;
  (e: "skip"): void;
}>();

const currentStep = ref(1);
const totalSteps = 4;

const steps = [
  {
    id: 1,
    icon: Sparkles,
    title: "Bienvenue sur ElevenDrive",
    description:
      "Votre logiciel de gestion d'auto-école est prêt. Cette brève introduction vous présente l'essentiel pour bien démarrer.",
    points: [
      "Inscriptions, dossiers et candidats",
      "Séances de code, conduite et examens",
      "Suivi de caisse et comptabilité journalière",
    ],
  },
  {
    id: 2,
    icon: WifiOff,
    title: "Fonctionnement Hors-Ligne",
    description:
      "ElevenDrive fonctionne avec une base locale embarquée. Votre guichet continue d'enregistrer les paiements et inscriptions même sans connexion.",
    points: [
      "Zéro blocage lors des coupures Internet",
      "Données sauvegardées instantanément sur votre PC",
      "Synchronisation automatique avec le Cloud dès le retour du réseau",
    ],
  },
  {
    id: 3,
    icon: Receipt,
    title: "Guichet et Impression Reçus",
    description:
      "Éditez vos reçus immédiatement et sécurisez vos clôtures de caisse grâce aux outils guichet intégrés.",
    points: [
      "Impression thermique ESC/POS directe (58 / 80 mm)",
      "Clôture de caisse quotidienne (Rapport Z)",
      "Raccourcis clavier au guichet : F1 (Candidat), F2 (Versement), F3 (Recherche)",
    ],
  },
  {
    id: 4,
    icon: ShieldCheck,
    title: "Prêt à démarrer",
    description:
      "Toutes les vérifications sont terminées. Vous pouvez maintenant vous connecter avec vos identifiants de travail.",
    points: [
      "Base de données locale initialisée",
      "Mises à jour automatiques configurées",
      "Espace sécurisé et prêt à l'emploi",
    ],
  },
];

const activeStep = computed(() => steps[currentStep.value - 1]);

const next = () => {
  if (currentStep.value < totalSteps) {
    currentStep.value++;
  } else {
    finish();
  }
};

const prev = () => {
  if (currentStep.value > 1) {
    currentStep.value--;
  }
};

const finish = () => {
  try {
    localStorage.setItem("elevendrive_desktop_onboarded_v1", "true");
  } catch {}
  emit("complete");
};

const skip = () => {
  try {
    localStorage.setItem("elevendrive_desktop_onboarded_v1", "true");
  } catch {}
  emit("skip");
};

const handleKeyDown = (e: KeyboardEvent) => {
  if (e.key === "ArrowRight" || e.key === "Enter") next();
  else if (e.key === "ArrowLeft") prev();
  else if (e.key === "Escape") skip();
};

onMounted(() => window.addEventListener("keydown", handleKeyDown));
onUnmounted(() => window.removeEventListener("keydown", handleKeyDown));
</script>

<template>
  <div
    class="fixed inset-0 z-50 flex items-center justify-center bg-background p-4 select-none font-sans"
  >
    <!-- Carte centrale sobre et épurée -->
    <div
      class="w-full max-w-md rounded-[6px] border border-border bg-card p-6 shadow-xs flex flex-col"
    >
      <!-- En-tête : Logo + Étape + Bouton fermer -->
      <div class="flex items-center justify-between pb-4 mb-4 border-b border-border">
        <div class="flex items-center gap-2">
          <img
            :src="appLogo"
            alt="ElevenDrive Logo"
            class="h-5 w-5 object-contain"
          />
          <span class="text-xs font-semibold text-foreground tracking-tight">
            ElevenDrive
          </span>
        </div>

        <div class="flex items-center gap-3">
          <span class="text-[11px] font-mono text-muted-foreground">
            {{ currentStep }} / {{ totalSteps }}
          </span>
          <button
            @click="skip"
            type="button"
            class="text-muted-foreground hover:text-foreground transition p-0.5 rounded-[3px] hover:bg-muted"
            title="Passer l'introduction (Échap)"
          >
            <X class="h-3.5 w-3.5" />
          </button>
        </div>
      </div>

      <!-- Corps : Icône simple + Titre + Descriptif + Liste -->
      <div class="flex flex-col py-2">
        <div class="flex items-center gap-2.5 mb-2.5">
          <div
            class="h-7 w-7 rounded-[4px] border border-border bg-muted/60 flex items-center justify-center text-brand shrink-0"
          >
            <component :is="activeStep.icon" class="h-3.5 w-3.5" />
          </div>
          <h2 class="text-sm font-semibold text-foreground tracking-tight">
            {{ activeStep.title }}
          </h2>
        </div>

        <p class="text-xs text-muted-foreground leading-relaxed mb-4">
          {{ activeStep.description }}
        </p>

        <!-- Points clés sobres -->
        <ul class="space-y-2 mb-4 p-3 rounded-[4px] border border-border bg-muted/20">
          <li
            v-for="(point, idx) in activeStep.points"
            :key="idx"
            class="flex items-start gap-2 text-[11px] text-muted-foreground"
          >
            <span class="h-1 w-1 rounded-full bg-brand mt-1.5 shrink-0"></span>
            <span class="leading-tight">{{ point }}</span>
          </li>
        </ul>
      </div>

      <!-- Pied de carte : Stepper discret + Boutons simples -->
      <div class="flex items-center justify-between pt-4 mt-2 border-t border-border">
        <!-- Points de progression discrets -->
        <div class="flex items-center gap-1.5">
          <span
            v-for="stepNum in totalSteps"
            :key="stepNum"
            class="h-1.5 rounded-full transition-all duration-150"
            :class="[
              stepNum === currentStep
                ? 'w-4 bg-brand'
                : 'w-1.5 bg-muted-foreground/30',
            ]"
          ></span>
        </div>

        <!-- Boutons d'action -->
        <div class="flex items-center gap-2">
          <button
            v-if="currentStep > 1"
            @click="prev"
            type="button"
            class="h-7 px-2.5 rounded-[4px] text-xs font-medium text-muted-foreground hover:text-foreground hover:bg-muted transition flex items-center gap-1"
          >
            <ArrowLeft class="h-3 w-3" />
            <span>Retour</span>
          </button>

          <button
            v-if="currentStep < totalSteps"
            @click="next"
            type="button"
            class="h-7 px-3 rounded-[4px] text-xs font-medium bg-primary text-primary-foreground hover:bg-primary/90 transition flex items-center gap-1 shadow-xs"
          >
            <span>Suivant</span>
            <ArrowRight class="h-3 w-3" />
          </button>

          <button
            v-else
            @click="finish"
            type="button"
            class="h-7 px-3.5 rounded-[4px] text-xs font-medium bg-brand text-white hover:bg-brand/90 transition flex items-center gap-1 shadow-xs"
          >
            <span>Commencer</span>
            <ArrowRight class="h-3 w-3" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
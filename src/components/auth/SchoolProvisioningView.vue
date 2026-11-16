<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import type { AuthSessionData } from "../../types";
import appLogo from "@root/app-icon.png";
import { PROVISION_STEPS } from "./provisioningSteps";
import ProvisioningConsoleGrid from "./ProvisioningConsoleGrid.vue";
import {
  Cpu,
  HardDrive,
  Activity,
  CheckCircle2,
  Sparkles,
  ArrowRight,
  Clock,
  Server,
  Zap,
} from "lucide-vue-next";

defineProps<{
  sessionData: AuthSessionData;
}>();

const emit = defineEmits<{
  (e: "finish"): void;
}>();

const TOTAL_DURATION_SECONDS = 60;
const elapsedSeconds = ref(0);
const progressPercent = ref(0);
const currentSpeed = ref("4.8 Mo/s");
const currentIOPS = ref(480);
const latencyMs = ref(14);
const isCompleted = ref(false);

let timerInterval: ReturnType<typeof setInterval> | null = null;
let statsInterval: ReturnType<typeof setInterval> | null = null;

const consoleLogs = ref<string[]>([
  "[$] Démarrage du moteur de provisioning ElevenDrive v2026...",
  "[$] Handshake SSL sécurisé avec le cluster PostgreSQL Neon Serverless...",
]);

const currentStep = computed(() => {
  for (let i = PROVISION_STEPS.length - 1; i >= 0; i--) {
    if (progressPercent.value >= PROVISION_STEPS[i].targetPercent) {
      return PROVISION_STEPS[i];
    }
  }
  return PROVISION_STEPS[0];
});

const remainingSeconds = computed(() =>
  Math.max(0, TOTAL_DURATION_SECONDS - elapsedSeconds.value),
);

onMounted(() => {
  const tickIntervalMs = 250;
  const incrementPerTick =
    100 / ((TOTAL_DURATION_SECONDS * 1000) / tickIntervalMs);

  timerInterval = setInterval(() => {
    if (progressPercent.value < 100) {
      progressPercent.value = Math.min(
        100,
        +(progressPercent.value + incrementPerTick).toFixed(1),
      );
      elapsedSeconds.value = Math.min(
        TOTAL_DURATION_SECONDS,
        Math.floor((progressPercent.value / 100) * TOTAL_DURATION_SECONDS),
      );

      if (Math.random() > 0.65 && consoleLogs.value.length < 25) {
        const randStep = currentStep.value;
        consoleLogs.value.push(
          `[OK] ${randStep.title} -> ${randStep.detail.substring(0, 45)}... [${currentSpeed.value}]`,
        );
      }
    } else {
      progressPercent.value = 100;
      elapsedSeconds.value = TOTAL_DURATION_SECONDS;
      isCompleted.value = true;
      if (timerInterval) clearInterval(timerInterval);
      if (statsInterval) clearInterval(statsInterval);
      consoleLogs.value.push(
        "[SUCCESS] Auto-École déployée et validée par le Super-Admin avec succès !",
      );
      consoleLogs.value.push(
        "[LICENCE] Période d'essai Super-Admin accordée : 30 jours gratuits activés.",
      );
    }
  }, tickIntervalMs);

  statsInterval = setInterval(() => {
    if (!isCompleted.value) {
      const speed = (4.2 + Math.random() * 4.5).toFixed(1);
      currentSpeed.value = `${speed} Mo/s`;
      currentIOPS.value = Math.floor(420 + Math.random() * 160);
      latencyMs.value = Math.floor(12 + Math.random() * 6);
    }
  }, 1200);
});

onUnmounted(() => {
  if (timerInterval) clearInterval(timerInterval);
  if (statsInterval) clearInterval(statsInterval);
});
</script>

<template>
  <div
    class="fixed inset-0 z-[100] flex flex-col justify-between bg-[#080a0f] text-white p-4 sm:p-8 select-none overflow-y-auto animate-in fade-in duration-300 font-sans"
  >
    <!-- En-tête -->
    <div
      class="w-full flex items-center justify-between border-b border-white/10 pb-4 max-w-5xl mx-auto"
    >
      <div class="flex items-center gap-3">
        <div
          class="h-9 w-9 rounded-lg bg-brand/20 border border-brand/40 flex items-center justify-center p-1.5 shadow-lg shadow-brand/10"
        >
          <img
            :src="appLogo"
            alt="ElevenDrive"
            class="h-full w-full object-contain"
          />
        </div>
        <div>
          <div class="flex items-center gap-2">
            <span class="text-sm font-bold tracking-tight text-white"
              >ElevenDrive Provisioning Engine</span
            >
            <span
              class="px-2 py-0.2 rounded text-[10px] font-mono font-semibold uppercase bg-brand/20 text-brand border border-brand/40"
            >
              Live Setup
            </span>
          </div>
          <div class="text-[11px] text-white/50">
            Initialisation de l'établissement :
            <strong class="text-white/90">{{
              sessionData.auto_ecole.nom_commercial
            }}</strong>
          </div>
        </div>
      </div>

      <div
        class="flex items-center gap-2 px-3 py-1.5 rounded-full bg-white/5 border border-white/10 font-mono text-xs"
      >
        <Clock class="w-3.5 h-3.5 text-brand animate-pulse" />
        <span v-if="!isCompleted" class="text-white/80">
          Temps estimé :
          <strong class="text-brand">{{ remainingSeconds }}s</strong> / 60s
        </span>
        <span
          v-else
          class="text-emerald-400 font-semibold flex items-center gap-1"
        >
          <CheckCircle2 class="w-3.5 h-3.5" />
          <span>Configuration terminée !</span>
        </span>
      </div>
    </div>

    <!-- Corps Central -->
    <div class="max-w-5xl w-full mx-auto my-6 space-y-6">
      <div class="text-center space-y-2 max-w-2xl mx-auto">
        <h2
          class="text-xl sm:text-2xl font-extrabold text-white tracking-tight"
        >
          Traitement des données & création de votre Auto-École
        </h2>
        <p class="text-xs sm:text-sm text-white/60 leading-relaxed">
          Veuillez patienter pendant le déploiement sécurisé de votre
          environnement multi-tenant, l'injection des référentiels ministériels
          et la synchronisation locale. Cette étape automatique ne prend que
          <strong class="text-brand font-semibold">60 secondes</strong>.
        </p>
      </div>

      <!-- Progression & Métriques -->
      <div
        class="p-5 rounded-xl bg-white/[0.03] border border-white/10 backdrop-blur-xl space-y-4 shadow-2xl"
      >
        <div class="flex items-center justify-between text-xs font-mono">
          <span class="text-white/70 flex items-center gap-2">
            <Zap class="w-4 h-4 text-brand animate-bounce" />
            <span>Progression globale :</span>
            <strong class="text-white font-bold"
              >{{ Math.round(progressPercent) }}%</strong
            >
          </span>
          <span class="text-white/50">
            Étape en cours :
            <span class="text-white font-semibold">{{
              currentStep.title
            }}</span>
          </span>
        </div>

        <div
          class="relative w-full h-3 rounded-full bg-black/60 overflow-hidden border border-white/10"
        >
          <div
            class="h-full bg-gradient-to-r from-brand via-emerald-400 to-cyan-400 transition-all duration-300 ease-out shadow-[0_0_20px_rgba(16,185,129,0.7)]"
            :style="{ width: `${progressPercent}%` }"
          ></div>
        </div>

        <div class="grid grid-cols-2 sm:grid-cols-4 gap-3 pt-2">
          <div
            class="p-3 rounded-lg bg-black/40 border border-white/5 space-y-1"
          >
            <div
              class="text-[10px] uppercase font-mono text-white/40 flex items-center gap-1.5"
            >
              <Activity class="w-3.5 h-3.5 text-cyan-400" />
              <span>Vitesse (Data)</span>
            </div>
            <div class="text-base font-bold font-mono text-cyan-300">
              {{ currentSpeed }}
            </div>
          </div>

          <div
            class="p-3 rounded-lg bg-black/40 border border-white/5 space-y-1"
          >
            <div
              class="text-[10px] uppercase font-mono text-white/40 flex items-center gap-1.5"
            >
              <Cpu class="w-3.5 h-3.5 text-amber-400" />
              <span>Débit IOPS</span>
            </div>
            <div class="text-base font-bold font-mono text-amber-300">
              {{ currentIOPS }} IOPS
            </div>
          </div>

          <div
            class="p-3 rounded-lg bg-black/40 border border-white/5 space-y-1"
          >
            <div
              class="text-[10px] uppercase font-mono text-white/40 flex items-center gap-1.5"
            >
              <Server class="w-3.5 h-3.5 text-emerald-400" />
              <span>Latence Base</span>
            </div>
            <div class="text-base font-bold font-mono text-emerald-300">
              {{ latencyMs }} ms
            </div>
          </div>

          <div
            class="p-3 rounded-lg bg-black/40 border border-white/5 space-y-1"
          >
            <div
              class="text-[10px] uppercase font-mono text-white/40 flex items-center gap-1.5"
            >
              <HardDrive class="w-3.5 h-3.5 text-purple-400" />
              <span>Structures</span>
            </div>
            <div class="text-base font-bold font-mono text-purple-300">
              {{ Math.min(24, Math.floor((progressPercent / 100) * 24)) }} / 24
              prêtes
            </div>
          </div>
        </div>
      </div>

      <ProvisioningConsoleGrid
        :steps="PROVISION_STEPS"
        :progress-percent="progressPercent"
        :console-logs="consoleLogs"
      />

      <div
        v-if="isCompleted"
        class="p-5 rounded-xl bg-gradient-to-r from-emerald-950/60 via-brand-dark/40 to-black/80 border border-emerald-500/40 text-center space-y-3 animate-in zoom-in-95 duration-300 shadow-2xl"
      >
        <div
          class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-emerald-500/20 text-emerald-300 border border-emerald-500/40 text-xs font-bold uppercase tracking-wider"
        >
          <Sparkles class="w-4 h-4 text-amber-400" />
          <span>Attribution Super-Admin Officielle</span>
        </div>

        <h3 class="text-lg font-bold text-white">
          🎉 Félicitations ! Votre Auto-École est prête avec 30 jours gratuits
        </h3>

        <p class="text-xs text-white/80 max-w-xl mx-auto leading-relaxed">
          Le système Super-Administrateur ElevenDrive vous a accordé une
          <strong class="text-emerald-300"
            >période d'essai gratuit de 30 jours complets</strong
          >
          avec tous les modules débloqués (Guichet, Caisse, Examens, Flotte et
          Portails Mobiles).
        </p>

        <div class="pt-2">
          <button
            type="button"
            @click="emit('finish')"
            class="inline-flex items-center gap-2 px-6 py-2.5 rounded-[4px] bg-emerald-500 hover:bg-emerald-600 text-black font-extrabold text-xs shadow-xl transition cursor-pointer hover:scale-105"
          >
            <span>Accéder à mon tableau de bord Directeur</span>
            <ArrowRight class="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>

    <div
      class="w-full text-center text-white/30 text-[11px] font-mono border-t border-white/10 pt-3 max-w-5xl mx-auto flex items-center justify-between"
    >
      <span>AutoSchool & ElevenDrive Architecture · Protocole Sécurisé</span>
      <span
        >Session Directeur Active · {{ sessionData.user.prenom }}
        {{ sessionData.user.nom }}</span
      >
    </div>
  </div>
</template>

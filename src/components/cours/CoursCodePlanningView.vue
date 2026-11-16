<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import {
  useCoursCodeStore,
  isSeanceCodeExpired,
} from "../../stores/coursCodeStore";
import { useAuthStore } from "../../stores/authStore";
import { realtimeSocket } from "../../services/realtimeSocketService";
import CoursCodePlanningModal from "./CoursCodePlanningModal.vue";
import ClasseVirtuelleModal from "./ClasseVirtuelleModal.vue";
import type { SeanceCode } from "../../types";
import {
  Video,
  Calendar,
  Clock,
  Plus,
  Users,
  CheckCircle2,
  ExternalLink,
  BookOpen,
  Eye,
  Layers,
  MapPin,
  Car,
  Radio,
} from "lucide-vue-next";

const coursStore = useCoursCodeStore();
const authStore = useAuthStore();
const isWsConnected = realtimeSocket.isConnected;

const isPlanningModalOpen = ref(false);
const seanceToReprogram = ref<SeanceCode | null>(null);
const seanceToJoin = ref<SeanceCode | null>(null);
const canReprogram = computed(
  () =>
    authStore.currentRole === "secretaire" ||
    authStore.currentRole === "manager",
);

const openReprogramModal = (s: SeanceCode) => {
  if (!canReprogram.value) return;
  seanceToReprogram.value = s;
  isPlanningModalOpen.value = true;
};

const closePlanningModal = () => {
  isPlanningModalOpen.value = false;
  seanceToReprogram.value = null;
};
const nowTick = ref(Date.now());
let tickTimer: ReturnType<typeof setInterval> | null = null;
let unsubscribeWs: (() => void) | null = null;

onMounted(() => {
  coursStore.loadSeances();
  tickTimer = setInterval(() => {
    nowTick.value = Date.now();
  }, 30000);
  unsubscribeWs = realtimeSocket.subscribe((msg) => {
    if (msg.type === "COURS_CODE_UPDATED") {
      coursStore.loadSeances();
    }
  });
});

onUnmounted(() => {
  if (tickTimer) {
    clearInterval(tickTimer);
    tickTimer = null;
  }
  if (unsubscribeWs) {
    unsubscribeWs();
    unsubscribeWs = null;
  }
});

const stats = computed(() => {
  const all = coursStore.seances;
  return {
    total: all.length,
    enLigne: all.filter((s) => s.type_seance === "EN_LIGNE").length,
    enSalle: all.filter((s) => s.type_seance === "PRESENTIEL").length,
    aujourdhui: coursStore.seancesAujourdhui.length,
  };
});
</script>

<template>
  <div class="space-y-5 p-4 sm:p-6 max-w-7xl mx-auto">
    <!-- Onglet de bascule Pédagogie (Théorique Code / Pratique Conduite) -->
    <div
      class="flex items-center justify-between flex-wrap gap-3 pb-2 border-b border-border"
    >
      <div
        class="flex items-center gap-1.5 p-1 bg-muted/60 rounded-[6px] border border-border text-xs font-semibold"
      >
        <button
          type="button"
          @click="authStore.setActiveSection('cours-planning')"
          class="px-3 py-1.5 rounded-[4px] flex items-center gap-2 transition cursor-pointer bg-card text-foreground shadow-xs font-bold border border-border/80"
        >
          <Video class="w-3.5 h-3.5 text-brand" />
          <span>📚 Cours Théoriques (Code & Visio)</span>
        </button>

        <button
          type="button"
          @click="authStore.setActiveSection('conduite-planning')"
          class="px-3 py-1.5 rounded-[4px] flex items-center gap-2 transition cursor-pointer text-muted-foreground hover:text-foreground"
        >
          <Car class="w-3.5 h-3.5 text-brand" />
          <span>🚗 Conduite Pratique (Moniteurs & Flotte)</span>
        </button>
      </div>

      <div class="flex items-center space-x-2">
        <button
          type="button"
          @click="isPlanningModalOpen = true"
          class="flex items-center space-x-1.5 px-3.5 py-2 rounded-[4px] bg-brand hover:bg-brand-dark text-primary font-bold text-xs uppercase tracking-wider transition shadow-xs cursor-pointer"
        >
          <Plus class="w-4 h-4" />
          <span>Planifier un cours de code</span>
        </button>
      </div>
    </div>

    <!-- En-tête de section -->
    <header
      class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 pb-4 border-b border-border"
    >
      <div class="space-y-1">
        <div class="flex items-center space-x-2">
          <div
            class="w-8 h-8 rounded-[4px] bg-brand-light border border-brand-border text-brand-dark flex items-center justify-center"
          >
            <Video class="w-4 h-4 text-brand" />
          </div>
          <h1 class="text-base sm:text-lg font-bold text-foreground">
            Planning des Cours Théoriques & Visioconférences
          </h1>
          <span
            class="text-[11px] font-semibold px-2 py-0.5 rounded-[4px] bg-brand-light text-brand-dark border border-brand-border"
          >
            Pédagogie Code
          </span>
          <span
            class="inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-bold"
            :class="
              isWsConnected
                ? 'bg-emerald-500/15 text-emerald-700 dark:text-emerald-300'
                : 'bg-amber-500/15 text-amber-700'
            "
          >
            <Radio class="h-3 w-3 animate-pulse" />
            {{ isWsConnected ? "Temps réel WS" : "Reconnexion WS..." }}
          </span>
        </div>
        <p class="text-xs text-muted-foreground">
          Gérez l'emploi du temps des séances de code en salle et planifiez les
          cours en ligne interactifs pour vos apprenants.
        </p>
      </div>
    </header>

    <!-- Cartes d'indicateurs -->
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
      <div class="p-3.5 rounded-[5px] bg-card border border-border">
        <span class="text-[11px] text-muted-foreground font-medium"
          >Total Séances</span
        >
        <p class="text-xl font-bold tabular-nums text-foreground mt-1">
          {{ stats.total }}
        </p>
      </div>

      <div class="p-3.5 rounded-[5px] bg-card border border-border">
        <span class="text-[11px] text-sky-600 dark:text-sky-400 font-semibold"
          >💻 Cours en Ligne (Visio)</span
        >
        <p class="text-xl font-bold tabular-nums text-foreground mt-1">
          {{ stats.enLigne }}
        </p>
      </div>

      <div class="p-3.5 rounded-[5px] bg-card border border-border">
        <span
          class="text-[11px] text-emerald-600 dark:text-emerald-400 font-semibold"
          >🏫 Cours en Salle</span
        >
        <p class="text-xl font-bold tabular-nums text-foreground mt-1">
          {{ stats.enSalle }}
        </p>
      </div>

      <div class="p-3.5 rounded-[5px] bg-card border border-border">
        <span
          class="text-[11px] text-amber-600 dark:text-amber-400 font-semibold"
          >Prévus Aujourd'hui</span
        >
        <p class="text-xl font-bold tabular-nums text-foreground mt-1">
          {{ stats.aujourdhui }}
        </p>
      </div>
    </div>

    <!-- Barre de filtrage -->
    <div class="flex items-center justify-between border-b border-border pb-3">
      <div class="flex items-center space-x-2">
        <button
          type="button"
          @click="coursStore.filtreType = 'ALL'"
          class="px-3 py-1.5 rounded-[4px] text-xs font-semibold transition cursor-pointer border"
          :class="
            coursStore.filtreType === 'ALL'
              ? 'bg-foreground text-background border-foreground font-bold'
              : 'bg-card text-muted-foreground border-border hover:text-foreground'
          "
        >
          Toutes les sessions ({{ coursStore.seances.length }})
        </button>

        <button
          type="button"
          @click="coursStore.filtreType = 'EN_LIGNE'"
          class="px-3 py-1.5 rounded-[4px] text-xs font-semibold transition cursor-pointer border flex items-center gap-1.5"
          :class="
            coursStore.filtreType === 'EN_LIGNE'
              ? 'bg-sky-500 text-white border-sky-500 font-bold'
              : 'bg-card text-muted-foreground border-border hover:text-foreground'
          "
        >
          <span>💻 En Ligne / Visio</span>
          <span class="text-[10px] px-1 rounded-full bg-card/40">{{
            stats.enLigne
          }}</span>
        </button>

        <button
          type="button"
          @click="coursStore.filtreType = 'PRESENTIEL'"
          class="px-3 py-1.5 rounded-[4px] text-xs font-semibold transition cursor-pointer border flex items-center gap-1.5"
          :class="
            coursStore.filtreType === 'PRESENTIEL'
              ? 'bg-emerald-600 text-white border-emerald-600 font-bold'
              : 'bg-card text-muted-foreground border-border hover:text-foreground'
          "
        >
          <span>🏫 En Salle (Présentiel)</span>
          <span class="text-[10px] px-1 rounded-full bg-card/40">{{
            stats.enSalle
          }}</span>
        </button>
      </div>
    </div>

    <!-- Liste / Grille des séances -->
    <div
      v-if="coursStore.seancesFiltrees.length === 0"
      class="p-12 text-center rounded-[5px] border border-dashed border-border bg-card space-y-3"
    >
      <Video class="w-10 h-10 text-muted-foreground/40 mx-auto" />
      <div>
        <h3 class="text-sm font-semibold text-foreground">
          Aucune séance planifiée
        </h3>
        <p class="text-xs text-muted-foreground mt-0.5">
          Programmez une nouvelle séance de code théorique pour permettre aux
          élèves et moniteurs d'y participer.
        </p>
      </div>
      <button
        type="button"
        @click="isPlanningModalOpen = true"
        class="inline-flex items-center space-x-1.5 px-3 py-1.5 rounded-[4px] bg-brand-light text-primary border border-brand-border text-xs font-semibold hover:bg-brand-light/80 transition cursor-pointer"
      >
        <Plus class="w-3.5 h-3.5" />
        <span>Planifier une séance</span>
      </button>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="s in coursStore.seancesFiltrees"
        :key="s.id"
        class="p-4 rounded-[5px] bg-card border border-border hover:border-brand-border/60 transition shadow-xs flex flex-col justify-between space-y-3"
      >
        <!-- En-tête de la séance -->
        <div class="space-y-2">
          <div class="flex items-center justify-between gap-2">
            <span
              class="inline-flex items-center gap-1 text-[10px] font-bold px-2 py-0.5 rounded-[4px] border uppercase tracking-wider"
              :class="
                s.type_seance === 'EN_LIGNE'
                  ? 'bg-sky-500/10 text-sky-600 border-sky-500/20'
                  : 'bg-emerald-500/10 text-emerald-600 border-emerald-500/20'
              "
            >
              {{
                s.type_seance === "EN_LIGNE"
                  ? "💻 Cours en Ligne"
                  : "🏫 En Salle"
              }}
            </span>

            <span
              class="text-[10px] font-mono text-muted-foreground bg-muted px-2 py-0.5 rounded"
            >
              {{ s.date_seance }}
            </span>
          </div>

          <h3 class="text-xs sm:text-sm font-bold text-foreground leading-snug">
            {{ s.theme_traite }}
          </h3>

          <p class="text-xs text-muted-foreground flex items-center gap-1.5">
            <span class="font-medium text-foreground">Formateur :</span>
            {{ s.formateur_nom }}
          </p>
        </div>

        <!-- Détails horaires & visio -->
        <div
          class="pt-2 border-t border-border space-y-1.5 text-[11px] text-muted-foreground"
        >
          <div class="flex items-center justify-between">
            <span class="flex items-center gap-1 text-foreground font-semibold">
              <Clock class="w-3.5 h-3.5 text-brand" /> {{ s.heure_debut }} -
              {{ s.heure_fin }}
            </span>
            <span class="text-[10px]">Crédit : <strong>+2h code</strong></span>
          </div>

          <div
            v-if="s.type_seance === 'EN_LIGNE' && s.lien_visio"
            class="flex items-center justify-between text-[10px] pt-1"
          >
            <span class="truncate max-w-[200px] text-sky-600 font-mono">
              {{ s.lien_visio.replace("https://", "") }}
            </span>
            <span
              v-if="isSeanceCodeExpired(s, nowTick)"
              class="text-muted-foreground font-bold flex items-center gap-1"
            >
              <span class="w-1.5 h-1.5 rounded-full bg-muted-foreground" />
              Terminée
            </span>
            <span
              v-else
              class="text-emerald-600 font-bold flex items-center gap-1"
            >
              <span
                class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"
              />
              Salon actif
            </span>
          </div>
        </div>

        <!-- Boutons d'action -->
        <div
          class="pt-2 flex items-center justify-between gap-2 border-t border-border"
        >
          <button
            v-if="s.type_seance === 'EN_LIGNE'"
            type="button"
            :disabled="isSeanceCodeExpired(s, nowTick)"
            @click="!isSeanceCodeExpired(s, nowTick) && (seanceToJoin = s)"
            class="flex-1 flex items-center justify-center space-x-1.5 px-3 py-1.5 rounded-[4px] font-bold text-xs transition"
            :class="
              isSeanceCodeExpired(s, nowTick)
                ? 'bg-muted text-muted-foreground border border-border cursor-not-allowed'
                : 'bg-brand hover:bg-brand-dark text-primary cursor-pointer'
            "
          >
            <Video class="w-3.5 h-3.5" />
            <span>{{
              isSeanceCodeExpired(s, nowTick)
                ? "Séance terminée"
                : "Rejoindre la classe"
            }}</span>
          </button>

          <span
            v-else
            class="text-[11px] text-muted-foreground italic flex-1 flex items-center gap-1"
          >
            <MapPin
              class="w-3.5 h-3.5"
              :class="
                isSeanceCodeExpired(s, nowTick)
                  ? 'text-muted-foreground'
                  : 'text-brand'
              "
            />
            {{
              isSeanceCodeExpired(s, nowTick)
                ? "Séance en salle terminée"
                : "Salle de cours agence"
            }}
          </span>

          <button
            v-if="isSeanceCodeExpired(s, nowTick) && canReprogram"
            type="button"
            @click="openReprogramModal(s)"
            class="inline-flex items-center gap-1 px-2.5 py-1.5 rounded-[4px] border border-brand-border bg-brand-light text-brand-dark hover:bg-brand hover:text-white text-xs font-semibold transition cursor-pointer shrink-0"
            title="Reprogrammer ce cours expiré (Directeur / Secrétaire uniquement)"
          >
            <Calendar class="w-3.5 h-3.5" />
            <span>Reprogrammer</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Modale de planification / reprogrammation -->
    <CoursCodePlanningModal
      :isOpen="isPlanningModalOpen"
      :seanceToReprogram="seanceToReprogram"
      @close="closePlanningModal"
    />

    <!-- Modale de classe virtuelle interactive -->
    <ClasseVirtuelleModal :seance="seanceToJoin" @close="seanceToJoin = null" />
  </div>
</template>

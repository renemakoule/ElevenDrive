<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { useCoursCodeStore, isSeanceCodeExpired } from "../../stores/coursCodeStore";
import { useGuichetStore } from "../../stores/guichetStore";
import { realtimeSocket } from "../../services/realtimeSocketService";
import ClasseVirtuelleModal from "./ClasseVirtuelleModal.vue";
import type { SeanceCode } from "../../types";
import {
  Video,
  Clock,
  Calendar,
  CheckCircle2,
  Sparkles,
  ExternalLink,
  BookOpen,
  AlertCircle,
  HelpCircle,
  Laptop,
  Building2,
  Shuffle,
  Radio,
} from "lucide-vue-next";

const props = defineProps<{
  candidat: {
    id: string;
    dossierId: string;
    nom: string;
    heuresCode: number;
    heuresCodeMax: number;
    categoriePermis: string;
    modeTheorique?: "PRESENTIEL" | "EN_LIGNE" | "HYBRIDE";
  };
}>();

const coursStore = useCoursCodeStore();
const guichetStore = useGuichetStore();
const activeSeanceVisio = ref<SeanceCode | null>(null);
const filtreMode = ref<"ALL" | "EN_LIGNE" | "PRESENTIEL">("ALL");
const isWsConnected = realtimeSocket.isConnected;
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

const seancesEnLigne = computed(() => {
  if (filtreMode.value === "ALL") return coursStore.seances;
  return coursStore.seances.filter((s) => s.type_seance === filtreMode.value);
});

const modeConfig = computed(() => {
  const m = props.candidat.modeTheorique || "PRESENTIEL";
  switch (m) {
    case "EN_LIGNE":
      return {
        label: "Formation 100% En Ligne",
        icon: Laptop,
        badgeClass: "bg-sky-500/10 text-sky-600 border-sky-500/20",
        desc: "Votre formule prévoit l'intégralité des cours théoriques à distance via nos classes virtuelles interactives.",
      };
    case "HYBRIDE":
      return {
        label: "Formation Hybride (Salle & Ligne)",
        icon: Shuffle,
        badgeClass: "bg-purple-500/10 text-purple-600 border-purple-500/20",
        desc: "Vous pouvez participer aux séances théoriques en agence comme aux cours interactifs en direct depuis chez vous.",
      };
    default:
      return {
        label: "Formation En Salle (Présentiel)",
        icon: Building2,
        badgeClass: "bg-emerald-500/10 text-emerald-600 border-emerald-500/20",
        desc: "Votre formule standard prévoit les cours en salle, mais vous avez également accès aux visioconférences de révision.",
      };
  }
});

const pctHeures = computed(() => {
  const max = props.candidat.heuresCodeMax || 20;
  return Math.min(100, Math.round((props.candidat.heuresCode / max) * 100));
});

const rejoindreSession = (s: SeanceCode) => {
  if (isSeanceCodeExpired(s, nowTick.value)) {
    guichetStore.showToast("Ce cours est terminé. L'accès au salon est fermé.", "info");
    return;
  }
  activeSeanceVisio.value = s;
};
</script>

<template>
  <div class="space-y-4">
    <!-- ═════════ 1. Carte Statut de Modalité et Progression ═════════ -->
    <section
      class="rounded-[5px] border border-border bg-card p-4 sm:p-5 space-y-4"
    >
      <div
        class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3"
      >
        <div class="flex items-center space-x-3">
          <div
            class="w-10 h-10 rounded-lg bg-sky-500/10 border border-sky-500/20 text-sky-600 flex items-center justify-center shrink-0"
          >
            <component :is="modeConfig.icon" class="w-5 h-5" />
          </div>
          <div>
            <div class="flex items-center gap-2">
              <h2 class="text-sm sm:text-base font-bold text-foreground">
                {{ modeConfig.label }}
              </h2>
              <span
                class="text-[10px] font-bold px-2 py-0.5 rounded-full border uppercase"
                :class="modeConfig.badgeClass"
              >
                Permis {{ candidat.categoriePermis }}
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
            <p class="text-xs text-muted-foreground mt-0.5">
              {{ modeConfig.desc }}
            </p>
          </div>
        </div>

        <div
          class="flex sm:flex-col items-center sm:items-end justify-between border-t sm:border-t-0 pt-2 sm:pt-0 border-border"
        >
          <span class="text-[11px] text-muted-foreground">Volume validé :</span>
          <span
            class="text-sm sm:text-base font-bold text-foreground tabular-nums"
          >
            {{ candidat.heuresCode }} / {{ candidat.heuresCodeMax }} h
          </span>
        </div>
      </div>

      <!-- Barre de progression -->
      <div class="space-y-1.5 pt-1">
        <div class="flex items-center justify-between text-[11px]">
          <span class="text-muted-foreground flex items-center gap-1">
            <Sparkles class="w-3.5 h-3.5 text-brand" /> +2h créditées
            automatiquement par séance suivie
          </span>
          <span class="font-bold text-brand-dark tabular-nums"
            >{{ pctHeures }}% effectué</span
          >
        </div>
        <div class="w-full h-2 rounded-full bg-muted overflow-hidden">
          <div
            class="h-full bg-brand-dark transition-all duration-500 rounded-full"
            :style="{ width: `${pctHeures}%` }"
          />
        </div>
      </div>
    </section>

    <!-- ═════════ 2. Cours Théoriques & Classes Virtuelles Programmées ═════════ -->
    <div class="space-y-3">
      <div class="flex flex-wrap items-center justify-between gap-2">
        <div>
          <h3 class="text-sm font-bold text-foreground flex items-center gap-2">
            <Video class="w-4 h-4 text-sky-600" />
            Planning des Cours Théoriques & Classes Virtuelles (Visio)
          </h3>
          <p class="text-xs text-muted-foreground">
            Retrouvez en temps réel toutes les séances de code en salle et les
            visioconférences programmées par la direction et le secrétariat.
          </p>
        </div>

        <div class="flex items-center gap-1.5">
          <button
            type="button"
            @click="filtreMode = 'ALL'"
            class="px-2.5 py-1 rounded text-[11px] font-semibold border transition cursor-pointer"
            :class="
              filtreMode === 'ALL'
                ? 'bg-foreground text-background border-foreground'
                : 'bg-card text-muted-foreground border-border'
            "
          >
            Tous ({{ coursStore.seances.length }})
          </button>
          <button
            type="button"
            @click="filtreMode = 'EN_LIGNE'"
            class="px-2.5 py-1 rounded text-[11px] font-semibold border transition cursor-pointer"
            :class="
              filtreMode === 'EN_LIGNE'
                ? 'bg-sky-600 text-white border-sky-600'
                : 'bg-card text-muted-foreground border-border'
            "
          >
            💻 En Ligne ({{
              coursStore.seances.filter((s) => s.type_seance === "EN_LIGNE")
                .length
            }})
          </button>
          <button
            type="button"
            @click="filtreMode = 'PRESENTIEL'"
            class="px-2.5 py-1 rounded text-[11px] font-semibold border transition cursor-pointer"
            :class="
              filtreMode === 'PRESENTIEL'
                ? 'bg-emerald-600 text-white border-emerald-600'
                : 'bg-card text-muted-foreground border-border'
            "
          >
            🏫 En Salle ({{
              coursStore.seances.filter((s) => s.type_seance === "PRESENTIEL")
                .length
            }})
          </button>
        </div>
      </div>

      <!-- Aucune session -->
      <div
        v-if="seancesEnLigne.length === 0"
        class="p-8 text-center rounded-[5px] border border-dashed border-border bg-card space-y-2"
      >
        <Video class="w-8 h-8 text-muted-foreground/40 mx-auto" />
        <h4 class="text-xs sm:text-sm font-semibold text-foreground">
          Aucun cours théorique ou classe virtuelle en attente
        </h4>
        <p class="text-xs text-muted-foreground max-w-md mx-auto">
          Les moniteurs et la direction publient régulièrement le planning des
          cours en salle et en direct. Consultez également vos entraînements
          Quiz et Annales PDF.
        </p>
      </div>

      <!-- Grille des cours théoriques & en ligne -->
      <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-3.5">
        <div
          v-for="s in seancesEnLigne"
          :key="s.id"
          class="p-4 rounded-[5px] border border-border bg-card hover:border-sky-500/50 transition flex flex-col justify-between space-y-3"
        >
          <div class="space-y-2">
            <div class="flex items-center justify-between">
              <span
                v-if="s.type_seance === 'EN_LIGNE'"
                class="inline-flex items-center gap-1.5 text-[10px] font-bold px-2 py-0.5 rounded-[4px] bg-sky-500/10 text-sky-600 border border-sky-500/20 uppercase"
              >
                <span
                  class="w-1.5 h-1.5 rounded-full bg-sky-500 animate-pulse"
                />
                💻 Visioconférence Jitsi
              </span>
              <span
                v-else
                class="inline-flex items-center gap-1.5 text-[10px] font-bold px-2 py-0.5 rounded-[4px] bg-emerald-500/10 text-emerald-700 dark:text-emerald-300 border border-emerald-500/20 uppercase"
              >
                <span class="w-1.5 h-1.5 rounded-full bg-emerald-500" />
                🏫 Cours Théorique en Salle
              </span>

              <span
                class="text-[11px] font-mono text-muted-foreground bg-muted px-2 py-0.5 rounded"
              >
                {{ s.date_seance }}
              </span>
            </div>

            <h4 class="text-sm font-bold text-foreground leading-snug">
              {{ s.theme_traite }}
            </h4>

            <p class="text-xs text-muted-foreground">
              Formateur :
              <strong class="text-foreground">{{ s.formateur_nom }}</strong>
            </p>
          </div>

          <div class="pt-2 border-t border-border space-y-2 text-xs">
            <div
              class="flex items-center justify-between text-muted-foreground"
            >
              <span
                class="flex items-center gap-1 text-foreground font-semibold"
              >
                <Clock class="w-3.5 h-3.5 text-brand" /> {{ s.heure_debut }} -
                {{ s.heure_fin }}
              </span>
              <span
                class="text-[11px] font-medium text-emerald-600 bg-emerald-500/10 px-2 py-0.5 rounded"
              >
                Crédit +2 heures
              </span>
            </div>

            <button
              v-if="s.type_seance === 'EN_LIGNE'"
              type="button"
              :disabled="isSeanceCodeExpired(s, nowTick)"
              @click="rejoindreSession(s)"
              class="w-full flex items-center justify-center space-x-2 px-3 py-2 rounded-lg font-bold text-xs transition"
              :class="
                isSeanceCodeExpired(s, nowTick)
                  ? 'bg-muted text-muted-foreground border border-border cursor-not-allowed'
                  : 'bg-sky-600 hover:bg-sky-700 text-white cursor-pointer shadow-xs active:scale-[0.99]'
              "
            >
              <Video class="w-4 h-4" />
              <span>{{
                isSeanceCodeExpired(s, nowTick)
                  ? "Cours terminé — Accès fermé"
                  : "Entrer dans la classe virtuelle"
              }}</span>
            </button>
            <div
              v-else
              class="w-full flex items-center justify-center space-x-2 px-3 py-2 rounded-lg border font-bold text-xs"
              :class="
                isSeanceCodeExpired(s, nowTick)
                  ? 'bg-muted border-border text-muted-foreground'
                  : 'bg-emerald-500/10 border-emerald-500/20 text-emerald-700 dark:text-emerald-300'
              "
            >
              <Building2 class="w-4 h-4" />
              <span>{{
                isSeanceCodeExpired(s, nowTick)
                  ? "Séance en salle terminée"
                  : "Séance programmée en salle de code"
              }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ═════════ 3. Guide de Réussite en Ligne ═════════ -->
    <div
      class="rounded-[5px] border border-border bg-muted/20 p-4 space-y-2 text-xs"
    >
      <div class="font-bold text-foreground flex items-center gap-1.5">
        <HelpCircle class="w-4 h-4 text-brand" />
        Comment se déroule la validation de votre formation théorique ?
      </div>
      <ul class="space-y-1 text-muted-foreground list-disc list-inside">
        <li>
          Chaque participation active à une classe virtuelle crédite
          <strong>+2 heures</strong> sur votre livret d'apprentissage.
        </li>
        <li>
          Les cours sont interactifs : micro et caméra peuvent être activés pour
          poser vos questions au moniteur.
        </li>
        <li>
          Une fois votre quota d'heures théoriques atteint et un score ≥ 35/40
          aux tests blancs, vous êtes présenté à l'examen officiel.
        </li>
      </ul>
    </div>

    <!-- Modale de classe virtuelle interactive pour l'élève -->
    <ClasseVirtuelleModal
      :seance="activeSeanceVisio"
      @close="activeSeanceVisio = null"
    />
  </div>
</template>

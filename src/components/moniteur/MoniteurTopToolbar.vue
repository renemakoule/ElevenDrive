<script setup lang="ts">
import type { SeanceCode } from "../../types";
import {
  Video,
  BookOpen,
  FileText,
  Brain,
  RefreshCw,
  CloudOff,
} from "lucide-vue-next";

defineProps<{
  activeSection: string;
  coursCount: number;
  supportsCount: number;
  epreuvesCount: number;
  quizSoumissionsCount: number;
  isLoading: boolean;
  enAttenteSync: number;
  seanceCodeEnLigneAujourdhui: SeanceCode | null;
}>();

const emit = defineEmits<{
  navigate: [section: string];
  refresh: [];
  openClasseVirtuelle: [seance: SeanceCode];
}>();
</script>

<template>
  <div class="space-y-4">
    <!-- ═════════ En-tête ═════════ -->
    <header
      class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
    >
      <div class="min-w-0">
        <p class="text-xs text-muted-foreground">
          <template v-if="activeSection === 'livret'">
            Évaluation REMC instantanée, suivi kilométrique et signature
            numérique synchronisés en direct par WebSocket.
          </template>
          <template v-else-if="activeSection === 'vehicule'">
            État technique, compteur kilométrique en temps réel, assurances et
            disponibilité de votre véhicule d'apprentissage.
          </template>
          <template v-else>
            Planning des leçons, feuille de route et émargement tactile
            synchronisés en temps réel.
          </template>
        </p>
      </div>

      <div class="flex flex-wrap items-center gap-2 text-xs">
        <!-- Raccourcis pédagogiques (masqués sur mobile : gérés par MobileBottomNav) -->
        <button
          type="button"
          @click="emit('navigate', 'cours-planning')"
          class="hidden md:inline-flex h-8 items-center gap-1.5 px-2.5 rounded-[3px] border border-border bg-card text-xs font-semibold text-foreground hover:bg-muted transition cursor-pointer"
          title="Consulter les cours théoriques et visioconférences"
        >
          <Video class="w-3.5 h-3.5 text-sky-600" />
          <span>Cours & Visio</span>
          <span
            class="rounded-full bg-sky-500/15 px-1.5 py-0.2 text-[10px] font-bold text-sky-700 dark:text-sky-300"
          >
            {{ coursCount }}
          </span>
        </button>

        <button
          type="button"
          @click="emit('navigate', 'supports-cours')"
          class="hidden md:inline-flex h-8 items-center gap-1.5 px-2.5 rounded-[3px] border border-border bg-card text-xs font-semibold text-foreground hover:bg-muted transition cursor-pointer"
          title="Consulter et écouter les supports de cours théoriques"
        >
          <BookOpen class="w-3.5 h-3.5 text-brand" />
          <span>Supports & Audio</span>
          <span
            class="rounded-full bg-brand/15 px-1.5 py-0.2 text-[10px] font-bold text-brand-dark"
          >
            {{ supportsCount }}
          </span>
        </button>

        <button
          type="button"
          @click="emit('navigate', 'epreuves')"
          class="hidden md:inline-flex h-8 items-center gap-1.5 px-2.5 rounded-[3px] border border-border bg-card text-xs font-semibold text-foreground hover:bg-muted transition cursor-pointer"
          title="Consulter et écouter les épreuves officielles"
        >
          <FileText class="w-3.5 h-3.5 text-sky-600" />
          <span>Annales & Sujets</span>
          <span
            class="rounded-full bg-sky-500/15 px-1.5 py-0.2 text-[10px] font-bold text-sky-700 dark:text-sky-300"
          >
            {{ epreuvesCount }}
          </span>
        </button>

        <button
          type="button"
          @click="emit('navigate', 'quiz-forms')"
          class="hidden md:inline-flex h-8 items-center gap-1.5 px-2.5 rounded-[3px] border border-brand-border bg-brand-light/50 text-xs font-bold text-brand-dark hover:bg-brand-light transition cursor-pointer"
          title="Créer un formulaire de Quiz (style Google Forms) et consulter les réponses des élèves"
        >
          <Brain class="w-3.5 h-3.5 text-brand" />
          <span>Quiz & Réponses</span>
          <span
            class="rounded-full bg-brand px-1.5 py-0.2 text-[10px] font-bold text-brand-dark"
          >
            {{ quizSoumissionsCount }}
          </span>
        </button>

        <button
          type="button"
          @click="emit('refresh')"
          :disabled="isLoading"
          class="inline-flex h-9 w-9 md:h-8 md:w-8 items-center justify-center rounded-[3px] border border-border bg-card text-muted-foreground hover:bg-muted hover:text-foreground transition disabled:opacity-50"
          title="Rafraîchir le planning"
        >
          <RefreshCw
            class="h-3.5 w-3.5"
            :class="{ 'animate-spin': isLoading }"
          />
        </button>

        <span class="flex items-center gap-1.5 text-muted-foreground">
          <CloudOff class="h-3.5 w-3.5 text-muted-foreground" />
          <span class="hidden sm:inline">Mode hors-ligne</span>
          <span
            v-if="enAttenteSync > 0"
            class="rounded-full bg-amber-500/10 px-1.5 py-0.5 text-[10px] font-medium tabular-nums text-amber-600 dark:text-amber-400"
            :title="`${enAttenteSync} événement(s) à synchroniser`"
          >
            {{ enAttenteSync }} à sync.
          </span>
        </span>
      </div>
    </header>

    <!-- Alerte Séance de Code Théorique / Visioconférence Planifiée -->
    <div
      v-if="seanceCodeEnLigneAujourdhui"
      class="rounded-[4px] border border-sky-500/30 bg-sky-500/5 p-3.5 flex flex-col sm:flex-row sm:items-center justify-between gap-3"
    >
      <div class="flex items-center space-x-3">
        <div
          class="w-9 h-9 rounded-[4px] bg-sky-500/10 text-sky-600 flex items-center justify-center shrink-0"
        >
          <Video class="w-4 h-4" />
        </div>
        <div>
          <div class="flex items-center gap-2">
            <span
              class="text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-full bg-sky-500/20 text-sky-700 dark:text-sky-300"
            >
              {{
                seanceCodeEnLigneAujourdhui.type_seance === "EN_LIGNE"
                  ? "Visioconférence Code"
                  : "Cours de Code en Salle"
              }}
            </span>
            <span class="text-xs text-muted-foreground font-mono">
              {{ seanceCodeEnLigneAujourdhui.date_seance }} ·
              {{ seanceCodeEnLigneAujourdhui.heure_debut }} -
              {{ seanceCodeEnLigneAujourdhui.heure_fin }}
            </span>
          </div>
          <p class="text-xs font-bold text-foreground mt-0.5">
            {{ seanceCodeEnLigneAujourdhui.theme_traite }} (Formateur :
            {{ seanceCodeEnLigneAujourdhui.formateur_nom }})
          </p>
        </div>
      </div>

      <div class="flex flex-wrap items-center gap-2">
        <button
          v-if="seanceCodeEnLigneAujourdhui.type_seance === 'EN_LIGNE'"
          type="button"
          @click="emit('openClasseVirtuelle', seanceCodeEnLigneAujourdhui)"
          class="min-h-[40px] sm:min-h-0 flex-1 sm:flex-initial inline-flex items-center justify-center gap-1.5 px-3 py-1.5 rounded-[3px] bg-sky-600 hover:bg-sky-700 text-white font-bold text-xs transition cursor-pointer"
        >
          <Video class="w-3.5 h-3.5" />
          <span>Démarrer le direct</span>
        </button>
        <button
          type="button"
          @click="emit('navigate', 'cours-planning')"
          class="min-h-[40px] sm:min-h-0 flex-1 sm:flex-initial inline-flex items-center justify-center gap-1.5 px-3 py-1.5 rounded-[3px] border border-sky-500/30 bg-card hover:bg-muted text-foreground font-bold text-xs transition cursor-pointer"
        >
          <span>Voir les cours ({{ coursCount }})</span>
        </button>
      </div>
    </div>
  </div>
</template>

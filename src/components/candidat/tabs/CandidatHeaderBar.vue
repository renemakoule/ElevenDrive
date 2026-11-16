<script setup lang="ts">
import type { Candidat } from "../../../types";
import { Radio, BookOpen, Brain, CalendarPlus } from "lucide-vue-next";

defineProps<{
  isRecusSection: boolean;
  isWsConnected: boolean;
  lastRealtimeSyncAt: string;
  searchResults: Candidat[];
  selectedCandId: string;
  hasActiveCand: boolean;
  epreuvesCount: number;
  quizzesCount: number;
}>();

const emit = defineEmits<{
  (e: "update:selectedCandId", id: string): void;
  (e: "open-epreuves"): void;
  (e: "open-quiz"): void;
  (e: "open-reservation"): void;
}>();
</script>

<template>
  <header
    class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
  >
    <div class="min-w-0">
      <div class="flex items-center gap-2">
        <h1 class="text-base font-semibold text-foreground">
          <template v-if="isRecusSection">
            Mes Reçus Officiels & Solde en Temps Réel
          </template>
          <template v-else>Mon espace apprenant</template>
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
          {{ isWsConnected ? "Temps réel WS" : "Reconnexion WS..." }}
        </span>
        <span
          v-if="lastRealtimeSyncAt"
          class="hidden md:inline text-[10px] font-mono text-muted-foreground"
        >
          Sync: {{ lastRealtimeSyncAt }}
        </span>
      </div>
      <p class="text-xs text-muted-foreground">
        Suivi instantané de votre formation, progression pédagogique, livret
        REMC, reçus et leçons en temps réel.
      </p>
    </div>

    <div class="flex flex-wrap items-center gap-2">
      <select
        v-if="searchResults && searchResults.length > 0"
        :value="selectedCandId"
        @change="
          emit(
            'update:selectedCandId',
            ($event.target as HTMLSelectElement).value,
          )
        "
        aria-label="Sélectionner le candidat"
        class="h-10 sm:h-8 w-full sm:w-auto rounded-[3px] border border-border bg-card px-2.5 text-xs text-foreground focus:border-brand-dark focus:outline-none focus:ring-2 focus:ring-brand-dark/20 transition"
      >
        <option v-for="c in searchResults" :key="c.id" :value="c.id">
          {{ c.nom }} {{ c.prenoms }} · {{ c.numero_dossier || "DOS-ATTENTE" }}
        </option>
      </select>

      <button
        type="button"
        @click="emit('open-epreuves')"
        :disabled="!hasActiveCand"
        class="hidden sm:inline-flex h-8 items-center gap-1.5 rounded-[3px] border border-brand-border bg-brand-light px-3 text-xs font-semibold text-brand-dark transition hover:bg-brand-light/80 disabled:opacity-50"
      >
        <BookOpen class="h-3.5 w-3.5 text-brand" />
        Épreuves PDF ({{ epreuvesCount }})
      </button>

      <button
        type="button"
        @click="emit('open-quiz')"
        :disabled="!hasActiveCand"
        class="hidden sm:inline-flex h-8 items-center gap-1.5 rounded-[3px] border border-border bg-card px-3 text-xs font-medium text-foreground transition hover:bg-muted disabled:opacity-50"
      >
        <Brain class="h-3.5 w-3.5 text-brand-dark" />
        Quiz Moniteur ({{ quizzesCount }})
      </button>

      <button
        type="button"
        @click="emit('open-reservation')"
        :disabled="!hasActiveCand"
        class="w-full sm:w-auto inline-flex h-10 sm:h-8 items-center justify-center gap-1.5 rounded-[3px] bg-brand-dark px-3 text-xs font-medium text-white transition hover:bg-brand disabled:opacity-50"
      >
        <CalendarPlus class="h-3.5 w-3.5" />
        Réserver une leçon
      </button>
    </div>
  </header>
</template>

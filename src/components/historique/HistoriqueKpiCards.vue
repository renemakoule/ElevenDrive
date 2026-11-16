<script setup lang="ts">
import { History, Receipt, Brain, Car, Video } from "lucide-vue-next";
import { fcfa } from "./types";

defineProps<{
  currentRole: string;
  kpiTotalActions: number;
  kpiTotalMontant: number;
  quizCount: number;
  kpiTotalConduite: number;
  kpiTotalCode: number;
}>();
</script>

<template>
  <div class="grid grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4">
    <div
      class="rounded-[4px] border border-border bg-card p-4 flex flex-col justify-between"
    >
      <div class="flex items-center justify-between text-muted-foreground">
        <span class="text-xs font-medium">Total Événements</span>
        <History class="h-4 w-4 text-brand" />
      </div>
      <div class="mt-2">
        <span class="text-xl font-bold text-foreground tabular-nums">
          {{ kpiTotalActions }}
        </span>
        <p class="text-[11px] text-muted-foreground mt-0.5">
          Historique complet synchronisé
        </p>
      </div>
    </div>

    <div
      v-if="
        currentRole === 'secretaire' ||
        currentRole === 'manager' ||
        currentRole === 'caissier'
      "
      class="rounded-[4px] border border-border bg-card p-4 flex flex-col justify-between"
    >
      <div class="flex items-center justify-between text-muted-foreground">
        <span class="text-xs font-medium">Paiements Validés</span>
        <Receipt class="h-4 w-4 text-emerald-500" />
      </div>
      <div class="mt-2">
        <span class="text-xl font-bold text-foreground tabular-nums">
          {{ fcfa(kpiTotalMontant) }} FCFA
        </span>
        <p class="text-[11px] text-emerald-600 dark:text-emerald-400 mt-0.5">
          Total transactions enregistrées
        </p>
      </div>
    </div>
    <div
      v-else
      class="rounded-[4px] border border-border bg-card p-4 flex flex-col justify-between"
    >
      <div class="flex items-center justify-between text-muted-foreground">
        <span class="text-xs font-medium">Séries Quiz Effectuées</span>
        <Brain class="h-4 w-4 text-purple-500" />
      </div>
      <div class="mt-2">
        <span class="text-xl font-bold text-foreground tabular-nums">
          {{ quizCount }}
        </span>
        <p class="text-[11px] text-purple-600 dark:text-purple-400 mt-0.5">
          Évaluations et entraînements
        </p>
      </div>
    </div>

    <div
      class="rounded-[4px] border border-border bg-card p-4 flex flex-col justify-between"
    >
      <div class="flex items-center justify-between text-muted-foreground">
        <span class="text-xs font-medium">Leçons Pratiques</span>
        <Car class="h-4 w-4 text-blue-500" />
      </div>
      <div class="mt-2">
        <span class="text-xl font-bold text-foreground tabular-nums">
          {{ kpiTotalConduite }}
        </span>
        <p class="text-[11px] text-blue-600 dark:text-blue-400 mt-0.5">
          Séances sur piste & circulation
        </p>
      </div>
    </div>

    <div
      class="rounded-[4px] border border-border bg-card p-4 flex flex-col justify-between"
    >
      <div class="flex items-center justify-between text-muted-foreground">
        <span class="text-xs font-medium">Séances Code & Visio</span>
        <Video class="h-4 w-4 text-indigo-500" />
      </div>
      <div class="mt-2">
        <span class="text-xl font-bold text-foreground tabular-nums">
          {{ kpiTotalCode }}
        </span>
        <p class="text-[11px] text-indigo-600 dark:text-indigo-400 mt-0.5">
          Présentiel & classes virtuelles
        </p>
      </div>
    </div>
  </div>
</template>


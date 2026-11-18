<script setup lang="ts">
import Skeleton from "./Skeleton.vue";
import KpiGridSkeleton from "./KpiGridSkeleton.vue";

withDefaults(
  defineProps<{
    showLicence?: boolean;
  }>(),
  {
    showLicence: true,
  },
);
</script>

<template>
  <div class="space-y-4">
    <!-- En-tête de la vue -->
    <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
      <div class="space-y-1.5">
        <Skeleton class="h-5 w-48" />
        <Skeleton class="h-3.5 w-72 sm:w-96" />
      </div>
      <div class="flex items-center gap-2">
        <Skeleton class="h-8 w-20 rounded-[3px]" />
        <Skeleton class="h-8 w-32 rounded-[3px]" />
        <Skeleton class="h-8 w-28 rounded-[3px]" />
        <Skeleton class="h-8 w-20 rounded-[3px]" />
      </div>
    </div>

    <!-- Bandeau Licence Skeleton -->
    <div
      v-if="showLicence"
      class="rounded-[3px] border border-border bg-card p-3.5 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3"
    >
      <div class="flex items-center gap-3">
        <Skeleton class="h-8 w-8 rounded-full shrink-0" />
        <div class="space-y-1">
          <Skeleton class="h-3.5 w-44" />
          <Skeleton class="h-2.5 w-60" />
        </div>
      </div>
      <Skeleton class="h-7 w-28 rounded-[3px]" />
    </div>

    <!-- 4 Cartes KPI Skeleton -->
    <KpiGridSkeleton :count="4" />

    <!-- 2 Graphiques Chart.js Skeleton -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-3">
      <!-- Graphique Donut (Répartition Formules) -->
      <div class="lg:col-span-5 rounded-[3px] border border-border bg-card p-4 space-y-4">
        <div class="flex items-center justify-between">
          <Skeleton class="h-4 w-40" />
          <Skeleton class="h-3 w-16" />
        </div>
        <div class="flex justify-center items-center py-6">
          <Skeleton class="h-44 w-44 rounded-full" />
        </div>
        <div class="grid grid-cols-2 gap-2 pt-2 border-t border-border">
          <div v-for="c in 4" :key="c" class="flex items-center gap-2">
            <Skeleton class="h-2.5 w-2.5 rounded-full" />
            <Skeleton class="h-3 w-20" />
          </div>
        </div>
      </div>

      <!-- Graphique Bar (Recettes mensuelles) -->
      <div class="lg:col-span-7 rounded-[3px] border border-border bg-card p-4 space-y-4">
        <div class="flex items-center justify-between">
          <Skeleton class="h-4 w-48" />
          <Skeleton class="h-6 w-24 rounded" />
        </div>
        <div class="h-48 flex items-end justify-between gap-2 pt-4 px-2">
          <div
            v-for="h in [40, 65, 30, 85, 55, 90, 70]"
            :key="h"
            class="flex-1 flex flex-col items-center gap-2 h-full justify-end"
          >
            <Skeleton
              class="w-full rounded-t-[3px]"
              :style="{ height: `${h}%` }"
            />
            <Skeleton class="h-2.5 w-6" />
          </div>
        </div>
      </div>
    </div>

    <!-- Encaissements récents Skeleton -->
    <div class="rounded-[3px] border border-border bg-card p-4 space-y-3">
      <div class="flex items-center justify-between border-b border-border pb-3">
        <Skeleton class="h-4 w-44" />
        <Skeleton class="h-3 w-24" />
      </div>
      <div class="divide-y divide-border">
        <div
          v-for="i in 4"
          :key="i"
          class="py-2.5 flex items-center justify-between gap-3"
        >
          <div class="flex items-center gap-3">
            <Skeleton class="h-8 w-8 rounded-full" />
            <div class="space-y-1">
              <Skeleton class="h-3.5 w-32" />
              <Skeleton class="h-2.5 w-20 font-mono" />
            </div>
          </div>
          <div class="text-right space-y-1">
            <Skeleton class="h-3.5 w-24 ml-auto" />
            <Skeleton class="h-2.5 w-16 ml-auto" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

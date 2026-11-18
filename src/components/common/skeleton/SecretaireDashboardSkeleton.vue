<script setup lang="ts">
import Skeleton from "./Skeleton.vue";
import KpiGridSkeleton from "./KpiGridSkeleton.vue";
import TableSkeleton from "./TableSkeleton.vue";
</script>

<template>
  <div class="space-y-4 bg-background">
    <!-- En-tête -->
    <header class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div class="min-w-0 space-y-1">
        <Skeleton class="h-5 w-44" />
        <Skeleton class="h-3 w-72 sm:w-80" />
      </div>
      <div class="flex flex-wrap items-center gap-1.5 sm:gap-2">
        <Skeleton class="h-9 w-9 sm:h-8 sm:w-8 rounded-[3px]" />
        <Skeleton class="h-9 sm:h-8 w-28 rounded-[3px]" />
        <Skeleton class="h-9 sm:h-8 w-32 rounded-[3px]" />
        <Skeleton class="h-9 sm:h-8 w-36 rounded-[3px]" />
      </div>
    </header>

    <!-- 4 KPI cards divisées -->
    <KpiGridSkeleton :count="4" :divided="true" />

    <!-- Tableau dossiers récents -->
    <TableSkeleton :rows="6" :header-cols="5" :show-search="true" />
  </div>
</template>

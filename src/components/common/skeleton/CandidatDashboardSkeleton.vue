<script setup lang="ts">
import Skeleton from "./Skeleton.vue";
</script>

<template>
  <div class="space-y-4">
    <!-- Profil & En-tête Candidat -->
    <div class="rounded-[5px] border border-border bg-card p-4 sm:p-5">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div class="flex items-center gap-3.5">
          <Skeleton class="h-14 w-14 rounded-full shrink-0" />
          <div class="space-y-1.5">
            <div class="flex items-center gap-2">
              <Skeleton class="h-5 w-40" />
              <Skeleton class="h-4 w-20 rounded-full" />
            </div>
            <div class="flex items-center gap-2">
              <Skeleton class="h-3.5 w-24 font-mono" />
              <Skeleton class="h-3.5 w-28" />
            </div>
          </div>
        </div>
        <div class="flex items-center gap-2">
          <Skeleton class="h-8 w-28 rounded-[4px]" />
          <Skeleton class="h-8 w-24 rounded-[4px]" />
        </div>
      </div>
    </div>

    <!-- Double Jauges de Progression (Code Théorique & Conduite Pratique) -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
      <div
        v-for="j in 2"
        :key="j"
        class="rounded-[5px] border border-border bg-card p-4 space-y-3"
      >
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <Skeleton class="h-7 w-7 rounded-lg" />
            <Skeleton class="h-4 w-32" />
          </div>
          <Skeleton class="h-5 w-14 rounded-full" />
        </div>
        <Skeleton class="h-3 w-full rounded-full" />
        <div class="flex items-center justify-between pt-1">
          <Skeleton class="h-3 w-28" />
          <Skeleton class="h-3 w-16" />
        </div>
      </div>
    </div>

    <!-- État Financier & Reste à Payer -->
    <div class="rounded-[5px] border border-border bg-card p-4 space-y-3">
      <div class="flex items-center justify-between border-b border-border pb-3">
        <Skeleton class="h-4 w-40" />
        <Skeleton class="h-5 w-24 rounded-full" />
      </div>
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
        <div v-for="c in 3" :key="c" class="p-3 rounded-[4px] bg-muted/30 space-y-1.5">
          <Skeleton class="h-3 w-20" />
          <Skeleton class="h-5 w-28 font-mono" />
        </div>
      </div>
    </div>
  </div>
</template>

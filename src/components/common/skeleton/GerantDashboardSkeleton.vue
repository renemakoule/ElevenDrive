<script setup lang="ts">
import Skeleton from "./Skeleton.vue";
</script>

<template>
  <div class="flex flex-1 flex-col gap-4 overflow-y-auto">
    <!-- En-tête -->
    <header class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div class="flex items-center gap-3">
        <Skeleton class="h-8 w-20 rounded-[4px]" />
        <div class="space-y-1">
          <Skeleton class="h-5 w-48" />
          <Skeleton class="h-3 w-64" />
        </div>
      </div>
    </header>

    <!-- 2 KPIs créances et assidus -->
    <section class="grid grid-cols-1 divide-y divide-border rounded-[5px] border border-border bg-card sm:grid-cols-2 sm:divide-x sm:divide-y-0">
      <div v-for="i in 2" :key="i" class="p-4 space-y-1.5">
        <Skeleton class="h-3 w-28" />
        <Skeleton class="h-7 w-36" />
        <Skeleton class="h-2.5 w-44" />
      </div>
    </section>

    <!-- Onglets navigation -->
    <div class="flex gap-2 border-b border-border pb-1">
      <Skeleton class="h-8 w-32 rounded-t" />
      <Skeleton class="h-8 w-28 rounded-t" />
      <Skeleton class="h-8 w-24 rounded-t" />
    </div>

    <!-- Contenu liste -->
    <div class="rounded-[5px] border border-border bg-card p-4 space-y-3">
      <div
        v-for="i in 5"
        :key="i"
        class="py-3 flex items-center justify-between border-b border-border last:border-b-0"
      >
        <div class="flex items-center gap-3">
          <Skeleton class="h-9 w-9 rounded-full" />
          <div class="space-y-1.5">
            <Skeleton class="h-4 w-36" />
            <Skeleton class="h-2.5 w-24 font-mono" />
          </div>
        </div>
        <div class="flex items-center gap-3">
          <div class="text-right space-y-1">
            <Skeleton class="h-4 w-24 ml-auto" />
            <Skeleton class="h-2.5 w-16 ml-auto" />
          </div>
          <Skeleton class="h-8 w-20 rounded-[4px]" />
        </div>
      </div>
    </div>
  </div>
</template>

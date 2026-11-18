<script setup lang="ts">
import Skeleton from "./Skeleton.vue";

withDefaults(
  defineProps<{
    rows?: number;
    showSearch?: boolean;
    headerCols?: number;
  }>(),
  {
    rows: 5,
    showSearch: true,
    headerCols: 5,
  },
);
</script>

<template>
  <div class="rounded-[3px] border border-border bg-card overflow-hidden">
    <!-- Barre de recherche / filtre -->
    <div
      v-if="showSearch"
      class="flex flex-col gap-2.5 px-4 py-3 sm:flex-row sm:items-center sm:justify-between border-b border-border"
    >
      <div class="flex items-center gap-2">
        <Skeleton class="h-4 w-32" />
        <Skeleton class="h-3 w-20" />
      </div>
      <div class="flex items-center gap-2">
        <Skeleton class="h-7 w-48 rounded-[4px]" />
        <Skeleton class="h-7 w-24 rounded-[4px]" />
      </div>
    </div>

    <!-- Vue Mobile (< md) -->
    <div class="md:hidden divide-y divide-border">
      <div
        v-for="r in rows"
        :key="r"
        class="w-full px-3.5 py-3 flex items-center justify-between gap-3"
      >
        <div class="flex items-center gap-2.5 min-w-0 flex-1">
          <Skeleton class="h-9 w-9 rounded-full shrink-0" />
          <div class="space-y-1.5 flex-1">
            <Skeleton class="h-3.5 w-3/4" />
            <Skeleton class="h-2.5 w-1/2" />
          </div>
        </div>
        <Skeleton class="h-5 w-16 rounded-full shrink-0" />
      </div>
    </div>

    <!-- Vue Desktop (>= md) -->
    <div class="hidden md:block overflow-x-auto">
      <table class="w-full text-left text-xs border-collapse">
        <thead>
          <tr class="border-b border-border bg-muted/40">
            <th v-for="c in headerCols" :key="c" class="px-4 py-2.5">
              <Skeleton class="h-3 w-20" />
            </th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border">
          <tr v-for="r in rows" :key="r" class="px-4 py-3">
            <td class="px-4 py-3">
              <div class="flex items-center gap-2.5">
                <Skeleton class="h-8 w-8 rounded-full shrink-0" />
                <div class="space-y-1">
                  <Skeleton class="h-3.5 w-28" />
                  <Skeleton class="h-2.5 w-16" />
                </div>
              </div>
            </td>
            <td class="px-3 py-3">
              <div class="space-y-1">
                <Skeleton class="h-3.5 w-20 font-mono" />
                <Skeleton class="h-2.5 w-24" />
              </div>
            </td>
            <td v-if="headerCols > 2" class="px-3 py-3">
              <Skeleton class="h-3.5 w-28" />
            </td>
            <td v-if="headerCols > 3" class="px-3 py-3">
              <Skeleton class="h-5 w-20 rounded-full" />
            </td>
            <td v-if="headerCols > 4" class="px-4 py-3 text-right">
              <Skeleton class="h-4 w-4 ml-auto rounded" />
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

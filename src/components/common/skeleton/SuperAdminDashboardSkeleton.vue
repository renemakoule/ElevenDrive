<script setup lang="ts">
import Skeleton from "./Skeleton.vue";

withDefaults(
  defineProps<{
    rows?: number;
  }>(),
  {
    rows: 7,
  },
);
</script>

<template>
  <div class="space-y-3">
    <!-- Barre de sous-filtres -->
    <div class="h-10 border-b border-border bg-card px-4 flex items-center justify-between">
      <div class="flex items-center gap-3">
        <Skeleton class="h-4 w-36" />
        <Skeleton class="h-5 w-40 rounded-[4px]" />
      </div>
      <Skeleton class="h-7 w-28 rounded-[4px]" />
    </div>

    <!-- Table Console Console Neon SuperAdmin -->
    <div class="border border-border rounded-[4px] bg-card overflow-hidden">
      <table class="w-full text-left text-xs border-collapse">
        <thead class="bg-muted/80 border-b border-border">
          <tr>
            <th class="py-2.5 px-3"><Skeleton class="h-3 w-14" /></th>
            <th class="py-2.5 px-3"><Skeleton class="h-3 w-32" /></th>
            <th class="py-2.5 px-3"><Skeleton class="h-3 w-20" /></th>
            <th class="py-2.5 px-3 text-center"><Skeleton class="h-3 w-24 mx-auto" /></th>
            <th class="py-2.5 px-3"><Skeleton class="h-3 w-28" /></th>
            <th class="py-2.5 px-3"><Skeleton class="h-3 w-20" /></th>
            <th class="py-2.5 px-3 text-right"><Skeleton class="h-3 w-16 ml-auto" /></th>
            <th class="py-2.5 px-3 text-right"><Skeleton class="h-3 w-16 ml-auto" /></th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border">
          <tr v-for="r in rows" :key="r" class="px-3 py-2.5">
            <td class="py-2 px-3">
              <Skeleton class="h-5 w-16 rounded-[3px]" />
            </td>
            <td class="py-2 px-3">
              <div class="space-y-1">
                <Skeleton class="h-4 w-36 font-semibold" />
                <Skeleton class="h-2.5 w-24" />
              </div>
            </td>
            <td class="py-2 px-3">
              <Skeleton class="h-3.5 w-20 font-mono" />
            </td>
            <td class="py-2 px-3 text-center">
              <Skeleton class="h-5 w-14 mx-auto rounded-full" />
            </td>
            <td class="py-2 px-3">
              <div class="space-y-1">
                <Skeleton class="h-3.5 w-24" />
                <Skeleton class="h-2.5 w-16" />
              </div>
            </td>
            <td class="py-2 px-3">
              <Skeleton class="h-3.5 w-24 font-mono" />
            </td>
            <td class="py-2 px-3 text-right">
              <Skeleton class="h-3.5 w-10 ml-auto" />
            </td>
            <td class="py-2 px-3 text-right">
              <Skeleton class="h-3.5 w-12 ml-auto" />
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

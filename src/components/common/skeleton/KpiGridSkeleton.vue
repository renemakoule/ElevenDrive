<script setup lang="ts">
import Skeleton from "./Skeleton.vue";

withDefaults(
  defineProps<{
    count?: number;
    divided?: boolean;
  }>(),
  {
    count: 4,
    divided: false,
  },
);
</script>

<template>
  <div
    v-if="!divided"
    class="flex overflow-x-auto snap-x snap-mandatory gap-2.5 no-scrollbar sm:grid sm:grid-cols-2 lg:grid-cols-4 sm:gap-3"
  >
    <div
      v-for="i in count"
      :key="i"
      class="min-w-[74%] sm:min-w-0 snap-start shrink-0 rounded-[3px] border border-border bg-card px-4 py-3 space-y-2.5"
    >
      <div class="flex items-center justify-between">
        <Skeleton class="h-7 w-7 rounded-lg" />
        <Skeleton class="h-4 w-16 rounded" />
      </div>
      <Skeleton class="h-3 w-28 mt-2" />
      <div class="flex items-baseline justify-between gap-2 mt-1">
        <Skeleton class="h-6 w-32" />
        <Skeleton class="h-3.5 w-10 rounded" />
      </div>
    </div>
  </div>

  <section
    v-else
    class="flex overflow-x-auto snap-x snap-mandatory gap-2.5 no-scrollbar sm:grid sm:grid-cols-2 lg:grid-cols-4 sm:gap-0 sm:divide-x sm:divide-border sm:rounded-[3px] sm:border sm:border-border sm:bg-card"
  >
    <div
      v-for="i in count"
      :key="i"
      class="min-w-[74%] sm:min-w-0 snap-start shrink-0 rounded-[3px] border border-border bg-card sm:border-0 sm:rounded-none p-3.5 sm:p-4 space-y-2"
    >
      <Skeleton class="h-3 w-28" />
      <Skeleton class="h-6 w-24 my-1" />
      <Skeleton class="h-3 w-36" />
    </div>
  </section>
</template>

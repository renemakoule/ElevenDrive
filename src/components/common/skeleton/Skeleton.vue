<script setup lang="ts">
withDefaults(
  defineProps<{
    className?: string;
  }>(),
  {
    className: "",
  },
);
</script>

<template>
  <div
    :class="[
      'animate-pulse rounded-[3px] bg-muted/70 dark:bg-muted/40 transition-colors',
      className,
    ]"
    aria-hidden="true"
  />
</template>

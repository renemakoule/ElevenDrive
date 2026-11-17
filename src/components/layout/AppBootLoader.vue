<script setup lang="ts">
import { Loader2 } from "lucide-vue-next";
</script>

<template>
  <div
    class="flex-1 w-full h-full min-h-[70vh] flex flex-col items-center justify-center bg-background text-foreground px-4 select-none"
    role="status"
    aria-live="polite"
  >
    <div class="flex flex-col items-center justify-center gap-3 text-center">
      <div
        class="flex h-11 w-11 items-center justify-center rounded-[3px] border border-brand-border bg-brand-light text-brand-dark"
      >
        <Loader2 class="h-5 w-5 animate-spin text-brand" />
      </div>
      <div class="flex flex-col items-center justify-center text-center w-full">
        <p class="text-sm font-semibold tracking-tight text-foreground text-center w-full">
          ElevenDrive
        </p>
        <p class="text-[11px] text-muted-foreground text-center w-full mt-0.5">
          Chargement de l'interface...
        </p>
      </div>
    </div>
  </div>
</template>

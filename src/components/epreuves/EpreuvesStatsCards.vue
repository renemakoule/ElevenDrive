<script setup lang="ts">
defineProps<{
  stats: {
    total: number;
    permisB: number;
    permisA: number;
    poidsLourd: number;
    avecCorrige: number;
    totalDownloads: number;
  };
}>();
</script>

<template>
  <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-3">
    <div class="p-3.5 rounded-[4px] bg-card border border-border">
      <span class="text-[11px] text-muted-foreground font-medium">Total Sujets</span>
      <p class="text-xl font-bold tabular-nums text-foreground mt-1">{{ stats.total }}</p>
    </div>

    <div class="p-3.5 rounded-[4px] bg-card border border-border">
      <span class="text-[11px] text-emerald-600 dark:text-emerald-400 font-semibold">Permis B (Auto)</span>
      <p class="text-xl font-bold tabular-nums text-foreground mt-1">{{ stats.permisB }}</p>
    </div>

    <div class="p-3.5 rounded-[4px] bg-card border border-border">
      <span class="text-[11px] text-amber-600 dark:text-amber-400 font-semibold">Permis A (Motos)</span>
      <p class="text-xl font-bold tabular-nums text-foreground mt-1">{{ stats.permisA }}</p>
    </div>

    <div class="p-3.5 rounded-[4px] bg-card border border-border">
      <span class="text-[11px] text-blue-600 dark:text-blue-400 font-semibold">Poids Lourd (C/D/E)</span>
      <p class="text-xl font-bold tabular-nums text-foreground mt-1">{{ stats.poidsLourd }}</p>
    </div>

    <div class="p-3.5 rounded-[4px] bg-card border border-border">
      <span class="text-[11px] text-muted-foreground font-medium">Avec Corrigé</span>
      <p class="text-xl font-bold tabular-nums text-foreground mt-1">{{ stats.avecCorrige }}</p>
    </div>

    <div class="p-3.5 rounded-[4px] bg-card border border-border">
      <span class="text-[11px] text-muted-foreground font-medium">Téléchargements</span>
      <p class="text-xl font-bold tabular-nums text-foreground mt-1">{{ stats.totalDownloads }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Star } from "lucide-vue-next";

defineProps<{
  averageStars: string | number;
  totalCount: number;
  filteredCount: number;
  starCounts: Record<1 | 2 | 3 | 4 | 5, number>;
  selectedFilterStars: number | "all";
}>();

const emit = defineEmits<{
  (e: "update:selectedFilterStars", val: number | "all"): void;
}>();
</script>

<template>
  <div class="space-y-4">
    <!-- Synthèse des Notations -->
    <div class="grid grid-cols-1 md:grid-cols-12 gap-4">
      <div
        class="md:col-span-4 p-5 rounded-[3px] border border-border bg-card flex flex-col items-center justify-center text-center space-y-2"
      >
        <span
          class="text-[11px] font-bold text-muted-foreground uppercase tracking-wider"
        >
          Indice de Satisfaction Global
        </span>
        <div class="flex items-baseline gap-2">
          <span
            class="text-4xl font-extrabold font-mono text-foreground tracking-tight"
          >
            {{ averageStars }}
          </span>
          <span class="text-sm text-muted-foreground font-semibold">/ 5.0</span>
        </div>
        <div class="flex items-center gap-1 text-amber-400">
          <Star
            v-for="i in 5"
            :key="i"
            class="w-4 h-4"
            :class="
              i <= Math.round(Number(averageStars))
                ? 'fill-amber-400 text-amber-400'
                : 'text-muted-foreground/30'
            "
          />
        </div>
        <span class="text-[11px] text-muted-foreground">
          Basé sur {{ totalCount }} avis certifiés
        </span>
      </div>

      <div
        class="md:col-span-8 p-5 rounded-[3px] border border-border bg-card space-y-2"
      >
        <span
          class="text-[11px] font-bold text-muted-foreground uppercase tracking-wider block mb-2"
        >
          Répartition des Notes
        </span>
        <div
          v-for="star in [5, 4, 3, 2, 1] as const"
          :key="star"
          class="flex items-center gap-3 text-xs"
        >
          <span
            class="w-10 font-bold flex items-center gap-1 shrink-0 text-foreground"
          >
            {{ star }} <Star class="w-3 h-3 fill-amber-400 text-amber-400" />
          </span>
          <div class="flex-1 h-2 rounded-full bg-muted overflow-hidden">
            <div
              class="h-full bg-amber-400 rounded-full transition-all duration-300"
              :style="{
                width: totalCount
                  ? `${(starCounts[star] / totalCount) * 100}%`
                  : '0%',
              }"
            ></div>
          </div>
          <span
            class="w-8 text-right font-mono text-muted-foreground text-[11px] shrink-0"
          >
            {{ starCounts[star] }}
          </span>
        </div>
      </div>
    </div>

    <!-- Filtres par étoiles -->
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-1.5 text-xs">
        <span class="text-muted-foreground text-[11px] font-semibold mr-1"
          >Filtrer par note :</span
        >
        <button
          type="button"
          @click="emit('update:selectedFilterStars', 'all')"
          class="px-2.5 py-1 rounded-[4px] text-xs font-medium transition cursor-pointer"
          :class="
            selectedFilterStars === 'all'
              ? 'bg-brand text-white font-semibold'
              : 'bg-muted text-muted-foreground hover:text-foreground'
          "
        >
          Tous ({{ totalCount }})
        </button>
        <button
          v-for="s in [5, 4, 3, 2, 1]"
          :key="s"
          type="button"
          @click="emit('update:selectedFilterStars', s)"
          class="px-2 py-1 rounded-[4px] text-xs font-medium transition cursor-pointer flex items-center gap-1"
          :class="
            selectedFilterStars === s
              ? 'bg-brand text-white font-semibold'
              : 'bg-muted text-muted-foreground hover:text-foreground'
          "
        >
          <span>{{ s }}</span>
          <Star class="w-3 h-3 fill-amber-400 text-amber-400" />
        </button>
      </div>

      <span class="text-xs text-muted-foreground">
        {{ filteredCount }} avis affiché{{ filteredCount > 1 ? "s" : "" }}
      </span>
    </div>
  </div>
</template>

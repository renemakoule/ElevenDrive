<script setup lang="ts">
import {
  LayoutList,
  CalendarRange,
  ChevronLeft,
  ChevronRight,
  Calendar,
} from "lucide-vue-next";

defineProps<{
  viewMode: "jour" | "semaine";
  selectedDate: string;
  todayIso: string;
  dateLongue: string;
  weekDays: Array<{
    date: string;
    dayName: string;
    dayNumber: number;
    isToday: boolean;
    isSelected: boolean;
    count: number;
  }>;
  stats: Array<{
    key: string;
    label: string;
    value: string;
    tag?: string;
    progress?: number;
    caption: string;
  }>;
}>();

const emit = defineEmits<{
  "update:viewMode": [mode: "jour" | "semaine"];
  prevDay: [];
  nextDay: [];
  resetToday: [];
  selectDate: [iso: string];
}>();
</script>

<template>
  <div class="space-y-4">
    <!-- ═════════ BARRE DE NAVIGATION TEMPORELLE & CALENDRIER ═════════ -->
    <div class="space-y-2.5 rounded-[4px] border border-border bg-card p-3">
      <!-- Ligne 1 : Commandes principales et Commutateur Jour / Semaine -->
      <div class="flex flex-wrap items-center justify-between gap-2.5">
        <!-- Sélecteur de mode d'affichage -->
        <div
          class="inline-flex rounded-[4px] border border-border bg-muted/50 p-0.5"
        >
          <button
            type="button"
            @click="emit('update:viewMode', 'jour')"
            class="inline-flex items-center gap-1.5 rounded-[4px] px-3 py-1 text-xs font-medium transition"
            :class="
              viewMode === 'jour'
                ? 'bg-background text-foreground font-semibold'
                : 'text-muted-foreground hover:text-foreground'
            "
          >
            <LayoutList class="h-3.5 w-3.5" />
            <span>Vue Journée</span>
          </button>
          <button
            type="button"
            @click="emit('update:viewMode', 'semaine')"
            class="inline-flex items-center gap-1.5 rounded-[4px] px-3 py-1 text-xs font-medium transition"
            :class="
              viewMode === 'semaine'
                ? 'bg-background text-foreground font-semibold'
                : 'text-muted-foreground hover:text-foreground'
            "
          >
            <CalendarRange class="h-3.5 w-3.5" />
            <span>Vue Semaine</span>
          </button>
        </div>

        <!-- Contrôles de dates : Précédent, Suivant, Aujourd'hui, Input Date -->
        <div class="flex items-center gap-1.5 text-xs">
          <button
            type="button"
            @click="emit('prevDay')"
            title="Jour précédent"
            class="inline-flex h-8 w-8 items-center justify-center rounded-[4px] border border-border bg-background text-foreground hover:bg-muted transition"
          >
            <ChevronLeft class="h-4 w-4" />
          </button>

          <button
            v-if="selectedDate !== todayIso"
            type="button"
            @click="emit('resetToday')"
            class="inline-flex h-8 items-center gap-1 rounded-[4px] border border-brand-border bg-brand-light px-2.5 text-xs font-semibold text-brand-dark hover:bg-brand-light/80 transition"
          >
            Aujourd'hui
          </button>

          <!-- Input Date natif déguisé pour choisir n'importe quel jour -->
          <label
            class="relative inline-flex h-8 items-center gap-1.5 rounded-[4px] border border-border bg-background px-2.5 font-medium text-foreground hover:bg-muted cursor-pointer transition"
          >
            <Calendar class="h-3.5 w-3.5 text-brand" />
            <span class="capitalize text-xs">{{ dateLongue }}</span>
            <input
              type="date"
              :value="selectedDate"
              @change="
                (e) => emit('selectDate', (e.target as HTMLInputElement).value)
              "
              class="absolute inset-0 opacity-0 cursor-pointer w-full h-full"
            />
          </label>

          <button
            type="button"
            @click="emit('nextDay')"
            title="Jour suivant"
            class="inline-flex h-8 w-8 items-center justify-center rounded-[4px] border border-border bg-background text-foreground hover:bg-muted transition"
          >
            <ChevronRight class="h-4 w-4" />
          </button>
        </div>
      </div>

      <!-- Ligne 2 : Carrousel des 7 jours de la semaine (Pills interactives) -->
      <div class="grid grid-cols-7 gap-1.5 pt-1 border-t border-border/50">
        <button
          v-for="d in weekDays"
          :key="d.date"
          type="button"
          @click="emit('selectDate', d.date)"
          class="relative flex flex-col items-center justify-center rounded-[4px] py-2 px-1 text-center transition select-none"
          :class="[
            d.isSelected
              ? 'bg-brand text-white ring-2 ring-brand/40'
              : 'border border-border/70 bg-muted/20 hover:bg-muted hover:border-foreground/20 text-foreground',
          ]"
        >
          <!-- Badge Aujourd'hui -->
          <span
            v-if="d.isToday && !d.isSelected"
            class="absolute -top-1 rounded-full bg-brand px-1 text-[8px] font-bold text-white uppercase"
          >
            Auj.
          </span>

          <span
            class="text-[10px] font-semibold tracking-wider uppercase"
            :class="d.isSelected ? 'text-white/90' : 'text-muted-foreground'"
          >
            {{ d.dayName }}
          </span>

          <span
            class="text-sm font-bold tabular-nums mt-0.5"
            :class="d.isSelected ? 'text-white' : 'text-foreground'"
          >
            {{ d.dayNumber }}
          </span>

          <!-- Indicateur de séances -->
          <span
            v-if="d.count > 0"
            class="mt-1 h-1.5 w-1.5 rounded-full"
            :class="d.isSelected ? 'bg-white' : 'bg-brand'"
          ></span>
          <span
            v-else
            class="mt-1 h-1.5 w-1.5 rounded-full bg-transparent"
          ></span>
        </button>
      </div>
    </div>

    <!-- ═════════ Indicateurs KPIs Réels (Carrousel horizontal tactile sur mobile) ═════════ -->
    <section
      class="flex overflow-x-auto snap-x snap-mandatory gap-2.5 no-scrollbar md:grid md:grid-cols-3 md:gap-0 md:divide-x md:divide-border md:rounded-[3px] md:border md:border-border md:bg-card"
    >
      <div
        v-for="s in stats"
        :key="s.key"
        class="min-w-[78%] sm:min-w-[260px] md:min-w-0 snap-start rounded-[3px] border border-border bg-card md:border-0 md:rounded-none p-3.5 md:p-4 shrink-0"
      >
        <div class="flex items-center justify-between gap-2">
          <span class="text-xs text-muted-foreground">{{ s.label }}</span>
          <span
            v-if="s.tag"
            class="rounded bg-muted px-1.5 py-0.5 font-mono text-[11px] font-medium text-foreground ring-1 ring-border"
          >
            {{ s.tag }}
          </span>
        </div>

        <p
          class="mt-1 truncate text-lg font-semibold tabular-nums text-foreground"
        >
          {{ s.value }}
        </p>

        <div
          v-if="s.progress !== undefined"
          class="mt-2.5 h-1 overflow-hidden rounded-full bg-muted"
          role="progressbar"
          :aria-valuenow="s.progress"
          aria-valuemin="0"
          aria-valuemax="100"
          :aria-label="s.label"
        >
          <div
            class="h-full rounded-full bg-brand transition-all duration-500"
            :style="{ width: s.progress + '%' }"
          />
        </div>

        <p
          class="text-[11px] text-muted-foreground"
          :class="s.progress !== undefined ? 'mt-2' : 'mt-3.5'"
        >
          {{ s.caption }}
        </p>
      </div>
    </section>
  </div>
</template>

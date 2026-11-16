<script setup lang="ts">
import type { HistoriqueItem } from "./types";
import {
  History,
  Receipt,
  Car,
  Video,
  Brain,
  ShieldCheck,
  Users,
  Sparkles,
  Search,
} from "lucide-vue-next";

defineProps<{
  items: HistoriqueItem[];
  searchQuery: string;
  selectedFilter: string;
  currentRole: string;
}>();

const emit = defineEmits<{
  (e: "update:searchQuery", val: string): void;
  (e: "update:selectedFilter", val: string): void;
}>();

const FILTER_TABS = [
  { id: "ALL", label: "Tous", activeClass: "bg-brand text-brand-foreground" },
  {
    id: "PAIEMENT",
    label: "Paiements",
    activeClass: "bg-emerald-600 text-white",
  },
  { id: "CONDUITE", label: "Conduite", activeClass: "bg-blue-600 text-white" },
  {
    id: "CODE",
    label: "Code & Visio",
    activeClass: "bg-indigo-600 text-white",
  },
  { id: "QUIZ", label: "Quiz", activeClass: "bg-purple-600 text-white" },
];
</script>

<template>
  <div class="space-y-4">
    <!-- Barre de recherche et filtres -->
    <div
      class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3 bg-card border border-border p-3 rounded-[4px]"
    >
      <div class="relative flex-1">
        <Search
          class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground"
        />
        <input
          :value="searchQuery"
          @input="
            emit(
              'update:searchQuery',
              ($event.target as HTMLInputElement).value,
            )
          "
          type="text"
          placeholder="Rechercher par candidat, montant, référence, thème..."
          class="w-full pl-9 pr-3 py-1.5 text-xs rounded-[4px] border border-border bg-background text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-brand"
        />
      </div>

      <div class="flex items-center gap-1.5 overflow-x-auto pb-1 sm:pb-0">
        <button
          v-for="tab in FILTER_TABS"
          :key="tab.id"
          type="button"
          @click="emit('update:selectedFilter', tab.id)"
          class="px-2.5 py-1 rounded-[4px] text-xs font-medium whitespace-nowrap transition"
          :class="
            selectedFilter === tab.id
              ? tab.activeClass
              : 'bg-muted/50 text-muted-foreground hover:bg-muted hover:text-foreground'
          "
        >
          {{ tab.label }}
        </button>

        <button
          v-if="currentRole === 'manager' || currentRole === 'directeur'"
          type="button"
          @click="emit('update:selectedFilter', 'DEROGATION')"
          class="px-2.5 py-1 rounded-[4px] text-xs font-medium whitespace-nowrap transition"
          :class="
            selectedFilter === 'DEROGATION'
              ? 'bg-amber-600 text-white'
              : 'bg-muted/50 text-muted-foreground hover:bg-muted hover:text-foreground'
          "
        >
          Dérogations
        </button>
      </div>
    </div>

    <!-- Flux chronologique -->
    <div class="rounded-[4px] border border-border bg-card overflow-hidden">
      <div
        class="p-3 bg-muted/40 border-b border-border flex items-center justify-between text-xs font-semibold text-muted-foreground"
      >
        <span>Événements & Actions Chronologiques</span>
        <span>{{ items.length }} résultat(s)</span>
      </div>

      <div v-if="items.length === 0" class="p-12 text-center">
        <History class="h-10 w-10 text-muted-foreground/40 mx-auto mb-3" />
        <p class="text-sm font-medium text-foreground">
          Aucun événement d'activité trouvé
        </p>
        <p class="text-xs text-muted-foreground mt-1">
          Aucune action ne correspond à vos filtres de recherche actuels.
        </p>
      </div>

      <div v-else class="divide-y divide-border">
        <div
          v-for="item in items"
          :key="item.id"
          class="p-4 transition hover:bg-muted/30 flex items-start gap-3.5 relative"
          :class="{
            'bg-emerald-500/5 dark:bg-emerald-500/10': item.isRealtimeLive,
          }"
        >
          <div
            class="h-9 w-9 rounded-[4px] shrink-0 flex items-center justify-center font-bold text-white"
            :class="{
              'bg-emerald-600': item.type === 'PAIEMENT',
              'bg-blue-600': item.type === 'CONDUITE',
              'bg-indigo-600': item.type === 'CODE',
              'bg-purple-600': item.type === 'QUIZ',
              'bg-amber-600':
                item.type === 'DEROGATION' || item.type === 'CAISSE',
              'bg-orange-600': item.type === 'FLOTTE',
              'bg-brand': item.type === 'INSCRIPTION' || item.type === 'EXAMEN',
            }"
          >
            <Receipt v-if="item.type === 'PAIEMENT'" class="h-4 w-4" />
            <Car v-else-if="item.type === 'CONDUITE'" class="h-4 w-4" />
            <Video v-else-if="item.type === 'CODE'" class="h-4 w-4" />
            <Brain v-else-if="item.type === 'QUIZ'" class="h-4 w-4" />
            <ShieldCheck
              v-else-if="item.type === 'DEROGATION'"
              class="h-4 w-4"
            />
            <Car v-else-if="item.type === 'FLOTTE'" class="h-4 w-4" />
            <History v-else class="h-4 w-4" />
          </div>

          <div class="flex-1 min-w-0">
            <div
              class="flex flex-wrap items-center justify-between gap-1.5 mb-1"
            >
              <div class="flex items-center gap-2 flex-wrap">
                <h3 class="text-xs sm:text-sm font-semibold text-foreground">
                  {{ item.titre }}
                </h3>
                <span
                  v-if="item.isRealtimeLive"
                  class="inline-flex items-center gap-1 px-1.5 py-0.5 rounded text-[10px] font-bold bg-emerald-500 text-white animate-pulse"
                >
                  <Sparkles class="h-2.5 w-2.5" />
                  EN DIRECT
                </span>
              </div>
              <span class="text-[11px] text-muted-foreground whitespace-nowrap">
                {{ item.dateStr }} • {{ item.heureStr }}
              </span>
            </div>

            <p class="text-xs text-muted-foreground">
              {{ item.description }}
            </p>

            <div class="flex flex-wrap items-center gap-2 mt-2 pt-1">
              <span
                class="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-medium border"
                :class="item.badgeClass"
              >
                {{ item.badgeLabel }}
              </span>

              <span
                v-if="item.acteur"
                class="text-[11px] text-muted-foreground flex items-center gap-1"
              >
                <Users class="h-3 w-3" />
                {{ item.acteur }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

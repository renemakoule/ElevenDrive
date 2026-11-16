<script setup lang="ts">
import type { Vehicule } from "../../../types";
import { CATEGORIES_PERMIS } from "./planningConduiteUtils";
import {
  Calendar,
  Users,
  Car,
  CheckCircle2,
  Search,
} from "lucide-vue-next";

defineProps<{
  stats: {
    total: number;
    totalHeures: number;
    uniqueMoniteurs: number;
    uniqueVehicules: number;
    effectuees: number;
  };
  searchQuery: string;
  selectedCategorie: string;
  selectedMoniteurId: string;
  selectedVehiculeId: string;
  selectedStatut: string;
  moniteurs: { id: string; nomComplet: string; role: string; telephone?: string }[];
  vehicules: Vehicule[];
  filteredCount: number;
}>();

const emit = defineEmits<{
  (e: "update:searchQuery", val: string): void;
  (e: "update:selectedCategorie", val: string): void;
  (e: "update:selectedMoniteurId", val: string): void;
  (e: "update:selectedVehiculeId", val: string): void;
  (e: "update:selectedStatut", val: string): void;
}>();
</script>

<template>
  <div class="space-y-4">
    <!-- Cartes d'indicateurs KPI -->
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
      <div class="p-3.5 rounded-[4px] bg-card border border-border">
        <span
          class="text-[11px] text-muted-foreground font-medium flex items-center gap-1"
        >
          <Calendar class="w-3.5 h-3.5 text-brand" /> Leçons Programmées
        </span>
        <p class="text-xl font-bold tabular-nums text-foreground mt-1">
          {{ stats.total }}
          <span class="text-xs font-normal text-muted-foreground"
            >({{ stats.totalHeures }}h)</span
          >
        </p>
      </div>

      <div class="p-3.5 rounded-[4px] bg-card border border-border">
        <span
          class="text-[11px] text-sky-600 dark:text-sky-400 font-semibold flex items-center gap-1"
        >
          <Users class="w-3.5 h-3.5" /> Moniteurs Mobilisés
        </span>
        <p class="text-xl font-bold tabular-nums text-foreground mt-1">
          {{ stats.uniqueMoniteurs }}
          <span class="text-xs font-normal text-muted-foreground"
            >instructeurs</span
          >
        </p>
      </div>

      <div class="p-3.5 rounded-[4px] bg-card border border-border">
        <span
          class="text-[11px] text-emerald-600 dark:text-emerald-400 font-semibold flex items-center gap-1"
        >
          <Car class="w-3.5 h-3.5" /> Véhicules en Rotation
        </span>
        <p class="text-xl font-bold tabular-nums text-foreground mt-1">
          {{ stats.uniqueVehicules }}
          <span class="text-xs font-normal text-muted-foreground"
            >voitures</span
          >
        </p>
      </div>

      <div class="p-3.5 rounded-[4px] bg-card border border-border">
        <span
          class="text-[11px] text-amber-600 dark:text-amber-400 font-semibold flex items-center gap-1"
        >
          <CheckCircle2 class="w-3.5 h-3.5" /> Séances Effectuées
        </span>
        <p class="text-xl font-bold tabular-nums text-foreground mt-1">
          {{ stats.effectuees }} / {{ stats.total }}
        </p>
      </div>
    </div>

    <!-- Barre de filtres avancés -->
    <div
      class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 p-3 rounded-[4px] bg-card border border-border"
    >
      <div class="flex flex-wrap items-center gap-2 flex-1">
        <div class="relative min-w-[200px] flex-1 sm:flex-initial">
          <Search
            class="w-3.5 h-3.5 absolute left-2.5 top-2.5 text-muted-foreground"
          />
          <input
            type="text"
            :value="searchQuery"
            @input="
              emit('update:searchQuery', ($event.target as HTMLInputElement).value)
            "
            placeholder="Rechercher élève, moniteur..."
            class="w-full pl-8 pr-3 py-1.5 text-xs rounded-[4px] bg-background border border-border focus:border-brand outline-none"
          />
        </div>

        <select
          :value="selectedCategorie"
          @change="
            emit(
              'update:selectedCategorie',
              ($event.target as HTMLSelectElement).value,
            )
          "
          class="px-2.5 py-1.5 text-xs rounded-[4px] bg-background border border-border focus:border-brand outline-none cursor-pointer"
        >
          <option
            v-for="cat in CATEGORIES_PERMIS"
            :key="cat.code"
            :value="cat.code"
          >
            {{ cat.label }}
          </option>
        </select>

        <select
          :value="selectedMoniteurId"
          @change="
            emit(
              'update:selectedMoniteurId',
              ($event.target as HTMLSelectElement).value,
            )
          "
          class="px-2.5 py-1.5 text-xs rounded-[4px] bg-background border border-border focus:border-brand outline-none cursor-pointer"
        >
          <option value="">Tous les moniteurs ({{ moniteurs.length }})</option>
          <option v-for="m in moniteurs" :key="m.id" :value="m.id">
            👤 {{ m.nomComplet }}
          </option>
        </select>

        <select
          :value="selectedVehiculeId"
          @change="
            emit(
              'update:selectedVehiculeId',
              ($event.target as HTMLSelectElement).value,
            )
          "
          class="px-2.5 py-1.5 text-xs rounded-[4px] bg-background border border-border focus:border-brand outline-none cursor-pointer"
        >
          <option value="">Tous les véhicules ({{ vehicules.length }})</option>
          <option v-for="v in vehicules" :key="v.id" :value="v.id">
            🚗 {{ v.immatriculation }} ({{ v.marque_modele }})
          </option>
        </select>

        <select
          :value="selectedStatut"
          @change="
            emit(
              'update:selectedStatut',
              ($event.target as HTMLSelectElement).value,
            )
          "
          class="px-2.5 py-1.5 text-xs rounded-[4px] bg-background border border-border focus:border-brand outline-none cursor-pointer"
        >
          <option value="ALL">Tous les statuts</option>
          <option value="PROGRAMME">Programmées</option>
          <option value="EN_COURS">En cours</option>
          <option value="EFFECTUE">Effectuées</option>
          <option value="ANNULE_ELEVE">Annulées</option>
        </select>
      </div>

      <div
        class="text-[11px] text-muted-foreground font-mono self-end sm:self-center"
      >
        {{ filteredCount }} séance(s) trouvée(s)
      </div>
    </div>
  </div>
</template>


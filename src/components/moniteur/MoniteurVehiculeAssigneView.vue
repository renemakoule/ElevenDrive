<script setup lang="ts">
import type { Vehicule } from "../../types";
import { Car } from "lucide-vue-next";

defineProps<{
  vehiculesDisponibles: Vehicule[];
  vehiculeAssigne: Vehicule | null;
  vehiculeInfo: {
    id?: string;
    immatriculation: string;
    modele: string;
    visiteTechnique: string;
    assurance: string;
    kilometrage: number;
  };
  vehiculeChoisiId: string;
  isUpdatingStatut: boolean;
}>();

const emit = defineEmits<{
  "update:vehiculeChoisiId": [id: string];
  changeStatut: [
    vehiculeId: string,
    statut: "DISPONIBLE" | "EN_MAINTENANCE" | "EN_PANNE",
  ];
}>();

const num = (n: number) => n.toLocaleString("fr-FR");
</script>

<template>
  <div class="space-y-4">
    <div class="grid grid-cols-1 gap-4 lg:grid-cols-12">
      <!-- Détails du véhicule actif et changement de statut temps réel -->
      <section
        class="rounded-[4px] border border-border bg-card p-5 space-y-4 lg:col-span-7"
      >
        <div
          class="flex flex-wrap items-center justify-between gap-3 border-b border-border pb-3"
        >
          <div class="flex items-center gap-2.5">
            <div
              class="flex h-10 w-10 items-center justify-center rounded-[4px] bg-brand-light text-primary"
            >
              <Car class="h-5 w-5" />
            </div>
            <div>
              <h2 class="text-sm font-bold text-foreground">
                {{ vehiculeInfo.modele }}
              </h2>
              <p class="font-mono text-xs text-muted-foreground">
                Immatriculation : {{ vehiculeInfo.immatriculation }}
              </p>
            </div>
          </div>

          <select
            v-if="vehiculesDisponibles.length > 1"
            :value="vehiculeChoisiId || vehiculeAssigne?.id || ''"
            @change="
              (e) =>
                emit(
                  'update:vehiculeChoisiId',
                  (e.target as HTMLSelectElement).value,
                )
            "
            class="h-8 rounded-[4px] border border-border bg-background px-2.5 text-xs font-medium text-foreground"
          >
            <option v-for="v in vehiculesDisponibles" :key="v.id" :value="v.id">
              {{ v.immatriculation }} — {{ v.marque_modele }} ({{ v.statut }})
            </option>
          </select>
        </div>

        <div
          v-if="vehiculeAssigne"
          class="grid grid-cols-2 gap-3 sm:grid-cols-3"
        >
          <div class="rounded-[4px] border border-border bg-muted/20 p-3">
            <span class="text-[11px] text-muted-foreground"
              >Compteur Kilométrique</span
            >
            <p class="mt-1 font-mono text-lg font-bold text-foreground">
              {{ num(vehiculeInfo.kilometrage) }} km
            </p>
          </div>
          <div class="rounded-xl border border-border bg-muted/20 p-3">
            <span class="text-[11px] text-muted-foreground"
              >Visite Technique</span
            >
            <p class="mt-1 text-xs font-semibold text-foreground">
              {{ vehiculeInfo.visiteTechnique }}
            </p>
          </div>
          <div class="rounded-[4px] border border-border bg-muted/20 p-3">
            <span class="text-[11px] text-muted-foreground">Assurance</span>
            <p class="mt-1 text-xs font-semibold text-brand-dark">
              {{ vehiculeInfo.assurance }}
            </p>
          </div>
        </div>

        <!-- Commutateur de statut en temps réel (WebSocket + PostgreSQL) -->
        <div
          v-if="vehiculeAssigne"
          class="space-y-2 pt-2 border-t border-border"
        >
          <span class="text-xs font-semibold text-foreground block">
            Mettre à jour la disponibilité du véhicule en temps réel :
          </span>
          <div class="grid grid-cols-3 gap-2">
            <button
              type="button"
              :disabled="isUpdatingStatut"
              @click="emit('changeStatut', vehiculeAssigne.id, 'DISPONIBLE')"
              class="rounded-[4px] border px-3 py-2 text-xs font-bold transition cursor-pointer"
              :class="
                vehiculeAssigne.statut === 'DISPONIBLE'
                  ? 'border-emerald-600 bg-emerald-500/15 text-emerald-700 dark:text-emerald-300'
                  : 'border-border bg-background text-muted-foreground hover:bg-muted'
              "
            >
              ✅ Disponible
            </button>
            <button
              type="button"
              :disabled="isUpdatingStatut"
              @click="
                emit('changeStatut', vehiculeAssigne.id, 'EN_MAINTENANCE')
              "
              class="rounded-[4px] border px-3 py-2 text-xs font-bold transition cursor-pointer"
              :class="
                vehiculeAssigne.statut === 'EN_MAINTENANCE'
                  ? 'border-amber-500 bg-amber-500/15 text-amber-700 dark:text-amber-300'
                  : 'border-border bg-background text-muted-foreground hover:bg-muted'
              "
            >
              🔧 En maintenance
            </button>
            <button
              type="button"
              :disabled="isUpdatingStatut"
              @click="emit('changeStatut', vehiculeAssigne.id, 'EN_PANNE')"
              class="rounded-[4px] border px-3 py-2 text-xs font-bold transition cursor-pointer"
              :class="
                vehiculeAssigne.statut === 'EN_PANNE'
                  ? 'border-destructive bg-destructive/15 text-destructive'
                  : 'border-border bg-background text-muted-foreground hover:bg-muted'
              "
            >
              ⚠️ En panne
            </button>
          </div>
        </div>
      </section>

      <!-- Liste complète du parc automobile de l'auto-école en temps réel -->
      <section
        class="rounded-[4px] border border-border bg-card p-5 space-y-3 lg:col-span-5"
      >
        <div class="flex items-center justify-between">
          <h3
            class="text-xs font-bold uppercase tracking-wider text-foreground"
          >
            Flotte de l'auto-école ({{ vehiculesDisponibles.length }})
          </h3>
          <span class="text-[11px] text-muted-foreground"
            >Sélectionner un véhicule</span
          >
        </div>

        <div
          v-if="vehiculesDisponibles.length"
          class="space-y-2 max-h-[360px] overflow-y-auto pr-1"
        >
          <button
            v-for="v in vehiculesDisponibles"
            :key="v.id"
            type="button"
            @click="emit('update:vehiculeChoisiId', v.id)"
            class="w-full text-left rounded-[4px] border p-3 transition cursor-pointer flex items-center justify-between gap-3"
            :class="
              vehiculeAssigne?.id === v.id
                ? 'border-brand bg-brand-light/25 ring-1 ring-brand'
                : 'border-border bg-muted/20 hover:bg-muted/50'
            "
          >
            <div>
              <div class="text-xs font-bold text-foreground">
                {{ v.marque_modele }}
              </div>
              <div class="font-mono text-[11px] text-muted-foreground">
                {{ v.immatriculation }} · {{ v.type_boite }} ·
                {{ num(v.kilometrage_actuel || 0) }} km
              </div>
            </div>
            <span
              class="rounded-full px-2 py-0.5 text-[10px] font-bold"
              :class="
                v.statut === 'DISPONIBLE'
                  ? 'bg-emerald-500/15 text-emerald-700 dark:text-emerald-300'
                  : v.statut === 'EN_MAINTENANCE'
                    ? 'bg-amber-500/15 text-amber-700'
                    : 'bg-destructive/15 text-destructive'
              "
            >
              {{ v.statut }}
            </span>
          </button>
        </div>
        <div
          v-else
          class="rounded-[4px] border border-dashed border-border p-6 text-center text-xs text-muted-foreground"
        >
          Aucun véhicule enregistré dans la base de données.
        </div>
      </section>
    </div>
  </div>
</template>

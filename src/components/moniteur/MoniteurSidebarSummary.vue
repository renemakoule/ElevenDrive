<script setup lang="ts">
import type { Vehicule } from "../../types";
import { Phone, Car, Wrench, ShieldCheck } from "lucide-vue-next";

defineProps<{
  activeMoniteur: {
    id: string;
    nom: string;
    titre: string;
    telephone: string;
  };
  initiales: string;
  vehiculesDisponibles: Vehicule[];
  vehiculeChoisiId: string;
  vehiculeInfo: {
    id?: string;
    immatriculation: string;
    modele: string;
    visiteTechnique: string;
    assurance: string;
    kilometrage: number;
  };
}>();

const emit = defineEmits<{
  "update:vehiculeChoisiId": [id: string];
}>();
</script>

<template>
  <aside class="space-y-4 lg:col-span-4">
    <!-- Carte Moniteur -->
    <section class="rounded-[4px] border border-border bg-card p-4">
      <div class="flex items-center gap-3">
        <div
          class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-brand-light text-xs font-semibold text-brand-dark"
          aria-hidden="true"
        >
          {{ initiales }}
        </div>
        <div class="min-w-0 flex-1">
          <h2 class="truncate text-sm font-semibold text-foreground">
            {{ activeMoniteur.nom }}
          </h2>
          <p class="truncate text-xs text-muted-foreground">
            {{ activeMoniteur.titre }}
          </p>
        </div>
      </div>

      <dl class="mt-3 space-y-2 border-t border-border pt-3 text-xs">
        <div class="flex items-center justify-between gap-3">
          <dt class="flex items-center gap-1.5 text-muted-foreground">
            <Phone class="h-3.5 w-3.5" />
            Téléphone
          </dt>
          <dd>
            <a
              :href="`tel:${activeMoniteur.telephone}`"
              class="font-medium text-foreground hover:text-brand-dark"
            >
              {{ activeMoniteur.telephone }}
            </a>
          </dd>
        </div>
      </dl>
    </section>

    <!-- Carte Véhicule Assigné Réel -->
    <section class="rounded-[4px] border border-border bg-card p-4 space-y-3">
      <div class="flex items-center justify-between">
        <h3
          class="text-xs font-semibold text-foreground flex items-center gap-1.5"
        >
          <Car class="h-3.5 w-3.5 text-brand" />
          Véhicule d'apprentissage
        </h3>

        <!-- Sélecteur si plusieurs véhicules -->
        <select
          v-if="vehiculesDisponibles.length > 1"
          :value="vehiculeChoisiId || vehiculeInfo.id || ''"
          @change="
            (e) =>
              emit(
                'update:vehiculeChoisiId',
                (e.target as HTMLSelectElement).value,
              )
          "
          class="h-6 rounded border border-border bg-muted/40 px-1 text-[11px] text-foreground focus:border-brand focus:outline-none"
        >
          <option v-for="v in vehiculesDisponibles" :key="v.id" :value="v.id">
            {{ v.immatriculation }} - {{ v.marque_modele }}
          </option>
        </select>
      </div>

      <dl class="space-y-2 text-xs">
        <div class="flex items-center justify-between gap-3">
          <dt class="text-muted-foreground">Modèle</dt>
          <dd class="truncate font-medium text-foreground">
            {{ vehiculeInfo.modele }}
          </dd>
        </div>
        <div class="flex items-center justify-between gap-3">
          <dt class="text-muted-foreground">Immatriculation</dt>
          <dd class="font-mono font-medium text-foreground">
            {{ vehiculeInfo.immatriculation }}
          </dd>
        </div>
        <div class="flex items-center justify-between gap-3">
          <dt class="flex items-center gap-1.5 text-muted-foreground">
            <Wrench class="h-3.5 w-3.5" />
            Visite technique
          </dt>
          <dd class="font-medium text-foreground">
            {{ vehiculeInfo.visiteTechnique }}
          </dd>
        </div>
        <div class="flex items-center justify-between gap-3">
          <dt class="flex items-center gap-1.5 text-muted-foreground">
            <ShieldCheck class="h-3.5 w-3.5 text-brand" />
            Assurance
          </dt>
          <dd class="truncate font-medium text-brand-dark">
            {{ vehiculeInfo.assurance }}
          </dd>
        </div>
      </dl>
    </section>
  </aside>
</template>

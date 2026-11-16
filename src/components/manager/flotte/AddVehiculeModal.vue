<script setup lang="ts">
import { ref, watch } from "vue";
import type { StatutVehicule } from "../../../types";
import { Car, X } from "lucide-vue-next";

const props = defineProps<{
  isOpen: boolean;
  isSubmitting: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (
    e: "submit",
    payload: {
      immatriculation: string;
      marque_modele: string;
      type_boite: "MANUELLE" | "AUTOMATIQUE";
      kilometrage_actuel: number;
      statut: StatutVehicule;
      date_visite_technique: string;
      date_expiration_assurance: string;
    },
  ): void;
}>();

const immatriculation = ref("");
const marqueModele = ref("");
const typeBoite = ref<"MANUELLE" | "AUTOMATIQUE">("MANUELLE");
const kilometrage = ref(0);
const statut = ref<StatutVehicule>("DISPONIBLE");
const visiteExpire = ref("");
const assuranceExpire = ref("");

watch(
  () => props.isOpen,
  (open) => {
    if (!open) {
      immatriculation.value = "";
      marqueModele.value = "";
      kilometrage.value = 0;
      visiteExpire.value = "";
      assuranceExpire.value = "";
      statut.value = "DISPONIBLE";
    }
  },
);

const handleSubmit = () => {
  emit("submit", {
    immatriculation: immatriculation.value.trim().toUpperCase(),
    marque_modele: marqueModele.value.trim(),
    type_boite: typeBoite.value,
    kilometrage_actuel: Number(kilometrage.value) || 0,
    statut: statut.value,
    date_visite_technique: visiteExpire.value,
    date_expiration_assurance: assuranceExpire.value,
  });
};
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/50 p-0 sm:p-4"
    role="dialog"
    aria-modal="true"
    @click.self="emit('close')"
  >
    <div
      class="w-full sm:max-w-md max-h-[92dvh] flex flex-col overflow-hidden rounded-t-[3px] sm:rounded-[3px] border-t sm:border border-border bg-card pb-safe animate-in fade-in"
    >
      <div
        class="sm:hidden mx-auto mt-2 h-1 w-9 rounded-full bg-border shrink-0"
      />
      <header
        class="flex items-center justify-between border-b border-border px-5 py-3.5"
      >
        <div class="flex items-center space-x-2">
          <Car class="w-4 h-4 text-brand" />
          <h2 class="text-sm font-bold text-foreground">
            Enregistrer un véhicule de l'auto-école
          </h2>
        </div>
        <button
          type="button"
          @click="emit('close')"
          class="text-muted-foreground hover:text-foreground"
        >
          <X class="w-4 h-4" />
        </button>
      </header>

      <form @submit.prevent="handleSubmit" class="p-5 space-y-3.5 text-xs">
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="block font-medium text-foreground mb-1"
              >Immatriculation *</label
            >
            <input
              type="text"
              v-model="immatriculation"
              required
              placeholder="Ex: 1245-HJ-01"
              class="w-full h-8 px-2.5 rounded-lg border border-border bg-background text-foreground font-mono uppercase focus:border-brand focus:outline-none"
            />
          </div>
          <div>
            <label class="block font-medium text-foreground mb-1"
              >Marque & Modèle *</label
            >
            <input
              type="text"
              v-model="marqueModele"
              required
              placeholder="Peugeot 208 double pédale"
              class="w-full h-8 px-2.5 rounded-lg border border-border bg-background text-foreground focus:border-brand focus:outline-none"
            />
          </div>
        </div>

        <div class="grid grid-cols-3 gap-3">
          <div>
            <label class="block font-medium text-foreground mb-1"
              >Transmission</label
            >
            <select
              v-model="typeBoite"
              class="w-full h-8 px-2 rounded-lg border border-border bg-background text-foreground"
            >
              <option value="MANUELLE">Manuelle</option>
              <option value="AUTOMATIQUE">Automatique</option>
            </select>
          </div>
          <div>
            <label class="block font-medium text-foreground mb-1"
              >Kilométrage initial</label
            >
            <input
              type="number"
              v-model="kilometrage"
              min="0"
              class="w-full h-8 px-2.5 rounded-lg border border-border bg-background text-foreground tabular-nums focus:border-brand focus:outline-none"
            />
          </div>
          <div>
            <label class="block font-medium text-foreground mb-1"
              >Statut initial</label
            >
            <select
              v-model="statut"
              class="w-full h-8 px-2 rounded-lg border border-border bg-background text-foreground"
            >
              <option value="DISPONIBLE">Disponible</option>
              <option value="EN_ENTRETIEN">En entretien</option>
              <option value="EN_PANNE">En panne</option>
            </select>
          </div>
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="block font-medium text-foreground mb-1"
              >Expiration Visite Technique *</label
            >
            <input
              type="date"
              v-model="visiteExpire"
              required
              class="w-full h-8 px-2 rounded-lg border border-border bg-background text-foreground focus:border-brand focus:outline-none"
            />
          </div>
          <div>
            <label class="block font-medium text-foreground mb-1"
              >Expiration Assurance *</label
            >
            <input
              type="date"
              v-model="assuranceExpire"
              required
              class="w-full h-8 px-2 rounded-lg border border-border bg-background text-foreground focus:border-brand focus:outline-none"
            />
          </div>
        </div>

        <div
          class="pt-3 border-t border-border flex items-center justify-end space-x-2"
        >
          <button
            type="button"
            @click="emit('close')"
            :disabled="isSubmitting"
            class="px-3 py-1.5 rounded-lg border border-border text-muted-foreground hover:bg-muted font-medium transition"
          >
            Annuler
          </button>
          <button
            type="submit"
            :disabled="isSubmitting"
            class="px-4 py-1.5 rounded-lg bg-brand text-brand-foreground font-semibold shadow-xs hover:bg-brand/90 transition disabled:opacity-50"
          >
            {{ isSubmitting ? "Enregistrement..." : "Enregistrer en base" }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

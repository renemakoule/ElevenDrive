<script setup lang="ts">
import { ref } from "vue";
import type { SchoolItemSummary } from "../../../types";
import { Send, CheckCircle2 } from "lucide-vue-next";

defineProps<{
  isOpen: boolean;
  isSuccess: boolean;
  schools: SchoolItemSummary[];
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (
    e: "submit",
    payload: {
      cibleType: "all" | "single";
      cibleSchoolId: string;
      objet: string;
      message: string;
    },
  ): void;
}>();

const requestForm = ref({
  cibleType: "all" as "all" | "single",
  cibleSchoolId: "",
  objet: "Évaluation de satisfaction ElevenDrive 2026",
  message:
    "Chère Direction, votre avis compte énormément pour continuer à perfectionner la plateforme. Prenez 1 minute pour noter le système et nous faire part de vos suggestions !",
});
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-xs p-4 select-none animate-in fade-in"
    @click.self="emit('close')"
  >
    <div
      class="bg-card border border-border rounded-[3px] max-w-lg w-full p-5 space-y-4 text-xs"
    >
      <div
        class="flex items-center justify-between border-b border-border/60 pb-3"
      >
        <span
          class="font-bold text-sm text-foreground flex items-center gap-1.5"
        >
          <Send class="w-4 h-4 text-brand" />
          <span>Campagne de Collecte d'Avis Réseau</span>
        </span>
        <button
          type="button"
          @click="emit('close')"
          class="text-muted-foreground hover:text-foreground p-1"
        >
          ✕
        </button>
      </div>

      <div v-if="isSuccess" class="py-6 text-center space-y-2">
        <CheckCircle2 class="w-10 h-10 text-emerald-500 mx-auto" />
        <div class="font-bold text-sm text-foreground">
          Campagne transmise avec succès !
        </div>
        <p class="text-xs text-muted-foreground">
          L'invitation de notation en étoiles a été envoyée aux établissements
          ciblés.
        </p>
      </div>

      <form
        v-else
        @submit.prevent="emit('submit', { ...requestForm })"
        class="space-y-3"
      >
        <div>
          <label class="block font-semibold text-muted-foreground mb-1"
            >Destinataires</label
          >
          <div class="grid grid-cols-2 gap-2 mb-2">
            <button
              type="button"
              @click="requestForm.cibleType = 'all'"
              class="p-1.5 rounded border text-xs font-medium cursor-pointer"
              :class="
                requestForm.cibleType === 'all'
                  ? 'bg-brand text-white border-brand'
                  : 'bg-muted text-muted-foreground'
              "
            >
              Toutes les auto-écoles
            </button>
            <button
              type="button"
              @click="requestForm.cibleType = 'single'"
              class="p-1.5 rounded border text-xs font-medium cursor-pointer"
              :class="
                requestForm.cibleType === 'single'
                  ? 'bg-brand text-white border-brand'
                  : 'bg-muted text-muted-foreground'
              "
            >
              Auto-école spécifique
            </button>
          </div>
          <select
            v-if="requestForm.cibleType === 'single'"
            v-model="requestForm.cibleSchoolId"
            required
            class="w-full p-2 rounded border border-border bg-background"
          >
            <option value="" disabled>Sélectionner l'établissement...</option>
            <option v-for="s in schools" :key="s.id" :value="s.id">
              {{ s.nom_commercial || (s as any).nom }} ({{
                s.ville || "Actif"
              }})
            </option>
          </select>
        </div>

        <div>
          <label class="block font-semibold text-muted-foreground mb-1"
            >Objet de la notification</label
          >
          <input
            v-model="requestForm.objet"
            type="text"
            required
            class="w-full p-2 rounded border border-border bg-background"
          />
        </div>

        <div>
          <label class="block font-semibold text-muted-foreground mb-1"
            >Message d'invitation</label
          >
          <textarea
            v-model="requestForm.message"
            rows="3"
            required
            class="w-full p-2 rounded border border-border bg-background"
          ></textarea>
        </div>

        <div class="flex justify-end gap-2 pt-2">
          <button
            type="button"
            @click="emit('close')"
            class="px-3 py-1.5 rounded border border-border hover:bg-muted"
          >
            Annuler
          </button>
          <button
            type="submit"
            class="px-4 py-1.5 rounded bg-brand hover:bg-brand-dark text-white font-semibold flex items-center gap-1.5 cursor-pointer"
          >
            <Send class="w-3.5 h-3.5" />
            <span>Diffuser l'invitation</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

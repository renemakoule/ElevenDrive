<script setup lang="ts">
import { ref, watch } from "vue";
import type { SchoolItemSummary } from "../../../types";
import type { SchoolFeedbackItem } from "./types";
import { Star } from "lucide-vue-next";

const props = defineProps<{
  isOpen: boolean;
  isSubmitting: boolean;
  schools: SchoolItemSummary[];
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (
    e: "submit",
    payload: {
      autoEcoleId: string;
      auteur: string;
      role: string;
      etoiles: number;
      categorie: SchoolFeedbackItem["categorie"];
      commentaire: string;
    },
  ): void;
}>();

const form = ref({
  autoEcoleId: "",
  auteur: "Directeur Gérant",
  role: "Directeur",
  etoiles: 5,
  categorie: "ergonomie" as SchoolFeedbackItem["categorie"],
  commentaire: "",
});

watch(
  () => props.isOpen,
  (open) => {
    if (open) form.value.commentaire = "";
  },
);

const submit = () => {
  if (!form.value.commentaire.trim()) return;
  emit("submit", { ...form.value });
};
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-xs p-4 select-none animate-in fade-in"
    @click.self="emit('close')"
  >
    <div
      class="bg-card border border-border rounded-[3px] max-w-lg w-full p-5 space-y-4 shadow-xl text-xs"
    >
      <div
        class="flex items-center justify-between border-b border-border pb-3"
      >
        <span
          class="font-bold text-sm text-foreground flex items-center gap-1.5"
        >
          <Star class="w-4 h-4 text-amber-400 fill-amber-400" />
          <span>Simulateur d'Avis Auto-École</span>
        </span>
        <button
          type="button"
          @click="emit('close')"
          class="text-muted-foreground hover:text-foreground p-1"
        >
          ✕
        </button>
      </div>

      <form @submit.prevent="submit" class="space-y-3">
        <div>
          <label class="block font-semibold text-muted-foreground mb-1"
            >Auto-école émettrice</label
          >
          <select
            v-model="form.autoEcoleId"
            class="w-full p-2 rounded border border-border bg-background"
          >
            <option value="">Auto-École de test</option>
            <option v-for="s in schools" :key="s.id" :value="s.id">
              {{ s.nom_commercial || (s as any).nom }} ({{
                s.ville || "Actif"
              }})
            </option>
          </select>
        </div>

        <div class="grid grid-cols-2 gap-2">
          <div>
            <label class="block font-semibold text-muted-foreground mb-1"
              >Auteur</label
            >
            <input
              v-model="form.auteur"
              type="text"
              required
              class="w-full p-2 rounded border border-border bg-background"
            />
          </div>
          <div>
            <label class="block font-semibold text-muted-foreground mb-1"
              >Rôle</label
            >
            <select
              v-model="form.role"
              class="w-full p-2 rounded border border-border bg-background"
            >
              <option value="Directeur">Directeur</option>
              <option value="Secrétaire">Secrétaire</option>
              <option value="Moniteur">Moniteur</option>
            </select>
          </div>
        </div>

        <div>
          <label class="block font-semibold text-muted-foreground mb-1"
            >Note en étoiles (1 à 5)</label
          >
          <div
            class="flex items-center gap-2 p-2 bg-muted/30 rounded border border-border"
          >
            <button
              v-for="s in 5"
              :key="s"
              type="button"
              @click="form.etoiles = s"
              class="p-1 cursor-pointer transition hover:scale-110"
            >
              <Star
                class="w-6 h-6"
                :class="
                  s <= form.etoiles
                    ? 'fill-amber-400 text-amber-400'
                    : 'text-muted-foreground/30'
                "
              />
            </button>
            <span class="ml-2 font-bold text-sm text-foreground"
              >{{ form.etoiles }} / 5 étoiles</span
            >
          </div>
        </div>

        <div>
          <label class="block font-semibold text-muted-foreground mb-1"
            >Catégorie</label
          >
          <select
            v-model="form.categorie"
            class="w-full p-2 rounded border border-border bg-background"
          >
            <option value="ergonomie">Ergonomie & Design</option>
            <option value="guichet">Guichet & Reçus Thermiques</option>
            <option value="pwa_mobile">Application PWA Moniteurs/Élèves</option>
            <option value="examens">Examens & Listes d'Attente</option>
            <option value="support">Qualité du Support Client</option>
            <option value="general">Général</option>
          </select>
        </div>

        <div>
          <label class="block font-semibold text-muted-foreground mb-1"
            >Commentaire / Suggestion</label
          >
          <textarea
            v-model="form.commentaire"
            rows="3"
            required
            placeholder="Donnez votre avis sur l'application ElevenDrive..."
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
            :disabled="isSubmitting"
            class="px-4 py-1.5 rounded bg-brand hover:bg-brand-dark text-white font-semibold cursor-pointer disabled:opacity-50"
          >
            {{ isSubmitting ? "Enregistrement..." : "Enregistrer l'avis" }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

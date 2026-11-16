<script setup lang="ts">
import { ref, watch } from "vue";
import type { SchoolItemSummary } from "../../../types";
import { X, CheckSquare, Square, Save, Loader2, Layers } from "lucide-vue-next";

const props = defineProps<{
  isOpen: boolean;
  school: SchoolItemSummary | null;
  isLoading: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "save", schoolId: string, modules: string[]): void;
}>();

const availableModules = [
  {
    id: "guichet",
    nom: "Guichet & Enrôlement Candidats",
    description:
      "Inscriptions rapides, génération des dossiers et attestations",
    icon: "Receipt",
  },
  {
    id: "caisse",
    nom: "Caisse & Émargement Financier",
    description: "Reçus thermiques, clôture de caisse Z, validation manager",
    icon: "CreditCard",
  },
  {
    id: "conduite",
    nom: "Livret de Conduite & Moniteurs",
    description:
      "Suivi des compétences, émargement mobile et signatures numériques",
    icon: "Car",
  },
  {
    id: "examens",
    nom: "Sessions & Examens d'État",
    description:
      "Éligibilité ministérielle, bordereaux officiels et publication des résultats",
    icon: "GraduationCap",
  },
  {
    id: "flotte",
    nom: "Flotte Automobile & Maintenance",
    description:
      "Alertes vidange, assurances, visites techniques et suivi kilométrique",
    icon: "Truck",
  },
  {
    id: "archivage_cloud",
    nom: "Archivage Numérique Cloud R2",
    description:
      "Stockage sécurisé Cloudflare R2 des CNI, photos et pièces d'examen",
    icon: "Cloud",
  },
];

const selectedModules = ref<string[]>([]);

watch(
  () => props.school,
  (s) => {
    if (s && s.modules_actifs) {
      try {
        const parsed = JSON.parse(s.modules_actifs);
        if (Array.isArray(parsed)) {
          selectedModules.value = parsed;
          return;
        }
      } catch {}
    }
    selectedModules.value = availableModules.map((m) => m.id);
  },
  { immediate: true },
);

const toggleModule = (id: string) => {
  if (selectedModules.value.includes(id)) {
    selectedModules.value = selectedModules.value.filter((m) => m !== id);
  } else {
    selectedModules.value.push(id);
  }
};

const selectAll = () => {
  selectedModules.value = availableModules.map((m) => m.id);
};

const deselectAll = () => {
  selectedModules.value = [];
};

const handleSubmit = () => {
  if (!props.school) return;
  emit("save", props.school.id, selectedModules.value);
};
</script>

<template>
  <div
    v-if="isOpen && school"
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm"
  >
    <div
      class="w-full max-w-lg bg-card border border-border rounded-[8px] shadow-2xl overflow-hidden animate-in fade-in zoom-in-95 duration-150"
    >
      <!-- En-tête -->
      <div
        class="h-12 border-b border-border bg-muted/50 px-4 flex items-center justify-between"
      >
        <div
          class="flex items-center gap-2 text-foreground font-semibold text-sm"
        >
          <Layers class="w-4 h-4 text-brand" />
          <span
            >Droits d'Accès & Modules Système ·
            {{ school.nom_commercial }}</span
          >
        </div>
        <button
          type="button"
          @click="emit('close')"
          class="text-muted-foreground hover:text-foreground p-1 rounded-[4px] transition-colors cursor-pointer"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Corps -->
      <form
        @submit.prevent="handleSubmit"
        class="p-5 space-y-4 text-xs font-sans"
      >
        <div class="flex items-center justify-between text-muted-foreground">
          <span
            >Activez ou restreignez les fonctionnalités selon l'offre souscrite
            :</span
          >
          <div class="flex items-center gap-2">
            <button
              type="button"
              @click="selectAll"
              class="text-brand hover:underline font-medium cursor-pointer"
            >
              Tout cocher
            </button>
            <span>·</span>
            <button
              type="button"
              @click="deselectAll"
              class="text-muted-foreground hover:underline font-medium cursor-pointer"
            >
              Tout décocher
            </button>
          </div>
        </div>

        <div class="space-y-2 max-h-72 overflow-y-auto pr-1">
          <div
            v-for="mod in availableModules"
            :key="mod.id"
            @click="toggleModule(mod.id)"
            :class="[
              'p-3 rounded-[6px] border flex items-start gap-3 transition cursor-pointer select-none',
              selectedModules.includes(mod.id)
                ? 'bg-brand/10 border-brand/40 shadow-xs'
                : 'bg-card border-border hover:bg-muted/50 opacity-80 hover:opacity-100',
            ]"
          >
            <div class="mt-0.5">
              <CheckSquare
                v-if="selectedModules.includes(mod.id)"
                class="w-4 h-4 text-brand"
              />
              <Square v-else class="w-4 h-4 text-muted-foreground" />
            </div>

            <div class="flex-1">
              <div
                class="font-semibold text-foreground flex items-center justify-between"
              >
                <span>{{ mod.nom }}</span>
                <span
                  class="text-[10px] font-mono px-1.5 py-0.5 rounded-[3px] font-semibold"
                  :class="
                    selectedModules.includes(mod.id)
                      ? 'bg-brand/20 text-brand'
                      : 'bg-muted text-muted-foreground'
                  "
                >
                  {{
                    selectedModules.includes(mod.id) ? "ACTIVÉ" : "VERROUILLÉ"
                  }}
                </span>
              </div>
              <p
                class="text-[11px] text-muted-foreground mt-0.5 leading-relaxed"
              >
                {{ mod.description }}
              </p>
            </div>
          </div>
        </div>

        <!-- Actions -->
        <div
          class="pt-3 border-t border-border flex items-center justify-end gap-2.5"
        >
          <button
            type="button"
            @click="emit('close')"
            class="px-3 py-1.5 rounded-[4px] bg-muted hover:bg-muted/80 text-foreground border border-border font-medium transition cursor-pointer"
          >
            Annuler
          </button>
          <button
            type="submit"
            :disabled="isLoading"
            class="px-4 py-1.5 rounded-[4px] bg-brand-600 hover:bg-brand-700 text-white font-medium flex items-center gap-1.5 transition cursor-pointer disabled:opacity-50 shadow-xs"
          >
            <Loader2
              v-if="isLoading"
              class="w-3.5 h-3.5 animate-spin text-white"
            />
            <Save v-else class="w-3.5 h-3.5" />
            <span>Appliquer les Droits</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

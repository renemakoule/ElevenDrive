<script setup lang="ts">
import { ref } from 'vue';
import { usePlanningConduiteConfigStore } from '../../stores/planningConduiteConfigStore';
import {
  X,
  Clock,
  Calendar,
  Plus,
  Trash2,
  RotateCcw,
  CheckCircle2,
  Sliders,
  Sparkles,
  Info
} from 'lucide-vue-next';

defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
}>();

const configStore = usePlanningConduiteConfigStore();

// Formulaire d'ajout d'un nouveau créneau
const newLabel = ref('');
const newDebut = ref('08:00');
const newFin = ref('10:00');
const addError = ref<string | null>(null);

const ajouterNouveauCreneau = () => {
  addError.value = null;
  if (!newDebut.value || !newFin.value) {
    addError.value = 'Veuillez renseigner les heures de début et de fin.';
    return;
  }
  if (newDebut.value >= newFin.value) {
    addError.value = "L'heure de début doit être antérieure à l'heure de fin.";
    return;
  }

  configStore.ajouterCreneau({
    label: newLabel.value.trim() || `Créneau ${newDebut.value}`,
    debut: newDebut.value,
    fin: newFin.value,
  });

  newLabel.value = '';
};

const handleReset = () => {
  if (window.confirm('Voulez-vous réinitialiser tous les créneaux et jours aux réglages standards ?')) {
    configStore.resetToDefaults();
  }
};
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-50 flex items-center justify-center p-3 sm:p-5 bg-black/80 backdrop-blur-sm animate-in fade-in"
    role="dialog"
    aria-modal="true"
  >
    <div
      class="bg-card border border-border w-full max-w-2xl rounded-[8px] flex flex-col overflow-hidden text-foreground animate-in zoom-in-95 duration-150"
    >
      <!-- En-tête -->
      <header class="px-5 py-3.5 border-b border-border flex items-center justify-between bg-muted/40">
        <div class="flex items-center space-x-2.5">
          <div class="w-8 h-8 rounded-[5px] bg-brand/10 border border-brand/20 text-primary flex items-center justify-center">
            <Sliders class="w-4 h-4" />
          </div>
          <div>
            <h2 class="text-sm font-bold text-foreground">
              Configuration des Créneaux & Jours de Passage
            </h2>
            <p class="text-[11px] text-muted-foreground">
              Définissez votre propre grille horaire de conduite et activez les jours de formation de votre établissement.
            </p>
          </div>
        </div>

        <button
          type="button"
          @click="emit('close')"
          class="p-1 rounded-[4px] text-muted-foreground hover:text-foreground transition cursor-pointer"
        >
          <X class="w-5 h-5" />
        </button>
      </header>

      <!-- Contenu -->
      <div class="p-5 space-y-5 overflow-y-auto max-h-[75vh]">
        <!-- SECTION 1 : Jours de passage autorisés -->
        <div class="space-y-2">
          <div class="flex items-center justify-between">
            <label class="text-xs font-bold text-foreground flex items-center gap-1.5">
              <Calendar class="w-3.5 h-3.5 text-brand" />
              <span>1. Jours de passage des leçons pratiques (Semaine)</span>
            </label>
            <span class="text-[10px] text-muted-foreground font-mono">
              {{ configStore.joursActifs.length }} / 7 jours cochés
            </span>
          </div>

          <p class="text-[11px] text-muted-foreground">
            Cochez les jours où les moniteurs effectuent les séances sur piste ou circulation :
          </p>

          <div class="grid grid-cols-2 sm:grid-cols-4 gap-2 pt-1">
            <button
              v-for="jour in configStore.joursPassage"
              :key="jour.id"
              type="button"
              @click="configStore.toggleJour(jour.id)"
              class="p-2.5 rounded-[6px] border text-xs font-semibold flex items-center justify-between transition cursor-pointer"
              :class="
                jour.actif
                  ? 'bg-brand/15 border-brand text-primary'
                  : 'bg-muted/30 border-border text-muted-foreground hover:bg-muted'
              "
            >
              <div class="flex items-center gap-2">
                <span class="w-4 h-4 rounded border flex items-center justify-center text-[10px]"
                  :class="jour.actif ? 'bg-brand text-primary border-brand font-bold' : 'border-border bg-background'"
                >
                  <span v-if="jour.actif">✓</span>
                </span>
                <span>{{ jour.label }}</span>
              </div>
              <span class="text-[9px] font-mono uppercase" :class="jour.actif ? 'text-primary font-bold' : 'text-muted-foreground'">
                {{ jour.code }}
              </span>
            </button>
          </div>
        </div>

        <!-- SECTION 2 : Grille des Créneaux Horaires Personnalisés -->
        <div class="space-y-3 pt-2 border-t border-border">
          <div class="flex items-center justify-between">
            <label class="text-xs font-bold text-foreground flex items-center gap-1.5">
              <Clock class="w-3.5 h-3.5 text-brand" />
              <span>2. Créneaux horaires personnalisés (Grille du planning)</span>
            </label>
            <button
              type="button"
              @click="handleReset"
              class="text-[11px] text-muted-foreground hover:text-foreground flex items-center gap-1 transition cursor-pointer"
              title="Rétablir les 4 créneaux par défaut"
            >
              <RotateCcw class="w-3 h-3" />
              <span>Réinitialiser par défaut</span>
            </button>
          </div>

          <p class="text-[11px] text-muted-foreground">
            Ajustez les heures de début et fin de chaque plage ou ajoutez vos propres créneaux selon votre organisation :
          </p>

          <!-- Liste des créneaux existants -->
          <div class="space-y-2">
            <div
              v-for="cr in configStore.creneaux"
              :key="cr.id"
              class="p-2.5 rounded-[6px] bg-muted/40 border border-border flex items-center justify-between gap-3"
            >
              <div class="flex items-center gap-2.5 flex-1 min-w-0">
                <input
                  type="text"
                  v-model="cr.label"
                  @change="configStore.save()"
                  placeholder="Intitulé créneau..."
                  class="px-2 py-1 text-xs rounded bg-background border border-border focus:border-brand outline-none flex-1 font-semibold text-foreground min-w-[120px]"
                />

                <div class="flex items-center gap-1.5 font-mono text-xs">
                  <input
                    type="time"
                    v-model="cr.debut"
                    @change="configStore.save()"
                    class="px-1.5 py-1 text-xs rounded bg-background border border-border focus:border-brand outline-none text-foreground font-bold"
                  />
                  <span class="text-muted-foreground">à</span>
                  <input
                    type="time"
                    v-model="cr.fin"
                    @change="configStore.save()"
                    class="px-1.5 py-1 text-xs rounded bg-background border border-border focus:border-brand outline-none text-foreground font-bold"
                  />
                </div>
              </div>

              <!-- Bouton supprimer -->
              <button
                type="button"
                @click="configStore.supprimerCreneau(cr.id)"
                class="p-1.5 rounded text-muted-foreground hover:text-red-500 hover:bg-red-500/10 transition cursor-pointer"
                title="Supprimer ce créneau"
              >
                <Trash2 class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>

          <!-- Formulaire d'ajout rapide d'un créneau -->
          <div class="p-3 rounded-[6px] bg-background border border-dashed border-border space-y-2">
            <span class="text-[11px] font-bold text-foreground flex items-center gap-1">
              <Plus class="w-3.5 h-3.5 text-brand" /> Ajouter un nouveau créneau horaire
            </span>

            <div v-if="addError" class="text-[11px] text-red-500 font-medium">
              {{ addError }}
            </div>

            <div class="flex flex-wrap items-center gap-2">
              <input
                type="text"
                v-model="newLabel"
                placeholder="Libellé (ex: Début de matinée, Midi...)"
                class="px-2.5 py-1.5 text-xs rounded bg-muted/40 border border-border focus:border-brand outline-none flex-1 min-w-[160px]"
              />

              <div class="flex items-center gap-1.5 font-mono text-xs">
                <input
                  type="time"
                  v-model="newDebut"
                  class="px-2 py-1.5 text-xs rounded bg-muted/40 border border-border focus:border-brand outline-none font-bold"
                />
                <span class="text-muted-foreground">à</span>
                <input
                  type="time"
                  v-model="newFin"
                  class="px-2 py-1.5 text-xs rounded bg-muted/40 border border-border focus:border-brand outline-none font-bold"
                />
              </div>

              <button
                type="button"
                @click="ajouterNouveauCreneau"
                class="px-3 py-1.5 rounded bg-brand hover:bg-brand-dark text-primary font-bold text-xs uppercase tracking-wider transition cursor-pointer flex items-center gap-1 shadow-xs"
              >
                <Plus class="w-3.5 h-3.5" />
                <span>Ajouter</span>
              </button>
            </div>
          </div>
        </div>

        <!-- Note d'information -->
        <div class="p-3 rounded-[6px] bg-sky-500/10 border border-sky-500/20 text-sky-700 dark:text-sky-300 text-xs flex items-start gap-2">
          <Info class="w-4 h-4 shrink-0 mt-0.5" />
          <div class="text-[11px] leading-relaxed">
            Ces créneaux personnalisés s'appliqueront instantanément sur la grille du Directeur et de la Secrétaire. Les heures de conduite libres restent également programmables manuellement.
          </div>
        </div>
      </div>

      <!-- Pied de page -->
      <footer class="px-5 py-3 border-t border-border flex items-center justify-end bg-muted/20">
        <button
          type="button"
          @click="emit('close')"
          class="px-5 py-2 rounded-[5px] bg-brand hover:bg-brand-dark text-primary font-bold text-xs uppercase tracking-wider transition cursor-pointer shadow-xs"
        >
          Valider & Appliquer
        </button>
      </footer>
    </div>
  </div>
</template>

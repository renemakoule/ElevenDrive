<script setup lang="ts">
import { ref, computed, watch } from "vue";
import { useGuichetStore } from "../../../stores/guichetStore";
import type { TranchePaiementConfig } from "../../../types";
import { formatFcfa, type TranchesMode } from "./tranchesUtils";
import TranchesPaymentEditor from "./TranchesPaymentEditor.vue";
import { PlusCircle, X } from "lucide-vue-next";

const props = defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
}>();

const guichetStore = useGuichetStore();

const newCategorie = ref("B");
const newIntitule = ref("");
const newDescription = ref("");
const newHeuresTheoriques = ref(20);
const newHeuresPratiques = ref(20);
const newPrixTotal = ref<number | null>(null);
const newAcompteMinimum = ref<number | null>(null);
const newTranchesMode = ref<TranchesMode>("3");
const newTranches = ref<TranchePaiementConfig[]>([]);

watch(
  () => props.isOpen,
  (open) => {
    if (open) {
      newCategorie.value = "B";
      newIntitule.value = "";
      newDescription.value = "";
      newHeuresTheoriques.value = 20;
      newHeuresPratiques.value = 20;
      newPrixTotal.value = null;
      newAcompteMinimum.value = null;
      newTranchesMode.value = "3";
      newTranches.value = [];
    }
  },
);

const newTranchesTotal = computed(() =>
  newTranches.value.reduce((acc, t) => acc + (Number(t.montant) || 0), 0),
);

const newTranchesEcart = computed(() => {
  const prix = Number(newPrixTotal.value) || 0;
  return newTranchesTotal.value - prix;
});

const handleCreateFormule = async () => {
  if (!newIntitule.value.trim()) {
    guichetStore.showToast(
      "L'intitulé de la formule est obligatoire.",
      "error",
    );
    return;
  }
  const prix = Number(newPrixTotal.value);
  if (!prix || prix <= 0) {
    guichetStore.showToast(
      "Veuillez saisir un prix total valide supérieur à 0.",
      "error",
    );
    return;
  }
  if (newTranches.value.length > 0 && newTranchesEcart.value !== 0) {
    guichetStore.showToast(
      `La somme des tranches (${formatFcfa(newTranchesTotal.value)}) doit correspondre exactement au prix total (${formatFcfa(prix)}).`,
      "error",
    );
    return;
  }

  const acompte =
    Number(newAcompteMinimum.value) || newTranches.value[0]?.montant || 0;
  if (acompte > prix) {
    guichetStore.showToast(
      "L'acompte minimum ne peut excéder le prix total.",
      "error",
    );
    return;
  }

  try {
    await guichetStore.creerFormule({
      categorie_permis: newCategorie.value,
      intitule: newIntitule.value.trim(),
      description: newDescription.value.trim() || undefined,
      heures_theoriques: Number(newHeuresTheoriques.value) || 0,
      heures_pratiques: Number(newHeuresPratiques.value) || 0,
      prix_total: prix,
      acompte_minimum: acompte,
      tranches_paiement:
        newTranches.value.length > 0 ? newTranches.value : undefined,
    });
    emit("close");
  } catch {
    // Le toast d'erreur est géré par le store
  }
};
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 backdrop-blur-sm overflow-y-auto"
    role="dialog"
    aria-modal="true"
  >
    <div
      class="w-full max-w-xl max-h-[92vh] overflow-y-auto rounded-[4px] border border-border bg-card shadow-xl animate-in fade-in zoom-in-95 my-auto"
    >
      <header
        class="sticky top-0 z-10 flex items-center justify-between border-b border-border px-5 py-3.5 bg-card/95 backdrop-blur-sm"
      >
        <div class="flex items-center space-x-2">
          <PlusCircle class="w-4 h-4 text-brand" />
          <div>
            <h2 class="text-sm font-bold text-foreground">
              Nouvelle Formule & Modalités Tarifaires
            </h2>
            <p class="text-[10px] text-muted-foreground">
              Fixez le tarif officiel et définissez l'échéancier des tranches de
              paiement
            </p>
          </div>
        </div>
        <button
          type="button"
          @click="emit('close')"
          class="text-muted-foreground hover:text-foreground"
        >
          <X class="w-4 h-4" />
        </button>
      </header>

      <form @submit.prevent="handleCreateFormule" class="p-5 space-y-4 text-xs">
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="block font-medium text-foreground mb-1"
              >Catégorie de Permis *</label
            >
            <select
              v-model="newCategorie"
              class="w-full h-8 px-2 rounded-[4px] border border-border bg-background text-foreground"
            >
              <option value="A">Permis A (Motos & Gros Cubes)</option>
              <option value="A1">Permis A1 (Léger 2 roues)</option>
              <option value="B">Permis B (Véhicule Léger / Tourisme)</option>
              <option value="C">
                Permis C (Poids Lourd / Transport Marchandises)
              </option>
              <option value="D">
                Permis D (Transport en Commun / Autocar)
              </option>
              <option value="E">Permis E (Remorques & Semi-remorques)</option>
              <option value="REC">Recyclage / Performance (Perfectionnement)</option>
            </select>
          </div>
          <div>
            <label class="block font-medium text-foreground mb-1"
              >Intitulé Commercial *</label
            >
            <input
              type="text"
              v-model="newIntitule"
              required
              placeholder="Ex: Permis B VIP Accéléré"
              class="w-full h-8 px-2 rounded-[4px] border border-border bg-background text-foreground"
            />
          </div>
        </div>

        <div>
          <label class="block font-medium text-foreground mb-1"
            >Description pédagogique</label
          >
          <input
            type="text"
            v-model="newDescription"
            placeholder="Ex: Formule intensive, cours de code illimités et simulateur inclus"
            class="w-full h-8 px-2 rounded-[4px] border border-border bg-background text-foreground"
          />
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="block font-medium text-foreground mb-1"
              >Heures Théorie (Code)</label
            >
            <input
              type="number"
              v-model="newHeuresTheoriques"
              min="0"
              class="w-full h-8 px-2 rounded-[4px] border border-border bg-background text-foreground tabular-nums"
            />
          </div>
          <div>
            <label class="block font-medium text-foreground mb-1"
              >Heures Pratique (Conduite)</label
            >
            <input
              type="number"
              v-model="newHeuresPratiques"
              min="0"
              class="w-full h-8 px-2 rounded-[4px] border border-border bg-background text-foreground tabular-nums"
            />
          </div>
        </div>

        <div class="grid grid-cols-2 gap-3 pt-2 border-t border-border">
          <div>
            <label class="block font-bold text-foreground mb-1">
              Prix Total Fixé (FCFA) *
            </label>
            <input
              type="number"
              v-model="newPrixTotal"
              min="0"
              step="any"
              required
              placeholder="Saisir le tarif (ex: 150000)"
              class="w-full h-8 px-2.5 rounded-[4px] border border-brand bg-background text-foreground font-bold tabular-nums focus:ring-1 focus:ring-brand"
            />
            <span class="text-[10px] text-muted-foreground mt-0.5 block">
              Montant global net de la formation
            </span>
          </div>
          <div>
            <label class="block font-medium text-foreground mb-1">
              Acompte Minimum Requis (FCFA) *
            </label>
            <input
              type="number"
              v-model="newAcompteMinimum"
              min="0"
              step="any"
              required
              placeholder="Montant 1ère tranche"
              class="w-full h-8 px-2.5 rounded-[4px] border border-border bg-background text-foreground tabular-nums font-semibold"
            />
            <span class="text-[10px] text-muted-foreground mt-0.5 block">
              Synchronisé avec la 1ère tranche
            </span>
          </div>
        </div>

        <!-- Éditeur de tranches réutilisable -->
        <TranchesPaymentEditor
          v-model="newTranches"
          v-model:mode="newTranchesMode"
          :prix-total="newPrixTotal"
          :show-empty-prompt="true"
          @sync-acompte="(val) => (newAcompteMinimum = val)"
        />

        <div
          class="pt-3 border-t border-border flex items-center justify-end space-x-2"
        >
          <button
            type="button"
            @click="emit('close')"
            class="px-3 py-1.5 rounded-[4px] border border-border text-muted-foreground hover:bg-muted font-medium transition"
          >
            Annuler
          </button>
          <button
            type="submit"
            :disabled="
              guichetStore.loading ||
              (newTranches.length > 0 && newTranchesEcart !== 0)
            "
            class="px-4 py-1.5 rounded-[4px] bg-brand text-brand-foreground font-semibold shadow-xs hover:bg-brand/90 transition disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{
              guichetStore.loading
                ? "Enregistrement..."
                : "Fixer le Tarif & Enregistrer"
            }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

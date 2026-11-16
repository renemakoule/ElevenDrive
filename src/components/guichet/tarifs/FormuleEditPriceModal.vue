<script setup lang="ts">
import { ref, computed, watch } from "vue";
import { useGuichetStore } from "../../../stores/guichetStore";
import type { FormulePermis, TranchePaiementConfig } from "../../../types";
import {
  formatFcfa,
  parseTranches,
  genererTranchesPredefinies,
  type TranchesMode,
} from "./tranchesUtils";
import TranchesPaymentEditor from "./TranchesPaymentEditor.vue";
import { Edit3, X } from "lucide-vue-next";

const props = defineProps<{
  isOpen: boolean;
  formule: FormulePermis | null;
}>();

const emit = defineEmits<{
  (e: "close"): void;
}>();

const guichetStore = useGuichetStore();

const nouveauPrix = ref<number | null>(null);
const nouvelAcompte = ref<number | null>(null);
const motifModification = ref("");
const editTranches = ref<TranchePaiementConfig[]>([]);
const editTranchesMode = ref<TranchesMode>("custom");

watch(
  () => [props.isOpen, props.formule] as const,
  ([open, f]) => {
    if (open && f) {
      nouveauPrix.value = f.prix_total;
      nouvelAcompte.value = f.acompte_minimum;
      motifModification.value = "";
      const existantes = parseTranches(f.tranches_paiement);
      if (existantes && existantes.length > 0) {
        editTranches.value = JSON.parse(JSON.stringify(existantes));
        editTranchesMode.value = "custom";
      } else {
        editTranches.value = genererTranchesPredefinies(f.prix_total, 3);
        editTranchesMode.value = "3";
      }
    }
  },
  { immediate: true },
);

const editTranchesTotal = computed(() =>
  editTranches.value.reduce((acc, t) => acc + (Number(t.montant) || 0), 0),
);

const editTranchesEcart = computed(() => {
  const prix = Number(nouveauPrix.value) || 0;
  return editTranchesTotal.value - prix;
});

const handleUpdatePrice = async () => {
  if (!props.formule) return;
  const prix = Number(nouveauPrix.value);
  if (!prix || prix <= 0) {
    guichetStore.showToast(
      "Veuillez saisir un nouveau prix valide supérieur à 0.",
      "error",
    );
    return;
  }

  if (editTranches.value.length > 0 && editTranchesEcart.value !== 0) {
    guichetStore.showToast(
      `La somme des tranches (${formatFcfa(editTranchesTotal.value)}) doit correspondre au prix total (${formatFcfa(prix)}).`,
      "error",
    );
    return;
  }

  const acompte =
    Number(nouvelAcompte.value) || editTranches.value[0]?.montant || 0;
  if (acompte > prix) {
    guichetStore.showToast(
      "L'acompte minimum ne peut excéder le prix total.",
      "error",
    );
    return;
  }

  try {
    await guichetStore.modifierPrixFormule(
      props.formule.id,
      prix,
      acompte,
      motifModification.value.trim() || undefined,
      editTranches.value.length > 0 ? editTranches.value : undefined,
    );
    emit("close");
  } catch {
    // Le toast d'erreur est géré par le store
  }
};
</script>

<template>
  <div
    v-if="isOpen && formule"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 backdrop-blur-sm overflow-y-auto"
    role="dialog"
    aria-modal="true"
  >
    <div
      class="w-full max-w-xl max-h-[92vh] overflow-y-auto rounded-[4px] border border-border bg-card shadow-xl animate-in fade-in zoom-in-95 my-auto"
    >
      <header
        class="sticky top-0 z-10 flex items-center justify-between border-b border-border px-5 py-3.5 bg-brand-light/40 backdrop-blur-sm"
      >
        <div class="flex items-center space-x-2">
          <Edit3 class="w-4 h-4 text-brand" />
          <div>
            <h2 class="text-sm font-bold text-foreground">
              Ajustement Tarifaire & Échéancier (Direction)
            </h2>
            <p class="text-[10px] text-muted-foreground">
              {{ formule.intitule }} (Permis {{ formule.categorie_permis }})
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

      <form @submit.prevent="handleUpdatePrice" class="p-5 space-y-4 text-xs">
        <div class="rounded-[4px] bg-muted p-3 border border-border">
          <div class="flex justify-between items-center text-muted-foreground">
            <span>Prix actuellement en vigueur :</span>
            <span class="font-semibold text-foreground line-through tabular-nums">
              {{ formatFcfa(formule.prix_total) }}
            </span>
          </div>
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="block font-bold text-foreground mb-1"
              >Nouveau Prix Officiel (FCFA) *</label
            >
            <input
              type="number"
              v-model="nouveauPrix"
              min="0"
              step="any"
              required
              class="w-full h-9 px-3 rounded-[4px] border border-brand bg-background text-foreground font-bold text-sm tabular-nums focus:ring-2 focus:ring-brand/40"
            />
          </div>

          <div>
            <label class="block font-medium text-foreground mb-1"
              >Nouvel Acompte Minimum (FCFA)</label
            >
            <input
              type="number"
              v-model="nouvelAcompte"
              min="0"
              step="any"
              class="w-full h-9 px-3 rounded-[4px] border border-border bg-background text-foreground tabular-nums font-semibold"
            />
          </div>
        </div>

        <!-- Éditeur de tranches réutilisable -->
        <TranchesPaymentEditor
          v-model="editTranches"
          v-model:mode="editTranchesMode"
          :prix-total="nouveauPrix"
          title="Réajustement des Tranches de Paiement"
          subtitle="Adaptez l'échelonnement au nouveau tarif officiel"
          :compact-buttons="true"
          @sync-acompte="(val) => (nouvelAcompte = val)"
        />

        <div>
          <label class="block font-medium text-foreground mb-1"
            >Motif de la modification (Audit & Décision de Direction)</label
          >
          <input
            type="text"
            v-model="motifModification"
            placeholder="Ex: Décision tarifaire rentrée 2026, nouvel aménagement d'échéances"
            class="w-full h-8 px-2.5 rounded-[4px] border border-border bg-background text-foreground"
          />
        </div>

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
              (editTranches.length > 0 && editTranchesEcart !== 0)
            "
            class="px-4 py-1.5 rounded-[4px] bg-brand text-brand-foreground font-semibold shadow-xs hover:bg-brand/90 transition disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{
              guichetStore.loading
                ? "Mise à jour..."
                : "Confirmer Nouveau Tarif & Tranches"
            }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>


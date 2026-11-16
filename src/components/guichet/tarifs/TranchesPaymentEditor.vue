<script setup lang="ts">
import { computed, watch } from "vue";
import type { TranchePaiementConfig } from "../../../types";
import {
  formatFcfa,
  genererTranchesPredefinies,
  type TranchesMode,
} from "./tranchesUtils";
import { Coins, Trash2, Plus, RefreshCw, AlertTriangle } from "lucide-vue-next";

const props = withDefaults(
  defineProps<{
    modelValue: TranchePaiementConfig[];
    mode: TranchesMode;
    prixTotal: number | null;
    title?: string;
    subtitle?: string;
    compactButtons?: boolean;
    showEmptyPrompt?: boolean;
  }>(),
  {
    title: "Définition des Tranches de Paiement",
    subtitle: "Organisez l'échelonnement des versements pour les élèves",
    compactButtons: false,
    showEmptyPrompt: false,
  },
);

const emit = defineEmits<{
  (e: "update:modelValue", val: TranchePaiementConfig[]): void;
  (e: "update:mode", val: TranchesMode): void;
  (e: "sync-acompte", montant: number): void;
}>();

const tranchesTotal = computed(() =>
  props.modelValue.reduce((acc, t) => acc + (Number(t.montant) || 0), 0),
);

const tranchesEcart = computed(() => {
  const prix = Number(props.prixTotal) || 0;
  return tranchesTotal.value - prix;
});

const applyMode = (mode: TranchesMode) => {
  emit("update:mode", mode);
  const prix = Number(props.prixTotal) || 0;
  if (prix <= 0) return;

  if (mode !== "custom") {
    const nb = parseInt(mode, 10);
    const next = genererTranchesPredefinies(prix, nb);
    emit("update:modelValue", next);
    if (next.length > 0) {
      emit("sync-acompte", next[0].montant);
    }
  } else if (props.modelValue.length === 0) {
    const next: TranchePaiementConfig[] = [
      {
        numero: 1,
        intitule: "1ère Tranche (Acompte)",
        montant: Math.round(prix * 0.5),
        echeance: "À l'inscription",
      },
      {
        numero: 2,
        intitule: "2ème Tranche (Solde)",
        montant: Math.round(prix * 0.5),
        echeance: "Avant examen",
      },
    ];
    emit("update:modelValue", next);
    emit("sync-acompte", next[0].montant);
  }
};

watch(
  () => props.prixTotal,
  (nouveauPrix) => {
    const prix = Number(nouveauPrix) || 0;
    if (prix > 0 && props.mode !== "custom") {
      const nb = parseInt(props.mode, 10) || 3;
      const next = genererTranchesPredefinies(prix, nb);
      emit("update:modelValue", next);
      if (next.length > 0) {
        emit("sync-acompte", next[0].montant);
      }
    }
  },
);

const ajouterTranche = () => {
  emit("update:mode", "custom");
  const num = props.modelValue.length + 1;
  const soldeRestant = Math.max(0, -tranchesEcart.value);
  const next = [
    ...props.modelValue,
    {
      numero: num,
      intitule: `${num}e Tranche`,
      montant: soldeRestant > 0 ? soldeRestant : 0,
      echeance: `Échéance ${num}`,
    },
  ];
  emit("update:modelValue", next);
};

const supprimerTranche = (index: number) => {
  emit("update:mode", "custom");
  const next = props.modelValue
    .filter((_, idx) => idx !== index)
    .map((t, i) => ({ ...t, numero: i + 1 }));
  emit("update:modelValue", next);
  if (next.length > 0) {
    emit("sync-acompte", Number(next[0].montant) || 0);
  }
};

const handleMontantInput = (idx: number, rawVal: any) => {
  const numVal = Number(rawVal) || 0;
  const next = props.modelValue.map((t, i) =>
    i === idx ? { ...t, montant: numVal } : t,
  );
  emit("update:modelValue", next);
  if (idx === 0) {
    emit("sync-acompte", numVal);
  }
};

const handleFieldUpdate = (
  idx: number,
  field: "intitule" | "echeance",
  val: string,
) => {
  const next = props.modelValue.map((t, i) =>
    i === idx ? { ...t, [field]: val } : t,
  );
  emit("update:modelValue", next);
};

const equilibrerTranches = () => {
  const prix = Number(props.prixTotal) || 0;
  if (prix <= 0 || props.modelValue.length === 0) return;
  const n = props.modelValue.length;
  const part = Math.round(prix / n / 1000) * 1000;
  let somme = 0;
  const next = props.modelValue.map((t, i) => {
    const montant = i === n - 1 ? Math.max(0, prix - somme) : part;
    if (i < n - 1) somme += part;
    return { ...t, montant };
  });
  emit("update:modelValue", next);
  emit("sync-acompte", next[0].montant);
};

const modeOptions: Array<{
  id: TranchesMode;
  label: string;
  compact: string;
}> = [
  { id: "1", label: "1x (Comptant)", compact: "1x" },
  { id: "2", label: "2 Tranches", compact: "2 Tr." },
  { id: "3", label: "3 Tranches", compact: "3 Tr." },
  { id: "4", label: "4 Tranches", compact: "4 Tr." },
  { id: "custom", label: "Libre", compact: "Libre" },
];
</script>

<template>
  <div class="pt-3 border-t border-border space-y-3">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div>
        <h4 class="font-bold text-foreground flex items-center gap-1.5 text-xs">
          <Coins class="w-3.5 h-3.5 text-brand" />
          <span>{{ title }}</span>
        </h4>
        <p class="text-[10px] text-muted-foreground">
          {{ subtitle }}
        </p>
      </div>

      <div class="flex items-center gap-1 bg-muted p-0.5 rounded-[4px]">
        <button
          v-for="opt in modeOptions"
          :key="opt.id"
          type="button"
          @click="applyMode(opt.id)"
          :class="[
            'px-2 py-0.5 rounded text-[10px] font-semibold transition',
            mode === opt.id
              ? 'bg-brand text-white shadow-xs'
              : 'text-muted-foreground hover:text-foreground',
          ]"
        >
          {{ compactButtons ? opt.compact : opt.label }}
        </button>
      </div>
    </div>

    <div class="space-y-2 bg-muted/40 rounded-[4px] p-3 border border-border">
      <div
        v-if="showEmptyPrompt && modelValue.length === 0"
        class="text-center py-4 text-muted-foreground text-[11px]"
      >
        Saisissez d'abord le Prix Total ci-dessus pour générer les tranches de
        paiement.
      </div>

      <div
        v-for="(tranche, idx) in modelValue"
        :key="tranche.numero"
        class="flex items-center gap-2 bg-card p-2 rounded-[4px] border border-border"
      >
        <span
          class="w-5 h-5 rounded-full bg-brand/10 text-brand font-bold text-[10px] flex items-center justify-center shrink-0"
        >
          {{ tranche.numero }}
        </span>

        <div class="flex-1">
          <input
            type="text"
            :value="tranche.intitule"
            @input="
              handleFieldUpdate(
                idx,
                'intitule',
                ($event.target as HTMLInputElement).value,
              )
            "
            placeholder="Libellé de la tranche"
            class="w-full h-7 px-2 rounded border border-border bg-background text-foreground text-[11px]"
          />
        </div>

        <div class="w-32">
          <input
            type="number"
            :value="tranche.montant"
            min="0"
            step="any"
            placeholder="Montant FCFA"
            @input="
              handleMontantInput(idx, ($event.target as HTMLInputElement).value)
            "
            class="w-full h-7 px-2 rounded border border-border bg-background text-foreground font-bold text-[11px] tabular-nums"
          />
        </div>

        <div class="w-36 sm:w-40">
          <input
            type="text"
            :value="tranche.echeance"
            @input="
              handleFieldUpdate(
                idx,
                'echeance',
                ($event.target as HTMLInputElement).value,
              )
            "
            placeholder="Échéance"
            class="w-full h-7 px-2 rounded border border-border bg-background text-muted-foreground text-[10px]"
          />
        </div>

        <button
          type="button"
          @click="supprimerTranche(idx)"
          class="text-muted-foreground hover:text-rose-500 p-1"
          title="Supprimer cette tranche"
        >
          <Trash2 class="w-3.5 h-3.5" />
        </button>
      </div>

      <div
        v-if="modelValue.length > 0"
        class="flex items-center justify-between pt-2 text-[11px]"
      >
        <div class="flex items-center gap-2">
          <button
            type="button"
            @click="ajouterTranche"
            class="inline-flex items-center gap-1 text-brand hover:underline font-semibold text-[11px]"
          >
            <Plus class="w-3 h-3" />
            <span>Ajouter une tranche</span>
          </button>
          <span class="text-muted-foreground">·</span>
          <button
            type="button"
            @click="equilibrerTranches"
            class="inline-flex items-center gap-1 text-muted-foreground hover:text-foreground text-[11px]"
          >
            <RefreshCw class="w-3 h-3" />
            <span>Répartir équitablement</span>
          </button>
        </div>

        <div class="flex items-center gap-2 font-semibold">
          <span>Total tranches :</span>
          <span
            :class="[
              'tabular-nums text-xs font-bold',
              tranchesEcart === 0
                ? 'text-emerald-600 dark:text-emerald-400'
                : 'text-amber-600 dark:text-amber-400',
            ]"
          >
            {{ formatFcfa(tranchesTotal) }}
          </span>
        </div>
      </div>

      <div
        v-if="modelValue.length > 0 && tranchesEcart !== 0"
        class="rounded bg-amber-50 dark:bg-amber-950/40 border border-amber-300 dark:border-amber-800 p-2 text-[10px] text-amber-800 dark:text-amber-300 flex items-center justify-between"
      >
        <div class="flex items-center gap-1.5">
          <AlertTriangle class="w-3.5 h-3.5 shrink-0" />
          <span>
            {{
              tranchesEcart < 0
                ? `Reste à répartir sur les tranches : ${formatFcfa(-tranchesEcart)}`
                : `Dépassement du prix total : +${formatFcfa(tranchesEcart)}`
            }}
          </span>
        </div>
        <button
          type="button"
          @click="equilibrerTranches"
          class="font-bold underline text-[10px]"
        >
          Ajuster automatiquement
        </button>
      </div>
    </div>
  </div>
</template>

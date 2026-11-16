<script setup lang="ts">
import { computed } from "vue";
import type { FormulePermis } from "../../../types";
import { formatFcfa, parseTranches } from "./tranchesUtils";
import {
  Clock,
  Car,
  Layers,
  CheckCircle2,
  Trash2,
  Edit3,
  Lock,
} from "lucide-vue-next";

const props = defineProps<{
  formule: FormulePermis;
  canModify: boolean;
}>();

const emit = defineEmits<{
  (e: "edit", formule: FormulePermis): void;
  (e: "delete", formule: FormulePermis): void;
}>();

const tranches = computed(() => parseTranches(props.formule.tranches_paiement));
</script>

<template>
  <div
    class="relative flex flex-col justify-between rounded-[4px] border border-border bg-card p-5 hover:border-brand/40 transition shadow-xs"
  >
    <!-- En-tête Formule -->
    <div>
      <div class="flex items-start justify-between gap-3">
        <div>
          <span
            class="inline-flex items-center px-2 py-0.5 rounded text-[10px] font-bold bg-muted text-foreground border border-border"
          >
            {{
              formule.categorie_permis === "REC" ||
              formule.categorie_permis === "RECYCLAGE"
                ? "Recyclage / Performance"
                : `Permis ${formule.categorie_permis}`
            }}
          </span>
          <h3 class="text-base font-bold text-foreground mt-2 leading-snug">
            {{ formule.intitule }}
          </h3>
        </div>
        <div class="text-right">
          <span class="text-lg font-bold text-brand tabular-nums">
            {{ formatFcfa(formule.prix_total) }}
          </span>
          <span class="block text-[10px] text-muted-foreground">
            Acompte min. {{ formatFcfa(formule.acompte_minimum) }}
          </span>
        </div>
      </div>

      <p v-if="formule.description" class="text-xs text-muted-foreground mt-2">
        {{ formule.description }}
      </p>

      <!-- Détails heures théoriques & pratiques -->
      <div class="mt-4 grid grid-cols-2 gap-2 text-xs">
        <div
          class="flex items-center space-x-2 rounded-[4px] bg-muted/50 p-2.5 border border-border"
        >
          <Clock class="w-4 h-4 text-brand" />
          <div>
            <span class="block font-bold text-foreground"
              >{{ formule.heures_theoriques }}h</span
            >
            <span class="text-[10px] text-muted-foreground">Code en salle</span>
          </div>
        </div>
        <div
          class="flex items-center space-x-2 rounded-[4px] bg-muted/50 p-2.5 border border-border"
        >
          <Car class="w-4 h-4 text-emerald-600" />
          <div>
            <span class="block font-bold text-foreground"
              >{{ formule.heures_pratiques }}h</span
            >
            <span class="text-[10px] text-muted-foreground"
              >Conduite réelle</span
            >
          </div>
        </div>
      </div>

      <!-- Section Échéancier des Tranches de Paiement -->
      <div class="mt-4 pt-3 border-t border-border/70 space-y-2">
        <div class="flex items-center justify-between">
          <span
            class="text-[11px] font-semibold text-foreground flex items-center gap-1.5"
          >
            <Layers class="w-3.5 h-3.5 text-brand" />
            <span>Modalités de Paiement Échelonné</span>
          </span>
          <span
            v-if="tranches.length > 0"
            class="text-[10px] font-bold px-1.5 py-0.5 rounded bg-emerald-50 text-emerald-700 dark:bg-emerald-950/40 dark:text-emerald-300 border border-emerald-200 dark:border-emerald-800"
          >
            {{ tranches.length }} Tranche{{ tranches.length > 1 ? "s" : "" }}
          </span>
          <span v-else class="text-[10px] text-muted-foreground">
            Acompte + Solde standard
          </span>
        </div>

        <div
          v-if="tranches.length > 0"
          class="space-y-1.5 text-xs bg-muted/30 rounded-[4px] p-2.5 border border-border/60"
        >
          <div
            v-for="t in tranches"
            :key="t.numero"
            class="flex items-center justify-between text-[11px] pb-1 border-b border-border/40 last:border-b-0 last:pb-0"
          >
            <div class="flex items-center gap-2">
              <span
                class="w-4 h-4 rounded-full bg-brand/10 text-brand font-bold text-[9px] flex items-center justify-center shrink-0"
              >
                {{ t.numero }}
              </span>
              <div>
                <span class="font-medium text-foreground">{{
                  t.intitule
                }}</span>
                <span
                  v-if="t.echeance"
                  class="block text-[9px] text-muted-foreground"
                >
                  {{ t.echeance }}
                </span>
              </div>
            </div>
            <span class="font-bold text-foreground tabular-nums">
              {{ formatFcfa(t.montant) }}
            </span>
          </div>
        </div>

        <div
          v-else
          class="text-[11px] text-muted-foreground bg-muted/20 rounded-[4px] p-2.5 border border-dashed border-border"
        >
          <div class="flex justify-between">
            <span>1. Acompte initial obligatoire :</span>
            <strong class="text-foreground tabular-nums">{{
              formatFcfa(formule.acompte_minimum)
            }}</strong>
          </div>
          <div class="flex justify-between mt-1">
            <span>2. Solde restant :</span>
            <strong class="text-foreground tabular-nums">{{
              formatFcfa(
                Math.max(0, formule.prix_total - formule.acompte_minimum),
              )
            }}</strong>
          </div>
          <p
            v-if="canModify"
            class="text-[10px] text-brand hover:underline cursor-pointer mt-1.5 font-medium"
            @click="emit('edit', formule)"
          >
            + Configurer l'échéancier détaillé en 2, 3 ou 4 tranches
          </p>
        </div>
      </div>
    </div>

    <!-- Pied de Carte : Contrôles d'Habilitation -->
    <div
      class="mt-5 pt-3 border-t border-border flex items-center justify-between"
    >
      <span
        class="text-[11px] text-muted-foreground flex items-center space-x-1"
      >
        <CheckCircle2 class="w-3.5 h-3.5 text-emerald-500" />
        <span>Catalogue Actif</span>
      </span>

      <div class="flex items-center gap-1.5">
        <button
          v-if="canModify"
          @click="emit('delete', formule)"
          type="button"
          class="inline-flex items-center p-1.5 rounded-[4px] text-muted-foreground hover:text-rose-600 hover:bg-rose-50 dark:hover:bg-rose-950/40 border border-transparent hover:border-rose-200 dark:hover:border-rose-900 transition active:scale-95 cursor-pointer"
          title="Supprimer cette formule du catalogue"
        >
          <Trash2 class="w-3.5 h-3.5" />
        </button>

        <button
          v-if="canModify"
          @click="emit('edit', formule)"
          type="button"
          class="inline-flex items-center space-x-1.5 px-3 py-1.5 rounded-[4px] bg-muted hover:bg-brand-light hover:text-brand-dark border border-border hover:border-brand-border text-xs font-semibold text-foreground transition active:scale-95 cursor-pointer"
          title="Modifier le tarif officiel ou les tranches de ce permis"
        >
          <Edit3 class="w-3.5 h-3.5 text-brand" />
          <span>Modifier Tarif & Tranches</span>
        </button>

        <button
          v-else
          @click="emit('edit', formule)"
          type="button"
          class="inline-flex items-center space-x-1.5 px-2.5 py-1.5 rounded-[4px] bg-muted/50 text-muted-foreground border border-border/60 text-xs font-medium cursor-not-allowed hover:bg-muted transition"
          title="Seul le Directeur peut modifier un tarif fixé"
        >
          <Lock class="w-3.5 h-3.5 text-muted-foreground" />
          <span>Tarif fixé (Directeur)</span>
        </button>
      </div>
    </div>
  </div>
</template>

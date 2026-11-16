<script setup lang="ts">
import { computed } from "vue";
import type { Paiement } from "../../../types";
import { Printer } from "lucide-vue-next";

const props = defineProps<{
  paiements?: Paiement[] | null;
}>();

const listePaiements = computed(() =>
  Array.isArray(props.paiements) ? props.paiements : [],
);

const emit = defineEmits<{
  (e: "reprint", id: string): void;
}>();

const fcfa = (n: number) => n.toLocaleString("fr-FR");

const modeLabels: Record<string, string> = {
  ESPECES: "Espèces",
  WAVE: "Wave",
  ORANGE_MONEY: "Orange Money",
  MTN_MOMO: "MTN MoMo",
};
const modeLabel = (m: string) => modeLabels[m] ?? m;

const formatDateHeure = (iso: string) => {
  const d = new Date(iso);
  return `${d.toLocaleDateString("fr-FR")} · ${d.toLocaleTimeString("fr-FR", {
    hour: "2-digit",
    minute: "2-digit",
  })}`;
};
</script>

<template>
  <section class="p-4">
    <div class="mb-1 flex items-baseline justify-between">
      <h3 class="text-sm font-semibold text-foreground">Encaissements</h3>
      <span class="text-[11px] text-muted-foreground"
        >{{ listePaiements.length }} versement{{
          listePaiements.length > 1 ? "s" : ""
        }}</span
      >
    </div>

    <ul v-if="listePaiements.length" class="divide-y divide-border">
      <li
        v-for="p in listePaiements"
        :key="p.id"
        class="flex items-center justify-between gap-3 py-2.5 text-xs"
      >
        <div class="min-w-0">
          <div class="flex flex-wrap items-center gap-x-2">
            <span class="font-mono font-medium text-foreground">{{
              p.reference_recu
            }}</span>
            <span class="text-muted-foreground">{{
              modeLabel(p.mode_reglement)
            }}</span>
          </div>
          <div class="truncate text-[11px] text-muted-foreground">
            {{ formatDateHeure(p.date_paiement) }}
            <template v-if="p.reference_transaction">
              · Réf. {{ p.reference_transaction }}</template
            >
          </div>
        </div>

        <div class="flex shrink-0 items-center gap-3">
          <div class="text-right">
            <div class="font-medium tabular-nums text-brand-dark">
              +{{ fcfa(p.montant) }} F
            </div>
            <div class="text-[11px] tabular-nums text-muted-foreground">
              Solde {{ fcfa(p.solde_apres_paiement) }} F
            </div>
          </div>

          <button
            type="button"
            @click="emit('reprint', p.id)"
            title="Réimprimer le ticket officiel (marqué DUPLICATA)"
            class="inline-flex h-8 items-center gap-1.5 rounded-lg border border-border bg-card px-2.5 text-[11px] font-medium text-muted-foreground transition hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40"
          >
            <Printer class="h-3.5 w-3.5 text-muted-foreground" />
            Duplicata
          </button>
        </div>
      </li>
    </ul>

    <p v-else class="py-4 text-center text-xs text-muted-foreground">
      Aucun encaissement enregistré.
    </p>
  </section>
</template>


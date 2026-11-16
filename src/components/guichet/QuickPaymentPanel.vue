<script setup lang="ts">
import { ref, computed, watch } from "vue";
import { useGuichetStore } from "../../stores/guichetStore";
import { useAuthStore } from "../../stores/authStore";
import { CheckCircle2, Printer } from "lucide-vue-next";

const store = useGuichetStore();
const authStore = useAuthStore();

const MONTANT_PAR_DEFAUT = 50000;
const PRESETS = [25000, 50000, 100000];

const nomCaissier = computed(() => {
  if (authStore.currentUser.prenom || authStore.currentUser.nom) {
    return `${authStore.currentUser.prenom} ${authStore.currentUser.nom}`.trim();
  }
  return "Opérateur Guichet";
});

const modes = [
  { value: "ESPECES", label: "Espèces", dot: "bg-brand" },
  { value: "WAVE", label: "Wave", dot: "bg-sky-500" },
  { value: "ORANGE_MONEY", label: "Orange Money", dot: "bg-orange-500" },
  { value: "MTN_MOMO", label: "MTN MoMo", dot: "bg-yellow-400" },
  { value: "MOOV", label: "Moov", dot: "bg-blue-500" },
];

const montant = ref<number>(MONTANT_PAR_DEFAUT);
const modeReglement = ref<string>("ESPECES");
const referenceTransaction = ref<string>("");
const especesRemises = ref<number>(MONTANT_PAR_DEFAUT);

const detail = computed(() => store.activeDossierDetail);
const solde = computed(() => detail.value?.dossier.solde_restant ?? 0);

const fcfa = (n: number) => (n ?? 0).toLocaleString("fr-FR");
const today = new Date().toLocaleDateString("fr-FR");

// Le montant proposé ne dépasse jamais le solde restant
watch(
  () => detail.value?.dossier.solde_restant,
  (newSolde) => {
    if (newSolde !== undefined && newSolde > 0) {
      montant.value = Math.min(MONTANT_PAR_DEFAUT, newSolde);
      especesRemises.value = montant.value;
    }
  },
  { immediate: true },
);

watch(montant, (m) => {
  if (especesRemises.value < m) especesRemises.value = m;
});

const montantInvalide = computed(
  () => !montant.value || montant.value <= 0 || montant.value > solde.value,
);

const renduMonnaie = computed(() => {
  if (modeReglement.value !== "ESPECES") return 0;
  return Math.max(0, (especesRemises.value || 0) - (montant.value || 0));
});

const manque = computed(() => {
  if (modeReglement.value !== "ESPECES") return 0;
  return Math.max(0, (montant.value || 0) - (especesRemises.value || 0));
});

const soldeApres = computed(() =>
  Math.max(0, solde.value - (montant.value || 0)),
);
const modeLabel = computed(
  () =>
    modes.find((m) => m.value === modeReglement.value)?.label ??
    modeReglement.value,
);

const setPresetAmount = (val: number) => {
  if (!detail.value) return;
  montant.value = Math.min(val, solde.value);
  especesRemises.value = montant.value;
};

const payFullBalance = () => {
  if (!detail.value) return;
  montant.value = solde.value;
  especesRemises.value = montant.value;
};

const submitPayment = () => {
  if (!detail.value) return;
  if (montant.value <= 0) {
    store.showToast("Saisissez un montant supérieur à 0 FCFA.", "error");
    return;
  }
  if (montant.value > solde.value) {
    store.showToast("Le montant saisi dépasse le solde restant.", "error");
    return;
  }
  if (modeReglement.value === "ESPECES" && manque.value > 0) {
    store.showToast(
      `Espèces insuffisantes : il manque ${fcfa(manque.value)} FCFA.`,
      "error",
    );
    return;
  }

  store.processPayment({
    montant: montant.value,
    mode_reglement: modeReglement.value,
    reference_transaction: referenceTransaction.value.trim() || undefined,
  });
};

const field =
  "h-8 w-full rounded-lg border border-border bg-card px-2.5 text-xs text-foreground transition placeholder:text-muted-foreground focus:border-brand focus:outline-none focus:ring-2 focus:ring-brand/20";
</script>

<template>
  <div
    class="flex flex-1 min-h-0 md:h-full flex-col overflow-hidden rounded-[5px] border border-border bg-card"
  >
    <header
      class="flex items-center justify-between gap-3 border-b border-border px-5 py-3.5 shrink-0"
    >
      <h2 class="flex items-center gap-2 text-sm font-semibold text-foreground">
        Encaissement
        <kbd
          class="rounded bg-muted px-1.5 py-0.5 font-mono text-[10px] font-normal text-muted-foreground"
          >F2</kbd
        >
      </h2>
      <span
        class="flex items-center gap-1.5 text-[11px] font-medium text-muted-foreground"
      >
        <span class="h-1.5 w-1.5 rounded-full bg-brand" />
        Prêt
      </span>
    </header>

    <div
      v-if="detail"
      class="flex-1 min-h-0 divide-y divide-border overflow-y-auto"
    >
      <!-- Dossier -->
      <section
        class="flex items-center justify-between gap-4 bg-muted/40 px-5 py-4"
      >
        <div class="min-w-0">
          <p class="truncate text-xs font-semibold text-foreground">
            {{ detail.candidat.nom }} {{ detail.candidat.prenoms }}
          </p>
          <p class="font-mono text-[11px] text-muted-foreground">
            {{ detail.dossier.numero_dossier }}
          </p>
        </div>
        <div class="shrink-0 text-right">
          <p class="text-[11px] text-muted-foreground">Reste à payer</p>
          <p
            class="text-sm font-semibold tabular-nums"
            :class="solde > 0 ? 'text-amber-600' : 'text-brand-dark'"
          >
            {{ fcfa(solde) }} FCFA
          </p>
        </div>
      </section>

      <!-- Dossier soldé -->
      <section
        v-if="solde === 0"
        class="flex flex-col items-center gap-1.5 px-4 py-10 text-center"
      >
        <CheckCircle2 class="h-6 w-6 text-brand" />
        <p class="text-xs font-semibold text-foreground">
          Dossier entièrement soldé
        </p>
        <p class="max-w-xs text-[11px] text-muted-foreground">
          Aucun paiement supplémentaire n'est dû pour cette formation.
        </p>
      </section>

      <template v-else>
        <!-- Montant -->
        <section class="space-y-3 px-5 py-5">
          <label
            for="montant-input"
            class="block text-xs font-semibold text-foreground"
            >Montant encaissé</label
          >

          <div class="relative">
            <input
              id="montant-input"
              type="number"
              v-model.number="montant"
              step="1000"
              min="1000"
              :max="solde"
              :aria-invalid="montantInvalide"
              class="h-10 w-full rounded-[4px] border bg-card pl-3 pr-14 text-right text-base font-semibold tabular-nums text-foreground transition focus:outline-none focus:ring-2"
              :class="
                montantInvalide
                  ? 'border-destructive/40 focus:border-destructive focus:ring-destructive/20'
                  : 'border-border focus:border-brand focus:ring-brand/20'
              "
            />
            <span
              class="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 text-[11px] text-muted-foreground"
              >FCFA</span
            >
          </div>

          <p v-if="montant > solde" class="text-[11px] text-destructive">
            Le montant dépasse le solde restant.
          </p>

          <div class="flex flex-wrap items-center gap-1.5">
            <button
              v-for="p in PRESETS"
              :key="p"
              type="button"
              @click="setPresetAmount(p)"
              class="h-7 rounded-[4px] border border-border bg-card px-2.5 text-[11px] font-medium tabular-nums text-foreground transition hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40"
            >
              {{ fcfa(p) }}
            </button>
            <button
              type="button"
              @click="payFullBalance"
              class="ml-auto h-7 rounded-[4px] border border-amber-200 bg-amber-50 px-2.5 text-[11px] font-medium tabular-nums text-amber-800 transition hover:bg-amber-100 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-500/40"
            >
              Tout régler · {{ fcfa(solde) }} F
            </button>
          </div>
        </section>

        <!-- Mode de règlement -->
        <section class="space-y-3 px-5 py-5">
          <h3 id="mode-label" class="text-xs font-semibold text-foreground">
            Mode de règlement
          </h3>

          <div
            class="grid grid-cols-3 gap-2"
            role="radiogroup"
            aria-labelledby="mode-label"
          >
            <button
              v-for="m in modes"
              :key="m.value"
              type="button"
              role="radio"
              :aria-checked="modeReglement === m.value"
              @click="modeReglement = m.value"
              class="flex h-9 items-center justify-center gap-1.5 rounded-[4px] border px-2.5 text-xs font-medium transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40"
              :class="
                modeReglement === m.value
                  ? 'border-brand bg-brand-light text-brand-dark font-semibold'
                  : 'border-border bg-card text-foreground/80 hover:bg-muted'
              "
            >
              <span class="h-1.5 w-1.5 shrink-0 rounded-full" :class="m.dot" />
              <span class="truncate">{{ m.label }}</span>
            </button>
          </div>

          <div v-if="modeReglement !== 'ESPECES'">
            <label
              for="ref-transaction"
              class="mb-1 block text-xs text-muted-foreground"
              >Référence de transaction (SMS)</label
            >
            <input
              id="ref-transaction"
              type="text"
              v-model="referenceTransaction"
              placeholder="TX-CI-9842104"
              :class="[field, 'font-mono']"
            />
          </div>
        </section>

        <!-- Espèces reçues -->
        <section v-if="modeReglement === 'ESPECES'" class="space-y-3 px-5 py-5">
          <h3 class="text-xs font-semibold text-foreground">
            Espèces reçues et rendu
          </h3>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label
                for="especes-remises"
                class="mb-1 block text-xs text-muted-foreground"
                >Montant remis</label
              >
              <div class="relative">
                <input
                  id="especes-remises"
                  type="number"
                  v-model.number="especesRemises"
                  step="1000"
                  min="0"
                  class="h-8 w-full rounded-[4px] border border-border bg-card pl-2.5 pr-8 text-right text-xs font-semibold tabular-nums text-foreground transition focus:border-brand focus:outline-none focus:ring-2 focus:ring-brand/20"
                />
                <span
                  class="pointer-events-none absolute right-2.5 top-1/2 -translate-y-1/2 text-[11px] text-muted-foreground"
                  >F</span
                >
              </div>
            </div>

            <div>
              <p class="mb-1 text-xs text-muted-foreground">
                {{ manque > 0 ? "Manquant" : "Rendu" }}
              </p>
              <p
                class="flex h-8 items-center justify-end rounded-[4px] bg-muted px-2.5 text-xs font-semibold tabular-nums"
                :class="manque > 0 ? 'text-destructive' : 'text-brand-dark'"
                aria-live="polite"
              >
                {{ fcfa(manque > 0 ? manque : renduMonnaie) }} FCFA
              </p>
            </div>
          </div>
        </section>

        <!-- Validation -->
        <section class="px-5 py-5">
          <button
            type="button"
            @click="submitPayment"
            :disabled="store.loading"
            class="inline-flex h-10 w-full items-center justify-center gap-2 rounded-[4px] bg-brand-dark px-4 text-xs font-semibold text-white transition hover:bg-brand focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40 focus-visible:ring-offset-1 disabled:cursor-not-allowed disabled:opacity-50"
          >
            <Printer class="h-3.5 w-3.5" />
            Valider et imprimer le reçu
            <kbd
              class="rounded bg-white/20 px-1 font-mono text-[10px] text-white/90"
              >Entrée</kbd
            >
          </button>
        </section>
      </template>

      <!-- Aperçu du ticket -->
      <section class="px-4 py-4">
        <div class="mb-2 flex items-baseline justify-between">
          <h3 class="text-xs font-semibold text-foreground">
            Aperçu du ticket
          </h3>
          <span class="text-[11px] text-muted-foreground">80 mm · ESC/POS</span>
        </div>

        <div
          class="mx-auto max-w-xs select-text space-y-2 rounded-[5px] border border-border bg-card p-3.5 font-mono text-[10px] leading-snug text-foreground"
        >
          <div class="border-b border-dashed border-border pb-2 text-center">
            <div class="text-xs font-bold">
              {{ authStore.schoolName.toUpperCase() }}
            </div>
            <div class="text-[9px] text-muted-foreground">
              {{ authStore.schoolAddress }}
            </div>
            <div
              v-if="authStore.schoolPhone"
              class="text-[9px] text-muted-foreground"
            >
              Tél. {{ authStore.schoolPhone }}
            </div>
            <div class="mt-1 text-[9px] font-bold">REÇU DE CAISSE</div>
          </div>

          <div class="space-y-0.5 text-[9px]">
            <div class="flex justify-between gap-2">
              <span>Date : {{ today }}</span>
              <span>Caissière : {{ nomCaissier }}</span>
            </div>
            <div class="flex justify-between gap-2">
              <span>Dossier : {{ detail.dossier.numero_dossier }}</span>
              <span class="font-bold"
                >Permis {{ detail.formule.categorie_permis }}</span
              >
            </div>
            <div>
              Élève : {{ detail.candidat.nom }} {{ detail.candidat.prenoms }}
            </div>
          </div>

          <div class="space-y-1 border-y border-dashed border-border py-1.5">
            <div class="flex justify-between gap-2 text-[11px] font-bold">
              <span>ENCAISSÉ</span>
              <span class="tabular-nums">{{ fcfa(montant || 0) }} FCFA</span>
            </div>
            <div
              class="flex justify-between gap-2 text-[9px] text-muted-foreground"
            >
              <span>Mode : {{ modeLabel }}</span>
              <span class="tabular-nums">Solde : {{ fcfa(soldeApres) }} F</span>
            </div>
          </div>

          <div class="text-center text-[9px] text-muted-foreground">
            <div>Code de sécurité généré à l'impression</div>
            <div class="mt-0.5 italic">Merci de votre confiance.</div>
          </div>
        </div>
      </section>
    </div>

    <div
      v-else
      class="flex flex-1 items-center justify-center p-8 text-center text-xs text-muted-foreground"
    >
      Sélectionnez un dossier pour encaisser une tranche.
    </div>
  </div>
</template>

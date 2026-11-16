<script setup lang="ts">
import { ref, computed, watch } from "vue";
import { useGuichetStore } from "../../stores/guichetStore";
import { useAuthStore } from "../../stores/authStore";
import { X, Printer } from "lucide-vue-next";

const store = useGuichetStore();
const authStore = useAuthStore();

const activeCaissiere = computed(() => {
  if (store.sessionCaisseDetail?.caissier_nom) {
    return store.sessionCaisseDetail.caissier_nom;
  }
  if (authStore.currentUser?.nom) {
    return `${authStore.currentUser.prenom} ${authStore.currentUser.nom} (${authStore.currentUser.roleLabel || "Guichet"})`;
  }
  return "Agent Guichet";
});

/* ───────────────────────── État ───────────────────────── */
const fondDeCaisse = computed(
  () =>
    store.sessionCaisseDetail?.fond_de_caisse ??
    store.stats.fond_de_caisse ??
    0,
);
const especesTheoriques = computed(
  () =>
    store.sessionCaisseDetail?.total_especes_theorique ??
    store.stats.total_especes_theorique ??
    0,
);
const totalTheoriqueAttendu = computed(
  () => fondDeCaisse.value + especesTheoriques.value,
);

const comptagePhysique = ref<number>(totalTheoriqueAttendu.value);
const justification = ref("");

const ecartCaisse = computed(
  () => (comptagePhysique.value || 0) - totalTheoriqueAttendu.value,
);

// À chaque ouverture : recharger la session et repartir d'un comptage propre
watch(
  () => store.isCashCloseModalOpen,
  async (open) => {
    if (!open) return;
    await store.loadSessionCaisse();
    comptagePhysique.value = totalTheoriqueAttendu.value;
    justification.value = "";
  },
  { immediate: true },
);

/* ───────────────────────── Affichage ───────────────────────── */
const fcfa = (n: number) => (n || 0).toLocaleString("fr-FR");
const today = new Date().toLocaleDateString("fr-FR");

const ventilation = computed(() => store.sessionCaisseDetail?.ventilation);

const lignes = computed(() => {
  const base = [
    {
      label: "Espèces (tiroir)",
      dot: "bg-brand",
      value:
        ventilation.value?.especes ?? store.stats.total_especes_theorique ?? 0,
      always: true,
    },
    {
      label: "Wave",
      dot: "bg-sky-500",
      value: ventilation.value?.wave ?? 0,
      always: true,
    },
    {
      label: "Orange Money",
      dot: "bg-orange-500",
      value: ventilation.value?.orange_money ?? 0,
      always: true,
    },
    {
      label: "MTN MoMo",
      dot: "bg-yellow-400",
      value: ventilation.value?.mtn_momo ?? 0,
      always: true,
    },
    {
      label: "Moov Money",
      dot: "bg-blue-500",
      value: ventilation.value?.moov ?? 0,
      always: false,
    },
    {
      label: "Virement",
      dot: "bg-slate-500",
      value: ventilation.value?.virement ?? 0,
      always: false,
    },
  ];
  return base.filter((l) => l.always || l.value > 0);
});

const totalEncaisse = computed(() => ventilation.value?.total_encaisse ?? 0);
const nbTransactions = computed(
  () => store.sessionCaisseDetail?.nombre_transactions ?? 0,
);

const ecartInfo = computed(() => {
  if (ecartCaisse.value === 0)
    return {
      label: "Conforme",
      box: "bg-brand-light text-brand-dark",
      dot: "bg-brand",
    };
  if (ecartCaisse.value > 0)
    return {
      label: "Surplus",
      box: "bg-sky-50 text-sky-800",
      dot: "bg-sky-500",
    };
  return {
    label: "Déficit",
    box: "bg-rose-50 text-rose-800",
    dot: "bg-rose-500",
  };
});

const justificationManquante = computed(
  () => ecartCaisse.value !== 0 && !justification.value.trim(),
);

/* ───────────────────────── Actions ───────────────────────── */
const isSubmitting = ref(false);

const close = () => {
  store.isCashCloseModalOpen = false;
};

const triggerPointX = () => {
  store.showToast(
    "Point X intermédiaire généré. Caisse maintenue active.",
    "info",
  );
  store.isCashCloseModalOpen = false;
  store.isRapportZModalOpen = true;
};

const submitClotureZ = async () => {
  if (justificationManquante.value) {
    store.showToast(
      "Un écart de caisse a été constaté. La justification est obligatoire.",
      "error",
    );
    return;
  }

  isSubmitting.value = true;
  try {
    await store.cloturerCaisse({
      total_especes_compte: Number(comptagePhysique.value) || 0,
      observations: justification.value.trim() || undefined,
    });
    store.isCashCloseModalOpen = false;
    store.isRapportZModalOpen = true;
  } catch (err) {
    // Toast déjà affiché
  } finally {
    isSubmitting.value = false;
  }
};
</script>

<template>
  <div
    v-if="store.isCashCloseModalOpen"
    class="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/50 p-0 sm:p-4"
    role="dialog"
    aria-modal="true"
    aria-labelledby="cash-close-title"
    @click.self="close"
  >
    <div
      class="flex max-h-[92dvh] w-full sm:max-w-lg flex-col overflow-hidden rounded-t-[3px] sm:rounded-[3px] border-t sm:border border-border bg-card pb-safe"
    >
      <div
        class="sm:hidden mx-auto mt-2 h-1 w-9 rounded-full bg-border shrink-0"
      />
      <!-- ═════════ En-tête ═════════ -->
      <header
        class="flex items-center justify-between gap-3 border-b border-border px-4 py-3"
      >
        <div class="min-w-0">
          <h2
            id="cash-close-title"
            class="flex items-center gap-2 text-sm font-semibold text-foreground"
          >
            Gestion de caisse
            <kbd
              class="rounded bg-muted px-1 font-mono text-[10px] font-normal text-muted-foreground"
              >F4</kbd
            >
          </h2>
          <p class="text-[11px] text-muted-foreground">
            Point X intermédiaire ou clôture Z définitive.
          </p>
        </div>
        <button
          type="button"
          @click="close"
          class="flex h-7 w-7 items-center justify-center rounded-[4px] text-muted-foreground transition hover:bg-muted hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40"
          aria-label="Fermer"
        >
          <X class="h-4 w-4" />
        </button>
      </header>

      <!-- ═════════ Corps ═════════ -->
      <div class="flex-1 divide-y divide-border overflow-y-auto">
        <!-- Session -->
        <section
          class="flex items-center justify-between gap-3 px-4 py-3 text-xs"
        >
          <div class="min-w-0">
            <p class="text-muted-foreground">Session active</p>
            <p class="truncate font-medium text-foreground">
              {{ activeCaissiere }}
            </p>
          </div>
          <div class="shrink-0 text-right">
            <p class="text-muted-foreground">Date</p>
            <p class="font-medium tabular-nums text-foreground">{{ today }}</p>
          </div>
        </section>

        <!-- 1. Recettes -->
        <section class="px-4 py-3">
          <div class="mb-2 flex items-baseline justify-between">
            <h3 class="text-xs font-semibold text-foreground">
              1. Recettes du jour
            </h3>
            <span class="text-[11px] text-muted-foreground"
              >{{ nbTransactions }} transactions</span
            >
          </div>

          <dl class="space-y-1.5 text-xs">
            <div
              v-for="l in lignes"
              :key="l.label"
              class="flex items-center justify-between gap-3"
            >
              <dt class="flex items-center gap-2 text-muted-foreground">
                <span class="h-1.5 w-1.5 rounded-full" :class="l.dot" />
                {{ l.label }}
              </dt>
              <dd class="font-medium tabular-nums text-foreground">
                {{ fcfa(l.value) }} FCFA
              </dd>
            </div>
          </dl>

          <div
            class="mt-2.5 flex items-baseline justify-between border-t border-border pt-2.5 text-xs"
          >
            <span class="font-medium text-foreground"
              >Total encaissé · tous modes</span
            >
            <span class="text-sm font-semibold tabular-nums text-foreground"
              >{{ fcfa(totalEncaisse) }} FCFA</span
            >
          </div>
        </section>

        <!-- 2. Comptage physique -->
        <section class="space-y-3 px-4 py-3">
          <h3 class="text-xs font-semibold text-foreground">
            2. Comptage du tiroir-caisse
          </h3>

          <dl class="grid grid-cols-2 gap-x-4 text-xs">
            <div>
              <dt class="text-muted-foreground">Fond de caisse</dt>
              <dd class="font-medium tabular-nums text-foreground">
                {{ fcfa(fondDeCaisse) }} FCFA
              </dd>
            </div>
            <div>
              <dt class="text-muted-foreground">Espèces attendues</dt>
              <dd class="font-medium tabular-nums text-brand-dark">
                {{ fcfa(totalTheoriqueAttendu) }} FCFA
              </dd>
            </div>
          </dl>

          <div class="relative">
            <label for="comptage-physique" class="sr-only"
              >Montant compté dans le tiroir en FCFA</label
            >
            <input
              id="comptage-physique"
              type="number"
              v-model.number="comptagePhysique"
              step="500"
              min="0"
              class="h-10 w-full rounded-[4px] border border-border bg-card pl-3 pr-14 text-right text-base font-semibold tabular-nums text-foreground transition focus:border-brand focus:outline-none focus:ring-2 focus:ring-brand/20"
            />
            <span
              class="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 text-[11px] text-muted-foreground"
            >
              FCFA
            </span>
          </div>

          <!-- Écart en direct -->
          <div
            class="flex items-center justify-between rounded-[4px] px-3 py-2 text-xs"
            :class="ecartInfo.box"
            aria-live="polite"
          >
            <span class="flex items-center gap-2 font-medium">
              <span class="h-1.5 w-1.5 rounded-full" :class="ecartInfo.dot" />
              Écart · {{ ecartInfo.label }}
            </span>
            <span class="font-semibold tabular-nums">
              {{ ecartCaisse > 0 ? "+" : "" }}{{ fcfa(ecartCaisse) }} FCFA
            </span>
          </div>

          <!-- Justification obligatoire -->
          <div v-if="ecartCaisse !== 0">
            <label
              for="justification-ecart"
              class="mb-1 block text-xs font-medium text-foreground"
            >
              Justification de l'écart <span class="text-destructive">*</span>
            </label>
            <textarea
              id="justification-ecart"
              v-model="justification"
              rows="2"
              placeholder="Cause du surplus ou du déficit…"
              :aria-invalid="justificationManquante"
              class="w-full resize-none rounded-[4px] border bg-card p-2.5 text-xs text-foreground transition placeholder:text-muted-foreground focus:outline-none focus:ring-2"
              :class="
                justificationManquante
                  ? 'border-destructive/40 focus:border-destructive focus:ring-destructive/20'
                  : 'border-border focus:border-brand focus:ring-brand/20'
              "
            ></textarea>
          </div>
        </section>
      </div>

      <!-- ═════════ Actions ═════════ -->
      <footer
        class="flex items-center justify-between gap-2 border-t border-border px-4 py-3"
      >
        <button
          type="button"
          @click="close"
          class="inline-flex h-8 items-center gap-2 rounded-[4px] px-3 text-xs font-medium text-muted-foreground transition hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40"
        >
          Fermer
          <kbd
            class="rounded bg-muted px-1 font-mono text-[10px] text-muted-foreground"
            >Échap</kbd
          >
        </button>

        <div class="flex items-center gap-2">
          <button
            type="button"
            @click="triggerPointX"
            title="Imprimer un ticket intermédiaire sans fermer la caisse"
            class="inline-flex h-8 items-center gap-1.5 rounded-[4px] border border-border bg-card px-3 text-xs font-medium text-foreground transition hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40"
          >
            <Printer class="h-3.5 w-3.5 text-muted-foreground" />
            Point X
          </button>

          <button
            type="button"
            @click="submitClotureZ"
            :disabled="isSubmitting"
            class="inline-flex h-8 items-center gap-1.5 rounded-[4px] bg-brand-dark px-3 text-xs font-medium text-white transition hover:bg-brand focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40 focus-visible:ring-offset-1 disabled:opacity-50"
          >
            <Printer class="h-3.5 w-3.5" />
            {{ isSubmitting ? "Clôture en cours..." : "Clôture Z" }}
            <kbd
              v-if="!isSubmitting"
              class="rounded bg-white/20 px-1 font-mono text-[10px] text-white/90"
              >Entrée</kbd
            >
          </button>
        </div>
      </footer>
    </div>
  </div>
</template>

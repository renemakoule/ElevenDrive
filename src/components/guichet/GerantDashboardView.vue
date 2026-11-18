<script setup lang="ts">
import { ref, computed, watch, onMounted } from "vue";
import { useGuichetStore } from "../../stores/guichetStore";
import { useAuthStore } from "../../stores/authStore";
import { realtimeSocket } from "../../services/realtimeSocketService";
import GerantImpayesTab from "./gerant/GerantImpayesTab.vue";
import GerantAssidusTab from "./gerant/GerantAssidusTab.vue";
import GerantOverridePinModal from "./gerant/GerantOverridePinModal.vue";
import GerantDashboardSkeleton from "../common/skeleton/GerantDashboardSkeleton.vue";
import { ArrowLeft, Check, CheckCircle2, Radio } from "lucide-vue-next";

const store = useGuichetStore();
const authStore = useAuthStore();
const isWsConnected = realtimeSocket.isConnected;
type TabKey = "impayes" | "assidus" | "recette";
const activeTab = ref<TabKey>("impayes");

const isOverrideModalOpen = ref(false);
const candToOverride = ref<any>(null);

const isInitialLoading = ref(!store.gerantDashboardData);

onMounted(async () => {
  try {
    await store.loadGerantDashboard();
  } finally {
    isInitialLoading.value = false;
  }
});

watch(
  () => authStore.activeSection,
  (sec) => {
    if (sec === "derogations") {
      activeTab.value = "assidus";
    }
  },
  { immediate: true },
);

const data = computed(() => store.gerantDashboardData);
const impayes = computed(() => data.value?.impayes ?? []);
const assidus = computed(() => data.value?.eleves_assidus_code ?? []);
const totalCreances = computed(() => data.value?.total_creances_impayes ?? 0);

const fcfa = (n: number) => (n ?? 0).toLocaleString("fr-FR");

const tabs = computed(() => [
  {
    key: "impayes" as TabKey,
    label: "Impayés et relances",
    count: impayes.value.length,
  },
  {
    key: "assidus" as TabKey,
    label: "Élèves assidus",
    count: assidus.value.length,
  },
  {
    key: "recette" as TabKey,
    label: "Recette MVP",
    count: undefined as number | undefined,
  },
]);

const stats = computed(() => [
  {
    key: "creances",
    label: "Créances impayées",
    value: `${fcfa(totalCreances.value)} FCFA`,
    tone: "text-amber-600",
    caption: `${impayes.value.length} dossiers concernés`,
  },
  {
    key: "assidus",
    label: "Élèves assidus au code",
    value: String(assidus.value.length),
    tone: "text-foreground",
    caption: "Prêts pour la conduite (≥ 80 % du code)",
  },
]);

const CHECKLIST = [
  {
    title: "Résilience des données et performance",
    items: [
      "Base SQLite locale en mode WAL (tolérance aux coupures brutales)",
      "Empreinte RAM minimale (< 60 Mo), adaptée aux PC d'agence modestes",
      "Identifiants UUID v4 (aucun conflit de clés)",
      "File d'événements prête pour la phase 2 (journal_sync_events)",
    ],
  },
  {
    title: "Ergonomie et rigueur comptable",
    items: [
      "Saisie 100 % clavier, raccourcis F1 à F5 en moins de 20 secondes",
      "Ticket thermique 80 mm ESC/POS et reçu A5 administratif",
      "Montant en toutes lettres conforme aux règles OHADA",
      "Détection des réimpressions (marquage DUPLICATA officiel)",
    ],
  },
];

const backToGuichet = () => {
  store.activeView = "guichet";
};

const sendWhatsAppRelance = (cand: any) => {
  const phone = cand.telephone_whatsapp.replace(/[^0-9]/g, "");
  const msg =
    `Bonjour ${cand.candidat_nom},\n\n` +
    `L'établissement ${authStore.schoolName} vous informe que le solde restant sur votre dossier *${cand.numero_dossier}* (${cand.formule}) est de *${fcfa(cand.solde_restant)} FCFA*.\n\n` +
    `Pour planifier vos leçons pratiques de conduite ou votre présentation à l'examen, nous vous invitons à régulariser votre tranche au guichet de l'agence ou par Wave/Orange Money.\n\n` +
    `Bien cordialement,\n` +
    `Direction ${authStore.schoolName}${authStore.schoolPhone ? " • Tél : " + authStore.schoolPhone : ""}`;

  window.open(
    `https://api.whatsapp.com/send?phone=${phone}&text=${encodeURIComponent(msg)}`,
    "_blank",
    "noopener,noreferrer",
  );
  store.showToast(
    `Relance WhatsApp préparée pour ${cand.candidat_nom}`,
    "success",
  );
};

const openOverride = (cand: any) => {
  candToOverride.value = cand;
  isOverrideModalOpen.value = true;
};

const submitOverride = async ({
  pin,
  motif,
}: {
  pin: string;
  motif: string;
}) => {
  if (!candToOverride.value || !pin) return;
  const success = await store.unlockConduiteWithPin(
    pin,
    candToOverride.value.dossier_id,
    motif || "Dérogation accordée depuis la gestion Dérogations Manager",
  );
  if (success) {
    isOverrideModalOpen.value = false;
    candToOverride.value = null;
    await store.loadGerantDashboard();
  }
};
</script>

<template>
  <!-- Skeleton Loader -->
  <GerantDashboardSkeleton
    v-if="isInitialLoading && !store.gerantDashboardData"
  />

  <div v-else class="flex flex-1 flex-col gap-4 overflow-y-auto">
    <header
      class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
    >
      <div class="flex min-w-0 items-center gap-3">
        <button
          type="button"
          @click="backToGuichet"
          class="inline-flex h-8 shrink-0 items-center gap-1.5 rounded-[4px] border border-border bg-card px-3 text-xs font-medium text-foreground transition hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-dark/40"
        >
          <ArrowLeft class="h-3.5 w-3.5 text-muted-foreground" />
          Guichet
          <kbd
            class="rounded bg-muted px-1 font-mono text-[10px] text-muted-foreground"
            >Échap</kbd
          >
        </button>

        <div class="min-w-0">
          <div class="flex items-center gap-2">
            <h1 class="truncate text-base font-semibold text-foreground">
              {{
                authStore.activeSection === "derogations"
                  ? "Dérogations Manager (PIN) & Déblocages"
                  : "Tableau de bord du gérant"
              }}
            </h1>
            <span
              class="inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-bold"
              :class="
                isWsConnected
                  ? 'bg-emerald-500/15 text-emerald-700 dark:text-emerald-300'
                  : 'bg-amber-500/15 text-amber-700'
              "
            >
              <Radio class="h-3 w-3 animate-pulse" />
              {{ isWsConnected ? "Temps réel WS" : "Reconnexion WS..." }}
            </span>
          </div>
          <p class="truncate text-xs text-muted-foreground">
            {{
              authStore.activeSection === "derogations"
                ? "Débloquez exceptionnellement l'accès aux leçons pratiques pour les élèves en retard de paiement."
                : "Impayés, assiduité au code et conformité hors ligne."
            }}
          </p>
        </div>
      </div>
    </header>

    <section
      class="grid grid-cols-1 divide-y divide-border rounded-[5px] border border-border bg-card sm:grid-cols-2 sm:divide-x sm:divide-y-0"
    >
      <div v-for="s in stats" :key="s.key" class="p-4">
        <span class="text-xs text-muted-foreground">{{ s.label }}</span>
        <p class="mt-1 text-lg font-semibold tabular-nums" :class="s.tone">
          {{ s.value }}
        </p>
        <p class="mt-1.5 text-[11px] text-muted-foreground">{{ s.caption }}</p>
      </div>
    </section>

    <nav
      class="flex gap-1 border-b border-border"
      role="tablist"
      aria-label="Sections du tableau de bord"
    >
      <button
        v-for="t in tabs"
        :key="t.key"
        type="button"
        role="tab"
        :aria-selected="activeTab === t.key"
        @click="activeTab = t.key"
        class="-mb-px inline-flex h-9 items-center gap-2 border-b-2 px-3 text-xs font-medium transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-dark/40"
        :class="
          activeTab === t.key
            ? 'border-brand-dark text-white'
            : 'border-transparent text-muted-foreground hover:text-foreground'
        "
      >
        {{ t.label }}
        <span
          v-if="t.count !== undefined"
          class="rounded-full bg-muted px-1.5 text-[10px] tabular-nums text-muted-foreground"
        >
          {{ t.count }}
        </span>
      </button>
    </nav>

    <section
      class="flex flex-1 flex-col overflow-hidden rounded-[5px] border border-border bg-card"
      role="tabpanel"
    >
      <GerantImpayesTab
        v-if="activeTab === 'impayes'"
        :impayes="impayes"
        @relance="sendWhatsAppRelance"
      />

      <GerantAssidusTab
        v-else-if="activeTab === 'assidus'"
        :assidus="assidus"
        :derogations-accorde="store.derogationsAccorde"
        @override="openOverride"
      />

      <div v-else class="flex-1 space-y-5 overflow-y-auto p-4">
        <div class="flex items-start gap-3">
          <div
            class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-brand-light text-brand-dark"
          >
            <CheckCircle2 class="h-4 w-4" />
          </div>
          <div class="min-w-0">
            <h2 class="text-sm font-semibold text-foreground">
              Jalon 1 atteint · MVP guichet d'agence validé
            </h2>
            <p class="mt-1 text-xs leading-relaxed text-muted-foreground">
              Le logiciel desktop est 100 % autonome. Inscriptions, photos
              webcam, encaissements par tranche, reçus thermiques 80 mm et
              présences de code fonctionnent sans Internet, même lors de
              coupures d'électricité prolongées.
            </p>
          </div>
        </div>

        <div
          class="grid grid-cols-1 gap-x-8 gap-y-5 border-t border-border pt-4 md:grid-cols-2"
        >
          <div v-for="g in CHECKLIST" :key="g.title">
            <h3 class="mb-2 text-xs font-semibold text-foreground">
              {{ g.title }}
            </h3>
            <ul class="space-y-1.5">
              <li
                v-for="item in g.items"
                :key="item"
                class="flex items-start gap-2 text-xs text-muted-foreground"
              >
                <Check class="mt-0.5 h-3.5 w-3.5 shrink-0 text-brand-dark" />
                <span>{{ item }}</span>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </section>

    <GerantOverridePinModal
      :is-open="isOverrideModalOpen"
      :cand="candToOverride"
      @close="isOverrideModalOpen = false"
      @submit="submitOverride"
    />
  </div>
</template>

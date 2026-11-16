<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { useAuthStore } from "../../stores/authStore";
import { useGuichetStore } from "../../stores/guichetStore";
import { useLanguageStore } from "../../stores/languageStore";
import { realtimeSocket } from "../../services/realtimeSocketService";
import ManagerLicenceBanner from "./home/ManagerLicenceBanner.vue";
import ManagerChartsSection from "./home/ManagerChartsSection.vue";
import ManagerRecentTransactions from "./home/ManagerRecentTransactions.vue";
import {
  AlertTriangle,
  ArrowUpRight,
  Award,
  Car,
  FileText,
  Users,
  RefreshCw,
  Calendar,
} from "lucide-vue-next";

const authStore = useAuthStore();
const guichetStore = useGuichetStore();
const languageStore = useLanguageStore();

const liveLivretActivity = ref<{
  candidatNom: string;
  kmParcourus: number;
  commentaire: string;
  updatedAt: string;
} | null>(null);
let unsubscribeHomeWs: (() => void) | null = null;

const t = (txt: string) => languageStore.t(txt);
const fmt = (n: number) => (n || 0).toLocaleString("fr-FR");

const SEUIL_ALERTE_JOURS = 14;
const JOUR_MS = 86_400_000;

const alertesFlotte = computed(() => {
  const aujourdhui = new Date();
  aujourdhui.setHours(0, 0, 0, 0);

  return guichetStore.vehicules.filter((v) => {
    const jVisite = Math.ceil(
      (new Date(v.date_visite_technique).getTime() - aujourdhui.getTime()) /
        JOUR_MS,
    );
    const jAssur = Math.ceil(
      (new Date(v.date_expiration_assurance).getTime() - aujourdhui.getTime()) /
        JOUR_MS,
    );
    return (
      jVisite <= SEUIL_ALERTE_JOURS ||
      jAssur <= SEUIL_ALERTE_JOURS ||
      v.statut === "EN_PANNE"
    );
  }).length;
});

const totalRecettesAffiche = computed(
  () =>
    guichetStore.gerantDashboardData?.ventilation?.total_encaisse ??
    (guichetStore.stats.total_especes_theorique || 0),
);

const totalDossiersFormules = computed(() => {
  const stats = guichetStore.gerantDashboardData?.inscriptions_formules;
  if (stats && stats.length > 0) {
    return stats.reduce((acc, curr) => acc + curr.inscrits, 0);
  }
  return (
    guichetStore.gerantDashboardData?.total_candidats_actifs ??
    (guichetStore.stats.total_candidats || 0)
  );
});

const indicateurs = computed(() => [
  {
    label: t("Chiffre d'affaires Caisse"),
    periode: t("Total encaissé"),
    valeur: `${fmt(totalRecettesAffiche.value)} FCFA`,
    progression: t("Actif"),
    icone: ArrowUpRight,
    ton: "emerald",
  },
  {
    label: t("Créances à recouvrer"),
    periode: `${guichetStore.gerantDashboardData?.nombre_dossiers_impayes || 0} ${t("dossiers")}`,
    valeur: `${fmt(guichetStore.gerantDashboardData?.total_creances_impayes || 0)} FCFA`,
    progression: t("En cours"),
    icone: AlertTriangle,
    ton: "amber",
  },
  {
    label: t("Taux de réussite Permis"),
    periode: t("Sessions d'examen"),
    valeur: guichetStore.gerantDashboardData?.taux_reussite_examen
      ? `${guichetStore.gerantDashboardData.taux_reussite_examen} %`
      : t("En attente"),
    progression: t("Officiel"),
    icone: Award,
    ton: "emerald",
  },
  {
    label: t("Candidats inscrits"),
    periode: t("Base active"),
    valeur: `${guichetStore.gerantDashboardData?.total_candidats_actifs ?? guichetStore.stats.total_candidats} ${t("élèves")}`,
    progression: `${guichetStore.searchResults.length} ${t("chargés")}`,
    icone: Users,
    ton: "sky",
  },
]);

const tonIcone: Record<string, string> = {
  emerald: "bg-transparent text-brand-dark",
  amber: "bg-amber-500/10 text-amber-600 dark:text-amber-400",
  sky: "bg-sky-500/10 text-sky-600 dark:text-sky-400",
};

onMounted(async () => {
  await Promise.all([
    guichetStore.loadGerantDashboard(),
    guichetStore.loadVehicules(),
    guichetStore.loadFormules(),
    authStore.loadProfileByRole(),
  ]);
  unsubscribeHomeWs = realtimeSocket.subscribe((msg) => {
    if (msg.type === "LIVRET_EMARGEMENT_LIVE" && msg.data) {
      liveLivretActivity.value = {
        candidatNom: msg.data.candidatNom || "Élève en conduite",
        kmParcourus: Math.max(
          0,
          (msg.data.kmFin ?? 0) - (msg.data.kmDebut ?? 0),
        ),
        commentaire: msg.data.commentaire || "",
        updatedAt: new Date().toLocaleTimeString("fr-FR"),
      };
    } else if (msg.type === "LIVRET_EMARGEMENT_UPDATED") {
      liveLivretActivity.value = null;
    }
  });
});

onUnmounted(() => {
  if (unsubscribeHomeWs) {
    unsubscribeHomeWs();
    unsubscribeHomeWs = null;
  }
});
</script>

<template>
  <div class="space-y-4">
    <!-- En-tête de la vue -->
    <div
      class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between"
    >
      <div>
        <h1 class="text-sm font-semibold leading-tight text-foreground">
          {{ t("Tableau de bord Direction") }}
        </h1>
        <p class="text-xs text-muted-foreground">
          {{
            t(
              "Analyses des recettes, effectifs candidats et suivi de la flotte automobile en temps réel",
            )
          }}
        </p>
      </div>

      <div class="flex flex-wrap items-center gap-1.5 sm:gap-2">
        <button
          type="button"
          @click="
            guichetStore.loadGerantDashboard();
            authStore.loadProfileByRole();
          "
          class="min-h-[38px] sm:min-h-0 flex items-center gap-1.5 rounded-[3px] border border-border bg-card px-2.5 py-1.5 text-xs font-medium text-muted-foreground hover:bg-muted hover:text-foreground transition"
          :title="t('Actualiser les données')"
        >
          <RefreshCw
            class="h-3.5 w-3.5"
            :class="{ 'animate-spin': guichetStore.loading }"
          />
          <span class="hidden sm:inline">{{ t("Actualiser") }}</span>
        </button>

        <button
          type="button"
          @click="authStore.setActiveSection('conduite-planning')"
          class="min-h-[38px] sm:min-h-0 flex items-center gap-1.5 rounded-[3px] border border-border bg-card px-3 py-1.5 text-xs font-medium text-foreground hover:bg-muted transition"
        >
          <Calendar class="h-3.5 w-3.5 text-brand" />
          {{ t("Planning Conduite") }}
        </button>

        <button
          type="button"
          @click="authStore.setActiveSection('finance')"
          class="min-h-[38px] sm:min-h-0 flex items-center gap-1.5 rounded-[3px] border border-border bg-card px-3 py-1.5 text-xs font-medium text-foreground hover:bg-muted transition"
        >
          <FileText class="h-3.5 w-3.5 text-muted-foreground" />
          {{ t("Audit financier") }}
        </button>

        <button
          type="button"
          @click="authStore.setActiveSection('flotte')"
          class="min-h-[38px] sm:min-h-0 flex items-center gap-1.5 rounded-[3px] bg-brand-dark px-3 py-1.5 text-xs font-medium text-white hover:bg-brand transition"
        >
          <Car class="h-3.5 w-3.5" />
          {{ t("Flotte") }}
          <span
            v-if="alertesFlotte"
            class="rounded-full bg-rose-500 px-1.5 py-0.2 text-[10px] font-bold text-white"
            >{{ alertesFlotte }}</span
          >
        </button>
      </div>
    </div>

    <!-- Bandeaux Licence & Activité Live -->
    <ManagerLicenceBanner :live-livret-activity="liveLivretActivity" />

    <!-- 4 Cartes KPI (Carrousel horizontal tactile sur mobile) -->
    <div
      class="flex overflow-x-auto snap-x snap-mandatory gap-2.5 no-scrollbar sm:grid sm:grid-cols-2 lg:grid-cols-4 sm:gap-3"
    >
      <div
        v-for="i in indicateurs"
        :key="i.label"
        class="min-w-[74%] sm:min-w-0 snap-start shrink-0 rounded-[3px] border border-border bg-card px-4 py-3 transition hover:border-brand/40"
      >
        <div class="flex items-center justify-between">
          <span
            class="flex h-7 w-7 items-center justify-center rounded-lg"
            :class="tonIcone[i.ton]"
          >
            <component :is="i.icone" class="h-4 w-4" />
          </span>
          <span
            class="rounded bg-muted px-1.5 py-0.5 text-[10px] font-medium text-muted-foreground"
            >{{ i.periode }}</span
          >
        </div>
        <div class="mt-2 text-xs text-muted-foreground">{{ i.label }}</div>
        <div class="flex items-baseline justify-between gap-1 mt-0.5">
          <span
            class="text-base font-bold leading-tight text-foreground tabular-nums"
            >{{ i.valeur }}</span
          >
          <span class="text-[10px] font-semibold text-brand">{{
            i.progression
          }}</span>
        </div>
      </div>
    </div>

    <!-- Graphiques Chart.js -->
    <ManagerChartsSection
      :total-recettes-affiche="totalRecettesAffiche"
      :total-dossiers-formules="totalDossiersFormules"
    />

    <!-- Encaissements récents -->
    <ManagerRecentTransactions />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, nextTick } from "vue";
import { useAuthStore } from "../../../stores/authStore";
import { useGuichetStore } from "../../../stores/guichetStore";
import { useLanguageStore } from "../../../stores/languageStore";
import { TrendingUp, PieChart, BarChart3 } from "lucide-vue-next";
import Chart from "chart.js/auto";

const props = defineProps<{
  totalRecettesAffiche: number;
  totalDossiersFormules: number;
}>();

const authStore = useAuthStore();
const guichetStore = useGuichetStore();
const languageStore = useLanguageStore();

const t = (txt: string) => languageStore.t(txt);
const fmt = (n: number) => (n || 0).toLocaleString("fr-FR");

const canvasRevenue = ref<HTMLCanvasElement | null>(null);
const canvasFormules = ref<HTMLCanvasElement | null>(null);
const canvasPaymentModes = ref<HTMLCanvasElement | null>(null);

let chartRevenueInstance: Chart | null = null;
let chartFormulesInstance: Chart | null = null;
let chartPaymentModesInstance: Chart | null = null;

const getThemeColors = () => {
  const isDark = authStore.isDark;
  return {
    textColor: isDark ? "#f8fafc" : "#0f172a",
    mutedText: isDark ? "#94a3b8" : "#64748b",
    gridColor: isDark ? "rgba(255, 255, 255, 0.06)" : "rgba(0, 0, 0, 0.05)",
    tooltipBg: isDark ? "#172033" : "#ffffff",
    tooltipBorder: isDark ? "#283548" : "#e2e8f0",
    tooltipText: isDark ? "#f8fafc" : "#0f172a",
  };
};

const renderRevenueChart = () => {
  if (!canvasRevenue.value) return;
  if (chartRevenueInstance) chartRevenueInstance.destroy();

  const ctx = canvasRevenue.value.getContext("2d");
  if (!ctx) return;

  const theme = getThemeColors();
  const gradient = ctx.createLinearGradient(0, 0, 0, 260);
  gradient.addColorStop(0, "rgba(16, 185, 129, 0.35)");
  gradient.addColorStop(1, "rgba(16, 185, 129, 0.0)");

  const rawMois = guichetStore.gerantDashboardData?.recettes_mensuelles || [];
  const moisNoms = [
    "Janv",
    "Févr",
    "Mars",
    "Avr",
    "Mai",
    "Juin",
    "Juil",
    "Août",
    "Sept",
    "Oct",
    "Nov",
    "Déc",
  ];

  let labels: string[] = [];
  let dataActual: number[] = [];

  if (rawMois.length > 0) {
    labels = rawMois.map((r) => {
      const parts = r.mois.split("-");
      if (parts.length === 2) {
        const idx = parseInt(parts[1], 10) - 1;
        return `${moisNoms[idx] || parts[1]} ${parts[0]}`;
      }
      return r.mois;
    });
    dataActual = rawMois.map((r) => r.montant);
  } else {
    const d = new Date();
    labels = [`${moisNoms[d.getMonth()]} ${d.getFullYear()}`];
    dataActual = [props.totalRecettesAffiche];
  }

  chartRevenueInstance = new Chart(ctx, {
    type: "line",
    data: {
      labels,
      datasets: [
        {
          label: t("Recettes réelles encaissées"),
          data: dataActual,
          borderColor: "#10b981",
          borderWidth: 2.5,
          backgroundColor: gradient,
          fill: true,
          tension: 0.38,
          pointBackgroundColor: "#10b981",
          pointBorderColor: "#ffffff",
          pointBorderWidth: 1.5,
          pointRadius: 4,
          pointHoverRadius: 6,
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      interaction: { mode: "index", intersect: false },
      plugins: {
        legend: { display: false },
        tooltip: {
          backgroundColor: theme.tooltipBg,
          borderColor: theme.tooltipBorder,
          borderWidth: 1,
          titleColor: theme.tooltipText,
          bodyColor: theme.mutedText,
          padding: 10,
          callbacks: {
            label: (ctx) => ` ${t("Recettes")} : ${fmt(Number(ctx.raw))} FCFA`,
          },
        },
      },
      scales: {
        x: {
          grid: { display: false },
          ticks: {
            color: theme.mutedText,
            font: { family: "Inter", size: 10 },
          },
          border: { display: false },
        },
        y: {
          grid: { color: theme.gridColor, drawTicks: false },
          ticks: {
            color: theme.mutedText,
            font: { family: "Inter", size: 10 },
            callback: (val) => `${fmt(Number(val))} F`,
          },
          border: { display: false },
        },
      },
    },
  });
};

const renderFormulesChart = () => {
  if (!canvasFormules.value) return;
  if (chartFormulesInstance) chartFormulesInstance.destroy();

  const ctx = canvasFormules.value.getContext("2d");
  if (!ctx) return;

  const theme = getThemeColors();
  const stats = guichetStore.gerantDashboardData?.inscriptions_formules || [];
  let labels: string[] = [];
  let dataCounts: number[] = [];

  if (stats && stats.length > 0) {
    labels = stats.map((f) => f.intitule);
    dataCounts = stats.map((f) => f.inscrits);
  } else if (guichetStore.formules && guichetStore.formules.length > 0) {
    labels = guichetStore.formules.map((f) => f.intitule);
    dataCounts = labels.map(() => 0);
  } else {
    labels = ["Permis B Standard (Véhicule Léger)"];
    dataCounts = [0];
  }

  chartFormulesInstance = new Chart(ctx, {
    type: "bar",
    data: {
      labels,
      datasets: [
        {
          label: t("Candidats inscrits"),
          data: dataCounts,
          backgroundColor: [
            "rgba(16, 185, 129, 0.85)",
            "rgba(14, 165, 233, 0.85)",
            "rgba(245, 158, 11, 0.85)",
            "rgba(139, 92, 246, 0.85)",
            "rgba(100, 116, 139, 0.85)",
          ],
          borderColor: ["#10b981", "#0ea5e9", "#f59e0b", "#8b5cf6", "#64748b"],
          borderWidth: 1,
          borderRadius: 6,
          borderSkipped: false,
          barThickness: 24,
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: { display: false },
        tooltip: {
          backgroundColor: theme.tooltipBg,
          borderColor: theme.tooltipBorder,
          borderWidth: 1,
          titleColor: theme.tooltipText,
          bodyColor: theme.mutedText,
          padding: 10,
          callbacks: {
            label: (ctx) => ` ${ctx.raw} ${t("candidats enrôlés")}`,
          },
        },
      },
      scales: {
        x: {
          grid: { display: false },
          ticks: {
            color: theme.mutedText,
            font: { family: "Inter", size: 10 },
          },
          border: { display: false },
        },
        y: {
          grid: { color: theme.gridColor, drawTicks: false },
          ticks: {
            color: theme.mutedText,
            stepSize: 1,
            font: { family: "Inter", size: 10 },
          },
          border: { display: false },
        },
      },
    },
  });
};

const renderPaymentModesChart = () => {
  if (!canvasPaymentModes.value) return;
  if (chartPaymentModesInstance) chartPaymentModesInstance.destroy();

  const ctx = canvasPaymentModes.value.getContext("2d");
  if (!ctx) return;

  const theme = getThemeColors();
  const vent = guichetStore.gerantDashboardData?.ventilation;

  const modesData = [
    {
      label: t("Espèces"),
      val: vent
        ? vent.especes
        : guichetStore.stats.total_especes_theorique || 0,
      color: "#10b981",
    },
    { label: "Wave", val: vent ? vent.wave : 0, color: "#0ea5e9" },
    {
      label: "Orange Money",
      val: vent ? vent.orange_money : 0,
      color: "#f97316",
    },
    { label: "MTN MoMo", val: vent ? vent.mtn_momo : 0, color: "#eab308" },
    { label: "Moov Money", val: vent ? vent.moov : 0, color: "#2563eb" },
    { label: t("Virement"), val: vent ? vent.virement : 0, color: "#9333ea" },
  ].filter((m) => m.val > 0);

  const hasData = modesData.length > 0;
  const labels = hasData
    ? modesData.map((m) => m.label)
    : [t("Aucun paiement")];
  const data = hasData ? modesData.map((m) => m.val) : [1];
  const colors = hasData
    ? modesData.map((m) => m.color)
    : [authStore.isDark ? "#334155" : "#e2e8f0"];

  chartPaymentModesInstance = new Chart(ctx, {
    type: "doughnut",
    data: {
      labels,
      datasets: [
        {
          data,
          backgroundColor: colors,
          borderColor: theme.tooltipBg,
          borderWidth: 2,
          hoverOffset: 6,
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      cutout: "72%",
      plugins: {
        legend: {
          position: "bottom",
          labels: {
            boxWidth: 10,
            boxHeight: 10,
            usePointStyle: true,
            pointStyle: "circle",
            color: theme.mutedText,
            font: { family: "Inter", size: 10, weight: 500 },
            padding: 12,
          },
        },
        tooltip: {
          backgroundColor: theme.tooltipBg,
          borderColor: theme.tooltipBorder,
          borderWidth: 1,
          titleColor: theme.tooltipText,
          bodyColor: theme.mutedText,
          padding: 10,
          callbacks: {
            label: (ctx) =>
              hasData
                ? ` ${ctx.label} : ${fmt(Number(ctx.raw))} FCFA`
                : " 0 FCFA",
          },
        },
      },
    },
  });
};

const refreshAllCharts = async () => {
  await nextTick();
  renderRevenueChart();
  renderFormulesChart();
  renderPaymentModesChart();
};

onMounted(() => {
  refreshAllCharts();
});

onUnmounted(() => {
  if (chartRevenueInstance) chartRevenueInstance.destroy();
  if (chartFormulesInstance) chartFormulesInstance.destroy();
  if (chartPaymentModesInstance) chartPaymentModesInstance.destroy();
});

watch(
  () => [
    authStore.isDark,
    languageStore.currentLanguage,
    languageStore.translationsRevision,
    guichetStore.gerantDashboardData,
  ],
  () => {
    refreshAllCharts();
  },
);
</script>

<template>
  <div class="grid grid-cols-1 gap-3 sm:gap-4 lg:grid-cols-3">
    <!-- Graphique 1 : Évolution des Recettes Réelles -->
    <div
      class="rounded-[3px] border border-border bg-card p-3.5 sm:p-4 lg:col-span-2 flex flex-col justify-between"
    >
      <div
        class="flex flex-wrap items-center justify-between gap-2 pb-3 border-b border-border"
      >
        <div class="flex items-center gap-2">
          <span
            class="flex h-6 w-6 items-center justify-center rounded-[3px] bg-brand/10 text-brand-dark"
          >
            <TrendingUp class="h-3.5 w-3.5" />
          </span>
          <div>
            <h2 class="text-xs font-semibold text-foreground">
              {{ t("Évolution des Recettes Réelles") }}
            </h2>
            <p class="text-[10px] text-muted-foreground">
              {{ t("Total des encaissements par période mensuelle") }}
            </p>
          </div>
        </div>
        <span
          class="text-xs font-semibold text-emerald-600 dark:text-emerald-400 tabular-nums"
        >
          {{ t("Total") }} : {{ fmt(totalRecettesAffiche) }} FCFA
        </span>
      </div>

      <div class="relative h-52 sm:h-64 w-full pt-3">
        <canvas ref="canvasRevenue"></canvas>
      </div>
    </div>

    <!-- Graphique 2 : Canaux de Règlement -->
    <div
      class="rounded-[3px] border border-border bg-card p-3.5 sm:p-4 flex flex-col justify-between"
    >
      <div class="flex items-center gap-2 pb-3 border-b border-border">
        <span
          class="flex h-6 w-6 items-center justify-center rounded-[3px] bg-sky-500/10 text-sky-600 dark:text-sky-400"
        >
          <PieChart class="h-3.5 w-3.5" />
        </span>
        <div>
          <h2 class="text-xs font-semibold text-foreground">
            {{ t("Canaux de Règlement") }}
          </h2>
          <p class="text-[10px] text-muted-foreground">
            {{ t("Mobile Money vs Espèces") }}
          </p>
        </div>
      </div>

      <div
        class="relative h-52 sm:h-60 w-full pt-2 flex items-center justify-center"
      >
        <canvas ref="canvasPaymentModes"></canvas>
        <div
          class="pointer-events-none absolute inset-0 flex flex-col items-center justify-center text-center pb-8"
        >
          <span
            class="text-[10px] font-medium text-muted-foreground uppercase tracking-wide"
            >{{ t("Recettes") }}</span
          >
          <span class="text-sm font-extrabold text-foreground tabular-nums"
            >{{ fmt(totalRecettesAffiche) }} F</span
          >
          <span class="text-[9px] text-brand font-medium">{{
            t("100% réel")
          }}</span>
        </div>
      </div>
    </div>

    <!-- Graphique 3 : Inscriptions par Formule de Permis -->
    <div
      class="rounded-[3px] border border-border bg-card p-3.5 sm:p-4 lg:col-span-3"
    >
      <div
        class="flex flex-wrap items-center justify-between gap-2 pb-3 border-b border-border"
      >
        <div class="flex items-center gap-2">
          <span
            class="flex h-6 w-6 items-center justify-center rounded-[3px] bg-amber-500/10 text-amber-600 dark:text-amber-400"
          >
            <BarChart3 class="h-3.5 w-3.5" />
          </span>
          <div>
            <h2 class="text-xs font-semibold text-foreground">
              {{ t("Inscriptions Réelles par Formule de Permis") }}
            </h2>
            <p class="text-[10px] text-muted-foreground">
              {{ t("Nombre de dossiers enrôlés par formule") }}
            </p>
          </div>
        </div>
        <span class="text-[11px] font-semibold text-foreground tabular-nums">
          {{ t("Total") }} : {{ totalDossiersFormules }}
          {{ t("dossier(s) souscrit(s)") }}
        </span>
      </div>

      <div class="relative h-48 sm:h-56 w-full pt-3">
        <canvas ref="canvasFormules"></canvas>
      </div>
    </div>
  </div>
</template>

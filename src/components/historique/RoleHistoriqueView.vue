<script lang="ts">
export type { HistoriqueItem } from "./types";
</script>

<script setup lang="ts">
import { ref, computed } from "vue";
import {
  History,
  RefreshCw,
  Printer,
  Download,
  Loader2,
} from "lucide-vue-next";
import { useAuthStore } from "../../stores/authStore";
import { useGuichetStore } from "../../stores/guichetStore";
import { useRoleHistorique } from "./useRoleHistorique";
import { exportHistoriqueToPdf } from "../../utils/pdfExport";
import type { HistoriqueItem } from "./types";
import HistoriqueKpiCards from "./HistoriqueKpiCards.vue";
import HistoriqueTimelineList from "./HistoriqueTimelineList.vue";

const {
  isWsConnected,
  lastSyncAt,
  isRefreshing,
  currentRole,
  roleTitle,
  quizCount,
  allHistoryItems,
  kpiTotalActions,
  kpiTotalMontant,
  kpiTotalConduite,
  kpiTotalCode,
  refreshAllData,
} = useRoleHistorique();

const searchQuery = ref("");
const selectedFilter = ref<string>("ALL");

const filteredItems = computed<HistoriqueItem[]>(() => {
  return allHistoryItems.value.filter((item) => {
    if (selectedFilter.value !== "ALL" && item.type !== selectedFilter.value) {
      return false;
    }
    if (searchQuery.value.trim()) {
      const q = searchQuery.value.toLowerCase().trim();
      const matchTitre = item.titre.toLowerCase().includes(q);
      const matchDesc = item.description.toLowerCase().includes(q);
      const matchActeur = (item.acteur || "").toLowerCase().includes(q);
      const matchBadge = item.badgeLabel.toLowerCase().includes(q);
      return matchTitre || matchDesc || matchActeur || matchBadge;
    }
    return true;
  });
});

const authStore = useAuthStore();
const guichetStore = useGuichetStore();
const isExporting = ref(false);

const printHistory = () => {
  window.print();
};

const exportPdfHistory = async () => {
  if (isExporting.value || filteredItems.value.length === 0) return;
  isExporting.value = true;
  try {
    const res = await exportHistoriqueToPdf({
      schoolName: authStore.currentSchool?.nom_commercial || "AutoSchool",
      roleTitle: roleTitle.value,
      totalActions: kpiTotalActions.value,
      totalMontant: kpiTotalMontant.value,
      items: filteredItems.value,
    });
    guichetStore.showToast(
      res.savedPath
        ? `PDF enregistré dans Téléchargements : ${res.filename}`
        : `PDF exporté : ${res.filename}`,
      "success",
    );
  } catch (err: any) {
    guichetStore.showToast(
      "Erreur export PDF : " + (err?.message || "Impossible de générer le PDF"),
      "error",
    );
  } finally {
    isExporting.value = false;
  }
};
</script>

<template>
  <div class="space-y-5 p-4 sm:p-6 max-w-7xl mx-auto">
    <!-- En-tête -->
    <div
      class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 border-b border-border pb-4"
    >
      <div>
        <div class="flex items-center gap-2">
          <div
            class="h-8 w-8 rounded-[4px] bg-brand/10 text-brand flex items-center justify-center font-bold"
          >
            <History class="h-4 w-4" />
          </div>
          <h1 class="text-lg font-bold text-foreground">
            Historique & Activités en Temps Réel
          </h1>
          <span
            class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-semibold border"
            :class="
              isWsConnected
                ? 'bg-emerald-500/10 text-emerald-700 border-emerald-500/30 dark:text-emerald-400'
                : 'bg-amber-500/10 text-amber-700 border-amber-500/30 dark:text-amber-400'
            "
          >
            <span
              class="h-2 w-2 rounded-full"
              :class="
                isWsConnected ? 'bg-emerald-500 animate-pulse' : 'bg-amber-500'
              "
            />
            {{ isWsConnected ? "Flux WebSocket Connecté" : "Mode Hors-Ligne" }}
          </span>
        </div>
        <p class="text-xs text-muted-foreground mt-1">
          Journal d'audit et suivi des opérations pour le rôle :
          <span class="font-semibold text-foreground">{{ roleTitle }}</span>
          <span v-if="lastSyncAt" class="ml-2">
            • Dernière mise à jour : {{ lastSyncAt }}
          </span>
        </p>
      </div>

      <div class="flex items-center gap-2 shrink-0">
        <button
          type="button"
          @click="refreshAllData"
          :disabled="isRefreshing"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-[4px] border border-border bg-card text-xs font-medium text-foreground hover:bg-muted transition active:scale-95 disabled:opacity-50"
          title="Actualiser les données"
        >
          <RefreshCw
            class="h-3.5 w-3.5"
            :class="{ 'animate-spin': isRefreshing }"
          />
          Actualiser
        </button>

        <button
          type="button"
          @click="exportPdfHistory"
          :disabled="isExporting || filteredItems.length === 0"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-[4px] bg-brand-dark text-xs font-medium text-white hover:bg-brand transition active:scale-95 disabled:opacity-50"
          title="Exporter l'historique en PDF"
        >
          <Loader2 v-if="isExporting" class="h-3.5 w-3.5 animate-spin" />
          <Download v-else class="h-3.5 w-3.5" />
          {{ isExporting ? "Export..." : "Exporter" }}
        </button>

        <button
          type="button"
          @click="printHistory"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-[4px] border border-border bg-card text-xs font-medium text-foreground hover:bg-muted transition active:scale-95"
          title="Imprimer l'historique"
        >
          <Printer class="h-3.5 w-3.5" />
          Imprimer
        </button>
      </div>
    </div>

    <!-- KPIs -->
    <HistoriqueKpiCards
      :current-role="currentRole"
      :kpi-total-actions="kpiTotalActions"
      :kpi-total-montant="kpiTotalMontant"
      :quiz-count="quizCount"
      :kpi-total-conduite="kpiTotalConduite"
      :kpi-total-code="kpiTotalCode"
    />

    <!-- Filtres & Flux -->
    <HistoriqueTimelineList
      :items="filteredItems"
      v-model:search-query="searchQuery"
      v-model:selected-filter="selectedFilter"
      :current-role="currentRole"
    />
  </div>
</template>

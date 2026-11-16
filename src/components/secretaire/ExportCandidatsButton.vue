<script setup lang="ts">
import { ref } from "vue";
import { Download, FileSpreadsheet, FileText, Loader2, ChevronDown } from "lucide-vue-next";
import type { Candidat } from "../../types";
import { useAuthStore } from "../../stores/authStore";
import { useGuichetStore } from "../../stores/guichetStore";
import {
  exportCandidatsToExcel,
  exportCandidatsToPdf,
} from "../../utils/candidatsExport";

const props = defineProps<{
  candidats: Candidat[];
}>();

const authStore = useAuthStore();
const guichetStore = useGuichetStore();

const isOpen = ref(false);
const isExporting = ref(false);

const runExport = async (format: "pdf" | "excel") => {
  isOpen.value = false;
  if (isExporting.value || props.candidats.length === 0) return;
  isExporting.value = true;
  const schoolName = authStore.currentSchool?.nom_commercial || "AutoSchool";
  try {
    const res =
      format === "pdf"
        ? await exportCandidatsToPdf(props.candidats, schoolName)
        : await exportCandidatsToExcel(props.candidats, schoolName);

    guichetStore.showToast(
      res.savedPath
        ? `Fichier enregistré dans Téléchargements : ${res.filename}`
        : `Export réussi : ${res.filename}`,
      "success",
    );
  } catch (err: any) {
    guichetStore.showToast(
      "Erreur export : " + (err?.message || "Échec de génération"),
      "error",
    );
  } finally {
    isExporting.value = false;
  }
};
</script>

<template>
  <div class="relative shrink-0">
    <button
      type="button"
      @click="isOpen = !isOpen"
      :disabled="isExporting || candidats.length === 0"
      class="inline-flex h-7 items-center gap-1.5 rounded-[4px] border border-border bg-card px-2.5 text-xs font-medium text-foreground transition hover:bg-muted disabled:opacity-50"
      title="Exporter les candidats et leurs paiements"
    >
      <Loader2 v-if="isExporting" class="h-3.5 w-3.5 animate-spin text-brand" />
      <Download v-else class="h-3.5 w-3.5 text-brand" />
      <span>{{ isExporting ? "Export..." : "Exporter" }}</span>
      <ChevronDown class="h-3 w-3 text-muted-foreground" />
    </button>

    <div
      v-if="isOpen"
      class="fixed inset-0 z-40"
      @click="isOpen = false"
    />

    <div
      v-if="isOpen"
      class="absolute right-0 mt-1.5 w-48 rounded-[4px] border border-border bg-card py-1 shadow-lg z-50"
    >
      <button
        type="button"
        @click="runExport('pdf')"
        class="flex w-full items-center gap-2 px-3 py-2 text-left text-xs text-foreground hover:bg-muted transition"
      >
        <FileText class="h-3.5 w-3.5 text-rose-600 shrink-0" />
        <span>Exporter en PDF (.pdf)</span>
      </button>
      <button
        type="button"
        @click="runExport('excel')"
        class="flex w-full items-center gap-2 px-3 py-2 text-left text-xs text-foreground hover:bg-muted transition"
      >
        <FileSpreadsheet class="h-3.5 w-3.5 text-emerald-600 shrink-0" />
        <span>Exporter en Excel (.xlsx)</span>
      </button>
    </div>
  </div>
</template>


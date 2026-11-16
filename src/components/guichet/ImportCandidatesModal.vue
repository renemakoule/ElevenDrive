<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { useGuichetStore } from "../../stores/guichetStore";
import {
  downloadCandidatesTemplate,
  inspectExcelFile,
  parseWithMappings,
} from "../../utils/excelCandidates";
import { buildInitialFormulaValueMap } from "../../utils/excelColumnMatcher";
import type {
  ExcelColumnMapping,
  ImporterCandidatsBatchResult,
  ParsedExcelCandidateRow,
  SystemImportField,
} from "../../types";
import ImportMappingPanel from "./enrollment/ImportMappingPanel.vue";
import ImportPreviewTable from "./enrollment/ImportPreviewTable.vue";
import ImportResultSummary from "./enrollment/ImportResultSummary.vue";
import { Download, FileSpreadsheet, Upload, X } from "lucide-vue-next";

const store = useGuichetStore();

const fileInputRef = ref<HTMLInputElement | null>(null);
const selectedFile = ref<File | null>(null);
const defaultFormuleId = ref("");
const ignorerDoublons = ref(true);
const impacterCaisseJour = ref(false);
const useDefaultIfUnknownFormule = ref(false);
const generateMissingPiece = ref(false);
const parsing = ref(false);
const parseError = ref("");

const excelHeaders = ref<string[]>([]);
const rawRows = ref<Record<string, unknown>[]>([]);
const columnMapping = ref<ExcelColumnMapping>({} as ExcelColumnMapping);
const formulaLabels = ref<string[]>([]);
const formulaValueMap = ref<Record<string, string>>({});
const includedFormulaLabels = ref<string[]>([]);
const importResult = ref<ImporterCandidatsBatchResult | null>(null);

watch(
  () => store.isImportCandidatesModalOpen,
  async (isOpen) => {
    if (isOpen) {
      if (!store.formules || store.formules.length === 0) {
        await store.loadFormules();
      }
      if (!defaultFormuleId.value && store.formules.length > 0) {
        defaultFormuleId.value = store.formules[0].id;
      }
    }
  },
);

watch(
  () => store.formules,
  (list) => {
    if (!defaultFormuleId.value && list.length > 0) {
      defaultFormuleId.value = list[0].id;
    }
  },
  { immediate: true },
);

const parsedRows = computed<ParsedExcelCandidateRow[]>(() => {
  if (rawRows.value.length === 0) return [];
  return parseWithMappings(
    rawRows.value,
    columnMapping.value,
    formulaValueMap.value,
    store.formules,
    {
      defaultFormuleId: defaultFormuleId.value,
      useDefaultIfUnknownFormule: useDefaultIfUnknownFormule.value,
      generateMissingPiece: generateMissingPiece.value,
    },
  );
});

const includedRows = computed<ParsedExcelCandidateRow[]>(() => {
  if (formulaLabels.value.length === 0) return parsedRows.value;
  const allowed = new Set(includedFormulaLabels.value);
  return parsedRows.value.filter(
    (r) => !r.rawFormuleLabel || allowed.has(r.rawFormuleLabel),
  );
});

const excludedCount = computed(
  () => parsedRows.value.length - includedRows.value.length,
);

const validRows = computed(() => includedRows.value.filter((r) => r.isValid));

const resetState = () => {
  selectedFile.value = null;
  parseError.value = "";
  excelHeaders.value = [];
  rawRows.value = [];
  formulaLabels.value = [];
  formulaValueMap.value = {};
  includedFormulaLabels.value = [];
  importResult.value = null;
  if (fileInputRef.value) fileInputRef.value.value = "";
};

const close = () => {
  store.isImportCandidatesModalOpen = false;
  resetState();
};

const refreshFormulaLabels = () => {
  const { labels, valueMap } = buildInitialFormulaValueMap(
    rawRows.value,
    columnMapping.value.formule,
    store.formules,
  );
  formulaLabels.value = labels;
  formulaValueMap.value = valueMap;
  includedFormulaLabels.value = [...labels];
};

const processFile = async (file: File) => {
  selectedFile.value = file;
  parseError.value = "";
  importResult.value = null;
  parsing.value = true;
  try {
    const inspected = await inspectExcelFile(file);
    excelHeaders.value = inspected.headers;
    rawRows.value = inspected.rawRows;
    columnMapping.value = inspected.autoMapping;
    refreshFormulaLabels();
  } catch (err: any) {
    rawRows.value = [];
    parseError.value = err.message || "Impossible de lire ce fichier Excel.";
  } finally {
    parsing.value = false;
  }
};

const onColumnChange = (field: SystemImportField, header: string) => {
  columnMapping.value = { ...columnMapping.value, [field]: header };
  if (field === "formule") {
    refreshFormulaLabels();
  }
};

const onFormulaValueChange = (label: string, formuleId: string) => {
  formulaValueMap.value = { ...formulaValueMap.value, [label]: formuleId };
};

const onFileChange = (e: Event) => {
  const file = (e.target as HTMLInputElement).files?.[0];
  if (file) processFile(file);
};

const onDropFile = (e: DragEvent) => {
  const file = e.dataTransfer?.files?.[0];
  if (file) processFile(file);
};

const submitImport = async () => {
  if (validRows.value.length === 0 || store.loading) return;
  importResult.value = await store.importCandidatesBatch({
    candidats: validRows.value.map((r) => r.payload),
    ignorer_doublons: ignorerDoublons.value,
    impacter_caisse_jour: impacterCaisseJour.value,
  });
};
</script>

<template>
  <div
    v-if="store.isImportCandidatesModalOpen"
    class="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/50 p-0 sm:p-4"
    role="dialog"
    aria-modal="true"
    aria-labelledby="import-excel-title"
    @click.self="close"
  >
    <div
      class="flex max-h-[92dvh] w-full sm:max-w-4xl flex-col overflow-hidden rounded-t-[3px] sm:rounded-[3px] border-t sm:border border-border bg-card pb-safe"
    >
      <header
        class="flex items-center justify-between gap-3 border-b border-border px-4 py-3"
      >
        <div class="flex items-center gap-2 min-w-0">
          <FileSpreadsheet class="h-4 w-4 text-brand shrink-0" />
          <div class="min-w-0">
            <h2
              id="import-excel-title"
              class="text-sm font-semibold text-foreground"
            >
              Importer des élèves depuis Excel
            </h2>
            <p class="text-[11px] text-muted-foreground">
              Téléversez n'importe quel fichier .xlsx, .xls ou .csv et associez vos colonnes ou formules si besoin.
            </p>
          </div>
        </div>
        <button
          type="button"
          @click="close"
          class="flex h-7 w-7 items-center justify-center rounded-[4px] text-muted-foreground transition hover:bg-muted hover:text-foreground"
          aria-label="Fermer"
        >
          <X class="h-4 w-4" />
        </button>
      </header>

      <div class="flex-1 space-y-3 overflow-y-auto p-4">
        <div
          class="flex flex-wrap items-center justify-between gap-2 rounded-[3px] border border-border bg-muted/40 p-3"
        >
          <div class="text-xs text-foreground">
            <p class="font-medium">Besoin du format attendu ?</p>
            <p class="text-[11px] text-muted-foreground">
              Le modèle inclut les colonnes requises et la liste de vos formules actuelles.
            </p>
          </div>
          <button
            type="button"
            @click="downloadCandidatesTemplate(store.formules)"
            class="inline-flex h-8 items-center gap-1.5 rounded-[4px] border border-border bg-card px-3 text-xs font-medium text-foreground transition hover:bg-muted"
          >
            <Download class="h-3.5 w-3.5 text-brand" />
            Télécharger le modèle (.xlsx)
          </button>
        </div>

        <template v-if="!importResult">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 text-xs">
            <div>
              <label class="mb-1 block text-[11px] font-medium text-muted-foreground">
                Formule par défaut (si colonne vide)
              </label>
              <select
                v-model="defaultFormuleId"
                class="h-8 w-full rounded-[4px] border border-border bg-background px-2 text-xs text-foreground"
              >
                <option v-for="f in store.formules" :key="f.id" :value="f.id">
                  {{ f.intitule }} ({{ f.categorie_permis }})
                </option>
              </select>
            </div>

            <div class="flex flex-col justify-center gap-1.5 pt-1">
              <label class="flex items-center gap-2 text-[11px] text-foreground cursor-pointer">
                <input
                  v-model="ignorerDoublons"
                  type="checkbox"
                  class="rounded border-border"
                />
                Ignorer les élèves déjà inscrits (même N° pièce)
              </label>
              <label class="flex items-center gap-2 text-[11px] text-foreground cursor-pointer">
                <input
                  v-model="impacterCaisseJour"
                  type="checkbox"
                  class="rounded border-border"
                />
                Comptabiliser les acomptes dans la caisse du jour
              </label>
            </div>
          </div>

          <div
            @dragover.prevent
            @drop.prevent="onDropFile"
            @click="fileInputRef?.click()"
            class="flex cursor-pointer flex-col items-center justify-center rounded-[4px] border border-dashed border-border bg-muted/20 p-4 text-center transition hover:border-brand hover:bg-brand-light/20"
          >
            <input
              ref="fileInputRef"
              type="file"
              accept=".xlsx,.xls,.csv"
              class="hidden"
              @change="onFileChange"
            />
            <Upload class="mb-1.5 h-5 w-5 text-brand" />
            <p class="text-xs font-medium text-foreground">
              {{
                selectedFile
                  ? selectedFile.name
                  : "Cliquez ou glissez-déposez votre fichier Excel ici"
              }}
            </p>
            <p class="mt-0.5 text-[11px] text-muted-foreground">
              Formats supportés : .xlsx, .xls, .csv (max 1000 lignes)
            </p>
          </div>

          <div
            v-if="parseError"
            class="rounded-[3px] border border-rose-200 bg-rose-50 p-2.5 text-xs text-rose-700"
          >
            {{ parseError }}
          </div>

          <div v-if="parsing" class="py-4 text-center text-xs text-muted-foreground">
            Analyse du fichier Excel en cours...
          </div>

          <template v-else-if="rawRows.length > 0">
            <ImportMappingPanel
              :headers="excelHeaders"
              :column-mapping="columnMapping"
              :formula-labels="formulaLabels"
              :formula-value-map="formulaValueMap"
              v-model:included-formula-labels="includedFormulaLabels"
              :formules="store.formules"
              v-model:use-default-if-unknown-formule="useDefaultIfUnknownFormule"
              v-model:generate-missing-piece="generateMissingPiece"
              @update-column="onColumnChange"
              @update-formula-value="onFormulaValueChange"
            />
            <ImportPreviewTable
              :rows="includedRows"
              :excluded-count="excludedCount"
            />
          </template>
        </template>

        <ImportResultSummary v-else :result="importResult" />
      </div>

      <footer
        class="flex items-center justify-between gap-2 border-t border-border px-4 py-3"
      >
        <button
          type="button"
          @click="importResult ? resetState() : close()"
          class="inline-flex h-8 items-center rounded-[4px] border border-border bg-card px-3 text-xs font-medium text-foreground transition hover:bg-muted"
        >
          {{ importResult ? "Importer un autre fichier" : "Fermer" }}
        </button>

        <button
          v-if="!importResult"
          type="button"
          :disabled="store.loading || validRows.length === 0"
          @click="submitImport"
          class="inline-flex h-8 items-center gap-2 rounded-[4px] bg-brand-dark px-3 text-xs font-medium text-white transition hover:bg-brand disabled:cursor-not-allowed disabled:opacity-50"
        >
          {{
            store.loading
              ? "Importation..."
              : `Inscrire ${validRows.length} élève(s) valide(s)`
          }}
        </button>
        <button
          v-else
          type="button"
          @click="close"
          class="inline-flex h-8 items-center rounded-[4px] bg-brand-dark px-3 text-xs font-medium text-white transition hover:bg-brand"
        >
          Terminer
        </button>
      </footer>
    </div>
  </div>
</template>

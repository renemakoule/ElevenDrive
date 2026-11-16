<script setup lang="ts">
import { computed, ref } from "vue";
import type {
  ExcelColumnMapping,
  FormulePermis,
  SystemImportField,
} from "../../../types";
import { SYSTEM_FIELD_LABELS } from "../../../utils/excelColumnMatcher";
import { ChevronDown, ChevronUp, SlidersHorizontal } from "lucide-vue-next";

const props = defineProps<{
  headers: string[];
  columnMapping: ExcelColumnMapping;
  formulaLabels: string[];
  formulaValueMap: Record<string, string>;
  includedFormulaLabels: string[];
  formules: FormulePermis[];
  useDefaultIfUnknownFormule: boolean;
  generateMissingPiece: boolean;
}>();

const emit = defineEmits<{
  (e: "update-column", field: SystemImportField, header: string): void;
  (e: "update-formula-value", label: string, formuleId: string): void;
  (e: "update:includedFormulaLabels", val: string[]): void;
  (e: "update:useDefaultIfUnknownFormule", val: boolean): void;
  (e: "update:generateMissingPiece", val: boolean): void;
}>();

const missingRequiredFields = computed(() =>
  SYSTEM_FIELD_LABELS.filter((def) => {
    if (!def.required) return false;
    if (
      (def.field === "nom" || def.field === "prenoms") &&
      props.columnMapping.nom_complet
    ) {
      return false;
    }
    return !props.columnMapping[def.field];
  }),
);

const showColumnMapper = ref(missingRequiredFields.value.length > 0);

const unmappedExcelHeaders = computed(() => {
  const mapped = new Set(Object.values(props.columnMapping).filter(Boolean));
  return props.headers.filter((h) => !mapped.has(h));
});

const isIncluded = (label: string) =>
  props.includedFormulaLabels.includes(label);

const toggleFormulaIncluded = (label: string) => {
  if (isIncluded(label)) {
    emit(
      "update:includedFormulaLabels",
      props.includedFormulaLabels.filter((l) => l !== label),
    );
  } else {
    emit("update:includedFormulaLabels", [
      ...props.includedFormulaLabels,
      label,
    ]);
  }
};

const selectAllFormulas = () => {
  emit("update:includedFormulaLabels", [...props.formulaLabels]);
};

const clearAllFormulas = () => {
  emit("update:includedFormulaLabels", []);
};

const unmappedSelectedCount = computed(
  () =>
    props.includedFormulaLabels.filter((lbl) => !props.formulaValueMap[lbl])
      .length,
);
</script>

<template>
  <div class="space-y-2.5 rounded-[3px] border border-border bg-muted/25 p-3">
    <!-- 1. Correspondance et filtre des formules détectées dans l'Excel -->
    <div v-if="formulaLabels.length > 0" class="space-y-1.5">
      <div class="flex flex-wrap items-center justify-between gap-2">
        <div class="flex items-center gap-2">
          <span class="text-[11px] font-semibold text-foreground">
            Formules à inscrire ({{ includedFormulaLabels.length }}/{{
              formulaLabels.length
            }})
          </span>
          <span
            v-if="unmappedSelectedCount > 0"
            class="text-[10px] font-medium text-amber-700"
          >
            {{ unmappedSelectedCount }} à associer
          </span>
        </div>

        <div class="flex items-center gap-2 text-[10px]">
          <button
            type="button"
            @click="selectAllFormulas"
            class="font-medium text-brand hover:underline"
          >
            Tout cocher
          </button>
          <span class="text-border">|</span>
          <button
            type="button"
            @click="clearAllFormulas"
            class="font-medium text-muted-foreground hover:text-foreground hover:underline"
          >
            Tout décocher
          </button>
        </div>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
        <div
          v-for="label in formulaLabels"
          :key="label"
          :class="[
            'flex items-center justify-between gap-2 rounded-[3px] border px-2 py-1.5 text-[11px] transition',
            isIncluded(label)
              ? 'border-border bg-card'
              : 'border-border/60 bg-muted/40 opacity-60',
          ]"
        >
          <label
            class="flex min-w-0 flex-1 cursor-pointer items-center gap-1.5"
          >
            <input
              type="checkbox"
              :checked="isIncluded(label)"
              @change="toggleFormulaIncluded(label)"
              class="rounded border-border shrink-0"
            />
            <span class="truncate font-medium text-foreground" :title="label">
              {{ label }}
            </span>
          </label>
          <select
            :value="formulaValueMap[label] || ''"
            :disabled="!isIncluded(label)"
            @change="
              emit(
                'update-formula-value',
                label,
                ($event.target as HTMLSelectElement).value,
              )
            "
            :class="[
              'h-7 max-w-[52%] rounded-[3px] border px-1.5 text-[11px] disabled:cursor-not-allowed',
              !isIncluded(label) || formulaValueMap[label]
                ? 'border-border bg-background text-foreground'
                : 'border-amber-400 bg-amber-50 text-amber-900',
            ]"
          >
            <option value="">-- Non associée --</option>
            <option v-for="f in formules" :key="f.id" :value="f.id">
              {{ f.intitule }} ({{ f.categorie_permis }})
            </option>
          </select>
        </div>
      </div>
    </div>

    <!-- 2. Options de tolérance sur les données manquantes -->
    <div class="flex flex-wrap items-center gap-x-4 gap-y-1.5 pt-1 text-[11px]">
      <label class="flex items-center gap-1.5 text-foreground cursor-pointer">
        <input
          :checked="useDefaultIfUnknownFormule"
          @change="
            emit(
              'update:useDefaultIfUnknownFormule',
              ($event.target as HTMLInputElement).checked,
            )
          "
          type="checkbox"
          class="rounded border-border"
        />
        Utiliser la formule par défaut si non associée
      </label>
      <label class="flex items-center gap-1.5 text-foreground cursor-pointer">
        <input
          :checked="generateMissingPiece"
          @change="
            emit(
              'update:generateMissingPiece',
              ($event.target as HTMLInputElement).checked,
            )
          "
          type="checkbox"
          class="rounded border-border"
        />
        Générer N° provisoire (PROV-...) si pièce absente
      </label>
    </div>

    <!-- 3. Bouton & Grille de correspondance des colonnes Excel -->
    <div class="border-t border-border pt-2">
      <button
        type="button"
        @click="showColumnMapper = !showColumnMapper"
        class="flex w-full items-center justify-between text-[11px] font-medium text-foreground hover:text-brand"
      >
        <span class="flex items-center gap-1.5">
          <SlidersHorizontal class="h-3.5 w-3.5 text-brand" />
          Correspondance des colonnes Excel
          <span
            v-if="missingRequiredFields.length > 0"
            class="rounded bg-amber-100 px-1.5 py-0.5 text-[10px] font-semibold text-amber-800"
          >
            {{ missingRequiredFields.length }} colonne(s) requise(s) à vérifier
          </span>
          <span
            v-else-if="unmappedExcelHeaders.length > 0"
            class="text-[10px] text-muted-foreground"
          >
            ({{ unmappedExcelHeaders.length }} colonne(s) supplémentaire(s)
            ignorée(s))
          </span>
        </span>
        <ChevronUp v-if="showColumnMapper" class="h-3.5 w-3.5" />
        <ChevronDown v-else class="h-3.5 w-3.5" />
      </button>

      <div v-if="showColumnMapper" class="mt-2.5 space-y-2">
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-2">
          <div
            v-for="def in SYSTEM_FIELD_LABELS"
            :key="def.field"
            class="flex flex-col gap-0.5 text-[10px]"
          >
            <label class="font-medium text-muted-foreground">
              {{ def.label }}
              <span v-if="def.required" class="text-rose-600">*</span>
            </label>
            <select
              :value="columnMapping[def.field]"
              @change="
                emit(
                  'update-column',
                  def.field,
                  ($event.target as HTMLSelectElement).value,
                )
              "
              class="h-7 rounded-[3px] border border-border bg-card px-1.5 text-[11px] text-foreground"
            >
              <option value="">-- Non mappée / Défaut --</option>
              <option v-for="h in headers" :key="h" :value="h">{{ h }}</option>
            </select>
          </div>
        </div>

        <p
          v-if="unmappedExcelHeaders.length > 0"
          class="text-[10px] text-muted-foreground"
        >
          Colonnes du fichier ignorées (absentes du système) :
          <span class="font-mono text-foreground">{{
            unmappedExcelHeaders.join(", ")
          }}</span>
        </p>
      </div>
    </div>
  </div>
</template>

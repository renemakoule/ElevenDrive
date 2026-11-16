<script setup lang="ts">
import type { Candidat } from "../../types";
import { FolderArchive, CheckCircle2, Archive } from "lucide-vue-next";

interface Props {
  candidates: Candidat[];
  selectedIds: string[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "toggle-select-all"): void;
  (e: "update:selectedIds", ids: string[]): void;
  (e: "batch-archive"): void;
  (e: "archive-single", candidatId: string): void;
}>();

const fcfa = (n: number) => (n || 0).toLocaleString("fr-FR");

const onCheckboxChange = (id: string, checked: boolean) => {
  const current = [...props.selectedIds];
  if (checked) {
    if (!current.includes(id)) current.push(id);
  } else {
    const idx = current.indexOf(id);
    if (idx !== -1) current.splice(idx, 1);
  }
  emit("update:selectedIds", current);
};
</script>

<template>
  <div class="space-y-4">
    <div
      class="flex items-center justify-between bg-card border border-border p-3 rounded-[4px]"
    >
      <div class="flex items-center gap-2">
        <input
          type="checkbox"
          :checked="
            selectedIds.length > 0 && selectedIds.length === candidates.length
          "
          @change="emit('toggle-select-all')"
          class="rounded border-border text-brand focus:ring-brand h-4 w-4"
        />
        <span class="text-xs font-medium text-foreground">
          Tout sélectionner ({{ selectedIds.length }} sélectionné(s))
        </span>
      </div>

      <button
        v-if="selectedIds.length > 0"
        type="button"
        @click="emit('batch-archive')"
        class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-[4px] bg-brand text-brand-dark font-bold text-xs hover:bg-brand-dark hover:text-white transition"
      >
        <FolderArchive class="h-4 w-4" />
        Archiver la sélection ({{ selectedIds.length }})
      </button>
    </div>

    <div class="rounded-[4px] border border-border bg-card overflow-hidden">
      <div v-if="candidates.length === 0" class="p-12 text-center">
        <CheckCircle2 class="h-10 w-10 text-emerald-500/60 mx-auto mb-3" />
        <p class="text-sm font-medium text-foreground">
          Tous les dossiers terminés sont à jour dans les archives !
        </p>
      </div>

      <div v-else class="divide-y divide-border">
        <div
          v-for="cand in candidates"
          :key="cand.id"
          class="p-4 transition hover:bg-muted/30 flex items-center justify-between gap-4"
        >
          <div class="flex items-center gap-3 min-w-0">
            <input
              type="checkbox"
              :checked="selectedIds.includes(cand.id)"
              @change="
                (e: Event) =>
                  onCheckboxChange(
                    cand.id,
                    (e.target as HTMLInputElement).checked,
                  )
              "
              class="rounded border-border text-brand focus:ring-brand h-4 w-4 shrink-0"
            />

            <div class="min-w-0">
              <div class="flex items-center gap-2 mb-0.5 flex-wrap">
                <h4
                  class="text-xs sm:text-sm font-semibold text-foreground truncate"
                >
                  {{ cand.nom }} {{ cand.prenoms }}
                </h4>
                <span class="font-mono text-[11px] text-brand-dark font-bold">
                  N° {{ cand.numero_dossier || "DOS-ATTENTE" }}
                </span>
                <span
                  class="text-[10px] px-1.5 py-0.2 rounded font-medium"
                  :class="
                    cand.solde_restant === 0
                      ? 'bg-emerald-500/10 text-emerald-700 dark:text-emerald-400'
                      : 'bg-amber-500/10 text-amber-700 dark:text-amber-400'
                  "
                >
                  {{
                    cand.solde_restant === 0
                      ? "Solde 100%"
                      : "Solde restant: " + fcfa(cand.solde_restant || 0) + " F"
                  }}
                </span>
              </div>

              <p class="text-xs text-muted-foreground">
                Permis {{ cand.categorie_permis || "B" }} ({{
                  cand.formule_nom || "Formation Permis"
                }}) • WhatsApp: {{ cand.telephone_whatsapp }}
              </p>
            </div>
          </div>

          <button
            type="button"
            @click="emit('archive-single', cand.id)"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-[4px] border border-brand bg-brand/10 text-brand-dark dark:text-brand font-semibold text-xs hover:bg-brand hover:text-brand-dark transition shrink-0"
          >
            <Archive class="h-3.5 w-3.5" />
            Archiver
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

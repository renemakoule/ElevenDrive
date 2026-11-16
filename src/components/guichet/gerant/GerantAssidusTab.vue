<script setup lang="ts">
import {
  Phone,
  CheckCircle2,
  AlertTriangle,
  ShieldAlert,
} from "lucide-vue-next";

defineProps<{
  assidus: any[];
  derogationsAccorde: string[];
}>();

const emit = defineEmits<{
  (e: "override", cand: any): void;
}>();

const pct = (n: number) => Math.max(0, Math.min(100, Math.round(n ?? 0)));
</script>

<template>
  <div class="flex flex-1 flex-col overflow-hidden">
    <div class="flex items-center justify-between gap-3 px-4 py-3">
      <h2 class="text-sm font-semibold text-foreground">
        Au moins 80 % du volume de code
      </h2>
      <span class="text-[11px] text-muted-foreground"
        >Éligibles pour débuter la conduite</span
      >
    </div>

    <ul
      v-if="assidus.length"
      class="flex-1 divide-y divide-border overflow-y-auto border-t border-border"
    >
      <li
        v-for="cand in assidus"
        :key="cand.dossier_id"
        class="flex flex-col gap-2 px-4 py-3 transition hover:bg-muted/50 sm:flex-row sm:items-center sm:justify-between"
      >
        <div class="min-w-0 space-y-1">
          <div class="flex flex-wrap items-center gap-x-2">
            <span class="text-xs font-semibold text-foreground">{{
              cand.candidat_nom
            }}</span>
            <span class="font-mono text-[11px] font-medium text-brand-dark">{{
              cand.numero_dossier
            }}</span>
          </div>
          <div
            class="flex flex-wrap items-center gap-x-3 text-[11px] text-muted-foreground"
          >
            <span>
              Code
              <span class="font-medium tabular-nums text-foreground">
                {{ cand.heures_code_effectuees }} /
                {{ cand.heures_theoriques }} h
              </span>
            </span>
            <span class="inline-flex items-center gap-1">
              <Phone class="h-3 w-3" />
              {{ cand.telephone_whatsapp }}
            </span>
          </div>
        </div>

        <div class="flex items-center gap-2 shrink-0">
          <span
            class="inline-flex shrink-0 items-center gap-1.5 text-[11px] font-medium"
            :class="
              Boolean(cand.conduite_debloquee_financierement) ||
              derogationsAccorde.includes(cand.dossier_id)
                ? 'text-brand-dark'
                : 'text-amber-600'
            "
          >
            <CheckCircle2
              v-if="
                Boolean(cand.conduite_debloquee_financierement) ||
                derogationsAccorde.includes(cand.dossier_id)
              "
              class="h-3.5 w-3.5"
            />
            <AlertTriangle v-else class="h-3.5 w-3.5" />
            {{
              Boolean(cand.conduite_debloquee_financierement) ||
              derogationsAccorde.includes(cand.dossier_id)
                ? "Conduite débloquée"
                : `En attente tranche 2 (${pct(cand.taux_financier_pct)} % payé)`
            }}
          </span>

          <button
            v-if="
              !(
                Boolean(cand.conduite_debloquee_financierement) ||
                derogationsAccorde.includes(cand.dossier_id)
              )
            "
            type="button"
            @click="emit('override', cand)"
            class="inline-flex h-9 sm:h-7 items-center gap-1.5 rounded-[3px] bg-brand-light px-3 sm:px-2.5 text-xs font-semibold text-brand-dark border border-brand-border hover:bg-brand-light/80 transition cursor-pointer"
            title="Débloquer la conduite par dérogation avec code PIN"
          >
            <ShieldAlert class="h-3.5 w-3.5 text-amber-600" />
            <span>Débloquer (PIN)</span>
          </button>
        </div>
      </li>
    </ul>
    <p
      v-else
      class="border-t border-border px-4 py-10 text-center text-xs text-muted-foreground"
    >
      Aucun élève n'a encore atteint 80 % du code.
    </p>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from "../../../stores/authStore";
import { ShieldCheck, Calendar, Monitor, Radio } from "lucide-vue-next";

defineProps<{
  liveLivretActivity: {
    candidatNom: string;
    kmParcourus: number;
    commentaire: string;
    updatedAt: string;
  } | null;
}>();

const authStore = useAuthStore();
</script>

<template>
  <div class="space-y-4">
    <!-- Bandeau Synchronisé Licence & Postes Autorisés -->
    <div
      v-if="authStore.currentSchool"
      class="rounded-[5px] border border-border bg-card px-4 py-2.5 flex flex-col sm:flex-row sm:items-center justify-between gap-3 text-xs"
    >
      <div class="flex flex-wrap items-center gap-3">
        <div class="flex items-center gap-1.5">
          <ShieldCheck class="w-4 h-4 text-emerald-500" />
          <span class="text-muted-foreground">Formule :</span>
          <span class="font-mono font-bold text-foreground">
            {{ authStore.currentSchool.type_licence || "TRIAL" }}
          </span>
        </div>

        <span class="hidden sm:inline text-border">|</span>

        <div class="flex items-center gap-1.5">
          <Calendar class="w-3.5 h-3.5 text-sky-500" />
          <span class="text-muted-foreground">Validité :</span>
          <span
            class="font-mono font-bold"
            :class="
              authStore.trialDaysRemaining <= 15
                ? 'text-amber-500'
                : 'text-emerald-600 dark:text-emerald-400'
            "
          >
            {{ authStore.trialDaysRemaining }}j restants
          </span>
          <span
            v-if="authStore.currentSchool.date_expiration_licence"
            class="text-[11px] font-mono text-muted-foreground"
          >
            (Fin :
            {{
              new Date(
                authStore.currentSchool.date_expiration_licence,
              ).toLocaleDateString("fr-FR")
            }})
          </span>
        </div>
      </div>

      <div class="flex items-center gap-3 min-w-[210px]">
        <div class="flex items-center gap-1.5 shrink-0">
          <Monitor class="w-3.5 h-3.5 text-brand" />
          <span class="text-muted-foreground">Postes :</span>
          <span class="font-mono font-bold text-foreground">
            {{ authStore.currentSchool.machines_actives || 0 }} /
            {{ authStore.currentSchool.max_machines || 3 }}
          </span>
        </div>
        <div
          class="flex-1 h-1.5 min-w-[70px] bg-muted rounded-full overflow-hidden border border-border/50"
        >
          <div
            class="h-full rounded-full transition-all duration-300"
            :class="
              (authStore.currentSchool.machines_actives || 0) >=
              (authStore.currentSchool.max_machines || 3)
                ? 'bg-amber-500'
                : 'bg-emerald-500'
            "
            :style="{
              width: `${Math.min(100, Math.max(0, ((authStore.currentSchool.machines_actives || 0) / (authStore.currentSchool.max_machines || 3)) * 100))}%`,
            }"
          ></div>
        </div>
      </div>
    </div>

    <!-- Alerte Activité Moniteur en Direct -->
    <div
      v-if="liveLivretActivity"
      class="rounded-[5px] border border-emerald-500/30 bg-emerald-500/5 px-4 py-2.5 flex flex-wrap items-center justify-between gap-2 text-xs"
    >
      <div class="flex items-center gap-2">
        <span
          class="inline-flex items-center gap-1 rounded-full bg-emerald-500/20 px-2 py-0.5 text-[10px] font-bold text-emerald-700 dark:text-emerald-300"
        >
          <Radio class="h-3 w-3 animate-pulse" />
          Leçon en direct
        </span>
        <span class="font-semibold text-foreground">{{
          liveLivretActivity.candidatNom
        }}</span>
        <span
          v-if="liveLivretActivity.commentaire"
          class="text-muted-foreground italic"
        >
          — « {{ liveLivretActivity.commentaire }} »
        </span>
      </div>
      <span class="font-mono text-[11px] text-muted-foreground">
        +{{ liveLivretActivity.kmParcourus }} km · Sync
        {{ liveLivretActivity.updatedAt }}
      </span>
    </div>
  </div>
</template>


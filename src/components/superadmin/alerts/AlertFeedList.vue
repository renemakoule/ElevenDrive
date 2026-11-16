<script setup lang="ts">
import type { NetworkAlertItem } from "./types";
import { Building2, Clock, Trash2, Loader2 } from "lucide-vue-next";

defineProps<{
  alerts: NetworkAlertItem[];
  isLoading: boolean;
}>();

const emit = defineEmits<{
  (e: "toggle", id: string): void;
  (e: "delete", id: string): void;
}>();

const formatDate = (iso: string) => {
  try {
    const d = new Date(iso);
    return `${d.toLocaleDateString("fr-FR")} à ${d.toLocaleTimeString("fr-FR", { hour: "2-digit", minute: "2-digit" })}`;
  } catch {
    return iso;
  }
};
</script>

<template>
  <div class="lg:col-span-7 space-y-3">
    <div class="flex items-center justify-between">
      <span
        class="text-xs font-bold text-foreground uppercase tracking-wider flex items-center gap-2"
      >
        <span>Alertes en cours & Historique ({{ alerts.length }})</span>
        <Loader2 v-if="isLoading" class="w-3.5 h-3.5 animate-spin text-brand" />
      </span>
    </div>

    <div
      v-if="isLoading && alerts.length === 0"
      class="p-8 text-center text-xs text-muted-foreground border border-border rounded-[3px] bg-card flex flex-col items-center justify-center gap-2"
    >
      <Loader2 class="w-5 h-5 animate-spin text-brand" />
      <span>Chargement des alertes depuis la base de données...</span>
    </div>

    <div
      v-else-if="alerts.length === 0"
      class="p-8 text-center text-xs text-muted-foreground border border-border rounded-[3px] bg-card"
    >
      Aucune alerte n'a été émise pour l'instant.
    </div>

    <div
      v-for="a in alerts"
      :key="a.id"
      class="border rounded-[3px] p-4 bg-card transition space-y-2.5"
      :class="[
        a.estActive
          ? a.severite === 'critical'
            ? 'border-rose-500/50 bg-rose-500/5'
            : a.severite === 'warning'
              ? 'border-amber-500/50 bg-amber-500/5'
              : 'border-blue-500/50 bg-blue-500/5'
          : 'border-border opacity-70 bg-muted/20',
      ]"
    >
      <div class="flex items-start justify-between gap-2">
        <div class="flex items-center gap-2 flex-wrap">
          <span
            class="text-[10px] font-bold px-2 py-0.5 rounded-[4px] uppercase"
            :class="
              a.severite === 'critical'
                ? 'bg-rose-600 text-white'
                : a.severite === 'warning'
                  ? 'bg-amber-600 text-white'
                  : 'bg-blue-600 text-white'
            "
          >
            {{
              a.severite === "critical"
                ? "Critique"
                : a.severite === "warning"
                  ? "Avertissement"
                  : "Info"
            }}
          </span>

          <span
            class="text-[11px] font-semibold text-foreground flex items-center gap-1 bg-muted px-2 py-0.5 rounded-[4px] border border-border"
          >
            <Building2 class="w-3 h-3 text-muted-foreground" />
            <span>{{ a.cibleSchoolName }}</span>
          </span>

          <span
            class="text-[10px] font-mono px-1.5 py-0.5 rounded"
            :class="
              a.estActive
                ? 'text-emerald-600 bg-emerald-500/10'
                : 'text-muted-foreground bg-muted'
            "
          >
            {{ a.estActive ? "Active" : "Révoquée" }}
          </span>

          <span
            v-if="a.bloquante"
            class="text-[10px] text-rose-600 bg-rose-500/10 px-1.5 py-0.5 rounded font-medium"
          >
            Bloquante
          </span>
        </div>

        <div class="flex items-center gap-1 shrink-0">
          <button
            type="button"
            @click="emit('toggle', a.id)"
            class="px-2 py-1 rounded text-[11px] font-medium border border-border hover:bg-muted transition"
            :title="
              a.estActive ? 'Désactiver cette alerte' : 'Réactiver cette alerte'
            "
          >
            {{ a.estActive ? "Révoquer" : "Réactiver" }}
          </button>
          <button
            type="button"
            @click="emit('delete', a.id)"
            class="p-1 rounded text-muted-foreground hover:text-rose-600 hover:bg-rose-500/10 transition"
            title="Supprimer"
          >
            <Trash2 class="w-3.5 h-3.5" />
          </button>
        </div>
      </div>

      <div>
        <h4 class="text-xs font-bold text-foreground">{{ a.titre }}</h4>
        <p class="text-xs text-muted-foreground mt-1 leading-relaxed">
          {{ a.message }}
        </p>
      </div>

      <div
        class="text-[10px] text-muted-foreground/80 flex items-center gap-1 pt-1 border-t border-border"
      >
        <Clock class="w-3 h-3" />
        <span>Émise le {{ formatDate(a.dateEnvoi) }}</span>
      </div>
    </div>
  </div>
</template>

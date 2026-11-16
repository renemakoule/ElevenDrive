<script setup lang="ts">
import { ref, onMounted } from "vue";
import type { AuditLogEntry } from "../../../types";
import { superAdminApi } from "../../../services/superAdminService";

const logs = ref<AuditLogEntry[]>([]);
const isLoading = ref(false);

const loadLogs = async () => {
  isLoading.value = true;
  try {
    logs.value = await superAdminApi.getAuditLogs();
  } catch (err) {
    console.error("Erreur chargement logs audit:", err);
  } finally {
    isLoading.value = false;
  }
};

const getActionBadgeClass = (action: string) => {
  if (action.includes("DELETE"))
    return "bg-destructive/10 text-destructive border-destructive/30";
  if (action.includes("SUSPEND"))
    return "bg-amber-500/10 text-amber-600 dark:text-amber-400 border-amber-500/30";
  if (action.includes("ACTIVATE") || action.includes("LICENCE"))
    return "bg-brand/10 text-brand border-brand/30";
  return "bg-sky-500/10 text-sky-600 dark:text-sky-400 border-sky-500/30";
};

onMounted(() => {
  loadLogs();
});

defineExpose({
  reload: loadLogs,
});
</script>

<template>
  <main class="flex-1 flex flex-col overflow-hidden bg-background">
    <!-- Barre de Filtres Secondaire -->
    <div
      class="h-11 border-b border-border bg-card px-4 flex items-center justify-between shrink-0"
    >
      <div class="flex items-center gap-2 text-xs">
        <span class="font-semibold text-foreground"
          >Journal d'Audit Système & Sécurité</span
        >
        <span class="text-muted-foreground font-mono"
          >({{ logs.length }} événements récents)</span
        >
      </div>

      <div class="text-[11px] font-mono text-muted-foreground">
        Traçabilité: <span class="text-brand font-medium">Immuable & Signé</span>
      </div>
    </div>

    <!-- Tableau des logs d'audit -->
    <div class="flex-1 overflow-auto">
      <table class="w-full text-left text-xs border-collapse">
        <thead
          class="sticky top-0 bg-muted/80 backdrop-blur-xs border-b border-border text-[11px] text-muted-foreground font-mono uppercase tracking-wider select-none z-10"
        >
          <tr>
            <th class="py-2.5 px-3">Date & Heure</th>
            <th class="py-2.5 px-3">Action Système</th>
            <th class="py-2.5 px-3">Entité Cible</th>
            <th class="py-2.5 px-3">Détails de l'Opération</th>
            <th class="py-2.5 px-3 text-right">Adresse IP</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border font-mono text-foreground">
          <tr
            v-for="log in logs"
            :key="log.id"
            class="hover:bg-muted/50 transition-colors"
          >
            <!-- Date & Heure -->
            <td class="py-2.5 px-3 text-muted-foreground">
              {{ new Date(log.created_at).toLocaleString("fr-FR") }}
            </td>

            <!-- Action -->
            <td class="py-2.5 px-3">
              <span
                :class="[
                  'px-2 py-0.5 rounded-[3px] text-[10px] font-bold border',
                  getActionBadgeClass(log.action),
                ]"
              >
                {{ log.action }}
              </span>
            </td>

            <!-- Cible -->
            <td class="py-2.5 px-3 text-foreground font-medium">
              <span class="text-muted-foreground text-[10px]"
                >[{{ log.cible_type }}]</span
              >
              {{ log.cible_id.slice(0, 12) }}...
            </td>

            <!-- Détails -->
            <td class="py-2.5 px-3 font-sans text-foreground">
              {{ log.details }}
            </td>

            <!-- IP -->
            <td class="py-2.5 px-3 text-right text-muted-foreground font-mono">
              {{ log.ip_address || "127.0.0.1" }}
            </td>
          </tr>

          <tr v-if="logs.length === 0 && !isLoading">
            <td colspan="5" class="text-center py-12 text-muted-foreground font-sans">
              Aucune action enregistrée dans le journal pour le moment.
            </td>
          </tr>
          <tr v-if="isLoading && logs.length === 0">
            <td
              colspan="5"
              class="text-center py-12 text-muted-foreground font-sans font-mono animate-pulse"
            >
              Chargement du journal d'audit...
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </main>
</template>

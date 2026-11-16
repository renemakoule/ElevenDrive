<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";
import type { SuperAdminGlobalCandidat } from "../../../types";
import { superAdminApi } from "../../../services/superAdminService";
import { Building2 } from "lucide-vue-next";

const props = defineProps<{
  searchQuery: string;
}>();

const candidats = ref<SuperAdminGlobalCandidat[]>([]);
const isLoading = ref(false);
const statusFilter = ref<string>("ALL");

const loadCandidats = async () => {
  isLoading.value = true;
  try {
    candidats.value = await superAdminApi.getCandidatsData(props.searchQuery);
  } catch (err) {
    console.error("Erreur chargement candidats globaux:", err);
  } finally {
    isLoading.value = false;
  }
};

watch(
  () => props.searchQuery,
  () => {
    loadCandidats();
  },
);

const filteredCandidats = computed(() => {
  return candidats.value.filter((c) => {
    if (statusFilter.value === "ALL") return true;
    return c.statut_dossier.toUpperCase() === statusFilter.value;
  });
});

const formatFCFA = (val: number) => {
  return (val || 0).toLocaleString("fr-FR");
};

const getStatusBadgeClass = (status: string) => {
  switch (status.toUpperCase()) {
    case "ADMIS":
      return "bg-brand/10 text-brand border-brand/30";
    case "EN_COURS":
      return "bg-sky-500/10 text-sky-600 dark:text-sky-400 border-sky-500/30";
    case "PREPARE_EXAMEN":
      return "bg-amber-500/10 text-amber-600 dark:text-amber-400 border-amber-500/30";
    case "AJOURNE":
      return "bg-rose-500/10 text-rose-600 dark:text-rose-400 border-rose-500/30";
    default:
      return "bg-muted text-muted-foreground border-border";
  }
};

onMounted(() => {
  loadCandidats();
});

defineExpose({
  reload: loadCandidats,
});
</script>

<template>
  <main class="flex-1 flex flex-col overflow-hidden bg-background">
    <!-- Barre de Filtres Secondaire -->
    <div
      class="h-11 border-b border-border bg-card px-4 flex items-center justify-between shrink-0"
    >
      <div class="flex items-center gap-2 text-xs">
        <span class="font-semibold text-foreground">Candidats Inscrits</span>
        <span class="text-muted-foreground font-mono"
          >({{ filteredCandidats.length }})</span
        >

        <div class="h-3 w-px bg-border mx-2"></div>

        <div
          class="flex items-center rounded-[4px] bg-muted p-0.5 border border-border"
        >
          <button
            type="button"
            @click="statusFilter = 'ALL'"
            :class="[
              'px-2 py-0.5 text-[11px] rounded-[3px] transition-colors cursor-pointer',
              statusFilter === 'ALL'
                ? 'bg-card text-foreground font-semibold shadow-xs'
                : 'text-muted-foreground hover:text-foreground',
            ]"
          >
            Tous
          </button>
          <button
            type="button"
            @click="statusFilter = 'EN_COURS'"
            :class="[
              'px-2 py-0.5 text-[11px] rounded-[3px] transition-colors cursor-pointer',
              statusFilter === 'EN_COURS'
                ? 'bg-sky-500/15 text-sky-700 dark:text-sky-300 font-semibold shadow-xs'
                : 'text-muted-foreground hover:text-foreground',
            ]"
          >
            En formation
          </button>
          <button
            type="button"
            @click="statusFilter = 'PREPARE_EXAMEN'"
            :class="[
              'px-2 py-0.5 text-[11px] rounded-[3px] transition-colors cursor-pointer',
              statusFilter === 'PREPARE_EXAMEN'
                ? 'bg-amber-500/15 text-amber-700 dark:text-amber-300 font-semibold shadow-xs'
                : 'text-muted-foreground hover:text-foreground',
            ]"
          >
            Prêt examen
          </button>
          <button
            type="button"
            @click="statusFilter = 'ADMIS'"
            :class="[
              'px-2 py-0.5 text-[11px] rounded-[3px] transition-colors cursor-pointer',
              statusFilter === 'ADMIS'
                ? 'bg-brand/15 text-brand font-semibold shadow-xs'
                : 'text-muted-foreground hover:text-foreground',
            ]"
          >
            Admis
          </button>
          <button
            type="button"
            @click="statusFilter = 'AJOURNE'"
            :class="[
              'px-2 py-0.5 text-[11px] rounded-[3px] transition-colors cursor-pointer',
              statusFilter === 'AJOURNE'
                ? 'bg-rose-500/15 text-rose-700 dark:text-rose-300 font-semibold shadow-xs'
                : 'text-muted-foreground hover:text-foreground',
            ]"
          >
            Ajournés
          </button>
        </div>
      </div>

      <div class="text-[11px] font-mono text-muted-foreground">
        Supervision:
        <span class="text-brand font-medium">Parcours & Apprenants Permis</span>
      </div>
    </div>

    <!-- Tableau Interactif Style Console Neon / MaryamE -->
    <div class="flex-1 overflow-auto">
      <table class="w-full text-left text-xs border-collapse">
        <thead
          class="sticky top-0 bg-muted/80 backdrop-blur-sm border-b border-border text-[11px] text-muted-foreground font-mono uppercase tracking-wider select-none z-10"
        >
          <tr>
            <th class="py-2.5 px-3">Candidat</th>
            <th class="py-2.5 px-3">N° Dossier</th>
            <th class="py-2.5 px-3">Auto-École</th>
            <th class="py-2.5 px-3">WhatsApp / Téléphone</th>
            <th class="py-2.5 px-3">Statut Dossier</th>
            <th class="py-2.5 px-3 text-right">Forfait</th>
            <th class="py-2.5 px-3 text-right">Solde Restant</th>
            <th class="py-2.5 px-3 text-right">Date Inscription</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border font-mono text-foreground">
          <tr
            v-for="candidat in filteredCandidats"
            :key="candidat.id"
            class="hover:bg-muted/50 transition-colors"
          >
            <!-- Nom, Prénom, Genre -->
            <td class="py-2.5 px-3 font-sans">
              <div class="flex items-center gap-2">
                <div
                  class="w-6 h-6 rounded-full bg-brand/10 border border-brand/20 flex items-center justify-center text-[10px] font-mono font-bold text-brand shrink-0"
                >
                  {{
                    (
                      candidat.prenoms?.[0] ||
                      candidat.nom?.[0] ||
                      "C"
                    ).toUpperCase()
                  }}
                </div>
                <div>
                  <div class="font-semibold text-foreground">
                    {{ candidat.nom }} {{ candidat.prenoms }}
                    <span
                      class="text-[10px] text-muted-foreground font-normal ml-1"
                      >({{ candidat.genre }})</span
                    >
                  </div>
                </div>
              </div>
            </td>

            <!-- Dossier -->
            <td class="py-2.5 px-3 font-mono text-muted-foreground">
              {{ candidat.numero_dossier || "EN ATTENTE" }}
            </td>

            <!-- Auto-école -->
            <td class="py-2.5 px-3 font-sans text-foreground">
              <div class="flex items-center gap-1.5">
                <Building2 class="w-3.5 h-3.5 text-muted-foreground shrink-0" />
                <span class="truncate max-w-[170px]">{{
                  candidat.auto_ecole_nom
                }}</span>
              </div>
            </td>

            <!-- Téléphone WhatsApp -->
            <td class="py-2.5 px-3 text-muted-foreground font-mono">
              {{ candidat.telephone_whatsapp }}
            </td>

            <!-- Statut Dossier -->
            <td class="py-2.5 px-3">
              <span
                :class="[
                  'inline-flex items-center gap-1 px-1.5 py-0.5 rounded-[4px] text-[10px] font-semibold border',
                  getStatusBadgeClass(candidat.statut_dossier),
                ]"
              >
                {{ candidat.statut_dossier }}
              </span>
            </td>

            <!-- Forfait total -->
            <td class="py-2.5 px-3 text-right font-medium text-foreground">
              {{ formatFCFA(candidat.montant_total) }}
              <span class="text-[10px] text-muted-foreground">FCFA</span>
            </td>

            <!-- Solde restant -->
            <td
              class="py-2.5 px-3 text-right font-bold"
              :class="
                candidat.solde_restant > 0
                  ? 'text-amber-600 dark:text-amber-400'
                  : 'text-brand'
              "
            >
              {{ formatFCFA(candidat.solde_restant) }}
            </td>

            <!-- Date -->
            <td class="py-2.5 px-3 text-right text-muted-foreground">
              {{ new Date(candidat.created_at).toLocaleDateString("fr-FR") }}
            </td>
          </tr>

          <tr v-if="filteredCandidats.length === 0 && !isLoading">
            <td
              colspan="8"
              class="text-center py-12 text-muted-foreground font-sans"
            >
              Aucun candidat trouvé pour cette sélection.
            </td>
          </tr>
          <tr v-if="isLoading && filteredCandidats.length === 0">
            <td
              colspan="8"
              class="text-center py-12 text-muted-foreground font-sans font-mono animate-pulse"
            >
              Chargement des candidats...
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </main>
</template>

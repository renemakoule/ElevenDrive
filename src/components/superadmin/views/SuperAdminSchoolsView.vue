<script setup lang="ts">
import { ref, computed } from "vue";
import type {
  SchoolItemSummary,
  UpdateSchoolAdminPayload,
  UpdateSchoolLicencePayload,
} from "../../../types";
import { superAdminApi } from "../../../services/superAdminService";
import SchoolInspectionPanel from "../panels/SchoolInspectionPanel.vue";
import EditSchoolModal from "../modals/EditSchoolModal.vue";
import ManageLicenceModal from "../modals/ManageLicenceModal.vue";
import ManageModulesModal from "../modals/ManageModulesModal.vue";
import DeleteSchoolModal from "../modals/DeleteSchoolModal.vue";
import SuperAdminDashboardSkeleton from "../../common/skeleton/SuperAdminDashboardSkeleton.vue";
import { Monitor } from "lucide-vue-next";

const props = defineProps<{
  schools: SchoolItemSummary[];
  isLoading: boolean;
}>();

const emit = defineEmits<{
  (e: "toggle-status", school: SchoolItemSummary): void;
  (e: "reload"): void;
}>();

const selectedFilter = ref<"all" | "active" | "suspended">("all");
const selectedSchool = ref<SchoolItemSummary | null>(null);

// État des modales
const isEditModalOpen = ref(false);
const isLicenceModalOpen = ref(false);
const isModulesModalOpen = ref(false);
const isDeleteModalOpen = ref(false);
const activeModalSchool = ref<SchoolItemSummary | null>(null);
const isActionLoading = ref(false);

const filteredSchools = computed(() => {
  return props.schools.filter((s) => {
    if (selectedFilter.value === "active") return s.est_actif;
    if (selectedFilter.value === "suspended") return !s.est_actif;
    return true;
  });
});

const formatFCFA = (val: number) => {
  return (val || 0).toLocaleString("fr-FR");
};

const selectSchool = (school: SchoolItemSummary) => {
  selectedSchool.value = school;
};

// Actions Modales
const openEditModal = (school: SchoolItemSummary) => {
  activeModalSchool.value = school;
  isEditModalOpen.value = true;
};

const openLicenceModal = (school: SchoolItemSummary) => {
  activeModalSchool.value = school;
  isLicenceModalOpen.value = true;
};

const openModulesModal = (school: SchoolItemSummary) => {
  activeModalSchool.value = school;
  isModulesModalOpen.value = true;
};

const openDeleteModal = (school: SchoolItemSummary) => {
  activeModalSchool.value = school;
  isDeleteModalOpen.value = true;
};

const handleSaveEdit = async (payload: UpdateSchoolAdminPayload) => {
  isActionLoading.value = true;
  try {
    await superAdminApi.updateSchool(payload);
    isEditModalOpen.value = false;
    emit("reload");
  } catch (err: any) {
    alert("Erreur lors de la mise à jour : " + err?.message);
  } finally {
    isActionLoading.value = false;
  }
};

const handleSaveLicence = async (payload: UpdateSchoolLicencePayload) => {
  isActionLoading.value = true;
  try {
    await superAdminApi.updateLicence(payload);
    isLicenceModalOpen.value = false;
    emit("reload");
  } catch (err: any) {
    alert("Erreur configuration licence : " + err?.message);
  } finally {
    isActionLoading.value = false;
  }
};

const handleSaveModules = async (schoolId: string, modules: string[]) => {
  isActionLoading.value = true;
  try {
    await superAdminApi.updateModules(schoolId, modules);
    isModulesModalOpen.value = false;
    emit("reload");
  } catch (err: any) {
    alert("Erreur configuration des modules : " + err?.message);
  } finally {
    isActionLoading.value = false;
  }
};

const handleDeleteConfirm = async (schoolId: string, confirmName: string) => {
  isActionLoading.value = true;
  try {
    await superAdminApi.deleteSchool(schoolId, confirmName);
    isDeleteModalOpen.value = false;
    if (selectedSchool.value?.id === schoolId) {
      selectedSchool.value = null;
    }
    emit("reload");
  } catch (err: any) {
    alert("Erreur lors de la suppression : " + err?.message);
  } finally {
    isActionLoading.value = false;
  }
};
</script>

<template>
  <main class="flex-1 flex flex-col overflow-hidden bg-background">
    <!-- Barre de Filtres Secondaire -->
    <div
      class="h-11 border-b border-border bg-card px-4 flex items-center justify-between shrink-0"
    >
      <div class="flex items-center gap-2 text-xs">
        <span class="font-semibold text-foreground">Auto-Écoles Enregistrées</span>
        <span class="text-muted-foreground font-mono"
          >({{ filteredSchools.length }})</span
        >

        <div class="h-3 w-px bg-border mx-2"></div>

        <div
          class="flex items-center rounded-[4px] bg-muted p-0.5 border border-border"
        >
          <button
            type="button"
            @click="selectedFilter = 'all'"
            :class="[
              'px-2 py-0.5 text-[11px] rounded-[3px] transition-colors cursor-pointer',
              selectedFilter === 'all'
                ? 'bg-card text-foreground font-semibold shadow-xs'
                : 'text-muted-foreground hover:text-foreground',
            ]"
          >
            Toutes
          </button>
          <button
            type="button"
            @click="selectedFilter = 'active'"
            :class="[
              'px-2 py-0.5 text-[11px] rounded-[3px] transition-colors cursor-pointer',
              selectedFilter === 'active'
                ? 'bg-card text-foreground font-semibold shadow-xs'
                : 'text-muted-foreground hover:text-foreground',
            ]"
          >
            Actives
          </button>
          <button
            type="button"
            @click="selectedFilter = 'suspended'"
            :class="[
              'px-2 py-0.5 text-[11px] rounded-[3px] transition-colors cursor-pointer',
              selectedFilter === 'suspended'
                ? 'bg-card text-foreground font-semibold shadow-xs'
                : 'text-muted-foreground hover:text-foreground',
            ]"
          >
            Suspendues
          </button>
        </div>
      </div>

      <div class="text-[11px] font-mono text-muted-foreground">
        Supervision:
        <span class="text-brand font-medium">Flotte & Droits d'Accès</span>
      </div>
    </div>

    <!-- Tableau Interactif Style Console Neon / Base de Données -->
    <div class="flex-1 flex overflow-hidden">
      <!-- État Skeleton Chargement Initial -->
      <div v-if="isLoading && filteredSchools.length === 0" class="flex-1 p-3 overflow-auto">
        <SuperAdminDashboardSkeleton :rows="8" />
      </div>

      <!-- Grille des enregistrements -->
      <div v-else class="flex-1 overflow-auto">
        <table class="w-full text-left text-xs border-collapse">
          <thead
            class="sticky top-0 bg-muted/80 backdrop-blur-xs border-b border-border text-[11px] text-muted-foreground font-mono uppercase tracking-wider select-none z-10"
          >
            <tr>
              <th class="py-2.5 px-3">Statut</th>
              <th class="py-2.5 px-3">Établissement / Nom</th>
              <th class="py-2.5 px-3">Slug (Code)</th>
              <th class="py-2.5 px-3 text-center">Postes / Machines</th>
              <th class="py-2.5 px-3">Licence & Échéance</th>
              <th class="py-2.5 px-3">Téléphone</th>
              <th class="py-2.5 px-3 text-right">Agences</th>
              <th class="py-2.5 px-3 text-right">Candidats</th>
              <th class="py-2.5 px-3 text-right">Volume</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-border font-mono text-foreground">
            <tr
              v-for="school in filteredSchools"
              :key="school.id"
              @click="selectSchool(school)"
              :class="[
                'hover:bg-muted/50 cursor-pointer transition-colors',
                selectedSchool?.id === school.id
                  ? 'bg-muted text-foreground font-semibold'
                  : '',
              ]"
            >
              <!-- Statut -->
              <td class="py-2 px-3">
                <span
                  :class="[
                    'inline-flex items-center gap-1 px-1.5 py-0.5 rounded-[3px] text-[10px] font-semibold',
                    school.est_actif
                      ? 'bg-brand/10 text-brand border border-brand/30'
                      : 'bg-destructive/10 text-destructive border border-destructive/30',
                  ]"
                >
                  <span
                    class="w-1.5 h-1.5 rounded-full"
                    :class="school.est_actif ? 'bg-brand' : 'bg-destructive'"
                  ></span>
                  {{ school.est_actif ? "ACTIF" : "SUSPENDU" }}
                </span>
              </td>

              <!-- Nom Commercial -->
              <td class="py-2 px-3 font-sans font-semibold text-foreground">
                <div>{{ school.nom_commercial }}</div>
                <div class="text-[10px] text-muted-foreground font-mono">
                  {{ school.ville || "Douala" }}, {{ school.pays }}
                </div>
              </td>

              <!-- Slug unique -->
              <td class="py-2 px-3 text-brand font-mono font-medium">
                /{{ school.code_unique }}
              </td>

              <!-- Quota Machines -->
              <td class="py-2 px-3 text-center">
                <span
                  class="inline-flex items-center gap-1 px-2 py-0.5 rounded-[3px] bg-muted border border-border text-[11px] font-mono"
                >
                  <Monitor class="w-3 h-3 text-brand" />
                  <strong class="text-foreground">{{
                    school.machines_actives || 0
                  }}</strong>
                  <span class="text-muted-foreground">/</span>
                  <span class="text-muted-foreground"
                    >{{ school.max_machines || 3 }} max</span
                  >
                </span>
              </td>

              <!-- Licence & Échéance -->
              <td class="py-2 px-3">
                <div class="flex items-center gap-1.5">
                  <span
                    :class="[
                      'w-1.5 h-1.5 rounded-full',
                      (school.jours_restants ?? 365) > 15
                        ? 'bg-brand'
                        : 'bg-amber-500',
                    ]"
                  ></span>
                  <span class="text-foreground font-medium">{{
                    school.type_licence || "ANNUEL"
                  }}</span>
                  <span class="text-[10px] text-muted-foreground"
                    >({{ school.jours_restants ?? 365 }}j)</span
                  >
                </div>
              </td>

              <!-- Téléphone -->
              <td class="py-2 px-3 text-muted-foreground">
                {{ school.telephone }}
              </td>

              <!-- Agences -->
              <td class="py-2 px-3 text-right text-foreground">
                {{ school.agences_count }}
              </td>

              <!-- Candidats -->
              <td class="py-2 px-3 text-right text-brand font-medium">
                {{ school.candidats_count }}
              </td>

              <!-- Volume -->
              <td class="py-2 px-3 text-right font-bold text-foreground">
                {{ formatFCFA(school.total_volume) }}
              </td>
            </tr>

            <tr v-if="filteredSchools.length === 0 && !isLoading">
              <td
                colspan="9"
                class="text-center py-12 text-muted-foreground font-sans"
              >
                Aucun établissement trouvé pour cette recherche.
              </td>
            </tr>
            <tr v-if="isLoading && filteredSchools.length === 0">
              <td
                colspan="9"
                class="text-center py-12 text-muted-foreground font-sans font-mono animate-pulse"
              >
                Chargement des auto-écoles...
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Volet d'Inspection Latéral Réutilisable -->
      <SchoolInspectionPanel
        v-if="selectedSchool"
        :school="selectedSchool"
        @toggle-status="emit('toggle-status', $event)"
        @edit="openEditModal($event)"
        @manage-licence="openLicenceModal($event)"
        @manage-modules="openModulesModal($event)"
        @delete="openDeleteModal($event)"
        @close="selectedSchool = null"
      />
    </div>

    <!-- Modales d'Action Super-Admin -->
    <EditSchoolModal
      :is-open="isEditModalOpen"
      :school="activeModalSchool"
      :is-loading="isActionLoading"
      @close="isEditModalOpen = false"
      @save="handleSaveEdit"
    />

    <ManageLicenceModal
      :is-open="isLicenceModalOpen"
      :school="activeModalSchool"
      :is-loading="isActionLoading"
      @close="isLicenceModalOpen = false"
      @save="handleSaveLicence"
    />

    <ManageModulesModal
      :is-open="isModulesModalOpen"
      :school="activeModalSchool"
      :is-loading="isActionLoading"
      @close="isModulesModalOpen = false"
      @save="handleSaveModules"
    />

    <DeleteSchoolModal
      :is-open="isDeleteModalOpen"
      :school="activeModalSchool"
      :is-loading="isActionLoading"
      @close="isDeleteModalOpen = false"
      @confirm="handleDeleteConfirm"
    />
  </main>
</template>

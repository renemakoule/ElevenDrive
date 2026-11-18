<script setup lang="ts">
import { ref, computed } from "vue";
import type {
  SchoolItemSummary,
  UpdateSchoolLicencePayload,
} from "../../../types";
import { superAdminApi } from "../../../services/superAdminService";
import ManageLicenceModal from "../modals/ManageLicenceModal.vue";
import Skeleton from "../../common/skeleton/Skeleton.vue";
import { ShieldCheck, Monitor, Calendar } from "lucide-vue-next";

const props = defineProps<{
  schools: SchoolItemSummary[];
  isLoading: boolean;
}>();

const emit = defineEmits<{
  (e: "reload"): void;
}>();

const selectedSchool = ref<SchoolItemSummary | null>(null);
const isLicenceModalOpen = ref(false);
const isActionLoading = ref(false);
const filterStatus = ref<"ALL" | "ACTIVE" | "EXPIREE" | "EXPIRING_SOON">("ALL");

const filteredList = computed(() => {
  return props.schools.filter((s) => {
    const days = s.jours_restants ?? 365;
    if (filterStatus.value === "ACTIVE") return s.est_actif && days > 0;
    if (filterStatus.value === "EXPIREE") return !s.est_actif || days <= 0;
    if (filterStatus.value === "EXPIRING_SOON") return days >= 0 && days <= 15;
    return true;
  });
});

const openLicenceModal = (school: SchoolItemSummary) => {
  selectedSchool.value = school;
  isLicenceModalOpen.value = true;
};

const handleSaveLicence = async (payload: UpdateSchoolLicencePayload) => {
  isActionLoading.value = true;
  try {
    await superAdminApi.updateLicence(payload);
    isLicenceModalOpen.value = false;
    emit("reload");
  } catch (err: any) {
    alert("Erreur mise à jour licence : " + err?.message);
  } finally {
    isActionLoading.value = false;
  }
};

const quickExtend = async (
  school: SchoolItemSummary,
  days: number,
  label: string,
) => {
  if (
    !confirm(`Prolonger la licence de ${school.nom_commercial} de ${label} ?`)
  )
    return;
  try {
    await superAdminApi.updateLicence({
      school_id: school.id,
      max_machines: school.max_machines || 3,
      type_licence: label,
      duree_valeur: days,
      duree_unite: "JOURS",
      mode: "EXTEND",
    });
    emit("reload");
  } catch (err: any) {
    alert("Erreur prolongation : " + err?.message);
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
        <span class="font-semibold text-foreground"
          >Postes & Abonnements Établissements</span
        >
        <span class="text-muted-foreground font-mono"
          >({{ filteredList.length }})</span
        >

        <div class="h-3 w-px bg-border mx-2"></div>

        <div
          class="flex items-center rounded-[4px] bg-muted p-0.5 border border-border"
        >
          <button
            type="button"
            @click="filterStatus = 'ALL'"
            :class="[
              'px-2 py-0.5 text-[11px] rounded-[3px] transition-colors cursor-pointer',
              filterStatus === 'ALL'
                ? 'bg-card text-foreground font-semibold shadow-xs'
                : 'text-muted-foreground hover:text-foreground',
            ]"
          >
            Toutes
          </button>
          <button
            type="button"
            @click="filterStatus = 'ACTIVE'"
            :class="[
              'px-2 py-0.5 text-[11px] rounded-[3px] transition-colors cursor-pointer',
              filterStatus === 'ACTIVE'
                ? 'bg-card text-foreground font-semibold shadow-xs'
                : 'text-muted-foreground hover:text-foreground',
            ]"
          >
            Actives
          </button>
          <button
            type="button"
            @click="filterStatus = 'EXPIRING_SOON'"
            :class="[
              'px-2 py-0.5 text-[11px] rounded-[3px] transition-colors cursor-pointer',
              filterStatus === 'EXPIRING_SOON'
                ? 'bg-card text-amber-600 dark:text-amber-400 font-semibold shadow-xs'
                : 'text-muted-foreground hover:text-foreground',
            ]"
          >
            Expire bientôt (&lt;15j)
          </button>
          <button
            type="button"
            @click="filterStatus = 'EXPIREE'"
            :class="[
              'px-2 py-0.5 text-[11px] rounded-[3px] transition-colors cursor-pointer',
              filterStatus === 'EXPIREE'
                ? 'bg-card text-destructive font-semibold shadow-xs'
                : 'text-muted-foreground hover:text-foreground',
            ]"
          >
            Expirées / Suspendues
          </button>
        </div>
      </div>

      <div class="text-[11px] font-mono text-muted-foreground">
        Règle:
        <span class="text-brand font-medium"
          >Enforcement Quotas & Expiration</span
        >
      </div>
    </div>

    <!-- Grille des Licences & Quotas -->
    <div class="flex-1 overflow-auto p-4">
      <!-- État Skeleton Chargement -->
      <div
        v-if="isLoading && filteredList.length === 0"
        class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4"
      >
        <div
          v-for="i in 6"
          :key="i"
          class="border border-border rounded-[4px] p-4 bg-card space-y-3"
        >
          <div class="flex items-center justify-between">
            <Skeleton class="h-4 w-32" />
            <Skeleton class="h-5 w-16 rounded-full" />
          </div>
          <Skeleton class="h-3 w-48" />
          <div class="space-y-2 pt-2 border-t border-border">
            <Skeleton class="h-3 w-28" />
            <Skeleton class="h-4 w-full rounded" />
            <div class="flex items-center justify-between pt-1">
              <Skeleton class="h-3 w-20" />
              <Skeleton class="h-6 w-20 rounded-[3px]" />
            </div>
          </div>
        </div>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div
          v-for="school in filteredList"
          :key="school.id"
          class="bg-card border border-border rounded-[8px] p-4 flex flex-col justify-between hover:border-brand/40 transition shadow-xs"
        >
          <div class="space-y-3">
            <!-- Nom et Statut -->
            <div class="flex items-start justify-between">
              <div>
                <h4 class="font-bold text-foreground text-sm">
                  {{ school.nom_commercial }}
                </h4>
                <div class="text-[10px] text-muted-foreground font-mono">
                  /{{ school.code_unique }} · {{ school.ville || "Douala" }}
                </div>
              </div>

              <span
                :class="[
                  'px-2 py-0.5 rounded-[3px] text-[10px] font-mono font-bold',
                  school.statut_licence === 'ACTIVE'
                    ? 'bg-brand/10 text-brand border border-brand/30'
                    : 'bg-destructive/10 text-destructive border border-destructive/30',
                ]"
              >
                {{ school.statut_licence || "ACTIVE" }}
              </span>
            </div>

            <!-- Jauge de validité temporelle -->
            <div
              class="p-3 rounded-[4px] bg-muted/50 border border-border space-y-2"
            >
              <div class="flex items-center justify-between text-xs">
                <span class="text-muted-foreground flex items-center gap-1.5">
                  <Calendar class="w-3.5 h-3.5 text-sky-500" />
                  <span>Durée restante :</span>
                </span>
                <span
                  class="font-mono font-bold"
                  :class="[
                    (school.jours_restants ?? 365) <= 0
                      ? 'text-destructive'
                      : (school.jours_restants ?? 365) <= 15
                        ? 'text-amber-500'
                        : 'text-brand',
                  ]"
                >
                  {{
                    (school.jours_restants ?? 365) <= 0
                      ? "Expirée"
                      : `${school.jours_restants ?? 365} jours`
                  }}
                </span>
              </div>

              <div
                class="text-[10px] text-muted-foreground font-mono flex items-center justify-between"
              >
                <span>Formule: {{ school.type_licence || "ANNUEL" }}</span>
                <span
                  >Fin:
                  {{
                    school.date_expiration_licence
                      ? new Date(
                          school.date_expiration_licence,
                        ).toLocaleDateString("fr-FR")
                      : "31/12/2026"
                  }}</span
                >
              </div>
            </div>

            <!-- Quota de Machines / Postes Desktop -->
            <div
              class="p-3 rounded-[4px] bg-muted/50 border border-border space-y-1.5"
            >
              <div class="flex items-center justify-between text-xs">
                <span class="text-muted-foreground flex items-center gap-1.5">
                  <Monitor class="w-3.5 h-3.5 text-brand" />
                  <span>Postes installés :</span>
                </span>
                <span class="font-mono font-bold text-foreground">
                  {{ school.machines_actives || 0 }} /
                  {{ school.max_machines || 3 }} max
                </span>
              </div>

              <!-- Barre de progression Postes synchronisée -->
              <div class="h-1.5 w-full bg-muted/80 rounded-full overflow-hidden">
                <div
                  class="h-full rounded-full transition-all duration-300"
                  :class="[
                    (school.machines_actives || 0) >= (school.max_machines || 3)
                      ? 'bg-amber-500'
                      : 'bg-brand',
                  ]"
                  :style="{
                    width: `${Math.min(100, Math.max(0, ((school.machines_actives || 0) / (school.max_machines || 3)) * 100))}%`,
                  }"
                ></div>
              </div>
            </div>
          </div>

          <!-- Actions de Gestion de Licence -->
          <div class="pt-3 border-t border-border mt-3 space-y-2">
            <div class="grid grid-cols-3 gap-1 text-[11px] font-mono">
              <button
                type="button"
                @click="quickExtend(school, 30, 'MENSUEL (+30j)')"
                class="py-1 rounded-[3px] bg-muted hover:bg-muted/80 text-foreground text-center border border-border transition cursor-pointer"
              >
                +30 Jours
              </button>
              <button
                type="button"
                @click="quickExtend(school, 90, 'TRIMESTRIEL (+90j)')"
                class="py-1 rounded-[3px] bg-muted hover:bg-muted/80 text-foreground text-center border border-border transition cursor-pointer"
              >
                +3 Mois
              </button>
              <button
                type="button"
                @click="quickExtend(school, 365, 'ANNUEL (+1 an)')"
                class="py-1 rounded-[3px] bg-muted hover:bg-muted/80 text-brand font-semibold text-center border border-border transition cursor-pointer"
              >
                +1 An
              </button>
            </div>

            <button
              type="button"
              @click="openLicenceModal(school)"
              class="w-full py-1.5 px-3 rounded-[4px] text-xs font-semibold bg-brand/10 hover:bg-brand/20 text-brand border border-brand/30 flex items-center justify-center gap-1.5 transition cursor-pointer"
            >
              <ShieldCheck class="w-3.5 h-3.5" />
              <span>Configurer Postes & Durée Précise</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Modale de Licence -->
    <ManageLicenceModal
      :is-open="isLicenceModalOpen"
      :school="selectedSchool"
      :is-loading="isActionLoading"
      @close="isLicenceModalOpen = false"
      @save="handleSaveLicence"
    />
  </main>
</template>

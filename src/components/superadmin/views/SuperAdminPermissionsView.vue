<script setup lang="ts">
import { ref, computed } from "vue";
import type { SchoolItemSummary } from "../../../types";
import { superAdminApi } from "../../../services/superAdminService";
import { Check, Lock } from "lucide-vue-next";

const props = defineProps<{
  schools: SchoolItemSummary[];
  isLoading: boolean;
}>();

const emit = defineEmits<{
  (e: "reload"): void;
}>();

const modulesConfig = [
  { id: "guichet", label: "Guichet", icon: "FileText" },
  { id: "caisse", label: "Caisse & Reçus", icon: "CreditCard" },
  { id: "conduite", label: "Livret Conduite", icon: "Car" },
  { id: "examens", label: "Examens État", icon: "GraduationCap" },
  { id: "flotte", label: "Flotte Véhicules", icon: "Truck" },
  { id: "archivage_cloud", label: "Archivage Cloud R2", icon: "Cloud" },
];

const parseSchoolModules = (modulesJSON?: string): string[] => {
  if (!modulesJSON) {
    return [
      "guichet",
      "caisse",
      "conduite",
      "examens",
      "flotte",
      "archivage_cloud",
    ];
  }
  try {
    const parsed = JSON.parse(modulesJSON);
    if (Array.isArray(parsed)) return parsed;
  } catch {}
  return [
    "guichet",
    "caisse",
    "conduite",
    "examens",
    "flotte",
    "archivage_cloud",
  ];
};

const isModuleActive = (
  school: SchoolItemSummary,
  moduleId: string,
): boolean => {
  const activeList = parseSchoolModules(school.modules_actifs);
  return activeList.includes(moduleId);
};

const toggleSchoolModule = async (
  school: SchoolItemSummary,
  moduleId: string,
) => {
  let activeList = [...parseSchoolModules(school.modules_actifs)];
  if (activeList.includes(moduleId)) {
    activeList = activeList.filter((m) => m !== moduleId);
  } else {
    activeList.push(moduleId);
  }
  try {
    await superAdminApi.updateModules(school.id, activeList);
    school.modules_actifs = JSON.stringify(activeList);
  } catch (err: any) {
    alert("Erreur lors de la mise à jour des modules : " + err?.message);
  }
};

const areAllActive = (school: SchoolItemSummary): boolean => {
  const activeList = parseSchoolModules(school.modules_actifs);
  return modulesConfig.every((m) => activeList.includes(m.id));
};

const toggleAllForSchool = async (school: SchoolItemSummary) => {
  const all = areAllActive(school);
  const nextList = all ? [] : modulesConfig.map((m) => m.id);
  try {
    await superAdminApi.updateModules(school.id, nextList);
    school.modules_actifs = JSON.stringify(nextList);
    emit("reload");
  } catch (err: any) {
    alert("Erreur lors de la mise à jour : " + err?.message);
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
          >Matrice des Droits d'Accès & Modules Système</span
        >
        <span class="text-muted-foreground font-mono"
          >({{ schools.length }} établissements)</span
        >
      </div>

      <div class="text-[11px] font-mono text-muted-foreground">
        Règle:
        <span class="text-brand font-medium"
          >Contrôle RLS & Modules Dédiés</span
        >
      </div>
    </div>

    <!-- Tableau Matriciel des Droits d'Accès -->
    <div class="flex-1 overflow-auto">
      <table class="w-full text-left text-xs border-collapse font-sans">
        <thead
          class="sticky top-0 bg-muted/80 backdrop-blur-xs border-b border-border text-[11px] text-muted-foreground font-mono uppercase tracking-wider select-none z-10"
        >
          <tr>
            <th class="py-2.5 px-3 min-w-[200px]">Établissement</th>
            <th
              v-for="mod in modulesConfig"
              :key="mod.id"
              class="py-2.5 px-3 text-center"
            >
              {{ mod.label }}
            </th>
            <th class="py-2.5 px-3 text-right">Action Rapide</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border font-mono text-foreground">
          <tr
            v-for="school in schools"
            :key="school.id"
            class="hover:bg-muted/50 transition-colors"
          >
            <!-- Auto-école -->
            <td class="py-2.5 px-3 font-sans">
              <div class="font-semibold text-foreground">
                {{ school.nom_commercial }}
              </div>
              <div class="text-[10px] text-muted-foreground font-mono">
                /{{ school.code_unique }} · {{ school.ville || "Douala" }}
              </div>
            </td>

            <!-- Colonne par module -->
            <td
              v-for="mod in modulesConfig"
              :key="mod.id"
              class="py-2.5 px-3 text-center"
            >
              <button
                type="button"
                @click="toggleSchoolModule(school, mod.id)"
                :title="
                  isModuleActive(school, mod.id)
                    ? 'Désactiver le module'
                    : 'Activer le module'
                "
                :class="[
                  'inline-flex items-center justify-center p-1.5 rounded-[4px] transition cursor-pointer',
                  isModuleActive(school, mod.id)
                    ? 'bg-brand/10 text-brand border border-brand/30 hover:bg-brand/20'
                    : 'bg-muted text-muted-foreground border border-border hover:bg-muted/80',
                ]"
              >
                <Check
                  v-if="isModuleActive(school, mod.id)"
                  class="w-3.5 h-3.5"
                />
                <Lock v-else class="w-3.5 h-3.5" />
              </button>
            </td>

            <!-- Tout activer / Tout verrouiller -->
            <td class="py-2.5 px-3 text-right">
              <button
                type="button"
                @click="toggleAllForSchool(school)"
                class="px-2 py-1 rounded-[3px] border text-[10px] transition cursor-pointer font-medium"
                :class="[
                  areAllActive(school)
                    ? 'bg-rose-500/10 hover:bg-rose-500/20 text-rose-600 dark:text-rose-400 border-rose-500/30'
                    : 'bg-muted hover:bg-brand/10 hover:text-brand text-muted-foreground border-border',
                ]"
              >
                {{
                  areAllActive(school)
                    ? "Tout verrouiller"
                    : "Tout déverrouiller"
                }}
              </button>
            </td>
          </tr>

          <tr v-if="schools.length === 0 && !isLoading">
            <td
              :colspan="modulesConfig.length + 2"
              class="text-center py-12 text-muted-foreground"
            >
              Aucune auto-école enregistrée.
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </main>
</template>

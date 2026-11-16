<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";
import type { SuperAdminGlobalUser } from "../../../types";
import { superAdminApi } from "../../../services/superAdminService";
import { Building2 } from "lucide-vue-next";

const props = defineProps<{
  searchQuery: string;
}>();

const users = ref<SuperAdminGlobalUser[]>([]);
const isLoading = ref(false);
const roleFilter = ref<string>("ALL");

const loadUsers = async () => {
  isLoading.value = true;
  try {
    users.value = await superAdminApi.getUsersData(props.searchQuery);
  } catch (err) {
    console.error("Erreur chargement utilisateurs globaux:", err);
  } finally {
    isLoading.value = false;
  }
};

watch(
  () => props.searchQuery,
  () => {
    loadUsers();
  },
);

const filteredUsers = computed(() => {
  return users.value.filter((u) => {
    if (roleFilter.value === "ALL") return true;
    return u.role.toUpperCase() === roleFilter.value;
  });
});

const getRoleBadgeClass = (role: string) => {
  switch (role.toUpperCase()) {
    case "SUPER_ADMIN":
      return "bg-purple-500/10 text-purple-600 dark:text-purple-400 border-purple-500/30";
    case "DIRECTEUR":
      return "bg-brand/10 text-brand border-brand/30";
    case "SECRETAIRE":
      return "bg-sky-500/10 text-sky-600 dark:text-sky-400 border-sky-500/30";
    case "MONITEUR":
      return "bg-amber-500/10 text-amber-600 dark:text-amber-400 border-amber-500/30";
    case "CAISSIER":
      return "bg-blue-500/10 text-blue-600 dark:text-blue-400 border-blue-500/30";
    default:
      return "bg-muted text-muted-foreground border-border";
  }
};

onMounted(() => {
  loadUsers();
});

defineExpose({
  reload: loadUsers,
});
</script>

<template>
  <main class="flex-1 flex flex-col overflow-hidden bg-background">
    <!-- Barre de Filtres Secondaire -->
    <div
      class="h-11 border-b border-border bg-card px-4 flex items-center justify-between shrink-0"
    >
      <div class="flex items-center gap-2 text-xs">
        <span class="font-semibold text-foreground">Utilisateurs Globaux</span>
        <span class="text-muted-foreground font-mono"
          >({{ filteredUsers.length }})</span
        >

        <div class="h-3 w-px bg-border mx-2"></div>

        <div
          class="flex items-center rounded-[4px] bg-muted p-0.5 border border-border"
        >
          <button
            type="button"
            @click="roleFilter = 'ALL'"
            :class="[
              'px-2 py-0.5 text-[11px] rounded-[3px] transition-colors cursor-pointer',
              roleFilter === 'ALL'
                ? 'bg-card text-foreground font-semibold shadow-xs'
                : 'text-muted-foreground hover:text-foreground',
            ]"
          >
            Tous
          </button>
          <button
            type="button"
            @click="roleFilter = 'DIRECTEUR'"
            :class="[
              'px-2 py-0.5 text-[11px] rounded-[3px] transition-colors cursor-pointer',
              roleFilter === 'DIRECTEUR'
                ? 'bg-card text-foreground font-semibold shadow-xs'
                : 'text-muted-foreground hover:text-foreground',
            ]"
          >
            Directeurs
          </button>
          <button
            type="button"
            @click="roleFilter = 'SECRETAIRE'"
            :class="[
              'px-2 py-0.5 text-[11px] rounded-[3px] transition-colors cursor-pointer',
              roleFilter === 'SECRETAIRE'
                ? 'bg-card text-foreground font-semibold shadow-xs'
                : 'text-muted-foreground hover:text-foreground',
            ]"
          >
            Secrétaires
          </button>
          <button
            type="button"
            @click="roleFilter = 'MONITEUR'"
            :class="[
              'px-2 py-0.5 text-[11px] rounded-[3px] transition-colors cursor-pointer',
              roleFilter === 'MONITEUR'
                ? 'bg-card text-foreground font-semibold shadow-xs'
                : 'text-muted-foreground hover:text-foreground',
            ]"
          >
            Moniteurs
          </button>
          <button
            type="button"
            @click="roleFilter = 'CAISSIER'"
            :class="[
              'px-2 py-0.5 text-[11px] rounded-[3px] transition-colors cursor-pointer',
              roleFilter === 'CAISSIER'
                ? 'bg-card text-foreground font-semibold shadow-xs'
                : 'text-muted-foreground hover:text-foreground',
            ]"
          >
            Caissiers
          </button>
        </div>
      </div>

      <div class="text-[11px] font-mono text-muted-foreground">
        Supervision:
        <span class="text-brand font-medium">Collaborateurs Établissements</span>
      </div>
    </div>

    <!-- Tableau Interactif Style Console Neon -->
    <div class="flex-1 overflow-auto">
      <table class="w-full text-left text-xs border-collapse">
        <thead
          class="sticky top-0 bg-muted/80 backdrop-blur-xs border-b border-border text-[11px] text-muted-foreground font-mono uppercase tracking-wider select-none z-10"
        >
          <tr>
            <th class="py-2.5 px-3">Collaborateur</th>
            <th class="py-2.5 px-3">Rôle</th>
            <th class="py-2.5 px-3">Auto-École Rattachée</th>
            <th class="py-2.5 px-3">Téléphone</th>
            <th class="py-2.5 px-3">Email</th>
            <th class="py-2.5 px-3">Statut</th>
            <th class="py-2.5 px-3 text-right">Date de création</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border font-mono text-foreground">
          <tr
            v-for="user in filteredUsers"
            :key="user.id"
            class="hover:bg-muted/50 transition-colors"
          >
            <!-- Nom & Prénoms avec Initiale -->
            <td class="py-2.5 px-3 font-sans">
              <div class="flex items-center gap-2.5">
                <div
                  class="w-6 h-6 rounded-full bg-muted border border-border flex items-center justify-center text-[10px] font-mono font-bold text-muted-foreground shrink-0"
                >
                  {{
                    (user.prenoms?.[0] || user.nom?.[0] || "U").toUpperCase()
                  }}
                </div>
                <div>
                  <div class="font-semibold text-foreground">
                    {{ user.nom }} {{ user.prenoms }}
                  </div>
                  <div class="text-[10px] text-muted-foreground font-mono">
                    ID: {{ user.id.slice(0, 8) }}...
                  </div>
                </div>
              </div>
            </td>

            <!-- Rôle avec badge -->
            <td class="py-2.5 px-3">
              <span
                :class="[
                  'inline-flex items-center gap-1 px-1.5 py-0.5 rounded-[3px] text-[10px] font-semibold border',
                  getRoleBadgeClass(user.role),
                ]"
              >
                {{ user.role }}
              </span>
            </td>

            <!-- Auto-école -->
            <td class="py-2.5 px-3 font-sans text-foreground">
              <div class="flex items-center gap-1.5">
                <Building2 class="w-3.5 h-3.5 text-muted-foreground shrink-0" />
                <span class="truncate max-w-[180px]">{{
                  user.auto_ecole_nom
                }}</span>
              </div>
            </td>

            <!-- Téléphone -->
            <td class="py-2.5 px-3 text-muted-foreground">
              {{ user.telephone }}
            </td>

            <!-- Email -->
            <td class="py-2.5 px-3 text-muted-foreground">
              {{ user.email || "—" }}
            </td>

            <!-- Statut -->
            <td class="py-2.5 px-3">
              <span
                :class="[
                  'inline-flex items-center gap-1 px-1.5 py-0.5 rounded-[3px] text-[10px] font-semibold',
                  user.statut === 'ACTIF'
                    ? 'bg-brand/10 text-brand border border-brand/30'
                    : 'bg-destructive/10 text-destructive border border-destructive/30',
                ]"
              >
                <span
                  class="w-1.5 h-1.5 rounded-full"
                  :class="
                    user.statut === 'ACTIF' ? 'bg-brand' : 'bg-destructive'
                  "
                ></span>
                {{ user.statut }}
              </span>
            </td>

            <!-- Date -->
            <td class="py-2.5 px-3 text-right text-muted-foreground">
              {{ new Date(user.created_at).toLocaleDateString("fr-FR") }}
            </td>
          </tr>

          <tr v-if="filteredUsers.length === 0 && !isLoading">
            <td colspan="7" class="text-center py-12 text-muted-foreground font-sans">
              Aucun utilisateur trouvé pour cette sélection.
            </td>
          </tr>
          <tr v-if="isLoading && filteredUsers.length === 0">
            <td
              colspan="7"
              class="text-center py-12 text-muted-foreground font-sans font-mono animate-pulse"
            >
              Chargement des utilisateurs globaux...
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </main>
</template>

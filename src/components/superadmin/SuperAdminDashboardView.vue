<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import { useAuthStore } from "../../stores/authStore";
import { superAdminApi } from "../../services/superAdminService";
import { nativeNotificationEngine } from "../../services/nativeNotificationService";
import { authBackend } from "../../services/authBackendService";
import type { SuperAdminStats, SchoolItemSummary } from "../../types";

// Composants factorisés
import SuperAdminHeader from "./SuperAdminHeader.vue";
import SuperAdminSidebar, { type SuperAdminTab } from "./SuperAdminSidebar.vue";
import SuperAdminSchoolsView from "./views/SuperAdminSchoolsView.vue";
import SuperAdminLicencesView from "./views/SuperAdminLicencesView.vue";
import SuperAdminPermissionsView from "./views/SuperAdminPermissionsView.vue";
import SuperAdminUsersView from "./views/SuperAdminUsersView.vue";
import SuperAdminCandidatsView from "./views/SuperAdminCandidatsView.vue";
import SuperAdminAuditView from "./views/SuperAdminAuditView.vue";
import SuperAdminAlertsView from "./views/SuperAdminAlertsView.vue";
import SuperAdminAnnouncementsView from "./views/SuperAdminAnnouncementsView.vue";
import SuperAdminFeedbacksView from "./views/SuperAdminFeedbacksView.vue";

const authStore = useAuthStore();

// État de navigation
const activeTab = ref<SuperAdminTab>("schools");
const searchQuery = ref("");
const isLoading = ref(true);
let unsubscribeWs: (() => void) | null = null;

const usersViewRef = ref<InstanceType<typeof SuperAdminUsersView> | null>(null);
const candidatsViewRef = ref<InstanceType<
  typeof SuperAdminCandidatsView
> | null>(null);
const auditViewRef = ref<InstanceType<typeof SuperAdminAuditView> | null>(null);

// Statistiques d'infrastructure
const stats = ref<SuperAdminStats>({
  total_schools: 0,
  active_schools: 0,
  total_users: 0,
  total_candidats: 0,
  total_volume_paiement: 0,
  cloud_status: "ONLINE",
  cluster_latency_ms: 24,
});

// Données des auto-écoles
const schools = ref<SchoolItemSummary[]>([]);

const loadData = async () => {
  isLoading.value = true;
  try {
    const data = await superAdminApi.getDashboardData(searchQuery.value);
    stats.value = data.stats;
    schools.value = data.schools || [];
  } catch (err: any) {
    console.error("Erreur chargement super-admin:", err);
  } finally {
    isLoading.value = false;
  }
};

const refreshCurrentView = async () => {
  await loadData();
  if (activeTab.value === "users" && usersViewRef.value) {
    usersViewRef.value.reload();
  } else if (activeTab.value === "candidats" && candidatsViewRef.value) {
    candidatsViewRef.value.reload();
  } else if (activeTab.value === "audit" && auditViewRef.value) {
    auditViewRef.value.reload();
  }
};

const handleToggleStatus = async (school: SchoolItemSummary) => {
  const newStatus = !school.est_actif;
  try {
    await superAdminApi.toggleSchoolStatus(school.id, newStatus);
    school.est_actif = newStatus;
  } catch (err: any) {
    alert("Erreur lors de la modification du statut : " + err?.message);
  }
};

const handleLogout = () => {
  authStore.logoutWeb();
  authStore.webPortalView = "super_admin_login";
  if (
    typeof window !== "undefined" &&
    window.location.pathname.startsWith("/super-admin")
  ) {
    window.history.pushState({}, "", "/super-admin/login");
  }
};

onMounted(() => {
  const session = authBackend.getActiveSession();
  if (session?.token) {
    nativeNotificationEngine.connect(session.token);
  }
  loadData();

  // Écoute temps réel des événements Super-Admin (ex: nouvelle auto-école inscrite, modification licence ou statut)
  unsubscribeWs = nativeNotificationEngine.subscribe((event) => {
    if (
      event.type === "SCHOOL_REGISTERED" ||
      event.type === "SUPER_ADMIN_REFRESH"
    ) {
      loadData();
    }
  });
});

onUnmounted(() => {
  if (unsubscribeWs) {
    unsubscribeWs();
    unsubscribeWs = null;
  }
});
</script>

<template>
  <div
    class="flex-1 w-full flex flex-col bg-background text-foreground font-sans select-none antialiased overflow-hidden"
  >
    <!-- 1. En-tête Principal Topbar Unique (Recherche, Télémétrie, Actions) -->
    <SuperAdminHeader
      :search-query="searchQuery"
      :is-loading="isLoading"
      :latency-ms="stats.cluster_latency_ms"
      :active-tab="activeTab"
      @update:search-query="searchQuery = $event"
      @refresh="refreshCurrentView"
      @logout="handleLogout"
    />

    <!-- 2. Corps Principal : Sidebar Unique + Connexion dynamique des Vues -->
    <div class="flex-1 flex overflow-hidden">
      <!-- Composant Unique : Sidebar Latérale de Navigation -->
      <SuperAdminSidebar
        :active-tab="activeTab"
        :stats="stats"
        @select-tab="activeTab = $event"
      />

      <!-- Contenu Principal des Vues Système -->
      <main class="flex-1 flex flex-col overflow-hidden bg-card px-4 pb-2 pt-2">
        <div
          class="flex-1 border border-border rounded-[4px] overflow-y-auto bg-background p-4 sm:p-5"
        >
          <!-- Onglet 1 : Auto-Écoles & Flotte (Édition, Suspension, Suppression) -->
          <SuperAdminSchoolsView
            v-if="activeTab === 'schools'"
            :schools="schools"
            :is-loading="isLoading"
            @toggle-status="handleToggleStatus"
            @reload="loadData"
          />

          <!-- Onglet 2 : Quotas Machines & Licences (Jours, Mois, Années, Machines) -->
          <SuperAdminLicencesView
            v-else-if="activeTab === 'licences'"
            :schools="schools"
            :is-loading="isLoading"
            @reload="loadData"
          />

          <!-- Onglet 3 : Matrice des Droits d'Accès & Modules Dédiés -->
          <SuperAdminPermissionsView
            v-else-if="activeTab === 'permissions'"
            :schools="schools"
            :is-loading="isLoading"
            @reload="loadData"
          />

          <!-- Onglet 4 : Collaborateurs & Utilisateurs Globaux -->
          <SuperAdminUsersView
            v-else-if="activeTab === 'users'"
            ref="usersViewRef"
            :search-query="searchQuery"
          />

          <!-- Onglet 5 : Candidats Inscrits & Parcours Permis -->
          <SuperAdminCandidatsView
            v-else-if="activeTab === 'candidats'"
            ref="candidatsViewRef"
            :search-query="searchQuery"
          />

          <!-- Onglet 6 : Journal d'Audit & Sécurité Système -->
          <SuperAdminAuditView
            v-else-if="activeTab === 'audit'"
            ref="auditViewRef"
          />

          <!-- Onglet 7 : Diffusion d'Alertes Réseau -->
          <SuperAdminAlertsView
            v-else-if="activeTab === 'alerts'"
            :schools="schools"
            :search-query="searchQuery"
          />

          <!-- Onglet 8 : Annonces & Événements Réseau -->
          <SuperAdminAnnouncementsView
            v-else-if="activeTab === 'announcements'"
            :schools="schools"
            :search-query="searchQuery"
          />

          <!-- Onglet 9 : Avis & Satisfaction des Auto-Écoles -->
          <SuperAdminFeedbacksView
            v-else-if="activeTab === 'feedbacks'"
            :schools="schools"
            :search-query="searchQuery"
          />
        </div>
      </main>
    </div>
  </div>
</template>

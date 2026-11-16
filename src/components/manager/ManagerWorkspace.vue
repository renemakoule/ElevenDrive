<script setup lang="ts">
import { onMounted, onUnmounted } from "vue";
import { useAuthStore } from "../../stores/authStore";
import { useGuichetStore } from "../../stores/guichetStore";
import { realtimeSocket } from "../../services/realtimeSocketService";
import ManagerHomeView from "./ManagerHomeView.vue";
import ManagerFinanceView from "./ManagerFinanceView.vue";
import ManagerFlotteView from "./ManagerFlotteView.vue";
import ManagerEquipeView from "./ManagerEquipeView.vue";
import GerantDashboardView from "../guichet/GerantDashboardView.vue";
import FormulesTarifsView from "../guichet/FormulesTarifsView.vue";
import EpreuvesManagementView from "../epreuves/EpreuvesManagementView.vue";
import CoursCodePlanningView from "../cours/CoursCodePlanningView.vue";
import PlanningConduiteView from "../cours/PlanningConduiteView.vue";
import SupportsCoursManagementView from "../supports/SupportsCoursManagementView.vue";
import TutorielsManagementView from "../tutoriels/TutorielsManagementView.vue";
import ExamensView from "../examens/ExamensView.vue";
import RoleHistoriqueView from "../historique/RoleHistoriqueView.vue";
import CentreArchivageView from "../archivage/CentreArchivageView.vue";
import LandingPricing from "../landing/LandingPricing.vue";

const authStore = useAuthStore();
const guichetStore = useGuichetStore();
let unsubscribeWs: (() => void) | null = null;

onMounted(() => {
  realtimeSocket.connect();
  unsubscribeWs = realtimeSocket.subscribe(async (msg) => {
    switch (msg.type) {
      case "PAIEMENT_ENREGISTRE":
      case "CANDIDAT_DOSSIER_UPDATED":
      case "CAISSE_SESSION_UPDATED":
      case "EXAMENS_SESSION_UPDATED":
      case "QUIZ_RESULTAT_UPDATED":
      case "LIVRET_EMARGEMENT_UPDATED":
        await Promise.all([
          guichetStore.loadGerantDashboard(),
          guichetStore.loadStats(),
          guichetStore.search(guichetStore.searchQuery || ""),
        ]);
        break;
      case "FORMULES_TARIFS_UPDATED":
        await Promise.all([
          guichetStore.loadFormules(),
          guichetStore.loadGerantDashboard(),
        ]);
        break;
      case "FLOTTE_VEHICULES_UPDATED":
        await guichetStore.loadVehicules();
        break;
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
  <div class="flex-1 overflow-y-auto">
    <!-- 1. Vue d'ensemble (style MaryamE-logistic) -->
    <template v-if="authStore.activeSection === 'overview'">
      <ManagerHomeView />
    </template>

    <!-- 1a0. Planning Conduite Pratique (Moniteurs & Flotte) -->
    <template v-else-if="authStore.activeSection === 'conduite-planning'">
      <PlanningConduiteView />
    </template>

    <!-- 1a. Planning Cours de Code Théorique & Visio -->
    <template v-else-if="authStore.activeSection === 'cours-planning'">
      <CoursCodePlanningView />
    </template>

    <!-- 1a2. Supports de Cours & Manuels PDF -->
    <template v-else-if="authStore.activeSection === 'supports-cours'">
      <SupportsCoursManagementView />
    </template>

    <!-- 1a3. Bibliothèque de Tutoriels & Guides Vidéo -->
    <template v-else-if="authStore.activeSection === 'tutoriels'">
      <TutorielsManagementView />
    </template>

    <!-- 1b. Catalogue des Tarifs Permis (Directeur Habilité à Modifier les Prix) -->
    <template v-else-if="authStore.activeSection === 'tarifs'">
      <FormulesTarifsView />
    </template>

    <!-- 1c. Banque d'Épreuves & Annales PDF (Pédagogie) -->
    <template v-else-if="authStore.activeSection === 'epreuves'">
      <EpreuvesManagementView />
    </template>

    <!-- 1d. Équipe & Collaborateurs (Directeur) -->
    <template v-else-if="authStore.activeSection === 'equipe'">
      <ManagerEquipeView />
    </template>

    <!-- 2. Audit & Journal Financier -->
    <template v-else-if="authStore.activeSection === 'finance'">
      <ManagerFinanceView />
    </template>

    <!-- 3. Flotte Automobile & Alertes -->
    <template v-else-if="authStore.activeSection === 'flotte'">
      <ManagerFlotteView />
    </template>

    <!-- 4. Taux de Réussite & Examens d'État Officiels (Bordereaux DGTTC, PV d'admis) -->
    <template v-else-if="authStore.activeSection === 'stats'">
      <ExamensView />
    </template>

    <!-- 5. Dérogations Manager (PIN) & Gestion des Impayés -->
    <template v-else-if="authStore.activeSection === 'derogations'">
      <GerantDashboardView />
    </template>

    <!-- 6. Historique & Journal Audit Global (Directeur) -->
    <template v-else-if="authStore.activeSection === 'historique'">
      <RoleHistoriqueView />
    </template>

    <!-- 7. Centre d'Archivage & Rétention Documentaire (Directeur) -->
    <template v-else-if="authStore.activeSection === 'archivage'">
      <CentreArchivageView />
    </template>

    <!-- 8. Bon Plan / Grille Tarifaire ElevenDrive (Directeur uniquement) -->
    <template v-else-if="authStore.activeSection === 'bon-plan'">
      <LandingPricing />
    </template>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue';
import { useAuthStore } from '../../stores/authStore';
import { useGuichetStore } from '../../stores/guichetStore';
import { useCoursCodeStore } from '../../stores/coursCodeStore';
import { useSupportsCoursStore } from '../../stores/supportsCoursStore';
import { useEpreuvesStore } from '../../stores/epreuvesStore';
import { useQuizFormsStore } from '../../stores/quizFormsStore';
import { realtimeSocket } from '../../services/realtimeSocketService';
import MoniteurPwaView from './MoniteurPwaView.vue';
import CoursCodePlanningView from '../cours/CoursCodePlanningView.vue';
import PlanningConduiteView from '../cours/PlanningConduiteView.vue';
import SupportsCoursManagementView from '../supports/SupportsCoursManagementView.vue';
import EpreuvesManagementView from '../epreuves/EpreuvesManagementView.vue';
import MoniteurQuizFormsView from '../quiz/MoniteurQuizFormsView.vue';
import RoleHistoriqueView from '../historique/RoleHistoriqueView.vue';

const authStore = useAuthStore();
const guichetStore = useGuichetStore();
const coursStore = useCoursCodeStore();
const supportsStore = useSupportsCoursStore();
const epreuvesStore = useEpreuvesStore();
const quizStore = useQuizFormsStore();

let unsubscribeWs: (() => void) | null = null;

onMounted(() => {
  realtimeSocket.connect();
  coursStore.loadSeances();
  supportsStore.loadSupports();
  epreuvesStore.loadEpreuves();
  guichetStore.loadVehicules();
  quizStore.loadAll();

  unsubscribeWs = realtimeSocket.subscribe(async (msg) => {
    switch (msg.type) {
      case 'COURS_CODE_UPDATED':
        await coursStore.loadSeances();
        break;
      case 'SUPPORTS_COURS_UPDATED':
        await supportsStore.loadSupports();
        break;
      case 'EPREUVES_UPDATED':
        await epreuvesStore.loadEpreuves();
        break;
      case 'FLOTTE_VEHICULES_UPDATED':
        await guichetStore.loadVehicules();
        break;
      case 'QUIZ_FORM_UPDATED':
        await quizStore.loadFormulaires();
        break;
      case 'QUIZ_SUBMISSION_RECEIVED':
        if (msg.data?.id) {
          quizStore.registerIncomingSubmission(msg.data);
          guichetStore.showToast(
            `Nouvelle réponse au Quiz de ${msg.data.candidat_nom} (${msg.data.pourcentage}%)`,
            'success'
          );
        } else {
          await quizStore.loadSoumissions();
        }
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
    <template v-if="authStore.activeSection === 'supports-cours'">
      <SupportsCoursManagementView />
    </template>
    <template v-else-if="authStore.activeSection === 'epreuves'">
      <EpreuvesManagementView />
    </template>
    <template v-else-if="authStore.activeSection === 'quiz-forms'">
      <MoniteurQuizFormsView />
    </template>
    <template v-else-if="authStore.activeSection === 'cours-planning'">
      <CoursCodePlanningView />
    </template>
    <template v-else-if="authStore.activeSection === 'conduite-planning'">
      <PlanningConduiteView />
    </template>
    <template v-else-if="authStore.activeSection === 'historique'">
      <RoleHistoriqueView />
    </template>
    <template v-else>
      <MoniteurPwaView />
    </template>
  </div>
</template>

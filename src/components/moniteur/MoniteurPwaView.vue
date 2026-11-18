<script setup lang="ts">
import { ref, computed } from "vue";
import { useGuichetStore } from "../../stores/guichetStore";
import { useAuthStore } from "../../stores/authStore";
import { useCoursCodeStore } from "../../stores/coursCodeStore";
import { useSupportsCoursStore } from "../../stores/supportsCoursStore";
import { useEpreuvesStore } from "../../stores/epreuvesStore";
import { useQuizFormsStore } from "../../stores/quizFormsStore";
import { api } from "../../db/api";
import type { SeanceConduite, CompetenceItem, SeanceCode } from "../../types";
import ConduiteSessionModal from "./ConduiteSessionModal.vue";
import ClasseVirtuelleModal from "../cours/ClasseVirtuelleModal.vue";
import MoniteurTopToolbar from "./MoniteurTopToolbar.vue";
import MoniteurLivretTactileView from "./MoniteurLivretTactileView.vue";
import MoniteurVehiculeAssigneView from "./MoniteurVehiculeAssigneView.vue";
import MoniteurCalendarBar from "./MoniteurCalendarBar.vue";
import MoniteurPlanningListSection from "./MoniteurPlanningListSection.vue";
import MoniteurSidebarSummary from "./MoniteurSidebarSummary.vue";
import MoniteurDashboardSkeleton from "../common/skeleton/MoniteurDashboardSkeleton.vue";
import { useMoniteurPlanning } from "./useMoniteurPlanning";

const store = useGuichetStore();
const authStore = useAuthStore();
const coursStore = useCoursCodeStore();
const supportsStore = useSupportsCoursStore();
const epreuvesStore = useEpreuvesStore();
const quizStore = useQuizFormsStore();

const activeSession = ref<SeanceConduite | null>(null);
const isSessionModalOpen = ref(false);
const vehiculeChoisiId = ref<string>("");
const seanceCodeVisioActive = ref<SeanceCode | null>(null);
const isUpdatingVehiculeStatut = ref(false);

const seanceCodeEnLigneAujourdhui = computed(
  () =>
    coursStore.seancesAujourdhui.find((s) => s.type_seance === "EN_LIGNE") ||
    coursStore.prochaineSeanceEnLigne ||
    coursStore.seancesAujourdhui[0] ||
    coursStore.seances[0] ||
    null,
);

const openClasseVirtuelle = (s: SeanceCode) => {
  seanceCodeVisioActive.value = s;
};

const activeMoniteur = computed(() => {
  const u = authStore.currentUser;
  const nomComplet =
    `${u?.prenom || ""} ${u?.nom || ""}`.trim() || "Moniteur Principal";
  return {
    id: u?.id || "",
    nom: nomComplet,
    titre: u?.roleLabel || "Moniteur titulaire · Enseignement pratique",
    telephone: u?.telephone || "+225 01 02 03 04 05",
  };
});

const {
  isWsConnected,
  rawSessions,
  isLoading,
  todayIso,
  selectedDate,
  viewMode,
  weekSessionsMap,
  weekDays,
  dateAfficheeLongue,
  loadData,
  selectDate,
  goToPrevDay,
  goToNextDay,
  resetToToday,
  handleSelectDayFromWeek,
} = useMoniteurPlanning(() => activeMoniteur.value.id);

const vehiculesDisponibles = computed(() => store.vehicules || []);

const vehiculeAssigne = computed(() => {
  if (vehiculeChoisiId.value) {
    const v = vehiculesDisponibles.value.find(
      (x) => x.id === vehiculeChoisiId.value,
    );
    if (v) return v;
  }
  return vehiculesDisponibles.value[0] || null;
});

const vehiculeInfo = computed(() => {
  const v = vehiculeAssigne.value;
  if (!v) {
    return {
      immatriculation: "Non assigné",
      modele: "Aucun véhicule",
      visiteTechnique: "—",
      assurance: "—",
      kilometrage: 0,
    };
  }
  const vt = v.date_visite_technique
    ? new Date(v.date_visite_technique).toLocaleDateString("fr-FR", {
        day: "numeric",
        month: "long",
        year: "numeric",
      })
    : "Non renseignée";
  const ass = v.date_expiration_assurance
    ? `Expire le ${new Date(v.date_expiration_assurance).toLocaleDateString("fr-FR")}`
    : "À jour";

  return {
    id: v.id,
    immatriculation: v.immatriculation,
    modele: `${v.marque_modele} (${v.type_boite === "AUTOMATIQUE" ? "BVA" : "BVM"})`,
    visiteTechnique: vt,
    assurance: ass,
    kilometrage: v.kilometrage_actuel || 0,
  };
});

const sessions = computed<SeanceConduite[]>(() => rawSessions.value);
const dateLongue = computed(() => dateAfficheeLongue.value);

const pct = (value: number, max: number) =>
  max > 0 ? Math.min(100, Math.round((value / max) * 100)) : 0;
const num = (n: number) => n.toLocaleString("fr-FR");
const formatDuree = (min: number) => {
  const h = Math.floor(min / 60);
  const m = min % 60;
  return m ? `${h} h ${String(m).padStart(2, "0")}` : `${h} h`;
};

const initiales = computed(
  () =>
    activeMoniteur.value.nom
      .split(" ")
      .filter(Boolean)
      .slice(0, 2)
      .map((p) => p.charAt(0))
      .join("")
      .toUpperCase() || "M",
);

const nbTerminees = computed(
  () => sessions.value.filter((s) => s.statut === "EFFECTUE").length,
);
const nbEnCours = computed(
  () => sessions.value.filter((s) => s.statut === "EN_COURS").length,
);
const nbAVenir = computed(
  () => sessions.value.filter((s) => s.statut === "PROGRAMME").length,
);

const totalMinutes = computed(() =>
  sessions.value.reduce((sum, s) => sum + s.duree_minutes, 0),
);
const minutesFaites = computed(() =>
  sessions.value
    .filter((s) => s.statut === "EFFECTUE")
    .reduce((sum, s) => sum + s.duree_minutes, 0),
);

const compteur = computed(() => vehiculeInfo.value.kilometrage || 0);
const enAttenteSync = computed(
  () => store.stats?.evenements_sync_en_attente ?? 0,
);

const stats = computed(() => [
  {
    key: "vehicule",
    label: "Véhicule assigné",
    value: vehiculeInfo.value.modele,
    tag: vehiculeInfo.value.immatriculation,
    caption: `Compteur : ${num(compteur.value)} km`,
  },
  {
    key: "heures",
    label: "Heures du jour",
    value: formatDuree(totalMinutes.value),
    progress: pct(minutesFaites.value, totalMinutes.value),
    caption: `${formatDuree(minutesFaites.value)} effectuées · ${formatDuree(totalMinutes.value - minutesFaites.value)} restantes`,
  },
  {
    key: "emargement",
    label: "Leçons émargées",
    value: `${nbTerminees.value} / ${sessions.value.length}`,
    progress: pct(nbTerminees.value, sessions.value.length),
    caption: `${nbEnCours.value} en cours · ${nbAVenir.value} à venir`,
  },
]);

const openSession = (session: SeanceConduite) => {
  activeSession.value = session;
  isSessionModalOpen.value = true;
};

const handleSaveSession = async (payload: {
  sessionId: string;
  dossierId: string;
  dureeHeures: number;
  kmDebut: number;
  kmFin: number;
  commentaire: string;
  competences: CompetenceItem[];
  signatureEleve: string;
}) => {
  try {
    await api.validerSeanceConduite({
      ...payload,
      vehiculeId: vehiculeAssigne.value?.id,
    });

    const target = rawSessions.value.find(
      (s) => s.id === payload.sessionId || s.dossier_id === payload.dossierId,
    );
    if (target) {
      target.statut = "EFFECTUE";
      target.km_debut = payload.kmDebut;
      target.km_fin = payload.kmFin;
      target.commentaire_moniteur = payload.commentaire;
      target.competences = payload.competences;
      target.heures_deja_effectuees += payload.dureeHeures;
      target.signature_eleve_url = payload.signatureEleve;
    }

    isSessionModalOpen.value = false;
    await store.loadVehicules();
    store.showToast(
      `Leçon validée : ${payload.dureeHeures ?? 2} h décomptées du forfait !`,
      "success",
    );
  } catch (err: any) {
    store.showToast("Erreur validation leçon: " + err.message, "error");
  }
};

const handleQuickValidate = async (session: SeanceConduite) => {
  const kmDeb = session.km_debut ?? vehiculeInfo.value.kilometrage ?? 0;
  const kmF =
    session.km_fin && session.km_fin > kmDeb ? session.km_fin : kmDeb + 15;
  const duree = Math.max(1, Math.round((session.duree_minutes || 120) / 60));

  await handleSaveSession({
    sessionId: session.id,
    dossierId: session.dossier_id,
    dureeHeures: duree,
    kmDebut: kmDeb,
    kmFin: kmF,
    commentaire:
      session.commentaire_moniteur ||
      "Présence validée en 1 clic par le moniteur",
    competences:
      session.competences && session.competences.length > 0
        ? session.competences
        : [],
    signatureEleve: session.signature_eleve_url || "PRESENCE_VALIDEE_1_CLIC",
  });
};

const handleChangerStatutVehicule = async (
  vehiculeId: string,
  statut: "DISPONIBLE" | "EN_MAINTENANCE" | "EN_PANNE",
) => {
  if (!vehiculeId) return;
  isUpdatingVehiculeStatut.value = true;
  try {
    await api.updateVehiculeStatut(vehiculeId, statut);
    await store.loadVehicules();
    store.showToast(
      `Statut du véhicule mis à jour en temps réel : ${statut}`,
      "success",
    );
  } catch (err: any) {
    store.showToast("Erreur mise à jour véhicule : " + err.message, "error");
  } finally {
    isUpdatingVehiculeStatut.value = false;
  }
};
</script>

<template>
  <div class="space-y-4">
    <MoniteurTopToolbar
      :active-section="authStore.activeSection"
      :cours-count="coursStore.seances.length"
      :supports-count="supportsStore.supportsList.length"
      :epreuves-count="epreuvesStore.epreuvesList.length"
      :quiz-soumissions-count="quizStore.soumissions.length"
      :is-loading="isLoading"
      :en-attente-sync="enAttenteSync"
      :seance-code-en-ligne-aujourdhui="seanceCodeEnLigneAujourdhui"
      @navigate="(section) => authStore.setActiveSection(section as any)"
      @refresh="loadData"
      @open-classe-virtuelle="openClasseVirtuelle"
    />

    <MoniteurLivretTactileView
      v-if="authStore.activeSection === 'livret'"
      :sessions="sessions"
      :selected-date="selectedDate"
      :date-longue="dateLongue"
      :default-kilometrage="vehiculeInfo.kilometrage"
      :is-ws-connected="isWsConnected"
      @prev-day="goToPrevDay"
      @next-day="goToNextDay"
      @select-date="selectDate"
      @save-session="handleSaveSession"
    />

    <MoniteurVehiculeAssigneView
      v-else-if="authStore.activeSection === 'vehicule'"
      :vehicules-disponibles="vehiculesDisponibles"
      :vehicule-assigne="vehiculeAssigne"
      :vehicule-info="vehiculeInfo"
      v-model:vehicule-choisi-id="vehiculeChoisiId"
      :is-updating-statut="isUpdatingVehiculeStatut"
      @change-statut="handleChangerStatutVehicule"
    />

    <!-- Skeleton Planning Moniteur -->
    <MoniteurDashboardSkeleton v-else-if="isLoading && rawSessions.length === 0" />

    <template v-else>
      <MoniteurCalendarBar
        v-model:view-mode="viewMode"
        :selected-date="selectedDate"
        :today-iso="todayIso"
        :date-longue="dateLongue"
        :week-days="weekDays"
        :stats="stats"
        @prev-day="goToPrevDay"
        @next-day="goToNextDay"
        @reset-today="resetToToday"
        @select-date="selectDate"
      />

      <div class="grid grid-cols-1 items-start gap-4 lg:grid-cols-12">
        <MoniteurPlanningListSection
          :view-mode="viewMode"
          :sessions="sessions"
          :selected-date="selectedDate"
          :today-iso="todayIso"
          :date-affichee-longue="dateAfficheeLongue"
          :default-kilometrage="vehiculeInfo.kilometrage"
          :week-days="weekDays"
          :week-sessions-map="weekSessionsMap"
          @open-session="openSession"
          @quick-validate="handleQuickValidate"
          @reset-today="resetToToday"
          @select-day-from-week="handleSelectDayFromWeek"
        />

        <MoniteurSidebarSummary
          :active-moniteur="activeMoniteur"
          :initiales="initiales"
          :vehicules-disponibles="vehiculesDisponibles"
          v-model:vehicule-choisi-id="vehiculeChoisiId"
          :vehicule-info="vehiculeInfo"
        />
      </div>
    </template>

    <ConduiteSessionModal
      v-if="isSessionModalOpen && activeSession"
      :session="activeSession"
      @close="isSessionModalOpen = false"
      @save="handleSaveSession"
    />

    <ClasseVirtuelleModal
      :seance="seanceCodeVisioActive"
      @close="seanceCodeVisioActive = null"
    />
  </div>
</template>

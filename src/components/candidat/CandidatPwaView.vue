<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from "vue";
import { useGuichetStore } from "../../stores/guichetStore";
import { useAuthStore } from "../../stores/authStore";
import { useEpreuvesStore } from "../../stores/epreuvesStore";
import { useCoursCodeStore } from "../../stores/coursCodeStore";
import { useSupportsCoursStore } from "../../stores/supportsCoursStore";
import { useTutorielsStore } from "../../stores/tutorielsStore";
import { useOfflineDownloadsStore } from "../../stores/offlineDownloadsStore";
import { useQuizFormsStore } from "../../stores/quizFormsStore";
import { realtimeSocket } from "../../services/realtimeSocketService";
import type { SeanceConduite } from "../../types";
import { api } from "../../db/api";
import CodeQuizModal from "./CodeQuizModal.vue";
import ReservationConduiteModal from "./ReservationConduiteModal.vue";
import EpreuveViewerModal from "../epreuves/EpreuveViewerModal.vue";
import CandidatCoursEnLigneView from "../cours/CandidatCoursEnLigneView.vue";
import CandidatSupportsCoursView from "../supports/CandidatSupportsCoursView.vue";
import CandidatTutorielsView from "../tutoriels/CandidatTutorielsView.vue";
import CandidatQuizFormsView from "../quiz/CandidatQuizFormsView.vue";
import OfflineDownloadsView from "./OfflineDownloadsView.vue";
import RoleHistoriqueView from "../historique/RoleHistoriqueView.vue";
import CandidatDashboardTab, {
  type CandidatProfileSummary,
  type LiveLivretEval,
} from "./tabs/CandidatDashboardTab.vue";
import CandidatEpreuvesTab from "./tabs/CandidatEpreuvesTab.vue";
import CandidatNavTabs, {
  type CandidatMainTab,
} from "./tabs/CandidatNavTabs.vue";
import CandidatHeaderBar from "./tabs/CandidatHeaderBar.vue";
import { FileText } from "lucide-vue-next";

const store = useGuichetStore();
const authStore = useAuthStore();
const epreuvesStore = useEpreuvesStore();
const coursStore = useCoursCodeStore();
const supportsStore = useSupportsCoursStore();
const tutorielsStore = useTutorielsStore();
const offlineStore = useOfflineDownloadsStore();
const quizStore = useQuizFormsStore();

const activeMainTab = ref<CandidatMainTab>("dashboard");
const selectedCandId = ref<string>("");
const isQuizModalOpen = ref(false);
const isReservationModalOpen = ref(false);
const lastQuizScore = ref<number | null>(null);

const isWsConnected = realtimeSocket.isConnected;
const lastRealtimeSyncAt = ref<string>("");
const seancesConduiteCandidat = ref<SeanceConduite[]>([]);
const liveLivretEvaluation = ref<LiveLivretEval | null>(null);
let unsubscribeWs: (() => void) | null = null;

const SEUIL_EXAMEN = 35;

const loadCandidatDrivingSessions = async () => {
  try {
    const allSessions = await api.getPlanningConduite("ALL");
    const dosId = store.activeDossierDetail?.dossier?.id;
    const candNom = store.activeDossierDetail
      ? `${store.activeDossierDetail.candidat.nom} ${store.activeDossierDetail.candidat.prenoms}`
          .trim()
          .toLowerCase()
      : "";
    if (dosId) {
      const matched = (allSessions || []).filter(
        (s) =>
          s.dossier_id === dosId ||
          (candNom && s.candidat_nom.toLowerCase().includes(candNom)),
      );
      seancesConduiteCandidat.value =
        matched.length > 0 ? matched : allSessions || [];
    } else {
      seancesConduiteCandidat.value = allSessions || [];
    }
  } catch {
    // ignore
  }
};

onMounted(async () => {
  realtimeSocket.connect();
  offlineStore.initListeners();
  offlineStore.loadDownloads();
  await Promise.all([
    epreuvesStore.loadEpreuves(),
    coursStore.loadSeances(),
    supportsStore.loadSupports(),
    tutorielsStore.loadTutoriels(),
    quizStore.loadAll(),
  ]);
  if (!store.searchResults || store.searchResults.length === 0) {
    await store.search("");
  }
  if (!store.vehicules || store.vehicules.length === 0) {
    await store.loadVehicules();
  }
  if (store.searchResults.length > 0) {
    const currentUser = authStore.currentUser;
    let initialId = store.selectedCandidatId || store.searchResults[0].id;
    if (currentUser && authStore.currentRole === "eleve") {
      const matched = store.searchResults.find(
        (c) =>
          (currentUser.telephone &&
            c.telephone_whatsapp === currentUser.telephone) ||
          (`${c.nom} ${c.prenoms}`
            .toLowerCase()
            .includes((currentUser.nom || "").toLowerCase()) &&
            `${c.nom} ${c.prenoms}`
              .toLowerCase()
              .includes((currentUser.prenom || "").toLowerCase())),
      );
      if (matched) initialId = matched.id;
    }
    selectedCandId.value = initialId;
    if (
      !store.activeDossierDetail ||
      store.activeDossierDetail.candidat.id !== initialId
    ) {
      await store.selectCandidat(initialId);
    }
  }
  await loadCandidatDrivingSessions();

  unsubscribeWs = realtimeSocket.subscribe(async (msg) => {
    lastRealtimeSyncAt.value = new Date().toLocaleTimeString("fr-FR");
    const currentCandId =
      selectedCandId.value || store.activeDossierDetail?.candidat?.id;
    const currentDossierId = store.activeDossierDetail?.dossier?.id;

    switch (msg.type) {
      case "PAIEMENT_ENREGISTRE":
      case "CANDIDAT_DOSSIER_UPDATED":
      case "QUIZ_RESULTAT_UPDATED":
        await store.search("");
        if (currentCandId) await store.selectCandidat(currentCandId);
        break;
      case "PLANNING_CONDUITE_UPDATED":
      case "LIVRET_EMARGEMENT_UPDATED":
        if (currentCandId) await store.selectCandidat(currentCandId);
        await loadCandidatDrivingSessions();
        if (msg.type === "LIVRET_EMARGEMENT_UPDATED") {
          liveLivretEvaluation.value = null;
        }
        break;
      case "LIVRET_EMARGEMENT_LIVE":
        if (
          msg.data &&
          (!currentDossierId ||
            msg.data.dossierId === currentDossierId ||
            msg.data.candidatNom === activeCand.value?.nom)
        ) {
          liveLivretEvaluation.value = {
            sessionId: msg.data.sessionId,
            dossierId: msg.data.dossierId,
            candidatNom: msg.data.candidatNom,
            kmDebut: msg.data.kmDebut ?? 0,
            kmFin: msg.data.kmFin ?? 0,
            commentaire: msg.data.commentaire ?? "",
            competences: Array.isArray(msg.data.competences)
              ? msg.data.competences
              : [],
            hasSigned: Boolean(msg.data.hasSigned),
            updatedAt: new Date().toLocaleTimeString("fr-FR"),
          };
        }
        break;
      case "COURS_CODE_UPDATED":
        await coursStore.loadSeances();
        if (currentCandId) await store.selectCandidat(currentCandId);
        break;
      case "SUPPORTS_COURS_UPDATED":
        await supportsStore.loadSupports();
        break;
      case "TUTORIELS_UPDATED":
        await tutorielsStore.loadTutoriels();
        break;
      case "EPREUVES_UPDATED":
        await epreuvesStore.loadEpreuves();
        break;
      case "QUIZ_FORM_UPDATED":
        await quizStore.loadFormulaires();
        break;
      case "FLOTTE_VEHICULES_UPDATED":
        await store.loadVehicules();
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

watch(selectedCandId, async (newId) => {
  if (
    newId &&
    (!store.activeDossierDetail ||
      store.activeDossierDetail.candidat.id !== newId)
  ) {
    lastQuizScore.value = null;
    await store.selectCandidat(newId);
    await loadCandidatDrivingSessions();
  }
});

const activeCand = computed<CandidatProfileSummary | null>(() => {
  const d = store.activeDossierDetail;
  if (!d) {
    const u = authStore.currentUser;
    if (!u) return null;
    return {
      id: u.id,
      dossierId: "",
      nom: `${u.nom || ""} ${u.prenom || ""}`.trim() || "Apprenant Connecté",
      numeroDossier: "DOS-ELEVE",
      telephone: u.telephone || "",
      formule: "Formation Permis B",
      categoriePermis: "B",
      agence: authStore.schoolName,
      montantTotal: 0,
      dejaPaye: 0,
      soldeRestant: 0,
      heuresCode: 0,
      heuresCodeMax: 20,
      heuresConduite: 0,
      heuresConduiteMax: 20,
      dernierTest: 0,
      dateInscription: new Date().toLocaleDateString("fr-FR"),
      recus: [],
      modeTheorique: "HYBRIDE",
    };
  }

  const c = d.candidat;
  const dos = d.dossier;
  const f = d.formule;
  const dejaPaye = Math.max(0, dos.montant_total - dos.solde_restant);
  const recus = (d.paiements || []).map((p) => ({
    numero: p.reference_recu || `REC-${p.id.slice(0, 8).toUpperCase()}`,
    date: new Date(p.date_paiement).toLocaleDateString("fr-FR"),
    heure: new Date(p.date_paiement).toLocaleTimeString("fr-FR", {
      hour: "2-digit",
      minute: "2-digit",
    }),
    montant: p.montant,
    mode: p.mode_reglement,
    hash: p.reference_recu || "P-VALIDE",
  }));

  return {
    id: c.id,
    dossierId: dos.id,
    nom: `${c.nom} ${c.prenoms}`.trim(),
    numeroDossier: dos.numero_dossier,
    telephone: c.telephone_whatsapp || "",
    formule: f.intitule,
    categoriePermis: f.categorie_permis || "B",
    agence: authStore.schoolName,
    montantTotal: dos.montant_total,
    dejaPaye,
    soldeRestant: dos.solde_restant,
    heuresCode: dos.heures_code_effectuees,
    heuresCodeMax: f.heures_theoriques,
    heuresConduite: dos.heures_conduite_effectuees,
    heuresConduiteMax: f.heures_pratiques,
    dernierTest: d.dernier_score_code || 0,
    dateInscription: new Date(
      dos.date_inscription || c.created_at,
    ).toLocaleDateString("fr-FR"),
    recus,
    modeTheorique: dos.mode_formation_theorique || "PRESENTIEL",
  };
});

const openQuiz = () => (isQuizModalOpen.value = true);
const openReservation = () => (isReservationModalOpen.value = true);

const handleQuizFinish = async (s: number) => {
  lastQuizScore.value = s;
  if (activeCand.value) {
    try {
      await api.enregistrerQuizScore({
        candidat_id: activeCand.value.id,
        score: s,
      });
      await store.selectCandidat(activeCand.value.id);
      store.showToast(
        `Quiz enregistré : ${s} / 40 (${s >= SEUIL_EXAMEN ? "Éligible pour l'examen officiel !" : "Poursuivez vos entraînements."})`,
        "success",
      );
    } catch (err: any) {
      store.showToast(
        `Quiz terminé : ${s} / 40 (Erreur enregistrement: ${err.message})`,
        "info",
      );
    }
  }
};

const handleReserved = async (payload: {
  dossierId: string;
  date: string;
  heureDebut: string;
  heureFin: string;
  vehiculeId?: string;
  vehiculeLabel: string;
}) => {
  try {
    await api.reserverSeanceConduite({
      dossier_id: payload.dossierId,
      date_seance: payload.date,
      heure_debut: payload.heureDebut,
      heure_fin: payload.heureFin,
      vehicule_id: payload.vehiculeId,
    });
    isReservationModalOpen.value = false;
    store.showToast(
      `Leçon de conduite réservée le ${new Date(payload.date).toLocaleDateString("fr-FR")} (${payload.heureDebut} – ${payload.heureFin})`,
      "success",
    );
    if (activeCand.value) await store.selectCandidat(activeCand.value.id);
  } catch (err: any) {
    store.showToast(`Erreur lors de la réservation : ${err.message}`, "error");
  }
};

const startPayment = (provider: string) => {
  if (!activeCand.value) return;
  const amount = activeCand.value.soldeRestant.toLocaleString("fr-FR");
  store.showToast(
    `Paiement de ${amount} FCFA via ${provider} : redirection vers la passerelle sécurisée...`,
    "info",
  );
};

const studentCategory = computed(() =>
  (activeCand.value?.categoriePermis || "B").toUpperCase().trim(),
);
const studentEpreuves = computed(() =>
  epreuvesStore.getEpreuvesForStudent(studentCategory.value),
);
const studentSupports = computed(() =>
  supportsStore.getSupportsForStudent(studentCategory.value),
);
const studentVideos = computed(() =>
  tutorielsStore.getTutorielsForStudent(studentCategory.value),
);
const studentQuizzes = computed(() =>
  quizStore.getQuizForStudent(studentCategory.value),
);

watch(
  () => authStore.activeSection,
  (section) => {
    if (section === "epreuves") activeMainTab.value = "epreuves";
    else if (section === "cours-online") activeMainTab.value = "cours";
    else if (section === "supports-cours") activeMainTab.value = "supports";
    else if (section === "tutoriels") activeMainTab.value = "tutoriels";
    else if (section === "telechargements" || section === "downloads")
      activeMainTab.value = "downloads";
    else if (section === "quiz") activeMainTab.value = "quiz";
    else if (section === "reservation") openReservation();
    else if (section === "historique") activeMainTab.value = "historique";
    else if (section === "overview" || section === "recus")
      activeMainTab.value = "dashboard";
  },
  { immediate: true },
);
</script>

<template>
  <div class="space-y-4">
    <CandidatHeaderBar
      :is-recus-section="authStore.activeSection === 'recus'"
      :is-ws-connected="isWsConnected"
      :last-realtime-sync-at="lastRealtimeSyncAt"
      :search-results="store.searchResults || []"
      v-model:selected-cand-id="selectedCandId"
      :has-active-cand="Boolean(activeCand)"
      :epreuves-count="studentEpreuves.length"
      :quizzes-count="studentQuizzes.length"
      @open-epreuves="activeMainTab = 'epreuves'"
      @open-quiz="activeMainTab = 'quiz'"
      @open-reservation="openReservation"
    />

    <CandidatNavTabs
      v-if="activeCand"
      v-model="activeMainTab"
      :student-category="studentCategory"
      :cours-count="coursStore.seances.length"
      :supports-count="studentSupports.length"
      :videos-count="studentVideos.length"
      :epreuves-count="studentEpreuves.length"
      :quizzes-count="studentQuizzes.length"
      :downloads-count="offlineStore.stats.total"
    />

    <!-- Chargement -->
    <div
      v-if="store.loading && !activeCand"
      class="flex flex-col items-center justify-center gap-2 rounded-[5px] border border-border bg-card p-12 text-center"
    >
      <div
        class="h-5 w-5 animate-spin rounded-full border-2 border-brand-dark border-t-transparent"
      />
      <span class="text-xs text-muted-foreground"
        >Chargement du dossier apprenant...</span
      >
    </div>

    <!-- Aucun dossier -->
    <div
      v-else-if="!activeCand"
      class="flex flex-col items-center justify-center gap-2 rounded-[5px] border border-border bg-card p-12 text-center"
    >
      <FileText class="h-8 w-8 text-muted-foreground/60" />
      <div class="text-sm font-semibold text-foreground">
        Aucun dossier apprenant actif
      </div>
      <p class="max-w-xs text-xs text-muted-foreground">
        Aucun candidat n'a été trouvé pour cette auto-école. Veuillez inscrire
        un nouveau candidat au guichet.
      </p>
    </div>

    <template v-else>
      <CandidatDashboardTab
        v-if="activeMainTab === 'dashboard'"
        :candidat="activeCand"
        :live-livret-evaluation="liveLivretEvaluation"
        :seances-conduite="seancesConduiteCandidat"
        :last-quiz-score="lastQuizScore"
        @open-cours="activeMainTab = 'cours'"
        @open-reservation="openReservation"
        @start-payment="startPayment"
      />

      <CandidatEpreuvesTab
        v-else-if="activeMainTab === 'epreuves'"
        :student-category="studentCategory"
        :epreuves="studentEpreuves"
      />

      <CandidatCoursEnLigneView
        v-else-if="activeMainTab === 'cours'"
        :candidat="activeCand"
      />

      <CandidatSupportsCoursView
        v-else-if="activeMainTab === 'supports'"
        :student-category="studentCategory"
      />

      <CandidatTutorielsView
        v-else-if="activeMainTab === 'tutoriels'"
        :student-category="studentCategory"
      />

      <CandidatQuizFormsView
        v-else-if="activeMainTab === 'quiz'"
        :candidat-id="activeCand.id"
        :candidat-nom="activeCand.nom"
        :numero-dossier="activeCand.numeroDossier"
        :student-category="studentCategory"
        @open-simulator="openQuiz"
        @score-saved="
          (sc) => {
            lastQuizScore = sc;
            if (activeCand) store.selectCandidat(activeCand.id);
          }
        "
      />
      <OfflineDownloadsView
        v-else-if="activeMainTab === 'downloads'"
        @navigate="(tab) => (activeMainTab = tab)"
      />
      <RoleHistoriqueView v-else-if="activeMainTab === 'historique'" />

      <CodeQuizModal
        v-if="isQuizModalOpen"
        @close="isQuizModalOpen = false"
        @finish="handleQuizFinish"
      />
      <ReservationConduiteModal
        v-if="isReservationModalOpen"
        :candidat-nom="activeCand.nom"
        :dossier-id="activeCand.dossierId"
        :numero-dossier="activeCand.numeroDossier"
        :telephone="activeCand.telephone"
        :solde-restant="activeCand.soldeRestant"
        :montant-total="activeCand.montantTotal"
        :heures-effectuees="activeCand.heuresConduite"
        :heures-total="activeCand.heuresConduiteMax"
        @close="isReservationModalOpen = false"
        @reserved="handleReserved"
      />
      <EpreuveViewerModal
        :epreuve="epreuvesStore.activePdfToPreview"
        :hide-download="true"
        @close="epreuvesStore.closePreview"
      />
    </template>
  </div>
</template>

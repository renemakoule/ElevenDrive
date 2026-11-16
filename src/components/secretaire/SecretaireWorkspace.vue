<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from "vue";
import { useAuthStore } from "../../stores/authStore";
import { useGuichetStore } from "../../stores/guichetStore";
import { realtimeSocket } from "../../services/realtimeSocketService";
import SecretaireHomeView from "./SecretaireHomeView.vue";
import SearchAndCandidatePanel from "../guichet/SearchAndCandidatePanel.vue";
import QuickPaymentPanel from "../guichet/QuickPaymentPanel.vue";
import CodeAttendanceView from "../guichet/CodeAttendanceView.vue";
import ExamensView from "../examens/ExamensView.vue";
import FormulesTarifsView from "../guichet/FormulesTarifsView.vue";
import EpreuvesManagementView from "../epreuves/EpreuvesManagementView.vue";
import CoursCodePlanningView from "../cours/CoursCodePlanningView.vue";
import PlanningConduiteView from "../cours/PlanningConduiteView.vue";
import SupportsCoursManagementView from "../supports/SupportsCoursManagementView.vue";
import TutorielsManagementView from "../tutoriels/TutorielsManagementView.vue";
import RoleHistoriqueView from "../historique/RoleHistoriqueView.vue";
import CentreArchivageView from "../archivage/CentreArchivageView.vue";
import {
  UserPlus,
  FileSpreadsheet,
  Receipt,
  Search,
  FileText,
  CheckSquare,
  Award,
  Tag,
  ArrowLeft,
  Calculator,
  BookOpen,
  Radio,
} from "lucide-vue-next";

const authStore = useAuthStore();
const guichetStore = useGuichetStore();
const isWsConnected = realtimeSocket.isConnected;
const lastRealtimeSyncAt = ref<string>("");
let unsubscribeWs: (() => void) | null = null;

onMounted(() => {
  realtimeSocket.connect();
  unsubscribeWs = realtimeSocket.subscribe(async (msg) => {
    lastRealtimeSyncAt.value = new Date().toLocaleTimeString("fr-FR");
    switch (msg.type) {
      case "PAIEMENT_ENREGISTRE":
      case "CANDIDAT_DOSSIER_UPDATED":
      case "CAISSE_SESSION_UPDATED":
      case "QUIZ_RESULTAT_UPDATED":
      case "LIVRET_EMARGEMENT_UPDATED":
        await Promise.all([
          guichetStore.loadStats(),
          guichetStore.search(guichetStore.searchQuery || ""),
        ]);
        if (guichetStore.selectedCandidatId) {
          await guichetStore.selectCandidat(guichetStore.selectedCandidatId);
        }
        break;
      case "FORMULES_TARIFS_UPDATED":
        await guichetStore.loadFormules();
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

const triggerNewCandidate = () => {
  guichetStore.isNewCandidateModalOpen = true;
};

const triggerImportCandidates = () => {
  guichetStore.isImportCandidatesModalOpen = true;
};

const triggerPayment = () => {
  authStore.setActiveSection("guichet");
  setTimeout(() => {
    const el = document.getElementById("montant-input");
    if (el) el.focus();
  }, 50);
};

const triggerSearch = () => {
  authStore.setActiveSection("guichet");
  setTimeout(() => {
    const el = document.getElementById("search-input");
    if (el) el.focus();
  }, 50);
};

const triggerCashClose = () => {
  guichetStore.isCashCloseModalOpen = true;
};

// Si l'utilisateur clique sur "Point Caisse (Rapport Z)" depuis la barre latérale
watch(
  () => authStore.activeSection,
  (section) => {
    if (section === "caisse") {
      guichetStore.loadSessionCaisse();
      guichetStore.loadStats();
      guichetStore.isCashCloseModalOpen = true;
    }
  },
);

const fcfa = (n: number) => (n || 0).toLocaleString("fr-FR");
</script>

<template>
  <div
    class="flex-1 min-h-0 h-full flex flex-col overflow-hidden bg-background"
  >
    <!-- Bandeau Mode Compagnon Mobile / Tablette (< md) -->
    <div
      v-if="authStore.platform === 'pwa'"
      class="md:hidden mb-2.5 px-3 py-2 rounded-[3px] border border-brand-border bg-brand-light/70 text-brand-dark text-[11px] flex items-center justify-between gap-2 shrink-0"
    >
      <span class="font-medium leading-snug">
        <strong>Mode Compagnon Mobile</strong> · Consultation dossiers, pointage
        Code et plannings en direct.
      </span>
    </div>

    <!-- Barre de Raccourcis Guichet Ultra-Compacte et Épurée (style minimaliste) -->
    <div
      class="bg-background pb-2.5 flex items-center justify-between text-xs shrink-0 select-none"
    >
      <div class="flex items-center space-x-1.5 overflow-x-auto no-scrollbar">
        <span
          class="text-[10px] uppercase font-bold text-muted-foreground tracking-wider mr-1 hidden sm:inline"
        >
          Accès Clavier :
        </span>

        <!-- F1 Inscription -->
        <button
          type="button"
          @click="triggerNewCandidate"
          class="min-h-[38px] md:min-h-0 flex items-center space-x-1.5 px-3 md:px-2 py-1.5 md:py-1 rounded-[3px] bg-muted hover:bg-muted/80 border border-border text-foreground text-xs md:text-[11px] font-medium transition active:scale-95 shrink-0"
        >
          <kbd class="hidden sm:inline-flex">F1</kbd>
          <UserPlus class="w-3.5 h-3.5 md:w-3 md:h-3 text-brand" />
          <span>Inscription</span>
        </button>

        <!-- Import Excel -->
        <button
          type="button"
          @click="triggerImportCandidates"
          class="min-h-[38px] md:min-h-0 flex items-center space-x-1.5 px-3 md:px-2 py-1.5 md:py-1 rounded-[3px] bg-muted hover:bg-muted/80 border border-border text-foreground text-xs md:text-[11px] font-medium transition active:scale-95 shrink-0"
        >
          <FileSpreadsheet class="w-3.5 h-3.5 md:w-3 md:h-3 text-emerald-600" />
          <span>Import Excel</span>
        </button>

        <!-- F2 Encaissement -->
        <!-- <button
          type="button"
          @click="triggerPayment"
          class="flex items-center space-x-1.5 px-2 py-1 rounded bg-brand-light hover:bg-brand-light/80 border border-brand-border text-brand-dark text-[11px] font-semibold transition active:scale-95"
        >
          <kbd class="border-brand-border text-brand-dark">F2</kbd>
          <Receipt class="w-3 h-3 text-brand" />
          <span>Encaissement</span>
        </button> -->

        <!-- F3 Recherche -->
        <button
          type="button"
          @click="triggerSearch"
          class="min-h-[38px] md:min-h-0 flex items-center space-x-1.5 px-3 md:px-2 py-1.5 md:py-1 rounded-[3px] bg-muted hover:bg-muted/80 border border-border text-foreground text-xs md:text-[11px] font-medium transition active:scale-95 shrink-0"
        >
          <kbd class="hidden sm:inline-flex">F3</kbd>
          <Search class="w-3.5 h-3.5 md:w-3 md:h-3 text-muted-foreground" />
          <span>Recherche</span>
        </button>

        <!-- F4 Caisse Z -->
        <!-- <button
          type="button"
          @click="triggerCashClose"
          class="flex items-center space-x-1.5 px-2 py-1 rounded bg-muted hover:bg-muted/80 border border-border text-foreground text-[11px] font-medium transition active:scale-95"
        >
          <kbd>F4</kbd>
          <FileText class="w-3 h-3 text-amber-600" />
          <span>Point Z</span>
        </button> -->

        <!-- F5 Présences Code -->
        <!-- <button
          type="button"
          @click="authStore.setActiveSection('code')"
          :class="[
            'flex items-center space-x-1.5 px-2 py-1 rounded border text-[11px] font-medium transition active:scale-95',
            authStore.activeSection === 'code'
              ? 'bg-purple-50 border-purple-200 text-purple-800 dark:bg-purple-950/40 dark:border-purple-800 dark:text-purple-300'
              : 'bg-muted hover:bg-muted/80 border-border text-foreground',
          ]"
        >
          <kbd>F5</kbd>
          <CheckSquare class="w-3 h-3 text-purple-600" />
          <span>Code</span>
        </button> -->

        <!-- F9 Examens -->
        <!-- <button
          type="button"
          @click="authStore.setActiveSection('examens')"
          :class="[
            'flex items-center space-x-1.5 px-2 py-1 rounded border text-[11px] font-medium transition active:scale-95',
            authStore.activeSection === 'examens'
              ? 'bg-amber-50 border-amber-200 text-amber-800 dark:bg-amber-950/40 dark:border-amber-800 dark:text-amber-300'
              : 'bg-muted hover:bg-muted/80 border-border text-foreground',
          ]"
        >
          <kbd>F9</kbd>
          <Award class="w-3 h-3 text-amber-600" />
          <span>Examens A4</span>
        </button> -->
      </div>

      <!-- Indicateur Temps Réel WebSocket + Astuce Clavier -->
      <!-- <div
        class="hidden lg:flex items-center space-x-2 text-[10px] text-muted-foreground"
      >
        <span
          class="inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-bold"
          :class="
            isWsConnected
              ? 'bg-emerald-500/15 text-emerald-700 dark:text-emerald-300'
              : 'bg-amber-500/15 text-amber-700'
          "
        >
          <Radio class="h-3 w-3 animate-pulse" />
          {{ isWsConnected ? "Temps réel WS" : "Reconnexion WS..." }}
        </span>
        <span
          v-if="lastRealtimeSyncAt"
          class="font-mono text-[10px] text-muted-foreground"
        >
          ({{ lastRealtimeSyncAt }})
        </span>
        <span><kbd>Entrée</kbd> Valider</span>
        <span><kbd>Échap</kbd> Retour</span>
      </div> -->
    </div>

    <!-- Contenu Principal selon l'onglet actif -->
    <div
      :class="
        authStore.activeSection === 'guichet'
          ? 'flex-1 min-h-0 overflow-y-auto md:overflow-hidden pt-2'
          : 'flex-1 min-h-0 overflow-y-auto pt-2'
      "
    >
      <!-- 1. Vue d'ensemble (style MaryamE-logistic) -->
      <template v-if="authStore.activeSection === 'overview'">
        <SecretaireHomeView />
      </template>

      <!-- 2. Guichet & Caisse (Écran Scindé 50/50 sur Desktop, Empilé scrollable sur Mobile) -->
      <template v-else-if="authStore.activeSection === 'guichet'">
        <div
          class="md:h-full min-h-0 flex flex-col md:flex-row gap-5 overflow-y-auto md:overflow-hidden"
        >
          <SearchAndCandidatePanel />
          <QuickPaymentPanel />
        </div>
      </template>

      <!-- 3. Présences Code en Salle (F5) -->
      <template v-else-if="authStore.activeSection === 'code'">
        <CodeAttendanceView />
      </template>

      <!-- 4. Examens d'État & Bordereau Ministériel (F9) -->
      <template v-else-if="authStore.activeSection === 'examens'">
        <ExamensView />
      </template>

      <!-- 4b. Catalogue des Tarifs Permis -->
      <template v-else-if="authStore.activeSection === 'tarifs'">
        <FormulesTarifsView />
      </template>

      <!-- 4c. Banque d'Épreuves & Annales PDF (Pédagogie) -->
      <template v-else-if="authStore.activeSection === 'epreuves'">
        <EpreuvesManagementView />
      </template>

      <!-- 4c2. Planning Conduite Pratique Moniteurs & Flotte -->
      <template v-else-if="authStore.activeSection === 'conduite-planning'">
        <PlanningConduiteView />
      </template>

      <!-- 4d. Planning Cours de Code Théorique & Visio -->
      <template v-else-if="authStore.activeSection === 'cours-planning'">
        <CoursCodePlanningView />
      </template>

      <!-- 4e. Supports de Cours & Manuels PDF -->
      <template v-else-if="authStore.activeSection === 'supports-cours'">
        <SupportsCoursManagementView />
      </template>

      <!-- 4f. Bibliothèque de Tutoriels & Guides Vidéo -->
      <template v-else-if="authStore.activeSection === 'tutoriels'">
        <TutorielsManagementView />
      </template>

      <!-- 5. Point Caisse Z Dédié -->
      <template v-else-if="authStore.activeSection === 'caisse'">
        <div class="max-w-2xl mx-auto space-y-4 py-4 px-2 sm:px-0">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <button
                type="button"
                @click="authStore.setActiveSection('overview')"
                class="inline-flex h-8 w-8 items-center justify-center rounded-lg border border-border bg-card text-muted-foreground hover:bg-muted hover:text-foreground transition"
                title="Retour à la vue d'ensemble"
              >
                <ArrowLeft class="h-4 w-4" />
              </button>
              <div>
                <h2 class="text-sm font-semibold text-foreground">
                  Point de Caisse Z & Clôture
                </h2>
                <p class="text-xs text-muted-foreground">
                  Session de caisse active du guichet.
                </p>
              </div>
            </div>

            <button
              type="button"
              @click="triggerCashClose"
              class="inline-flex h-8 items-center gap-1.5 rounded-lg bg-brand-dark px-3 text-xs font-medium text-white transition hover:bg-brand"
            >
              <Calculator class="h-3.5 w-3.5" />
              Ouvrir Pointage Z
              <kbd
                class="rounded bg-white/20 px-1 font-mono text-[10px] text-white/90"
                >F4</kbd
              >
            </button>
          </div>

          <div
            class="rounded-xl border border-border bg-card p-5 space-y-4 shadow-sm"
          >
            <div
              class="flex items-center justify-between border-b border-border pb-3"
            >
              <div>
                <span class="text-xs text-muted-foreground"
                  >Caissière / Agent en poste</span
                >
                <p class="font-medium text-sm text-foreground">
                  {{ authStore.currentUser.prenom }}
                  {{ authStore.currentUser.nom }} ({{
                    authStore.currentUser.roleLabel || "Guichet"
                  }})
                </p>
              </div>
              <span
                class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium bg-emerald-50 text-emerald-700 border border-emerald-200 dark:bg-emerald-950/40 dark:text-emerald-300 dark:border-emerald-800"
              >
                <span
                  class="h-1.5 w-1.5 rounded-full bg-emerald-500 animate-pulse"
                />
                Session active
              </span>
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
              <div class="p-3.5 rounded-lg border border-border bg-muted/30">
                <span class="text-xs text-muted-foreground"
                  >Fond de caisse initial</span
                >
                <p
                  class="mt-1 text-base font-semibold tabular-nums text-foreground"
                >
                  {{
                    fcfa(
                      guichetStore.sessionCaisseDetail?.fond_de_caisse ??
                        guichetStore.stats.fond_de_caisse,
                    )
                  }}
                  FCFA
                </p>
              </div>

              <div class="p-3.5 rounded-lg border border-border bg-muted/30">
                <span class="text-xs text-muted-foreground"
                  >Espèces encaissées (tiroir)</span
                >
                <p
                  class="mt-1 text-base font-semibold tabular-nums text-brand-dark"
                >
                  {{
                    fcfa(
                      guichetStore.sessionCaisseDetail?.ventilation?.especes ??
                        guichetStore.stats.total_especes_theorique,
                    )
                  }}
                  FCFA
                </p>
              </div>

              <div class="p-3.5 rounded-lg border border-border bg-muted/30">
                <span class="text-xs text-muted-foreground"
                  >Recette totale ({{
                    guichetStore.sessionCaisseDetail?.nombre_transactions ?? 0
                  }}
                  tx)</span
                >
                <p
                  class="mt-1 text-base font-semibold tabular-nums text-foreground"
                >
                  {{
                    fcfa(
                      guichetStore.sessionCaisseDetail?.ventilation
                        ?.total_encaisse ??
                        guichetStore.stats.total_especes_theorique,
                    )
                  }}
                  FCFA
                </p>
              </div>
            </div>

            <div
              class="flex items-center justify-between pt-2 border-t border-border text-xs"
            >
              <span class="text-muted-foreground"
                >Total théorique attendu au tiroir :</span
              >
              <span class="text-base font-bold tabular-nums text-foreground">
                {{
                  fcfa(
                    (guichetStore.sessionCaisseDetail?.fond_de_caisse ??
                      guichetStore.stats.fond_de_caisse) +
                      (guichetStore.sessionCaisseDetail?.ventilation?.especes ??
                        guichetStore.stats.total_especes_theorique),
                  )
                }}
                FCFA
              </span>
            </div>
          </div>
        </div>
      </template>

      <!-- 6. Historique des Opérations & Activités en Temps Réel -->
      <template v-else-if="authStore.activeSection === 'historique'">
        <RoleHistoriqueView />
      </template>

      <!-- 7. Centre d'Archivage & Documents Importants -->
      <template v-else-if="authStore.activeSection === 'archivage'">
        <CentreArchivageView />
      </template>
    </div>
  </div>
</template>

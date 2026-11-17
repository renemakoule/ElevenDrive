<script setup lang="ts">
import "./assets/main.css";
import { ref, onMounted, onUnmounted, computed, watch } from "vue";
import { useGuichetStore } from "./stores/guichetStore";
import { useAuthStore } from "./stores/authStore";
import { useLanguageStore } from "./stores/languageStore";
import { useNotificationStore } from "./stores/notificationStore";
import { domAutoTranslator } from "./services/domAutoTranslator";
import { realtimeSocket } from "./services/realtimeSocketService";
import { authBackend } from "./services/authBackendService";
import { checkForAppUpdate } from "./services/tauriUpdaterService";
import { isTauri, invokeTauri } from "./db/httpClient";
import IdleLockOverlay from "./components/common/IdleLockOverlay.vue";

// Layout & Global Components (style MaryamE-logistic)
import TauriDesktopHeader from "./components/layout/TauriDesktopHeader.vue";
import TauriUpdaterBanner from "./components/layout/TauriUpdaterBanner.vue";
import AppBootLoader from "./components/layout/AppBootLoader.vue";
import AppHeader from "./components/layout/AppHeader.vue";
import AppSidebar from "./components/layout/AppSidebar.vue";
import MobileBottomNav from "./components/layout/MobileBottomNav.vue";
import HelpFloatingButton from "./components/layout/HelpFloatingButton.vue";
import PlatformGuardView from "./components/common/PlatformGuardView.vue";
import LicenceExpiredGuardView from "./components/common/LicenceExpiredGuardView.vue";
import SchoolPricingModal from "./components/pricing/SchoolPricingModal.vue";

// Role-Specific Workspaces
import SecretaireWorkspace from "./components/secretaire/SecretaireWorkspace.vue";
import ManagerWorkspace from "./components/manager/ManagerWorkspace.vue";
import MoniteurWorkspace from "./components/moniteur/MoniteurWorkspace.vue";
import EleveWorkspace from "./components/eleve/EleveWorkspace.vue";

// Transaction Modals
import NewCandidateModal from "./components/guichet/NewCandidateModal.vue";
import ImportCandidatesModal from "./components/guichet/ImportCandidatesModal.vue";
import CashPointModal from "./components/guichet/CashPointModal.vue";
import ThermalReceiptModal from "./components/guichet/ThermalReceiptModal.vue";
import ThermalRapportZModal from "./components/guichet/ThermalRapportZModal.vue";
import ManagerOverrideModal from "./components/guichet/ManagerOverrideModal.vue";
import DesktopAuthPortal from "./components/auth/DesktopAuthPortal.vue";
import DesktopOnboardingTour from "./components/onboarding/DesktopOnboardingTour.vue";
import WebPwaPortal from "./components/pwa/WebPwaPortal.vue";
import LandingView from "./components/landing/LandingView.vue";
import SuperAdminDashboardView from "./components/superadmin/SuperAdminDashboardView.vue";
import SuperAdminLoginView from "./components/superadmin/SuperAdminLoginView.vue";

import { CheckCircle2, AlertCircle, Info } from "lucide-vue-next";

const guichetStore = useGuichetStore();
const authStore = useAuthStore();
const languageStore = useLanguageStore();
const notifStore = useNotificationStore();
const isPricingModalOpen = ref(false);
const isAppInitializing = ref(true);
let unsubscribeGlobalWs: (() => void) | null = null;

watch(
  () => [authStore.isAuthenticated, authStore.currentRole, authStore.currentUser?.id, authStore.currentSchool?.id] as const,
  ([auth, role, userId, schoolId]) => {
    if (!isTauri) return;
    if (auth) {
      invokeTauri("set_local_session_context", {
        userId: userId || "",
        autoEcoleId: schoolId || "",
        role: role || "secretaire",
      }).catch(() => {});
    } else {
      invokeTauri("clear_local_session_context").catch(() => {});
    }
  },
  { immediate: true },
);

const handleGlobalKeyDown = (e: KeyboardEvent) => {
  // Raccourcis clavier uniquement actifs pour la secrétaire sur Desktop
  if (
    authStore.currentRole === "secretaire" &&
    authStore.isCurrentRoleAllowed
  ) {
    if (e.key === "F1") {
      e.preventDefault();
      guichetStore.isNewCandidateModalOpen = true;
    } else if (e.key === "F2") {
      e.preventDefault();
      authStore.setActiveSection("guichet");
      setTimeout(() => {
        const el = document.getElementById("montant-input");
        if (el) el.focus();
      }, 50);
    } else if (e.key === "F3") {
      e.preventDefault();
      authStore.setActiveSection("guichet");
      setTimeout(() => {
        const el = document.getElementById("search-input");
        if (el) el.focus();
      }, 50);
    } else if (e.key === "F4") {
      e.preventDefault();
      guichetStore.isCashCloseModalOpen = true;
    } else if (e.key === "F5") {
      e.preventDefault();
      authStore.setActiveSection(
        authStore.activeSection === "code" ? "overview" : "code",
      );
    } else if (e.key === "F9") {
      e.preventDefault();
      authStore.setActiveSection(
        authStore.activeSection === "examens" ? "overview" : "examens",
      );
    }
  }

  // Touche Échap globale pour fermer les modales ou revenir à l'accueil
  if (e.key === "Escape") {
    if (guichetStore.isReceiptModalOpen) {
      guichetStore.isReceiptModalOpen = false;
    } else if (guichetStore.isRapportZModalOpen) {
      guichetStore.isRapportZModalOpen = false;
    } else if (guichetStore.isManagerOverrideModalOpen) {
      guichetStore.isManagerOverrideModalOpen = false;
    } else if (guichetStore.isNewCandidateModalOpen) {
      guichetStore.isNewCandidateModalOpen = false;
    } else if (guichetStore.isImportCandidatesModalOpen) {
      guichetStore.isImportCandidatesModalOpen = false;
    } else if (guichetStore.isCashCloseModalOpen) {
      guichetStore.isCashCloseModalOpen = false;
    } else if (authStore.activeSection !== "overview") {
      authStore.setActiveSection("overview");
    }
  }
};

// Détection du premier lancement Desktop pour la visite guidée
const isDesktopTourCompleted = ref(true);

const checkDesktopTourState = () => {
  try {
    const onboarded = localStorage.getItem("elevendrive_desktop_onboarded_v1");
    isDesktopTourCompleted.value = onboarded === "true";
  } catch {
    isDesktopTourCompleted.value = true;
  }
};

const handleTourComplete = () => {
  isDesktopTourCompleted.value = true;
  try {
    localStorage.setItem("elevendrive_desktop_onboarded_v1", "true");
  } catch {}
};

// Détection stricte de l'environnement natif Tauri (absent sur navigateur Web)
const isDesktopTauri = computed(() => {
  return (
    typeof window !== "undefined" &&
    ("__TAURI_INTERNALS__" in window || "__TAURI__" in window)
  );
});

// Détection stricte du mode console Super-Admin (évite toute superposition de LandingView)
const isSuperAdminMode = computed(() => {
  return (
    (authStore.isAuthenticated && authStore.currentRole === "super_admin") ||
    authStore.webPortalView === "super_admin_dashboard" ||
    authStore.webPortalView === "super_admin_login"
  );
});

const handlePopState = () => {
  authStore.initFromEnvironment();
};

const initAuthenticatedSync = async () => {
  if (!authStore.isAuthenticated) return;
  const session = authBackend.getActiveSession();
  const token = session?.token;
  if (token) {
    notifStore.initRealtime(token);
    notifStore.loadNotifications();
  }
  realtimeSocket.connect();
  if (!unsubscribeGlobalWs) {
    unsubscribeGlobalWs = realtimeSocket.subscribe((msg) => {
      if (msg && msg.type && msg.type.startsWith("SCHOOL_")) {
        authStore.handleRealtimeSchoolEvent({ type: msg.type, data: msg.data });
      }
    });
  }
};

const handleVisibilitySync = () => {
  if (
    typeof document !== "undefined" &&
    document.visibilityState === "visible" &&
    authStore.isAuthenticated &&
    authStore.currentRole !== "super_admin"
  ) {
    authStore.loadProfileByRole().catch(() => {});
  }
};

onMounted(async () => {
  try {
    languageStore.init();
    checkDesktopTourState();
    await authStore.initFromEnvironment();
    if (authStore.isAuthenticated) {
      await guichetStore.init();
      await initAuthenticatedSync();
    }
    domAutoTranslator.init();
    window.addEventListener("keydown", handleGlobalKeyDown);
    window.addEventListener("popstate", handlePopState);
    document.addEventListener("visibilitychange", handleVisibilitySync);
    if (isDesktopTauri.value) {
      setTimeout(() => checkForAppUpdate(true), 3000);
    }
  } finally {
    isAppInitializing.value = false;
  }
});

watch(
  () => authStore.isAuthenticated,
  async (isAuth) => {
    if (isAuth) {
      guichetStore.init();
      await initAuthenticatedSync();
    } else if (unsubscribeGlobalWs) {
      unsubscribeGlobalWs();
      unsubscribeGlobalWs = null;
    }
  },
);

watch(
  () => authStore.activeSection,
  () => {
    window.scrollTo(0, 0);
    if (document.documentElement) document.documentElement.scrollTop = 0;
    if (document.body) document.body.scrollTop = 0;
  },
);

onUnmounted(() => {
  domAutoTranslator.destroy();
  window.removeEventListener("keydown", handleGlobalKeyDown);
  window.removeEventListener("popstate", handlePopState);
  document.removeEventListener("visibilitychange", handleVisibilitySync);
  if (unsubscribeGlobalWs) {
    unsubscribeGlobalWs();
    unsubscribeGlobalWs = null;
  }
});
</script>

<template>
  <div
    class="h-screen w-screen flex flex-col bg-background text-foreground overflow-hidden font-sans select-none antialiased"
  >
    <!-- 0. Header Desktop Tauri (Barre de titre système, Logo contrasté, Menus & Contrôles) -->
    <TauriDesktopHeader v-if="isDesktopTauri" />
    <TauriUpdaterBanner v-if="isDesktopTauri" />

    <!-- 0b. État de chargement initial (évite tout flash vide ou mauvaise vue avant init) -->
    <AppBootLoader v-if="isAppInitializing" />

    <!-- ========================================================================= -->
    <!-- CONSOLE SUPER-ADMIN DÉDIÉE (Prioritaire : aucun composant tiers affiché)  -->
    <!-- ========================================================================= -->
    <template v-else-if="isSuperAdminMode">
      <!-- Dashboard Super-Admin -->
      <SuperAdminDashboardView
        v-if="
          (authStore.isAuthenticated &&
            authStore.currentRole === 'super_admin') ||
          authStore.webPortalView === 'super_admin_dashboard'
        "
      />

      <!-- Authentification Super-Admin (/super-admin/login) -->
      <SuperAdminLoginView
        v-else-if="authStore.webPortalView === 'super_admin_login'"
        @login-success="authStore.webPortalView = 'super_admin_dashboard'"
      />
    </template>

    <!-- ========================================================================= -->
    <!-- ESPACE AUTO-ÉCOLES STANDARD & VISITEURS (Directeur, Secrétaire, Élève...) -->
    <!-- ========================================================================= -->
    <template v-else>
      <!-- 1. En-tête Principal de l'Application (uniquement si connecté avec compte standard) -->
      <AppHeader v-if="authStore.isAuthenticated" />

      <!-- 2. Corps de l'Application (Sidebar Gauche + Cadre Dashboard Principal, uniquement si connecté) -->
      <div
        v-if="authStore.isAuthenticated"
        class="flex-1 min-h-0 flex overflow-hidden"
      >
        <!-- Barre Latérale de Navigation Dynamique (sans border-r) -->
        <AppSidebar />

        <!-- Zone de Contenu Principal : Bord-à-bord sur mobile, encadrée sur Desktop -->
        <main
          class="flex-1 min-h-0 flex flex-col overflow-hidden bg-card px-0 pb-0 pt-0 md:px-4 md:pb-2 md:pt-2"
        >
          <div
            class="flex-1 min-h-0 flex flex-col border-0 md:border md:border-border rounded-none md:rounded-[4px] overflow-y-auto bg-background p-3 sm:p-4 md:p-5"
          >
            <!-- Règle de restriction de plateforme (Secrétaire = Desktop, Moniteur/Élève = PWA) -->
            <template v-if="!authStore.isCurrentRoleAllowed">
              <PlatformGuardView />
            </template>

            <!-- Garde-fou Sécurité : Expiration de la période d'essai de 30 jours (Cryptage Direction) -->
            <template
              v-else-if="
                authStore.isLicenceExpired &&
                (authStore.currentRole === 'manager' ||
                  authStore.currentRole === 'secretaire')
              "
            >
              <LicenceExpiredGuardView
                @open-pricing="isPricingModalOpen = true"
              />
            </template>

            <!-- Espace Secrétaire (Desktop Exclusif) -->
            <template v-else-if="authStore.currentRole === 'secretaire'">
              <SecretaireWorkspace />
            </template>

            <!-- Espace Manager (Desktop & PWA) -->
            <template v-else-if="authStore.currentRole === 'manager'">
              <ManagerWorkspace />
            </template>

            <!-- Espace Moniteur (PWA Mobile Exclusif) -->
            <template v-else-if="authStore.currentRole === 'moniteur'">
              <MoniteurWorkspace />
            </template>

            <!-- Espace Élève (PWA Mobile Exclusif) -->
            <template v-else-if="authStore.currentRole === 'eleve'">
              <EleveWorkspace />
            </template>
          </div>
        </main>
      </div>

      <!-- 2b. Barre de Navigation Tactile Mobile (< md) -->
      <MobileBottomNav
        v-if="
          authStore.isAuthenticated &&
          authStore.isCurrentRoleAllowed &&
          !(
            authStore.isLicenceExpired &&
            (authStore.currentRole === 'manager' ||
              authStore.currentRole === 'secretaire')
          )
        "
      />

      <!-- 3. Bouton Flottant "Besoin d'aide ?" (style MaryamE-logistic) -->
      <HelpFloatingButton />

      <!-- 4. Modales Globales de Gestion Guichet & Reçus (uniquement si authentifié standard) -->
      <template v-if="authStore.isAuthenticated">
        <NewCandidateModal />
        <ImportCandidatesModal />
        <CashPointModal />
        <ThermalReceiptModal />
        <ThermalRapportZModal
          @close="guichetStore.isRapportZModalOpen = false"
        />
        <ManagerOverrideModal />
        <!-- Modale Bon Plan & Tarifs Abonnements -->
        <SchoolPricingModal
          :is-open="isPricingModalOpen"
          @close="isPricingModalOpen = false"
        />
      </template>

      <!-- 5. Si Non Authentifié : Portails de connexion Auto-École ou Visiteur -->
      <template v-if="!authStore.isAuthenticated">
        <!-- Visite Guidée Première Installation (Desktop Tauri) -->
        <DesktopOnboardingTour
          v-if="isDesktopTauri && !isDesktopTourCompleted"
          @complete="handleTourComplete"
          @skip="handleTourComplete"
        />

        <!-- Mode Desktop Tauri natif : Portail d'authentification Desktop -->
        <DesktopAuthPortal
          v-else-if="authStore.platform === 'desktop'"
          @reopen-tour="isDesktopTourCompleted = false"
        />

        <!-- Mode Web avec auto-école spécifiée (QR Code / lien direct auto-école) -->
        <WebPwaPortal v-else-if="authStore.hasExplicitSchoolInUrl" />

        <!-- Mode Web : Connexion ou inscription demandée -> bascule sur DesktopAuthPortal -->
        <DesktopAuthPortal
          v-else-if="authStore.webPortalView === 'desktop_auth'"
          @close="authStore.setWebPortalView('landing')"
        />

        <!-- Mode Web par défaut : Interface d'accueil par défaut de la plateforme (Landing Stripe) -->
        <LandingView
          v-else
          @navigate-auth="authStore.setWebPortalView('desktop_auth')"
        />
      </template>
    </template>

    <!-- 6. Toast Flottant de Notification Épuré (commun à tous les modes) -->
    <IdleLockOverlay />
    <transition
      enter-active-class="transform transition ease-out duration-300"
      enter-from-class="translate-y-4 opacity-0"
      enter-to-class="translate-y-0 opacity-100"
      leave-active-class="transition ease-in duration-200"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="guichetStore.toastMessage"
        :class="[
          'fixed bottom-20 md:bottom-5 right-3 md:right-5 left-3 sm:left-auto z-50 flex items-center gap-2.5 px-3.5 py-2.5 rounded-[3px] text-xs font-semibold border backdrop-blur-md',
          guichetStore.toastType === 'success'
            ? 'bg-brand-light text-brand-dark border-brand-border'
            : guichetStore.toastType === 'error'
              ? 'bg-rose-50 text-rose-900 border-rose-300'
              : 'bg-card text-card-foreground border-border',
        ]"
      >
        <CheckCircle2
          v-if="guichetStore.toastType === 'success'"
          class="w-4 h-4 text-brand"
        />
        <AlertCircle
          v-else-if="guichetStore.toastType === 'error'"
          class="w-4 h-4 text-rose-600"
        />
        <Info v-else class="w-4 h-4 text-sky-600" />
        <span>{{ guichetStore.toastMessage }}</span>
      </div>
    </transition>
  </div>
</template>

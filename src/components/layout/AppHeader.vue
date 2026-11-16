<script setup lang="ts">
import appLogo from "@root/app-icon.png";
import { ref, computed, onMounted, onUnmounted } from "vue";
import { useAuthStore } from "../../stores/authStore";
import { useGuichetStore } from "../../stores/guichetStore";
import { useLanguageStore } from "../../stores/languageStore";
import { neonSyncEngine, type SyncState } from "../../services/neonSyncService";
import SyncStatusModal from "../guichet/SyncStatusModal.vue";
import type { UserRole, AppPlatform } from "../../types";
import {
  PanelLeft,
  CheckCircle2,
  Bell,
  ChevronDown,
  Monitor,
  Smartphone,
  ShieldCheck,
  RefreshCw,
  LogOut,
  Sparkles,
  Sun,
  Moon,
  QrCode,
  Building2,
  MapPin,
  Settings,
  ChevronRight,
  WifiOff,
  WifiIcon,
  Wifi,
} from "lucide-vue-next";
import SchoolQrHubModal from "../common/SchoolQrHubModal.vue";
import UserProfileAndSchoolModal from "../profile/UserProfileAndSchoolModal.vue";
import NetworkNotificationsModal from "./NetworkNotificationsModal.vue";
import NotificationDrawer from "../notifications/NotificationDrawer.vue";
import SchoolPricingModal from "../pricing/SchoolPricingModal.vue";
import { useNotificationStore } from "../../stores/notificationStore";
import { authBackend } from "../../services/authBackendService";

const authStore = useAuthStore();
const guichetStore = useGuichetStore();
const languageStore = useLanguageStore();
const notifStore = useNotificationStore();

const isProfileModalOpen = ref(false);
const isSyncModalOpen = ref(false);
const isQrModalOpen = ref(false);
const isNotificationsModalOpen = ref(false);
const isDrawerOpen = ref(false);
const isPricingOpen = ref(false);
const syncState = ref<SyncState>(neonSyncEngine.getState());

let syncUnsubscribe: (() => void) | null = null;

const openProfileModal = () => {
  isProfileModalOpen.value = true;
};

const openNotificationDrawer = () => {
  isDrawerOpen.value = true;
  if (notifStore.unreadCount > 0) {
    notifStore.markAllAsRead();
  }
};

const handleOpenProfileEvent = () => {
  isProfileModalOpen.value = true;
};
const handleOpenQrEvent = () => {
  isQrModalOpen.value = true;
};

onMounted(() => {
  syncUnsubscribe = neonSyncEngine.subscribe((state) => {
    syncState.value = state;
  });
  // Déclencher une vérification de connectivité réelle dès le montage du composant
  neonSyncEngine.verifierConnectiviteReelle();

  // Initialisation du flux WebSocket et chargement des notifications
  const token = authBackend.getActiveSession()?.token;
  if (token) {
    notifStore.initRealtime(token);
    notifStore.loadNotifications();
  }

  window.addEventListener("elevendrive:open-profile", handleOpenProfileEvent);
  window.addEventListener("elevendrive:open-qr", handleOpenQrEvent);
});

onUnmounted(() => {
  if (syncUnsubscribe) syncUnsubscribe();
  window.removeEventListener("elevendrive:open-profile", handleOpenProfileEvent);
  window.removeEventListener("elevendrive:open-qr", handleOpenQrEvent);
});

const syncBadge = computed(() => {
  const isHorsLigne =
    !syncState.value.isOnline ||
    (typeof navigator !== "undefined" && !navigator.onLine);

  if (isHorsLigne) {
    return {
      label: languageStore.t("Hors ligne"),
      box: "bg-transparent border-rose-500/30 text-rose-600 dark:text-rose-400 font-bold",
      dot: "bg-transparent",
      isOffline: true,
    };
  }
  if (syncState.value.isSyncing) {
    return {
      label: languageStore.t("Synchronisation..."),
      box: "bg-transparent border-sky-500/20 text-sky-600 dark:text-sky-400 animate-pulse",
      dot: "bg-transparent",
      isOffline: false,
    };
  }
  if (!syncState.value.neonConnected || syncState.value.pendingCount > 0) {
    return {
      label: languageStore.t("En attente"),
      box: "bg-transparent border-amber-500/20 text-amber-600 dark:text-amber-400",
      dot: "bg-transparent",
      isOffline: false,
    };
  }
  return {
    label: languageStore.t("En ligne"),
    box: "bg-transparent border-brand-border text-brand-dark hover:bg-brand-light/80",
    dot: "bg-transparent",
    isOffline: false,
  };
});

const togglePlatform = () => {
  const newPlatform: AppPlatform =
    authStore.platform === "desktop" ? "pwa" : "desktop";
  authStore.setPlatform(newPlatform);
};
</script>

<template>
  <header
    class="h-12 md:h-14 pt-safe bg-card px-3 md:px-4 border-b border-border md:border-b-0 flex items-center justify-between select-none z-20 shrink-0"
  >
    <!-- Section Gauche : Logo + Toggle Sidebar + Salutation -->
    <div class="flex items-center space-x-3 md:space-x-4 min-w-0">
      <!-- Logo AutoSchool. (style MaryamE-logistic) -->
      <div class="flex items-center space-x-2 md:space-x-2.5 mr-1 md:mr-2 min-w-0">
        <img
          :src="appLogo"
          alt="AutoSchool Logo"
          class="h-6 w-6 object-contain rounded-[3px] shrink-0"
        />
        <span class="text-sm md:text-base font-extrabold tracking-tight text-foreground truncate">
          {{ authStore.schoolName }}<span class="text-brand">.</span>
        </span>
        <button
          type="button"
          @click="authStore.toggleSidebar()"
          class="hidden md:inline-flex p-1 rounded-[3px] text-muted-foreground hover:text-foreground hover:bg-muted transition"
          title="Réduire / Agrandir la barre latérale"
        >
          <PanelLeft class="w-4 h-4" />
        </button>
      </div>

      <!-- Salutation utilisateur & Badge Vérifié (Clic direct pour ouvrir le profil) -->
      <button
        type="button"
        @click="openProfileModal"
        class="hidden sm:flex flex-col text-left hover:opacity-80 transition cursor-pointer focus:outline-none"
        title="Cliquer pour ouvrir le profil & l'établissement"
      >
        <div class="flex items-center space-x-1.5">
          <span class="text-xs font-bold text-foreground">
            {{ authStore.greetingMessage }}
          </span>
          <span
            class="inline-flex items-center text-brand"
            title="Compte Actif & Vérifié"
          >
            <CheckCircle2 class="w-3.5 h-3.5 fill-brand-light" />
          </span>
        </div>
        <span class="text-[10px] text-muted-foreground font-medium">
          {{ languageStore.t(authStore.welcomeSubtitle) }}
        </span>
      </button>
    </div>

    <!-- Section Droite : Statut Sync, Rôle Switcher, Notifications & Avatar -->
    <div class="flex items-center space-x-1.5 sm:space-x-2.5 shrink-0">
      <!-- Bouton Accès Web & QR Code de Connexion Unique -->
      <button
        type="button"
        @click="isQrModalOpen = true"
        class="hidden sm:flex items-center space-x-1.5 px-2.5 py-1 rounded-[3px] text-[11px] font-semibold border border-brand-border bg-transparent text-brand-dark hover:bg-brand-light/80 transition"
        title="Afficher le QR Code et le lien d'accès au portail web de l'auto-école"
      >
        <QrCode class="w-3.5 h-3.5 text-brand" />
        <span class="hidden sm:inline">{{ languageStore.t("Accès Web") }}</span>
      </button>

      <!-- Badge Synchronisation Cloud Neontech -->
      <button
        type="button"
        @click="isSyncModalOpen = true"
        class="flex items-center space-x-1.5 px-2.5 py-1 rounded-[3px] text-[11px] font-medium border transition cursor-pointer"
        :class="syncBadge.box"
        :title="
          syncBadge.isOffline
            ? 'Mode 100% Hors-ligne : Connexion Internet coupée'
            : 'État de la synchronisation Cloud Neontech'
        "
      >
        <WifiOff
          v-if="syncBadge.isOffline"
          class="w-3.5 h-3.5 text-rose-600 dark:text-rose-400"
        />
        <Wifi
          v-else
          class="w-3.5 h-3.5 text-cyan-600 dark:text-cyan-400"
          :class="syncBadge.dot"
        />
        <!-- <span
          v-else
          class="w-1.5 h-1.5 rounded-full"
          :class="syncBadge.dot"
        ></span> -->
        <!-- Toujours visible si Hors ligne, masqué sur petits écrans uniquement si connecté -->
        <span
          class="text-[10px]"
          :class="syncBadge.isOffline ? 'inline font-bold' : 'hidden lg:inline'"
        >
          {{ syncBadge.label }}
        </span>
      </button>

      <!-- Bouton / Badge Bon Plan pour le Directeur -->
      <button
        v-if="authStore.currentRole === 'manager'"
        type="button"
        @click="isPricingOpen = true"
        class="hidden sm:inline-flex items-center gap-1.5 px-2.5 py-1 rounded-[3px] bg-amber-500/10 hover:bg-amber-500/20 text-amber-600 dark:text-amber-400 border border-amber-500/30 text-xs font-bold transition cursor-pointer"
        title="Consulter les formules d'abonnement & Bons Plans"
      >
        <Sparkles class="w-3.5 h-3.5 text-amber-500 animate-pulse" />
        <span>Bon Plan</span>
        <span
          v-if="authStore.trialDaysRemaining > 0"
          class="text-[9px] px-1 py-0.2 rounded bg-amber-500/20 font-mono font-semibold"
        >
          {{ authStore.trialDaysRemaining }}j
        </span>
      </button>

      <!-- Bascule Thème Sombre / Clair -->
      <button
        type="button"
        @click="authStore.toggleDarkMode()"
        class="p-1.5 rounded-lg text-muted-foreground hover:text-foreground hover:bg-muted transition"
        :title="
          authStore.isDark ? 'Passer en mode clair' : 'Passer en mode sombre'
        "
      >
        <Sun v-if="authStore.isDark" class="w-4 h-4 text-amber-400" />
        <Moon v-else class="w-4 h-4 text-muted-foreground" />
      </button>

      <!-- Cloche de Notifications -->
      <button
        type="button"
        @click="openNotificationDrawer"
        class="relative p-1.5 rounded-[3px] text-muted-foreground hover:text-foreground hover:bg-muted transition cursor-pointer"
        title="Centre de notifications et alertes"
      >
        <Bell class="w-4 h-4" />
        <!-- Indicateur chiffré affiché uniquement en cas de nouvelles notifications non lues -->
        <span
          v-if="notifStore.unreadCount > 0"
          class="absolute -top-1 -right-1 min-w-[16px] h-4 px-1 rounded-full bg-rose-600 text-white font-bold text-[9px] flex items-center justify-center animate-pulse shadow-xs"
        >
          {{ notifStore.unreadCount > 9 ? "9+" : notifStore.unreadCount }}
        </span>
      </button>

      <!-- Profil Utilisateur : Clic direct pour ouvrir l'interface Profil & Établissement -->
      <div class="flex items-center space-x-1.5">
        <button
          type="button"
          @click="openProfileModal"
          class="flex items-center space-x-2 pl-2.5 pr-1 py-1 rounded-[3px] bg-muted/60 hover:bg-muted hover:border-brand/50 transition cursor-pointer select-none group focus:outline-none focus:ring-2 focus:ring-brand/40"
          :title="`Ouvrir le profil de ${authStore.currentUser.prenom} ${authStore.currentUser.nom}`"
        >
          <!-- Rôle & Nom -->
          <span
            class="hidden md:inline-flex items-center space-x-1.5 text-xs font-semibold text-foreground"
          >
            <span
              >{{ authStore.currentUser.prenom }}
              {{ authStore.currentUser.nom }}</span
            >
          </span>

          <!-- Rond Avatar avec initiale et pastille de rôle -->
          <div
            class="w-6 h-6 rounded-full bg-brand-light text-brand-dark flex items-center justify-center font-bold text-xs uppercase shadow-xs group-hover:scale-105 transition-transform relative border border-brand-border"
          >
            {{ authStore.currentUser.initial }}
            <span
              class="absolute -bottom-0.5 -right-0.5 w-2 h-2 rounded-full border border-card"
              :class="[
                authStore.currentRole === 'secretaire'
                  ? 'bg-brand'
                  : authStore.currentRole === 'manager'
                    ? 'bg-indigo-500'
                    : authStore.currentRole === 'moniteur'
                      ? 'bg-sky-500'
                      : 'bg-purple-500',
              ]"
            ></span>
          </div>
        </button>

        <!-- Bouton Déconnexion / Verrouiller rapide -->
        <button
          type="button"
          @click="authStore.logout()"
          class="hidden sm:inline-flex p-1.5 rounded-[3px] text-rose-500 hover:bg-rose-500/10 hover:text-rose-600 transition"
          title="Verrouiller la session / Se déconnecter"
        >
          <LogOut class="w-4 h-4" />
        </button>
      </div>
    </div>
  </header>

  <!-- Modale État Synchro Neon -->
  <SyncStatusModal v-if="isSyncModalOpen" @close="isSyncModalOpen = false" />

  <!-- Modale Hub QR Codes Mobiles -->
  <SchoolQrHubModal :is-open="isQrModalOpen" @close="isQrModalOpen = false" />

  <!-- Modale Complète de Gestion du Profil Utilisateur, de l'Auto-École et des Agences -->
  <UserProfileAndSchoolModal
    v-if="isProfileModalOpen"
    @close="isProfileModalOpen = false"
  />

  <!-- Modale Notifications & Annonces Réseau Super-Admin -->
  <NetworkNotificationsModal
    :is-open="isNotificationsModalOpen"
    @close="isNotificationsModalOpen = false"
  />

  <!-- Tiroir de Notifications Temps Réel Multi-Tenant -->
  <NotificationDrawer
    :is-open="isDrawerOpen"
    @close="isDrawerOpen = false"
    @open-network-modal="
      isDrawerOpen = false;
      isNotificationsModalOpen = true;
    "
  />

  <!-- Modale Bons Plans & Tarifs d'Abonnement -->
  <SchoolPricingModal :is-open="isPricingOpen" @close="isPricingOpen = false" />
</template>

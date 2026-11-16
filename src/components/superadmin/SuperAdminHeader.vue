<script setup lang="ts">
import { ref, onMounted } from "vue";
import { Search, RefreshCw, Power, Sun, Moon, Bell } from "lucide-vue-next";
import { useAuthStore } from "../../stores/authStore";
import { useNotificationStore } from "../../stores/notificationStore";
import { authBackend } from "../../services/authBackendService";
import NotificationDrawer from "../notifications/NotificationDrawer.vue";
import type { SuperAdminTab } from "./SuperAdminSidebar.vue";

defineProps<{
  searchQuery: string;
  isLoading: boolean;
  latencyMs: number;
  activeTab: SuperAdminTab;
}>();

const emit = defineEmits<{
  (e: "update:searchQuery", val: string): void;
  (e: "refresh"): void;
  (e: "logout"): void;
}>();

const authStore = useAuthStore();
const notifStore = useNotificationStore();
const isDrawerOpen = ref(false);

onMounted(() => {
  const token = authBackend.getActiveSession()?.token;
  if (token) {
    notifStore.initRealtime(token);
    notifStore.loadNotifications();
  }
});

const getPlaceholder = (tab: SuperAdminTab) => {
  switch (tab) {
    case "schools":
      return "Rechercher une auto-école, code, agrément...";
    case "licences":
      return "Filtrer les licences, quotas de machines...";
    case "permissions":
      return "Rechercher un établissement pour les modules...";
    case "users":
      return "Rechercher un utilisateur, email, téléphone...";
    case "candidats":
      return "Rechercher un candidat, dossier, contact...";
    case "audit":
      return "Rechercher dans le journal d'audit...";
    case "alerts":
      return "Filtrer les alertes réseau diffusées...";
    case "announcements":
      return "Rechercher une annonce ou un événement...";
    case "feedbacks":
      return "Rechercher un avis, auto-école, note...";
    default:
      return "Rechercher dans la console Super-Admin...";
  }
};
</script>

<template>
  <header
    class="h-16 bg-card px-4 flex items-center justify-between shrink-0 z-20"
  >
    <!-- Logo & Badges Console -->
    <div class="flex items-center gap-3">
      <div class="flex items-center gap-2">
        <div
          class="w-6 h-6 rounded-[4px] bg-brand/10 border border-brand/30 flex items-center justify-center text-brand font-bold text-xs"
        >
          ▲
        </div>
        <span class="font-bold text-sm tracking-tight text-foreground"
          >ElevenDrive</span
        >
        <span
          class="text-[10px] px-1.5 py-0.5 rounded-[4px] bg-muted text-muted-foreground border border-border font-mono font-medium"
        >
          CONSOLE
        </span>
      </div>

      <span class="text-muted-foreground/40 mx-1">/</span>

      <!-- Badge Environnement -->
      <div
        class="flex items-center gap-1.5 px-2 py-0.5 rounded-full bg-brand/10 border border-brand/20 text-brand text-xs font-mono font-medium"
      >
        <span class="w-1.5 h-1.5 rounded-full bg-brand animate-pulse"></span>
        <span>production</span>
      </div>
    </div>

    <!-- Barre de Recherche Globale avec Raccourci Ctrl+K -->
    <div class="flex-1 max-w-md mx-6 hidden sm:block">
      <div class="relative">
        <Search
          class="w-3.5 h-3.5 text-muted-foreground absolute left-3 top-2.5 pointer-events-none"
        />
        <input
          :value="searchQuery"
          @input="
            emit(
              'update:searchQuery',
              ($event.target as HTMLInputElement).value,
            )
          "
          type="text"
          :placeholder="getPlaceholder(activeTab)"
          class="w-full h-8 pl-8 pr-12 rounded-[4px] bg-background border border-input text-xs text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-ring focus:border-ring transition-colors font-mono"
        />
        <kbd
          class="absolute right-2 top-1.5 text-[10px] font-mono px-1.5 py-0.5 rounded-[3px] bg-muted text-muted-foreground border border-border"
        >
          Ctrl K
        </kbd>
      </div>
    </div>

    <!-- Télémétrie, Bascule Thème & Actions Profil -->
    <div class="flex items-center gap-2">
      <!-- Badge Télémétrie Neontech -->
      <div
        class="flex items-center gap-2 px-2.5 py-1 rounded-[4px] bg-muted border border-border text-xs font-mono text-foreground"
      >
        <span class="text-brand font-semibold">Neontech</span>
        <span class="text-muted-foreground">·</span>
        <span class="text-muted-foreground">{{ latencyMs }}ms</span>
      </div>

      <!-- Bascule Thème Clair / Sombre -->
      <button
        type="button"
        @click="authStore.toggleDarkMode"
        class="p-1.5 rounded-[4px] text-muted-foreground hover:text-foreground hover:bg-muted transition cursor-pointer"
        :title="
          authStore.isDark ? 'Passer en mode clair' : 'Passer en mode sombre'
        "
      >
        <Sun v-if="authStore.isDark" class="w-4 h-4 text-amber-400" />
        <Moon v-else class="w-4 h-4 text-muted-foreground" />
      </button>

      <!-- Actualisation des données -->
      <button
        type="button"
        @click="emit('refresh')"
        title="Actualiser les données"
        class="p-1.5 rounded-[4px] hover:bg-muted text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
      >
        <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': isLoading }" />
      </button>

      <!-- Cloche de Notifications Réseau & Événements -->
      <button
        type="button"
        @click="isDrawerOpen = true"
        title="Centre de notifications et alertes"
        class="relative p-1.5 rounded-[4px] hover:bg-muted text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
      >
        <Bell class="w-3.5 h-3.5" />
        <span
          v-if="notifStore.unreadCount > 0"
          class="absolute -top-1 -right-1 min-w-[15px] h-3.5 px-0.5 rounded-full bg-rose-600 text-white font-bold text-[8px] flex items-center justify-center animate-pulse"
        >
          {{ notifStore.unreadCount > 9 ? '9+' : notifStore.unreadCount }}
        </span>
      </button>

      <div class="h-4 w-px bg-border mx-1"></div>

      <!-- Profil Super-Admin & Déconnexion -->
      <div class="flex items-center gap-1.5">
        <div
          class="w-7 h-7 rounded-full bg-brand/10 border border-brand/30 flex items-center justify-center text-xs font-mono font-bold text-brand"
        >
          SA
        </div>
        <button
          type="button"
          @click="emit('logout')"
          title="Se déconnecter de la console"
          class="p-1.5 rounded-[4px] hover:bg-destructive/10 text-muted-foreground hover:text-destructive transition-colors cursor-pointer"
        >
          <Power class="w-3.5 h-3.5" />
        </button>
      </div>
    </div>
  </header>

  <!-- Tiroir de Notifications Système Super-Admin -->
  <NotificationDrawer
    :is-open="isDrawerOpen"
    @close="isDrawerOpen = false"
  />
</template>

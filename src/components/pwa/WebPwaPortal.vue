<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useAuthStore } from '../../stores/authStore';
import type { UserRole } from '../../types';
import WebPwaHomeView from './WebPwaHomeView.vue';
import WebPwaAuthView from './WebPwaAuthView.vue';

const authStore = useAuthStore();

const currentView = ref<'home' | 'auth'>('home');
const targetRole = ref<UserRole>('eleve');
const deferredPrompt = ref<any>(null);
const canInstallPwa = ref(false);

onMounted(() => {
  // Détection du mode initial depuis authStore ou URL
  if (authStore.webPortalView === 'auth') {
    currentView.value = 'auth';
  } else {
    currentView.value = 'home';
  }

  // Écoute de l'événement PWA avant installation
  if (typeof window !== 'undefined') {
    window.addEventListener('beforeinstallprompt', (e: any) => {
      e.preventDefault();
      deferredPrompt.value = e;
      canInstallPwa.value = true;
    });

    window.addEventListener('appinstalled', () => {
      deferredPrompt.value = null;
      canInstallPwa.value = false;
    });
  }
});

const handleNavigateAuth = (role?: UserRole) => {
  if (role) {
    targetRole.value = role;
  }
  currentView.value = 'auth';
  authStore.setWebPortalView('auth');
};

const handleNavigateHome = () => {
  currentView.value = 'home';
  authStore.setWebPortalView('home');
};

const handleInstallPwa = async () => {
  if (deferredPrompt.value) {
    deferredPrompt.value.prompt();
    const { outcome } = await deferredPrompt.value.userChoice;
    if (outcome === 'accepted') {
      canInstallPwa.value = false;
    }
    deferredPrompt.value = null;
  } else {
    // Si l'invite native n'est pas disponible (ex: iOS Safari), basculer vers le guide
    const homeBtn = document.querySelector('[title="Guide d\'installation"]');
    if (homeBtn) (homeBtn as HTMLElement).click();
  }
};

const handleLoginSuccess = () => {
  // L'utilisateur est désormais authentifié via authStore.isAuthenticated = true
  // App.vue prendra automatiquement le relais pour afficher l'espace dédié
};
</script>

<template>
  <div class="fixed inset-0 z-50 overflow-y-auto bg-background text-foreground font-sans select-none antialiased">
    <transition
      name="fade"
      mode="out-in"
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0 scale-98"
      enter-to-class="opacity-100 scale-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100 scale-100"
      leave-to-class="opacity-0 scale-98"
    >
      <!-- VUE 1 : ACCUEIL PUBLIC DE L'AUTO-ÉCOLE (CONSERVÉ TEL QUEL) -->
      <WebPwaHomeView
        v-if="currentView === 'home'"
        :can-install-pwa="canInstallPwa"
        @navigate-auth="handleNavigateAuth"
        @install-pwa="handleInstallPwa"
      />

      <!-- VUE 2 : AUTHENTIFICATION SÉCURISÉE PWA -->
      <WebPwaAuthView
        v-else-if="currentView === 'auth'"
        :initial-role="targetRole"
        @navigate-home="handleNavigateHome"
        @login-success="handleLoginSuccess"
      />
    </transition>
  </div>
</template>

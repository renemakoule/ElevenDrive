<script setup lang="ts">
import { ref, computed } from 'vue';
import { useAuthStore } from '../../stores/authStore';
import type { UserRole } from '../../types';

// Sous-composants modulaires réutilisables
import PwaHomeHeader from './home/PwaHomeHeader.vue';
import PwaHomeHero from './home/PwaHomeHero.vue';
import PwaHomeRoleSelector from './home/PwaHomeRoleSelector.vue';
import PwaHomeFeatures from './home/PwaHomeFeatures.vue';
import PwaHomeContact from './home/PwaHomeContact.vue';
import PwaHomeFooter from './home/PwaHomeFooter.vue';
import PwaInstallHelpModal from './home/PwaInstallHelpModal.vue';

interface Props {
  canInstallPwa?: boolean;
}

defineProps<Props>();

const emit = defineEmits<{
  (e: 'navigate-auth', role?: UserRole): void;
  (e: 'install-pwa'): void;
}>();

const authStore = useAuthStore();
const isInstallHelpOpen = ref(false);

const schoolInitials = computed(() => {
  const name = authStore.schoolName || 'ElevenDrive';
  return name
    .split(' ')
    .filter(Boolean)
    .slice(0, 2)
    .map((w) => w[0]?.toUpperCase())
    .join('') || 'AE';
});

const whatsappUrl = computed(() => {
  const phone = (authStore.schoolPhone || '+237651727932').replace(/[^0-9]/g, '');
  const msg = encodeURIComponent(
    `Bonjour ${authStore.schoolName}, je vous contacte depuis votre portail web officiel pour obtenir des informations.`
  );
  return `https://wa.me/${phone}?text=${msg}`;
});

const callUrl = computed(() => {
  const phone = (authStore.schoolPhone || '').replace(/\s+/g, '');
  return phone ? `tel:${phone}` : '#';
});

const handleNavigateAuth = (role?: UserRole) => {
  emit('navigate-auth', role);
};

const handleInstallPwa = () => {
  emit('install-pwa');
};
</script>

<template>
  <div class="w-full min-h-screen flex flex-col bg-background text-foreground transition-colors selection:bg-brand selection:text-white">
    <!-- 1. En-tête supérieur officiel de l'Auto-École -->
    <PwaHomeHeader
      class="shrink-0"
      :school-name="authStore.schoolName"
      :school-agrement="authStore.schoolAgrement"
      :school-address="authStore.schoolAddress"
      :school-initials="schoolInitials"
      @navigate-auth="handleNavigateAuth()"
    />

    <!-- 2. Bannière Hero de bienvenue et CTA -->
    <PwaHomeHero
      class="shrink-0"
      :school-name="authStore.schoolName"
      :school-address="authStore.schoolAddress"
      :school-agrement="authStore.schoolAgrement"
      :agences-count="authStore.agencesList.length"
      @navigate-auth="handleNavigateAuth()"
      @install-pwa="handleInstallPwa()"
    />

    <!-- 3. Sélection d'espace dédié par profil (Candidat, Moniteur, Direction) -->
    <PwaHomeRoleSelector
      class="shrink-0"
      @navigate-auth="handleNavigateAuth"
    />

    <!-- 4. Atouts pédagogiques & digitaux de l'Auto-École -->
    <PwaHomeFeatures class="shrink-0" />

    <!-- 5. Coordonnées, contact direct WhatsApp & liste des agences -->
    <PwaHomeContact
      class="shrink-0"
      :school-name="authStore.schoolName"
      :school-phone="authStore.schoolPhone"
      :school-address="authStore.schoolAddress"
      :agences="authStore.agencesList"
      :whatsapp-url="whatsappUrl"
      :call-url="callUrl"
    />

    <!-- 6. Pied de page professionnel & lien vers le guide d'installation -->
    <PwaHomeFooter
      class="mt-auto shrink-0"
      :school-name="authStore.schoolName"
      :school-slug="authStore.schoolSlug"
      @open-install-help="isInstallHelpOpen = true"
    />

    <!-- 7. Modale explicative du guide d'installation PWA -->
    <PwaInstallHelpModal
      :is-open="isInstallHelpOpen"
      @close="isInstallHelpOpen = false"
    />
  </div>
</template>

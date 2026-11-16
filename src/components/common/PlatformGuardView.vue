<script setup lang="ts">
import { ref, watch, onMounted, computed } from 'vue';
import { useAuthStore } from '../../stores/authStore';
import { ShieldAlert, Monitor, Smartphone, ArrowRight, RefreshCw, QrCode } from 'lucide-vue-next';
import { generateQrDataUrl } from '../../utils/qrGenerator';

const authStore = useAuthStore();
const qrDataUrl = ref<string>('');

const currentRoleUrl = computed(() => {
  const links = authStore.schoolAccessLinks;
  if (!links) return '';
  if (authStore.currentRole === 'moniteur') return links.moniteur_link;
  if (authStore.currentRole === 'eleve') return links.eleve_link;
  return links.portal_url;
});

const loadQr = async () => {
  if (currentRoleUrl.value) {
    qrDataUrl.value = await generateQrDataUrl(currentRoleUrl.value, {
      width: 240,
      margin: 1,
      darkColor: authStore.currentRole === 'moniteur' ? '#0284c7' : '#7c3aed',
      lightColor: '#ffffff',
    });
  }
};

watch(() => currentRoleUrl.value, loadQr, { immediate: true });

const switchToDesktop = () => {
  authStore.setPlatform('desktop');
};

const switchToPwa = () => {
  authStore.setPlatform('pwa');
};
</script>

<template>
  <div class="flex-1 flex items-center justify-center p-6 bg-background">
    <div class="bg-card border border-border rounded-xl max-w-lg w-full p-8 text-center space-y-5 shadow-xl">
      <!-- Icône Statut -->
      <div class="mx-auto w-14 h-14 rounded-2xl flex items-center justify-center"
        :class="authStore.currentRole === 'secretaire' ? 'bg-amber-500/10 text-amber-600 border border-amber-500/20' : 'bg-sky-500/10 text-sky-600 border border-sky-500/20'"
      >
        <ShieldAlert v-if="authStore.currentRole === 'secretaire'" class="w-7 h-7" />
        <Smartphone v-else class="w-7 h-7" />
      </div>

      <!-- Titre et Explication -->
      <div class="space-y-2">
        <h2 class="text-base font-bold text-foreground">
          {{ authStore.currentRole === 'secretaire' 
            ? "Accès Restreint : Interface Secrétaire Réservée au Desktop" 
            : `Interface ${authStore.currentUser.roleLabel} Exclusivement Mobile (PWA)` 
          }}
        </h2>

        <p v-if="authStore.currentRole === 'secretaire'" class="text-xs text-muted-foreground leading-relaxed max-w-md mx-auto">
          Pour des impératifs stricts de sécurité financière, de gestion du tiroir-caisse et de liaison matérielle directe avec l'imprimante thermique 80 mm, <strong>le rôle Secrétaire Guichet n'est pas accessible sur le Web</strong> et fonctionne exclusivement sur l'ordinateur Desktop de l'agence.
        </p>

        <p v-else class="text-xs text-muted-foreground leading-relaxed max-w-md mx-auto">
          L'espace <strong>{{ authStore.currentUser.roleLabel }}</strong> est spécialement calibré pour les smartphones et tablettes des moniteurs et apprenants en mobilité (voiture de conduite, entraînement code 24/7). <strong>Il n'est pas prévu pour l'environnement Desktop de bureau</strong>.
        </p>
      </div>

      <!-- Aperçu QR Code pour scan smartphone immédiat -->
      <div v-if="authStore.currentRole !== 'secretaire' && qrDataUrl" class="my-2 flex flex-col items-center">
        <div class="relative group rounded-2xl border-2 border-brand-border/60 bg-white p-2.5 shadow-md">
          <img :src="qrDataUrl" alt="QR Code PWA" class="h-32 w-32 object-contain rounded" />
        </div>
        <span class="mt-2 text-[11px] font-medium text-foreground flex items-center gap-1.5">
          <QrCode class="h-3 w-3 text-brand" />
          Scannez avec l'appareil photo du smartphone 📸
        </span>
      </div>

      <!-- Action Recommandée -->
      <div class="pt-2 flex flex-col sm:flex-row items-center justify-center gap-3">
        <button
          v-if="authStore.currentRole === 'secretaire'"
          type="button"
          @click="switchToDesktop"
          class="w-full sm:w-auto px-4 py-2 rounded-lg bg-foreground hover:bg-foreground/90 text-background font-semibold text-xs flex items-center justify-center gap-2 transition active:scale-95 shadow-sm"
        >
          <Monitor class="w-3.5 h-3.5" />
          <span>Basculer sur l'Application Desktop</span>
        </button>

        <button
          v-else
          type="button"
          @click="switchToPwa"
          class="w-full sm:w-auto px-4 py-2 rounded-lg bg-brand-dark hover:bg-brand text-white font-semibold text-xs flex items-center justify-center gap-2 transition active:scale-95 shadow-sm"
        >
          <Smartphone class="w-3.5 h-3.5" />
          <span>Ouvrir en Mode PWA Mobile</span>
        </button>

        <button
          type="button"
          @click="authStore.setRole('manager')"
          class="w-full sm:w-auto px-4 py-2 rounded-lg border border-border bg-card hover:bg-muted text-foreground font-medium text-xs transition"
        >
          Changer de profil (Manager)
        </button>
      </div>
    </div>
  </div>
</template>

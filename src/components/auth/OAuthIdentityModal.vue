<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { X, ArrowRight, ShieldCheck, User } from 'lucide-vue-next';
import { authBackend } from '../../services/authBackendService';

type ProviderId = 'google' | 'facebook' | 'linkedin' | 'apple';

interface Props {
  isOpen: boolean;
  provider: ProviderId | null;
  isLoading?: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'confirm', data: { provider: ProviderId; email: string; nom: string; prenoms: string }): void;
}>();

const email = ref('');
const nom = ref('');
const prenoms = ref('');

const detectedAccounts = computed(() => {
  return authBackend.getKnownAccounts();
});

const selectExistingAccount = (account: any) => {
  if (account.email) email.value = account.email;
  if (account.nom) nom.value = account.nom;
  if (account.prenom) prenoms.value = account.prenom;
};

const providerDetails = computed(() => {
  switch (props.provider) {
    case 'google':
      return {
        name: 'Google',
        color: '#EA4335',
        domain: 'gmail.com',
        placeholderEmail: 'votre.adresse@gmail.com',
      };
    case 'facebook':
      return {
        name: 'Facebook',
        color: '#1877F2',
        domain: 'facebook.com',
        placeholderEmail: 'votre.adresse@facebook.com',
      };
    case 'linkedin':
      return {
        name: 'LinkedIn',
        color: '#0A66C2',
        domain: 'linkedin.com',
        placeholderEmail: 'votre.email.pro@domaine.com',
      };
    case 'apple':
      return {
        name: 'Apple',
        color: '#000000',
        domain: 'icloud.com',
        placeholderEmail: 'identifiant@icloud.com',
      };
    default:
      return {
        name: 'Fournisseur',
        color: '#059669',
        domain: 'email.com',
        placeholderEmail: 'votre.adresse@domaine.com',
      };
  }
});

// Réinitialiser les champs à vide lors du choix d'un provider pour que l'utilisateur saisisse ses propres identifiants réels
watch(
  () => props.provider,
  (p) => {
    if (p) {
      email.value = '';
      nom.value = '';
      prenoms.value = '';
    }
  },
  { immediate: true }
);

const handleSubmit = () => {
  if (!props.provider || !email.value.trim()) return;
  emit('confirm', {
    provider: props.provider,
    email: email.value.trim().toLowerCase(),
    nom: nom.value.trim().toUpperCase(),
    prenoms: prenoms.value.trim(),
  });
};
</script>

<template>
  <div
    v-if="isOpen && provider"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-sm animate-in fade-in duration-200"
    role="dialog"
    aria-modal="true"
  >
    <div
      class="w-full max-w-sm overflow-hidden rounded-2xl border border-border bg-card shadow-2xl animate-in zoom-in-95 duration-200"
    >
      <!-- En-tête avec identité visuelle du provider -->
      <header class="flex items-center justify-between border-b border-border px-5 py-4 bg-muted/40">
        <div class="flex items-center gap-2.5">
          <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-background border border-border shadow-sm">
            <!-- Google -->
            <svg v-if="provider === 'google'" class="h-4 w-4" viewBox="0 0 48 48">
              <path fill="#EA4335" d="M24 9.5c3.54 0 6.71 1.22 9.21 3.6l6.85-6.85C35.9 2.38 30.47 0 24 0 14.62 0 6.51 5.38 2.56 13.22l7.98 6.19C12.43 13.72 17.74 9.5 24 9.5z"/>
              <path fill="#4285F4" d="M46.98 24.55c0-1.57-.15-3.09-.38-4.55H24v9.02h12.94c-.58 2.96-2.26 5.48-4.78 7.18l7.73 6c4.51-4.18 7.09-10.36 7.09-17.65z"/>
              <path fill="#FBBC05" d="M10.53 28.59c-.48-1.45-.76-2.99-.76-4.59s.27-3.14.76-4.59l-7.98-6.19C.92 16.46 0 20.12 0 24c0 3.88.92 7.54 2.56 10.78l7.97-6.19z"/>
              <path fill="#34A853" d="M24 48c6.48 0 11.93-2.13 15.89-5.81l-7.73-6c-2.15 1.45-4.92 2.3-8.16 2.3-6.26 0-11.57-4.22-13.47-9.91l-7.98 6.19C6.51 42.62 14.62 48 24 48z"/>
            </svg>
            <!-- Facebook -->
            <svg v-else-if="provider === 'facebook'" class="h-4 w-4" viewBox="0 0 24 24" fill="#1877F2">
              <path d="M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z"/>
            </svg>
            <!-- LinkedIn -->
            <svg v-else-if="provider === 'linkedin'" class="h-4 w-4" viewBox="0 0 24 24" fill="#0A66C2">
              <path d="M19 3a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h14m-.5 15.5v-5.3a3.26 3.26 0 0 0-3.26-3.26c-.85 0-1.84.52-2.28 1.3v-1.11h-2.79v8.37h2.79v-4.93c0-.77.62-1.4 1.39-1.4a1.4 1.4 0 0 1 1.4 1.4v4.93h2.75M6.46 8.76a1.62 1.62 0 1 0 0-3.24 1.62 1.62 0 0 0 0 3.24m1.39 9.74v-8.37H5.07v8.37h2.78z"/>
            </svg>
            <!-- Apple -->
            <svg v-else class="h-4 w-4" viewBox="0 0 24 24" fill="currentColor">
              <path d="M17.05 20.28c-.98.95-2.05.8-3.08.35-1.09-.46-2.09-.48-3.24 0-1.44.62-2.2.44-3.06-.35C2.79 15.25 3.51 7.59 9.05 7.31c1.35.07 2.29.74 3.08.8 1.18-.24 2.31-.93 3.57-.84 1.51.12 2.65.72 3.4 1.8-3.12 1.87-2.38 5.98.48 7.13-.57 1.5-1.31 2.99-2.54 4.09zM12.03 7.25c-.15-2.23 1.66-4.07 3.74-4.25.29 2.58-2.34 4.5-3.74 4.25z"/>
            </svg>
          </div>
          <div>
            <h3 class="text-sm font-semibold text-foreground">
              Continuer avec {{ providerDetails.name }}
            </h3>
            <p class="text-[11px] text-muted-foreground">Authentification sécurisée</p>
          </div>
        </div>
        <button
          type="button"
          @click="emit('close')"
          class="rounded-lg p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground transition"
          aria-label="Fermer"
        >
          <X class="h-4 w-4" />
        </button>
      </header>

      <!-- Formulaire de confirmation du compte social -->
      <form @submit.prevent="handleSubmit" class="p-5 space-y-4 text-xs">
        <div class="rounded-lg border border-border/80 bg-muted/30 p-2.5 text-[11px] text-muted-foreground flex items-center gap-2">
          <ShieldCheck class="h-4 w-4 text-brand shrink-0" />
          <span>
            Si vous avez déjà un compte, vous serez connecté automatiquement. Sinon, vos informations seront pré-remplies pour finaliser votre inscription.
          </span>
        </div>

        <!-- Comptes existants mémorisés (pour test rapide) -->
        <div v-if="detectedAccounts.length > 0" class="rounded-lg border border-border/60 bg-background/50 p-2.5">
          <span class="text-[10px] uppercase font-semibold text-muted-foreground tracking-wider block mb-1.5">
            Comptes existants détectés :
          </span>
          <div class="flex flex-wrap gap-1.5">
            <button
              v-for="acc in detectedAccounts"
              :key="acc.email"
              type="button"
              @click="selectExistingAccount(acc)"
              class="inline-flex items-center gap-1.5 px-2 py-1 rounded text-[11px] border transition"
              :class="email.toLowerCase() === acc.email?.toLowerCase() ? 'bg-brand-light border-brand-border text-brand-dark font-medium' : 'bg-muted/40 border-border text-muted-foreground hover:text-foreground'"
            >
              <span class="w-1.5 h-1.5 rounded-full bg-brand"></span>
              <span>{{ acc.prenom }} {{ acc.nom }} ({{ acc.email }})</span>
            </button>
          </div>
        </div>

        <div>
          <label class="block font-medium text-foreground mb-1">
            Adresse email {{ providerDetails.name }} *
          </label>
          <input
            v-model="email"
            type="email"
            required
            :placeholder="'Ex: ' + providerDetails.placeholderEmail"
            class="w-full h-9 px-3 rounded-lg border border-border bg-background text-foreground text-xs focus:border-brand focus:outline-none"
          />
        </div>

        <div class="grid grid-cols-2 gap-2.5">
          <div>
            <label class="block font-medium text-foreground mb-1">Prénom</label>
            <input
              v-model="prenoms"
              type="text"
              placeholder="Votre prénom"
              class="w-full h-9 px-3 rounded-lg border border-border bg-background text-foreground text-xs focus:border-brand focus:outline-none"
            />
          </div>
          <div>
            <label class="block font-medium text-foreground mb-1">Nom</label>
            <input
              v-model="nom"
              type="text"
              placeholder="Votre nom de famille"
              class="w-full h-9 px-3 rounded-lg border border-border bg-background text-foreground text-xs uppercase focus:border-brand focus:outline-none"
            />
          </div>
        </div>

        <!-- Boutons d'action -->
        <div class="flex items-center gap-2 pt-2">
          <button
            type="button"
            @click="emit('close')"
            class="h-9 flex-1 rounded-lg border border-border bg-background text-xs font-medium text-foreground hover:bg-muted transition"
          >
            Annuler
          </button>
          <button
            type="submit"
            :disabled="isLoading || !email.trim()"
            class="h-9 flex-[2] inline-flex items-center justify-center gap-1.5 rounded-lg bg-brand text-xs font-semibold text-white hover:bg-brand-dark transition disabled:opacity-50"
          >
            <span v-if="isLoading">Vérification...</span>
            <template v-else>
              <span>Continuer</span>
              <ArrowRight class="h-3.5 w-3.5" />
            </template>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

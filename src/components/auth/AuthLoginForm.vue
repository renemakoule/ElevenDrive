<script setup lang="ts">
import { ref, reactive } from 'vue';
import { Phone, Lock, Eye, EyeOff, ArrowRight } from 'lucide-vue-next';
import type { LoginCredentials } from '../../types';
import OAuthProviderButtons from './OAuthProviderButtons.vue';

type ProviderId = 'google' | 'facebook' | 'linkedin' | 'apple';

interface Props {
  isLoading?: boolean;
}

defineProps<Props>();

const emit = defineEmits<{
  (e: 'submit', credentials: LoginCredentials): void;
  (e: 'oauth', provider: ProviderId): void;
  (e: 'clear-error'): void;
}>();

const loginForm = reactive<LoginCredentials>({
  identifiant: '',
  mot_de_passe: '',
});

const showPassword = ref(false);

const inputBase =
  'h-10 w-full rounded-[5px] border border-border bg-background text-sm text-foreground ' +
  'placeholder:text-muted-foreground/60 transition-colors hover:border-foreground/30 ' +
  'focus:outline-none focus:border-brand focus:ring-2 focus:ring-brand/20';
const labelClass = 'mb-1.5 block text-xs font-medium text-foreground';
const primaryBtn =
  'inline-flex h-10 w-full items-center justify-center gap-2 rounded-[5px] bg-brand text-sm font-semibold text-white ' +
  'transition hover:bg-brand-dark focus:outline-none focus-visible:ring-2 focus-visible:ring-brand/40 ' +
  'focus-visible:ring-offset-2 focus-visible:ring-offset-background active:scale-[0.99] disabled:opacity-50 disabled:pointer-events-none';

const handleSubmit = () => {
  if (!loginForm.identifiant || !loginForm.mot_de_passe) return;
  emit('submit', { ...loginForm, identifiant: loginForm.identifiant.trim() });
};
</script>

<template>
  <div class="mt-8">
    <!-- Grille des 4 fournisseurs OAuth -->
    <OAuthProviderButtons
      :disabled="isLoading"
      :show-last-badge="true"
      @select="emit('oauth', $event)"
    />

    <!-- Séparateur visuel -->
    <div class="my-5 flex items-center gap-3">
      <span class="h-px flex-1 bg-border"></span>
      <span class="text-[11px] font-medium text-muted-foreground uppercase tracking-wider">
        Ou avec vos accès
      </span>
      <span class="h-px flex-1 bg-border"></span>
    </div>

    <!-- Formulaire classique -->
    <form @submit.prevent="handleSubmit" class="space-y-4">
      <div>
        <label for="login-id" :class="labelClass">Identifiant professionnel</label>
        <div class="relative">
          <Phone class="pointer-events-none absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
          <input
            id="login-id"
            v-model="loginForm.identifiant"
            @input="emit('clear-error')"
            type="text"
            required
            autocomplete="username"
            placeholder="Ex: +237 651 727 932 ou email"
            :class="[inputBase, 'pl-9 pr-3']"
          />
        </div>
        <p class="mt-1 text-[10px] text-muted-foreground">
          Numéro WhatsApp ou adresse email utilisé lors de votre enregistrement.
        </p>
      </div>

      <div>
        <label for="login-password" :class="labelClass">Mot de passe</label>
        <div class="relative">
          <Lock class="pointer-events-none absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
          <input
            id="login-password"
            v-model="loginForm.mot_de_passe"
            @input="emit('clear-error')"
            :type="showPassword ? 'text' : 'password'"
            required
            autocomplete="current-password"
            placeholder="Votre mot de passe"
            :class="[inputBase, 'pl-9 pr-10']"
          />
          <button
            type="button"
            @click="showPassword = !showPassword"
            :aria-label="showPassword ? 'Masquer le mot de passe' : 'Afficher le mot de passe'"
            class="absolute right-2 top-2 rounded p-1 text-muted-foreground transition hover:text-foreground focus:outline-none focus-visible:ring-2 focus-visible:ring-brand/30"
          >
            <EyeOff v-if="showPassword" class="h-4 w-4" />
            <Eye v-else class="h-4 w-4" />
          </button>
        </div>
      </div>

      <button
        type="submit"
        :disabled="isLoading || !loginForm.identifiant || !loginForm.mot_de_passe"
        :class="[primaryBtn, 'mt-6']"
      >
        <span v-if="isLoading">Vérification des accès...</span>
        <template v-else>
          <span>Se connecter</span>
          <ArrowRight class="h-4 w-4" />
        </template>
      </button>
    </form>
  </div>
</template>

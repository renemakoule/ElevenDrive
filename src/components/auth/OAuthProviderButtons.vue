<script setup lang="ts">
import { ref } from 'vue';

type ProviderId = 'google' | 'facebook' | 'linkedin' | 'apple';

interface Props {
  disabled?: boolean;
  showLastBadge?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
  showLastBadge: true,
});

const emit = defineEmits<{
  (e: 'select', provider: ProviderId): void;
}>();

interface OAuthProvider {
  id: ProviderId;
  label: string;
}

const providers: OAuthProvider[] = [
  { id: 'google', label: 'Google' },
  { id: 'facebook', label: 'Facebook' },
  { id: 'linkedin', label: 'LinkedIn' },
  { id: 'apple', label: 'Apple' },
];

const LAST_PROVIDER_KEY = 'autoschool:last_oauth_provider';
const lastProvider = ref<ProviderId | null>(
  typeof localStorage !== 'undefined'
    ? (localStorage.getItem(LAST_PROVIDER_KEY) as ProviderId | null)
    : null
);

const handleSelect = (provider: ProviderId) => {
  lastProvider.value = provider;
  if (typeof localStorage !== 'undefined') {
    try {
      localStorage.setItem(LAST_PROVIDER_KEY, provider);
    } catch {}
  }
  emit('select', provider);
};
</script>

<template>
  <div class="grid grid-cols-2 gap-2.5">
    <button
      v-for="provider in providers"
      :key="provider.id"
      type="button"
      :disabled="disabled"
      @click="handleSelect(provider.id)"
      class="relative inline-flex h-10 items-center justify-center gap-2 rounded-[5px] border border-border bg-card px-3 text-sm font-medium text-foreground transition hover:bg-muted/10 hover:border-foreground/20 focus:outline-none focus-visible:ring-2 focus-visible:ring-brand/30 active:scale-[0.99] disabled:opacity-50 disabled:pointer-events-none"
    >
      <!-- Google -->
      <svg v-if="provider.id === 'google'" class="h-[18px] w-[18px] shrink-0" viewBox="0 0 48 48" aria-hidden="true">
        <path fill="#EA4335" d="M24 9.5c3.54 0 6.71 1.22 9.21 3.6l6.85-6.85C35.9 2.38 30.47 0 24 0 14.62 0 6.51 5.38 2.56 13.22l7.98 6.19C12.43 13.72 17.74 9.5 24 9.5z"/>
        <path fill="#4285F4" d="M46.98 24.55c0-1.57-.15-3.09-.38-4.55H24v9.02h12.94c-.58 2.96-2.26 5.48-4.78 7.18l7.73 6c4.51-4.18 7.09-10.36 7.09-17.65z"/>
        <path fill="#FBBC05" d="M10.53 28.59c-.48-1.45-.76-2.99-.76-4.59s.27-3.14.76-4.59l-7.98-6.19C.92 16.46 0 20.12 0 24c0 3.88.92 7.54 2.56 10.78l7.97-6.19z"/>
        <path fill="#34A853" d="M24 48c6.48 0 11.93-2.13 15.89-5.81l-7.73-6c-2.15 1.45-4.92 2.3-8.16 2.3-6.26 0-11.57-4.22-13.47-9.91l-7.98 6.19C6.51 42.62 14.62 48 24 48z"/>
      </svg>

      <!-- Facebook -->
      <svg v-else-if="provider.id === 'facebook'" class="h-[18px] w-[18px] shrink-0" viewBox="0 0 24 24" fill="#1877F2" aria-hidden="true">
        <path d="M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z"/>
      </svg>

      <!-- LinkedIn -->
      <svg v-else-if="provider.id === 'linkedin'" class="h-[18px] w-[18px] shrink-0" viewBox="0 0 24 24" fill="#0A66C2" aria-hidden="true">
        <path d="M19 3a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h14m-.5 15.5v-5.3a3.26 3.26 0 0 0-3.26-3.26c-.85 0-1.84.52-2.28 1.3v-1.11h-2.79v8.37h2.79v-4.93c0-.77.62-1.4 1.39-1.4a1.4 1.4 0 0 1 1.4 1.4v4.93h2.75M6.46 8.76a1.62 1.62 0 1 0 0-3.24 1.62 1.62 0 0 0 0 3.24m1.39 9.74v-8.37H5.07v8.37h2.78z"/>
      </svg>

      <!-- Apple -->
      <svg v-else class="h-[18px] w-[18px] shrink-0" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
        <path d="M17.05 20.28c-.98.95-2.05.8-3.08.35-1.09-.46-2.09-.48-3.24 0-1.44.62-2.2.44-3.06-.35C2.79 15.25 3.51 7.59 9.05 7.31c1.35.07 2.29.74 3.08.8 1.18-.24 2.31-.93 3.57-.84 1.51.12 2.65.72 3.4 1.8-3.12 1.87-2.38 5.98.48 7.13-.57 1.5-1.31 2.99-2.54 4.09zM12.03 7.25c-.15-2.23 1.66-4.07 3.74-4.25.29 2.58-2.34 4.5-3.74 4.25z"/>
      </svg>

      <span>{{ provider.label }}</span>

      <!-- Badge dernier fournisseur utilisé -->
      <span
        v-if="showLastBadge && lastProvider === provider.id"
        class="absolute -top-2 right-2 rounded-[5px] bg-brand px-1.5 py-0.5 text-[9px] font-semibold leading-none text-white"
      >
        Dernier
      </span>
    </button>
  </div>
</template>

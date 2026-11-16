<script setup lang="ts">
import { HelpCircle } from 'lucide-vue-next';

interface Props {
  schoolName: string;
  schoolSlug: string;
}

defineProps<Props>();

const emit = defineEmits<{
  (e: 'open-install-help'): void;
}>();

// Année en cours, mise à jour automatiquement
const currentYear = new Date().getFullYear();
</script>

<template>
  <footer class="mt-auto shrink-0 w-full border-t border-border bg-card py-6">
    <div class="max-w-6xl mx-auto px-4 sm:px-6 flex flex-col sm:flex-row items-center justify-between gap-4 text-xs text-muted-foreground text-center sm:text-left">
      <div>
        <span class="font-bold text-foreground">{{ schoolName }}</span>
        <span class="mx-1.5">·</span>
        <span>Identifiant portail : <code class="font-mono text-[11px] text-brand font-semibold">#{{ schoolSlug || 'auto-ecole' }}</code></span>
      </div>

      <div class="flex flex-wrap items-center justify-center gap-x-3 gap-y-1">
        <span>© {{ currentYear }} · Propulsé par <strong class="text-foreground">GAS</strong></span>
        <span class="text-muted-foreground/60">·</span>
        <button
          type="button"
          @click="emit('open-install-help')"
          class="hover:text-brand hover:underline inline-flex items-center gap-1 font-medium"
        >
          <HelpCircle class="w-3.5 h-3.5" />
          <span>Guide PWA</span>
        </button>
      </div>
    </div>
  </footer>
</template>
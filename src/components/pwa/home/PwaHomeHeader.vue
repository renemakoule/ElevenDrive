<script setup lang="ts">
import { ref } from 'vue';
import { useAuthStore } from '../../../stores/authStore';
import { useLanguageStore } from '../../../stores/languageStore';
import type { SupportedLanguage } from '../../../services/deepLTranslationService';
import { CheckCircle2, Globe, Sun, Moon, ArrowRight } from 'lucide-vue-next';

interface Props {
  schoolName: string;
  schoolAgrement: string;
  schoolAddress: string;
  schoolInitials: string;
}

defineProps<Props>();

const emit = defineEmits<{
  (e: 'navigate-auth'): void;
}>();

const authStore = useAuthStore();
const languageStore = useLanguageStore();

const isLangMenuOpen = ref(false);

const languages: Array<{ code: SupportedLanguage; label: string; flag: string }> = [
  { code: 'fr', label: 'Français', flag: '🇫🇷' },
  { code: 'en', label: 'English', flag: '🇬🇧' },
  { code: 'es', label: 'Español', flag: '🇪🇸' },
  { code: 'ar', label: 'العربية', flag: '🇸🇦' },
];

const selectLanguage = (code: SupportedLanguage) => {
  languageStore.setLanguage(code);
  isLangMenuOpen.value = false;
};
</script>

<template>
  <!-- Conteneur sticky : ne capte pas les clics autour de la barre -->
  <header class="sticky top-0 z-40 px-3 sm:px-6 pt-3 pointer-events-none">
    <!-- Barre centrée, bordures arrondies 4px -->
    <div
      class="pointer-events-auto max-w-5xl mx-auto h-16 px-3 sm:px-4 flex items-center justify-between gap-3 rounded-[4px] border border-border bg-card/90 backdrop-blur-md"
    >
      <!-- Logo & Identité officielle de l'Auto-école -->
      <div class="flex items-center gap-3 min-w-0">
        <div class="relative shrink-0">
          <div class="w-10 h-10 rounded-full bg-gradient-to-br from-brand to-brand-dark text-white flex items-center justify-center font-bold text-sm tracking-wider border border-brand-border">
            {{ schoolInitials }}
          </div>
          <span class="absolute -bottom-1 -right-1 w-3.5 h-3.5 rounded-full bg-emerald-500 border-2 border-card" title="Établissement vérifié"></span>
        </div>

        <div class="min-w-0">
          <div class="flex items-center gap-1.5 flex-wrap">
            <h1 class="text-sm sm:text-base font-bold text-foreground truncate max-w-[200px] sm:max-w-xs">
              {{ schoolName }}
            </h1>
            <span class="hidden sm:inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-semibold bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 border border-emerald-500/20">
              <CheckCircle2 class="w-3 h-3" />
              <span>Agréée</span>
            </span>
          </div>
        </div>
      </div>

      <!-- Contrôles : Langue, Thème & Connexion -->
      <div class="flex items-center gap-2">
        <!-- Sélecteur de langue -->
        <div class="relative">
          <button
            type="button"
            @click="isLangMenuOpen = !isLangMenuOpen"
            class="inline-flex items-center gap-1 px-2.5 py-1.5 rounded-[4px] border border-border bg-background text-xs font-medium text-foreground hover:bg-muted transition"
            title="Changer la langue"
          >
            <Globe class="w-3.5 h-3.5 text-muted-foreground" />
            <span class="text-[11px] font-mono uppercase">{{ languageStore.currentLanguage }}</span>
          </button>

          <!-- Dropdown langues -->
          <div
            v-if="isLangMenuOpen"
            class="absolute right-0 mt-1.5 w-36 rounded-[4px] border border-border bg-card py-1 z-50 animate-in fade-in zoom-in-95 duration-150"
          >
            <button
              v-for="l in languages"
              :key="l.code"
              type="button"
              @click="selectLanguage(l.code)"
              class="w-full flex items-center justify-between px-3 py-1.5 text-xs text-left text-foreground hover:bg-muted transition"
              :class="languageStore.currentLanguage === l.code ? 'font-bold text-brand bg-brand/5' : ''"
            >
              <span>{{ l.label }}</span>
              <span>{{ l.flag }}</span>
            </button>
          </div>
        </div>

        <!-- Bascule Thème Sombre / Clair -->
        <button
          type="button"
          @click="authStore.toggleDarkMode()"
          class="p-2 rounded-[4px] border border-border bg-background text-foreground hover:bg-muted transition"
          :title="authStore.isDark ? 'Mode clair' : 'Mode sombre'"
        >
          <Sun v-if="authStore.isDark" class="w-4 h-4 text-amber-400" />
          <Moon v-else class="w-4 h-4 text-slate-600" />
        </button>

        <!-- Bouton Connexion -->
        <button
          type="button"
          @click="emit('navigate-auth')"
          class="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-[4px] bg-brand text-xs font-semibold text-white hover:bg-brand-dark transition shadow-xs"
        >
          <span>Connexion</span>
          <ArrowRight class="w-3.5 h-3.5" />
        </button>
      </div>
    </div>
  </header>
</template>
<script setup lang="ts">
import { ref } from "vue";
import { useAuthStore } from "../../stores/authStore";
import { useLanguageStore } from "../../stores/languageStore";
import type { SupportedLanguage } from "../../services/deepLTranslationService";
import {
  Sparkles,
  DownloadCloud,
  Settings,
  Sun,
  Moon,
  Check,
} from "lucide-vue-next";

const emit = defineEmits<{
  (e: "select-nav", id: string): void;
}>();

const authStore = useAuthStore();
const languageStore = useLanguageStore();

const isLangMenuOpen = ref(false);

const chooseLanguage = (lang: SupportedLanguage) => {
  languageStore.setLanguage(lang);
  isLangMenuOpen.value = false;
};

const closeLangMenu = () => {
  isLangMenuOpen.value = false;
};
</script>

<template>
  <div
    v-if="isLangMenuOpen"
    class="fixed inset-0 z-40 bg-transparent"
    @click="closeLangMenu"
  ></div>

  <div
    class="shrink-0 p-2 space-y-1 bg-card border-t border-border/40 relative z-30"
  >
    <!-- Bouton Bon Plan (Directeur exclusivement) -->
    <button
      v-if="authStore.currentRole === 'manager'"
      type="button"
      @click="emit('select-nav', 'bon-plan')"
      class="w-full flex items-center space-x-2 px-3 py-2 rounded-[4px] text-[11px] font-semibold transition active:scale-[0.98]"
      :class="
        authStore.activeSection === 'bon-plan'
          ? 'bg-amber-500/20 text-amber-700 dark:text-amber-300 border border-amber-400 dark:border-amber-600 font-bold shadow-xs'
          : 'bg-amber-500/10 hover:bg-amber-500/20 text-amber-800 dark:text-amber-300 border border-amber-300/80 dark:border-amber-700/60'
      "
      :title="languageStore.t('Découvrir les offres & tarifs ElevenDrive')"
    >
      <Sparkles
        class="w-3.5 h-3.5 shrink-0 text-amber-600 dark:text-amber-400"
      />
      <span v-if="!authStore.isSidebarCollapsed" class="truncate">
        {{ languageStore.t("Bon Plan") }}
      </span>
      <span
        v-if="!authStore.isSidebarCollapsed"
        class="ml-auto text-[9px] font-bold px-1.5 py-0.5 rounded bg-amber-600 text-white leading-none"
      >
        -15%
      </span>
    </button>

    <!-- Bouton Installer l'App -->
    <button
      v-else
      type="button"
      @click="authStore.setPlatform('pwa')"
      class="w-full flex items-center space-x-2 px-3 py-2 rounded-[4px] text-[11px] font-medium bg-brand-light hover:bg-brand-light/80 text-brand-dark border border-brand-border transition"
      :title="languageStore.t('Installer ElevenDrive PWA sur mobile')"
    >
      <DownloadCloud class="w-3.5 h-3.5 shrink-0 text-brand" />
      <span
        v-if="!authStore.isSidebarCollapsed"
        class="truncate font-semibold"
      >
        {{ languageStore.t("Installer l'App") }}
      </span>
    </button>

    <!-- Paramètres -->
    <button
      type="button"
      class="w-full flex items-center space-x-2.5 px-3 py-1.5 rounded-[4px] text-muted-foreground hover:text-foreground hover:bg-muted text-[11px] transition"
    >
      <Settings class="w-3.5 h-3.5 text-muted-foreground" />
      <span v-if="!authStore.isSidebarCollapsed">{{
        languageStore.t("Paramètres")
      }}</span>
    </button>

    <!-- Thème -->
    <button
      type="button"
      @click="authStore.toggleDarkMode()"
      class="w-full flex items-center space-x-2.5 px-3 py-1.5 rounded-[4px] text-muted-foreground hover:text-foreground hover:bg-muted text-[11px] transition"
      :title="
        authStore.isDark
          ? languageStore.t('Passer en mode clair')
          : languageStore.t('Passer en mode sombre')
      "
    >
      <Sun v-if="authStore.isDark" class="w-3.5 h-3.5 text-amber-400" />
      <Moon v-else class="w-3.5 h-3.5 text-muted-foreground" />
      <span v-if="!authStore.isSidebarCollapsed">
        {{
          authStore.isDark
            ? languageStore.t("Thème Sombre")
            : languageStore.t("Thème Clair")
        }}
      </span>
    </button>

    <!-- Langue -->
    <div class="relative">
      <button
        type="button"
        @click="isLangMenuOpen = !isLangMenuOpen"
        class="w-full flex items-center space-x-2.5 px-3 py-1.5 rounded-[4px] text-muted-foreground hover:text-foreground hover:bg-muted text-[11px] transition cursor-pointer"
        :class="{ 'bg-muted text-foreground': isLangMenuOpen }"
        :title="languageStore.t('Langue')"
      >
        <span class="text-sm shrink-0">{{
          languageStore.currentMeta.flag
        }}</span>
        <span v-if="!authStore.isSidebarCollapsed" class="truncate">{{
          languageStore.t("Langue")
        }}</span>
        <span
          v-if="!authStore.isSidebarCollapsed"
          class="text-[10px] font-mono text-muted-foreground ml-auto uppercase"
        >
          {{ languageStore.currentLanguage }} ▾
        </span>
      </button>

      <div
        v-if="isLangMenuOpen"
        class="absolute z-50 rounded-lg border border-border bg-card p-1 shadow-2xl text-xs animate-in fade-in zoom-in-95 duration-100"
        :class="
          authStore.isSidebarCollapsed
            ? 'left-full bottom-0 ml-2 w-44'
            : 'bottom-full left-0 mb-1 w-full'
        "
      >
        <div
          class="px-2 py-1 text-[10px] font-semibold text-muted-foreground border-b border-border mb-1"
        >
          {{ languageStore.t("Langue") }}
        </div>
        <button
          v-for="lang in languageStore.availableLanguages"
          :key="lang.code"
          type="button"
          @click="chooseLanguage(lang.code)"
          class="w-full flex items-center justify-between px-2.5 py-1.5 rounded text-left transition hover:bg-muted text-foreground cursor-pointer"
          :class="{
            'font-semibold text-brand bg-brand-light/40':
              languageStore.currentLanguage === lang.code,
          }"
        >
          <div class="flex items-center gap-2 text-xs">
            <span class="text-base">{{ lang.flag }}</span>
            <span>{{ lang.name }}</span>
          </div>
          <Check
            v-if="languageStore.currentLanguage === lang.code"
            class="h-3.5 w-3.5 text-brand shrink-0"
          />
        </button>
      </div>
    </div>
  </div>
</template>


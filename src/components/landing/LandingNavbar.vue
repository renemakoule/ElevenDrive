<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { Menu, X, Sun, Moon, ChevronDown, Check } from "lucide-vue-next";
import appLogo from "@root/app-icon.png";
import { useAuthStore } from "../../stores/authStore";
import { useLanguageStore } from "../../stores/languageStore";
import type { SupportedLanguage } from "../../services/deepLTranslationService";

const emit = defineEmits<{
  (e: "navigate-auth"): void;
  (e: "scroll-to", id: string): void;
}>();

const authStore = useAuthStore();
const languageStore = useLanguageStore();

const isMenuOpen = ref(false);
const isLangMenuOpen = ref(false);

const links = computed(() => [
  { id: "about", label: languageStore.t("À propos") },
  { id: "solutions", label: languageStore.t("Solutions") },
  { id: "pricing", label: languageStore.t("Tarifs") },
]);

const goTo = (id: string) => {
  isMenuOpen.value = false;
  isLangMenuOpen.value = false;
  emit("scroll-to", id);
};

const goAuth = () => {
  isMenuOpen.value = false;
  isLangMenuOpen.value = false;
  emit("navigate-auth");
};

const toggleLangMenu = () => {
  isLangMenuOpen.value = !isLangMenuOpen.value;
};

const selectLanguage = (code: SupportedLanguage) => {
  languageStore.setLanguage(code);
  isLangMenuOpen.value = false;
};

const handleDocumentClick = (e: MouseEvent) => {
  const target = e.target as HTMLElement | null;
  if (target && !target.closest("#landing-lang-dropdown")) {
    isLangMenuOpen.value = false;
  }
};

onMounted(() => {
  document.addEventListener("click", handleDocumentClick);
});

onUnmounted(() => {
  document.removeEventListener("click", handleDocumentClick);
});
</script>

<template>
  <header
    class="sticky top-0 z-40 border-b border-border bg-background/80 backdrop-blur-md"
  >
    <nav
      class="mx-auto flex h-14 max-w-6xl items-center justify-between px-5 sm:px-6"
      aria-label="Navigation principale"
    >
      <!-- Logo -->
      <button
        type="button"
        class="flex items-center gap-2 focus:outline-none"
        @click="goTo('top')"
        aria-label="ElevenDrive, retour en haut"
      >
        <img
          :src="appLogo"
          alt=""
          class="h-6 w-6 rounded-[4px] object-contain"
        />
        <span class="text-[15px] font-semibold tracking-tight text-foreground">
          ElevenDrive<span class="text-brand">.</span>
        </span>
      </button>

      <!-- Liens (desktop) -->
      <ul class="hidden items-center gap-1 md:flex">
        <li v-for="link in links" :key="link.id">
          <button
            type="button"
            class="rounded-[4px] px-3 py-1.5 text-[13px] font-medium text-muted-foreground transition hover:text-foreground"
            @click="goTo(link.id)"
          >
            {{ link.label }}
          </button>
        </li>
      </ul>

      <!-- Actions -->
      <div class="flex items-center gap-1.5">
        <!-- 1. Sélecteur de Langue (Desktop) -->
        <div id="landing-lang-dropdown" class="relative hidden sm:block">
          <button
            type="button"
            class="flex h-8 items-center gap-1.5 rounded-[4px] border border-border bg-card/60 px-2 text-[12px] font-medium text-foreground transition hover:bg-muted focus:outline-none"
            :aria-expanded="isLangMenuOpen"
            title="Changer la langue d'affichage"
            @click.stop="toggleLangMenu"
          >
            <span class="text-xs">{{ languageStore.currentMeta.flag }}</span>
            <span
              class="uppercase tracking-wider text-[11px] font-semibold text-muted-foreground"
            >
              {{ languageStore.currentLanguage }}
            </span>
            <ChevronDown
              class="h-3 w-3 text-muted-foreground transition-transform duration-150"
              :class="{ 'rotate-180': isLangMenuOpen }"
            />
          </button>

          <!-- Dropdown langues -->
          <div
            v-if="isLangMenuOpen"
            class="absolute right-0 top-full mt-1.5 w-44 rounded-[4px] border border-border bg-card p-1 shadow-lg text-xs z-50 animate-in fade-in zoom-in-95 duration-100"
          >
            <button
              v-for="lang in languageStore.availableLanguages"
              :key="lang.code"
              type="button"
              @click="selectLanguage(lang.code)"
              class="flex w-full items-center justify-between rounded-[3px] px-2.5 py-1.5 text-left transition hover:bg-muted"
              :class="{
                'font-semibold text-brand bg-brand-50/50 dark:bg-brand-950/40':
                  languageStore.currentLanguage === lang.code,
              }"
            >
              <div class="flex items-center gap-2">
                <span class="text-sm">{{ lang.flag }}</span>
                <span class="text-[12px] text-foreground">{{ lang.name }}</span>
              </div>
              <Check
                v-if="languageStore.currentLanguage === lang.code"
                class="h-3.5 w-3.5 text-brand shrink-0"
              />
            </button>
          </div>
        </div>

        <!-- 2. Bascule Mode Clair / Sombre -->
        <button
          type="button"
          class="hidden h-8 w-8 items-center justify-center rounded-[4px] text-muted-foreground transition hover:bg-muted hover:text-foreground sm:flex"
          :aria-label="authStore.isDark ? 'Mode clair' : 'Mode sombre'"
          @click="authStore.toggleDarkMode()"
        >
          <Sun v-if="authStore.isDark" class="h-4 w-4" />
          <Moon v-else class="h-4 w-4" />
        </button>

        <!-- 3. Bouton Connexion -->
        <button
          type="button"
          class="rounded-[4px] bg-brand-600 px-3.5 py-1.5 text-[13px] font-medium text-white transition hover:bg-brand-700 active:scale-[0.98]"
          @click="goAuth"
        >
          {{ languageStore.t("Connexion") }}
        </button>

        <!-- 4. Toggle Menu Mobile -->
        <button
          type="button"
          class="flex h-8 w-8 items-center justify-center rounded-[4px] text-muted-foreground hover:bg-muted md:hidden"
          :aria-expanded="isMenuOpen"
          aria-label="Menu"
          @click="isMenuOpen = !isMenuOpen"
        >
          <X v-if="isMenuOpen" class="h-4 w-4" />
          <Menu v-else class="h-4 w-4" />
        </button>
      </div>
    </nav>

    <!-- Menu mobile -->
    <div
      v-if="isMenuOpen"
      class="border-t border-border bg-background px-5 py-3 md:hidden space-y-2 animate-in fade-in duration-100"
    >
      <button
        v-for="link in links"
        :key="link.id"
        type="button"
        class="block w-full rounded-[4px] px-2 py-2 text-left text-sm font-medium text-foreground hover:bg-muted"
        @click="goTo(link.id)"
      >
        {{ link.label }}
      </button>

      <!-- Sélecteur de langue en mode mobile -->
      <div class="border-t border-border/70 pt-2.5">
        <p
          class="px-2 text-[10px] font-semibold text-muted-foreground uppercase tracking-wider mb-2"
        >
          {{ languageStore.t("Langue") }}
        </p>
        <div class="grid grid-cols-2 gap-1.5 px-1">
          <button
            v-for="lang in languageStore.availableLanguages"
            :key="lang.code"
            type="button"
            @click="selectLanguage(lang.code)"
            class="flex items-center justify-between rounded-[4px] border border-border/70 bg-card px-2.5 py-1.5 text-[12px]"
            :class="{
              'border-brand text-brand font-semibold bg-brand-50/50 dark:bg-brand-950/40':
                languageStore.currentLanguage === lang.code,
            }"
          >
            <span class="flex items-center gap-1.5">
              <span>{{ lang.flag }}</span>
              <span>{{ lang.name }}</span>
            </span>
            <Check
              v-if="languageStore.currentLanguage === lang.code"
              class="h-3 w-3 text-brand"
            />
          </button>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { useAuthStore } from "../../stores/authStore";
import { useGuichetStore } from "../../stores/guichetStore";
import { useLanguageStore } from "../../stores/languageStore";
import type { SupportedLanguage } from "../../services/deepLTranslationService";
import appLogo from "@root/app-icon.png";
import {
  SunMoon,
  Sun,
  Moon,
  Monitor,
  Sliders,
  Palette,
  Sparkles,
  ChevronRight,
  Check,
  HelpCircle,
  Keyboard,
  Info,
  Minus,
  Square,
  LucideCopy,
  X,
  Maximize2,
  Terminal,
  Megaphone,
} from "lucide-vue-next";
import {
  CUSTOM_COLOR_THEMES,
  CUSTOM_BACKGROUND_STYLES,
  type BackgroundStyleId,
} from "../../config/themes";
import NetworkNotificationsModal from "./NetworkNotificationsModal.vue";

const authStore = useAuthStore();
const guichetStore = useGuichetStore();
const languageStore = useLanguageStore();

// État des menus ouverts
const activeMenu = ref<"affichage" | "aide" | "langue" | null>(null);
const activeSubmenu = ref<"theme" | "personnaliser" | null>(null);
const activePersonnaliserSub = ref<"color-theme" | "background" | null>(null);

// Modales annexes
const isAboutModalOpen = ref(false);
const isShortcutsModalOpen = ref(false);
const isNetworkAdsModalOpen = ref(false);

// Thème actuel, Thème de couleur & Background
const currentTheme = ref<"light" | "dark" | "system">("system");
const currentColorTheme = ref<string>("default");
const currentBgStyle = ref<BackgroundStyleId>("default");

const isTauri =
  typeof window !== "undefined" &&
  ("__TAURI_INTERNALS__" in window || "__TAURI__" in window);

const isMenuOpen = computed(() => activeMenu.value !== null);

const toggleMenu = (menu: "affichage" | "aide" | "langue") => {
  if (activeMenu.value === menu) {
    activeMenu.value = null;
    activeSubmenu.value = null;
    activePersonnaliserSub.value = null;
  } else {
    activeMenu.value = menu;
    activeSubmenu.value = null;
    activePersonnaliserSub.value = null;
  }
};

/**
 * Lorsque un menu est ouvert, survoler un autre menu bascule directement dessus
 */
const handleMenuHover = (menu: "affichage" | "aide" | "langue") => {
  if (activeMenu.value !== null && activeMenu.value !== menu) {
    activeMenu.value = menu;
    activeSubmenu.value = null;
    activePersonnaliserSub.value = null;
  }
};

const handleSubmenuHover = (sub: "theme" | "personnaliser") => {
  activeSubmenu.value = sub;
  if (sub !== "personnaliser") {
    activePersonnaliserSub.value = null;
  }
};

const closeAllMenus = () => {
  activeMenu.value = null;
  activeSubmenu.value = null;
  activePersonnaliserSub.value = null;
};

/**
 * Gestion du thème (Clair / Sombre / Système)
 */
const setTheme = (mode: "light" | "dark" | "system") => {
  currentTheme.value = mode;
  if (mode === "system") {
    localStorage.removeItem("autoschool_theme");
    const prefersDark =
      window.matchMedia &&
      window.matchMedia("(prefers-color-scheme: dark)").matches;
    authStore.applyTheme(prefersDark);
    authStore.isDark = prefersDark;
  } else {
    const isDark = mode === "dark";
    localStorage.setItem("autoschool_theme", mode);
    authStore.applyTheme(isDark);
    authStore.isDark = isDark;
  }
  closeAllMenus();
  guichetStore.showToast(
    `Thème d'affichage : ${mode === "dark" ? "Sombre" : mode === "light" ? "Clair" : "Automatique (Système)"}`,
    "info",
  );
};

const activeColorThemeMeta = computed(
  () =>
    CUSTOM_COLOR_THEMES.find((t) => t.id === currentColorTheme.value) ||
    CUSTOM_COLOR_THEMES[0],
);

const setColorTheme = (themeId: string, label: string) => {
  currentColorTheme.value = themeId;
  if (themeId === "default") {
    document.documentElement.removeAttribute("data-color-theme");
    localStorage.removeItem("autoschool_color_theme");
  } else {
    document.documentElement.setAttribute("data-color-theme", themeId);
    localStorage.setItem("autoschool_color_theme", themeId);
  }
  closeAllMenus();
  guichetStore.showToast(`Thème de couleur : ${label}`, "info");
};

const setBackgroundStyle = (bgId: BackgroundStyleId, label: string) => {
  currentBgStyle.value = bgId;
  if (bgId === "default") {
    document.documentElement.removeAttribute("data-bg-style");
    localStorage.removeItem("autoschool_bg_style");
  } else {
    document.documentElement.setAttribute("data-bg-style", bgId);
    localStorage.setItem("autoschool_bg_style", bgId);
  }
  closeAllMenus();
  guichetStore.showToast(`Background : ${label}`, "info");
};

/**
 * Gestion de la langue (DeepL API réactive via languageStore)
 */
const setLanguage = (lang: SupportedLanguage) => {
  languageStore.setLanguage(lang);
  closeAllMenus();
};

/**
 * Personnalisation : réduction de la barre latérale
 */
const toggleSidebar = () => {
  authStore.toggleSidebar();
  closeAllMenus();
};

const isMaximized = ref(false);

const updateMaximizedState = async () => {
  if (!isTauri) return;
  try {
    const { getCurrentWindow } = await import("@tauri-apps/api/window");
    isMaximized.value = await getCurrentWindow().isMaximized();
  } catch {
    // ignore
  }
};

const handleGlobalKeydown = (e: KeyboardEvent) => {
  if (e.key === "Escape" && activeMenu.value !== null) {
    closeAllMenus();
  }
};

onMounted(async () => {
  const savedTheme = localStorage.getItem("autoschool_theme");
  if (savedTheme === "light" || savedTheme === "dark") {
    currentTheme.value = savedTheme;
  } else {
    currentTheme.value = "system";
  }

  const savedColorTheme = localStorage.getItem("autoschool_color_theme");
  if (savedColorTheme) {
    currentColorTheme.value = savedColorTheme;
  }

  const savedBgStyle = localStorage.getItem("autoschool_bg_style") as BackgroundStyleId | null;
  if (savedBgStyle) {
    currentBgStyle.value = savedBgStyle;
  }

  languageStore.init();

  window.addEventListener("keydown", handleGlobalKeydown);

  if (isTauri) {
    await updateMaximizedState();
    try {
      const { getCurrentWindow } = await import("@tauri-apps/api/window");
      getCurrentWindow().onResized(async () => {
        await updateMaximizedState();
      });
    } catch {
      // ignore
    }
  }
});

onUnmounted(() => {
  window.removeEventListener("keydown", handleGlobalKeydown);
});

/**
 * Déplacement natif de la fenêtre Tauri (Drag)
 */
const handleMouseDown = async (e: MouseEvent) => {
  // Uniquement le clic gauche déclenche le drag
  if (e.button !== 0) return;
  const target = e.target as HTMLElement | null;
  if (!target) return;
  // Ne pas déplacer si on clique sur un bouton, menu, input ou modale
  if (
    target.closest("button") ||
    target.closest("nav") ||
    target.closest('[role="menu"]') ||
    target.closest('[role="menuitem"]') ||
    target.closest(".no-drag")
  ) {
    return;
  }

  if (isTauri) {
    try {
      const { invoke } = await import("@tauri-apps/api/core");
      await invoke("window_start_drag");
    } catch {
      try {
        const { getCurrentWindow } = await import("@tauri-apps/api/window");
        await getCurrentWindow().startDragging();
      } catch (err) {
        console.warn("Tauri startDragging error:", err);
      }
    }
  }
};

/**
 * Double-clic sur l'en-tête pour basculer plein écran / restauration
 */
const handleHeaderDblClick = async (e: MouseEvent) => {
  const target = e.target as HTMLElement | null;
  if (!target) return;
  if (
    target.closest("button") ||
    target.closest("nav") ||
    target.closest('[role="menu"]') ||
    target.closest('[role="menuitem"]') ||
    target.closest(".no-drag")
  ) {
    return;
  }
  await toggleMaximizeWindow();
};

/**
 * Contrôles de la fenêtre Tauri (Minimize, Maximize, Close)
 */
const minimizeWindow = async () => {
  if (!isTauri) return;
  try {
    const { invoke } = await import("@tauri-apps/api/core");
    await invoke("window_minimize");
  } catch {
    try {
      const { getCurrentWindow } = await import("@tauri-apps/api/window");
      await getCurrentWindow().minimize();
    } catch (err) {
      console.warn("Tauri window minimize error:", err);
    }
  }
};

const toggleMaximizeWindow = async () => {
  if (!isTauri) return;
  try {
    const { invoke } = await import("@tauri-apps/api/core");
    await invoke("window_toggle_maximize");
    await updateMaximizedState();
  } catch {
    try {
      const { getCurrentWindow } = await import("@tauri-apps/api/window");
      await getCurrentWindow().toggleMaximize();
      await updateMaximizedState();
    } catch (err) {
      console.warn("Tauri window toggleMaximize error:", err);
    }
  }
};

const closeWindow = async () => {
  if (!isTauri) return;
  try {
    const { invoke } = await import("@tauri-apps/api/core");
    await invoke("window_close");
  } catch {
    try {
      const { getCurrentWindow } = await import("@tauri-apps/api/window");
      await getCurrentWindow().close();
    } catch (err) {
      console.warn("Tauri window close error:", err);
    }
  }
};
</script>

<template>
  <header
    class="relative h-9 w-full bg-background border-b border-border flex items-center justify-between select-none z-[40] shrink-0 text-foreground cursor-default"
  >
    <!-- Overlay invisible pour fermer le menu au clic extérieur -->
    <div
      v-if="isMenuOpen"
      class="fixed inset-0 z-40 bg-transparent no-drag"
      @click="closeAllMenus"
    ></div>

    <!-- Section Gauche : Logo Visible + Menus (Affichage, Aide, Langue) -->
    <div class="flex items-center h-full pl-2.5 z-50 no-drag" @mousedown.stop>
      <!-- Logo bien visible et contrasté -->
      <div
        class="flex items-center gap-2 mr-3 px-1.5 py-0.5 rounded cursor-default"
        title="ElevenDrive"
      >
        <img
          :src="appLogo"
          alt="ElevenDrive"
          class="h-5 w-5 object-contain rounded-xs shrink-0 select-none pointer-events-none"
        />
      </div>

      <!-- Menu Bar Desktop -->
      <nav class="flex items-center h-full text-xs" role="menubar">
        <!-- 1. Menu Affichage -->
        <div class="relative h-full flex items-center">
          <button
            type="button"
            role="menuitem"
            @click.stop="toggleMenu('affichage')"
            @mouseenter="handleMenuHover('affichage')"
            class="px-2.5 py-1 rounded text-[11px] transition cursor-default flex items-center gap-1"
            :class="
              activeMenu === 'affichage'
                ? 'bg-muted text-foreground'
                : 'text-foreground/80 hover:text-foreground hover:bg-muted/70'
            "
          >
            <span>{{ languageStore.t("Affichage") }}</span>
          </button>

          <!-- Sheet / Menu déroulant Affichage -->
          <div
            v-if="activeMenu === 'affichage'"
            @click.stop
            class="absolute left-0 top-full mt-0.5 w-48 rounded-[3px] border border-border bg-background p-1 text-xs z-50 animate-in fade-in zoom-in-95 duration-100 shadow-xl"
          >
            <!-- Sous-menu Thème -->
            <div
              class="relative flex items-center justify-between px-2.5 py-1.5 rounded-[3px] hover:bg-muted cursor-pointer transition text-foreground"
              @mouseenter="handleSubmenuHover('theme')"
              @click.stop="
                activeSubmenu = activeSubmenu === 'theme' ? null : 'theme'
              "
            >
              <div class="flex text-xs items-center gap-2">
                <SunMoon class="h-3.5 w-3.5 text-brand" />
                <span>{{ languageStore.t("Thème") }}</span>
              </div>
              <ChevronRight class="h-3.5 w-3.5 text-muted-foreground" />

              <!-- Petit Sheet de côté pour les modes de thème -->
              <div
                v-if="activeSubmenu === 'theme'"
                class="absolute left-full top-0 ml-1 w-36 rounded-[3px] border border-border bg-background p-1 text-xs z-50 animate-in fade-in slide-in-from-left-1 duration-100"
              >
                <button
                  type="button"
                  @click="setTheme('light')"
                  class="w-full text-xs flex items-center justify-between px-2 py-1.5 rounded-[3px] hover:bg-muted transition text-left"
                  :class="{
                    'font-semibold text-brand bg-brand-light/40':
                      currentTheme === 'light',
                  }"
                >
                  <div class="flex text-xs items-center gap-2">
                    <Sun class="h-3.5 w-3.5 text-amber-500" />
                    <span>{{ languageStore.t("Soleil") }}</span>
                  </div>
                  <Check
                    v-if="currentTheme === 'light'"
                    class="h-3.5 w-3.5 text-brand"
                  />
                </button>

                <button
                  type="button"
                  @click="setTheme('dark')"
                  class="w-full flex text-xs items-center justify-between px-2 py-1.5 rounded-[3px] hover:bg-muted transition text-left"
                  :class="{
                    'font-semibold text-brand bg-brand-light/40':
                      currentTheme === 'dark',
                  }"
                >
                  <div class="flex text-xs items-center gap-2">
                    <Moon class="h-3.5 w-3.5 text-indigo-400" />
                    <span>{{ languageStore.t("Nuit") }}</span>
                  </div>
                  <Check
                    v-if="currentTheme === 'dark'"
                    class="h-3.5 w-3.5 text-brand"
                  />
                </button>

                <button
                  type="button"
                  @click="setTheme('system')"
                  class="w-full text-xs flex items-center justify-between px-2 py-1.5 rounded-[3px] hover:bg-muted transition text-left"
                  :class="{
                    'font-semibold text-brand bg-brand-light/40':
                      currentTheme === 'system',
                  }"
                >
                  <div class="flex text-xs items-center gap-2">
                    <Monitor class="h-3.5 w-3.5 text-sky-500" />
                    <span>{{ languageStore.t("Jour") }}</span>
                  </div>
                  <Check
                    v-if="currentTheme === 'system'"
                    class="h-3.5 w-3.5 text-brand"
                  />
                </button>
              </div>
            </div>

            <!-- Sous-menu Personnaliser -->
            <div
              class="relative flex items-center justify-between px-2.5 py-1.5 rounded-[3px] hover:bg-muted cursor-pointer transition text-foreground"
              @mouseenter="handleSubmenuHover('personnaliser')"
              @click.stop="
                activeSubmenu =
                  activeSubmenu === 'personnaliser' ? null : 'personnaliser'
              "
            >
              <div class="flex text-xs items-center gap-2">
                <Sliders class="h-3.5 w-3.5 text-muted-foreground" />
                <span>{{ languageStore.t("Personnaliser") }}</span>
              </div>
              <ChevronRight class="h-3.5 w-3.5 text-muted-foreground" />

              <!-- Sous-menu Personnaliser (Thème de couleur & Background) -->
              <div
                v-if="activeSubmenu === 'personnaliser'"
                class="absolute left-full top-0 ml-1 w-48 rounded-[3px] border border-border bg-background p-1 text-xs z-50 animate-in fade-in slide-in-from-left-1 duration-100 shadow-xl"
              >
                <!-- Option 1 : Thème de couleur -->
                <div
                  class="relative flex items-center justify-between px-2.5 py-1.5 rounded-[3px] hover:bg-muted cursor-pointer transition text-foreground"
                  @mouseenter="activePersonnaliserSub = 'color-theme'"
                  @click.stop="
                    activePersonnaliserSub =
                      activePersonnaliserSub === 'color-theme'
                        ? null
                        : 'color-theme'
                  "
                >
                  <div class="flex text-xs items-center gap-2 min-w-0">
                    <Palette class="h-3.5 w-3.5 text-primary shrink-0" />
                    <span class="truncate">{{
                      languageStore.t("Thème de couleur")
                    }}</span>
                  </div>
                  <div class="flex items-center gap-1.5 shrink-0">
                    <span
                      class="h-2.5 w-2.5 rounded-full border border-border"
                      :style="{ backgroundColor: activeColorThemeMeta.primary }"
                    ></span>
                    <ChevronRight class="h-3.5 w-3.5 text-muted-foreground" />
                  </div>

                  <!-- Liste des thèmes de couleurs personnalisés -->
                  <div
                    v-if="activePersonnaliserSub === 'color-theme'"
                    class="absolute left-full top-0 ml-1 w-56 max-h-80 overflow-y-auto rounded-[3px] border border-border bg-background p-1 text-xs z-50 animate-in fade-in slide-in-from-left-1 duration-100 shadow-xl"
                  >
                    <button
                      v-for="theme in CUSTOM_COLOR_THEMES"
                      :key="theme.id"
                      type="button"
                      @click.stop="setColorTheme(theme.id, theme.label)"
                      class="w-full text-xs flex items-center justify-between gap-2 px-2 py-1.5 rounded-[3px] hover:bg-muted transition text-left"
                      :class="{
                        'font-semibold text-primary bg-muted/60':
                          currentColorTheme === theme.id,
                      }"
                    >
                      <div class="flex items-center gap-2 min-w-0">
                        <div class="flex items-center -space-x-1 shrink-0">
                          <span
                            class="h-3 w-3 rounded-full border border-background"
                            :style="{ backgroundColor: theme.primary }"
                          ></span>
                          <span
                            class="h-3 w-3 rounded-full border border-background"
                            :style="{ backgroundColor: theme.secondary }"
                          ></span>
                          <span
                            class="h-3 w-3 rounded-full border border-background"
                            :style="{ backgroundColor: theme.accent }"
                          ></span>
                        </div>
                        <span class="truncate">{{ theme.label }}</span>
                      </div>
                      <Check
                        v-if="currentColorTheme === theme.id"
                        class="h-3.5 w-3.5 text-primary shrink-0"
                      />
                    </button>
                  </div>
                </div>

                <!-- Option 2 : Background -->
                <div
                  class="relative flex items-center justify-between px-2.5 py-1.5 rounded-[3px] hover:bg-muted cursor-pointer transition text-foreground"
                  @mouseenter="activePersonnaliserSub = 'background'"
                  @click.stop="
                    activePersonnaliserSub =
                      activePersonnaliserSub === 'background'
                        ? null
                        : 'background'
                  "
                >
                  <div class="flex text-xs items-center gap-2">
                    <Sparkles class="h-3.5 w-3.5 text-amber-500" />
                    <span>{{ languageStore.t("Background") }}</span>
                  </div>
                  <ChevronRight class="h-3.5 w-3.5 text-muted-foreground" />

                  <!-- Liste des backgrounds -->
                  <div
                    v-if="activePersonnaliserSub === 'background'"
                    class="absolute left-full top-0 ml-1 w-48 rounded-[3px] border border-border bg-background p-1 text-xs z-50 animate-in fade-in slide-in-from-left-1 duration-100 shadow-xl"
                  >
                    <button
                      v-for="bg in CUSTOM_BACKGROUND_STYLES"
                      :key="bg.id"
                      type="button"
                      @click.stop="setBackgroundStyle(bg.id, bg.label)"
                      class="w-full text-xs flex items-center justify-between gap-2 px-2 py-1.5 rounded-[3px] hover:bg-muted transition text-left"
                      :class="{
                        'font-semibold text-primary bg-muted/60':
                          currentBgStyle === bg.id,
                      }"
                    >
                      <div class="min-w-0">
                        <div class="truncate">{{ bg.label }}</div>
                        <div class="text-[10px] text-muted-foreground truncate">
                          {{ bg.description }}
                        </div>
                      </div>
                      <Check
                        v-if="currentBgStyle === bg.id"
                        class="h-3.5 w-3.5 text-primary shrink-0"
                      />
                    </button>
                  </div>
                </div>

                <div class="my-1 border-t border-border"></div>

                <!-- Option Barre latérale -->
                <button
                  type="button"
                  @click="toggleSidebar"
                  class="w-full flex items-center justify-between px-2.5 py-1.5 rounded-[3px] hover:bg-muted transition text-left text-foreground"
                >
                  <span class="text-xs">{{
                    languageStore.t("Barre latérale")
                  }}</span>
                  <span class="text-[10px] text-muted-foreground font-mono"
                    >Sidebar</span
                  >
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- 2. Menu Aide -->
        <div class="relative h-full flex items-center">
          <button
            type="button"
            role="menuitem"
            @click.stop="toggleMenu('aide')"
            @mouseenter="handleMenuHover('aide')"
            class="px-2.5 py-1 rounded text-[11px] transition cursor-default flex items-center gap-1"
            :class="
              activeMenu === 'aide'
                ? 'bg-muted text-foreground'
                : 'text-foreground/80 hover:text-foreground hover:bg-muted/70'
            "
          >
            <span>{{ languageStore.t("Aide") }}</span>
          </button>

          <!-- Sheet / Menu déroulant Aide -->
          <div
            v-if="activeMenu === 'aide'"
            @click.stop
            class="absolute left-0 top-full mt-0.5 w-56 rounded-[3px] border border-border bg-background p-1 text-xs z-50 animate-in fade-in zoom-in-95 duration-100 shadow-xl"
          >
            <button
              type="button"
              @click="
                authStore.isHelpModalOpen = true;
                closeAllMenus();
              "
              class="w-full flex text-xs items-center gap-2 px-2.5 py-1.5 rounded-[3px] hover:bg-muted transition text-left text-foreground"
            >
              <HelpCircle class="h-3.5 w-3.5 text-brand" />
              <span>{{ languageStore.t("Centre d'aide & Guichet") }}</span>
            </button>

            <button
              type="button"
              @click="
                isShortcutsModalOpen = true;
                closeAllMenus();
              "
              class="w-full flex items-center justify-between px-2.5 py-1.5 rounded-[3px] hover:bg-muted transition text-left text-foreground"
            >
              <div class="flex text-xs items-center gap-2">
                <Keyboard class="h-3.5 w-3.5 text-muted-foreground" />
                <span>{{ languageStore.t("Raccourcis clavier") }}</span>
              </div>
              <span class="text-[10px] text-muted-foreground font-mono"
                >F1-F9</span
              >
            </button>

            <div class="my-1 border-t border-border"></div>

            <button
              type="button"
              @click="
                isNetworkAdsModalOpen = true;
                closeAllMenus();
              "
              class="w-full text-xs flex items-center justify-between px-2.5 py-1.5 rounded-[3px] hover:bg-muted transition text-left text-foreground"
            >
              <div class="flex items-center gap-2">
                <Megaphone class="h-3.5 w-3.5 text-brand" />
                <span>{{ languageStore.t("Annonces Réseau & Avis") }}</span>
              </div>
              <span
                class="text-[9px] uppercase font-mono px-1 py-0.2 rounded bg-amber-500/10 text-amber-600 font-semibold"
              >
                Ads ★
              </span>
            </button>

            <button
              type="button"
              @click="
                isAboutModalOpen = true;
                closeAllMenus();
              "
              class="w-full text-xs flex items-center gap-2 px-2.5 py-1.5 rounded-[3px] hover:bg-muted transition text-left text-foreground"
            >
              <Info class="h-3.5 w-3.5 text-muted-foreground" />
              <span>{{ languageStore.t("À propos d'AutoSchool") }}</span>
            </button>

            <button
              type="button"
              @click="
                authStore.webPortalView = 'super_admin_login';
                closeAllMenus();
              "
              class="w-full text-xs flex items-center gap-2 px-2.5 py-1.5 rounded-[3px] hover:bg-emerald-500/10 hover:text-emerald-500 transition text-left text-foreground font-mono"
            >
              <Terminal class="h-3.5 w-3.5 text-emerald-500" />
              <span>Console Super-Admin</span>
            </button>
          </div>
        </div>

        <!-- 3. Menu Langue -->
        <div class="relative h-full flex items-center">
          <button
            type="button"
            role="menuitem"
            @click.stop="toggleMenu('langue')"
            @mouseenter="handleMenuHover('langue')"
            class="px-2.5 py-1 rounded text-[11px] transition cursor-default flex items-center gap-1.5"
            :class="
              activeMenu === 'langue'
                ? 'bg-muted text-foreground '
                : 'text-foreground/80 hover:text-foreground hover:bg-muted/70'
            "
          >
            <span class="text-xs">{{ languageStore.currentMeta.flag }}</span>
            <span>{{ languageStore.t("Langue") }}</span>
          </button>

          <!-- Sheet / Menu déroulant Langue (Direct selection) -->
          <div
            v-if="activeMenu === 'langue'"
            @click.stop
            class="absolute left-0 top-full mt-0.5 w-48 rounded-[3px] border border-border bg-background p-1 text-xs z-50 animate-in fade-in zoom-in-95 duration-100 shadow-xl"
          >
            <button
              v-for="lang in languageStore.availableLanguages"
              :key="lang.code"
              type="button"
              @click="setLanguage(lang.code)"
              class="w-full flex items-center justify-between px-2.5 py-2 rounded-[3px] hover:bg-muted transition text-left text-foreground"
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
      </nav>
    </div>

    <!-- Section Centrale : Zone déplaçable de la fenêtre Tauri (Drag Region) -->
    <div
      data-tauri-drag-region
      class="flex-1 h-full flex items-center justify-center cursor-default pointer-events-auto"
      :title="
        languageStore.t(
          'Glisser pour déplacer · Double-clic pour agrandir / restaurer',
        )
      "
    >
      <span
        class="text-[11px] font-medium text-muted-foreground/60 tracking-wider pointer-events-none select-none"
      >
        {{ authStore.schoolName }}
      </span>
    </div>

    <!-- Section Droite : Boutons de contrôle de la fenêtre (Réduire, Agrandir, Fermer) -->
    <div class="flex items-center h-full z-50 no-drag" @mousedown.stop>
      <!-- Réduire -->
      <button
        type="button"
        @click="minimizeWindow"
        class="h-full px-3.5 inline-flex items-center justify-center text-muted-foreground hover:text-foreground hover:bg-muted transition"
        :title="languageStore.t('Réduire')"
        :aria-label="languageStore.t('Réduire')"
      >
        <Minus class="w-3.5 h-3.5" />
      </button>

      <!-- Agrandir / Restaurer -->
      <button
        type="button"
        @click="toggleMaximizeWindow"
        class="h-full px-3.5 inline-flex items-center justify-center text-muted-foreground hover:text-foreground hover:bg-muted transition"
        :title="
          isMaximized
            ? languageStore.t('Restaurer')
            : languageStore.t('Agrandir')
        "
        :aria-label="
          isMaximized
            ? languageStore.t('Restaurer')
            : languageStore.t('Agrandir')
        "
      >
        <svg
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          v-if="isMaximized"
          class="w-3.5 h-3.5 rotate-180"
        >
          {/*carre du fond */}
          <rect x="4" y="4" width="14" height="14" rx="1" />
          {/*carre du devant avec fond */}
          <rect x="8" y="8" width="14" height="14" rx="1" />
        </svg>
        <!-- <LucideCopy v-if="isMaximized" class="w-3.5 h-3.5 rotate-180" /> -->
        <svg
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          v-else
          class="w-3.5 h-3.5"
        >
          {/*carre du fond */}
          <rect x="4" y="4" width="14" height="14" rx="1" />
        </svg>
        <!-- <Square v-else class="w-3.5 h-3.5" /> -->
      </button>

      <!-- Fermer -->
      <button
        type="button"
        @click="closeWindow"
        class="h-full px-3.5 inline-flex items-center justify-center text-muted-foreground hover:text-white hover:bg-rose-600 transition"
        :title="languageStore.t('Fermer')"
        :aria-label="languageStore.t('Fermer')"
      >
        <X class="w-3.5 h-3.5" />
      </button>
    </div>
  </header>

  <!-- Modale Raccourcis Clavier -->
  <div
    v-if="isShortcutsModalOpen"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 backdrop-blur-xs animate-in fade-in duration-150"
    role="dialog"
  >
    <div
      class="w-full max-w-md rounded-[4px] border border-border bg-card p-5 text-foreground text-xs space-y-4"
    >
      <div
        class="flex items-center justify-between border-b border-border pb-3"
      >
        <div class="flex items-center gap-2">
          <Keyboard class="h-4 w-4 text-brand" />
          <h3 class="font-bold text-sm">
            {{ languageStore.t("Raccourcis Clavier Guichet") }}
          </h3>
        </div>
        <button
          type="button"
          @click="isShortcutsModalOpen = false"
          class="rounded p-1 text-muted-foreground hover:bg-muted hover:text-foreground"
        >
          <X class="h-4 w-4" />
        </button>
      </div>

      <div class="space-y-2 font-mono">
        <div class="flex items-center justify-between p-2 rounded bg-muted/40">
          <span class="text-muted-foreground font-sans">{{
            languageStore.t("Nouveau Candidat")
          }}</span>
          <span
            class="px-2 py-0.5 rounded bg-card border border-border text-brand font-bold"
            >F1</span
          >
        </div>
        <div class="flex items-center justify-between p-2 rounded bg-muted/40">
          <span class="text-muted-foreground font-sans">{{
            languageStore.t("Encaisser Paiement")
          }}</span>
          <span
            class="px-2 py-0.5 rounded bg-card border border-border text-brand font-bold"
            >F2</span
          >
        </div>
        <div class="flex items-center justify-between p-2 rounded bg-muted/40">
          <span class="text-muted-foreground font-sans">{{
            languageStore.t("Rechercher Candidat")
          }}</span>
          <span
            class="px-2 py-0.5 rounded bg-card border border-border text-brand font-bold"
            >F3</span
          >
        </div>
        <div class="flex items-center justify-between p-2 rounded bg-muted/40">
          <span class="text-muted-foreground font-sans">{{
            languageStore.t("Clôture de Caisse (Rapport Z)")
          }}</span>
          <span
            class="px-2 py-0.5 rounded bg-card border border-border text-brand font-bold"
            >F4</span
          >
        </div>
        <div class="flex items-center justify-between p-2 rounded bg-muted/40">
          <span class="text-muted-foreground font-sans">{{
            languageStore.t("Vue Code de la route")
          }}</span>
          <span
            class="px-2 py-0.5 rounded bg-card border border-border text-brand font-bold"
            >F5</span
          >
        </div>
        <div class="flex items-center justify-between p-2 rounded bg-muted/40">
          <span class="text-muted-foreground font-sans">{{
            languageStore.t("Examens Permis")
          }}</span>
          <span
            class="px-2 py-0.5 rounded bg-card border border-border text-brand font-bold"
            >F9</span
          >
        </div>
      </div>

      <button
        type="button"
        @click="isShortcutsModalOpen = false"
        class="w-full h-8 rounded-[4px] bg-brand text-white font-semibold hover:bg-brand-dark transition cursor-pointer"
      >
        {{ languageStore.t("Fermer") }}
      </button>
    </div>
  </div>

  <!-- Modale À Propos -->
  <div
    v-if="isAboutModalOpen"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 backdrop-blur-xs animate-in fade-in duration-150"
    role="dialog"
  >
    <div
      class="w-full max-w-sm rounded-[4px] border border-border bg-card p-5 text-foreground text-xs space-y-4 text-center"
    >
      <div class="flex justify-end">
        <button
          type="button"
          @click="isAboutModalOpen = false"
          class="rounded p-1 text-muted-foreground hover:bg-muted hover:text-foreground"
        >
          <X class="h-4 w-4" />
        </button>
      </div>

      <div class="flex flex-col items-center gap-2">
        <div
          class="flex h-14 w-14 items-center justify-center rounded-2xl bg-card border border-border shadow-md p-2"
        >
          <img
            :src="appLogo"
            alt="AutoSchool Logo"
            class="h-full w-full object-contain select-none"
          />
        </div>
        <h3 class="font-bold text-base">AutoSchool Desktop</h3>
        <span class="text-xs text-muted-foreground font-mono"
          >Version 0.1.0 · ElevenDriver</span
        >
      </div>

      <p class="text-muted-foreground text-[11px] leading-relaxed">
        Système de gestion intégrée pour Auto-Écoles : encaissement guichet,
        livrets mobiles moniteurs, et portail candidats.
      </p>

      <div
        class="rounded-[4px] border border-border bg-muted/30 p-2.5 text-[10px] text-muted-foreground space-y-1"
      >
        <div>
          Moteur : <strong>Tauri v2 + Vue 3 + SQLite / Neon PostgreSQL</strong>
        </div>
        <div>Plateforme : <strong>Desktop Client (Windows)</strong></div>
      </div>

      <button
        type="button"
        @click="isAboutModalOpen = false"
        class="w-full h-8 rounded-[4px] bg-brand text-white font-semibold hover:bg-brand-dark transition cursor-pointer"
      >
        {{ languageStore.t("Fermer") }}
      </button>
    </div>
  </div>

  <!-- Modale Panneaux Ads & Retours Avis -->
  <NetworkNotificationsModal
    :is-open="isNetworkAdsModalOpen"
    @close="isNetworkAdsModalOpen = false"
  />
</template>

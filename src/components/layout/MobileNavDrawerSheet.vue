<script setup lang="ts">
import { computed } from "vue";
import { useAuthStore } from "../../stores/authStore";
import { useLanguageStore } from "../../stores/languageStore";
import { getRoleNavGroups, type NavGroup } from "./sidebarNavConfig";
import {
  X,
  Lock,
  User,
  QrCode,
  HelpCircle,
  Sun,
  Moon,
  LogOut,
  Globe,
  Sparkles,
} from "lucide-vue-next";

defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "select-section", sectionId: string): void;
}>();

const authStore = useAuthStore();
const languageStore = useLanguageStore();
const t = (txt: string) => languageStore.t(txt);

const navGroups = computed<NavGroup[]>(() => {
  return getRoleNavGroups(authStore.currentRole).map((group) => ({
    ...group,
    title: t(group.title),
    items: group.items.map((item) => ({
      ...item,
      label: t(item.label),
    })),
  }));
});

const handlePickSection = (sectionId: string) => {
  if (!authStore.isModuleEnabled(sectionId)) {
    return;
  }
  emit("select-section", sectionId);
  emit("close");
};

const openProfile = () => {
  emit("close");
  window.dispatchEvent(new CustomEvent("elevendrive:open-profile"));
};

const openQr = () => {
  emit("close");
  window.dispatchEvent(new CustomEvent("elevendrive:open-qr"));
};

const openHelp = () => {
  emit("close");
  window.dispatchEvent(new CustomEvent("elevendrive:open-help"));
};

const toggleLang = () => {
  const next = languageStore.currentLanguage === "fr" ? "en" : "fr";
  languageStore.setLanguage(next);
};

const handleLogout = () => {
  emit("close");
  authStore.logout();
};
</script>

<template>
  <transition
    enter-active-class="transition-opacity duration-150 ease-out"
    enter-from-class="opacity-0"
    enter-to-class="opacity-100"
    leave-active-class="transition-opacity duration-150 ease-in"
    leave-from-class="opacity-100"
    leave-to-class="opacity-0"
  >
    <div
      v-if="isOpen"
      class="fixed inset-0 z-50 md:hidden flex flex-col justify-end bg-black/45"
      @click.self="emit('close')"
    >
      <!-- Tiroir / Sheet : arrondi fin 3px sans aucune shadow -->
      <div
        class="w-full max-h-[86dvh] bg-card border-t border-x border-border rounded-t-[3px] flex flex-col overflow-hidden pb-safe"
      >
        <!-- Barre de préhension + en-tête profil -->
        <div class="shrink-0 px-4 pt-2 pb-3 border-b border-border bg-muted/20">
          <div class="mx-auto mb-2.5 h-1 w-9 rounded-full bg-border" />
          <div class="flex items-center justify-between gap-2">
            <button
              type="button"
              @click="openProfile"
              class="flex items-center gap-2.5 min-w-0 text-left rounded-[3px] py-1 px-1 -ml-1 active:bg-muted transition"
            >
              <div
                class="w-9 h-9 rounded-full bg-brand-light text-brand-dark border border-brand-border flex items-center justify-center font-bold text-xs uppercase shrink-0"
              >
                {{ authStore.currentUser.initial }}
              </div>
              <div class="min-w-0">
                <div class="text-xs font-bold text-foreground truncate">
                  {{ authStore.currentUser.prenom }}
                  {{ authStore.currentUser.nom }}
                </div>
                <div class="text-[11px] text-muted-foreground truncate">
                  {{ authStore.schoolName }} ·
                  {{ authStore.currentUser.roleLabel }}
                </div>
              </div>
            </button>

            <button
              type="button"
              @click="emit('close')"
              class="h-10 w-10 inline-flex items-center justify-center rounded-[3px] border border-border bg-background text-muted-foreground hover:text-foreground active:bg-muted transition shrink-0"
              aria-label="Fermer le menu"
            >
              <X class="w-4 h-4" />
            </button>
          </div>
        </div>

        <!-- Sections complètes du rôle -->
        <div class="flex-1 min-h-0 overflow-y-auto px-3 py-3 space-y-4">
          <div v-for="group in navGroups" :key="group.id" class="space-y-1.5">
            <div
              class="px-1 text-[10px] font-bold uppercase tracking-wider text-muted-foreground"
            >
              {{ group.title }}
            </div>
            <div class="grid grid-cols-1 gap-1">
              <button
                v-for="item in group.items"
                :key="item.id"
                type="button"
                @click="handlePickSection(item.id)"
                class="w-full min-h-[44px] flex items-center justify-between gap-3 px-3 py-2 rounded-[3px] border text-left text-xs transition active:scale-[0.99]"
                :class="[
                  !authStore.isModuleEnabled(item.id)
                    ? 'opacity-45 border-border/60 bg-muted/20 text-muted-foreground cursor-not-allowed'
                    : authStore.activeSection === item.id
                      ? 'bg-brand-light border-brand-border text-brand-dark font-semibold'
                      : 'bg-background border-border text-foreground hover:bg-muted/50',
                ]"
              >
                <div class="flex items-center gap-2.5 min-w-0">
                  <component
                    :is="item.icon"
                    class="w-4 h-4 shrink-0"
                    :class="
                      authStore.activeSection === item.id
                        ? 'text-brand'
                        : 'text-muted-foreground'
                    "
                  />
                  <span class="truncate">{{ item.label }}</span>
                </div>
                <Lock
                  v-if="!authStore.isModuleEnabled(item.id)"
                  class="w-3.5 h-3.5 shrink-0 text-muted-foreground"
                />
              </button>
            </div>
          </div>

          <!-- Actions rapides & paramètres -->
          <div class="pt-2 border-t border-border space-y-2">
            <div
              class="px-1 text-[10px] font-bold uppercase tracking-wider text-muted-foreground"
            >
              {{ t("Préférences & Outils") }}
            </div>
            <div class="grid grid-cols-2 gap-1.5">
              <button
                type="button"
                @click="openProfile"
                class="min-h-[44px] flex items-center gap-2 px-3 py-2 rounded-[3px] border border-border bg-background text-xs font-medium text-foreground active:bg-muted transition"
              >
                <User class="w-4 h-4 text-brand shrink-0" />
                <span class="truncate">{{ t("Mon Profil") }}</span>
              </button>

              <button
                v-if="
                  authStore.currentRole === 'manager' ||
                  authStore.currentRole === 'secretaire'
                "
                type="button"
                @click="openQr"
                class="min-h-[44px] flex items-center gap-2 px-3 py-2 rounded-[3px] border border-border bg-background text-xs font-medium text-foreground active:bg-muted transition"
              >
                <QrCode class="w-4 h-4 text-brand shrink-0" />
                <span class="truncate">{{ t("Accès Web & QR") }}</span>
              </button>

              <button
                v-if="authStore.currentRole === 'manager'"
                type="button"
                @click="handlePickSection('bon-plan')"
                class="min-h-[44px] flex items-center gap-2 px-3 py-2 rounded-[3px] border border-amber-500/30 bg-amber-500/10 text-xs font-semibold text-amber-700 dark:text-amber-300 active:opacity-80 transition"
              >
                <Sparkles class="w-4 h-4 shrink-0" />
                <span class="truncate">{{ t("Bon Plan Licence") }}</span>
              </button>

              <button
                type="button"
                @click="authStore.toggleDarkMode()"
                class="min-h-[44px] flex items-center gap-2 px-3 py-2 rounded-[3px] border border-border bg-background text-xs font-medium text-foreground active:bg-muted transition"
              >
                <Sun
                  v-if="authStore.isDark"
                  class="w-4 h-4 text-amber-400 shrink-0"
                />
                <Moon v-else class="w-4 h-4 text-muted-foreground shrink-0" />
                <span class="truncate">{{
                  authStore.isDark ? t("Thème Clair") : t("Thème Sombre")
                }}</span>
              </button>

              <button
                type="button"
                @click="toggleLang"
                class="min-h-[44px] flex items-center gap-2 px-3 py-2 rounded-[3px] border border-border bg-background text-xs font-medium text-foreground active:bg-muted transition"
              >
                <Globe class="w-4 h-4 text-muted-foreground shrink-0" />
                <span class="truncate">{{
                  languageStore.currentLanguage === "fr"
                    ? "Français (FR)"
                    : "English (EN)"
                }}</span>
              </button>

              <button
                type="button"
                @click="openHelp"
                class="min-h-[44px] flex items-center gap-2 px-3 py-2 rounded-[3px] border border-border bg-background text-xs font-medium text-foreground active:bg-muted transition"
              >
                <HelpCircle class="w-4 h-4 text-brand shrink-0" />
                <span class="truncate">{{ t("Besoin d'aide ?") }}</span>
              </button>
            </div>

            <button
              type="button"
              @click="handleLogout"
              class="w-full min-h-[44px] flex items-center justify-center gap-2 px-3 py-2 rounded-[3px] border border-rose-500/30 bg-rose-500/5 text-xs font-semibold text-rose-600 dark:text-rose-400 active:bg-rose-500/15 transition"
            >
              <LogOut class="w-4 h-4 shrink-0" />
              <span>{{ t("Se déconnecter") }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </transition>
</template>

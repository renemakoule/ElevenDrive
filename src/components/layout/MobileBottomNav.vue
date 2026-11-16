<script setup lang="ts">
import { ref, computed } from "vue";
import { useAuthStore } from "../../stores/authStore";
import { useLanguageStore } from "../../stores/languageStore";
import { getMobilePrimaryTabs } from "./mobileNavConfig";
import MobileNavDrawerSheet from "./MobileNavDrawerSheet.vue";
import { Menu } from "lucide-vue-next";

const authStore = useAuthStore();
const languageStore = useLanguageStore();
const isDrawerOpen = ref(false);

const t = (txt: string) => languageStore.t(txt);

const primaryTabs = computed(() =>
  getMobilePrimaryTabs(authStore.currentRole).map((tab) => ({
    ...tab,
    label: t(tab.label),
  })),
);

const isOverflowSectionActive = computed(() => {
  return !primaryTabs.value.some((tab) => tab.id === authStore.activeSection);
});

const selectTab = (sectionId: string) => {
  if (!authStore.isModuleEnabled(sectionId)) {
    return;
  }
  authStore.setActiveSection(sectionId);
};
</script>

<template>
  <nav
    v-if="primaryTabs.length > 0"
    aria-label="Navigation mobile"
    class="shrink-0 md:hidden w-full bg-card border-t border-border pb-safe select-none z-30"
  >
    <div class="grid grid-cols-5 h-[56px]">
      <button
        v-for="tab in primaryTabs"
        :key="tab.id"
        type="button"
        @click="selectTab(tab.id)"
        class="relative flex flex-col items-center justify-center gap-0.5 px-1 text-[10px] transition active:bg-muted/60 focus:outline-none"
        :class="[
          !authStore.isModuleEnabled(tab.id)
            ? 'opacity-40 text-muted-foreground'
            : authStore.activeSection === tab.id
              ? 'text-brand font-semibold'
              : 'text-muted-foreground font-medium',
        ]"
      >
        <span
          v-if="authStore.activeSection === tab.id"
          class="absolute top-0 left-1/2 -translate-x-1/2 h-[2px] w-7 rounded-b-[2px] bg-brand"
        />
        <component :is="tab.icon" class="w-[18px] h-[18px] shrink-0" />
        <span class="truncate max-w-full leading-tight">{{ tab.label }}</span>
      </button>

      <!-- Bouton Menu complet (Sheet tiroir 3px sans shadow) -->
      <button
        type="button"
        @click="isDrawerOpen = true"
        class="relative flex flex-col items-center justify-center gap-0.5 px-1 text-[10px] transition active:bg-muted/60 focus:outline-none"
        :class="
          isDrawerOpen || isOverflowSectionActive
            ? 'text-brand font-semibold'
            : 'text-muted-foreground font-medium'
        "
      >
        <span
          v-if="isOverflowSectionActive && !isDrawerOpen"
          class="absolute top-0 left-1/2 -translate-x-1/2 h-[2px] w-7 rounded-b-[2px] bg-brand"
        />
        <Menu class="w-[18px] h-[18px] shrink-0" />
        <span class="truncate max-w-full leading-tight">{{ t("Menu") }}</span>
      </button>
    </div>

    <MobileNavDrawerSheet
      :is-open="isDrawerOpen"
      @close="isDrawerOpen = false"
      @select-section="selectTab"
    />
  </nav>
</template>

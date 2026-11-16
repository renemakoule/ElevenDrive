<script setup lang="ts">
import { ref, computed, watch } from "vue";
import { useAuthStore } from "../../stores/authStore";
import { useLanguageStore } from "../../stores/languageStore";
import { getRoleNavGroups, type NavGroup } from "./sidebarNavConfig";
import SidebarFooterControls from "./SidebarFooterControls.vue";
import { ChevronDown, ChevronRight, Lock } from "lucide-vue-next";

const authStore = useAuthStore();
const languageStore = useLanguageStore();

const collapsedGroups = ref<Record<string, boolean>>({});

const toggleGroup = (groupId: string) => {
  collapsedGroups.value[groupId] = !collapsedGroups.value[groupId];
};

const isGroupCollapsed = (groupId: string) => !!collapsedGroups.value[groupId];

const navGroups = computed<NavGroup[]>(() => {
  const groups = getRoleNavGroups(authStore.currentRole);
  return groups.map((group) => ({
    ...group,
    title: languageStore.t(group.title),
    items: group.items.map((item) => ({
      ...item,
      label: languageStore.t(item.label),
    })),
  }));
});

watch(
  () => authStore.activeSection,
  (newSection) => {
    for (const group of navGroups.value) {
      if (group.items.some((item) => item.id === newSection)) {
        collapsedGroups.value[group.id] = false;
      }
    }
  },
  { immediate: true },
);

const selectNav = (id: string) => {
  if (!authStore.isModuleEnabled(id)) {
    alert(
      "🔒 Ce module a été verrouillé par le Super-Administrateur pour votre établissement.",
    );
    return;
  }
  authStore.setActiveSection(id);
};

watch(
  () => authStore.isModuleEnabled(authStore.activeSection),
  (enabled) => {
    if (!enabled && authStore.activeSection !== "overview") {
      authStore.setActiveSection("overview");
    }
  },
  { immediate: true },
);
</script>

<template>
  <aside
    class="hidden md:flex h-full bg-card flex-col justify-between transition-all duration-200 shrink-0 select-none z-10 relative overflow-hidden"
    :class="authStore.isSidebarCollapsed ? 'w-16' : 'w-56'"
  >
    <!-- Section Navigation Principale -->
    <div
      class="flex-1 min-h-0 py-2.5 flex flex-col overflow-y-auto sidebar-scroll"
    >
      <template v-if="!authStore.isSidebarCollapsed">
        <div v-for="(group, gIdx) in navGroups" :key="group.id" class="mb-1">
          <div
            v-if="gIdx > 0"
            class="my-1.5 border-t border-border mx-2.5"
          ></div>

          <button
            type="button"
            @click="toggleGroup(group.id)"
            class="w-full flex items-center justify-between px-3 py-1 text-[10px] font-bold text-muted-foreground/80 hover:text-foreground tracking-wider uppercase transition rounded-[3px] select-none group/gh"
            :title="
              isGroupCollapsed(group.id)
                ? languageStore.t('Déplier le groupe')
                : languageStore.t('Replier le groupe')
            "
          >
            <span class="truncate">{{ group.title }}</span>
            <span class="flex items-center space-x-1 shrink-0 ml-1">
              <span
                class="text-[9px] font-medium text-muted-foreground/60 group-hover/gh:text-muted-foreground"
              >
                {{ group.items.length }}
              </span>
              <component
                :is="isGroupCollapsed(group.id) ? ChevronRight : ChevronDown"
                class="w-3 h-3 text-muted-foreground/60 group-hover/gh:text-foreground transition-transform duration-150"
              />
            </span>
          </button>

          <nav
            v-show="!isGroupCollapsed(group.id)"
            class="space-y-0.5 px-2 mt-0.5"
          >
            <button
              v-for="item in group.items"
              :key="item.id"
              type="button"
              @click="selectNav(item.id)"
              class="w-full flex items-center justify-between px-3 py-1.5 rounded-[3px] text-xs transition active:scale-[0.98]"
              :class="[
                !authStore.isModuleEnabled(item.id)
                  ? 'opacity-40 hover:opacity-70 cursor-not-allowed text-muted-foreground'
                  : authStore.activeSection === item.id
                    ? 'bg-muted text-brand-dark font-semibold'
                    : 'text-muted-foreground hover:text-foreground hover:bg-brand-light',
              ]"
              :title="
                !authStore.isModuleEnabled(item.id)
                  ? `${item.label} (Verrouillé Super-Admin)`
                  : item.label
              "
            >
              <div class="flex items-center space-x-2.5 min-w-0">
                <component
                  :is="item.icon"
                  class="w-4 h-4 shrink-0"
                  :class="
                    !authStore.isModuleEnabled(item.id)
                      ? 'text-muted-foreground/60'
                      : authStore.activeSection === item.id
                        ? 'text-brand'
                        : 'text-muted-foreground'
                  "
                />
                <span class="truncate text-[12px]">
                  {{ item.label }}
                </span>
              </div>

              <div class="shrink-0 ml-1">
                <span
                  v-if="!authStore.isModuleEnabled(item.id)"
                  class="p-0.5 text-muted-foreground"
                >
                  <Lock class="w-3 h-3" />
                </span>
                <span
                  v-else-if="item.shortcut"
                  class="px-1.5 py-0.5 rounded text-[10px] font-mono font-bold bg-muted text-muted-foreground border border-border"
                >
                  {{ item.shortcut }}
                </span>
                <span
                  v-else-if="item.badge"
                  class="px-1.5 py-0.5 rounded-full text-[9px] font-medium"
                  :class="
                    item.badge.includes('alerte')
                      ? 'bg-amber-500/10 text-amber-600 dark:text-amber-400 border border-amber-500/20'
                      : 'bg-brand-light text-brand-dark'
                  "
                >
                  {{ item.badge }}
                </span>
              </div>
            </button>
          </nav>
        </div>
      </template>

      <template v-else>
        <div v-for="(group, gIdx) in navGroups" :key="group.id" class="px-2">
          <div v-if="gIdx > 0" class="my-1.5 border-t border-border mx-1"></div>
          <div class="space-y-1">
            <button
              v-for="item in group.items"
              :key="item.id"
              type="button"
              @click="selectNav(item.id)"
              class="w-full flex items-center justify-center p-2 rounded-[3px] text-xs transition active:scale-[0.98]"
              :class="[
                !authStore.isModuleEnabled(item.id)
                  ? 'opacity-40 cursor-not-allowed text-muted-foreground'
                  : authStore.activeSection === item.id
                    ? 'bg-muted text-brand-dark font-semibold'
                    : 'text-muted-foreground hover:text-foreground hover:bg-brand-light',
              ]"
              :title="
                !authStore.isModuleEnabled(item.id)
                  ? `${group.title} • ${item.label} (Verrouillé Super-Admin)`
                  : `${group.title} • ${item.label}`
              "
            >
              <component
                :is="!authStore.isModuleEnabled(item.id) ? Lock : item.icon"
                class="w-4 h-4 shrink-0"
                :class="
                  !authStore.isModuleEnabled(item.id)
                    ? 'text-muted-foreground/60'
                    : authStore.activeSection === item.id
                      ? 'text-brand'
                      : 'text-muted-foreground'
                "
              />
            </button>
          </div>
        </div>
      </template>
    </div>

    <SidebarFooterControls @select-nav="selectNav" />
  </aside>
</template>

<style scoped>
.sidebar-scroll {
  scrollbar-width: thin;
  scrollbar-color: transparent transparent;
  transition: scrollbar-color 0.2s ease;
}
.sidebar-scroll::-webkit-scrollbar {
  width: 4px;
}
.sidebar-scroll::-webkit-scrollbar-track {
  background: transparent;
}
.sidebar-scroll::-webkit-scrollbar-thumb {
  background: transparent;
  border-radius: 4px;
  transition: background-color 0.2s ease;
}
aside:hover .sidebar-scroll {
  scrollbar-color: var(--border) transparent;
}
aside:hover .sidebar-scroll::-webkit-scrollbar-thumb {
  background: var(--border);
}
aside:hover .sidebar-scroll::-webkit-scrollbar-thumb:hover {
  background: var(--muted-foreground);
}
</style>

<script setup lang="ts">
import {
  Building2,
  Users,
  GraduationCap,
  ShieldCheck,
  Layers,
  History,
  BellRing,
  Megaphone,
  Star,
} from "lucide-vue-next";
import type { SuperAdminStats } from "../../types";

export type SuperAdminTab =
  | "schools"
  | "licences"
  | "permissions"
  | "users"
  | "candidats"
  | "audit"
  | "alerts"
  | "announcements"
  | "feedbacks";

defineProps<{
  activeTab: SuperAdminTab;
  stats: SuperAdminStats;
}>();

const emit = defineEmits<{
  (e: "select-tab", tab: SuperAdminTab): void;
}>();

const formatFCFA = (val: number) => {
  return (val || 0).toLocaleString("fr-FR");
};
</script>

<template>
  <aside
    class="h-full w-56 bg-card flex flex-col justify-between shrink-0 overflow-hidden select-none"
  >
    <!-- Zone de navigation et sections interne avec scrollbar au survol uniquement -->
    <div class="flex-1 min-h-0 p-3 space-y-4 overflow-y-auto sidebar-scroll">
      <div>
        <span
          class="text-[10px] uppercase font-bold text-muted-foreground tracking-wider px-2"
        >
          Gestion Plateforme
        </span>
        <nav class="mt-2 space-y-1">
          <!-- 1. Menu Auto-Écoles -->
          <button
            type="button"
            @click="emit('select-tab', 'schools')"
            :class="[
              'w-full flex items-center gap-2 px-2.5 py-1.5 rounded-[4px] text-xs font-medium transition-colors cursor-pointer',
              activeTab === 'schools'
                ? 'bg-muted text-foreground border border-border shadow-xs'
                : 'text-muted-foreground hover:bg-muted/60 hover:text-foreground border border-transparent',
            ]"
          >
            <Building2
              class="w-4 h-4"
              :class="
                activeTab === 'schools' ? 'text-brand' : 'text-muted-foreground'
              "
            />
            <span>Auto-Écoles</span>
            <span
              :class="[
                'ml-auto text-[10px] font-mono px-1.5 py-0.5 rounded-[3px]',
                activeTab === 'schools'
                  ? 'bg-brand/10 text-brand font-semibold'
                  : 'bg-muted text-muted-foreground',
              ]"
            >
              {{ stats.total_schools }}
            </span>
          </button>

          <!-- 2. Menu Licences & Quotas Machines -->
          <button
            type="button"
            @click="emit('select-tab', 'licences')"
            :class="[
              'w-full flex items-center gap-2 px-2.5 py-1.5 rounded-[4px] text-xs font-medium transition-colors cursor-pointer',
              activeTab === 'licences'
                ? 'bg-muted text-foreground border border-border shadow-xs'
                : 'text-muted-foreground hover:bg-muted/60 hover:text-foreground border border-transparent',
            ]"
          >
            <ShieldCheck
              class="w-4 h-4"
              :class="
                activeTab === 'licences'
                  ? 'text-brand'
                  : 'text-muted-foreground'
              "
            />
            <span>Postes & Licences</span>
          </button>

          <!-- 3. Menu Contrôle d'Accès & Droits -->
          <button
            type="button"
            @click="emit('select-tab', 'permissions')"
            :class="[
              'w-full flex items-center gap-2 px-2.5 py-1.5 rounded-[4px] text-xs font-medium transition-colors cursor-pointer',
              activeTab === 'permissions'
                ? 'bg-muted text-foreground border border-border shadow-xs'
                : 'text-muted-foreground hover:bg-muted/60 hover:text-foreground border border-transparent',
            ]"
          >
            <Layers
              class="w-4 h-4"
              :class="
                activeTab === 'permissions'
                  ? 'text-brand'
                  : 'text-muted-foreground'
              "
            />
            <span>Contrôle d'Accès</span>
          </button>

          <!-- 4. Menu Utilisateurs Globaux -->
          <button
            type="button"
            @click="emit('select-tab', 'users')"
            :class="[
              'w-full flex items-center gap-2 px-2.5 py-1.5 rounded-[4px] text-xs font-medium transition-colors cursor-pointer',
              activeTab === 'users'
                ? 'bg-muted text-foreground border border-border shadow-xs'
                : 'text-muted-foreground hover:bg-muted/60 hover:text-foreground border border-transparent',
            ]"
          >
            <Users
              class="w-4 h-4"
              :class="
                activeTab === 'users' ? 'text-brand' : 'text-muted-foreground'
              "
            />
            <span>Collaborateurs</span>
            <span
              :class="[
                'ml-auto text-[10px] font-mono px-1.5 py-0.5 rounded-[3px]',
                activeTab === 'users'
                  ? 'bg-brand/10 text-brand font-semibold'
                  : 'bg-muted text-muted-foreground',
              ]"
            >
              {{ stats.total_users }}
            </span>
          </button>

          <!-- 5. Menu Candidats Inscrits -->
          <button
            type="button"
            @click="emit('select-tab', 'candidats')"
            :class="[
              'w-full flex items-center gap-2 px-2.5 py-1.5 rounded-[4px] text-xs font-medium transition-colors cursor-pointer',
              activeTab === 'candidats'
                ? 'bg-muted text-foreground border border-border shadow-xs'
                : 'text-muted-foreground hover:bg-muted/60 hover:text-foreground border border-transparent',
            ]"
          >
            <GraduationCap
              class="w-4 h-4"
              :class="
                activeTab === 'candidats'
                  ? 'text-brand'
                  : 'text-muted-foreground'
              "
            />
            <span>Candidats Inscrits</span>
            <span
              :class="[
                'ml-auto text-[10px] font-mono px-1.5 py-0.5 rounded-[3px]',
                activeTab === 'candidats'
                  ? 'bg-brand/10 text-brand font-semibold'
                  : 'bg-muted text-muted-foreground',
              ]"
            >
              {{ stats.total_candidats }}
            </span>
          </button>

          <!-- 6. Menu Journal d'Audit Système -->
          <button
            type="button"
            @click="emit('select-tab', 'audit')"
            :class="[
              'w-full flex items-center gap-2 px-2.5 py-1.5 rounded-[4px] text-xs font-medium transition-colors cursor-pointer',
              activeTab === 'audit'
                ? 'bg-muted text-foreground border border-border shadow-xs'
                : 'text-muted-foreground hover:bg-muted/60 hover:text-foreground border border-transparent',
            ]"
          >
            <History
              class="w-4 h-4"
              :class="
                activeTab === 'audit' ? 'text-brand' : 'text-muted-foreground'
              "
            />
            <span>Journal d'Audit</span>
          </button>
        </nav>
      </div>

      <!-- Section Diffusion & Retours Réseau (Alerte, Annonce, Avis) -->
      <div
        class="p-3 rounded-[4px] bg-background border border-border space-y-2.5"
      >
        <div class="flex items-center justify-between">
          <span
            class="text-[10px] uppercase font-bold text-muted-foreground tracking-wider"
          >
            Diffusion & Retours
          </span>
          <span class="flex h-2 w-2 relative">
            <span
              class="animate-ping absolute inline-flex h-full w-full rounded-full bg-brand-400 opacity-75"
            ></span>
            <span
              class="relative inline-flex rounded-full h-2 w-2 bg-brand-500"
            ></span>
          </span>
        </div>

        <nav class="space-y-1">
          <!-- 1. Alerte -->
          <button
            type="button"
            @click="emit('select-tab', 'alerts')"
            :class="[
              'w-full flex items-center gap-2 px-2.5 py-1.5 rounded-[4px] text-xs font-medium transition-colors cursor-pointer',
              activeTab === 'alerts'
                ? 'bg-muted text-foreground border border-border shadow-xs font-semibold'
                : 'text-muted-foreground hover:bg-muted/60 hover:text-foreground border border-transparent',
            ]"
          >
            <BellRing
              class="w-4 h-4"
              :class="
                activeTab === 'alerts'
                  ? 'text-amber-500'
                  : 'text-muted-foreground'
              "
            />
            <span>Alerte</span>
            <span
              :class="[
                'ml-auto text-[9px] px-1.5 py-0.5 rounded-[3px] font-semibold',
                activeTab === 'alerts'
                  ? 'bg-amber-500/20 text-amber-600 dark:text-amber-400'
                  : 'bg-muted text-muted-foreground',
              ]"
            >
              Direct
            </span>
          </button>

          <!-- 2. Annonce -->
          <button
            type="button"
            @click="emit('select-tab', 'announcements')"
            :class="[
              'w-full flex items-center gap-2 px-2.5 py-1.5 rounded-[4px] text-xs font-medium transition-colors cursor-pointer',
              activeTab === 'announcements'
                ? 'bg-muted text-foreground border border-border shadow-xs font-semibold'
                : 'text-muted-foreground hover:bg-muted/60 hover:text-foreground border border-transparent',
            ]"
          >
            <Megaphone
              class="w-4 h-4"
              :class="
                activeTab === 'announcements'
                  ? 'text-blue-500'
                  : 'text-muted-foreground'
              "
            />
            <span>Annonce</span>
            <span
              :class="[
                'ml-auto text-[9px] px-1.5 py-0.5 rounded-[3px] font-semibold',
                activeTab === 'announcements'
                  ? 'bg-blue-500/20 text-blue-600 dark:text-blue-400'
                  : 'bg-muted text-muted-foreground',
              ]"
            >
              Médias
            </span>
          </button>

          <!-- 3. Avis -->
          <button
            type="button"
            @click="emit('select-tab', 'feedbacks')"
            :class="[
              'w-full flex items-center gap-2 px-2.5 py-1.5 rounded-[4px] text-xs font-medium transition-colors cursor-pointer',
              activeTab === 'feedbacks'
                ? 'bg-muted text-foreground border border-border shadow-xs font-semibold'
                : 'text-muted-foreground hover:bg-muted/60 hover:text-foreground border border-transparent',
            ]"
          >
            <Star
              class="w-4 h-4"
              :class="
                activeTab === 'feedbacks'
                  ? 'text-amber-400 fill-amber-400'
                  : 'text-muted-foreground'
              "
            />
            <span>Avis</span>
            <span
              :class="[
                'ml-auto text-[9px] font-mono px-1.5 py-0.5 rounded-[3px] font-bold',
                activeTab === 'feedbacks'
                  ? 'bg-brand/10 text-brand'
                  : 'bg-muted text-muted-foreground',
              ]"
            >
              4.8 ★
            </span>
          </button>
        </nav>
      </div>

      <!-- Métriques Globales Résumées -->
      <div
        class="p-3 rounded-[4px] bg-background border border-border space-y-2"
      >
        <span
          class="text-[10px] uppercase font-bold text-muted-foreground tracking-wider"
        >
          Volume Transigé
        </span>
        <div class="text-sm font-bold font-mono text-foreground">
          {{ formatFCFA(stats.total_volume_paiement) }}
          <span class="text-[10px] text-brand font-semibold">FCFA</span>
        </div>
        <div
          class="text-[11px] text-muted-foreground flex items-center justify-between"
        >
          <span>Établissements actifs</span>
          <span class="font-mono text-brand font-medium">
            {{ stats.active_schools }} / {{ stats.total_schools }}
          </span>
        </div>
      </div>
    </div>

    <!-- Footer Système toujours ancré en bas sans scroll -->
    <div
      class="shrink-0 p-3 border-t border-border/40 text-[10px] text-muted-foreground font-mono bg-card"
    >
      <span>ElevenDrive Core v1.2</span>
    </div>
  </aside>
</template>

<style scoped>
/* Scrollbar invisible par défaut */
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

/* Apparition de la barre de scroll UNIQUEMENT au survol de la souris dans la sidebar */
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

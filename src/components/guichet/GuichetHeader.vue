<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { useGuichetStore } from "../../stores/guichetStore";
import { useAuthStore } from "../../stores/authStore";
import { neonSyncEngine, type SyncState } from "../../services/neonSyncService";
import SyncStatusModal from "./SyncStatusModal.vue";
import { WifiOff, Clock, Cloud, RefreshCw } from "lucide-vue-next";

const store = useGuichetStore();
const authStore = useAuthStore();
const currentTime = ref("");
const isSyncModalOpen = ref(false);
const syncState = ref<SyncState>(neonSyncEngine.getState());

const currentUserName = computed(() => {
  if (authStore.currentUser?.nom) {
    return `${authStore.currentUser.prenom} ${authStore.currentUser.nom}`;
  }
  return "Agent Guichet";
});
const schoolName = computed(() => authStore.schoolName || "AutoSchool");
const agenceName = computed(
  () => authStore.currentAgence?.nom || "Agence Guichet",
);
const schoolInitials = computed(() =>
  (schoolName.value || "AUT").slice(0, 3).toUpperCase(),
);

let syncUnsubscribe: (() => void) | null = null;
let timer: number;

const updateClock = () => {
  currentTime.value = new Date().toLocaleTimeString("fr-FR", {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });
};

onMounted(() => {
  updateClock();
  timer = window.setInterval(updateClock, 1000);
  syncUnsubscribe = neonSyncEngine.subscribe((state) => {
    syncState.value = state;
  });
});

onUnmounted(() => {
  clearInterval(timer);
  if (syncUnsubscribe) syncUnsubscribe();
});

/* ───────────────────────── Affichage ───────────────────────── */
const enAttente = computed(() => store.stats.evenements_sync_en_attente ?? 0);
const tiroir = computed(() =>
  (store.stats.total_especes_theorique ?? 0).toLocaleString("fr-FR"),
);

const sync = computed(() => {
  if (syncState.value.isSyncing) {
    return {
      label: "Synchronisation…",
      box: "border-sky-200 bg-sky-50 text-sky-800",
      icon: RefreshCw,
      spin: true,
    };
  }
  if (!syncState.value.isOnline) {
    return {
      label: "Hors ligne",
      box: "border-rose-200 bg-rose-50 text-rose-800",
      icon: WifiOff,
      spin: false,
    };
  }
  if (enAttente.value > 0) {
    return {
      label: "Cloud",
      box: "border-amber-200 bg-amber-50 text-amber-800",
      icon: Cloud,
      spin: false,
    };
  }
  return {
    label: "Synchronisé",
    box: "border-brand-border bg-brand-light text-brand-dark",
    icon: Cloud,
    spin: false,
  };
});
</script>

<template>
  <header
    class="flex items-center justify-between gap-4 border-b border-border bg-card px-4 py-2 text-xs text-foreground"
  >
    <!-- ═════════ Identité de l'agence ═════════ -->
    <div class="flex min-w-0 items-center gap-3">
      <div
        class="flex h-8 w-8 shrink-0 items-center justify-center rounded-[5px] bg-brand-dark text-[11px] font-semibold tracking-wide text-white"
        aria-hidden="true"
      >
        {{ schoolInitials }}
      </div>
      <div class="min-w-0">
        <div class="truncate text-sm font-semibold text-foreground">
          {{ schoolName }}
        </div>
        <div class="truncate text-[11px] text-muted-foreground">
          {{ agenceName }}
        </div>
      </div>
    </div>

    <!-- ═════════ État du guichet ═════════ -->
    <div class="flex shrink-0 items-center gap-4">
      <!-- Agent connecté -->
      <div class="hidden items-center gap-1.5 text-muted-foreground md:flex">
        <span class="h-1.5 w-1.5 rounded-full bg-brand" />
        Agent
        <span class="font-medium text-foreground">{{ currentUserName }}</span>
      </div>

      <span class="hidden h-4 w-px bg-border md:block" />

      <!-- Tiroir-caisse -->
      <div class="hidden items-center gap-1.5 text-muted-foreground sm:flex">
        Tiroir
        <span class="font-semibold tabular-nums text-foreground"
          >{{ tiroir }} FCFA</span
        >
      </div>

      <span class="hidden h-4 w-px bg-border sm:block" />

      <!-- Synchronisation -->
      <button
        type="button"
        @click="isSyncModalOpen = true"
        title="Gérer la synchronisation cloud"
        class="inline-flex h-7 items-center gap-1.5 rounded-full border px-2.5 font-medium transition hover:brightness-95 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40"
        :class="sync.box"
      >
        <component
          :is="sync.icon"
          class="h-3.5 w-3.5"
          :class="{ 'animate-spin': sync.spin }"
        />
        {{ sync.label }}
        <span v-if="enAttente > 0" class="tabular-nums opacity-80"
          >· {{ enAttente }} en attente</span
        >
      </button>

      <!-- Horloge -->
      <!-- <div class="flex items-center gap-1.5 text-muted-foreground">
        <Clock class="h-3.5 w-3.5 text-muted-foreground" />
        <span class="font-medium tabular-nums text-foreground">{{ currentTime }}</span>
      </div> -->
    </div>

    <!-- Modale de synchronisation -->
    <SyncStatusModal v-if="isSyncModalOpen" @close="isSyncModalOpen = false" />
  </header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { neonSyncEngine, type SyncState } from "../../services/neonSyncService";
import { useGuichetStore } from "../../stores/guichetStore";
import { RefreshCw, X, Zap, Eye, EyeOff, ExternalLink } from "lucide-vue-next";

defineProps<{
  isOpen?: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
}>();

const store = useGuichetStore();
const syncState = ref<SyncState>(neonSyncEngine.getState());
const customUrl = ref(syncState.value.neonEndpoint);
const showUrl = ref(false);
const isTesting = ref(false);
const testResult = ref<{
  success: boolean;
  latencyMs: number;
  message: string;
} | null>(null);

let unsubscribe: (() => void) | null = null;

const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === "Escape") emit("close");
};

onMounted(() => {
  unsubscribe = neonSyncEngine.subscribe((state) => {
    syncState.value = state;
  });
  window.addEventListener("keydown", handleKeydown);
});

onUnmounted(() => {
  if (unsubscribe) unsubscribe();
  window.removeEventListener("keydown", handleKeydown);
});

const enAttente = computed(() => store.stats.evenements_sync_en_attente ?? 0);

const lastSync = computed(() =>
  syncState.value.lastSyncTimestamp
    ? new Date(syncState.value.lastSyncTimestamp).toLocaleTimeString("fr-FR")
    : "aucune",
);

const status = computed(() => [
  {
    label: "Réseau",
    value: syncState.value.isOnline ? "En ligne" : "Hors ligne",
    dot: syncState.value.isOnline ? "bg-brand" : "bg-destructive",
  },
  {
    label: "Cloud",
    value: syncState.value.neonConnected ? "Connecté" : "Déconnecté",
    dot: syncState.value.neonConnected ? "bg-brand" : "bg-amber-500",
  },
  {
    label: "File locale",
    value: `${enAttente.value} en attente`,
    dot: enAttente.value > 0 ? "bg-amber-500" : "bg-brand",
  },
]);

const schema = [
  { label: "Tables", value: "14 tables relationnelles" },
  { label: "Isolation", value: "RLS multi-agence" },
  { label: "Ingestion", value: "process_sync_batch()" },
  { label: "Finance", value: "Append-only, sans conflit" },
];

const handleSaveUrl = () => {
  neonSyncEngine.setNeonEndpoint(customUrl.value);
  store.showToast("Chaîne de connexion mise à jour.", "success");
};

const handleTestConnection = async () => {
  isTesting.value = true;
  testResult.value = null;
  try {
    testResult.value = await neonSyncEngine.testNeonConnection(customUrl.value);
  } finally {
    isTesting.value = false;
  }
};

const handleForceSync = async () => {
  const res = await neonSyncEngine.triggerAutoSync();
  if (res.statut === "SUCCES") {
    store.stats.evenements_sync_en_attente = 0;
    store.showToast(res.message || "Synchronisation réussie.", "success");
  } else {
    store.showToast(res.message || "Erreur de synchronisation.", "error");
  }
};
</script>

<template>
  <div
    v-if="isOpen !== false"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 backdrop-blur-sm"
    role="dialog"
    aria-modal="true"
    aria-labelledby="sync-title"
    @click.self="emit('close')"
  >
    <div
      class="flex max-h-[92vh] w-full max-w-xl flex-col overflow-hidden rounded-[5px] border border-border bg-card"
    >
      <header
        class="flex items-center justify-between gap-3 border-b border-border px-4 py-3"
      >
        <div class="min-w-0">
          <h2 id="sync-title" class="text-sm font-semibold text-foreground">
            Synchronisation cloud
          </h2>
          <p class="text-[11px] text-muted-foreground">
            Réplication des données locales vers PostgreSQL (Neon).
          </p>
        </div>
        <button
          type="button"
          @click="emit('close')"
          class="flex h-7 w-7 items-center justify-center rounded-[4px] text-muted-foreground transition hover:bg-muted hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-dark/40"
          aria-label="Fermer"
        >
          <X class="h-4 w-4" />
        </button>
      </header>

      <div class="flex-1 divide-y divide-border overflow-y-auto">
        <!-- État -->
        <section class="grid grid-cols-3 divide-x divide-border">
          <div v-for="s in status" :key="s.label" class="px-4 py-3">
            <p class="text-[11px] text-muted-foreground">{{ s.label }}</p>
            <p
              class="mt-0.5 flex items-center gap-1.5 text-xs font-medium text-foreground"
            >
              <span class="h-1.5 w-1.5 rounded-full" :class="s.dot" />
              {{ s.value }}
            </p>
          </div>
        </section>

        <!-- Connexion -->
        <section class="space-y-3 px-4 py-4">
          <div class="flex items-center justify-between gap-3">
            <label for="sync-url" class="text-xs font-semibold text-foreground"
              >Chaîne de connexion</label
            >
            <a
              href="https://console.neon.tech"
              target="_blank"
              rel="noopener noreferrer"
              class="inline-flex items-center gap-1 text-[11px] font-medium text-brand-dark hover:underline"
            >
              Console Neon
              <ExternalLink class="h-3 w-3" />
            </a>
          </div>

          <div class="relative">
            <input
              id="sync-url"
              v-model="customUrl"
              :type="showUrl ? 'text' : 'password'"
              autocomplete="off"
              spellcheck="false"
              placeholder="postgresql://user:password@ep-xyz.eu-central-1.aws.neon.tech/neondb"
              class="h-8 w-full rounded-[4px] border border-border bg-card pl-2.5 pr-9 font-mono text-xs text-foreground transition placeholder:text-muted-foreground focus:border-brand-dark focus:outline-none focus:ring-2 focus:ring-brand-dark/20"
            />
            <button
              type="button"
              @click="showUrl = !showUrl"
              class="absolute right-1 top-1/2 flex h-6 w-6 -translate-y-1/2 items-center justify-center rounded text-muted-foreground transition hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-dark/40"
              :aria-label="showUrl ? 'Masquer la chaîne' : 'Afficher la chaîne'"
            >
              <EyeOff v-if="showUrl" class="h-3.5 w-3.5" />
              <Eye v-else class="h-3.5 w-3.5" />
            </button>
          </div>

          <div class="flex flex-wrap items-center justify-between gap-2">
            <div class="flex items-center gap-2">
              <button
                type="button"
                @click="handleTestConnection"
                :disabled="isTesting"
                class="inline-flex h-8 items-center gap-1.5 rounded-[4px] border border-border bg-card px-3 text-xs font-medium text-foreground transition hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-dark/40 disabled:opacity-50"
              >
                <Zap class="h-3.5 w-3.5 text-muted-foreground" />
                {{ isTesting ? "Test en cours…" : "Tester" }}
              </button>
              <button
                type="button"
                @click="handleSaveUrl"
                class="inline-flex h-8 items-center rounded-[4px] border border-border bg-card px-3 text-xs font-medium text-foreground transition hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-dark/40"
              >
                Enregistrer
              </button>
            </div>

            <button
              type="button"
              @click="handleForceSync"
              :disabled="syncState.isSyncing"
              class="inline-flex h-8 items-center gap-1.5 rounded-[4px] bg-brand-dark px-3 text-xs font-medium text-white transition hover:bg-brand focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-dark/40 focus-visible:ring-offset-1 disabled:cursor-not-allowed disabled:opacity-60"
            >
              <RefreshCw
                class="h-3.5 w-3.5"
                :class="{ 'animate-spin': syncState.isSyncing }"
              />
              {{ syncState.isSyncing ? "Synchronisation…" : "Synchroniser" }}
            </button>
          </div>

          <p
            v-if="testResult"
            role="status"
            class="flex items-center gap-2 rounded-lg px-3 py-2 text-xs"
            :class="
              testResult.success
                ? 'bg-brand/10 text-brand-dark border border-brand-border'
                : 'bg-destructive/10 text-destructive border border-destructive/20'
            "
          >
            <span
              class="h-1.5 w-1.5 shrink-0 rounded-full"
              :class="testResult.success ? 'bg-brand' : 'bg-destructive'"
            />
            <span class="min-w-0 flex-1">{{ testResult.message }}</span>
            <span
              v-if="testResult.success"
              class="shrink-0 tabular-nums opacity-70"
              >{{ testResult.latencyMs }} ms</span
            >
          </p>
        </section>

        <!-- Schéma -->
        <section class="px-4 py-4">
          <h3 class="mb-2 text-xs font-semibold text-foreground">
            Schéma déployé
          </h3>
          <dl class="grid grid-cols-2 gap-x-6 gap-y-1.5 text-xs">
            <div
              v-for="i in schema"
              :key="i.label"
              class="flex items-baseline justify-between gap-3"
            >
              <dt class="text-muted-foreground">{{ i.label }}</dt>
              <dd class="truncate font-medium text-foreground">
                {{ i.value }}
              </dd>
            </div>
          </dl>
        </section>

        <!-- Historique -->
        <section class="px-4 py-4">
          <div class="mb-2 flex items-baseline justify-between">
            <h3 class="text-xs font-semibold text-foreground">Lots transmis</h3>
            <span class="text-[11px] text-muted-foreground"
              >Dernière synchro : {{ lastSync }}</span
            >
          </div>

          <p
            v-if="syncState.syncHistory.length === 0"
            class="py-3 text-center text-xs text-muted-foreground"
          >
            Aucun lot synchronisé pour le moment.
          </p>

          <ul v-else class="max-h-40 divide-y divide-border overflow-y-auto">
            <li
              v-for="item in syncState.syncHistory"
              :key="item.id"
              class="flex items-center justify-between gap-3 py-2 text-xs"
            >
              <span class="flex items-center gap-2">
                <span class="font-mono font-medium text-foreground">{{
                  item.timestamp
                }}</span>
                <span class="text-muted-foreground"
                  >{{ item.count }} transaction(s)</span
                >
              </span>
              <span class="flex items-center gap-2">
                <span class="truncate text-[11px] text-muted-foreground">{{
                  item.details
                }}</span>
                <span
                  class="flex items-center gap-1.5 text-[11px] font-medium text-brand-dark"
                >
                  <span class="h-1.5 w-1.5 rounded-full bg-brand" />
                  {{ item.statut }}
                </span>
              </span>
            </li>
          </ul>
        </section>
      </div>

      <footer class="flex justify-end border-t border-border px-4 py-3">
        <button
          type="button"
          @click="emit('close')"
          class="inline-flex h-8 items-center gap-2 rounded-[4px] border border-border bg-card px-3 text-xs font-medium text-foreground transition hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-dark/40"
        >
          Fermer
          <kbd
            class="rounded bg-muted px-1 font-mono text-[10px] text-muted-foreground"
            >Échap</kbd
          >
        </button>
      </footer>
    </div>
  </div>
</template>

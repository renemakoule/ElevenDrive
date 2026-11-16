<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";
import { supportService } from "../../services/supportService";
import type { NetworkAlert, NetworkAnnouncement } from "../../types";
import {
  X,
  AlertTriangle,
  Megaphone,
  Loader2,
  Star,
  RotateCcw,
} from "lucide-vue-next";
import NetworkAnnouncementsTab from "./network/NetworkAnnouncementsTab.vue";
import NetworkAlertsTab from "./network/NetworkAlertsTab.vue";
import NetworkFeedbackTab from "./network/NetworkFeedbackTab.vue";
import NetworkImageViewer, {
  type PreviewImageData,
} from "./network/NetworkImageViewer.vue";

const props = defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
}>();

const alerts = ref<NetworkAlert[]>([]);
const announcements = ref<NetworkAnnouncement[]>([]);
const isLoading = ref(false);
const activeTab = ref<"announcements" | "alerts" | "feedback">("announcements");

// Gestion des panneaux ads masqués individuellement
const dismissedItems = ref<Set<string>>(new Set());

// Visionneuse Plein Écran (Lightbox)
const previewImage = ref<PreviewImageData | null>(null);
const currentImageIndex = ref(0);

const loadDismissed = () => {
  try {
    const raw = localStorage.getItem("autoschool_dismissed_network_items");
    if (raw) {
      const arr = JSON.parse(raw);
      if (Array.isArray(arr)) {
        dismissedItems.value = new Set(arr);
      }
    }
  } catch {
    // ignore
  }
};

const saveDismissed = () => {
  try {
    localStorage.setItem(
      "autoschool_dismissed_network_items",
      JSON.stringify(Array.from(dismissedItems.value)),
    );
  } catch {
    // ignore
  }
};

const dismissPanel = (id: string) => {
  dismissedItems.value.add(id);
  saveDismissed();
};

const restoreDismissedPanels = () => {
  dismissedItems.value.clear();
  saveDismissed();
};

const visibleAnnouncements = computed(() => {
  return announcements.value.filter((a) => !dismissedItems.value.has(a.id));
});

const visibleAlerts = computed(() => {
  return alerts.value.filter((a) => !dismissedItems.value.has(a.id));
});

const totalDismissedCount = computed(() => {
  let count = 0;
  for (const a of announcements.value) {
    if (dismissedItems.value.has(a.id)) count++;
  }
  for (const alt of alerts.value) {
    if (dismissedItems.value.has(alt.id)) count++;
  }
  return count;
});

const imageAnnouncements = computed(() => {
  return visibleAnnouncements.value.filter((a) => Boolean(a.image_url));
});

const formatDate = (iso: string) => {
  try {
    const d = new Date(iso);
    return `${d.toLocaleDateString("fr-FR")} à ${d.toLocaleTimeString("fr-FR", { hour: "2-digit", minute: "2-digit" })}`;
  } catch {
    return iso;
  }
};

const openImagePreview = (ann: NetworkAnnouncement) => {
  if (!ann.image_url) return;
  const idx = imageAnnouncements.value.findIndex((a) => a.id === ann.id);
  currentImageIndex.value = idx >= 0 ? idx : 0;
  previewImage.value = {
    url: ann.image_url,
    title: ann.titre,
    category: ann.categorie,
    content: ann.contenu,
    date: formatDate(ann.date_publication),
    actionUrl: ann.lien_action,
  };
};

const closeImagePreview = () => {
  previewImage.value = null;
};

const goToPreviousImage = () => {
  if (imageAnnouncements.value.length <= 1) return;
  const newIdx =
    (currentImageIndex.value - 1 + imageAnnouncements.value.length) %
    imageAnnouncements.value.length;
  currentImageIndex.value = newIdx;
  const target = imageAnnouncements.value[newIdx];
  previewImage.value = {
    url: target.image_url!,
    title: target.titre,
    category: target.categorie,
    content: target.contenu,
    date: formatDate(target.date_publication),
    actionUrl: target.lien_action,
  };
};

const goToNextImage = () => {
  if (imageAnnouncements.value.length <= 1) return;
  const newIdx =
    (currentImageIndex.value + 1) % imageAnnouncements.value.length;
  currentImageIndex.value = newIdx;
  const target = imageAnnouncements.value[newIdx];
  previewImage.value = {
    url: target.image_url!,
    title: target.titre,
    category: target.categorie,
    content: target.contenu,
    date: formatDate(target.date_publication),
    actionUrl: target.lien_action,
  };
};

const loadData = async () => {
  isLoading.value = true;
  try {
    loadDismissed();
    const [rawAlerts, rawAnns] = await Promise.all([
      supportService.getActiveAlerts(),
      supportService.getAnnouncements(),
    ]);
    alerts.value = rawAlerts || [];
    announcements.value = rawAnns || [];
  } catch {
    // ignore
  } finally {
    isLoading.value = false;
  }
};

onMounted(() => {
  loadData();
});

watch(
  () => props.isOpen,
  (open) => {
    if (open) {
      loadData();
    }
  },
);
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-xs p-3 sm:p-4 select-none animate-in fade-in duration-150"
    @click.self="emit('close')"
  >
    <div
      class="bg-card border border-border rounded-[4px] max-w-4xl w-full overflow-hidden flex flex-col max-h-[90vh] shadow-2xl animate-in zoom-in-95 duration-200"
    >
      <!-- En-tête de la Modal -->
      <div
        class="p-4 bg-muted/40 border-b border-border flex items-center justify-between shrink-0"
      >
        <div class="flex items-center space-x-3">
          <div
            class="w-9 h-9 rounded-[4px] bg-brand-light text-brand-dark flex items-center justify-center font-bold shadow-xs"
          >
            <Megaphone class="w-4 h-4 text-brand" />
          </div>
          <div>
            <div
              class="font-bold text-foreground text-sm flex items-center gap-2"
            >
              <span>Communications Réseau & Retours</span>
              <span
                class="text-[9px] uppercase font-mono px-1.5 py-0.2 rounded bg-brand/10 text-brand font-semibold"
              >
                ElevenDrive Ads
              </span>
            </div>
            <div class="text-[11px] text-muted-foreground">
              Annonces officielles, alertes système et notation du réseau
            </div>
          </div>
        </div>

        <button
          type="button"
          @click="emit('close')"
          class="p-1.5 rounded-[4px] text-muted-foreground hover:text-foreground hover:bg-muted transition cursor-pointer"
          title="Fermer la fenêtre"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Barre d'onglets principale -->
      <div class="flex border-b border-border bg-muted/20 shrink-0 text-xs">
        <button
          type="button"
          @click="activeTab = 'announcements'"
          class="flex-1 py-2.5 text-center font-semibold border-b-2 transition flex items-center justify-center gap-1.5 cursor-pointer"
          :class="
            activeTab === 'announcements'
              ? 'border-brand text-brand bg-card'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          "
        >
          <Megaphone class="w-3.5 h-3.5" />
          <span>Panneaux Ads ({{ visibleAnnouncements.length }})</span>
        </button>

        <button
          type="button"
          @click="activeTab = 'alerts'"
          class="flex-1 py-2.5 text-center font-semibold border-b-2 transition flex items-center justify-center gap-1.5 cursor-pointer"
          :class="
            activeTab === 'alerts'
              ? 'border-brand text-brand bg-card'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          "
        >
          <AlertTriangle class="w-3.5 h-3.5" />
          <span>Alertes Système ({{ visibleAlerts.length }})</span>
        </button>

        <button
          type="button"
          @click="activeTab = 'feedback'"
          class="flex-1 py-2.5 text-center font-semibold border-b-2 transition flex items-center justify-center gap-1.5 cursor-pointer"
          :class="
            activeTab === 'feedback'
              ? 'border-brand text-brand bg-card'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          "
        >
          <Star class="w-3.5 h-3.5 text-amber-500 fill-amber-500" />
          <span>Donner mon avis</span>
        </button>
      </div>

      <!-- Corps défilable -->
      <div class="p-4 overflow-y-auto space-y-4 flex-1 text-xs">
        <div
          v-if="isLoading"
          class="py-16 flex flex-col items-center justify-center text-muted-foreground gap-2"
        >
          <Loader2 class="w-6 h-6 animate-spin text-brand" />
          <span>Synchronisation des panneaux en direct...</span>
        </div>

        <!-- 1. Onglet Panneaux Ads -->
        <NetworkAnnouncementsTab
          v-else-if="activeTab === 'announcements'"
          :announcements="visibleAnnouncements"
          :dismissed-count="totalDismissedCount"
          @preview="openImagePreview"
          @dismiss="dismissPanel"
          @restore="restoreDismissedPanels"
        />

        <!-- 2. Onglet Alertes Système -->
        <NetworkAlertsTab
          v-else-if="activeTab === 'alerts'"
          :alerts="visibleAlerts"
          :dismissed-count="totalDismissedCount"
          @dismiss="dismissPanel"
          @restore="restoreDismissedPanels"
        />

        <!-- 3. Onglet Donner mon avis -->
        <NetworkFeedbackTab v-else />
      </div>

      <!-- Pied de modal avec réinitialisation et fermeture -->
      <div
        class="p-3 bg-muted/40 border-t border-border flex items-center justify-between text-xs"
      >
        <div>
          <button
            v-if="totalDismissedCount > 0 && activeTab !== 'feedback'"
            type="button"
            @click="restoreDismissedPanels"
            class="text-[11px] text-muted-foreground hover:text-foreground flex items-center gap-1 transition cursor-pointer"
          >
            <RotateCcw class="w-3 h-3" />
            <span
              >Réinitialiser les {{ totalDismissedCount }} annonces masquées</span
            >
          </button>
        </div>

        <button
          type="button"
          @click="emit('close')"
          class="px-4 py-1.5 text-xs font-semibold rounded-[3px] bg-muted hover:bg-muted/80 text-foreground transition cursor-pointer"
        >
          Fermer
        </button>
      </div>
    </div>

    <!-- Visionneuse Plein Écran (Lightbox Premium) -->
    <NetworkImageViewer
      :image="previewImage"
      :current-index="currentImageIndex"
      :total-images="imageAnnouncements.length"
      @close="closeImagePreview"
      @prev="goToPreviousImage"
      @next="goToNextImage"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useNotificationStore } from "../../stores/notificationStore";
import type { SystemNotificationItem } from "../../types";
import {
  Bell,
  X,
  Check,
  CheckCheck,
  Volume2,
  VolumeX,
  AlertTriangle,
  Info,
  ShieldAlert,
  CheckCircle2,
  Clock,
  Sparkles,
  ExternalLink,
  Loader2,
  Layers,
  Megaphone,
} from "lucide-vue-next";

defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "openNetworkModal"): void;
}>();

const notifStore = useNotificationStore();
const activeFilter = ref<"all" | "unread" | "alert">("all");
const isActivatingPush = ref(false);
const pushActivatedSuccess = ref(false);

const handleActivatePush = async () => {
  isActivatingPush.value = true;
  try {
    const granted = await notifStore.requestNativePermission();
    if (granted) {
      pushActivatedSuccess.value = true;
      setTimeout(() => {
        pushActivatedSuccess.value = false;
      }, 4000);
    }
  } finally {
    isActivatingPush.value = false;
  }
};

onMounted(() => {
  notifStore.loadNotifications();
});

const filteredNotifications = computed(() => {
  if (activeFilter.value === "unread") {
    return notifStore.notifications.filter((n) => !n.est_lu);
  }
  if (activeFilter.value === "alert") {
    return notifStore.notifications.filter(
      (n) => n.type === "critical" || n.type === "warning",
    );
  }
  return notifStore.notifications;
});

const formatDate = (iso: string) => {
  try {
    const d = new Date(iso);
    return `${d.toLocaleDateString("fr-FR")} à ${d.toLocaleTimeString("fr-FR", { hour: "2-digit", minute: "2-digit" })}`;
  } catch {
    return iso;
  }
};
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-40 flex items-center justify-end bg-black/40 backdrop-blur-xs select-none animate-in fade-in duration-150"
    @click.self="emit('close')"
  >
    <div
      class="bg-card border-l border-border h-full w-full max-w-md flex flex-col animate-in slide-in-from-right duration-200"
    >
      <!-- En-tête -->
      <div
        class="p-4 border-b border-border bg-muted/30 flex items-center justify-between shrink-0"
      >
        <div class="flex items-center space-x-2">
          <div
            class="w-8 h-8 rounded-[4px] bg-brand-light text-brand-dark flex items-center justify-center font-bold"
          >
            <Bell class="w-4 h-4 text-brand" />
          </div>
          <div>
            <div
              class="font-bold text-foreground text-sm flex items-center gap-2"
            >
              <span>Centre de Notifications</span>
              <span
                v-if="notifStore.unreadCount > 0"
                class="text-[10px] px-1.5 py-0.2 rounded-full bg-rose-600 text-white font-bold"
              >
                {{ notifStore.unreadCount }}
              </span>
            </div>
            <div class="text-[11px] text-muted-foreground">
              Flux d'événements et alertes temps réel
            </div>
          </div>
        </div>

        <div class="flex items-center gap-1">
          <!-- Toggle Son -->
          <button
            type="button"
            @click="notifStore.toggleSound"
            class="p-1.5 rounded-[4px] hover:bg-muted text-muted-foreground hover:text-foreground transition cursor-pointer"
            :title="
              notifStore.isSoundEnabled
                ? 'Couper le son des alertes'
                : 'Activer le carillon sonore'
            "
          >
            <Volume2
              v-if="notifStore.isSoundEnabled"
              class="w-4 h-4 text-brand"
            />
            <VolumeX v-else class="w-4 h-4 text-muted-foreground" />
          </button>

          <!-- Fermer -->
          <button
            type="button"
            @click="emit('close')"
            class="p-1.5 rounded-[4px] text-muted-foreground hover:text-foreground hover:bg-muted transition cursor-pointer"
          >
            <X class="w-4 h-4" />
          </button>
        </div>
      </div>

      <!-- Bandeau d'accès direct aux Panneaux Ads et Notation Réseau -->
      <div
        class="px-4 py-2 bg-gradient-to-r from-brand/10 via-brand/5 to-transparent border-b border-border flex items-center justify-between text-xs"
      >
        <button
          type="button"
          @click="emit('openNetworkModal')"
          class="flex items-center gap-1.5 text-brand font-bold text-[11px] hover:underline cursor-pointer"
        >
          <Megaphone class="w-3.5 h-3.5" />
          <span>Panneaux Ads Réseau & Noter le système</span>
        </button>
        <span
          class="text-[9px] uppercase font-mono px-1.5 py-0.5 rounded bg-amber-500/10 text-amber-600 dark:text-amber-400 font-semibold flex items-center gap-1"
        >
          <span>Avis</span>
          <span>★</span>
        </span>
      </div>

      <!-- Bannière d'autorisation des push OS si pas encore accordée -->
      <div
        v-if="!notifStore.isPermissionGranted"
        class="p-3 bg-brand-light border-b border-brand-border/60 text-xs flex items-center justify-between gap-2 transition-all"
      >
        <div class="text-[11px] text-brand-dark font-medium leading-tight">
          Recevez les alertes sur votre bureau même lorsque l'application est
          réduite.
        </div>
        <button
          type="button"
          :disabled="isActivatingPush"
          @click="handleActivatePush"
          class="shrink-0 px-2.5 py-1 rounded-[3px] bg-brand text-white font-bold text-[10px] shadow-xs hover:bg-brand-dark transition cursor-pointer flex items-center gap-1.5"
        >
          <Loader2 v-if="isActivatingPush" class="w-3 h-3 animate-spin" />
          <span>{{
            isActivatingPush ? "Activation..." : "Activer le push"
          }}</span>
        </button>
      </div>

      <!-- Confirmation immédiate de succès -->
      <div
        v-else-if="pushActivatedSuccess"
        class="p-2.5 bg-emerald-500/10 border-b border-emerald-500/30 text-xs flex items-center gap-2 text-emerald-700 dark:text-emerald-300 animate-in fade-in"
      >
        <CheckCircle2 class="w-4 h-4 shrink-0 text-emerald-600" />
        <span class="text-[11px] font-semibold"
          >Notifications bureau activées avec succès !</span
        >
      </div>

      <!-- Barre d'outils & Onglets -->
      <div
        class="px-4 py-2 border-b border-border bg-muted/10 flex items-center justify-between shrink-0 text-xs"
      >
        <div class="flex items-center gap-1">
          <button
            type="button"
            @click="activeFilter = 'all'"
            class="px-2 py-1 rounded-[3px] text-[11px] font-semibold transition cursor-pointer"
            :class="
              activeFilter === 'all'
                ? 'bg-brand text-white'
                : 'text-muted-foreground hover:text-foreground hover:bg-muted'
            "
          >
            Toutes ({{ notifStore.notifications.length }})
          </button>
          <button
            type="button"
            @click="activeFilter = 'unread'"
            class="px-2 py-1 rounded-[3px] text-[11px] font-semibold transition cursor-pointer"
            :class="
              activeFilter === 'unread'
                ? 'bg-brand text-white'
                : 'text-muted-foreground hover:text-foreground hover:bg-muted'
            "
          >
            Non lues ({{ notifStore.unreadCount }})
          </button>
          <button
            type="button"
            @click="activeFilter = 'alert'"
            class="px-2 py-1 rounded-[3px] text-[11px] font-semibold transition cursor-pointer"
            :class="
              activeFilter === 'alert'
                ? 'bg-brand text-white'
                : 'text-muted-foreground hover:text-foreground hover:bg-muted'
            "
          >
            Urgences
          </button>
        </div>

        <button
          v-if="notifStore.unreadCount > 0"
          type="button"
          @click="notifStore.markAllAsRead"
          class="text-[11px] font-medium text-brand hover:underline flex items-center gap-1 cursor-pointer"
        >
          <CheckCheck class="w-3.5 h-3.5" />
          <span>Tout lire</span>
        </button>
      </div>

      <!-- Liste défilable des notifications -->
      <div class="flex-1 overflow-y-auto p-3 space-y-2 text-xs">
        <div
          v-if="notifStore.isLoading && notifStore.notifications.length === 0"
          class="py-16 text-center text-muted-foreground flex flex-col items-center justify-center gap-2"
        >
          <Loader2 class="w-6 h-6 animate-spin text-brand" />
          <span>Chargement des notifications en direct...</span>
        </div>

        <div
          v-else-if="filteredNotifications.length === 0"
          class="py-16 text-center text-muted-foreground space-y-2"
        >
          <CheckCircle2 class="w-10 h-10 text-emerald-500 mx-auto opacity-70" />
          <div class="font-bold text-foreground">Aucune notification</div>
          <p class="text-[11px]">Vous êtes complètement à jour.</p>
        </div>

        <div
          v-for="n in filteredNotifications"
          :key="n.id"
          @click="notifStore.markAsRead(n.id)"
          class="p-3 rounded-[3px] border transition cursor-pointer space-y-1 relative group"
          :class="[
            n.est_lu
              ? 'border-border/60 bg-card text-muted-foreground opacity-80 hover:opacity-100'
              : 'border-border bg-card shadow-xs ring-1 ring-brand/10',
            n.type === 'critical' ? 'border-l-4 border-l-rose-500' : '',
            n.type === 'warning' ? 'border-l-4 border-l-amber-500' : '',
            n.type === 'success' ? 'border-l-4 border-l-emerald-500' : '',
            n.type === 'info' ? 'border-l-4 border-l-blue-500' : '',
          ]"
        >
          <!-- Pastille non lu -->
          <div
            v-if="!n.est_lu"
            class="absolute top-3 right-3 w-2 h-2 rounded-full bg-brand"
          ></div>

          <div class="flex items-center gap-2 flex-wrap pr-4">
            <!-- Badge Module -->
            <span
              class="text-[9px] uppercase font-bold px-1.5 py-0.2 rounded bg-muted text-foreground border border-border/60"
            >
              {{ n.module }}
            </span>

            <span class="text-[11px] font-bold text-foreground">
              {{ n.titre }}
            </span>
          </div>

          <p class="text-xs text-foreground/90 leading-relaxed">
            {{ n.message }}
          </p>

          <div
            class="flex items-center justify-between pt-1 border-t border-border/40 text-[10px] text-muted-foreground"
          >
            <span class="flex items-center gap-1">
              <Clock class="w-3 h-3" />
              <span>{{ formatDate(n.created_at) }}</span>
            </span>

            <button
              v-if="!n.est_lu"
              type="button"
              @click.stop="notifStore.markAsRead(n.id)"
              class="opacity-0 group-hover:opacity-100 hover:text-brand flex items-center gap-0.5 transition"
              title="Marquer comme lu"
            >
              <Check class="w-3 h-3" />
              <span>Marquer lu</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Pied de tiroir -->
      <div
        class="p-3 border-t border-border bg-muted/30 flex items-center justify-between text-[11px] text-muted-foreground shrink-0"
      >
        <div class="flex items-center gap-1.5">
          <span
            class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"
          ></span>
          <span>Connecté</span>
        </div>
        <button
          type="button"
          @click="emit('close')"
          class="font-semibold text-foreground hover:underline cursor-pointer"
        >
          Fermer
        </button>
      </div>
    </div>
  </div>
</template>

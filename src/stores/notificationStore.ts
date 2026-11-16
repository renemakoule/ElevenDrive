import { defineStore } from "pinia";
import { ref, computed } from "vue";
import type { SystemNotificationItem } from "../types";
import { notificationApi } from "../services/notificationService";
import { nativeNotificationEngine } from "../services/nativeNotificationService";

export const useNotificationStore = defineStore("notification", () => {
  const notifications = ref<SystemNotificationItem[]>([]);
  const unreadCount = ref(0);
  const isLoading = ref(false);
  const isSoundEnabled = ref(nativeNotificationEngine.isSoundEnabled());
  const isPermissionGranted = ref(
    nativeNotificationEngine.isPermissionGranted(),
  );
  let isWsConnected = false;
  let unsubscribeWs: (() => void) | null = null;

  const loadNotifications = async (unreadOnly = false) => {
    isLoading.value = true;
    try {
      const data = await notificationApi.getNotifications(unreadOnly, 50);
      notifications.value = data.notifications || [];
      unreadCount.value = data.unread_count || 0;
    } catch (err) {
      // ignore
    } finally {
      isLoading.value = false;
    }
  };

  const markAsRead = async (id: string) => {
    const target = notifications.value.find((n) => n.id === id);
    if (target && !target.est_lu) {
      target.est_lu = true;
      unreadCount.value = Math.max(0, unreadCount.value - 1);
      await notificationApi.markAsRead(id);
    }
  };

  const markAllAsRead = async () => {
    notifications.value.forEach((n) => {
      n.est_lu = true;
    });
    unreadCount.value = 0;
    await notificationApi.markAllAsRead();
  };

  const toggleSound = () => {
    const next = !isSoundEnabled.value;
    isSoundEnabled.value = next;
    nativeNotificationEngine.setSoundEnabled(next);
  };

  const requestNativePermission = async () => {
    const granted = await nativeNotificationEngine.requestPermission();
    isPermissionGranted.value = granted;
    return granted;
  };

  const initRealtime = (token: string) => {
    if (!token) return;

    nativeNotificationEngine.connect(token);
    isWsConnected = true;

    if (!unsubscribeWs) {
      unsubscribeWs = nativeNotificationEngine.subscribe((event) => {
        // 1. Prise en charge des notifications système
        if (event.type === "NOTIFICATION_NEW" && event.data) {
          const newNotif = event.data as SystemNotificationItem;
          // Éviter les doublons et actualiser réactivement
          const existingIdx = notifications.value.findIndex(
            (n) => n.id === newNotif.id,
          );
          if (existingIdx === -1) {
            notifications.value.unshift(newNotif);
            if (!newNotif.est_lu) {
              unreadCount.value++;
            }
          } else {
            const wasUnread = !notifications.value[existingIdx].est_lu;
            notifications.value[existingIdx] = newNotif;
            if (wasUnread && newNotif.est_lu) {
              unreadCount.value = Math.max(0, unreadCount.value - 1);
            }
          }
        }

        // 2. Orchestration en temps réel des directives Super-Admin (licences, statuts, modules)
        if (
          event.type.startsWith("SCHOOL_") ||
          (event.type === "NOTIFICATION_NEW" &&
            (event.data?.module === "abonnement" ||
              event.data?.module === "super_admin"))
        ) {
          import("./authStore")
            .then(({ useAuthStore }) => {
              const authStore = useAuthStore();
              if (event.type.startsWith("SCHOOL_")) {
                authStore.handleRealtimeSchoolEvent(event);
              } else {
                authStore.loadProfileByRole();
              }
            })
            .catch((err) => {
              console.warn("[NotificationStore] Sync AuthStore error:", err);
            });
        }
      });
    }
  };

  const cleanup = () => {
    if (unsubscribeWs) {
      unsubscribeWs();
      unsubscribeWs = null;
    }
    nativeNotificationEngine.disconnect();
    isWsConnected = false;
  };

  return {
    notifications,
    unreadCount,
    isLoading,
    isSoundEnabled,
    isPermissionGranted,
    loadNotifications,
    markAsRead,
    markAllAsRead,
    toggleSound,
    requestNativePermission,
    initRealtime,
    cleanup,
  };
});

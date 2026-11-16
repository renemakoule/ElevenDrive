import { ref, onMounted, onUnmounted, nextTick, type Ref } from "vue";
import { useAuthStore } from "../../stores/authStore";
import { useGuichetStore } from "../../stores/guichetStore";
import { authBackend } from "../../services/authBackendService";

export type ProviderId = "google" | "facebook" | "linkedin" | "apple";

export const providerNames: Record<ProviderId, string> = {
  google: "Google",
  facebook: "Facebook",
  linkedin: "LinkedIn",
  apple: "Apple",
};

export function useDesktopOAuthBridge(options: {
  isDesktopTauri: Ref<boolean>;
  onNeedsRegistration: (profile: {
    provider: ProviderId;
    email: string;
    nom: string;
    prenoms: string;
  }) => Promise<void>;
}) {
  const authStore = useAuthStore();
  const guichetStore = useGuichetStore();

  const isOAuthDesktopWaiting = ref(false);
  const oauthWaitingProvider = ref<ProviderId>("google");
  let oauthPollInterval: ReturnType<typeof setInterval> | null = null;

  const bringTauriWindowToFront = async () => {
    if (!options.isDesktopTauri.value) return;
    try {
      const { getCurrentWindow } = await import("@tauri-apps/api/window");
      const win = getCurrentWindow();
      await win.unminimize();
      await win.show();
      await win.setFocus();
    } catch {
      // Ignorer si non supporté
    }
  };

  const cancelOAuthDesktop = () => {
    if (oauthPollInterval) {
      clearInterval(oauthPollInterval);
      oauthPollInterval = null;
    }
    isOAuthDesktopWaiting.value = false;
  };

  const triggerOAuthFlow = async (provider: ProviderId) => {
    authStore.authError = null;
    oauthWaitingProvider.value = provider;
    isOAuthDesktopWaiting.value = true;

    const sessionId =
      "desk_" +
      Date.now().toString(36) +
      Math.random().toString(36).slice(2, 7);

    try {
      await authBackend.startOAuthDesktop(provider, sessionId);
    } catch (err: any) {
      cancelOAuthDesktop();
      guichetStore.showToast(
        typeof err === "string"
          ? err
          : err?.message || "Impossible d'ouvrir le navigateur.",
        "error",
      );
      return;
    }

    if (oauthPollInterval) clearInterval(oauthPollInterval);
    const startTime = Date.now();

    oauthPollInterval = setInterval(async () => {
      if (Date.now() - startTime > 180000) {
        cancelOAuthDesktop();
        guichetStore.showToast(
          "Délai d'authentification dépassé. Veuillez réessayer.",
          "error",
        );
        return;
      }

      try {
        const res = await authBackend.pollOAuthSession(sessionId);
        if (res.status === "authenticated" && res.token) {
          cancelOAuthDesktop();
          await bringTauriWindowToFront();
          await authStore.applyOAuthSession(
            res.token,
            res.refresh_token,
            res.role as any,
          );
          guichetStore.showToast(
            `Bienvenue, ${authStore.currentUser?.prenom || "Directeur"} ! Session active.`,
            "success",
          );
        } else if (res.status === "needs_registration") {
          cancelOAuthDesktop();
          await bringTauriWindowToFront();
          await options.onNeedsRegistration({
            provider: (res.provider as ProviderId) || provider,
            email: res.email || "",
            nom: res.nom || "",
            prenoms: res.prenoms || "",
          });
          guichetStore.showToast(
            `Identité certifiée (${res.email}). Finalisez le nom de votre établissement.`,
            "info",
          );
        } else if (res.status === "error") {
          cancelOAuthDesktop();
          await bringTauriWindowToFront();
          guichetStore.showToast(
            res.error || "Erreur d'authentification.",
            "error",
          );
        }
      } catch {
        // continuer polling
      }
    }, 1000);
  };

  const handleOAuthMessage = async (event: MessageEvent) => {
    if (event.data?.type !== "AUTOSCHOOL_OAUTH_RESULT") return;
    const result = event.data.payload;
    if (!result) return;

    if (result.status === "authenticated") {
      cancelOAuthDesktop();
      await bringTauriWindowToFront();
      if (result.token) {
        await authStore.applyOAuthSession(
          result.token,
          result.refresh_token,
          result.role,
        );
        guichetStore.showToast(
          `Bienvenue, ${authStore.currentUser?.prenom || result.user?.prenom || "Directeur"} ! Session active.`,
          "success",
        );
      }
    } else if (result.status === "needs_registration") {
      cancelOAuthDesktop();
      await bringTauriWindowToFront();
      await options.onNeedsRegistration({
        provider: (result.provider as ProviderId) || "google",
        email: result.email || "",
        nom: result.nom || "",
        prenoms: result.prenoms || "",
      });
      guichetStore.showToast(
        `Identité Google certifiée (${result.email}). Finalisez le nom de votre établissement.`,
        "info",
      );
    } else if (result.status === "error") {
      cancelOAuthDesktop();
      await bringTauriWindowToFront();
      guichetStore.showToast(
        result.error || "Erreur d'authentification avec le fournisseur.",
        "error",
      );
    }
  };

  onMounted(() => {
    window.addEventListener("message", handleOAuthMessage);
    if (authStore.pendingOAuthRegistration) {
      const pending = authStore.pendingOAuthRegistration;
      authStore.pendingOAuthRegistration = null;
      nextTick(() => {
        options.onNeedsRegistration({
          provider: (pending.provider as ProviderId) || "google",
          email: pending.email,
          nom: pending.nom,
          prenoms: pending.prenoms,
        });
      });
    }
  });

  onUnmounted(() => {
    window.removeEventListener("message", handleOAuthMessage);
    if (oauthPollInterval) {
      clearInterval(oauthPollInterval);
      oauthPollInterval = null;
    }
  });

  return {
    isOAuthDesktopWaiting,
    oauthWaitingProvider,
    cancelOAuthDesktop,
    triggerOAuthFlow,
  };
}

<script setup lang="ts">
import { ref, computed, nextTick, onMounted } from "vue";
import { useAuthStore } from "../../stores/authStore";
import { useGuichetStore } from "../../stores/guichetStore";
import type { RegisterSchoolPayload, LoginCredentials } from "../../types";
import appLogo from "@root/app-icon.png";
import {
  AlertTriangle,
  Building2,
  X,
  ArrowLeft,
  Terminal,
  Sparkles,
  Server,
  Wifi,
  WifiOff,
  RotateCw,
  Settings2,
} from "lucide-vue-next";
import { BACKEND_ORIGIN, API_BASE_URL } from "../../config/env";

import AuthBrandingAside from "./AuthBrandingAside.vue";
import AuthLoginForm from "./AuthLoginForm.vue";
import RegisterSchoolWizard from "./RegisterSchoolWizard.vue";
import RegisterSuccessView from "./RegisterSuccessView.vue";
import SchoolProvisioningView from "./SchoolProvisioningView.vue";
import OAuthIdentityModal from "./OAuthIdentityModal.vue";
import DesktopOAuthWaitingOverlay from "./DesktopOAuthWaitingOverlay.vue";
import QrCodeModal from "../common/QrCodeModal.vue";
import {
  useDesktopOAuthBridge,
  type ProviderId,
} from "./useDesktopOAuthBridge";

const emit = defineEmits<{
  (e: "close"): void;
  (e: "reopen-tour"): void;
}>();

const authStore = useAuthStore();
const guichetStore = useGuichetStore();

const activeTab = ref<"login" | "register">("login");
const isProvisioning = ref(false);
const currentRegisterStep = ref<1 | 2 | 3 | 4 | 5>(1);
const registerWizardRef = ref<InstanceType<typeof RegisterSchoolWizard> | null>(
  null,
);

const isOAuthModalOpen = ref(false);
const activeOAuthProvider = ref<ProviderId | null>(null);

const isQrModalOpen = ref(false);
const selectedQrLink = ref<{
  key: string;
  title: string;
  caption: string;
  shareTitle: string;
  url: string;
  dot: string;
} | null>(null);

const isDesktopTauri = computed(
  () =>
    typeof window !== "undefined" &&
    ("__TAURI_INTERNALS__" in window || "__TAURI__" in window),
);

// --- Détection & Diagnostic Connectivité Serveur Cloud ---
const serverStatus = ref<"checking" | "online" | "offline">("checking");
const activeBackendOrigin = computed(() => BACKEND_ORIGIN || "https://elevendrive.onrender.com");
const isConfigServerModalOpen = ref(false);
const customUrlInput = ref(
  (typeof localStorage !== "undefined" &&
    localStorage.getItem("elevendrive_custom_api_url")) ||
    "",
);

const checkServerConnection = async () => {
  serverStatus.value = "checking";
  try {
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), 4000);
    const res = await fetch(`${API_BASE_URL}/sync/ping`, {
      method: "GET",
      signal: controller.signal,
    }).catch(() => null);
    clearTimeout(timeout);

    if (res && res.ok) {
      serverStatus.value = "online";
    } else {
      const resRoot = await fetch(`${BACKEND_ORIGIN}/api`, {
        method: "GET",
      }).catch(() => null);
      serverStatus.value = resRoot && resRoot.ok ? "online" : "offline";
    }
  } catch {
    serverStatus.value = "offline";
  }
};

const saveCustomServerUrl = () => {
  const trimmed = customUrlInput.value.trim();
  if (trimmed) {
    localStorage.setItem("elevendrive_custom_api_url", trimmed);
  } else {
    localStorage.removeItem("elevendrive_custom_api_url");
  }
  isConfigServerModalOpen.value = false;
  window.location.reload();
};

onMounted(() => {
  checkServerConnection();
});

const {
  isOAuthDesktopWaiting,
  oauthWaitingProvider,
  cancelOAuthDesktop,
  triggerOAuthFlow,
} = useDesktopOAuthBridge({
  isDesktopTauri,
  onNeedsRegistration: async (profile) => {
    activeTab.value = "register";
    currentRegisterStep.value = 1;
    await nextTick();
    registerWizardRef.value?.applyOAuthProfile(profile);
  },
});

const accessLinks = computed(() => {
  const data = authStore.registrationSuccessData;
  if (!data) return [];
  const list = [
    {
      key: "moniteur",
      title: "Moniteurs",
      caption: "Livret de conduite mobile & émargement",
      shareTitle: "Portail Moniteurs · " + data.auto_ecole.nom_commercial,
      url: data.auto_ecole.links?.moniteur_link || "",
      dot: "bg-sky-500",
    },
    {
      key: "eleve",
      title: "Élèves & Candidats",
      caption: "Espace apprenant, solde & quiz 24/7",
      shareTitle: "Espace Apprenant · " + data.auto_ecole.nom_commercial,
      url: data.auto_ecole.links?.eleve_link || "",
      dot: "bg-purple-500",
    },
  ];

  if (data.auto_ecole.links?.inscription_link) {
    list.push({
      key: "inscription",
      title: "Pré-inscription en Ligne",
      caption: "Formulaire public pour nouveaux apprenants",
      shareTitle: "Inscription en Ligne · " + data.auto_ecole.nom_commercial,
      url: data.auto_ecole.links.inscription_link,
      dot: "bg-emerald-500",
    });
  }

  return list;
});

const clearError = () => {
  if (authStore.authError) {
    authStore.authError = null;
  }
};

const switchTab = (tab: "login" | "register") => {
  activeTab.value = tab;
  clearError();
  if (tab === "register") {
    currentRegisterStep.value = 1;
    registerWizardRef.value?.resetStep();
  }
};

const handleLogin = async (credentials: LoginCredentials) => {
  clearError();
  try {
    await authStore.login(credentials);
    guichetStore.showToast(
      `Bienvenue, ${authStore.currentUser.prenom} ! Session guichet active.`,
      "success",
    );
  } catch {
    // Erreur interceptée via authStore.authError
  }
};

const handleRegister = async (payload: RegisterSchoolPayload) => {
  clearError();
  try {
    await authStore.registerSchool(payload);
    isProvisioning.value = true;
    guichetStore.showToast(
      "Auto-École enregistrée. Lancement du traitement et déploiement...",
      "info",
    );
  } catch (err: any) {
    const errorMsg =
      authStore.authError ||
      err?.message ||
      "Erreur lors de l'enregistrement de l'auto-école.";
    guichetStore.showToast(errorMsg, "error");
  }
};

const handleFinishProvisioning = () => {
  isProvisioning.value = false;
  authStore.registrationSuccessData = null;
  guichetStore.showToast(
    "Bienvenue ! Période d'essai Super-Admin de 30 jours activée avec succès.",
    "success",
  );
};

const handleOAuthConfirm = async (data: {
  provider: ProviderId;
  email: string;
  nom: string;
  prenoms: string;
}) => {
  isOAuthModalOpen.value = false;
  clearError();
  try {
    const res = await authStore.loginWithOAuth(data.provider, {
      email: data.email,
      nom: data.nom,
      prenoms: data.prenoms,
    });

    if (res.status === "authenticated") {
      guichetStore.showToast(
        `Bienvenue, ${authStore.currentUser?.prenom || data.prenoms || "Directeur"} ! Session active.`,
        "success",
      );
    } else if (res.status === "needs_registration") {
      activeTab.value = "register";
      currentRegisterStep.value = 1;
      await nextTick();
      registerWizardRef.value?.applyOAuthProfile({
        provider: res.provider,
        email: res.profile?.email || data.email,
        nom: res.profile?.nom || data.nom,
        prenoms: res.profile?.prenoms || data.prenoms,
      });
      guichetStore.showToast(
        `Compte ${data.provider.toUpperCase()} vérifié ! Complétez les informations de votre établissement pour finaliser la création.`,
        "info",
      );
    }
  } catch {
    // Erreur interceptée via authStore.authError
  }
};

const openQrModal = (link: (typeof accessLinks.value)[0]) => {
  selectedQrLink.value = link;
  isQrModalOpen.value = true;
};

const finishRegistrationAndEnter = () => {
  authStore.registrationSuccessData = null;
};
</script>

<template>
  <div
    class="fixed inset-x-0 bottom-0 z-50 grid overflow-hidden bg-background lg:grid-cols-2"
    :class="isDesktopTauri ? 'top-9' : 'top-0'"
  >
    <AuthBrandingAside />

    <main class="relative flex flex-col overflow-y-auto">
      <button
        v-if="!isDesktopTauri"
        type="button"
        class="absolute right-5 top-5 z-20 flex items-center gap-1.5 rounded-[4px] border border-border bg-card/90 px-3 py-1.5 text-[11px] font-medium text-muted-foreground shadow-xs backdrop-blur-sm transition hover:bg-muted hover:text-foreground"
        title="Retour à l'accueil"
        @click="emit('close')"
      >
        <ArrowLeft class="h-3.5 w-3.5" />
        <span>Accueil</span>
      </button>

      <button
        v-else
        type="button"
        class="absolute right-5 top-5 z-20 flex items-center gap-1.5 rounded-full border border-border bg-card/80 px-3 py-1.5 text-[11px] font-medium text-muted-foreground shadow-xs backdrop-blur-sm transition hover:bg-muted hover:text-foreground hover:border-emerald-500/40"
        title="Revoir la visite guidée d'accueil"
        @click="emit('reopen-tour')"
      >
        <Sparkles class="h-3.5 w-3.5 text-emerald-500" />
        <span>Visite guidée</span>
      </button>

      <DesktopOAuthWaitingOverlay
        :is-waiting="isOAuthDesktopWaiting"
        :provider="oauthWaitingProvider"
        @cancel="cancelOAuthDesktop"
      />

      <div class="flex flex-1 items-center justify-center px-6 py-10 sm:px-10">
        <div
          class="w-full animate-in fade-in zoom-in-95 duration-200"
          :class="
            authStore.registrationSuccessData
              ? 'max-w-xl'
              : activeTab === 'register'
                ? 'max-w-md'
                : 'max-w-sm'
          "
        >
          <div class="mb-8 flex items-center justify-center gap-2.5 lg:hidden">
            <img
              :src="appLogo"
              alt="AutoSchool Logo"
              class="h-8 w-8 object-contain rounded-lg shadow-xs"
            />
            <span
              class="text-base font-semibold tracking-tight text-foreground"
            >
              ElevenDrive<span class="text-brand">.</span>
            </span>
          </div>

          <SchoolProvisioningView
            v-if="isProvisioning && authStore.registrationSuccessData"
            :session-data="authStore.registrationSuccessData"
            @finish="handleFinishProvisioning"
          />

          <RegisterSuccessView
            v-else-if="authStore.registrationSuccessData"
            :data="authStore.registrationSuccessData"
            :access-links="accessLinks"
            @open-qr="openQrModal"
            @finish="finishRegistrationAndEnter"
          />

          <template v-else>
            <div class="text-center">
              <h1 class="text-xl font-semibold tracking-tight text-foreground">
                {{
                  activeTab === "login" ? "Connexion" : "Créer votre auto-école"
                }}
              </h1>
              <p class="mt-2 text-xs text-muted-foreground">
                <template v-if="activeTab === 'login'">
                  Accédez à votre espace Auto-ecole.
                </template>
                <template v-else-if="currentRegisterStep === 1">
                  Étape 1 sur 5 : Démarrage rapide ou identité de
                  l'établissement.
                </template>
                <template v-else-if="currentRegisterStep === 2">
                  Étape 2 sur 5 : Localisation géographique et siège social.
                </template>
                <template v-else-if="currentRegisterStep === 3">
                  Étape 3 sur 5 : Coordonnées du Directeur / Administrateur.
                </template>
                <template v-else-if="currentRegisterStep === 4">
                  Étape 4 sur 5 : Sécurité du compte et code PIN managérial.
                </template>
                <template v-else>
                  Étape 5 sur 5 : Récapitulatif et validation finale de
                  l'auto-école.
                </template>
              </p>

              <div
                v-if="activeTab === 'login' && authStore.currentSchool"
                class="mt-3 inline-flex items-center gap-1.5 rounded-full border border-brand-border bg-brand-light/70 px-3 py-1 text-xs font-medium text-brand-dark"
              >
                <Building2 class="h-3.5 w-3.5 text-brand" />
                <span>{{ authStore.currentSchool.nom_commercial }}</span>
              </div>
            </div>

            <div
              v-if="authStore.authError"
              role="alert"
              class="mt-6 flex items-start gap-2.5 rounded-[5px] border border-rose-500/20 bg-rose-500/10 p-3 text-xs text-rose-600 dark:text-rose-400 animate-in fade-in duration-150"
            >
              <AlertTriangle class="mt-0.5 h-4 w-4 shrink-0" />
              <div class="flex-1">
                <span class="font-semibold">Erreur d'authentification</span>
                <p class="mt-0.5">{{ authStore.authError }}</p>
              </div>
              <button
                type="button"
                @click="clearError"
                class="rounded p-0.5 text-rose-500 hover:bg-rose-500/20"
                aria-label="Fermer"
              >
                <X class="h-4 w-4" />
              </button>
            </div>

            <AuthLoginForm
              v-if="activeTab === 'login'"
              :is-loading="authStore.isAuthLoading"
              @submit="handleLogin"
              @oauth="triggerOAuthFlow"
              @clear-error="clearError"
            />

            <RegisterSchoolWizard
              v-else
              ref="registerWizardRef"
              :is-loading="authStore.isAuthLoading"
              :error-message="authStore.authError"
              @submit="handleRegister"
              @oauth="triggerOAuthFlow"
              @clear-error="clearError"
              @step-change="currentRegisterStep = $event"
            />

            <p class="mt-8 text-center text-sm text-muted-foreground">
              <template v-if="activeTab === 'login'">
                Nouvel établissement ?
                <button
                  type="button"
                  @click="switchTab('register')"
                  class="font-medium text-brand underline-offset-4 transition hover:underline focus:outline-none focus-visible:underline"
                >
                  Enregistrer mon auto-école
                </button>
              </template>
              <template v-else>
                Déjà enregistré ?
                <button
                  type="button"
                  @click="switchTab('login')"
                  class="font-medium text-brand underline-offset-4 transition hover:underline focus:outline-none focus-visible:underline"
                >
                  Se connecter
                </button>
              </template>
            </p>

            <div class="mt-6 pt-4 border-t border-border/40 flex flex-col sm:flex-row items-center justify-between gap-2 text-xs">
              <!-- Indicateur Connectivité Serveur Cloud -->
              <div class="flex items-center gap-2">
                <div
                  class="flex items-center gap-1.5 px-2 py-0.5 rounded-[4px] border text-[11px] font-mono cursor-pointer transition hover:opacity-80"
                  :class="[
                    serverStatus === 'online'
                      ? 'border-emerald-500/30 bg-emerald-500/10 text-emerald-600 dark:text-emerald-400'
                      : serverStatus === 'offline'
                        ? 'border-rose-500/30 bg-rose-500/10 text-rose-600 dark:text-rose-400'
                        : 'border-border bg-muted text-muted-foreground',
                  ]"
                  @click="checkServerConnection"
                  :title="'Cliquez pour retester la connexion à ' + activeBackendOrigin"
                >
                  <span
                    class="h-1.5 w-1.5 rounded-full"
                    :class="[
                      serverStatus === 'online'
                        ? 'bg-emerald-500 animate-pulse'
                        : serverStatus === 'offline'
                          ? 'bg-rose-500'
                          : 'bg-amber-500 animate-ping',
                    ]"
                  ></span>
                  <span>{{ serverStatus === 'online' ? 'Cloud En Ligne' : serverStatus === 'offline' ? 'Cloud Injoignable' : 'Vérification...' }}</span>
                  <RotateCw class="w-3 h-3 ml-0.5 opacity-60 hover:opacity-100" />
                </div>

                <button
                  type="button"
                  @click="isConfigServerModalOpen = true"
                  class="p-1 text-muted-foreground hover:text-foreground transition rounded"
                  title="Configurer l'adresse du serveur API"
                >
                  <Settings2 class="w-3.5 h-3.5" />
                </button>
              </div>

              <button
                type="button"
                @click="authStore.webPortalView = 'super_admin_login'"
                class="inline-flex items-center gap-1.5 text-xs text-muted-foreground/70 hover:text-emerald-500 transition font-mono"
              >
                <Terminal class="w-3.5 h-3.5" />
                <span>Console Super-Admin</span>
              </button>
            </div>
          </template>
        </div>
      </div>
    </main>

    <OAuthIdentityModal
      :is-open="isOAuthModalOpen"
      :provider="activeOAuthProvider"
      :is-loading="authStore.isAuthLoading"
      @close="isOAuthModalOpen = false"
      @confirm="handleOAuthConfirm"
    />

    <QrCodeModal
      :is-open="isQrModalOpen"
      :title="selectedQrLink?.title || ''"
      :subtitle="selectedQrLink?.caption || ''"
      :url="selectedQrLink?.url || ''"
      :badge-text="selectedQrLink?.title || ''"
      :badge-dot="selectedQrLink?.dot || 'bg-brand'"
      :school-name="
        authStore.registrationSuccessData?.auto_ecole.nom_commercial || ''
      "
      @close="isQrModalOpen = false"
    />

    <!-- Modal Configuration Manuelle URL Serveur API -->
    <div
      v-if="isConfigServerModalOpen"
      class="fixed inset-0 z-[120] flex items-center justify-center bg-black/70 backdrop-blur-xs p-4 animate-in fade-in duration-150 select-none"
      @click.self="isConfigServerModalOpen = false"
    >
      <div
        class="bg-card border border-border rounded-[4px] max-w-md w-full p-5 space-y-4 shadow-xl text-foreground"
      >
        <div class="flex items-center justify-between pb-3 border-b border-border">
          <div class="flex items-center gap-2">
            <Server class="w-4 h-4 text-brand" />
            <h3 class="font-bold text-sm">Configuration Serveur Cloud</h3>
          </div>
          <button
            type="button"
            @click="isConfigServerModalOpen = false"
            class="p-1 rounded text-muted-foreground hover:text-foreground hover:bg-muted transition"
          >
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="space-y-3 text-xs">
          <p class="text-muted-foreground leading-relaxed">
            Par défaut, l'application se connecte au cluster Cloud Neon / Render hébergé.
            Si vous disposez d'une instance dédiée ou d'un serveur local, vous pouvez spécifier son URL ici.
          </p>

          <div>
            <label class="block font-medium mb-1 text-foreground">
              URL Serveur API Backend
            </label>
            <input
              v-model="customUrlInput"
              type="url"
              placeholder="https://elevendrive.onrender.com"
              class="w-full h-9 px-3 rounded-[3px] border border-border bg-background text-xs font-mono focus:border-brand focus:outline-none"
            />
            <span class="text-[10px] text-muted-foreground mt-1 block">
              Actuel : <span class="font-mono text-foreground">{{ activeBackendOrigin }}</span>
            </span>
          </div>

          <div class="flex items-center justify-between pt-2">
            <button
              type="button"
              @click="checkServerConnection"
              class="px-2.5 py-1.5 rounded-[3px] border border-border bg-muted hover:bg-muted/80 text-foreground text-[11px] font-semibold flex items-center gap-1.5 transition"
            >
              <RotateCw class="w-3 h-3" />
              <span>Tester la connexion</span>
            </button>

            <span
              v-if="serverStatus !== 'checking'"
              class="text-[11px] font-mono font-bold"
              :class="serverStatus === 'online' ? 'text-emerald-500' : 'text-rose-500'"
            >
              {{ serverStatus === 'online' ? '✓ Serveur Joignable' : '✗ Serveur Injoignable' }}
            </span>
          </div>
        </div>

        <div class="flex items-center justify-end gap-2 pt-3 border-t border-border">
          <button
            type="button"
            @click="isConfigServerModalOpen = false"
            class="px-3 py-1.5 rounded-[3px] bg-muted hover:bg-muted/80 text-foreground text-xs font-medium transition"
          >
            Annuler
          </button>
          <button
            type="button"
            @click="saveCustomServerUrl"
            class="px-3.5 py-1.5 rounded-[3px] bg-brand hover:bg-brand-dark text-white text-xs font-bold transition shadow-xs"
          >
            Enregistrer & Recharger
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

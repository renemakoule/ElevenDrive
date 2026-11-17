<script setup lang="ts">
import { ref, computed, nextTick } from "vue";
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
} from "lucide-vue-next";

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

            <div class="mt-6 pt-4 border-t border-border/40 text-center">
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
  </div>
</template>

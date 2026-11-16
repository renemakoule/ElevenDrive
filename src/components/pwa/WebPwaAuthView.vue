<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { useAuthStore } from "../../stores/authStore";
import { useGuichetStore } from "../../stores/guichetStore";
import { api } from "../../db/api";
import type { UserRole, LoginCredentials } from "../../types";
import {
  GraduationCap,
  Smartphone,
  ShieldCheck,
  Lock,
  Phone,
  Eye,
  EyeOff,
  ArrowRight,
  ArrowLeft,
  Loader2,
  AlertCircle,
  Building2,
  Sun,
  Moon,
  MessageCircle,
  Receipt,
} from "lucide-vue-next";
import OAuthIdentityModal from "../auth/OAuthIdentityModal.vue";

type SocialProvider = "google" | "facebook" | "linkedin" | "apple";

interface Props {
  initialRole?: UserRole;
}

const props = withDefaults(defineProps<Props>(), {
  initialRole: "eleve",
});

const emit = defineEmits<{
  (e: "navigate-home"): void;
  (e: "login-success"): void;
}>();

const authStore = useAuthStore();
const guichetStore = useGuichetStore();

const selectedRole = ref<UserRole>(props.initialRole || "eleve");
const identifiant = ref("");
const motDePasse = ref("");
const showPassword = ref(false);
const rememberMe = ref(true);
const isLoading = ref(false);
const socialLoading = ref<SocialProvider | null>(null);
const errorMessage = ref<string | null>(null);

// Connexion élève sans mot de passe (dossier / WhatsApp du guichet)
const isCandidateFastLogin = ref(false);

// Dernier provider utilisé (badge "Dernier")
const LAST_PROVIDER_KEY = "autoschool_last_provider";
const lastProvider = ref<SocialProvider | null>(null);

// Modal d'identité OAuth (Google, Facebook, LinkedIn, Apple)
const isOAuthModalOpen = ref(false);
const activeOAuthProvider = ref<SocialProvider | null>(null);

const handleOAuthMessage = async (event: MessageEvent) => {
  if (event.data?.type !== "AUTOSCHOOL_OAUTH_RESULT") return;
  const result = event.data.payload;
  if (!result) return;

  if (result.status === "authenticated" && result.token) {
    const session = await authStore.applyOAuthSession(
      result.token,
      result.refresh_token,
      result.role,
    );
    if (session && rememberMe.value && typeof localStorage !== "undefined") {
      localStorage.setItem(
        "autoschool_web_session",
        JSON.stringify({
          user: session.user,
          role: session.user.role,
          schoolSlug: authStore.schoolSlug,
          token: session.token,
        }),
      );
    }
    guichetStore.showToast(
      `Bienvenue, ${authStore.currentUser?.prenom || result.user?.prenom || "Élève"} ! Connexion réussie.`,
      "success",
    );
    emit("login-success");
  } else if (result.status === "needs_registration") {
    errorMessage.value = `Aucun compte élève associé à l'adresse "${result.email}". Rapprochez-vous de votre auto-école pour votre inscription.`;
  } else if (result.status === "error") {
    errorMessage.value = result.error || "Erreur d'authentification Google.";
  }
};

onMounted(() => {
  window.addEventListener("message", handleOAuthMessage);
  if (
    props.initialRole &&
    ["eleve", "moniteur", "manager"].includes(props.initialRole)
  ) {
    selectedRole.value = props.initialRole;
  }
  try {
    const stored = localStorage.getItem(
      LAST_PROVIDER_KEY,
    ) as SocialProvider | null;
    lastProvider.value = stored ?? "google";
  } catch {
    lastProvider.value = "google";
  }
});

onUnmounted(() => {
  window.removeEventListener("message", handleOAuthMessage);
});

const providers: { id: SocialProvider; label: string }[] = [
  { id: "google", label: "Google" },
  { id: "facebook", label: "Facebook" },
  { id: "linkedin", label: "LinkedIn" },
  { id: "apple", label: "Apple" },
];

const roleConfig = computed(() => {
  switch (selectedRole.value) {
    case "eleve":
      return {
        label: "Candidat & Élève",
        icon: GraduationCap,
        identifiantLabel: isCandidateFastLogin.value
          ? "N° WhatsApp ou numéro de dossier"
          : "Identifiant",
        identifiantPlaceholder: "Ex: +237 651 727 932 ou DOS-2024-001",
        hint: "Numéro WhatsApp ou numéro de dossier remis lors de votre inscription au guichet.",
      };
    case "moniteur":
      return {
        label: "Moniteur de conduite",
        icon: Smartphone,
        identifiantLabel: "Identifiant professionnel",
        identifiantPlaceholder: "Ex: +237 651 727 932 ou email",
        hint: "Mot de passe initial généré par la Direction dans la gestion des collaborateurs.",
      };
    case "secretaire":
      return {
        label: "Secrétariat (Mode Compagnon)",
        icon: Receipt,
        identifiantLabel: "Identifiant professionnel",
        identifiantPlaceholder: "Ex: +237 651 727 932 ou email",
        hint: "Accès compagnon mobile : dossiers élèves, émargement code en salle et plannings.",
      };
    case "manager":
      return {
        label: "Direction & Gérance",
        icon: ShieldCheck,
        identifiantLabel: "Identifiant professionnel",
        identifiantPlaceholder: "Ex: +237 651 727 932 ou email",
        hint: "Numéro WhatsApp ou adresse email utilisés lors de votre enregistrement.",
      };
    default:
      return {
        label: "Collaborateur",
        icon: GraduationCap,
        identifiantLabel: "Identifiant",
        identifiantPlaceholder: "Téléphone ou email",
        hint: "",
      };
  }
});

const roles: { id: UserRole; label: string; icon: any }[] = [
  { id: "eleve", label: "Élève", icon: GraduationCap },
  { id: "moniteur", label: "Moniteur", icon: Smartphone },
  { id: "secretaire", label: "Secrétaire", icon: Receipt },
  { id: "manager", label: "Direction", icon: ShieldCheck },
];

// Les providers sociaux ne concernent que l'espace élève
const showSocial = computed(() => selectedRole.value === "eleve");

const switchRole = (role: UserRole) => {
  selectedRole.value = role;
  errorMessage.value = null;
  if (role !== "eleve") {
    isCandidateFastLogin.value = false;
  }
};

const whatsappSupportUrl = computed(() => {
  const phone = (authStore.schoolPhone || "+237699001122").replace(
    /[^0-9]/g,
    "",
  );
  const msg = encodeURIComponent(
    `Bonjour, je n'arrive pas à me connecter à mon espace ${roleConfig.value.label} sur le portail ${authStore.schoolName}. Pouvez-vous m'assister ?`,
  );
  return `https://wa.me/${phone}?text=${msg}`;
});

/**
 * Déclenchement du flux OAuth pour l'élève (Google, Facebook, LinkedIn, Apple)
 */
const handleSocialLogin = (provider: SocialProvider) => {
  errorMessage.value = null;
  // Déclenchement direct du flux officiel OAuth 2.0 pour tous les providers
  authStore.startOAuthFlow(provider);
};

/**
 * Traitement de la confirmation de l'identité OAuth avec l'API Go
 */
const handleOAuthConfirm = async (data: {
  provider: SocialProvider;
  email: string;
  nom: string;
  prenoms: string;
}) => {
  isOAuthModalOpen.value = false;
  errorMessage.value = null;
  socialLoading.value = data.provider;

  try {
    localStorage.setItem(LAST_PROVIDER_KEY, data.provider);
  } catch {
    /* stockage indisponible : on ignore */
  }

  try {
    const res = await authStore.loginWithOAuth(data.provider, {
      email: data.email,
      nom: data.nom,
      prenoms: data.prenoms,
    });

    if (res.status === "authenticated" && res.session) {
      // Mémoriser la session web PWA pour la persistance locale sur mobile
      if (rememberMe.value && typeof localStorage !== "undefined") {
        localStorage.setItem(
          "autoschool_web_session",
          JSON.stringify({
            user: res.session.user,
            role: res.session.user.role,
            schoolSlug: authStore.schoolSlug,
            token: res.session.token,
          }),
        );
      }

      guichetStore.showToast(
        `Bienvenue, ${res.session.user.prenom || res.session.user.nom || ""} ! Connexion réussie.`,
        "success",
      );
      emit("login-success");
    } else {
      errorMessage.value = `Aucun compte élève associé à "${data.email}". Veuillez utiliser votre numéro WhatsApp ou vous rapprocher du guichet.`;
    }
  } catch (err: any) {
    errorMessage.value =
      err?.message ||
      `La connexion avec ${data.provider} a échoué. Réessayez ou utilisez vos accès.`;
  } finally {
    socialLoading.value = null;
  }
};

const handleLogin = async () => {
  errorMessage.value = null;
  const cleanId = identifiant.value.trim();

  if (!cleanId) {
    errorMessage.value = "Saisissez votre identifiant ou numéro de téléphone.";
    return;
  }

  if (!isCandidateFastLogin.value && !motDePasse.value) {
    errorMessage.value = "Saisissez votre mot de passe.";
    return;
  }

  isLoading.value = true;
  try {
    // CAS 1 : connexion standard via l'API Go (/api/login)
    if (!isCandidateFastLogin.value) {
      try {
        const credentials: LoginCredentials = {
          identifiant: cleanId,
          mot_de_passe: motDePasse.value,
        };
        const session = await authStore.login(
          credentials,
          authStore.schoolSlug,
        );

        if (rememberMe.value && typeof localStorage !== "undefined") {
          localStorage.setItem(
            "autoschool_web_session",
            JSON.stringify({
              user: session.user,
              role: session.user.role,
              schoolSlug: authStore.schoolSlug,
              token: session.token,
            }),
          );
        }

        guichetStore.showToast(
          `Bienvenue, ${session.user.prenom || session.user.nom} ! Connexion réussie.`,
          "success",
        );
        emit("login-success");
        return;
      } catch (loginErr: any) {
        if (selectedRole.value !== "eleve") {
          throw loginErr;
        }
      }
    }

    // CAS 2 : élève via dossier candidat ou numéro WhatsApp du guichet
    if (selectedRole.value === "eleve") {
      const candidates = await api.searchCandidats(cleanId);
      const cleanPhone = cleanId.replace(/[^0-9]/g, "");

      const foundCandidate =
        candidates.find(
          (c) =>
            (cleanPhone &&
              (c.telephone || "").replace(/[^0-9]/g, "") === cleanPhone) ||
            c.nom.toLowerCase().includes(cleanId.toLowerCase()) ||
            c.prenoms.toLowerCase().includes(cleanId.toLowerCase()) ||
            c.id.toLowerCase() === cleanId.toLowerCase(),
        ) || (candidates.length === 1 ? candidates[0] : null);

      if (foundCandidate) {
        authStore.loginWebAsCandidate({
          id: foundCandidate.id,
          nom: foundCandidate.nom,
          prenom: foundCandidate.prenoms,
          telephone: foundCandidate.telephone || cleanPhone || "",
          dossier: foundCandidate.id,
        });

        guichetStore.showToast(
          `Bienvenue ${foundCandidate.prenoms} ! Accès à votre livret d'apprentissage.`,
          "success",
        );
        emit("login-success");
        return;
      }
      throw new Error(
        `Aucun dossier trouvé pour « ${cleanId} » dans cette auto-école. Vérifiez votre numéro ou contactez le guichet.`,
      );
    }

    throw new Error("Identifiant ou mot de passe incorrect.");
  } catch (err: any) {
    errorMessage.value =
      err?.message || "Connexion impossible. Vérifiez vos identifiants.";
  } finally {
    isLoading.value = false;
  }
};
</script>

<template>
  <div
    class="min-h-screen w-full grid lg:grid-cols-2 bg-background text-foreground selection:bg-brand selection:text-white"
  >
    <!-- ============ PANNEAU GAUCHE : IDENTITÉ ============ -->
    <aside
      class="relative hidden lg:flex flex-col border-r border-border overflow-hidden bg-background"
    >
      <!-- Trame de points, estompée vers les bords -->
      <div
        class="login-dots pointer-events-none absolute inset-0"
        aria-hidden="true"
      ></div>

      <!-- Logo -->
      <div class="relative z-10 flex items-center gap-3 px-10 pt-9">
        <svg viewBox="0 0 32 32" class="w-8 h-8" aria-hidden="true">
          <defs>
            <linearGradient id="ed-g" x1="0" y1="0" x2="1" y2="1">
              <stop offset="0" stop-color="#a855f7" />
              <stop offset="1" stop-color="#38bdf8" />
            </linearGradient>
          </defs>
          <path
            d="M16 2 28 9v14L16 30 4 23V9z"
            fill="none"
            stroke="url(#ed-g)"
            stroke-width="1.6"
          />
          <path
            d="M16 2v28M4 9l24 14M28 9 4 23"
            stroke="url(#ed-g)"
            stroke-width="0.9"
            opacity="0.7"
          />
          <circle cx="16" cy="16" r="3.2" fill="url(#ed-g)" />
        </svg>
        <span class="text-[15px] font-semibold tracking-tight"
          >ElevenDrive<span class="text-brand">.</span></span
        >
      </div>

      <!-- Message central -->
      <div
        class="relative z-10 flex-1 flex flex-col items-center justify-center px-12 text-center"
      >
        <div
          class="w-[52px] h-[52px] border border-border bg-background flex items-center justify-center mb-24"
        >
          <svg viewBox="0 0 32 32" class="w-8 h-8" aria-hidden="true">
            <path
              d="M16 2 28 9v14L16 30 4 23V9z"
              fill="none"
              stroke="url(#ed-g)"
              stroke-width="1.6"
            />
            <path
              d="M16 2v28M4 9l24 14M28 9 4 23"
              stroke="url(#ed-g)"
              stroke-width="0.9"
              opacity="0.7"
            />
            <circle cx="16" cy="16" r="3.2" fill="url(#ed-g)" />
          </svg>
        </div>

        <h1
          class="max-w-md text-[34px] leading-[1.15] font-bold tracking-tight text-balance"
        >
          Pilotez votre auto-école, du guichet à la route.
        </h1>
        <p class="mt-4 text-sm text-muted-foreground">
          Plateforme intégrée pour les auto-écoles en Afrique.
        </p>
      </div>
    </aside>

    <!-- ============ PANNEAU DROIT : FORMULAIRE ============ -->
    <main
      class="relative flex flex-col items-center justify-center px-5 py-10 sm:px-10"
    >
      <!-- Barre du haut -->
      <div
        class="absolute top-4 left-4 right-4 flex items-center justify-between"
      >
        <button
          type="button"
          @click="emit('navigate-home')"
          class="inline-flex items-center gap-1.5 h-8 px-3 rounded-[4px] border border-border bg-card text-xs font-medium hover:bg-muted transition focus-visible:outline focus-visible:outline-2 focus-visible:outline-brand"
        >
          <ArrowLeft class="w-3.5 h-3.5" />
          <span>Retour à l'accueil</span>
        </button>

        <button
          type="button"
          @click="authStore.toggleDarkMode()"
          class="h-8 w-8 inline-flex items-center justify-center rounded-[4px] border border-border bg-card hover:bg-muted transition focus-visible:outline focus-visible:outline-2 focus-visible:outline-brand"
          aria-label="Changer de thème"
        >
          <Sun v-if="authStore.isDark" class="w-4 h-4 text-amber-400" />
          <Moon v-else class="w-4 h-4 text-slate-600" />
        </button>
      </div>

      <div class="w-full max-w-[385px]">
        <!-- En-tête -->
        <div class="text-center">
          <h2 class="text-[22px] font-bold tracking-tight">Connexion</h2>
          <p class="mt-1.5 text-xs text-muted-foreground">
            Accédez à votre espace Auto-école.
          </p>

          <div
            class="mt-4 inline-flex items-center gap-1.5 h-7 px-3.5 rounded-full text-xs font-semibold border border-brand/40 bg-brand/5 text-brand"
          >
            <Building2 class="w-3.5 h-3.5" />
            <span>{{ authStore.schoolName }}</span>
          </div>
        </div>

        <!-- Sélecteur de rôle -->
        <div
          class="mt-6 grid grid-cols-2 sm:grid-cols-4 gap-1 p-1 rounded-[4px] bg-muted/60 border border-border"
          role="tablist"
        >
          <button
            v-for="r in roles"
            :key="r.id"
            type="button"
            role="tab"
            :aria-selected="selectedRole === r.id"
            @click="switchRole(r.id)"
            class="flex items-center justify-center gap-1.5 h-8 rounded text-xs transition focus-visible:outline focus-visible:outline-2 focus-visible:outline-brand"
            :class="
              selectedRole === r.id
                ? 'bg-card text-foreground font-semibold border border-border'
                : 'text-muted-foreground hover:text-foreground border border-transparent'
            "
          >
            <component :is="r.icon" class="w-3.5 h-3.5" />
            <span>{{ r.label }}</span>
          </button>
        </div>

        <!-- Providers sociaux : espace élève uniquement -->
        <template v-if="showSocial">
          <div class="mt-5 grid grid-cols-2 gap-3">
            <button
              v-for="p in providers"
              :key="p.id"
              type="button"
              :disabled="!!socialLoading"
              @click="handleSocialLogin(p.id)"
              class="relative h-10 inline-flex items-center justify-center gap-2.5 rounded-[4px] border border-border bg-card text-[13px] font-semibold hover:bg-muted transition disabled:opacity-60 focus-visible:outline focus-visible:outline-2 focus-visible:outline-brand"
            >
              <span
                v-if="lastProvider === p.id"
                class="absolute -top-2 right-3 px-1.5 py-px rounded-[4px] bg-brand text-white text-[9px] font-semibold leading-tight"
              >
                Dernier
              </span>

              <Loader2
                v-if="socialLoading === p.id"
                class="w-4 h-4 animate-spin"
              />

              <!-- Google -->
              <svg
                v-else-if="p.id === 'google'"
                viewBox="0 0 24 24"
                class="w-[18px] h-[18px]"
                aria-hidden="true"
              >
                <path
                  fill="#EA4335"
                  d="M12 10.2v3.9h5.5c-.24 1.4-1.7 4.1-5.5 4.1a6.2 6.2 0 0 1 0-12.4c1.96 0 3.27.84 4.02 1.56l2.74-2.64C17 3.1 14.7 2 12 2a10 10 0 1 0 0 20c5.77 0 9.6-4.06 9.6-9.76 0-.66-.07-1.16-.16-1.64z"
                />
                <path
                  fill="#4285F4"
                  d="M21.6 12.24c0-.66-.07-1.16-.16-1.64H12v3.9h5.5a4.7 4.7 0 0 1-2 3.1l3.1 2.4c1.8-1.7 3-4.2 3-7.76z"
                />
                <path
                  fill="#FBBC05"
                  d="M5.8 14.3a6.2 6.2 0 0 1 0-4.6L2.6 7.2a10 10 0 0 0 0 9.6z"
                />
                <path
                  fill="#34A853"
                  d="M12 22c2.7 0 4.96-.9 6.6-2.4l-3.1-2.4c-.86.58-2 .98-3.5.98a6.2 6.2 0 0 1-5.8-4.3l-3.2 2.5A10 10 0 0 0 12 22z"
                />
              </svg>

              <!-- Facebook -->
              <svg
                v-else-if="p.id === 'facebook'"
                viewBox="0 0 24 24"
                class="w-[18px] h-[18px]"
                aria-hidden="true"
              >
                <circle cx="12" cy="12" r="11" fill="#1877F2" />
                <path
                  fill="#fff"
                  d="M13.4 22v-8h2.7l.4-3.2h-3.1V8.8c0-.9.26-1.5 1.6-1.5h1.6V4.4c-.3 0-1.3-.1-2.4-.1-2.4 0-4 1.4-4 4.1v2.4H7.5V14h2.7v8z"
                />
              </svg>

              <!-- LinkedIn -->
              <svg
                v-else-if="p.id === 'linkedin'"
                viewBox="0 0 24 24"
                class="w-[18px] h-[18px]"
                aria-hidden="true"
              >
                <rect
                  x="2"
                  y="2"
                  width="20"
                  height="20"
                  rx="3"
                  fill="#0A66C2"
                />
                <path
                  fill="#fff"
                  d="M7.3 10h2.4v7.7H7.3zM8.5 6.3a1.4 1.4 0 1 1 0 2.8 1.4 1.4 0 0 1 0-2.8zM11.2 10h2.3v1.05c.34-.64 1.18-1.3 2.4-1.3 2.5 0 2.95 1.65 2.95 3.8v4.15h-2.4v-3.68c0-.88-.02-2-1.22-2-1.22 0-1.4.95-1.4 1.94v3.74h-2.4z"
                />
              </svg>

              <!-- Apple -->
              <svg
                v-else
                viewBox="0 0 24 24"
                class="w-[18px] h-[18px] fill-current"
                aria-hidden="true"
              >
                <path
                  d="M16.37 12.6c-.02-2.2 1.8-3.25 1.88-3.3-1.02-1.5-2.62-1.7-3.19-1.73-1.35-.14-2.65.8-3.34.8-.7 0-1.76-.78-2.9-.76-1.5.02-2.87.87-3.64 2.2-1.55 2.7-.4 6.7 1.12 8.9.74 1.07 1.62 2.28 2.78 2.24 1.12-.04 1.54-.72 2.9-.72 1.35 0 1.73.72 2.9.7 1.2-.02 1.96-1.1 2.7-2.17.85-1.24 1.2-2.44 1.22-2.5-.03-.01-2.34-.9-2.4-3.56zM14.2 6.1c.6-.74 1.02-1.76.9-2.78-.88.04-1.94.58-2.57 1.32-.56.65-1.05 1.7-.92 2.7.98.08 1.98-.5 2.6-1.24z"
                />
              </svg>

              <span>{{ p.label }}</span>
            </button>
          </div>

          <!-- Séparateur -->
          <div
            class="mt-5 flex items-center gap-4 text-[11px] font-medium text-muted-foreground"
          >
            <span class="flex-1 h-px bg-border"></span>
            <span>OU AVEC VOS ACCÈS</span>
            <span class="flex-1 h-px bg-border"></span>
          </div>
        </template>

        <!-- Formulaire -->
        <form @submit.prevent="handleLogin" class="mt-5 space-y-4" novalidate>
          <div
            v-if="errorMessage"
            role="alert"
            class="flex items-start gap-2.5 p-3 rounded-[4px] bg-rose-500/10 border border-rose-500/20 text-rose-600 dark:text-rose-400 text-xs leading-relaxed"
          >
            <AlertCircle class="w-4 h-4 shrink-0 mt-px" />
            <span>{{ errorMessage }}</span>
          </div>

          <!-- Identifiant -->
          <div>
            <label for="login-id" class="block text-xs font-semibold mb-1.5">
              {{ roleConfig.identifiantLabel }}
            </label>
            <div class="relative flex items-center">
              <Phone class="absolute left-3 w-4 h-4 text-muted-foreground" />
              <input
                id="login-id"
                v-model="identifiant"
                type="text"
                autocomplete="username"
                :placeholder="roleConfig.identifiantPlaceholder"
                class="w-full h-10 pl-9 pr-3 rounded-[4px] border border-border bg-card text-[13px] placeholder:text-muted-foreground/70 focus:border-brand focus:ring-2 focus:ring-brand/20 focus:outline-none transition"
              />
            </div>
            <p
              v-if="roleConfig.hint"
              class="mt-1.5 text-[10.5px] text-muted-foreground leading-snug"
            >
              {{ roleConfig.hint }}
            </p>
          </div>

          <!-- Mot de passe -->
          <div v-if="!isCandidateFastLogin">
            <label for="login-pwd" class="block text-xs font-semibold mb-1.5"
              >Mot de passe</label
            >
            <div class="relative flex items-center">
              <Lock class="absolute left-3 w-4 h-4 text-muted-foreground" />
              <input
                id="login-pwd"
                v-model="motDePasse"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="current-password"
                placeholder="Votre mot de passe"
                class="w-full h-10 pl-9 pr-10 rounded-[4px] border border-border bg-card text-[13px] placeholder:text-muted-foreground/70 focus:border-brand focus:ring-2 focus:ring-brand/20 focus:outline-none transition"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute right-3 p-1 text-muted-foreground hover:text-foreground transition"
                :aria-label="
                  showPassword
                    ? 'Masquer le mot de passe'
                    : 'Afficher le mot de passe'
                "
              >
                <EyeOff v-if="showPassword" class="w-4 h-4" />
                <Eye v-else class="w-4 h-4" />
              </button>
            </div>
          </div>

          <!-- Alternative élève sans mot de passe -->
          <button
            v-if="selectedRole === 'eleve'"
            type="button"
            @click="isCandidateFastLogin = !isCandidateFastLogin"
            class="text-xs text-brand hover:underline font-medium"
          >
            {{
              isCandidateFastLogin
                ? "Se connecter avec un mot de passe"
                : "Inscription récente ? Connexion par n° WhatsApp ou dossier"
            }}
          </button>

          <!-- Session -->
          <label
            class="flex items-center gap-2 text-xs text-muted-foreground cursor-pointer select-none"
          >
            <input
              v-model="rememberMe"
              type="checkbox"
              class="w-4 h-4 rounded border-border text-brand focus:ring-brand"
            />
            Garder ma session active
          </label>

          <!-- Bouton -->
          <button
            type="submit"
            :disabled="isLoading"
            class="w-full h-11 rounded-[4px] bg-brand hover:bg-brand-dark text-white font-semibold text-sm transition flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-brand"
          >
            <Loader2 v-if="isLoading" class="w-4 h-4 animate-spin" />
            <span>{{
              isLoading ? "Connexion en cours…" : "Se connecter"
            }}</span>
            <ArrowRight v-if="!isLoading" class="w-4 h-4" />
          </button>
        </form>
      </div>
    </main>

    <!-- Modal d'identification OAuth (Google, Facebook, LinkedIn, Apple) -->
    <OAuthIdentityModal
      :is-open="isOAuthModalOpen"
      :provider="activeOAuthProvider"
      :is-loading="!!socialLoading"
      @close="isOAuthModalOpen = false"
      @confirm="handleOAuthConfirm"
    />
  </div>
</template>

<style scoped>
/* Trame de points qui s'estompe vers les bords (effet de la maquette) */
.login-dots {
  background-image: radial-gradient(currentColor 1px, transparent 1.4px);
  background-size: 14px 14px;
  color: hsl(var(--foreground, 0 0% 10%) / 0.55);
  -webkit-mask-image: radial-gradient(
    ellipse 75% 70% at 50% 50%,
    #000 15%,
    transparent 78%
  );
  mask-image: radial-gradient(
    ellipse 75% 70% at 50% 50%,
    #000 15%,
    transparent 78%
  );
}

@media (prefers-reduced-motion: reduce) {
  .animate-spin {
    animation-duration: 2s;
  }
}
</style>

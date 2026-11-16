<script setup lang="ts">
import { ref, computed } from "vue";
import { useAuthStore } from "../../stores/authStore";
import {
  Lock,
  ShieldAlert,
  Sparkles,
  ArrowRight,
  LogOut,
  Building2,
  AlertTriangle,
  KeyRound,
  FileCode,
  RefreshCw,
} from "lucide-vue-next";

const emit = defineEmits<{
  (e: "openPricing"): void;
}>();

const authStore = useAuthStore();
const isChecking = ref(false);

const isDirecteur = computed(() => authStore.currentRole === "manager");
const isSuspended = computed(
  () => authStore.currentSchool?.est_actif === false,
);
const isTrial = computed(() => {
  const t = (authStore.currentSchool?.type_licence || "TRIAL").toUpperCase();
  return t.includes("TRIAL") || t.includes("ESSAI") || t.includes("GRATUIT");
});

const checkStatusNow = async () => {
  isChecking.value = true;
  try {
    await authStore.loadProfileByRole();
  } finally {
    isChecking.value = false;
  }
};

const expirationDateFormatted = computed(() => {
  const raw = authStore.currentSchool?.date_expiration_licence;
  if (!raw) return "Échue";
  try {
    return new Date(raw).toLocaleDateString("fr-FR", {
      day: "2-digit",
      month: "long",
      year: "numeric",
    });
  } catch {
    return raw;
  }
});
</script>

<template>
  <div
    class="min-h-[75vh] w-full flex items-center justify-center p-4 select-none animate-in fade-in duration-300"
  >
    <div
      class="max-w-2xl w-full rounded-[3px] border border-rose-500/30 bg-gradient-to-b from-[#0f1118] via-[#090b10] to-[#06080c] text-white p-6 sm:p-10 relative overflow-hidden text-center space-y-6"
    >
      <!-- Halo de cryptage d'arrière-plan -->
      <div
        class="absolute -top-24 left-1/2 -translate-x-1/2 w-96 h-96 bg-rose-600/10 rounded-full blur-3xl pointer-events-none"
      ></div>

      <!-- Badge de verrouillage système -->
      <div class="relative z-10 flex flex-col items-center gap-3">
        <div
          class="h-16 w-16 rounded-[3px] bg-rose-500/10 border-2 border-rose-500/40 flex items-center justify-center shadow-lg shadow-rose-950/50"
        >
          <Lock class="h-8 w-8 text-rose-500 animate-pulse" />
        </div>

        <div
          class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-rose-500/15 border border-rose-500/30 text-[11px] font-mono font-bold uppercase tracking-wider text-rose-400"
        >
          <ShieldAlert class="w-3.5 h-3.5" />
          <span v-if="isSuspended"
            >Accès Suspendu par le Super-Administrateur</span
          >
          <span v-else-if="isTrial"
            >Accès Crypté · Période d'Essai 30 Jours Épuisée</span
          >
          <span v-else
            >Accès Verrouillé · Abonnement
            {{ authStore.currentSchool?.type_licence }} Expiré</span
          >
        </div>
      </div>

      <!-- Titre et explications -->
      <div class="relative z-10 space-y-2">
        <h2 class="text-xl sm:text-2xl font-black tracking-tight text-white">
          <template v-if="isSuspended">
            Établissement temporairement suspendu ({{ authStore.schoolName }})
          </template>
          <template v-else-if="isTrial">
            Période d'essai expirée pour {{ authStore.schoolName }}
          </template>
          <template v-else>
            Licence expirée pour {{ authStore.schoolName }}
          </template>
        </h2>
        <p
          class="text-xs sm:text-sm text-white/70 max-w-lg mx-auto leading-relaxed"
        >
          <template v-if="isSuspended">
            L'accès à votre établissement a été suspendu depuis la console
            centrale Super-Admin. Vos données et dossiers restent sauvegardés et
            intacts.
          </template>
          <template v-else-if="isTrial">
            Votre établissement a bénéficié de
            <strong class="text-amber-400 font-bold"
              >30 jours d'essai gratuit offerts par ElevenDrive</strong
            >
            (échéance dépassée le {{ expirationDateFormatted }}). Les données et
            dossiers sont entièrement sauvegardés et sécurisés.
          </template>
          <template v-else>
            Votre formule
            <strong class="text-amber-400 font-bold">{{
              authStore.currentSchool?.type_licence
            }}</strong>
            est arrivée à échéance le {{ expirationDateFormatted }}. Renouvelez
            votre licence pour continuer l'exploitation.
          </template>
        </p>
      </div>

      <!-- Bloc d'empreinte de sécurité / Cryptage -->
      <div
        class="relative z-10 p-3.5 rounded-lg bg-black/60 border border-white/10 font-mono text-[11px] text-white/50 text-left space-y-1.5"
      >
        <div
          class="flex items-center justify-between text-white/40 text-[10px] border-b border-white/10 pb-1"
        >
          <span class="flex items-center gap-1">
            <KeyRound class="w-3 h-3 text-rose-400" />
            <span>ÉTAT DU VERROUILLAGE SYSTÈME</span>
          </span>
          <span class="text-rose-400">STATUT : CRYPTÉ</span>
        </div>
        <div class="flex items-center justify-between">
          <span>Établissement :</span>
          <strong class="text-white"
            >{{ authStore.schoolName }} ({{ authStore.schoolSlug }})</strong
          >
        </div>
        <div class="flex items-center justify-between">
          <span>Utilisateur connecté :</span>
          <strong class="text-white"
            >{{ authStore.currentUser.prenom }}
            {{ authStore.currentUser.nom }} ({{
              authStore.rolePermissions.roleTitle
            }})</strong
          >
        </div>
        <div class="flex items-center justify-between">
          <span>Licence requise :</span>
          <span class="text-amber-400 font-bold"
            >Abonnement Eleven Pro ou Custom</span
          >
        </div>
      </div>

      <!-- Action spécifique selon le Rôle (Directeur vs Secrétaire) -->
      <div class="relative z-10 pt-2">
        <!-- 1. CAS DU DIRECTEUR (manager) : Affichage direct du Bouton "Bon Plan" -->
        <div v-if="isDirecteur" class="space-y-4">
          <div
            class="p-4 rounded-[4px] bg-gradient-to-r from-amber-500/10 via-brand/10 to-amber-500/10 border border-amber-500/30 space-y-2"
          >
            <div
              class="flex items-center justify-center gap-2 text-amber-400 font-bold text-xs uppercase tracking-wide"
            >
              <Sparkles class="w-4 h-4 text-amber-400 animate-spin" />
              <span>Offre Spéciale de Déblocage Immédiat</span>
            </div>
            <p class="text-xs text-white/80 max-w-md mx-auto">
              En tant que Directeur, vous pouvez renouveler et réactiver
              instantanément l'ensemble des modules.
            </p>
          </div>

          <!-- Bouton BON PLAN Majeur -->
          <button
            type="button"
            @click="emit('openPricing')"
            class="w-full sm:w-auto px-8 py-3.5 rounded-[4px] bg-gradient-to-r from-amber-500 via-amber-400 to-amber-500 text-black font-black text-sm uppercase tracking-wider shadow-2xl hover:scale-105 transition-transform duration-200 cursor-pointer inline-flex items-center justify-center gap-2.5 animate-pulse"
          >
            <Sparkles class="w-5 h-5 text-black" />
            <span>Bon Plan · Débloquer l'Auto-École</span>
            <ArrowRight class="w-4 h-4 text-black" />
          </button>
        </div>

        <!-- 2. CAS DU SECRÉTAIRE : Message d'invitation à contacter le Directeur -->
        <div v-else class="space-y-3">
          <div
            class="p-4 rounded-[4px] bg-amber-500/10 border border-amber-500/30 text-amber-300 text-xs text-center space-y-1"
          >
            <div class="font-bold flex items-center justify-center gap-1.5">
              <AlertTriangle class="w-4 h-4 text-amber-400" />
              <span>Privilèges réservés à la Direction Générale</span>
            </div>
            <p
              class="text-[11px] text-white/70 max-w-md mx-auto leading-relaxed"
            >
              L'accès au guichet est suspendu suite à l'épuisement des 30 jours
              d'essai. Veuillez inviter le
              <strong>Directeur de l'auto-école</strong> à se connecter pour
              souscrire au plan supérieur via le bouton "Bon Plan".
            </p>
          </div>
        </div>
      </div>

      <!-- Actualisation & Déconnexion sécurisée -->
      <div
        class="relative z-10 pt-3 border-t border-white/10 flex flex-wrap items-center justify-center gap-4"
      >
        <button
          type="button"
          @click="checkStatusNow"
          :disabled="isChecking"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-[4px] bg-white/10 hover:bg-white/15 border border-white/15 text-xs text-white transition cursor-pointer disabled:opacity-50"
        >
          <RefreshCw
            class="w-3.5 h-3.5"
            :class="isChecking ? 'animate-spin text-amber-400' : ''"
          />
          <span>{{
            isChecking ? "Synchronisation..." : "Vérifier le déblocage"
          }}</span>
        </button>

        <button
          type="button"
          @click="authStore.logout()"
          class="inline-flex items-center gap-1.5 text-xs text-white/50 hover:text-white transition cursor-pointer"
        >
          <LogOut class="w-3.5 h-3.5" />
          <span>Fermer la session / Changer de compte</span>
        </button>
      </div>
    </div>
  </div>
</template>

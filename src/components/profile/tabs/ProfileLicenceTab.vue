<script setup lang="ts">
import { ref, computed } from "vue";
import { useAuthStore } from "../../../stores/authStore";
import {
  ShieldCheck,
  Calendar,
  Monitor,
  Lock,
  CheckCircle2,
  Sparkles,
  RefreshCw,
} from "lucide-vue-next";

const emit = defineEmits<{
  (e: "open-pricing"): void;
}>();

const authStore = useAuthStore();
const perms = computed(() => authStore.rolePermissions);
const isRefreshing = ref(false);

const refreshLicenceData = async () => {
  isRefreshing.value = true;
  try {
    await authStore.loadProfileByRole();
  } finally {
    isRefreshing.value = false;
  }
};

const modulesList = [
  { key: "guichet", label: "Guichet & Inscriptions" },
  { key: "caisse", label: "Caisse & Reçus" },
  { key: "conduite", label: "Livret Conduite" },
  { key: "examens", label: "Examens d'État" },
  { key: "flotte", label: "Flotte Véhicules" },
  { key: "archivage_cloud", label: "Archivage Cloud R2" },
];
</script>

<template>
  <div class="space-y-5">
    <!-- Barre de synchronisation directe -->
    <div class="flex items-center justify-between text-xs">
      <div class="flex items-center gap-2 text-muted-foreground">
        <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
        <span>Synchronisé en temps réel avec la console Super-Admin & PostgreSQL</span>
      </div>
      <button
        type="button"
        @click="refreshLicenceData"
        :disabled="isRefreshing"
        class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-[4px] bg-muted hover:bg-muted/80 border border-border text-foreground font-medium text-[11px] transition cursor-pointer disabled:opacity-50"
      >
        <RefreshCw
          class="w-3 h-3"
          :class="isRefreshing ? 'animate-spin text-brand' : ''"
        />
        <span>{{ isRefreshing ? "Actualisation..." : "Actualiser" }}</span>
      </button>
    </div>

    <!-- Bannière Statut Formule Active & Période d'Essai -->
    <div
      class="p-4 rounded-[6px] border flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4"
      :class="
        authStore.isLicenceExpired
          ? 'bg-rose-500/10 border-rose-500/30'
          : authStore.trialDaysRemaining <= 15
            ? 'bg-amber-500/10 border-amber-500/30'
            : 'bg-emerald-500/10 border-emerald-500/30'
      "
    >
      <div class="flex items-start sm:items-center gap-3.5">
        <div
          class="w-11 h-11 rounded-[5px] flex items-center justify-center font-bold text-xl shrink-0"
          :class="
            authStore.isLicenceExpired
              ? 'bg-rose-500/20 text-rose-600 dark:text-rose-400 border border-rose-500/30'
              : 'bg-emerald-500/20 text-emerald-600 dark:text-emerald-400 border border-emerald-500/30'
          "
        >
          <ShieldCheck class="w-5 h-5" />
        </div>
        <div>
          <div class="flex items-center gap-2">
            <h3 class="text-sm font-bold text-foreground">
              Formule Active :
              <span class="font-mono text-brand">{{
                authStore.currentSchool?.type_licence || "TRIAL (30j gratuit)"
              }}</span>
            </h3>
            <span
              class="px-2 py-0.5 rounded-[3px] text-[10px] font-mono font-bold uppercase border"
              :class="
                authStore.isLicenceExpired
                  ? 'bg-rose-500/10 text-rose-600 border-rose-500/30'
                  : 'bg-emerald-500/10 text-emerald-600 border-emerald-500/30'
              "
            >
              {{ authStore.isLicenceExpired ? "Expirée" : "Active" }}
            </span>
          </div>
          <p class="text-xs text-muted-foreground mt-0.5">
            {{
              (authStore.currentSchool?.type_licence || "")
                .toUpperCase()
                .includes("TRIAL") ||
              (authStore.currentSchool?.type_licence || "")
                .toUpperCase()
                .includes("ESSAI")
                ? "Période d'essai offerte de 1 mois (30 jours complets) accordée à l'inscription."
                : "Abonnement professionnel ElevenDrive géré par le Super-Administrateur."
            }}
          </p>
        </div>
      </div>

      <!-- Compte à rebours / Durée restante -->
      <div class="text-left sm:text-right shrink-0">
        <div
          class="text-[11px] text-muted-foreground flex items-center gap-1.5 sm:justify-end"
        >
          <Calendar class="w-3.5 h-3.5 text-sky-500" />
          <span>Durée restante :</span>
        </div>
        <div
          class="font-mono font-extrabold text-base"
          :class="
            authStore.isLicenceExpired
              ? 'text-rose-600'
              : authStore.trialDaysRemaining <= 15
                ? 'text-amber-500'
                : 'text-emerald-600 dark:text-emerald-400'
          "
        >
          {{
            authStore.isLicenceExpired
              ? "Licence Échue"
              : `${authStore.trialDaysRemaining} jours restants`
          }}
        </div>
        <div class="text-[10px] text-muted-foreground font-mono mt-0.5">
          Échéance :
          {{
            authStore.currentSchool?.date_expiration_licence
              ? new Date(
                  authStore.currentSchool.date_expiration_licence,
                ).toLocaleDateString("fr-FR", {
                  day: "numeric",
                  month: "long",
                  year: "numeric",
                })
              : "Non définie"
          }}
        </div>
      </div>
    </div>

    <!-- Quota de Postes / Machines Desktop Autorisées -->
    <div class="p-4 rounded-[6px] bg-muted/40 border border-border space-y-3">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <Monitor class="w-4 h-4 text-brand" />
          <h4 class="font-bold text-xs text-foreground">
            Postes de Travail Desktop Autorisés
          </h4>
        </div>
        <span class="font-mono font-bold text-xs text-foreground">
          {{ authStore.currentSchool?.machines_actives || 0 }} /
          {{ authStore.currentSchool?.max_machines || 3 }} max autorisés
        </span>
      </div>

      <div
        class="h-2 w-full bg-muted/80 rounded-full overflow-hidden border border-border/50"
      >
        <div
          class="h-full rounded-full transition-all duration-300"
          :class="
            (authStore.currentSchool?.machines_actives || 0) >=
            (authStore.currentSchool?.max_machines || 3)
              ? 'bg-amber-500'
              : 'bg-emerald-500'
          "
          :style="{
            width: `${Math.min(100, Math.max(0, ((authStore.currentSchool?.machines_actives || 0) / (authStore.currentSchool?.max_machines || 3)) * 100))}%`,
          }"
        ></div>
      </div>

      <p class="text-[11px] text-muted-foreground leading-relaxed">
        Définit le nombre d'ordinateurs PC/Mac exécutant l'application ElevenDrive
        connectés simultanément à cet établissement.
      </p>
    </div>

    <!-- Matrice des Modules Système Autorisés par le Super-Admin -->
    <div class="p-4 rounded-[6px] bg-muted/40 border border-border space-y-3">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <Lock class="w-4 h-4 text-brand" />
          <h4 class="font-bold text-xs text-foreground">
            Modules Système & Droits d'Accès Débloqués
          </h4>
        </div>
        <span class="text-[10px] font-mono text-muted-foreground">
          Super-Admin Control
        </span>
      </div>

      <div class="grid grid-cols-2 sm:grid-cols-3 gap-2 pt-1 text-xs">
        <div
          v-for="mod in modulesList"
          :key="mod.key"
          class="p-2.5 rounded-[4px] border flex items-center justify-between"
          :class="
            authStore.isModuleEnabled(mod.key)
              ? 'bg-emerald-500/5 border-emerald-500/30 text-foreground'
              : 'bg-muted/60 border-border text-muted-foreground opacity-70'
          "
        >
          <span class="font-medium text-[11px]">{{ mod.label }}</span>
          <span
            v-if="authStore.isModuleEnabled(mod.key)"
            class="text-[10px] font-mono font-bold text-emerald-600 dark:text-emerald-400 flex items-center gap-1"
          >
            <CheckCircle2 class="w-3.5 h-3.5" />
            Actif
          </span>
          <span
            v-else
            class="text-[10px] font-mono text-muted-foreground flex items-center gap-1"
          >
            <Lock class="w-3.5 h-3.5" />
            Verrouillé
          </span>
        </div>
      </div>
    </div>

    <!-- Section Bon Plan & Surclassement (Directeur) -->
    <div
      v-if="perms.canManageSchool"
      class="p-4 rounded-[6px] bg-gradient-to-r from-amber-500/10 via-brand/10 to-sky-500/10 border border-amber-500/30 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3"
    >
      <div>
        <div class="flex items-center gap-2">
          <Sparkles class="w-4 h-4 text-amber-500 animate-pulse" />
          <h4 class="font-bold text-xs text-foreground">
            Gérer la Formule ou Ajouter des Postes
          </h4>
        </div>
        <p class="text-[11px] text-muted-foreground mt-0.5">
          Passez au plan supérieur, renouvelez votre licence ou débloquez plus de
          postes de travail.
        </p>
      </div>

      <button
        type="button"
        @click="emit('open-pricing')"
        class="px-4 py-2 rounded-[4px] bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-600 hover:to-amber-700 text-white font-bold text-xs flex items-center gap-1.5 transition cursor-pointer shadow-xs shrink-0"
      >
        <Sparkles class="w-3.5 h-3.5" />
        <span>Voir les Bons Plans & Tarifs</span>
      </button>
    </div>
  </div>
</template>

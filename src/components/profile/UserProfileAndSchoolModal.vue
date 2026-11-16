<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useAuthStore } from "../../stores/authStore";
import { useGuichetStore } from "../../stores/guichetStore";
import {
  X,
  User,
  Building2,
  MapPin,
  ShieldCheck,
  CheckCircle2,
  LogOut,
} from "lucide-vue-next";
import SchoolPricingModal from "../pricing/SchoolPricingModal.vue";
import ProfileUserTab, { type RoleBadgeMeta } from "./tabs/ProfileUserTab.vue";
import ProfileSchoolTab from "./tabs/ProfileSchoolTab.vue";
import ProfileAgencesTab from "./tabs/ProfileAgencesTab.vue";
import ProfileLicenceTab from "./tabs/ProfileLicenceTab.vue";

const emit = defineEmits<{
  (e: "close"): void;
}>();

const authStore = useAuthStore();
const guichetStore = useGuichetStore();

const activeTab = ref<"user" | "school" | "agences" | "licence">("user");
const isLoadingProfile = ref(false);
const isPricingOpen = ref(false);

const perms = computed(() => authStore.rolePermissions);

const roleBadgeMeta = computed<RoleBadgeMeta>(() => {
  switch (authStore.currentRole) {
    case "manager":
      return {
        label: "Directeur / Gérant",
        badgeClass:
          "bg-indigo-500/10 text-indigo-600 dark:text-indigo-400 border-indigo-500/20",
        description:
          "Administration générale, modification des tarifs, dérogations PIN & pilotage multi-agences",
        permissionsList: [
          "Gestion de l'établissement & agrément",
          "Création & suppression d'agences",
          "Modification des tarifs fixés",
          "Code PIN Manager & Dérogations",
          "Gestion de l'équipe & audit financier",
        ],
      };
    case "secretaire":
      return {
        label: authStore.currentUser.roleLabel || "Secrétaire Guichet",
        badgeClass:
          "bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/20",
        description:
          "Enrôlement des candidats, encaissement caisse, émargement code & bordereaux d'examen",
        permissionsList: [
          "Inscription candidats & création dossiers",
          "Encaissement & édition reçus caisse",
          "Fixation initiale des formules permis",
          "Point de caisse journalier (Rapport Z)",
          "Consultation réseau d'agences",
        ],
      };
    case "moniteur":
      return {
        label: "Moniteur de Conduite",
        badgeClass:
          "bg-sky-500/10 text-sky-600 dark:text-sky-400 border-sky-500/20",
        description:
          "Formation pratique et théorique, livret d'apprentissage mobile & suivi véhicule",
        permissionsList: [
          "Planning des séances de conduite",
          "Émargement tactile livret élève",
          "Animation cours de code & visio",
          "Suivi technique du véhicule assigné",
        ],
      };
    case "eleve":
    default:
      return {
        label: "Candidat / Élève Permis",
        badgeClass:
          "bg-purple-500/10 text-purple-600 dark:text-purple-400 border-purple-500/20",
        description:
          "Suivi pédagogique personnel, révision code 24/7, réservation conduite & reçus de paiement",
        permissionsList: [
          "Consultation dossier & solde financier",
          "Cours en ligne, PDF & tutoriels vidéo",
          "Séries Quiz Code & annales d'examen",
          "Réservation de créneaux de conduite",
        ],
      };
  }
});

const selectTab = async (tab: "user" | "school" | "agences" | "licence") => {
  activeTab.value = tab;
  if (tab === "licence") {
    await authStore.loadProfileByRole();
  }
};

onMounted(async () => {
  isLoadingProfile.value = true;
  try {
    await authStore.loadProfileByRole();
    if (
      authStore.currentRole === "eleve" &&
      !guichetStore.activeDossierDetail &&
      authStore.currentUser.id &&
      !authStore.currentUser.id.startsWith("usr_")
    ) {
      await guichetStore.selectCandidat(authStore.currentUser.id);
    }
  } finally {
    isLoadingProfile.value = false;
  }
});
</script>

<template>
  <div
    class="fixed inset-0 z-50 flex items-end sm:items-center justify-center p-0 sm:p-4 bg-black/60 animate-in fade-in duration-150"
    @click.self="emit('close')"
  >
    <div
      class="bg-card border-t sm:border border-border w-full sm:max-w-4xl max-h-[92dvh] rounded-t-[3px] sm:rounded-[3px] flex flex-col overflow-hidden text-foreground pb-safe"
      @keydown.esc="emit('close')"
    >
      <div
        class="sm:hidden mx-auto mt-2 h-1 w-9 rounded-full bg-border shrink-0"
      />
      <!-- En-tête Modale adapté au Rôle -->
      <div
        class="px-4 sm:px-5 py-3.5 sm:py-4 border-b border-border flex items-center justify-between bg-muted/30 shrink-0"
      >
        <div class="flex items-center space-x-3">
          <div
            class="w-10 h-10 rounded-[3px] bg-brand-light border border-brand-border text-brand-dark flex items-center justify-center font-black text-base"
          >
            {{ authStore.currentUser.initial }}
          </div>
          <div>
            <div class="flex flex-wrap items-center gap-2">
              <h2 class="text-base sm:text-sm font-bold text-foreground">
                Profil {{ roleBadgeMeta.label }}
              </h2>
              <span
                class="text-[10px] font-bold px-2 py-0.5 rounded-full border uppercase tracking-wider"
                :class="roleBadgeMeta.badgeClass"
              >
                {{ roleBadgeMeta.label }}
              </span>
            </div>
            <p class="text-xs text-muted-foreground mt-0.5">
              {{ authStore.schoolName }} · {{ roleBadgeMeta.description }}
            </p>
          </div>
        </div>

        <button
          type="button"
          @click="emit('close')"
          class="p-2 rounded-xl text-muted-foreground hover:text-foreground hover:bg-muted transition"
          title="Fermer (Échap)"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Navigation par Onglets -->
      <div
        class="flex border-b border-border bg-muted/15 px-3 sm:px-5 py-1 shrink-0 overflow-x-auto no-scrollbar gap-1"
      >
        <button
          type="button"
          @click="selectTab('user')"
          class="flex items-center space-x-1.5 py-2.5 px-2.5 text-xs font-semibold border-b-2 transition -mb-px shrink-0 cursor-pointer"
          :class="
            activeTab === 'user'
              ? 'border-brand text-brand-dark dark:text-brand'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          "
        >
          <User class="w-4 h-4" />
          <span>Mon Profil</span>
        </button>

        <button
          type="button"
          @click="selectTab('school')"
          class="flex items-center space-x-1.5 py-2.5 px-2.5 text-xs font-semibold border-b-2 transition -mb-px shrink-0 cursor-pointer"
          :class="
            activeTab === 'school'
              ? 'border-brand text-brand-dark dark:text-brand'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          "
        >
          <Building2 class="w-4 h-4" />
          <span>Mon Auto-École</span>
        </button>

        <button
          v-if="perms.canViewAgences"
          type="button"
          @click="selectTab('agences')"
          class="flex items-center space-x-1.5 py-2.5 px-2.5 text-xs font-semibold border-b-2 transition -mb-px shrink-0 cursor-pointer"
          :class="
            activeTab === 'agences'
              ? 'border-brand text-brand-dark dark:text-brand'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          "
        >
          <MapPin class="w-4 h-4" />
          <span>Agences ({{ authStore.agencesList.length }})</span>
        </button>

        <button
          type="button"
          @click="selectTab('licence')"
          class="flex items-center space-x-1.5 py-2.5 px-2.5 text-xs font-semibold border-b-2 transition -mb-px shrink-0 cursor-pointer"
          :class="
            activeTab === 'licence'
              ? 'border-brand text-brand-dark dark:text-brand'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          "
        >
          <ShieldCheck class="w-4 h-4 text-brand" />
          <span>Postes & Licence</span>
          <span
            v-if="authStore.trialDaysRemaining > 0"
            class="text-[10px] font-mono px-1.5 py-0.2 rounded-full bg-brand/10 text-brand font-bold border border-brand/30"
          >
            {{ authStore.trialDaysRemaining }}j
          </span>
        </button>
      </div>

      <!-- Corps Déroulant -->
      <div class="flex-1 overflow-y-auto p-5 sm:p-6">
        <ProfileUserTab
          v-if="activeTab === 'user'"
          :role-badge-meta="roleBadgeMeta"
        />
        <ProfileSchoolTab v-else-if="activeTab === 'school'" />
        <ProfileAgencesTab
          v-else-if="activeTab === 'agences' && perms.canViewAgences"
        />
        <ProfileLicenceTab
          v-else-if="activeTab === 'licence'"
          @open-pricing="isPricingOpen = true"
        />
      </div>

      <!-- Pied de Page Modale -->
      <div
        class="px-5 py-3 border-t border-border bg-muted/20 flex items-center justify-between text-xs text-muted-foreground shrink-0"
      >
        <div class="flex items-center gap-1.5">
          <CheckCircle2 class="w-4 h-4 text-emerald-500" />
          <span>
            Session {{ roleBadgeMeta.label }} synchronisée selon vos
            habilitations.
          </span>
        </div>

        <div class="flex items-center space-x-2">
          <button
            type="button"
            @click="
              authStore.logout();
              emit('close');
            "
            class="px-3.5 py-1.5 rounded-xl border border-rose-500/30 text-rose-600 dark:text-rose-400 hover:bg-rose-500/10 font-semibold transition flex items-center space-x-1.5"
          >
            <LogOut class="w-3.5 h-3.5" />
            <span>Se déconnecter</span>
          </button>
          <button
            type="button"
            @click="emit('close')"
            class="px-4 py-1.5 rounded-xl border border-border text-foreground hover:bg-muted font-semibold transition"
          >
            Fermer
          </button>
        </div>
      </div>
    </div>
  </div>

  <!-- Modal Bons Plans & Tarifs déclenchable depuis l'onglet Licence -->
  <SchoolPricingModal :is-open="isPricingOpen" @close="isPricingOpen = false" />
</template>

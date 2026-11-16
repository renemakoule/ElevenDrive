<script setup lang="ts">
import { ref, reactive, computed, watch } from "vue";
import { useAuthStore } from "../../../stores/authStore";
import { useGuichetStore } from "../../../stores/guichetStore";
import {
  Phone,
  Mail,
  MapPin,
  ShieldCheck,
  CheckCircle2,
  Award,
  Car,
  Lock,
  Save,
} from "lucide-vue-next";

export interface RoleBadgeMeta {
  label: string;
  badgeClass: string;
  description: string;
  permissionsList: string[];
}

const props = defineProps<{
  roleBadgeMeta: RoleBadgeMeta;
}>();

const authStore = useAuthStore();
const guichetStore = useGuichetStore();

const isSaving = ref(false);
const perms = computed(() => authStore.rolePermissions);
const eleveDossier = computed(() => guichetStore.activeDossierDetail);
const moniteurVehicule = computed(() => guichetStore.vehicules?.[0] || null);

const userForm = reactive({
  nom: authStore.currentUser.nom || "",
  prenom: authStore.currentUser.prenom || "",
  telephone: authStore.currentUser.telephone || "",
  email: authStore.currentUser.email || "",
  agence:
    authStore.currentAgence?.nom ||
    authStore.currentUser.agence ||
    "Agence Principale",
  dossier: authStore.currentUser.dossier || "",
  code_pin_manager: authStore.currentUser.code_pin_manager || "",
  nouveau_mot_de_passe: "",
  confirmer_mot_de_passe: "",
});

const syncFromStore = () => {
  userForm.nom = authStore.currentUser.nom || "";
  userForm.prenom = authStore.currentUser.prenom || "";
  userForm.telephone = authStore.currentUser.telephone || "";
  userForm.email = authStore.currentUser.email || "";
  userForm.agence =
    authStore.currentAgence?.nom ||
    authStore.currentUser.agence ||
    "Agence Principale";
  userForm.dossier =
    eleveDossier.value?.dossier?.numero_dossier ||
    authStore.currentUser.dossier ||
    "";
  userForm.code_pin_manager = authStore.currentUser.code_pin_manager || "";
};

watch(
  () => [authStore.currentUser, authStore.currentAgence, eleveDossier.value],
  () => syncFromStore(),
  { deep: true, immediate: true },
);

const saveUserProfile = async () => {
  if (!userForm.nom.trim() || !userForm.prenom.trim()) {
    guichetStore.showToast("Le nom et le prénom sont obligatoires.", "error");
    return;
  }
  if (
    perms.value.canChangePassword &&
    userForm.nouveau_mot_de_passe &&
    userForm.nouveau_mot_de_passe !== userForm.confirmer_mot_de_passe
  ) {
    guichetStore.showToast(
      "Les mots de passe saisis ne correspondent pas.",
      "error",
    );
    return;
  }
  if (
    perms.value.canConfigurePin &&
    userForm.code_pin_manager &&
    !/^\d{4}$/.test(userForm.code_pin_manager.trim())
  ) {
    guichetStore.showToast(
      "Le code PIN Manager doit être composé exactement de 4 chiffres.",
      "error",
    );
    return;
  }

  isSaving.value = true;
  try {
    await authStore.updateUserProfile({
      nom: userForm.nom.trim(),
      prenom: userForm.prenom.trim(),
      telephone: userForm.telephone.trim(),
      email: userForm.email.trim() || undefined,
      code_pin_manager: perms.value.canConfigurePin
        ? userForm.code_pin_manager.trim() || undefined
        : undefined,
      nouveau_mot_de_passe: perms.value.canChangePassword
        ? userForm.nouveau_mot_de_passe || undefined
        : undefined,
    });
    guichetStore.showToast(
      `Profil ${props.roleBadgeMeta.label} mis à jour avec succès.`,
      "success",
    );
    userForm.nouveau_mot_de_passe = "";
    userForm.confirmer_mot_de_passe = "";
  } catch (err: any) {
    guichetStore.showToast(
      err?.message || "Erreur lors de la sauvegarde du profil.",
      "error",
    );
  } finally {
    isSaving.value = false;
  }
};
</script>

<template>
  <div class="space-y-6">
    <!-- Carte d'identité visuelle & Permissions du Rôle -->
    <div
      class="p-4 rounded-[3px] bg-muted/40 border border-border flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4"
    >
      <div class="flex items-center space-x-4">
        <!-- <div
          class="w-14 h-14 rounded-[5px] bg-brand-light text-brand-dark flex items-center justify-center font-black text-xl uppercase border border-brand-border"
        >
          {{ userForm.prenom?.[0] || "U" }}
        </div> -->
        <div>
          <div
            class="font-bold text-base text-foreground flex items-center gap-2"
          >
            <span>{{ userForm.prenom }} {{ userForm.nom }}</span>
            <span
              class="text-[10px] font-bold px-2 py-0.5 rounded-full border uppercase"
              :class="roleBadgeMeta.badgeClass"
            >
              {{ roleBadgeMeta.label }}
            </span>
          </div>
          <div
            class="text-xs text-muted-foreground mt-1 flex flex-wrap gap-x-4 gap-y-1"
          >
            <span v-if="userForm.telephone" class="flex items-center gap-1">
              <Phone class="w-3.5 h-3.5 text-muted-foreground" />
              {{ userForm.telephone }}
            </span>
            <span v-if="userForm.email" class="flex items-center gap-1">
              <Mail class="w-3.5 h-3.5 text-muted-foreground" />
              {{ userForm.email }}
            </span>
            <span v-if="userForm.agence" class="flex items-center gap-1">
              <MapPin class="w-3.5 h-3.5 text-brand" />
              {{ userForm.agence }}
            </span>
          </div>
        </div>
      </div>

      <div class="text-right">
        <div class="text-[11px] text-muted-foreground">Habilitation Rôle</div>
        <div class="text-xs font-mono font-bold text-emerald-500">
          {{ authStore.currentRole.toUpperCase() }} · ACTIF
        </div>
      </div>
    </div>

    <!-- Encart Permissions du Rôle -->
    <div class="p-3.5 rounded-[4px] bg-card border border-border">
      <div
        class="text-[11px] font-bold uppercase tracking-wider text-muted-foreground mb-2 flex items-center gap-1.5"
      >
        <ShieldCheck class="w-3.5 h-3.5 text-brand" />
        Permissions accordées à votre profil ({{ roleBadgeMeta.label }})
      </div>
      <div class="flex flex-wrap gap-2">
        <span
          v-for="perm in roleBadgeMeta.permissionsList"
          :key="perm"
          class="inline-flex items-center gap-1 text-[11px] px-2.5 py-1 rounded-full bg-muted text-foreground border border-border"
        >
          <CheckCircle2 class="w-3 h-3 text-emerald-500 shrink-0" />
          {{ perm }}
        </span>
      </div>
    </div>

    <!-- SECTION SPÉCIFIQUE RÔLE ÉLÈVE : Dossier d'Inscription & Progression -->
    <div
      v-if="perms.isEleve"
      class="p-4 rounded-[4px] bg-purple-500/5 border border-purple-500/20 space-y-3"
    >
      <div class="flex items-center justify-between">
        <h3
          class="text-xs font-bold uppercase tracking-wider text-purple-600 dark:text-purple-400 flex items-center gap-1.5"
        >
          <Award class="w-4 h-4" /> Mon Dossier Candidat Permis
        </h3>
        <span
          class="text-xs font-mono font-bold px-2.5 py-0.5 rounded bg-purple-500/10 text-purple-600 dark:text-purple-300"
        >
          N° {{ userForm.dossier || "En cours d'attribution" }}
        </span>
      </div>

      <div
        v-if="eleveDossier"
        class="grid grid-cols-2 sm:grid-cols-4 gap-3 pt-1"
      >
        <div class="p-2.5 rounded bg-card border border-border">
          <div class="text-[10px] text-muted-foreground">Formule inscrite</div>
          <div class="text-xs font-bold text-foreground truncate">
            {{ eleveDossier.formule?.intitule || "Permis B" }}
          </div>
        </div>
        <div class="p-2.5 rounded bg-card border border-border">
          <div class="text-[10px] text-muted-foreground">Séances Code</div>
          <div class="text-xs font-bold text-foreground">
            {{ eleveDossier.dossier.heures_code_effectuees }} /
            {{ eleveDossier.formule?.heures_theoriques || 20 }}h
          </div>
        </div>
        <div class="p-2.5 rounded bg-card border border-border">
          <div class="text-[10px] text-muted-foreground">Heures Conduite</div>
          <div class="text-xs font-bold text-foreground">
            {{ eleveDossier.dossier.heures_conduite_effectuees }} /
            {{ eleveDossier.formule?.heures_pratiques || 20 }}h
          </div>
        </div>
        <div class="p-2.5 rounded bg-card border border-border">
          <div class="text-[10px] text-muted-foreground">Solde Restant</div>
          <div
            class="text-xs font-bold"
            :class="
              eleveDossier.dossier.solde_restant <= 0
                ? 'text-emerald-500'
                : 'text-amber-500'
            "
          >
            {{
              (eleveDossier.dossier.solde_restant || 0).toLocaleString("fr-FR")
            }}
            FCFA
          </div>
        </div>
      </div>
    </div>

    <!-- SECTION SPÉCIFIQUE RÔLE MONITEUR : Affectation Pédagogique & Flotte -->
    <div
      v-if="perms.isMoniteur"
      class="p-4 rounded-[4px] bg-sky-500/5 border border-sky-500/20 space-y-3"
    >
      <h3
        class="text-xs font-bold uppercase tracking-wider text-sky-600 dark:text-sky-400 flex items-center gap-1.5"
      >
        <Car class="w-4 h-4" /> Affectation Pédagogique & Flotte
      </h3>
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
        <div class="p-2.5 rounded bg-card border border-border">
          <div class="text-[10px] text-muted-foreground">
            Agence de rattachement
          </div>
          <div class="text-xs font-bold text-foreground">
            {{ userForm.agence }}
          </div>
        </div>
        <div class="p-2.5 rounded bg-card border border-border">
          <div class="text-[10px] text-muted-foreground">
            Véhicule d'instruction
          </div>
          <div class="text-xs font-bold text-foreground">
            {{
              moniteurVehicule
                ? `${moniteurVehicule.marque_modele} (${moniteurVehicule.immatriculation})`
                : "Flotte Auto-École"
            }}
          </div>
        </div>
        <div class="p-2.5 rounded bg-card border border-border">
          <div class="text-[10px] text-muted-foreground">
            Statut Instructeur
          </div>
          <div class="text-xs font-bold text-emerald-500">
            Habilité Code & Conduite
          </div>
        </div>
      </div>
    </div>

    <!-- Formulaire de Modification Utilisateur -->
    <form @submit.prevent="saveUserProfile" class="space-y-4">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div>
          <label
            class="block text-xs font-semibold text-muted-foreground mb-1.5"
          >
            Prénom(s) *
          </label>
          <input
            v-model="userForm.prenom"
            type="text"
            required
            class="w-full px-3 py-2 rounded-xl bg-card border border-border text-foreground text-sm focus:ring-2 focus:ring-brand/40 focus:border-brand outline-none transition"
            placeholder="Ex : René"
          />
        </div>

        <div>
          <label
            class="block text-xs font-semibold text-muted-foreground mb-1.5"
          >
            Nom de famille *
          </label>
          <input
            v-model="userForm.nom"
            type="text"
            required
            class="w-full px-3 py-2 rounded-xl bg-card border border-border text-foreground text-sm focus:ring-2 focus:ring-brand/40 focus:border-brand outline-none transition"
            placeholder="Ex : Makou"
          />
        </div>

        <div>
          <label
            class="block text-xs font-semibold text-muted-foreground mb-1.5"
          >
            Numéro de Téléphone (WhatsApp) *
          </label>
          <div class="relative">
            <Phone
              class="w-4 h-4 text-muted-foreground absolute left-3 top-2.5"
            />
            <input
              v-model="userForm.telephone"
              type="tel"
              required
              class="w-full pl-9 pr-3 py-2 rounded-xl bg-card border border-border text-foreground text-sm focus:ring-2 focus:ring-brand/40 focus:border-brand outline-none transition"
              placeholder="Ex : 651727932"
            />
          </div>
        </div>

        <div>
          <label
            class="block text-xs font-semibold text-muted-foreground mb-1.5"
          >
            Adresse Email
          </label>
          <div class="relative">
            <Mail
              class="w-4 h-4 text-muted-foreground absolute left-3 top-2.5"
            />
            <input
              v-model="userForm.email"
              type="email"
              class="w-full pl-9 pr-3 py-2 rounded-xl bg-card border border-border text-foreground text-sm focus:ring-2 focus:ring-brand/40 focus:border-brand outline-none transition"
              placeholder="contact@auto-ecole.com"
            />
          </div>
        </div>
      </div>

      <!-- Paramètres de Sécurité -->
      <div
        v-if="perms.canConfigurePin || perms.canChangePassword"
        class="pt-4 border-t border-border mt-4"
      >
        <h3
          class="text-xs font-bold uppercase text-muted-foreground tracking-wider mb-3 flex items-center gap-1.5"
        >
          <Lock class="w-3.5 h-3.5 text-brand" />
          {{
            perms.canConfigurePin
              ? "Sécurité du Compte & Code PIN Manager (Directeur)"
              : "Sécurité & Mot de passe de session"
          }}
        </h3>

        <div
          class="grid grid-cols-1 gap-4"
          :class="perms.canConfigurePin ? 'sm:grid-cols-3' : 'sm:grid-cols-2'"
        >
          <div v-if="perms.canConfigurePin">
            <label
              class="block text-xs font-semibold text-muted-foreground mb-1.5"
            >
              Code PIN Manager (4 chiffres)
            </label>
            <input
              v-model="userForm.code_pin_manager"
              type="password"
              maxlength="4"
              class="w-full px-3 py-2 rounded-xl bg-card border border-border text-foreground text-sm font-mono tracking-widest text-center focus:ring-2 focus:ring-brand/40 focus:border-brand outline-none transition"
              placeholder="••••"
            />
            <p class="text-[10px] text-muted-foreground mt-1">
              Réservé au Directeur pour autoriser les dérogations au guichet.
            </p>
          </div>

          <div v-if="perms.canChangePassword">
            <label
              class="block text-xs font-semibold text-muted-foreground mb-1.5"
            >
              Nouveau mot de passe
            </label>
            <input
              v-model="userForm.nouveau_mot_de_passe"
              type="password"
              class="w-full px-3 py-2 rounded-xl bg-card border border-border text-foreground text-sm focus:ring-2 focus:ring-brand/40 focus:border-brand outline-none transition"
              placeholder="Laisser vide si inchangé"
            />
          </div>

          <div v-if="perms.canChangePassword">
            <label
              class="block text-xs font-semibold text-muted-foreground mb-1.5"
            >
              Confirmer le mot de passe
            </label>
            <input
              v-model="userForm.confirmer_mot_de_passe"
              type="password"
              class="w-full px-3 py-2 rounded-xl bg-card border border-border text-foreground text-sm focus:ring-2 focus:ring-brand/40 focus:border-brand outline-none transition"
              placeholder="Répéter le mot de passe"
            />
          </div>
        </div>
      </div>

      <div class="flex justify-end pt-3">
        <button
          type="submit"
          :disabled="isSaving"
          class="flex items-center space-x-2 px-5 py-2.5 rounded-xl bg-brand hover:bg-brand-dark text-brand-dark font-bold text-xs uppercase tracking-wider transition shadow-md disabled:opacity-50"
        >
          <Save class="w-4 h-4" />
          <span>
            {{
              isSaving
                ? "Enregistrement..."
                : `Enregistrer mon profil ${roleBadgeMeta.label}`
            }}
          </span>
        </button>
      </div>
    </form>
  </div>
</template>

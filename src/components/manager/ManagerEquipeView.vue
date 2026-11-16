<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { useAuthStore } from "../../stores/authStore";
import { useGuichetStore } from "../../stores/guichetStore";
import { authBackend } from "../../services/authBackendService";
import { realtimeSocket } from "../../services/realtimeSocketService";
import EquipeMemberCard from "./equipe/EquipeMemberCard.vue";
import EquipeMemberFormModal from "./equipe/EquipeMemberFormModal.vue";
import EquipeDeleteModal from "./equipe/EquipeDeleteModal.vue";
import type { UserTeamMember, Agence, RegisterUserPayload } from "../../types";
import {
  Users,
  UserPlus,
  RefreshCw,
  Search,
  X,
  CheckCircle2,
  Copy,
  Share2,
} from "lucide-vue-next";

const authStore = useAuthStore();
const guichetStore = useGuichetStore();
let unsubscribeWs: (() => void) | null = null;

const team = ref<UserTeamMember[]>([]);
const agences = ref<Agence[]>([]);
const isLoading = ref(false);
const searchQuery = ref("");
const selectedRoleFilter = ref<string>("ALL");
const selectedAgenceFilter = ref<string>("ALL");

// État Modal d'Ajout / Modification
const isAddModalOpen = ref(false);
const modalMode = ref<"ADD" | "EDIT">("ADD");
const editingMember = ref<UserTeamMember | null>(null);
const isSubmitting = ref(false);
const modalError = ref<string | null>(null);

// État Modal de Suppression
const isDeleteModalOpen = ref(false);
const memberToDelete = ref<UserTeamMember | null>(null);
const isDeleting = ref(false);

// Dernier utilisateur créé avec succès pour partage
const lastCreatedUser = ref<{
  nom: string;
  prenoms: string;
  role: string;
  telephone: string;
  mot_de_passe: string;
  agence: string;
} | null>(null);

const loadData = async () => {
  isLoading.value = true;
  try {
    const [teamData, agencesData] = await Promise.all([
      authBackend.getSchoolUsers(),
      authBackend.getSchoolAgences(),
    ]);
    team.value = teamData;
    agences.value = agencesData;
  } catch (err: any) {
    guichetStore.showToast("Erreur chargement équipe: " + err.message, "error");
  } finally {
    isLoading.value = false;
  }
};

onMounted(async () => {
  await loadData();
  unsubscribeWs = realtimeSocket.subscribe((msg) => {
    if (msg.type === "EQUIPE_USERS_UPDATED") {
      loadData();
    }
  });
});

onUnmounted(() => {
  if (unsubscribeWs) {
    unsubscribeWs();
    unsubscribeWs = null;
  }
});

// Filtres
const filteredTeam = computed(() => {
  return team.value.filter((member) => {
    if (
      selectedRoleFilter.value !== "ALL" &&
      member.role !== selectedRoleFilter.value
    ) {
      return false;
    }
    if (
      selectedAgenceFilter.value !== "ALL" &&
      member.agence_id !== selectedAgenceFilter.value
    ) {
      return false;
    }
    if (searchQuery.value.trim()) {
      const q = searchQuery.value.toLowerCase().trim();
      const nomComplet = `${member.nom} ${member.prenoms}`.toLowerCase();
      const tel = member.telephone.toLowerCase();
      const email = (member.email || "").toLowerCase();
      return nomComplet.includes(q) || tel.includes(q) || email.includes(q);
    }
    return true;
  });
});

// KPIs
const totalCollaborateurs = computed(() => team.value.length);
const totalMoniteurs = computed(
  () => team.value.filter((m) => m.role === "MONITEUR").length,
);
const totalSecretaires = computed(
  () =>
    team.value.filter((m) => m.role === "SECRETAIRE" || m.role === "CAISSIER")
      .length,
);
const totalAgences = computed(() => agences.value.length || 1);

const handleOpenAddModal = () => {
  modalMode.value = "ADD";
  editingMember.value = null;
  modalError.value = null;
  isAddModalOpen.value = true;
};

const handleOpenEditModal = (member: UserTeamMember) => {
  modalMode.value = "EDIT";
  editingMember.value = member;
  modalError.value = null;
  isAddModalOpen.value = true;
};

const handleSaveUser = async (data: {
  form: RegisterUserPayload;
  statut: "ACTIF" | "SUSPENDU";
}) => {
  if (isSubmitting.value) return;
  modalError.value = null;
  isSubmitting.value = true;

  try {
    const formVal = data.form;
    const targetAgence = agences.value.find((a) => a.id === formVal.agence_id);
    const validAgenceId =
      formVal.agence_id &&
      !formVal.agence_id.startsWith("ag-") &&
      formVal.agence_id.length >= 32
        ? formVal.agence_id
        : targetAgence &&
            !targetAgence.id.startsWith("ag-") &&
            targetAgence.id.length >= 32
          ? targetAgence.id
          : undefined;

    const payload = {
      ...formVal,
      agence_id: validAgenceId,
    };

    if (modalMode.value === "ADD") {
      await authBackend.registerUser(payload);

      lastCreatedUser.value = {
        nom: formVal.nom.toUpperCase(),
        prenoms: formVal.prenoms,
        role: formVal.role,
        telephone: formVal.telephone,
        mot_de_passe: formVal.mot_de_passe,
        agence: targetAgence?.nom || "Agence Siège",
      };

      guichetStore.showToast(
        `Collaborateur ${formVal.prenoms} ajouté avec succès !`,
        "success",
      );
    } else if (modalMode.value === "EDIT" && editingMember.value) {
      await authBackend.updateUser(editingMember.value.id, {
        ...payload,
        statut: data.statut,
      });

      guichetStore.showToast(
        `Collaborateur ${formVal.prenoms} mis à jour avec succès !`,
        "success",
      );
    }

    isAddModalOpen.value = false;
    await loadData();
  } catch (err: any) {
    // Si l'erreur indique que le téléphone existe déjà, vérifier si l'utilisateur a en fait été créé lors d'une requête concurrente
    if (
      modalMode.value === "ADD" &&
      (err.message || "").includes("existe déjà")
    ) {
      await loadData();
      const cleanInputTel = data.form.telephone.replace(/\s+/g, "");
      const alreadyCreated = team.value.some(
        (m) => m.telephone.replace(/\s+/g, "") === cleanInputTel,
      );
      if (alreadyCreated) {
        guichetStore.showToast(
          `Collaborateur ${data.form.prenoms} ajouté avec succès !`,
          "success",
        );
        isAddModalOpen.value = false;
        return;
      }
    }
    modalError.value = err.message || "Erreur lors de l'enregistrement.";
  } finally {
    isSubmitting.value = false;
  }
};

const handleOpenDeleteModal = (member: UserTeamMember) => {
  memberToDelete.value = member;
  isDeleteModalOpen.value = true;
};

const handleConfirmDelete = async () => {
  if (!memberToDelete.value) return;
  isDeleting.value = true;
  try {
    await authBackend.deleteUser(memberToDelete.value.id);
    guichetStore.showToast(
      `Collaborateur ${memberToDelete.value.prenoms} ${memberToDelete.value.nom} supprimé`,
      "success",
    );
    isDeleteModalOpen.value = false;
    memberToDelete.value = null;
    await loadData();
  } catch (err: any) {
    guichetStore.showToast("Erreur suppression : " + err.message, "error");
  } finally {
    isDeleting.value = false;
  }
};

const copyCredentials = async () => {
  if (!lastCreatedUser.value) return;
  const text =
    `Identifiants AutoSchool (${lastCreatedUser.value.agence}) :\n` +
    `Utilisateur : ${lastCreatedUser.value.prenoms} ${lastCreatedUser.value.nom} (${lastCreatedUser.value.role})\n` +
    `Identifiant : ${lastCreatedUser.value.telephone}\n` +
    `Mot de passe : ${lastCreatedUser.value.mot_de_passe}\n` +
    `Lien d'accès : ${window.location.origin}/${authStore.currentSchool?.code_unique || ""}`;

  try {
    await navigator.clipboard.writeText(text);
    guichetStore.showToast(
      "Identifiants copiés dans le presse-papiers !",
      "success",
    );
  } catch {
    guichetStore.showToast("Identifiants copiés !", "success");
  }
};

const shareWhatsAppCredentials = () => {
  if (!lastCreatedUser.value) return;
  const text = encodeURIComponent(
    `Bonjour ${lastCreatedUser.value.prenoms} !\nVoici vos accès à la plateforme AutoSchool (${authStore.schoolName}) :\n\n` +
      `• Identifiant : ${lastCreatedUser.value.telephone}\n` +
      `• Mot de passe : ${lastCreatedUser.value.mot_de_passe}\n` +
      `• Agence : ${lastCreatedUser.value.agence}\n\n` +
      `Connectez-vous ici : ${window.location.origin}/${authStore.currentSchool?.code_unique || ""}`,
  );
  window.open(
    `https://wa.me/${lastCreatedUser.value.telephone.replace(/[^0-9]/g, "")}?text=${text}`,
    "_blank",
  );
};
</script>

<template>
  <div class="space-y-5 p-6 animate-in fade-in duration-200">
    <!-- En-tête -->
    <div
      class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4"
    >
      <div>
        <div class="flex items-center gap-2.5">
          <div
            class="flex h-9 w-9 items-center justify-center rounded-xl border border-brand-border bg-brand-light text-brand"
          >
            <Users class="h-5 w-5" />
          </div>
          <div>
            <h1 class="text-base font-bold text-foreground">
              Équipe & Collaborateurs
            </h1>
            <p class="text-xs text-muted-foreground">
              Gestion du personnel, rôles, affectations aux agences et accès
              applicatifs
            </p>
          </div>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <button
          type="button"
          @click="loadData"
          class="inline-flex h-9 items-center gap-1.5 rounded-lg border border-border bg-card px-3 text-xs font-medium text-foreground transition hover:bg-muted"
        >
          <RefreshCw
            class="h-3.5 w-3.5 text-muted-foreground"
            :class="{ 'animate-spin': isLoading }"
          />
          <span>Actualiser</span>
        </button>

        <button
          type="button"
          @click="handleOpenAddModal"
          class="inline-flex h-9 items-center gap-1.5 rounded-lg bg-brand px-3.5 text-xs font-semibold text-white transition hover:bg-brand-dark shadow-sm"
        >
          <UserPlus class="h-4 w-4" />
          <span>Nouveau Collaborateur</span>
        </button>
      </div>
    </div>

    <!-- Bannière récapitulative des identifiants créés -->
    <div
      v-if="lastCreatedUser"
      class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 rounded-[5px] border border-emerald-500/20 bg-emerald-500/10 p-4 text-xs text-emerald-900 dark:text-emerald-300 animate-in slide-in-from-top-2 duration-200"
    >
      <div class="flex items-start gap-2.5">
        <CheckCircle2
          class="h-5 w-5 text-emerald-600 dark:text-emerald-400 shrink-0 mt-0.5"
        />
        <div>
          <span class="font-semibold text-sm"
            >Compte créé pour {{ lastCreatedUser.prenoms }}
            {{ lastCreatedUser.nom }}</span
          >
          <p class="mt-0.5 text-emerald-700 dark:text-emerald-400">
            Identifiant :
            <strong class="font-mono">{{ lastCreatedUser.telephone }}</strong> ·
            Mot de passe :
            <strong class="font-mono">{{
              lastCreatedUser.mot_de_passe
            }}</strong>
            · Agence : <strong>{{ lastCreatedUser.agence }}</strong>
          </p>
        </div>
      </div>
      <div class="flex items-center gap-2 shrink-0">
        <button
          type="button"
          @click="copyCredentials"
          class="inline-flex h-8 items-center gap-1 rounded-md border border-emerald-600/30 bg-background px-2.5 text-xs font-medium text-foreground transition hover:bg-muted"
        >
          <Copy class="h-3 w-3 text-emerald-600" />
          <span>Copier</span>
        </button>
        <button
          type="button"
          @click="shareWhatsAppCredentials"
          class="inline-flex h-8 items-center gap-1 rounded-md bg-emerald-600 px-2.5 text-xs font-medium text-white transition hover:bg-emerald-700"
        >
          <Share2 class="h-3 w-3" />
          <span>WhatsApp</span>
        </button>
        <button
          type="button"
          @click="lastCreatedUser = null"
          class="rounded p-1 text-emerald-600 hover:text-emerald-800"
        >
          <X class="h-4 w-4" />
        </button>
      </div>
    </div>

    <!-- KPIs Équipe -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-3">
      <div class="rounded-[5px] border border-border bg-card p-4">
        <span class="text-xs font-medium text-muted-foreground"
          >Effectif total</span
        >
        <div class="mt-1 flex items-baseline gap-2">
          <span class="text-2xl font-bold text-foreground">{{
            totalCollaborateurs
          }}</span>
          <span class="text-xs text-muted-foreground">membres</span>
        </div>
      </div>
      <div class="rounded-[5px] border border-border bg-card p-4">
        <span class="text-xs font-medium text-muted-foreground"
          >Moniteurs de conduite</span
        >
        <div class="mt-1 flex items-baseline gap-2">
          <span class="text-2xl font-bold text-sky-600 dark:text-sky-400">{{
            totalMoniteurs
          }}</span>
          <span class="text-xs text-muted-foreground">actifs PWA</span>
        </div>
      </div>
      <div class="rounded-[5px] border border-border bg-card p-4">
        <span class="text-xs font-medium text-muted-foreground"
          >Secrétariat & Caisse</span
        >
        <div class="mt-1 flex items-baseline gap-2">
          <span
            class="text-2xl font-bold text-emerald-600 dark:text-emerald-400"
            >{{ totalSecretaires }}</span
          >
          <span class="text-xs text-muted-foreground">au guichet</span>
        </div>
      </div>
      <div class="rounded-[5px] border border-border bg-card p-4">
        <span class="text-xs font-medium text-muted-foreground"
          >Agences couvertes</span
        >
        <div class="mt-1 flex items-baseline gap-2">
          <span class="text-2xl font-bold text-foreground">{{
            totalAgences
          }}</span>
          <span class="text-xs text-muted-foreground">succursales</span>
        </div>
      </div>
    </div>

    <!-- Barre d'outils et filtres -->
    <div
      class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3 rounded-[5px] border border-border bg-card p-3"
    >
      <div class="relative flex-1 max-w-sm">
        <Search class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Rechercher par nom, téléphone, email..."
          class="h-9 w-full rounded-[8px] border border-border bg-muted/40 pl-9 pr-3 text-xs text-foreground placeholder:text-muted-foreground/60 focus:border-brand focus:outline-none"
        />
      </div>

      <div class="flex items-center gap-2 overflow-x-auto">
        <select
          v-model="selectedRoleFilter"
          class="h-9 rounded-[5px] border border-border bg-muted/40 px-2.5 text-xs text-foreground focus:border-brand focus:outline-none"
        >
          <option value="ALL">Tous les rôles</option>
          <option value="MONITEUR">Moniteurs</option>
          <option value="SECRETAIRE">Secrétaires</option>
          <option value="CAISSIER">Caissières</option>
          <option value="DIRECTEUR">Direction</option>
        </select>

        <select
          v-model="selectedAgenceFilter"
          class="h-9 rounded-[5px] border border-border bg-muted/40 px-2.5 text-xs text-foreground focus:border-brand focus:outline-none"
        >
          <option value="ALL">Toutes les agences ({{ agences.length }})</option>
          <option v-for="ag in agences" :key="ag.id" :value="ag.id">
            {{ ag.nom }} ({{
              ag.est_siege ? "Siège" : ag.commune || ag.ville || "Succursale"
            }})
          </option>
        </select>
      </div>
    </div>

    <!-- Liste des collaborateurs -->
    <div
      v-if="filteredTeam.length > 0"
      class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3"
    >
      <EquipeMemberCard
        v-for="member in filteredTeam"
        :key="member.id"
        :member="member"
        :is-current-user="member.id === authStore.currentUser?.id"
        @edit="handleOpenEditModal"
        @delete="handleOpenDeleteModal"
      />
    </div>

    <!-- Empty state -->
    <div
      v-else
      class="flex flex-col items-center justify-center rounded-[5px] border border-dashed border-border bg-card/50 p-12 text-center"
    >
      <div
        class="flex h-12 w-12 items-center justify-center rounded-full bg-muted text-muted-foreground"
      >
        <Users class="h-6 w-6" />
      </div>
      <h3 class="mt-3 text-sm font-semibold text-foreground">
        Aucun collaborateur trouvé
      </h3>
      <p class="mt-1 max-w-sm text-xs text-muted-foreground">
        {{
          searchQuery
            ? "Aucun membre ne correspond à votre recherche."
            : "Ajoutez votre premier moniteur ou secrétaire pour lui donner accès à la plateforme."
        }}
      </p>
      <button
        type="button"
        @click="handleOpenAddModal"
        class="mt-4 inline-flex items-center gap-1.5 rounded-lg bg-brand px-3.5 py-2 text-xs font-semibold text-white transition hover:bg-brand-dark"
      >
        <UserPlus class="h-4 w-4" />
        <span>Ajouter un collaborateur</span>
      </button>
    </div>

    <!-- Modales -->
    <EquipeMemberFormModal
      :is-open="isAddModalOpen"
      :mode="modalMode"
      :member="editingMember"
      :agences="agences"
      :is-submitting="isSubmitting"
      :error-msg="modalError"
      @close="isAddModalOpen = false"
      @save="handleSaveUser"
    />

    <EquipeDeleteModal
      :is-open="isDeleteModalOpen"
      :member="memberToDelete"
      :is-deleting="isDeleting"
      @close="isDeleteModalOpen = false"
      @confirm="handleConfirmDelete"
    />
  </div>
</template>

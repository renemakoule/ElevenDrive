<script setup lang="ts">
import { ref, reactive, computed } from "vue";
import { useAuthStore } from "../../../stores/authStore";
import { useGuichetStore } from "../../../stores/guichetStore";
import type { Agence } from "../../../types";
import { MapPin, Plus, Edit2, Trash2 } from "lucide-vue-next";

const authStore = useAuthStore();
const guichetStore = useGuichetStore();

const perms = computed(() => authStore.rolePermissions);
const isEditingAgence = ref(false);
const editingAgenceId = ref<string | null>(null);
const isSavingAgence = ref(false);

const agenceForm = reactive({
  nom: "",
  code_agence: "",
  adresse: "",
  ville: "",
  commune: "",
  telephone: "",
  est_siege: false,
});

const startAddAgence = () => {
  if (!perms.value.canManageAgences) return;
  isEditingAgence.value = true;
  editingAgenceId.value = null;
  agenceForm.nom = "";
  agenceForm.code_agence = `AG-${Math.floor(10 + Math.random() * 90)}`;
  agenceForm.adresse = "";
  agenceForm.ville = authStore.currentSchool?.ville || "Douala";
  agenceForm.commune = "";
  agenceForm.telephone = authStore.currentSchool?.telephone || "";
  agenceForm.est_siege = false;
};

const startEditAgence = (agence: Agence) => {
  if (!perms.value.canManageAgences) return;
  isEditingAgence.value = true;
  editingAgenceId.value = agence.id;
  agenceForm.nom = agence.nom;
  agenceForm.code_agence = agence.code_agence || "";
  agenceForm.adresse = agence.adresse || "";
  agenceForm.ville = agence.ville || "";
  agenceForm.commune = agence.commune || "";
  agenceForm.telephone = agence.telephone || "";
  agenceForm.est_siege = agence.est_siege;
};

const cancelAgenceEdit = () => {
  isEditingAgence.value = false;
  editingAgenceId.value = null;
};

const saveAgence = async () => {
  if (!perms.value.canManageAgences) return;
  if (!agenceForm.nom.trim()) {
    guichetStore.showToast("Le nom de l'agence est obligatoire.", "error");
    return;
  }

  isSavingAgence.value = true;
  try {
    if (editingAgenceId.value) {
      await authStore.updateAgence(editingAgenceId.value, {
        nom: agenceForm.nom.trim(),
        code_agence: agenceForm.code_agence.trim(),
        adresse: agenceForm.adresse.trim(),
        ville: agenceForm.ville.trim(),
        commune: agenceForm.commune.trim(),
        telephone: agenceForm.telephone.trim(),
        est_siege: agenceForm.est_siege,
      });
      guichetStore.showToast("L'agence a été modifiée avec succès.", "success");
    } else {
      await authStore.addAgence({
        auto_ecole_id: authStore.currentSchool?.id || "school-primary",
        nom: agenceForm.nom.trim(),
        code_agence: agenceForm.code_agence.trim(),
        adresse: agenceForm.adresse.trim(),
        ville: agenceForm.ville.trim(),
        commune: agenceForm.commune.trim(),
        telephone: agenceForm.telephone.trim(),
        est_siege: agenceForm.est_siege,
      });
      guichetStore.showToast("Nouvelle agence ajoutée avec succès.", "success");
    }

    cancelAgenceEdit();
  } catch (err: any) {
    guichetStore.showToast(
      err.message || "Erreur lors de l'enregistrement de l'agence.",
      "error",
    );
  } finally {
    isSavingAgence.value = false;
  }
};

const removeAgence = async (agence: Agence) => {
  if (!perms.value.canManageAgences) return;
  if (agence.est_siege) {
    guichetStore.showToast(
      "Impossible de supprimer l'agence siège principale.",
      "error",
    );
    return;
  }
  if (confirm(`Voulez-vous vraiment supprimer l'agence "${agence.nom}" ?`)) {
    try {
      await authStore.deleteAgence(agence.id);
      guichetStore.showToast("L'agence a été supprimée.", "info");
    } catch (err: any) {
      guichetStore.showToast(
        err.message || "Erreur lors de la suppression de l'agence.",
        "error",
      );
    }
  }
};
</script>

<template>
  <div class="space-y-5">
    <!-- Barre d'action supérieure -->
    <div class="flex items-center justify-between">
      <div>
        <h3 class="text-sm font-bold text-foreground">
          Réseau d'Agences de l'Auto-École
        </h3>
        <p class="text-xs text-muted-foreground">
          {{
            perms.canManageAgences
              ? "Gérez votre agence siège et vos succursales géographiques"
              : "Consultation des agences et succursales de l'établissement"
          }}
        </p>
      </div>

      <button
        v-if="perms.canManageAgences && !isEditingAgence"
        type="button"
        @click="startAddAgence"
        class="flex items-center space-x-1.5 px-3.5 py-2 rounded-xl bg-brand hover:bg-brand-dark text-brand-dark font-bold text-xs uppercase tracking-wider transition"
      >
        <Plus class="w-4 h-4" />
        <span>Ajouter une agence</span>
      </button>
    </div>

    <!-- Formulaire Ajout / Édition Agence (Directeur uniquement) -->
    <div
      v-if="perms.canManageAgences && isEditingAgence"
      class="p-4 sm:p-5 rounded-[4px] bg-muted/40 border border-brand-border/60 space-y-4 animate-in fade-in duration-150"
    >
      <div class="flex items-center justify-between pb-3 border-b border-border">
        <h4
          class="text-xs font-bold uppercase tracking-wider text-brand-dark dark:text-brand flex items-center gap-1.5"
        >
          <MapPin class="w-4 h-4" />
          {{
            editingAgenceId ? "Modifier l'agence" : "Créer une nouvelle agence"
          }}
        </h4>
        <button
          type="button"
          @click="cancelAgenceEdit"
          class="text-xs text-muted-foreground hover:text-foreground"
        >
          Annuler
        </button>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
        <div>
          <label class="block text-xs font-semibold text-muted-foreground mb-1">
            Nom de l'agence *
          </label>
          <input
            v-model="agenceForm.nom"
            type="text"
            required
            class="w-full px-3 py-2 rounded-xl bg-card border border-border text-foreground text-sm focus:ring-2 focus:ring-brand/40 outline-none"
            placeholder="Ex : Agence Makepe ou Succursale Akwa"
          />
        </div>

        <div>
          <label class="block text-xs font-semibold text-muted-foreground mb-1">
            Code Agence (interne)
          </label>
          <input
            v-model="agenceForm.code_agence"
            type="text"
            class="w-full px-3 py-2 rounded-xl bg-card border border-border text-foreground text-sm font-mono focus:ring-2 focus:ring-brand/40 outline-none"
            placeholder="Ex : AG-DLA-01"
          />
        </div>

        <div>
          <label class="block text-xs font-semibold text-muted-foreground mb-1">
            Ville
          </label>
          <input
            v-model="agenceForm.ville"
            type="text"
            class="w-full px-3 py-2 rounded-xl bg-card border border-border text-foreground text-sm focus:ring-2 focus:ring-brand/40 outline-none"
            placeholder="Ex : Douala"
          />
        </div>

        <div>
          <label class="block text-xs font-semibold text-muted-foreground mb-1">
            Commune / Arrondissement
          </label>
          <input
            v-model="agenceForm.commune"
            type="text"
            class="w-full px-3 py-2 rounded-xl bg-card border border-border text-foreground text-sm focus:ring-2 focus:ring-brand/40 outline-none"
            placeholder="Ex : Douala 5ème"
          />
        </div>

        <div class="sm:col-span-2">
          <label class="block text-xs font-semibold text-muted-foreground mb-1">
            Adresse physique précise
          </label>
          <input
            v-model="agenceForm.adresse"
            type="text"
            class="w-full px-3 py-2 rounded-xl bg-card border border-border text-foreground text-sm focus:ring-2 focus:ring-brand/40 outline-none"
            placeholder="Ex : Rond-point poulin, face pharmacie"
          />
        </div>

        <div>
          <label class="block text-xs font-semibold text-muted-foreground mb-1">
            Téléphone direct de l'agence
          </label>
          <input
            v-model="agenceForm.telephone"
            type="tel"
            class="w-full px-3 py-2 rounded-xl bg-card border border-border text-foreground text-sm focus:ring-2 focus:ring-brand/40 outline-none"
            placeholder="Ex : 651727932"
          />
        </div>

        <div class="flex items-center space-x-2 pt-5">
          <input
            id="est_siege"
            v-model="agenceForm.est_siege"
            type="checkbox"
            class="rounded text-brand focus:ring-brand h-4 w-4"
          />
          <label
            for="est_siege"
            class="text-xs font-semibold text-foreground cursor-pointer"
          >
            Définir comme Siège Social Principal
          </label>
        </div>
      </div>

      <div class="flex justify-end space-x-2 pt-2">
        <button
          type="button"
          @click="cancelAgenceEdit"
          class="px-4 py-2 rounded-xl border border-border text-xs font-semibold hover:bg-muted transition"
        >
          Annuler
        </button>
        <button
          type="button"
          @click="saveAgence"
          :disabled="isSavingAgence"
          class="px-5 py-2 rounded-xl bg-brand hover:bg-brand-dark text-brand-dark font-bold text-xs uppercase tracking-wider transition disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {{
            isSavingAgence
              ? "Enregistrement en cours..."
              : editingAgenceId
                ? "Enregistrer les modifications"
                : "Créer l'agence"
          }}
        </button>
      </div>
    </div>

    <!-- Liste des Agences -->
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
      <div
        v-for="ag in authStore.agencesList"
        :key="ag.id"
        class="p-4 rounded-[4px] border transition flex flex-col justify-between"
        :class="
          ag.est_siege
            ? 'bg-brand-light/20 border-brand-border'
            : 'bg-card border-border hover:border-muted-foreground/30'
        "
      >
        <div>
          <div class="flex items-start justify-between gap-2 mb-2">
            <div class="font-bold text-sm text-foreground flex items-center gap-1.5">
              <MapPin class="w-4 h-4 text-brand" />
              <span>{{ ag.nom }}</span>
            </div>
            <span
              class="text-[10px] font-bold px-2 py-0.5 rounded-full uppercase"
              :class="
                ag.est_siege
                  ? 'bg-brand/20 text-brand-dark dark:text-brand'
                  : 'bg-muted text-muted-foreground'
              "
            >
              {{ ag.est_siege ? "Siège Social" : "Succursale" }}
            </span>
          </div>

          <div class="text-xs text-muted-foreground space-y-1 mt-2">
            <div v-if="ag.adresse" class="text-foreground">
              📍 {{ ag.adresse }}
            </div>
            <div v-if="ag.ville || ag.commune">
              🏙️ {{ ag.commune ? `${ag.commune}, ` : "" }}{{ ag.ville }}
            </div>
            <div v-if="ag.telephone" class="font-mono text-[11px]">
              📞 {{ ag.telephone }}
            </div>
            <div
              v-if="ag.code_agence"
              class="text-[10px] font-mono text-muted-foreground"
            >
              Code : {{ ag.code_agence }}
            </div>
          </div>
        </div>

        <div
          v-if="perms.canManageAgences"
          class="flex items-center justify-end space-x-2 pt-3 mt-3 border-t border-border/50"
        >
          <button
            type="button"
            @click="startEditAgence(ag)"
            class="flex items-center space-x-1 px-2.5 py-1 rounded-lg text-xs font-medium text-foreground hover:bg-muted transition"
          >
            <Edit2 class="w-3.5 h-3.5 text-muted-foreground" />
            <span>Modifier</span>
          </button>

          <button
            v-if="!ag.est_siege"
            type="button"
            @click="removeAgence(ag)"
            class="p-1 rounded-lg text-rose-500 hover:bg-rose-500/10 transition"
            title="Supprimer cette succursale"
          >
            <Trash2 class="w-3.5 h-3.5" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

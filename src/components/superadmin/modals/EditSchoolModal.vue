<script setup lang="ts">
import { reactive, watch } from "vue";
import type {
  SchoolItemSummary,
  UpdateSchoolAdminPayload,
} from "../../../types";
import { X, Building2, Save, Loader2 } from "lucide-vue-next";

const props = defineProps<{
  isOpen: boolean;
  school: SchoolItemSummary | null;
  isLoading: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "save", payload: UpdateSchoolAdminPayload): void;
}>();

const form = reactive<UpdateSchoolAdminPayload>({
  school_id: "",
  nom_commercial: "",
  code_unique: "",
  numero_agrement: "",
  telephone: "",
  email: "",
  ville: "Douala",
  pays: "Cameroun",
  devise: "XAF",
  est_actif: true,
  max_machines: 3,
  type_licence: "ANNUEL",
  modules_actifs: "[]",
});

watch(
  () => props.school,
  (s) => {
    if (s) {
      form.school_id = s.id;
      form.nom_commercial = s.nom_commercial;
      form.code_unique = s.code_unique;
      form.numero_agrement = s.numero_agrement;
      form.telephone = s.telephone;
      form.email = s.email || "";
      form.ville = s.ville || "Douala";
      form.pays = s.pays || "Cameroun";
      form.devise = s.devise || "XAF";
      form.est_actif = s.est_actif;
      form.max_machines = s.max_machines || 3;
      form.type_licence = s.type_licence || "ANNUEL";
      form.modules_actifs = s.modules_actifs || "[]";
    }
  },
  { immediate: true },
);

const handleSubmit = () => {
  if (!form.nom_commercial.trim() || !form.code_unique.trim()) {
    alert("Le nom commercial et le slug unique sont obligatoires.");
    return;
  }
  emit("save", { ...form });
};
</script>

<template>
  <div
    v-if="isOpen && school"
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm"
  >
    <div
      class="w-full max-w-xl bg-card border border-border rounded-[8px] shadow-2xl overflow-hidden animate-in fade-in zoom-in-95 duration-150"
    >
      <!-- En-tête -->
      <div
        class="h-12 border-b border-border bg-muted/50 px-4 flex items-center justify-between"
      >
        <div class="flex items-center gap-2 text-foreground font-semibold text-sm">
          <Building2 class="w-4 h-4 text-brand" />
          <span>Modifier l'Établissement · {{ school.nom_commercial }}</span>
        </div>
        <button
          type="button"
          @click="emit('close')"
          class="text-muted-foreground hover:text-foreground p-1 rounded-[4px] transition-colors cursor-pointer"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Corps du formulaire -->
      <form
        @submit.prevent="handleSubmit"
        class="p-5 space-y-4 text-xs font-sans"
      >
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="block text-muted-foreground mb-1 font-medium"
              >Nom Commercial</label
            >
            <input
              v-model="form.nom_commercial"
              type="text"
              required
              class="w-full h-8 px-2.5 rounded-[4px] bg-background border border-input text-foreground focus:outline-none focus:ring-1 focus:ring-ring focus:border-ring font-medium"
            />
          </div>

          <div>
            <label class="block text-muted-foreground mb-1 font-medium"
              >Code Unique (Slug URL)</label
            >
            <input
              v-model="form.code_unique"
              type="text"
              required
              class="w-full h-8 px-2.5 rounded-[4px] bg-background border border-input text-brand font-mono focus:outline-none focus:ring-1 focus:ring-ring focus:border-ring"
            />
          </div>
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="block text-muted-foreground mb-1 font-medium"
              >Numéro d'Agrément Officiel</label
            >
            <input
              v-model="form.numero_agrement"
              type="text"
              required
              class="w-full h-8 px-2.5 rounded-[4px] bg-background border border-input text-foreground focus:outline-none focus:ring-1 focus:ring-ring focus:border-ring font-mono"
            />
          </div>

          <div>
            <label class="block text-muted-foreground mb-1 font-medium"
              >Téléphone & WhatsApp</label
            >
            <input
              v-model="form.telephone"
              type="text"
              required
              class="w-full h-8 px-2.5 rounded-[4px] bg-background border border-input text-foreground focus:outline-none focus:ring-1 focus:ring-ring focus:border-ring font-mono"
            />
          </div>
        </div>

        <div class="grid grid-cols-3 gap-3">
          <div>
            <label class="block text-muted-foreground mb-1 font-medium"
              >Email Établissement</label
            >
            <input
              v-model="form.email"
              type="email"
              placeholder="contact@autoecole.com"
              class="w-full h-8 px-2.5 rounded-[4px] bg-background border border-input text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-ring focus:border-ring font-mono"
            />
          </div>

          <div>
            <label class="block text-muted-foreground mb-1 font-medium">Ville</label>
            <input
              v-model="form.ville"
              type="text"
              class="w-full h-8 px-2.5 rounded-[4px] bg-background border border-input text-foreground focus:outline-none focus:ring-1 focus:ring-ring focus:border-ring"
            />
          </div>

          <div>
            <label class="block text-muted-foreground mb-1 font-medium">Pays</label>
            <input
              v-model="form.pays"
              type="text"
              class="w-full h-8 px-2.5 rounded-[4px] bg-background border border-input text-foreground focus:outline-none focus:ring-1 focus:ring-ring focus:border-ring"
            />
          </div>
        </div>

        <div class="grid grid-cols-3 gap-3 pt-2 border-t border-border">
          <div>
            <label class="block text-muted-foreground mb-1 font-medium"
              >Devise Monétaire</label
            >
            <select
              v-model="form.devise"
              class="w-full h-8 px-2 rounded-[4px] bg-background border border-input text-foreground focus:outline-none focus:ring-1 focus:ring-ring focus:border-ring font-mono"
            >
              <option value="XAF">XAF (FCFA CEMAC)</option>
              <option value="XOF">XOF (FCFA UEMOA)</option>
              <option value="EUR">EUR (€ Euro)</option>
              <option value="USD">USD ($ Dollar)</option>
            </select>
          </div>

          <div>
            <label class="block text-muted-foreground mb-1 font-medium"
              >Nombre Postes / Machines</label
            >
            <input
              v-model.number="form.max_machines"
              type="number"
              min="1"
              max="100"
              class="w-full h-8 px-2.5 rounded-[4px] bg-background border border-input text-brand font-mono font-bold focus:outline-none focus:ring-1 focus:ring-ring focus:border-ring"
            />
          </div>

          <div>
            <label class="block text-muted-foreground mb-1 font-medium"
              >Statut d'Accès</label
            >
            <select
              v-model="form.est_actif"
              class="w-full h-8 px-2 rounded-[4px] bg-background border border-input font-semibold focus:outline-none focus:ring-1 focus:ring-ring focus:border-ring"
              :class="form.est_actif ? 'text-brand' : 'text-rose-600 dark:text-rose-400'"
            >
              <option :value="true">ACTIF</option>
              <option :value="false">SUSPENDU</option>
            </select>
          </div>
        </div>

        <!-- Boutons d'action -->
        <div
          class="pt-4 border-t border-border flex items-center justify-end gap-2.5"
        >
          <button
            type="button"
            @click="emit('close')"
            class="px-3 py-1.5 rounded-[4px] bg-muted hover:bg-muted/80 text-foreground border border-border font-medium transition cursor-pointer"
          >
            Annuler
          </button>
          <button
            type="submit"
            :disabled="isLoading"
            class="px-4 py-1.5 rounded-[4px] bg-brand-600 hover:bg-brand-700 text-white font-medium flex items-center gap-1.5 transition cursor-pointer disabled:opacity-50 shadow-xs"
          >
            <Loader2
              v-if="isLoading"
              class="w-3.5 h-3.5 animate-spin text-white"
            />
            <Save v-else class="w-3.5 h-3.5" />
            <span>Enregistrer les modifications</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

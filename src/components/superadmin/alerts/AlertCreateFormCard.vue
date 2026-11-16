<script setup lang="ts">
import { ref } from "vue";
import type { SchoolItemSummary } from "../../../types";
import {
  AlertTriangle,
  Info,
  ShieldAlert,
  Send,
  CheckCircle2,
  Radio,
  Loader2,
  AlertCircle,
} from "lucide-vue-next";

defineProps<{
  schools: SchoolItemSummary[];
  isSubmitting: boolean;
  isFormSuccess: boolean;
  errorMessage: string | null;
}>();

export interface CreateAlertPayload {
  titre: string;
  message: string;
  severite: "info" | "warning" | "critical";
  cibleType: "all" | "single";
  cibleSchoolId: string;
  bloquante: boolean;
}

const emit = defineEmits<{
  (e: "submit", payload: CreateAlertPayload): void;
}>();

const form = ref<CreateAlertPayload>({
  titre: "",
  message: "",
  severite: "warning",
  cibleType: "all",
  cibleSchoolId: "",
  bloquante: false,
});

const handleSubmit = () => {
  if (!form.value.titre.trim() || !form.value.message.trim()) return;
  emit("submit", { ...form.value });
};

const resetFields = () => {
  form.value.titre = "";
  form.value.message = "";
  form.value.bloquante = false;
};

defineExpose({ resetFields });
</script>

<template>
  <div
    class="lg:col-span-5 border border-border rounded-[3px] bg-card p-5 space-y-4"
  >
    <div class="flex items-center justify-between border-b border-border pb-3">
      <span
        class="font-bold text-xs text-foreground uppercase tracking-wider flex items-center gap-1.5"
      >
        <Radio class="w-4 h-4 text-brand" />
        <span>Nouvelle Alerte</span>
      </span>
      <span
        v-if="isFormSuccess"
        class="text-[11px] text-emerald-600 dark:text-emerald-400 font-semibold flex items-center gap-1 animate-in fade-in"
      >
        <CheckCircle2 class="w-3.5 h-3.5" />
        <span>Diffusée avec succès !</span>
      </span>
    </div>

    <form @submit.prevent="handleSubmit" class="space-y-3.5 text-xs">
      <!-- Destinataires -->
      <div>
        <label class="block text-[11px] font-semibold text-muted-foreground mb-1">
          Cible de l'alerte
        </label>
        <div class="grid grid-cols-2 gap-2 mb-2">
          <button
            type="button"
            @click="form.cibleType = 'all'"
            class="px-2.5 py-1.5 rounded-[4px] text-xs font-medium border text-center transition cursor-pointer"
            :class="
              form.cibleType === 'all'
                ? 'bg-brand text-white border-brand font-semibold'
                : 'bg-muted text-muted-foreground border-border hover:text-foreground'
            "
          >
            Toutes les auto-écoles
          </button>
          <button
            type="button"
            @click="form.cibleType = 'single'"
            class="px-2.5 py-1.5 rounded-[4px] text-xs font-medium border text-center transition cursor-pointer"
            :class="
              form.cibleType === 'single'
                ? 'bg-brand text-white border-brand font-semibold'
                : 'bg-muted text-muted-foreground border-border hover:text-foreground'
            "
          >
            Auto-école spécifique
          </button>
        </div>

        <select
          v-if="form.cibleType === 'single'"
          v-model="form.cibleSchoolId"
          required
          class="w-full p-2 text-xs rounded-[4px] border border-border bg-background focus:outline-hidden focus:ring-1 focus:ring-brand animate-in fade-in"
        >
          <option value="" disabled>Sélectionner un établissement...</option>
          <option
            v-for="school in schools"
            :key="school.id"
            :value="school.id"
          >
            {{ school.nom_commercial || (school as any).nom }} ({{
              school.ville || "Actif"
            }})
          </option>
        </select>
      </div>

      <!-- Sévérité -->
      <div>
        <label class="block text-[11px] font-semibold text-muted-foreground mb-1">
          Niveau de sévérité
        </label>
        <div class="grid grid-cols-3 gap-2">
          <button
            type="button"
            @click="form.severite = 'info'"
            class="px-2 py-1.5 rounded-[4px] text-[11px] border text-center transition cursor-pointer flex items-center justify-center gap-1"
            :class="
              form.severite === 'info'
                ? 'bg-blue-500/20 text-blue-700 dark:text-blue-300 border-blue-500/50 font-bold'
                : 'bg-muted text-muted-foreground border-border hover:text-foreground'
            "
          >
            <Info class="w-3 h-3 text-blue-500" />
            <span>Info</span>
          </button>
          <button
            type="button"
            @click="form.severite = 'warning'"
            class="px-2 py-1.5 rounded-[4px] text-[11px] border text-center transition cursor-pointer flex items-center justify-center gap-1"
            :class="
              form.severite === 'warning'
                ? 'bg-amber-500/20 text-amber-700 dark:text-amber-300 border-amber-500/50 font-bold'
                : 'bg-muted text-muted-foreground border-border hover:text-foreground'
            "
          >
            <AlertTriangle class="w-3 h-3 text-amber-500" />
            <span>Avertissement</span>
          </button>
          <button
            type="button"
            @click="form.severite = 'critical'"
            class="px-2 py-1.5 rounded-[4px] text-[11px] border text-center transition cursor-pointer flex items-center justify-center gap-1"
            :class="
              form.severite === 'critical'
                ? 'bg-rose-500/20 text-rose-700 dark:text-rose-300 border-rose-500/50 font-bold'
                : 'bg-muted text-muted-foreground border-border hover:text-foreground'
            "
          >
            <ShieldAlert class="w-3 h-3 text-rose-500" />
            <span>Critique</span>
          </button>
        </div>
      </div>

      <!-- Titre -->
      <div>
        <label class="block text-[11px] font-semibold text-muted-foreground mb-1">
          Titre de l'alerte
        </label>
        <input
          v-model="form.titre"
          type="text"
          required
          placeholder="Ex: Maintenance système / Suspension temporaire..."
          class="w-full px-3 py-1.5 text-xs rounded-[4px] border border-border bg-background focus:outline-hidden focus:ring-1 focus:ring-brand"
        />
      </div>

      <!-- Message -->
      <div>
        <label class="block text-[11px] font-semibold text-muted-foreground mb-1">
          Message détaillé
        </label>
        <textarea
          v-model="form.message"
          rows="3"
          required
          placeholder="Saisissez le texte de l'alerte diffusée sur les écrans du personnel..."
          class="w-full p-2.5 text-xs rounded-[4px] border border-border bg-background focus:outline-hidden focus:ring-1 focus:ring-brand"
        ></textarea>
      </div>

      <div class="flex items-center space-x-2 pt-1">
        <input
          id="bloquante-check"
          v-model="form.bloquante"
          type="checkbox"
          class="rounded border-border text-brand focus:ring-brand h-4 w-4"
        />
        <label
          for="bloquante-check"
          class="text-xs text-foreground cursor-pointer select-none"
        >
          Affichage modal bloquant (impose une confirmation de lecture)
        </label>
      </div>

      <div
        v-if="errorMessage"
        class="p-2.5 rounded-[4px] bg-rose-500/10 border border-rose-500/20 text-rose-600 dark:text-rose-400 text-xs flex items-center gap-2"
      >
        <AlertCircle class="w-4 h-4 shrink-0" />
        <span>{{ errorMessage }}</span>
      </div>

      <div class="pt-2">
        <button
          type="submit"
          :disabled="isSubmitting"
          class="w-full py-2 px-4 rounded-[4px] bg-brand hover:bg-brand-dark text-white font-semibold text-xs transition flex items-center justify-center gap-2 cursor-pointer active:scale-[0.99] disabled:opacity-50"
        >
          <Loader2 v-if="isSubmitting" class="w-3.5 h-3.5 animate-spin" />
          <Send v-else class="w-3.5 h-3.5" />
          <span>{{
            isSubmitting
              ? "Diffusion en cours..."
              : "Diffuser l'alerte sur le réseau"
          }}</span>
        </button>
      </div>
    </form>
  </div>
</template>


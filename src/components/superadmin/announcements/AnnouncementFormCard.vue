<script setup lang="ts">
import { ref } from "vue";
import type { SchoolItemSummary } from "../../../types";
import {
  Image as ImageIcon,
  FileText,
  Send,
  CheckCircle2,
  Upload,
  Sparkles,
  Loader2,
  AlertCircle,
} from "lucide-vue-next";

defineProps<{
  schools: SchoolItemSummary[];
  isSubmitting: boolean;
  isFormSuccess: boolean;
  errorMessage: string | null;
}>();

export interface CreateAnnouncementFormPayload {
  titre: string;
  contenu: string;
  typeFormat: "text_only" | "text_image" | "image_only";
  imageUrl: string;
  categorie: "evenement" | "mise_a_jour" | "formation" | "general";
  cibleType: "all" | "single";
  cibleSchoolId: string;
  dateEvenement: string;
  lienAction: string;
}

const emit = defineEmits<{
  (e: "submit", payload: CreateAnnouncementFormPayload): void;
}>();

const fileInputRef = ref<HTMLInputElement | null>(null);
const form = ref<CreateAnnouncementFormPayload>({
  titre: "",
  contenu: "",
  typeFormat: "text_image",
  imageUrl: "",
  categorie: "evenement",
  cibleType: "all",
  cibleSchoolId: "",
  dateEvenement: "",
  lienAction: "",
});

const handleImageUpload = (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0];
  if (!file) return;

  const reader = new FileReader();
  reader.onload = (e) => {
    form.value.imageUrl = e.target?.result as string;
  };
  reader.readAsDataURL(file);
};

const handleSubmit = () => {
  if (form.value.typeFormat !== "image_only" && !form.value.titre.trim()) return;
  emit("submit", { ...form.value });
};

const resetFields = () => {
  form.value.titre = "";
  form.value.contenu = "";
  form.value.imageUrl = "";
  form.value.dateEvenement = "";
  form.value.lienAction = "";
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
        <Sparkles class="w-4 h-4 text-brand" />
        <span>Créer une Annonce</span>
      </span>
      <span
        v-if="isFormSuccess"
        class="text-[11px] text-emerald-600 dark:text-emerald-400 font-semibold flex items-center gap-1 animate-in fade-in"
      >
        <CheckCircle2 class="w-3.5 h-3.5" />
        <span>Publiée avec succès !</span>
      </span>
    </div>

    <form @submit.prevent="handleSubmit" class="space-y-3.5 text-xs">
      <!-- Format -->
      <div>
        <label class="block text-[11px] font-semibold text-muted-foreground mb-1">
          Format de l'annonce
        </label>
        <div class="grid grid-cols-3 gap-1.5">
          <button
            type="button"
            @click="form.typeFormat = 'text_only'"
            class="px-2 py-1.5 rounded-[4px] text-[11px] font-medium border text-center transition cursor-pointer flex items-center justify-center gap-1"
            :class="
              form.typeFormat === 'text_only'
                ? 'bg-brand text-white border-brand font-bold'
                : 'bg-muted text-muted-foreground border-border hover:text-foreground'
            "
          >
            <FileText class="w-3 h-3" />
            <span>Texte seul</span>
          </button>
          <button
            type="button"
            @click="form.typeFormat = 'text_image'"
            class="px-2 py-1.5 rounded-[4px] text-[11px] font-medium border text-center transition cursor-pointer flex items-center justify-center gap-1"
            :class="
              form.typeFormat === 'text_image'
                ? 'bg-brand text-white border-brand font-bold'
                : 'bg-muted text-muted-foreground border-border hover:text-foreground'
            "
          >
            <ImageIcon class="w-3 h-3" />
            <span>Texte + Image</span>
          </button>
          <button
            type="button"
            @click="form.typeFormat = 'image_only'"
            class="px-2 py-1.5 rounded-[4px] text-[11px] font-medium border text-center transition cursor-pointer flex items-center justify-center gap-1"
            :class="
              form.typeFormat === 'image_only'
                ? 'bg-brand text-white border-brand font-bold'
                : 'bg-muted text-muted-foreground border-border hover:text-foreground'
            "
          >
            <ImageIcon class="w-3 h-3" />
            <span>Image seule</span>
          </button>
        </div>
      </div>

      <!-- Audience -->
      <div>
        <label class="block text-[11px] font-semibold text-muted-foreground mb-1">
          Audience visée
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
            Établissement précis
          </button>
        </div>

        <select
          v-if="form.cibleType === 'single'"
          v-model="form.cibleSchoolId"
          required
          class="w-full p-2 text-xs rounded-[4px] border border-border bg-background"
        >
          <option value="" disabled>Sélectionner l'auto-école ciblée...</option>
          <option v-for="s in schools" :key="s.id" :value="s.id">
            {{ s.nom_commercial || (s as any).nom }} ({{ s.ville || "Actif" }})
          </option>
        </select>
      </div>

      <!-- Titre -->
      <div v-if="form.typeFormat !== 'image_only'">
        <label class="block text-[11px] font-semibold text-muted-foreground mb-1">
          Titre de l'annonce ou de l'événement
        </label>
        <input
          v-model="form.titre"
          type="text"
          required
          placeholder="Ex: Stage de recyclage moniteurs / Webinaire..."
          class="w-full px-3 py-1.5 text-xs rounded-[4px] border border-border bg-background focus:outline-hidden focus:ring-1 focus:ring-brand"
        />
      </div>

      <!-- Contenu texte -->
      <div v-if="form.typeFormat !== 'image_only'">
        <label class="block text-[11px] font-semibold text-muted-foreground mb-1">
          Message / Description
        </label>
        <textarea
          v-model="form.contenu"
          rows="3"
          required
          placeholder="Détaillez les informations de l'annonce ou de la manifestation..."
          class="w-full p-2.5 text-xs rounded-[4px] border border-border bg-background focus:outline-hidden focus:ring-1 focus:ring-brand"
        ></textarea>
      </div>

      <!-- Image attachée -->
      <div
        v-if="form.typeFormat === 'text_image' || form.typeFormat === 'image_only'"
        class="space-y-2"
      >
        <label class="block text-[11px] font-semibold text-muted-foreground">
          Visuel de l'annonce (Image / Affiche)
        </label>

        <div class="flex items-center gap-2">
          <input
            v-model="form.imageUrl"
            type="url"
            placeholder="URL de l'image (https://...)"
            class="flex-1 px-3 py-1.5 text-xs rounded-[4px] border border-border bg-background"
          />
          <button
            type="button"
            @click="fileInputRef?.click()"
            class="px-2.5 py-1.5 rounded-[4px] border border-border bg-muted hover:bg-muted/80 text-foreground flex items-center gap-1.5 text-xs cursor-pointer"
            title="Téléverser depuis l'ordinateur"
          >
            <Upload class="w-3.5 h-3.5" />
            <span>Fichier</span>
          </button>
          <input
            ref="fileInputRef"
            type="file"
            accept="image/*"
            class="hidden"
            @change="handleImageUpload"
          />
        </div>

        <div
          v-if="form.imageUrl"
          class="relative rounded-[3px] overflow-hidden border border-border max-h-36 bg-black/5"
        >
          <img
            :src="form.imageUrl"
            alt="Aperçu annonce"
            class="w-full h-36 object-cover"
          />
          <button
            type="button"
            @click="form.imageUrl = ''"
            class="absolute top-2 right-2 px-2 py-0.5 rounded bg-black/70 text-white text-[10px] hover:bg-black"
          >
            Retirer
          </button>
        </div>
      </div>

      <!-- Date & Lien optionnels -->
      <div class="grid grid-cols-2 gap-2 pt-1">
        <div>
          <label
            class="block text-[10px] font-semibold text-muted-foreground mb-1"
          >
            Date de l'événement (Optionnel)
          </label>
          <input
            v-model="form.dateEvenement"
            type="date"
            class="w-full px-2.5 py-1 text-xs rounded-[4px] border border-border bg-background"
          />
        </div>
        <div>
          <label
            class="block text-[10px] font-semibold text-muted-foreground mb-1"
          >
            Lien web d'action (Optionnel)
          </label>
          <input
            v-model="form.lienAction"
            type="url"
            placeholder="https://..."
            class="w-full px-2.5 py-1 text-xs rounded-[4px] border border-border bg-background"
          />
        </div>
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
              ? "Publication en cours..."
              : "Diffuser l'annonce maintenant"
          }}</span>
        </button>
      </div>
    </form>
  </div>
</template>


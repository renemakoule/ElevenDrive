<script setup lang="ts">
import { ref, watch } from "vue";
import type { Candidat, TypeDocumentImportant } from "../../types";
import { UploadCloud, X } from "lucide-vue-next";

interface Props {
  isOpen: boolean;
  candidates: Candidat[];
  initialCandidatId?: string;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "close"): void;
  (
    e: "submit",
    payload: {
      candidatId: string;
      type: TypeDocumentImportant;
      titre: string;
      nomFichier: string;
      base64: string;
      taille: number;
      notes: string;
    },
  ): void;
}>();

const uploadCandidatId = ref<string>("");
const uploadDocType = ref<TypeDocumentImportant>("PIECE_IDENTITE");
const uploadDocTitre = ref<string>("");
const uploadDocFile = ref<{ nom: string; base64: string; size: number } | null>(
  null,
);
const uploadDocNotes = ref<string>("");

watch(
  () => props.isOpen,
  (open) => {
    if (open) {
      uploadCandidatId.value =
        props.initialCandidatId || props.candidates[0]?.id || "";
      uploadDocType.value = "PIECE_IDENTITE";
      uploadDocTitre.value = "";
      uploadDocFile.value = null;
      uploadDocNotes.value = "";
    }
  },
  { immediate: true },
);

const handleDocFileInput = (e: Event) => {
  const target = e.target as HTMLInputElement;
  if (!target.files || target.files.length === 0) return;
  const file = target.files[0];
  const reader = new FileReader();
  reader.onload = () => {
    uploadDocFile.value = {
      nom: file.name,
      base64: reader.result as string,
      size: file.size,
    };
    if (!uploadDocTitre.value) {
      uploadDocTitre.value = file.name.replace(/\.[^/.]+$/, "");
    }
  };
  reader.readAsDataURL(file);
};

const onSubmit = () => {
  if (!uploadCandidatId.value || !uploadDocFile.value) return;
  emit("submit", {
    candidatId: uploadCandidatId.value,
    type: uploadDocType.value,
    titre: uploadDocTitre.value || uploadDocFile.value.nom,
    nomFichier: uploadDocFile.value.nom,
    base64: uploadDocFile.value.base64,
    taille: uploadDocFile.value.size,
    notes: uploadDocNotes.value,
  });
};
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 backdrop-blur-xs"
    role="dialog"
    aria-modal="true"
  >
    <div
      class="w-full max-w-md rounded-[4px] border border-border bg-card p-5 space-y-4"
    >
      <div
        class="flex items-center justify-between border-b border-border pb-3"
      >
        <div class="flex items-center gap-2">
          <UploadCloud class="h-5 w-5 text-brand" />
          <h2 class="text-base font-bold text-foreground">
            Téléverser un Document Important
          </h2>
        </div>
        <button
          type="button"
          @click="emit('close')"
          class="text-muted-foreground hover:text-foreground"
        >
          <X class="h-4 w-4" />
        </button>
      </div>

      <div class="space-y-3 text-xs">
        <div>
          <label class="block font-medium text-foreground mb-1">
            Apprenant Concerné
          </label>
          <select
            v-model="uploadCandidatId"
            class="w-full rounded-[4px] border border-border bg-background p-2 text-foreground focus:outline-none"
          >
            <option v-for="cand in candidates" :key="cand.id" :value="cand.id">
              {{ cand.nom }} {{ cand.prenoms }} (N°
              {{ cand.numero_dossier || "DOS-ATTENTE" }})
            </option>
          </select>
        </div>

        <div>
          <label class="block font-medium text-foreground mb-1">
            Type de Document
          </label>
          <select
            v-model="uploadDocType"
            class="w-full rounded-[4px] border border-border bg-background p-2 text-foreground focus:outline-none"
          >
            <option value="PIECE_IDENTITE">
              Pièce d'Identité (CNI / Passeport / Titre)
            </option>
            <option value="VISITE_MEDICALE">
              Certificat d'Aptitude Médicale
            </option>
            <option value="PHOTO_IDENTITE">Photo d'Identité 4x4 Récente</option>
            <option value="CONTRAT_FORMATION">
              Contrat de Formation d'Auto-École
            </option>
            <option value="ATTESTATION_EXAMEN">
              Attestation de Réussite / Permis Provisoire
            </option>
            <option value="AUTRE">Autre Document Justificatif</option>
          </select>
        </div>

        <div>
          <label class="block font-medium text-foreground mb-1">
            Titre / Description du document
          </label>
          <input
            v-model="uploadDocTitre"
            type="text"
            placeholder="Ex: CNI Recto-Verso KOUAME Marc"
            class="w-full rounded-[4px] border border-border bg-background p-2 text-foreground focus:outline-none"
          />
        </div>

        <div>
          <label class="block font-medium text-foreground mb-1">
            Fichier (PDF ou Image)
          </label>
          <input
            type="file"
            @change="handleDocFileInput"
            accept="image/*,.pdf"
            class="w-full text-xs text-muted-foreground file:mr-3 file:py-1.5 file:px-3 file:rounded-md file:border file:border-border file:bg-muted file:text-xs file:font-semibold hover:file:bg-brand hover:file:text-brand-dark"
          />
        </div>

        <div>
          <label class="block font-medium text-foreground mb-1">
            Notes administratives (Optionnel)
          </label>
          <input
            v-model="uploadDocNotes"
            type="text"
            placeholder="Validité, conformité, observations..."
            class="w-full rounded-[4px] border border-border bg-background p-2 text-foreground focus:outline-none"
          />
        </div>
      </div>

      <div
        class="flex items-center justify-end gap-2 border-t border-border pt-3"
      >
        <button
          type="button"
          @click="emit('close')"
          class="px-3 py-1.5 rounded-[4px] border border-border bg-card text-xs font-medium text-foreground hover:bg-muted"
        >
          Annuler
        </button>

        <button
          type="button"
          @click="onSubmit"
          class="px-4 py-1.5 rounded-[4px] bg-brand text-brand-dark font-bold text-xs hover:bg-brand-dark hover:text-white transition"
        >
          Téléverser et Sécuriser
        </button>
      </div>
    </div>
  </div>
</template>

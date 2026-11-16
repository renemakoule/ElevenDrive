<script setup lang="ts">
import { ref, watch } from "vue";
import type { Candidat, MotifArchivage, TypeDocumentImportant } from "../../types";
import { FolderArchive, UploadCloud, X } from "lucide-vue-next";

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
      motif: MotifArchivage;
      referenceBoite: string;
      observations: string;
      documentsFinaux: {
        type: TypeDocumentImportant;
        titre: string;
        nomFichier: string;
        base64: string;
      }[];
    },
  ): void;
}>();

const selectedCandidatId = ref<string>("");
const formMotif = ref<MotifArchivage>("PERMIS_OBTENU");
const formReferenceBoite = ref<string>("");
const formObservations = ref<string>("");
const attachedDocs = ref<
  {
    type: TypeDocumentImportant;
    titre: string;
    nomFichier: string;
    base64: string;
  }[]
>([]);

watch(
  () => props.isOpen,
  (open) => {
    if (open) {
      if (props.initialCandidatId) {
        selectedCandidatId.value = props.initialCandidatId;
        const cand = props.candidates.find((c) => c.id === props.initialCandidatId);
        if (cand) {
          formReferenceBoite.value = `BOITE-${new Date().getFullYear()}-${cand.numero_dossier || "001"}`;
        }
      } else if (props.candidates.length > 0) {
        selectedCandidatId.value = props.candidates[0].id;
        formReferenceBoite.value = `BOITE-${new Date().getFullYear()}-${props.candidates[0].numero_dossier || "001"}`;
      }
      attachedDocs.value = [];
      formObservations.value = "";
      formMotif.value = "PERMIS_OBTENU";
    }
  },
  { immediate: true },
);

const handleFileAttach = (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (!target.files || target.files.length === 0) return;
  const file = target.files[0];
  const reader = new FileReader();
  reader.onload = () => {
    attachedDocs.value.push({
      type: "AUTRE",
      titre: file.name.replace(/\.[^/.]+$/, ""),
      nomFichier: file.name,
      base64: reader.result as string,
    });
  };
  reader.readAsDataURL(file);
};

const removeAttachedDoc = (idx: number) => {
  attachedDocs.value.splice(idx, 1);
};

const onSubmit = () => {
  if (!selectedCandidatId.value) return;
  emit("submit", {
    candidatId: selectedCandidatId.value,
    motif: formMotif.value,
    referenceBoite: formReferenceBoite.value,
    observations: formObservations.value,
    documentsFinaux: attachedDocs.value,
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
      class="w-full max-w-lg rounded-[4px] border border-border bg-card p-5 space-y-4 max-h-[90vh] overflow-y-auto"
    >
      <div
        class="flex items-center justify-between border-b border-border pb-3"
      >
        <div class="flex items-center gap-2">
          <FolderArchive class="h-5 w-5 text-brand" />
          <h2 class="text-base font-bold text-foreground">
            Archiver un Dossier Candidat
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
        <!-- Sélection du candidat -->
        <div>
          <label class="block font-medium text-foreground mb-1">
            Candidat & Dossier
          </label>
          <select
            v-model="selectedCandidatId"
            class="w-full rounded-[4px] border border-border bg-background p-2 text-foreground focus:outline-none focus:ring-1 focus:ring-brand"
          >
            <option
              v-for="cand in candidates"
              :key="cand.id"
              :value="cand.id"
            >
              {{ cand.nom }} {{ cand.prenoms }} (N°
              {{ cand.numero_dossier || "DOS-ATTENTE" }} •
              {{ cand.formule_nom || "Permis" }})
            </option>
          </select>
        </div>

        <!-- Motif d'archivage -->
        <div>
          <label class="block font-medium text-foreground mb-1">
            Motif Officiel d'Archivage
          </label>
          <select
            v-model="formMotif"
            class="w-full rounded-[4px] border border-border bg-background p-2 text-foreground focus:outline-none focus:ring-1 focus:ring-brand"
          >
            <option value="PERMIS_OBTENU">
              Permis Obtenu avec Succès (Admis)
            </option>
            <option value="DOSSIER_SOLDE">
              Formation Achevée & Solde 100% Réglé
            </option>
            <option value="ABANDON_DESISTEMENT">
              Abandon ou Désistement du Candidat
            </option>
            <option value="PRESCRIPTION_LEGALE">
              Prescription Légale (Délai de Validité Expiré)
            </option>
            <option value="TRANSFERT_AUTRE_ECOLE">
              Transfert vers une Autre Auto-École
            </option>
            <option value="AUTRE">Autre Motif Administratif</option>
          </select>
        </div>

        <!-- Référence de boîte / rayon physique -->
        <div>
          <label class="block font-medium text-foreground mb-1">
            Référence Boîte / Emplacement Physique
          </label>
          <input
            v-model="formReferenceBoite"
            type="text"
            placeholder="Ex: BOITE-2026-A12, RAYON-3"
            class="w-full rounded-[4px] border border-border bg-background p-2 text-foreground focus:outline-none focus:ring-1 focus:ring-brand"
          />
        </div>

        <!-- Observations -->
        <div>
          <label class="block font-medium text-foreground mb-1">
            Observations / Bilan de Clôture
          </label>
          <textarea
            v-model="formObservations"
            rows="2"
            placeholder="Remarques, date d'obtention de l'attestation, justificatifs annexés..."
            class="w-full rounded-[4px] border border-border bg-background p-2 text-foreground focus:outline-none focus:ring-1 focus:ring-brand"
          ></textarea>
        </div>

        <!-- Pièces justificatives annexes -->
        <div>
          <label class="block font-medium text-foreground mb-1">
            Joindre Justificatif (Scan Permis, PV Examen)
          </label>
          <div class="flex items-center gap-2">
            <label
              class="cursor-pointer inline-flex items-center gap-1.5 px-3 py-1.5 rounded-md border border-border bg-muted/40 hover:bg-muted text-foreground transition text-xs font-medium"
            >
              <UploadCloud class="h-3.5 w-3.5" />
              <span>Sélectionner Fichier</span>
              <input
                type="file"
                @change="handleFileAttach"
                class="hidden"
                accept="image/*,.pdf"
              />
            </label>
            <span class="text-muted-foreground text-[11px]">
              {{ attachedDocs.length }} fichier(s) joint(s)
            </span>
          </div>

          <div v-if="attachedDocs.length > 0" class="mt-2 space-y-1">
            <div
              v-for="(doc, idx) in attachedDocs"
              :key="idx"
              class="flex items-center justify-between p-2 rounded bg-muted/30 border border-border text-[11px]"
            >
              <span class="truncate max-w-[280px]">{{ doc.nomFichier }}</span>
              <button
                type="button"
                @click="removeAttachedDoc(idx)"
                class="text-rose-500 hover:text-rose-700 font-bold"
              >
                <X class="h-3.5 w-3.5" />
              </button>
            </div>
          </div>
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
          Confirmer & Sceller l'Archive
        </button>
      </div>
    </div>
  </div>
</template>


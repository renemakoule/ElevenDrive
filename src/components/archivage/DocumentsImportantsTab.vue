<script setup lang="ts">
import type { DocumentArchive, TypeDocumentImportant } from "../../types";
import {
  FileText,
  UploadCloud,
  FileCheck,
  ExternalLink,
  Download,
} from "lucide-vue-next";

interface Props {
  documents: DocumentArchive[];
}

defineProps<Props>();

const emit = defineEmits<{
  (e: "open-upload"): void;
}>();

const docTypeLabels: Record<TypeDocumentImportant, string> = {
  PIECE_IDENTITE: "Pièce d'Identité (CNI / Passeport)",
  VISITE_MEDICALE: "Certificat de Visite Médicale",
  PHOTO_IDENTITE: "Photo d'Identité 4x4",
  CONTRAT_FORMATION: "Contrat de Formation Signé",
  RECU_PAIEMENT: "Reçu & Bilan de Règlement",
  ATTESTATION_EXAMEN: "Attestation / Résultat d'Examen",
  AUTRE: "Autre Document Justificatif",
};
</script>

<template>
  <div class="space-y-4">
    <div
      class="flex items-center justify-between bg-card border border-border p-3 rounded-[4px]"
    >
      <div>
        <h3 class="text-xs sm:text-sm font-semibold text-foreground">
          Stockage des Documents Importants & Pièces Justificatives
        </h3>
        <p class="text-xs text-muted-foreground">
          Numérisation des CNI, certificats médicaux et pièces requises pour l'examen.
        </p>
      </div>

      <button
        type="button"
        @click="emit('open-upload')"
        class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-[4px] bg-brand text-brand-dark font-semibold text-xs hover:bg-brand-dark hover:text-white transition"
      >
        <UploadCloud class="h-4 w-4" />
        Ajouter un Document
      </button>
    </div>

    <div class="rounded-[4px] border border-border bg-card overflow-hidden">
      <div
        v-if="documents.length === 0"
        class="p-12 text-center"
      >
        <FileText class="h-10 w-10 text-muted-foreground/40 mx-auto mb-3" />
        <p class="text-sm font-medium text-foreground">
          Aucun document téléversé actuellement
        </p>
        <p class="text-xs text-muted-foreground mt-1">
          Cliquez sur « Ajouter un Document » pour téléverser une pièce
          d'identité ou un certificat médical.
        </p>
      </div>

      <div v-else class="divide-y divide-border">
        <div
          v-for="doc in documents"
          :key="doc.id"
          class="p-4 transition hover:bg-muted/30 flex items-center justify-between gap-4"
        >
          <div class="flex items-center gap-3.5 min-w-0">
            <div
              class="h-9 w-9 rounded-[4px] shrink-0 flex items-center justify-center font-bold text-white bg-emerald-600"
            >
              <FileCheck class="h-5 w-5" />
            </div>

            <div class="min-w-0">
              <div class="flex items-center gap-2 mb-0.5 flex-wrap">
                <h4
                  class="text-xs sm:text-sm font-semibold text-foreground truncate"
                >
                  {{ doc.titre }}
                </h4>
                <span
                  class="text-[10px] px-2 py-0.5 rounded-full border border-border font-medium bg-muted text-muted-foreground"
                >
                  {{ docTypeLabels[doc.type_document] || doc.type_document }}
                </span>
                <span
                  v-if="doc.est_archive"
                  class="text-[10px] px-2 py-0.5 rounded-full bg-blue-500/10 text-blue-700 dark:text-blue-300 border border-blue-500/20 font-semibold"
                >
                  Scellé dans l'Archive
                </span>
              </div>

              <p class="text-xs text-muted-foreground">
                Fichier :
                <strong class="text-foreground">{{ doc.nom_fichier }}</strong>
                • Téléversé le
                {{ new Date(doc.date_televersement).toLocaleDateString("fr-FR") }}
                par {{ doc.televerse_par_nom }}
              </p>
            </div>
          </div>

          <div class="flex items-center gap-2 shrink-0">
            <a
              v-if="doc.fichier_url_ou_base64 && doc.fichier_url_ou_base64.startsWith('http')"
              :href="doc.fichier_url_ou_base64"
              target="_blank"
              rel="noopener noreferrer"
              class="inline-flex items-center gap-1.5 px-2.5 py-1.5 rounded-[4px] border border-border bg-card text-xs font-medium text-sky-600 hover:bg-sky-500/10 transition"
              title="Consulter directement sur Cloudflare R2"
            >
              <ExternalLink class="h-3.5 w-3.5" />
              Voir sur R2
            </a>

            <a
              v-if="doc.fichier_url_ou_base64"
              :href="doc.fichier_url_ou_base64"
              :download="doc.nom_fichier"
              class="inline-flex items-center gap-1.5 px-2.5 py-1.5 rounded-[4px] border border-border bg-card text-xs font-medium text-foreground hover:bg-muted transition"
              title="Télécharger le fichier"
            >
              <Download class="h-3.5 w-3.5" />
              Télécharger
            </a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>


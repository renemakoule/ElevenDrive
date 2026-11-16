<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useAuthStore } from "../../stores/authStore";
import { useGuichetStore } from "../../stores/guichetStore";
import { useArchivageStore } from "../../stores/archivageStore";
import type {
  DossierArchiveItem,
  MotifArchivage,
  TypeDocumentImportant,
} from "../../types";

// Composants modulaires du centre d'archivage
import CentreArchivageHeader from "./CentreArchivageHeader.vue";
import ArchivageTabsNav, { type ArchivageTab } from "./ArchivageTabsNav.vue";
import CoffreFortArchivesTab from "./CoffreFortArchivesTab.vue";
import DossiersEligiblesTab from "./DossiersEligiblesTab.vue";
import DocumentsImportantsTab from "./DocumentsImportantsTab.vue";
import ArchiverDossierModal from "./ArchiverDossierModal.vue";
import DetailArchiveModal from "./DetailArchiveModal.vue";
import TeleverserDocumentModal from "./TeleverserDocumentModal.vue";
import ConfirmationPurgeModal from "./ConfirmationPurgeModal.vue";

const authStore = useAuthStore();
const guichetStore = useGuichetStore();
const archivageStore = useArchivageStore();

// État des onglets
const activeTab = ref<ArchivageTab>("archives");

// État des modales
const isArchiveModalOpen = ref(false);
const isViewModalOpen = ref(false);
const isUploadDocModalOpen = ref(false);
const isConfirmDeleteOpen = ref(false);

// Candidat sélectionné pour une action
const selectedCandidatId = ref<string>("");
const archiveToDelete = ref<DossierArchiveItem | null>(null);

// Sélection par lot (Batch Archive)
const selectedEligibleIds = ref<string[]>([]);

onMounted(async () => {
  await Promise.all([archivageStore.loadArchives(), guichetStore.search("")]);
});

// Candidats éligibles à l'archivage
const eligibleCandidats = computed(() => {
  return guichetStore.searchResults.filter((c) => c.statut_dossier !== "ARCHIVE");
});

// Modale d'archivage
const openArchiveModal = (candidatId?: string) => {
  selectedCandidatId.value = candidatId || eligibleCandidats.value[0]?.id || "";
  isArchiveModalOpen.value = true;
};

const handleArchivageSubmit = async (payload: {
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
}) => {
  try {
    await archivageStore.archiverDossier({
      candidatId: payload.candidatId,
      dossierId: payload.candidatId,
      motif: payload.motif,
      referenceBoite: payload.referenceBoite,
      observations: payload.observations,
      documentsFinaux: payload.documentsFinaux,
    });
    isArchiveModalOpen.value = false;
    activeTab.value = "archives";
  } catch (err: any) {
    guichetStore.showToast(err.message || "Erreur d'archivage", "error");
  }
};

// Consultation détaillée d'une archive
const openViewArchive = (archive: DossierArchiveItem) => {
  archivageStore.selectedArchiveId = archive.id;
  isViewModalOpen.value = true;
};

// Désarchivage
const handleDesarchiver = async (archive: DossierArchiveItem) => {
  if (
    confirm(
      `Êtes-vous sûr de vouloir désarchiver le dossier ${archive.numero_dossier} (${archive.candidat_nom}) et le remettre en cours d'apprentissage actif ?`,
    )
  ) {
    try {
      await archivageStore.desarchiverDossier(archive.id);
    } catch (err: any) {
      guichetStore.showToast(
        err.message || "Erreur lors du désarchivage",
        "error",
      );
    }
  }
};

// Purge certifiée Directeur
const openConfirmDelete = (archive: DossierArchiveItem) => {
  archiveToDelete.value = archive;
  isConfirmDeleteOpen.value = true;
};

const executeDeleteArchive = async () => {
  if (!archiveToDelete.value) return;
  try {
    await archivageStore.supprimerArchive(archiveToDelete.value.id);
    isConfirmDeleteOpen.value = false;
    archiveToDelete.value = null;
  } catch (err: any) {
    guichetStore.showToast(
      err.message || "Erreur lors de la suppression",
      "error",
    );
  }
};

// Archivage par lot
const toggleSelectAllEligibles = () => {
  if (selectedEligibleIds.value.length === eligibleCandidats.value.length) {
    selectedEligibleIds.value = [];
  } else {
    selectedEligibleIds.value = eligibleCandidats.value.map((c) => c.id);
  }
};

const executeBatchArchive = async () => {
  if (selectedEligibleIds.value.length === 0) return;
  if (
    !confirm(
      `Confirmez-vous l'archivage groupé de ${selectedEligibleIds.value.length} dossier(s) ?`,
    )
  ) {
    return;
  }

  for (const candId of selectedEligibleIds.value) {
    const cand = eligibleCandidats.value.find((c) => c.id === candId);
    if (cand) {
      try {
        await archivageStore.archiverDossier({
          candidatId: cand.id,
          dossierId: cand.id,
          motif: "DOSSIER_SOLDE",
          referenceBoite: `LOT-${new Date().getFullYear()}`,
          observations: "Archivage groupé automatisé",
        });
      } catch {
        // continue
      }
    }
  }
  selectedEligibleIds.value = [];
  activeTab.value = "archives";
};

// Téléversement direct
const openUploadDocModal = (candidatId?: string) => {
  selectedCandidatId.value =
    candidatId || guichetStore.searchResults[0]?.id || "";
  isUploadDocModalOpen.value = true;
};

const handleUploadDocSubmit = async (payload: {
  candidatId: string;
  type: TypeDocumentImportant;
  titre: string;
  nomFichier: string;
  base64: string;
  taille: number;
  notes: string;
}) => {
  await archivageStore.televerserDocumentImportant({
    candidatId: payload.candidatId,
    dossierId: payload.candidatId,
    type: payload.type,
    titre: payload.titre,
    nomFichier: payload.nomFichier,
    base64: payload.base64,
    taille: payload.taille,
    notes: payload.notes,
  });
  isUploadDocModalOpen.value = false;
};

const printInventory = () => {
  window.print();
};
</script>

<template>
  <div class="space-y-5 p-4 sm:p-6 max-w-7xl mx-auto">
    <!-- 1. EN-TÊTE PRINCIPAL & STATISTIQUES -->
    <CentreArchivageHeader
      :total-archives="archivageStore.totalArchivesCount"
      :total-docs="archivageStore.totalDocumentsArchived"
      :total-fonds="archivageStore.totalFondsArchives"
      :eligible-count="eligibleCandidats.length"
      :bucket-name="archivageStore.cloudR2Statut.bucketName"
      @open-archive="openArchiveModal()"
      @open-upload="openUploadDocModal()"
      @print-inventory="printInventory"
    />

    <!-- 2. NAVIGATION DES ONGLETS -->
    <ArchivageTabsNav
      v-model:active-tab="activeTab"
      :archives-count="archivageStore.filteredArchives.length"
      :eligibles-count="eligibleCandidats.length"
      :documents-count="archivageStore.documentsImportants.length"
    />

    <!-- 3. CONTENU DES ONGLETS -->
    <!-- Onglet 1 : Coffre-Fort des Archives -->
    <CoffreFortArchivesTab
      v-if="activeTab === 'archives'"
      :is-manager="authStore.currentRole === 'manager'"
      @view-archive="openViewArchive"
      @desarchiver="handleDesarchiver"
      @delete-archive="openConfirmDelete"
    />

    <!-- Onglet 2 : Dossiers Éligibles à l'Archivage -->
    <DossiersEligiblesTab
      v-else-if="activeTab === 'eligibles'"
      :candidates="eligibleCandidats"
      :selected-ids="selectedEligibleIds"
      @toggle-select-all="toggleSelectAllEligibles"
      @update:selected-ids="(ids) => (selectedEligibleIds = ids)"
      @batch-archive="executeBatchArchive"
      @archive-single="openArchiveModal"
    />

    <!-- Onglet 3 : Pièces & Documents Numérisés -->
    <DocumentsImportantsTab
      v-else-if="activeTab === 'documents'"
      :documents="archivageStore.documentsImportants"
      @open-upload="openUploadDocModal()"
    />

    <!-- 4. MODALES MODULAIRES -->
    <!-- Modale 1 : Archiver un dossier -->
    <ArchiverDossierModal
      :is-open="isArchiveModalOpen"
      :candidates="eligibleCandidats"
      :initial-candidat-id="selectedCandidatId"
      @close="isArchiveModalOpen = false"
      @submit="handleArchivageSubmit"
    />

    <!-- Modale 2 : Consultation détaillée de l'archive -->
    <DetailArchiveModal
      :is-open="isViewModalOpen"
      :archive="archivageStore.selectedArchive"
      @close="isViewModalOpen = false"
    />

    <!-- Modale 3 : Téléverser un document important -->
    <TeleverserDocumentModal
      :is-open="isUploadDocModalOpen"
      :candidates="guichetStore.searchResults"
      :initial-candidat-id="selectedCandidatId"
      @close="isUploadDocModalOpen = false"
      @submit="handleUploadDocSubmit"
    />

    <!-- Modale 4 : Confirmation purge certifiée Directeur -->
    <ConfirmationPurgeModal
      :is-open="isConfirmDeleteOpen"
      :archive="archiveToDelete"
      @close="isConfirmDeleteOpen = false"
      @confirm="executeDeleteArchive"
    />
  </div>
</template>

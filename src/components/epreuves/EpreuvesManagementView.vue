<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import { useEpreuvesStore } from "../../stores/epreuvesStore";
import { useGuichetStore } from "../../stores/guichetStore";
import { realtimeSocket } from "../../services/realtimeSocketService";
import type { EpreuvePdf } from "../../types";
import EpreuvesStatsCards from "./EpreuvesStatsCards.vue";
import EpreuvesFilterBar from "./EpreuvesFilterBar.vue";
import EpreuveCard from "./EpreuveCard.vue";
import EpreuveUploadModal from "./EpreuveUploadModal.vue";
import EpreuveViewerModal from "./EpreuveViewerModal.vue";
import { BookOpen, Plus, FileText, Radio } from "lucide-vue-next";

const epreuvesStore = useEpreuvesStore();
const guichetStore = useGuichetStore();
const isWsConnected = realtimeSocket.isConnected;
let unsubscribeWs: (() => void) | null = null;

const isAddModalOpen = ref(false);

onMounted(() => {
  epreuvesStore.loadEpreuves();
  unsubscribeWs = realtimeSocket.subscribe((msg) => {
    if (msg.type === "EPREUVES_UPDATED") {
      epreuvesStore.loadEpreuves();
    }
  });
});

onUnmounted(() => {
  if (unsubscribeWs) {
    unsubscribeWs();
    unsubscribeWs = null;
  }
});

const handleDelete = async (ep: EpreuvePdf) => {
  if (
    confirm(
      `Voulez-vous vraiment supprimer définitivement l'épreuve "${ep.titre}" de la base de données ?`,
    )
  ) {
    try {
      await epreuvesStore.deleteEpreuve(ep.id);
      guichetStore.showToast(
        "Épreuve supprimée de la base de données.",
        "info",
      );
    } catch (err: any) {
      guichetStore.showToast(`Erreur suppression: ${err.message}`, "error");
    }
  }
};
</script>

<template>
  <div class="space-y-5 p-4 sm:p-6 max-w-7xl mx-auto">
    <!-- En-tête de section -->
    <header
      class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 pb-4 border-b border-border"
    >
      <div class="space-y-1">
        <div class="flex items-center space-x-2">
          <div
            class="w-8 h-8 rounded-[4px] bg-brand-light border border-border text-accent-foreground flex items-center justify-center"
          >
            <BookOpen class="w-4 h-4 text-accent-foreground" />
          </div>
          <h1 class="text-base sm:text-lg font-bold text-foreground">
            Banque d'Épreuves & Annales PDF
          </h1>
          <span
            class="text-[11px] font-semibold px-2 py-0.5 rounded-[4px] bg-brand-light text-accent-foreground border border-brand-border"
          >
            Pédagogie & Révision
          </span>
        </div>
        <p class="text-xs text-muted-foreground">
          Mettez à disposition des candidats des sujets officiels et examens
          blancs par catégorie de permis (Permis A, B, C, D, E).
        </p>
      </div>

      <div class="flex items-center space-x-2">
        <button
          type="button"
          @click="isAddModalOpen = true"
          class="flex items-center space-x-1.5 px-3.5 py-2 rounded-[4px] bg-brand hover:bg-brand-dark text-white font-bold text-xs uppercase tracking-wider transition cursor-pointer"
        >
          <Plus class="w-4 h-4" />
          <span>Déposer une épreuve PDF</span>
        </button>
      </div>
    </header>

    <!-- Indicateurs / Statistiques Pédagogiques -->
    <EpreuvesStatsCards :stats="epreuvesStore.stats" />

    <!-- Barre de Recherche et Filtrage -->
    <EpreuvesFilterBar
      v-model:searchQuery="epreuvesStore.searchQuery"
      v-model:selectedCategory="epreuvesStore.selectedCategorieFilter"
      v-model:onlyWithCorrection="epreuvesStore.onlyWithCorrection"
    />

    <!-- Liste / Grille des Épreuves PDF -->
    <div
      v-if="epreuvesStore.filteredEpreuves.length === 0"
      class="p-12 text-center rounded-[4px] border border-dashed border-border bg-card space-y-3"
    >
      <FileText class="w-10 h-10 text-muted-foreground/40 mx-auto" />
      <div>
        <h3 class="text-sm font-semibold text-foreground">
          Aucune épreuve trouvée
        </h3>
        <p class="text-xs text-muted-foreground mt-0.5">
          Modifiez vos critères de recherche ou cliquez sur "Déposer une épreuve
          PDF" pour ajouter un nouveau sujet.
        </p>
      </div>
      <button
        type="button"
        @click="isAddModalOpen = true"
        class="inline-flex items-center space-x-1.5 px-3 py-1.5 rounded-[4px] bg-brand-light text-accent-foreground border border-border text-xs font-semibold hover:bg-brand-light/80 transition cursor-pointer"
      >
        <Plus class="w-3.5 h-3.5" />
        <span>Ajouter une épreuve</span>
      </button>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <EpreuveCard
        v-for="ep in epreuvesStore.filteredEpreuves"
        :key="ep.id"
        :epreuve="ep"
        @preview="epreuvesStore.previewPdf(ep)"
        @listen="epreuvesStore.previewPdf(ep)"
        @download="epreuvesStore.downloadPdf(ep)"
        @delete="handleDelete(ep)"
      />
    </div>

    <!-- Modale de Dépôt / Ajout -->
    <EpreuveUploadModal
      :isOpen="isAddModalOpen"
      @close="isAddModalOpen = false"
      @success="isAddModalOpen = false"
    />

    <!-- Modale de Visualisation PDF -->
    <EpreuveViewerModal
      :epreuve="epreuvesStore.activePdfToPreview"
      @close="epreuvesStore.closePreview()"
      @download="(ep) => epreuvesStore.downloadPdf(ep)"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import { useSupportsCoursStore } from "../../stores/supportsCoursStore";
import { realtimeSocket } from "../../services/realtimeSocketService";
import SupportCoursUploadModal from "./SupportCoursUploadModal.vue";
import DocumentReaderModal from "./DocumentReaderModal.vue";
import AudioSpeechBar from "../common/AudioSpeechBar.vue";
import { speechAudioService } from "../../services/speechAudioService";
import { PdfTextExtractor } from "../../services/pdfTextExtractor";
import {
  openExternalLink,
  getDomainFromUrl,
  formatExternalUrl,
} from "../../utils/urlHelper";
import type { SupportCours } from "../../types";
import {
  FileText,
  Plus,
  Search,
  Download,
  Trash2,
  Eye,
  Link,
  BookOpen,
  Sparkles,
  ExternalLink,
  Tag,
  CheckCircle2,
  X,
  Volume2,
  Radio,
} from "lucide-vue-next";

const supportsStore = useSupportsCoursStore();
const isWsConnected = realtimeSocket.isConnected;
let unsubscribeWs: (() => void) | null = null;

const isUploadModalOpen = ref(false);
const activeSupportToPreview = ref<SupportCours | null>(null);

const isAudioBarOpen = ref(false);
const isExtracting = ref(false);
const currentAudioTitle = ref("");

onMounted(() => {
  supportsStore.loadSupports();
  unsubscribeWs = realtimeSocket.subscribe((msg) => {
    if (msg.type === "SUPPORTS_COURS_UPDATED") {
      supportsStore.loadSupports();
    }
  });
});

onUnmounted(() => {
  speechAudioService.arreter();
  if (unsubscribeWs) {
    unsubscribeWs();
    unsubscribeWs = null;
  }
});

const CATEGORIES_FILTERS = [
  { id: "ALL", label: "Tous les cours" },
  { id: "B", label: "Permis B (Voiture)" },
  { id: "A", label: "Permis A / A1 (Moto)" },
  { id: "C", label: "Permis C (Poids Lourd)" },
  { id: "D", label: "Permis D (Transport)" },
  { id: "E", label: "Permis E (Remorque)" },
  { id: "TOUTES", label: "Tronc Commun" },
];

const openSupport = (support: SupportCours) => {
  if (support.type_support === "LIEN_WEB" && support.url_externe) {
    openExternalLink(support.url_externe);
  } else {
    activeSupportToPreview.value = support;
  }
};

const ecouterSupport = (support: SupportCours) => {
  if (support.type_support === "LIEN_WEB" && support.url_externe) {
    openExternalLink(support.url_externe);
    return;
  }
  activeSupportToPreview.value = support;
};

const handleCloseAudio = () => {
  speechAudioService.arreter();
  isAudioBarOpen.value = false;
};

const handleClosePreview = () => {
  speechAudioService.arreter();
  isAudioBarOpen.value = false;
  activeSupportToPreview.value = null;
};

const handleDelete = async (support: SupportCours) => {
  if (
    confirm(`Confirmez-vous la suppression du support "${support.titre}" ?`)
  ) {
    await supportsStore.supprimerSupport(support.id);
  }
};
</script>

<template>
  <div class="space-y-5 p-4 sm:p-6 max-w-7xl mx-auto">
    <!-- En-tête -->
    <header
      class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 pb-4 border-b border-border"
    >
      <div class="space-y-1">
        <div class="flex items-center space-x-2">
          <div
            class="w-8 h-8 rounded-[4px] bg-brand-light border border-brand-border text-brand-dark flex items-center justify-center"
          >
            <BookOpen class="w-4 h-4 text-brand" />
          </div>
          <h1 class="text-base sm:text-lg font-bold text-foreground">
            Supports de Cours Théoriques & Manuels PDF
          </h1>
          <span
            class="inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-bold"
            :class="
              isWsConnected
                ? 'bg-emerald-500/15 text-emerald-700 dark:text-emerald-300'
                : 'bg-amber-500/15 text-amber-700'
            "
          >
            <Radio class="h-3 w-3 animate-pulse" />
            {{ isWsConnected ? "Temps réel WS" : "Reconnexion WS..." }}
          </span>
          <span
            class="text-[11px] font-semibold px-2 py-0.5 rounded-[4px] bg-brand-light text-brand-dark border border-brand-border"
          >
            Pédagogie Code
          </span>
        </div>
        <p class="text-xs text-muted-foreground">
          Téléversez des manuels PDF, fiches de révision et écoutez les cours en
          synthèse vocale avec vos apprenants.
        </p>
      </div>

      <div class="flex items-center space-x-2">
        <button
          type="button"
          @click="isUploadModalOpen = true"
          class="flex items-center space-x-1.5 px-3.5 py-2 rounded-[4px] bg-brand hover:bg-brand-dark text-primary font-bold text-xs uppercase tracking-wider transition cursor-pointer"
        >
          <Plus class="w-4 h-4" />
          <span>Ajouter un support de cours</span>
        </button>
      </div>
    </header>

    <!-- Indicateurs -->
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
      <div class="p-3.5 rounded-[5px] bg-card border border-border">
        <span class="text-[11px] text-muted-foreground font-medium"
          >Total Documents</span
        >
        <p class="text-xl font-bold tabular-nums text-foreground mt-1">
          {{ supportsStore.stats.total }}
        </p>
      </div>

      <div class="p-3.5 rounded-[5px] bg-card border border-border">
        <span
          class="text-[11px] text-emerald-600 dark:text-emerald-400 font-semibold"
          >📄 Manuels PDF</span
        >
        <p class="text-xl font-bold tabular-nums text-foreground mt-1">
          {{ supportsStore.stats.pdfCount }}
        </p>
      </div>

      <div class="p-3.5 rounded-[5px] bg-card border border-border">
        <span class="text-[11px] text-sky-600 dark:text-sky-400 font-semibold"
          >🔗 Liens Web de cours</span
        >
        <p class="text-xl font-bold tabular-nums text-foreground mt-1">
          {{ supportsStore.stats.lienCount }}
        </p>
      </div>

      <div class="p-3.5 rounded-[5px] bg-card border border-border">
        <span
          class="text-[11px] text-amber-600 dark:text-amber-400 font-semibold"
          >Consultations / Téléchargements</span
        >
        <p class="text-xl font-bold tabular-nums text-foreground mt-1">
          {{ supportsStore.stats.totalDownloads }}
        </p>
      </div>
    </div>

    <!-- Filtres & Recherche -->
    <div
      class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 border-b border-border pb-3"
    >
      <div class="flex items-center space-x-1.5 overflow-x-auto pb-1 sm:pb-0">
        <button
          v-for="cat in CATEGORIES_FILTERS"
          :key="cat.id"
          type="button"
          @click="supportsStore.selectedCategorieFilter = cat.id"
          class="px-3 py-1.5 rounded-[4px] text-xs font-semibold transition cursor-pointer border shrink-0"
          :class="
            supportsStore.selectedCategorieFilter === cat.id
              ? 'bg-foreground text-background border-foreground font-bold'
              : 'bg-card text-muted-foreground border-border hover:text-foreground'
          "
        >
          {{ cat.label }}
        </button>
      </div>

      <div class="flex items-center gap-2">
        <!-- Commutateur Type PDF / Lien -->
        <select
          v-model="supportsStore.selectedTypeFilter"
          class="h-8 px-2 rounded-[4px] border border-border bg-card text-foreground text-xs focus:outline-none focus:border-brand-dark"
        >
          <option value="ALL">Tous les formats</option>
          <option value="PDF">📄 PDF uniquement</option>
          <option value="LIEN_WEB">🔗 Liens web uniquement</option>
        </select>

        <div class="relative w-full sm:w-56">
          <input
            v-model="supportsStore.searchQuery"
            type="text"
            placeholder="Rechercher un manuel..."
            class="w-full h-8 pl-8 pr-3 rounded-[4px] border border-border bg-card text-foreground text-xs focus:outline-none focus:border-brand-dark"
          />
          <Search
            class="w-3.5 h-3.5 text-muted-foreground absolute left-2.5 top-2.5"
          />
        </div>
      </div>
    </div>

    <!-- Grille des Supports de Cours -->
    <div
      v-if="supportsStore.filteredSupports.length === 0"
      class="p-12 text-center rounded-[5px] border border-dashed border-border bg-card space-y-3"
    >
      <FileText class="w-10 h-10 text-muted-foreground/40 mx-auto" />
      <div>
        <h3 class="text-sm font-semibold text-foreground">
          Aucun support de cours disponible
        </h3>
        <p class="text-xs text-muted-foreground mt-0.5">
          Publiez des manuels de cours au format PDF ou des liens web pour
          permettre à vos élèves de réviser.
        </p>
      </div>
      <button
        type="button"
        @click="isUploadModalOpen = true"
        class="inline-flex items-center space-x-1.5 px-3 py-1.5 rounded-[4px] bg-brand-light text-primary border border-brand-border text-xs font-semibold hover:bg-brand-light/80 transition cursor-pointer"
      >
        <Plus class="w-3.5 h-3.5" />
        <span>Téléverser un premier cours</span>
      </button>
    </div>

    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="s in supportsStore.filteredSupports"
        :key="s.id"
        class="p-4 rounded-[5px] bg-card border border-border hover:border-brand-border transition shadow-xs flex flex-col justify-between space-y-3"
      >
        <div class="space-y-2">
          <div class="flex items-center justify-between gap-2">
            <span
              class="inline-flex items-center gap-1 text-[10px] font-bold px-2 py-0.5 rounded-[4px] border uppercase"
              :class="
                s.type_support === 'PDF'
                  ? 'bg-emerald-500/10 text-emerald-600 border-emerald-500/20'
                  : 'bg-sky-500/10 text-sky-600 border-sky-500/20'
              "
            >
              {{ s.type_support === "PDF" ? "📄 Document PDF" : "🔗 Lien Web" }}
            </span>

            <span
              class="text-[10px] font-mono text-muted-foreground bg-muted px-2 py-0.5 rounded"
            >
              Permis {{ s.categorie_permis }}
            </span>
          </div>

          <h3
            @click="openSupport(s)"
            class="text-xs sm:text-sm font-bold text-foreground leading-snug line-clamp-2 cursor-pointer hover:text-brand transition"
            :title="s.titre"
          >
            {{ s.titre }}
          </h3>

          <!-- Lien Web directement cliquable et visible -->
          <div
            v-if="s.type_support === 'LIEN_WEB' && s.url_externe"
            class="pt-0.5"
          >
            <a
              :href="formatExternalUrl(s.url_externe)"
              @click.stop.prevent="openExternalLink(s.url_externe)"
              class="inline-flex items-center gap-1.5 text-xs font-semibold text-sky-600 hover:text-sky-700 dark:text-sky-400 hover:underline max-w-full truncate cursor-pointer"
              :title="formatExternalUrl(s.url_externe)"
            >
              <ExternalLink class="w-3.5 h-3.5 shrink-0" />
              <span class="truncate">{{ s.url_externe }}</span>
            </a>
          </div>

          <p
            v-if="s.description"
            class="text-[11px] text-muted-foreground line-clamp-2"
          >
            {{ s.description }}
          </p>
        </div>

        <div class="pt-2 border-t border-border space-y-2 text-xs">
          <div
            class="flex items-center justify-between text-[11px] text-muted-foreground"
          >
            <span class="truncate max-w-[150px]"
              >Par {{ s.ajoute_par || "Direction" }}</span
            >
            <a
              v-if="s.type_support === 'LIEN_WEB' && s.url_externe"
              :href="formatExternalUrl(s.url_externe)"
              @click.stop.prevent="openExternalLink(s.url_externe)"
              class="font-mono text-[10px] text-sky-600 dark:text-sky-400 hover:underline cursor-pointer flex items-center gap-1"
              :title="formatExternalUrl(s.url_externe)"
            >
              <span>{{ getDomainFromUrl(s.url_externe) }}</span>
              <ExternalLink class="w-2.5 h-2.5 shrink-0" />
            </a>
            <span v-else class="font-mono text-[10px]">
              {{ s.fichier_taille ? s.fichier_taille : "PDF" }}
            </span>
          </div>

          <!-- Actions -->
          <div class="pt-1 flex items-center justify-between gap-2">
            <button
              type="button"
              @click="openSupport(s)"
              class="flex-1 flex items-center justify-center space-x-1.5 px-3 py-1.5 rounded-[4px] bg-brand hover:bg-brand-dark text-primary font-bold text-xs transition cursor-pointer"
            >
              <component
                :is="s.type_support === 'PDF' ? Eye : ExternalLink"
                class="w-3.5 h-3.5"
              />
              <span>{{
                s.type_support === "PDF" ? "Consulter le PDF" : "Ouvrir le lien"
              }}</span>
            </button>

            <!-- Bouton Écouter le cours en TTS -->
            <button
              type="button"
              @click="ecouterSupport(s)"
              class="p-1.5 rounded-[4px] border border-border hover:bg-brand/10 text-muted-foreground hover:text-brand transition cursor-pointer"
              title="Écouter la lecture vocale de ce support"
            >
              <Volume2 class="w-3.5 h-3.5" />
            </button>

            <button
              v-if="s.type_support === 'PDF'"
              type="button"
              @click="supportsStore.downloadPdf(s)"
              class="p-1.5 rounded-[4px] border border-border hover:bg-muted text-muted-foreground hover:text-foreground transition cursor-pointer"
              title="Télécharger le fichier"
            >
              <Download class="w-3.5 h-3.5" />
            </button>

            <button
              type="button"
              @click="handleDelete(s)"
              class="p-1.5 rounded-[4px] border border-border hover:bg-red-500/10 hover:border-red-500/30 text-muted-foreground hover:text-red-600 transition cursor-pointer"
              title="Supprimer ce support"
            >
              <Trash2 class="w-3.5 h-3.5" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Lecteur audio flottant global pour la page de gestion -->
    <div
      v-if="isAudioBarOpen && !activeSupportToPreview"
      class="fixed bottom-4 left-4 right-4 sm:left-auto sm:right-6 sm:w-96 z-40 shadow-2xl rounded-xl overflow-hidden border border-border"
    >
      <AudioSpeechBar
        :title="currentAudioTitle"
        :is-extracting="isExtracting"
        @close="handleCloseAudio"
      />
    </div>

    <!-- Modale d'ajout -->
    <SupportCoursUploadModal
      :isOpen="isUploadModalOpen"
      @close="isUploadModalOpen = false"
    />

    <!-- Visionneuse PDF & Lecteur Bilingue Avancé Modal -->
    <DocumentReaderModal
      v-if="activeSupportToPreview"
      :is-open="!!activeSupportToPreview"
      :titre="activeSupportToPreview.titre"
      :categorie-permis="activeSupportToPreview.categorie_permis"
      :fichier-data="activeSupportToPreview.fichier_data"
      :description="activeSupportToPreview.description"
      @download="supportsStore.downloadPdf(activeSupportToPreview)"
      @close="handleClosePreview"
    />
  </div>
</template>

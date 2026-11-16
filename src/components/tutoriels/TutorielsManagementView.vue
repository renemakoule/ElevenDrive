<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import { useTutorielsStore } from "../../stores/tutorielsStore";
import { realtimeSocket } from "../../services/realtimeSocketService";
import VideoPlayerModal from "./VideoPlayerModal.vue";
import TutorielUploadModal from "./TutorielUploadModal.vue";
import type { TutorielVideo } from "../../types";
import {
  Play,
  Plus,
  Search,
  Clock,
  Trash2,
  Eye,
  Film,
  Sparkles,
  ExternalLink,
  Tag,
  CheckCircle2,
  Radio,
} from "lucide-vue-next";

const tutorielsStore = useTutorielsStore();
const isWsConnected = realtimeSocket.isConnected;
let unsubscribeWs: (() => void) | null = null;

const isUploadModalOpen = ref(false);
const videoToPlay = ref<TutorielVideo | null>(null);

onMounted(() => {
  tutorielsStore.loadTutoriels();
  unsubscribeWs = realtimeSocket.subscribe((msg) => {
    if (msg.type === "TUTORIELS_UPDATED") {
      tutorielsStore.loadTutoriels();
    }
  });
});

onUnmounted(() => {
  if (unsubscribeWs) {
    unsubscribeWs();
    unsubscribeWs = null;
  }
});

const CATEGORIES_FILTERS = [
  { id: "ALL", label: "Toutes les vidéos" },
  { id: "B", label: "Permis B (Voiture)" },
  { id: "A", label: "Permis A / A1 (Moto)" },
  { id: "C", label: "Permis C (Poids Lourd)" },
  { id: "D", label: "Permis D (Transport)" },
  { id: "E", label: "Permis E (Remorque)" },
  { id: "TOUTES", label: "Tronc Commun" },
];

const openVideo = (video: TutorielVideo) => {
  videoToPlay.value = video;
  tutorielsStore.playVideo(video);
};

const handleDelete = async (video: TutorielVideo) => {
  if (confirm(`Confirmez-vous la suppression du tutoriel "${video.titre}" ?`)) {
    await tutorielsStore.supprimerTutoriel(video.id);
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
            class="w-8 h-8 rounded-[4px] bg-red-500/10 border border-red-500/20 text-red-600 flex items-center justify-center"
          >
            <Film class="w-4 h-4 text-red-600" />
          </div>
          <h1 class="text-base sm:text-lg font-bold text-foreground">
            Guides & Tutoriels Vidéo Pédagogiques
          </h1>
          <span
            class="text-[11px] font-semibold px-2 py-0.5 rounded-[4px] bg-red-500/10 text-red-600 border border-red-500/20"
          >
            Pédagogie Pratique
          </span>
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
        </div>
        <p class="text-xs text-muted-foreground">
          Mettez à disposition des vidéos d'apprentissage, repères de manœuvres
          et explications pratiques pour vos élèves.
        </p>
      </div>

      <div class="flex items-center space-x-2">
        <button
          type="button"
          @click="isUploadModalOpen = true"
          class="flex items-center space-x-1.5 px-3.5 py-2 rounded-[4px] bg-brand hover:bg-brand-dark text-brand-dark font-bold text-xs uppercase tracking-wider transition shadow-xs cursor-pointer"
        >
          <Plus class="w-4 h-4" />
          <span>Ajouter un tutoriel vidéo</span>
        </button>
      </div>
    </header>

    <!-- Indicateurs -->
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
      <div class="p-3.5 rounded-[5px] bg-card border border-border">
        <span class="text-[11px] text-muted-foreground font-medium"
          >Total Vidéos</span
        >
        <p class="text-xl font-bold tabular-nums text-foreground mt-1">
          {{ tutorielsStore.stats.total }}
        </p>
      </div>

      <div class="p-3.5 rounded-[5px] bg-card border border-border">
        <span class="text-[11px] text-sky-600 dark:text-sky-400 font-semibold"
          >🚗 Permis B (Auto)</span
        >
        <p class="text-xl font-bold tabular-nums text-foreground mt-1">
          {{ tutorielsStore.stats.permisB }}
        </p>
      </div>

      <div class="p-3.5 rounded-[5px] bg-card border border-border">
        <span
          class="text-[11px] text-amber-600 dark:text-amber-400 font-semibold"
          >🏍️ Permis A (Moto)</span
        >
        <p class="text-xl font-bold tabular-nums text-foreground mt-1">
          {{ tutorielsStore.stats.permisA }}
        </p>
      </div>

      <div class="p-3.5 rounded-[5px] bg-card border border-border">
        <span
          class="text-[11px] text-emerald-600 dark:text-emerald-400 font-semibold"
          >Vues Apprenants</span
        >
        <p class="text-xl font-bold tabular-nums text-foreground mt-1">
          {{ tutorielsStore.stats.totalVues }}
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
          @click="tutorielsStore.selectedCategorieFilter = cat.id"
          class="px-3 py-1.5 rounded-[4px] text-xs font-semibold transition cursor-pointer border shrink-0"
          :class="
            tutorielsStore.selectedCategorieFilter === cat.id
              ? 'bg-foreground text-background border-foreground font-bold'
              : 'bg-card text-muted-foreground border-border hover:text-foreground'
          "
        >
          {{ cat.label }}
        </button>
      </div>

      <div class="relative w-full sm:w-64">
        <input
          v-model="tutorielsStore.searchQuery"
          type="text"
          placeholder="Rechercher une manœuvre..."
          class="w-full h-8 pl-8 pr-3 rounded-[4px] border border-border bg-card text-foreground text-xs focus:outline-none focus:border-brand-dark"
        />
        <Search
          class="w-3.5 h-3.5 text-muted-foreground absolute left-2.5 top-2.5"
        />
      </div>
    </div>

    <!-- Grille des Vidéos -->
    <div
      v-if="tutorielsStore.filteredTutoriels.length === 0"
      class="p-12 text-center rounded-[5px] border border-dashed border-border bg-card space-y-3"
    >
      <Film class="w-10 h-10 text-muted-foreground/40 mx-auto" />
      <div>
        <h3 class="text-sm font-semibold text-foreground">
          Aucun tutoriel vidéo disponible
        </h3>
        <p class="text-xs text-muted-foreground mt-0.5">
          Ajoutez des liens YouTube, Vimeo ou fichiers MP4 pour enrichir le
          parcours d'apprentissage de vos élèves.
        </p>
      </div>
      <button
        type="button"
        @click="isUploadModalOpen = true"
        class="inline-flex items-center space-x-1.5 px-3 py-1.5 rounded-[4px] bg-brand-light text-brand-dark border border-brand-border text-xs font-semibold hover:bg-brand-light/80 transition cursor-pointer"
      >
        <Plus class="w-3.5 h-3.5" />
        <span>Publier une première vidéo</span>
      </button>
    </div>

    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="v in tutorielsStore.filteredTutoriels"
        :key="v.id"
        class="rounded-[5px] bg-card border border-border hover:border-red-500/50 transition shadow-xs flex flex-col justify-between overflow-hidden group"
      >
        <!-- Miniature avec bouton lecture overlay -->
        <div
          class="relative aspect-video bg-black/90 overflow-hidden cursor-pointer"
          @click="openVideo(v)"
        >
          <img
            v-if="v.miniature_url"
            :src="v.miniature_url"
            :alt="v.titre"
            class="w-full h-full object-cover group-hover:scale-105 transition duration-300 opacity-90 group-hover:opacity-100"
          />
          <div
            v-else
            class="w-full h-full flex items-center justify-center bg-muted/30"
          >
            <Film class="w-10 h-10 text-muted-foreground/30" />
          </div>

          <!-- Bouton Play central -->
          <div
            class="absolute inset-0 flex items-center justify-center bg-black/20 group-hover:bg-black/40 transition"
          >
            <div
              class="w-11 h-11 rounded-full bg-red-600 text-white flex items-center justify-center shadow-lg group-hover:scale-110 transition"
            >
              <Play class="w-5 h-5 fill-current ml-0.5" />
            </div>
          </div>

          <!-- Badges sur la miniature -->
          <div class="absolute top-2 left-2 flex items-center gap-1.5">
            <span
              class="text-[10px] font-bold px-2 py-0.5 rounded-[4px] bg-black/75 text-white backdrop-blur-xs uppercase"
            >
              Permis {{ v.categorie_permis }}
            </span>
          </div>

          <div
            class="absolute bottom-2 right-2 text-[10px] font-mono font-bold px-1.5 py-0.5 rounded bg-black/80 text-white flex items-center gap-1"
          >
            <Clock class="w-3 h-3 text-red-400" />
            {{ v.duree_minutes ? v.duree_minutes + " min" : "Vidéo" }}
          </div>
        </div>

        <!-- Contenu de la fiche -->
        <div class="p-3.5 space-y-2 flex-1 flex flex-col justify-between">
          <div>
            <h3
              class="text-xs sm:text-sm font-bold text-foreground leading-snug line-clamp-2"
            >
              {{ v.titre }}
            </h3>
            <p
              v-if="v.description"
              class="text-[11px] text-muted-foreground line-clamp-2 mt-1"
            >
              {{ v.description }}
            </p>
          </div>

          <div
            class="pt-2 border-t border-border flex items-center justify-between text-[10px] text-muted-foreground"
          >
            <span class="truncate max-w-[150px]"
              >Par {{ v.ajoute_par || "Direction" }}</span
            >
            <span class="flex items-center gap-1 font-mono">
              <Eye class="w-3 h-3" /> {{ v.vues_count || 0 }} vues
            </span>
          </div>
        </div>

        <!-- Actions -->
        <div
          class="px-3.5 py-2 border-t border-border bg-muted/10 flex items-center justify-between gap-2"
        >
          <button
            type="button"
            @click="openVideo(v)"
            class="flex-1 flex items-center justify-center space-x-1.5 px-3 py-1.5 rounded-[4px] bg-red-600 hover:bg-red-700 text-white font-bold text-xs transition cursor-pointer shadow-xs"
          >
            <Play class="w-3.5 h-3.5 fill-current" />
            <span>Lire la vidéo</span>
          </button>

          <button
            type="button"
            @click="handleDelete(v)"
            class="p-1.5 rounded-[4px] border border-border hover:bg-red-500/10 hover:border-red-500/30 text-muted-foreground hover:text-red-600 transition cursor-pointer"
            title="Supprimer ce tutoriel"
          >
            <Trash2 class="w-3.5 h-3.5" />
          </button>
        </div>
      </div>
    </div>

    <!-- Modale d'ajout -->
    <TutorielUploadModal
      :isOpen="isUploadModalOpen"
      @close="isUploadModalOpen = false"
    />

    <!-- Lecteur Vidéo Modal -->
    <VideoPlayerModal :video="videoToPlay" @close="videoToPlay = null" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useOfflineDownloadsStore } from '../../stores/offlineDownloadsStore';
import type { OfflineItem } from '../../services/offlineStorageService';
import OfflineItemViewerModal from './OfflineItemViewerModal.vue';
import {
  DownloadCloud,
  CheckCircle2,
  Trash2,
  Eye,
  Play,
  FileText,
  Film,
  BookOpen,
  Search,
  HardDrive,
  Wifi,
  WifiOff,
  Sparkles,
  ArrowRight
} from 'lucide-vue-next';

const emit = defineEmits<{
  (e: 'navigate', tab: 'cours' | 'supports' | 'tutoriels' | 'epreuves'): void;
}>();

const offlineStore = useOfflineDownloadsStore();
const activePreviewItem = ref<OfflineItem | null>(null);

onMounted(() => {
  offlineStore.initListeners();
  offlineStore.loadDownloads();
});

const openItem = (item: OfflineItem) => {
  activePreviewItem.value = item;
};

const handleDelete = async (item: OfflineItem) => {
  if (confirm(`Voulez-vous supprimer "${item.titre}" de votre espace hors-ligne ?`)) {
    await offlineStore.removeDownload(item.id);
  }
};

const handleClearAll = async () => {
  if (confirm("Confirmez-vous la suppression de TOUS vos téléchargements hors-ligne ?")) {
    await offlineStore.clearAllDownloads();
  }
};
</script>

<template>
  <div class="space-y-4">
    <!-- Bannière Mode Connectivité -->
    <div
      v-if="!offlineStore.isOnline"
      class="p-3 rounded-[5px] bg-amber-500/10 border border-amber-500/30 text-amber-700 dark:text-amber-300 flex items-center justify-between text-xs font-semibold animate-pulse"
    >
      <div class="flex items-center space-x-2">
        <WifiOff class="w-4 h-4 text-amber-500 shrink-0" />
        <span>Vous êtes actuellement <strong>hors-ligne</strong>. Vos cours et épreuves téléchargés restent 100% accessibles.</span>
      </div>
      <span class="px-2 py-0.5 rounded bg-amber-500/20 text-[10px] uppercase font-bold tracking-wider">
        Mode Hors-Ligne
      </span>
    </div>

    <!-- En-tête Espace Téléchargements (Style YouTube/Netflix Offline) -->
    <section class="p-4 rounded-[5px] border border-border bg-card flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
      <div class="flex items-center space-x-3">
        <div class="w-10 h-10 rounded-lg bg-emerald-500/10 border border-emerald-500/20 text-emerald-600 flex items-center justify-center shrink-0">
          <DownloadCloud class="w-5 h-5 text-emerald-600" />
        </div>
        <div>
          <div class="flex items-center gap-2">
            <h2 class="text-sm sm:text-base font-bold text-foreground">
              Mes Téléchargements Hors-Ligne
            </h2>
            <span class="inline-flex items-center gap-1 text-[10px] font-bold px-2 py-0.5 rounded-full bg-emerald-500/10 text-emerald-600 border border-emerald-500/20">
              <CheckCircle2 class="w-3 h-3" /> Espace Hors-Ligne
            </span>
          </div>
          <p class="text-xs text-muted-foreground mt-0.5">
            Vos cours, annales d'examens et guides vidéo stockés en local pour réviser sans connexion Internet.
          </p>
        </div>
      </div>

      <!-- Compteur stockage & Bouton Vider -->
      <div class="flex flex-wrap items-center gap-2.5 w-full md:w-auto justify-between md:justify-end">
        <div class="flex items-center gap-2 px-3 py-1.5 rounded-[5px] border border-border bg-muted/30 text-xs">
          <HardDrive class="w-3.5 h-3.5 text-muted-foreground" />
          <span class="text-muted-foreground">Espace utilisé :</span>
          <strong class="text-foreground font-mono">{{ offlineStore.stats.totalSizeFormatted }}</strong>
          <span class="text-muted-foreground">({{ offlineStore.stats.total }} élément(s))</span>
        </div>

        <button
          v-if="offlineStore.items.length > 0"
          type="button"
          @click="handleClearAll"
          class="flex items-center space-x-1 px-2.5 py-1.5 rounded-lg border border-red-500/30 bg-red-500/10 hover:bg-red-500/20 text-red-600 text-xs font-semibold transition cursor-pointer"
          title="Supprimer tous les fichiers pour libérer de l'espace"
        >
          <Trash2 class="w-3.5 h-3.5" />
          <span class="hidden sm:inline">Vider l'espace</span>
        </button>
      </div>
    </section>

    <!-- Filtres & Barre de recherche -->
    <div class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3">
      <!-- Filtres par type de document -->
      <div class="flex items-center space-x-1 overflow-x-auto p-1 bg-muted/20 border border-border rounded-[5px] shrink-0">
        <button
          type="button"
          @click="offlineStore.activeFilter = 'ALL'"
          class="px-3 py-1.5 rounded-md text-xs font-semibold transition shrink-0"
          :class="offlineStore.activeFilter === 'ALL'
            ? 'bg-card text-foreground border border-border'
            : 'text-muted-foreground hover:text-foreground'"
        >
          Tous ({{ offlineStore.stats.total }})
        </button>

        <button
          type="button"
          @click="offlineStore.activeFilter = 'SUPPORT'"
          class="px-3 py-1.5 rounded-md text-xs font-semibold transition shrink-0"
          :class="offlineStore.activeFilter === 'SUPPORT'
            ? 'bg-card text-foreground border border-border'
            : 'text-muted-foreground hover:text-foreground'"
        >
          Supports de Cours ({{ offlineStore.stats.supportsCount }})
        </button>

        <button
          type="button"
          @click="offlineStore.activeFilter = 'EPREUVE'"
          class="px-3 py-1.5 rounded-md text-xs font-semibold transition shrink-0"
          :class="offlineStore.activeFilter === 'EPREUVE'
            ? 'bg-card text-foreground border border-border'
            : 'text-muted-foreground hover:text-foreground'"
        >
          Épreuves PDF ({{ offlineStore.stats.epreuvesCount }})
        </button>

        <button
          type="button"
          @click="offlineStore.activeFilter = 'TUTORIEL'"
          class="px-3 py-1.5 rounded-md text-xs font-semibold transition shrink-0"
          :class="offlineStore.activeFilter === 'TUTORIEL'
            ? 'bg-card text-foreground border border-border'
            : 'text-muted-foreground hover:text-foreground'"
        >
          Vidéos ({{ offlineStore.stats.tutorielsCount }})
        </button>
      </div>

      <!-- Recherche instantanée -->
      <div class="relative w-full sm:w-64">
        <input
          v-model="offlineStore.searchQuery"
          type="text"
          placeholder="Rechercher hors-ligne..."
          class="w-full h-8 pl-8 pr-3 rounded-[5px] border border-border bg-background text-foreground text-xs focus:outline-none focus:border-brand-dark"
        />
        <Search class="w-3.5 h-3.5 text-muted-foreground absolute left-2.5 top-2.5" />
      </div>
    </div>

    <!-- Si aucun élément téléchargé -->
    <div
      v-if="offlineStore.filteredItems.length === 0"
      class="p-12 text-center rounded-[5px] border border-dashed border-border bg-card space-y-4"
    >
      <div class="w-16 h-16 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-600 flex items-center justify-center mx-auto">
        <DownloadCloud class="w-8 h-8" />
      </div>

      <div class="space-y-1 max-w-md mx-auto">
        <h3 class="text-sm sm:text-base font-bold text-foreground">
          Aucun document téléchargé hors-ligne
        </h3>
        <p class="text-xs text-muted-foreground">
          Téléchargez vos cours, fiches pratiques et sujets d'épreuves pour y accéder à tout moment, même dans les transports ou sans forfait data.
        </p>
      </div>

      <div class="flex flex-wrap items-center justify-center gap-2 pt-2">
        <button
          type="button"
          @click="$emit('navigate', 'supports')"
          class="inline-flex items-center gap-1.5 px-3 py-2 rounded-lg bg-emerald-600 hover:bg-emerald-700 text-white font-bold text-xs transition shadow-xs"
        >
          <FileText class="w-3.5 h-3.5" />
          <span>Explorer les supports</span>
        </button>

        <button
          type="button"
          @click="$emit('navigate', 'epreuves')"
          class="inline-flex items-center gap-1.5 px-3 py-2 rounded-lg border border-border hover:bg-muted text-foreground font-semibold text-xs transition"
        >
          <BookOpen class="w-3.5 h-3.5 text-brand" />
          <span>Voir les épreuves PDF</span>
        </button>

        <button
          type="button"
          @click="$emit('navigate', 'tutoriels')"
          class="inline-flex items-center gap-1.5 px-3 py-2 rounded-lg border border-border hover:bg-muted text-foreground font-semibold text-xs transition"
        >
          <Film class="w-3.5 h-3.5 text-red-500" />
          <span>Découvrir les tutoriels</span>
        </button>
      </div>
    </div>

    <!-- Grille des Éléments Téléchargés Hors-Ligne -->
    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3.5">
      <div
        v-for="item in offlineStore.filteredItems"
        :key="item.id"
        class="rounded-[5px] bg-card border border-border hover:border-emerald-500/50 transition p-4 flex flex-col justify-between space-y-3 group"
      >
        <div class="space-y-2.5">
          <!-- Badges type et statut hors-ligne -->
          <div class="flex items-center justify-between">
            <span
              class="inline-flex items-center gap-1 text-[10px] font-bold px-2 py-0.5 rounded-full border uppercase"
              :class="{
                'bg-emerald-500/10 text-emerald-600 border-emerald-500/20': item.type === 'SUPPORT',
                'bg-brand-light text-brand-dark border-brand-border': item.type === 'EPREUVE',
                'bg-red-500/10 text-red-600 border-red-500/20': item.type === 'TUTORIEL',
              }"
            >
              <component
                :is="item.type === 'TUTORIEL' ? Film : (item.type === 'EPREUVE' ? BookOpen : FileText)"
                class="w-3 h-3"
              />
              {{ item.type === 'SUPPORT' ? 'Support de cours' : (item.type === 'EPREUVE' ? 'Épreuve PDF' : 'Tutoriel Vidéo') }}
            </span>

            <span class="inline-flex items-center gap-1 text-[10px] font-semibold text-emerald-600 dark:text-emerald-400 bg-emerald-500/10 px-2 py-0.5 rounded-full">
              <CheckCircle2 class="w-3 h-3" /> Hors-Ligne
            </span>
          </div>

          <h3 class="text-xs sm:text-sm font-bold text-foreground leading-snug line-clamp-2">
            {{ item.titre }}
          </h3>

          <p v-if="item.description" class="text-[11px] text-muted-foreground line-clamp-2">
            {{ item.description }}
          </p>
        </div>

        <div class="pt-2.5 border-t border-border space-y-2.5">
          <!-- Métadonnées taille et date -->
          <div class="flex items-center justify-between text-[10px] text-muted-foreground">
            <span class="font-mono font-semibold text-foreground">{{ item.taille_formatee }}</span>
            <span>Permis <strong>{{ item.categorie_permis }}</strong></span>
            <span>{{ new Date(item.date_telechargement).toLocaleDateString('fr-FR') }}</span>
          </div>

          <!-- Actions de lecture et suppression -->
          <div class="flex items-center gap-2">
            <button
              type="button"
              @click="openItem(item)"
              class="flex-1 flex items-center justify-center space-x-1.5 px-3 py-2 rounded-lg bg-emerald-600 hover:bg-emerald-700 text-white font-bold text-xs transition cursor-pointer shadow-xs active:scale-[0.99]"
            >
              <component :is="item.type === 'TUTORIEL' ? Play : Eye" class="w-3.5 h-3.5 fill-current" />
              <span>{{ item.type === 'TUTORIEL' ? 'Regarder' : 'Consulter' }}</span>
            </button>

            <button
              type="button"
              @click="handleDelete(item)"
              class="p-2 rounded-lg border border-border hover:bg-red-500/10 text-muted-foreground hover:text-red-600 transition cursor-pointer"
              title="Supprimer du stockage hors-ligne"
            >
              <Trash2 class="w-3.5 h-3.5" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Visionneuse Hors-Ligne Modal -->
    <OfflineItemViewerModal
      :item="activePreviewItem"
      :is-online="offlineStore.isOnline"
      @close="activePreviewItem = null"
    />
  </div>
</template>

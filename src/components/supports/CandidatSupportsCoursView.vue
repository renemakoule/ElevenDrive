<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { useSupportsCoursStore } from "../../stores/supportsCoursStore";
import { useOfflineDownloadsStore } from "../../stores/offlineDownloadsStore";
import { realtimeSocket } from "../../services/realtimeSocketService";
import type { SupportCours } from "../../types";
import { speechAudioService } from "../../services/speechAudioService";
import { PdfTextExtractor } from "../../services/pdfTextExtractor";
import {
  openExternalLink,
  getDomainFromUrl,
  formatExternalUrl,
} from "../../utils/urlHelper";
import AudioSpeechBar from "../common/AudioSpeechBar.vue";
import DocumentReaderModal from "./DocumentReaderModal.vue";
import {
  FileText,
  BookOpen,
  DownloadCloud,
  Eye,
  Link,
  Search,
  CheckCircle2,
  ExternalLink,
  X,
  Volume2,
  Radio,
} from "lucide-vue-next";

const props = defineProps<{
  studentCategory: string;
}>();

const supportsStore = useSupportsCoursStore();
const offlineStore = useOfflineDownloadsStore();
const isWsConnected = realtimeSocket.isConnected;
let unsubscribeWs: (() => void) | null = null;
const searchQuery = ref("");
const activeSupportToPreview = ref<SupportCours | null>(null);

const isAudioBarOpen = ref(false);
const isExtracting = ref(false);
const currentAudioTitle = ref("");

onMounted(() => {
  supportsStore.loadSupports();
  offlineStore.loadDownloads();
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

const studentSupports = computed(() => {
  const base = supportsStore.getSupportsForStudent(props.studentCategory);
  if (!searchQuery.value.trim()) return base;
  const q = searchQuery.value.toLowerCase().trim();
  return base.filter(
    (s) =>
      s.titre.toLowerCase().includes(q) ||
      (s.description && s.description.toLowerCase().includes(q)),
  );
});

const openSupport = (s: SupportCours) => {
  if (s.type_support === "LIEN_WEB" && s.url_externe) {
    openExternalLink(s.url_externe);
  } else {
    activeSupportToPreview.value = s;
  }
};

const ecouterSupport = (support: SupportCours) => {
  if (support.type_support === "LIEN_WEB" && support.url_externe) {
    openExternalLink(support.url_externe);
    return;
  }
  activeSupportToPreview.value = support;
};

const fermerAudio = () => {
  speechAudioService.arreter();
  isAudioBarOpen.value = false;
};

const fermerModalPreview = () => {
  speechAudioService.arreter();
  isAudioBarOpen.value = false;
  activeSupportToPreview.value = null;
};
</script>

<template>
  <div class="space-y-4">
    <!-- Bannière explicative -->
    <section
      class="p-4 rounded-xl border border-border bg-card flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3 shadow-xs"
    >
      <div class="flex items-center space-x-3">
        <div
          class="w-10 h-10 rounded-lg bg-emerald-500/10 border border-emerald-500/20 text-emerald-600 flex items-center justify-center shrink-0"
        >
          <BookOpen class="w-5 h-5 text-emerald-600" />
        </div>
        <div>
          <div class="flex items-center gap-2">
            <h2 class="text-sm sm:text-base font-bold text-foreground">
              Manuels & Supports de Cours Théoriques
            </h2>
            <span
              class="inline-flex items-center gap-1 text-[10px] font-bold px-2 py-0.5 rounded-full bg-emerald-500/10 text-emerald-600 border border-emerald-500/20"
            >
              <CheckCircle2 class="w-3 h-3" /> Permis {{ studentCategory }}
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
          <p class="text-xs text-muted-foreground mt-0.5">
            Consultez vos livrets de révision, fiches mémo et écoutez les cours
            en synthèse vocale hors-ligne.
          </p>
        </div>
      </div>

      <div class="relative w-full sm:w-60 shrink-0">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Rechercher un cours..."
          class="w-full h-8 pl-8 pr-3 rounded-lg border border-border bg-background text-foreground text-xs focus:outline-none focus:border-brand-dark"
        />
        <Search
          class="w-3.5 h-3.5 text-muted-foreground absolute left-2.5 top-2.5"
        />
      </div>
    </section>

    <!-- Aucun cours -->
    <div
      v-if="studentSupports.length === 0"
      class="p-10 text-center rounded-xl border border-dashed border-border bg-card space-y-2"
    >
      <FileText class="w-8 h-8 text-muted-foreground/40 mx-auto" />
      <h3 class="text-sm font-semibold text-foreground">
        Aucun support de cours disponible
      </h3>
      <p class="text-xs text-muted-foreground max-w-sm mx-auto">
        Votre auto-école publiera prochainement les documents et liens de
        révision pour le Permis {{ studentCategory }}.
      </p>
    </div>

    <!-- Grille des Supports -->
    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3.5">
      <div
        v-for="s in studentSupports"
        :key="s.id"
        class="rounded-xl bg-card border border-border hover:border-emerald-500/50 transition shadow-xs p-4 flex flex-col justify-between space-y-3"
      >
        <div class="space-y-2">
          <div class="flex items-center justify-between">
            <span
              class="inline-flex items-center gap-1.5 text-[10px] font-bold px-2 py-0.5 rounded-full border uppercase"
              :class="
                s.type_support === 'PDF'
                  ? 'bg-emerald-500/10 text-emerald-600 border-emerald-500/20'
                  : 'bg-sky-500/10 text-sky-600 border-sky-500/20'
              "
            >
              {{ s.type_support === "PDF" ? "📄 Document PDF" : "🔗 Lien Web" }}
            </span>

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
            <span
              v-else
              class="text-[10px] font-mono text-muted-foreground bg-muted px-2 py-0.5 rounded"
            >
              {{ s.fichier_taille ? s.fichier_taille : "PDF" }}
            </span>
          </div>

          <h3
            @click="openSupport(s)"
            class="text-xs sm:text-sm font-bold text-foreground leading-snug line-clamp-2 cursor-pointer hover:text-emerald-600 transition"
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

        <div class="pt-2 border-t border-border space-y-2">
          <div
            class="flex items-center justify-between text-[10px] text-muted-foreground"
          >
            <span>Par {{ s.ajoute_par || "Votre auto-école" }}</span>
            <span class="font-mono"
              >{{ s.telechargements_count || 0 }} consultation(s)</span
            >
          </div>

          <div class="flex items-center gap-2">
            <button
              type="button"
              @click="openSupport(s)"
              class="flex-1 flex items-center justify-center space-x-1.5 px-3 py-2 rounded-lg bg-emerald-600 hover:bg-emerald-700 text-white font-bold text-xs transition cursor-pointer shadow-xs active:scale-[0.99]"
            >
              <component
                :is="s.type_support === 'PDF' ? Eye : ExternalLink"
                class="w-3.5 h-3.5"
              />
              <span>{{
                s.type_support === "PDF" ? "Lire le cours" : "Ouvrir le lien"
              }}</span>
            </button>

            <!-- Bouton Écouter directement la fiche de cours -->
            <button
              type="button"
              @click="ecouterSupport(s)"
              class="p-2 rounded-lg border border-border hover:bg-muted text-muted-foreground hover:text-emerald-600 transition cursor-pointer flex items-center justify-center shrink-0"
              title="Écouter la lecture vocale de ce support"
            >
              <Volume2 class="w-3.5 h-3.5" />
            </button>

            <!-- Téléchargement pour consultation hors-ligne (style YouTube offline) -->
            <button
              type="button"
              @click="offlineStore.downloadSupport(s)"
              class="p-2 rounded-lg border transition cursor-pointer flex items-center justify-center shrink-0"
              :class="
                offlineStore.isDownloaded(s.id)
                  ? 'border-emerald-500/40 bg-emerald-500/10 text-emerald-600'
                  : 'border-border hover:bg-muted text-muted-foreground hover:text-foreground'
              "
              :title="
                offlineStore.isDownloaded(s.id)
                  ? 'Disponible hors-ligne'
                  : 'Télécharger pour révision hors-ligne sans connexion'
              "
            >
              <component
                :is="
                  offlineStore.isDownloaded(s.id) ? CheckCircle2 : DownloadCloud
                "
                class="w-3.5 h-3.5"
              />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Lecteur audio flottant global pour la page des cours -->
    <div
      v-if="isAudioBarOpen && !activeSupportToPreview"
      class="fixed bottom-4 left-4 right-4 sm:left-auto sm:right-6 sm:w-96 z-40 shadow-2xl rounded-xl overflow-hidden border border-border"
    >
      <AudioSpeechBar
        :title="currentAudioTitle"
        :is-extracting="isExtracting"
        @close="fermerAudio"
      />
    </div>

    <!-- Visionneuse PDF & Lecteur Bilingue Avancé Modal -->
    <DocumentReaderModal
      v-if="activeSupportToPreview"
      :is-open="!!activeSupportToPreview"
      :titre="activeSupportToPreview.titre"
      :categorie-permis="activeSupportToPreview.categorie_permis"
      :fichier-data="activeSupportToPreview.fichier_data"
      :description="activeSupportToPreview.description"
      @download="offlineStore.downloadSupport(activeSupportToPreview)"
      @close="fermerModalPreview"
    />
  </div>
</template>

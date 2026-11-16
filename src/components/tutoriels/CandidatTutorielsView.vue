<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { useTutorielsStore } from "../../stores/tutorielsStore";
import { useOfflineDownloadsStore } from "../../stores/offlineDownloadsStore";
import { realtimeSocket } from "../../services/realtimeSocketService";
import VideoPlayerModal from "./VideoPlayerModal.vue";
import type { TutorielVideo } from "../../types";
import {
  Film,
  Play,
  Clock,
  Search,
  Tag,
  Sparkles,
  Eye,
  CheckCircle2,
  Lock,
  DownloadCloud,
  Radio,
} from "lucide-vue-next";

const props = defineProps<{
  studentCategory: string;
}>();

const tutorielsStore = useTutorielsStore();
const offlineStore = useOfflineDownloadsStore();
const isWsConnected = realtimeSocket.isConnected;
let unsubscribeWs: (() => void) | null = null;
const searchQuery = ref("");
const activeVideo = ref<TutorielVideo | null>(null);

onMounted(() => {
  tutorielsStore.loadTutoriels();
  offlineStore.loadDownloads();
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

const studentVideos = computed(() => {
  const base = tutorielsStore.getTutorielsForStudent(props.studentCategory);
  if (!searchQuery.value.trim()) return base;
  const q = searchQuery.value.toLowerCase().trim();
  return base.filter(
    (v) =>
      v.titre.toLowerCase().includes(q) ||
      (v.description && v.description.toLowerCase().includes(q)),
  );
});

const openVideo = (v: TutorielVideo) => {
  activeVideo.value = v;
  tutorielsStore.playVideo(v);
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
          class="w-10 h-10 rounded-lg bg-red-500/10 border border-red-500/20 text-red-600 flex items-center justify-center shrink-0"
        >
          <Film class="w-5 h-5 text-red-600" />
        </div>
        <div>
          <div class="flex items-center gap-2">
            <h2 class="text-sm sm:text-base font-bold text-foreground">
              Tutoriels & Guides Vidéo Pédagogiques
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
            Visionnez les vidéos d'apprentissage, guides de manœuvres et
            conseils de vos moniteurs pour préparer vos cours pratiques et votre
            examen.
          </p>
        </div>
      </div>

      <div class="relative w-full sm:w-60 shrink-0">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Rechercher une vidéo..."
          class="w-full h-8 pl-8 pr-3 rounded-lg border border-border bg-background text-foreground text-xs focus:outline-none focus:border-brand-dark"
        />
        <Search
          class="w-3.5 h-3.5 text-muted-foreground absolute left-2.5 top-2.5"
        />
      </div>
    </section>

    <!-- Aucune vidéo -->
    <div
      v-if="studentVideos.length === 0"
      class="p-10 text-center rounded-xl border border-dashed border-border bg-card space-y-2"
    >
      <Film class="w-8 h-8 text-muted-foreground/40 mx-auto" />
      <h3 class="text-sm font-semibold text-foreground">
        Aucun tutoriel vidéo disponible pour le moment
      </h3>
      <p class="text-xs text-muted-foreground max-w-sm mx-auto">
        Vos formateurs publient régulièrement de nouveaux guides vidéos pour la
        catégorie Permis {{ studentCategory }}.
      </p>
    </div>

    <!-- Grille des Vidéos -->
    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3.5">
      <div
        v-for="v in studentVideos"
        :key="v.id"
        class="rounded-xl bg-card border border-border hover:border-red-500/50 transition shadow-xs flex flex-col justify-between overflow-hidden group"
      >
        <!-- Miniature cliquable -->
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

          <div
            class="absolute inset-0 flex items-center justify-center bg-black/20 group-hover:bg-black/40 transition"
          >
            <div
              class="w-10 h-10 rounded-full bg-red-600 text-white flex items-center justify-center shadow-lg group-hover:scale-110 transition"
            >
              <Play class="w-4 h-4 fill-current ml-0.5" />
            </div>
          </div>

          <div
            class="absolute bottom-2 right-2 text-[10px] font-mono font-bold px-1.5 py-0.5 rounded bg-black/80 text-white flex items-center gap-1"
          >
            <Clock class="w-3 h-3 text-red-400" />
            {{ v.duree_minutes ? v.duree_minutes + " min" : "Vidéo" }}
          </div>
        </div>

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
            <span>Par {{ v.ajoute_par || "Votre auto-école" }}</span>
            <span class="flex items-center gap-1 font-mono">
              <Eye class="w-3 h-3" /> {{ v.vues_count || 0 }} consultation(s)
            </span>
          </div>

          <div class="flex items-center gap-2">
            <button
              type="button"
              @click="openVideo(v)"
              class="flex-1 flex items-center justify-center space-x-1.5 px-3 py-2 rounded-lg bg-red-600 hover:bg-red-700 text-white font-bold text-xs transition cursor-pointer shadow-xs active:scale-[0.99]"
            >
              <Play class="w-3.5 h-3.5 fill-current" />
              <span>Regarder</span>
            </button>

            <!-- Bouton téléchargement hors-ligne (style YouTube offline) -->
            <button
              type="button"
              @click="offlineStore.downloadTutoriel(v)"
              class="p-2 rounded-lg border transition cursor-pointer flex items-center justify-center shrink-0"
              :class="
                offlineStore.isDownloaded(v.id)
                  ? 'border-emerald-500/40 bg-emerald-500/10 text-emerald-600'
                  : 'border-border hover:bg-muted text-muted-foreground hover:text-foreground'
              "
              :title="
                offlineStore.isDownloaded(v.id)
                  ? 'Disponible hors-ligne'
                  : 'Conserver pour révision hors-ligne sans connexion'
              "
            >
              <component
                :is="
                  offlineStore.isDownloaded(v.id) ? CheckCircle2 : DownloadCloud
                "
                class="w-3.5 h-3.5"
              />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Lecteur Vidéo Modal -->
    <VideoPlayerModal :video="activeVideo" @close="activeVideo = null" />
  </div>
</template>

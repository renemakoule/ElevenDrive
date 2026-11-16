<script setup lang="ts">
import { computed } from 'vue';
import type { TutorielVideo } from '../../types';
import { X, ExternalLink, Play, Clock, User, ShieldCheck } from 'lucide-vue-next';

const props = defineProps<{
  video: TutorielVideo | null;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
}>();

/**
 * Normalise l'URL pour un affichage parfait dans un iframe (YouTube, Vimeo, MP4)
 */
const videoEmbedInfo = computed(() => {
  if (!props.video?.url_video) return null;
  const url = props.video.url_video.trim();

  // 1. YouTube
  if (url.includes('youtube.com') || url.includes('youtu.be')) {
    let videoId = '';
    if (url.includes('youtu.be/')) {
      videoId = url.split('youtu.be/')[1].split('?')[0];
    } else if (url.includes('watch?v=')) {
      videoId = url.split('watch?v=')[1].split('&')[0];
    } else if (url.includes('/embed/')) {
      videoId = url.split('/embed/')[1].split('?')[0];
    }
    return {
      type: 'IFRAME',
      src: `https://www.youtube-nocookie.com/embed/${videoId}?autoplay=1&rel=0`,
    };
  }

  // 2. Vimeo
  if (url.includes('vimeo.com')) {
    const parts = url.split('vimeo.com/');
    const id = parts[1]?.split('?')[0];
    return {
      type: 'IFRAME',
      src: `https://player.vimeo.com/video/${id}?autoplay=1`,
    };
  }

  // 3. Fichier vidéo MP4 direct ou WebM
  if (url.endsWith('.mp4') || url.endsWith('.webm') || url.endsWith('.ogg')) {
    return {
      type: 'HTML5_VIDEO',
      src: url,
    };
  }

  // 4. Par défaut : tentative iframe
  return {
    type: 'IFRAME',
    src: url,
  };
});

const openExternal = () => {
  if (props.video?.url_video) {
    window.open(props.video.url_video, '_blank', 'noopener,noreferrer');
  }
};
</script>

<template>
  <div
    v-if="video"
    class="fixed inset-0 z-50 flex items-center justify-center p-2 sm:p-4 bg-black/85 backdrop-blur-sm animate-in fade-in"
    role="dialog"
    aria-modal="true"
  >
    <div class="bg-card border border-border w-full max-w-4xl max-h-[94vh] rounded-[5px] shadow-2xl flex flex-col overflow-hidden text-foreground">
      <!-- En-tête -->
      <header class="px-4 py-3 border-b border-border flex items-center justify-between bg-muted/30">
        <div class="flex items-center space-x-2.5 min-w-0">
          <span class="inline-flex items-center gap-1.5 text-[10px] font-bold px-2 py-0.5 rounded-[4px] bg-red-500/10 text-red-600 border border-red-500/20 uppercase tracking-wider">
            <Play class="w-3 h-3 fill-current" />
            Tutoriel Vidéo
          </span>
          <div class="min-w-0">
            <h2 class="text-xs sm:text-sm font-bold text-foreground truncate">
              {{ video.titre }}
            </h2>
            <p class="text-[10px] text-muted-foreground">
              Permis <strong>{{ video.categorie_permis }}</strong> · {{ video.duree_minutes ? video.duree_minutes + ' min' : 'Vidéo de formation' }}
            </p>
          </div>
        </div>

        <div class="flex items-center space-x-1.5">
          <button
            type="button"
            @click="openExternal"
            class="flex items-center space-x-1 px-2.5 py-1 rounded-[4px] border border-border hover:bg-muted text-xs font-semibold text-foreground transition cursor-pointer"
            title="Ouvrir la vidéo sur la plateforme d'origine"
          >
            <ExternalLink class="w-3.5 h-3.5" />
            <span class="hidden sm:inline">Plein écran web</span>
          </button>

          <button
            type="button"
            @click="$emit('close')"
            class="p-1.5 rounded-[4px] hover:bg-muted text-muted-foreground hover:text-foreground transition cursor-pointer"
          >
            <X class="w-4 h-4" />
          </button>
        </div>
      </header>

      <!-- Lecteur Vidéo (Aspect 16:9 responsive) -->
      <div class="relative bg-black w-full aspect-video flex items-center justify-center overflow-hidden">
        <template v-if="videoEmbedInfo">
          <!-- Mode Iframe Embed (YouTube / Vimeo) -->
          <iframe
            v-if="videoEmbedInfo.type === 'IFRAME'"
            :src="videoEmbedInfo.src"
            class="w-full h-full border-0"
            allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
            allowfullscreen
          />

          <!-- Mode Lecteur HTML5 Natif (MP4) -->
          <video
            v-else-if="videoEmbedInfo.type === 'HTML5_VIDEO'"
            :src="videoEmbedInfo.src"
            controls
            autoplay
            class="w-full h-full object-contain"
          >
            Votre navigateur ne supporte pas la balise vidéo.
          </video>
        </template>
      </div>

      <!-- Informations & Description Pédagogique -->
      <footer class="p-4 border-t border-border bg-card space-y-2 text-xs">
        <div class="flex items-center justify-between text-[11px] text-muted-foreground">
          <span class="flex items-center gap-1.5 font-medium text-foreground">
            <User class="w-3.5 h-3.5 text-brand" /> Ajouté par : {{ video.ajoute_par || 'Direction Pédagogique' }}
          </span>
          <span class="flex items-center gap-1">
            <Clock class="w-3.5 h-3.5" /> {{ video.vues_count || 0 }} consultation(s)
          </span>
        </div>

        <p v-if="video.description" class="text-muted-foreground leading-relaxed">
          {{ video.description }}
        </p>
      </footer>
    </div>
  </div>
</template>

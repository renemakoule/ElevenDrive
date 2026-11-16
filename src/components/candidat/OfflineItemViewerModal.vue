<script setup lang="ts">
import { computed } from 'vue';
import type { OfflineItem } from '../../services/offlineStorageService';
import { openExternalLink } from '../../utils/urlHelper';
import {
  X,
  CheckCircle2,
  FileText,
  Film,
  ExternalLink,
  BookOpen,
  WifiOff,
  Wifi,
  Sparkles
} from 'lucide-vue-next';

const props = defineProps<{
  item: OfflineItem | null;
  isOnline: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
}>();

const isPdf = computed(() => {
  return props.item?.mime_type === 'application/pdf' || props.item?.blob_or_data?.startsWith('data:application/pdf');
});

const isNativeVideo = computed(() => {
  if (!props.item || props.item.type !== 'TUTORIEL') return false;
  const d = props.item.blob_or_data || '';
  return d.startsWith('data:video') || d.endsWith('.mp4') || d.endsWith('.webm');
});

const openExternal = (url?: string) => {
  if (url) {
    openExternalLink(url);
  }
};
</script>

<template>
  <div
    v-if="item"
    class="fixed inset-0 z-50 flex items-center justify-center p-2 sm:p-4 bg-black/85 backdrop-blur-sm animate-in fade-in"
    role="dialog"
    aria-modal="true"
  >
    <div class="bg-card border border-border w-full max-w-4xl h-[92vh] rounded-[5px] flex flex-col overflow-hidden text-foreground">
      <!-- En-tête -->
      <header class="px-4 py-2.5 border-b border-border flex items-center justify-between bg-muted/40">
        <div class="flex items-center space-x-2.5 min-w-0">
          <span class="inline-flex items-center gap-1 text-[10px] font-bold px-2 py-0.5 rounded-[4px] bg-emerald-500/10 text-emerald-600 border border-emerald-500/20">
            <CheckCircle2 class="w-3 h-3" />
            Lecture Hors-Ligne Sécurisée
          </span>

          <span class="text-[10px] font-bold px-2 py-0.5 rounded-[4px] bg-brand-light text-brand-dark border border-brand-border uppercase">
            Permis {{ item.categorie_permis }}
          </span>

          <h2 class="text-xs sm:text-sm font-bold text-foreground truncate max-w-sm">
            {{ item.titre }}
          </h2>
        </div>

        <div class="flex items-center space-x-2">
          <!-- Statut de connectivité -->
          <span
            class="hidden sm:inline-flex items-center gap-1 text-[10px] px-2 py-0.5 rounded-full border"
            :class="isOnline ? 'bg-sky-500/10 text-sky-600 border-sky-500/20' : 'bg-amber-500/10 text-amber-600 border-amber-500/20'"
          >
            <component :is="isOnline ? Wifi : WifiOff" class="w-3 h-3" />
            {{ isOnline ? 'Connecté' : 'Sans connexion (100% Hors-Ligne)' }}
          </span>

          <button
            type="button"
            @click="$emit('close')"
            class="p-1 rounded-[4px] text-muted-foreground hover:text-foreground transition cursor-pointer"
            title="Fermer la visionneuse"
          >
            <X class="w-5 h-5" />
          </button>
        </div>
      </header>

      <!-- Corps du visualiseur -->
      <div class="flex-1 bg-zinc-950 overflow-hidden relative flex flex-col justify-center">
        <!-- 1. CAS PDF (Épreuve ou Support de cours) -->
        <template v-if="isPdf">
          <iframe
            v-if="item.blob_or_data"
            :src="item.blob_or_data"
            class="w-full h-full border-none"
            title="Visionneuse Document Hors-Ligne"
          />
          <div v-else class="flex-1 flex flex-col items-center justify-center p-8 text-center text-muted-foreground">
            <FileText class="w-12 h-12 mb-3 text-muted-foreground/40" />
            <p class="text-sm font-semibold text-foreground">Fichier PDF introuvable en mémoire</p>
          </div>
        </template>

        <!-- 2. CAS VIDÉO NATIVE (MP4 / WebM stocké hors-ligne) -->
        <template v-else-if="isNativeVideo">
          <div class="flex-1 flex items-center justify-center p-4">
            <video
              controls
              autoplay
              playsinline
              class="w-full max-h-full rounded-[4px] bg-black"
              :src="item.blob_or_data"
            >
              Votre navigateur ne prend pas en charge la balise vidéo.
            </video>
          </div>
        </template>

        <!-- 3. CAS LIEN EXTERNE OU YOUTUBE HORS-LIGNE -->
        <template v-else>
          <div class="flex-1 flex flex-col items-center justify-center p-6 text-center max-w-lg mx-auto space-y-4">
            <div class="w-16 h-16 rounded-full bg-red-500/10 border border-red-500/20 text-red-500 flex items-center justify-center">
              <Film class="w-8 h-8" />
            </div>

            <div class="space-y-1">
              <h3 class="text-base font-bold text-white">{{ item.titre }}</h3>
              <p class="text-xs text-zinc-400">
                Guide pédagogique enregistré hors-ligne · Permis {{ item.categorie_permis }}
              </p>
            </div>

            <div v-if="item.description" class="p-3.5 rounded bg-zinc-900 border border-zinc-800 text-left text-xs text-zinc-300 w-full space-y-1">
              <strong class="text-zinc-100 block">Notes & Résumé Pédagogique :</strong>
              <p>{{ item.description }}</p>
            </div>

            <div class="text-[11px] text-zinc-400 bg-zinc-900/80 p-3 rounded border border-zinc-800 flex items-start gap-2 text-left">
              <Sparkles class="w-4 h-4 text-brand shrink-0 mt-0.5" />
              <span>
                Ce tutoriel utilise le lecteur sécurisé en ligne. Ses métadonnées, notes et résumés sont disponibles 100% hors-ligne. Une connexion Internet est requise pour le streaming direct du flux YouTube.
              </span>
            </div>

            <button
              v-if="item.url_externe"
              type="button"
              @click="openExternal(item.url_externe)"
              class="inline-flex items-center gap-1.5 px-4 py-2 rounded bg-red-600 hover:bg-red-700 text-white font-bold text-xs transition"
            >
              <ExternalLink class="w-3.5 h-3.5" />
              <span>Ouvrir sur la plateforme (si connecté)</span>
            </button>
          </div>
        </template>
      </div>

      <!-- Pied de page -->
      <footer class="px-4 py-2 bg-card border-t border-border flex items-center justify-between text-[11px] text-muted-foreground">
        <span>Taille stockée en local : <strong>{{ item.taille_formatee }}</strong></span>
        <span>Téléchargé le {{ new Date(item.date_telechargement).toLocaleDateString('fr-FR') }}</span>
      </footer>
    </div>
  </div>
</template>

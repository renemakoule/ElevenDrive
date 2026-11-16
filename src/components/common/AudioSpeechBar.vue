<script setup lang="ts">
import { ref, computed } from 'vue';
import { speechAudioService } from '../../services/speechAudioService';
import {
  Play,
  Pause,
  Square,
  SkipBack,
  SkipForward,
  Volume2,
  Volume1,
  VolumeX,
  X,
  Gauge,
  Sparkles,
  ChevronUp,
  ChevronDown,
} from 'lucide-vue-next';

const props = withDefaults(
  defineProps<{
    title?: string;
    isExtracting?: boolean;
    compact?: boolean;
  }>(),
  {
    title: 'Lecture Vocale',
    isExtracting: false,
    compact: false,
  }
);

const emit = defineEmits<{
  (e: 'close'): void;
}>();

const isExpanded = ref(true);
const vitessesDisponibles = [0.8, 1.0, 1.25, 1.5];

/* ---------- Volume ---------- */
const showVolume = ref(false);
const volume = ref<number>(
  (speechAudioService as any).currentVolume?.value ?? 1
);
const lastVolume = ref(1);

const VolumeIcon = computed(() => {
  if (volume.value === 0) return VolumeX;
  if (volume.value < 0.5) return Volume1;
  return Volume2;
});

const applyVolume = (v: number) => {
  volume.value = Math.min(1, Math.max(0, v));
  (speechAudioService as any).setVolume?.(volume.value);
};

const onVolumeInput = (e: Event) => {
  applyVolume(parseFloat((e.target as HTMLInputElement).value));
};

const toggleMute = () => {
  if (volume.value > 0) {
    lastVolume.value = volume.value;
    applyVolume(0);
  } else {
    applyVolume(lastVolume.value || 1);
  }
};

/* ---------- Progression fluide ---------- */
// Fraction à l'intérieur du segment courant (0 → 1), si le service l'expose
const segmentFraction = computed(() => {
  const f = (speechAudioService as any).segmentProgress?.value ?? 0;
  return Math.min(1, Math.max(0, f));
});

const progression = computed(() => {
  const total = speechAudioService.totalSegments.value;
  if (total <= 0) return 0;
  const done = speechAudioService.currentSegmentIndex.value + segmentFraction.value;
  return Math.min(100, (done / total) * 100);
});

/* ---------- Actions ---------- */
const handlePlayPause = () => {
  speechAudioService.basculerPlayPause();
};

// Stop : arrête la lecture immédiatement (sans fermer le lecteur)
const handleStopPlayback = () => {
  showVolume.value = false;
  speechAudioService.arreter();
};

// Fermer : arrête puis ferme la barre
const handleClose = () => {
  speechAudioService.arreter();
  emit('close');
};

const cycleSpeed = () => {
  const current = speechAudioService.currentRate.value;
  const nextIdx = (vitessesDisponibles.indexOf(current) + 1) % vitessesDisponibles.length;
  speechAudioService.setVitesse(vitessesDisponibles[nextIdx]);
};
</script>

<template>
  <div
    class="relative w-full bg-card/95 backdrop-blur-md border-t border-border transition-all duration-200 select-none text-foreground"
    :class="compact ? 'py-1.5 px-3' : 'py-2.5 px-4'"
  >
    <!-- Indicateur d'extraction en cours -->
    <div
      v-if="isExtracting"
      class="flex items-center justify-center space-x-2 py-1 text-xs text-brand font-medium animate-pulse"
    >
      <Sparkles class="w-4 h-4 animate-spin text-brand" />
      <span>Préparation de la lecture audio du document...</span>
    </div>

    <div v-else class="flex flex-col space-y-2">
      <div class="flex items-center justify-between gap-2 sm:gap-4">
        <!-- Titre, statut et volume -->
        <div class="flex items-center space-x-2.5 min-w-0 flex-1">
          <div class="relative shrink-0">
            <!-- Icône volume cliquable -->
            <button
              type="button"
              @click="showVolume = !showVolume"
              class="w-8 h-8 rounded-[5px] flex items-center justify-center border transition-colors cursor-pointer"
              :class="
                speechAudioService.isPlaying.value && !speechAudioService.isPaused.value
                  ? 'bg-brand/10 border-brand/30 text-brand'
                  : 'bg-muted border-border text-muted-foreground'
              "
              title="Régler le volume"
            >
              <div
                v-if="speechAudioService.isPlaying.value && !speechAudioService.isPaused.value && volume > 0"
                class="flex items-end gap-0.5 h-3.5"
              >
                <span class="w-0.5 bg-brand h-full animate-[pulse_0.6s_ease-in-out_infinite]" />
                <span class="w-0.5 bg-brand h-2 animate-[pulse_0.4s_ease-in-out_infinite_0.1s]" />
                <span class="w-0.5 bg-brand h-3 animate-[pulse_0.7s_ease-in-out_infinite_0.2s]" />
                <span class="w-0.5 bg-brand h-1.5 animate-[pulse_0.5s_ease-in-out_infinite_0.3s]" />
              </div>
              <component :is="VolumeIcon" v-else class="w-4 h-4" />
            </button>

            <!-- Fond pour fermer en cliquant à l'extérieur -->
            <div v-if="showVolume" class="fixed inset-0 z-40" @click="showVolume = false" />

            <!-- Popover volume -->
            <div
              v-if="showVolume"
              class="absolute bottom-full left-0 mb-2 z-50 flex items-center gap-2 px-3 py-2 rounded-lg bg-card border border-border shadow-lg"
            >
              <button
                type="button"
                @click="toggleMute"
                class="p-1 rounded hover:bg-muted text-muted-foreground hover:text-foreground cursor-pointer"
                :title="volume === 0 ? 'Réactiver le son' : 'Couper le son'"
              >
                <component :is="VolumeIcon" class="w-4 h-4" />
              </button>
              <input
                type="range"
                min="0"
                max="1"
                step="0.05"
                :value="volume"
                @input="onVolumeInput"
                class="w-28 h-1 cursor-pointer accent-[var(--color-brand,theme(colors.brand))]"
                aria-label="Volume"
              />
              <span class="text-[11px] font-mono w-8 text-right text-muted-foreground">
                {{ Math.round(volume * 100) }}%
              </span>
            </div>
          </div>

          <div class="min-w-0">
            <div class="flex items-center gap-2">
              <span class="text-xs font-bold text-foreground truncate max-w-[200px] sm:max-w-xs">
                {{ title }}
              </span>
              <span
                v-if="speechAudioService.totalSegments.value > 0"
                class="text-[10px] font-mono px-1.5 py-0.5 rounded bg-muted text-muted-foreground border border-border shrink-0"
              >
                {{ speechAudioService.currentSegmentIndex.value + 1 }} /
                {{ speechAudioService.totalSegments.value }}
              </span>
            </div>
            <p
              v-if="speechAudioService.currentSegmentText.value"
              class="text-[11px] text-muted-foreground truncate max-w-sm sm:max-w-md hidden sm:block italic"
            >
              « {{ speechAudioService.currentSegmentText.value }} »
            </p>
          </div>
        </div>

        <!-- Contrôles -->
        <div class="flex items-center space-x-1 sm:space-x-1.5 shrink-0">
          <button
            type="button"
            @click="speechAudioService.precedent()"
            :disabled="speechAudioService.currentSegmentIndex.value <= 0"
            class="p-1.5 rounded-md hover:bg-muted text-muted-foreground hover:text-foreground disabled:opacity-30 disabled:pointer-events-none transition cursor-pointer"
            title="Phrase précédente"
          >
            <SkipBack class="w-4 h-4" />
          </button>

          <button
            type="button"
            @click="handlePlayPause"
            class="flex items-center justify-center w-8 h-8 rounded-full bg-brand hover:bg-brand-dark text-white shadow-xs transition active:scale-95 cursor-pointer"
            :title="speechAudioService.isPaused.value ? 'Reprendre' : 'Mettre en pause'"
          >
            <Play
              v-if="speechAudioService.isPaused.value || !speechAudioService.isPlaying.value"
              class="w-4 h-4 fill-current ml-0.5"
            />
            <Pause v-else class="w-4 h-4 fill-current" />
          </button>

          <button
            type="button"
            @click="speechAudioService.suivant()"
            :disabled="
              speechAudioService.currentSegmentIndex.value >=
              speechAudioService.totalSegments.value - 1
            "
            class="p-1.5 rounded-md hover:bg-muted text-muted-foreground hover:text-foreground disabled:opacity-30 disabled:pointer-events-none transition cursor-pointer"
            title="Phrase suivante"
          >
            <SkipForward class="w-4 h-4" />
          </button>

          <!-- Stop -->
          <button
            type="button"
            @click="handleStopPlayback"
            :disabled="!speechAudioService.isPlaying.value"
            class="p-1.5 rounded-md hover:bg-muted text-muted-foreground hover:text-destructive disabled:opacity-30 disabled:pointer-events-none transition cursor-pointer"
            title="Arrêter la lecture"
          >
            <Square class="w-4 h-4" />
          </button>

          <button
            type="button"
            @click="cycleSpeed"
            class="flex items-center gap-1 px-2 py-1 rounded-md bg-muted hover:bg-muted/80 text-[11px] font-bold text-foreground border border-border transition cursor-pointer"
            title="Changer la vitesse d'élocution"
          >
            <Gauge class="w-3 h-3 text-muted-foreground" />
            <span>{{ speechAudioService.currentRate.value }}x</span>
          </button>

          <button
            type="button"
            @click="isExpanded = !isExpanded"
            class="p-1.5 rounded-md hover:bg-muted text-muted-foreground hover:text-foreground transition cursor-pointer hidden sm:inline-flex"
            :title="isExpanded ? 'Masquer la transcription' : 'Afficher la transcription'"
          >
            <component :is="isExpanded ? ChevronDown : ChevronUp" class="w-4 h-4" />
          </button>

          <button
            type="button"
            @click="handleClose"
            class="p-1.5 rounded-md hover:bg-muted text-muted-foreground hover:text-foreground transition cursor-pointer ml-1"
            title="Fermer le lecteur vocal"
          >
            <X class="w-4 h-4" />
          </button>
        </div>
      </div>

      <!-- Barre de progression fluide -->
      <div class="w-full h-1 bg-muted rounded-full overflow-hidden">
        <div
          class="h-full w-full bg-brand origin-left will-change-transform transition-transform duration-500 ease-linear"
          :style="{ transform: `scaleX(${progression / 100})` }"
        />
      </div>

      <!-- Transcription -->
      <div
        v-if="isExpanded && speechAudioService.currentSegmentText.value"
        class="bg-muted/30 border border-border rounded-[5px] p-2.5 text-xs text-foreground leading-relaxed animate-in fade-in flex items-start gap-2"
      >
        <span class="text-[10px] font-bold px-1.5 py-0.5 rounded bg-brand/10 text-brand border border-brand/20 shrink-0">
          En cours
        </span>
        <p class="font-medium text-foreground">
          {{ speechAudioService.currentSegmentText.value }}
        </p>
      </div>
    </div>
  </div>
</template>
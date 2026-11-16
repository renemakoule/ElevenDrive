<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from "vue";
import {
  X,
  ExternalLink,
  ZoomIn,
  ZoomOut,
  RotateCw,
  Download,
  Check,
  ChevronLeft,
  ChevronRight,
  Loader2,
} from "lucide-vue-next";

export interface PreviewImageData {
  url: string;
  title: string;
  category?: string;
  content?: string;
  date?: string;
  actionUrl?: string;
}

const props = withDefaults(
  defineProps<{
    image: PreviewImageData | null;
    currentIndex?: number;
    totalImages?: number;
  }>(),
  {
    currentIndex: 0,
    totalImages: 1,
  },
);

const emit = defineEmits<{
  (e: "close"): void;
  (e: "prev"): void;
  (e: "next"): void;
}>();

const zoomScale = ref(1);
const rotationDeg = ref(0);
const isDownloading = ref(false);
const downloadSuccess = ref(false);
let downloadSuccessTimer: any = null;

const resetZoom = () => {
  zoomScale.value = 1;
  rotationDeg.value = 0;
};

const zoomIn = () => {
  zoomScale.value = Math.min(3, +(zoomScale.value + 0.25).toFixed(2));
};

const zoomOut = () => {
  zoomScale.value = Math.max(0.5, +(zoomScale.value - 0.25).toFixed(2));
};

const rotateImage = () => {
  rotationDeg.value = (rotationDeg.value + 90) % 360;
};

const handleClose = () => {
  resetZoom();
  isDownloading.value = false;
  downloadSuccess.value = false;
  if (downloadSuccessTimer) clearTimeout(downloadSuccessTimer);
  emit("close");
};

const downloadCurrentImage = async () => {
  if (!props.image?.url || isDownloading.value) return;

  isDownloading.value = true;
  downloadSuccess.value = false;
  if (downloadSuccessTimer) clearTimeout(downloadSuccessTimer);

  try {
    const rawUrl = props.image.url;
    const response = await fetch(rawUrl, { mode: "cors" });
    if (!response.ok) throw new Error("Échec du téléchargement");

    const blob = await response.blob();
    const blobUrl = URL.createObjectURL(blob);

    let extension = "png";
    const contentType = response.headers.get("content-type");
    if (contentType?.includes("jpeg") || contentType?.includes("jpg")) {
      extension = "jpg";
    } else if (contentType?.includes("webp")) {
      extension = "webp";
    } else if (contentType?.includes("svg")) {
      extension = "svg";
    }

    const cleanTitle = (props.image.title || "image-annonce").replace(
      /[^a-zA-Z0-9_-]/g,
      "_",
    );
    const filename = `${cleanTitle}.${extension}`;

    const link = document.createElement("a");
    link.href = blobUrl;
    link.download = filename;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);

    setTimeout(() => URL.revokeObjectURL(blobUrl), 1500);

    downloadSuccess.value = true;
    downloadSuccessTimer = setTimeout(() => {
      downloadSuccess.value = false;
    }, 2500);
  } catch {
    // Fallback direct
    const link = document.createElement("a");
    link.href = props.image.url;
    link.download = `${props.image.title || "image-annonce"}.png`;
    link.target = "_blank";
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);

    downloadSuccess.value = true;
    downloadSuccessTimer = setTimeout(() => {
      downloadSuccess.value = false;
    }, 2500);
  } finally {
    isDownloading.value = false;
  }
};

const handleKeyDown = (e: KeyboardEvent) => {
  if (!props.image) return;

  if (e.key === "Escape") {
    handleClose();
  } else if (e.key === "ArrowRight") {
    emit("next");
  } else if (e.key === "ArrowLeft") {
    emit("prev");
  } else if (e.key === "+" || e.key === "=") {
    zoomIn();
  } else if (e.key === "-" || e.key === "_") {
    zoomOut();
  } else if (e.key.toLowerCase() === "r") {
    rotateImage();
  }
};

watch(
  () => props.image?.url,
  () => {
    resetZoom();
    isDownloading.value = false;
    downloadSuccess.value = false;
  },
);

onMounted(() => {
  if (typeof window !== "undefined") {
    window.addEventListener("keydown", handleKeyDown);
  }
});

onUnmounted(() => {
  if (typeof window !== "undefined") {
    window.removeEventListener("keydown", handleKeyDown);
  }
  if (downloadSuccessTimer) clearTimeout(downloadSuccessTimer);
});
</script>

<template>
  <Teleport to="body">
    <div
      v-if="image"
      class="fixed inset-0 z-[9999] bg-[#07080c]/98 backdrop-blur-2xl flex flex-col justify-between p-4 sm:p-6 select-none animate-in fade-in duration-200"
      @click.self="handleClose"
    >
      <!-- 1. Barre supérieure épurée -->
      <div
        class="w-full flex items-center justify-between text-white shrink-0 z-50"
      >
        <div class="flex items-center gap-3">
          <span
            v-if="image.category"
            class="px-2.5 py-0.5 rounded-[3px] text-[10px] uppercase font-bold bg-brand text-white tracking-wider"
          >
            {{ image.category }}
          </span>

          <div class="flex flex-col">
            <span
              class="font-bold text-sm sm:text-base leading-tight text-white/95"
            >
              {{ image.title }}
            </span>
            <span class="text-[11px] text-white/50 font-mono">
              <template v-if="totalImages > 1">
                {{ currentIndex + 1 }} sur {{ totalImages }} ·
              </template>
              {{ image.date }}
            </span>
          </div>
        </div>

        <div class="flex items-center gap-2">
          <!-- Lien d'action si disponible -->
          <a
            v-if="image.actionUrl"
            :href="image.actionUrl"
            target="_blank"
            rel="noopener"
            class="hidden sm:inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-white/10 hover:bg-white/20 text-white text-xs font-semibold backdrop-blur-md transition cursor-pointer"
          >
            <span>Ouvrir l'annonce</span>
            <ExternalLink class="w-3.5 h-3.5" />
          </a>

          <!-- Bouton Fermer -->
          <button
            type="button"
            @click="handleClose"
            class="p-2 rounded-full bg-white/10 hover:bg-white/25 text-white/90 hover:text-white transition backdrop-blur-md cursor-pointer"
            title="Fermer la visionneuse (Échap)"
          >
            <X class="w-5 h-5" />
          </button>
        </div>
      </div>

      <!-- 2. Zone d'affichage centrale -->
      <div
        class="relative flex-1 flex items-center justify-center overflow-hidden my-3"
        @click.self="handleClose"
      >
        <!-- Flèche Précédent -->
        <button
          v-if="totalImages > 1"
          type="button"
          @click.stop="emit('prev')"
          class="absolute left-2 sm:left-6 z-40 p-2.5 rounded-full bg-black/40 hover:bg-black/75 text-white/80 hover:text-white backdrop-blur-md transition cursor-pointer border border-white/10 shadow-lg"
          title="Image précédente (Flèche gauche)"
        >
          <ChevronLeft class="w-6 h-6" />
        </button>

        <!-- Fond d'ambiance flouté (Glow naturel) -->
        <img
          :src="image.url"
          aria-hidden="true"
          class="absolute max-w-[85vw] max-h-[75vh] object-contain blur-3xl opacity-25 scale-125 pointer-events-none transition-all duration-300"
        />

        <!-- Image principale avec Zoom et Rotation fluides -->
        <div
          class="relative z-20 flex items-center justify-center max-w-full max-h-full transition-transform duration-150 ease-out"
          :style="{
            transform: `scale(${zoomScale}) rotate(${rotationDeg}deg)`,
          }"
        >
          <img
            :src="image.url"
            :alt="image.title"
            class="max-w-[90vw] max-h-[72vh] object-contain rounded-[4px] shadow-2xl drop-shadow-[0_20px_35px_rgba(0,0,0,0.85)] cursor-zoom-in"
            @dblclick="resetZoom"
            title="Double-cliquer pour réinitialiser le zoom"
          />
        </div>

        <!-- Flèche Suivant -->
        <button
          v-if="totalImages > 1"
          type="button"
          @click.stop="emit('next')"
          class="absolute right-2 sm:right-6 z-40 p-2.5 rounded-full bg-black/40 hover:bg-black/75 text-white/80 hover:text-white backdrop-blur-md transition cursor-pointer border border-white/10 shadow-lg"
          title="Image suivante (Flèche droite)"
        >
          <ChevronRight class="w-6 h-6" />
        </button>
      </div>

      <!-- 3. Barre flottante de commandes minimaliste (Bottom Pill Bar) -->
      <div class="w-full flex flex-col items-center gap-2 shrink-0 z-50">
        <div
          class="flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-black/60 backdrop-blur-xl border border-white/15 shadow-2xl text-white text-xs"
        >
          <!-- Bouton Image Précédente si plusieurs -->
          <button
            v-if="totalImages > 1"
            type="button"
            @click.stop="emit('prev')"
            class="p-2 rounded-full hover:bg-white/15 text-white/80 hover:text-white transition cursor-pointer"
            title="Image précédente (←)"
          >
            <ChevronLeft class="w-4 h-4" />
          </button>

          <!-- Zoom Out -->
          <button
            type="button"
            @click.stop="zoomOut"
            class="p-2 rounded-full hover:bg-white/15 text-white/80 hover:text-white transition cursor-pointer"
            title="Zoom arrière (-)"
          >
            <ZoomOut class="w-4 h-4" />
          </button>

          <!-- Reset 100% -->
          <button
            type="button"
            @click.stop="resetZoom"
            class="px-2.5 py-1 rounded-full text-xs font-mono font-bold hover:bg-white/15 text-white/90 hover:text-white transition cursor-pointer"
            title="Réinitialiser (100%)"
          >
            {{ Math.round(zoomScale * 100) }}%
          </button>

          <!-- Zoom In -->
          <button
            type="button"
            @click.stop="zoomIn"
            class="p-2 rounded-full hover:bg-white/15 text-white/80 hover:text-white transition cursor-pointer"
            title="Zoom avant (+)"
          >
            <ZoomIn class="w-4 h-4" />
          </button>

          <div class="h-4 w-px bg-white/20 mx-1"></div>

          <!-- Rotation -->
          <button
            type="button"
            @click.stop="rotateImage"
            class="p-2 rounded-full hover:bg-white/15 text-white/80 hover:text-white transition cursor-pointer"
            title="Pivoter de 90° (R)"
          >
            <RotateCw class="w-4 h-4" />
          </button>

          <!-- Téléchargement avec indicateur réactif et validation -->
          <button
            type="button"
            :disabled="isDownloading"
            @click.stop="downloadCurrentImage"
            class="p-2 rounded-full transition cursor-pointer flex items-center justify-center relative"
            :class="[
              downloadSuccess
                ? 'bg-emerald-500/25 text-emerald-300 ring-1 ring-emerald-400/50'
                : isDownloading
                  ? 'bg-white/20 text-white cursor-wait'
                  : 'hover:bg-white/15 text-white/80 hover:text-white',
            ]"
            :title="
              downloadSuccess
                ? 'Téléchargement terminé !'
                : isDownloading
                  ? 'Téléchargement en cours...'
                  : 'Télécharger l\'image'
            "
          >
            <Loader2
              v-if="isDownloading"
              class="w-4 h-4 animate-spin text-cyan-300"
            />
            <Check
              v-else-if="downloadSuccess"
              class="w-4 h-4 text-emerald-400 animate-in zoom-in-75 duration-200"
            />
            <Download v-else class="w-4 h-4" />
          </button>

          <!-- Bouton Image Suivante si plusieurs -->
          <button
            v-if="totalImages > 1"
            type="button"
            @click.stop="emit('next')"
            class="p-2 rounded-full hover:bg-white/15 text-white/80 hover:text-white transition cursor-pointer"
            title="Image suivante (→)"
          >
            <ChevronRight class="w-4 h-4" />
          </button>
        </div>

        <!-- Description texte sous la pill bar si renseignée -->
        <div
          v-if="image.content"
          class="max-w-xl text-center text-xs text-white/75 px-4 line-clamp-2 leading-relaxed"
        >
          {{ image.content }}
        </div>
      </div>
    </div>
  </Teleport>
</template>

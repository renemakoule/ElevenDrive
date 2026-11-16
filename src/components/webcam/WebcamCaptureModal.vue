<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { Camera, X, RefreshCw, Check, Upload, AlertCircle } from 'lucide-vue-next';

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'captured', photoDataUrl: string): void;
}>();

const videoRef = ref<HTMLVideoElement | null>(null);
const canvasRef = ref<HTMLCanvasElement | null>(null);
const stream = ref<MediaStream | null>(null);
const capturedImage = ref<string | null>(null);
const cameraError = ref<string | null>(null);
const isLoading = ref(true);
const flashEffect = ref(false);

const startCamera = async () => {
  isLoading.value = true;
  cameraError.value = null;
  capturedImage.value = null;

  try {
    if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
      throw new Error("L'accès à la webcam n'est pas supporté par ce système.");
    }

    const mediaStream = await navigator.mediaDevices.getUserMedia({
      video: {
        width: { ideal: 1280 },
        height: { ideal: 720 },
        facingMode: 'user',
      },
      audio: false,
    });

    stream.value = mediaStream;
    if (videoRef.value) {
      videoRef.value.srcObject = mediaStream;
      videoRef.value.play();
    }
  } catch (err: any) {
    cameraError.value = err.message || "Impossible d'accéder à la webcam USB.";
  } finally {
    isLoading.value = false;
  }
};

const stopCamera = () => {
  if (stream.value) {
    stream.value.getTracks().forEach((track) => track.stop());
    stream.value = null;
  }
};

const capturePhoto = () => {
  if (!videoRef.value) return;

  // Déclencher l'effet flash visuel
  flashEffect.value = true;
  setTimeout(() => (flashEffect.value = false), 150);

  const video = videoRef.value;
  const canvas = document.createElement('canvas');
  
  // Format standard officiel photo 4x4 (carré haute résolution 600x600)
  const size = 600;
  canvas.width = size;
  canvas.height = size;
  const ctx = canvas.getContext('2d');
  if (!ctx) return;

  // Calcul du recadrage centré carré à partir du flux vidéo 16:9
  const videoW = video.videoWidth || 1280;
  const videoH = video.videoHeight || 720;
  const minDim = Math.min(videoW, videoH);
  const startX = (videoW - minDim) / 2;
  const startY = (videoH - minDim) / 2;

  // Dessin avec effet miroir horizontal naturel pour webcam
  ctx.translate(size, 0);
  ctx.scale(-1, 1);
  ctx.drawImage(video, startX, startY, minDim, minDim, 0, 0, size, size);

  capturedImage.value = canvas.toDataURL('image/jpeg', 0.88);
};

const retakePhoto = () => {
  capturedImage.value = null;
};

const confirmPhoto = () => {
  if (capturedImage.value) {
    emit('captured', capturedImage.value);
    stopCamera();
    emit('close');
  }
};

const handleFileUpload = (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0];
  if (!file) return;

  const reader = new FileReader();
  reader.onload = (e) => {
    capturedImage.value = e.target?.result as string;
  };
  reader.readAsDataURL(file);
};

const handleKeyDown = (e: KeyboardEvent) => {
  if (e.code === 'Space') {
    e.preventDefault();
    if (!capturedImage.value) {
      capturePhoto();
    }
  } else if (e.key === 'Enter') {
    e.preventDefault();
    if (capturedImage.value) {
      confirmPhoto();
    }
  } else if (e.key === 'Escape') {
    e.preventDefault();
    if (capturedImage.value) {
      retakePhoto();
    } else {
      stopCamera();
      emit('close');
    }
  }
};

onMounted(() => {
  startCamera();
  window.addEventListener('keydown', handleKeyDown);
});

onUnmounted(() => {
  stopCamera();
  window.removeEventListener('keydown', handleKeyDown);
});
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 select-none">
    <div class="bg-card border border-border rounded-[5px] max-w-lg w-full overflow-hidden flex flex-col">
      <!-- En-tête Modal Prise de Vue -->
      <header class="p-3 bg-card border-b border-border flex items-center justify-between">
        <div class="flex items-center space-x-2">
          <Camera class="w-4 h-4 text-brand-dark" />
          <span class="text-xs font-bold uppercase tracking-wider text-foreground">
            Prise de Photo Officielle 4x4 (Webcam USB)
          </span>
        </div>
        <button
          @click="stopCamera(); emit('close');"
          class="p-1 rounded-lg text-muted-foreground hover:text-foreground hover:bg-muted transition"
          aria-label="Fermer"
        >
          <X class="w-4 h-4" />
        </button>
      </header>

      <!-- Zone Principale Viseur Webcam / Aperçu Capturé -->
      <div class="p-6 bg-muted/30 flex flex-col items-center justify-center relative">
        <!-- Effet Flash Shutter -->
        <div
          v-if="flashEffect"
          class="absolute inset-0 bg-white z-40 transition-opacity duration-150 pointer-events-none"
        ></div>

        <!-- Mode 1 : Flux Vidéo Webcam en Direct avec Pochoir Cadrage 4x4 -->
        <div
          v-if="!capturedImage"
          class="relative w-72 h-72 rounded-[5px] overflow-hidden border-2 border-brand bg-black shadow-inner flex items-center justify-center"
        >
          <video
            ref="videoRef"
            autoplay
            playsinline
            muted
            class="w-full h-full object-cover transform -scale-x-100"
          ></video>

          <!-- Pochoir Visage & Cadrage Réglementaire 4x4 (Norme Ministère des Transports) -->
          <div class="absolute inset-0 pointer-events-none flex flex-col items-center justify-center">
            <!-- Ovale Tête / Yeux -->
            <div class="w-36 h-48 rounded-[20%] border-2 border-dashed border-brand/70 relative">
              <!-- Repère Ligne des Yeux -->
              <div class="absolute top-[38%] left-2 right-2 border-t border-brand/40"></div>
              <!-- Repère Menton -->
              <div class="absolute bottom-1 left-8 right-8 border-b-2 border-brand/60"></div>
            </div>
            <!-- Repère Épaules -->
            <div class="w-56 h-12 border-t-2 border-dashed border-brand/40 mt-1 rounded-t-full"></div>
            <!-- Badge 4x4 officiel -->
            <span class="absolute top-2 right-2 px-1.5 py-0.5 rounded bg-black/70 text-[10px] font-mono text-brand font-bold">
              Format 4x4
            </span>
          </div>

          <!-- Message d'erreur si la webcam n'est pas connectée -->
          <div
            v-if="cameraError"
            class="absolute inset-0 bg-card/95 flex flex-col items-center justify-center p-4 text-center"
          >
            <AlertCircle class="w-10 h-10 text-amber-500 mb-2" />
            <p class="text-xs text-foreground font-medium">{{ cameraError }}</p>
            <p class="text-[11px] text-muted-foreground mt-1">Vous pouvez importer un fichier image depuis votre ordinateur :</p>
            <label class="mt-3 px-3 py-1.5 rounded-lg bg-card hover:bg-muted text-xs text-foreground border border-border cursor-pointer flex items-center gap-1.5 transition">
              <Upload class="w-3.5 h-3.5 text-brand-dark" />
              <span>Choisir un fichier</span>
              <input type="file" accept="image/*" class="hidden" @change="handleFileUpload" />
            </label>
          </div>
        </div>

        <!-- Mode 2 : Image Capturée & Prête pour Validation -->
        <div
          v-else
          class="relative w-72 h-72 rounded-[5px] overflow-hidden border-4 border-brand bg-black"
        >
          <img :src="capturedImage" alt="Photo capturée" class="w-full h-full object-cover" />
          <div class="absolute bottom-2 left-2 right-2 bg-black/70 backdrop-blur-sm py-1 px-2 rounded text-center text-[10px] text-brand font-bold">
            ✓ Photo 4x4 Recadrée & Prête
          </div>
        </div>

        <!-- Consigne clavier claire sous le viseur -->
        <div class="mt-4 text-center">
          <p v-if="!capturedImage" class="text-xs text-muted-foreground font-medium">
            Placez le visage au centre du pochoir puis appuyez sur <kbd class="rounded bg-muted px-1.5 py-0.5 font-mono text-brand-dark border border-brand/40">Espace</kbd>
          </p>
          <p v-else class="text-xs text-muted-foreground font-medium">
            Appuyez sur <kbd class="rounded bg-muted px-1.5 py-0.5 font-mono text-brand-dark border border-brand/40">Entrée</kbd> pour valider ou <kbd class="rounded bg-muted px-1.5 py-0.5 font-mono text-muted-foreground">Échap</kbd> pour reprendre
          </p>
        </div>
      </div>

      <!-- Actions Pied de Modal -->
      <footer class="p-4 bg-card border-t border-border flex items-center justify-between">
        <!-- Option d'import alternatif de photo -->
        <label class="cursor-pointer text-xs text-muted-foreground hover:text-foreground flex items-center gap-1.5 transition">
          <Upload class="w-3.5 h-3.5 text-muted-foreground" />
          <span>Importer fichier</span>
          <input type="file" accept="image/*" class="hidden" @change="handleFileUpload" />
        </label>

        <!-- Boutons d'action -->
        <div class="flex items-center space-x-2">
          <template v-if="!capturedImage">
            <button
              type="button"
              @click="stopCamera(); emit('close');"
              class="px-3 py-1.5 rounded-lg border border-border bg-card hover:bg-muted text-foreground text-xs font-medium transition"
            >
              Annuler
            </button>
            <button
              type="button"
              @click="capturePhoto"
              class="px-4 py-2 rounded-lg bg-brand-dark hover:bg-brand text-white font-semibold text-xs flex items-center gap-2 shadow-sm transition"
            >
              <Camera class="w-4 h-4" />
              <span>Prendre la Photo</span>
              <kbd class="rounded bg-white/20 px-1 font-mono text-[10px] text-white/90">Espace</kbd>
            </button>
          </template>

          <template v-else>
            <button
              type="button"
              @click="retakePhoto"
              class="px-3 py-1.5 rounded-lg border border-border bg-card hover:bg-muted text-foreground text-xs font-medium flex items-center gap-1.5 transition"
            >
              <RefreshCw class="w-3.5 h-3.5 text-amber-500" />
              <span>Reprendre</span>
              <kbd class="rounded bg-muted px-1 font-mono text-[10px] text-muted-foreground">Échap</kbd>
            </button>
            <button
              type="button"
              @click="confirmPhoto"
              class="px-4 py-2 rounded-lg bg-brand-dark hover:bg-brand text-white font-semibold text-xs flex items-center gap-2 shadow-sm transition"
            >
              <Check class="w-4 h-4" />
              <span>Valider la Photo</span>
              <kbd class="rounded bg-white/20 px-1 font-mono text-[10px] text-white/90">Entrée</kbd>
            </button>
          </template>
        </div>
      </footer>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue';
import { useLanguageStore } from '../../stores/languageStore';
import { speechAudioService } from '../../services/speechAudioService';
import { PdfTextExtractor } from '../../services/pdfTextExtractor';
import { deepLTranslationService } from '../../services/deepLTranslationService';
import AudioSpeechBar from '../common/AudioSpeechBar.vue';
import {
  FileText,
  Languages,
  Volume2,
  Download,
  X,
  Sparkles,
  BookOpen,
  CheckCircle2,
  Play,
  RotateCw
} from 'lucide-vue-next';

const props = withDefaults(
  defineProps<{
    isOpen: boolean;
    titre: string;
    categoriePermis?: string;
    fichierData?: string;
    description?: string;
    allowDownload?: boolean;
  }>(),
  {
    isOpen: false,
    categoriePermis: 'B',
    fichierData: '',
    description: '',
    allowDownload: true,
  }
);

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'download'): void;
}>();

const languageStore = useLanguageStore();

// 'pdf' = Vue visionneuse PDF brute | 'text' = Vue texte bilingue transcrit & traduit
const viewMode = ref<'pdf' | 'text'>('pdf');

// Extraction & Traduction
const isExtracting = ref(false);
const isTranslating = ref(false);
const rawExtractedText = ref('');
const extractedPages = ref<string[]>([]);
const originalSegments = ref<string[]>([]);
const translatedSegments = ref<string[]>([]);
const translatedTitle = ref('');

// Barre audio
const isAudioBarOpen = ref(false);

// Cache mémoire local par langue pour éviter de réémettre des requêtes
const translationCache = new Map<string, { title: string; segments: string[] }>();

const isCurrentLangRTL = computed(() => languageStore.isRTL);
const currentLangMeta = computed(() => languageStore.currentMeta);

const activeDisplayTitle = computed(() => {
  if (languageStore.currentLanguage === 'fr') {
    return props.titre;
  }
  return translatedTitle.value || props.titre;
});

const currentPlayingSegmentIndex = computed(() => speechAudioService.currentSegmentIndex.value);
const isSpeechPlaying = computed(() => speechAudioService.isPlaying.value && !speechAudioService.isPaused.value);

// Extraction initiale du PDF
const extraireContenuPdf = async () => {
  if (!props.fichierData) return;

  isExtracting.value = true;
  try {
    const res = await PdfTextExtractor.extraireTexte(props.fichierData, {
      titre: props.titre,
      description: props.description,
      categorie: props.categoriePermis,
    });

    rawExtractedText.value = res.texte;
    extractedPages.value = res.pages || [res.texte];

    // Découpage en segments audibles
    const segments = speechAudioService.decouperTexteEnSegments(res.texte);
    originalSegments.value = segments.length > 0 ? segments : [res.texte];

    // Si la langue actuelle n'est pas le français, traduire immédiatement
    await appliquerTraduction(languageStore.currentLanguage);
  } catch (err) {
    console.error('Erreur lors de l extraction du PDF :', err);
    originalSegments.value = [props.titre, props.description || ''].filter(Boolean);
    translatedSegments.value = [...originalSegments.value];
  } finally {
    isExtracting.value = false;
  }
};

// Application de la traduction sur les segments et le titre
const appliquerTraduction = async (lang: string) => {
  if (lang === 'fr') {
    translatedSegments.value = [...originalSegments.value];
    translatedTitle.value = props.titre;
    return;
  }

  // Vérifier le cache
  const cached = translationCache.get(lang);
  if (cached) {
    translatedSegments.value = cached.segments;
    translatedTitle.value = cached.title;
    return;
  }

  if (originalSegments.value.length === 0) return;

  isTranslating.value = true;
  try {
    const targetLang = lang as 'en' | 'es' | 'ar';

    // 1. Traduire le titre du support
    let tradTitle = props.titre;
    try {
      tradTitle = await deepLTranslationService.translateAsync(props.titre, targetLang);
    } catch {
      tradTitle = props.titre;
    }
    translatedTitle.value = tradTitle;

    // 2. Traduire l'ensemble ordonné des segments du PDF
    const tradSegs = await deepLTranslationService.traduireSegments(
      originalSegments.value,
      targetLang
    );
    translatedSegments.value = tradSegs;

    // Mémoriser dans le cache de session
    translationCache.set(lang, {
      title: tradTitle,
      segments: tradSegs,
    });
  } catch (err) {
    console.warn('Erreur lors de la traduction du PDF :', err);
    translatedSegments.value = [...originalSegments.value];
    translatedTitle.value = props.titre;
  } finally {
    isTranslating.value = false;
  }
};

// Démarrer la lecture audio
const lancerLectureAudio = async (indexDepart = 0) => {
  isAudioBarOpen.value = true;

  // S'assurer que le contenu a été extrait
  if (originalSegments.value.length === 0) {
    await extraireContenuPdf();
  }

  const currentLang = languageStore.currentLanguage;
  if (currentLang !== 'fr' && translatedSegments.value.length === 0) {
    await appliquerTraduction(currentLang);
  }

  const segmentsALire = currentLang === 'fr' ? originalSegments.value : translatedSegments.value;

  speechAudioService.lireSegments(segmentsALire, indexDepart, {
    lang: currentLang,
  });
};

const allerAuSegment = (index: number) => {
  if (!isAudioBarOpen.value || !speechAudioService.isPlaying.value) {
    lancerLectureAudio(index);
  } else {
    speechAudioService.allerAuSegment(index);
  }
};

const fermerAudio = () => {
  speechAudioService.arreter();
  isAudioBarOpen.value = false;
};

const handleClose = () => {
  speechAudioService.arreter();
  isAudioBarOpen.value = false;
  emit('close');
};

const forcerRetraduction = async () => {
  translationCache.delete(languageStore.currentLanguage);
  await appliquerTraduction(languageStore.currentLanguage);
};

// Réagir aux changements de langue de l'application
watch(
  () => languageStore.currentLanguage,
  async (newLang) => {
    if (props.isOpen) {
      await appliquerTraduction(newLang);
      // Si la lecture audio est en cours, relancer avec la nouvelle langue et voix
      if (speechAudioService.isPlaying.value) {
        lancerLectureAudio(speechAudioService.currentSegmentIndex.value);
      }
    }
  }
);

// Réagir à l'ouverture de la modale
watch(
  () => props.isOpen,
  (open) => {
    if (open) {
      // Par défaut, si la langue n'est pas le français, basculer sur la vue texte traduit
      if (languageStore.currentLanguage !== 'fr') {
        viewMode.value = 'text';
      } else {
        viewMode.value = 'pdf';
      }
      extraireContenuPdf();
    } else {
      speechAudioService.arreter();
      isAudioBarOpen.value = false;
    }
  },
  { immediate: true }
);

onUnmounted(() => {
  speechAudioService.arreter();
});
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-50 flex items-center justify-center p-2 sm:p-4 bg-black/80 backdrop-blur-sm animate-in fade-in"
    role="dialog"
    aria-modal="true"
  >
    <div
      class="bg-card border border-border w-full max-w-5xl h-[92vh] rounded-[8px] shadow-2xl flex flex-col overflow-hidden text-foreground"
    >
      <!-- En-tête de la visionneuse avec onglets bilingues -->
      <header class="px-3 sm:px-5 py-2.5 border-b border-border flex flex-wrap items-center justify-between gap-2 bg-muted/40 shrink-0">
        <!-- Titre et catégorie -->
        <div class="flex items-center space-x-2 min-w-0 max-w-full sm:max-w-md">
          <span
            v-if="categoriePermis"
            class="text-[10px] font-bold px-2 py-0.5 rounded-[4px] bg-brand/10 text-brand border border-brand/20 uppercase shrink-0"
          >
            Permis {{ categoriePermis }}
          </span>
          <div class="min-w-0">
            <h2 class="text-xs sm:text-sm font-bold text-foreground truncate" :title="activeDisplayTitle">
              {{ activeDisplayTitle }}
            </h2>
            <p v-if="languageStore.currentLanguage !== 'fr' && activeDisplayTitle !== titre" class="text-[10px] text-muted-foreground truncate">
              Original : {{ titre }}
            </p>
          </div>
        </div>

        <!-- Sélecteur de vue : PDF brut vs Texte Traduit -->
        <div class="flex items-center rounded-lg p-0.5 bg-muted border border-border text-xs">
          <button
            type="button"
            @click="viewMode = 'pdf'"
            class="flex items-center space-x-1.5 px-2.5 py-1 rounded-[5px] font-semibold transition cursor-pointer"
            :class="
              viewMode === 'pdf'
                ? 'bg-card text-foreground shadow-xs font-bold'
                : 'text-muted-foreground hover:text-foreground'
            "
          >
            <FileText class="w-3.5 h-3.5" />
            <span class="hidden sm:inline">Document PDF</span>
          </button>

          <button
            type="button"
            @click="viewMode = 'text'"
            class="flex items-center space-x-1.5 px-2.5 py-1 rounded-[5px] font-semibold transition cursor-pointer"
            :class="
              viewMode === 'text'
                ? 'bg-card text-brand shadow-xs font-bold'
                : 'text-muted-foreground hover:text-foreground'
            "
          >
            <Sparkles class="w-3.5 h-3.5 text-brand" />
            <span>
              {{ languageStore.currentLanguage === 'fr' ? 'Transcription & Audio' : `Traduction (${currentLangMeta.flag} ${currentLangMeta.code.toUpperCase()})` }}
            </span>
          </button>
        </div>

        <!-- Actions : Écouter, Télécharger, Fermer -->
        <div class="flex items-center space-x-1.5 shrink-0">
          <button
            type="button"
            @click="isAudioBarOpen ? fermerAudio() : lancerLectureAudio(0)"
            class="flex items-center space-x-1.5 px-2.5 py-1.5 rounded-[5px] border text-xs font-semibold transition cursor-pointer"
            :class="
              isSpeechPlaying
                ? 'bg-brand text-white border-brand shadow-sm'
                : isAudioBarOpen
                ? 'bg-brand/10 text-brand border-brand/30'
                : 'bg-muted hover:bg-muted/80 text-foreground border-border'
            "
            :title="isAudioBarOpen ? 'Fermer la barre audio' : 'Écouter la lecture vocale complète du document'"
          >
            <Volume2 class="w-3.5 h-3.5" :class="isSpeechPlaying ? 'animate-pulse' : ''" />
            <span class="hidden sm:inline">{{ isSpeechPlaying ? 'En lecture...' : 'Écouter' }}</span>
          </button>

          <button
            v-if="allowDownload && fichierData"
            type="button"
            @click="emit('download')"
            class="flex items-center space-x-1.5 px-2.5 py-1.5 rounded-[5px] bg-muted hover:bg-muted/80 text-foreground border border-border text-xs font-semibold transition cursor-pointer"
            title="Télécharger le fichier PDF original"
          >
            <Download class="w-3.5 h-3.5" />
            <span class="hidden md:inline">Télécharger</span>
          </button>

          <button
            type="button"
            @click="handleClose"
            class="p-1.5 rounded-[5px] text-muted-foreground hover:text-foreground hover:bg-muted transition cursor-pointer"
            title="Fermer"
          >
            <X class="w-4 h-4" />
          </button>
        </div>
      </header>

      <!-- Corps du document -->
      <div class="flex-1 bg-zinc-950 relative overflow-hidden flex flex-col min-h-0">
        <!-- VUE 1 : IFRAME PDF BRUT -->
        <div v-show="viewMode === 'pdf'" class="w-full h-full relative">
          <iframe
            v-if="fichierData"
            :src="fichierData"
            class="w-full h-full border-0 bg-zinc-900"
            title="Visionneuse PDF"
          />
          <div v-else class="flex flex-col items-center justify-center h-full text-zinc-400 p-6 text-center">
            <FileText class="w-12 h-12 mb-3 text-zinc-600" />
            <p class="text-sm font-medium">Aucun fichier PDF disponible pour ce cours.</p>
          </div>
        </div>

        <!-- VUE 2 : TEXTE BILINGUE TRANSHRIT & TRADUIT -->
        <div
          v-show="viewMode === 'text'"
          class="w-full h-full overflow-y-auto p-4 sm:p-8 bg-card text-foreground"
          :dir="isCurrentLangRTL ? 'rtl' : 'ltr'"
        >
          <!-- Chargement de l'extraction ou de la traduction -->
          <div
            v-if="isExtracting || isTranslating"
            class="flex flex-col items-center justify-center py-16 space-y-4"
          >
            <div class="w-10 h-10 rounded-full border-2 border-brand border-t-transparent animate-spin" />
            <p class="text-sm font-semibold text-foreground">
              {{ isExtracting ? 'Extraction du texte haute précision depuis le PDF...' : `Traduction du document en ${currentLangMeta.name} via DeepL...` }}
            </p>
            <p class="text-xs text-muted-foreground max-w-sm text-center">
              Veuillez patienter quelques instants pendant le traitement automatique.
            </p>
          </div>

          <!-- Affichage du document textuel bilingue avec karaoké audio interactif -->
          <div v-else class="max-w-3xl mx-auto space-y-6 pb-20">
            <!-- Bandeau de statut de traduction -->
            <div class="flex items-center justify-between p-3 rounded-lg bg-muted/40 border border-border text-xs">
              <div class="flex items-center space-x-2">
                <span class="text-base">{{ currentLangMeta.flag }}</span>
                <span class="font-bold text-foreground">
                  Langue affichée : {{ currentLangMeta.name }}
                </span>
                <span v-if="languageStore.currentLanguage !== 'fr'" class="text-[10px] px-1.5 py-0.5 rounded bg-brand/10 text-brand font-medium">
                  Traduction active
                </span>
              </div>

              <div class="flex items-center space-x-2">
                <span class="text-[11px] text-muted-foreground">
                  {{ translatedSegments.length }} phrases
                </span>
                <button
                  v-if="languageStore.currentLanguage !== 'fr'"
                  type="button"
                  @click="forcerRetraduction"
                  class="flex items-center space-x-1 text-[11px] text-muted-foreground hover:text-brand transition cursor-pointer"
                  title="Actualiser la traduction"
                >
                  <RotateCw class="w-3 h-3" />
                  <span>Rafraîchir</span>
                </button>
              </div>
            </div>

            <!-- Titre du document -->
            <div class="border-b border-border pb-4 space-y-1">
              <h1 class="text-lg sm:text-xl font-extrabold text-foreground tracking-tight">
                {{ activeDisplayTitle }}
              </h1>
              <p v-if="description" class="text-xs text-muted-foreground italic">
                {{ description }}
              </p>
            </div>

            <!-- Liste des paragraphes / segments textuels avec interactivité audio karaoké -->
            <div class="space-y-2 text-sm leading-relaxed">
              <p class="text-[11px] text-muted-foreground uppercase font-semibold tracking-wider">
                Cliquez sur une phrase pour démarrer la lecture vocale :
              </p>

              <div class="space-y-1.5">
                <div
                  v-for="(seg, idx) in translatedSegments"
                  :key="idx"
                  @click="allerAuSegment(idx)"
                  class="p-2 sm:p-2.5 rounded-[5px] transition-all duration-150 cursor-pointer text-foreground group flex items-start gap-2"
                  :class="
                    isAudioBarOpen && currentPlayingSegmentIndex === idx
                      ? 'bg-brand/10 border-l-4 border-brand text-brand font-semibold shadow-xs'
                      : 'hover:bg-muted/60 border-l-2 border-transparent'
                  "
                >
                  <span
                    class="text-[10px] font-mono px-1.5 py-0.5 rounded bg-muted/60 text-muted-foreground shrink-0 select-none group-hover:bg-brand/10 group-hover:text-brand"
                  >
                    {{ idx + 1 }}
                  </span>
                  <span class="flex-1">
                    {{ seg }}
                  </span>
                  <Volume2
                    v-if="isAudioBarOpen && currentPlayingSegmentIndex === idx && isSpeechPlaying"
                    class="w-4 h-4 text-brand animate-pulse shrink-0 mt-0.5"
                  />
                  <Play
                    v-else
                    class="w-3.5 h-3.5 text-muted-foreground opacity-0 group-hover:opacity-100 transition shrink-0 mt-0.5"
                  />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Barre audio vocale ancrée en bas de la visionneuse -->
      <AudioSpeechBar
        v-if="isAudioBarOpen"
        :title="activeDisplayTitle"
        :is-extracting="isExtracting || isTranslating"
        @close="fermerAudio"
      />
    </div>
  </div>
</template>

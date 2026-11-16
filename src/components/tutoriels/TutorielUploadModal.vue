<script setup lang="ts">
import { ref, computed } from 'vue';
import { useTutorielsStore } from '../../stores/tutorielsStore';
import { X, Play, Video, Link, Sparkles, AlertCircle } from 'lucide-vue-next';

defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
}>();

const tutorielsStore = useTutorielsStore();

const formTitre = ref('');
const formCategorie = ref('B');
const formUrlVideo = ref('');
const formDuree = ref(8);
const formDescription = ref('');
const formMiniature = ref('');
const isSubmitting = ref(false);

const CATEGORIES_OPTIONS = [
  { value: 'TOUTES', label: 'Toutes Catégories (Général / Code)' },
  { value: 'B', label: 'Permis B (Voiture Légère)' },
  { value: 'A', label: 'Permis A / A1 (Moto & Deux-roues)' },
  { value: 'C', label: 'Permis C (Poids Lourd)' },
  { value: 'D', label: 'Permis D (Transport en Commun)' },
  { value: 'E', label: 'Permis E (Remorques & Attelages)' },
];

/**
 * Auto-génération de la miniature si lien YouTube
 */
const autoThumbnail = computed(() => {
  if (formMiniature.value.trim()) return formMiniature.value;
  const url = formUrlVideo.value.trim();
  if (url.includes('youtu.be/')) {
    const id = url.split('youtu.be/')[1].split('?')[0];
    return `https://img.youtube.com/vi/${id}/hqdefault.jpg`;
  }
  if (url.includes('watch?v=')) {
    const id = url.split('watch?v=')[1].split('&')[0];
    return `https://img.youtube.com/vi/${id}/hqdefault.jpg`;
  }
  return '';
});

const handleSubmit = async () => {
  if (!formTitre.value.trim() || !formUrlVideo.value.trim()) return;

  isSubmitting.value = true;
  try {
    await tutorielsStore.ajouterTutoriel({
      titre: formTitre.value.trim(),
      categorie_permis: formCategorie.value,
      url_video: formUrlVideo.value.trim(),
      duree_minutes: formDuree.value || 5,
      description: formDescription.value.trim() || undefined,
      miniature_url: autoThumbnail.value || undefined,
    });
    // Réinitialiser formulaire
    formTitre.value = '';
    formUrlVideo.value = '';
    formDescription.value = '';
    emit('close');
  } catch {
    // Géré par le store
  } finally {
    isSubmitting.value = false;
  }
};
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-50 flex items-center justify-center p-3 sm:p-4 bg-black/70 backdrop-blur-xs animate-in fade-in"
    role="dialog"
    aria-modal="true"
  >
    <div class="bg-card border border-border w-full max-w-xl max-h-[92vh] rounded-[5px] shadow-2xl flex flex-col overflow-hidden text-foreground">
      <!-- En-tête -->
      <header class="px-5 py-3.5 border-b border-border flex items-center justify-between bg-muted/20">
        <div class="flex items-center space-x-2">
          <div class="w-7 h-7 rounded-[4px] bg-red-500/10 border border-red-500/20 text-red-600 flex items-center justify-center">
            <Play class="w-3.5 h-3.5 fill-current" />
          </div>
          <div>
            <h2 class="text-sm font-bold text-foreground">Ajouter un tutoriel vidéo pédagogique</h2>
            <p class="text-[11px] text-muted-foreground">Publication d'un guide vidéo pour les apprenants</p>
          </div>
        </div>
        <button
          type="button"
          @click="$emit('close')"
          class="text-muted-foreground hover:text-foreground cursor-pointer"
        >
          <X class="w-4 h-4" />
        </button>
      </header>

      <!-- Formulaire -->
      <form @submit.prevent="handleSubmit" class="p-5 overflow-y-auto space-y-4 text-xs">
        <!-- Titre de la vidéo -->
        <div>
          <label class="block font-semibold text-foreground mb-1">Titre du tutoriel vidéo *</label>
          <input
            v-model="formTitre"
            type="text"
            required
            placeholder="Ex : Réussir son créneau à droite étape par étape"
            class="w-full h-8 px-2.5 rounded-[4px] border border-border bg-background text-foreground text-xs focus:outline-none focus:border-brand-dark"
          />
        </div>

        <!-- Catégorie & Durée -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div>
            <label class="block font-semibold text-foreground mb-1">Catégorie de Permis *</label>
            <select
              v-model="formCategorie"
              class="w-full h-8 px-2 rounded-[4px] border border-border bg-background text-foreground text-xs focus:outline-none focus:border-brand-dark font-medium"
            >
              <option v-for="c in CATEGORIES_OPTIONS" :key="c.value" :value="c.value">
                {{ c.label }}
              </option>
            </select>
          </div>

          <div>
            <label class="block font-semibold text-foreground mb-1">Durée approximative (minutes) *</label>
            <input
              v-model.number="formDuree"
              type="number"
              min="1"
              max="120"
              required
              class="w-full h-8 px-2 rounded-[4px] border border-border bg-background text-foreground text-xs focus:outline-none focus:border-brand-dark"
            />
          </div>
        </div>

        <!-- URL de la vidéo (YouTube, Vimeo, MP4 direct) -->
        <div>
          <label class="block font-semibold text-foreground mb-1">Lien de la vidéo (YouTube, Vimeo ou MP4) *</label>
          <div class="relative">
            <input
              v-model="formUrlVideo"
              type="url"
              required
              placeholder="https://www.youtube.com/watch?v=... ou https://youtu.be/..."
              class="w-full h-8 pl-8 pr-2.5 rounded-[4px] border border-border bg-background text-foreground text-xs focus:outline-none focus:border-brand-dark font-mono text-[11px]"
            />
            <Link class="w-3.5 h-3.5 text-muted-foreground absolute left-2.5 top-2.5" />
          </div>
          <p class="text-[10px] text-muted-foreground mt-1">
            Compatible avec YouTube, Vimeo ou les liens directs de fichiers vidéo (.mp4).
          </p>
        </div>

        <!-- Aperçu de miniature si détectée -->
        <div v-if="autoThumbnail" class="p-2.5 rounded-[4px] border border-border bg-muted/20 flex items-center gap-3">
          <img :src="autoThumbnail" alt="Aperçu" class="w-20 h-12 object-cover rounded border border-border" />
          <div>
            <span class="text-[11px] font-semibold text-foreground flex items-center gap-1">
              <Sparkles class="w-3.5 h-3.5 text-brand" /> Miniature détectée
            </span>
            <p class="text-[10px] text-muted-foreground">La couverture sera automatiquement affichée sur la fiche vidéo.</p>
          </div>
        </div>

        <!-- Description & repères pédagogiques -->
        <div>
          <label class="block font-semibold text-foreground mb-1">Description & Conseils clés (Facultatif)</label>
          <textarea
            v-model="formDescription"
            rows="3"
            placeholder="Points de repère, règles à respecter, conseils du moniteur..."
            class="w-full p-2 rounded-[4px] border border-border bg-background text-foreground text-xs focus:outline-none focus:border-brand-dark resize-none"
          />
        </div>

        <!-- Boutons d'action -->
        <div class="pt-3 border-t border-border flex items-center justify-end space-x-2">
          <button
            type="button"
            @click="$emit('close')"
            class="px-3 py-1.5 rounded-[4px] border border-border hover:bg-muted text-xs font-semibold text-foreground transition cursor-pointer"
          >
            Annuler
          </button>

          <button
            type="submit"
            :disabled="isSubmitting"
            class="flex items-center space-x-1.5 px-4 py-1.5 rounded-[4px] bg-brand hover:bg-brand-dark text-brand-dark font-bold text-xs uppercase tracking-wider transition cursor-pointer disabled:opacity-50"
          >
            <Play class="w-3.5 h-3.5 fill-current" />
            <span>{{ isSubmitting ? 'Enregistrement...' : 'Publier le tutoriel' }}</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

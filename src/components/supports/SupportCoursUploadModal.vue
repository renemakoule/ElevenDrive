<script setup lang="ts">
import { ref } from 'vue';
import { useSupportsCoursStore } from '../../stores/supportsCoursStore';
import { useGuichetStore } from '../../stores/guichetStore';
import { formatExternalUrl } from '../../utils/urlHelper';
import { X, FileText, UploadCloud, Link, FileCheck, AlertCircle } from 'lucide-vue-next';

defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
}>();

const supportsStore = useSupportsCoursStore();
const guichetStore = useGuichetStore();

const formTitre = ref('');
const formCategorie = ref('B');
const formTypeSupport = ref<'PDF' | 'LIEN_WEB'>('PDF');
const formDescription = ref('');
const formUrlExterne = ref('');

// Fichier PDF
const fichierNom = ref('');
const fichierTaille = ref('');
const fichierBase64 = ref('');
const isSubmitting = ref(false);

const CATEGORIES_OPTIONS = [
  { value: 'TOUTES', label: 'Toutes Catégories (Tronc Commun / Signalisation)' },
  { value: 'B', label: 'Permis B (Voiture Légère)' },
  { value: 'A', label: 'Permis A / A1 (Moto & Deux-roues)' },
  { value: 'C', label: 'Permis C (Poids Lourd)' },
  { value: 'D', label: 'Permis D (Transport en Commun)' },
  { value: 'E', label: 'Permis E (Remorques & Attelages)' },
];

const handleFileUpload = (e: Event) => {
  const target = e.target as HTMLInputElement;
  const file = target.files?.[0];
  if (!file) return;

  if (file.type !== 'application/pdf' && !file.name.toLowerCase().endsWith('.pdf')) {
    guichetStore.showToast('Veuillez sélectionner un fichier PDF valide.', 'error');
    target.value = '';
    return;
  }

  // Limite : 15 Mo
  if (file.size > 15 * 1024 * 1024) {
    guichetStore.showToast('Le fichier PDF dépasse la taille maximale autorisée (15 Mo).', 'error');
    target.value = '';
    return;
  }

  fichierNom.value = file.name;
  fichierTaille.value = (file.size / (1024 * 1024)).toFixed(2) + ' Mo';

  const reader = new FileReader();
  reader.onload = (uploadEvent) => {
    fichierBase64.value = uploadEvent.target?.result as string;
  };
  reader.readAsDataURL(file);
};

const handleSubmit = async () => {
  if (!formTitre.value.trim()) return;

  if (formTypeSupport.value === 'PDF' && !fichierBase64.value) {
    guichetStore.showToast('Veuillez déposer un document PDF pour ce support.', 'error');
    return;
  }

  if (formTypeSupport.value === 'LIEN_WEB' && !formUrlExterne.value.trim()) {
    guichetStore.showToast("Veuillez renseigner l'adresse web du cours.", 'error');
    return;
  }

  isSubmitting.value = true;
  try {
    await supportsStore.ajouterSupport({
      titre: formTitre.value.trim(),
      categorie_permis: formCategorie.value,
      description: formDescription.value.trim() || undefined,
      type_support: formTypeSupport.value,
      fichier_nom: formTypeSupport.value === 'PDF' ? fichierNom.value : undefined,
      fichier_taille: formTypeSupport.value === 'PDF' ? fichierTaille.value : undefined,
      fichier_data: formTypeSupport.value === 'PDF' ? fichierBase64.value : undefined,
      url_externe: formTypeSupport.value === 'LIEN_WEB' ? formatExternalUrl(formUrlExterne.value) : undefined,
    });

    // Réinitialisation
    formTitre.value = '';
    formDescription.value = '';
    formUrlExterne.value = '';
    fichierNom.value = '';
    fichierBase64.value = '';
    emit('close');
  } catch {
    // géré par le store
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
          <div class="w-7 h-7 rounded-[4px] bg-brand-light border border-brand-border text-brand-dark flex items-center justify-center">
            <FileText class="w-4 h-4 text-brand" />
          </div>
          <div>
            <h2 class="text-sm font-bold text-foreground">Ajouter un support de cours</h2>
            <p class="text-[11px] text-muted-foreground">Publication d'un manuel, livret PDF ou lien externe</p>
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
        <!-- Type de support : PDF vs Lien Web -->
        <div>
          <label class="block font-semibold text-foreground mb-1.5">Format du support *</label>
          <div class="grid grid-cols-2 gap-2">
            <button
              type="button"
              @click="formTypeSupport = 'PDF'"
              class="flex items-center justify-center gap-2 p-2.5 rounded-[4px] border text-xs font-semibold transition cursor-pointer"
              :class="formTypeSupport === 'PDF'
                ? 'border-brand bg-brand-light text-brand-dark shadow-xs'
                : 'border-border bg-card text-muted-foreground hover:bg-muted'"
            >
              <FileText class="w-4 h-4 text-brand" />
              <span>📄 Document PDF à télécharger</span>
            </button>

            <button
              type="button"
              @click="formTypeSupport = 'LIEN_WEB'"
              class="flex items-center justify-center gap-2 p-2.5 rounded-[4px] border text-xs font-semibold transition cursor-pointer"
              :class="formTypeSupport === 'LIEN_WEB'
                ? 'border-brand bg-brand-light text-brand-dark shadow-xs'
                : 'border-border bg-card text-muted-foreground hover:bg-muted'"
            >
              <Link class="w-4 h-4 text-brand" />
              <span>🔗 Lien Web externe (Drive, Web)</span>
            </button>
          </div>
        </div>

        <!-- Titre du cours -->
        <div>
          <label class="block font-semibold text-foreground mb-1">Titre du cours *</label>
          <input
            v-model="formTitre"
            type="text"
            required
            placeholder="Ex : Guide complet des Panneaux et Signaux Routiers"
            class="w-full h-8 px-2.5 rounded-[4px] border border-border bg-background text-foreground text-xs focus:outline-none focus:border-brand-dark"
          />
        </div>

        <!-- Catégorie de Permis -->
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

        <!-- Si mode PDF : Déposer le fichier -->
        <div v-if="formTypeSupport === 'PDF'" class="space-y-1">
          <label class="block font-semibold text-foreground mb-1">Fichier PDF du Cours *</label>
          <label class="border-2 border-dashed border-border hover:border-brand rounded-[5px] p-4 flex flex-col items-center justify-center gap-1.5 cursor-pointer bg-muted/20 hover:bg-brand-light/20 transition">
            <UploadCloud class="w-6 h-6 text-brand" />
            <span class="font-medium text-foreground text-xs">
              {{ fichierNom ? fichierNom : 'Cliquez pour sélectionner le fichier PDF' }}
            </span>
            <span class="text-[10px] text-muted-foreground">
              {{ fichierTaille ? `Taille : ${fichierTaille}` : 'Format PDF uniquement (Max 15 Mo)' }}
            </span>
            <input
              type="file"
              accept="application/pdf,.pdf"
              class="hidden"
              @change="handleFileUpload"
            />
          </label>
        </div>

        <!-- Si mode Lien Web : Saisir l'URL -->
        <div v-else class="space-y-1">
          <label class="block font-semibold text-foreground mb-1">Lien Web du support de cours *</label>
          <div class="relative">
            <input
              v-model="formUrlExterne"
              type="url"
              required
              placeholder="https://drive.google.com/... ou https://mon-cours.com/..."
              class="w-full h-8 pl-8 pr-2.5 rounded-[4px] border border-border bg-background text-foreground text-xs focus:outline-none focus:border-brand-dark font-mono text-[11px]"
            />
            <Link class="w-3.5 h-3.5 text-muted-foreground absolute left-2.5 top-2.5" />
          </div>
          <p class="text-[10px] text-muted-foreground">
            L'élève sera redirigé vers ce lien sécurisé pour consulter le support.
          </p>
        </div>

        <!-- Description -->
        <div>
          <label class="block font-semibold text-foreground mb-1">Description / Sommaire (Facultatif)</label>
          <textarea
            v-model="formDescription"
            rows="3"
            placeholder="Chapitres abordés, notions clés, recommandations d'étude..."
            class="w-full p-2 rounded-[4px] border border-border bg-background text-foreground text-xs focus:outline-none focus:border-brand-dark resize-none"
          />
        </div>

        <!-- Actions -->
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
            class="flex items-center space-x-1.5 px-4 py-1.5 rounded-[4px] bg-brand hover:bg-brand-dark text-primary font-bold text-xs uppercase tracking-wider transition cursor-pointer disabled:opacity-50"
          >
            <FileText class="w-3.5 h-3.5" />
            <span>{{ isSubmitting ? 'Téléversement...' : 'Publier le support' }}</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

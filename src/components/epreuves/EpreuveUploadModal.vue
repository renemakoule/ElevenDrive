<script setup lang="ts">
import { ref } from 'vue';
import { useEpreuvesStore } from '../../stores/epreuvesStore';
import { useAuthStore } from '../../stores/authStore';
import { useGuichetStore } from '../../stores/guichetStore';
import { CATEGORIES_PERMIS } from './epreuvesConstants';
import { BookOpen, X, UploadCloud } from 'lucide-vue-next';

defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'success'): void;
}>();

const epreuvesStore = useEpreuvesStore();
const authStore = useAuthStore();
const guichetStore = useGuichetStore();

const isSubmitting = ref(false);
const formTitre = ref('');
const formCategorie = ref('B');
const formDescription = ref('');
const formAnnee = ref(new Date().getFullYear());
const formSession = ref('Session Normale');
const formDuree = ref(45);
const formQuestions = ref(40);
const formAvecCorrige = ref(true);

const uploadedFileName = ref('');
const uploadedFileSize = ref('');
const uploadedFileData = ref('');
const fileInputRef = ref<HTMLInputElement | null>(null);

const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  if (!file) return;

  if (file.type !== 'application/pdf' && !file.name.toLowerCase().endsWith('.pdf')) {
    guichetStore.showToast('Veuillez sélectionner un document au format PDF (.pdf).', 'error');
    return;
  }

  uploadedFileName.value = file.name;
  const sizeMb = (file.size / (1024 * 1024)).toFixed(1);
  uploadedFileSize.value = `${sizeMb} Mo`;

  const reader = new FileReader();
  reader.onload = (e) => {
    uploadedFileData.value = e.target?.result as string;
  };
  reader.readAsDataURL(file);
};

const triggerFileInput = () => {
  fileInputRef.value?.click();
};

const handleSubmit = async () => {
  if (!formTitre.value.trim()) {
    guichetStore.showToast("Le titre de l'épreuve est requis.", 'error');
    return;
  }

  if (!uploadedFileData.value) {
    guichetStore.showToast('Veuillez sélectionner un fichier PDF à téléverser.', 'error');
    return;
  }

  isSubmitting.value = true;
  try {
    await epreuvesStore.addEpreuve({
      titre: formTitre.value.trim(),
      categorie_permis: formCategorie.value,
      description: formDescription.value.trim() || undefined,
      annee: formAnnee.value,
      session: formSession.value.trim() || "Session d'entraînement",
      duree_minutes: formDuree.value,
      nombre_questions: formQuestions.value,
      avec_corrige: formAvecCorrige.value,
      fichier_nom: uploadedFileName.value,
      fichier_taille: uploadedFileSize.value || '1.0 Mo',
      fichier_data: uploadedFileData.value,
      ajoute_par: authStore.currentUser.prenom + ' ' + authStore.currentUser.nom + ' (' + authStore.currentUser.roleLabel + ')',
    });

    guichetStore.showToast(`Épreuve pour Permis ${formCategorie.value} enregistrée avec succès en base de données !`, 'success');
    emit('success');
    emit('close');
  } catch (err: any) {
    guichetStore.showToast(`Erreur d'enregistrement : ${err.message}`, 'error');
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
    <div class="bg-card border border-border w-full max-w-xl max-h-[90vh] rounded-[5px] shadow-2xl flex flex-col overflow-hidden text-foreground">
      <!-- En-tête -->
      <header class="px-5 py-3.5 border-b border-border flex items-center justify-between bg-muted/20">
        <div class="flex items-center space-x-2">
          <BookOpen class="w-4 h-4 text-brand" />
          <h2 class="text-sm font-bold text-foreground">Déposer une nouvelle épreuve PDF</h2>
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
        <!-- Titre de l'épreuve -->
        <div>
          <label class="block font-semibold text-foreground mb-1">Titre de l'épreuve *</label>
          <input
            v-model="formTitre"
            type="text"
            required
            placeholder="Ex: Épreuve Blanche Permis B 2024 · Spécial Priorités"
            class="w-full h-8 px-2.5 rounded-[4px] border border-border bg-background text-foreground text-xs focus:outline-none focus:border-brand-dark"
          />
        </div>

        <!-- Catégorie de permis & Année -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div>
            <label class="block font-semibold text-foreground mb-1">Catégorie de Permis Cible *</label>
            <select
              v-model="formCategorie"
              class="w-full h-8 px-2 rounded-[4px] border border-border bg-background text-foreground text-xs focus:outline-none focus:border-brand-dark font-medium"
            >
              <option v-for="cat in CATEGORIES_PERMIS" :key="cat.value" :value="cat.value">
                {{ cat.label }}
              </option>
            </select>
            <p class="text-[10px] text-muted-foreground mt-0.5">
              Seuls les élèves inscrits à cette catégorie y auront accès.
            </p>
          </div>

          <div>
            <label class="block font-semibold text-foreground mb-1">Session / Examen</label>
            <input
              v-model="formSession"
              type="text"
              placeholder="Ex: Session Normale 2024"
              class="w-full h-8 px-2.5 rounded-[4px] border border-border bg-background text-foreground text-xs focus:outline-none focus:border-brand-dark"
            />
          </div>
        </div>

        <!-- Paramètres du test (durée, questions, corrigé) -->
        <div class="grid grid-cols-3 gap-3">
          <div>
            <label class="block font-semibold text-foreground mb-1">Durée (minutes)</label>
            <input
              v-model.number="formDuree"
              type="number"
              min="5"
              max="180"
              class="w-full h-8 px-2.5 rounded-[4px] border border-border bg-background text-foreground text-xs focus:outline-none focus:border-brand-dark"
            />
          </div>

          <div>
            <label class="block font-semibold text-foreground mb-1">Nb. Questions</label>
            <input
              v-model.number="formQuestions"
              type="number"
              min="10"
              max="100"
              class="w-full h-8 px-2.5 rounded-[4px] border border-border bg-background text-foreground text-xs focus:outline-none focus:border-brand-dark"
            />
          </div>

          <div class="flex flex-col justify-end">
            <label class="flex items-center space-x-1.5 h-8 text-xs font-semibold text-foreground cursor-pointer select-none">
              <input
                v-model="formAvecCorrige"
                type="checkbox"
                class="rounded border-border text-brand focus:ring-0"
              />
              <span>Corrigé inclus</span>
            </label>
          </div>
        </div>

        <!-- Description pédagogique -->
        <div>
          <label class="block font-semibold text-foreground mb-1">Description pédagogique ou consignes</label>
          <textarea
            v-model="formDescription"
            rows="2"
            placeholder="Ex: Épreuve d'entraînement conforme aux dernières questions du Ministère..."
            class="w-full p-2 rounded-[4px] border border-border bg-background text-foreground text-xs focus:outline-none focus:border-brand-dark"
          />
        </div>

        <!-- Dépôt du fichier PDF -->
        <div>
          <label class="block font-semibold text-foreground mb-1">Document PDF du Sujet *</label>

          <!-- Input file caché -->
          <input
            ref="fileInputRef"
            type="file"
            accept="application/pdf"
            class="hidden"
            @change="handleFileSelect"
          />

          <!-- Zone de dépôt -->
          <div
            @click="triggerFileInput"
            class="p-4 border-2 border-dashed border-border hover:border-brand-border rounded-[4px] bg-muted/20 hover:bg-muted/40 cursor-pointer transition text-center space-y-1.5"
          >
            <UploadCloud class="w-6 h-6 text-brand mx-auto" />
            <div v-if="uploadedFileName" class="text-xs">
              <span class="font-bold text-foreground">{{ uploadedFileName }}</span>
              <span class="text-muted-foreground block text-[10px]">Taille : {{ uploadedFileSize }} · Prêt pour publication</span>
            </div>
            <div v-else class="text-xs">
              <span class="font-semibold text-foreground">Cliquez pour téléverser votre fichier PDF</span>
              <p class="text-[10px] text-muted-foreground">Format PDF uniquement (max 10 Mo)</p>
            </div>
          </div>
        </div>

        <!-- Pied de formulaire -->
        <div class="pt-3 border-t border-border flex items-center justify-end space-x-2">
          <button
            type="button"
            @click="$emit('close')"
            class="px-3 py-1.5 rounded-[4px] border border-border text-muted-foreground hover:text-foreground text-xs font-semibold transition cursor-pointer"
          >
            Annuler
          </button>

          <button
            type="submit"
            :disabled="isSubmitting"
            class="px-4 py-1.5 rounded-[4px] bg-brand hover:bg-brand-dark text-brand-dark font-bold text-xs uppercase tracking-wider transition cursor-pointer disabled:opacity-50"
          >
            {{ isSubmitting ? 'Publication en cours...' : "Publier l'épreuve" }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

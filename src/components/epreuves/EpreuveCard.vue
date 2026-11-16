<script setup lang="ts">
import type { EpreuvePdf } from '../../types';
import { getCategoryBadgeClass } from './epreuvesConstants';
import {
  Clock,
  HelpCircle,
  Eye,
  Download,
  Trash2,
  CheckCircle2,
  Volume2
} from 'lucide-vue-next';

interface Props {
  epreuve: EpreuvePdf;
  canDelete?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  canDelete: true,
});

defineEmits<{
  (e: 'preview', ep: EpreuvePdf): void;
  (e: 'download', ep: EpreuvePdf): void;
  (e: 'delete', ep: EpreuvePdf): void;
  (e: 'listen', ep: EpreuvePdf): void;
}>();
</script>

<template>
  <div class="p-4 rounded-[4px] bg-card border border-border hover:border-brand-border/60 transition shadow-xs flex flex-col justify-between space-y-3">
    <!-- En-tête carte -->
    <div class="space-y-2">
      <div class="flex items-center justify-between gap-2">
        <span
          class="inline-flex items-center text-[10px] font-bold px-2 py-0.5 rounded-[4px] border uppercase tracking-wider"
          :class="getCategoryBadgeClass(epreuve.categorie_permis)"
        >
          Permis {{ epreuve.categorie_permis }}
        </span>

        <span
          v-if="epreuve.avec_corrige"
          class="inline-flex items-center gap-1 text-[10px] font-medium text-emerald-600 dark:text-emerald-400 bg-emerald-500/10 px-1.5 py-0.5 rounded border border-emerald-500/20"
        >
          <CheckCircle2 class="w-3 h-3" /> Corrigé
        </span>
        <span v-else class="text-[10px] text-muted-foreground bg-muted px-1.5 py-0.5 rounded">
          Sujet seul
        </span>
      </div>

      <h3 class="text-xs sm:text-sm font-bold text-foreground leading-snug line-clamp-2">
        {{ epreuve.titre }}
      </h3>

      <p v-if="epreuve.description" class="text-xs text-muted-foreground line-clamp-2">
        {{ epreuve.description }}
      </p>
    </div>

    <!-- Métadonnées & Pédagogie -->
    <div class="space-y-2 pt-2 border-t border-border text-[11px] text-muted-foreground">
      <div class="flex items-center justify-between">
        <span class="flex items-center gap-1">
          <Clock class="w-3.5 h-3.5" /> {{ epreuve.duree_minutes }} min
        </span>
        <span class="flex items-center gap-1">
          <HelpCircle class="w-3.5 h-3.5" /> {{ epreuve.nombre_questions }} questions
        </span>
        <span class="font-mono text-[10px]">
          {{ epreuve.fichier_taille || 'PDF' }}
        </span>
      </div>

      <div class="flex items-center justify-between text-[10px]">
        <span>Session : <strong class="text-foreground">{{ epreuve.session || epreuve.annee }}</strong></span>
        <span>{{ epreuve.telechargements_count || 0 }} consultation(s)</span>
      </div>
    </div>

    <!-- Boutons d'action -->
    <div class="pt-2 flex items-center justify-between gap-2 border-t border-border">
      <button
        type="button"
        @click="$emit('preview', epreuve)"
        class="flex-1 flex items-center justify-center space-x-1.5 px-2.5 py-1.5 rounded-[4px] bg-muted hover:bg-muted/80 text-foreground text-xs font-medium transition cursor-pointer"
      >
        <Eye class="w-3.5 h-3.5 text-brand" />
        <span>Aperçu PDF</span>
      </button>

      <!-- Écouter vocalement -->
      <button
        type="button"
        @click="$emit('listen', epreuve)"
        class="flex items-center justify-center p-1.5 rounded-[4px] border border-border hover:bg-brand/10 text-muted-foreground hover:text-brand transition cursor-pointer"
        title="Écouter la lecture vocale du sujet"
      >
        <Volume2 class="w-4 h-4" />
      </button>

      <button
        type="button"
        @click="$emit('download', epreuve)"
        class="flex items-center justify-center p-1.5 rounded-[4px] border border-border hover:bg-muted text-foreground transition cursor-pointer"
        title="Télécharger le fichier PDF"
      >
        <Download class="w-4 h-4 text-emerald-600 dark:text-emerald-400" />
      </button>

      <button
        v-if="canDelete"
        type="button"
        @click="$emit('delete', epreuve)"
        class="flex items-center justify-center p-1.5 rounded-[4px] border border-border hover:bg-rose-500/10 text-rose-600 transition cursor-pointer"
        title="Supprimer l'épreuve"
      >
        <Trash2 class="w-4 h-4" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { EpreuvePdf } from '../../types';
import DocumentReaderModal from '../supports/DocumentReaderModal.vue';

const props = withDefaults(
  defineProps<{
    epreuve: EpreuvePdf | null;
    hideDownload?: boolean;
  }>(),
  {
    hideDownload: false,
  }
);

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'download', ep: EpreuvePdf): void;
}>();
</script>

<template>
  <DocumentReaderModal
    v-if="epreuve"
    :is-open="!!epreuve"
    :titre="epreuve.titre"
    :categorie-permis="epreuve.categorie_permis"
    :fichier-data="epreuve.fichier_data"
    :description="epreuve.description"
    :allow-download="!hideDownload"
    @download="emit('download', epreuve!)"
    @close="emit('close')"
  />
</template>

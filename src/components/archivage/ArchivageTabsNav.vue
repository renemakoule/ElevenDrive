<script setup lang="ts">
import { FolderArchive, Layers, FileText } from "lucide-vue-next";

export type ArchivageTab = "archives" | "eligibles" | "documents";

interface Props {
  activeTab: ArchivageTab;
  archivesCount: number;
  eligiblesCount: number;
  documentsCount: number;
}

defineProps<Props>();

const emit = defineEmits<{
  (e: "update:activeTab", tab: ArchivageTab): void;
}>();
</script>

<template>
  <div
    class="flex items-center space-x-1.5 border-b border-border bg-muted/15 p-1 rounded-[4px]"
  >
    <button
      type="button"
      @click="emit('update:activeTab', 'archives')"
      class="flex items-center space-x-1.5 px-3 py-1.5 rounded-md text-xs font-semibold transition"
      :class="
        activeTab === 'archives'
          ? 'bg-card text-foreground border border-border shadow-xs'
          : 'text-muted-foreground hover:text-foreground'
      "
    >
      <FolderArchive class="w-4 h-4 text-brand" />
      <span>Coffre-Fort des Archives ({{ archivesCount }})</span>
    </button>

    <button
      type="button"
      @click="emit('update:activeTab', 'eligibles')"
      class="flex items-center space-x-1.5 px-3 py-1.5 rounded-[4px] text-xs font-semibold transition"
      :class="
        activeTab === 'eligibles'
          ? 'bg-card text-foreground border border-border shadow-xs'
          : 'text-muted-foreground hover:text-foreground'
      "
    >
      <Layers class="w-4 h-4 text-amber-500" />
      <span>Dossiers Éligibles à l'Archivage ({{ eligiblesCount }})</span>
    </button>

    <button
      type="button"
      @click="emit('update:activeTab', 'documents')"
      class="flex items-center space-x-1.5 px-3 py-1.5 rounded-[4px] text-xs font-semibold transition"
      :class="
        activeTab === 'documents'
          ? 'bg-card text-foreground border border-border shadow-xs'
          : 'text-muted-foreground hover:text-foreground'
      "
    >
      <FileText class="w-4 h-4 text-emerald-600" />
      <span>Pièces & Documents Téléversés ({{ documentsCount }})</span>
    </button>
  </div>
</template>

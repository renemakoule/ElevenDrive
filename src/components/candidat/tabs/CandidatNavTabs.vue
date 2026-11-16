<script setup lang="ts">
import { computed } from "vue";
import {
  LayoutGrid,
  Video,
  FileText,
  Film,
  BookOpen,
  Brain,
  DownloadCloud,
  History,
} from "lucide-vue-next";

export type CandidatMainTab =
  | "dashboard"
  | "cours"
  | "supports"
  | "tutoriels"
  | "epreuves"
  | "quiz"
  | "downloads"
  | "historique";

const props = defineProps<{
  modelValue: CandidatMainTab;
  studentCategory: string;
  coursCount: number;
  supportsCount: number;
  videosCount: number;
  epreuvesCount: number;
  quizzesCount: number;
  downloadsCount: number;
}>();

const emit = defineEmits<{
  (e: "update:modelValue", tab: CandidatMainTab): void;
}>();

const tabs = computed(() => [
  {
    id: "dashboard" as CandidatMainTab,
    label: "Tableau de Bord & Reçus",
    icon: LayoutGrid,
    iconClass: "text-brand",
    count: undefined,
    activeBadgeClass: "",
  },
  {
    id: "cours" as CandidatMainTab,
    label: "Cours & Visio",
    icon: Video,
    iconClass: "text-sky-500",
    count: props.coursCount > 0 ? props.coursCount : undefined,
    activeBadgeClass: "bg-sky-500/20 text-sky-600 dark:text-sky-300",
  },
  {
    id: "supports" as CandidatMainTab,
    label: "Supports & Livrets",
    icon: FileText,
    iconClass: "text-emerald-600",
    count: props.supportsCount,
    activeBadgeClass:
      props.supportsCount > 0
        ? "bg-emerald-500/20 text-emerald-700 dark:text-emerald-300"
        : "bg-muted text-muted-foreground",
  },
  {
    id: "tutoriels" as CandidatMainTab,
    label: "Tutoriels Vidéo",
    icon: Film,
    iconClass: "text-red-500",
    count: props.videosCount,
    activeBadgeClass:
      props.videosCount > 0
        ? "bg-red-500/20 text-red-700 dark:text-red-300"
        : "bg-muted text-muted-foreground",
  },
  {
    id: "epreuves" as CandidatMainTab,
    label: `Mes Épreuves PDF (Permis ${props.studentCategory})`,
    icon: BookOpen,
    iconClass: "text-brand",
    count: props.epreuvesCount,
    activeBadgeClass:
      props.epreuvesCount > 0
        ? "bg-brand text-brand-dark"
        : "bg-muted text-muted-foreground",
  },
  {
    id: "quiz" as CandidatMainTab,
    label: "Quiz du Moniteur",
    icon: Brain,
    iconClass: "text-brand",
    count: props.quizzesCount,
    activeBadgeClass:
      props.quizzesCount > 0
        ? "bg-brand text-brand-dark"
        : "bg-muted text-muted-foreground",
  },
  {
    id: "downloads" as CandidatMainTab,
    label: "Hors-Ligne",
    icon: DownloadCloud,
    iconClass: "text-emerald-600",
    count: props.downloadsCount,
    activeBadgeClass:
      props.downloadsCount > 0
        ? "bg-emerald-500/20 text-emerald-700 dark:text-emerald-300"
        : "bg-muted text-muted-foreground",
  },
  {
    id: "historique" as CandidatMainTab,
    label: "Mon Historique",
    icon: History,
    iconClass: "text-brand",
    count: undefined,
    activeBadgeClass: "",
  },
]);
</script>

<template>
  <div
    class="flex items-center space-x-1.5 border border-border bg-muted/15 p-1 rounded-[3px] shrink-0 overflow-x-auto no-scrollbar"
  >
    <button
      v-for="t in tabs"
      :key="t.id"
      type="button"
      @click="emit('update:modelValue', t.id)"
      class="flex items-center space-x-1.5 min-h-[38px] sm:min-h-0 px-3 py-1.5 rounded-[3px] text-xs font-semibold transition shrink-0"
      :class="
        modelValue === t.id
          ? 'bg-card text-foreground border border-border'
          : 'text-muted-foreground hover:text-foreground border border-transparent'
      "
    >
      <component :is="t.icon" class="w-3.5 h-3.5" :class="t.iconClass" />
      <span>{{ t.label }}</span>
      <span
        v-if="t.count !== undefined"
        class="px-1.5 py-0.2 rounded-full text-[10px] font-bold"
        :class="t.activeBadgeClass"
      >
        {{ t.count }}
      </span>
    </button>
  </div>
</template>

<script setup lang="ts">
import type { NetworkAnnouncementItem } from "./types";
import {
  Calendar,
  Trash2,
  Clock,
  Building2,
  ExternalLink,
  Loader2,
} from "lucide-vue-next";

defineProps<{
  announcements: NetworkAnnouncementItem[];
  isLoading: boolean;
}>();

const emit = defineEmits<{
  (e: "delete", id: string): void;
}>();

const formatDate = (iso: string) => {
  try {
    const d = new Date(iso);
    return `${d.toLocaleDateString("fr-FR")} à ${d.toLocaleTimeString("fr-FR", { hour: "2-digit", minute: "2-digit" })}`;
  } catch {
    return iso;
  }
};
</script>

<template>
  <div class="lg:col-span-7 space-y-3">
    <span
      class="text-xs font-bold text-foreground uppercase tracking-wider flex items-center gap-2"
    >
      <span>Fil des annonces publiées ({{ announcements.length }})</span>
      <Loader2 v-if="isLoading" class="w-3.5 h-3.5 animate-spin text-brand" />
    </span>

    <div
      v-if="isLoading && announcements.length === 0"
      class="p-8 text-center text-xs text-muted-foreground border border-border rounded-[3px] bg-card flex flex-col items-center justify-center gap-2"
    >
      <Loader2 class="w-5 h-5 animate-spin text-brand" />
      <span>Chargement des annonces depuis la base de données...</span>
    </div>

    <div
      v-else-if="announcements.length === 0"
      class="p-8 text-center text-xs text-muted-foreground border border-border rounded-[3px] bg-card"
    >
      Aucune annonce n'a encore été diffusée.
    </div>

    <div
      v-for="ann in announcements"
      :key="ann.id"
      class="border border-border/80 rounded-[3px] bg-card overflow-hidden hover:border-brand-500/50 transition space-y-0"
    >
      <div
        v-if="ann.imageUrl"
        class="w-full h-48 overflow-hidden relative bg-neutral-950 flex items-center justify-center border-b border-border/50"
      >
        <img
          :src="ann.imageUrl"
          aria-hidden="true"
          class="absolute inset-0 w-full h-full object-cover blur-xl opacity-30 scale-110 pointer-events-none"
        />
        <img
          :src="ann.imageUrl"
          :alt="ann.titre"
          class="relative z-10 max-h-full max-w-full object-contain p-2"
        />
        <div class="absolute top-2 left-2 z-20 flex items-center gap-1.5">
          <span
            class="text-[10px] font-bold px-2 py-0.5 rounded bg-black/70 text-white backdrop-blur-xs uppercase shadow-sm"
          >
            {{ ann.categorie }}
          </span>
          <span
            class="text-[10px] font-medium px-2 py-0.5 rounded bg-brand/90 text-white backdrop-blur-xs shadow-sm"
          >
            {{ ann.cibleSchoolName }}
          </span>
        </div>
      </div>

      <div class="p-4 space-y-2.5">
        <div v-if="!ann.imageUrl" class="flex items-center justify-between">
          <div class="flex items-center gap-1.5">
            <span
              class="text-[10px] font-bold px-2 py-0.5 rounded bg-blue-500/10 text-blue-700 dark:text-blue-300 border border-blue-500/20 uppercase"
            >
              {{ ann.categorie }}
            </span>
            <span
              class="text-[11px] font-semibold text-foreground flex items-center gap-1 bg-muted px-2 py-0.5 rounded border border-border"
            >
              <Building2 class="w-3 h-3 text-muted-foreground" />
              <span>{{ ann.cibleSchoolName }}</span>
            </span>
          </div>
          <button
            type="button"
            @click="emit('delete', ann.id)"
            class="p-1 rounded text-muted-foreground hover:text-rose-600 transition"
            title="Supprimer"
          >
            <Trash2 class="w-3.5 h-3.5" />
          </button>
        </div>

        <div>
          <div class="flex items-start justify-between">
            <h4 class="text-sm font-bold text-foreground">{{ ann.titre }}</h4>
            <button
              v-if="ann.imageUrl"
              type="button"
              @click="emit('delete', ann.id)"
              class="p-1 rounded text-muted-foreground hover:text-rose-600 transition"
              title="Supprimer"
            >
              <Trash2 class="w-3.5 h-3.5" />
            </button>
          </div>
          <p
            v-if="ann.contenu"
            class="text-xs text-muted-foreground mt-1.5 leading-relaxed whitespace-pre-line"
          >
            {{ ann.contenu }}
          </p>
        </div>

        <div
          class="flex items-center justify-between pt-2 border-t border-border/50 text-[11px]"
        >
          <div class="flex items-center gap-3 text-muted-foreground">
            <span
              v-if="ann.dateEvenement"
              class="flex items-center gap-1 text-brand font-semibold"
            >
              <Calendar class="w-3.5 h-3.5" />
              <span
                >Événement le
                {{
                  new Date(ann.dateEvenement).toLocaleDateString("fr-FR")
                }}</span
              >
            </span>
            <span class="flex items-center gap-1 text-[10px]">
              <Clock class="w-3 h-3" />
              <span>{{ formatDate(ann.datePublication) }}</span>
            </span>
          </div>

          <a
            v-if="ann.lienAction"
            :href="ann.lienAction"
            target="_blank"
            class="text-brand hover:underline font-semibold flex items-center gap-1"
          >
            <span>Accéder</span>
            <ExternalLink class="w-3 h-3" />
          </a>
        </div>
      </div>
    </div>
  </div>
</template>


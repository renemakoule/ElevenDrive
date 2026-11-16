<script setup lang="ts">
import type { NetworkAnnouncement } from "../../../types";
import {
  X,
  Clock,
  ZoomIn,
  Maximize2,
  ExternalLink,
  EyeOff,
} from "lucide-vue-next";

defineProps<{
  announcement: NetworkAnnouncement;
}>();

const emit = defineEmits<{
  (e: "preview", ann: NetworkAnnouncement): void;
  (e: "dismiss", id: string): void;
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
  <div
    class="group relative border border-border/80 rounded-[4px] bg-gradient-to-b from-card to-muted/20 overflow-hidden hover:shadow-md transition-all duration-200 flex flex-col justify-between"
  >
    <!-- Bouton "Fermer ce panneau Ads" directement en haut à droite -->
    <button
      type="button"
      @click.stop="emit('dismiss', announcement.id)"
      class="absolute top-2.5 right-2.5 z-20 p-1.5 rounded-full bg-black/60 hover:bg-black text-white/90 hover:text-white backdrop-blur-xs transition shadow-sm cursor-pointer group-hover:scale-105"
      title="Fermer et masquer cette annonce"
    >
      <X class="w-3.5 h-3.5" />
    </button>

    <div>
      <!-- Vignette Image Ads responsive, 100% visible (sans rognage) avec fond flouté d'ambiance -->
      <div
        v-if="announcement.image_url"
        @click="emit('preview', announcement)"
        class="relative w-full h-40 sm:h-44 shrink-0 overflow-hidden bg-neutral-950 flex items-center justify-center cursor-pointer group/img select-none border-b border-border"
        title="Cliquer pour afficher l'image complète en plein écran"
      >
        <!-- Fond d'ambiance flouté dérivé de l'image -->
        <img
          :src="announcement.image_url"
          aria-hidden="true"
          class="absolute inset-0 max-h-full max-w-full object-contain blur-xl opacity-35 scale-110 pointer-events-none"
        />

        <!-- Image réelle intégrale (ratio naturel préservé) -->
        <img
          :src="announcement.image_url"
          :alt="announcement.titre"
          class="relative z-10 max-h-full max-w-full object-contain p-2 transition-transform duration-300 group-hover/img:scale-[1.03]"
        />

        <!-- Overlay interactif au survol -->
        <div
          class="absolute inset-0 z-20 bg-black/40 opacity-0 group-hover/img:opacity-100 transition-opacity flex items-center justify-center gap-1.5 text-white font-semibold text-xs backdrop-blur-xs"
        >
          <ZoomIn class="w-4 h-4 text-white" />
          <span>Plein écran</span>
        </div>

        <!-- Badge catégorie sur l'image -->
        <div class="absolute bottom-2 left-2 z-20 pointer-events-none">
          <span
            class="px-2 py-0.5 rounded-[3px] bg-brand text-white font-bold text-[9px] uppercase tracking-wider shadow-md"
          >
            Annonce · {{ announcement.categorie }}
          </span>
        </div>

        <!-- Micro-icône d'agrandissement en bas à droite -->
        <div
          class="absolute bottom-2 right-2 z-20 p-1 rounded bg-black/70 text-white/90 pointer-events-none"
        >
          <Maximize2 class="w-3 h-3" />
        </div>
      </div>

      <!-- Contenu textuel de l'annonce -->
      <div class="p-3.5 space-y-2">
        <div class="flex items-center justify-between gap-2">
          <span
            v-if="!announcement.image_url"
            class="text-[9px] px-2 py-0.5 rounded-[3px] bg-brand/10 text-brand font-bold uppercase tracking-wider"
          >
            Annonce · {{ announcement.categorie }}
          </span>
          <span
            class="text-[10px] text-muted-foreground flex items-center gap-1 font-mono ml-auto"
          >
            <Clock class="w-3 h-3" />
            <span>{{ formatDate(announcement.date_publication) }}</span>
          </span>
        </div>

        <h4
          @click="announcement.image_url ? emit('preview', announcement) : null"
          class="font-extrabold text-foreground text-sm leading-snug"
          :class="
            announcement.image_url
              ? 'cursor-pointer hover:text-brand transition'
              : ''
          "
        >
          {{ announcement.titre }}
        </h4>

        <p class="text-xs text-muted-foreground leading-relaxed line-clamp-3">
          {{ announcement.contenu }}
        </p>
      </div>
    </div>

    <!-- Pied de la carte Ads avec CTA et bouton masquer -->
    <div
      class="p-3.5 pt-0 flex items-center justify-between border-t border-border text-xs mt-2"
    >
      <div class="flex items-center gap-2 pt-2">
        <a
          v-if="announcement.lien_action"
          :href="announcement.lien_action"
          target="_blank"
          rel="noopener"
          class="inline-flex items-center gap-1 px-2.5 py-1 rounded-[3px] bg-brand text-white font-bold text-[10px] shadow-xs hover:bg-brand-dark transition cursor-pointer"
        >
          <span>En savoir plus</span>
          <ExternalLink class="w-3 h-3" />
        </a>

        <button
          v-if="announcement.image_url"
          type="button"
          @click="emit('preview', announcement)"
          class="inline-flex items-center gap-1 text-[10px] text-brand hover:underline font-semibold cursor-pointer"
        >
          <ZoomIn class="w-3 h-3" />
          <span>Agrandir</span>
        </button>
        <span
          v-else-if="!announcement.lien_action"
          class="text-[10px] text-muted-foreground italic"
        >
          ElevenDrive
        </span>
      </div>

      <button
        type="button"
        @click="emit('dismiss', announcement.id)"
        class="pt-2 text-[10px] text-muted-foreground hover:text-foreground flex items-center gap-1 cursor-pointer transition ml-auto"
        title="Masquer cette annonce"
      >
        <EyeOff class="w-3 h-3" />
        <span>Masquer</span>
      </button>
    </div>
  </div>
</template>

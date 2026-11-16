<script setup lang="ts">
import { Download, X } from 'lucide-vue-next';

interface Props {
  isOpen: boolean;
}

defineProps<Props>();

const emit = defineEmits<{
  (e: 'close'): void;
}>();
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm animate-in fade-in"
    @click.self="emit('close')"
  >
    <div class="w-full max-w-md rounded-[4px] border border-border bg-card p-6 space-y-4 animate-in zoom-in-95 shadow-2xl">
      <div class="flex items-center justify-between pb-3 border-b border-border">
        <div class="flex items-center gap-2">
          <Download class="w-5 h-5 text-brand" />
          <h4 class="text-sm font-bold text-foreground">Comment installer cette application PWA ?</h4>
        </div>
        <button
          type="button"
          @click="emit('close')"
          class="p-1 rounded text-muted-foreground hover:bg-muted transition"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <div class="space-y-3 text-xs text-muted-foreground leading-relaxed">
        <p>
          Cette Progressive Web App fonctionne comme une application installée sur votre smartphone ou tablette, sans téléchargement lourd depuis les stores :
        </p>

        <div class="p-3 rounded-[4px] border border-border bg-muted/40 space-y-1">
          <strong class="text-foreground block font-semibold">Sur Android (Google Chrome) :</strong>
          <span>Appuyez sur les 3 points en haut à droite, puis sélectionnez <strong>« Ajouter à l'écran d'accueil »</strong> ou <strong>« Installer l'application »</strong>.</span>
        </div>

        <div class="p-3 rounded-[4px] border border-border bg-muted/40 space-y-1">
          <strong class="text-foreground block font-semibold">Sur iPhone / iPad (Apple Safari) :</strong>
          <span>Appuyez sur l'icône de partage (le rectangle avec la flèche vers le haut), puis faites défiler et choisissez <strong>« Sur l'écran d'accueil »</strong>.</span>
        </div>
      </div>

      <div class="pt-2 flex justify-end">
        <button
          type="button"
          @click="emit('close')"
          class="px-4 py-2 rounded-[4px] bg-brand text-white text-xs font-semibold hover:bg-brand-dark transition shadow-xs"
        >
          J'ai compris
        </button>
      </div>
    </div>
  </div>
</template>

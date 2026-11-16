<script setup lang="ts">
import { ref, watch } from "vue";
import type { SchoolFeedbackItem } from "./types";
import { Sparkles, Send } from "lucide-vue-next";

const props = defineProps<{
  isOpen: boolean;
  isSubmitting: boolean;
  feedback: SchoolFeedbackItem | null;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "submit", reponse: string): void;
}>();

const reponse = ref("");

watch(
  () => [props.isOpen, props.feedback] as const,
  ([open, fb]) => {
    if (open && fb) {
      reponse.value = fb.reponseAdmin || "";
    }
  },
  { immediate: true },
);

const submit = () => {
  if (!reponse.value.trim()) return;
  emit("submit", reponse.value.trim());
};
</script>

<template>
  <div
    v-if="isOpen && feedback"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-xs p-4 select-none animate-in fade-in"
    @click.self="emit('close')"
  >
    <div
      class="bg-card border border-border rounded-[3px] max-w-lg w-full p-5 space-y-4 shadow-xl text-xs"
    >
      <div
        class="flex items-center justify-between border-b border-border pb-3"
      >
        <span
          class="font-bold text-sm text-foreground flex items-center gap-1.5"
        >
          <Sparkles class="w-4 h-4 text-brand" />
          <span>Réponse Officielle ElevenDrive</span>
        </span>
        <button
          type="button"
          @click="emit('close')"
          class="text-muted-foreground hover:text-foreground p-1"
        >
          ✕
        </button>
      </div>

      <div class="space-y-2 bg-muted/20 p-3 rounded border border-border/50">
        <div class="font-bold text-foreground">
          {{ feedback.autoEcoleNom }} • {{ feedback.auteur }}
        </div>
        <p class="text-xs text-muted-foreground italic">
          « {{ feedback.commentaire }} »
        </p>
      </div>

      <form @submit.prevent="submit" class="space-y-3">
        <div>
          <label class="block font-semibold text-muted-foreground mb-1">
            Votre réponse officielle (visible par l'auto-école)
          </label>
          <textarea
            v-model="reponse"
            rows="4"
            required
            placeholder="Ex: Merci pour votre retour ! Cette fonctionnalité est prise en compte..."
            class="w-full p-2.5 rounded border border-border bg-background focus:ring-1 focus:ring-brand"
          ></textarea>
        </div>

        <div class="flex justify-end gap-2 pt-2">
          <button
            type="button"
            @click="emit('close')"
            class="px-3 py-1.5 rounded border border-border hover:bg-muted"
          >
            Annuler
          </button>
          <button
            type="submit"
            :disabled="isSubmitting"
            class="px-4 py-1.5 rounded bg-brand hover:bg-brand-dark text-white font-semibold cursor-pointer disabled:opacity-50 flex items-center gap-1.5"
          >
            <Send class="w-3.5 h-3.5" />
            <span>{{ isSubmitting ? "Envoi..." : "Publier la réponse" }}</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

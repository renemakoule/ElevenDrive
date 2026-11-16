<script setup lang="ts">
import { ref } from "vue";
import type { SchoolItemSummary } from "../../../types";
import { AlertTriangle, Trash2, X, Loader2 } from "lucide-vue-next";

const props = defineProps<{
  isOpen: boolean;
  school: SchoolItemSummary | null;
  isLoading: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "confirm", schoolId: string, confirmName: string): void;
}>();

const confirmInput = ref("");

const handleConfirm = () => {
  if (!props.school) return;
  if (
    confirmInput.value.trim() !== props.school.nom_commercial &&
    confirmInput.value.trim() !== props.school.code_unique
  ) {
    alert(
      `Veuillez saisir exactement "${props.school.nom_commercial}" pour confirmer.`,
    );
    return;
  }
  emit("confirm", props.school.id, confirmInput.value.trim());
};
</script>

<template>
  <div
    v-if="isOpen && school"
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm"
  >
    <div
      class="w-full max-w-md bg-card border border-destructive/30 rounded-[8px] shadow-2xl overflow-hidden animate-in fade-in zoom-in-95 duration-150"
    >
      <!-- En-tête de danger -->
      <div
        class="h-12 border-b border-destructive/20 bg-destructive/10 px-4 flex items-center justify-between"
      >
        <div
          class="flex items-center gap-2 text-destructive font-semibold text-sm"
        >
          <AlertTriangle class="w-4 h-4 text-destructive" />
          <span>Zone Dangereuse · Suppression Définitive</span>
        </div>
        <button
          type="button"
          @click="emit('close')"
          class="text-muted-foreground hover:text-foreground p-1 rounded-[4px] transition-colors cursor-pointer"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Corps -->
      <div class="p-5 space-y-4 text-xs font-sans">
        <div
          class="p-3 rounded-[6px] bg-destructive/10 border border-destructive/20 text-destructive leading-relaxed"
        >
          <p class="font-bold mb-1">
            Attention : cette action est irréversible !
          </p>
          <p>
            Vous êtes sur le point de supprimer l'auto-école
            <strong class="text-foreground font-semibold">{{ school.nom_commercial }}</strong>
            ({{ school.agences_count }} agences,
            {{ school.users_count }} utilisateurs,
            {{ school.candidats_count }} candidats).
          </p>
        </div>

        <div>
          <label class="block text-foreground mb-1.5 font-medium">
            Pour confirmer, veuillez saisir le nom exact :
            <span class="text-destructive font-mono font-bold">{{
              school.nom_commercial
            }}</span>
          </label>
          <input
            v-model="confirmInput"
            type="text"
            placeholder="Tapez le nom de l'auto-école..."
            class="w-full h-9 px-3 rounded-[4px] bg-background border border-destructive/40 text-foreground font-mono text-xs focus:outline-none focus:ring-1 focus:ring-destructive focus:border-destructive"
          />
        </div>

        <!-- Actions -->
        <div
          class="pt-3 border-t border-border flex items-center justify-end gap-2.5"
        >
          <button
            type="button"
            @click="emit('close')"
            class="px-3 py-1.5 rounded-[4px] bg-muted hover:bg-muted/80 text-foreground border border-border font-medium transition cursor-pointer"
          >
            Annuler
          </button>
          <button
            type="button"
            @click="handleConfirm"
            :disabled="
              isLoading ||
              (confirmInput.trim() !== school.nom_commercial &&
                confirmInput.trim() !== school.code_unique)
            "
            class="px-4 py-1.5 rounded-[4px] bg-destructive text-destructive-foreground hover:opacity-90 font-medium flex items-center gap-1.5 transition cursor-pointer disabled:opacity-40 disabled:pointer-events-none shadow-xs"
          >
            <Loader2
              v-if="isLoading"
              class="w-3.5 h-3.5 animate-spin text-white"
            />
            <Trash2 v-else class="w-3.5 h-3.5" />
            <span>Supprimer définitivement</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

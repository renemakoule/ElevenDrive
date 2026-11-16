<script setup lang="ts">
import type { SchoolFeedbackItem } from "./types";
import {
  Star,
  Building2,
  MessageSquare,
  Trash2,
  Sparkles,
  Clock,
} from "lucide-vue-next";

defineProps<{
  feedback: SchoolFeedbackItem;
}>();

const emit = defineEmits<{
  (e: "reply", fb: SchoolFeedbackItem): void;
  (
    e: "change-status",
    id: string,
    status: "nouveau" | "en_cours" | "resolu",
  ): void;
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
  <div
    class="border border-border rounded-[3px] p-4 bg-card transition space-y-3 hover:border-brand-500/40"
  >
    <div class="flex items-start justify-between gap-3 flex-wrap">
      <div class="space-y-1">
        <div class="flex items-center gap-2 flex-wrap">
          <span
            class="font-bold text-xs text-foreground flex items-center gap-1.5"
          >
            <Building2 class="w-3.5 h-3.5 text-brand" />
            <span>{{ feedback.autoEcoleNom }}</span>
          </span>
          <span class="text-[11px] text-muted-foreground">
            • {{ feedback.auteur }} ({{ feedback.role }})
          </span>
          <span
            class="text-[10px] font-bold uppercase px-2 py-0.5 rounded-[4px]"
            :class="
              feedback.statut === 'resolu'
                ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400'
                : feedback.statut === 'en_cours'
                  ? 'bg-amber-500/10 text-amber-600 dark:text-amber-400'
                  : 'bg-blue-500/10 text-blue-600 dark:text-blue-400'
            "
          >
            {{
              feedback.statut === "resolu"
                ? "Pris en compte"
                : feedback.statut === "en_cours"
                  ? "En étude"
                  : "Nouveau"
            }}
          </span>
        </div>

        <div class="flex items-center gap-1">
          <Star
            v-for="st in 5"
            :key="st"
            class="w-3.5 h-3.5"
            :class="
              st <= feedback.etoiles
                ? 'fill-amber-400 text-amber-400'
                : 'text-muted-foreground/30'
            "
          />
          <span class="ml-1 text-[11px] font-bold text-foreground"
            >{{ feedback.etoiles }}/5</span
          >
          <span
            class="ml-2 text-[10px] text-muted-foreground px-1.5 py-0.2 rounded bg-muted"
          >
            {{ feedback.categorie }}
          </span>
        </div>
      </div>

      <div class="flex items-center gap-1.5 shrink-0">
        <button
          type="button"
          @click="emit('reply', feedback)"
          class="px-2 py-1 rounded text-[11px] font-medium border border-border bg-card hover:bg-muted text-foreground transition flex items-center gap-1 cursor-pointer"
        >
          <MessageSquare class="w-3 h-3 text-brand" />
          <span>{{
            feedback.reponseAdmin ? "Modifier réponse" : "Répondre"
          }}</span>
        </button>

        <select
          :value="feedback.statut"
          @change="
            emit(
              'change-status',
              feedback.id,
              ($event.target as HTMLSelectElement).value as any,
            )
          "
          class="text-[11px] p-1 rounded border border-border bg-background cursor-pointer"
        >
          <option value="nouveau">Nouveau</option>
          <option value="en_cours">En étude</option>
          <option value="resolu">Pris en compte</option>
        </select>
        <button
          type="button"
          @click="emit('delete', feedback.id)"
          class="p-1 rounded text-muted-foreground hover:text-rose-600 transition"
          title="Supprimer"
        >
          <Trash2 class="w-3.5 h-3.5" />
        </button>
      </div>
    </div>

    <p
      class="text-xs text-foreground leading-relaxed bg-muted/20 p-3 rounded-[4px] border border-border/40"
    >
      « {{ feedback.commentaire }} »
    </p>

    <div
      v-if="feedback.reponseAdmin"
      class="p-2.5 rounded-[4px] bg-brand-500/5 border border-brand-500/20 text-xs space-y-1"
    >
      <div class="font-bold text-[11px] text-brand flex items-center gap-1">
        <Sparkles class="w-3 h-3" />
        <span>Réponse de l'équipe ElevenDrive Core :</span>
      </div>
      <p class="text-[11px] text-muted-foreground leading-relaxed">
        {{ feedback.reponseAdmin }}
      </p>
    </div>

    <div class="text-[10px] text-muted-foreground flex items-center gap-1 pt-1">
      <Clock class="w-3 h-3" />
      <span>Reçu le {{ formatDate(feedback.date) }}</span>
    </div>
  </div>
</template>

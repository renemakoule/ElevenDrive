<script setup lang="ts">
import { useQuizFormsStore } from "../../../stores/quizFormsStore";
import { Brain, Plus, Users, Trash2 } from "lucide-vue-next";

const emit = defineEmits<{
  (e: "create-new"): void;
  (e: "view-responses", quizId: string): void;
}>();

const quizStore = useQuizFormsStore();

const formatDateHeure = (iso: string) => {
  if (!iso) return "À l'instant";
  const d = new Date(iso);
  if (isNaN(d.getTime())) return iso;
  return d.toLocaleDateString("fr-FR", {
    day: "2-digit",
    month: "short",
    hour: "2-digit",
    minute: "2-digit",
  });
};
</script>

<template>
  <div class="space-y-4">
    <div
      v-if="quizStore.formulaires.length === 0"
      class="rounded-lg border border-dashed border-border bg-card p-12 text-center space-y-3"
    >
      <Brain class="h-10 w-10 text-muted-foreground/40 mx-auto" />
      <h3 class="text-sm font-bold text-foreground">
        Aucun formulaire de quiz publié
      </h3>
      <p class="text-xs text-muted-foreground max-w-md mx-auto">
        Utilisez l'éditeur « Créer un Quiz » pour concevoir votre premier
        formulaire interactif et l'envoyer instantanément aux élèves.
      </p>
      <button
        type="button"
        @click="emit('create-new')"
        class="inline-flex items-center gap-1.5 rounded-lg bg-brand-dark px-4 py-2 text-xs font-bold text-white hover:bg-brand transition cursor-pointer"
      >
        <Plus class="h-4 w-4" />
        <span>Créer un formulaire de quiz</span>
      </button>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div
        v-for="qf in quizStore.formulaires"
        :key="qf.id"
        class="rounded-lg border border-border bg-card p-4 flex flex-col justify-between gap-3 shadow-xs hover:border-brand/50 transition"
      >
        <div class="space-y-2">
          <div class="flex items-center justify-between gap-2">
            <span
              class="rounded bg-brand-light px-2 py-0.5 text-[10px] font-bold uppercase text-brand-dark border border-brand-border"
            >
              Permis {{ qf.categorie_permis }}
            </span>
            <span class="text-[11px] text-muted-foreground font-mono">
              {{ formatDateHeure(qf.created_at) }}
            </span>
          </div>

          <h3 class="text-sm font-bold text-foreground">{{ qf.titre }}</h3>
          <p v-if="qf.description" class="text-xs text-muted-foreground line-clamp-2">
            {{ qf.description }}
          </p>
        </div>

        <div
          class="grid grid-cols-3 gap-2 rounded-md bg-muted/30 p-2.5 text-center border border-border text-xs"
        >
          <div>
            <div class="text-[10px] text-muted-foreground">Questions</div>
            <div class="font-bold text-foreground">
              {{ qf.questions?.length || 0 }}
            </div>
          </div>
          <div>
            <div class="text-[10px] text-muted-foreground">Durée / Seuil</div>
            <div class="font-bold text-foreground">
              {{ qf.duree_minutes }}m · {{ qf.seuil_reussite_pct }}%
            </div>
          </div>
          <div>
            <div class="text-[10px] text-muted-foreground">Copies reçues</div>
            <div class="font-extrabold text-brand-dark">
              {{
                quizStore.getSoumissionsByQuiz(qf.id).length ||
                qf.reponses_count ||
                0
              }}
            </div>
          </div>
        </div>

        <div class="flex items-center justify-between gap-2 pt-2 border-t border-border">
          <button
            type="button"
            @click="emit('view-responses', qf.id)"
            class="flex-1 inline-flex items-center justify-center gap-1.5 rounded-md bg-brand-light px-3 py-1.5 text-xs font-bold text-brand-dark hover:bg-brand-light/80 transition cursor-pointer"
          >
            <Users class="h-3.5 w-3.5 text-brand" />
            <span>
              Voir les réponses ({{ quizStore.getSoumissionsByQuiz(qf.id).length }})
            </span>
          </button>

          <button
            type="button"
            @click="quizStore.supprimerFormulaire(qf.id)"
            class="p-1.5 rounded-md border border-border text-muted-foreground hover:text-rose-600 hover:border-rose-500/40 transition cursor-pointer"
            title="Supprimer ce formulaire"
          >
            <Trash2 class="h-4 w-4" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>


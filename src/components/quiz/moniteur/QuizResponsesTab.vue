<script setup lang="ts">
import { ref, computed } from "vue";
import { useQuizFormsStore } from "../../../stores/quizFormsStore";
import type { QuizSoumissionEleve } from "../../../types";
import {
  Search,
  FileCheck2,
  CheckCircle2,
  XCircle,
  Eye,
} from "lucide-vue-next";

const props = defineProps<{
  selectedQuizFilter: string;
}>();

const emit = defineEmits<{
  (e: "update:selectedQuizFilter", id: string): void;
  (e: "inspect", sub: QuizSoumissionEleve): void;
}>();

const quizStore = useQuizFormsStore();
const searchStudentTerm = ref("");

const filteredSubmissions = computed(() => {
  let list = quizStore.soumissions;
  if (props.selectedQuizFilter !== "ALL") {
    list = list.filter((s) => s.quiz_id === props.selectedQuizFilter);
  }
  const q = searchStudentTerm.value.trim().toLowerCase();
  if (q) {
    list = list.filter(
      (s) =>
        s.candidat_nom.toLowerCase().includes(q) ||
        s.numero_dossier.toLowerCase().includes(q) ||
        s.quiz_titre.toLowerCase().includes(q),
    );
  }
  return list;
});

const statsReponses = computed(() => {
  const list = filteredSubmissions.value;
  const total = list.length;
  if (total === 0) {
    return { total: 0, moyennePct: 0, tauxReussite: 0, admisCount: 0 };
  }
  const sumPct = list.reduce((acc, s) => acc + s.pourcentage, 0);
  const admisCount = list.filter((s) => s.est_reussi).length;
  return {
    total,
    moyennePct: Math.round(sumPct / total),
    tauxReussite: Math.round((admisCount / total) * 100),
    admisCount,
  };
});

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

const formatDuration = (sec: number) => {
  if (!sec || sec <= 0) return "< 1 min";
  const m = Math.floor(sec / 60);
  const s = sec % 60;
  return m > 0 ? `${m} min ${s}s` : `${s}s`;
};
</script>

<template>
  <div class="space-y-4">
    <!-- 4 Cartes KPI des réponses reçues -->
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
      <div class="rounded-lg border border-border bg-card p-3.5">
        <div class="text-[11px] text-muted-foreground">
          Copies d'élèves reçues
        </div>
        <div class="mt-1 text-xl font-extrabold text-foreground tabular-nums">
          {{ statsReponses.total }}
        </div>
      </div>

      <div class="rounded-lg border border-border bg-card p-3.5">
        <div class="text-[11px] text-muted-foreground">Moyenne générale</div>
        <div class="mt-1 text-xl font-extrabold text-brand-dark tabular-nums">
          {{ statsReponses.moyennePct }} %
        </div>
      </div>

      <div class="rounded-lg border border-border bg-card p-3.5">
        <div class="text-[11px] text-muted-foreground">Taux de réussite</div>
        <div class="mt-1 text-xl font-extrabold text-emerald-600 tabular-nums">
          {{ statsReponses.tauxReussite }} %
        </div>
      </div>

      <div class="rounded-lg border border-border bg-card p-3.5">
        <div class="text-[11px] text-muted-foreground">
          Élèves admis au seuil
        </div>
        <div class="mt-1 text-xl font-extrabold text-foreground tabular-nums">
          {{ statsReponses.admisCount }} / {{ statsReponses.total }}
        </div>
      </div>
    </div>

    <!-- Barre de filtres -->
    <div
      class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3 rounded-lg border border-border bg-card p-3"
    >
      <div class="flex items-center gap-2 flex-1">
        <Search class="h-4 w-4 text-muted-foreground shrink-0" />
        <input
          v-model="searchStudentTerm"
          type="text"
          placeholder="Rechercher un élève par nom, prénom, N° de dossier ou titre du quiz..."
          class="w-full bg-transparent text-xs text-foreground focus:outline-none"
        />
      </div>

      <div class="flex items-center gap-2">
        <select
          :value="selectedQuizFilter"
          @change="
            emit(
              'update:selectedQuizFilter',
              ($event.target as HTMLSelectElement).value,
            )
          "
          class="rounded-md border border-border bg-background px-2.5 py-1.5 text-xs font-semibold text-foreground focus:border-brand focus:outline-none"
        >
          <option value="ALL">Tous les formulaires de quiz</option>
          <option
            v-for="qf in quizStore.formulaires"
            :key="qf.id"
            :value="qf.id"
          >
            {{ qf.titre }}
          </option>
        </select>
      </div>
    </div>

    <!-- Aucune réponse reçue -->
    <div
      v-if="filteredSubmissions.length === 0"
      class="rounded-lg border border-dashed border-border bg-card p-12 text-center space-y-2"
    >
      <FileCheck2 class="h-10 w-10 text-muted-foreground/40 mx-auto" />
      <h3 class="text-sm font-bold text-foreground">
        Aucune réponse d'élève reçue pour le moment
      </h3>
      <p class="text-xs text-muted-foreground max-w-md mx-auto">
        Dès qu'un élève valide un formulaire de quiz depuis son espace apprenant,
        sa copie complète et ses réponses question par question s'affichent
        instantanément ici.
      </p>
    </div>

    <!-- Tableau des réponses individuelles -->
    <div v-else class="overflow-hidden rounded-lg border border-border bg-card shadow-xs">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs">
          <thead
            class="border-b border-border bg-muted/50 text-[11px] font-semibold text-muted-foreground"
          >
            <tr>
              <th class="px-4 py-2.5">Élève / Candidat</th>
              <th class="px-3 py-2.5">Formulaire de Quiz</th>
              <th class="px-3 py-2.5 text-center">Score & Note</th>
              <th class="px-3 py-2.5 text-center">Statut</th>
              <th class="px-3 py-2.5">Durée & Date</th>
              <th class="px-4 py-2.5 text-right">Détail des réponses</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-border">
            <tr
              v-for="sub in filteredSubmissions"
              :key="sub.id"
              class="hover:bg-muted/30 transition"
            >
              <td class="px-4 py-3">
                <div class="font-bold text-foreground">
                  {{ sub.candidat_nom }}
                </div>
                <div class="font-mono text-[10px] text-muted-foreground">
                  {{ sub.numero_dossier || "Dossier Actif" }} · Permis
                  {{ sub.categorie_permis }}
                </div>
              </td>

              <td class="px-3 py-3">
                <div class="font-semibold text-foreground line-clamp-1">
                  {{ sub.quiz_titre }}
                </div>
                <div class="text-[10px] text-muted-foreground">
                  {{ sub.details?.length || 0 }} question(s) répondues
                </div>
              </td>

              <td class="px-3 py-3 text-center">
                <span class="font-mono font-extrabold text-sm text-foreground">
                  {{ sub.score_obtenu }} / {{ sub.total_points }} pts
                </span>
                <span class="block text-[10px] font-bold text-muted-foreground">
                  ({{ sub.pourcentage }}%)
                </span>
              </td>

              <td class="px-3 py-3 text-center">
                <span
                  class="inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-[10px] font-bold"
                  :class="
                    sub.est_reussi
                      ? 'bg-emerald-500/15 text-emerald-700 dark:text-emerald-300'
                      : 'bg-amber-500/15 text-amber-700 dark:text-amber-300'
                  "
                >
                  <CheckCircle2 v-if="sub.est_reussi" class="h-3 w-3" />
                  <XCircle v-else class="h-3 w-3" />
                  {{ sub.est_reussi ? "Réussi" : "À renforcer" }}
                </span>
              </td>

              <td class="px-3 py-3 text-muted-foreground">
                <div class="font-mono text-[11px]">
                  {{ formatDuration(sub.duree_secondes) }}
                </div>
                <div class="text-[10px]">
                  {{ formatDateHeure(sub.submitted_at) }}
                </div>
              </td>

              <td class="px-4 py-3 text-right">
                <button
                  type="button"
                  @click="emit('inspect', sub)"
                  class="inline-flex items-center gap-1.5 rounded-md bg-brand-light border border-brand-border px-2.5 py-1 text-xs font-bold text-brand-dark hover:bg-brand hover:text-white transition cursor-pointer"
                >
                  <Eye class="h-3.5 w-3.5" />
                  <span>Voir les réponses</span>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>


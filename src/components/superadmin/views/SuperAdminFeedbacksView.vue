<script lang="ts">
export type { SchoolFeedbackItem } from "../feedbacks/types";
</script>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import type { SchoolItemSummary } from "../../../types";
import { superAdminApi } from "../../../services/superAdminService";
import type { SchoolFeedbackItem } from "../feedbacks/types";
import FeedbackSummaryKpi from "../feedbacks/FeedbackSummaryKpi.vue";
import FeedbackCardItem from "../feedbacks/FeedbackCardItem.vue";
import FeedbackRequestModal from "../feedbacks/FeedbackRequestModal.vue";
import FeedbackTestModal from "../feedbacks/FeedbackTestModal.vue";
import FeedbackReplyModal from "../feedbacks/FeedbackReplyModal.vue";
import { Star, Send, Sparkles, Loader2 } from "lucide-vue-next";

const props = defineProps<{
  schools: SchoolItemSummary[];
  searchQuery?: string;
}>();

const feedbacks = ref<SchoolFeedbackItem[]>([]);
const selectedFilterStars = ref<number | "all">("all");
const isRequestModalOpen = ref(false);
const isAddTestModalOpen = ref(false);
const isReplyModalOpen = ref(false);
const replyTarget = ref<SchoolFeedbackItem | null>(null);
const isRequestSuccess = ref(false);
const isLoading = ref(false);
const isSubmitting = ref(false);
const errorMessage = ref<string | null>(null);

const mapFromApi = (item: any): SchoolFeedbackItem => ({
  id: item.id,
  autoEcoleId: item.auto_ecole_id || item.autoEcoleId || undefined,
  autoEcoleNom: item.auto_ecole_nom || item.autoEcoleNom || "Auto-École Pilote",
  auteur: item.auteur || "Équipe Auto-École",
  role: item.role || "Directeur",
  etoiles: item.etoiles || 5,
  categorie: item.categorie || "general",
  commentaire: item.commentaire || "",
  statut: item.statut || "nouveau",
  reponseAdmin: item.reponse_admin || item.reponseAdmin || undefined,
  date: item.created_at || item.date || new Date().toISOString(),
});

const loadFeedbacks = async () => {
  isLoading.value = true;
  errorMessage.value = null;
  try {
    const raw = await superAdminApi.getFeedbacks();
    feedbacks.value = (raw || []).map(mapFromApi);
  } catch (e: any) {
    errorMessage.value = e.message || "Erreur lors du chargement des avis";
  } finally {
    isLoading.value = false;
  }
};

onMounted(() => {
  loadFeedbacks();
});

const averageStars = computed(() => {
  if (feedbacks.value.length === 0) return 0;
  const sum = feedbacks.value.reduce((acc, f) => acc + f.etoiles, 0);
  return (sum / feedbacks.value.length).toFixed(1);
});

const starCounts = computed(() => {
  const counts: Record<1 | 2 | 3 | 4 | 5, number> = {
    5: 0,
    4: 0,
    3: 0,
    2: 0,
    1: 0,
  };
  feedbacks.value.forEach((f) => {
    if (f.etoiles >= 1 && f.etoiles <= 5) {
      counts[f.etoiles as 1 | 2 | 3 | 4 | 5]++;
    }
  });
  return counts;
});

const filteredFeedbacks = computed(() => {
  return feedbacks.value.filter((f) => {
    if (
      selectedFilterStars.value !== "all" &&
      f.etoiles !== selectedFilterStars.value
    ) {
      return false;
    }
    if (props.searchQuery) {
      const q = props.searchQuery.toLowerCase();
      return (
        f.autoEcoleNom.toLowerCase().includes(q) ||
        f.commentaire.toLowerCase().includes(q) ||
        f.auteur.toLowerCase().includes(q)
      );
    }
    return true;
  });
});

const envoyerDemandeAvis = async (payload: {
  cibleType: "all" | "single";
  cibleSchoolId: string;
  objet: string;
  message: string;
}) => {
  isSubmitting.value = true;
  errorMessage.value = null;
  try {
    await superAdminApi.createAnnouncement({
      titre: payload.objet,
      contenu: payload.message,
      type_format: "text_only",
      categorie: "general",
      cible_type: payload.cibleType,
      cible_school_id:
        payload.cibleType === "single" ? payload.cibleSchoolId : undefined,
    });

    isRequestSuccess.value = true;
    setTimeout(() => {
      isRequestSuccess.value = false;
      isRequestModalOpen.value = false;
    }, 2000);
  } catch (err: any) {
    errorMessage.value =
      err.message || "Impossible de diffuser la sollicitation d'avis";
  } finally {
    isSubmitting.value = false;
  }
};

const enregistrerAvisTest = async (payload: {
  autoEcoleId: string;
  auteur: string;
  role: string;
  etoiles: number;
  categorie: SchoolFeedbackItem["categorie"];
  commentaire: string;
}) => {
  isSubmitting.value = true;
  errorMessage.value = null;
  try {
    const created = await superAdminApi.createFeedback({
      auto_ecole_id: payload.autoEcoleId || undefined,
      auteur: payload.auteur.trim(),
      role: payload.role,
      etoiles: payload.etoiles,
      categorie: payload.categorie,
      commentaire: payload.commentaire.trim(),
    });
    feedbacks.value.unshift(mapFromApi(created));
    isAddTestModalOpen.value = false;
  } catch (err: any) {
    errorMessage.value = err.message || "Impossible d'enregistrer l'avis";
  } finally {
    isSubmitting.value = false;
  }
};

const openReplyModal = (fb: SchoolFeedbackItem) => {
  replyTarget.value = fb;
  isReplyModalOpen.value = true;
};

const envoyerReponseAdmin = async (reponseText: string) => {
  if (!replyTarget.value) return;
  isSubmitting.value = true;
  try {
    await superAdminApi.replyFeedback(
      replyTarget.value.id,
      reponseText,
      "resolu",
    );
    const target = feedbacks.value.find((f) => f.id === replyTarget.value?.id);
    if (target) {
      target.reponseAdmin = reponseText;
      target.statut = "resolu";
    }
    isReplyModalOpen.value = false;
  } catch (err: any) {
    alert(
      "Erreur lors de l'enregistrement de la réponse: " +
        (err.message || "Erreur serveur"),
    );
  } finally {
    isSubmitting.value = false;
  }
};

const changerStatut = async (
  fbId: string,
  nouveauStatut: "nouveau" | "en_cours" | "resolu",
) => {
  try {
    await superAdminApi.updateFeedbackStatus(fbId, nouveauStatut);
    const target = feedbacks.value.find((f) => f.id === fbId);
    if (target) {
      target.statut = nouveauStatut;
    }
  } catch (err: any) {
    alert("Erreur modification statut: " + (err.message || "Erreur serveur"));
  }
};

const supprimerAvis = async (fbId: string) => {
  if (!confirm("Supprimer définitivement cet avis de la base de données ?"))
    return;
  try {
    await superAdminApi.deleteFeedback(fbId);
    feedbacks.value = feedbacks.value.filter((f) => f.id !== fbId);
  } catch (err: any) {
    alert("Erreur suppression: " + (err.message || "Erreur serveur"));
  }
};
</script>

<template>
  <div class="space-y-6">
    <!-- En-tête -->
    <div class="flex items-center justify-between border-b border-border pb-4">
      <div>
        <h2 class="text-base font-bold text-foreground flex items-center gap-2">
          <Star class="w-5 h-5 text-amber-400 fill-amber-400" />
          <span>Avis & Satisfaction des Auto-Écoles</span>
        </h2>
        <p class="text-xs text-muted-foreground mt-0.5">
          Collectez, analysez et suivez les retours d'expérience et notations en
          étoiles de vos établissements partenaires.
        </p>
      </div>

      <div class="flex items-center gap-2">
        <button
          type="button"
          @click="isAddTestModalOpen = true"
          class="px-3 py-1.5 rounded-[4px] border border-border bg-card hover:bg-muted text-foreground text-xs font-medium transition cursor-pointer flex items-center gap-1.5"
        >
          <Sparkles class="w-3.5 h-3.5 text-brand" />
          <span>Simuler un avis</span>
        </button>

        <button
          type="button"
          @click="isRequestModalOpen = true"
          class="px-3 py-1.5 rounded-[4px] bg-brand hover:bg-brand-dark text-white text-xs font-semibold shadow-2xs transition cursor-pointer flex items-center gap-1.5"
        >
          <Send class="w-3.5 h-3.5" />
          <span>Solliciter les avis réseau</span>
        </button>
      </div>
    </div>

    <!-- Synthèse & Filtres -->
    <FeedbackSummaryKpi
      :average-stars="averageStars"
      :total-count="feedbacks.length"
      :filtered-count="filteredFeedbacks.length"
      :star-counts="starCounts"
      v-model:selected-filter-stars="selectedFilterStars"
    />

    <!-- Liste des avis -->
    <div class="space-y-3">
      <div
        v-if="isLoading && feedbacks.length === 0"
        class="p-8 text-center text-xs text-muted-foreground border border-border rounded-[3px] bg-card flex flex-col items-center justify-center gap-2"
      >
        <Loader2 class="w-5 h-5 animate-spin text-brand" />
        <span
          >Chargement des avis et notations depuis la base de données...</span
        >
      </div>

      <div
        v-else-if="filteredFeedbacks.length === 0"
        class="p-8 text-center text-xs text-muted-foreground border border-border rounded-[3px] bg-card"
      >
        Aucun avis ne correspond à ce filtre.
      </div>

      <FeedbackCardItem
        v-for="fb in filteredFeedbacks"
        :key="fb.id"
        :feedback="fb"
        @reply="openReplyModal"
        @change-status="changerStatut"
        @delete="supprimerAvis"
      />
    </div>

    <!-- Modales -->
    <FeedbackRequestModal
      :is-open="isRequestModalOpen"
      :is-success="isRequestSuccess"
      :schools="schools"
      @close="isRequestModalOpen = false"
      @submit="envoyerDemandeAvis"
    />

    <FeedbackTestModal
      :is-open="isAddTestModalOpen"
      :is-submitting="isSubmitting"
      :schools="schools"
      @close="isAddTestModalOpen = false"
      @submit="enregistrerAvisTest"
    />

    <FeedbackReplyModal
      :is-open="isReplyModalOpen"
      :is-submitting="isSubmitting"
      :feedback="replyTarget"
      @close="isReplyModalOpen = false"
      @submit="envoyerReponseAdmin"
    />
  </div>
</template>

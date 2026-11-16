<script lang="ts">
export type { NetworkAnnouncementItem } from "../announcements/types";
</script>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import type { SchoolItemSummary } from "../../../types";
import { superAdminApi } from "../../../services/superAdminService";
import type { NetworkAnnouncementItem } from "../announcements/types";
import AnnouncementFormCard, {
  type CreateAnnouncementFormPayload,
} from "../announcements/AnnouncementFormCard.vue";
import AnnouncementFeedList from "../announcements/AnnouncementFeedList.vue";
import { Megaphone } from "lucide-vue-next";

defineProps<{
  schools: SchoolItemSummary[];
  searchQuery?: string;
}>();

const announcements = ref<NetworkAnnouncementItem[]>([]);
const isFormSuccess = ref(false);
const isLoading = ref(false);
const isSubmitting = ref(false);
const errorMessage = ref<string | null>(null);
const formCardRef = ref<InstanceType<typeof AnnouncementFormCard> | null>(null);

const mapFromApi = (item: any): NetworkAnnouncementItem => ({
  id: item.id,
  titre: item.titre,
  contenu: item.contenu || "",
  typeFormat: item.type_format || item.typeFormat || "text_only",
  imageUrl: item.image_url || item.imageUrl || undefined,
  categorie: item.categorie || "general",
  cibleType: item.cible_type || item.cibleType || "all",
  cibleSchoolId: item.cible_school_id || item.cibleSchoolId || undefined,
  cibleSchoolName:
    item.cible_school_name || item.cibleSchoolName || "Toutes les auto-écoles",
  dateEvenement: item.date_evenement || item.dateEvenement || undefined,
  lienAction: item.lien_action || item.lienAction || undefined,
  datePublication:
    item.date_publication || item.datePublication || item.created_at,
});

const loadAnnouncements = async () => {
  isLoading.value = true;
  errorMessage.value = null;
  try {
    const raw = await superAdminApi.getAnnouncements();
    announcements.value = (raw || []).map(mapFromApi);
  } catch (e: any) {
    errorMessage.value = e.message || "Erreur lors du chargement des annonces";
  } finally {
    isLoading.value = false;
  }
};

onMounted(() => {
  loadAnnouncements();
});

const publierAnnonce = async (payload: CreateAnnouncementFormPayload) => {
  isSubmitting.value = true;
  errorMessage.value = null;
  try {
    const created = await superAdminApi.createAnnouncement({
      titre: payload.titre.trim() || "Annonce Illustrée",
      contenu: payload.contenu.trim(),
      type_format: payload.typeFormat,
      image_url: payload.imageUrl || undefined,
      categorie: payload.categorie,
      cible_type: payload.cibleType,
      cible_school_id:
        payload.cibleType === "single" ? payload.cibleSchoolId : undefined,
      date_evenement: payload.dateEvenement || undefined,
      lien_action: payload.lienAction || undefined,
    });

    announcements.value.unshift(mapFromApi(created));
    formCardRef.value?.resetFields();
    isFormSuccess.value = true;
    setTimeout(() => {
      isFormSuccess.value = false;
    }, 3000);
  } catch (err: any) {
    errorMessage.value = err.message || "Impossible de publier l'annonce";
  } finally {
    isSubmitting.value = false;
  }
};

const supprimerAnnonce = async (id: string) => {
  if (
    !confirm("Supprimer définitivement cette annonce de la base de données ?")
  )
    return;
  try {
    await superAdminApi.deleteAnnouncement(id);
    announcements.value = announcements.value.filter((a) => a.id !== id);
  } catch (err: any) {
    alert("Erreur suppression: " + (err.message || "Erreur serveur"));
  }
};
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between border-b border-border pb-4">
      <div>
        <h2 class="text-base font-bold text-foreground flex items-center gap-2">
          <Megaphone class="w-5 h-5 text-blue-500" />
          <span>Annonces Réseau & Événements</span>
        </h2>
        <p class="text-xs text-muted-foreground mt-0.5">
          Publiez des annonces d'événements, nouveautés ou messages officiels
          avec ou sans illustration visuelle.
        </p>
      </div>
      <span
        class="text-xs px-2.5 py-1 rounded-[4px] bg-blue-500/10 text-blue-700 dark:text-blue-300 border border-blue-500/20 font-medium"
      >
        {{ announcements.length }} annonce{{
          announcements.length > 1 ? "s" : ""
        }}
        en ligne
      </span>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
      <AnnouncementFormCard
        ref="formCardRef"
        :schools="schools"
        :is-submitting="isSubmitting"
        :is-form-success="isFormSuccess"
        :error-message="errorMessage"
        @submit="publierAnnonce"
      />

      <AnnouncementFeedList
        :announcements="announcements"
        :is-loading="isLoading"
        @delete="supprimerAnnonce"
      />
    </div>
  </div>
</template>

<script lang="ts">
export type { NetworkAlertItem } from "../alerts/types";
</script>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import type { SchoolItemSummary } from "../../../types";
import { superAdminApi } from "../../../services/superAdminService";
import type { NetworkAlertItem } from "../alerts/types";
import AlertCreateFormCard, {
  type CreateAlertPayload,
} from "../alerts/AlertCreateFormCard.vue";
import AlertFeedList from "../alerts/AlertFeedList.vue";
import { BellRing } from "lucide-vue-next";

defineProps<{
  schools: SchoolItemSummary[];
  searchQuery?: string;
}>();

const alerts = ref<NetworkAlertItem[]>([]);
const isFormSuccess = ref(false);
const isLoading = ref(false);
const isSubmitting = ref(false);
const errorMessage = ref<string | null>(null);
const formCardRef = ref<InstanceType<typeof AlertCreateFormCard> | null>(null);

const mapFromApi = (item: any): NetworkAlertItem => ({
  id: item.id,
  titre: item.titre,
  message: item.message,
  severite: item.severite,
  cibleType: item.cible_type || item.cibleType || "all",
  cibleSchoolId: item.cible_school_id || item.cibleSchoolId,
  cibleSchoolName:
    item.cible_school_name || item.cibleSchoolName || "Toutes les auto-écoles",
  bloquante: Boolean(item.bloquante),
  dateEnvoi: item.date_envoi || item.dateEnvoi || item.created_at,
  estActive:
    item.est_active !== undefined
      ? Boolean(item.est_active)
      : item.estActive !== undefined
        ? Boolean(item.estActive)
        : true,
});

const loadAlerts = async () => {
  isLoading.value = true;
  errorMessage.value = null;
  try {
    const raw = await superAdminApi.getAlerts();
    alerts.value = (raw || []).map(mapFromApi);
  } catch (e: any) {
    errorMessage.value =
      e.message || "Erreur lors de la récupération des alertes";
  } finally {
    isLoading.value = false;
  }
};

onMounted(() => {
  loadAlerts();
});

const envoyerAlerte = async (payload: CreateAlertPayload) => {
  isSubmitting.value = true;
  errorMessage.value = null;
  try {
    const created = await superAdminApi.createAlert({
      titre: payload.titre.trim(),
      message: payload.message.trim(),
      severite: payload.severite,
      cible_type: payload.cibleType,
      cible_school_id:
        payload.cibleType === "single" ? payload.cibleSchoolId : undefined,
      bloquante: payload.bloquante,
    });

    alerts.value.unshift(mapFromApi(created));
    formCardRef.value?.resetFields();
    isFormSuccess.value = true;
    setTimeout(() => {
      isFormSuccess.value = false;
    }, 3000);
  } catch (err: any) {
    errorMessage.value = err.message || "Impossible de diffuser l'alerte";
  } finally {
    isSubmitting.value = false;
  }
};

const toggleAlertStatus = async (alertId: string) => {
  try {
    const newStatus = await superAdminApi.toggleAlert(alertId);
    const target = alerts.value.find((a) => a.id === alertId);
    if (target) {
      target.estActive = newStatus;
    }
  } catch (err: any) {
    alert(
      "Erreur lors de la bascule de statut: " +
        (err.message || "Erreur serveur"),
    );
  }
};

const supprimerAlerte = async (alertId: string) => {
  if (!confirm("Supprimer définitivement cette alerte de la base de données ?"))
    return;
  try {
    await superAdminApi.deleteAlert(alertId);
    alerts.value = alerts.value.filter((a) => a.id !== alertId);
  } catch (err: any) {
    alert(
      "Erreur lors de la suppression: " + (err.message || "Erreur serveur"),
    );
  }
};

const activeAlertsCount = computed(
  () => alerts.value.filter((a) => a.estActive).length,
);
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between border-b border-border pb-4">
      <div>
        <h2 class="text-base font-bold text-foreground flex items-center gap-2">
          <BellRing class="w-5 h-5 text-amber-500" />
          <span>Diffusion d'Alertes Réseau</span>
        </h2>
        <p class="text-xs text-muted-foreground mt-0.5">
          Émettez des alertes urgentes ou de maintenance vers l'ensemble du
          réseau ou vers une auto-école spécifique.
        </p>
      </div>
      <div class="flex items-center gap-2">
        <span
          class="text-xs px-2.5 py-1 rounded-[4px] bg-amber-500/10 text-amber-700 dark:text-amber-300 border border-amber-500/20 font-medium"
        >
          {{ activeAlertsCount }} alerte{{
            activeAlertsCount > 1 ? "s" : ""
          }}
          active{{ activeAlertsCount > 1 ? "s" : "" }}
        </span>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
      <AlertCreateFormCard
        ref="formCardRef"
        :schools="schools"
        :is-submitting="isSubmitting"
        :is-form-success="isFormSuccess"
        :error-message="errorMessage"
        @submit="envoyerAlerte"
      />

      <AlertFeedList
        :alerts="alerts"
        :is-loading="isLoading"
        @toggle="toggleAlertStatus"
        @delete="supprimerAlerte"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from "vue";
import { useAuthStore } from "../../stores/authStore";
import { useGuichetStore } from "../../stores/guichetStore";
import { realtimeSocket } from "../../services/realtimeSocketService";
import type { StatutVehicule } from "../../types";
import AddVehiculeModal from "./flotte/AddVehiculeModal.vue";
import {
  AlertTriangle,
  ArrowLeft,
  Plus,
  RefreshCw,
  Radio,
} from "lucide-vue-next";

const authStore = useAuthStore();
const guichetStore = useGuichetStore();
const isWsConnected = realtimeSocket.isConnected;

onMounted(async () => {
  await guichetStore.loadVehicules();
});

const isAddModalOpen = ref(false);
const isSubmitting = ref(false);

const SEUIL_ALERTE_JOURS = 14;
const JOUR_MS = 86_400_000;

const joursRestants = (iso: string) => {
  if (!iso) return 999;
  const aujourdhui = new Date();
  aujourdhui.setHours(0, 0, 0, 0);
  const cible = new Date(iso);
  if (isNaN(cible.getTime())) return 999;
  return Math.ceil((cible.getTime() - aujourdhui.getTime()) / JOUR_MS);
};

const dateCourte = (iso: string) => {
  if (!iso) return "—";
  const d = new Date(iso);
  if (isNaN(d.getTime())) return "—";
  return d.toLocaleDateString("fr-FR", {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
  });
};

const lignes = computed(() =>
  guichetStore.vehicules.map((v) => ({
    ...v,
    joursVisite: joursRestants(v.date_visite_technique),
    joursAssurance: joursRestants(v.date_expiration_assurance),
  })),
);

const alertes = computed(() =>
  lignes.value.filter(
    (v) =>
      v.joursVisite <= SEUIL_ALERTE_JOURS ||
      v.joursAssurance <= SEUIL_ALERTE_JOURS,
  ),
);

const texteDelai = (j: number) =>
  j < 0
    ? `expirée depuis ${-j} j`
    : j === 0
      ? "expire aujourd'hui"
      : `expire dans ${j} j`;

const handleAddVehicule = async (payload: {
  immatriculation: string;
  marque_modele: string;
  type_boite: "MANUELLE" | "AUTOMATIQUE";
  kilometrage_actuel: number;
  statut: StatutVehicule;
  date_visite_technique: string;
  date_expiration_assurance: string;
}) => {
  if (!payload.immatriculation || !payload.marque_modele) {
    guichetStore.showToast(
      "L'immatriculation et le modèle sont obligatoires.",
      "error",
    );
    return;
  }
  if (!payload.date_visite_technique || !payload.date_expiration_assurance) {
    guichetStore.showToast(
      "Les dates de visite technique et d'assurance sont obligatoires.",
      "error",
    );
    return;
  }

  isSubmitting.value = true;
  try {
    await guichetStore.ajouterVehicule(payload);
    isAddModalOpen.value = false;
  } catch {
    // Notifications gérées dans guichetStore
  } finally {
    isSubmitting.value = false;
  }
};

const handleStatutChange = async (vehiculeId: string, event: Event) => {
  const target = event.target as HTMLSelectElement;
  const nouveauStatut = target.value as StatutVehicule;
  try {
    await guichetStore.modifierStatutVehicule(vehiculeId, nouveauStatut);
  } catch {
    await guichetStore.loadVehicules();
  }
};
</script>

<template>
  <div class="space-y-4">
    <!-- En-tête -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2">
        <button
          type="button"
          @click="authStore.setActiveSection('overview')"
          class="rounded-lg border border-border bg-card p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground transition"
          aria-label="Retour"
        >
          <ArrowLeft class="h-4 w-4" />
        </button>
        <div>
          <div class="flex items-center gap-2">
            <h1 class="text-sm font-semibold leading-tight text-foreground">
              Flotte automobile
            </h1>
            <span
              class="inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-bold"
              :class="
                isWsConnected
                  ? 'bg-emerald-500/15 text-emerald-700 dark:text-emerald-300'
                  : 'bg-amber-500/15 text-amber-700'
              "
            >
              <Radio class="h-3 w-3 animate-pulse" />
              {{ isWsConnected ? "Temps réel WS" : "Reconnexion WS..." }}
            </span>
          </div>
          <p class="text-xs text-muted-foreground">
            Parc de véhicules, contrôle technique et assurances officielles
            synchronisés en direct
          </p>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <button
          type="button"
          @click="guichetStore.loadVehicules()"
          class="min-h-[38px] sm:min-h-0 flex items-center gap-1.5 rounded-[3px] border border-border bg-card px-2.5 py-1.5 text-xs font-medium text-muted-foreground hover:bg-muted hover:text-foreground transition"
          title="Actualiser la liste"
        >
          <RefreshCw
            class="h-3.5 w-3.5"
            :class="{ 'animate-spin': guichetStore.flotteLoading }"
          />
          <span class="hidden sm:inline">Actualiser</span>
        </button>
        <button
          type="button"
          @click="isAddModalOpen = true"
          class="flex-1 sm:flex-initial min-h-[38px] sm:min-h-0 flex items-center justify-center gap-1.5 rounded-[3px] bg-brand-dark px-3 py-1.5 text-xs font-medium text-white hover:bg-brand transition"
        >
          <Plus class="h-3.5 w-3.5" />
          Ajouter un véhicule
        </button>
      </div>
    </div>

    <!-- Alertes d'expiration -->
    <div v-if="alertes.length > 0" class="space-y-2">
      <div
        v-for="v in alertes"
        :key="v.id"
        class="flex items-center gap-2 rounded-[5px] border border-amber-500/20 bg-amber-500/10 px-3.5 py-2 text-xs text-amber-600 dark:text-amber-400"
      >
        <AlertTriangle
          class="h-4 w-4 shrink-0 text-amber-600 dark:text-amber-400"
        />
        <span>
          <span class="font-semibold"
            >{{ v.marque_modele }} ({{ v.immatriculation }})</span
          >
          :
          <template v-if="v.joursVisite <= SEUIL_ALERTE_JOURS">
            Visite technique {{ texteDelai(v.joursVisite) }} ({{
              dateCourte(v.date_visite_technique)
            }}).
          </template>
          <template v-if="v.joursAssurance <= SEUIL_ALERTE_JOURS">
            Assurance {{ texteDelai(v.joursAssurance) }} ({{
              dateCourte(v.date_expiration_assurance)
            }}).
          </template>
        </span>
      </div>
    </div>

    <!-- Tableau du Parc Automobile -->
    <div class="overflow-hidden rounded-[3px] border border-border bg-card">
      <div
        class="flex items-center justify-between border-b border-border px-4 py-3"
      >
        <div>
          <span class="text-xs font-semibold text-foreground"
            >Parc automobile</span
          >
          <span class="ml-1.5 text-xs text-muted-foreground"
            >{{ guichetStore.vehicules.length }} véhicules enregistrés</span
          >
        </div>
        <div
          v-if="guichetStore.flotteLoading"
          class="text-xs text-muted-foreground animate-pulse"
        >
          Chargement...
        </div>
      </div>

      <!-- Cartes Mobiles (< md) -->
      <div class="md:hidden divide-y divide-border">
        <div
          v-if="
            guichetStore.vehicules.length === 0 && !guichetStore.flotteLoading
          "
          class="px-4 py-8 text-center text-xs text-muted-foreground"
        >
          Aucun véhicule enregistré dans le parc automobile.
        </div>
        <div v-for="v in lignes" :key="v.id" class="p-3.5 space-y-2.5 text-xs">
          <div class="flex items-start justify-between gap-2">
            <div class="min-w-0">
              <div class="font-semibold text-foreground truncate">
                {{ v.marque_modele }}
              </div>
              <div class="mt-0.5 flex items-center gap-1.5">
                <span
                  class="font-mono font-semibold text-[11px] text-foreground"
                  >{{ v.immatriculation }}</span
                >
                <span
                  class="rounded-[3px] bg-muted px-1.5 py-0.5 font-mono text-[10px] text-muted-foreground"
                  >{{ v.type_boite }}</span
                >
              </div>
            </div>
            <div
              class="font-mono text-xs font-medium text-foreground tabular-nums shrink-0"
            >
              {{ (v.kilometrage_actuel || 0).toLocaleString("fr-FR") }} km
            </div>
          </div>

          <div class="grid grid-cols-2 gap-2 text-[11px]">
            <div
              class="rounded-[3px] border border-border bg-muted/20 px-2.5 py-1.5"
            >
              <span class="block text-muted-foreground">Visite tech.</span>
              <span
                class="font-mono font-semibold tabular-nums"
                :class="
                  v.joursVisite <= SEUIL_ALERTE_JOURS
                    ? 'text-rose-600 dark:text-rose-400'
                    : 'text-foreground'
                "
              >
                {{ dateCourte(v.date_visite_technique) }}
              </span>
            </div>
            <div
              class="rounded-[3px] border border-border bg-muted/20 px-2.5 py-1.5"
            >
              <span class="block text-muted-foreground">Assurance</span>
              <span
                class="font-mono font-semibold tabular-nums"
                :class="
                  v.joursAssurance <= SEUIL_ALERTE_JOURS
                    ? 'text-rose-600 dark:text-rose-400'
                    : 'text-foreground'
                "
              >
                {{ dateCourte(v.date_expiration_assurance) }}
              </span>
            </div>
          </div>

          <select
            :value="v.statut"
            @change="handleStatutChange(v.id, $event)"
            class="h-10 w-full rounded-[3px] border border-border bg-background px-2.5 text-xs font-medium text-foreground transition focus:outline-none focus:ring-1 focus:ring-brand"
          >
            <option value="DISPONIBLE">Disponible</option>
            <option value="EN_LECON">En leçon</option>
            <option value="EN_ENTRETIEN">En entretien</option>
            <option value="EN_PANNE">En panne</option>
          </select>
        </div>
      </div>

      <div class="hidden md:block overflow-x-auto">
        <table class="w-full text-left text-xs">
          <thead class="bg-muted text-muted-foreground">
            <tr>
              <th class="px-3 py-2 font-medium">Immatriculation</th>
              <th class="px-3 py-2 font-medium">Marque & Modèle</th>
              <th class="px-3 py-2 font-medium">Boîte</th>
              <th class="px-3 py-2 text-right font-medium">Kilométrage</th>
              <th class="px-3 py-2 font-medium">Visite technique</th>
              <th class="px-3 py-2 font-medium">Assurance</th>
              <th class="px-3 py-2 font-medium">Statut opérationnel</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-border">
            <tr
              v-if="
                guichetStore.vehicules.length === 0 &&
                !guichetStore.flotteLoading
              "
            >
              <td
                colspan="7"
                class="px-4 py-8 text-center text-xs text-muted-foreground"
              >
                Aucun véhicule enregistré dans le parc automobile. Cliquez sur
                "Ajouter un véhicule" pour enregistrer vos véhicules.
              </td>
            </tr>
            <tr
              v-for="v in lignes"
              :key="v.id"
              class="hover:bg-muted/50 transition"
            >
              <td class="px-3 py-2 font-mono font-semibold text-foreground">
                {{ v.immatriculation }}
              </td>
              <td class="px-3 py-2 font-medium text-foreground">
                {{ v.marque_modele }}
              </td>
              <td class="px-3 py-2 text-muted-foreground">
                <span
                  class="rounded bg-muted px-1.5 py-0.5 font-mono text-[10px]"
                  >{{ v.type_boite }}</span
                >
              </td>
              <td
                class="px-3 py-2 text-right font-mono text-foreground tabular-nums"
              >
                {{ (v.kilometrage_actuel || 0).toLocaleString("fr-FR") }} km
              </td>
              <td
                class="px-3 py-2 tabular-nums"
                :class="
                  v.joursVisite <= SEUIL_ALERTE_JOURS
                    ? 'font-semibold text-rose-600 dark:text-rose-400'
                    : 'text-muted-foreground'
                "
              >
                {{ dateCourte(v.date_visite_technique) }}
              </td>
              <td
                class="px-3 py-2 tabular-nums"
                :class="
                  v.joursAssurance <= SEUIL_ALERTE_JOURS
                    ? 'font-semibold text-rose-600 dark:text-rose-400'
                    : 'text-muted-foreground'
                "
              >
                {{ dateCourte(v.date_expiration_assurance) }}
              </td>
              <td class="px-3 py-2">
                <select
                  :value="v.statut"
                  @change="handleStatutChange(v.id, $event)"
                  class="rounded-md border border-border bg-background px-2 py-1 text-xs font-medium text-foreground transition focus:outline-none focus:ring-1 focus:ring-brand"
                >
                  <option value="DISPONIBLE">Disponible</option>
                  <option value="EN_LECON">En leçon</option>
                  <option value="EN_ENTRETIEN">En entretien</option>
                  <option value="EN_PANNE">En panne</option>
                </select>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <AddVehiculeModal
      :is-open="isAddModalOpen"
      :is-submitting="isSubmitting"
      @close="isAddModalOpen = false"
      @submit="handleAddVehicule"
    />
  </div>
</template>

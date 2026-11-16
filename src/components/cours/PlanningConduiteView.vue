<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from "vue";
import { api } from "../../db/api";
import { useAuthStore } from "../../stores/authStore";
import { useGuichetStore } from "../../stores/guichetStore";
import { usePlanningConduiteConfigStore } from "../../stores/planningConduiteConfigStore";
import { realtimeSocket } from "../../services/realtimeSocketService";
import PlanifierConduiteModal from "./PlanifierConduiteModal.vue";
import ConfigCreneauxModal from "./ConfigCreneauxModal.vue";
import PlanningConduiteFilters from "./planning/PlanningConduiteFilters.vue";
import PlanningConduiteGrid from "./planning/PlanningConduiteGrid.vue";
import PlanningConduiteTable from "./planning/PlanningConduiteTable.vue";
import { extraireCategorie } from "./planning/planningConduiteUtils";
import type { SeanceConduite } from "../../types";
import {
  Car,
  Video,
  Calendar,
  Plus,
  RefreshCw,
  ChevronLeft,
  ChevronRight,
  LayoutGrid,
  List,
  Sliders,
} from "lucide-vue-next";

const authStore = useAuthStore();
const guichetStore = useGuichetStore();
const configStore = usePlanningConduiteConfigStore();

// Navigation par dates
const selectedDate = ref(new Date().toISOString().slice(0, 10));
const viewMode = ref<"grid" | "table">("grid");

// Filtres
const selectedMoniteurId = ref("");
const selectedVehiculeId = ref("");
const selectedStatut = ref("ALL");
const selectedCategorie = ref("ALL");
const searchQuery = ref("");

// Modales
const isPlanningModalOpen = ref(false);
const isConfigModalOpen = ref(false);
const modalInitialDate = ref("");
const modalInitialHeureDebut = ref("08:00");
const modalInitialHeureFin = ref("10:00");
const modalInitialMoniteurId = ref("");

// Données
const seances = ref<SeanceConduite[]>([]);
const moniteurs = ref<
  { id: string; nomComplet: string; role: string; telephone?: string }[]
>([]);
const isLoading = ref(false);
const isCancellingId = ref<string | null>(null);

let unsubscribeWs: (() => void) | null = null;

onMounted(async () => {
  configStore.init();
  await Promise.all([
    chargerMoniteurs(),
    guichetStore.vehicules.length === 0
      ? guichetStore.loadVehicules()
      : Promise.resolve(),
  ]);
  await chargerPlanning();

  unsubscribeWs = realtimeSocket.subscribe((msg) => {
    if (
      msg.type === "PLANNING_CONDUITE_UPDATED" ||
      msg.type === "LIVRET_EMARGEMENT_UPDATED"
    ) {
      chargerPlanning();
    } else if (msg.type === "FLOTTE_VEHICULES_UPDATED") {
      guichetStore.loadVehicules();
    }
  });
});

onUnmounted(() => {
  if (unsubscribeWs) {
    unsubscribeWs();
    unsubscribeWs = null;
  }
});

watch(selectedDate, () => {
  chargerPlanning();
});

const chargerMoniteurs = async () => {
  try {
    moniteurs.value = await api.getMoniteurs();
  } catch (err) {
    console.warn("Erreur chargement moniteurs:", err);
  }
};

const chargerPlanning = async () => {
  isLoading.value = true;
  try {
    const data = await api.getPlanningConduite(
      selectedDate.value,
      selectedMoniteurId.value ? selectedMoniteurId.value : undefined,
    );
    seances.value = data || [];
  } catch (err: any) {
    guichetStore.showToast(
      err.message || "Erreur lors du chargement du planning conduite",
      "error",
    );
  } finally {
    isLoading.value = false;
  }
};

// Navigation par dates
const goToToday = () => {
  selectedDate.value = new Date().toISOString().slice(0, 10);
};

const jumpDay = (delta: number) => {
  const d = new Date(selectedDate.value);
  d.setDate(d.getDate() + delta);
  selectedDate.value = d.toISOString().slice(0, 10);
};

const isToday = computed(
  () => selectedDate.value === new Date().toISOString().slice(0, 10),
);

const formattedDateLong = computed(() => {
  try {
    const [year, month, day] = selectedDate.value.split("-").map(Number);
    const dateObj = new Date(year, month - 1, day);
    return dateObj.toLocaleDateString("fr-FR", {
      weekday: "long",
      day: "numeric",
      month: "long",
      year: "numeric",
    });
  } catch {
    return selectedDate.value;
  }
});

// Filtrage en mémoire
const seancesFiltrees = computed(() => {
  return seances.value.filter((s) => {
    if (
      selectedCategorie.value !== "ALL" &&
      extraireCategorie(s.formule_nom) !== selectedCategorie.value
    ) {
      return false;
    }
    if (
      selectedMoniteurId.value &&
      s.moniteur_id !== selectedMoniteurId.value
    ) {
      return false;
    }
    if (
      selectedVehiculeId.value &&
      s.vehicule_id !== selectedVehiculeId.value
    ) {
      return false;
    }
    if (selectedStatut.value !== "ALL" && s.statut !== selectedStatut.value) {
      return false;
    }
    if (searchQuery.value.trim()) {
      const q = searchQuery.value.toLowerCase().trim();
      const nom = (s.candidat_nom || "").toLowerCase();
      const tel = (s.telephone_whatsapp || "").toLowerCase();
      const mon = (s.moniteur_nom || "").toLowerCase();
      const veh =
        `${s.vehicule_immatriculation || ""} ${s.vehicule_modele || ""}`.toLowerCase();
      if (
        !nom.includes(q) &&
        !tel.includes(q) &&
        !mon.includes(q) &&
        !veh.includes(q)
      ) {
        return false;
      }
    }
    return true;
  });
});

// Statistiques du jour
const stats = computed(() => {
  const total = seancesFiltrees.value.length;
  return {
    total,
    totalHeures: total * 2,
    uniqueMoniteurs: new Set(
      seancesFiltrees.value.map((s) => s.moniteur_id).filter(Boolean),
    ).size,
    uniqueVehicules: new Set(
      seancesFiltrees.value.map((s) => s.vehicule_id).filter(Boolean),
    ).size,
    effectuees: seancesFiltrees.value.filter((s) => s.statut === "EFFECTUE")
      .length,
  };
});

// Groupement des séances selon les créneaux personnalisés
const creneauxGroupes = computed(() => {
  const list = seancesFiltrees.value;
  return configStore.creneauxActifs.map((creneau) => ({
    ...creneau,
    seances: list.filter(
      (s) =>
        s.heure_debut === creneau.debut ||
        s.heure_debut?.startsWith(creneau.debut.slice(0, 2)),
    ),
  }));
});

const seancesHorsCreneaux = computed(() => {
  const standardsDebuts = configStore.creneauxActifs.map((c) => c.debut);
  return seancesFiltrees.value.filter(
    (s) =>
      !standardsDebuts.some(
        (d) => s.heure_debut === d || s.heure_debut?.startsWith(d.slice(0, 2)),
      ),
  );
});

const ouvrirPlanificateur = (initialCreneau?: {
  debut: string;
  fin: string;
}) => {
  modalInitialDate.value = selectedDate.value;
  modalInitialHeureDebut.value = initialCreneau?.debut || "08:00";
  modalInitialHeureFin.value = initialCreneau?.fin || "10:00";
  modalInitialMoniteurId.value = selectedMoniteurId.value || "";
  isPlanningModalOpen.value = true;
};

const annulerSeance = async (s: SeanceConduite) => {
  const confirmMsg = `Êtes-vous sûr de vouloir annuler la leçon de ${s.candidat_nom} prévue à ${s.heure_debut} ? Ce créneau sera de nouveau libre.`;
  if (!window.confirm(confirmMsg)) return;

  isCancellingId.value = s.id;
  try {
    await api.annulerSeanceConduite(s.id);
    guichetStore.showToast(
      `Séance de conduite de ${s.candidat_nom} annulée.`,
      "info",
    );
    await chargerPlanning();
  } catch (err: any) {
    guichetStore.showToast(
      err.message || "Erreur lors de l'annulation",
      "error",
    );
  } finally {
    isCancellingId.value = null;
  }
};

const envoyerWhatsApp = (s: SeanceConduite) => {
  const cleanPhone = (s.telephone_whatsapp || "").replace(/[^0-9]/g, "");
  if (!cleanPhone) {
    guichetStore.showToast(
      "Numéro de téléphone non renseigné pour cet élève.",
      "info",
    );
    return;
  }
  const msg = `Bonjour ${s.candidat_nom}, nous vous rappelons votre leçon de conduite pratique prévue le ${s.date_seance} de ${s.heure_debut} à ${s.heure_fin} avec votre formateur ${s.moniteur_nom || "Auto-École"}. Merci de vous munir de votre livret d'apprentissage.`;
  window.open(
    `https://wa.me/${cleanPhone}?text=${encodeURIComponent(msg)}`,
    "_blank",
    "noopener,noreferrer",
  );
};
</script>

<template>
  <div class="space-y-5 p-4 sm:p-6 max-w-7xl mx-auto">
    <!-- Onglet de bascule Pédagogie -->
    <div
      class="flex items-center justify-between flex-wrap gap-3 pb-2 border-b border-border"
    >
      <div
        class="flex items-center gap-1.5 p-1 bg-muted/60 rounded-[6px] border border-border text-xs font-semibold"
      >
        <button
          type="button"
          @click="authStore.setActiveSection('cours-planning')"
          class="px-3 py-1.5 rounded-[4px] flex items-center gap-2 transition cursor-pointer text-muted-foreground hover:text-foreground"
        >
          <Video class="w-3.5 h-3.5 text-sky-600" />
          <span>📚 Cours Théoriques (Code & Visio)</span>
        </button>

        <button
          type="button"
          @click="authStore.setActiveSection('conduite-planning')"
          class="px-3 py-1.5 rounded-[4px] flex items-center gap-2 transition cursor-pointer bg-card text-foreground font-bold border border-border"
        >
          <Car class="w-3.5 h-3.5 text-brand" />
          <span>🚗 Conduite Pratique (Moniteurs & Flotte)</span>
        </button>
      </div>

      <div class="flex items-center space-x-2">
        <button
          type="button"
          @click="isConfigModalOpen = true"
          class="flex items-center space-x-1.5 px-3 py-2 rounded-[4px] bg-muted hover:bg-muted/80 text-foreground font-semibold text-xs transition border border-border cursor-pointer"
          title="Configurer vos créneaux horaires et jours de passage"
        >
          <Sliders class="w-3.5 h-3.5 text-brand" />
          <span class="hidden sm:inline">Créneaux & Jours de passage</span>
        </button>

        <button
          type="button"
          @click="ouvrirPlanificateur()"
          class="flex items-center space-x-1.5 px-3.5 py-2 rounded-[4px] bg-brand hover:bg-brand-dark text-white font-bold text-xs uppercase tracking-wider transition cursor-pointer"
        >
          <Plus class="w-4 h-4" />
          <span>Planifier une leçon</span>
        </button>
      </div>
    </div>

    <!-- En-tête de section + Date -->
    <header
      class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 pb-2"
    >
      <div class="space-y-1">
        <div class="flex items-center space-x-2">
          <div
            class="w-8 h-8 rounded-[4px] bg-brand/10 border border-border text-brand flex items-center justify-center"
          >
            <Car class="w-4 h-4" />
          </div>
          <h1 class="text-base sm:text-lg font-bold text-foreground">
            Planning & Supervision des Leçons de Conduite
          </h1>
        </div>
        <p class="text-xs text-muted-foreground">
          Visualisez la grille personnalisée des moniteurs, les jours de passage
          et planifiez les leçons par catégorie de permis.
        </p>
      </div>

      <div
        class="flex items-center space-x-2 bg-card p-1.5 rounded-[6px] border border-border"
      >
        <button
          type="button"
          @click="jumpDay(-1)"
          class="p-1.5 rounded text-muted-foreground hover:text-foreground hover:bg-muted transition cursor-pointer"
          title="Jour précédent"
        >
          <ChevronLeft class="w-4 h-4" />
        </button>

        <button
          type="button"
          @click="goToToday"
          class="px-2.5 py-1 text-xs rounded font-semibold transition cursor-pointer"
          :class="
            isToday
              ? 'bg-brand text-white'
              : 'bg-muted text-foreground hover:bg-muted/80'
          "
        >
          Aujourd'hui
        </button>

        <input
          type="date"
          v-model="selectedDate"
          class="px-2 py-1 text-xs rounded bg-background border border-border focus:border-brand outline-none text-foreground font-medium"
        />

        <button
          type="button"
          @click="jumpDay(1)"
          class="p-1.5 rounded text-muted-foreground hover:text-foreground hover:bg-muted transition cursor-pointer"
          title="Jour suivant"
        >
          <ChevronRight class="w-4 h-4" />
        </button>

        <button
          type="button"
          @click="chargerPlanning"
          :disabled="isLoading"
          class="p-1.5 rounded text-muted-foreground hover:text-foreground hover:bg-muted transition cursor-pointer"
          title="Actualiser"
        >
          <RefreshCw class="w-4 h-4" :class="{ 'animate-spin': isLoading }" />
        </button>
      </div>
    </header>

    <!-- Bandeau date courante -->
    <div
      class="flex items-center justify-between flex-wrap gap-2 px-4 py-2.5 rounded-[4px] bg-muted/40 border border-border"
    >
      <div class="flex items-center space-x-2 flex-wrap gap-1.5">
        <Calendar class="w-4 h-4 text-brand" />
        <span class="text-xs font-bold capitalize text-foreground">
          {{ formattedDateLong }}
        </span>
        <span
          v-if="isToday"
          class="text-[10px] px-1.5 py-0.5 rounded bg-brand text-white font-bold border border-border"
        >
          Aujourd'hui
        </span>
        <span
          class="text-[10px] font-semibold px-2 py-0.5 rounded-full flex items-center gap-1"
          :class="
            configStore.isJourDePassage(selectedDate)
              ? 'bg-emerald-500/10 text-emerald-600 border border-emerald-500/20'
              : 'bg-amber-500/10 text-amber-600 border border-amber-500/20'
          "
        >
          <span>{{
            configStore.isJourDePassage(selectedDate)
              ? "✓ Jour de formation actif"
              : "⚠️ Jour habituellement fermé"
          }}</span>
        </span>
      </div>

      <div class="flex items-center space-x-2">
        <button
          type="button"
          @click="isConfigModalOpen = true"
          class="text-[11px] text-muted-foreground hover:text-foreground transition flex items-center gap-1 font-medium cursor-pointer"
        >
          <Sliders class="w-3 h-3 text-brand" />
          <span>{{ configStore.joursActifsLabels }}</span>
        </button>

        <div
          class="flex items-center space-x-1 bg-card p-1 rounded-[4px] border border-border"
        >
          <button
            type="button"
            @click="viewMode = 'grid'"
            class="p-1.5 rounded text-xs transition cursor-pointer flex items-center gap-1.5"
            :class="
              viewMode === 'grid'
                ? 'bg-muted text-foreground font-bold'
                : 'text-muted-foreground hover:text-foreground'
            "
            title="Vue Grille des Créneaux"
          >
            <LayoutGrid class="w-3.5 h-3.5" />
            <span class="text-[11px] hidden sm:inline">Créneaux</span>
          </button>
          <button
            type="button"
            @click="viewMode = 'table'"
            class="p-1.5 rounded text-xs transition cursor-pointer flex items-center gap-1.5"
            :class="
              viewMode === 'table'
                ? 'bg-muted text-foreground font-bold'
                : 'text-muted-foreground hover:text-foreground'
            "
            title="Vue Liste Détaillée"
          >
            <List class="w-3.5 h-3.5" />
            <span class="text-[11px] hidden sm:inline">Tableau</span>
          </button>
        </div>
      </div>
    </div>

    <!-- KPIs & Filtres -->
    <PlanningConduiteFilters
      :stats="stats"
      v-model:search-query="searchQuery"
      v-model:selected-categorie="selectedCategorie"
      v-model:selected-moniteur-id="selectedMoniteurId"
      v-model:selected-vehicule-id="selectedVehiculeId"
      v-model:selected-statut="selectedStatut"
      :moniteurs="moniteurs"
      :vehicules="guichetStore.vehicules"
      :filtered-count="seancesFiltrees.length"
    />

    <!-- Vue Grille ou Vue Tableau -->
    <PlanningConduiteGrid
      v-if="viewMode === 'grid'"
      :creneaux-groupes="creneauxGroupes"
      :seances-hors-creneaux="seancesHorsCreneaux"
      :is-cancelling-id="isCancellingId"
      @assign="ouvrirPlanificateur"
      @whatsapp="envoyerWhatsApp"
      @cancel="annulerSeance"
    />

    <PlanningConduiteTable
      v-else
      :seances="seancesFiltrees"
      :is-cancelling-id="isCancellingId"
      @whatsapp="envoyerWhatsApp"
      @cancel="annulerSeance"
    />

    <!-- Modales -->
    <PlanifierConduiteModal
      :isOpen="isPlanningModalOpen"
      :initialDate="modalInitialDate"
      :initialHeureDebut="modalInitialHeureDebut"
      :initialHeureFin="modalInitialHeureFin"
      :initialMoniteurId="modalInitialMoniteurId"
      @close="isPlanningModalOpen = false"
      @saved="chargerPlanning"
    />

    <ConfigCreneauxModal
      :isOpen="isConfigModalOpen"
      @close="isConfigModalOpen = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";
import { api } from "../../db/api";
import { useGuichetStore } from "../../stores/guichetStore";
import { usePlanningConduiteConfigStore } from "../../stores/planningConduiteConfigStore";
import ConfigCreneauxModal from "./ConfigCreneauxModal.vue";
import PlanifierMultiCandidatSelector from "./planning/PlanifierMultiCandidatSelector.vue";
import type { Candidat, SeanceConduite } from "../../types";
import {
  X,
  Calendar,
  Clock,
  Car,
  User,
  AlertCircle,
  Sliders,
} from "lucide-vue-next";

const props = withDefaults(
  defineProps<{
    isOpen: boolean;
    initialDate?: string;
    initialHeureDebut?: string;
    initialHeureFin?: string;
    initialMoniteurId?: string;
  }>(),
  {
    isOpen: false,
    initialDate: "",
    initialHeureDebut: "08:00",
    initialHeureFin: "10:00",
    initialMoniteurId: "",
  },
);

const emit = defineEmits<{
  (e: "close"): void;
  (e: "saved", seance: SeanceConduite): void;
}>();

const guichetStore = useGuichetStore();
const configStore = usePlanningConduiteConfigStore();

// Sélection multi-élèves
const allCandidats = ref<Candidat[]>([]);
const selectedCandidats = ref<Candidat[]>([]);
const dossiersMap = ref<Record<string, any>>({});
const isConfigModalOpen = ref(false);

// Formulaire
const formDate = ref(
  props.initialDate || new Date().toISOString().slice(0, 10),
);
const isHoraireLibre = ref(false);
const selectedCreneauId = ref<string>("");
const formHeureDebut = ref(props.initialHeureDebut || "08:00");
const formHeureFin = ref(props.initialHeureFin || "10:00");
const formMoniteurId = ref(props.initialMoniteurId || "");
const formVehiculeId = ref("");
const isSubmitting = ref(false);
const formError = ref<string | null>(null);

// Moniteurs disponibles
interface MoniteurItem {
  id: string;
  nomComplet: string;
  role: string;
  telephone?: string;
}
const moniteurs = ref<MoniteurItem[]>([]);
const isLoadingMoniteurs = ref(false);

const extraireCategorie = (nom?: string) => {
  if (!nom) return "B";
  const upper = nom.toUpperCase();
  if (
    upper.includes("PERMIS A") ||
    upper.includes("MOTO") ||
    upper.includes("A1") ||
    upper.includes("A2")
  )
    return "A";
  if (
    upper.includes("PERMIS C") ||
    upper.includes("POIDS LOURD") ||
    upper.includes("CAMION")
  )
    return "C";
  if (
    upper.includes("PERMIS D") ||
    upper.includes("TRANSPORT") ||
    upper.includes("BUS")
  )
    return "D";
  if (upper.includes("PERMIS E") || upper.includes("REMORQUE")) return "E";
  return "B";
};

const chargerDonnees = async () => {
  configStore.init();
  isLoadingMoniteurs.value = true;
  try {
    const [monList, candList] = await Promise.all([
      api.getMoniteurs(),
      api.searchCandidats(""),
      guichetStore.vehicules.length === 0
        ? guichetStore.loadVehicules()
        : Promise.resolve(),
    ]);
    moniteurs.value = monList;
    allCandidats.value = candList || [];

    if (monList.length > 0 && !formMoniteurId.value) {
      formMoniteurId.value = monList[0].id;
    }
    if (guichetStore.vehicules.length > 0 && !formVehiculeId.value) {
      formVehiculeId.value = guichetStore.vehicules[0].id;
    }

    // Aligner avec les créneaux configurés
    if (configStore.creneauxActifs.length > 0) {
      const match = configStore.creneauxActifs.find(
        (c) => c.debut === formHeureDebut.value && c.fin === formHeureFin.value,
      );
      if (match) {
        selectedCreneauId.value = match.id;
      } else {
        selectedCreneauId.value = configStore.creneauxActifs[0].id;
        formHeureDebut.value = configStore.creneauxActifs[0].debut;
        formHeureFin.value = configStore.creneauxActifs[0].fin;
      }
    }
  } catch (err) {
    console.warn("Erreur chargement données planification:", err);
  } finally {
    isLoadingMoniteurs.value = false;
  }
};

onMounted(() => {
  chargerDonnees();
});

watch(
  () => props.isOpen,
  (open) => {
    if (open) {
      configStore.init();
      formDate.value =
        props.initialDate || new Date().toISOString().slice(0, 10);
      formHeureDebut.value = props.initialHeureDebut || "08:00";
      formHeureFin.value = props.initialHeureFin || "10:00";
      if (props.initialMoniteurId) {
        formMoniteurId.value = props.initialMoniteurId;
      }
      formError.value = null;

      // Détecter si l'horaire fourni correspond à un créneau existant
      const match = configStore.creneauxActifs.find(
        (c) => c.debut === formHeureDebut.value && c.fin === formHeureFin.value,
      );
      if (match) {
        selectedCreneauId.value = match.id;
        isHoraireLibre.value = false;
      } else {
        isHoraireLibre.value = true;
      }
    }
  },
);

const toggleCandidat = async (c: Candidat) => {
  const exists = selectedCandidats.value.some((item) => item.id === c.id);
  if (exists) {
    removeCandidat(c.id);
    return;
  }

  selectedCandidats.value = [...selectedCandidats.value, c];
  if (!dossiersMap.value[c.id]) {
    try {
      const detail = await api.getCandidatDossier(c.id);
      if (detail) {
        dossiersMap.value = { ...dossiersMap.value, [c.id]: detail };
      }

      if (selectedCandidats.value.length === 1) {
        const cat =
          c.categorie_permis ||
          extraireCategorie(detail?.formule?.intitule || c.formule_nom);
        const vehiculesAdaptes = guichetStore.vehicules.filter((v) => {
          const immatEtModele =
            `${v.marque_modele} ${v.immatriculation}`.toLowerCase();
          if (cat === "A")
            return (
              immatEtModele.includes("moto") ||
              immatEtModele.includes("yamaha") ||
              immatEtModele.includes("honda")
            );
          if (cat === "C" || cat === "D" || cat === "E")
            return (
              immatEtModele.includes("camion") ||
              immatEtModele.includes("bus") ||
              immatEtModele.includes("poids")
            );
          return !immatEtModele.includes("moto");
        });
        if (vehiculesAdaptes.length > 0) {
          formVehiculeId.value = vehiculesAdaptes[0].id;
        }
      }
    } catch (err) {
      console.warn("Erreur chargement dossier:", err);
    }
  }
};

const removeCandidat = (id: string) => {
  selectedCandidats.value = selectedCandidats.value.filter((c) => c.id !== id);
};

const clearAllCandidats = () => {
  selectedCandidats.value = [];
};

// Choix de créneau horaire
const choisirCreneauPredefini = (cr: {
  id: string;
  debut: string;
  fin: string;
}) => {
  selectedCreneauId.value = cr.id;
  isHoraireLibre.value = false;
  formHeureDebut.value = cr.debut;
  formHeureFin.value = cr.fin;
};

const activerHoraireLibre = () => {
  isHoraireLibre.value = true;
  selectedCreneauId.value = "";
};

// Calcul de la durée en clair
const dureeCalculee = computed(() => {
  try {
    const [hD, mD] = formHeureDebut.value.split(":").map(Number);
    const [hF, mF] = formHeureFin.value.split(":").map(Number);
    const debutMin = hD * 60 + mD;
    const finMin = hF * 60 + mF;
    const diff = finMin - debutMin;
    if (diff <= 0) return "Horaires invalides";
    const h = Math.floor(diff / 60);
    const m = diff % 60;
    if (m === 0) return `${h} heure${h > 1 ? "s" : ""}`;
    return `${h}h${m.toString().padStart(2, "0")}`;
  } catch {
    return "2 heures";
  }
});

// Vérification si le jour sélectionné est un jour de passage autorisé
const isJourDePassageActif = computed(() => {
  return configStore.isJourDePassage(formDate.value);
});

const nomJourDate = computed(() => {
  return configStore.getNomJour(formDate.value);
});

const handleSubmit = async () => {
  formError.value = null;

  if (selectedCandidats.value.length === 0) {
    formError.value = "Veuillez sélectionner au moins un élève.";
    return;
  }
  if (!formDate.value) {
    formError.value = "Veuillez sélectionner une date.";
    return;
  }
  if (!formHeureDebut.value || !formHeureFin.value) {
    formError.value = "Veuillez renseigner les horaires de la leçon.";
    return;
  }
  if (formHeureDebut.value >= formHeureFin.value) {
    formError.value = "L'heure de début doit être antérieure à l'heure de fin.";
    return;
  }
  if (!formMoniteurId.value) {
    formError.value = "Veuillez affecter un moniteur référent.";
    return;
  }

  isSubmitting.value = true;
  try {
    const createdList: SeanceConduite[] = [];
    const errorsList: string[] = [];

    for (const cand of selectedCandidats.value) {
      let detail = dossiersMap.value[cand.id];
      if (!detail) {
        try {
          detail = await api.getCandidatDossier(cand.id);
        } catch {
          detail = null;
        }
      }
      const dossierId = detail?.dossier?.id;
      if (!dossierId) {
        errorsList.push(`${cand.nom} ${cand.prenoms} : dossier introuvable`);
        continue;
      }

      try {
        const seance = await api.reserverSeanceConduite({
          dossier_id: dossierId,
          date_seance: formDate.value,
          heure_debut: formHeureDebut.value,
          heure_fin: formHeureFin.value,
          moniteur_id: formMoniteurId.value || undefined,
          vehicule_id: formVehiculeId.value || undefined,
        });
        createdList.push(seance);
      } catch (e: any) {
        errorsList.push(
          `${cand.nom} ${cand.prenoms} : ${e?.message || "erreur"}`,
        );
      }
    }

    if (createdList.length > 0) {
      guichetStore.showToast(
        createdList.length === 1
          ? `Leçon de conduite programmée pour ${selectedCandidats.value[0].nom} (${formHeureDebut.value} - ${formHeureFin.value}) !`
          : `${createdList.length} leçon(s) de conduite programmée(s) (${formHeureDebut.value} - ${formHeureFin.value}) !`,
        errorsList.length > 0 ? "info" : "success",
      );
      emit("saved", createdList[0]);
      if (errorsList.length === 0) {
        selectedCandidats.value = [];
        emit("close");
      } else {
        formError.value = `Certains élèves n'ont pas pu être planifiés : ${errorsList.join(" ; ")}`;
      }
    } else {
      formError.value =
        errorsList.join(" ; ") || "Erreur lors de la réservation de la leçon.";
    }
  } finally {
    isSubmitting.value = false;
  }
};
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-50 flex items-center justify-center p-3 sm:p-5 bg-black/80 backdrop-blur-sm animate-in fade-in"
    role="dialog"
    aria-modal="true"
  >
    <div
      class="bg-card border border-border w-full max-w-2xl rounded-[4px] flex flex-col overflow-hidden text-foreground animate-in zoom-in-95 duration-150"
    >
      <!-- En-tête -->
      <header
        class="px-5 py-3 border-b border-border flex items-center justify-between bg-muted"
      >
        <div class="flex items-center space-x-2.5">
          <div
            class="w-8 h-8 rounded-[4px] bg-brand/10 border border-border text-white flex items-center justify-center"
          >
            <Car class="w-4 h-4" />
          </div>
          <div>
            <h2 class="text-sm font-bold text-foreground">
              Planifier une Leçon de Conduite Pratique
            </h2>
            <p class="text-[11px] text-muted-foreground">
              Affectation personnalisée : élève par catégorie de permis, créneau
              horaire, moniteur et véhicule.
            </p>
          </div>
        </div>

        <button
          type="button"
          @click="emit('close')"
          class="p-1 rounded-[4px] text-muted-foreground hover:text-foreground transition cursor-pointer"
        >
          <X class="w-5 h-5" />
        </button>
      </header>

      <!-- Corps du formulaire -->
      <form
        @submit.prevent="handleSubmit"
        class="p-5 space-y-4 overflow-y-auto max-h-[75vh]"
      >
        <!-- Alerte erreur -->
        <div
          v-if="formError"
          class="p-3 rounded-[4px] bg-red-500/10 border border-red-500/20 text-red-600 dark:text-red-400 text-xs flex items-start gap-2"
        >
          <AlertCircle class="w-4 h-4 shrink-0 mt-0.5" />
          <span>{{ formError }}</span>
        </div>

        <!-- 1. SÉLECTION D'UN OU PLUSIEURS ÉLÈVES -->
        <PlanifierMultiCandidatSelector
          :all-candidats="allCandidats"
          :selected-candidats="selectedCandidats"
          :dossiers-map="dossiersMap"
          @toggle-candidat="toggleCandidat"
          @remove-candidat="removeCandidat"
          @clear-all="clearAllCandidats"
        />

        <!-- ══════════════════════════════════════════════════════
             2. DATE & VÉRIFICATION DES JOURS DE PASSAGE
             ══════════════════════════════════════════════════════ -->
        <div
          class="grid grid-cols-1 sm:grid-cols-2 gap-3 pt-2 border-t border-border"
        >
          <!-- Date -->
          <div class="space-y-1">
            <label
              class="block text-xs font-bold text-foreground flex items-center gap-1.5"
            >
              <Calendar class="w-3.5 h-3.5 text-brand" />
              <span>Date de la séance <span class="text-red-500">*</span></span>
            </label>
            <input
              type="date"
              v-model="formDate"
              required
              class="w-full px-3 py-2 text-xs rounded-[4px] bg-background border border-border focus:border-brand outline-none"
            />

            <!-- Alerte jour de passage non coché -->
            <div
              v-if="!isJourDePassageActif"
              class="text-[10px] text-amber-600 dark:text-amber-400 flex items-center gap-1 font-medium pt-0.5"
            >
              <AlertCircle class="w-3 h-3 shrink-0" />
              <span
                >{{ nomJourDate }} n'est pas coché dans vos jours de passage
                habituels (dérogation accordée).</span
              >
            </div>
          </div>

          <!-- Moniteur -->
          <div class="space-y-1">
            <label
              class="block text-xs font-bold text-foreground flex items-center gap-1.5"
            >
              <User class="w-3.5 h-3.5 text-brand" />
              <span
                >Moniteur en charge <span class="text-red-500">*</span></span
              >
            </label>
            <select
              v-model="formMoniteurId"
              required
              class="w-full px-3 py-2 text-xs rounded-[4px] bg-background border border-border focus:border-brand outline-none cursor-pointer"
            >
              <option v-for="m in moniteurs" :key="m.id" :value="m.id">
                {{ m.nomComplet }} ({{ m.role }})
              </option>
            </select>
          </div>
        </div>

        <!-- ══════════════════════════════════════════════════════
             3. CRÉNEAUX PERSONNALISÉS & HORAIRES LIBRES
             ══════════════════════════════════════════════════════ -->
        <div class="space-y-2 pt-2 border-t border-border">
          <div class="flex items-center justify-between">
            <label
              class="block text-xs font-bold text-foreground flex items-center gap-1.5"
            >
              <Clock class="w-3.5 h-3.5 text-brand" />
              <span>Créneaux horaires de l'établissement</span>
            </label>

            <button
              type="button"
              @click="isConfigModalOpen = true"
              class="text-[11px] text-brand hover:underline flex items-center gap-1 font-semibold cursor-pointer"
            >
              <Sliders class="w-3 h-3" />
              <span>Configurer mes créneaux & jours</span>
            </button>
          </div>

          <!-- Grille des créneaux configurés par l'utilisateur -->
          <div class="grid grid-cols-2 sm:grid-cols-4 gap-2">
            <button
              v-for="cr in configStore.creneauxActifs"
              :key="cr.id"
              type="button"
              @click="choisirCreneauPredefini(cr)"
              class="p-2 rounded-[4px] border text-center transition cursor-pointer flex flex-col items-center justify-center"
              :class="
                !isHoraireLibre && selectedCreneauId === cr.id
                  ? 'bg-brand border-border text-white font-bold shadow-xs'
                  : 'bg-muted border-border text-foreground hover:bg-muted'
              "
            >
              <span class="text-xs font-mono font-bold"
                >{{ cr.debut }} – {{ cr.fin }}</span
              >
              <span class="text-[9px] text-muted-foreground mt-0.5">{{
                cr.label
              }}</span>
            </button>

            <!-- Option Horaire libre -->
            <button
              type="button"
              @click="activerHoraireLibre"
              class="p-2 rounded-[4px] border text-center transition cursor-pointer flex flex-col items-center justify-center"
              :class="
                isHoraireLibre
                  ? 'bg-sky-500/15 border-sky-500 text-sky-600 font-bold'
                  : 'bg-muted/40 border-border text-foreground hover:bg-muted'
              "
            >
              <span class="text-xs font-bold">⏰ Horaire libre</span>
              <span class="text-[9px] text-muted-foreground mt-0.5"
                >Personnalisé</span
              >
            </button>
          </div>

          <!-- Champs de saisie si horaire libre choisi -->
          <div
            v-if="isHoraireLibre"
            class="p-3 rounded-[4px] bg-sky-500/5 border border-sky-500/20 flex flex-wrap items-center gap-3"
          >
            <div class="flex items-center gap-2 text-xs">
              <span class="font-bold text-foreground">Heure début :</span>
              <input
                type="time"
                v-model="formHeureDebut"
                class="px-2 py-1 text-xs rounded bg-background border border-border font-bold font-mono outline-none"
              />
            </div>

            <div class="flex items-center gap-2 text-xs">
              <span class="font-bold text-foreground">Heure fin :</span>
              <input
                type="time"
                v-model="formHeureFin"
                class="px-2 py-1 text-xs rounded bg-background border border-border font-bold font-mono outline-none"
              />
            </div>

            <span class="text-[11px] font-semibold text-brand ml-auto">
              Durée calculée : {{ dureeCalculee }}
            </span>
          </div>
        </div>

        <!-- ══════════════════════════════════════════════════════
             4. VÉHICULE AUTO-ÉCOLE ASSIGNÉ
             ══════════════════════════════════════════════════════ -->
        <div class="space-y-1.5 pt-2 border-t border-border">
          <label
            class="block text-xs font-bold text-foreground flex items-center gap-1.5"
          >
            <Car class="w-3.5 h-3.5 text-brand" />
            <span>Véhicule École assigné</span>
          </label>

          <select
            v-model="formVehiculeId"
            class="w-full px-3 py-2 text-xs rounded-[4px] bg-background border border-border focus:border-brand outline-none cursor-pointer"
          >
            <option value="">
              Attribution automatique (premier véhicule disponible)
            </option>
            <option
              v-for="v in guichetStore.vehicules"
              :key="v.id"
              :value="v.id"
            >
              {{ v.marque_modele }} ({{ v.immatriculation }}) · Boîte
              {{
                v.type_boite === "AUTOMATIQUE"
                  ? "Automatique (BVA)"
                  : "Manuelle (BVM)"
              }}
            </option>
          </select>
        </div>

        <!-- 5. Autorisation administrative -->
        <div
          class="p-3 rounded-[4px] bg-muted border border-border flex items-start gap-2.5 text-xs"
        >
          <ShieldCheck class="w-4 h-4 text-emerald-600 shrink-0 mt-0.5" />
          <div class="text-[11px] text-muted-foreground leading-relaxed">
            <span class="font-bold text-foreground"
              >Autorisation Directe de Direction & Guichet :</span
            >
            Cette planification est automatiquement synchronisée avec la feuille
            de route du moniteur.
          </div>
        </div>

        <!-- Boutons d'action -->
        <div
          class="pt-3 border-t border-border flex items-center justify-end space-x-2"
        >
          <button
            type="button"
            @click="emit('close')"
            class="px-4 py-2 rounded-[4px] bg-muted hover:bg-muted/80 text-foreground text-xs font-semibold transition cursor-pointer"
          >
            Annuler
          </button>

          <button
            type="submit"
            :disabled="isSubmitting || selectedCandidats.length === 0"
            class="px-5 py-2 rounded-[4px] bg-brand hover:bg-brand-dark text-white font-bold text-xs uppercase tracking-wider transition disabled:opacity-50 cursor-pointer flex items-center gap-2"
          >
            <span
              v-if="isSubmitting"
              class="w-3.5 h-3.5 rounded-full border-2 border-border border-t-transparent animate-spin"
            />
            <span>{{
              isSubmitting
                ? "Programmation..."
                : selectedCandidats.length > 1
                  ? `Confirmer pour ${selectedCandidats.length} élèves`
                  : "Confirmer la leçon"
            }}</span>
          </button>
        </div>
      </form>
    </div>

    <!-- Modale de configuration des créneaux & jours -->
    <ConfigCreneauxModal
      :isOpen="isConfigModalOpen"
      @close="isConfigModalOpen = false"
    />
  </div>
</template>

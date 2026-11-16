<script setup lang="ts">
import { ref, reactive, watch, computed } from "vue";
import type {
  SchoolItemSummary,
  UpdateSchoolLicencePayload,
} from "../../../types";
import {
  X,
  ShieldCheck,
  Monitor,
  Calendar,
  Loader2,
  Check,
} from "lucide-vue-next";

const props = defineProps<{
  isOpen: boolean;
  school: SchoolItemSummary | null;
  isLoading: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "save", payload: UpdateSchoolLicencePayload): void;
}>();

const form = reactive<UpdateSchoolLicencePayload>({
  school_id: "",
  max_machines: 3,
  type_licence: "ANNUEL",
  duree_valeur: 1,
  duree_unite: "ANNEES",
  mode: "KEEP_DATE",
  date_expiration_licence: undefined,
});

const isCustomDate = ref(false);
const customDateStr = ref("");
const applicationMode = ref<"REPLACE" | "EXTEND" | "KEEP_DATE">("KEEP_DATE");

watch(
  () => props.school,
  (s) => {
    if (s) {
      form.school_id = s.id;
      form.max_machines = s.max_machines || 3;
      isCustomDate.value = false;

      // Si l'école a déjà une échéance active, on commence en mode "Conserver l'échéance" (pratique pour ajuster les postes sans casser la date)
      const hasValidFutureExp =
        s.date_expiration_licence &&
        new Date(s.date_expiration_licence).getTime() > Date.now();
      applicationMode.value = hasValidFutureExp ? "KEEP_DATE" : "REPLACE";

      // Détection automatique de la durée active pour pré-sélectionner le bon preset
      const rawType = (s.type_licence || "").toUpperCase();
      if (
        rawType.includes("TRIAL") ||
        rawType.includes("ESSAI") ||
        rawType.includes("GRATUIT")
      ) {
        form.type_licence = "TRIAL";
        form.duree_valeur = 30;
        form.duree_unite = "JOURS";
      } else if (rawType.includes("MENSUEL") || rawType.includes("30J")) {
        form.type_licence = "MENSUEL";
        form.duree_valeur = 1;
        form.duree_unite = "MOIS";
      } else if (rawType.includes("TRIMESTRIEL") || rawType.includes("90J")) {
        form.type_licence = "TRIMESTRIEL";
        form.duree_valeur = 3;
        form.duree_unite = "MOIS";
      } else if (rawType.includes("SEMESTRIEL") || rawType.includes("180J")) {
        form.type_licence = "SEMESTRIEL";
        form.duree_valeur = 6;
        form.duree_unite = "MOIS";
      } else if (rawType.includes("3 ANS") || rawType.includes("PLURIANNUEL")) {
        form.type_licence = "3 ANS";
        form.duree_valeur = 3;
        form.duree_unite = "ANNEES";
      } else {
        // Formule Annuelle par défaut
        form.type_licence = "ANNUEL";
        form.duree_valeur = 1;
        form.duree_unite = "ANNEES";
      }

      if (s.date_expiration_licence) {
        customDateStr.value = s.date_expiration_licence.split("T")[0];
      }
    }
  },
  { immediate: true },
);

// Presets rapides de durée : applique la nouvelle formule à partir d'aujourd'hui (ou selon le mode choisi)
const applyPreset = (
  val: number,
  unite: "JOURS" | "MOIS" | "ANNEES",
  typeLabel: string,
) => {
  form.duree_valeur = val;
  form.duree_unite = unite;
  form.type_licence = typeLabel;
  isCustomDate.value = false;
  if (applicationMode.value === "KEEP_DATE") {
    applicationMode.value = "REPLACE";
  }
};

// Date projetée calculée avec exactitude selon le mode choisi
const projectedExpirationDate = computed(() => {
  if (isCustomDate.value && customDateStr.value) {
    return new Date(customDateStr.value).toLocaleDateString("fr-FR", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  }
  if (
    applicationMode.value === "KEEP_DATE" &&
    props.school?.date_expiration_licence
  ) {
    return new Date(props.school.date_expiration_licence).toLocaleDateString(
      "fr-FR",
      {
        year: "numeric",
        month: "long",
        day: "numeric",
      },
    );
  }
  let base = new Date();
  if (
    applicationMode.value === "EXTEND" &&
    props.school?.date_expiration_licence
  ) {
    const current = new Date(props.school.date_expiration_licence);
    if (current.getTime() > Date.now()) {
      base = new Date(current.getTime());
    }
  }
  const date = new Date(base);
  if (form.duree_unite === "ANNEES") {
    date.setFullYear(date.getFullYear() + form.duree_valeur);
  } else if (form.duree_unite === "MOIS") {
    date.setMonth(date.getMonth() + form.duree_valeur);
  } else {
    date.setDate(date.getDate() + form.duree_valeur);
  }
  return date.toLocaleDateString("fr-FR", {
    year: "numeric",
    month: "long",
    day: "numeric",
  });
});

const handleSubmit = () => {
  if (form.max_machines < 1) {
    alert("Le nombre de machines doit être au minimum de 1.");
    return;
  }
  const payload: UpdateSchoolLicencePayload = {
    school_id: form.school_id,
    max_machines: form.max_machines,
    type_licence: form.type_licence,
    duree_valeur: form.duree_valeur,
    duree_unite: form.duree_unite,
    mode: isCustomDate.value ? "REPLACE" : applicationMode.value,
    date_expiration_licence:
      isCustomDate.value && customDateStr.value
        ? new Date(customDateStr.value).toISOString()
        : undefined,
  };
  emit("save", payload);
};
</script>

<template>
  <div
    v-if="isOpen && school"
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm"
  >
    <div
      class="w-full max-w-lg bg-card border border-border rounded-[8px] shadow-2xl overflow-hidden animate-in fade-in zoom-in-95 duration-150"
    >
      <!-- En-tête -->
      <div
        class="h-12 border-b border-border bg-muted/50 px-4 flex items-center justify-between"
      >
        <div
          class="flex items-center gap-2 text-foreground font-semibold text-sm"
        >
          <ShieldCheck class="w-4 h-4 text-brand" />
          <span
            >Gestion des Postes & Licence · {{ school.nom_commercial }}</span
          >
        </div>
        <button
          type="button"
          @click="emit('close')"
          class="text-muted-foreground hover:text-foreground p-1 rounded-[4px] transition-colors cursor-pointer"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Corps du formulaire -->
      <form
        @submit.prevent="handleSubmit"
        class="p-5 space-y-4 text-xs font-sans"
      >
        <!-- 1. Quota de Postes / Machines Desktop autorisées -->
        <div
          class="p-3.5 rounded-[6px] bg-muted/40 border border-border space-y-2"
        >
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <Monitor class="w-4 h-4 text-brand" />
              <span class="font-semibold text-foreground"
                >Nombre de Machines Autorisées (Postes Desktop)</span
              >
            </div>
            <span class="text-[11px] font-mono text-muted-foreground">
              Actuellement actives : {{ school.machines_actives || 0 }} /
              {{ school.max_machines || 3 }}
            </span>
          </div>
          <p class="text-[11px] text-muted-foreground leading-relaxed">
            Définit le nombre maximal d'ordinateurs pouvant faire tourner
            l'application desktop connectée à cette auto-école.
          </p>
          <div class="flex flex-wrap items-center gap-2 pt-1">
            <button
              v-for="qty in [1, 2, 3, 5, 10, 20]"
              :key="qty"
              type="button"
              @click="form.max_machines = qty"
              :class="[
                'px-2.5 py-1 rounded-[4px] text-xs font-mono font-bold transition cursor-pointer border',
                form.max_machines === qty
                  ? 'bg-brand-600 text-white border-brand-600 shadow-xs'
                  : 'bg-background border-border text-foreground hover:bg-muted',
              ]"
            >
              {{ qty }} {{ qty === 1 ? "poste" : "postes" }}
            </button>
            <input
              v-model.number="form.max_machines"
              type="number"
              min="1"
              max="200"
              class="w-16 h-7 px-2 rounded-[4px] bg-background border border-input text-brand font-mono font-bold text-center focus:outline-none focus:ring-1 focus:ring-ring focus:border-ring"
            />
          </div>
        </div>

        <!-- 2. Durée d'utilisation du système (Jours / Mois / Années) -->
        <div
          class="p-3.5 rounded-[6px] bg-muted/40 border border-border space-y-3"
        >
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <Calendar class="w-4 h-4 text-brand" />
              <span class="font-semibold text-foreground"
                >Période & Durée d'Utilisation du Système</span
              >
            </div>
            <span class="text-[11px] font-mono text-brand font-medium">
              Formule: {{ form.type_licence }}
            </span>
          </div>

          <!-- Raccourcis de durée avec présélection réactive -->
          <div class="grid grid-cols-3 gap-2">
            <button
              type="button"
              @click="applyPreset(30, 'JOURS', 'TRIAL')"
              :class="[
                'p-2 rounded-[4px] border text-left transition cursor-pointer',
                form.duree_valeur === 30 &&
                form.duree_unite === 'JOURS' &&
                !isCustomDate
                  ? 'bg-brand/10 border-brand text-brand font-semibold shadow-xs'
                  : 'bg-background border-border text-foreground hover:bg-muted',
              ]"
            >
              <div class="font-bold text-xs">30 Jours</div>
              <div class="text-[10px] text-muted-foreground">
                Période d'essai (1 Mois)
              </div>
            </button>

            <button
              type="button"
              @click="applyPreset(1, 'MOIS', 'MENSUEL')"
              :class="[
                'p-2 rounded-[4px] border text-left transition cursor-pointer',
                form.duree_valeur === 1 &&
                form.duree_unite === 'MOIS' &&
                !isCustomDate
                  ? 'bg-brand/10 border-brand text-brand font-semibold shadow-xs'
                  : 'bg-background border-border text-foreground hover:bg-muted',
              ]"
            >
              <div class="font-bold text-xs">1 Mois</div>
              <div class="text-[10px] text-muted-foreground">
                Abonnement 30j
              </div>
            </button>

            <button
              type="button"
              @click="applyPreset(3, 'MOIS', 'TRIMESTRIEL')"
              :class="[
                'p-2 rounded-[4px] border text-left transition cursor-pointer',
                form.duree_valeur === 3 &&
                form.duree_unite === 'MOIS' &&
                !isCustomDate
                  ? 'bg-brand/10 border-brand text-brand font-semibold shadow-xs'
                  : 'bg-background border-border text-foreground hover:bg-muted',
              ]"
            >
              <div class="font-bold text-xs">3 Mois</div>
              <div class="text-[10px] text-muted-foreground">Trimestriel</div>
            </button>

            <button
              type="button"
              @click="applyPreset(6, 'MOIS', 'SEMESTRIEL')"
              :class="[
                'p-2 rounded-[4px] border text-left transition cursor-pointer',
                form.duree_valeur === 6 &&
                form.duree_unite === 'MOIS' &&
                !isCustomDate
                  ? 'bg-brand/10 border-brand text-brand font-semibold shadow-xs'
                  : 'bg-background border-border text-foreground hover:bg-muted',
              ]"
            >
              <div class="font-bold text-xs">6 Mois</div>
              <div class="text-[10px] text-muted-foreground">Semestriel</div>
            </button>

            <button
              type="button"
              @click="applyPreset(1, 'ANNEES', 'ANNUEL')"
              :class="[
                'p-2 rounded-[4px] border text-left transition cursor-pointer',
                form.duree_valeur === 1 &&
                form.duree_unite === 'ANNEES' &&
                !isCustomDate
                  ? 'bg-brand/10 border-brand text-brand font-semibold shadow-xs'
                  : 'bg-background border-border text-foreground hover:bg-muted',
              ]"
            >
              <div class="font-bold text-xs">1 An</div>
              <div class="text-[10px] text-muted-foreground">
                Annuel complet
              </div>
            </button>

            <button
              type="button"
              @click="applyPreset(3, 'ANNEES', '3 ANS')"
              :class="[
                'p-2 rounded-[4px] border text-left transition cursor-pointer',
                form.duree_valeur === 3 &&
                form.duree_unite === 'ANNEES' &&
                !isCustomDate
                  ? 'bg-brand/10 border-brand text-brand font-semibold shadow-xs'
                  : 'bg-background border-border text-foreground hover:bg-muted',
              ]"
            >
              <div class="font-bold text-xs">3 Ans</div>
              <div class="text-[10px] text-muted-foreground">Long terme</div>
            </button>
          </div>

          <!-- Mode d'application de la durée -->
          <div class="pt-2 border-t border-border space-y-2">
            <div class="flex flex-wrap items-center gap-1.5 text-[11px]">
              <button
                type="button"
                @click="
                  applicationMode = 'REPLACE';
                  isCustomDate = false;
                "
                :class="[
                  'px-2 py-1 rounded-[4px] border font-medium transition cursor-pointer',
                  applicationMode === 'REPLACE' && !isCustomDate
                    ? 'bg-brand/15 border-brand text-brand'
                    : 'bg-background border-border text-muted-foreground hover:text-foreground',
                ]"
              >
                Définir dès aujourd'hui
              </button>
              <button
                type="button"
                @click="
                  applicationMode = 'EXTEND';
                  isCustomDate = false;
                "
                :class="[
                  'px-2 py-1 rounded-[4px] border font-medium transition cursor-pointer',
                  applicationMode === 'EXTEND' && !isCustomDate
                    ? 'bg-brand/15 border-brand text-brand'
                    : 'bg-background border-border text-muted-foreground hover:text-foreground',
                ]"
              >
                Cumuler à l'échéance
              </button>
              <button
                v-if="school.date_expiration_licence"
                type="button"
                @click="
                  applicationMode = 'KEEP_DATE';
                  isCustomDate = false;
                "
                :class="[
                  'px-2 py-1 rounded-[4px] border font-medium transition cursor-pointer',
                  applicationMode === 'KEEP_DATE' && !isCustomDate
                    ? 'bg-brand/15 border-brand text-brand'
                    : 'bg-background border-border text-muted-foreground hover:text-foreground',
                ]"
              >
                Postes uniquement (date inchangée)
              </button>
            </div>

            <!-- Option Date Personnalisée -->
            <div class="flex items-center justify-between pt-1">
              <label
                class="flex items-center gap-2 cursor-pointer text-foreground"
              >
                <input
                  v-model="isCustomDate"
                  type="checkbox"
                  class="rounded-[3px] bg-background border-input text-brand focus:ring-0"
                />
                <span class="text-xs"
                  >Définir une date d'expiration exacte</span
                >
              </label>
              <input
                v-if="isCustomDate"
                v-model="customDateStr"
                type="date"
                class="h-7 px-2 rounded-[4px] bg-background border border-input text-foreground text-xs font-mono focus:outline-none focus:ring-1 focus:ring-ring focus:border-ring"
              />
            </div>
          </div>

          <!-- Résumé visuel d'expiration -->
          <div
            class="p-2.5 rounded-[4px] bg-background border border-border flex items-center justify-between text-xs font-mono"
          >
            <span class="text-muted-foreground"
              >Nouvelle date d'échéance :</span
            >
            <span class="font-bold text-brand">{{
              projectedExpirationDate
            }}</span>
          </div>
        </div>

        <!-- Boutons d'action -->
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
            type="submit"
            :disabled="isLoading"
            class="px-4 py-1.5 rounded-[4px] bg-brand-600 hover:bg-brand-700 text-white font-medium flex items-center gap-1.5 transition cursor-pointer disabled:opacity-50 shadow-xs"
          >
            <Loader2
              v-if="isLoading"
              class="w-3.5 h-3.5 animate-spin text-white"
            />
            <Check v-else class="w-3.5 h-3.5" />
            <span>Appliquer la Licence</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

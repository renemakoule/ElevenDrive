<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from "vue";
import { useAuthStore } from "../../stores/authStore";
import { useGuichetStore } from "../../stores/guichetStore";
import { realtimeSocket } from "../../services/realtimeSocketService";
import {
  X,
  Check,
  AlertTriangle,
  Car,
  Calendar,
  Clock,
  UserCheck,
  Radio,
} from "lucide-vue-next";

const authStore = useAuthStore();
const guichetStore = useGuichetStore();
const isWsConnected = realtimeSocket.isConnected;
let unsubscribeWs: (() => void) | null = null;

const props = defineProps<{
  candidatNom: string;
  dossierId: string;
  numeroDossier: string;
  telephone: string;
  soldeRestant: number;
  montantTotal: number;
  heuresEffectuees: number;
  heuresTotal: number;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (
    e: "reserved",
    payload: {
      dossierId: string;
      date: string;
      heureDebut: string;
      heureFin: string;
      vehiculeId?: string;
      vehiculeLabel: string;
    },
  ): void;
}>();

const SEUIL_SOLDE = 0.4;
const NB_JOURS = 6;

const HEURES_DISPO = [
  { label: "08:00 – 10:00", debut: "08:00", fin: "10:00" },
  { label: "10:30 – 12:30", debut: "10:30", fin: "12:30" },
  { label: "14:00 – 16:00", debut: "14:00", fin: "16:00" },
  { label: "16:30 – 18:30", debut: "16:30", fin: "18:30" },
];

onMounted(async () => {
  await guichetStore.loadVehicules();
  if (vehiculesDisponibles.value.length > 0 && !selectedVehiculeId.value) {
    selectedVehiculeId.value = vehiculesDisponibles.value[0].id;
  }
  unsubscribeWs = realtimeSocket.subscribe(async (msg) => {
    if (
      msg.type === "FLOTTE_VEHICULES_UPDATED" ||
      msg.type === "PLANNING_CONDUITE_UPDATED"
    ) {
      await guichetStore.loadVehicules();
    }
  });
});

onUnmounted(() => {
  if (unsubscribeWs) {
    unsubscribeWs();
    unsubscribeWs = null;
  }
});

const vehiculesDisponibles = computed(() => {
  if (guichetStore.vehicules && guichetStore.vehicules.length > 0) {
    return guichetStore.vehicules
      .filter((v) => !v.statut || v.statut === "DISPONIBLE")
      .map((v) => ({
        id: v.id,
        label: `${v.marque_modele} (${v.immatriculation}) · ${v.type_boite === "AUTOMATIQUE" ? "BVA" : "BVM"}`,
      }));
  }
  return [{ id: "", label: "Véhicule École (Attribution automatique)" }];
});

const majuscule = (s: string) => s.charAt(0).toUpperCase() + s.slice(1);

// Prochains jours ouvrés (dimanche exclu), à partir de demain ou aujourd'hui
const jours = (() => {
  const liste: { iso: string; valeur: string; court: string }[] = [];
  const jour = new Date();
  while (liste.length < NB_JOURS) {
    if (jour.getDay() !== 0) {
      const annee = jour.getFullYear();
      const mois = String(jour.getMonth() + 1).padStart(2, "0");
      const d = String(jour.getDate()).padStart(2, "0");
      liste.push({
        iso: `${annee}-${mois}-${d}`,
        valeur: majuscule(
          jour.toLocaleDateString("fr-FR", {
            weekday: "long",
            day: "numeric",
            month: "long",
            year: "numeric",
          }),
        ),
        court: majuscule(
          jour.toLocaleDateString("fr-FR", {
            weekday: "short",
            day: "numeric",
            month: "short",
          }),
        ),
      });
    }
    jour.setDate(jour.getDate() + 1);
  }
  return liste;
})();

const selectedDate = ref(jours[0]);
const selectedHeure = ref(HEURES_DISPO[1]);
const selectedVehiculeId = ref("");

// Synchroniser le premier véhicule dès qu'il est chargé
if (vehiculesDisponibles.value.length > 0) {
  selectedVehiculeId.value = vehiculesDisponibles.value[0].id;
}

const isEligible = computed(
  () =>
    props.soldeRestant <= 0 ||
    props.soldeRestant <= props.montantTotal * SEUIL_SOLDE,
);

const isSubmitting = ref(false);

const handleConfirm = () => {
  if (!isEligible.value || isSubmitting.value) return;
  isSubmitting.value = true;
  const vChoisi = vehiculesDisponibles.value.find(
    (v) => v.id === selectedVehiculeId.value,
  );
  emit("reserved", {
    dossierId: props.dossierId,
    date: selectedDate.value.iso,
    heureDebut: selectedHeure.value.debut,
    heureFin: selectedHeure.value.fin,
    vehiculeId: selectedVehiculeId.value || undefined,
    vehiculeLabel: vChoisi?.label || "Véhicule École",
  });
};

const bouton = (actif: boolean) =>
  actif
    ? "border-brand-dark bg-brand/15 text-brand-dark font-semibold ring-1 ring-brand-dark"
    : "border-border bg-card text-muted-foreground hover:bg-muted hover:text-foreground transition";
</script>

<template>
  <div
    class="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/50 p-0 sm:p-3 select-none"
    @click.self="emit('close')"
  >
    <div
      class="flex max-h-[92dvh] w-full sm:max-w-md flex-col overflow-hidden rounded-t-[3px] sm:rounded-[3px] border-t sm:border border-border bg-card pb-safe"
    >
      <div class="sm:hidden mx-auto mt-2 h-1 w-9 rounded-full bg-border shrink-0" />
      <!-- En-tête -->
      <header
        class="flex items-center justify-between border-b border-border px-4 py-3"
      >
        <div>
          <div
            class="flex items-center gap-2 text-xs font-semibold leading-tight text-foreground"
          >
            <span>Réserver une leçon de conduite</span>
            <span
              class="inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-bold"
              :class="
                isWsConnected
                  ? 'bg-emerald-500/15 text-emerald-700 dark:text-emerald-300'
                  : 'bg-amber-500/15 text-amber-700'
              "
            >
              <Radio class="h-2.5 w-2.5 animate-pulse" />
              {{ isWsConnected ? "Temps réel WS" : "Sync..." }}
            </span>
          </div>
          <div class="text-[11px] text-muted-foreground">
            {{ candidatNom }} · dossier {{ numeroDossier }} · 2 h
          </div>
        </div>
        <button
          @click="emit('close')"
          class="rounded p-1 text-muted-foreground hover:bg-muted hover:text-foreground transition"
          aria-label="Fermer"
        >
          <X class="h-4 w-4" />
        </button>
      </header>

      <!-- Corps -->
      <div class="flex-1 space-y-3 overflow-y-auto p-4 text-xs">
        <!-- Éligibilité -->
        <div
          :class="[
            'flex items-start gap-2 rounded-xl border px-3 py-2 leading-snug',
            isEligible
              ? 'border-brand-border bg-brand/10 text-brand-dark'
              : 'border-amber-500/20 bg-amber-500/10 text-amber-600 dark:text-amber-400',
          ]"
        >
          <Check
            v-if="isEligible"
            class="mt-0.5 h-3.5 w-3.5 shrink-0 text-brand-dark"
          />
          <AlertTriangle
            v-else
            class="mt-0.5 h-3.5 w-3.5 shrink-0 text-amber-500"
          />
          <span v-if="isEligible"
            >Règlement à jour : les créneaux au volant peuvent être
            réservés.</span
          >
          <span v-else>
            Solde restant de {{ soldeRestant.toLocaleString("fr-FR") }} FCFA :
            la 2e tranche (60% réglés) doit être acquittée avant de réserver une
            leçon pratique.
          </span>
        </div>

        <!-- Jour -->
        <div>
          <div
            class="mb-1.5 flex items-center gap-1.5 text-xs font-medium text-muted-foreground"
          >
            <Calendar class="h-3.5 w-3.5 text-brand-dark" />
            <span>Sélectionnez un jour</span>
          </div>
          <div class="grid grid-cols-3 gap-2">
            <button
              v-for="j in jours"
              :key="j.iso"
              type="button"
              @click="selectedDate = j"
              :class="[
                'rounded-lg border px-2 py-1.5 text-center text-xs font-medium transition',
                bouton(selectedDate.iso === j.iso),
              ]"
            >
              {{ j.court }}
            </button>
          </div>
        </div>

        <!-- Créneau -->
        <div>
          <div
            class="mb-1.5 flex items-center gap-1.5 text-xs font-medium text-muted-foreground"
          >
            <Clock class="h-3.5 w-3.5 text-brand-dark" />
            <span>Créneau horaire (2 heures)</span>
          </div>
          <div class="grid grid-cols-2 gap-2 sm:grid-cols-4">
            <button
              v-for="h in HEURES_DISPO"
              :key="h.debut"
              type="button"
              @click="selectedHeure = h"
              :class="[
                'rounded-lg border px-1.5 py-1.5 text-center font-mono text-[11px] transition',
                bouton(selectedHeure.debut === h.debut),
              ]"
            >
              {{ h.label }}
            </button>
          </div>
        </div>

        <!-- Véhicule et moniteur réels -->
        <div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
          <label class="block">
            <span
              class="mb-1 flex items-center gap-1.5 text-xs font-medium text-muted-foreground"
            >
              <Car class="h-3.5 w-3.5 text-brand-dark" />
              Véhicule auto-école
            </span>
            <select
              v-model="selectedVehiculeId"
              class="w-full rounded-[5px] border border-border bg-card px-2.5 py-1.5 text-xs text-foreground focus:border-brand-dark focus:outline-none focus:ring-2 focus:ring-brand-dark/20 transition"
            >
              <option
                v-for="v in vehiculesDisponibles"
                :key="v.id"
                :value="v.id"
              >
                {{ v.label }}
              </option>
            </select>
          </label>
          <label class="block">
            <span
              class="mb-1 flex items-center gap-1.5 text-xs font-medium text-muted-foreground"
            >
              <UserCheck class="h-3.5 w-3.5 text-brand-dark" />
              Moniteur instructeur
            </span>
            <div
              class="w-full rounded-[5px] border border-border bg-muted/40 px-2.5 py-1.5 text-xs text-foreground"
            >
              Moniteur Référent Auto-École
            </div>
          </label>
        </div>

        <!-- Récapitulatif -->
        <div
          class="space-y-1.5 rounded-[5px] border border-border bg-muted/30 px-3 py-2.5"
        >
          <div class="flex justify-between gap-2">
            <span class="text-muted-foreground">Rendez-vous</span>
            <span class="text-right font-medium text-foreground"
              >{{ selectedDate.valeur }}, {{ selectedHeure.label }}</span
            >
          </div>
          <div class="flex justify-between gap-2">
            <span class="text-muted-foreground">Point de départ</span>
            <span class="text-foreground">{{
              authStore.schoolAddress || authStore.schoolName
            }}</span>
          </div>
          <div class="flex justify-between gap-2">
            <span class="text-muted-foreground">Volume de conduite</span>
            <span class="font-mono text-foreground"
              >{{ heuresEffectuees }} / {{ heuresTotal }} h (+2h
              réservées)</span
            >
          </div>
        </div>
      </div>

      <!-- Actions -->
      <footer
        class="flex items-center justify-end gap-2 border-t border-border px-4 py-3"
      >
        <button
          type="button"
          @click="emit('close')"
          class="inline-flex h-10 sm:h-8 items-center rounded-[3px] border border-border bg-card px-3 text-xs font-medium text-foreground hover:bg-muted transition"
        >
          Annuler
        </button>
        <button
          type="button"
          @click="handleConfirm"
          :disabled="!isEligible || isSubmitting"
          class="flex-1 sm:flex-initial inline-flex h-10 sm:h-8 items-center justify-center gap-1.5 rounded-[3px] bg-brand-dark px-3 text-xs font-semibold text-white hover:bg-brand disabled:cursor-not-allowed disabled:opacity-40 transition"
        >
          <Check class="h-3.5 w-3.5" />
          {{ isSubmitting ? "Réservation..." : "Confirmer la réservation" }}
        </button>
      </footer>
    </div>
  </div>
</template>

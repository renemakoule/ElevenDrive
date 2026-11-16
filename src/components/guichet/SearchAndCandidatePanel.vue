<script setup lang="ts">
import { ref, computed } from "vue";
import { useGuichetStore } from "../../stores/guichetStore";
import { useAuthStore } from "../../stores/authStore";
import WebcamCaptureModal from "../webcam/WebcamCaptureModal.vue";
import RegistrationContractModal from "./RegistrationContractModal.vue";
import CandidateProgressSection from "./dossier/CandidateProgressSection.vue";
import CandidatePaymentsHistory from "./dossier/CandidatePaymentsHistory.vue";
import {
  Search,
  User,
  Phone,
  CreditCard,
  CalendarDays,
  HeartPulse,
  Camera,
  FileText,
  Archive,
} from "lucide-vue-next";

const store = useGuichetStore();
const authStore = useAuthStore();
const searchInput = ref("");
const isWebcamOpen = ref(false);
const isContractModalOpen = ref(false);

const detail = computed(() => store.activeDossierDetail);

const encaisse = computed(() =>
  detail.value
    ? detail.value.dossier.montant_total - detail.value.dossier.solde_restant
    : 0,
);
const pctPaye = computed(() => {
  const total = detail.value?.dossier.montant_total ?? 0;
  return total > 0
    ? Math.min(100, Math.round((encaisse.value / total) * 100))
    : 0;
});

const clamp = (n: number | undefined) =>
  Math.max(0, Math.min(100, Math.round(n ?? 0)));
const pctCode = computed(() => clamp(detail.value?.progression_code_pct));
const pctConduite = computed(() =>
  clamp(detail.value?.progression_conduite_pct),
);

const derogationActive = computed(() =>
  detail.value
    ? Boolean(detail.value.dossier.derogation_conduite_accordee) ||
      store.derogationsAccorde.includes(detail.value.dossier.id)
    : false,
);
const conduiteVerrouillee = computed(
  () => detail.value?.dossier.heures_conduite_effectuees === 0,
);

const infos = computed(() => {
  const c = detail.value?.candidat;
  if (!c) return [];
  return [
    { icon: Phone, label: "WhatsApp", value: c.telephone_whatsapp },
    { icon: CreditCard, label: c.type_piece_id, value: c.numero_piece_id },
    { icon: CalendarDays, label: "Né(e) le", value: c.date_naissance },
    {
      icon: HeartPulse,
      label: "Visite médicale",
      value: c.visite_medicale_statut,
    },
  ];
});

const fcfa = (n: number) => n.toLocaleString("fr-FR");

const statutLabel = (s: string) => {
  const t = (s ?? "").replace(/_/g, " ").toLowerCase();
  return t.charAt(0).toUpperCase() + t.slice(1);
};

const onSearchInput = (e: Event) => {
  const target = e.target as HTMLInputElement;
  searchInput.value = target.value;
  store.search(target.value);
};

const openWebcam = () => {
  if (store.activeDossierDetail) {
    isWebcamOpen.value = true;
  }
};

const onPhotoUpdated = (newPhotoUrl: string) => {
  if (store.activeDossierDetail) {
    store.updateCandidatePhoto(
      store.activeDossierDetail.candidat.id,
      newPhotoUrl,
    );
  }
};
</script>

<template>
  <div
    class="flex flex-1 min-h-0 md:h-full flex-col overflow-hidden rounded-[5px] border border-border bg-card"
  >
    <!-- Recherche -->
    <div class="border-b border-border p-4 shrink-0">
      <div class="relative">
        <Search
          class="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground"
        />
        <input
          id="search-input"
          type="text"
          :value="searchInput"
          @input="onSearchInput"
          aria-label="Recherche rapide"
          placeholder="Nom, téléphone, N° CNI ou N° dossier"
          class="h-10 w-full rounded-[4px] border border-border bg-muted pl-9 pr-10 text-xs text-foreground transition placeholder:text-muted-foreground focus:border-brand focus:bg-card focus:outline-none focus:ring-2 focus:ring-brand/20"
        />
        <kbd
          class="absolute right-2.5 top-1/2 -translate-y-1/2 rounded bg-muted px-1.5 py-0.5 font-mono text-[10px] text-muted-foreground"
        >
          F3
        </kbd>
      </div>

      <div
        v-if="store.searchResults.length > 1"
        class="mt-3 flex items-center gap-2 overflow-x-auto pb-1"
      >
        <span class="shrink-0 text-[11px] font-medium text-muted-foreground"
          >{{ store.searchResults.length }} résultats</span
        >
        <button
          v-for="cand in store.searchResults"
          :key="cand.id"
          type="button"
          @click="store.selectCandidat(cand.id)"
          :aria-pressed="store.selectedCandidatId === cand.id"
          class="shrink-0 whitespace-nowrap rounded-full border px-3 py-1 text-[11px] font-medium transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40"
          :class="
            store.selectedCandidatId === cand.id
              ? 'border-brand bg-brand-light text-brand-dark font-semibold'
              : 'border-border bg-card text-foreground/80 hover:bg-muted'
          "
        >
          {{ cand.nom }} {{ cand.prenoms }}
        </button>
      </div>
    </div>

    <!-- Fiche dossier -->
    <div
      v-if="detail"
      class="flex-1 min-h-0 divide-y divide-border overflow-y-auto"
    >
      <!-- 1. Identité -->
      <section class="flex items-start gap-5 p-5">
        <button
          type="button"
          @click="openWebcam"
          class="group relative h-24 w-20 shrink-0 overflow-hidden rounded-[4px] border border-border bg-muted transition hover:border-brand focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40"
          title="Prendre ou mettre à jour la photo 4x4 via webcam"
        >
          <img
            v-if="detail.candidat.photo_url"
            :src="detail.candidat.photo_url"
            alt="Photo d'identité"
            class="h-full w-full object-cover"
          />
          <span
            v-else
            class="flex h-full flex-col items-center justify-center gap-1 text-muted-foreground"
          >
            <User class="h-6 w-6" />
            <span class="text-[10px]">Ajouter</span>
          </span>

          <span
            class="absolute inset-0 flex flex-col items-center justify-center gap-0.5 bg-black/60 text-[10px] font-medium text-white opacity-0 transition group-hover:opacity-100 group-focus-visible:opacity-100"
          >
            <Camera class="h-4 w-4" />
            Changer
          </span>

          <span
            class="absolute bottom-1 right-1 rounded bg-card/90 px-1 text-[9px] font-medium text-foreground"
          >
            {{ detail.candidat.genre }}
          </span>
        </button>

        <div class="min-w-0 flex-1">
          <div
            class="flex flex-wrap items-start justify-between gap-x-3 gap-y-2"
          >
            <div class="min-w-0">
              <h2 class="truncate text-base font-semibold text-foreground">
                {{ detail.candidat.nom }} {{ detail.candidat.prenoms }}
              </h2>
              <p
                class="mt-0.5 flex flex-wrap items-center gap-x-2 text-xs text-muted-foreground"
              >
                <span class="font-mono font-medium text-brand-dark"
                  >N° {{ detail.dossier.numero_dossier }}</span
                >
                <span class="text-border">·</span>
                <span>{{ detail.formule.intitule }}</span>
                <span class="text-border">·</span>
                <span
                  class="inline-flex items-center gap-1 px-1.5 py-0.2 rounded text-[10px] font-bold"
                  :class="
                    detail.dossier.mode_formation_theorique === 'EN_LIGNE'
                      ? 'bg-sky-500/10 text-sky-600 border border-sky-500/20'
                      : detail.dossier.mode_formation_theorique === 'HYBRIDE'
                        ? 'bg-purple-500/10 text-purple-600 border border-purple-500/20'
                        : 'bg-emerald-500/10 text-emerald-600 border border-emerald-500/20'
                  "
                >
                  {{
                    detail.dossier.mode_formation_theorique === "EN_LIGNE"
                      ? "💻 En Ligne"
                      : detail.dossier.mode_formation_theorique === "HYBRIDE"
                        ? "🔄 Hybride"
                        : "🏫 En Salle"
                  }}
                </span>
              </p>
            </div>

            <div class="flex shrink-0 items-center gap-2">
              <span
                class="flex items-center gap-1.5 text-[11px] font-medium text-brand-dark"
              >
                <span class="h-1.5 w-1.5 rounded-full bg-brand" />
                {{ statutLabel(detail.dossier.statut_dossier) }}
              </span>
              <button
                type="button"
                @click="isContractModalOpen = true"
                title="Générer et imprimer le contrat d'apprentissage et la fiche d'inscription"
                class="inline-flex h-8 items-center gap-1.5 rounded-lg border border-border bg-card px-2.5 text-xs font-medium text-foreground transition hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40"
              >
                <FileText class="h-3.5 w-3.5 text-muted-foreground" />
                Contrat A4
              </button>

              <button
                type="button"
                @click="authStore.setActiveSection('archivage')"
                title="Accéder au centre d'archivage pour ce dossier"
                class="inline-flex h-8 items-center gap-1.5 rounded-lg border border-border bg-card px-2.5 text-xs font-medium text-brand hover:bg-brand/10 transition"
              >
                <Archive class="h-3.5 w-3.5" />
                Archivage
              </button>
            </div>
          </div>

          <dl
            class="mt-4 grid grid-cols-1 gap-x-8 gap-y-2.5 text-xs sm:grid-cols-2"
          >
            <div
              v-for="i in infos"
              :key="i.label"
              class="flex items-center justify-between gap-3 min-w-0"
            >
              <dt
                class="flex items-center gap-1.5 text-muted-foreground shrink-0"
              >
                <component :is="i.icon" class="h-3.5 w-3.5" />
                {{ i.label }}
              </dt>
              <dd
                class="truncate font-medium text-foreground"
                :title="String(i.value)"
              >
                {{ i.value }}
              </dd>
            </div>
          </dl>
        </div>
      </section>

      <!-- 2. Situation financière -->
      <section class="p-5">
        <div class="mb-3 flex items-baseline justify-between">
          <h3 class="text-sm font-semibold text-foreground">
            Situation financière
          </h3>
          <span class="text-[11px] text-muted-foreground"
            >{{ pctPaye }} % réglé · FCFA</span
          >
        </div>

        <div class="grid grid-cols-3 divide-x divide-border">
          <div class="pr-3">
            <p class="text-xs text-muted-foreground">Forfait</p>
            <p
              class="mt-0.5 text-sm font-semibold tabular-nums text-foreground"
            >
              {{ fcfa(detail.dossier.montant_total) }}
            </p>
          </div>
          <div class="px-3">
            <p class="text-xs text-muted-foreground">Encaissé</p>
            <p
              class="mt-0.5 text-sm font-semibold tabular-nums text-brand-dark"
            >
              {{ fcfa(encaisse) }}
            </p>
          </div>
          <div class="pl-3">
            <p class="text-xs text-muted-foreground">
              {{
                detail.dossier.solde_restant > 0
                  ? "Reste à payer"
                  : "Solde réglé"
              }}
            </p>
            <p
              class="mt-0.5 text-sm font-semibold tabular-nums"
              :class="
                detail.dossier.solde_restant > 0
                  ? 'text-amber-600'
                  : 'text-brand-dark'
              "
            >
              {{ fcfa(detail.dossier.solde_restant) }}
            </p>
          </div>
        </div>

        <div
          class="mt-3 h-1 overflow-hidden rounded-full bg-muted"
          role="progressbar"
          :aria-valuenow="pctPaye"
          aria-valuemin="0"
          aria-valuemax="100"
          aria-label="Part du forfait réglée"
        >
          <div
            class="h-full rounded-full bg-brand transition-all duration-500"
            :style="{ width: pctPaye + '%' }"
          />
        </div>
      </section>

      <!-- 3. Progression pédagogique -->
      <CandidateProgressSection
        :detail="detail"
        :pct-code="pctCode"
        :pct-conduite="pctConduite"
        :conduite-verrouillee="conduiteVerrouillee"
        :derogation-active="derogationActive"
        @open-override="store.isManagerOverrideModalOpen = true"
      />

      <!-- 4. Historique des encaissements -->
      <CandidatePaymentsHistory
        :paiements="detail.paiements || []"
        @reprint="store.reprintReceipt"
      />
    </div>

    <!-- Aucun candidat sélectionné -->
    <div
      v-else
      class="flex flex-1 flex-col items-center justify-center p-8 text-center"
    >
      <User class="mb-2 h-8 w-8 text-muted-foreground/40" />
      <p class="text-sm font-medium text-foreground">
        Aucun dossier sélectionné
      </p>
      <p class="mt-1 text-xs text-muted-foreground">
        <kbd
          class="rounded bg-muted px-1 font-mono text-[10px] text-muted-foreground"
          >F3</kbd
        >
        pour rechercher un élève ·
        <kbd
          class="rounded bg-muted px-1 font-mono text-[10px] text-muted-foreground"
          >F1</kbd
        >
        pour une nouvelle inscription
      </p>
    </div>

    <!-- Modales -->
    <WebcamCaptureModal
      v-if="isWebcamOpen"
      @close="isWebcamOpen = false"
      @captured="onPhotoUpdated"
    />

    <RegistrationContractModal
      v-if="isContractModalOpen"
      @close="isContractModalOpen = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from "vue";
import { useGuichetStore } from "../../stores/guichetStore";
import WebcamCaptureModal from "../webcam/WebcamCaptureModal.vue";
import CandidateFormuleStep from "./enrollment/CandidateFormuleStep.vue";
import CandidateIdentityStep from "./enrollment/CandidateIdentityStep.vue";
import CandidateContactAcompteStep from "./enrollment/CandidateContactAcompteStep.vue";
import { FileSpreadsheet, X } from "lucide-vue-next";

const store = useGuichetStore();

const nom = ref("");
const prenoms = ref("");
const genre = ref<"M" | "F">("M");
const dateNaissance = ref("2000-01-01");
const lieuNaissance = ref("Abidjan");
const telephoneWhatsapp = ref("+225 ");
const typePiece = ref("CNI");
const numeroPiece = ref("");
const formuleId = ref(store.formules[0]?.id || "");
const acompteInitial = ref(store.formules[0]?.acompte_minimum || 0);
const modeReglement = ref("ESPECES");
const visiteMedicale = ref<
  "VALIDE_SANS_RESERVE" | "VALIDE_AVEC_VERRES" | "EN_ATTENTE"
>("VALIDE_SANS_RESERVE");
const photoUrl = ref<string | null>(null);
const isWebcamOpen = ref(false);
const modeFormationTheorique = ref<"PRESENTIEL" | "EN_LIGNE" | "HYBRIDE">(
  "PRESENTIEL",
);

watch(
  () => store.isNewCandidateModalOpen,
  async (isOpen) => {
    if (isOpen && (!store.formules || store.formules.length === 0)) {
      await store.loadFormules();
    }
    if (store.formules.length > 0 && !formuleId.value) {
      formuleId.value = store.formules[0].id;
      acompteInitial.value = store.formules[0].acompte_minimum;
    }
  },
);

watch(
  () => store.formules,
  (f) => {
    if (f.length > 0 && !formuleId.value) {
      formuleId.value = f[0].id;
      acompteInitial.value = f[0].acompte_minimum;
    }
  },
  { immediate: true },
);

const selectedFormule = computed(
  () =>
    store.formules.find((f) => f.id === formuleId.value) || store.formules[0],
);

const ageCalcule = computed(() => {
  if (!dateNaissance.value) return 0;
  const birth = new Date(dateNaissance.value);
  const now = new Date();
  let age = now.getFullYear() - birth.getFullYear();
  const m = now.getMonth() - birth.getMonth();
  if (m < 0 || (m === 0 && now.getDate() < birth.getDate())) age--;
  return age;
});

const ageMinimum = computed(() =>
  selectedFormule.value?.categorie_permis === "A" ? 16 : 18,
);

const isAgeLegal = computed(() => {
  if (!selectedFormule.value) return true;
  return ageCalcule.value >= ageMinimum.value;
});

const acompteTropBas = computed(
  () =>
    !!selectedFormule.value &&
    acompteInitial.value < selectedFormule.value.acompte_minimum,
);

const soldeApresAcompte = computed(() =>
  selectedFormule.value
    ? Math.max(
        0,
        selectedFormule.value.prix_total - (acompteInitial.value || 0),
      )
    : 0,
);

const fcfa = (n: number) => n.toLocaleString("fr-FR");

const selectFormule = (id: string) => {
  formuleId.value = id;
  if (selectedFormule.value) {
    acompteInitial.value = selectedFormule.value.acompte_minimum;
  }
};

const onPhotoCaptured = (dataUrl: string) => {
  photoUrl.value = dataUrl;
  store.showToast("Photo 4x4 enregistrée.", "success");
};

const close = () => {
  store.isNewCandidateModalOpen = false;
};

const openExcelImport = () => {
  store.isNewCandidateModalOpen = false;
  store.isImportCandidatesModalOpen = true;
};

const submitEnrollment = () => {
  if (!nom.value.trim() || !prenoms.value.trim()) {
    store.showToast("Saisissez le nom et les prénoms de l'élève.", "error");
    return;
  }
  if (!telephoneWhatsapp.value.trim() || telephoneWhatsapp.value.length < 8) {
    store.showToast("Saisissez un numéro WhatsApp valide.", "error");
    return;
  }
  if (!numeroPiece.value.trim()) {
    store.showToast("Saisissez le numéro de la pièce d'identité.", "error");
    return;
  }
  if (!isAgeLegal.value) {
    store.showToast(
      `Âge insuffisant (${ageCalcule.value} ans). Minimum requis : ${ageMinimum.value} ans pour cette formule.`,
      "error",
    );
    return;
  }
  if (selectedFormule.value && acompteTropBas.value) {
    store.showToast(
      `Acompte minimum pour cette formule : ${fcfa(selectedFormule.value.acompte_minimum)} FCFA.`,
      "error",
    );
    return;
  }

  store.createCandidate({
    nom: nom.value.trim().toUpperCase(),
    prenoms: prenoms.value.trim(),
    genre: genre.value,
    date_naissance: dateNaissance.value,
    lieu_naissance: lieuNaissance.value,
    telephone_whatsapp: telephoneWhatsapp.value,
    type_piece_id: typePiece.value,
    numero_piece_id: numeroPiece.value,
    formule_id: formuleId.value,
    acompte_initial: acompteInitial.value,
    mode_reglement: modeReglement.value,
    photo_url: photoUrl.value || undefined,
    visite_medicale_statut: visiteMedicale.value,
    mode_formation_theorique: modeFormationTheorique.value,
  });
};
</script>

<template>
  <div
    v-if="store.isNewCandidateModalOpen"
    class="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/50 p-0 sm:p-4"
    role="dialog"
    aria-modal="true"
    aria-labelledby="enrollment-title"
    @click.self="close"
  >
    <div
      class="flex max-h-[92dvh] w-full sm:max-w-2xl flex-col overflow-hidden rounded-t-[3px] sm:rounded-[3px] border-t sm:border border-border bg-card pb-safe"
    >
      <div
        class="sm:hidden mx-auto mt-2 h-1 w-9 rounded-full bg-border shrink-0"
      />
      <header
        class="flex items-center justify-between gap-3 border-b border-border px-4 py-3"
      >
        <div class="min-w-0">
          <h2
            id="enrollment-title"
            class="flex items-center gap-2 text-sm font-semibold text-foreground"
          >
            Nouvelle inscription
            <kbd
              class="hidden sm:inline rounded bg-muted px-1 font-mono text-[10px] font-normal text-muted-foreground"
              >F1</kbd
            >
          </h2>
          <p class="text-[11px] text-muted-foreground">
            Enrôlement du candidat et création du dossier.
          </p>
        </div>
        <div class="flex items-center gap-2 shrink-0">
          <button
            type="button"
            @click="openExcelImport"
            class="inline-flex h-7 items-center gap-1.5 rounded-[4px] border border-border bg-muted/60 px-2.5 text-[11px] font-medium text-foreground transition hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40"
          >
            <FileSpreadsheet class="h-3.5 w-3.5 text-emerald-600" />
            <span>Importer Excel</span>
          </button>
          <button
            type="button"
            @click="close"
            class="flex h-7 w-7 items-center justify-center rounded-[4px] text-muted-foreground transition hover:bg-muted hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40"
            aria-label="Fermer"
          >
            <X class="h-4 w-4" />
          </button>
        </div>
      </header>

      <form
        @submit.prevent="submitEnrollment"
        class="flex min-h-0 flex-1 flex-col"
      >
        <div class="flex-1 divide-y divide-border overflow-y-auto">
          <CandidateFormuleStep
            :formules="store.formules"
            :formule-id="formuleId"
            v-model:mode-formation-theorique="modeFormationTheorique"
            @select-formule="selectFormule"
          />

          <CandidateIdentityStep
            v-model:nom="nom"
            v-model:prenoms="prenoms"
            v-model:genre="genre"
            v-model:date-naissance="dateNaissance"
            v-model:lieu-naissance="lieuNaissance"
            :photo-url="photoUrl"
            :age-calcule="ageCalcule"
            :age-minimum="ageMinimum"
            :is-age-legal="isAgeLegal"
            :categorie-permis="selectedFormule?.categorie_permis"
            @open-webcam="isWebcamOpen = true"
            @remove-photo="photoUrl = null"
          />

          <CandidateContactAcompteStep
            v-model:telephone-whatsapp="telephoneWhatsapp"
            v-model:type-piece="typePiece"
            v-model:numero-piece="numeroPiece"
            v-model:visite-medicale="visiteMedicale"
            v-model:acompte-initial="acompteInitial"
            v-model:mode-reglement="modeReglement"
            :acompte-trop-bas="acompteTropBas"
            :solde-apres-acompte="soldeApresAcompte"
            :selected-formule="selectedFormule"
          />
        </div>

        <footer
          class="flex items-center justify-between gap-2 border-t border-border px-4 py-3"
        >
          <button
            type="button"
            @click="close"
            class="inline-flex h-8 items-center gap-2 rounded-[4px] border border-border bg-card px-3 text-xs font-medium text-foreground transition hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40"
          >
            Annuler
            <kbd
              class="rounded bg-muted px-1 font-mono text-[10px] text-muted-foreground"
              >Échap</kbd
            >
          </button>

          <button
            type="submit"
            :disabled="
              store.loading || !isAgeLegal || store.formules.length === 0
            "
            class="inline-flex h-8 items-center gap-2 rounded-[4px] bg-brand-dark px-3 text-xs font-medium text-white transition hover:bg-brand focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40 focus-visible:ring-offset-1 disabled:cursor-not-allowed disabled:opacity-50"
          >
            Valider l'inscription
            <kbd
              class="rounded bg-white/20 px-1 font-mono text-[10px] text-white/90"
              >Entrée</kbd
            >
          </button>
        </footer>
      </form>
    </div>

    <WebcamCaptureModal
      v-if="isWebcamOpen"
      @close="isWebcamOpen = false"
      @captured="onPhotoCaptured"
    />
  </div>
</template>

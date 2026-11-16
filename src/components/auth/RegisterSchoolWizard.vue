<script setup lang="ts">
import { ref, reactive, watch } from "vue";
import { ArrowRight, ArrowLeft, Sparkles, CheckCircle2 } from "lucide-vue-next";
import type { RegisterSchoolPayload } from "../../types";
import OAuthProviderButtons from "./OAuthProviderButtons.vue";
import RegisterWizardDirectorStep from "./wizard/RegisterWizardDirectorStep.vue";
import RegisterWizardSummaryStep from "./wizard/RegisterWizardSummaryStep.vue";

type ProviderId = "google" | "facebook" | "linkedin" | "apple";

defineProps<{
  isLoading?: boolean;
  errorMessage?: string | null;
}>();

const emit = defineEmits<{
  (e: "submit", payload: RegisterSchoolPayload): void;
  (e: "oauth", provider: ProviderId): void;
  (e: "clear-error"): void;
  (e: "step-change", step: 1 | 2 | 3 | 4 | 5): void;
}>();

const registerStep = ref<1 | 2 | 3 | 4 | 5>(1);
const registerFormEl = ref<HTMLFormElement | null>(null);
const connectedOAuthProvider = ref<ProviderId | null>(null);

const registerForm = reactive<RegisterSchoolPayload>({
  nom_commercial: "",
  numero_agrement: "",
  pays: "Cameroun",
  ville: "",
  commune: "",
  adresse: "",
  telephone_ecole: "",
  email_ecole: "",
  devise: "FCFA (XAF)",
  directeur_nom: "",
  directeur_prenoms: "",
  directeur_telephone: "",
  directeur_email: "",
  directeur_mot_de_passe: "",
  code_pin_manager: "",
});

watch(
  () => registerForm.pays,
  (newPays) => {
    if (newPays === "Cameroun") {
      registerForm.devise = "FCFA (XAF)";
    } else if (newPays === "Guinée") {
      registerForm.devise = "GNF";
    } else {
      registerForm.devise = "FCFA (XOF)";
    }
  },
  { immediate: true },
);

watch(registerStep, (newStep) => {
  emit("step-change", newStep);
});

const inputBase =
  "h-10 w-full rounded-[5px] border border-border bg-background text-sm text-foreground " +
  "placeholder:text-muted-foreground/60 transition-colors hover:border-foreground/30 " +
  "focus:outline-none focus:border-brand focus:ring-2 focus:ring-brand/20";
const labelClass = "mb-1.5 block text-xs font-medium text-foreground";
const primaryBtn =
  "inline-flex h-10 w-full items-center justify-center gap-2 rounded-[5px] bg-brand text-sm font-semibold text-white " +
  "transition hover:bg-brand-dark focus:outline-none focus-visible:ring-2 focus-visible:ring-brand/40 " +
  "focus-visible:ring-offset-2 focus-visible:ring-offset-background active:scale-[0.99] disabled:opacity-50 disabled:pointer-events-none";
const ghostBtn =
  "inline-flex h-10 items-center justify-center gap-2 rounded-[5px] border border-border bg-transparent px-4 text-sm font-medium text-foreground " +
  "transition hover:bg-muted focus:outline-none focus-visible:ring-2 focus-visible:ring-brand/30";

const nextStep = () => {
  emit("clear-error");
  const fields = registerFormEl.value?.querySelectorAll<
    HTMLInputElement | HTMLSelectElement
  >(`[data-step="${registerStep.value}"]`);
  let valid = true;
  fields?.forEach((field) => {
    if (valid && !field.reportValidity()) valid = false;
  });
  if (valid && registerStep.value < 5) {
    registerStep.value = (registerStep.value + 1) as 1 | 2 | 3 | 4 | 5;
  }
};

const prevStep = () => {
  emit("clear-error");
  if (registerStep.value > 1) {
    registerStep.value = (registerStep.value - 1) as 1 | 2 | 3 | 4 | 5;
  }
};

const goToStep = (step: 1 | 2 | 3 | 4 | 5) => {
  emit("clear-error");
  registerStep.value = step;
};

const handleFormSubmit = () => {
  if (registerStep.value < 5) {
    nextStep();
  } else {
    emit("submit", { ...registerForm });
  }
};

const applyOAuthProfile = (info: {
  provider: ProviderId;
  email: string;
  nom: string;
  prenoms: string;
}) => {
  connectedOAuthProvider.value = info.provider;
  if (info.email) registerForm.directeur_email = info.email;
  if (info.nom) registerForm.directeur_nom = info.nom;
  if (info.prenoms) registerForm.directeur_prenoms = info.prenoms;
  registerStep.value = 1;
};

defineExpose({
  resetStep: () => {
    registerStep.value = 1;
    connectedOAuthProvider.value = null;
  },
  applyOAuthProfile,
  goToStep,
  currentStep: registerStep,
});
</script>

<template>
  <form ref="registerFormEl" @submit.prevent="handleFormSubmit" class="mt-8">
    <!-- Badge OAuth -->
    <div
      v-if="connectedOAuthProvider"
      class="mb-5 flex items-center justify-between rounded-[5px] border border-brand-border bg-brand-light/70 p-3 text-xs text-brand-dark animate-in fade-in duration-200"
    >
      <div class="flex items-center gap-2.5">
        <CheckCircle2 class="h-4 w-4 text-brand shrink-0" />
        <div>
          <span class="font-semibold text-xs text-brand-dark">
            Identité {{ connectedOAuthProvider.toUpperCase() }} vérifiée
          </span>
          <p class="text-[11px] text-muted-foreground">
            {{ registerForm.directeur_email || "Profil social rattaché" }} · Vos
            coordonnées sont pré-remplies.
          </p>
        </div>
      </div>
      <button
        type="button"
        @click="connectedOAuthProvider = null"
        class="text-[10px] font-medium text-muted-foreground hover:text-foreground underline underline-offset-2"
      >
        Détacher
      </button>
    </div>

    <!-- Progression -->
    <div class="mb-6 flex items-center gap-3">
      <div class="flex flex-1 gap-1.5">
        <span
          v-for="s in 5"
          :key="s"
          class="h-1 flex-1 rounded-full transition-all duration-300"
          :class="s <= registerStep ? 'bg-brand' : 'bg-border'"
        ></span>
      </div>
      <span class="text-xs font-medium text-muted-foreground tabular-nums"
        >Étape {{ registerStep }} sur 5</span
      >
    </div>

    <!-- ÉTAPE 1 : Établissement / OAuth -->
    <div v-show="registerStep === 1" class="space-y-4">
      <template v-if="!connectedOAuthProvider">
        <div>
          <span class="mb-2 block text-xs font-medium text-foreground">
            Créer rapidement avec un compte :
          </span>
          <OAuthProviderButtons
            :disabled="isLoading"
            :show-last-badge="false"
            @select="emit('oauth', $event)"
          />
        </div>

        <div class="my-4 flex items-center gap-3">
          <span class="h-px flex-1 bg-border"></span>
          <span
            class="text-[11px] font-medium text-muted-foreground uppercase tracking-wider"
          >
            Ou inscription classique
          </span>
          <span class="h-px flex-1 bg-border"></span>
        </div>
      </template>

      <div
        v-else
        class="rounded-[5px] border border-brand-border/80 bg-brand-light/50 p-3 text-xs text-brand-dark flex items-start gap-2.5 mb-2"
      >
        <Sparkles class="h-4 w-4 text-brand shrink-0 mt-0.5" />
        <div>
          <span class="font-semibold text-xs"
            >Compte {{ connectedOAuthProvider.toUpperCase() }} connecté</span
          >
          <p class="text-[11px] text-muted-foreground mt-0.5">
            Renseignez maintenant le nom et l'agrément de votre auto-école pour
            finaliser la création de votre établissement.
          </p>
        </div>
      </div>

      <div>
        <label for="reg-nom" :class="labelClass"
          >Nom commercial de l'établissement</label
        >
        <input
          id="reg-nom"
          v-model="registerForm.nom_commercial"
          @input="emit('clear-error')"
          data-step="1"
          type="text"
          required
          placeholder="Ex : Auto-École de l'Avenir"
          :class="[inputBase, 'px-3']"
        />
      </div>

      <div>
        <label for="reg-agrement" :class="labelClass"
          >Numéro d'agrément ministériel</label
        >
        <input
          id="reg-agrement"
          v-model="registerForm.numero_agrement"
          @input="emit('clear-error')"
          data-step="1"
          type="text"
          required
          placeholder="Ex : AGR-MIN-TRANS-2026/045"
          :class="[inputBase, 'px-3 font-mono uppercase']"
        />
      </div>

      <button type="button" @click="nextStep" :class="[primaryBtn, 'mt-4']">
        <span>Continuer</span>
        <ArrowRight class="h-4 w-4" />
      </button>
    </div>

    <!-- ÉTAPE 2 : Localisation & Siège -->
    <div v-show="registerStep === 2" class="space-y-4">
      <div>
        <label for="reg-pays" :class="labelClass">Pays et Devise</label>
        <select
          id="reg-pays"
          v-model="registerForm.pays"
          data-step="2"
          :class="[inputBase, 'px-3']"
        >
          <option value="Côte d'Ivoire">Côte d'Ivoire (XOF)</option>
          <option value="Sénégal">Sénégal (XOF)</option>
          <option value="Cameroun">Cameroun (XAF)</option>
          <option value="Bénin">Bénin (XOF)</option>
          <option value="Togo">Togo (XOF)</option>
          <option value="Burkina Faso">Burkina Faso (XOF)</option>
          <option value="Mali">Mali (XOF)</option>
          <option value="Guinée">Guinée (GNF)</option>
        </select>
      </div>

      <div class="grid grid-cols-2 gap-3">
        <div>
          <label for="reg-ville" :class="labelClass">Ville</label>
          <input
            id="reg-ville"
            v-model="registerForm.ville"
            @input="emit('clear-error')"
            data-step="2"
            type="text"
            required
            placeholder="Ex : Douala"
            :class="[inputBase, 'px-3']"
          />
        </div>
        <div>
          <label for="reg-commune" :class="labelClass">Commune du siège</label>
          <input
            id="reg-commune"
            v-model="registerForm.commune"
            @input="emit('clear-error')"
            data-step="2"
            type="text"
            required
            placeholder="Ex : Makepe"
            :class="[inputBase, 'px-3']"
          />
        </div>
      </div>

      <div>
        <label for="reg-adresse" :class="labelClass">
          Adresse physique / Rue
          <span class="font-normal text-muted-foreground">(facultatif)</span>
        </label>
        <input
          id="reg-adresse"
          v-model="registerForm.adresse"
          data-step="2"
          type="text"
          placeholder="Ex : Rond-poulin, Boulangerie Saker"
          :class="[inputBase, 'px-3']"
        />
      </div>

      <div class="flex gap-3 pt-2">
        <button
          type="button"
          @click="prevStep"
          :class="[ghostBtn, 'shrink-0']"
          aria-label="Retour à l'étape 1"
        >
          <ArrowLeft class="h-4 w-4" />
        </button>
        <button type="button" @click="nextStep" :class="primaryBtn">
          <span>Continuer vers le directeur</span>
          <ArrowRight class="h-4 w-4" />
        </button>
      </div>
    </div>

    <!-- ÉTAPES 3 & 4 -->
    <RegisterWizardDirectorStep
      :step="registerStep"
      :form="registerForm"
      :connected-o-auth-provider="connectedOAuthProvider"
      :input-base="inputBase"
      :label-class="labelClass"
      :primary-btn="primaryBtn"
      :ghost-btn="ghostBtn"
      @next="nextStep"
      @prev="prevStep"
      @clear-error="emit('clear-error')"
    />

    <!-- ÉTAPE 5 -->
    <RegisterWizardSummaryStep
      v-show="registerStep === 5"
      :form="registerForm"
      :connected-o-auth-provider="connectedOAuthProvider"
      :is-loading="isLoading"
      :error-message="errorMessage"
      :primary-btn="primaryBtn"
      :ghost-btn="ghostBtn"
      @go-to-step="goToStep"
      @prev="prevStep"
    />
  </form>
</template>

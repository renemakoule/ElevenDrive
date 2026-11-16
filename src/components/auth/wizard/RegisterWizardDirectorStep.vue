<script setup lang="ts">
import { ref } from "vue";
import {
  ArrowLeft,
  ArrowRight,
  CheckCircle2,
  Eye,
  EyeOff,
} from "lucide-vue-next";
import type { RegisterSchoolPayload } from "../../../types";

defineProps<{
  step: 1 | 2 | 3 | 4 | 5;
  form: RegisterSchoolPayload;
  connectedOAuthProvider: string | null;
  inputBase: string;
  labelClass: string;
  primaryBtn: string;
  ghostBtn: string;
}>();

const emit = defineEmits<{
  (e: "next"): void;
  (e: "prev"): void;
  (e: "clear-error"): void;
}>();

const showDirectorPassword = ref(false);
</script>

<template>
  <div>
    <!-- ÉTAPE 3 : Directeur / Administrateur -->
    <div v-show="step === 3" class="space-y-4">
      <div class="grid grid-cols-2 gap-3">
        <div>
          <label for="reg-dir-nom" :class="labelClass">Nom de famille</label>
          <input
            id="reg-dir-nom"
            v-model="form.directeur_nom"
            @input="emit('clear-error')"
            data-step="3"
            type="text"
            required
            placeholder="Ex : MAKOULE"
            :class="[inputBase, 'px-3 uppercase']"
          />
        </div>
        <div>
          <label for="reg-dir-prenoms" :class="labelClass">Prénoms</label>
          <input
            id="reg-dir-prenoms"
            v-model="form.directeur_prenoms"
            @input="emit('clear-error')"
            data-step="3"
            type="text"
            required
            placeholder="Ex : Rene Boris"
            :class="[inputBase, 'px-3']"
          />
        </div>
      </div>

      <div>
        <label for="reg-dir-tel" :class="labelClass"
          >Numéro WhatsApp du directeur</label
        >
        <input
          id="reg-dir-tel"
          v-model="form.directeur_telephone"
          @input="emit('clear-error')"
          data-step="3"
          type="tel"
          required
          placeholder="Ex : +237 651 727 932"
          :class="[inputBase, 'px-3']"
        />
        <p class="mt-1 text-[11px] text-muted-foreground">
          Sera utilisé comme identifiant de connexion guichet.
        </p>
      </div>

      <div>
        <div class="flex items-center justify-between mb-1.5">
          <label
            for="reg-dir-email"
            class="block text-xs font-medium text-foreground"
          >
            Adresse email
            <span class="font-normal text-muted-foreground">(facultatif)</span>
          </label>
          <span
            v-if="connectedOAuthProvider"
            class="inline-flex items-center gap-1 rounded bg-brand-light px-1.5 py-0.5 text-[10px] font-medium text-brand-dark"
          >
            <CheckCircle2 class="h-3 w-3 text-brand" />
            Vérifié via {{ connectedOAuthProvider.toUpperCase() }}
          </span>
        </div>
        <input
          id="reg-dir-email"
          v-model="form.directeur_email"
          @input="emit('clear-error')"
          data-step="3"
          type="email"
          placeholder="Ex : direction@monautoecole.ci"
          :class="[inputBase, 'px-3']"
        />
      </div>

      <div class="flex gap-3 pt-2">
        <button
          type="button"
          @click="emit('prev')"
          :class="[ghostBtn, 'shrink-0']"
          aria-label="Retour à l'étape 2"
        >
          <ArrowLeft class="h-4 w-4" />
        </button>
        <button type="button" @click="emit('next')" :class="primaryBtn">
          <span>Continuer vers la sécurité</span>
          <ArrowRight class="h-4 w-4" />
        </button>
      </div>
    </div>

    <!-- ÉTAPE 4 : Sécurité & Code PIN Managérial -->
    <div v-show="step === 4" class="space-y-4">
      <div
        class="rounded-[5px] border border-brand-border/60 bg-brand-light/40 p-3 text-xs text-foreground"
      >
        <div class="flex items-center justify-between">
          <span class="font-semibold text-brand-dark">{{
            form.nom_commercial || "Votre établissement"
          }}</span>
          <span class="text-[11px] text-muted-foreground"
            >{{ form.ville || "Siège" }}, {{ form.pays }}</span
          >
        </div>
        <div class="mt-1 text-muted-foreground">
          Directeur :
          <span class="font-medium text-foreground"
            >{{ form.directeur_prenoms }} {{ form.directeur_nom }}</span
          >
          <span v-if="form.directeur_telephone" class="ml-1 font-mono"
            >({{ form.directeur_telephone }})</span
          >
        </div>
      </div>

      <div>
        <label for="reg-dir-password" :class="labelClass"
          >Mot de passe (6 caractères minimum)</label
        >
        <div class="relative">
          <input
            id="reg-dir-password"
            v-model="form.directeur_mot_de_passe"
            @input="emit('clear-error')"
            data-step="4"
            :type="showDirectorPassword ? 'text' : 'password'"
            required
            minlength="6"
            autocomplete="new-password"
            placeholder="••••••••"
            :class="[inputBase, 'pl-3 pr-10']"
          />
          <button
            type="button"
            @click="showDirectorPassword = !showDirectorPassword"
            :aria-label="
              showDirectorPassword
                ? 'Masquer le mot de passe'
                : 'Afficher le mot de passe'
            "
            class="absolute right-2 top-2 rounded p-1 text-muted-foreground transition hover:text-foreground focus:outline-none focus-visible:ring-2 focus-visible:ring-brand/30"
          >
            <EyeOff v-if="showDirectorPassword" class="h-4 w-4" />
            <Eye v-else class="h-4 w-4" />
          </button>
        </div>
      </div>

      <div>
        <label for="reg-pin" :class="labelClass"
          >Code PIN managérial (4 chiffres)</label
        >
        <input
          id="reg-pin"
          v-model="form.code_pin_manager"
          @input="emit('clear-error')"
          data-step="4"
          type="password"
          inputmode="numeric"
          maxlength="4"
          required
          placeholder="••••"
          :class="[inputBase, 'px-3 text-center font-mono tracking-[0.5em] text-lg']"
        />
        <p class="mt-1.5 text-xs text-muted-foreground">
          Code secret autorisant les remises spéciales, annulations et
          dérogations au guichet.
        </p>
      </div>

      <div class="flex gap-3 pt-2">
        <button
          type="button"
          @click="emit('prev')"
          :class="[ghostBtn, 'shrink-0']"
          aria-label="Retour à l'étape 3"
        >
          <ArrowLeft class="h-4 w-4" />
        </button>
        <button type="button" @click="emit('next')" :class="primaryBtn">
          <span>Continuer vers le récapitulatif</span>
          <ArrowRight class="h-4 w-4" />
        </button>
      </div>
    </div>
  </div>
</template>


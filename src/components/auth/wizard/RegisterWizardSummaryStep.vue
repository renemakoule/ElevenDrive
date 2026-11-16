<script setup lang="ts">
import {
  ArrowLeft,
  Sparkles,
  Building2,
  MapPin,
  User,
  ShieldCheck,
  Edit3,
  AlertTriangle,
} from "lucide-vue-next";
import type { RegisterSchoolPayload } from "../../../types";

defineProps<{
  form: RegisterSchoolPayload;
  connectedOAuthProvider: string | null;
  isLoading?: boolean;
  errorMessage?: string | null;
  primaryBtn: string;
  ghostBtn: string;
}>();

const emit = defineEmits<{
  (e: "go-to-step", step: 1 | 2 | 3 | 4 | 5): void;
  (e: "prev"): void;
}>();
</script>

<template>
  <div class="space-y-4">
    <div
      class="rounded-[5px] border border-brand-border/70 bg-brand-light/30 p-3.5 text-xs text-foreground space-y-3"
    >
      <!-- 1. Établissement -->
      <div
        class="flex items-start justify-between border-b border-border/60 pb-3"
      >
        <div class="flex items-start gap-2.5">
          <div
            class="flex h-7 w-7 items-center justify-center rounded-lg bg-brand/10 text-brand shrink-0 mt-0.5"
          >
            <Building2 class="h-4 w-4" />
          </div>
          <div>
            <span
              class="text-[10px] font-semibold text-muted-foreground uppercase tracking-wider block"
              >Établissement</span
            >
            <span class="font-bold text-sm text-foreground block">{{
              form.nom_commercial || "—"
            }}</span>
            <span class="font-mono text-xs text-muted-foreground mt-0.5 block"
              >Agrément : {{ form.numero_agrement || "—" }}</span
            >
          </div>
        </div>
        <button
          type="button"
          @click="emit('go-to-step', 1)"
          class="inline-flex items-center gap-1 rounded px-2 py-1 text-[11px] font-medium text-brand hover:bg-brand-light transition"
        >
          <Edit3 class="h-3 w-3" />
          <span>Modifier</span>
        </button>
      </div>

      <!-- 2. Localisation & Siège -->
      <div
        class="flex items-start justify-between border-b border-border/60 pb-3"
      >
        <div class="flex items-start gap-2.5">
          <div
            class="flex h-7 w-7 items-center justify-center rounded-lg bg-brand/10 text-brand shrink-0 mt-0.5"
          >
            <MapPin class="h-4 w-4" />
          </div>
          <div>
            <span
              class="text-[10px] font-semibold text-muted-foreground uppercase tracking-wider block"
              >Siège & Devise</span
            >
            <span class="font-medium text-xs text-foreground block">
              {{ form.commune ? form.commune + ", " : "" }}{{ form.ville }} ({{
                form.pays
              }})
            </span>
            <span
              v-if="form.adresse"
              class="text-[11px] text-muted-foreground block mt-0.5"
            >
              {{ form.adresse }}
            </span>
            <span
              class="inline-block mt-1 rounded bg-muted px-2 py-0.5 text-[10px] font-mono font-medium text-muted-foreground"
            >
              Devise : {{ form.devise }}
            </span>
          </div>
        </div>
        <button
          type="button"
          @click="emit('go-to-step', 2)"
          class="inline-flex items-center gap-1 rounded px-2 py-1 text-[11px] font-medium text-brand hover:bg-brand-light transition"
        >
          <Edit3 class="h-3 w-3" />
          <span>Modifier</span>
        </button>
      </div>

      <!-- 3. Directeur / Administrateur -->
      <div
        class="flex items-start justify-between border-b border-border/60 pb-3"
      >
        <div class="flex items-start gap-2.5">
          <div
            class="flex h-7 w-7 items-center justify-center rounded-lg bg-brand/10 text-brand shrink-0 mt-0.5"
          >
            <User class="h-4 w-4" />
          </div>
          <div>
            <span
              class="text-[10px] font-semibold text-muted-foreground uppercase tracking-wider block"
              >Directeur / Gérant</span
            >
            <span class="font-medium text-xs text-foreground block">
              {{ form.directeur_prenoms }} {{ form.directeur_nom }}
            </span>
            <span class="font-mono text-xs text-muted-foreground block mt-0.5">
              WhatsApp : {{ form.directeur_telephone || "—" }}
            </span>
            <div
              v-if="form.directeur_email"
              class="flex items-center gap-1.5 mt-0.5"
            >
              <span class="text-[11px] text-muted-foreground">{{
                form.directeur_email
              }}</span>
              <span
                v-if="connectedOAuthProvider"
                class="rounded bg-brand-light border border-brand-border px-1.5 py-0.2 text-[9px] font-semibold text-brand-dark"
              >
                {{ connectedOAuthProvider.toUpperCase() }}
              </span>
            </div>
          </div>
        </div>
        <button
          type="button"
          @click="emit('go-to-step', 3)"
          class="inline-flex items-center gap-1 rounded px-2 py-1 text-[11px] font-medium text-brand hover:bg-brand-light transition"
        >
          <Edit3 class="h-3 w-3" />
          <span>Modifier</span>
        </button>
      </div>

      <!-- 4. Sécurité & PIN -->
      <div class="flex items-start justify-between">
        <div class="flex items-start gap-2.5">
          <div
            class="flex h-7 w-7 items-center justify-center rounded-lg bg-brand/10 text-brand shrink-0 mt-0.5"
          >
            <ShieldCheck class="h-4 w-4" />
          </div>
          <div>
            <span
              class="text-[10px] font-semibold text-muted-foreground uppercase tracking-wider block"
              >Sécurité du compte</span
            >
            <span class="text-xs text-muted-foreground block">
              Mot de passe :
              <span
                class="font-mono tracking-widest text-foreground font-semibold"
                >••••••••</span
              >
            </span>
            <span class="text-xs text-muted-foreground block mt-0.5">
              Code PIN Manager :
              <span class="font-mono tracking-widest text-brand font-semibold"
                >•••• (4 chiffres)</span
              >
            </span>
          </div>
        </div>
        <button
          type="button"
          @click="emit('go-to-step', 4)"
          class="inline-flex items-center gap-1 rounded px-2 py-1 text-[11px] font-medium text-brand hover:bg-brand-light transition"
        >
          <Edit3 class="h-3 w-3" />
          <span>Modifier</span>
        </button>
      </div>
    </div>

    <p class="text-center text-[11px] text-muted-foreground">
      Vérifiez vos données avant confirmation. Vous pourrez modifier ces
      paramètres depuis votre espace Direction.
    </p>

    <div
      v-if="errorMessage"
      role="alert"
      class="flex items-start gap-2.5 rounded-[5px] border border-rose-500/20 bg-rose-500/10 p-3 text-xs text-rose-600 dark:text-rose-400 animate-in fade-in duration-150"
    >
      <AlertTriangle class="h-4 w-4 shrink-0 mt-0.5" />
      <span class="leading-relaxed">{{ errorMessage }}</span>
    </div>

    <div class="flex gap-3 pt-2">
      <button
        type="button"
        @click="emit('prev')"
        :class="[ghostBtn, 'shrink-0']"
        aria-label="Retour à l'étape 4"
      >
        <ArrowLeft class="h-4 w-4" />
      </button>
      <button type="submit" :disabled="isLoading" :class="primaryBtn">
        <span v-if="isLoading">Création de l'établissement...</span>
        <template v-else>
          <span>Confirmer et créer l'auto-école</span>
          <Sparkles class="h-4 w-4" />
        </template>
      </button>
    </div>
  </div>
</template>

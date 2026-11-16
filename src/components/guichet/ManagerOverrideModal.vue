<script setup lang="ts">
import { ref, nextTick } from "vue";
import { useGuichetStore } from "../../stores/guichetStore";
import { X, ShieldAlert } from "lucide-vue-next";

const store = useGuichetStore();

const pin = ref("");
const motif = ref("");
const errorMessage = ref("");
const pinInput = ref<HTMLInputElement | null>(null);

const close = () => {
  store.isManagerOverrideModalOpen = false;
  pin.value = "";
  motif.value = "";
  errorMessage.value = "";
};

const submitOverride = async () => {
  if (!store.activeDossierDetail) return;

  if (!pin.value.trim()) {
    errorMessage.value = "Saisissez le code PIN du manager.";
    return;
  }

  const success = await store.unlockConduiteWithPin(
    pin.value,
    store.activeDossierDetail.dossier.id,
    motif.value.trim(),
  );
  if (success) {
    pin.value = "";
    motif.value = "";
    errorMessage.value = "";
    return;
  }

  errorMessage.value =
    "Code PIN incorrect ou droits de dérogation insuffisants.";
  pin.value = "";
  await nextTick();
  pinInput.value?.focus();
};
</script>

<template>
  <div
    v-if="store.isManagerOverrideModalOpen && store.activeDossierDetail"
    class="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/50 p-0 sm:p-4"
    role="dialog"
    aria-modal="true"
    aria-labelledby="override-title"
    @click.self="close"
  >
    <div
      class="w-full sm:max-w-sm overflow-hidden rounded-t-[3px] sm:rounded-[3px] border-t sm:border border-border bg-card pb-safe"
    >
      <div
        class="sm:hidden mx-auto mt-2 h-1 w-9 rounded-full bg-border shrink-0"
      />
      <header
        class="flex items-center justify-between gap-3 border-b border-border px-4 py-3"
      >
        <div class="min-w-0">
          <h2 id="override-title" class="text-sm font-semibold text-foreground">
            Dérogation manager
          </h2>
          <p class="text-[11px] text-muted-foreground">
            Autorisation spéciale · code PIN requis
          </p>
        </div>
        <button
          type="button"
          @click="close"
          class="flex h-7 w-7 items-center justify-center rounded-[4px] text-muted-foreground transition hover:bg-muted hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-dark/40"
          aria-label="Fermer"
        >
          <X class="h-4 w-4" />
        </button>
      </header>

      <form @submit.prevent="submitOverride" class="space-y-4 px-4 py-4">
        <div>
          <p class="text-[11px] text-muted-foreground">Dossier concerné</p>
          <p class="text-xs font-semibold text-foreground">
            {{ store.activeDossierDetail.candidat.nom }}
            {{ store.activeDossierDetail.candidat.prenoms }}
            <span class="ml-1 font-mono font-medium text-brand-dark">
              {{ store.activeDossierDetail.dossier.numero_dossier }}
            </span>
          </p>
          <p
            class="mt-2 flex items-start gap-1.5 rounded-[4px] bg-amber-500/10 px-3 py-2 text-[11px] leading-snug text-amber-600 dark:text-amber-400"
          >
            <ShieldAlert
              class="mt-0.5 h-3.5 w-3.5 shrink-0 text-amber-600 dark:text-amber-400"
            />
            Lève le verrou : accès immédiat à la conduite avant d'atteindre 60 %
            du forfait réglé.
          </p>
        </div>

        <div>
          <label
            for="override-pin"
            class="mb-1 block text-xs text-muted-foreground"
          >
            Code PIN directeur ou superviseur
            <span class="text-destructive">*</span>
          </label>
          <input
            id="override-pin"
            ref="pinInput"
            type="password"
            inputmode="numeric"
            autocomplete="off"
            maxlength="6"
            v-model="pin"
            autofocus
            :aria-invalid="!!errorMessage"
            aria-describedby="override-error"
            class="h-10 w-full rounded-[4px] border bg-card text-center font-mono text-lg font-semibold tracking-[0.4em] text-foreground transition focus:outline-none focus:ring-2"
            :class="
              errorMessage
                ? 'border-destructive focus:border-destructive focus:ring-destructive/20'
                : 'border-border focus:border-brand-dark focus:ring-brand-dark/20'
            "
          />
          <p
            v-if="errorMessage"
            id="override-error"
            role="alert"
            class="mt-1.5 text-[11px] font-medium text-destructive"
          >
            {{ errorMessage }}
          </p>
        </div>

        <div>
          <label
            for="override-motif"
            class="mb-1 block text-xs text-muted-foreground"
            >Motif de la dérogation</label
          >
          <input
            id="override-motif"
            type="text"
            v-model="motif"
            placeholder="Ex: Échelonnement autorisé, avance tuteur reçue..."
            class="h-8 w-full rounded-[4px] border border-border bg-card px-2.5 text-xs text-foreground transition focus:border-brand-dark focus:outline-none focus:ring-2 focus:ring-brand-dark/20"
          />
        </div>

        <div
          class="flex items-center justify-between gap-2 border-t border-border pt-3"
        >
          <button
            type="button"
            @click="close"
            class="inline-flex h-8 items-center gap-2 rounded-[4px] border border-border bg-card px-3 text-xs font-medium text-foreground transition hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-dark/40"
          >
            Annuler
            <kbd
              class="rounded bg-muted px-1 font-mono text-[10px] text-muted-foreground"
              >Échap</kbd
            >
          </button>

          <button
            type="submit"
            class="inline-flex h-8 items-center gap-2 rounded-[4px] bg-brand-dark px-3 text-xs font-medium text-white transition hover:bg-brand focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-dark/40 focus-visible:ring-offset-1"
          >
            Valider la dérogation
            <kbd
              class="rounded bg-white/20 px-1 font-mono text-[10px] text-white/90"
              >Entrée</kbd
            >
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

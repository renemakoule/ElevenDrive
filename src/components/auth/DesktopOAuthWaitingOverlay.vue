<script setup lang="ts">
import { Globe, Loader2 } from "lucide-vue-next";
import { providerNames, type ProviderId } from "./useDesktopOAuthBridge";

defineProps<{
  isWaiting: boolean;
  provider: ProviderId;
}>();

const emit = defineEmits<{
  (e: "cancel"): void;
}>();
</script>

<template>
  <div
    v-if="isWaiting"
    class="absolute inset-0 z-50 bg-background/95 backdrop-blur-md flex flex-col items-center justify-center p-6 text-center animate-in fade-in duration-200"
  >
    <div class="relative mb-6 flex items-center justify-center">
      <div
        class="h-20 w-20 rounded-2xl bg-brand/10 border-2 border-brand/30 flex items-center justify-center shadow-lg shadow-brand/10"
      >
        <Globe class="h-10 w-10 text-brand animate-pulse" />
      </div>
      <div
        class="absolute -bottom-1 -right-1 h-7 w-7 rounded-full bg-card border border-border flex items-center justify-center shadow-md"
      >
        <Loader2 class="h-4 w-4 text-brand animate-spin" />
      </div>
    </div>

    <h3 class="text-lg font-bold text-foreground mb-2">
      Validation {{ providerNames[provider] }} dans votre navigateur
    </h3>

    <p class="text-xs text-muted-foreground max-w-sm leading-relaxed mb-6">
      Votre navigateur par défaut s'est ouvert pour vous connecter en toute
      sécurité avec votre compte
      <strong class="text-foreground font-semibold">{{
        providerNames[provider]
      }}</strong
      >. <br /><br />
      Dès que vous avez validé sur le navigateur, vous serez connecté
      automatiquement ici sans aucune saisie.
    </p>

    <button
      type="button"
      @click="emit('cancel')"
      class="px-5 py-2 rounded-[4px] border border-border bg-card hover:bg-muted text-xs font-semibold text-foreground transition active:scale-95 shadow-xs"
    >
      Annuler la connexion
    </button>
  </div>
</template>

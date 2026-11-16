<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from "vue";
import { Lock, LogOut, Loader2, ShieldCheck } from "lucide-vue-next";
import { useAuthStore } from "../../stores/authStore";
import { authBackend } from "../../services/authBackendService";
import { api } from "../../db/api";

const IDLE_TIMEOUT_MS = 15 * 60 * 1000; // 15 minutes d'inactivité

const authStore = useAuthStore();
const isLocked = ref(false);
const secretInput = ref("");
const errorMsg = ref("");
const isChecking = ref(false);
let idleTimer: ReturnType<typeof setTimeout> | null = null;

const resetTimer = () => {
  if (idleTimer) clearTimeout(idleTimer);
  if (!authStore.isAuthenticated || isLocked.value) return;
  idleTimer = setTimeout(() => {
    if (authStore.isAuthenticated) {
      isLocked.value = true;
      secretInput.value = "";
      errorMsg.value = "";
    }
  }, IDLE_TIMEOUT_MS);
};

const ACTIVITY_EVENTS = [
  "mousemove",
  "keydown",
  "mousedown",
  "touchstart",
] as const;

onMounted(() => {
  ACTIVITY_EVENTS.forEach((ev) =>
    window.addEventListener(ev, resetTimer, { passive: true }),
  );
  resetTimer();
});

onUnmounted(() => {
  ACTIVITY_EVENTS.forEach((ev) => window.removeEventListener(ev, resetTimer));
  if (idleTimer) clearTimeout(idleTimer);
});

watch(
  () => authStore.isAuthenticated,
  (auth) => {
    if (!auth) {
      isLocked.value = false;
      if (idleTimer) clearTimeout(idleTimer);
    } else {
      resetTimer();
    }
  },
);

const unlock = async () => {
  const val = secretInput.value.trim();
  if (!val || isChecking.value) return;
  isChecking.value = true;
  errorMsg.value = "";

  try {
    // 1. Vérifier si un PIN Manager local correspond
    if (/^\d{4,8}$/.test(val)) {
      const ok = await api.verifierCodeManager(val).catch(() => false);
      if (ok) {
        isLocked.value = false;
        secretInput.value = "";
        resetTimer();
        return;
      }
    }

    // 2. Sinon vérifier le mot de passe du compte actif via le backend
    const ident =
      authStore.currentUser?.telephone || authStore.currentUser?.email || "";
    const schoolCode = authStore.currentSchool?.code_unique || "";
    if (!ident) {
      throw new Error("Identifiant introuvable, veuillez vous reconnecter.");
    }
    await authBackend.login(
      { identifiant: ident, mot_de_passe: val },
      schoolCode || undefined,
    );
    isLocked.value = false;
    secretInput.value = "";
    resetTimer();
  } catch (err: any) {
    errorMsg.value = err?.message || "Mot de passe ou code PIN incorrect.";
  } finally {
    isChecking.value = false;
  }
};

const logoutNow = () => {
  isLocked.value = false;
  secretInput.value = "";
  authStore.logout();
};
</script>

<template>
  <div
    v-if="isLocked && authStore.isAuthenticated"
    class="fixed inset-0 z-[100] flex items-center justify-center bg-black/75 backdrop-blur-md p-4 select-none"
    role="dialog"
    aria-modal="true"
  >
    <div
      class="w-full max-w-sm rounded-[6px] border border-border bg-card p-6 shadow-2xl space-y-4"
    >
      <div class="flex items-center gap-3">
        <div
          class="h-10 w-10 rounded-[5px] bg-brand-light border border-brand-border flex items-center justify-center text-brand-dark shrink-0"
        >
          <Lock class="h-5 w-5" />
        </div>
        <div class="min-w-0">
          <h2 class="text-sm font-bold text-foreground">Poste verrouillé</h2>
          <p class="text-xs text-muted-foreground truncate">
            {{ authStore.currentUser?.prenom }}
            {{ authStore.currentUser?.nom }} ·
            {{ authStore.currentUser?.roleLabel }}
          </p>
        </div>
      </div>

      <p class="text-xs text-muted-foreground">
        Session verrouillée après 15 minutes d'inactivité pour protéger les
        données du guichet.
      </p>

      <form @submit.prevent="unlock" class="space-y-3">
        <input
          v-model="secretInput"
          type="password"
          autofocus
          placeholder="Mot de passe ou Code PIN..."
          class="h-9 w-full rounded-[4px] border border-border bg-background px-3 text-xs text-foreground placeholder:text-muted-foreground focus:border-brand focus:outline-none focus:ring-2 focus:ring-brand/20"
        />

        <p v-if="errorMsg" class="text-[11px] font-medium text-rose-600">
          {{ errorMsg }}
        </p>

        <div class="flex items-center justify-between gap-2 pt-1">
          <button
            type="button"
            @click="logoutNow"
            class="inline-flex h-8 items-center gap-1.5 rounded-[4px] border border-border bg-muted/40 px-3 text-xs font-medium text-muted-foreground hover:text-foreground transition cursor-pointer"
          >
            <LogOut class="h-3.5 w-3.5" />
            <span>Déconnexion</span>
          </button>

          <button
            type="submit"
            :disabled="isChecking || !secretInput.trim()"
            class="inline-flex h-8 items-center gap-1.5 rounded-[4px] bg-brand-dark px-4 text-xs font-semibold text-white hover:bg-brand transition disabled:opacity-50 cursor-pointer"
          >
            <Loader2 v-if="isChecking" class="h-3.5 w-3.5 animate-spin" />
            <ShieldCheck v-else class="h-3.5 w-3.5" />
            <span>Déverrouiller</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import type { SeanceCode } from "../../types";
import {
  Video,
  ExternalLink,
  X,
  Copy,
  Check,
  Users,
  Clock,
  ShieldCheck,
  UserCheck,
  ChevronRight,
  ChevronDown,
  Loader2,
  Sparkles,
} from "lucide-vue-next";
import { useGuichetStore } from "../../stores/guichetStore";
import { useAuthStore } from "../../stores/authStore";
import { isSeanceCodeExpired } from "../../stores/coursCodeStore";
import { realtimeSocket } from "../../services/realtimeSocketService";
import { api } from "../../db/api";

const props = defineProps<{
  seance: SeanceCode | null;
}>();

const emit = defineEmits<{
  (e: "close"): void;
}>();

const guichetStore = useGuichetStore();
const authStore = useAuthStore();

const hasCopied = ref(false);
const showStudentsDrawer = ref(true);
const isSavingPresences = ref(false);

interface ConnectedStudent {
  candidatId: string;
  nomComplet: string;
  numeroDossier: string;
  telephone?: string;
  joinedAt: string;
  estPresent: boolean;
}

const connectedStudents = ref<ConnectedStudent[]>([]);
let unsubscribeWs: (() => void) | null = null;

// Détection du rôle formateur/modérateur
const isHost = computed(() => {
  return (
    authStore.currentRole === "moniteur" || authStore.currentRole === "manager"
  );
});

const isExpired = computed(() =>
  props.seance ? isSeanceCodeExpired(props.seance) : false,
);

// URL de visioconférence enrichie avec identité certifiée imposée dans Jitsi
const computedVisioUrl = computed(() => {
  if (!props.seance?.lien_visio) return "";
  const base = props.seance.lien_visio.split("#")[0];

  const user = authStore.currentUser;
  const roleLabel = isHost.value
    ? authStore.currentRole === "manager"
      ? "Directeur"
      : "Moniteur"
    : "Élève";

  const nomComplet = user?.prenom
    ? `${user.nom || ""} ${user.prenom || ""}`.trim()
    : user?.nom || "Apprenant";

  const dossierInfo = user?.dossier ? ` · ${user.dossier}` : "";
  const displayName = `${nomComplet} (${roleLabel}${dossierInfo})`;

  const params: string[] = [
    `userInfo.displayName="${displayName}"`,
    `config.startWithAudioMuted=${isHost.value ? "false" : "true"}`,
    `config.startWithVideoMuted=false`,
    `config.prejoinPageEnabled=false`,
    `config.disableDeepLinking=true`,
  ];

  return `${base}#${params.join("&")}`;
});

const copyLink = async () => {
  if (!props.seance?.lien_visio) return;
  try {
    await navigator.clipboard.writeText(props.seance.lien_visio);
    hasCopied.value = true;
    guichetStore.showToast(
      "Lien du salon copié dans le presse-papiers.",
      "info",
    );
    setTimeout(() => {
      hasCopied.value = false;
    }, 2000);
  } catch {
    guichetStore.showToast(props.seance.lien_visio, "info");
  }
};

const openExternal = () => {
  if (computedVisioUrl.value) {
    window.open(computedVisioUrl.value, "_blank", "noopener,noreferrer");
  }
};

const notifyJoin = () => {
  if (!props.seance?.id) return;
  const user = authStore.currentUser;

  realtimeSocket.sendEvent("CLASSE_VIRTUELLE_JOIN", {
    seanceId: props.seance.id,
    candidatId:
      user?.id || guichetStore.selectedCandidatId || "cand-visio-" + Date.now(),
    nomComplet:
      `${user?.nom || ""} ${user?.prenom || ""}`.trim() || "Élève Connecté",
    numeroDossier:
      user?.dossier ||
      guichetStore.activeDossierDetail?.dossier?.numero_dossier ||
      "DOS-PERMIS",
    telephone: user?.telephone || "",
    userRole: authStore.currentRole,
    joinedAt: new Date().toLocaleTimeString("fr-FR", {
      hour: "2-digit",
      minute: "2-digit",
    }),
  });
};

const notifyLeave = () => {
  if (!props.seance?.id) return;
  const user = authStore.currentUser;

  realtimeSocket.sendEvent("CLASSE_VIRTUELLE_LEAVE", {
    seanceId: props.seance.id,
    candidatId: user?.id,
    nomComplet: `${user?.nom || ""} ${user?.prenom || ""}`.trim(),
  });
};

onMounted(() => {
  // 1. Émettre l'événement de connexion de l'utilisateur
  notifyJoin();

  // 2. Écouter les arrivées des participants en direct via WebSocket
  unsubscribeWs = realtimeSocket.subscribe((msg) => {
    if (!props.seance) return;

    if (
      msg.type === "CLASSE_VIRTUELLE_JOIN" &&
      msg.data?.seanceId === props.seance.id
    ) {
      // Ignorer si c'est le formateur lui-même
      if (msg.data.userRole === "moniteur" || msg.data.userRole === "manager") {
        return;
      }

      const existing = connectedStudents.value.find(
        (s) => s.candidatId === msg.data.candidatId,
      );
      if (!existing) {
        connectedStudents.value.push({
          candidatId: msg.data.candidatId,
          nomComplet: msg.data.nomComplet,
          numeroDossier: msg.data.numeroDossier,
          telephone: msg.data.telephone,
          joinedAt:
            msg.data.joinedAt ||
            new Date().toLocaleTimeString("fr-FR", {
              hour: "2-digit",
              minute: "2-digit",
            }),
          estPresent: true,
        });

        if (isHost.value) {
          guichetStore.showToast(
            `🟢 ${msg.data.nomComplet} a rejoint la classe virtuelle.`,
            "info",
          );
        }
      }
    }

    if (
      msg.type === "CLASSE_VIRTUELLE_LEAVE" &&
      msg.data?.seanceId === props.seance.id
    ) {
      connectedStudents.value = connectedStudents.value.filter(
        (s) => s.candidatId !== msg.data.candidatId,
      );
    }
  });
});

onUnmounted(() => {
  notifyLeave();
  if (unsubscribeWs) {
    unsubscribeWs();
    unsubscribeWs = null;
  }
});

const handleValidateOnlinePresences = async () => {
  if (!props.seance) return;
  const activeStudents = connectedStudents.value.filter((s) => s.estPresent);

  if (activeStudents.length === 0) {
    guichetStore.showToast(
      "Aucun élève actuellement coché présent à valider.",
      "info",
    );
    return;
  }

  isSavingPresences.value = true;
  try {
    const presences = activeStudents.map((s) => ({
      candidat_id: s.candidatId,
      est_present: true,
      score_test: 0,
    }));

    await api.enregistrerSeanceCode({
      theme_traite: props.seance.theme_traite,
      formateur_nom: props.seance.formateur_nom,
      heure_debut: props.seance.heure_debut,
      heure_fin: props.seance.heure_fin,
      presences,
    });

    guichetStore.showToast(
      `Présence validée pour ${activeStudents.length} élève(s) connecté(s) ! +2h créditées sur leur dossier.`,
      "success",
    );
  } catch (err: any) {
    guichetStore.showToast(
      `Erreur enregistrement présence : ${err.message}`,
      "error",
    );
  } finally {
    isSavingPresences.value = false;
  }
};
</script>

<template>
  <div
    v-if="seance"
    class="fixed inset-0 z-50 flex items-center justify-center p-2 sm:p-4 bg-black/85 backdrop-blur-sm animate-in fade-in select-none"
    role="dialog"
    aria-modal="true"
  >
    <div
      class="bg-card border border-border w-full max-w-6xl h-[94vh] rounded-[6px] flex flex-col overflow-hidden text-foreground shadow-2xl"
    >
      <!-- En-tête de la classe virtuelle -->
      <header
        class="px-4 py-2.5 border-b border-border flex items-center justify-between bg-muted/40"
      >
        <div class="flex items-center space-x-2.5 min-w-0">
          <span
            class="inline-flex items-center gap-1.5 text-[10px] font-bold px-2 py-0.5 rounded-[4px] bg-rose-500/15 text-rose-600 border border-rose-500/30 uppercase tracking-wider"
          >
            <span class="w-2 h-2 rounded-full bg-rose-600 animate-pulse" />
            Direct Visio Jitsi
          </span>

          <div class="min-w-0 truncate">
            <h2
              class="text-xs sm:text-sm font-bold text-foreground truncate flex items-center gap-2"
            >
              <span>{{ seance.theme_traite }}</span>
              <span
                v-if="isHost"
                class="text-[10px] font-semibold px-1.5 py-0.2 rounded bg-brand-light text-brand-dark border border-brand-border"
              >
                Hôte Modérateur
              </span>
            </h2>
            <p class="text-[10px] text-muted-foreground truncate">
              Formateur : <strong>{{ seance.formateur_nom }}</strong> ·
              {{ seance.heure_debut }} - {{ seance.heure_fin }}
            </p>
          </div>
        </div>

        <div class="flex items-center space-x-1.5">
          <!-- Bouton panneau des élèves connectés pour le formateur -->
          <button
            v-if="isHost"
            type="button"
            @click="showStudentsDrawer = !showStudentsDrawer"
            class="flex items-center space-x-1.5 px-2.5 py-1 rounded-[4px] border text-xs font-semibold transition cursor-pointer"
            :class="
              showStudentsDrawer
                ? 'bg-brand-light text-brand-dark border-brand-border font-bold'
                : 'border-border bg-card hover:bg-muted text-foreground'
            "
            title="Afficher les élèves connectés"
          >
            <Users class="w-3.5 h-3.5 text-brand" />
            <span>Connectés ({{ connectedStudents.length }})</span>
          </button>

          <button
            type="button"
            @click="copyLink"
            class="flex items-center space-x-1 px-2.5 py-1 rounded-[4px] border border-border bg-card hover:bg-muted text-xs font-semibold text-foreground transition cursor-pointer"
            title="Copier le lien d'invitation"
          >
            <component
              :is="hasCopied ? Check : Copy"
              class="w-3.5 h-3.5 text-brand"
            />
            <span class="hidden sm:inline">{{
              hasCopied ? "Copié !" : "Partager"
            }}</span>
          </button>

          <button
            type="button"
            @click="openExternal"
            class="flex items-center space-x-1 px-2.5 py-1 rounded-[4px] bg-brand hover:bg-brand-dark text-white text-xs font-semibold transition cursor-pointer"
            title="Ouvrir en plein écran dans un nouvel onglet"
          >
            <ExternalLink class="w-3.5 h-3.5" />
            <span class="hidden sm:inline">Plein écran</span>
          </button>

          <button
            type="button"
            @click="$emit('close')"
            class="p-1.5 rounded-[4px] text-muted-foreground hover:text-foreground transition cursor-pointer ml-1"
            title="Quitter la classe virtuelle"
          >
            <X class="w-5 h-5" />
          </button>
        </div>
      </header>

      <!-- Zone centrale : Visioconférence Jitsi + Panneau latéral des élèves en direct -->
      <div class="flex-1 bg-zinc-950 overflow-hidden relative flex flex-row">
        <!-- Zone Iframe Jitsi Meet -->
        <div class="flex-1 h-full relative flex flex-col bg-black">
          <div
            v-if="isExpired"
            class="flex flex-col items-center justify-center h-full text-center p-6 space-y-2"
          >
            <Clock class="w-10 h-10 text-amber-400" />
            <p class="text-sm font-semibold text-white">
              Séance terminée — Accès au salon fermé
            </p>
            <p class="text-xs text-zinc-400">
              L'horaire de fin ({{ seance?.heure_fin }}) est dépassé.
            </p>
          </div>
          <iframe
            v-else-if="computedVisioUrl"
            :src="computedVisioUrl"
            class="w-full h-full border-none"
            allow="camera; microphone; fullscreen; display-capture; autoplay"
            title="Salle de classe virtuelle interactive"
          />
          <div
            v-else
            class="flex flex-col items-center justify-center h-full text-center p-6 space-y-2"
          >
            <Video class="w-10 h-10 text-muted-foreground/40" />
            <p class="text-sm font-semibold text-white">
              Lien de visioconférence non configuré
            </p>
          </div>
        </div>

        <!-- Panneau latéral Moniteur : Suivi des présences en temps réel -->
        <aside
          v-if="isHost && showStudentsDrawer"
          class="w-72 sm:w-80 border-l border-border bg-card flex flex-col h-full z-10 animate-in slide-in-from-right-4 duration-200"
        >
          <div
            class="p-3 border-b border-border flex items-center justify-between bg-muted/30"
          >
            <div>
              <h3
                class="text-xs font-bold text-foreground flex items-center gap-1.5"
              >
                <Users class="w-3.5 h-3.5 text-brand" />
                <span>Apprenants en direct</span>
              </h3>
              <p class="text-[10px] text-muted-foreground">
                Détection automatique WebSocket
              </p>
            </div>
            <span
              class="text-xs font-bold px-2 py-0.5 rounded-full bg-emerald-500/15 text-emerald-700 dark:text-emerald-400"
            >
              {{ connectedStudents.length }} connecté(s)
            </span>
          </div>

          <!-- Liste des élèves connectés -->
          <div class="flex-1 overflow-y-auto p-2 space-y-2">
            <div
              v-if="connectedStudents.length === 0"
              class="h-full flex flex-col items-center justify-center text-center p-4 text-muted-foreground space-y-2"
            >
              <Users class="w-8 h-8 opacity-40" />
              <p class="text-xs">En attente de connexion des élèves…</p>
              <p class="text-[10px] opacity-75">
                Dès qu'un élève clique sur « Rejoindre la classe », il apparaît
                automatiquement ici.
              </p>
            </div>

            <div
              v-for="student in connectedStudents"
              :key="student.candidatId"
              class="p-2.5 rounded-[4px] border border-border bg-background/80 hover:bg-muted/40 transition flex items-center justify-between gap-2"
            >
              <div class="min-w-0 flex-1">
                <div class="flex items-center gap-1.5">
                  <span
                    class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse shrink-0"
                  />
                  <span class="text-xs font-semibold text-foreground truncate">
                    {{ student.nomComplet }}
                  </span>
                </div>
                <div
                  class="text-[10px] text-muted-foreground flex items-center gap-2 mt-0.5"
                >
                  <span class="font-mono text-primary font-medium">{{
                    student.numeroDossier
                  }}</span>
                  <span>· Connecté à {{ student.joinedAt }}</span>
                </div>
              </div>

              <label
                class="flex items-center cursor-pointer shrink-0"
                title="Valider la présence"
              >
                <input
                  v-model="student.estPresent"
                  type="checkbox"
                  class="w-4 h-4 rounded border-border text-brand focus:ring-brand"
                />
              </label>
            </div>
          </div>

          <!-- Action d'émargement en 1 clic pour le formateur -->
          <div class="p-3 border-t border-border bg-muted/20 space-y-2">
            <button
              type="button"
              @click="handleValidateOnlinePresences"
              :disabled="connectedStudents.length === 0 || isSavingPresences"
              class="w-full h-9 rounded-[4px] bg-emerald-600 hover:bg-emerald-700 text-white font-semibold text-xs flex items-center justify-center gap-1.5 transition disabled:opacity-50 disabled:cursor-not-allowed shadow-xs cursor-pointer active:scale-95"
            >
              <Loader2
                v-if="isSavingPresences"
                class="w-3.5 h-3.5 animate-spin"
              />
              <Check v-else class="w-3.5 h-3.5 stroke-[3]" />
              <span>Valider présences (+2h créditées)</span>
            </button>
            <p class="text-[10px] text-muted-foreground text-center">
              Crédite automatiquement 2h officielles de code sur le dossier des
              élèves cochés.
            </p>
          </div>
        </aside>
      </div>

      <!-- Pied de session -->
      <footer
        class="px-4 py-2 border-t border-border bg-card flex items-center justify-between text-[11px] text-muted-foreground"
      >
        <span
          >Partage d'écran & Micro activés · Utilisez des écouteurs pour un
          meilleur confort sonore.</span
        >
        <span
          class="text-emerald-600 dark:text-emerald-400 font-medium flex items-center gap-1"
        >
          <ShieldCheck class="w-3.5 h-3.5" />
          Identification sécurisée certifiée
        </span>
      </footer>
    </div>
  </div>
</template>

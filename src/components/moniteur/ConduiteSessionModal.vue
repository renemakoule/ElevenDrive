<script setup lang="ts">
import { computed, nextTick, ref, watch, onMounted, onUnmounted } from "vue";
import type { SeanceConduite, CompetenceItem } from "../../types";
import { speechAudioService } from "../../services/speechAudioService";
import { realtimeSocket } from "../../services/realtimeSocketService";
import {
  X,
  Car,
  Navigation,
  Award,
  ShieldCheck,
  RotateCcw,
  Check,
  Volume2,
  Radio,
} from "lucide-vue-next";

const props = defineProps<{
  session: SeanceConduite;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (
    e: "save",
    payload: {
      sessionId: string;
      dossierId: string;
      dureeHeures: number;
      kmDebut: number;
      kmFin: number;
      commentaire: string;
      competences: CompetenceItem[];
      signatureEleve: string;
    },
  ): void;
}>();

const DUREE_SEANCE_H = 2;

type Niveau = CompetenceItem["niveau"];
type Onglet = "competences" | "vehicule" | "signature";

const onglets: { id: Onglet; label: string; icone: typeof Award }[] = [
  { id: "competences", label: "Compétences", icone: Award },
  { id: "vehicule", label: "Véhicule et km", icone: Navigation },
  { id: "signature", label: "Signature", icone: ShieldCheck },
];

const niveaux: { n: Niveau; label: string; actif: string }[] = [
  {
    n: 1,
    label: "Non acquis",
    actif: "border-destructive bg-destructive/15 text-destructive font-bold",
  },
  {
    n: 2,
    label: "En cours",
    actif:
      "border-amber-500 bg-amber-500/15 text-amber-600 dark:text-amber-400 font-bold",
  },
  {
    n: 3,
    label: "Acquis",
    actif:
      "border-sky-500 bg-sky-500/15 text-sky-600 dark:text-sky-400 font-bold",
  },
  {
    n: 4,
    label: "Maîtrisé",
    actif: "border-brand bg-brand/15 text-brand-dark font-bold",
  },
];

const localCompetences = ref<CompetenceItem[]>(
  JSON.parse(JSON.stringify(props.session.competences || [])),
);
const kmDebut = ref<number>(props.session.km_debut ?? 0);
const kmFin = ref<number>(
  props.session.km_fin && props.session.km_fin >= kmDebut.value
    ? props.session.km_fin
    : kmDebut.value + 15,
);
const commentaire = ref<string>(props.session.commentaire_moniteur ?? "");
const ongletActif = ref<Onglet>("competences");

const numeroLecon = computed(
  () => Math.floor(props.session.heures_deja_effectuees / DUREE_SEANCE_H) + 1,
);
const distance = computed(() => Math.max(0, kmFin.value - kmDebut.value));
const kmValide = computed(() => kmFin.value >= kmDebut.value);

const heuresApres = computed(
  () => props.session.heures_deja_effectuees + DUREE_SEANCE_H,
);
const progression = computed(() => {
  const total = props.session.heures_forfait_total;
  return total > 0
    ? Math.min(100, Math.round((heuresApres.value / total) * 100))
    : 0;
});

const setLevel = (comp: CompetenceItem, niveau: Niveau) => {
  comp.niveau = niveau;
  broadcastLiveUpdate("COMPETENCE");
};

const isWsConnected = realtimeSocket.isConnected;
const lastLiveSyncAt = ref<string>("");
let unsubscribeWs: (() => void) | null = null;

const broadcastLiveUpdate = (field: string) => {
  lastLiveSyncAt.value = new Date().toLocaleTimeString("fr-FR");
  realtimeSocket.sendEvent("LIVRET_EMARGEMENT_LIVE", {
    sessionId: props.session.id,
    dossierId: props.session.dossier_id,
    candidatNom: props.session.candidat_nom,
    field,
    kmDebut: kmDebut.value,
    kmFin: kmFin.value,
    commentaire: commentaire.value,
    competences: localCompetences.value,
    hasSigned: hasSigned.value,
  });
};

watch([kmDebut, kmFin, commentaire], () => {
  broadcastLiveUpdate("BILAN_KM");
});

onMounted(() => {
  unsubscribeWs = realtimeSocket.subscribe((msg) => {
    if (
      msg.type === "LIVRET_EMARGEMENT_LIVE" &&
      msg.data?.sessionId === props.session.id
    ) {
      if (Array.isArray(msg.data.competences)) {
        localCompetences.value = msg.data.competences;
      }
      if (typeof msg.data.kmDebut === "number")
        kmDebut.value = msg.data.kmDebut;
      if (typeof msg.data.kmFin === "number") kmFin.value = msg.data.kmFin;
      if (typeof msg.data.commentaire === "string")
        commentaire.value = msg.data.commentaire;
      lastLiveSyncAt.value = new Date().toLocaleTimeString("fr-FR");
    }
  });
});

const vocaliserCompetence = (comp: CompetenceItem) => {
  const niveauTrouve = niveaux.find((n) => n.n === comp.niveau);
  const niveauTxt = niveauTrouve ? niveauTrouve.label : "Non évalué";
  const texte = `Compétence ${comp.code} : ${comp.libelle}. Évaluation actuelle : ${niveauTxt}.`;
  speechAudioService.parlerCourt(texte);
};

const vocaliserBilan = () => {
  const texte = commentaire.value.trim()
    ? `Bilan de la leçon : ${commentaire.value}`
    : `Leçon numéro ${numeroLecon.value} pour ${props.session.candidat_nom}. Distance parcourue : ${distance.value} kilomètres.`;
  speechAudioService.parlerCourt(texte);
};

const handleClose = () => {
  speechAudioService.arreter();
  emit("close");
};

onUnmounted(() => {
  speechAudioService.arreter();
  if (unsubscribeWs) {
    unsubscribeWs();
    unsubscribeWs = null;
  }
});

/* Signature */
const canvasRef = ref<HTMLCanvasElement | null>(null);
const isDrawing = ref(false);
const hasSigned = ref(Boolean(props.session.signature_eleve_url));
let canvasPret = false;

watch(ongletActif, async (onglet) => {
  if (onglet !== "signature" || canvasPret) return;
  await nextTick();
  const canvas = canvasRef.value;
  if (!canvas) return;
  canvas.width = canvas.clientWidth;
  canvas.height = canvas.clientHeight;
  canvasPret = true;
});

const startDrawing = (e: PointerEvent) => {
  isDrawing.value = true;
  canvasRef.value?.setPointerCapture(e.pointerId);
  draw(e);
};

const stopDrawing = () => {
  isDrawing.value = false;
  canvasRef.value?.getContext("2d")?.beginPath();
  broadcastLiveUpdate("SIGNATURE");
};

const draw = (e: PointerEvent) => {
  const canvas = canvasRef.value;
  const ctx = canvas?.getContext("2d");
  if (!isDrawing.value || !canvas || !ctx) return;

  const rect = canvas.getBoundingClientRect();
  const x = e.clientX - rect.left;
  const y = e.clientY - rect.top;

  ctx.lineWidth = 2;
  ctx.lineCap = "round";
  ctx.strokeStyle = "#059669";
  ctx.lineTo(x, y);
  ctx.stroke();
  ctx.beginPath();
  ctx.moveTo(x, y);
  hasSigned.value = true;
};

const clearSignature = () => {
  const canvas = canvasRef.value;
  canvas?.getContext("2d")?.clearRect(0, 0, canvas.width, canvas.height);
  hasSigned.value = false;
};

const handleSave = () => {
  if (!kmValide.value) return;

  const sig =
    canvasRef.value && hasSigned.value
      ? canvasRef.value.toDataURL("image/png")
      : props.session.signature_eleve_url || "";

  speechAudioService.arreter();

  emit("save", {
    sessionId: props.session.id,
    dossierId: props.session.dossier_id,
    dureeHeures: DUREE_SEANCE_H,
    kmDebut: kmDebut.value,
    kmFin: kmFin.value,
    commentaire: commentaire.value,
    competences: localCompetences.value,
    signatureEleve: sig,
  });
};

const handleQuickPresent = () => {
  speechAudioService.arreter();
  emit("save", {
    sessionId: props.session.id,
    dossierId: props.session.dossier_id,
    dureeHeures: DUREE_SEANCE_H,
    kmDebut: kmDebut.value,
    kmFin: kmFin.value > kmDebut.value ? kmFin.value : kmDebut.value + 15,
    commentaire:
      commentaire.value.trim() || "Présence validée en 1 clic par le moniteur",
    competences: localCompetences.value,
    signatureEleve:
      props.session.signature_eleve_url || "PRESENCE_VALIDEE_1_CLIC",
  });
};
</script>

<template>
  <div
    class="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/50 p-0 sm:p-3 select-none"
    @click.self="handleClose"
  >
    <div
      class="flex max-h-[92dvh] w-full sm:max-w-md flex-col overflow-hidden rounded-t-[3px] sm:rounded-[3px] border-t sm:border border-border bg-card pb-safe"
    >
      <div
        class="sm:hidden mx-auto mt-2 h-1 w-9 rounded-full bg-border shrink-0"
      />
      <!-- En-tête -->
      <header
        class="flex items-center justify-between border-b border-border px-4 py-3"
      >
        <div class="flex items-center gap-2">
          <Car class="h-4 w-4 text-brand-dark" />
          <div>
            <div
              class="flex items-center gap-2 text-xs font-semibold leading-tight text-foreground"
            >
              <span>Leçon de conduite n° {{ numeroLecon }}</span>
              <span
                class="inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-bold"
                :class="
                  isWsConnected
                    ? 'bg-emerald-500/15 text-emerald-700 dark:text-emerald-300'
                    : 'bg-amber-500/15 text-amber-700'
                "
              >
                <Radio class="h-2.5 w-2.5 animate-pulse" />
                {{ isWsConnected ? "Direct WS" : "Sync locale" }}
              </span>
            </div>
            <div class="text-[11px] text-muted-foreground">
              {{ session.candidat_nom }} · {{ session.heure_debut }} –
              {{ session.heure_fin }}
              <span
                v-if="lastLiveSyncAt"
                class="text-[10px] text-emerald-600 dark:text-emerald-400 font-mono ml-1"
              >
                (Sync {{ lastLiveSyncAt }})
              </span>
            </div>
          </div>
        </div>
        <button
          type="button"
          @click="handleClose"
          class="rounded p-1 text-muted-foreground hover:bg-muted hover:text-foreground transition cursor-pointer"
          aria-label="Fermer"
        >
          <X class="h-4 w-4" />
        </button>
      </header>

      <!-- Navigation onglets -->
      <nav class="flex border-b border-border bg-muted/40 text-xs">
        <button
          v-for="o in onglets"
          :key="o.id"
          type="button"
          @click="ongletActif = o.id"
          class="flex flex-1 items-center justify-center gap-1.5 py-2.5 font-medium transition focus-visible:outline-none cursor-pointer"
          :class="
            ongletActif === o.id
              ? 'border-b-2 border-brand-dark bg-card font-semibold text-primary'
              : 'text-muted-foreground hover:text-foreground'
          "
        >
          <component :is="o.icone" class="h-3.5 w-3.5" />
          {{ o.label }}
        </button>
      </nav>

      <!-- Corps modal -->
      <div class="flex-1 overflow-y-auto p-4 text-xs">
        <!-- Compétences du livret -->
        <div v-show="ongletActif === 'competences'" class="space-y-3">
          <div
            v-for="comp in localCompetences"
            :key="comp.id"
            class="space-y-1.5 rounded-[4px] border border-border bg-card p-3"
          >
            <div class="flex items-center justify-between gap-2">
              <div class="flex items-center gap-1.5">
                <span
                  class="rounded bg-muted px-1.5 py-0.5 font-mono text-[10px] font-semibold text-foreground"
                >
                  {{ comp.code }}
                </span>
                <!-- Bouton d'écoute vocale de la compétence -->
                <button
                  type="button"
                  @click="vocaliserCompetence(comp)"
                  class="p-1 rounded hover:bg-muted text-muted-foreground hover:text-brand transition cursor-pointer"
                  title="Écouter la compétence à voix haute"
                >
                  <Volume2 class="w-3.5 h-3.5" />
                </button>
              </div>
              <span class="text-[10px] text-muted-foreground">{{
                comp.categorie
              }}</span>
            </div>
            <div class="font-medium text-foreground leading-snug">
              {{ comp.libelle }}
            </div>

            <!-- Boutons de niveau 1 à 4 -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-2 sm:gap-1.5 pt-1">
              <button
                v-for="niv in niveaux"
                :key="niv.n"
                type="button"
                @click="setLevel(comp, niv.n)"
                class="min-h-[40px] sm:min-h-0 rounded-[3px] border px-2 py-1.5 sm:py-1 text-xs sm:text-[11px] font-medium transition focus-visible:outline-none cursor-pointer"
                :class="
                  comp.niveau === niv.n
                    ? niv.actif
                    : 'border-border text-muted-foreground hover:bg-muted'
                "
              >
                {{ niv.label }}
              </button>
            </div>
          </div>
        </div>

        <!-- Véhicule et km -->
        <div v-show="ongletActif === 'vehicule'" class="space-y-4">
          <div
            class="rounded-[3px] border border-border bg-card p-3 space-y-1.5"
          >
            <div class="text-[11px] text-muted-foreground">
              Véhicule assigné
            </div>
            <div class="font-semibold text-sm text-foreground">
              {{ session.vehicule_modele }}
            </div>
            <div class="font-mono text-xs text-muted-foreground">
              {{ session.vehicule_immatriculation }}
            </div>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <label class="block">
              <span class="mb-1 block text-xs text-muted-foreground"
                >Km au départ</span
              >
              <input
                v-model.number="kmDebut"
                type="number"
                inputmode="numeric"
                min="0"
                class="min-h-[42px] sm:min-h-0 w-full rounded-[3px] border border-border bg-card px-2.5 py-2 sm:py-1.5 font-mono text-xs text-foreground focus:border-brand-dark focus:outline-none focus:ring-2 focus:ring-brand-dark/20 transition"
              />
            </label>
            <label class="block">
              <span class="mb-1 block text-xs text-muted-foreground"
                >Km à l'arrivée</span
              >
              <input
                v-model.number="kmFin"
                type="number"
                inputmode="numeric"
                min="0"
                :class="[
                  'min-h-[42px] sm:min-h-0 w-full rounded-[3px] border bg-card px-2.5 py-2 sm:py-1.5 font-mono text-xs text-foreground focus:outline-none focus:ring-2 transition',
                  kmValide
                    ? 'border-border focus:border-brand-dark focus:ring-brand-dark/20'
                    : 'border-destructive focus:border-destructive focus:ring-destructive/20',
                ]"
              />
            </label>
          </div>

          <div
            class="flex items-center justify-between rounded-[4px] border border-border bg-muted/30 px-3 py-2"
          >
            <span class="text-xs text-muted-foreground"
              >Distance parcourue</span
            >
            <span class="font-mono text-xs font-medium text-foreground"
              >{{ distance }} km</span
            >
          </div>
          <div v-if="!kmValide" class="text-xs text-destructive">
            Le kilométrage d'arrivée doit être supérieur ou égal à celui de
            départ.
          </div>
        </div>

        <!-- Signature et clôture -->
        <div v-show="ongletActif === 'signature'" class="space-y-4">
          <label class="block">
            <div class="flex items-center justify-between mb-1">
              <span class="text-xs text-muted-foreground"
                >Bilan et observations pédagogiques</span
              >
              <button
                v-if="commentaire"
                type="button"
                @click="vocaliserBilan"
                class="flex items-center gap-1 text-[11px] font-medium text-brand hover:underline cursor-pointer"
                title="Écouter le bilan de la séance à voix haute"
              >
                <Volume2 class="w-3 h-3" />
                <span>Écouter le bilan</span>
              </button>
            </div>
            <textarea
              v-model="commentaire"
              rows="3"
              placeholder="Ex: Bonne maîtrise du volant, travailler le démarrage en côte..."
              class="w-full resize-none rounded-[4px] border border-border bg-card p-2 text-xs leading-relaxed text-foreground focus:border-brand-dark focus:outline-none focus:ring-2 focus:ring-brand-dark/20 transition"
            ></textarea>
          </label>

          <div class="space-y-1.5">
            <div class="flex items-center justify-between">
              <span class="text-xs text-muted-foreground"
                >Signature tactile de l'élève</span
              >
              <button
                type="button"
                @click="clearSignature"
                class="flex items-center gap-1 text-xs text-muted-foreground hover:text-destructive transition cursor-pointer"
              >
                <RotateCcw class="h-3 w-3" />
                Effacer
              </button>
            </div>

            <div
              class="relative overflow-hidden rounded-[4px] border border-border bg-card"
            >
              <img
                v-if="session.signature_eleve_url && !hasSigned"
                :src="session.signature_eleve_url"
                alt="Signature enregistrée"
                class="absolute inset-0 h-full w-full object-contain p-2 pointer-events-none"
              />
              <canvas
                ref="canvasRef"
                class="block h-[130px] w-full touch-none"
                @pointerdown="startDrawing"
                @pointermove="draw"
                @pointerup="stopDrawing"
                @pointercancel="stopDrawing"
                @pointerleave="stopDrawing"
              ></canvas>
              <div
                v-if="!hasSigned && !session.signature_eleve_url"
                class="pointer-events-none absolute inset-0 flex items-center justify-center text-xs text-muted-foreground/60"
              >
                Signer ici au doigt ou au stylet
              </div>
            </div>
          </div>

          <div
            class="space-y-1.5 rounded-[4px] border border-border bg-muted/30 px-3 py-2"
          >
            <div class="flex items-center justify-between">
              <span class="text-foreground font-medium"
                >+{{ DUREE_SEANCE_H }} h décomptées du forfait</span
              >
              <span class="font-mono text-muted-foreground"
                >{{ heuresApres }} / {{ session.heures_forfait_total }} h</span
              >
            </div>
            <div class="h-1.5 overflow-hidden rounded-full bg-muted">
              <div
                class="h-full rounded-full bg-brand transition-all duration-500"
                :style="{ width: progression + '%' }"
              ></div>
            </div>
          </div>
        </div>
      </div>

      <!-- Actions -->
      <footer
        class="flex flex-wrap sm:flex-nowrap items-center justify-between gap-2 border-t border-border px-4 py-3"
      >
        <button
          v-if="session.statut !== 'EFFECTUE'"
          type="button"
          @click="handleQuickPresent"
          class="flex-1 sm:flex-initial inline-flex h-10 sm:h-8 items-center justify-center gap-1.5 rounded-[3px] bg-emerald-600 hover:bg-emerald-700 text-white px-3 text-xs font-semibold transition cursor-pointer active:scale-95"
          title="Marquer l'élève présent et valider la leçon en 1 clic sans formalités"
        >
          <Check class="h-3.5 w-3.5 stroke-[3]" />
          <span>Présent (1 clic)</span>
        </button>
        <div v-else class="hidden sm:block"></div>

        <div class="flex flex-1 sm:flex-initial items-center justify-end gap-2">
          <button
            type="button"
            @click="handleClose"
            class="inline-flex h-10 sm:h-8 items-center rounded-[3px] border border-border bg-card px-3 text-xs font-medium text-foreground hover:bg-muted transition cursor-pointer"
          >
            Annuler
          </button>
          <button
            type="button"
            @click="handleSave"
            :disabled="!kmValide"
            class="flex-1 sm:flex-initial inline-flex h-10 sm:h-8 items-center justify-center gap-1.5 rounded-[3px] bg-brand-dark px-3 text-xs font-semibold text-white hover:bg-brand disabled:cursor-not-allowed disabled:opacity-40 transition cursor-pointer"
          >
            <Check class="h-3.5 w-3.5" />
            Valider la leçon
          </button>
        </div>
      </footer>
    </div>
  </div>
</template>

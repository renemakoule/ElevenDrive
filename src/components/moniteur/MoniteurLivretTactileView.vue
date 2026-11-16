<script setup lang="ts">
import { ref, computed, watch, nextTick, onMounted, onUnmounted } from "vue";
import { useGuichetStore } from "../../stores/guichetStore";
import { realtimeSocket } from "../../services/realtimeSocketService";
import type { SeanceConduite, CompetenceItem } from "../../types";
import {
  ChevronLeft,
  ChevronRight,
  Radio,
  Award,
  RotateCcw,
  Check,
} from "lucide-vue-next";

const props = defineProps<{
  sessions: SeanceConduite[];
  selectedDate: string;
  dateLongue: string;
  defaultKilometrage: number;
  isWsConnected: boolean;
}>();

const emit = defineEmits<{
  prevDay: [];
  nextDay: [];
  selectDate: [iso: string];
  saveSession: [
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
  ];
}>();

const store = useGuichetStore();

const statusMap = {
  EFFECTUE: { label: "Terminé", dot: "bg-brand", text: "text-brand-dark" },
  EN_COURS: {
    label: "En cours",
    dot: "bg-sky-500 animate-pulse",
    text: "text-sky-600 dark:text-sky-400",
  },
  PROGRAMME: {
    label: "À venir",
    dot: "bg-muted-foreground/40",
    text: "text-muted-foreground",
  },
} as const;

const statusOf = (s: SeanceConduite) =>
  statusMap[s.statut as keyof typeof statusMap] ?? statusMap.PROGRAMME;

const pct = (value: number, max: number) =>
  max > 0 ? Math.min(100, Math.round((value / max) * 100)) : 0;

const selectedLivretSessionId = ref<string>("");
const selectedLivretSession = computed<SeanceConduite | null>(() => {
  if (!props.sessions.length) return null;
  if (selectedLivretSessionId.value) {
    const found = props.sessions.find(
      (s) => s.id === selectedLivretSessionId.value,
    );
    if (found) return found;
  }
  return (
    props.sessions.find((s) => s.statut === "EN_COURS") ||
    props.sessions.find((s) => s.statut === "PROGRAMME") ||
    props.sessions[0] ||
    null
  );
});

const livretCompetences = ref<CompetenceItem[]>([]);
const livretKmDebut = ref<number>(0);
const livretKmFin = ref<number>(15);
const livretCommentaire = ref<string>("");
const livretCanvasRef = ref<HTMLCanvasElement | null>(null);
const livretIsDrawing = ref(false);
const livretHasSigned = ref(false);

const niveauxLivret: {
  n: CompetenceItem["niveau"];
  label: string;
  actif: string;
}[] = [
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

const initLivretCanvas = () => {
  const canvas = livretCanvasRef.value;
  if (!canvas) return;
  canvas.width = canvas.clientWidth || 320;
  canvas.height = canvas.clientHeight || 140;
};

const syncLivretStateFromSession = async (s: SeanceConduite | null) => {
  if (!s) return;
  selectedLivretSessionId.value = s.id;
  livretCompetences.value = JSON.parse(JSON.stringify(s.competences || []));
  const startKm = s.km_debut ?? props.defaultKilometrage ?? 0;
  livretKmDebut.value = startKm;
  livretKmFin.value = s.km_fin && s.km_fin >= startKm ? s.km_fin : startKm + 15;
  livretCommentaire.value = s.commentaire_moniteur || "";
  livretHasSigned.value = Boolean(s.signature_eleve_url);
  await nextTick();
  initLivretCanvas();
};

watch(
  selectedLivretSession,
  (s) => {
    syncLivretStateFromSession(s);
  },
  { immediate: true },
);

let unsubscribeWs: (() => void) | null = null;

onMounted(async () => {
  await nextTick();
  initLivretCanvas();

  unsubscribeWs = realtimeSocket.subscribe((msg) => {
    if (
      msg.type === "LIVRET_EMARGEMENT_LIVE" &&
      msg.data?.sessionId &&
      selectedLivretSession.value?.id === msg.data.sessionId
    ) {
      if (Array.isArray(msg.data.competences))
        livretCompetences.value = msg.data.competences;
      if (typeof msg.data.kmDebut === "number")
        livretKmDebut.value = msg.data.kmDebut;
      if (typeof msg.data.kmFin === "number")
        livretKmFin.value = msg.data.kmFin;
      if (typeof msg.data.commentaire === "string")
        livretCommentaire.value = msg.data.commentaire;
    }
  });
});

onUnmounted(() => {
  if (unsubscribeWs) {
    unsubscribeWs();
    unsubscribeWs = null;
  }
});

const broadcastInlineLivretLive = (field: string) => {
  const s = selectedLivretSession.value;
  if (!s) return;
  realtimeSocket.sendEvent("LIVRET_EMARGEMENT_LIVE", {
    sessionId: s.id,
    dossierId: s.dossier_id,
    candidatNom: s.candidat_nom,
    field,
    kmDebut: livretKmDebut.value,
    kmFin: livretKmFin.value,
    commentaire: livretCommentaire.value,
    competences: livretCompetences.value,
    hasSigned: livretHasSigned.value,
  });
};

const setInlineLivretLevel = (
  comp: CompetenceItem,
  niv: CompetenceItem["niveau"],
) => {
  comp.niveau = niv;
  broadcastInlineLivretLive("COMPETENCE");
};

const startLivretDraw = (e: PointerEvent) => {
  const canvas = livretCanvasRef.value;
  if (!canvas) return;
  if (canvas.width === 0) initLivretCanvas();
  livretIsDrawing.value = true;
  canvas.setPointerCapture(e.pointerId);
  drawLivret(e);
};

const stopLivretDraw = () => {
  if (!livretIsDrawing.value) return;
  livretIsDrawing.value = false;
  livretCanvasRef.value?.getContext("2d")?.beginPath();
  broadcastInlineLivretLive("SIGNATURE");
};

const drawLivret = (e: PointerEvent) => {
  const canvas = livretCanvasRef.value;
  const ctx = canvas?.getContext("2d");
  if (!livretIsDrawing.value || !canvas || !ctx) return;
  const rect = canvas.getBoundingClientRect();
  const x = e.clientX - rect.left;
  const y = e.clientY - rect.top;
  ctx.lineWidth = 2.2;
  ctx.lineCap = "round";
  ctx.strokeStyle = "#059669";
  ctx.lineTo(x, y);
  ctx.stroke();
  ctx.beginPath();
  ctx.moveTo(x, y);
  livretHasSigned.value = true;
};

const clearLivretSignature = () => {
  const canvas = livretCanvasRef.value;
  if (canvas) {
    canvas.getContext("2d")?.clearRect(0, 0, canvas.width, canvas.height);
  }
  livretHasSigned.value = false;
  broadcastInlineLivretLive("CLEAR_SIGNATURE");
};

const handleSaveInlineLivret = () => {
  const s = selectedLivretSession.value;
  if (!s) return;
  if (livretKmFin.value < livretKmDebut.value) {
    store.showToast(
      "Le kilométrage final doit être supérieur ou égal au kilométrage de départ.",
      "error",
    );
    return;
  }
  const sig =
    livretCanvasRef.value && livretHasSigned.value
      ? livretCanvasRef.value.toDataURL("image/png")
      : s.signature_eleve_url || "";

  emit("saveSession", {
    sessionId: s.id,
    dossierId: s.dossier_id,
    dureeHeures: Math.max(1, Math.round((s.duree_minutes || 120) / 60)),
    kmDebut: livretKmDebut.value,
    kmFin: livretKmFin.value,
    commentaire: livretCommentaire.value,
    competences: livretCompetences.value,
    signatureEleve: sig,
  });
};

const handleQuickInlinePresent = () => {
  const s = selectedLivretSession.value;
  if (!s) return;
  const kmDeb = livretKmDebut.value;
  const kmF = livretKmFin.value > kmDeb ? livretKmFin.value : kmDeb + 15;

  emit("saveSession", {
    sessionId: s.id,
    dossierId: s.dossier_id,
    dureeHeures: Math.max(1, Math.round((s.duree_minutes || 120) / 60)),
    kmDebut: kmDeb,
    kmFin: kmF,
    commentaire:
      livretCommentaire.value.trim() ||
      "Présence validée en 1 clic par le moniteur",
    competences: livretCompetences.value,
    signatureEleve: s.signature_eleve_url || "PRESENCE_VALIDEE_1_CLIC",
  });
};
</script>

<template>
  <div class="grid grid-cols-1 gap-4 lg:grid-cols-12">
    <!-- Colonne gauche : Sélecteur d'élève / séance du jour -->
    <aside class="space-y-3 lg:col-span-4">
      <div class="rounded-[4px] border border-border bg-card p-3.5 space-y-3">
        <div class="flex items-center justify-between">
          <h2
            class="text-xs font-bold uppercase tracking-wider text-foreground"
          >
            Séances du {{ dateLongue }}
          </h2>
          <span
            class="rounded-full bg-brand-light px-2 py-0.5 text-[10px] font-bold text-brand-dark"
          >
            {{ sessions.length }} élève(s)
          </span>
        </div>

        <!-- Sélecteur de date rapide -->
        <div class="flex items-center gap-1.5">
          <button
            type="button"
            @click="emit('prevDay')"
            class="inline-flex h-7 w-7 items-center justify-center rounded border border-border bg-background hover:bg-muted"
          >
            <ChevronLeft class="h-3.5 w-3.5" />
          </button>
          <input
            type="date"
            :value="selectedDate"
            @change="
              (e) => emit('selectDate', (e.target as HTMLInputElement).value)
            "
            class="h-7 flex-1 rounded border border-border bg-background px-2 text-xs text-foreground"
          />
          <button
            type="button"
            @click="emit('nextDay')"
            class="inline-flex h-7 w-7 items-center justify-center rounded border border-border bg-background hover:bg-muted"
          >
            <ChevronRight class="h-3.5 w-3.5" />
          </button>
        </div>

        <div
          v-if="sessions.length"
          class="flex lg:block gap-2 overflow-x-auto lg:overflow-y-auto no-scrollbar snap-x lg:max-h-[520px] lg:space-y-2 lg:pr-1"
        >
          <button
            v-for="s in sessions"
            :key="s.id"
            type="button"
            @click="syncLivretStateFromSession(s)"
            class="min-w-[230px] lg:min-w-0 lg:w-full shrink-0 snap-start text-left rounded-[3px] border p-3 transition cursor-pointer"
            :class="
              selectedLivretSession?.id === s.id
                ? 'border-brand bg-brand-light/25 ring-1 ring-brand'
                : 'border-border bg-muted/20 hover:bg-muted/50'
            "
          >
            <div class="flex items-center justify-between gap-2">
              <span class="font-semibold text-xs text-foreground truncate">{{
                s.candidat_nom
              }}</span>
              <span class="font-mono text-[10px] text-muted-foreground"
                >{{ s.heure_debut }} - {{ s.heure_fin }}</span
              >
            </div>
            <div class="flex items-center justify-between mt-1.5 text-[11px]">
              <span class="text-muted-foreground truncate">{{
                s.formule_nom
              }}</span>
              <span
                class="flex items-center gap-1 font-medium"
                :class="statusOf(s).text"
              >
                <span
                  class="h-1.5 w-1.5 rounded-full"
                  :class="statusOf(s).dot"
                />
                {{ statusOf(s).label }}
              </span>
            </div>
            <div class="mt-2 flex items-center gap-2">
              <div class="h-1 flex-1 overflow-hidden rounded-full bg-muted">
                <div
                  class="h-full rounded-full bg-brand"
                  :style="{
                    width:
                      pct(s.heures_deja_effectuees, s.heures_forfait_total) +
                      '%',
                  }"
                />
              </div>
              <span class="text-[10px] font-mono text-muted-foreground">
                {{ s.heures_deja_effectuees }}/{{ s.heures_forfait_total }}h
              </span>
            </div>
          </button>
        </div>

        <div
          v-else
          class="rounded-[4px] border border-dashed border-border p-6 text-center text-xs text-muted-foreground"
        >
          Aucune séance planifiée à cette date.
        </div>
      </div>
    </aside>

    <!-- Colonne droite : Fiche d'évaluation REMC & Signature Tactile en direct -->
    <section class="lg:col-span-8">
      <div
        v-if="selectedLivretSession"
        class="rounded-[4px] border border-border bg-card overflow-hidden"
      >
        <div
          class="flex flex-wrap items-center justify-between gap-3 border-b border-border bg-muted/30 px-4 py-3"
        >
          <div>
            <div class="flex items-center gap-2">
              <h3 class="text-sm font-bold text-foreground">
                {{ selectedLivretSession.candidat_nom }}
              </h3>
              <span
                class="rounded bg-brand-light px-2 py-0.5 text-[10px] font-bold text-brand-dark"
              >
                {{ selectedLivretSession.formule_nom }}
              </span>
            </div>
            <p class="text-xs text-muted-foreground">
              Créneau : {{ selectedLivretSession.heure_debut }} –
              {{ selectedLivretSession.heure_fin }} · Progression :
              {{ selectedLivretSession.heures_deja_effectuees }} /
              {{ selectedLivretSession.heures_forfait_total }} h
            </p>
          </div>

          <div class="flex items-center gap-2">
            <span
              class="inline-flex items-center gap-1 rounded-full px-2.5 py-1 text-[10px] font-bold"
              :class="
                isWsConnected
                  ? 'bg-emerald-500/15 text-emerald-700 dark:text-emerald-300'
                  : 'bg-amber-500/15 text-amber-700'
              "
            >
              <Radio class="h-3 w-3 animate-pulse" />
              Diffusion WebSocket instantanée
            </span>
          </div>
        </div>

        <div class="p-4 space-y-5">
          <!-- 1. Grille des compétences REMC -->
          <div class="space-y-2.5">
            <h4
              class="text-xs font-bold uppercase tracking-wider text-muted-foreground flex items-center gap-1.5"
            >
              <Award class="h-3.5 w-3.5 text-brand" />
              1. Évaluation des compétences REMC en direct
            </h4>
            <div class="grid gap-2.5 sm:grid-cols-2">
              <div
                v-for="comp in livretCompetences"
                :key="comp.id"
                class="rounded-[4px] border border-border bg-background p-3 space-y-2"
              >
                <div class="flex items-center justify-between gap-2">
                  <span
                    class="rounded bg-muted px-1.5 py-0.5 font-mono text-[10px] font-bold text-foreground"
                  >
                    {{ comp.code }}
                  </span>
                  <span class="text-[10px] text-muted-foreground">{{
                    comp.categorie
                  }}</span>
                </div>
                <p class="text-xs font-medium text-foreground leading-snug">
                  {{ comp.libelle }}
                </p>
                <div class="grid grid-cols-2 gap-1.5 pt-1">
                  <button
                    v-for="niv in niveauxLivret"
                    :key="niv.n"
                    type="button"
                    @click="setInlineLivretLevel(comp, niv.n)"
                    class="min-h-[40px] sm:min-h-0 rounded-[3px] border px-2 py-1.5 sm:py-1 text-xs sm:text-[11px] font-medium transition cursor-pointer"
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
          </div>

          <!-- 2. Relevé Kilométrique & Observations -->
          <div
            class="grid grid-cols-2 gap-3 sm:grid-cols-3 pt-2 border-t border-border"
          >
            <label class="block">
              <span class="mb-1 block text-xs font-medium text-muted-foreground"
                >Kilométrage départ</span
              >
              <input
                v-model.number="livretKmDebut"
                @input="broadcastInlineLivretLive('KM_DEBUT')"
                type="number"
                inputmode="numeric"
                min="0"
                class="min-h-[42px] sm:min-h-0 w-full rounded-[3px] border border-border bg-background px-3 py-2 sm:py-1.5 font-mono text-xs text-foreground"
              />
            </label>
            <label class="block">
              <span class="mb-1 block text-xs font-medium text-muted-foreground"
                >Kilométrage arrivée</span
              >
              <input
                v-model.number="livretKmFin"
                @input="broadcastInlineLivretLive('KM_FIN')"
                type="number"
                inputmode="numeric"
                min="0"
                class="min-h-[42px] sm:min-h-0 w-full rounded-[3px] border border-border bg-background px-3 py-2 sm:py-1.5 font-mono text-xs text-foreground"
              />
            </label>
            <div class="col-span-2 sm:col-span-1 flex flex-col justify-end">
              <div
                class="min-h-[40px] sm:min-h-0 rounded-[3px] border border-border bg-muted/30 px-3 py-2 sm:py-1.5 flex items-center justify-between"
              >
                <span class="text-xs text-muted-foreground">Distance :</span>
                <span class="font-mono text-xs font-bold text-brand-dark">
                  +{{ Math.max(0, livretKmFin - livretKmDebut) }} km
                </span>
              </div>
            </div>
          </div>

          <label class="block">
            <span class="mb-1 block text-xs font-medium text-muted-foreground">
              Bilan pédagogique et observations du moniteur (synchronisé en
              direct)
            </span>
            <textarea
              v-model="livretCommentaire"
              @input="broadcastInlineLivretLive('COMMENTAIRE')"
              rows="2"
              placeholder="Observations sur la maîtrise du véhicule, trajectoires, priorités..."
              class="w-full rounded-[4px] border border-border bg-background p-2.5 text-xs text-foreground"
            ></textarea>
          </label>

          <!-- 3. Pavé de Signature Tactile -->
          <div class="space-y-2 pt-2 border-t border-border">
            <div class="flex items-center justify-between">
              <span class="text-xs font-bold text-foreground"
                >Signature tactile de l'apprenant</span
              >
              <button
                type="button"
                @click="clearLivretSignature"
                class="inline-flex items-center gap-1 text-xs text-muted-foreground hover:text-destructive cursor-pointer"
              >
                <RotateCcw class="h-3 w-3" />
                Effacer la signature
              </button>
            </div>

            <div
              class="relative overflow-hidden rounded-[4px] border border-border bg-background"
            >
              <img
                v-if="
                  selectedLivretSession.signature_eleve_url && !livretHasSigned
                "
                :src="selectedLivretSession.signature_eleve_url"
                alt="Signature élève"
                class="absolute inset-0 h-full w-full object-contain p-2 pointer-events-none"
              />
              <canvas
                ref="livretCanvasRef"
                class="block h-[140px] w-full touch-none cursor-crosshair"
                @pointerdown="startLivretDraw"
                @pointermove="drawLivret"
                @pointerup="stopLivretDraw"
                @pointercancel="stopLivretDraw"
                @pointerleave="stopLivretDraw"
              ></canvas>
              <div
                v-if="
                  !livretHasSigned && !selectedLivretSession.signature_eleve_url
                "
                class="pointer-events-none absolute inset-0 flex items-center justify-center text-xs text-muted-foreground/60"
              >
                Signez ici au doigt, au stylet ou à la souris pour émarger la
                séance
              </div>
            </div>
          </div>

          <!-- Bouton de validation finale -->
          <div
            class="flex flex-wrap sm:flex-nowrap items-center justify-between gap-2 pt-2 border-t border-border"
          >
            <button
              v-if="selectedLivretSession.statut !== 'EFFECTUE'"
              type="button"
              @click="handleQuickInlinePresent"
              class="flex-1 sm:flex-initial min-h-[42px] sm:min-h-0 inline-flex items-center justify-center gap-1.5 rounded-[3px] bg-emerald-600 hover:bg-emerald-700 text-white px-3 py-2 text-xs font-semibold transition cursor-pointer active:scale-95"
              title="Marquer présent et valider la séance immédiatement"
            >
              <Check class="h-4 w-4 stroke-[3]" />
              <span>Présent (1 clic)</span>
            </button>
            <div v-else class="hidden sm:block"></div>

            <button
              type="button"
              @click="handleSaveInlineLivret"
              class="flex-1 sm:flex-initial min-h-[42px] sm:min-h-0 inline-flex items-center justify-center gap-2 rounded-[3px] bg-brand-dark px-4 py-2 text-xs font-bold text-white hover:bg-brand transition cursor-pointer"
            >
              <Check class="h-4 w-4" />
              <span>{{
                selectedLivretSession.statut === "EFFECTUE"
                  ? "Mettre à jour l'émargement"
                  : "Valider et émarger avec signature"
              }}</span>
            </button>
          </div>
        </div>
      </div>

      <div
        v-else
        class="rounded-[4px] border border-border bg-card p-10 text-center space-y-2"
      >
        <p class="text-sm font-semibold text-foreground">
          Aucune séance sélectionnée
        </p>
        <p class="text-xs text-muted-foreground">
          Choisissez une séance de conduite à gauche ou changez de date.
        </p>
      </div>
    </section>
  </div>
</template>

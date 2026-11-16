<script setup lang="ts">
import { computed } from "vue";
import { useCoursCodeStore } from "../../../stores/coursCodeStore";
import type {
  SeanceConduite,
  CompetenceItem,
  ModeFormationTheorique,
} from "../../../types";
import {
  Radio,
  Video,
  Car,
  ShieldCheck,
  Phone,
  MapPin,
  CalendarDays,
} from "lucide-vue-next";

export interface CandidatRecuItem {
  numero: string;
  date: string;
  heure: string;
  montant: number;
  mode: string;
  hash: string;
}

export interface CandidatProfileSummary {
  id: string;
  dossierId: string;
  nom: string;
  numeroDossier: string;
  telephone: string;
  formule: string;
  categoriePermis: string;
  agence: string;
  montantTotal: number;
  dejaPaye: number;
  soldeRestant: number;
  heuresCode: number;
  heuresCodeMax: number;
  heuresConduite: number;
  heuresConduiteMax: number;
  dernierTest: number;
  dateInscription: string;
  recus: CandidatRecuItem[];
  modeTheorique: ModeFormationTheorique;
}

export interface LiveLivretEval {
  sessionId: string;
  dossierId: string;
  candidatNom: string;
  kmDebut: number;
  kmFin: number;
  commentaire: string;
  competences: CompetenceItem[];
  hasSigned: boolean;
  updatedAt: string;
}

const props = defineProps<{
  candidat: CandidatProfileSummary;
  liveLivretEvaluation: LiveLivretEval | null;
  seancesConduite: SeanceConduite[];
  lastQuizScore: number | null;
}>();

const emit = defineEmits<{
  (e: "open-cours"): void;
  (e: "open-reservation"): void;
  (e: "start-payment", provider: string): void;
}>();

const coursStore = useCoursCodeStore();

const SEUIL_EXAMEN = 35;

const modeLabels: Record<string, string> = {
  ESPECES: "Espèces",
  WAVE: "Wave",
  ORANGE_MONEY: "Orange Money",
  MTN_MOMO: "MTN MoMo",
  MOOV: "Moov Money",
  VIREMENT: "Virement",
};

const paymentProviders = [
  { label: "Wave", dot: "bg-sky-500" },
  { label: "Orange Money", dot: "bg-orange-500" },
  { label: "MTN MoMo", dot: "bg-yellow-400" },
];

const pct = (value: number, max: number) =>
  max > 0 ? Math.round((value / max) * 100) : 0;
const fcfa = (n: number) => n.toLocaleString("fr-FR");

const initiales = computed(() =>
  props.candidat.nom
    .split(" ")
    .slice(0, 2)
    .map((p) => p.charAt(0))
    .join(""),
);

const score = computed(
  () => props.lastQuizScore ?? props.candidat.dernierTest ?? 0,
);
const eligible = computed(() => score.value >= SEUIL_EXAMEN);

const stats = computed(() => {
  const c = props.candidat;
  return [
    {
      key: "solde",
      label: "Solde à régler",
      value: `${fcfa(c.soldeRestant)} FCFA`,
      progress: pct(c.dejaPaye, c.montantTotal),
      caption: `${fcfa(c.dejaPaye)} réglés sur ${fcfa(c.montantTotal)}`,
    },
    {
      key: "code",
      label: "Code de la route",
      value: `${c.heuresCode} / ${c.heuresCodeMax} h`,
      progress: pct(c.heuresCode, c.heuresCodeMax),
      caption:
        score.value > 0
          ? `Dernier test : ${score.value}/40`
          : "Aucun test blanc effectué",
      badge:
        score.value > 0
          ? eligible.value
            ? "Éligible examen"
            : "À retravailler"
          : "Entraînement",
      badgeOk: eligible.value,
    },
    {
      key: "conduite",
      label: "Conduite",
      value: `${c.heuresConduite} / ${c.heuresConduiteMax} h`,
      progress: pct(c.heuresConduite, c.heuresConduiteMax),
      caption: `${Math.max(0, c.heuresConduiteMax - c.heuresConduite)} h restantes au programme`,
    },
  ];
});

const infos = computed(() => {
  const c = props.candidat;
  return [
    {
      icon: Phone,
      label: "Téléphone WhatsApp",
      value: c.telephone || "Non renseigné",
    },
    { icon: MapPin, label: "Auto-École", value: c.agence },
    {
      icon: CalendarDays,
      label: "Inscrit le",
      value: c.dateInscription,
    },
  ];
});
</script>

<template>
  <div class="space-y-4">
    <!-- Bannière Temps Réel : Émargement & Évaluation REMC en direct par le Moniteur -->
    <div
      v-if="liveLivretEvaluation"
      class="rounded-[3px] border border-emerald-500/40 bg-emerald-500/5 p-3.5 md:p-4 space-y-3"
    >
      <div class="flex flex-wrap items-center justify-between gap-2">
        <div class="flex items-center gap-2">
          <span
            class="inline-flex items-center gap-1.5 rounded-full bg-emerald-500/20 px-2.5 py-0.5 text-[10px] font-bold uppercase tracking-wider text-emerald-700 dark:text-emerald-300"
          >
            <Radio class="h-3 w-3 animate-pulse" />
            Émargement & Livret REMC en direct
          </span>
          <span class="text-xs font-semibold text-foreground">
            Votre moniteur évalue actuellement votre séance de conduite
          </span>
        </div>
        <span class="font-mono text-[11px] text-muted-foreground">
          Mis à jour à {{ liveLivretEvaluation.updatedAt }} · +{{
            Math.max(
              0,
              liveLivretEvaluation.kmFin - liveLivretEvaluation.kmDebut,
            )
          }}
          km parcourus
        </span>
      </div>

      <div
        v-if="liveLivretEvaluation.competences.length"
        class="grid grid-cols-1 gap-2 sm:grid-cols-2 lg:grid-cols-3"
      >
        <div
          v-for="comp in liveLivretEvaluation.competences"
          :key="comp.id"
          class="flex items-center justify-between gap-2 rounded border border-border bg-card px-2.5 py-1.5 text-xs"
        >
          <span class="truncate font-medium text-foreground">
            {{ comp.code }} · {{ comp.libelle }}
          </span>
          <span
            class="shrink-0 rounded px-1.5 py-0.5 text-[10px] font-bold"
            :class="{
              'bg-destructive/15 text-destructive': comp.niveau === 1,
              'bg-amber-500/15 text-amber-600': comp.niveau === 2,
              'bg-sky-500/15 text-sky-600': comp.niveau === 3,
              'bg-brand-light text-brand-dark': comp.niveau === 4,
            }"
          >
            Niv. {{ comp.niveau }}/4
          </span>
        </div>
      </div>

      <p
        v-if="liveLivretEvaluation.commentaire"
        class="text-xs italic text-muted-foreground"
      >
        Observation du moniteur : « {{ liveLivretEvaluation.commentaire }} »
      </p>
    </div>

    <!-- Bannière Prochain Cours Théorique ou Classe Virtuelle si disponible -->
    <div
      v-if="coursStore.prochaineSeanceEnLigne || coursStore.seances[0]"
      class="rounded-[5px] border border-sky-500/30 bg-sky-500/5 p-4 flex flex-col sm:flex-row sm:items-center justify-between gap-3"
    >
      <div class="flex items-center space-x-3">
        <div
          class="w-10 h-10 rounded-lg bg-sky-500/10 text-sky-600 flex items-center justify-center shrink-0"
        >
          <Video class="w-5 h-5" />
        </div>
        <div>
          <div class="flex items-center gap-2">
            <span
              class="text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-full bg-sky-500/20 text-sky-700 dark:text-sky-300"
            >
              {{
                (coursStore.prochaineSeanceEnLigne || coursStore.seances[0])
                  .type_seance === "EN_LIGNE"
                  ? "Cours de Code en Direct (Visio)"
                  : "Cours de Code en Salle"
              }}
            </span>
            <span class="text-xs text-muted-foreground font-mono">
              {{
                (coursStore.prochaineSeanceEnLigne || coursStore.seances[0])
                  .date_seance
              }}
              ·
              {{
                (coursStore.prochaineSeanceEnLigne || coursStore.seances[0])
                  .heure_debut
              }}
              -
              {{
                (coursStore.prochaineSeanceEnLigne || coursStore.seances[0])
                  .heure_fin
              }}
            </span>
          </div>
          <h4 class="text-xs sm:text-sm font-bold text-foreground mt-0.5">
            {{
              (coursStore.prochaineSeanceEnLigne || coursStore.seances[0])
                .theme_traite
            }}
            (avec
            {{
              (coursStore.prochaineSeanceEnLigne || coursStore.seances[0])
                .formateur_nom
            }})
          </h4>
        </div>
      </div>

      <button
        type="button"
        @click="emit('open-cours')"
        class="w-full sm:w-auto min-h-[40px] sm:min-h-0 inline-flex items-center justify-center gap-1.5 px-3.5 py-2 rounded-[3px] bg-sky-600 hover:bg-sky-700 text-white font-bold text-xs transition cursor-pointer"
      >
        <Video class="w-3.5 h-3.5" />
        <span>Voir les cours ({{ coursStore.seances.length }})</span>
      </button>
    </div>

    <!-- Indicateurs 3 colonnes (Carrousel horizontal tactile sur mobile) -->
    <section
      class="flex overflow-x-auto snap-x snap-mandatory gap-2.5 no-scrollbar md:grid md:grid-cols-3 md:gap-0 md:divide-x md:divide-border md:rounded-[3px] md:border md:border-border md:bg-card"
    >
      <div
        v-for="s in stats"
        :key="s.key"
        class="min-w-[78%] sm:min-w-[260px] md:min-w-0 snap-start rounded-[3px] border border-border bg-card md:border-0 md:rounded-none p-3.5 md:p-4 shrink-0"
      >
        <div class="flex items-center justify-between gap-2">
          <span class="text-xs text-muted-foreground">{{ s.label }}</span>
          <span
            v-if="s.badge"
            class="rounded-full px-2 py-0.5 text-[10px] font-medium"
            :class="
              s.badgeOk
                ? 'bg-brand-light text-brand-dark'
                : 'bg-amber-500/10 text-amber-600 dark:text-amber-400'
            "
          >
            {{ s.badge }}
          </span>
        </div>

        <p class="mt-1 text-lg font-semibold tabular-nums text-foreground">
          {{ s.value }}
        </p>

        <div
          class="mt-2.5 h-1 overflow-hidden rounded-full bg-muted"
          role="progressbar"
          :aria-valuenow="s.progress"
          aria-valuemin="0"
          aria-valuemax="100"
          :aria-label="s.label"
        >
          <div
            class="h-full rounded-full bg-brand transition-all duration-500"
            :style="{ width: s.progress + '%' }"
          />
        </div>

        <p class="mt-2 text-[11px] text-muted-foreground">
          {{ s.caption }}
        </p>
      </div>
    </section>

    <!-- Contenu : reçus & conduite (gauche) · dossier & paiement (droite) -->
    <div class="grid grid-cols-1 items-start gap-4 lg:grid-cols-12">
      <!-- Reçus de paiement réels & séances conduite -->
      <section
        class="overflow-hidden rounded-[5px] border border-border bg-card lg:col-span-7"
      >
        <div class="flex items-center justify-between px-4 py-3">
          <h2 class="text-sm font-semibold text-foreground">
            Reçus de paiement officiels
          </h2>
          <span class="text-xs text-muted-foreground"
            >{{ candidat.recus.length }} émis</span
          >
        </div>

        <!-- Cartes reçus sur mobile (< md) -->
        <div class="md:hidden divide-y divide-border border-t border-border">
          <div
            v-if="candidat.recus.length === 0"
            class="px-4 py-6 text-center text-xs text-muted-foreground"
          >
            Aucun paiement enregistré pour ce dossier.
          </div>
          <div
            v-for="rec in candidat.recus"
            :key="rec.numero"
            class="px-3.5 py-3 flex items-center justify-between gap-2 text-xs"
          >
            <div class="min-w-0 space-y-0.5">
              <div class="font-mono font-semibold text-foreground truncate">
                {{ rec.numero }}
              </div>
              <div class="text-[11px] text-muted-foreground">
                {{ rec.date }} · {{ rec.heure }} ·
                {{ modeLabels[rec.mode] ?? rec.mode }}
              </div>
            </div>
            <div class="text-right shrink-0">
              <div class="font-semibold tabular-nums text-foreground">
                {{ fcfa(rec.montant) }} F
              </div>
              <div
                class="font-mono text-[10px] text-muted-foreground truncate max-w-[90px]"
              >
                {{ rec.hash }}
              </div>
            </div>
          </div>
        </div>

        <div class="hidden md:block overflow-x-auto">
          <table class="w-full text-left text-xs">
            <thead>
              <tr
                class="border-y border-border bg-muted/50 text-[11px] text-muted-foreground"
              >
                <th class="px-4 py-2 font-medium">Reçu</th>
                <th class="px-3 py-2 font-medium">Mode</th>
                <th class="px-3 py-2 text-right font-medium">Montant</th>
                <th class="hidden px-4 py-2 font-medium sm:table-cell">
                  Réf. Transaction
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-border">
              <tr v-if="candidat.recus.length === 0">
                <td
                  colspan="4"
                  class="px-4 py-6 text-center text-xs text-muted-foreground"
                >
                  Aucun paiement enregistré pour ce dossier.
                </td>
              </tr>
              <tr
                v-for="rec in candidat.recus"
                :key="rec.numero"
                class="transition hover:bg-muted/50"
              >
                <td class="px-4 py-2.5">
                  <div class="font-mono font-medium text-foreground">
                    {{ rec.numero }}
                  </div>
                  <div class="text-[11px] text-muted-foreground">
                    {{ rec.date }} · {{ rec.heure }}
                  </div>
                </td>
                <td class="px-3 py-2.5 text-foreground">
                  {{ modeLabels[rec.mode] ?? rec.mode }}
                </td>
                <td
                  class="px-3 py-2.5 text-right font-medium tabular-nums text-foreground"
                >
                  {{ fcfa(rec.montant) }} F
                </td>
                <td class="hidden px-4 py-2.5 sm:table-cell">
                  <span
                    class="block max-w-[120px] truncate font-mono text-[10px] text-muted-foreground"
                    :title="rec.hash"
                  >
                    {{ rec.hash }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Séances de conduite réservées & Livret en temps réel -->
        <div class="border-t border-border p-4 space-y-2.5 bg-muted/10">
          <div class="flex items-center justify-between">
            <h3
              class="text-xs font-bold uppercase tracking-wider text-foreground flex items-center gap-1.5"
            >
              <Car class="h-3.5 w-3.5 text-brand" />
              Séances de Conduite & Livret Temps Réel
            </h3>
            <button
              type="button"
              @click="emit('open-reservation')"
              class="text-[11px] font-semibold text-brand-dark hover:underline cursor-pointer"
            >
              + Réserver un créneau
            </button>
          </div>

          <div v-if="seancesConduite.length > 0" class="space-y-2">
            <div
              v-for="sc in seancesConduite"
              :key="sc.id"
              class="rounded-[4px] border border-border bg-card p-3 flex flex-col sm:flex-row sm:items-center justify-between gap-2 text-xs"
            >
              <div class="space-y-0.5">
                <div class="flex flex-wrap items-center gap-2">
                  <span
                    class="rounded bg-muted px-1.5 py-0.5 font-mono text-[10px] font-bold text-foreground"
                    >{{ sc.date_seance }}</span
                  >
                  <span class="font-mono font-bold text-foreground"
                    >{{ sc.heure_debut }} – {{ sc.heure_fin }}</span
                  >
                  <span
                    class="rounded-full px-2 py-0.5 text-[10px] font-bold"
                    :class="
                      sc.statut === 'EFFECTUE'
                        ? 'bg-emerald-500/15 text-emerald-700 dark:text-emerald-300'
                        : 'bg-sky-500/15 text-sky-600'
                    "
                  >
                    {{
                      sc.statut === "EFFECTUE"
                        ? "Émargée & Validée"
                        : "Programmée"
                    }}
                  </span>
                </div>
                <p class="text-[11px] text-muted-foreground">
                  Moniteur : {{ sc.moniteur_nom || "Moniteur Assigné" }} ·
                  Véhicule : {{ sc.vehicule_modele || "Véhicule École" }}
                  <span v-if="sc.vehicule_immatriculation">
                    ({{ sc.vehicule_immatriculation }})
                  </span>
                  <span v-if="sc.commentaire_moniteur">
                    · « {{ sc.commentaire_moniteur }} »</span
                  >
                </p>
              </div>
              <div
                class="text-right font-mono text-[11px] text-muted-foreground"
              >
                <span v-if="sc.km_fin && sc.km_debut">
                  +{{ Math.max(0, sc.km_fin - sc.km_debut) }} km
                </span>
                <span v-else>{{ sc.duree_minutes }} min</span>
              </div>
            </div>
          </div>
          <p v-else class="text-xs text-muted-foreground">
            Aucune séance de conduite planifiée aujourd'hui. Cliquez sur «
            Réserver une leçon » pour choisir un créneau.
          </p>
        </div>
      </section>

      <!-- Colonne droite : Profil & Solde -->
      <div class="space-y-4 lg:col-span-5">
        <section class="rounded-[5px] border border-border bg-card p-4">
          <div class="flex items-center gap-3">
            <div
              class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-brand-light text-xs font-semibold text-brand-dark"
              aria-hidden="true"
            >
              {{ initiales }}
            </div>
            <div class="min-w-0 flex-1">
              <h2 class="truncate text-sm font-semibold text-foreground">
                {{ candidat.nom }}
              </h2>
              <p class="truncate text-xs text-muted-foreground">
                {{ candidat.formule }} ·
                <span class="font-mono font-medium">{{
                  candidat.numeroDossier
                }}</span>
              </p>
            </div>
            <span
              class="inline-flex shrink-0 items-center gap-1 rounded-full bg-brand-light px-2 py-0.5 text-[10px] font-medium text-brand-dark"
            >
              <ShieldCheck class="h-3 w-3" />
              Conforme
            </span>
          </div>

          <dl class="mt-3 space-y-2 border-t border-border pt-3">
            <div
              v-for="i in infos"
              :key="i.label"
              class="flex items-center justify-between gap-3 text-xs"
            >
              <dt class="flex items-center gap-1.5 text-muted-foreground">
                <component :is="i.icon" class="h-3.5 w-3.5" />
                {{ i.label }}
              </dt>
              <dd class="truncate font-medium text-foreground">
                {{ i.value }}
              </dd>
            </div>
          </dl>
        </section>

        <section class="rounded-[5px] border border-border bg-card p-4">
          <template v-if="candidat.soldeRestant > 0">
            <div class="flex items-baseline justify-between gap-2">
              <h2 class="text-sm font-semibold text-foreground">
                Régler le solde
              </h2>
              <span class="text-sm font-semibold tabular-nums text-foreground">
                {{ fcfa(candidat.soldeRestant) }} FCFA
              </span>
            </div>

            <div class="mt-3 grid grid-cols-3 gap-2">
              <button
                v-for="p in paymentProviders"
                :key="p.label"
                type="button"
                @click="emit('start-payment', p.label)"
                class="flex min-h-[42px] sm:min-h-0 sm:h-9 items-center justify-center gap-1.5 rounded-[3px] border border-border bg-card px-2 text-xs font-medium text-foreground transition hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-dark/40"
              >
                <span
                  class="h-1.5 w-1.5 shrink-0 rounded-full"
                  :class="p.dot"
                />
                <span class="truncate">{{ p.label }}</span>
              </button>
            </div>
          </template>

          <div v-else class="flex items-center gap-2 text-xs text-brand-dark">
            <ShieldCheck class="h-4 w-4 shrink-0 text-brand-dark" />
            <span>Votre dossier est entièrement soldé. Félicitations !</span>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

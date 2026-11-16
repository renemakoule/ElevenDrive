<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useGuichetStore } from "../../stores/guichetStore";
import { useAuthStore } from "../../stores/authStore";
import { realtimeSocket } from "../../services/realtimeSocketService";
import type { Candidat } from "../../types";
import ExportCandidatsButton from "./ExportCandidatsButton.vue";
import {
  UserPlus,
  ChevronRight,
  Phone,
  RefreshCw,
  Search,
  CreditCard,
  CheckCircle2,
  Clock,
  AlertCircle,
  ShieldCheck,
  Car,
  Radio,
} from "lucide-vue-next";

const guichetStore = useGuichetStore();
const authStore = useAuthStore();
const isWsConnected = realtimeSocket.isConnected;

const isRefreshing = ref(false);
const filtreLocal = ref("");

/* ───────────────────────── Initialisation & Rafraîchissement ───────────────────────── */
const refreshData = async () => {
  isRefreshing.value = true;
  try {
    await Promise.all([
      guichetStore.init(),
      guichetStore.search(filtreLocal.value),
    ]);
  } finally {
    isRefreshing.value = false;
  }
};

onMounted(async () => {
  if (
    guichetStore.searchResults.length === 0 ||
    !guichetStore.formules.length
  ) {
    await guichetStore.init();
  }
});

const handleBackendSearch = async () => {
  await guichetStore.search(filtreLocal.value);
};

/* ───────────────────────── Helpers Réels (Données Dynamiques du Backend) ───────────────────────── */
const fcfa = (n: number) => (n || 0).toLocaleString("fr-FR");

const numeroDossier = (c: Candidat) => c.numero_dossier || "En attente";

const formuleLabel = (c: Candidat) => c.formule_nom || "Non assignée";

const reglement = (c: Candidat) => {
  if (c.solde_restant !== undefined && c.solde_restant !== null) {
    if (c.solde_restant <= 0) {
      return {
        label: "Soldé",
        dot: "bg-emerald-500",
        text: "text-emerald-700 dark:text-emerald-400",
        badge:
          "bg-emerald-50 text-emerald-700 border-emerald-200 dark:bg-emerald-950/40 dark:text-emerald-300 dark:border-emerald-800",
      };
    }
    return {
      label: `Reste ${fcfa(c.solde_restant)} F`,
      dot: "bg-amber-500",
      text: "text-amber-700 dark:text-amber-400",
      badge:
        "bg-amber-50 text-amber-700 border-amber-200 dark:bg-amber-950/40 dark:text-amber-300 dark:border-amber-800",
    };
  }
  if (c.statut_dossier) {
    const libelles: Record<string, string> = {
      EN_COURS: "En cours",
      PREPARE_EXAMEN: "Prépa examen",
      ADMIS: "Admis",
      AJOURNE: "Ajourné",
      ABANDONNE: "Abandonné",
    };
    return {
      label: libelles[c.statut_dossier] || c.statut_dossier.replace(/_/g, " "),
      dot: "bg-blue-500",
      text: "text-blue-700 dark:text-blue-400",
      badge:
        "bg-blue-50 text-blue-700 border-blue-200 dark:bg-blue-950/40 dark:text-blue-300 dark:border-blue-800",
    };
  }
  return {
    label: "Dossier actif",
    dot: "bg-emerald-500",
    text: "text-emerald-700 dark:text-emerald-400",
    badge:
      "bg-emerald-50 text-emerald-700 border-emerald-200 dark:bg-emerald-950/40 dark:text-emerald-300 dark:border-emerald-800",
  };
};

const candidats = computed(() => guichetStore.searchResults ?? []);

const dossiersAffiches = computed(() => {
  const query = filtreLocal.value.trim().toLowerCase();
  if (!query) return candidats.value.slice(0, 8);
  return candidats.value
    .filter((c) => {
      const nomComplet = `${c.nom || ""} ${c.prenoms || ""}`.toLowerCase();
      const dos = (c.numero_dossier || "").toLowerCase();
      const tel = (c.telephone_whatsapp || "").toLowerCase();
      const piece = (c.numero_piece_id || "").toLowerCase();
      return (
        nomComplet.includes(query) ||
        dos.includes(query) ||
        tel.includes(query) ||
        piece.includes(query)
      );
    })
    .slice(0, 8);
});

const initiales = (c: { nom?: string; prenoms?: string }) =>
  `${c.nom?.charAt(0) ?? ""}${c.prenoms?.charAt(0) ?? ""}`.toUpperCase() || "?";

/* ───────────────────────── Indicateurs KPIs Réels ───────────────────────── */
const stats = computed(() => [
  {
    key: "dossiers",
    label: "Candidats enregistrés",
    value: String(
      guichetStore.stats?.total_candidats ?? candidats.value.length,
    ),
    caption: "Inscriptions en base",
    live: false,
  },
  {
    key: "caisse",
    label: "Espèces au tiroir",
    value: `${fcfa(guichetStore.stats?.total_especes_theorique ?? 0)} F`,
    caption: `Fond de caisse : ${fcfa(guichetStore.stats?.fond_de_caisse ?? 0)} F`,
    live: true,
  },
  {
    key: "code",
    label: "Formules au catalogue",
    value: String(guichetStore.formules.length),
    caption: "Tarifs officiels actifs",
    live: false,
  },
  {
    key: "sync",
    label: "Connectivité Cloud",
    value:
      guichetStore.stats?.mode_connectivite === "CLOUD_POSTGRES_CONNECTED"
        ? "En ligne"
        : "Local / Prêt",
    caption:
      (guichetStore.stats?.evenements_sync_en_attente ?? 0) > 0
        ? `${guichetStore.stats.evenements_sync_en_attente} sync en attente`
        : "Base de données synchronisée",
    live: guichetStore.stats?.mode_connectivite === "CLOUD_POSTGRES_CONNECTED",
  },
]);

/* ───────────────────────── Actions ───────────────────────── */
const openNewCandidate = () => {
  guichetStore.isNewCandidateModalOpen = true;
};

const goToGuichet = () => {
  authStore.setActiveSection("guichet");
};

const selectCandidate = async (id: string) => {
  await guichetStore.selectCandidat(id);
  authStore.setActiveSection("guichet");
};
</script>

<template>
  <div class="space-y-4 bg-background">
    <!-- ═════════ En-tête ═════════ -->
    <header
      class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
    >
      <div class="min-w-0">
        <div class="flex items-center gap-2">
          <h1 class="text-base font-semibold text-foreground">
            Vue d'ensemble Secrétariat
          </h1>
          <!-- <span
            class="inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-bold"
            :class="
              isWsConnected
                ? 'bg-emerald-500/15 text-emerald-700 dark:text-emerald-300'
                : 'bg-amber-500/15 text-amber-700'
            "
          >
            <Radio class="h-3 w-3 animate-pulse" />
            {{ isWsConnected ? "Temps réel WS" : "Reconnexion WS..." }}
          </span> -->
        </div>
        <p class="text-xs text-muted-foreground">
          Guichet, inscriptions des candidats et tiroir-caisse de l'agence
          synchronisés en direct.
        </p>
      </div>

      <div class="flex flex-wrap items-center gap-1.5 sm:gap-2">
        <button
          type="button"
          @click="refreshData"
          :disabled="isRefreshing"
          class="inline-flex h-9 w-9 sm:h-8 sm:w-8 items-center justify-center rounded-[3px] border border-border bg-card text-muted-foreground transition hover:bg-muted hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40 disabled:opacity-50"
          title="Rafraîchir les données du guichet"
        >
          <RefreshCw
            class="h-3.5 w-3.5"
            :class="{ 'animate-spin': isRefreshing }"
          />
        </button>

        <button
          type="button"
          @click="goToGuichet"
          class="inline-flex h-9 sm:h-8 items-center gap-2 rounded-[3px] border border-border bg-card px-3 text-xs font-medium text-foreground transition hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40"
        >
          Guichet rapide
          <kbd
            class="hidden sm:inline rounded bg-muted px-1 font-mono text-[10px] text-muted-foreground"
            >F2</kbd
          >
        </button>

        <button
          type="button"
          @click="authStore.setActiveSection('conduite-planning')"
          class="inline-flex h-9 sm:h-8 items-center gap-1.5 rounded-[3px] border border-border bg-card px-3 text-xs font-medium text-foreground transition hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40"
        >
          <Car class="h-3.5 w-3.5 text-brand" />
          Planning Conduite
        </button>

        <button
          type="button"
          @click="openNewCandidate"
          class="flex-1 sm:flex-initial inline-flex h-9 sm:h-8 items-center justify-center gap-1.5 rounded-[3px] bg-brand-dark px-3 text-xs font-medium text-white transition hover:bg-brand focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40 focus-visible:ring-offset-1"
        >
          <UserPlus class="h-3.5 w-3.5" />
          Nouvelle inscription
          <kbd
            class="hidden sm:inline rounded bg-white/20 px-1 font-mono text-[10px] text-white/90"
            >F1</kbd
          >
        </button>
      </div>
    </header>

    <!-- ═════════ Indicateurs KPIs Réels (Carrousel horizontal tactile sur mobile) ═════════ -->
    <section
      class="flex overflow-x-auto snap-x snap-mandatory gap-2.5 no-scrollbar sm:grid sm:grid-cols-2 lg:grid-cols-4 sm:gap-0 sm:divide-x sm:divide-border sm:rounded-[3px] sm:border sm:border-border sm:bg-card"
    >
      <div
        v-for="s in stats"
        :key="s.key"
        class="min-w-[74%] sm:min-w-0 snap-start shrink-0 rounded-[3px] border border-border bg-card sm:border-0 sm:rounded-none p-3.5 sm:p-4"
      >
        <span class="text-xs text-muted-foreground">{{ s.label }}</span>
        <p class="mt-1 text-lg font-semibold tabular-nums text-foreground">
          {{ s.value }}
        </p>
        <p
          class="mt-1.5 flex items-center gap-1.5 text-[11px] text-muted-foreground"
        >
          <span
            v-if="s.live"
            class="h-1.5 w-1.5 rounded-full bg-brand animate-pulse"
          />
          {{ s.caption }}
        </p>
      </div>
    </section>

    <!-- ═════════ Dossiers récents ═════════ -->
    <section class="rounded-[3px] border border-border bg-card">
      <div
        class="flex flex-col gap-2.5 px-4 py-3 sm:flex-row sm:items-center sm:justify-between border-b border-border"
      >
        <div class="flex min-w-0 items-baseline gap-2">
          <h2 class="text-sm font-semibold text-foreground">
            Dossiers récents
          </h2>
          <span class="truncate text-xs text-muted-foreground"
            >{{ candidats.length }} candidats enregistrés</span
          >
        </div>

        <div class="flex items-center gap-2">
          <!-- Filtre rapide local / recherche -->
          <div class="relative w-full sm:w-60">
            <Search
              class="pointer-events-none absolute left-2.5 top-1/2 -translate-y-1/2 h-3.5 w-3.5 text-muted-foreground"
            />
            <input
              v-model="filtreLocal"
              type="text"
              placeholder="Filtrer nom, dossier..."
              class="h-7 w-full rounded-[4px] border border-border bg-muted/40 pl-8 pr-2 text-xs text-foreground placeholder:text-muted-foreground focus:border-brand focus:outline-none focus:ring-1 focus:ring-brand/30"
              @keyup.enter="handleBackendSearch"
            />
          </div>

          <ExportCandidatsButton :candidats="candidats" />

          <button
            type="button"
            @click="goToGuichet"
            class="inline-flex shrink-0 items-center gap-0.5 text-xs font-medium text-brand transition hover:text-brand-dark focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40"
          >
            Recherche et caisse
            <kbd
              class="hidden sm:inline ml-1 rounded bg-muted px-1 font-mono text-[10px] text-muted-foreground"
              >F3</kbd
            >
            <ChevronRight class="h-3.5 w-3.5" />
          </button>
        </div>
      </div>

      <template v-if="dossiersAffiches.length">
        <!-- Cartes Mobiles (< md) -->
        <div class="md:hidden divide-y divide-border">
          <button
            v-for="c in dossiersAffiches"
            :key="c.id"
            type="button"
            @click="selectCandidate(c.id)"
            class="w-full text-left px-3.5 py-3 flex items-center justify-between gap-2.5 active:bg-muted/60 transition"
          >
            <div class="flex items-center gap-2.5 min-w-0">
              <div
                class="flex h-9 w-9 shrink-0 items-center justify-center overflow-hidden rounded-full bg-muted text-xs font-semibold text-muted-foreground"
              >
                <img
                  v-if="c.photo_url"
                  :src="c.photo_url"
                  :alt="c.nom"
                  class="h-full w-full object-cover"
                />
                <span v-else>{{ initiales(c) }}</span>
              </div>
              <div class="min-w-0 space-y-0.5">
                <div class="truncate text-xs font-semibold text-foreground">
                  {{ c.nom }} {{ c.prenoms }}
                </div>
                <div class="truncate text-[11px] text-muted-foreground">
                  <span class="font-mono">{{ numeroDossier(c) }}</span> ·
                  {{ formuleLabel(c) }}
                </div>
              </div>
            </div>

            <div class="flex items-center gap-1.5 shrink-0">
              <span
                class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-medium border"
                :class="reglement(c).badge"
              >
                <span
                  class="h-1.5 w-1.5 rounded-full"
                  :class="reglement(c).dot"
                />
                {{ reglement(c).label }}
              </span>
              <ChevronRight class="h-4 w-4 text-muted-foreground" />
            </div>
          </button>
        </div>

        <div class="hidden md:block overflow-x-auto">
          <table class="w-full text-left text-xs">
            <thead>
              <tr
                class="border-b border-border bg-muted/40 text-[11px] text-muted-foreground"
              >
                <th class="px-4 py-2 font-medium">Candidat</th>
                <th class="px-3 py-2 font-medium">Dossier</th>
                <th class="hidden px-3 py-2 font-medium md:table-cell">
                  Téléphone
                </th>
                <th class="px-3 py-2 font-medium">Règlement</th>
                <th class="w-8 px-4 py-2" />
              </tr>
            </thead>

            <tbody class="divide-y divide-border">
              <tr
                v-for="c in dossiersAffiches"
                :key="c.id"
                tabindex="0"
                class="group cursor-pointer transition hover:bg-muted/60 focus-visible:bg-muted/60 focus-visible:outline-none"
                @click="selectCandidate(c.id)"
                @keydown.enter="selectCandidate(c.id)"
              >
                <!-- Candidat -->
                <td class="px-4 py-2.5">
                  <div class="flex items-center gap-2.5">
                    <div
                      class="flex h-8 w-8 shrink-0 items-center justify-center overflow-hidden rounded-full bg-muted text-[11px] font-semibold text-muted-foreground"
                    >
                      <img
                        v-if="c.photo_url"
                        :src="c.photo_url"
                        :alt="c.nom"
                        class="h-full w-full object-cover"
                      />
                      <span v-else>{{ initiales(c) }}</span>
                    </div>
                    <div class="min-w-0">
                      <div class="truncate font-medium text-foreground">
                        {{ c.nom }} {{ c.prenoms }}
                      </div>
                      <div class="truncate text-[11px] text-muted-foreground">
                        {{ c.type_piece_id || "Pièce" }}
                        {{ c.numero_piece_id || "-" }}
                      </div>
                    </div>
                  </div>
                </td>

                <!-- Dossier + formule -->
                <td class="px-3 py-2.5">
                  <div class="font-mono font-medium text-foreground">
                    {{ numeroDossier(c) }}
                  </div>
                  <div
                    class="truncate text-[11px] text-muted-foreground max-w-[200px]"
                    :title="formuleLabel(c)"
                  >
                    {{ formuleLabel(c) }}
                  </div>
                </td>

                <!-- Téléphone -->
                <td class="hidden px-3 py-2.5 md:table-cell">
                  <span
                    class="inline-flex items-center gap-1.5 tabular-nums text-foreground"
                  >
                    <Phone class="h-3 w-3 text-muted-foreground" />
                    {{ c.telephone_whatsapp || "-" }}
                  </span>
                </td>

                <!-- Règlement -->
                <td class="px-3 py-2.5">
                  <span
                    class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-[11px] font-medium border"
                    :class="reglement(c).badge"
                  >
                    <span
                      class="h-1.5 w-1.5 rounded-full"
                      :class="reglement(c).dot"
                    />
                    {{ reglement(c).label }}
                  </span>
                </td>

                <!-- Ouvrir -->
                <td class="px-4 py-2.5 text-right">
                  <ChevronRight
                    class="inline h-4 w-4 text-muted-foreground transition group-hover:translate-x-0.5 group-hover:text-brand"
                  />
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </template>

      <!-- État vide -->
      <div
        v-else
        class="flex flex-col items-center gap-3 px-4 py-10 text-center"
      >
        <p class="text-xs text-muted-foreground">
          {{
            filtreLocal
              ? "Aucun résultat trouvé pour cette recherche."
              : "Aucun dossier enregistré pour le moment."
          }}
        </p>
        <button
          v-if="!filtreLocal"
          type="button"
          @click="openNewCandidate"
          class="inline-flex h-8 items-center gap-1.5 rounded-[4px] bg-brand-dark px-3 text-xs font-medium text-white transition hover:bg-brand"
        >
          <UserPlus class="h-3.5 w-3.5" />
          Nouvelle inscription
        </button>
        <button
          v-else
          type="button"
          @click="
            filtreLocal = '';
            refreshData();
          "
          class="inline-flex h-8 items-center gap-1.5 rounded-[4px] border border-border bg-card px-3 text-xs font-medium text-foreground transition hover:bg-muted"
        >
          Réinitialiser le filtre
        </button>
      </div>
    </section>
  </div>
</template>

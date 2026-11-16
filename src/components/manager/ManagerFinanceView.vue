<script setup lang="ts">
import { computed, ref, onMounted } from "vue";
import { useAuthStore } from "../../stores/authStore";
import { useGuichetStore } from "../../stores/guichetStore";
import { realtimeSocket } from "../../services/realtimeSocketService";
import {
  ArrowLeft,
  Download,
  Phone,
  RefreshCw,
  CheckCircle2,
  Radio,
} from "lucide-vue-next";

const authStore = useAuthStore();
const guichetStore = useGuichetStore();
const isWsConnected = realtimeSocket.isConnected;

const fmt = (n: number) => (n || 0).toLocaleString("fr-FR");

onMounted(async () => {
  await guichetStore.loadGerantDashboard();
});

interface ModePaiementStat {
  label: string;
  montant: number;
  dot: string;
}

interface ImpayeItem {
  id: string;
  candidat: string;
  dossier: string;
  formule: string;
  montantTotal: number;
  dejaPaye: number;
  solde: number;
  tauxReglement: number;
  telephone: string;
  retard: string;
}

const modes = computed<ModePaiementStat[]>(() => {
  const v = guichetStore.gerantDashboardData?.ventilation;
  const especes = v
    ? v.especes
    : guichetStore.stats.total_especes_theorique || 0;
  const wave = v ? v.wave : 0;
  const om = v ? v.orange_money : 0;
  const mtn = v ? v.mtn_momo : 0;
  const moov = v ? v.moov : 0;
  const vir = v ? v.virement : 0;

  const base: ModePaiementStat[] = [
    { label: "Espèces", montant: especes, dot: "bg-brand" },
    { label: "Wave", montant: wave, dot: "bg-sky-500" },
    { label: "Orange Money", montant: om, dot: "bg-amber-500" },
    { label: "MTN MoMo", montant: mtn, dot: "bg-yellow-500" },
  ];

  if (moov > 0) {
    base.push({ label: "Moov Money", montant: moov, dot: "bg-blue-600" });
  }
  if (vir > 0) {
    base.push({ label: "Virement", montant: vir, dot: "bg-purple-600" });
  }

  return base;
});

const totalEncaisse = computed(() =>
  modes.value.reduce((s: number, m: ModePaiementStat) => s + m.montant, 0),
);

const repartition = computed(() =>
  modes.value.map((m) => ({
    ...m,
    part:
      totalEncaisse.value > 0
        ? ((m.montant / totalEncaisse.value) * 100).toLocaleString("fr-FR", {
            maximumFractionDigits: 1,
          })
        : "0",
  })),
);

const calculerRetard = (dateIso?: string) => {
  if (!dateIso) return "En cours";
  const dateInsc = new Date(dateIso);
  if (isNaN(dateInsc.getTime())) return "En cours";
  const now = new Date();
  const diffDays = Math.floor(
    (now.getTime() - dateInsc.getTime()) / (1000 * 60 * 60 * 24),
  );
  if (diffDays <= 0) return "Aujourd'hui";
  if (diffDays === 1) return "1 jour";
  if (diffDays < 30) return `${diffDays} j`;
  const mois = Math.floor(diffDays / 30);
  const resteJours = diffDays % 30;
  if (resteJours === 0) return `${mois} mois`;
  return `${mois} m ${resteJours} j`;
};

const impayes = computed<ImpayeItem[]>(() => {
  const raw = guichetStore.gerantDashboardData?.impayes || [];
  return raw.map(
    (d: any, idx: number): ImpayeItem => ({
      id: d.dossier_id || `deb-${idx}`,
      candidat: d.candidat_nom || "Candidat",
      dossier: d.numero_dossier || "Dossier",
      formule: d.formule || "Permis",
      montantTotal: Number(d.montant_total) || 0,
      dejaPaye: Number(d.deja_paye) || 0,
      solde: Number(d.solde_restant) || 0,
      tauxReglement: Number(d.taux_reglement_pct) || 0,
      telephone: d.telephone_whatsapp || "",
      retard: calculerRetard(d.date_inscription),
    }),
  );
});

const totalImpayes = computed(() =>
  impayes.value.reduce((s: number, d: ImpayeItem) => s + d.solde, 0),
);

const relancer = (d: ImpayeItem) => {
  const msg = encodeURIComponent(
    `Bonjour ${d.candidat}, l'établissement ${authStore.schoolName} vous informe qu'un solde de ${fmt(d.solde)} FCFA reste à régler sur votre dossier ${d.dossier}, afin de planifier vos leçons.`,
  );
  window.open(
    `https://wa.me/${d.telephone.replace(/\D/g, "")}?text=${msg}`,
    "_blank",
  );
};

const exporterCsv = () => {
  const lignes = [
    [
      "Section",
      "Libellé",
      "Dossier",
      "Montant Total",
      "Déjà Payé",
      "Solde Restant",
      "Retard",
    ],
    ...modes.value.map((m) => [
      "Encaissements",
      m.label,
      "",
      "",
      "",
      String(m.montant),
      "",
    ]),
    ...impayes.value.map((d) => [
      "Impayés",
      d.candidat,
      d.dossier,
      String(d.montantTotal),
      String(d.dejaPaye),
      String(d.solde),
      d.retard,
    ]),
  ];
  const contenu = "\uFEFF" + lignes.map((l) => l.join(";")).join("\n");
  const lien = document.createElement("a");
  lien.href = URL.createObjectURL(
    new Blob([contenu], { type: "text/csv;charset=utf-8" }),
  );
  lien.download = `audit-financier-${new Date().toISOString().slice(0, 10)}.csv`;
  lien.click();
  URL.revokeObjectURL(lien.href);
};
</script>

<template>
  <div class="space-y-4">
    <!-- En-tête -->
    <div
      class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
    >
      <div class="flex items-center gap-2">
        <button
          type="button"
          @click="authStore.setActiveSection('overview')"
          class="min-h-[38px] min-w-[38px] sm:min-h-0 sm:min-w-0 inline-flex items-center justify-center rounded-[3px] border border-border bg-card p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground transition"
          aria-label="Retour"
        >
          <ArrowLeft class="h-4 w-4" />
        </button>
        <div>
          <h1 class="text-sm font-semibold leading-tight text-foreground">
            Audit financier & Caisse
          </h1>
          <p class="text-xs text-muted-foreground">
            Ventilation réelle des encaissements et recouvrement des créances
          </p>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <button
          type="button"
          @click="guichetStore.loadGerantDashboard()"
          class="min-h-[38px] sm:min-h-0 flex items-center gap-1.5 rounded-[3px] border border-border bg-card px-2.5 py-1.5 text-xs font-medium text-muted-foreground hover:bg-muted hover:text-foreground transition"
          title="Actualiser les données financières"
        >
          <RefreshCw
            class="h-3.5 w-3.5"
            :class="{ 'animate-spin': guichetStore.loading }"
          />
          Actualiser
        </button>
        <button
          type="button"
          @click="exporterCsv"
          class="flex-1 sm:flex-initial min-h-[38px] sm:min-h-0 flex items-center justify-center gap-1.5 rounded-[3px] border border-border bg-card px-3 py-1.5 text-xs font-medium text-foreground hover:bg-muted transition"
        >
          <Download class="h-3.5 w-3.5 text-muted-foreground" />
          Exporter en CSV
        </button>
      </div>
    </div>

    <!-- Ventilation Réelle par Mode de Paiement -->
    <div class="grid grid-cols-2 gap-2.5 md:grid-cols-4">
      <div
        v-for="m in repartition"
        :key="m.label"
        class="rounded-[3px] border border-border bg-card px-3.5 py-3"
      >
        <div class="flex items-center gap-1.5 text-xs text-muted-foreground">
          <span class="h-2 w-2 rounded-full" :class="m.dot"></span>
          {{ m.label }}
        </div>
        <div class="mt-1 flex items-baseline justify-between gap-1">
          <span
            class="text-sm sm:text-base font-semibold text-foreground tabular-nums truncate"
            >{{ fmt(m.montant) }} F</span
          >
          <span class="text-xs text-muted-foreground tabular-nums shrink-0"
            >{{ m.part }} %</span
          >
        </div>
      </div>
    </div>

    <!-- Tableau de Recouvrement des Impayés -->
    <div class="overflow-hidden rounded-[3px] border border-border bg-card">
      <div
        class="flex flex-wrap items-center justify-between gap-2 border-b border-border px-4 py-3"
      >
        <div>
          <span class="text-xs font-semibold text-foreground"
            >Créances à recouvrer</span
          >
          <span class="ml-1.5 text-xs text-muted-foreground">
            {{ impayes.length }} dossiers débiteurs
          </span>
        </div>
        <div
          class="text-xs font-semibold text-amber-600 dark:text-amber-400 tabular-nums"
        >
          Total : {{ fmt(totalImpayes) }} FCFA
        </div>
      </div>

      <!-- Cartes Mobiles (< md) -->
      <div class="md:hidden divide-y divide-border">
        <div
          v-if="impayes.length === 0 && !guichetStore.loading"
          class="px-4 py-8 text-center text-xs text-muted-foreground flex flex-col items-center gap-1.5"
        >
          <CheckCircle2 class="h-6 w-6 text-emerald-500" />
          <span class="font-medium text-foreground"
            >Aucun impayé en attente</span
          >
          <span
            >Toutes les souscriptions sont intégralement soldées ou à
            jour.</span
          >
        </div>
        <div v-for="d in impayes" :key="d.id" class="p-3.5 space-y-2.5 text-xs">
          <div class="flex items-start justify-between gap-2">
            <div class="min-w-0">
              <div class="font-semibold text-foreground truncate">
                {{ d.candidat }}
              </div>
              <div class="text-[11px] text-muted-foreground truncate">
                <span class="font-mono">{{ d.dossier }}</span> · {{ d.formule }}
              </div>
            </div>
            <div class="text-right shrink-0">
              <div
                class="font-mono font-bold text-rose-600 dark:text-rose-400 tabular-nums"
              >
                -{{ fmt(d.solde) }} F
              </div>
              <div
                class="text-[10px] font-medium text-amber-600 dark:text-amber-400"
              >
                {{ d.retard }}
              </div>
            </div>
          </div>
          <div class="flex items-center justify-between gap-2 pt-1">
            <span class="text-[11px] text-muted-foreground tabular-nums">
              Versé :
              <strong class="text-emerald-600 dark:text-emerald-400"
                >{{ fmt(d.dejaPaye) }} F</strong
              >
              / {{ fmt(d.montantTotal) }} F
            </span>
            <button
              type="button"
              @click="relancer(d)"
              :disabled="!d.telephone"
              class="min-h-[38px] inline-flex items-center gap-1.5 rounded-[3px] border border-border bg-background px-3 text-xs font-medium text-foreground hover:bg-muted transition disabled:opacity-40"
            >
              <Phone class="h-3.5 w-3.5 text-brand-dark" />
              Relancer
            </button>
          </div>
        </div>
      </div>

      <div class="hidden md:block overflow-x-auto">
        <table class="w-full text-left text-xs">
          <thead class="bg-muted text-muted-foreground">
            <tr>
              <th class="px-3 py-2 font-medium">Candidat</th>
              <th class="px-3 py-2 font-medium">Dossier</th>
              <th class="px-3 py-2 font-medium">Formule</th>
              <th class="px-3 py-2 text-right font-medium">Forfait</th>
              <th class="px-3 py-2 text-right font-medium">Versé</th>
              <th class="px-3 py-2 text-right font-medium">Solde dû</th>
              <th class="px-3 py-2 font-medium">Ancienneté</th>
              <th class="px-3 py-2 text-right"></th>
            </tr>
          </thead>
          <tbody class="divide-y divide-border">
            <tr v-if="impayes.length === 0 && !guichetStore.loading">
              <td
                colspan="8"
                class="px-4 py-8 text-center text-xs text-muted-foreground"
              >
                <div class="flex flex-col items-center justify-center gap-1.5">
                  <CheckCircle2 class="h-6 w-6 text-emerald-500" />
                  <span class="font-medium text-foreground"
                    >Aucun impayé en attente</span
                  >
                  <span
                    >Toutes les souscriptions sont intégralement soldées ou à
                    jour.</span
                  >
                </div>
              </td>
            </tr>
            <tr
              v-for="d in impayes"
              :key="d.id"
              class="hover:bg-muted/50 transition"
            >
              <td class="px-3 py-2 font-medium text-foreground">
                {{ d.candidat }}
              </td>
              <td class="px-3 py-2 font-mono text-muted-foreground">
                {{ d.dossier }}
              </td>
              <td class="px-3 py-2 text-muted-foreground">{{ d.formule }}</td>
              <td
                class="px-3 py-2 text-right font-mono text-muted-foreground tabular-nums"
              >
                {{ fmt(d.montantTotal) }} F
              </td>
              <td
                class="px-3 py-2 text-right font-mono text-emerald-600 dark:text-emerald-400 tabular-nums"
              >
                {{ fmt(d.dejaPaye) }} F
              </td>
              <td
                class="px-3 py-2 text-right font-mono font-semibold text-rose-600 dark:text-rose-400 tabular-nums"
              >
                {{ fmt(d.solde) }} F
              </td>
              <td
                class="px-3 py-2 text-xs font-medium text-amber-600 dark:text-amber-400 tabular-nums"
              >
                {{ d.retard }}
              </td>
              <td class="px-3 py-2 text-right">
                <button
                  type="button"
                  @click="relancer(d)"
                  :disabled="!d.telephone"
                  class="inline-flex items-center gap-1 rounded-lg border border-border bg-card px-2.5 py-1 text-xs font-medium text-foreground hover:bg-muted transition disabled:opacity-40"
                  title="Envoyer un rappel sur WhatsApp"
                >
                  <Phone class="h-3 w-3 text-brand-dark" />
                  Relancer
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

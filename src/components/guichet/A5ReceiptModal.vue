<script setup lang="ts">
import { computed } from "vue";
import { useGuichetStore } from "../../stores/guichetStore";
import { useAuthStore } from "../../stores/authStore";
import { montantEnLettresFCFA } from "../../utils/numberToWords";
import { X, Printer } from "lucide-vue-next";

const emit = defineEmits<{
  (e: "close"): void;
}>();

const store = useGuichetStore();
const authStore = useAuthStore();

/* ───────────────────────── Helpers ───────────────────────── */
const fcfa = (n: number) => n.toLocaleString("fr-FR");

const schoolMeta = computed(() => ({
  nomCommercial: authStore.currentSchool?.nom_commercial || "AutoSchool",
  numeroAgrement:
    authStore.currentSchool?.numero_agrement || "Agrément Ministériel",
  agence: authStore.currentAgence?.nom || "Agence Guichet",
  adresse: authStore.currentAgence?.adresse || "Siège social",
  telephone:
    authStore.currentAgence?.telephone ||
    authStore.currentSchool?.telephone ||
    "",
}));

const modeLabels: Record<string, string> = {
  ESPECES: "Espèces",
  WAVE: "Wave",
  ORANGE_MONEY: "Orange Money",
  MTN_MOMO: "MTN MoMo",
};
const modeLabel = (m: string) => modeLabels[m] ?? m;

const receipt = computed(() => store.lastReceipt);

const dateHeure = computed(() => {
  if (!receipt.value) return "";
  const d = new Date(receipt.value.date_paiement);
  const date = d.toLocaleDateString("fr-FR");
  const heure = d.toLocaleTimeString("fr-FR", {
    hour: "2-digit",
    minute: "2-digit",
  });
  return `${date} à ${heure}`;
});

const identite = computed(() => {
  const r = receipt.value;
  if (!r) return [];
  return [
    { label: "Dossier", value: r.numero_dossier },
    { label: "Candidat", value: r.candidat_nom },
    { label: "Prestation", value: r.formule_nom },
    { label: "Règlement", value: modeLabel(r.mode_reglement) },
  ];
});

const printA5 = () => {
  window.print();
};
</script>

<template>
  <div
    v-if="receipt"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 backdrop-blur-sm"
    role="dialog"
    aria-modal="true"
    aria-labelledby="receipt-title"
  >
    <div
      class="flex max-h-[94vh] w-full max-w-xl flex-col overflow-hidden rounded-[4px] border border-border bg-card"
    >
      <!-- ═════════ En-tête ═════════ -->
      <header
        class="flex items-center justify-between gap-3 border-b border-border px-4 py-3"
      >
        <div class="min-w-0">
          <h2 id="receipt-title" class="text-sm font-semibold text-foreground">
            Reçu officiel
          </h2>
          <p class="text-[11px] text-muted-foreground">
            Format A5 · exemplaire entreprise et candidat
          </p>
        </div>
        <button
          type="button"
          @click="emit('close')"
          class="flex h-7 w-7 items-center justify-center rounded-[4px] text-muted-foreground transition hover:bg-muted hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-dark/40"
          aria-label="Fermer"
        >
          <X class="h-4 w-4" />
        </button>
      </header>

      <!-- ═════════ Aperçu du document ═════════ -->
      <div class="flex-1 overflow-y-auto bg-muted/40 p-4">
        <!-- L'id est conservé : il peut être ciblé par les styles d'impression -->
        <article
          id="contract-a4-container"
          class="relative select-text space-y-4 rounded-[4px] border border-slate-200 bg-white p-5 text-[11px] leading-normal text-slate-900"
        >
          <!-- Filigrane duplicata -->
          <div
            v-if="receipt.est_duplicata"
            class="pointer-events-none absolute inset-0 flex items-center justify-center overflow-hidden opacity-10"
          >
            <span
              class="-rotate-[25deg] border-4 border-red-600 px-6 py-1.5 text-4xl font-black text-red-600"
            >
              DUPLICATA N° {{ receipt.nombre_impressions || 2 }}
            </span>
          </div>

          <!-- Émetteur & numéro -->
          <header
            class="flex items-start justify-between gap-4 border-b border-slate-900 pb-3"
          >
            <div class="space-y-0.5">
              <div class="text-sm font-bold tracking-tight uppercase">
                {{ schoolMeta.nomCommercial }}
              </div>
              <div class="text-[10px] font-medium text-slate-700">
                Agrément ministériel : {{ schoolMeta.numeroAgrement }}
              </div>
              <div class="text-[10px] text-slate-500">
                {{ schoolMeta.agence }} · {{ schoolMeta.adresse }}
              </div>
              <div
                v-if="schoolMeta.telephone"
                class="text-[10px] text-slate-500"
              >
                Tél. {{ schoolMeta.telephone }}
              </div>
            </div>

            <div class="shrink-0 text-right">
              <div
                class="inline-block rounded border border-slate-900 px-2 py-0.5 font-mono text-[11px] font-semibold"
              >
                Reçu N° {{ receipt.numero_recu }}
              </div>
              <div class="mt-1 text-[10px] text-slate-500">{{ dateHeure }}</div>
            </div>
          </header>

          <h3
            class="text-center text-[11px] font-semibold uppercase tracking-wider"
          >
            Reçu d'encaissement · Attestation de paiement
          </h3>

          <!-- Identité -->
          <dl class="grid grid-cols-2 gap-x-6 gap-y-1.5">
            <div v-for="i in identite" :key="i.label" class="flex gap-2">
              <dt class="w-16 shrink-0 text-slate-500">{{ i.label }}</dt>
              <dd class="min-w-0 flex-1 font-medium">{{ i.value }}</dd>
            </div>
          </dl>

          <!-- Détail du règlement -->
          <table class="w-full border-collapse text-left">
            <thead>
              <tr
                class="border-y border-slate-300 text-[10px] uppercase tracking-wide text-slate-500"
              >
                <th class="py-1.5 pr-2 font-medium">Désignation</th>
                <th class="px-2 py-1.5 text-center font-medium">Mode</th>
                <th class="py-1.5 pl-2 text-right font-medium">
                  Montant encaissé
                </th>
              </tr>
            </thead>
            <tbody>
              <tr class="border-b border-slate-200 align-top">
                <td class="py-2 pr-2">
                  Versement acompte · formation conduite automobile
                  <span
                    class="mt-0.5 block break-all font-mono text-[9px] text-slate-400"
                  >
                    Sécurité : {{ receipt.code_securite_hash }}
                  </span>
                </td>
                <td class="px-2 py-2 text-center">
                  {{ modeLabel(receipt.mode_reglement) }}
                </td>
                <td
                  class="whitespace-nowrap py-2 pl-2 text-right text-xs font-bold tabular-nums"
                >
                  {{ fcfa(receipt.montant_verse) }} FCFA
                </td>
              </tr>
            </tbody>
          </table>

          <!-- Montant en lettres -->
          <div class="rounded border border-slate-200 bg-slate-50 px-3 py-2">
            <span class="block text-[10px] text-slate-500"
              >Arrêté le présent reçu à la somme de</span
            >
            <strong class="font-serif text-xs italic">{{
              montantEnLettresFCFA(receipt.montant_verse)
            }}</strong>
          </div>

          <!-- Situation du dossier & signature -->
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-1.5">
              <div class="text-[10px] uppercase tracking-wide text-slate-500">
                Situation du dossier
              </div>
              <div class="flex items-baseline justify-between gap-2">
                <span>Reste à payer</span>
                <strong
                  class="tabular-nums"
                  :class="
                    receipt.solde_restant > 0
                      ? 'text-red-700'
                      : 'text-emerald-700'
                  "
                >
                  {{ fcfa(receipt.solde_restant) }} FCFA
                </strong>
              </div>
              <p class="text-[10px] text-slate-500">
                {{
                  receipt.solde_restant === 0
                    ? "Dossier entièrement soldé"
                    : "Solde restant exigible"
                }}
              </p>
            </div>

            <div
              class="flex h-24 flex-col justify-between rounded border border-slate-200 p-2.5 text-right"
            >
              <span class="text-[10px] uppercase tracking-wide text-slate-500"
                >L'auto-école · la caissière</span
              >
              <span
                class="inline-block self-end rounded border border-emerald-700 px-1.5 py-0.5 text-[9px] font-semibold text-emerald-800"
              >
                TAMPON OFFICIEL · ENCAISSÉ
              </span>
              <span class="font-mono text-[10px] text-slate-700">{{
                receipt.caissier_nom
              }}</span>
            </div>
          </div>
        </article>
      </div>

      <!-- ═════════ Actions ═════════ -->
      <footer
        class="flex items-center justify-between gap-2 border-t border-border px-4 py-3"
      >
        <button
          type="button"
          @click="emit('close')"
          class="inline-flex h-8 items-center gap-2 rounded-[4px] border border-border bg-card px-3 text-xs font-medium text-foreground transition hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-dark/40"
        >
          Fermer
          <kbd
            class="rounded bg-muted px-1 font-mono text-[10px] text-muted-foreground"
            >Échap</kbd
          >
        </button>

        <button
          type="button"
          @click="printA5"
          class="inline-flex h-8 items-center gap-1.5 rounded-[4px] bg-brand-dark px-3 text-xs font-medium text-white transition hover:bg-brand focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-dark/40 focus-visible:ring-offset-1"
        >
          <Printer class="h-3.5 w-3.5" />
          Imprimer le reçu A5
          <kbd
            class="rounded bg-white/20 px-1 font-mono text-[10px] text-white/90"
            >Entrée</kbd
          >
        </button>
      </footer>
    </div>
  </div>
</template>

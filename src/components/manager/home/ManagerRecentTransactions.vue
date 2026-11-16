<script setup lang="ts">
import { computed } from "vue";
import { useAuthStore } from "../../../stores/authStore";
import { useGuichetStore } from "../../../stores/guichetStore";
import { useLanguageStore } from "../../../stores/languageStore";
import { ChevronRight } from "lucide-vue-next";

const authStore = useAuthStore();
const guichetStore = useGuichetStore();
const languageStore = useLanguageStore();

const t = (txt: string) => languageStore.t(txt);
const fmt = (n: number) => (n || 0).toLocaleString("fr-FR");

const libelleMode: Record<string, string> = {
  WAVE: "Wave",
  ESPECES: "Espèces",
  ORANGE_MONEY: "Orange Money",
  MTN_MOMO: "MTN MoMo",
  MOOV: "Moov Money",
  VIREMENT: "Virement",
};

const pointMode: Record<string, string> = {
  WAVE: "bg-sky-500",
  ESPECES: "bg-transparent",
  ORANGE_MONEY: "bg-amber-500",
  MTN_MOMO: "bg-yellow-500",
  MOOV: "bg-blue-600",
  VIREMENT: "bg-purple-600",
};

const transactions = computed(() => {
  const paiements = guichetStore.gerantDashboardData?.derniers_paiements;
  if (!paiements || paiements.length === 0) return [];
  return paiements.map((p) => ({
    id: p.id,
    dossier: p.numero_dossier,
    candidat: p.candidat_nom,
    formule: p.formule_nom,
    montant: p.montant,
    mode: p.mode_reglement,
    date: p.date_paiement
      ? new Date(p.date_paiement).toLocaleDateString("fr-FR", {
          day: "2-digit",
          month: "2-digit",
          hour: "2-digit",
          minute: "2-digit",
        })
      : "Aujourd'hui",
    recu: p.numero_recu || "Reçu caisse",
  }));
});
</script>

<template>
  <div class="overflow-hidden rounded-[3px] border border-border bg-card">
    <div
      class="flex flex-wrap items-center justify-between gap-2 border-b border-border px-4 py-3"
    >
      <div class="flex items-center gap-2">
        <span class="text-xs font-semibold text-foreground">{{
          t("Derniers encaissements enregistrés")
        }}</span>
        <span
          class="rounded-full bg-brand/10 px-2 py-0.5 text-[10px] font-semibold text-brand"
          >{{ t("Base de données") }}</span
        >
      </div>
      <button
        type="button"
        @click="authStore.setActiveSection('finance')"
        class="flex items-center gap-0.5 text-xs font-medium text-brand-dark hover:underline"
      >
        {{ t("Audit complet") }}
        <ChevronRight class="h-3.5 w-3.5" />
      </button>
    </div>

    <!-- Vue Cartes Mobile (< md) -->
    <div class="md:hidden divide-y divide-border">
      <div
        v-if="transactions.length === 0"
        class="px-4 py-6 text-center text-xs text-muted-foreground"
      >
        {{
          t(
            "Aucun encaissement enregistré pour le moment. Dès qu'un paiement est saisi au guichet, il apparaîtra ici en temps réel.",
          )
        }}
      </div>
      <div
        v-for="tx in transactions"
        :key="tx.id"
        class="px-3.5 py-3 flex items-center justify-between gap-3 text-xs"
      >
        <div class="min-w-0 space-y-0.5">
          <div class="font-medium text-foreground truncate">
            {{ tx.candidat }}
          </div>
          <div
            class="flex flex-wrap items-center gap-1.5 text-[11px] text-muted-foreground"
          >
            <span class="font-mono">{{ tx.recu }}</span>
            <span>·</span>
            <span
              class="inline-flex items-center gap-1 text-foreground font-medium"
            >
              <span
                class="h-1.5 w-1.5 rounded-full"
                :class="pointMode[tx.mode] || 'bg-slate-400'"
              />
              {{ libelleMode[tx.mode] || tx.mode }}
            </span>
          </div>
          <div class="text-[10px] text-muted-foreground">
            {{ tx.formule }} · {{ tx.date }}
          </div>
        </div>
        <div
          class="shrink-0 text-right font-mono font-bold text-foreground tabular-nums"
        >
          +{{ fmt(tx.montant) }} F
        </div>
      </div>
    </div>

    <div class="hidden md:block overflow-x-auto">
      <table class="w-full text-left text-xs">
        <thead class="bg-muted text-muted-foreground">
          <tr>
            <th class="px-3 py-2 font-medium">{{ t("Reçu N°") }}</th>
            <th class="px-3 py-2 font-medium">{{ t("Candidat") }}</th>
            <th class="px-3 py-2 font-medium">{{ t("Formule") }}</th>
            <th class="px-3 py-2 font-medium">{{ t("Mode") }}</th>
            <th class="px-3 py-2 font-medium">{{ t("Date & Heure") }}</th>
            <th class="px-3 py-2 text-right font-medium">
              {{ t("Montant versé") }}
            </th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border">
          <tr v-if="transactions.length === 0">
            <td colspan="6" class="px-3 py-6 text-center text-muted-foreground">
              {{
                t(
                  "Aucun encaissement enregistré pour le moment. Dès qu'un paiement est saisi au guichet, il apparaîtra ici en temps réel.",
                )
              }}
            </td>
          </tr>
          <tr
            v-for="tx in transactions"
            :key="tx.id"
            class="hover:bg-muted/50 transition"
          >
            <td class="px-3 py-2 font-mono text-[11px] text-muted-foreground">
              {{ tx.recu }}
            </td>
            <td class="px-3 py-2">
              <div class="font-medium text-foreground">{{ tx.candidat }}</div>
              <div class="font-mono text-[10px] text-muted-foreground">
                {{ tx.dossier }}
              </div>
            </td>
            <td class="px-3 py-2 text-muted-foreground">{{ tx.formule }}</td>
            <td class="px-3 py-2">
              <span
                class="inline-flex items-center gap-1.5 font-medium text-foreground"
              >
                <span
                  class="h-1.5 w-1.5 rounded-full"
                  :class="pointMode[tx.mode] || 'bg-slate-400'"
                ></span>
                {{ libelleMode[tx.mode] || tx.mode }}
              </span>
            </td>
            <td class="px-3 py-2 text-muted-foreground">{{ tx.date }}</td>
            <td
              class="px-3 py-2 text-right font-mono font-bold text-foreground tabular-nums"
            >
              +{{ fmt(tx.montant) }} FCFA
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

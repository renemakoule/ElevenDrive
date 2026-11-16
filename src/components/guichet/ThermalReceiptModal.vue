<script setup lang="ts">
import { computed, ref } from "vue";
import { useGuichetStore } from "../../stores/guichetStore";
import { useAuthStore } from "../../stores/authStore";
import { montantEnLettresFCFA } from "../../utils/numberToWords";
import A5ReceiptModal from "./A5ReceiptModal.vue";
import { Printer, X, FileText } from "lucide-vue-next";

const store = useGuichetStore();
const authStore = useAuthStore();
const isA5ModalOpen = ref(false);

const close = () => {
  store.isReceiptModalOpen = false;
};

const printTicket = () => window.print();

const fmt = (n: number) => n.toLocaleString("fr-FR");

const recu = computed(() => store.lastReceipt);

const dateTexte = computed(() => {
  if (!recu.value) return "";
  const d = new Date(recu.value.date_paiement);
  const jour = d.toLocaleDateString("fr-FR");
  const heure = d.toLocaleTimeString("fr-FR", {
    hour: "2-digit",
    minute: "2-digit",
  });
  return `${jour} à ${heure}`;
});

const solde = computed(() => recu.value?.solde_restant ?? 0);

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
</script>

<template>
  <div
    v-if="store.isReceiptModalOpen && recu"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 select-none backdrop-blur-sm"
  >
    <div
      class="flex max-h-[95vh] w-full max-w-xl flex-col overflow-hidden rounded-[5px] border border-border bg-card"
    >
      <!-- En-tête -->
      <header
        class="flex items-center justify-between border-b border-border px-4 py-3"
      >
        <div class="flex items-center gap-2 text-foreground">
          <Printer class="h-4 w-4 text-brand-dark" />
          <span class="text-xs font-semibold">Reçu de paiement</span>
        </div>
        <button
          @click="close"
          class="rounded p-1 text-muted-foreground hover:bg-muted hover:text-foreground transition"
          aria-label="Fermer"
        >
          <X class="h-4 w-4" />
        </button>
      </header>

      <!-- Reçu -->
      <div class="flex justify-center overflow-y-auto bg-muted/40 p-4">
        <div
          id="thermal-ticket-container"
          class="w-full max-w-md select-text rounded-[5px] border border-slate-300 bg-white p-6 font-sans text-[10px] leading-snug text-gray-900"
        >
          <!-- Identité et numéro -->
          <div class="flex items-start justify-between gap-4">
            <div>
              <div
                class="text-[13px] font-bold leading-tight tracking-tight uppercase"
              >
                {{ schoolMeta.nomCommercial }}
              </div>
              <div class="mt-1 text-[8.5px] font-medium text-gray-600">
                Agrément ministériel : {{ schoolMeta.numeroAgrement }}
              </div>
              <div class="mt-0.5 text-[8.5px] text-gray-500">
                {{ schoolMeta.agence }} · {{ schoolMeta.adresse }}
              </div>
              <div
                v-if="schoolMeta.telephone"
                class="text-[8.5px] text-gray-500"
              >
                Tél. {{ schoolMeta.telephone }}
              </div>
            </div>

            <div class="shrink-0 text-right">
              <div
                class="inline-block rounded-[3px] border-[1.5px] border-gray-900 px-2 py-1 font-mono text-[9px] font-medium tracking-wide"
              >
                Reçu N° {{ recu.numero_recu }}
              </div>
              <div class="mt-1.5 text-[8.5px] text-gray-500">
                {{ dateTexte }}
              </div>
            </div>
          </div>

          <div class="mt-3 border-t border-gray-900"></div>

          <!-- Titre -->
          <div class="py-3 text-center">
            <div class="text-[9px] font-semibold uppercase tracking-[0.12em]">
              Reçu d'encaissement · Attestation de paiement
            </div>
            <div
              v-if="recu.est_duplicata"
              class="mx-auto mt-1.5 inline-block rounded-[3px] border border-dashed border-red-600 px-2 py-0.5 text-[8px] font-semibold text-red-700"
            >
              Duplicata n° {{ recu.nombre_impressions || 2 }} · réimpression
              certifiée conforme
            </div>
          </div>

          <!-- Informations -->
          <div class="grid grid-cols-2 gap-x-6 gap-y-1.5">
            <div class="flex gap-3">
              <span class="w-14 shrink-0 text-gray-500">Dossier</span>
              <span class="font-medium">{{ recu.numero_dossier }}</span>
            </div>
            <div class="flex gap-3">
              <span class="w-14 shrink-0 text-gray-500">Candidat</span>
              <span class="font-medium">{{ recu.candidat_nom }}</span>
            </div>
            <div class="flex gap-3">
              <span class="w-14 shrink-0 text-gray-500">Prestation</span>
              <span class="font-medium">{{ recu.formule_nom }}</span>
            </div>
            <div class="flex gap-3">
              <span class="w-14 shrink-0 text-gray-500">Règlement</span>
              <span class="font-medium">{{ recu.mode_reglement }}</span>
            </div>
          </div>

          <!-- Désignation -->
          <div class="mt-4 border-y border-gray-200">
            <div
              class="grid grid-cols-[1fr_auto_auto] gap-4 border-b border-gray-200 py-1.5 text-[8px] uppercase tracking-wide text-gray-500"
            >
              <span>Désignation</span>
              <span class="w-16 text-center">Mode</span>
              <span class="w-24 text-right">Montant encaissé</span>
            </div>
            <div class="grid grid-cols-[1fr_auto_auto] items-start gap-4 py-2">
              <div>
                <div>Versement · formation conduite automobile</div>
                <div class="mt-0.5 font-mono text-[7px] text-gray-400">
                  Sécurité : {{ recu.code_securite_hash.substring(0, 32) }}
                </div>
              </div>
              <span class="w-16 text-center">{{ recu.mode_reglement }}</span>
              <span class="w-24 text-right text-[11px] font-bold"
                >{{ fmt(recu.montant_verse) }} FCFA</span
              >
            </div>
          </div>

          <!-- Montant en lettres -->
          <div
            class="mt-3 rounded-md border border-gray-200 bg-gray-50 px-3 py-2"
          >
            <div class="text-[8.5px] text-gray-500">
              Arrêté le présent reçu à la somme de
            </div>
            <div
              class="mt-0.5 font-serif text-[12px] font-bold italic leading-tight"
            >
              {{ montantEnLettresFCFA(recu.montant_verse) }}
            </div>
          </div>

          <!-- Situation et cachet -->
          <div class="mt-3 grid grid-cols-2 items-stretch gap-3">
            <div class="flex flex-col justify-between py-1">
              <div>
                <div class="text-[8px] uppercase tracking-wide text-gray-500">
                  Situation du dossier
                </div>
                <div class="mt-1.5 flex items-baseline justify-between gap-2">
                  <span>Reste à payer</span>
                  <span
                    :class="[
                      'text-[10px] font-bold',
                      solde === 0 ? 'text-emerald-700' : 'text-gray-900',
                    ]"
                  >
                    {{ fmt(solde) }} FCFA
                  </span>
                </div>
              </div>
              <div class="mt-2 text-[8px] text-gray-500">
                {{
                  solde === 0
                    ? "Dossier entièrement soldé"
                    : "Solde à régler avant la prochaine tranche"
                }}
              </div>
            </div>

            <div
              class="rounded-[4px] border border-gray-300 px-3 py-2 text-right"
            >
              <div class="text-[8px] uppercase tracking-wide text-gray-500">
                L'auto-école · La caissière
              </div>
              <div
                class="mt-1.5 inline-block rounded-[3px] border border-emerald-700 px-1.5 py-0.5 text-[7.5px] font-semibold uppercase tracking-wide text-emerald-700"
              >
                Tampon officiel · Encaissé
              </div>
              <div class="mt-2.5 font-mono text-[8.5px] text-gray-700">
                {{ recu.caissier_nom }}
              </div>
            </div>
          </div>

          <div class="mt-3 text-center text-[7.5px] text-gray-400">
            Les acomptes ne sont pas remboursables. Conservez ce reçu.
          </div>
        </div>
      </div>

      <!-- Actions -->
      <footer
        class="flex items-center justify-between border-t border-border px-4 py-3"
      >
        <button
          type="button"
          @click="isA5ModalOpen = true"
          class="inline-flex items-center gap-1.5 rounded-[4px] border border-border bg-card px-3 py-1.5 text-xs font-medium text-foreground hover:bg-muted transition"
          title="Reçu A5 avec en-tête et tampon"
        >
          <FileText class="h-3.5 w-3.5 text-sky-500" />
          Format A5
        </button>

        <div class="flex items-center gap-2">
          <button
            type="button"
            @click="close"
            class="inline-flex items-center gap-1 rounded-[4px] border border-border bg-card px-3 py-1.5 text-xs text-foreground hover:bg-muted transition"
          >
            Fermer
            <kbd
              class="rounded bg-muted px-1 font-mono text-[10px] text-muted-foreground"
              >Échap</kbd
            >
          </button>
          <button
            type="button"
            @click="printTicket"
            class="inline-flex items-center gap-1.5 rounded-lg bg-brand-dark px-3 py-1.5 text-xs font-semibold text-white hover:bg-brand transition"
          >
            <Printer class="h-3.5 w-3.5" />
            Imprimer
            <kbd
              class="rounded bg-white/20 px-1 font-mono text-[10px] text-white/90"
              >Entrée</kbd
            >
          </button>
        </div>
      </footer>
    </div>

    <A5ReceiptModal v-if="isA5ModalOpen" @close="isA5ModalOpen = false" />
  </div>
</template>

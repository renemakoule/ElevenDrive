<script setup lang="ts">
import { computed } from "vue";
import { useGuichetStore } from "../../stores/guichetStore";
import { useAuthStore } from "../../stores/authStore";
import { montantEnLettresFCFA } from "../../utils/numberToWords";
import { FileText, X, Printer, Share2 } from "lucide-vue-next";

const emit = defineEmits<{
  (e: "close"): void;
}>();

const store = useGuichetStore();
const authStore = useAuthStore();

const meta = computed(() => {
  const school = authStore.currentSchool;
  const agence = authStore.currentAgence;
  const user = authStore.currentUser;
  const caissier =
    store.sessionCaisseDetail?.caissier_nom ||
    (user?.nom ? `${user.prenom} ${user.nom}` : "Agent Guichet");

  return {
    nomCommercial: school?.nom_commercial || "AutoSchool",
    numeroAgrement: school?.numero_agrement || "Agrément Ministériel",
    agence: agence?.nom || "Agence Guichet",
    adresse: agence?.adresse || "Siège social",
    telephone: agence?.telephone || school?.telephone || "",
    caissiere: caissier,
    guichet: "Guichet 1",
    session: store.sessionCaisseDetail?.id
      ? store.sessionCaisseDetail.id.slice(-8)
      : `SES-${new Date().toISOString().slice(0, 10)}`,
    recus: `${store.sessionCaisseDetail?.nombre_transactions ?? 0} émis`,
  };
});

const fmt = (n: number) => n.toLocaleString("fr-FR");

const now = new Date();
const dateStr = now.toLocaleDateString("fr-FR");
const heureStr = now.toLocaleTimeString("fr-FR", {
  hour: "2-digit",
  minute: "2-digit",
});

const ventilation = computed(() => store.sessionCaisseDetail?.ventilation);

const especes = computed(
  () => ventilation.value?.especes ?? store.stats.total_especes_theorique ?? 0,
);

const lignes = computed(() =>
  [
    { label: "Espèces", montant: especes.value, toujours: true },
    { label: "Wave", montant: ventilation.value?.wave ?? 0 },
    { label: "Orange Money", montant: ventilation.value?.orange_money ?? 0 },
    { label: "MTN MoMo", montant: ventilation.value?.mtn_momo ?? 0 },
    { label: "Moov Money", montant: ventilation.value?.moov ?? 0 },
  ].filter((l) => l.toujours || l.montant > 0),
);

const total = computed(
  () =>
    ventilation.value?.total_encaisse ??
    lignes.value.reduce((s, l) => s + l.montant, 0),
);

const fond = computed(() => store.stats.fond_de_caisse);
const theorique = computed(() => fond.value + especes.value);
const compte = computed(() => theorique.value);
const ecart = computed(() => compte.value - theorique.value);

const nbTransactions = computed(
  () => store.sessionCaisseDetail?.nombre_transactions ?? 0,
);

const printRapportZ = () => window.print();

const shareWhatsApp = () => {
  const detail = lignes.value
    .map((l) => `${l.label} : ${fmt(l.montant)} FCFA`)
    .join("\n");
  const texte = [
    "Rapport de clôture de caisse (Z)",
    `${dateStr} à ${heureStr}`,
    `Agence ${meta.value.agence}`,
    `Caissière : ${meta.value.caissiere}`,
    "",
    detail,
    `*Total du jour : ${fmt(total.value)} FCFA*`,
    "",
    `Fond de caisse : ${fmt(fond.value)} FCFA`,
    `Espèces encaissées : ${fmt(especes.value)} FCFA`,
    `Total en caisse : ${fmt(compte.value)} FCFA`,
    ecart.value === 0 ? "Écart : aucun" : `Écart : ${fmt(ecart.value)} FCFA`,
    "",
    "Session clôturée et verrouillée.",
  ].join("\n");

  window.open(`https://wa.me/?text=${encodeURIComponent(texte)}`, "_blank");
  store.showToast("Rapport prêt à être envoyé sur WhatsApp", "success");
};
</script>

<template>
  <div
    v-if="store.isRapportZModalOpen"
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
          <FileText class="h-4 w-4 text-brand-dark" />
          <span class="text-xs font-semibold"
            >Clôture de caisse (Rapport Z)</span
          >
        </div>
        <button
          @click="emit('close')"
          class="rounded p-1 text-muted-foreground hover:bg-muted hover:text-foreground transition"
          aria-label="Fermer"
        >
          <X class="h-4 w-4" />
        </button>
      </header>

      <!-- Rapport -->
      <div class="flex justify-center overflow-y-auto bg-muted/40 p-4">
        <div
          id="thermal-ticket-container"
          class="w-full max-w-[520px] select-text rounded-[5px] border border-gray-200 bg-white p-6 font-sans text-[10px] leading-snug text-gray-900"
        >
          <!-- Identité et numéro -->
          <div class="flex items-start justify-between gap-4">
            <div>
              <div
                class="text-[13px] font-bold leading-tight tracking-tight uppercase"
              >
                {{ meta.nomCommercial }}
              </div>
              <div class="mt-1 text-[8.5px] font-medium text-gray-600">
                Agrément ministériel : {{ meta.numeroAgrement }}
              </div>
              <div class="mt-0.5 text-[8.5px] text-gray-500">
                {{ meta.agence }} · {{ meta.adresse }}
              </div>
              <div v-if="meta.telephone" class="text-[8.5px] text-gray-500">
                Tél. {{ meta.telephone }}
              </div>
            </div>

            <div class="shrink-0 text-right">
              <div
                class="inline-block rounded-[3px] border-[1.5px] border-gray-900 px-2 py-1 font-mono text-[9px] font-medium tracking-wide"
              >
                Rapport Z N° {{ meta.session }}
              </div>
              <div class="mt-1.5 text-[8.5px] text-gray-500">
                {{ dateStr }} à {{ heureStr }}
              </div>
            </div>
          </div>

          <div class="mt-3 border-t border-gray-900"></div>

          <!-- Titre -->
          <div class="py-3 text-center">
            <div class="text-[9px] font-semibold uppercase tracking-[0.12em]">
              Rapport de clôture de caisse · Z
            </div>
          </div>

          <!-- Informations -->
          <div class="grid grid-cols-2 gap-x-6 gap-y-1.5">
            <div class="flex gap-3">
              <span class="w-16 shrink-0 text-gray-500">Caissière</span>
              <span class="font-medium">{{ meta.caissiere }}</span>
            </div>
            <div class="flex gap-3">
              <span class="w-16 shrink-0 text-gray-500">Guichet</span>
              <span class="font-medium">{{ meta.guichet }}</span>
            </div>
            <div class="flex gap-3">
              <span class="w-16 shrink-0 text-gray-500">Session</span>
              <span class="font-medium">{{ meta.session }}</span>
            </div>
            <div class="flex gap-3">
              <span class="w-16 shrink-0 text-gray-500">Opérations</span>
              <span class="font-medium"
                >{{ nbTransactions }} transaction{{
                  nbTransactions > 1 ? "s" : ""
                }}</span
              >
            </div>
          </div>

          <!-- Recettes -->
          <div class="mt-4 border-y border-gray-200">
            <div
              class="grid grid-cols-[1fr_auto] gap-4 border-b border-gray-200 py-1.5 text-[8px] uppercase tracking-wide text-gray-500"
            >
              <span>Mode de règlement</span>
              <span class="w-28 text-right">Montant encaissé</span>
            </div>
            <div
              v-for="l in lignes"
              :key="l.label"
              class="grid grid-cols-[1fr_auto] gap-4 border-b border-gray-100 py-1.5 last:border-b-0"
            >
              <span>{{ l.label }}</span>
              <span class="w-28 text-right font-medium"
                >{{ fmt(l.montant) }} FCFA</span
              >
            </div>
            <div
              class="grid grid-cols-[1fr_auto] gap-4 border-t border-gray-300 py-2"
            >
              <span class="font-semibold">Total des recettes du jour</span>
              <span class="w-28 text-right text-[11px] font-bold"
                >{{ fmt(total) }} FCFA</span
              >
            </div>
          </div>

          <!-- Total en lettres -->
          <div
            class="mt-3 rounded-[5px] border border-gray-200 bg-gray-50 px-3 py-2"
          >
            <div class="text-[8.5px] text-gray-500">
              Recettes du jour arrêtées à la somme de
            </div>
            <div
              class="mt-0.5 font-serif text-[12px] font-bold italic leading-tight"
            >
              {{ montantEnLettresFCFA(total) }}
            </div>
          </div>

          <!-- Tiroir et visa -->
          <div class="mt-3 grid grid-cols-2 items-stretch gap-3">
            <div class="py-1">
              <div class="text-[8px] uppercase tracking-wide text-gray-500">
                Contrôle du tiroir-caisse
              </div>
              <div class="mt-1.5 space-y-0.5">
                <div class="flex justify-between gap-2">
                  <span>Fond de départ</span><span>{{ fmt(fond) }} FCFA</span>
                </div>
                <div class="flex justify-between gap-2">
                  <span>Espèces encaissées</span
                  ><span>{{ fmt(especes) }} FCFA</span>
                </div>
                <div class="flex justify-between gap-2">
                  <span>Attendu</span><span>{{ fmt(theorique) }} FCFA</span>
                </div>
                <div class="flex justify-between gap-2">
                  <span>Compté</span><span>{{ fmt(compte) }} FCFA</span>
                </div>
                <div class="flex justify-between gap-2">
                  <span class="font-medium">Écart</span>
                  <span
                    :class="[
                      'font-bold',
                      ecart === 0 ? 'text-emerald-700' : 'text-rose-700',
                    ]"
                  >
                    {{ ecart === 0 ? "Aucun" : fmt(ecart) + " FCFA" }}
                  </span>
                </div>
              </div>
              <div class="mt-2 text-[8px] text-gray-500">
                Reçus {{ meta.recus }}
              </div>
            </div>

            <div
              class="flex flex-col justify-between rounded-md border border-gray-300 px-3 py-2 text-right"
            >
              <div>
                <div class="text-[8px] uppercase tracking-wide text-gray-500">
                  La caissière · La direction
                </div>
                <div
                  class="mt-1.5 inline-block rounded-[3px] border border-emerald-700 px-1.5 py-0.5 text-[7.5px] font-semibold uppercase tracking-wide text-emerald-700"
                >
                  Caisse clôturée · Verrouillée
                </div>
              </div>
              <div class="mt-4">
                <div class="font-mono text-[8.5px] text-gray-700">
                  {{ meta.caissiere }} ({{ meta.guichet }})
                </div>
                <div
                  class="mt-5 border-t border-gray-300 pt-1 text-[8px] text-gray-500"
                >
                  Visa de la direction
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Actions -->
      <footer
        class="flex items-center justify-between border-t border-border px-4 py-3"
      >
        <button
          type="button"
          @click="shareWhatsApp"
          class="inline-flex items-center gap-1.5 rounded-[4px] border border-border bg-card px-3 py-1.5 text-xs font-medium text-foreground hover:bg-muted transition"
        >
          <Share2 class="h-3.5 w-3.5 text-brand" />
          WhatsApp
        </button>

        <div class="flex items-center gap-2">
          <button
            type="button"
            @click="emit('close')"
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
            @click="printRapportZ"
            class="inline-flex items-center gap-1.5 rounded-[4px] bg-brand-dark px-3 py-1.5 text-xs font-semibold text-white hover:bg-brand transition"
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
  </div>
</template>

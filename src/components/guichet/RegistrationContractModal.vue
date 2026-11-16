<script setup lang="ts">
import { computed } from "vue";
import { useGuichetStore } from "../../stores/guichetStore";
import { useAuthStore } from "../../stores/authStore";
import { X, Printer, User } from "lucide-vue-next";

const emit = defineEmits<{
  (e: "close"): void;
}>();

const store = useGuichetStore();
const authStore = useAuthStore();
const detail = computed(() => store.activeDossierDetail);

const fcfa = (n: number) => (n ?? 0).toLocaleString("fr-FR");

const medicalLabels: Record<string, string> = {
  VALIDE_SANS_RESERVE: "Apte sans réserve",
  VALIDE_AVEC_VERRES: "Apte, port de verres obligatoire",
  EN_ATTENTE: "Certificat en attente",
};

const verse = computed(() =>
  detail.value
    ? detail.value.dossier.montant_total - detail.value.dossier.solde_restant
    : 0,
);

const identite = computed(() => {
  const c = detail.value?.candidat;
  if (!c) return [];
  return [
    { label: "Nom et prénoms", value: `${c.nom} ${c.prenoms}` },
    { label: "Sexe", value: c.genre === "M" ? "Masculin" : "Féminin" },
    { label: "Naissance", value: `${c.date_naissance} à ${c.lieu_naissance}` },
    { label: "Nationalité", value: c.nationalite },
    {
      label: "Pièce d'identité",
      value: `${c.type_piece_id} N° ${c.numero_piece_id}`,
    },
    { label: "WhatsApp", value: c.telephone_whatsapp },
    {
      label: "Aptitude médicale",
      value:
        medicalLabels[c.visite_medicale_statut] ?? c.visite_medicale_statut,
    },
    { label: "Quartier", value: c.quartier_residence || "Non précisé" },
  ];
});

const printContract = () => {
  window.print();
};
</script>

<template>
  <div
    v-if="detail"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 backdrop-blur-sm"
    role="dialog"
    aria-modal="true"
    aria-labelledby="contract-title"
  >
    <div
      class="flex max-h-[94vh] w-full max-w-3xl flex-col overflow-hidden rounded-[5px] border border-border bg-card"
    >
      <header
        class="flex items-center justify-between gap-3 border-b border-border px-4 py-3"
      >
        <div class="min-w-0">
          <h2 id="contract-title" class="text-sm font-semibold text-foreground">
            Contrat de formation
          </h2>
          <p class="text-[11px] text-muted-foreground">
            Fiche d'enrôlement et contrat · format A4
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

      <div class="flex-1 overflow-y-auto bg-muted/40 p-4">
        <article
          id="contract-a4-container"
          class="mx-auto w-full max-w-2xl select-text space-y-4 rounded-[5px] border border-slate-200 bg-white p-7 text-[11px] leading-normal text-slate-900 shadow-sm"
        >
          <!-- Émetteur -->
          <header
            class="flex items-start justify-between gap-4 border-b border-slate-900 pb-3"
          >
            <div class="space-y-0.5">
              <h1 class="text-base font-bold tracking-tight">
                {{ authStore.schoolName.toUpperCase() }}
              </h1>
              <p class="text-[10px] font-medium text-slate-700">
                Agrément ministériel : {{ authStore.schoolAgrement }}
              </p>
              <p class="text-[10px] text-slate-500">
                {{ authStore.schoolAddress }}
              </p>
              <p
                v-if="authStore.schoolPhone"
                class="text-[10px] text-slate-500"
              >
                Tél. {{ authStore.schoolPhone }}
              </p>
            </div>

            <div class="shrink-0 text-right">
              <div
                class="inline-block rounded border border-slate-900 px-2 py-0.5 font-mono text-[11px] font-semibold"
              >
                Dossier {{ detail.dossier.numero_dossier }}
              </div>
              <p class="mt-1 text-[10px] text-slate-500">
                Enrôlé le {{ detail.dossier.date_inscription }}
              </p>
              <p class="text-[10px] font-medium">
                Permis catégorie {{ detail.formule.categorie_permis }}
              </p>
            </div>
          </header>

          <h3
            class="text-center text-[11px] font-semibold uppercase tracking-wider"
          >
            Contrat d'apprentissage à la conduite automobile
          </h3>

          <!-- Candidat -->
          <section class="flex items-start gap-4">
            <div
              class="flex h-28 w-24 shrink-0 items-center justify-center overflow-hidden rounded border border-slate-300 bg-slate-50"
            >
              <img
                v-if="detail.candidat.photo_url"
                :src="detail.candidat.photo_url"
                alt="Photo d'identité"
                class="h-full w-full object-cover"
              />
              <span v-else class="flex flex-col items-center text-slate-400">
                <User class="h-7 w-7" />
                <span class="text-[9px]">Photo 4x4</span>
              </span>
            </div>

            <dl class="grid flex-1 grid-cols-2 gap-x-4 gap-y-2">
              <div v-for="i in identite" :key="i.label" class="min-w-0">
                <dt class="text-[10px] text-slate-500">{{ i.label }}</dt>
                <dd class="font-medium">{{ i.value }}</dd>
              </div>
            </dl>
          </section>

          <!-- Prestation -->
          <section class="border-t border-slate-200 pt-3">
            <h4
              class="mb-1.5 text-[10px] font-semibold uppercase tracking-wide text-slate-500"
            >
              Prestation souscrite
            </h4>
            <dl class="grid grid-cols-3 gap-4">
              <div>
                <dt class="text-[10px] text-slate-500">Formule</dt>
                <dd class="font-medium">{{ detail.formule.intitule }}</dd>
              </div>
              <div>
                <dt class="text-[10px] text-slate-500">Code théorique</dt>
                <dd class="font-medium tabular-nums">
                  {{ detail.formule.heures_theoriques }} heures
                </dd>
              </div>
              <div>
                <dt class="text-[10px] text-slate-500">Conduite</dt>
                <dd class="font-medium tabular-nums">
                  {{ detail.formule.heures_pratiques }} heures
                </dd>
              </div>
            </dl>
          </section>

          <!-- Conditions financières -->
          <section class="border-t border-slate-200 pt-3">
            <h4
              class="mb-1.5 text-[10px] font-semibold uppercase tracking-wide text-slate-500"
            >
              Conditions financières
            </h4>
            <dl class="grid grid-cols-3 gap-4">
              <div>
                <dt class="text-[10px] text-slate-500">Montant total</dt>
                <dd class="font-semibold tabular-nums">
                  {{ fcfa(detail.dossier.montant_total) }} FCFA
                </dd>
              </div>
              <div>
                <dt class="text-[10px] text-slate-500">Déjà versé</dt>
                <dd class="font-semibold tabular-nums text-emerald-700">
                  {{ fcfa(verse) }} FCFA
                </dd>
              </div>
              <div>
                <dt class="text-[10px] text-slate-500">Solde restant</dt>
                <dd
                  class="font-semibold tabular-nums"
                  :class="
                    detail.dossier.solde_restant > 0
                      ? 'text-red-700'
                      : 'text-emerald-700'
                  "
                >
                  {{ fcfa(detail.dossier.solde_restant) }} FCFA
                </dd>
              </div>
            </dl>

            <p class="mt-2 text-[10px] italic leading-snug text-slate-500">
              Le versement d'au moins 60 % du montant total est obligatoire
              avant le début des leçons de conduite. Le solde doit être
              intégralement réglé avant la convocation à l'examen officiel du
              permis.
            </p>
          </section>

          <!-- Signatures -->
          <section class="grid grid-cols-2 gap-6 pt-2">
            <div
              class="flex h-28 flex-col justify-between rounded border border-slate-300 p-3"
            >
              <span
                class="text-[10px] font-semibold uppercase tracking-wide text-slate-600"
                >Signature de l'élève</span
              >
              <span class="text-[9px] italic text-slate-400"
                >Mention « Lu et approuvé »</span
              >
            </div>

            <div
              class="flex h-28 flex-col justify-between rounded border border-slate-300 p-3 text-right"
            >
              <span
                class="text-[10px] font-semibold uppercase tracking-wide text-slate-600"
                >{{ authStore.schoolName }}</span
              >
              <span class="text-[9px] text-slate-400"
                >Cachet et signature de la direction</span
              >
            </div>
          </section>
        </article>
      </div>

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
          @click="printContract"
          class="inline-flex h-8 items-center gap-1.5 rounded-[4px] bg-brand-dark px-3 text-xs font-medium text-white transition hover:bg-brand focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-dark/40 focus-visible:ring-offset-1"
        >
          <Printer class="h-3.5 w-3.5" />
          Imprimer le contrat A4
          <kbd
            class="rounded bg-white/20 px-1 font-mono text-[10px] text-white/90"
            >Entrée</kbd
          >
        </button>
      </footer>
    </div>
  </div>
</template>

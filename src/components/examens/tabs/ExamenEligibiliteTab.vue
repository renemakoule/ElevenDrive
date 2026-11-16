<script setup lang="ts">
import type { CandidatExamenItem, DossierNonEligible } from "../../../types";
import { AlertTriangle, ShieldAlert } from "lucide-vue-next";

defineProps<{
  candidats: CandidatExamenItem[];
  autresDossiers: DossierNonEligible[];
  hasActiveSession: boolean;
  canForce: boolean;
  isLoading: boolean;
}>();

const emit = defineEmits<{
  (e: "force", dossier: DossierNonEligible): void;
}>();

const dateCourte = (iso: string) => {
  if (!iso) return "—";
  const d = new Date(`${iso}T12:00:00`);
  return Number.isNaN(d.getTime()) ? iso : d.toLocaleDateString("fr-FR");
};

const criteres = (c: CandidatExamenItem) => [
  {
    label:
      c.solde_restant <= 0
        ? "Solde réglé"
        : `Solde ${c.solde_restant.toLocaleString("fr-FR")} F`,
    ok: c.solde_restant <= 0,
  },
  {
    label: `Code ${c.heures_code_effectuees}/${c.heures_code_requises} h`,
    ok: c.heures_code_effectuees >= c.heures_code_requises,
  },
  {
    label: `Conduite ${c.heures_conduite_effectuees}/${c.heures_conduite_requises} h`,
    ok: c.heures_conduite_effectuees >= c.heures_conduite_requises,
  },
  {
    label: "Aptitude médicale",
    ok:
      c.visite_medicale_statut === "VALIDE_SANS_RESERVE" ||
      c.visite_medicale_statut === "VALIDE_AVEC_VERRES",
  },
];
</script>

<template>
  <div class="flex-1 space-y-3 overflow-y-auto p-3 text-[10px]">
    <div class="text-muted-foreground">
      Critères réglementaires d'inscription au bordereau : solde soldé (0 FCFA),
      quotas d'heures de code et de conduite atteints, certificat médical
      valide.
    </div>

    <div>
      <div class="mb-1 flex items-center justify-between">
        <span class="text-[11px] font-semibold text-foreground">
          {{
            hasActiveSession
              ? "Candidats retenus (Éligibles au bordereau)"
              : "Candidats prêts pour le prochain examen (Éligibles)"
          }}
        </span>
        <span class="text-muted-foreground"
          >{{ candidats.length }} candidat(s)</span
        >
      </div>

      <div
        class="divide-y divide-border overflow-hidden rounded-[5px] border border-border"
      >
        <div
          v-if="candidats.length === 0"
          class="py-8 px-4 text-center text-muted-foreground"
        >
          {{
            isLoading
              ? "Analyse de l'éligibilité des dossiers en cours..."
              : "Aucun candidat n'est actuellement éligible. Les dossiers doivent avoir un solde soldé, les heures requises validées et un certificat médical conforme."
          }}
        </div>
        <div
          v-for="c in candidats"
          :key="c.id"
          class="flex items-center justify-between gap-3 px-3 py-1.5 hover:bg-muted/60"
        >
          <div class="flex min-w-0 items-center gap-2.5">
            <span
              class="rounded bg-muted px-1.5 py-0.5 font-mono text-[9px] font-medium text-foreground"
              >{{ c.numero_table }}</span
            >
            <div class="min-w-0">
              <div
                class="flex flex-wrap items-center gap-1.5 font-medium text-foreground"
              >
                <span>{{ c.nom }} {{ c.prenoms }}</span>
                <span
                  class="font-mono text-[9px] font-normal text-muted-foreground"
                  >{{ c.numero_dossier }}</span
                >
                <span
                  v-if="c.est_derogation"
                  class="inline-flex items-center gap-1 rounded bg-amber-500/15 border border-amber-300 dark:border-amber-700/60 px-1.5 py-0.5 text-[8px] font-bold text-amber-700 dark:text-amber-300 uppercase tracking-wide"
                  title="Candidat inscrit sous dérogation exceptionnelle de la Direction"
                >
                  <ShieldAlert
                    class="h-2.5 w-2.5 text-amber-600 dark:text-amber-400"
                  />
                  Dérogation Direction
                </span>
              </div>
              <div class="text-[9px] text-muted-foreground">
                Né(e) le {{ dateCourte(c.date_naissance) }} ·
                {{ c.type_piece_id }} {{ c.numero_piece_id }} ·
                {{ c.telephone_whatsapp }}
              </div>
            </div>
          </div>

          <div class="flex shrink-0 items-center gap-2.5">
            <span
              v-for="k in criteres(c)"
              :key="k.label"
              class="flex items-center gap-1 text-[9px] text-muted-foreground"
            >
              <span
                class="h-1.5 w-1.5 rounded-full"
                :class="k.ok ? 'bg-brand' : 'bg-destructive'"
              ></span>
              {{ k.label }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <div>
      <div class="mb-1 text-[11px] font-semibold text-foreground">
        Dossiers non éligibles (Critères non remplis)
      </div>
      <div
        class="divide-y divide-border overflow-hidden rounded-[5px] border border-border"
      >
        <div
          v-if="autresDossiers.length === 0"
          class="py-4 px-3 text-center text-muted-foreground"
        >
          Aucun dossier bloqué. Tous les candidats enregistrés sont conformes !
        </div>
        <div
          v-for="d in autresDossiers"
          :key="d.numero_dossier"
          class="flex items-center justify-between gap-3 px-3 py-2 hover:bg-muted/40 transition"
        >
          <div class="min-w-0">
            <div class="font-medium text-foreground">
              {{ d.nom }}
              <span
                class="ml-1 font-mono text-[9px] font-normal text-muted-foreground"
                >{{ d.numero_dossier }}</span
              >
            </div>
            <div
              class="mt-0.5 flex flex-wrap items-center gap-1 text-[9px] text-amber-700 dark:text-amber-400"
            >
              <AlertTriangle class="h-2.5 w-2.5 shrink-0 text-amber-500" />
              {{ d.motifs.join(" · ") }}
            </div>
          </div>
          <div class="flex items-center gap-2 shrink-0">
            <span class="text-[9px] text-muted-foreground hidden sm:inline"
              >En cours de formation</span
            >
            <button
              v-if="hasActiveSession && canForce"
              type="button"
              @click="emit('force', d)"
              class="inline-flex items-center gap-1 rounded-md bg-amber-500/15 hover:bg-amber-500/25 text-amber-800 dark:text-amber-300 border border-amber-300 dark:border-amber-700/70 px-2 py-1 text-[10px] font-semibold transition shadow-xs"
              title="Accorder une dérogation exceptionnelle pour ajouter cet élève au bordereau d'État"
            >
              <ShieldAlert
                class="h-3 w-3 text-amber-600 dark:text-amber-400"
              />
              <span>Forcer au bordereau</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>


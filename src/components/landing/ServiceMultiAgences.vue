<script setup lang="ts">
import { computed } from "vue";
import {
  Building2,
  TrendingUp,
  Wrench,
  Users2,
  CheckCircle2,
  PieChart,
  BarChart3,
  ArrowRight,
  ShieldAlert,
} from "lucide-vue-next";
import { useLanguageStore } from "../../stores/languageStore";

const emit = defineEmits<{
  (e: "navigate-auth"): void;
}>();

const languageStore = useLanguageStore();
const t = (text: string) => languageStore.t(text);

const agencyPoints = computed(() => [
  {
    title: t("Consolidation financière multi-succursales"),
    desc: t("Visualisez les recettes, impayés et marges nettes de chaque agence depuis un cockpit unique."),
  },
  {
    title: t("Suivi du parc de véhicules & alertes"),
    desc: t("Échéances d'assurances, contrôles techniques, vidanges et consommations de carburant."),
  },
  {
    title: t("Cloisonnement des accès & rôles stricts"),
    desc: t("La secrétaire n'accède qu'à sa caisse locale, le directeur supervise tout le réseau."),
  },
  {
    title: t("Taux de réussite aux examens par moniteur"),
    desc: t("Indicateurs de performance pédagogique objectifs pour optimiser l'encadrement des candidats."),
  },
]);
</script>

<template>
  <section
    id="service-multi-agences"
    class="relative border-t border-border bg-muted/20 py-16 sm:py-20"
  >
    <div class="mx-auto max-w-6xl px-5 sm:px-6">
      <div class="grid items-center gap-10 lg:grid-cols-12 lg:gap-14">
        <!-- Colonne aperçu Stripe-like minimaliste -->
        <div class="order-2 lg:order-1 lg:col-span-6">
          <div
            class="relative rounded-[6px] border border-border bg-card p-4 shadow-sm sm:p-5"
          >
            <!-- Header cockpit multi-agences -->
            <div
              class="flex items-center justify-between border-b border-border pb-3"
            >
              <div class="flex items-center gap-2">
                <div
                  class="flex h-7 w-7 items-center justify-center rounded-[4px] bg-indigo-50 text-indigo-700 dark:bg-indigo-950/60 dark:text-indigo-300"
                >
                  <Building2 class="h-4 w-4" />
                </div>
                <div>
                  <p class="text-[12px] font-semibold text-foreground">
                    {{ t("Cockpit de Direction · Réseau Auto-Écoles") }}
                  </p>
                  <p class="text-[10px] text-muted-foreground">
                    {{ t("3 agences actives · 12 véhicules en rotation") }}
                  </p>
                </div>
              </div>
              <span
                class="rounded-[4px] bg-indigo-50 px-2 py-0.5 text-[10px] font-semibold text-indigo-700 border border-indigo-200/60 dark:bg-indigo-950/40 dark:text-indigo-300 dark:border-indigo-800/40"
              >
                {{ t("Vue Réseau") }}
              </span>
            </div>

            <!-- Liste des agences avec métriques compactes -->
            <div class="my-4 space-y-2 text-[12px]">
              <div
                class="flex items-center justify-between rounded-[4px] border border-border bg-background/50 p-2.5"
              >
                <div>
                  <div class="flex items-center gap-1.5">
                    <span class="h-2 w-2 rounded-full bg-emerald-500" />
                    <p class="text-[11px] font-semibold text-foreground">
                      {{ t("Agence Principale (Centre)") }}
                    </p>
                  </div>
                  <p class="text-[10px] text-muted-foreground">
                    {{ t("42 candidats actifs · 3 moniteurs") }}
                  </p>
                </div>
                <div class="text-right">
                  <p class="text-[11px] font-semibold text-foreground">
                    4 850 000 FCFA
                  </p>
                  <p class="text-[10px] text-emerald-600 font-medium">
                    {{ t("+18% ce mois") }}
                  </p>
                </div>
              </div>

              <div
                class="flex items-center justify-between rounded-[4px] border border-border bg-background/50 p-2.5"
              >
                <div>
                  <div class="flex items-center gap-1.5">
                    <span class="h-2 w-2 rounded-full bg-emerald-500" />
                    <p class="text-[11px] font-semibold text-foreground">
                      {{ t("Succursale Nord (Akwa)") }}
                    </p>
                  </div>
                  <p class="text-[10px] text-muted-foreground">
                    {{ t("28 candidats actifs · 2 moniteurs") }}
                  </p>
                </div>
                <div class="text-right">
                  <p class="text-[11px] font-semibold text-foreground">
                    2 940 000 FCFA
                  </p>
                  <p class="text-[10px] text-emerald-600 font-medium">
                    {{ t("+12% ce mois") }}
                  </p>
                </div>
              </div>

              <div
                class="flex items-center justify-between rounded-[4px] border border-border bg-background/50 p-2.5"
              >
                <div>
                  <div class="flex items-center gap-1.5">
                    <span class="h-2 w-2 rounded-full bg-emerald-500" />
                    <p class="text-[11px] font-semibold text-foreground">
                      {{ t("Succursale Sud (Bonapriso)") }}
                    </p>
                  </div>
                  <p class="text-[10px] text-muted-foreground">
                    {{ t("19 candidats actifs · 2 moniteurs") }}
                  </p>
                </div>
                <div class="text-right">
                  <p class="text-[11px] font-semibold text-foreground">
                    1 720 000 FCFA
                  </p>
                  <p class="text-[10px] text-emerald-600 font-medium">
                    {{ t("+8% ce mois") }}
                  </p>
                </div>
              </div>
            </div>

            <!-- Mini statut entretien flotte -->
            <div class="rounded-[4px] bg-muted/40 p-2.5">
              <div class="flex items-center justify-between">
                <span
                  class="flex items-center gap-1.5 text-[11px] font-medium text-foreground"
                >
                  <Wrench class="h-3.5 w-3.5 text-brand" />
                  {{ t("État du parc automobile") }}
                </span>
                <span class="text-[10px] font-medium text-emerald-600">
                  {{ t("11/12 véhicules opérationnels") }}
                </span>
              </div>
              <div
                class="mt-1.5 flex items-center justify-between text-[10px] text-muted-foreground"
              >
                <span>{{ t("Véhicule N°03 : Vidange programmée le 05 Nov") }}</span>
                <span
                  class="rounded-[2px] bg-amber-50 px-1 text-amber-700 dark:bg-amber-950/40 dark:text-amber-300"
                  >{{ t("Entretien planifié") }}</span
                >
              </div>
            </div>

            <!-- Footer KPI globale -->
            <div
              class="mt-3 flex items-center justify-between pt-2 border-t border-border text-[10px] text-muted-foreground"
            >
              <span class="inline-flex items-center gap-1">
                <TrendingUp class="h-3 w-3 text-emerald-600" />
                {{ t("Taux de réussite global au permis") }}
              </span>
              <span class="font-semibold text-foreground"
                >{{ t("89.4% (Moyenne réseau)") }}</span
              >
            </div>
          </div>
        </div>

        <!-- Colonne texte -->
        <div class="order-1 lg:order-2 lg:col-span-6">
          <span
            class="inline-flex items-center gap-1.5 rounded-[4px] border border-indigo-500/20 bg-indigo-50/60 px-2.5 py-1 text-[11px] font-semibold uppercase tracking-wider text-indigo-700 dark:bg-indigo-950/40 dark:text-indigo-300"
          >
            <Building2 class="h-3 w-3" />
            {{ t("Direction & Multi-Agences") }}
          </span>

          <h2
            class="mt-3 text-[24px] font-semibold leading-tight tracking-tight text-foreground sm:text-[30px]"
          >
            {{ t("Pilotez tout votre réseau d'agences d'un coup d'œil.") }}
          </h2>

          <p class="mt-3 text-[14px] leading-relaxed text-muted-foreground">
            {{
              t(
                "Que vous gériez une seule antenne ou dix succursales réparties sur le territoire, ElevenDrive unifie la gestion de vos équipes, de vos liquidités et de votre flotte automobile."
              )
            }}
          </p>

          <ul class="mt-6 space-y-3">
            <li
              v-for="pt in agencyPoints"
              :key="pt.title"
              class="flex items-start gap-2.5"
            >
              <CheckCircle2
                class="mt-0.5 h-4 w-4 shrink-0 text-indigo-600 dark:text-indigo-400"
              />
              <div>
                <p class="text-[13px] font-semibold text-foreground">
                  {{ pt.title }}
                </p>
                <p class="text-[12px] text-muted-foreground">{{ pt.desc }}</p>
              </div>
            </li>
          </ul>

          <div class="mt-6 flex items-center gap-3">
            <button
              type="button"
              class="inline-flex items-center gap-1.5 rounded-[4px] bg-foreground px-3.5 py-1.5 text-[12px] font-medium text-background transition hover:opacity-90 active:scale-[0.98]"
              @click="emit('navigate-auth')"
            >
              {{ t("Découvrir le cockpit") }}
              <ArrowRight class="h-3 w-3" />
            </button>
            <span
              class="inline-flex items-center gap-1 text-[11px] text-muted-foreground"
            >
              <BarChart3 class="h-3.5 w-3.5 text-indigo-600" />
              {{ t("Tableau de bord gérant temps réel") }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

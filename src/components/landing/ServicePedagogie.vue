<script setup lang="ts">
import { computed } from "vue";
import {
  GraduationCap,
  BookOpen,
  Car,
  Clock,
  CheckCircle2,
  CalendarCheck,
  Award,
  ArrowRight,
} from "lucide-vue-next";
import { useLanguageStore } from "../../stores/languageStore";

const emit = defineEmits<{
  (e: "navigate-auth"): void;
}>();

const languageStore = useLanguageStore();
const t = (text: string) => languageStore.t(text);

const pedagogicalPoints = computed(() => [
  {
    title: t("Livret numérique officiel"),
    desc: t("Suivi des 4 étapes de formation, validation des compétences et bilans de progression."),
  },
  {
    title: t("Planning de conduite interactif"),
    desc: t("Affectation rapide des moniteurs et des véhicules, gestion des créneaux et notifications."),
  },
  {
    title: t("Entraînements au Code de la Route"),
    desc: t("Séries thématiques, tests blancs chronométrés et analyse fine des fautes récurrentes."),
  },
  {
    title: t("Application moniteur mobile PWA"),
    desc: t("Émargement direct dans le véhicule, signature électronique et pointage en direct."),
  },
]);
</script>

<template>
  <section
    id="service-pedagogie"
    class="relative border-t border-border bg-muted/20 py-16 sm:py-20"
  >
    <div class="mx-auto max-w-6xl px-5 sm:px-6">
      <div class="grid items-center gap-10 lg:grid-cols-12 lg:gap-14">
        <!-- Colonne aperçu Stripe-like minimaliste (inversé sur desktop pour dynamiser la mise en page) -->
        <div class="order-2 lg:order-1 lg:col-span-6">
          <div
            class="relative rounded-[6px] border border-border bg-card p-4 shadow-sm sm:p-5"
          >
            <!-- Header fiche pédagogique -->
            <div
              class="flex items-center justify-between border-b border-border pb-3"
            >
              <div class="flex items-center gap-2">
                <div
                  class="flex h-7 w-7 items-center justify-center rounded-[4px] bg-sky-50 text-sky-700 dark:bg-sky-950/60 dark:text-sky-300"
                >
                  <GraduationCap class="h-4 w-4" />
                </div>
                <div>
                  <p class="text-[12px] font-semibold text-foreground">
                    {{ t("Livret d'Apprentissage Permis B") }}
                  </p>
                  <p class="text-[10px] text-muted-foreground">
                    {{ t("Moniteur référent : Koffi Serge · Véhicule : Toyota Yaris N° 04") }}
                  </p>
                </div>
              </div>
              <span
                class="rounded-[4px] bg-brand-50 px-2 py-0.5 text-[10px] font-semibold text-brand-700 border border-brand-200/60 dark:bg-brand-950/40 dark:text-brand-300 dark:border-brand-800/40"
              >
                {{ t("Étape 3 / 4 en cours") }}
              </span>
            </div>

            <!-- Jauges de progression concises -->
            <div class="my-4 grid grid-cols-2 gap-3 text-[12px]">
              <div
                class="rounded-[4px] border border-border bg-background/50 p-2.5"
              >
                <div class="flex items-center justify-between">
                  <span
                    class="flex items-center gap-1 text-[11px] font-medium text-muted-foreground"
                  >
                    <BookOpen class="h-3 w-3 text-sky-500" />
                    {{ t("Code théorique") }}
                  </span>
                  <span class="text-[11px] font-semibold text-foreground"
                    >88%</span
                  >
                </div>
                <div
                  class="mt-2 h-1.5 w-full rounded-[2px] bg-muted overflow-hidden"
                >
                  <div class="h-full bg-sky-500" style="width: 88%" />
                </div>
                <p class="mt-1 text-[10px] text-muted-foreground">
                  {{ t("37/40 au dernier examen blanc") }}
                </p>
              </div>

              <div
                class="rounded-[4px] border border-border bg-background/50 p-2.5"
              >
                <div class="flex items-center justify-between">
                  <span
                    class="flex items-center gap-1 text-[11px] font-medium text-muted-foreground"
                  >
                    <Car class="h-3 w-3 text-brand" />
                    {{ t("Heures de conduite") }}
                  </span>
                  <span class="text-[11px] font-semibold text-foreground"
                    >16 / 20h</span
                  >
                </div>
                <div
                  class="mt-2 h-1.5 w-full rounded-[2px] bg-muted overflow-hidden"
                >
                  <div class="h-full bg-brand-600" style="width: 80%" />
                </div>
                <p class="mt-1 text-[10px] text-muted-foreground">
                  {{ t("Prochaine séance : Demain à 14:00") }}
                </p>
              </div>
            </div>

            <!-- Liste des dernières compétences validées -->
            <div
              class="space-y-1.5 rounded-[4px] bg-muted/40 p-2.5 text-[11px]"
            >
              <div
                class="flex items-center justify-between text-muted-foreground"
              >
                <span
                  class="flex items-center gap-1.5 text-foreground font-medium"
                >
                  <CheckCircle2 class="h-3.5 w-3.5 text-emerald-600" />
                  {{ t("Trajectoires en agglomération & ronds-points") }}
                </span>
                <span class="text-[10px] text-emerald-600 font-semibold"
                  >{{ t("Validé") }}</span
                >
              </div>
              <div
                class="flex items-center justify-between text-muted-foreground"
              >
                <span
                  class="flex items-center gap-1.5 text-foreground font-medium"
                >
                  <CheckCircle2 class="h-3.5 w-3.5 text-emerald-600" />
                  {{ t("Créneau arrière et stationnement en bataille") }}
                </span>
                <span class="text-[10px] text-emerald-600 font-semibold"
                  >{{ t("Validé") }}</span
                >
              </div>
              <div
                class="flex items-center justify-between text-muted-foreground"
              >
                <span
                  class="flex items-center gap-1.5 text-foreground font-medium"
                >
                  <Clock class="h-3.5 w-3.5 text-amber-500" />
                  {{ t("Insertion sur voie rapide & allure soutenue") }}
                </span>
                <span class="text-[10px] text-amber-600 font-medium"
                  >{{ t("En apprentissage") }}</span
                >
              </div>
            </div>

            <!-- Footer certification & examen -->
            <div
              class="mt-3 flex items-center justify-between pt-2 border-t border-border text-[10px] text-muted-foreground"
            >
              <span class="inline-flex items-center gap-1">
                <Award class="h-3 w-3 text-brand" />
                {{ t("Dossier éligible à la présentation au permis") }}
              </span>
              <span class="font-medium text-foreground"
                >{{ t("Session Décembre 2024") }}</span
              >
            </div>
          </div>
        </div>

        <!-- Colonne texte -->
        <div class="order-1 lg:order-2 lg:col-span-6">
          <span
            class="inline-flex items-center gap-1.5 rounded-[4px] border border-sky-500/20 bg-sky-50/60 px-2.5 py-1 text-[11px] font-semibold uppercase tracking-wider text-sky-700 dark:bg-sky-950/40 dark:text-sky-300"
          >
            <GraduationCap class="h-3 w-3" />
            {{ t("Pédagogie & Conduite") }}
          </span>

          <h2
            class="mt-3 text-[24px] font-semibold leading-tight tracking-tight text-foreground sm:text-[30px]"
          >
            {{ t("Suivez la progression de chaque candidat, séance après séance.") }}
          </h2>

          <p class="mt-3 text-[14px] leading-relaxed text-muted-foreground">
            {{
              t(
                "Finis les livrets papier égarés ou détériorés. Les moniteurs enregistrent les heures de conduite et valident les compétences directement depuis leur smartphone, en toute simplicité."
              )
            }}
          </p>

          <ul class="mt-6 space-y-3">
            <li
              v-for="pt in pedagogicalPoints"
              :key="pt.title"
              class="flex items-start gap-2.5"
            >
              <CheckCircle2
                class="mt-0.5 h-4 w-4 shrink-0 text-sky-600 dark:text-sky-400"
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
              {{ t("Découvrir le livret") }}
              <ArrowRight class="h-3 w-3" />
            </button>
            <span
              class="inline-flex items-center gap-1 text-[11px] text-muted-foreground"
            >
              <CalendarCheck class="h-3.5 w-3.5 text-sky-600" />
              {{ t("Planning synchronisé en direct") }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

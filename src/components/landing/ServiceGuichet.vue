<script setup lang="ts">
import { computed } from "vue";
import {
  CreditCard,
  Receipt,
  FileSpreadsheet,
  CheckCircle2,
  ShieldCheck,
  Smartphone,
  Banknote,
  ArrowRight,
} from "lucide-vue-next";
import { useLanguageStore } from "../../stores/languageStore";

const emit = defineEmits<{
  (e: "navigate-auth"): void;
}>();

const languageStore = useLanguageStore();
const t = (text: string) => languageStore.t(text);

const points = computed(() => [
  {
    title: t("Enregistrement minute"),
    desc: t("Saisie rapide des pièces d'identité, photos et choix de la formule de permis."),
  },
  {
    title: t("Paiements échelonnés & Mobile Money"),
    desc: t("Tranches personnalisées, versements en espèces, Orange Money, MTN MoMo et Wave."),
  },
  {
    title: t("Tickets thermiques & reçus conformes"),
    desc: t("Impression instantanée sur imprimantes POS 80mm/58mm et reçus PDF dématérialisés."),
  },
  {
    title: t("Clôture de caisse certifiée (Rapport Z)"),
    desc: t("Audit mathématique automatique pour éliminer toute fuite de trésorerie au guichet."),
  },
]);
</script>

<template>
  <section
    id="service-guichet"
    class="relative border-t border-border py-16 sm:py-20"
  >
    <div class="mx-auto max-w-6xl px-5 sm:px-6">
      <div class="grid items-center gap-10 lg:grid-cols-12 lg:gap-14">
        <!-- Colonne texte -->
        <div class="lg:col-span-6">
          <span
            class="inline-flex items-center gap-1.5 rounded-[4px] border border-brand/20 bg-brand-50/60 px-2.5 py-1 text-[11px] font-semibold uppercase tracking-wider text-brand-700 dark:bg-brand-950/40 dark:text-brand-300"
          >
            <CreditCard class="h-3 w-3" />
            {{ t("Guichet & Encaissement") }}
          </span>

          <h2
            class="mt-3 text-[24px] font-semibold leading-tight tracking-tight text-foreground sm:text-[30px]"
          >
            {{ t("Une caisse sécurisée, rapide et infaillible au comptoir.") }}
          </h2>

          <p class="mt-3 text-[14px] leading-relaxed text-muted-foreground">
            {{
              t(
                "Conçu pour que vos secrétaires encaissent en quelques secondes sans jamais faire d'erreur de calcul. Tout est journalisé avec précision horodatée et contrôles hiérarchiques."
              )
            }}
          </p>

          <ul class="mt-6 space-y-3">
            <li
              v-for="pt in points"
              :key="pt.title"
              class="flex items-start gap-2.5"
            >
              <CheckCircle2
                class="mt-0.5 h-4 w-4 shrink-0 text-brand-600 dark:text-brand-400"
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
              {{ t("Faites un tour") }}
              <ArrowRight class="h-3 w-3" />
            </button>
            <span
              class="inline-flex items-center gap-1 text-[11px] text-muted-foreground"
            >
              <ShieldCheck class="h-3.5 w-3.5 text-emerald-600" />
              {{ t("Conforme règles ministérielles") }}
            </span>
          </div>
        </div>

        <!-- Colonne aperçu Stripe-like minimaliste -->
        <div class="lg:col-span-6">
          <div
            class="relative rounded-[6px] border border-border bg-card p-4 sm:p-5"
          >
            <!-- Header du reçu / caisse -->
            <div
              class="flex items-center justify-between border-b border-border pb-3"
            >
              <div class="flex items-center gap-2">
                <div
                  class="flex h-7 w-7 items-center justify-center rounded-[4px] bg-brand-50 text-brand-700 dark:bg-brand-950/60 dark:text-brand-300"
                >
                  <Receipt class="h-4 w-4" />
                </div>
                <div>
                  <p class="text-[12px] font-semibold text-foreground">
                    {{ t("Reçu de versement N° REC-2024-089") }}
                  </p>
                  <p class="text-[10px] text-muted-foreground">
                    {{ t("Guichet Central · Candidat : Amadou Diallo") }}
                  </p>
                </div>
              </div>
              <span
                class="rounded-[4px] bg-emerald-50 px-2 py-0.5 text-[10px] font-semibold text-emerald-700 border border-emerald-200/60 dark:bg-emerald-950/40 dark:text-emerald-300 dark:border-emerald-800/40"
              >
                {{ t("Payé & Certifié") }}
              </span>
            </div>

            <!-- Détail du versement -->
            <div class="my-4 space-y-2 text-[12px]">
              <div
                class="flex justify-between py-1 border-b border-border text-muted-foreground"
              >
                <span>{{ t("Formule sélectionnée") }}</span>
                <span class="font-medium text-foreground"
                  >{{ t("Permis B (Accéléré + Code)") }}</span
                >
              </div>
              <div
                class="flex justify-between py-1 border-b border-border text-muted-foreground"
              >
                <span>{{ t("Total de la formation") }}</span>
                <span class="font-medium text-foreground">185 000 FCFA</span>
              </div>
              <div
                class="flex justify-between py-1 border-b border-border text-muted-foreground"
              >
                <span>{{ t("Versement du jour") }}</span>
                <span class="font-semibold text-brand-600 dark:text-brand-400"
                  >+ 65 000 FCFA</span
                >
              </div>
              <div class="flex justify-between py-1 text-muted-foreground">
                <span>{{ t("Reste à payer") }}</span>
                <span class="font-medium text-foreground">40 000 FCFA</span>
              </div>
            </div>

            <!-- Badges méthodes de règlement acceptées -->
            <div class="rounded-[4px] bg-muted/50 p-2.5">
              <p
                class="text-[10px] font-medium text-muted-foreground uppercase tracking-wide"
              >
                {{ t("Moyens de paiement acceptés en temps réel") }}
              </p>
              <div class="mt-2 flex flex-wrap gap-1.5 text-[11px]">
                <span
                  class="inline-flex items-center gap-1 rounded-[3px] border border-border bg-card px-2 py-0.5 font-medium text-foreground"
                >
                  <Banknote class="h-3 w-3 text-emerald-600" /> {{ t("Espèces") }}
                </span>
                <span
                  class="inline-flex items-center gap-1 rounded-[3px] border border-border bg-card px-2 py-0.5 font-medium text-foreground"
                >
                  <Smartphone class="h-3 w-3 text-orange-500" /> Orange Money
                </span>
                <span
                  class="inline-flex items-center gap-1 rounded-[3px] border border-border bg-card px-2 py-0.5 font-medium text-foreground"
                >
                  <Smartphone class="h-3 w-3 text-yellow-500" /> MTN MoMo
                </span>
                <span
                  class="inline-flex items-center gap-1 rounded-[3px] border border-border bg-card px-2 py-0.5 font-medium text-foreground"
                >
                  <Smartphone class="h-3 w-3 text-sky-500" /> Wave
                </span>
              </div>
            </div>

            <!-- Mini footer clôture Z -->
            <div
              class="mt-3 flex items-center justify-between pt-2 border-t border-border text-[10px] text-muted-foreground"
            >
              <span class="inline-flex items-center gap-1">
                <FileSpreadsheet class="h-3 w-3 text-brand" />
                {{ t("Journal de caisse synchronisé") }}
              </span>
              <span>{{ t("Clôture Z : Conforme (0,00 FCFA d'écart)") }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

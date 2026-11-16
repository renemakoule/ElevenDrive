<script setup lang="ts">
import { computed } from "vue";
import {
  Database,
  Cloud,
  FileCheck2,
  Lock,
  WifiOff,
  RefreshCw,
  Server,
  ArrowRight,
  Shield,
} from "lucide-vue-next";
import { useLanguageStore } from "../../stores/languageStore";

const emit = defineEmits<{
  (e: "navigate-auth"): void;
}>();

const languageStore = useLanguageStore();
const t = (text: string) => languageStore.t(text);

const archivePoints = computed(() => [
  {
    title: t("Fonctionnement 100% hors-ligne (Local-First)"),
    desc: t("Vos postes guichets continuent de tourner à plein régime même en cas de panne Internet totale."),
  },
  {
    title: t("Synchronisation automatique sans conflit"),
    desc: t("Dès que la connexion revient, les écritures SQLite locales sont fusionnées en toute sécurité."),
  },
  {
    title: t("Coffre-fort Cloudflare R2 / S3 chiffré"),
    desc: t("Archivage sécurisé des pièces d'identité, photos d'identité officielles et certificats médicaux."),
  },
  {
    title: t("Dossiers d'examen prêts pour les inspecteurs"),
    desc: t("Bordereaux d'inscription certifiés, états de présentation et fiches de suivi en un clic."),
  },
]);
</script>

<template>
  <section
    id="service-archivage"
    class="relative border-t border-border py-16 sm:py-20"
  >
    <div class="mx-auto max-w-6xl px-5 sm:px-6">
      <div class="grid items-center gap-10 lg:grid-cols-12 lg:gap-14">
        <!-- Colonne texte -->
        <div class="lg:col-span-6">
          <span
            class="inline-flex items-center gap-1.5 rounded-[4px] border border-emerald-500/20 bg-emerald-50/60 px-2.5 py-1 text-[11px] font-semibold uppercase tracking-wider text-emerald-700 dark:bg-emerald-950/40 dark:text-emerald-300"
          >
            <Database class="h-3 w-3" />
            {{ t("Local-First & Archivage R2") }}
          </span>

          <h2
            class="mt-3 text-[24px] font-semibold leading-tight tracking-tight text-foreground sm:text-[30px]"
          >
            {{ t("Résilience totale aux pannes réseau et archivage pérenne.") }}
          </h2>

          <p class="mt-3 text-[14px] leading-relaxed text-muted-foreground">
            {{
              t(
                "L'activité d'une auto-école ne doit jamais être paralysée par une coupure de réseau. Grâce à son moteur local SQLite couplé au cloud R2, ElevenDrive garantit une fluidité continue."
              )
            }}
          </p>

          <ul class="mt-6 space-y-3">
            <li
              v-for="pt in archivePoints"
              :key="pt.title"
              class="flex items-start gap-2.5"
            >
              <FileCheck2
                class="mt-0.5 h-4 w-4 shrink-0 text-emerald-600 dark:text-emerald-400"
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
              {{ t("Voir la technologie") }}
              <ArrowRight class="h-3 w-3" />
            </button>
            <span
              class="inline-flex items-center gap-1 text-[11px] text-muted-foreground"
            >
              <Shield class="h-3.5 w-3.5 text-emerald-600" />
              {{ t("Chiffrement AES-256 au repos") }}
            </span>
          </div>
        </div>

        <!-- Colonne aperçu Stripe-like minimaliste -->
        <div class="lg:col-span-6">
          <div
            class="relative rounded-[6px] border border-border bg-card p-4 shadow-sm sm:p-5"
          >
            <!-- Header statut de réplication -->
            <div
              class="flex items-center justify-between border-b border-border pb-3"
            >
              <div class="flex items-center gap-2">
                <div
                  class="flex h-7 w-7 items-center justify-center rounded-[4px] bg-emerald-50 text-emerald-700 dark:bg-emerald-950/60 dark:text-emerald-300"
                >
                  <Server class="h-4 w-4" />
                </div>
                <div>
                  <p class="text-[12px] font-semibold text-foreground">
                    {{ t("Moteur de Réplication Local-First") }}
                  </p>
                  <p class="text-[10px] text-muted-foreground">
                    {{ t("Poste Guichet #1 · autoschool.db · SQLite v3") }}
                  </p>
                </div>
              </div>
              <span
                class="inline-flex items-center gap-1 rounded-[4px] bg-emerald-50 px-2 py-0.5 text-[10px] font-semibold text-emerald-700 border border-emerald-200/60 dark:bg-emerald-950/40 dark:text-emerald-300 dark:border-emerald-800/40"
              >
                <span
                  class="h-1.5 w-1.5 rounded-full bg-emerald-500 animate-pulse"
                />
                {{ t("Synchronisé") }}
              </span>
            </div>

            <!-- Simulation de cycle hors-ligne / cloud -->
            <div class="my-4 space-y-2 text-[12px]">
              <div
                class="flex items-center justify-between rounded-[4px] border border-border bg-background/50 p-2.5"
              >
                <div class="flex items-center gap-2">
                  <div
                    class="flex h-6 w-6 items-center justify-center rounded-[3px] bg-muted text-muted-foreground"
                  >
                    <Database class="h-3.5 w-3.5 text-foreground" />
                  </div>
                  <div>
                    <p class="text-[11px] font-medium text-foreground">
                      {{ t("Base locale") }}
                    </p>
                    <p class="text-[10px] text-muted-foreground">
                      {{ t("2 480 transactions enregistrées localement") }}
                    </p>
                  </div>
                </div>
                <span class="text-[10px] font-mono text-muted-foreground"
                  >0.8 ms</span
                >
              </div>

              <div
                class="flex items-center justify-between rounded-[4px] border border-border bg-background/50 p-2.5"
              >
                <div class="flex items-center gap-2">
                  <div
                    class="flex h-6 w-6 items-center justify-center rounded-[3px] bg-muted text-muted-foreground"
                  >
                    <Cloud class="h-3.5 w-3.5 text-sky-500" />
                  </div>
                  <div>
                    <p class="text-[11px] font-medium text-foreground">
                      {{ t("Coffre") }}
                    </p>
                    <p class="text-[10px] text-muted-foreground">
                      {{ t("Photos d'identité, CNI & certificats sauvegardés") }}
                    </p>
                  </div>
                </div>
                <span class="text-[10px] font-semibold text-emerald-600"
                  >{{ t("Sécurisé") }}</span
                >
              </div>
            </div>

            <!-- Preuve de résilience hors ligne -->
            <div class="rounded-[4px] bg-muted/40 p-2.5">
              <div class="flex items-center justify-between">
                <span
                  class="flex items-center gap-1.5 text-[11px] font-medium text-foreground"
                >
                  <WifiOff class="h-3.5 w-3.5 text-amber-500" />
                  {{ t("Mode coupure Internet activé ?") }}
                </span>
                <span
                  class="rounded-[3px] bg-background px-1.5 py-0.5 text-[10px] font-medium text-muted-foreground border border-border"
                >
                  {{ t("100% opérationnel") }}
                </span>
              </div>
              <p class="mt-1 text-[10px] text-muted-foreground">
                {{
                  t(
                    "Les reçus sont émis et numérotés sans interruption. La file d'attente se synchronise automatiquement dès le rétablissement de la 4G/Fibre."
                  )
                }}
              </p>
            </div>

            <!-- Footer hash intégrité -->
            <div
              class="mt-3 flex items-center justify-between pt-2 border-t border-border text-[10px] text-muted-foreground"
            >
              <span class="inline-flex items-center gap-1">
                <Lock class="h-3 w-3 text-emerald-600" />
                {{ t("Journal d'audit inviolable") }}
              </span>
              <span class="font-mono text-[9px] text-muted-foreground"
                >SHA-256: 9e4f...b72a</span
              >
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

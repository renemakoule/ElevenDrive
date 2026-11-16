<script setup lang="ts">
import { computed } from "vue";
import {
  ShieldCheck,
  Zap,
  Globe2,
  Users,
  CheckCircle,
  ArrowRight,
} from "lucide-vue-next";
import { useLanguageStore } from "../../stores/languageStore";

const emit = defineEmits<{
  (e: "navigate-auth"): void;
}>();

const languageStore = useLanguageStore();
const t = (text: string) => languageStore.t(text);

const stats = computed(() => [
  {
    value: "100%",
    label: t("Opérationnel hors-ligne"),
    sub: t("Aucun blocage de caisse"),
  },
  {
    value: "< 15 min",
    label: t("Prise en main"),
    sub: t("Interface pensée pour le guichet"),
  },
  {
    value: "0 FCFA",
    label: t("Écart de caisse"),
    sub: t("Rapprochement Z certifié"),
  },
  {
    value: "AES-256",
    label: t("Chiffrement des dossiers"),
    sub: t("Archivage cloud R2 sécurisé"),
  },
]);

const pillars = computed(() => [
  {
    icon: Zap,
    title: t("Vitesse d'exécution sans latence"),
    desc: t(
      "Développé en local pour une réactivité instantanée, sans dépendre du temps de chargement d'un navigateur classique."
    ),
  },
  {
    icon: ShieldCheck,
    title: t("Conformité & souveraineté des données"),
    desc: t(
      "Vos données restent physiquement dans votre établissement sur vos disques, tout en étant dupliquées de façon chiffrée sur votre coffre distant."
    ),
  },
  {
    icon: Globe2,
    title: t("Adapté aux spécificités locales"),
    desc: t(
      "Prise en charge intégrée des devises régionales (FCFA, EUR, USD), des numéros d'agrément officiels et des bordereaux ministériels."
    ),
  },
  {
    icon: Users,
    title: t("Écosystème unifié pour tous les acteurs"),
    desc: t(
      "Directeurs, secrétaires, moniteurs et Apprenants sur mobile partagent le même dossier sans friction."
    ),
  },
]);
</script>

<template>
  <section
    id="about"
    class="relative border-t border-border bg-muted/15 py-16 sm:py-20"
  >
    <div class="mx-auto max-w-6xl px-5 sm:px-6">
      <!-- En-tête -->
      <div class="mx-auto max-w-2xl text-center">
        <span
          class="inline-flex items-center gap-1.5 rounded-[4px] border border-border bg-card px-2.5 py-1 text-[11px] font-medium text-muted-foreground"
        >
          {{ t("Notre engagement & vision") }}
        </span>

        <h2
          class="mt-3 text-[26px] font-semibold leading-tight tracking-tight text-foreground sm:text-[32px]"
        >
          {{ t("Construit pour la réalité du terrain.") }}
        </h2>

        <p class="mt-3 text-[14px] leading-relaxed text-muted-foreground">
          {{
            t(
              "ElevenDrive est né d'un constat simple : la plupart des logiciels traditionnels s'arrêtent de fonctionner dès que la connexion Internet flanche. Nous avons conçu un système d'exploitation moderne, d'une fiabilité absolue."
            )
          }}
        </p>
      </div>

      <!-- Métriques d'impact -->
      <div class="mt-12 grid grid-cols-2 gap-3 sm:grid-cols-4 sm:gap-4">
        <div
          v-for="stat in stats"
          :key="stat.label"
          class="rounded-[4px] border border-border bg-card p-4 text-center"
        >
          <p
            class="text-[22px] font-semibold tracking-tight text-foreground sm:text-[26px]"
          >
            {{ stat.value }}
          </p>
          <p class="mt-1 text-[12px] font-medium text-foreground">
            {{ stat.label }}
          </p>
          <p class="mt-0.5 text-[10px] text-muted-foreground">{{ stat.sub }}</p>
        </div>
      </div>

      <!-- 4 Piliers fondamentaux -->
      <div class="mt-12 grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
        <div
          v-for="pillar in pillars"
          :key="pillar.title"
          class="rounded-[4px] border border-border bg-card p-4 transition hover:border-foreground/20"
        >
          <div
            class="flex h-8 w-8 items-center justify-center rounded-[4px] bg-muted text-foreground"
          >
            <component :is="pillar.icon" class="h-4 w-4" />
          </div>
          <h3 class="mt-3 text-[13px] font-semibold text-foreground">
            {{ pillar.title }}
          </h3>
          <p class="mt-1.5 text-[12px] leading-relaxed text-muted-foreground">
            {{ pillar.desc }}
          </p>
        </div>
      </div>

      <!-- Bloc mission & signature -->
      <div class="mt-12 rounded-[6px] border border-border bg-card p-6 sm:p-8">
        <div
          class="mx-auto flex flex-col items-center justify-between gap-6 md:flex-row"
        >
          <div class="max-w-xl">
            <p class="text-[14px] font-medium text-foreground">
              «
              {{
                t(
                  "Notre ambition est d'offrir à chaque auto-école, qu'elle soit en plein cœur d'une métropole ou dans une région enclavée, le même niveau d'excellence logicielle que les plus grandes entreprises de la tech mondiale."
                )
              }}
              »
            </p>
            <p class="mt-2 text-[11px] text-muted-foreground">
              {{ t("Équipe ElevenDrive · Ingénierie & Systèmes Embarqués GAS") }}
            </p>
          </div>

          <button
            type="button"
            class="shrink-0 inline-flex items-center gap-1.5 rounded-[4px] bg-brand-600 px-4 py-2 text-[12px] font-medium text-white transition hover:bg-brand-700 active:scale-[0.98]"
            @click="emit('navigate-auth')"
          >
            {{ t("Rejoindre la plateforme") }}
            <ArrowRight class="h-3 w-3" />
          </button>
        </div>
      </div>
    </div>
  </section>
</template>

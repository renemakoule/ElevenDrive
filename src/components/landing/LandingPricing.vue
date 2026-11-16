<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import {
  Check,
  Sparkles,
  ArrowRight,
  MessageSquare,
  Mail,
  X,
} from "lucide-vue-next";
import { useLanguageStore } from "../../stores/languageStore";
import { useAuthStore } from "../../stores/authStore";
import { useGuichetStore } from "../../stores/guichetStore";

const emit = defineEmits<{
  (e: "navigate-auth"): void;
}>();

const authStore = useAuthStore();
const guichetStore = useGuichetStore();
const languageStore = useLanguageStore();
const t = (text: string) => languageStore.t(text);

const isAnnual = ref(true);
const contactModalOuvert = ref(false);

const contactForm = ref({
  nom: "",
  autoEcole: "",
  telephone: "",
  email: "",
  besoins: "",
});

const plans = computed(() => [
  {
    id: "eleven",
    name: "Eleven",
    badge: t("Recommandé"),
    desc: t(
      "La solution complète tout-en-un pour les auto-écoles actives avec moniteurs sur le terrain.",
    ),
    priceMonthly: 12000,
    priceAnnual: 10200, // 12 000 FCFA avec -15% = 10 200 FCFA/mois (122 400 FCFA/an)
    popular: true,
    isCustom: false,
    features: [
      t("Jusqu'à 3 agences incluses"),
      t("Candidats et dossiers illimités"),
      t("Application PWA mobile pour moniteurs & élèves"),
      t("Livret d'apprentissage numérique interactif"),
      t("Archivage Cloudflare R2 certifié (photos & CNI)"),
      t("Gestion de la flotte automobile & alertes entretien"),
      t("Paiements Mobile Money (Orange, MTN, Wave)"),
      t("Support prioritaire WhatsApp & Téléphone 6j/7"),
    ],
    cta: t("Choisir l'offre Eleven"),
  },
  {
    id: "customeleven",
    name: "CustomEleven",
    badge: t("Sur-mesure"),
    desc: t(
      "Personnalisation avancée et déploiement sur-mesure du service ElevenDrive pour votre établissement ou réseau.",
    ),
    priceMonthly: null,
    priceAnnual: null,
    popular: false,
    isCustom: true,
    features: [
      t("Personnalisation complète de la plateforme ElevenDrive"),
      t("Nombre d'agences et succursales sur-mesure"),
      t("Intégrations spécifiques (ERP, terminaux & banques locales)"),
      t("Hébergement dédié ou infrastructure locale sécurisée"),
      t("Formation sur-mesure des équipes et directeurs pédagogiques"),
      t("Gestionnaire de compte dédié & SLA prioritaire 24/7"),
    ],
    cta: t("Contacter l'équipe ElevenDrive"),
  },
]);

const formatCurrency = (amount: number | null) => {
  if (amount === null || amount === undefined) return t("Sur devis");
  return new Intl.NumberFormat("fr-FR").format(amount) + " FCFA";
};

const handlePlanClick = (plan: (typeof plans.value)[number]) => {
  if (plan.isCustom) {
    contactModalOuvert.value = true;
  } else if (authStore.isAuthenticated) {
    guichetStore.showToast(
      "Votre établissement est sous licence ElevenDrive. Pour renouveler ou ajouter des agences, contactez l'équipe.",
      "info",
    );
    contactModalOuvert.value = true;
  } else {
    emit("navigate-auth");
  }
};

const contacterWhatsApp = () => {
  const message = [
    "Bonjour l'équipe ElevenDrive,",
    "",
    "Je souhaite obtenir des informations sur l'offre CustomEleven (personnalisation du service ElevenDrive) :",
    `- Nom : ${contactForm.value.nom || "Non renseigné"}`,
    `- Établissement : ${contactForm.value.autoEcole || "Non renseigné"}`,
    `- Téléphone : ${contactForm.value.telephone || "Non renseigné"}`,
    contactForm.value.email ? `- Email : ${contactForm.value.email}` : "",
    `- Besoins : ${contactForm.value.besoins || "Personnalisation globale du service"}`,
  ]
    .filter(Boolean)
    .join("\n");

  window.open(
    `https://wa.me/2250700000000?text=${encodeURIComponent(message)}`,
    "_blank",
  );
};

const contacterEmail = () => {
  const subject = encodeURIComponent(
    `Demande de personnalisation CustomEleven - ${
      contactForm.value.autoEcole || "Auto-École"
    }`,
  );
  const body = encodeURIComponent(
    `Bonjour l'équipe ElevenDrive,\n\nJe souhaite échanger sur la formule sur-mesure CustomEleven :\nNom : ${
      contactForm.value.nom || ""
    }\nAuto-École : ${contactForm.value.autoEcole || ""}\nTéléphone : ${
      contactForm.value.telephone || ""
    }\nBesoins spécifiques :\n${contactForm.value.besoins || ""}\n\nMerci !`,
  );
  window.open(
    `mailto:contact@elevendrive.com?subject=${subject}&body=${body}`,
    "_blank",
  );
};

onMounted(() => {
  if (authStore.isAuthenticated) {
    window.scrollTo(0, 0);
  }
});
</script>

<template>
  <section
    :id="authStore.isAuthenticated ? undefined : 'pricing'"
    class="relative"
    :class="
      authStore.isAuthenticated
        ? 'py-8 sm:py-12 border-0'
        : 'border-t border-border py-16 sm:py-20'
    "
  >
    <div class="mx-auto max-w-6xl px-5 sm:px-6">
      <!-- En-tête de section -->
      <div class="mx-auto max-w-2xl text-center">
        <span
          class="inline-flex items-center gap-1.5 rounded-[4px] border border-border bg-card px-2.5 py-1 text-[11px] font-medium text-muted-foreground"
        >
          {{ t("Tarification claire et transparente") }}
        </span>

        <h2
          class="mt-3 text-[26px] font-semibold leading-tight tracking-tight text-foreground sm:text-[32px]"
        >
          {{ t("Des tarifs simples, sans mauvaise surprise.") }}
        </h2>

        <p class="mt-3 text-[14px] leading-relaxed text-muted-foreground">
          {{
            t(
              "Aucun frais d'installation caché. Optez pour la formule Eleven standard ou personnalisez l'ensemble de votre solution avec CustomEleven.",
            )
          }}
        </p>

        <!-- Toggle Mensuel / Annuel -->
        <div
          class="mt-7 inline-flex items-center rounded-[4px] border border-border bg-card p-0.5"
        >
          <button
            type="button"
            class="rounded-[3px] px-3 py-1 text-[12px] font-medium transition"
            :class="
              !isAnnual
                ? 'bg-foreground text-background'
                : 'text-muted-foreground hover:text-foreground'
            "
            @click="isAnnual = false"
          >
            {{ t("Facturation mensuelle") }}
          </button>
          <button
            type="button"
            class="flex items-center gap-1.5 rounded-[3px] px-3 py-1 text-[12px] font-medium transition"
            :class="
              isAnnual
                ? 'bg-foreground text-background'
                : 'text-muted-foreground hover:text-foreground'
            "
            @click="isAnnual = true"
          >
            <span>{{ t("Facturation annuelle") }}</span>
            <span
              class="rounded-[2px] bg-brand-100 px-1 py-0.2 text-[10px] font-semibold text-brand-800 dark:bg-brand-950 dark:text-brand-300"
            >
              -15%
            </span>
          </button>
        </div>
      </div>

      <!-- Grille des offres (2 colonnes centrées) -->
      <div class="mt-12 grid gap-6 max-w-4xl mx-auto md:grid-cols-2">
        <div
          v-for="plan in plans"
          :key="plan.id"
          class="relative flex flex-col justify-between rounded-[8px] border bg-card p-6 sm:p-7 transition hover:border-foreground/30 shadow-xs"
          :class="
            plan.popular
              ? 'border-brand-500/70 ring-1 ring-brand-500/30'
              : 'border-border'
          "
        >
          <!-- Badge Recommandé / Sur-mesure -->
          <div v-if="plan.badge" class="absolute -top-2.5 right-4">
            <span
              class="inline-flex items-center gap-1 rounded-[3px] px-2 py-0.5 text-[10px] font-medium text-white shadow-xs"
              :class="
                plan.popular ? 'bg-brand-600' : 'bg-slate-700 dark:bg-slate-600'
              "
            >
              <Sparkles class="h-2.5 w-2.5" />
              {{ plan.badge }}
            </span>
          </div>

          <div>
            <div class="flex items-baseline justify-between">
              <h3 class="text-[18px] font-semibold text-foreground">
                {{ plan.name }}
              </h3>
            </div>

            <p
              class="mt-2 min-h-[38px] text-[12px] leading-relaxed text-muted-foreground"
            >
              {{ plan.desc }}
            </p>

            <!-- Prix -->
            <div class="mt-5 pb-5 border-b border-border">
              <template v-if="!plan.isCustom">
                <div class="flex items-baseline gap-1">
                  <span
                    class="text-[30px] font-bold tracking-tight text-foreground"
                  >
                    {{
                      formatCurrency(
                        isAnnual ? plan.priceAnnual : plan.priceMonthly,
                      )
                    }}
                  </span>
                  <span class="text-[12px] text-muted-foreground"
                    >/ {{ t("mois") }}</span
                  >
                </div>
                <p class="mt-1 text-[11px] text-muted-foreground">
                  {{
                    isAnnual
                      ? t(
                          "Facturé annuellement (122 400 FCFA/an avec 15% de réduction)",
                        )
                      : t("Sans engagement de durée (12 000 FCFA/mois)")
                  }}
                </p>
              </template>
              <template v-else>
                <div class="flex items-baseline gap-1.5">
                  <span
                    class="text-[30px] font-bold tracking-tight text-foreground"
                  >
                    {{ t("Sur devis") }}
                  </span>
                  <span class="text-[12px] text-muted-foreground"
                    >({{ t("Sur-mesure") }})</span
                  >
                </div>
                <p class="mt-1 text-[11px] text-muted-foreground">
                  {{ t("Solution taillée selon votre cahier des charges") }}
                </p>
              </template>
            </div>

            <!-- Liste des fonctionnalités -->
            <ul class="mt-5 space-y-2.5">
              <li
                v-for="(feat, idx) in plan.features"
                :key="idx"
                class="flex items-start gap-2.5 text-[12px] text-foreground"
              >
                <Check
                  class="mt-0.5 h-3.5 w-3.5 shrink-0 text-brand-600 dark:text-brand-400"
                />
                <span class="leading-tight">{{ feat }}</span>
              </li>
            </ul>
          </div>

          <!-- Bouton CTA -->
          <div class="mt-8 pt-4">
            <button
              type="button"
              class="flex w-full items-center justify-center gap-1.5 rounded-[4px] px-3.5 py-2.5 text-[12px] font-medium transition active:scale-[0.98]"
              :class="
                plan.popular
                  ? 'bg-brand-600 text-white hover:bg-brand-700 shadow-xs'
                  : 'border border-border bg-background text-foreground hover:bg-muted'
              "
              @click="handlePlanClick(plan)"
            >
              <span>{{ plan.cta }}</span>
              <component
                :is="plan.isCustom ? MessageSquare : ArrowRight"
                class="h-3.5 w-3.5"
              />
            </button>
          </div>
        </div>
      </div>

      <!-- Mini rassurance -->
      <div class="mt-12 text-center text-[12px] text-muted-foreground">
        {{
          t(
            "Besoin d'un accompagnement personnalisé ou d'une intégration spécifique ?",
          )
        }}
        <button
          type="button"
          class="font-medium text-foreground underline underline-offset-4 hover:text-brand transition ml-1"
          @click="contactModalOuvert = true"
        >
          {{ t("Contactez l'équipe commerciale ElevenDrive") }}
        </button>
      </div>
    </div>

    <!-- Modale de Contact CustomEleven -->
    <div
      v-if="contactModalOuvert"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 select-none backdrop-blur-xs"
    >
      <div
        class="w-full max-w-lg overflow-hidden rounded-xl border border-border bg-card shadow-2xl animate-in fade-in zoom-in-95 duration-150"
      >
        <!-- En-tête -->
        <div
          class="flex items-center justify-between border-b border-border bg-muted/30 px-5 py-4"
        >
          <div class="flex items-center gap-2.5">
            <div
              class="flex h-8 w-8 items-center justify-center rounded-lg bg-brand-100 dark:bg-brand-950 text-brand-600"
            >
              <Sparkles class="h-4 w-4" />
            </div>
            <div>
              <h3 class="text-sm font-semibold text-foreground">
                {{ t("Personnalisation CustomEleven") }}
              </h3>
              <p class="text-[11px] text-muted-foreground">
                {{ t("Échangez directement avec l'équipe ElevenDrive") }}
              </p>
            </div>
          </div>
          <button
            type="button"
            @click="contactModalOuvert = false"
            class="rounded-md p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground transition"
          >
            <X class="h-4 w-4" />
          </button>
        </div>

        <!-- Corps -->
        <form
          @submit.prevent="contacterWhatsApp"
          class="p-5 space-y-3.5 text-xs"
        >
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div>
              <label class="block text-[11px] font-medium text-foreground mb-1">
                {{ t("Votre nom complet") }}
              </label>
              <input
                type="text"
                v-model="contactForm.nom"
                placeholder="Ex : MAKOULE Rene"
                class="w-full h-8 px-2.5 rounded-md border border-border bg-background text-foreground text-xs focus:outline-none focus:ring-1 focus:ring-brand"
              />
            </div>
            <div>
              <label class="block text-[11px] font-medium text-foreground mb-1">
                {{ t("Nom de l'Auto-École / Réseau") }}
              </label>
              <input
                type="text"
                v-model="contactForm.autoEcole"
                placeholder="Ex : Auto-École du Progrès"
                class="w-full h-8 px-2.5 rounded-md border border-border bg-background text-foreground text-xs focus:outline-none focus:ring-1 focus:ring-brand"
              />
            </div>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div>
              <label class="block text-[11px] font-medium text-foreground mb-1">
                {{ t("Téléphone / WhatsApp") }}
              </label>
              <input
                type="tel"
                v-model="contactForm.telephone"
                placeholder="+237 651..."
                class="w-full h-8 px-2.5 rounded-md border border-border bg-background text-foreground text-xs focus:outline-none focus:ring-1 focus:ring-brand"
              />
            </div>
            <div>
              <label class="block text-[11px] font-medium text-foreground mb-1">
                {{ t("Adresse Email") }}
              </label>
              <input
                type="email"
                v-model="contactForm.email"
                placeholder="contact@autoecole.com"
                class="w-full h-8 px-2.5 rounded-md border border-border bg-background text-foreground text-xs focus:outline-none focus:ring-1 focus:ring-brand"
              />
            </div>
          </div>

          <div>
            <label class="block text-[11px] font-medium text-foreground mb-1">
              {{ t("Vos besoins de personnalisation") }}
            </label>
            <textarea
              v-model="contactForm.besoins"
              rows="3"
              placeholder="Ex : Nous gérons 5 agences, besoin d'un domaine personnalisé et d'une formation de 12 moniteurs..."
              class="w-full p-2.5 rounded-md border border-border bg-background text-foreground text-xs focus:outline-none focus:ring-1 focus:ring-brand resize-none"
            ></textarea>
          </div>

          <!-- Canaux de contact direct rapide -->
          <div class="pt-2 flex flex-col sm:flex-row items-center gap-2">
            <button
              type="button"
              @click="contacterWhatsApp"
              class="w-full inline-flex items-center justify-center gap-1.5 rounded-md bg-emerald-600 hover:bg-emerald-700 px-3 py-2 text-xs font-semibold text-white transition shadow-xs"
            >
              <MessageSquare class="h-3.5 w-3.5" />
              <span>{{ t("Contacter via WhatsApp") }}</span>
            </button>
            <button
              type="button"
              @click="contacterEmail"
              class="w-full inline-flex items-center justify-center gap-1.5 rounded-md border border-border bg-muted/60 hover:bg-muted px-3 py-2 text-xs font-semibold text-foreground transition"
            >
              <Mail class="h-3.5 w-3.5 text-muted-foreground" />
              <span>{{ t("Envoyer par Email") }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import { useAuthStore } from "../../stores/authStore";
import { useGuichetStore } from "../../stores/guichetStore";
import {
  X,
  Check,
  CheckCircle2,
  Sparkles,
  ArrowRight,
  ShieldCheck,
  Zap,
  CreditCard,
  Building2,
  Phone,
  MessageSquare,
  Clock,
  AlertTriangle,
} from "lucide-vue-next";

defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
}>();

const authStore = useAuthStore();
const guichetStore = useGuichetStore();

const isAnnual = ref(true);
const selectedPlanId = ref<"eleven" | "customeleven">("eleven");

// --- Analyse de la licence de l'auto-école courante ---
const school = computed(() => authStore.currentSchool);
const isLicenceExpired = computed(() => authStore.isLicenceExpired);
const daysRemaining = computed(() => authStore.trialDaysRemaining);

const rawTypeLicence = computed(() =>
  (school.value?.type_licence || "TRIAL").toUpperCase(),
);

const isTrial = computed(() => {
  const t = rawTypeLicence.value;
  return t.includes("TRIAL") || t.includes("ESSAI") || t.includes("GRATUIT");
});

const isCustomSchool = computed(() => {
  const t = rawTypeLicence.value;
  return t.includes("CUSTOM") || t.includes("SUR-MESURE");
});

// Le plan est considéré comme actif et en cours si la licence n'est pas expirée
// et qu'il reste au moins 1 jour de validité (essai 30j ou abonnement souscrit)
const isElevenProPlanActive = computed(() => {
  if (!school.value) return false;
  return !isCustomSchool.value && !isLicenceExpired.value && daysRemaining.value > 0;
});

const isCustomPlanActive = computed(() => {
  if (!school.value) return false;
  return isCustomSchool.value && !isLicenceExpired.value && daysRemaining.value > 0;
});

const expirationDateFormatted = computed(() => {
  const raw = school.value?.date_expiration_licence;
  if (!raw) return "";
  try {
    return new Date(raw).toLocaleDateString("fr-FR", {
      day: "2-digit",
      month: "long",
      year: "numeric",
    });
  } catch {
    return "";
  }
});

const daysRemainingLabel = computed(() => {
  const d = daysRemaining.value;
  if (d <= 0) return "Expiré";
  if (d === 1) return "1 jour restant";
  return `${d} jours restants`;
});

const plans = computed(() => [
  {
    id: "eleven" as const,
    name: "Eleven Pro",
    badge: "Recommandé · Bon Plan",
    desc: "La solution tout-en-un pour gérer votre auto-école, vos agences et moniteurs sans limitation.",
    priceMonthly: 12000,
    priceAnnual: 10200, // -15% en annuel
    popular: true,
    isActive: isElevenProPlanActive.value,
    isCurrentPlan: !isCustomSchool.value && Boolean(school.value),
    features: [
      "Jusqu'à 3 agences physiques incluses",
      "Candidats et dossiers illimités",
      "Application PWA mobile pour moniteurs & élèves",
      "Livret d'apprentissage numérique interactif",
      "Archivage Cloudflare R2 certifié (photos & CNI)",
      "Gestion de la flotte automobile & alertes entretien",
      "Paiements Mobile Money (Orange, MTN, Wave)",
      "Support prioritaire WhatsApp & Téléphone 6j/7",
    ],
  },
  {
    id: "customeleven" as const,
    name: "Custom Eleven",
    badge: "Sur-Mesure",
    desc: "Déploiement personnalisé pour réseaux d'auto-écoles et flottes multi-régionales.",
    priceMonthly: null,
    priceAnnual: null,
    popular: false,
    isActive: isCustomPlanActive.value,
    isCurrentPlan: isCustomSchool.value && Boolean(school.value),
    features: [
      "Agences et succursales illimitées",
      "Intégration bancaire directe & terminaux sur site",
      "Hébergement dédié et souveraineté locale",
      "Formation personnalisée de vos équipes",
      "Directeur de compte dédié 24/7",
    ],
  },
]);

const formatPrice = (amount: number | null) => {
  if (amount === null) return "Sur devis";
  return new Intl.NumberFormat("fr-FR").format(amount) + " FCFA";
};

const handleSelectPlan = (planId: "eleven" | "customeleven", isActive: boolean) => {
  // Si le plan est encore actif, le clic est désactivé
  if (isActive) return;

  selectedPlanId.value = planId;
  const schoolName = authStore.schoolName || "Mon Auto-École";
  const planName = planId === "eleven" ? "Eleven Pro" : "Custom Eleven";
  const billingType = isAnnual.value ? "Annuelle (-15%)" : "Mensuelle";
  const msg = encodeURIComponent(
    `Bonjour ElevenDrive,\nJe souhaite souscrire / renouveler l'abonnement pour mon auto-école "${schoolName}".\nFormule : ${planName} (Facturation ${billingType}).`,
  );
  window.open(`https://wa.me/237651727932?text=${msg}`, "_blank");
};
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-[110] flex items-center justify-center bg-black/80 backdrop-blur-md p-3 sm:p-5 select-none animate-in fade-in duration-200"
    @click.self="emit('close')"
  >
    <div
      class="bg-card border border-border rounded-[3px] max-w-4xl w-full overflow-hidden flex flex-col max-h-[92vh] animate-in zoom-in-95 duration-200 text-foreground"
    >
      <!-- En-tête -->
      <div
        class="p-4 bg-muted/40 border-b border-border flex items-center justify-between shrink-0"
      >
        <div class="flex items-center gap-2.5">
          <div
            class="w-9 h-9 rounded-lg bg-amber-500/15 text-amber-500 flex items-center justify-center font-bold shadow-xs border border-amber-500/30"
          >
            <Sparkles class="w-5 h-5 text-amber-500" />
          </div>
          <div>
            <div
              class="font-bold text-foreground text-sm flex items-center gap-2"
            >
              <span>Bons Plans & Abonnements ElevenDrive</span>
              <span
                v-if="!isLicenceExpired && school"
                class="text-[9px] uppercase font-mono px-1.5 py-0.2 rounded bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 font-semibold border border-emerald-500/30 flex items-center gap-1"
              >
                <CheckCircle2 class="w-2.5 h-2.5" />
                <span>{{ isTrial ? 'Essai 30j Actif' : 'Abonnement Actif' }}</span>
              </span>
              <span
                v-else-if="isLicenceExpired && school"
                class="text-[9px] uppercase font-mono px-1.5 py-0.2 rounded bg-rose-500/15 text-rose-600 dark:text-rose-400 font-semibold border border-rose-500/30 flex items-center gap-1"
              >
                <AlertTriangle class="w-2.5 h-2.5" />
                <span>Abonnement Expiré</span>
              </span>
              <span
                v-else
                class="text-[9px] uppercase font-mono px-1.5 py-0.2 rounded bg-amber-500/20 text-amber-600 dark:text-amber-400 font-semibold border border-amber-500/30"
              >
                Passez au plan supérieur
              </span>
            </div>
            <div class="text-[11px] text-muted-foreground">
              Choisissez la formule idéale pour débloquer votre auto-école sans
              interruption
            </div>
          </div>
        </div>

        <button
          type="button"
          @click="emit('close')"
          class="p-1.5 rounded-[4px] text-muted-foreground hover:text-foreground hover:bg-muted transition cursor-pointer"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Corps -->
      <div class="p-5 overflow-y-auto space-y-5 flex-1 text-xs">
        <!-- Bandeau d'information spécifique au statut de licence de l'auto-école -->
        <div
          v-if="school && (isElevenProPlanActive || isCustomPlanActive)"
          class="rounded-[4px] border border-emerald-500/30 bg-emerald-500/10 dark:bg-emerald-500/10 p-3 flex flex-col sm:flex-row sm:items-center justify-between gap-2 text-xs"
        >
          <div class="flex items-center gap-2">
            <ShieldCheck class="w-4 h-4 text-emerald-600 dark:text-emerald-400 shrink-0" />
            <div>
              <span class="font-bold text-foreground">
                {{ isTrial ? 'Période d\'essai gratuite (30 jours) en cours' : 'Votre abonnement est actuellement actif' }} :
              </span>
              <span class="text-muted-foreground ml-1">
                Formule {{ rawTypeLicence }} sur votre auto-école {{ authStore.schoolName }}.
              </span>
            </div>
          </div>
          <div class="flex items-center gap-1.5 font-mono text-[11px] font-bold text-emerald-600 dark:text-emerald-400 shrink-0">
            <Clock class="w-3.5 h-3.5" />
            <span>{{ daysRemainingLabel }}</span>
            <span v-if="expirationDateFormatted" class="font-normal text-muted-foreground text-[10px]">
              (Jusqu'au {{ expirationDateFormatted }})
            </span>
          </div>
        </div>

        <div
          v-else-if="school && isLicenceExpired"
          class="rounded-[4px] border border-rose-500/30 bg-rose-500/10 p-3 flex flex-col sm:flex-row sm:items-center justify-between gap-2 text-xs"
        >
          <div class="flex items-center gap-2 text-rose-700 dark:text-rose-400">
            <AlertTriangle class="w-4 h-4 shrink-0" />
            <div>
              <span class="font-bold">Votre période d'abonnement / essai gratuit est arrivée à échéance.</span>
              <div class="text-[11px] opacity-90 mt-0.5">
                Le blocage a été levé sur les formules : vous pouvez souscrire ou renouveler immédiatement ci-dessous.
              </div>
            </div>
          </div>
        </div>

        <!-- Basculeur Mensuel / Annuel -->
        <div class="flex items-center justify-center">
          <div
            class="inline-flex items-center p-1 rounded-[3px] bg-muted border border-border text-xs font-semibold"
          >
            <button
              type="button"
              @click="isAnnual = false"
              class="px-3.5 py-1.5 rounded-[3px] transition cursor-pointer"
              :class="
                !isAnnual
                  ? 'bg-card text-foreground font-bold'
                  : 'text-muted-foreground hover:text-foreground'
              "
            >
              Facturation Mensuelle
            </button>
            <button
              type="button"
              @click="isAnnual = true"
              class="px-3.5 py-1.5 rounded-[3px] transition flex items-center gap-1.5 cursor-pointer"
              :class="
                isAnnual
                  ? 'bg-brand text-white font-bold'
                  : 'text-muted-foreground hover:text-foreground'
              "
            >
              <span>Annuelle</span>
              <span
                class="px-1.5 py-0.2 rounded-full text-[9px] font-mono bg-white/20 text-white font-bold"
              >
                -15%
              </span>
            </button>
          </div>
        </div>

        <!-- Cartes de Plans -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 items-stretch">
          <div
            v-for="plan in plans"
            :key="plan.id"
            class="relative rounded-[6px] border p-5 flex flex-col justify-between transition-all duration-200"
            :class="[
              plan.isActive
                ? 'border-border/80 bg-muted/20 opacity-85 shadow-none'
                : plan.popular
                  ? 'border-brand bg-gradient-to-b from-brand/5 via-card to-card ring-1 ring-brand/50 shadow-2xl'
                  : 'border-border bg-card hover:border-border/80',
            ]"
          >
            <!-- Badge supérieur (Recommandé / Bon Plan ou Abonnement en cours) -->
            <!-- 1. CAS OÙ LE PLAN EST ACTIF : Marqué grisé / actif avec la période restante -->
            <div
              v-if="plan.isActive"
              class="absolute -top-2.5 right-4 px-2.5 py-0.5 rounded-full bg-muted border border-emerald-500/40 text-emerald-600 dark:text-emerald-400 text-[9px] font-extrabold uppercase tracking-wider shadow-xs flex items-center gap-1"
            >
              <CheckCircle2 class="w-3 h-3 text-emerald-500" />
              <span>Abonnement Actif · {{ daysRemainingLabel }}</span>
            </div>

            <!-- 2. CAS OÙ LE PLAN EST EXPIRÉ : Alerte de renouvellement -->
            <div
              v-else-if="plan.isCurrentPlan && isLicenceExpired"
              class="absolute -top-2.5 right-4 px-2.5 py-0.5 rounded-full bg-rose-600 text-white text-[9px] font-extrabold uppercase tracking-wider shadow-sm flex items-center gap-1 animate-pulse"
            >
              <AlertTriangle class="w-2.5 h-2.5" />
              <span>Expiré · À renouveler</span>
            </div>

            <!-- 3. CAS STANDARD : Badge Bon Plan classique -->
            <div
              v-else-if="plan.popular"
              class="absolute -top-2.5 right-4 px-2.5 py-0.5 rounded-full bg-brand text-white text-[9px] font-extrabold uppercase tracking-wider shadow-sm flex items-center gap-1"
            >
              <Sparkles class="w-2.5 h-2.5" />
              <span>{{ plan.badge }}</span>
            </div>

            <div class="space-y-4">
              <div>
                <div class="flex items-center gap-2">
                  <h3 class="text-base font-extrabold text-foreground">
                    {{ plan.name }}
                  </h3>
                  <span
                    v-if="plan.isActive"
                    class="text-[9px] font-mono px-1.5 py-0.2 rounded bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 font-bold border border-emerald-500/30"
                  >
                    PLAN EN COURS
                  </span>
                </div>
                <p
                  class="text-[11px] text-muted-foreground mt-1 leading-relaxed"
                >
                  {{ plan.desc }}
                </p>
              </div>

              <!-- Période restante mise en valeur si le plan est en cours d'utilisation -->
              <div
                v-if="plan.isActive"
                class="p-2.5 rounded-[4px] bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-between text-xs"
              >
                <div class="flex items-center gap-1.5 text-emerald-700 dark:text-emerald-300 font-semibold">
                  <Clock class="w-3.5 h-3.5" />
                  <span>Temps restant :</span>
                </div>
                <div class="font-mono font-bold text-emerald-600 dark:text-emerald-400">
                  {{ daysRemainingLabel }}
                </div>
              </div>

              <!-- Tarification -->
              <div class="py-2 border-y border-border/60">
                <div class="flex items-baseline gap-1">
                  <span class="text-2xl font-black text-foreground">
                    {{
                      formatPrice(
                        isAnnual ? plan.priceAnnual : plan.priceMonthly,
                      )
                    }}
                  </span>
                  <span
                    v-if="plan.priceMonthly"
                    class="text-[11px] text-muted-foreground font-medium"
                  >
                    / mois
                  </span>
                </div>
                <div
                  v-if="plan.priceMonthly && isAnnual"
                  class="text-[10px] text-emerald-600 dark:text-emerald-400 font-medium mt-0.5"
                >
                  Facturé annuellement (économie de 21 600 FCFA/an)
                </div>
              </div>

              <!-- Fonctionnalités incluses -->
              <div class="space-y-2">
                <div
                  class="text-[10px] uppercase font-bold text-muted-foreground tracking-wider"
                >
                  Inclus dans la formule :
                </div>
                <ul class="space-y-1.5">
                  <li
                    v-for="(feat, idx) in plan.features"
                    :key="idx"
                    class="flex items-start gap-2 text-xs text-foreground/90"
                  >
                    <Check class="w-3.5 h-3.5 text-brand shrink-0 mt-0.5" />
                    <span>{{ feat }}</span>
                  </li>
                </ul>
              </div>
            </div>

            <!-- Bouton d'action -->
            <div class="pt-5 mt-4 border-t border-border/40">
              <!-- Si le plan est actif : bouton grisé et désactivé avec jours restants -->
              <button
                v-if="plan.isActive"
                type="button"
                disabled
                class="w-full py-2.5 rounded-[4px] font-bold text-xs transition flex items-center justify-center gap-1.5 bg-muted/60 text-muted-foreground border border-border cursor-not-allowed opacity-75 select-none"
              >
                <CheckCircle2 class="w-3.5 h-3.5 text-emerald-500" />
                <span>Abonnement en cours ({{ daysRemainingLabel }})</span>
              </button>

              <!-- Si le plan a expiré : bouton réactivé pour renouvellement immédiat -->
              <button
                v-else-if="plan.isCurrentPlan && isLicenceExpired"
                type="button"
                @click="handleSelectPlan(plan.id, false)"
                class="w-full py-2.5 rounded-[4px] font-extrabold text-xs shadow-md transition cursor-pointer flex items-center justify-center gap-1.5 bg-brand hover:bg-brand-dark text-white ring-2 ring-brand/30"
              >
                <span>Renouveler ce Bon Plan</span>
                <ArrowRight class="w-3.5 h-3.5" />
              </button>

              <!-- Cas standard pour souscription -->
              <button
                v-else
                type="button"
                @click="handleSelectPlan(plan.id, false)"
                class="w-full py-2.5 rounded-[4px] font-extrabold text-xs shadow-xs transition cursor-pointer flex items-center justify-center gap-1.5"
                :class="
                  plan.popular
                    ? 'bg-brand hover:bg-brand-dark text-white'
                    : 'bg-muted hover:bg-muted/80 text-foreground'
                "
              >
                <span>Souscrire ce Bon Plan</span>
                <ArrowRight class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
        </div>

        <!-- Section Contact Direct & Moyens de paiement -->
        <div
          class="p-4 rounded-[6px] bg-muted/30 border border-border flex flex-col sm:flex-row items-center justify-between gap-3 text-xs"
        >
          <div class="flex items-center gap-2.5">
            <CreditCard class="w-5 h-5 text-brand shrink-0" />
            <div>
              <div class="font-bold text-foreground">
                Moyens de règlement acceptés
              </div>
              <div class="text-[11px] text-muted-foreground">
                Mobile Money (Orange Money, MTN MoMo, Wave), Virement Bancaire &
                Carte
              </div>
            </div>
          </div>

          <a
            href="https://wa.me/237651727932?text=Bonjour,%20je%20souhaite%20activer%20mon%20plan%20ElevenDrive"
            target="_blank"
            rel="noopener"
            class="px-4 py-2 rounded-[4px] bg-emerald-600 hover:bg-emerald-700 text-white font-bold text-xs flex items-center gap-1.5 transition shrink-0 cursor-pointer shadow-xs"
          >
            <MessageSquare class="w-3.5 h-3.5" />
            <span>Assistance WhatsApp Directe</span>
          </a>
        </div>
      </div>

      <!-- Pied de modal -->
      <div
        class="p-3 bg-muted/40 border-t border-border flex items-center justify-between text-xs"
      >
        <div class="flex items-center gap-2">
          <span class="text-[11px] text-muted-foreground">
            Auto-École : <strong>{{ authStore.schoolName }}</strong>
          </span>
          <span
            v-if="school"
            class="text-[10px] font-mono px-1.5 py-0.2 rounded bg-muted text-muted-foreground border border-border"
          >
            {{ rawTypeLicence }} · {{ daysRemainingLabel }}
          </span>
        </div>
        <button
          type="button"
          @click="emit('close')"
          class="px-4 py-1.5 text-xs font-semibold rounded-[3px] bg-muted hover:bg-muted/80 text-foreground transition cursor-pointer"
        >
          Fermer
        </button>
      </div>
    </div>
  </div>
</template>

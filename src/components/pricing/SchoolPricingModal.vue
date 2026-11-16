<script setup lang="ts">
import { ref, computed } from "vue";
import { useAuthStore } from "../../stores/authStore";
import { useGuichetStore } from "../../stores/guichetStore";
import {
  X,
  Check,
  Sparkles,
  ArrowRight,
  ShieldCheck,
  Zap,
  CreditCard,
  Building2,
  Phone,
  MessageSquare,
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
const isContactModalOpen = ref(false);
const selectedPlanId = ref<"eleven" | "customeleven">("eleven");

const plans = computed(() => [
  {
    id: "eleven",
    name: "Eleven Pro",
    badge: "Recommandé · Bon Plan",
    desc: "La solution tout-en-un pour gérer votre auto-école, vos agences et moniteurs sans limitation.",
    priceMonthly: 12000,
    priceAnnual: 10200, // -15% en annuel
    popular: true,
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
    id: "customeleven",
    name: "Custom Eleven",
    badge: "Sur-Mesure",
    desc: "Déploiement personnalisé pour réseaux d'auto-écoles et flottes multi-régionales.",
    priceMonthly: null,
    priceAnnual: null,
    popular: false,
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

const handleSelectPlan = (planId: "eleven" | "customeleven") => {
  selectedPlanId.value = planId;
  isContactModalOpen.value = true;
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
      <div class="p-5 overflow-y-auto space-y-6 flex-1 text-xs">
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
                  ? 'bg-brand text-white  font-bold'
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
              plan.popular
                ? 'border-brand bg-gradient-to-b from-brand/5 via-card to-card ring-1 ring-brand/50 shadow-2xl'
                : 'border-border bg-card hover:border-border/80',
            ]"
          >
            <!-- Badge populaire -->
            <div
              v-if="plan.popular"
              class="absolute -top-2.5 right-4 px-2.5 py-0.5 rounded-full bg-brand text-white text-[9px] font-extrabold uppercase tracking-wider shadow-sm flex items-center gap-1"
            >
              <Sparkles class="w-2.5 h-2.5" />
              <span>{{ plan.badge }}</span>
            </div>

            <div class="space-y-4">
              <div>
                <h3 class="text-base font-extrabold text-foreground">
                  {{ plan.name }}
                </h3>
                <p
                  class="text-[11px] text-muted-foreground mt-1 leading-relaxed"
                >
                  {{ plan.desc }}
                </p>
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
              <button
                type="button"
                @click="handleSelectPlan(plan.id as any)"
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
        <span class="text-[11px] text-muted-foreground">
          Auto-École : <strong>{{ authStore.schoolName }}</strong>
        </span>
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

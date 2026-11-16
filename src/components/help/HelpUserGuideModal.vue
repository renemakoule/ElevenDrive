<script setup lang="ts">
import { ref, computed } from "vue";
import { useLanguageStore } from "../../stores/languageStore";
import {
  BookOpen,
  X,
  Search,
  UserPlus,
  Receipt,
  Calendar,
  Award,
  FileText,
  Cloud,
  ChevronRight,
  Sparkles,
} from "lucide-vue-next";

defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
}>();

const languageStore = useLanguageStore();
const t = (text: string) => languageStore.t(text);

const activeChapter = ref("inscription");
const searchQuery = ref("");

interface Chapter {
  id: string;
  title: string;
  icon: any;
  summary: string;
  steps: { title: string; desc: string }[];
  proTip: string;
}

const chapters: Chapter[] = [
  {
    id: "inscription",
    title: "1. Enrôlement & Dossier Candidat",
    icon: UserPlus,
    summary:
      "Créer rapidement le dossier d'un nouvel élève avec attribution de formule et gestion de l'acompte initial.",
    steps: [
      {
        title: "Ouvrir le formulaire d'inscription",
        desc: "Appuyez sur la touche F1 ou cliquez sur '+ Nouveau Candidat' sur votre écran d'accueil.",
      },
      {
        title: "Sélectionner la Formule Permis",
        desc: "Choisissez parmi les formules définies par votre Directeur (ex: Permis A, B, C, D) avec les volumes d'heures prévus.",
      },
      {
        title: "Renseigner l'identité et les coordonnées",
        desc: "Nom, prénom, téléphone et date de naissance. Un matricule unique lui sera automatiquement attribué.",
      },
      {
        title: "Enregistrer l'acompte obligatoire",
        desc: "Le système exige l'acompte minimum fixé pour valider le contrat d'inscription.",
      },
    ],
    proTip:
      "Un SMS et un ticket de caisse thermique avec QR code d'accès peuvent être imprimés dès la validation.",
  },
  {
    id: "encaissement",
    title: "2. Encaissement & Tickets ESC/POS",
    icon: Receipt,
    summary:
      "Encaisser les acomptes et soldes restants avec traçabilité intégrale et impression instantanée.",
    steps: [
      {
        title: "Rechercher le candidat",
        desc: "Appuyez sur F3 ou F2, tapez les premières lettres du nom ou le matricule de l'élève.",
      },
      {
        title: "Saisir le montant perçu",
        desc: "Indiquez la somme versée en espèces, Wave, Orange Money ou MTN MoMo.",
      },
      {
        title: "Valider l'opération",
        desc: "La caisse s'incrémente automatiquement et le solde restant de l'élève est recalculé en temps réel.",
      },
      {
        title: "Délivrer le reçu thermique",
        desc: "Le ticket officiel 80mm est généré avec numéro séquentiel et signature électronique.",
      },
    ],
    proTip:
      "Toute tentative d'encaissement inférieur au minimum requis nécessite le PIN de dérogation du Directeur.",
  },
  {
    id: "planning",
    title: "3. Planning Conduite & Cours Code",
    icon: Calendar,
    summary:
      "Planifier les créneaux de conduite pratique avec moniteurs et véhicules assignés.",
    steps: [
      {
        title: "Consulter la disponibilité des moniteurs",
        desc: "Accédez à l'onglet 'Planning Conduite' pour visualiser le calendrier interactif par véhicule et moniteur.",
      },
      {
        title: "Réserver un créneau de 2 heures",
        desc: "Sélectionnez l'élève et validez. Le système vérifie que son solde et ses heures théoriques sont suffisants.",
      },
      {
        title: "Émargement tactile en fin de session",
        desc: "Sur mobile PWA ou tablette, le moniteur valide les compétences abordées et fait signer l'élève.",
      },
    ],
    proTip:
      "Les moniteurs ont accès à leur planning hors-ligne et reçoivent une notification pour chaque nouvelle réservation.",
  },
  {
    id: "examens",
    title: "4. Examens d'État & Listes d'Attente",
    icon: Award,
    summary:
      "Gérer les sessions officielles DGTTC, l'inscription des candidats soldés et le forçage manuel.",
    steps: [
      {
        title: "Créer ou ouvrir une session d'examen",
        desc: "Définissez la date d'examen officiel et le centre d'examen dans l'onglet 'Examens'.",
      },
      {
        title: "Ajout automatique des candidats soldés",
        desc: "Tous les élèves ayant entièrement payé leurs frais de formation sont immédiatement inscrits sur la liste officielle.",
      },
      {
        title: "Gestion de la liste d'attente (Impayés)",
        desc: "Les candidats avec un reste à payer sont placés en attente jusqu'au règlement de leur solde.",
      },
      {
        title: "Bouton de forçage manuel (Dérogation)",
        desc: "Le Directeur peut transférer manuellement un candidat d'attente vers la liste officielle avec motif tracé.",
      },
    ],
    proTip:
      "Le bordereau officiel A4 prêt pour la DGTTC peut être exporté ou imprimé en un clic.",
  },
  {
    id: "caisse",
    title: "5. Clôture de Caisse (Rapport Z)",
    icon: FileText,
    summary:
      "Fermer la session journalière du guichet, réconcilier les fonds et imprimer le rapport financier.",
    steps: [
      {
        title: "Déclencher la clôture",
        desc: "En fin de journée, appuyez sur la touche F4 ou rendez-vous sur 'Point Caisse (Rapport Z)'.",
      },
      {
        title: "Comptage physique des espèces",
        desc: "Saisissez les montants réels en caisse. Le système calcule automatiquement les écarts éventuels.",
      },
      {
        title: "Verrouillage de la session",
        desc: "Une fois validée, la session est clôturée et archivée de façon inaltérable dans le journal d'audit.",
      },
      {
        title: "Impression du ticket Z de synthèse",
        desc: "Sortez le récapitulatif complet sur l'imprimante thermique pour remise au Directeur.",
      },
    ],
    proTip:
      "Le Directeur reçoit instantanément la notification de clôture avec le chiffre d'affaires du jour.",
  },
  {
    id: "cloud",
    title: "6. Synchronisation Neon Cloud & Offline",
    icon: Cloud,
    summary:
      "Continuer de travailler même sans connexion Internet grâce au stockage local chiffré.",
    steps: [
      {
        title: "Fonctionnement 100% Hors-Ligne",
        desc: "Si Internet est coupé, vous pouvez continuer d'inscrire des élèves et d'enregistrer des paiements.",
      },
      {
        title: "File d'attente locale sécurisée",
        desc: "Toutes les opérations sont conservées localement dans IndexedDB avec horodatage strict.",
      },
      {
        title: "Synchronisation automatique au retour réseau",
        desc: "Dès que la connexion réapparaît, vos données fusionnent avec le serveur Neon Postgres sans doublon.",
      },
    ],
    proTip:
      "Vérifiez le badge 'Cloud Connecté' en haut de votre écran pour connaître l'état du serveur.",
  },
];

const selectedChapter = computed(() => {
  return chapters.find((c) => c.id === activeChapter.value) || chapters[0];
});
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-xs p-4 select-none animate-in fade-in duration-150"
    @click.self="emit('close')"
  >
    <div
      class="bg-card border border-border rounded-[3px] max-w-3xl w-full overflow-hidden flex flex-col h-[85vh]"
    >
      <!-- En-tête -->
      <div
        class="p-4 bg-muted/40 border-b border-border flex items-center justify-between shrink-0"
      >
        <div class="flex items-center space-x-2.5">
          <div
            class="w-8 h-8 rounded-[4px] bg-brand-light text-brand-dark flex items-center justify-center font-bold"
          >
            <BookOpen class="w-4 h-4 text-brand" />
          </div>
          <div>
            <div
              class="font-bold text-foreground text-sm flex items-center gap-2"
            >
              <span>{{ t("Guide Pratique d'Utilisation") }}</span>
              <span
                class="text-[10px] px-1.5 py-0.2 rounded bg-brand-light text-brand-dark font-medium"
              >
                {{ t("ElevenDrive Documentation") }}
              </span>
            </div>
            <div class="text-[11px] text-muted-foreground">
              {{
                t(
                  "Toutes les étapes pour maîtriser l'exploitation de votre auto-école",
                )
              }}
            </div>
          </div>
        </div>
        <button
          type="button"
          @click="emit('close')"
          class="p-1 rounded-[4px] text-muted-foreground hover:text-foreground hover:bg-muted transition"
          :title="t('Fermer')"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Corps : Colonne de gauche (Chapitres) + Colonne de droite (Contenu pas-à-pas) -->
      <div class="flex-1 min-h-0 flex overflow-hidden">
        <!-- Sommaire latéral -->
        <div
          class="w-64 border-r border-border bg-muted/20 p-2 overflow-y-auto space-y-1 shrink-0"
        >
          <div
            class="px-2 py-1 text-[10px] font-bold text-muted-foreground uppercase tracking-wider"
          >
            {{ t("Modules du Guide") }}
          </div>
          <button
            v-for="ch in chapters"
            :key="ch.id"
            type="button"
            @click="activeChapter = ch.id"
            class="w-full text-left p-2 rounded-[4px] text-xs transition flex items-center space-x-2.5"
            :class="
              activeChapter === ch.id
                ? 'bg-brand text-white font-semibold'
                : 'text-muted-foreground hover:text-foreground hover:bg-muted/60'
            "
          >
            <component
              :is="ch.icon"
              class="w-4 h-4 shrink-0"
              :class="activeChapter === ch.id ? 'text-white' : 'text-brand'"
            />
            <span class="truncate text-[11px]">{{ t(ch.title) }}</span>
          </button>
        </div>

        <!-- Détail du chapitre actif -->
        <div class="flex-1 p-5 overflow-y-auto space-y-4">
          <!-- Titre & Résumé -->
          <div class="border-b border-border/60 pb-3">
            <h3
              class="text-base font-bold text-foreground flex items-center gap-2"
            >
              <component
                :is="selectedChapter.icon"
                class="w-5 h-5 text-brand"
              />
              <span>{{ t(selectedChapter.title) }}</span>
            </h3>
            <p class="text-xs text-muted-foreground mt-1 leading-relaxed">
              {{ t(selectedChapter.summary) }}
            </p>
          </div>

          <!-- Étapes chronologiques -->
          <div class="space-y-3">
            <div
              v-for="(st, idx) in selectedChapter.steps"
              :key="idx"
              class="flex items-start space-x-3 p-3 rounded-[4px] border border-border/70 bg-card"
            >
              <div
                class="w-6 h-6 rounded-full bg-brand-light text-brand-dark font-bold text-xs flex items-center justify-center shrink-0 mt-0.5"
              >
                {{ idx + 1 }}
              </div>
              <div class="space-y-1">
                <div class="text-xs font-semibold text-foreground">
                  {{ t(st.title) }}
                </div>
                <div class="text-[11px] text-muted-foreground leading-relaxed">
                  {{ t(st.desc) }}
                </div>
              </div>
            </div>
          </div>

          <!-- Astuce Pro -->
          <div
            class="p-3.5 rounded-[4px] bg-amber-500/10 border border-amber-500/30 text-amber-800 dark:text-amber-300 flex items-start space-x-2.5"
          >
            <Sparkles
              class="w-4 h-4 shrink-0 text-amber-600 dark:text-amber-400 mt-0.5"
            />
            <div class="text-xs leading-relaxed">
              <span class="font-bold">{{
                t("Conseil Pro ElevenDrive :")
              }}</span>
              {{ t(selectedChapter.proTip) }}
            </div>
          </div>
        </div>
      </div>

      <!-- Pied de page -->
      <div
        class="p-3 bg-muted/40 border-t border-border flex items-center justify-between text-xs text-muted-foreground shrink-0"
      >
        <span class="text-[11px]">
          {{ t("Besoin d'une formation d'équipe ? Contactez notre support.") }}
        </span>
        <button
          type="button"
          @click="emit('close')"
          class="px-3 py-1.5 rounded-[4px] bg-muted hover:bg-muted/80 text-foreground font-medium text-xs transition"
        >
          {{ t("Fermer le Guide") }}
        </button>
      </div>
    </div>
  </div>
</template>

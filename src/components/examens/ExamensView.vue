<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from "vue";
import { useGuichetStore } from "../../stores/guichetStore";
import { useAuthStore } from "../../stores/authStore";
import { api } from "../../db/api";
import { realtimeSocket } from "../../services/realtimeSocketService";
import BordereauMinisterielModal from "./BordereauMinisterielModal.vue";
import SessionExamenModal from "./SessionExamenModal.vue";
import ExamenEligibiliteTab from "./tabs/ExamenEligibiliteTab.vue";
import ExamenPvTab from "./tabs/ExamenPvTab.vue";
import ExamenNotificationsTab from "./tabs/ExamenNotificationsTab.vue";
import ForcerCandidatModal from "./ForcerCandidatModal.vue";
import type {
  SessionExamen,
  CandidatExamenItem,
  DossierNonEligible,
} from "../../types";
import {
  Award,
  Printer,
  Send,
  ShieldCheck,
  RefreshCw,
  Plus,
  Edit2,
  CalendarPlus,
  AlertCircle,
} from "lucide-vue-next";

const store = useGuichetStore();
const authStore = useAuthStore();
let unsubscribeWs: (() => void) | null = null;

type Onglet = "eligibilite" | "pv" | "notifications";
type Resultat = CandidatExamenItem["resultat_code"];
type ChampResultat = "resultat_code" | "resultat_conduite";

const onglets: { id: Onglet; label: string; icone: typeof Award }[] = [
  { id: "eligibilite", label: "Éligibilité", icone: ShieldCheck },
  { id: "pv", label: "Procès-verbal", icone: Award },
  { id: "notifications", label: "Notifications", icone: Send },
];

const HEURE_CONVOCATION = "07h30";
const telephoneAgence = computed(
  () => authStore.schoolPhone || "Contact Auto-École",
);

const ongletActif = ref<Onglet>("eligibilite");
const bordereauOuvert = ref(false);
const sessionModalOuverte = ref(false);
const isLoading = ref(false);

const session = ref<SessionExamen | null>(null);
const rawCandidats = ref<CandidatExamenItem[]>([]);
const rawIneligibles = ref<DossierNonEligible[]>([]);

const loadExamData = async () => {
  isLoading.value = true;
  try {
    const data = await api.getSessionExamenActive();
    session.value = data.session;
    rawCandidats.value = data.candidats;
    rawIneligibles.value = data.ineligibles;
  } catch (err: any) {
    store.showToast(
      "Erreur chargement session d'examen: " + err.message,
      "error",
    );
  } finally {
    isLoading.value = false;
  }
};

onMounted(async () => {
  await loadExamData();
  unsubscribeWs = realtimeSocket.subscribe((msg) => {
    if (
      msg.type === "EXAMENS_SESSION_UPDATED" ||
      msg.type === "PAIEMENT_ENREGISTRE" ||
      msg.type === "CANDIDAT_DOSSIER_UPDATED" ||
      msg.type === "LIVRET_EMARGEMENT_UPDATED" ||
      msg.type === "COURS_CODE_UPDATED"
    ) {
      loadExamData();
    }
  });
});

onUnmounted(() => {
  if (unsubscribeWs) {
    unsubscribeWs();
    unsubscribeWs = null;
  }
});

const candidats = computed<CandidatExamenItem[]>(() => rawCandidats.value);
const autresDossiers = computed<DossierNonEligible[]>(
  () => rawIneligibles.value,
);
const publie = computed(
  () => session.value?.statut_session === "RESULTATS_PUBLIES",
);
const canForce = computed(
  () =>
    authStore.currentRole === "manager" ||
    authStore.currentRole === "super_admin",
);

const dateSession = computed(() =>
  session.value?.date_examen
    ? new Date(`${session.value.date_examen}T12:00:00`)
    : null,
);

const dateLongue = computed(() =>
  dateSession.value
    ? dateSession.value.toLocaleDateString("fr-FR", {
        weekday: "long",
        day: "numeric",
        month: "long",
        year: "numeric",
      })
    : "Date non définie",
);

const dateNumerique = computed(() =>
  dateSession.value
    ? dateSession.value.toLocaleDateString("fr-FR")
    : "Non définie",
);

const ouvrirPlanification = () => {
  sessionModalOuverte.value = true;
};

const onSessionSauvegardee = async () => {
  sessionModalOuverte.value = false;
  await loadExamData();
};

const onSessionSupprimee = async () => {
  sessionModalOuverte.value = false;
  session.value = null;
  await loadExamData();
};

/* Forçage Manuel / Dérogation Direction */
const forcerModalOuvert = ref(false);
const dossierAForcer = ref<DossierNonEligible | null>(null);
const isSubmittingForcage = ref(false);
const erreurForcage = ref("");

const ouvrirForcage = (d: DossierNonEligible) => {
  dossierAForcer.value = d;
  erreurForcage.value = "";
  forcerModalOuvert.value = true;
};

const fermerForcage = () => {
  forcerModalOuvert.value = false;
  dossierAForcer.value = null;
  erreurForcage.value = "";
};

const validerForcage = async ({
  motif,
  codePin,
}: {
  motif: string;
  codePin: string;
}) => {
  if (!dossierAForcer.value) return;
  if (!session.value?.id) {
    store.showToast("Aucune session active", "error");
    return;
  }
  if (!codePin || codePin.length < 4) {
    erreurForcage.value =
      "Veuillez saisir votre code PIN Directeur (4 chiffres).";
    return;
  }

  isSubmittingForcage.value = true;
  erreurForcage.value = "";
  try {
    const pinValide = await api.verifierCodeManager(codePin);
    if (!pinValide) {
      erreurForcage.value = "Code PIN Directeur incorrect.";
      return;
    }

    const dossierId =
      dossierAForcer.value.dossier_id || dossierAForcer.value.numero_dossier;

    await api.forcerCandidatExamen({
      session_examen_id: session.value.id,
      dossier_id: dossierId,
      motif,
      code_pin: codePin,
    });

    store.showToast(
      `Dérogation accordée : ${dossierAForcer.value.nom} a été forcé et inscrit au bordereau d'État.`,
      "success",
    );
    fermerForcage();
    await loadExamData();
  } catch (err: any) {
    erreurForcage.value = err.message || "Erreur lors du forçage de l'élève.";
  } finally {
    isSubmittingForcage.value = false;
  }
};

/* Résultats & PV */
const libelleResultat: Record<string, string> = {
  ADMIS: "Admis",
  AJOURNE: "Ajourné",
  ABSENT: "Absent",
  EN_ATTENTE: "En attente",
};

const SUITE: Partial<Record<Resultat, Resultat>> = {
  ADMIS: "AJOURNE",
  AJOURNE: "ABSENT",
};

const recalculerFinal = (c: CandidatExamenItem) => {
  if (c.resultat_code === "ADMIS" && c.resultat_conduite === "ADMIS")
    c.resultat_final = "ADMIS";
  else if (c.resultat_code === "ABSENT" || c.resultat_conduite === "ABSENT")
    c.resultat_final = "ABSENT";
  else c.resultat_final = "AJOURNE";
};

const changerResultat = (c: CandidatExamenItem, champ: ChampResultat) => {
  c[champ] = SUITE[c[champ]] ?? "ADMIS";
  recalculerFinal(c);
};

const validerPv = async () => {
  if (!session.value?.id) {
    store.showToast(
      "Veuillez planifier une session d'examen avant de valider le procès-verbal.",
      "info",
    );
    return;
  }

  try {
    const payload = {
      session_examen_id: session.value.id,
      resultats: candidats.value.map((c) => ({
        candidat_examen_id: c.dossier_id,
        resultat_code: c.resultat_code,
        resultat_conduite: c.resultat_conduite,
        resultat_final: c.resultat_final,
      })),
    };
    await api.enregistrerResultatsExamen(payload);
    session.value.statut_session = "RESULTATS_PUBLIES";
    store.showToast(
      "Procès-verbal officiel validé et résultats enregistrés en base !",
      "success",
    );
  } catch (err: any) {
    store.showToast("Erreur validation procès-verbal: " + err.message, "error");
  }
};

/* WhatsApp */
const ouvrirWhatsApp = (c: CandidatExamenItem, message: string) => {
  const numero = c.telephone_whatsapp.replace(/\D/g, "");
  window.open(
    `https://wa.me/${numero}?text=${encodeURIComponent(message)}`,
    "_blank",
  );
};

const envoyerConvocation = (c: CandidatExamenItem) => {
  if (!session.value) {
    store.showToast(
      "Aucune session active pour générer les convocations.",
      "info",
    );
    return;
  }

  ouvrirWhatsApp(
    c,
    [
      "Convocation à l'examen officiel du permis de conduire",
      "",
      `Bonjour ${c.nom} ${c.prenoms},`,
      `L'établissement ${authStore.schoolName} vous convoque à l'examen officiel du permis (catégorie ${session.value.categorie_permis}) :`,
      "",
      `Date : ${dateLongue.value} à ${HEURE_CONVOCATION}`,
      `Centre : ${session.value.centre_examen}`,
      `N° de table : *${c.numero_table}*`,
      `Dossier : ${c.numero_dossier}`,
      "",
      `Présentez-vous impérativement avec votre pièce d'identité originale (${c.type_piece_id} n° ${c.numero_piece_id}).`,
      "",
      `${authStore.schoolName} · ${telephoneAgence.value}`,
    ].join("\n"),
  );
  store.showToast(`Convocation prête pour ${c.nom}`, "success");
};

const envoyerResultat = (c: CandidatExamenItem) => {
  if (!session.value) return;

  const entete = `Bonjour ${c.nom} ${c.prenoms},`;
  const agenceNom = authStore.currentUser?.agence || "l'agence centrale";
  const message =
    c.resultat_final === "ADMIS"
      ? [
          "Résultat de l'examen officiel du permis de conduire",
          "",
          entete,
          `Félicitations, vous êtes déclaré(e) *ADMIS(E)* à l'examen du permis (catégorie ${session.value.categorie_permis}, session du ${dateNumerique.value}) !`,
          "",
          `Épreuve théorique (code) : ${libelleResultat[c.resultat_code]}`,
          `Épreuve pratique (conduite) : ${libelleResultat[c.resultat_conduite]}`,
          "",
          `Présentez-vous à ${agenceNom} muni(e) de 2 photos d'identité pour retirer votre attestation provisoire de conduite.`,
          "",
          authStore.schoolName,
        ].join("\n")
      : [
          "Résultat de l'examen officiel du permis de conduire",
          "",
          entete,
          `La commission d'examen vous a déclaré(e) *AJOURNÉ(E)* lors de la session du ${dateNumerique.value}.`,
          "",
          "Passez au guichet pour planifier une séance de révision et vous présenter à la prochaine session ministérielle.",
          "",
          authStore.schoolName,
        ].join("\n");

  ouvrirWhatsApp(c, message);
  store.showToast(`Notification de résultat prête pour ${c.nom}`, "success");
};
</script>

<template>
  <div class="flex flex-1 select-none flex-col space-y-3 overflow-y-auto p-1">
    <!-- En-tête -->
    <div
      class="flex flex-wrap items-center justify-between gap-2.5 rounded-[5px] border border-border bg-card px-3 py-2.5"
    >
      <div>
        <h2 class="text-sm font-semibold leading-tight text-foreground">
          Examens et bordereaux d'État
        </h2>
        <p class="text-[10px] text-muted-foreground">
          Éligibilité ministérielle, bordereau de présentation A4 et
          procès-verbal
        </p>
      </div>

      <div class="flex items-center gap-2">
        <button
          type="button"
          @click="loadExamData"
          :disabled="isLoading"
          class="inline-flex h-7 w-7 items-center justify-center rounded-md border border-border bg-card text-muted-foreground hover:bg-muted hover:text-foreground transition disabled:opacity-50"
          title="Rafraîchir les données"
        >
          <RefreshCw
            class="h-3.5 w-3.5"
            :class="{ 'animate-spin': isLoading }"
          />
        </button>

        <template v-if="session && session.id">
          <div class="text-right leading-tight hidden sm:block">
            <div class="text-[9px] text-muted-foreground">Session</div>
            <div class="font-mono text-[10px] font-medium text-foreground">
              {{ session.code_session }}
            </div>
          </div>
          <div class="text-right leading-tight hidden sm:block">
            <div class="text-[9px] text-muted-foreground">Date</div>
            <div class="font-mono text-[10px] font-medium text-foreground">
              {{ dateNumerique }}
            </div>
          </div>
          <button
            type="button"
            @click="ouvrirPlanification"
            class="flex items-center gap-1 rounded-md border border-border bg-card px-2 py-1 text-[10px] font-medium text-foreground hover:bg-muted transition"
            title="Modifier les détails de cette session"
          >
            <Edit2 class="h-3 w-3 text-brand" />
            Modifier
          </button>
          <button
            type="button"
            @click="bordereauOuvert = true"
            class="flex items-center gap-1 rounded-md bg-brand-dark px-2.5 py-1 text-[10px] font-medium text-white hover:bg-brand transition"
          >
            <Printer class="h-3 w-3" />
            Aperçu A4
          </button>
        </template>

        <template v-else>
          <span
            class="rounded bg-muted px-2 py-1 font-mono text-[10px] text-muted-foreground"
          >
            Aucune session active
          </span>
          <button
            type="button"
            @click="ouvrirPlanification"
            class="flex items-center gap-1.5 rounded-md bg-brand-dark px-2.5 py-1 text-[10px] font-semibold text-white hover:bg-brand transition shadow-sm"
          >
            <CalendarPlus class="h-3.5 w-3.5" />
            Planifier une session d'examen
          </button>
        </template>
      </div>
    </div>

    <!-- Bannière aucune session -->
    <div
      v-if="!session || !session.id"
      class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3 rounded-lg border border-dashed border-amber-300/80 bg-amber-50/70 p-3.5 text-xs text-amber-900 dark:border-amber-800/60 dark:bg-amber-950/20 dark:text-amber-200"
    >
      <div class="flex items-start gap-2.5">
        <AlertCircle
          class="h-4 w-4 shrink-0 text-amber-600 dark:text-amber-400 mt-0.5"
        />
        <div>
          <div class="font-semibold text-amber-950 dark:text-amber-100">
            Aucune session d'examen planifiée
          </div>
          <p class="text-[11px] text-amber-800/90 dark:text-amber-300/90">
            Fixez la date officielle de passage, le centre d'examen et
            l'inspecteur pour générer le bordereau ministériel et valider le
            procès-verbal.
          </p>
        </div>
      </div>
      <button
        type="button"
        @click="ouvrirPlanification"
        class="inline-flex shrink-0 items-center gap-1.5 rounded-md bg-amber-700 hover:bg-amber-800 px-3 py-1.5 text-xs font-semibold text-white transition shadow-sm"
      >
        <Plus class="h-3.5 w-3.5" />
        Planifier maintenant
      </button>
    </div>

    <!-- Onglets -->
    <div class="flex items-center gap-1.5">
      <button
        v-for="o in onglets"
        :key="o.id"
        type="button"
        @click="ongletActif = o.id"
        :class="[
          'flex items-center gap-1 rounded-md border px-2 py-1 text-[10px] font-medium transition',
          ongletActif === o.id
            ? 'border-brand-border bg-brand-light text-brand-dark font-semibold'
            : 'border-border bg-card text-muted-foreground hover:bg-muted',
        ]"
      >
        <component
          :is="o.icone"
          class="h-3 w-3"
          :class="ongletActif === o.id ? 'text-brand' : 'text-muted-foreground'"
        />
        {{ o.label }}
      </button>
    </div>

    <!-- Contenu -->
    <div
      class="flex flex-1 flex-col overflow-hidden rounded-[5px] border border-border bg-card"
    >
      <ExamenEligibiliteTab
        v-if="ongletActif === 'eligibilite'"
        :candidats="candidats"
        :autres-dossiers="autresDossiers"
        :has-active-session="Boolean(session && session.id)"
        :can-force="canForce"
        :is-loading="isLoading"
        @force="ouvrirForcage"
      />

      <ExamenPvTab
        v-else-if="ongletActif === 'pv'"
        :has-active-session="Boolean(session && session.id)"
        :candidats="candidats"
        :publie="publie"
        :libelle-resultat="libelleResultat"
        @planifier="ouvrirPlanification"
        @toggle-result="changerResultat"
        @validate-pv="validerPv"
      />

      <ExamenNotificationsTab
        v-else
        :has-active-session="Boolean(session && session.id)"
        :candidats="candidats"
        :publie="publie"
        :libelle-resultat="libelleResultat"
        @planifier="ouvrirPlanification"
        @send-convocation="envoyerConvocation"
        @send-result="envoyerResultat"
      />
    </div>

    <!-- Modales -->
    <SessionExamenModal
      v-if="sessionModalOuverte"
      :session="session"
      @close="sessionModalOuverte = false"
      @saved="onSessionSauvegardee"
      @deleted="onSessionSupprimee"
    />

    <BordereauMinisterielModal
      v-if="bordereauOuvert && session"
      :session="session"
      :candidats="candidats"
      @close="bordereauOuvert = false"
    />

    <ForcerCandidatModal
      :is-open="forcerModalOuvert"
      :dossier="dossierAForcer"
      :is-submitting="isSubmittingForcage"
      :error-msg="erreurForcage"
      @close="fermerForcage"
      @confirm="validerForcage"
    />
  </div>
</template>

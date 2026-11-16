<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { useAuthStore } from "../../../stores/authStore";
import { authBackend } from "../../../services/authBackendService";
import { supportService } from "../../../services/supportService";
import type { CreateFeedbackPayload } from "../../../types";
import {
  Star,
  Sparkles,
  Building2,
  User,
  Tag,
  AlertTriangle,
  Send,
  Loader2,
  CheckCircle2,
} from "lucide-vue-next";

const authStore = useAuthStore();

// Étoiles : 0 par défaut (vides) si aucune note n'est encore mémorisée
const selectedStars = ref(0);
const hoverStars = ref(0);
const feedbackCategory = ref<
  "ergonomie" | "guichet" | "pwa_mobile" | "examens" | "support" | "general"
>("general");
const feedbackComment = ref("");
const feedbackSchoolName = ref("");
const feedbackAuthor = ref("");
const feedbackRole = ref("");
const isSubmittingFeedback = ref(false);
const feedbackSuccess = ref(false);
const feedbackError = ref("");

const getStorageKey = () => {
  const session = authBackend.getActiveSession();
  const userId =
    session?.user?.id ||
    authStore.currentUser?.id ||
    session?.user?.email ||
    authStore.currentUser?.email ||
    feedbackAuthor.value.trim();
  const schoolId =
    session?.auto_ecole?.id ||
    authStore.currentSchool?.id ||
    authStore.currentSchool?.code_unique ||
    feedbackSchoolName.value.trim();

  return `autoschool_feedback_stars_${schoolId || "global"}_${userId || "user"}`;
};

const loadSavedRating = () => {
  try {
    const key = getStorageKey();
    let saved = localStorage.getItem(key);
    if (!saved) {
      saved = localStorage.getItem("autoschool_last_feedback_stars");
    }
    if (saved) {
      const parsed = parseInt(saved, 10);
      if (parsed >= 1 && parsed <= 5) {
        selectedStars.value = parsed;
        return;
      }
    }
  } catch {
    // ignore
  }
  selectedStars.value = 0; // vide par défaut
};

const setRating = (star: number) => {
  selectedStars.value = star;
  feedbackError.value = "";
  try {
    const key = getStorageKey();
    localStorage.setItem(key, star.toString());
    localStorage.setItem("autoschool_last_feedback_stars", star.toString());
  } catch {
    // ignore
  }
};

const resolveIdentity = () => {
  const session = authBackend.getActiveSession();

  // 1. Établissement exact
  const schoolName =
    session?.auto_ecole?.nom_commercial ||
    authStore.currentSchool?.nom_commercial ||
    (authStore.schoolName !== "ElevenDrive" ? authStore.schoolName : "") ||
    session?.auto_ecole?.code_unique ||
    authStore.schoolName ||
    "Auto-École Partenaire";

  feedbackSchoolName.value = schoolName;

  // 2. Auteur connecté réel
  let fullName = "";
  let roleTitle = "";

  if (session?.user) {
    const p = session.user.prenom || "";
    const n = session.user.nom || "";
    fullName = `${p} ${n}`.trim();
    roleTitle = session.user.roleLabel || session.user.role || "";
  } else if (authStore.customUserProfile) {
    const p = authStore.customUserProfile.prenom || "";
    const n = authStore.customUserProfile.nom || "";
    fullName = `${p} ${n}`.trim();
    roleTitle = authStore.customUserProfile.roleLabel || "";
  } else if (authStore.currentUser) {
    const p = authStore.currentUser.prenom || "";
    const n = authStore.currentUser.nom || "";
    fullName = `${p} ${n}`.trim();
    roleTitle = authStore.currentUser.roleLabel || "";
  }

  if (!roleTitle) {
    roleTitle =
      authStore.rolePermissions?.roleTitle ||
      (authStore.currentRole === "manager" ? "Directeur" : "Secrétaire");
  }

  if (!fullName || fullName === "Guichet Secrétaire") {
    fullName = "Secrétariat Guichet";
  } else if (fullName === "Gérant Directeur") {
    fullName = "Direction Générale";
  }

  feedbackAuthor.value = fullName;
  feedbackRole.value = roleTitle;
};

const getStarLabel = (stars: number) => {
  switch (stars) {
    case 1:
      return "Très insatisfait · Système difficile ou bloquant";
    case 2:
      return "Insatisfait · Des lacunes importantes";
    case 3:
      return "Moyen · Fonctionnel mais perfectible";
    case 4:
      return "Très bien · Agréable et efficace";
    case 5:
      return "Excellent ! · Solution complète et performante";
    default:
      return "Aucune note sélectionnée · Cliquez sur les étoiles pour noter (1 à 5)";
  }
};

const submitReview = async () => {
  if (selectedStars.value < 1) {
    feedbackError.value =
      "Veuillez attribuer au moins 1 étoile pour donner votre avis.";
    return;
  }

  if (!feedbackComment.value.trim()) {
    feedbackError.value = "Veuillez saisir votre commentaire ou suggestion.";
    return;
  }

  feedbackError.value = "";
  isSubmittingFeedback.value = true;

  try {
    const session = authBackend.getActiveSession();
    const schoolID =
      authStore.currentSchool?.id ||
      session?.auto_ecole?.id ||
      authStore.currentUser?.auto_ecole_id;

    const schoolLabel = feedbackSchoolName.value.trim();
    const authorLabel = feedbackAuthor.value.trim();
    const composedAuthor =
      authorLabel && schoolLabel
        ? `${authorLabel} · ${schoolLabel}`
        : authorLabel || schoolLabel || "Responsable Auto-École";

    const payload: CreateFeedbackPayload = {
      auto_ecole_id: schoolID,
      auteur: composedAuthor,
      role: feedbackRole.value || authStore.currentRole || "collaborateur",
      etoiles: selectedStars.value,
      categorie: feedbackCategory.value,
      commentaire: feedbackComment.value.trim(),
    };

    await supportService.submitFeedback(payload);
    // Mémoriser la note accordée
    setRating(selectedStars.value);
    feedbackSuccess.value = true;
    feedbackComment.value = "";
  } catch (err: any) {
    feedbackError.value =
      err?.message || "Impossible de transmettre votre avis pour le moment.";
  } finally {
    isSubmittingFeedback.value = false;
  }
};

const resetFeedbackForm = () => {
  feedbackSuccess.value = false;
  feedbackComment.value = "";
  feedbackError.value = "";
  loadSavedRating();
};

onMounted(() => {
  resolveIdentity();
  loadSavedRating();
});

watch(
  () => [authStore.currentSchool, authStore.currentUser, authStore.currentRole],
  () => {
    resolveIdentity();
    loadSavedRating();
  },
);
</script>

<template>
  <div>
    <!-- Écran de confirmation de l'avis -->
    <div
      v-if="feedbackSuccess"
      class="py-12 px-6 text-center space-y-4 bg-muted/20 border border-emerald-500/30 rounded-[4px] animate-in zoom-in-95"
    >
      <div
        class="w-14 h-14 mx-auto rounded-full bg-emerald-500/10 text-emerald-600 flex items-center justify-center"
      >
        <CheckCircle2 class="w-8 h-8" />
      </div>

      <div class="space-y-1">
        <h3 class="font-bold text-foreground text-base">
          Merci pour votre avis précieux !
        </h3>
        <p class="text-xs text-muted-foreground max-w-md mx-auto">
          Votre note de
          <span class="font-bold text-amber-500"
            >{{ selectedStars }} / 5 étoiles</span
          >
          et vos commentaires ont été transmis directement au
          Super-Administrateur pour continuer d'améliorer la plateforme.
        </p>
      </div>

      <button
        type="button"
        @click="resetFeedbackForm"
        class="px-4 py-2 rounded-[3px] bg-brand text-white font-bold text-xs hover:bg-brand-dark transition cursor-pointer"
      >
        Envoyer un autre avis
      </button>
    </div>

    <!-- Formulaire de notation -->
    <form v-else @submit.prevent="submitReview" class="space-y-4 p-1">
      <div class="space-y-1">
        <div
          class="font-bold text-foreground text-sm flex items-center gap-1.5"
        >
          <Sparkles class="w-4 h-4 text-amber-500" />
          <span>Votre avis compte pour faire évoluer ElevenDrive</span>
        </div>
        <p class="text-[11px] text-muted-foreground">
          Partagez votre expérience sur la rapidité, l'ergonomie, les
          fonctionnalités guichet et vos suggestions de perfectionnement.
        </p>
      </div>

      <!-- Widget d'étoiles interactif (1 à 5) -->
      <div
        class="p-3.5 rounded-[4px] border border-border bg-muted/30 space-y-2"
      >
        <label
          class="block text-[11px] font-bold text-foreground uppercase tracking-wide flex items-center justify-between"
        >
          <span>Votre note globale :</span>
          <span
            v-if="selectedStars > 0"
            class="text-[10px] text-amber-600 dark:text-amber-400 font-semibold normal-case flex items-center gap-1 bg-amber-500/10 px-1.5 py-0.5 rounded-[3px]"
          >
            <Sparkles class="w-3 h-3" />
            <span>Note mémorisée</span>
          </span>
        </label>

        <div class="flex items-center gap-2">
          <div class="flex items-center gap-1" @mouseleave="hoverStars = 0">
            <button
              v-for="star in 5"
              :key="star"
              type="button"
              @mouseenter="hoverStars = star"
              @click="setRating(star)"
              class="p-1 transition-transform hover:scale-125 focus:outline-none cursor-pointer"
              :title="getStarLabel(star)"
            >
              <Star
                class="w-6 h-6 transition-colors duration-150"
                :class="
                  star <= (hoverStars || selectedStars)
                    ? 'text-amber-500 fill-amber-500'
                    : 'text-muted-foreground/30 hover:text-amber-400/80'
                "
              />
            </button>
          </div>

          <span class="text-xs font-bold text-foreground ml-2 font-mono">
            {{
              (hoverStars || selectedStars) > 0
                ? `${hoverStars || selectedStars} / 5`
                : "0 / 5"
            }}
          </span>
        </div>

        <div
          class="text-[11px] font-medium transition-colors"
          :class="
            (hoverStars || selectedStars) > 0
              ? 'text-amber-600 dark:text-amber-400'
              : 'text-muted-foreground italic'
          "
        >
          {{ getStarLabel(hoverStars || selectedStars) }}
        </div>
      </div>

      <!-- Grille Établissement, Auteur & Module -->
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
        <!-- 1. Établissement exact -->
        <div class="space-y-1">
          <label
            class="block text-[11px] font-bold text-foreground flex items-center justify-between"
          >
            <span class="flex items-center gap-1">
              <Building2 class="w-3.5 h-3.5 text-brand" />
              <span>Établissement</span>
            </span>
            <span
              class="text-[9px] text-emerald-600 dark:text-emerald-400 font-semibold flex items-center gap-0.5"
            >
              <CheckCircle2 class="w-2.5 h-2.5" />
              <span>Actif</span>
            </span>
          </label>
          <input
            v-model="feedbackSchoolName"
            type="text"
            required
            placeholder="Nom de votre auto-école"
            class="w-full px-2.5 py-1.5 rounded-[3px] border border-border bg-muted/40 font-semibold text-foreground text-xs focus:ring-1 focus:ring-brand outline-none"
          />
        </div>

        <!-- 2. Auteur connecté & Rôle -->
        <div class="space-y-1">
          <label
            class="block text-[11px] font-bold text-foreground flex items-center justify-between"
          >
            <span class="flex items-center gap-1">
              <User class="w-3.5 h-3.5 text-muted-foreground" />
              <span>Auteur connecté</span>
            </span>
            <span
              class="text-[9px] px-1 py-0.2 rounded bg-muted text-muted-foreground font-mono font-medium"
            >
              {{ feedbackRole }}
            </span>
          </label>
          <input
            v-model="feedbackAuthor"
            type="text"
            required
            placeholder="Nom & Prénom de l'auteur"
            class="w-full px-2.5 py-1.5 rounded-[3px] border border-border bg-background text-foreground text-xs focus:ring-1 focus:ring-brand outline-none"
          />
        </div>

        <!-- 3. Module concerné -->
        <div class="space-y-1">
          <label
            class="block text-[11px] font-bold text-foreground flex items-center gap-1"
          >
            <Tag class="w-3.5 h-3.5 text-muted-foreground" />
            <span>Module concerné</span>
          </label>
          <select
            v-model="feedbackCategory"
            class="w-full px-2 py-1.5 rounded-[3px] border border-border bg-background text-foreground text-xs focus:ring-1 focus:ring-brand outline-none"
          >
            <option value="general">Général & Expérience globale</option>
            <option value="guichet">Guichet & Encaissement Caisse</option>
            <option value="ergonomie">Design & Ergonomie Desktop</option>
            <option value="examens">Examens & Progression Candidats</option>
            <option value="pwa_mobile">Application Mobile / PWA</option>
            <option value="support">Assistance & Support Client</option>
          </select>
        </div>
      </div>

      <!-- Champ Commentaire -->
      <div class="space-y-1">
        <label class="block text-[11px] font-bold text-foreground">
          Votre retour détaillé & suggestions d'amélioration
        </label>
        <textarea
          v-model="feedbackComment"
          rows="4"
          required
          placeholder="Ex : Le nouveau module de reçu thermique est très rapide. Pourrions-nous ajouter un filtre supplémentaire par date sur le guichet ?..."
          class="w-full px-3 py-2 rounded-[3px] border border-border bg-background text-foreground text-xs focus:ring-1 focus:ring-brand outline-none resize-none leading-relaxed"
        ></textarea>
      </div>

      <!-- Message d'erreur éventuel -->
      <div
        v-if="feedbackError"
        class="p-2.5 rounded-[3px] bg-rose-500/10 border border-rose-500/30 text-rose-600 dark:text-rose-400 text-[11px] flex items-center gap-2"
      >
        <AlertTriangle class="w-4 h-4 shrink-0" />
        <span>{{ feedbackError }}</span>
      </div>

      <!-- Bouton d'envoi -->
      <div class="flex items-center justify-between pt-2">
        <span class="text-[10px] text-muted-foreground">
          Vos retours sont analysés en temps réel par l'administration.
        </span>

        <button
          type="submit"
          :disabled="isSubmittingFeedback"
          class="px-4 py-2 rounded-[3px] bg-brand text-white font-bold text-xs shadow-xs hover:bg-brand-dark transition cursor-pointer flex items-center gap-1.5 disabled:opacity-50"
        >
          <Loader2
            v-if="isSubmittingFeedback"
            class="w-3.5 h-3.5 animate-spin"
          />
          <Send v-else class="w-3.5 h-3.5" />
          <span>{{
            isSubmittingFeedback ? "Transmission..." : "Transmettre mon avis"
          }}</span>
        </button>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import type { UserRole } from "../../../types";
import PwaHomeRoleCard from "./PwaHomeRoleCard.vue";
import { GraduationCap, Smartphone, ShieldCheck } from "lucide-vue-next";

const emit = defineEmits<{
  (e: "navigate-auth", role?: UserRole): void;
}>();

const handleSelect = (role: UserRole) => {
  emit("navigate-auth", role);
};
</script>

<template>
  <section
    class="relative shrink-0 w-full py-10 max-w-6xl mx-auto px-4 sm:px-6"
  >
    <div class="border border-primary rounded-[4px] py-8 px-10">
      <div class="text-center mb-8">
        <h3 class="text-lg sm:text-xl font-bold text-foreground">
          Choisissez votre espace de connexion
        </h3>
        <p class="text-xs sm:text-sm text-muted-foreground mt-1">
          Accédez directement à vos outils pédagogiques et administratifs selon
          votre profil.
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-5">
        <!-- 1. Espace Candidat & Élève -->
        <PwaHomeRoleCard
          role="eleve"
          title="Espace Candidat & Élève"
          tag="Apprenant"
          description="Consultez votre livret d'apprentissage, réservez vos séances de conduite, suivez votre progression aux 29 compétences et révisez les séries d'examens d'État 24h/24."
          :icon="GraduationCap"
          color="purple"
          :features="[
            'Livret tactile & notation en temps réel',
            'Planning des leçons & réservation en ligne',
            'Historique des paiements & reçus guichet',
          ]"
          action-label="Connexion Candidat"
          @select="handleSelect"
        />

        <!-- 2. Espace Moniteur de Conduite -->
        <PwaHomeRoleCard
          role="moniteur"
          title="Espace Moniteur de Conduite"
          tag="Instructeur"
          description="Votre livret de conduite tactile sur smartphone et tablette PWA. Évaluez les élèves, validez les heures de conduite et faites signer électroniquement sur écran."
          :icon="Smartphone"
          color="sky"
          :features="[
            'Planning du jour & gestion des créneaux',
            'Émargement tactile candidat & instructeur',
            'Mode 100% hors-ligne pour la route',
          ]"
          action-label="Connexion Moniteur"
          @select="handleSelect"
        />

        <!-- 3. Espace Direction & Supervision -->
        <PwaHomeRoleCard
          role="manager"
          title="Espace Direction & Supervision"
          tag="Direction"
          description="Supervisez les opérations de votre établissement, l'activité de vos moniteurs, les encaissements guichet, et validez les dérogations PIN depuis n'importe où."
          :icon="ShieldCheck"
          color="indigo"
          :features="[
            'Tableau de bord financier & guichets',
            'Gestion d\'équipe, moniteurs & secrétaires',
            'Contrôle des bordereaux et examens',
          ]"
          action-label="Connexion Direction"
          @select="handleSelect"
        />
      </div>
    </div>
  </section>
</template>

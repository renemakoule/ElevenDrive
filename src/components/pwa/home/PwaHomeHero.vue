<script setup lang="ts">
import {
  Sparkles,
  MapPin,
  Building2,
  Award,
  ArrowRight,
  Download,
} from "lucide-vue-next";
import { motion, useReducedMotion } from "motion-v";

interface Props {
  schoolName: string;
  schoolAddress: string;
  schoolAgrement: string;
  agencesCount: number;
}

defineProps<Props>();

const emit = defineEmits<{
  (e: "navigate-auth"): void;
  (e: "install-pwa"): void;
}>();

// Respect de prefers-reduced-motion
const reduceMotion = useReducedMotion();

// Nombre de lignes de points (16px chacune) : ajustez selon la hauteur du hero
const ROWS = 32;
</script>

<template>
  <section
    class="relative z-10 shrink-0 w-full overflow-hidden border-b border-border bg-gradient-to-b from-brand/5 via-background to-background py-10 sm:py-16"
  >
    <!-- Trame de fond géométrique estompée -->
    <div
      class="pointer-events-none absolute inset-0 overflow-hidden text-foreground/20"
      style="
        mask-image: radial-gradient(
          ellipse at center,
          black 20%,
          transparent 75%
        );
        -webkit-mask-image: radial-gradient(
          ellipse at center,
          black 20%,
          transparent 75%
        );
      "
    >
      <motion.div
        v-for="i in ROWS"
        :key="i"
        class="w-full"
        :style="{
          height: '16px',
          backgroundImage: 'radial-gradient(currentColor 1px, transparent 1px)',
          backgroundSize: '16px 16px',
        }"
        :initial="{ y: 0, opacity: 0.35 }"
        :animate="
          reduceMotion
            ? { opacity: 0.6 }
            : { y: [0, -6, 0, 6, 0], opacity: [0.3, 1, 0.3] }
        "
        :transition="{
          duration: 4,
          repeat: Infinity,
          ease: 'easeInOut',
          delay: i * 0.12,
        }"
      />
    </div>

    <div
      class="pointer-events-none absolute -left-24 -top-24 h-80 w-80 rounded-full bg-brand/15 blur-3xl"
    ></div>
    <div
      class="pointer-events-none absolute -bottom-32 -right-24 h-96 w-96 rounded-full bg-brand/10 blur-3xl"
    ></div>

    <!-- Contenu Principal Hero -->
    <div class="relative z-10 max-w-4xl mx-auto px-4 sm:px-6 text-center">
      <!-- Tag Portail officiel -->
      <div
        class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-muted-foreground text-muted border border-brand-border mb-4"
      >
        <Sparkles class="w-3.5 h-3.5 text-muted" />
        <span>Portail Numérique Officiel</span>
      </div>

      <h2
        class="text-2xl sm:text-4xl mb-8 font-extrabold text-foreground tracking-tight leading-tight"
      >
        Bienvenue à {{ schoolName }}
      </h2>

      <div class="border border-primary bg-background rounded-[4px] py-2 mb-8">
        <p
          class="mt-3 max-w-2xl mx-auto text-xs sm:text-base text-muted-foreground leading-relaxed"
        >
          Votre formation au permis de conduire modernisée. Accédez à votre
          livret tactile de conduite, réservez vos séances, révisez vos cours de
          code et consultez vos résultats d'examen directement depuis votre
          smartphone.
        </p>
      </div>

      <!-- Pastilles coordonnées -->
      <div
        class="mt-6 flex flex-wrap items-center justify-center gap-2 sm:gap-3 text-xs text-muted"
      >
        <div
          class="inline-flex items-center gap-1.5 rounded-[4px] border border-border bg-muted-foreground px-3 py-1.5"
        >
          <MapPin class="w-3.5 h-3.5 text-primary" />
          <span>{{ schoolAddress || "Siège Principal" }}</span>
        </div>
        <div
          class="inline-flex items-center gap-1.5 rounded-[4px] border border-border bg-muted-foreground px-3 py-1.5"
        >
          <Building2 class="w-3.5 h-3.5 text-primary" />
          <span>{{ agencesCount || 1 }} Agence(s) connectée(s)</span>
        </div>
        <div
          class="inline-flex items-center gap-1.5 rounded-[4px] border border-border bg-muted-foreground px-3 py-1.5 font-mono"
        >
          <Award class="w-3.5 h-3.5 text-primary" />
          <span>{{ schoolAgrement }}</span>
        </div>
      </div>

      <!-- Boutons d'Action Principaux (CTA) -->
      <div
        class="mt-8 flex flex-col sm:flex-row items-center justify-center gap-3"
      >
        <button
          type="button"
          @click="emit('navigate-auth')"
          class="w-full sm:w-auto inline-flex items-center justify-center gap-2 rounded-[4px] bg-brand px-6 py-3.5 text-sm font-bold text-white hover:bg-brand-dark hover:shadow-lg transition transform active:scale-98 shadow-sm"
        >
          <span>Se connecter à mon espace</span>
          <ArrowRight class="w-4 h-4" />
        </button>

        <button
          type="button"
          @click="emit('install-pwa')"
          class="w-full sm:w-auto inline-flex items-center justify-center gap-2 rounded-[4px] border border-border bg-card px-5 py-3.5 text-sm font-semibold text-foreground hover:bg-muted hover:border-brand/40 transition shadow-xs"
        >
          <Download class="w-4 h-4 text-brand" />
          <span>Installer l'application PWA</span>
        </button>
      </div>
    </div>
  </section>
</template>

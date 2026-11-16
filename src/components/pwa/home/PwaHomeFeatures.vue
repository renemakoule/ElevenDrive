<script setup lang="ts">
import { Smartphone, Calendar, BookOpen, Receipt } from 'lucide-vue-next';

const features = [
  {
    title: 'Livret PWA 100% Tactile',
    description: 'Fini le livret papier égaré : votre fiche officielle est synchronisée en temps réel et consultable hors-ligne.',
    icon: Smartphone,
    boxClass: 'bg-brand/10 text-brand',
  },
  {
    title: 'Créneaux de Conduite Flexibles',
    description: 'Consultez les disponibilités de votre moniteur et planifiez vos leçons de conduite pratique en quelques clics.',
    icon: Calendar,
    boxClass: 'bg-sky-500/10 text-sky-600 dark:text-sky-400',
  },
  {
    title: 'Annales & Quiz Officiels',
    description: 'Entraînez-vous sur les vrais sujets de l\'examen d\'État avec explications audio traduites dans votre langue.',
    icon: BookOpen,
    boxClass: 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400',
  },
  {
    title: 'Transparence des Reçus',
    description: 'Historique certifié de chaque acompte, émission instantanée de reçus thermiques et contrôle du reste à payer.',
    icon: Receipt,
    boxClass: 'bg-amber-500/10 text-amber-600 dark:text-amber-400',
  },
];
</script>

<template>
  <section class="shrink-0 w-full py-10 bg-muted/30 border-y border-border">
    <div class="max-w-6xl mx-auto px-4 sm:px-6">
      <div class="text-center mb-8">
        <h3 class="text-lg sm:text-xl font-bold text-foreground">
          L'excellence pédagogique au service de votre réussite
        </h3>
        <p class="text-xs sm:text-sm text-muted-foreground mt-1">
          Tous les outils nécessaires pour obtenir votre permis du premier coup.
        </p>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div
          v-for="(item, idx) in features"
          :key="idx"
          class="rounded-[4px] border border-border bg-card p-4 space-y-2 hover:border-brand/40 transition shadow-2xs"
        >
          <div
            class="w-9 h-9 rounded-lg flex items-center justify-center font-bold"
            :class="item.boxClass"
          >
            <component :is="item.icon" class="w-5 h-5" />
          </div>
          <h5 class="text-sm font-semibold text-foreground">{{ item.title }}</h5>
          <p class="text-xs text-muted-foreground leading-relaxed">
            {{ item.description }}
          </p>
        </div>
      </div>
    </div>
  </section>
</template>

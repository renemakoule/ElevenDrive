<script setup lang="ts">
import LandingNavbar from "./LandingNavbar.vue";
import LandingHero from "./LandingHero.vue";
import ServiceGuichet from "./ServiceGuichet.vue";
import ServicePedagogie from "./ServicePedagogie.vue";
import ServiceCloudArchivage from "./ServiceCloudArchivage.vue";
import ServiceMultiAgences from "./ServiceMultiAgences.vue";
import LandingAbout from "./LandingAbout.vue";
import LandingPricing from "./LandingPricing.vue";
import LandingFooter from "./LandingFooter.vue";

const emit = defineEmits<{
  (e: "navigate-auth"): void;
}>();

const handleScrollTo = (id: string) => {
  if (id === "top") {
    window.scrollTo({ top: 0, behavior: "smooth" });
    return;
  }
  const target = document.getElementById(id);
  if (target) {
    target.scrollIntoView({ behavior: "smooth", block: "start" });
  }
};
</script>

<template>
  <div
    class="fixed inset-0 z-50 overflow-y-auto bg-background text-foreground antialiased selection:bg-brand-500 selection:text-white"
  >
    <!-- 1. En-tête de navigation -->
    <LandingNavbar
      @navigate-auth="emit('navigate-auth')"
      @scroll-to="handleScrollTo"
    />

    <!-- 2. Section Hero principale -->
    <main>
      <LandingHero
        @navigate-auth="emit('navigate-auth')"
        @scroll-to="handleScrollTo"
      />

      <!-- 3. Sections modulaires de présentation des services -->
      <div id="solutions" class="scroll-mt-14">
        <!-- Service 1 : Guichet & Caisse -->
        <ServiceGuichet @navigate-auth="emit('navigate-auth')" />

        <!-- Service 2 : Pédagogie & Conduite -->
        <ServicePedagogie @navigate-auth="emit('navigate-auth')" />

        <!-- Service 3 : Local-First & Archivage Cloud R2 -->
        <ServiceCloudArchivage @navigate-auth="emit('navigate-auth')" />

        <!-- Service 4 : Direction & Multi-Agences -->
        <ServiceMultiAgences @navigate-auth="emit('navigate-auth')" />
      </div>

      <!-- 4. Section À propos / Vision & Piliers -->
      <div id="about" class="scroll-mt-14">
        <LandingAbout @navigate-auth="emit('navigate-auth')" />
      </div>

      <!-- 5. Section Tarifs clairs & transparents -->
      <div id="pricing" class="scroll-mt-14">
        <LandingPricing @navigate-auth="emit('navigate-auth')" />
      </div>
    </main>

    <!-- 6. Pied de page épuré -->
    <LandingFooter
      @navigate-auth="emit('navigate-auth')"
      @scroll-to="handleScrollTo"
    />
  </div>
</template>

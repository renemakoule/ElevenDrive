<script setup lang="ts">
import appLogo from "@root/app-icon.png";
</script>

<template>
  <aside
    class="relative hidden flex-col justify-between overflow-hidden border-r border-border bg-muted/30 p-10 lg:flex"
  >
    <!-- Trame de fond géométrique estompée -->
    <div
      class="pointer-events-none absolute inset-0 text-foreground/10"
      style="
        background-image: radial-gradient(currentColor 1px, transparent 1px);
        background-size: 16px 16px;
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
    ></div>
    <div
      class="pointer-events-none absolute -left-24 -top-24 h-80 w-80 rounded-full bg-brand/15 blur-3xl"
    ></div>
    <div
      class="pointer-events-none absolute -bottom-32 -right-24 h-96 w-96 rounded-full bg-brand/10 blur-3xl"
    ></div>

    <!-- En-tête de marque -->
    <div class="relative flex items-center gap-2.5">
      <img
        :src="appLogo"
        alt="ElevenDrive Logo"
        class="h-8 w-8 object-contain"
      />
      <span class="text-base font-semibold tracking-tight text-foreground">
        ElevenDrive<span class="text-brand">.</span>
      </span>
    </div>

    <!-- Proposition de valeur centrale -->
    <div class="relative mx-auto max-w-sm text-center">
      <div
        class="mx-auto mb-8 flex h-28 w-28 items-center justify-center border border-border/80 bg-card/60 p-3 backdrop-blur-xs"
      >
        <img
          :src="appLogo"
          alt="AutoSchool Logo"
          class="h-full w-full object-contain filter drop-shadow-sm"
        />
      </div>
      <h2
        class="text-3xl font-semibold leading-tight tracking-tight text-foreground"
      >
        Pilotez votre auto-école, du guichet à la route.
      </h2>
      <p class="mt-3 text-sm leading-relaxed text-muted-foreground">
        Plateforme intégrée pour les auto-écoles en Afrique.
      </p>
    </div>

    <!-- Pied de page -->
    <footer class="px-6 py-4 text-center text-[11px] text-muted-foreground">
      Powered by GAS.
    </footer>
  </aside>
</template>

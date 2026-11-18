<script setup lang="ts">
import Skeleton from "./Skeleton.vue";
</script>

<template>
  <div class="space-y-4">
    <!-- Barre Calendrier -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 rounded-[3px] border border-border bg-card p-3">
      <div class="flex items-center gap-2">
        <Skeleton class="h-8 w-8 rounded" />
        <Skeleton class="h-5 w-40" />
        <Skeleton class="h-8 w-8 rounded" />
      </div>
      <div class="flex items-center gap-1.5 overflow-x-auto">
        <Skeleton v-for="d in 6" :key="d" class="h-10 w-12 rounded-[4px] shrink-0" />
      </div>
    </div>

    <!-- Contenu Principal : Grille Planning + Sidebar -->
    <div class="grid grid-cols-1 items-start gap-4 lg:grid-cols-12">
      <!-- Planning des séances -->
      <div class="lg:col-span-8 space-y-3">
        <div class="rounded-[3px] border border-border bg-card p-4 space-y-3">
          <div class="flex items-center justify-between border-b border-border pb-3">
            <Skeleton class="h-4 w-36" />
            <Skeleton class="h-4 w-20" />
          </div>
          <div class="space-y-3">
            <div
              v-for="s in 4"
              :key="s"
              class="rounded-[4px] border border-border p-3.5 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 bg-muted/20"
            >
              <div class="flex items-center gap-3">
                <Skeleton class="h-10 w-14 rounded font-mono" />
                <div class="space-y-1.5">
                  <div class="flex items-center gap-2">
                    <Skeleton class="h-4 w-32" />
                    <Skeleton class="h-4 w-16 rounded-full" />
                  </div>
                  <Skeleton class="h-3 w-48" />
                </div>
              </div>
              <div class="flex items-center gap-2">
                <Skeleton class="h-8 w-24 rounded-[3px]" />
                <Skeleton class="h-8 w-8 rounded-[3px]" />
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Résumé Latéral Moniteur / Véhicule -->
      <div class="lg:col-span-4 space-y-3">
        <div class="rounded-[3px] border border-border bg-card p-4 space-y-4">
          <div class="flex items-center gap-3">
            <Skeleton class="h-12 w-12 rounded-full shrink-0" />
            <div class="space-y-1.5 flex-1">
              <Skeleton class="h-4 w-32" />
              <Skeleton class="h-3 w-40" />
            </div>
          </div>
          <div class="border-t border-border pt-3 space-y-2">
            <Skeleton class="h-3.5 w-28" />
            <Skeleton class="h-9 w-full rounded-[4px]" />
            <div class="grid grid-cols-2 gap-2 pt-2">
              <Skeleton class="h-12 rounded-[3px]" />
              <Skeleton class="h-12 rounded-[3px]" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ProvisionStep } from "./provisioningSteps";
import { CheckCircle2, Layers, FileCode } from "lucide-vue-next";

defineProps<{
  steps: ProvisionStep[];
  progressPercent: number;
  consoleLogs: string[];
}>();
</script>

<template>
  <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
    <!-- Étapes séquentielles -->
    <div
      class="p-4 rounded-xl bg-white/[0.02] border border-white/10 space-y-2 max-h-60 overflow-y-auto"
    >
      <div
        class="text-xs font-bold text-white/80 uppercase font-mono tracking-wider flex items-center gap-2 mb-3"
      >
        <Layers class="w-3.5 h-3.5 text-brand" />
        <span>Modules & Structures déployés</span>
      </div>

      <div
        v-for="step in steps"
        :key="step.id"
        class="flex items-start gap-2.5 p-2 rounded-lg transition"
        :class="
          progressPercent >= step.targetPercent
            ? 'bg-emerald-500/10 border border-emerald-500/20'
            : 'opacity-40'
        "
      >
        <div class="mt-0.5">
          <CheckCircle2
            v-if="progressPercent >= step.targetPercent"
            class="w-4 h-4 text-emerald-400"
          />
          <div
            v-else
            class="w-4 h-4 rounded-full border border-white/30"
          ></div>
        </div>
        <div class="flex-1 text-xs">
          <div
            class="font-bold"
            :class="
              progressPercent >= step.targetPercent
                ? 'text-white'
                : 'text-white/60'
            "
          >
            {{ step.title }}
          </div>
          <div class="text-[11px] text-white/50">
            {{ step.detail }}
          </div>
        </div>
      </div>
    </div>

    <!-- Terminal de streaming -->
    <div
      class="p-4 rounded-xl bg-black/80 border border-white/10 font-mono text-xs flex flex-col justify-between max-h-60 overflow-hidden shadow-inner"
    >
      <div
        class="flex items-center justify-between border-b border-white/10 pb-2 mb-2 text-white/40 text-[10px]"
      >
        <span class="flex items-center gap-1.5">
          <FileCode class="w-3 h-3 text-cyan-400" />
          <span>TERMINAL DE DÉPLOIEMENT</span>
        </span>
        <span class="text-emerald-400 animate-pulse">● LIVE STREAM</span>
      </div>

      <div
        class="space-y-1.5 overflow-y-auto pr-1 text-[11px] leading-relaxed select-text"
      >
        <div
          v-for="(log, idx) in consoleLogs"
          :key="idx"
          :class="
            log.includes('SUCCESS')
              ? 'text-emerald-400 font-bold'
              : log.includes('LICENCE')
                ? 'text-amber-300 font-bold'
                : 'text-white/70'
          "
        >
          {{ log }}
        </div>
      </div>

      <div
        class="pt-2 text-[10px] text-white/40 flex items-center justify-between border-t border-white/10 mt-2"
      >
        <span>Réseau : Neon Serverless Multi-Region</span>
        <span>Cryptage : AES-256 GCM</span>
      </div>
    </div>
  </div>
</template>


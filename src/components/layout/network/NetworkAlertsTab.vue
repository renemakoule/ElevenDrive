<script setup lang="ts">
import type { NetworkAlert } from "../../../types";
import { CheckCircle2, RotateCcw, X, Clock, EyeOff } from "lucide-vue-next";

defineProps<{
  alerts: NetworkAlert[];
  dismissedCount: number;
}>();

const emit = defineEmits<{
  (e: "dismiss", id: string): void;
  (e: "restore"): void;
}>();

const formatDate = (iso: string) => {
  try {
    const d = new Date(iso);
    return `${d.toLocaleDateString("fr-FR")} à ${d.toLocaleTimeString("fr-FR", { hour: "2-digit", minute: "2-digit" })}`;
  } catch {
    return iso;
  }
};
</script>

<template>
  <div>
    <!-- État nominal / Aucune alerte active -->
    <div
      v-if="alerts.length === 0"
      class="py-14 text-center text-muted-foreground space-y-2"
    >
      <CheckCircle2 class="w-10 h-10 text-emerald-500 mx-auto opacity-80" />
      <div class="font-bold text-foreground text-sm">Système Nominal</div>
      <p class="text-[11px]">
        Aucune alerte technique ou opérationnelle n'est active actuellement.
      </p>
      <button
        v-if="dismissedCount > 0"
        type="button"
        @click="emit('restore')"
        class="mt-2 inline-flex items-center gap-1.5 text-[11px] text-brand hover:underline font-semibold cursor-pointer"
      >
        <RotateCcw class="w-3.5 h-3.5" />
        <span>Réafficher les alertes masquées</span>
      </button>
    </div>

    <!-- Liste des alertes réseau -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-3.5 items-stretch">
      <div
        v-for="alt in alerts"
        :key="alt.id"
        class="relative p-4 rounded-[4px] border space-y-2 group transition-all"
        :class="[
          alt.severite === 'critical'
            ? 'border-rose-500/50 bg-rose-500/5'
            : alt.severite === 'warning'
              ? 'border-amber-500/50 bg-amber-500/5'
              : 'border-blue-500/50 bg-blue-500/5',
        ]"
      >
        <!-- Bouton fermer directement l'alerte -->
        <button
          type="button"
          @click="emit('dismiss', alt.id)"
          class="absolute top-3 right-3 p-1 rounded hover:bg-muted text-muted-foreground hover:text-foreground transition cursor-pointer"
          title="Masquer cette alerte"
        >
          <X class="w-3.5 h-3.5" />
        </button>

        <div class="flex items-center gap-2 pr-6">
          <span
            class="text-[9px] uppercase font-bold px-2 py-0.5 rounded-[3px]"
            :class="
              alt.severite === 'critical'
                ? 'bg-rose-600 text-white'
                : alt.severite === 'warning'
                  ? 'bg-amber-600 text-white'
                  : 'bg-blue-600 text-white'
            "
          >
            {{
              alt.severite === "critical"
                ? "Alerte Critique"
                : alt.severite === "warning"
                  ? "Avertissement"
                  : "Information Réseau"
            }}
          </span>

          <span
            class="text-[10px] text-muted-foreground flex items-center gap-1 font-mono"
          >
            <Clock class="w-3 h-3" />
            <span>{{ formatDate(alt.date_envoi) }}</span>
          </span>
        </div>

        <div class="font-bold text-foreground text-sm">{{ alt.titre }}</div>
        <p
          class="text-xs text-muted-foreground leading-relaxed whitespace-pre-line"
        >
          {{ alt.message }}
        </p>

        <div class="pt-1 flex items-center justify-end">
          <button
            type="button"
            @click="emit('dismiss', alt.id)"
            class="text-[11px] text-muted-foreground hover:text-foreground font-medium flex items-center gap-1 cursor-pointer transition"
          >
            <EyeOff class="w-3 h-3" />
            <span>Masquer cette alerte</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

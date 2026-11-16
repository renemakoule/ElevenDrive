<script setup lang="ts">
import type { FormulePermis } from "../../../types";

defineProps<{
  formules: FormulePermis[];
  formuleId: string;
  modeFormationTheorique: "PRESENTIEL" | "EN_LIGNE" | "HYBRIDE";
}>();

const emit = defineEmits<{
  (e: "select-formule", id: string): void;
  (
    e: "update:modeFormationTheorique",
    val: "PRESENTIEL" | "EN_LIGNE" | "HYBRIDE",
  ): void;
}>();

const fcfa = (n: number) => n.toLocaleString("fr-FR");

const modesFormationOptions = [
  {
    value: "PRESENTIEL",
    label: "100% En Salle",
    desc: "Cours réguliers à l'agence",
    icon: "🏫",
    badge: "Présentiel",
  },
  {
    value: "EN_LIGNE",
    label: "100% En Ligne",
    desc: "Visioconférence & supports à distance",
    icon: "💻",
    badge: "Virtuel",
  },
  {
    value: "HYBRIDE",
    label: "Hybride (Mixte)",
    desc: "Salle & Visio au choix de l'élève",
    icon: "🔄",
    badge: "Recommandé",
  },
] as const;
</script>

<template>
  <div>
    <!-- Formule -->
    <section class="px-4 py-4">
      <h3 class="mb-2 text-xs font-semibold text-foreground">Formule</h3>
      <div
        v-if="formules.length > 0"
        class="grid grid-cols-1 gap-2 sm:grid-cols-3"
        role="radiogroup"
        aria-label="Formule de permis"
      >
        <button
          v-for="f in formules"
          :key="f.id"
          type="button"
          role="radio"
          :aria-checked="formuleId === f.id"
          @click="emit('select-formule', f.id)"
          class="flex flex-col justify-between rounded-[4px] border p-3 text-left transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40"
          :class="
            formuleId === f.id
              ? 'border-brand bg-brand-light'
              : 'border-border bg-card hover:border-border hover:bg-muted'
          "
        >
          <span>
            <span class="block text-xs font-semibold text-foreground">{{
              f.intitule
            }}</span>
            <span class="mt-0.5 block text-[11px] text-muted-foreground">
              Code {{ f.heures_theoriques }} h · Conduite
              {{ f.heures_pratiques }} h
            </span>
          </span>
          <span class="mt-2.5 block border-t border-border pt-2">
            <span
              class="block text-sm font-semibold tabular-nums text-foreground"
              >{{ fcfa(f.prix_total) }} FCFA</span
            >
            <span class="block text-[11px] text-muted-foreground"
              >Acompte min. {{ fcfa(f.acompte_minimum) }} FCFA</span
            >
          </span>
        </button>
      </div>

      <div
        v-else
        class="rounded-[4px] border border-dashed border-amber-300 dark:border-amber-800 bg-amber-50/50 dark:bg-amber-950/20 p-4 text-center text-xs space-y-1.5"
      >
        <p class="font-semibold text-amber-800 dark:text-amber-300">
          Aucune formule tarifaire n'est encore configurée.
        </p>
        <p class="text-[11px] text-muted-foreground">
          Rendez-vous dans la section
          <strong>"Tarifs & Permis"</strong> pour créer vos formules officielles
          avant d'enregistrer des candidats.
        </p>
      </div>
    </section>

    <!-- Modalité de formation théorique -->
    <section class="space-y-2.5 border-t border-border px-4 py-3.5 bg-muted/20">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-1.5">
          <span class="text-xs font-semibold text-foreground"
            >Modalité des cours théoriques (Code)</span
          >
          <span
            class="text-[10px] text-white bg-brand-light px-1.5 py-0.2 rounded font-semibold border border-brand-border"
            >Nouveau</span
          >
        </div>
        <span class="text-[11px] text-muted-foreground"
          >Option choisie par l'élève</span
        >
      </div>

      <div
        class="grid grid-cols-1 gap-2 sm:grid-cols-3"
        role="radiogroup"
        aria-label="Modalité de formation théorique"
      >
        <button
          v-for="m in modesFormationOptions"
          :key="m.value"
          type="button"
          role="radio"
          :aria-checked="modeFormationTheorique === m.value"
          @click="emit('update:modeFormationTheorique', m.value)"
          class="flex flex-col justify-between rounded-[4px] border p-2.5 text-left transition cursor-pointer"
          :class="
            modeFormationTheorique === m.value
              ? 'border-brand bg-brand-light text-white'
              : 'border-border bg-card text-muted-foreground hover:bg-muted'
          "
        >
          <div>
            <div class="flex items-center justify-between">
              <span class="text-base">{{ m.icon }}</span>
              <span
                class="text-[9px] font-bold uppercase tracking-wider px-1.5 py-0.2 rounded"
                :class="
                  modeFormationTheorique === m.value
                    ? 'bg-brand text-white'
                    : 'bg-muted text-muted-foreground'
                "
              >
                {{ m.badge }}
              </span>
            </div>
            <span class="mt-1 block text-xs font-bold text-foreground">{{
              m.label
            }}</span>
          </div>
          <span
            class="mt-1 block text-[10px] text-muted-foreground leading-tight"
            >{{ m.desc }}</span
          >
        </button>
      </div>
    </section>
  </div>
</template>

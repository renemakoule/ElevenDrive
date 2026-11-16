<script setup lang="ts">
import { computed } from "vue";
import {
  LayoutDashboard,
  CreditCard,
  Users,
  GraduationCap,
  Car,
  FolderArchive,
} from "lucide-vue-next";
import { useLanguageStore } from "../../stores/languageStore";

const languageStore = useLanguageStore();
const t = (text: string) => languageStore.t(text);

// Aperçu illustratif de l'interface (données fictives de démonstration)
const nav = [
  { icon: LayoutDashboard, active: true },
  { icon: CreditCard },
  { icon: Users },
  { icon: GraduationCap },
  { icon: Car },
  { icon: FolderArchive },
];

const stats = computed(() => [
  { label: t("Encaissements du jour"), value: "485 000 FCFA", trend: "+12%" },
  { label: t("Candidats actifs"), value: "128", trend: "+4" },
  { label: t("Séances aujourd'hui"), value: "24", trend: "" },
  { label: t("Taux de réussite"), value: "87%", trend: "+3%" },
]);

const bars = [38, 52, 44, 63, 58, 72, 66, 80, 74, 88, 82, 94];

const rows = computed(() => [
  {
    who: t("Candidat A."),
    what: t("Permis B · Acompte"),
    amount: "50 000",
    mode: t("Espèces"),
  },
  {
    who: t("Candidat M."),
    what: t("Permis B · Solde"),
    amount: "120 000",
    mode: t("Mobile"),
  },
  {
    who: t("Candidat S."),
    what: t("Permis A · Inscription"),
    amount: "35 000",
    mode: t("Espèces"),
  },
  {
    who: t("Candidat K."),
    what: t("Code · Séance"),
    amount: "10 000",
    mode: t("Espèces"),
  },
]);
</script>

<template>
  <div
    class="overflow-hidden rounded-[8px] border border-border bg-card"
    aria-hidden="true"
  >
    <!-- Barre de fenêtre -->
    <div
      class="flex items-center gap-1.5 border-b border-border bg-muted/50 px-3.5 py-2.5"
    >
      <span class="h-2 w-2 rounded-full bg-border" />
      <span class="h-2 w-2 rounded-full bg-border" />
      <span class="h-2 w-2 rounded-full bg-border" />
      <span
        class="ml-3 rounded-[3px] bg-background px-2.5 py-0.5 text-[10px] text-muted-foreground"
      >
        {{ t("Tableau de bord · Direction") }}
      </span>
    </div>

    <div class="flex">
      <!-- Sidebar -->
      <div
        class="hidden w-12 shrink-0 flex-col items-center gap-1 border-r border-border bg-muted/30 py-3 sm:flex"
      >
        <span
          v-for="(item, i) in nav"
          :key="i"
          class="flex h-8 w-8 items-center justify-center rounded-[4px]"
          :class="
            item.active ? 'bg-brand/10 text-brand-600' : 'text-muted-foreground'
          "
        >
          <component :is="item.icon" class="h-4 w-4" />
        </span>
      </div>

      <!-- Contenu -->
      <div class="min-w-0 flex-1 space-y-3 p-3 sm:p-4">
        <div class="grid grid-cols-2 gap-2.5 lg:grid-cols-4">
          <div
            v-for="s in stats"
            :key="s.label"
            class="rounded-[6px] border border-border bg-background p-3"
          >
            <p class="truncate text-[10px] text-muted-foreground">
              {{ s.label }}
            </p>
            <div class="mt-1 flex items-baseline gap-1.5">
              <span
                class="text-[15px] font-semibold tabular-nums text-foreground"
              >
                {{ s.value }}
              </span>
              <span
                v-if="s.trend"
                class="text-[10px] font-medium text-brand-600"
              >
                {{ s.trend }}
              </span>
            </div>
          </div>
        </div>

        <div class="grid gap-2.5 lg:grid-cols-5">
          <!-- Graphique -->
          <div
            class="rounded-[6px] border border-border bg-background p-3 lg:col-span-3"
          >
            <p class="text-[10px] font-medium text-foreground">
              {{ t("Recettes · 12 derniers jours") }}
            </p>
            <div class="mt-3 flex h-24 items-end gap-1.5">
              <span
                v-for="(h, i) in bars"
                :key="i"
                class="flex-1 rounded-t-[2px] bg-gradient-to-t from-brand-500/60 to-brand-500"
                :style="{
                  height: h + '%',
                  opacity: 0.55 + (i / bars.length) * 0.45,
                }"
              />
            </div>
          </div>

          <!-- Paiements récents -->
          <div
            class="rounded-[6px] border border-border bg-background p-3 lg:col-span-2"
          >
            <p class="text-[10px] font-medium text-foreground">
              {{ t("Paiements récents") }}
            </p>
            <ul class="mt-2 divide-y divide-border">
              <li
                v-for="r in rows"
                :key="r.who + r.what"
                class="flex items-center justify-between gap-2 py-1.5"
              >
                <div class="min-w-0">
                  <p class="truncate text-[10px] font-medium text-foreground">
                    {{ r.who }}
                  </p>
                  <p class="truncate text-[9px] text-muted-foreground">
                    {{ r.what }}
                  </p>
                </div>
                <div class="shrink-0 text-right">
                  <p
                    class="text-[10px] font-semibold tabular-nums text-foreground"
                  >
                    {{ r.amount }}
                  </p>
                  <p class="text-[9px] text-muted-foreground">{{ r.mode }}</p>
                </div>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

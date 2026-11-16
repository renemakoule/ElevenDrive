<script setup lang="ts">
import { ref, computed } from "vue";
import { useLanguageStore } from "../../stores/languageStore";
import { useAuthStore } from "../../stores/authStore";
import {
  Keyboard,
  X,
  Search,
  CheckCircle2,
  Sliders,
  Sparkles,
} from "lucide-vue-next";

defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
}>();

const languageStore = useLanguageStore();
const authStore = useAuthStore();
const t = (text: string) => languageStore.t(text);

const searchQuery = ref("");
const selectedRole = ref<"all" | "secretaire" | "manager" | "global">("all");

interface ShortcutItem {
  key: string;
  action: string;
  category: "guichet" | "navigation" | "systeme";
  roles: string[];
  description: string;
}

const shortcuts: ShortcutItem[] = [
  {
    key: "F1",
    action: "Nouvelle Inscription Express",
    category: "guichet",
    roles: ["secretaire", "manager"],
    description:
      "Ouvre le formulaire d'enrôlement rapide d'un nouveau candidat.",
  },
  {
    key: "F2",
    action: "Encaissement Guichet Direct",
    category: "guichet",
    roles: ["secretaire", "manager"],
    description:
      "Bascule sur le guichet et positionne le curseur sur le montant.",
  },
  {
    key: "F3",
    action: "Rechercher un Candidat",
    category: "guichet",
    roles: ["secretaire", "manager"],
    description:
      "Focus immédiat sur la barre de recherche globale des dossiers.",
  },
  {
    key: "F4",
    action: "Clôture de Caisse (Rapport Z)",
    category: "guichet",
    roles: ["secretaire", "manager"],
    description:
      "Ouvre la modale de réconciliation d'espèces et impression ticket Z.",
  },
  {
    key: "F5",
    action: "Feuille de Présences Code",
    category: "navigation",
    roles: ["secretaire", "manager"],
    description:
      "Bascule rapidement vers le module d'émargement en salle de code.",
  },
  {
    key: "F9",
    action: "Examens d'État & Bordereaux",
    category: "navigation",
    roles: ["secretaire", "manager"],
    description:
      "Accès direct à la gestion des sessions d'examens et listes DGTTC.",
  },
  {
    key: "Échap",
    action: "Fermer Modale / Retour Accueil",
    category: "systeme",
    roles: ["secretaire", "manager", "moniteur", "eleve"],
    description: "Ferme toute modale active ou revient à la vue d'ensemble.",
  },
  {
    key: "Ctrl + P",
    action: "Imprimer Ticket ou Document",
    category: "systeme",
    roles: ["secretaire", "manager"],
    description: "Déclenche l'impression ESC/POS ou A4 du document affiché.",
  },
];

const filteredShortcuts = computed(() => {
  return shortcuts.filter((s) => {
    const matchesSearch =
      s.action.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      s.key.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      s.description.toLowerCase().includes(searchQuery.value.toLowerCase());

    const matchesRole =
      selectedRole.value === "all" ||
      (selectedRole.value === "secretaire" && s.roles.includes("secretaire")) ||
      (selectedRole.value === "manager" && s.roles.includes("manager")) ||
      (selectedRole.value === "global" && s.roles.includes("eleve"));

    return matchesSearch && matchesRole;
  });
});
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-xs p-4 select-none animate-in fade-in duration-150"
    @click.self="emit('close')"
  >
    <div
      class="bg-card border border-border rounded-[3px] max-w-xl w-full overflow-hidden flex flex-col max-h-[90vh]"
    >
      <!-- En-tête -->
      <div
        class="p-4 bg-muted/40 border-b border-border flex items-center justify-between shrink-0"
      >
        <div class="flex items-center space-x-2.5">
          <div
            class="w-8 h-8 rounded-[4px] bg-brand-light text-brand-dark flex items-center justify-center font-bold"
          >
            <Keyboard class="w-4 h-4 text-brand" />
          </div>
          <div>
            <div
              class="font-bold text-foreground text-sm flex items-center gap-2"
            >
              <span>{{ t("Raccourcis Clavier Rapides") }}</span>
              <span
                class="text-[10px] px-1.5 py-0.2 rounded bg-brand-light text-brand-dark font-medium"
              >
                {{ t("Ergonomie Guichet") }}
              </span>
            </div>
            <div class="text-[11px] text-muted-foreground">
              {{ t("Gagnez en rapidité lors de vos opérations au quotidien") }}
            </div>
          </div>
        </div>
        <button
          type="button"
          @click="emit('close')"
          class="p-1 rounded-[3px] text-muted-foreground hover:text-foreground hover:bg-muted transition"
          :title="t('Fermer')"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Filtres et Recherche -->
      <div class="p-4 border-b border-border/50 bg-muted/20 space-y-3 shrink-0">
        <div class="relative">
          <Search
            class="w-3.5 h-3.5 text-muted-foreground absolute left-3 top-1/2 -translate-y-1/2"
          />
          <input
            v-model="searchQuery"
            type="text"
            :placeholder="
              t('Rechercher un raccourci (ex: F1, Caisse, Inscription)...')
            "
            class="w-full pl-9 pr-3 py-1.5 text-xs rounded-[4px] border border-border bg-background focus:outline-hidden focus:ring-1 focus:ring-brand"
          />
        </div>

        <div class="flex items-center gap-1.5 text-xs">
          <button
            type="button"
            @click="selectedRole = 'all'"
            class="px-2.5 py-1 rounded-[4px] text-[11px] font-medium transition"
            :class="
              selectedRole === 'all'
                ? 'bg-brand text-white font-semibold'
                : 'bg-muted text-muted-foreground hover:text-foreground'
            "
          >
            {{ t("Tous les raccourcis") }}
          </button>
          <button
            type="button"
            @click="selectedRole = 'secretaire'"
            class="px-2.5 py-1 rounded-[4px] text-[11px] font-medium transition"
            :class="
              selectedRole === 'secretaire'
                ? 'bg-brand text-white font-semibold'
                : 'bg-muted text-muted-foreground hover:text-foreground'
            "
          >
            {{ t("Secrétaire (Guichet)") }}
          </button>
          <button
            type="button"
            @click="selectedRole = 'manager'"
            class="px-2.5 py-1 rounded-[4px] text-[11px] font-medium transition"
            :class="
              selectedRole === 'manager'
                ? 'bg-brand text-white font-semibold'
                : 'bg-muted text-muted-foreground hover:text-foreground'
            "
          >
            {{ t("Directeur") }}
          </button>
        </div>
      </div>

      <!-- Liste des Raccourcis -->
      <div class="p-4 overflow-y-auto space-y-2 flex-1 min-h-0">
        <div
          v-if="filteredShortcuts.length === 0"
          class="text-center py-8 text-muted-foreground text-xs"
        >
          {{ t("Aucun raccourci ne correspond à votre recherche.") }}
        </div>

        <div
          v-for="s in filteredShortcuts"
          :key="s.key"
          class="flex items-center justify-between p-2.5 rounded-[4px] border border-border/70 bg-card hover:bg-muted/40 transition group"
        >
          <div class="space-y-0.5 pr-3">
            <div
              class="text-xs font-semibold text-foreground flex items-center gap-2"
            >
              <span>{{ t(s.action) }}</span>
            </div>
            <div class="text-[11px] text-muted-foreground">
              {{ t(s.description) }}
            </div>
          </div>
          <div class="shrink-0">
            <kbd
              class="px-2.5 py-1 text-xs font-mono font-bold bg-muted text-foreground border border-border rounded"
            >
              {{ s.key }}
            </kbd>
          </div>
        </div>
      </div>

      <!-- Pied de page -->
      <div
        class="p-3 bg-muted/40 border-t border-border flex items-center justify-between text-xs text-muted-foreground shrink-0"
      >
        <span class="flex items-center gap-1.5 text-[11px]">
          <Sparkles class="w-3.5 h-3.5 text-amber-500" />
          <span>{{
            t(
              "Astuce : Les touches F1 à F9 sont actives sur le clavier physique.",
            )
          }}</span>
        </span>
        <button
          type="button"
          @click="emit('close')"
          class="px-3 py-1.5 rounded-[4px] bg-muted hover:bg-muted/80 text-foreground font-medium text-xs transition"
        >
          {{ t("Fermer") }}
        </button>
      </div>
    </div>
  </div>
</template>

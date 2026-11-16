<script setup lang="ts">
import { computed, ref } from "vue";
import type { Candidat } from "../../../types";
import { Check, Search, Users, X } from "lucide-vue-next";

const props = defineProps<{
  allCandidats: Candidat[];
  selectedCandidats: Candidat[];
  dossiersMap: Record<string, any>;
}>();

const emit = defineEmits<{
  (e: "toggle-candidat", c: Candidat): void;
  (e: "remove-candidat", id: string): void;
  (e: "clear-all"): void;
}>();

const CATEGORIES_PERMIS = [
  { code: "ALL", label: "Toutes", icon: "🌐" },
  { code: "B", label: "Permis B", icon: "🚗" },
  { code: "A", label: "Permis A", icon: "🏍️" },
  { code: "C", label: "Permis C", icon: "🚚" },
  { code: "D", label: "Permis D", icon: "🚌" },
  { code: "E", label: "Permis E", icon: "🚜" },
  { code: "REC", label: "Recyclage", icon: "🔄" },
];

const selectedCategorie = ref("ALL");
const searchQuery = ref("");

const extraireCategorie = (nom?: string) => {
  if (!nom) return "B";
  const upper = nom.toUpperCase();
  if (upper.includes("RECYCLAGE") || upper.includes("PERFECTIONNEMENT"))
    return "REC";
  if (
    upper.includes("PERMIS A") ||
    upper.includes("MOTO") ||
    upper.includes("A1") ||
    upper.includes("A2")
  )
    return "A";
  if (
    upper.includes("PERMIS C") ||
    upper.includes("POIDS LOURD") ||
    upper.includes("CAMION")
  )
    return "C";
  if (
    upper.includes("PERMIS D") ||
    upper.includes("TRANSPORT") ||
    upper.includes("BUS")
  )
    return "D";
  if (upper.includes("PERMIS E") || upper.includes("REMORQUE")) return "E";
  return "B";
};

const isSelected = (id: string) =>
  props.selectedCandidats.some((c) => c.id === id);

const countsParCategorie = computed(() => {
  const counts: Record<string, number> = {
    ALL: props.allCandidats.length,
    B: 0,
    A: 0,
    C: 0,
    D: 0,
    E: 0,
    REC: 0,
  };
  for (const c of props.allCandidats) {
    const rawCat = (
      c.categorie_permis || extraireCategorie(c.formule_nom)
    ).toUpperCase();
    if (rawCat === "A1" || rawCat === "A2") counts.A++;
    else if (rawCat === "RECYCLAGE") counts.REC++;
    else if (counts[rawCat] !== undefined) counts[rawCat]++;
  }
  return counts;
});

const candidatsAffiches = computed(() => {
  let list = props.allCandidats;

  if (selectedCategorie.value !== "ALL") {
    list = list.filter((c) => {
      const cat = (
        c.categorie_permis || extraireCategorie(c.formule_nom)
      ).toUpperCase();
      if (selectedCategorie.value === "A")
        return cat === "A" || cat === "A1" || cat === "A2";
      if (selectedCategorie.value === "REC")
        return cat === "REC" || cat === "RECYCLAGE";
      return cat === selectedCategorie.value;
    });
  }

  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase().trim();
    list = list.filter((c) => {
      const nom = `${c.nom || ""} ${c.prenoms || ""}`.toLowerCase();
      const tel = (c.telephone_whatsapp || c.telephone || "").toLowerCase();
      const dos = (c.numero_dossier || "").toLowerCase();
      return nom.includes(q) || tel.includes(q) || dos.includes(q);
    });
  }

  return list.slice(0, 25);
});
</script>

<template>
  <div class="space-y-2">
    <div class="flex items-center justify-between">
      <label
        class="flex items-center gap-1.5 text-xs font-bold text-foreground"
      >
        <Users class="h-3.5 w-3.5 text-brand" />
        <span>1. Sélectionner un ou plusieurs élèves</span>
        <span class="text-red-500">*</span>
      </label>
      <div class="flex items-center gap-2 text-[10px]">
        <span class="font-semibold text-brand">
          {{ selectedCandidats.length }} sélectionné(s)
        </span>
        <button
          v-if="selectedCandidats.length > 0"
          type="button"
          @click="emit('clear-all')"
          class="text-muted-foreground hover:text-rose-600 underline cursor-pointer"
        >
          Tout effacer
        </button>
      </div>
    </div>

    <!-- Chips des élèves sélectionnés -->
    <div
      v-if="selectedCandidats.length > 0"
      class="flex flex-wrap gap-1.5 rounded-[4px] border border-brand/30 bg-brand-light/40 p-2"
    >
      <span
        v-for="sel in selectedCandidats"
        :key="sel.id"
        class="inline-flex items-center gap-1.5 rounded-[3px] border border-brand/30 bg-card px-2 py-1 text-[11px] font-medium text-foreground"
      >
        <span>{{ sel.nom }} {{ sel.prenoms }}</span>
        <span
          v-if="dossiersMap[sel.id]?.formule"
          class="text-[10px] text-muted-foreground"
        >
          ({{
            dossiersMap[sel.id].dossier?.heures_conduite_effectuees || 0
          }}h/{{ dossiersMap[sel.id].formule.heures_pratiques || 20 }}h)
        </span>
        <button
          type="button"
          @click="emit('remove-candidat', sel.id)"
          class="rounded p-0.5 text-muted-foreground hover:bg-muted hover:text-rose-600"
          aria-label="Retirer"
        >
          <X class="h-3 w-3" />
        </button>
      </span>
    </div>

    <!-- Filtre catégorie -->
    <div class="flex items-center gap-1.5 overflow-x-auto pb-1 text-xs">
      <button
        v-for="cat in CATEGORIES_PERMIS"
        :key="cat.code"
        type="button"
        @click="selectedCategorie = cat.code"
        class="px-2.5 py-1 rounded-[4px] border text-[11px] font-semibold whitespace-nowrap transition cursor-pointer flex items-center gap-1"
        :class="
          selectedCategorie === cat.code
            ? 'bg-brand border-brand text-white font-bold'
            : 'bg-muted border-border text-muted-foreground hover:text-foreground'
        "
      >
        <span>{{ cat.icon }}</span>
        <span>{{ cat.label }}</span>
        <span class="font-mono text-[9px] opacity-80">
          ({{ countsParCategorie[cat.code] || 0 }})
        </span>
      </button>
    </div>

    <!-- Champ de recherche -->
    <div class="relative">
      <Search class="w-4 h-4 absolute left-3 top-2.5 text-muted-foreground" />
      <input
        type="text"
        v-model="searchQuery"
        placeholder="Rechercher par nom, prénom, téléphone ou dossier..."
        class="w-full pl-9 pr-3 py-2 text-xs rounded-[4px] bg-background border border-border focus:border-brand focus:ring-1 focus:ring-brand outline-none"
      />
    </div>

    <!-- Liste sélectionnable -->
    <div
      class="border border-border rounded-[4px] divide-y divide-border max-h-44 overflow-y-auto bg-muted/20"
    >
      <div
        v-if="candidatsAffiches.length === 0"
        class="p-4 text-center text-xs text-muted-foreground italic"
      >
        Aucun élève trouvé pour cette catégorie ou recherche.
      </div>

      <button
        v-for="c in candidatsAffiches"
        :key="c.id"
        type="button"
        @click="emit('toggle-candidat', c)"
        class="w-full px-3 py-2 text-left hover:bg-muted text-xs flex items-center justify-between transition cursor-pointer"
        :class="isSelected(c.id) ? 'bg-brand-light/50' : ''"
      >
        <div class="flex items-center space-x-2.5 min-w-0">
          <div
            class="w-5 h-5 rounded-[3px] border flex items-center justify-center shrink-0"
            :class="
              isSelected(c.id)
                ? 'bg-brand border-brand text-white'
                : 'border-border bg-card'
            "
          >
            <Check v-if="isSelected(c.id)" class="w-3.5 h-3.5" />
          </div>
          <div class="min-w-0">
            <p class="font-bold text-foreground truncate">
              {{ c.nom }} {{ c.prenoms }}
            </p>
            <p class="text-[10px] text-muted-foreground truncate">
              {{ c.telephone_whatsapp || c.telephone || "Non renseigné" }}
              <span v-if="c.formule_nom" class="font-medium text-brand">
                · {{ c.formule_nom }}
              </span>
            </p>
          </div>
        </div>

        <span
          class="text-[10px] px-1.5 py-0.5 rounded font-mono font-bold uppercase bg-sky-500/10 text-sky-600 border border-sky-500/20 shrink-0 ml-2"
        >
          {{ c.categorie_permis || extraireCategorie(c.formule_nom) }}
        </span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { nextTick } from "vue";
import { useGuichetStore } from "../../stores/guichetStore";

const store = useGuichetStore();

type View = "code" | "gerant" | "moniteur" | "candidat" | "examens";

const focusById = async (id: string) => {
  store.activeView = "guichet";
  await nextTick();
  document.getElementById(id)?.focus();
};

const toggleView = (view: View) => {
  store.activeView = store.activeView === view ? "guichet" : view;
};

const shortcuts: Array<{
  key: string;
  label: string;
  group: "actions" | "vues";
  view?: View;
  run: () => void;
}> = [
  {
    key: "F1",
    label: "Inscription",
    group: "actions",
    run: () => (store.isNewCandidateModalOpen = true),
  },
  {
    key: "F2",
    label: "Encaissement",
    group: "actions",
    run: () => focusById("montant-input"),
  },
  {
    key: "F3",
    label: "Recherche",
    group: "actions",
    run: () => focusById("search-input"),
  },
  {
    key: "F4",
    label: "Caisse Z",
    group: "actions",
    run: () => (store.isCashCloseModalOpen = true),
  },
  {
    key: "F5",
    label: "Présences code",
    group: "vues",
    view: "code",
    run: () => toggleView("code"),
  },
  {
    key: "F6",
    label: "Gérant",
    group: "vues",
    view: "gerant",
    run: () => toggleView("gerant"),
  },
  {
    key: "F7",
    label: "Moniteur",
    group: "vues",
    view: "moniteur",
    run: () => toggleView("moniteur"),
  },
  {
    key: "F8",
    label: "Candidat",
    group: "vues",
    view: "candidat",
    run: () => toggleView("candidat"),
  },
  {
    key: "F9",
    label: "Examens",
    group: "vues",
    view: "examens",
    run: () => toggleView("examens"),
  },
];

const isActive = (s: (typeof shortcuts)[number]) =>
  !!s.view && store.activeView === s.view;
</script>

<template>
  <nav
    class="flex items-center justify-between gap-4 border-b border-border bg-card px-4 py-1.5"
    aria-label="Raccourcis"
  >
    <ul class="flex min-w-0 items-center gap-1 overflow-x-auto">
      <template v-for="(s, i) in shortcuts" :key="s.key">
        <li
          v-if="i > 0 && shortcuts[i - 1].group !== s.group"
          class="mx-1.5 h-4 w-px shrink-0 bg-border"
          aria-hidden="true"
        />
        <li class="shrink-0">
          <button
            type="button"
            @click="s.run"
            :aria-keyshortcuts="s.key"
            :aria-pressed="s.view ? isActive(s) : undefined"
            class="inline-flex h-7 items-center gap-1.5 whitespace-nowrap rounded-[4px] px-2 text-xs font-medium transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40"
            :class="
              isActive(s)
                ? 'bg-brand-light text-white font-semibold'
                : 'text-muted-foreground hover:bg-muted hover:text-foreground'
            "
          >
            <kbd
              class="rounded px-1 font-mono text-[10px]"
              :class="
                isActive(s)
                  ? 'bg-brand-light text-white'
                  : 'bg-muted text-muted-foreground'
              "
            >
              {{ s.key }}
            </kbd>
            {{ s.label }}
          </button>
        </li>
      </template>
    </ul>

    <div
      class="hidden shrink-0 items-center gap-3 text-[11px] text-muted-foreground xl:flex"
    >
      <span
        ><kbd
          class="rounded bg-muted px-1 font-mono text-[10px] text-muted-foreground"
          >Tab</kbd
        >
        champs</span
      >
      <span
        ><kbd
          class="rounded bg-muted px-1 font-mono text-[10px] text-muted-foreground"
          >Entrée</kbd
        >
        valider</span
      >
      <span
        ><kbd
          class="rounded bg-muted px-1 font-mono text-[10px] text-muted-foreground"
          >Échap</kbd
        >
        annuler</span
      >
    </div>
  </nav>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useAuthStore } from "../../stores/authStore";
import { Printer, X } from "lucide-vue-next";

const authStore = useAuthStore();

type Student = {
  id: string;
  nom: string;
  numeroDossier: string;
  present: boolean;
  note: number;
};

const props = defineProps<{
  theme: string;
  formateur: string;
  heureDebut: string;
  heureFin: string;
  students: Student[];
}>();

const emit = defineEmits<{
  (e: "close"): void;
}>();

/* ───────────────────────── Calculs ───────────────────────── */
const today = new Date().toLocaleDateString("fr-FR");

const nbPresents = computed(
  () => props.students.filter((s) => s.present).length,
);
const nbAbsents = computed(() => props.students.length - nbPresents.value);

const duree = computed(() => {
  const toMin = (t: string) => {
    const [h, m] = (t || "").split(":").map(Number);
    return h * 60 + (m || 0);
  };
  const diff = toMin(props.heureFin) - toMin(props.heureDebut);
  if (!Number.isFinite(diff) || diff <= 0) return "";
  const h = Math.floor(diff / 60);
  const m = diff % 60;
  return m ? `${h} h ${String(m).padStart(2, "0")}` : `${h} h`;
});

const meta = computed(() => [
  { label: "Thème traité", value: props.theme },
  { label: "Formateur", value: props.formateur },
  {
    label: "Horaires",
    value: `${props.heureDebut} – ${props.heureFin}${duree.value ? ` (${duree.value})` : ""}`,
  },
]);

const printSheet = () => {
  window.print();
};
</script>

<template>
  <div
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 backdrop-blur-sm"
    role="dialog"
    aria-modal="true"
    aria-labelledby="sheet-title"
  >
    <div
      class="flex max-h-[94vh] w-full max-w-3xl flex-col overflow-hidden rounded-[4px] border border-border bg-card"
    >
      <!-- ═════════ En-tête ═════════ -->
      <header
        class="flex items-center justify-between gap-3 border-b border-border px-4 py-3"
      >
        <div class="min-w-0">
          <h2 id="sheet-title" class="text-sm font-semibold text-foreground">
            Feuille d'émargement
          </h2>
          <p class="text-[11px] text-muted-foreground">
            Registre officiel de code · format A4
          </p>
        </div>
        <button
          type="button"
          @click="emit('close')"
          class="flex h-7 w-7 items-center justify-center rounded-[4px] text-muted-foreground transition hover:bg-muted hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-dark/40"
          aria-label="Fermer"
        >
          <X class="h-4 w-4" />
        </button>
      </header>

      <!-- ═════════ Aperçu du document ═════════ -->
      <div class="flex-1 overflow-y-auto bg-muted p-4">
        <!-- L'id est conservé : il peut être ciblé par les styles d'impression -->
        <article
          id="contract-a4-container"
          class="select-text space-y-4 rounded-[4px] border border-slate-200 bg-white p-6 text-[11px] leading-normal text-slate-900"
        >
          <!-- Émetteur -->
          <header
            class="flex items-start justify-between gap-4 border-b border-slate-900 pb-3"
          >
            <div class="space-y-0.5">
              <div class="text-sm font-bold tracking-tight">
                {{ authStore.schoolName.toUpperCase() }}
              </div>
              <div class="text-[10px] font-medium text-slate-700">
                Agrément ministériel : {{ authStore.schoolAgrement }}
              </div>
              <div class="text-[10px] text-slate-500">
                {{ authStore.schoolAddress }}
              </div>
            </div>

            <div class="shrink-0 text-right">
              <div
                class="inline-block rounded border border-slate-900 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wide"
              >
                Registre de code théorique
              </div>
              <div class="mt-1 text-[10px] text-slate-500">
                Date : {{ today }}
              </div>
            </div>
          </header>

          <h3
            class="text-center text-[11px] font-semibold uppercase tracking-wider"
          >
            Feuille d'émargement et de contrôle pédagogique
          </h3>

          <!-- Séance -->
          <dl class="grid grid-cols-3 gap-4">
            <div v-for="m in meta" :key="m.label">
              <dt class="text-[10px] text-slate-500">{{ m.label }}</dt>
              <dd class="font-medium">{{ m.value }}</dd>
            </div>
          </dl>

          <!-- Émargements -->
          <table class="w-full border-collapse text-left">
            <thead>
              <tr
                class="border-y border-slate-300 bg-slate-50 text-[10px] uppercase tracking-wide text-slate-500"
              >
                <th
                  class="w-10 border-r border-slate-200 px-2 py-1.5 text-center font-medium"
                >
                  N°
                </th>
                <th class="border-r border-slate-200 px-2 py-1.5 font-medium">
                  Candidat
                </th>
                <th class="border-r border-slate-200 px-2 py-1.5 font-medium">
                  Dossier
                </th>
                <th
                  class="w-20 border-r border-slate-200 px-2 py-1.5 text-center font-medium"
                >
                  Présence
                </th>
                <th
                  class="w-16 border-r border-slate-200 px-2 py-1.5 text-center font-medium"
                >
                  Score
                </th>
                <th class="w-32 px-2 py-1.5 text-center font-medium">
                  Signature
                </th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="(s, index) in students"
                :key="s.id"
                class="border-b border-slate-200"
              >
                <td
                  class="border-r border-slate-200 px-2 py-2 text-center font-mono text-slate-500"
                >
                  {{ index + 1 }}
                </td>
                <td class="border-r border-slate-200 px-2 py-2 font-semibold">
                  {{ s.nom }}
                </td>
                <td
                  class="border-r border-slate-200 px-2 py-2 font-mono text-[10px] text-slate-600"
                >
                  {{ s.numeroDossier }}
                </td>
                <td
                  class="border-r border-slate-200 px-2 py-2 text-center font-medium"
                  :class="s.present ? 'text-slate-900' : 'text-slate-400'"
                >
                  {{ s.present ? "Présent" : "Absent" }}
                </td>
                <td
                  class="border-r border-slate-200 px-2 py-2 text-center font-mono font-semibold"
                >
                  {{ s.present ? `${s.note}/40` : "–" }}
                </td>
                <td
                  class="px-2 py-2 text-center text-[10px] italic text-slate-400"
                >
                  {{ s.present ? "[Émargé]" : "" }}
                </td>
              </tr>
            </tbody>
          </table>

          <!-- Synthèse & visa -->
          <div class="grid grid-cols-2 gap-6 pt-1">
            <dl class="space-y-1.5">
              <div class="flex justify-between gap-2">
                <dt class="text-slate-500">Total inscrits</dt>
                <dd class="font-semibold tabular-nums">
                  {{ students.length }}
                </dd>
              </div>
              <div class="flex justify-between gap-2">
                <dt class="text-slate-500">Présents en salle</dt>
                <dd class="font-semibold tabular-nums text-emerald-700">
                  {{ nbPresents }}
                </dd>
              </div>
              <div class="flex justify-between gap-2">
                <dt class="text-slate-500">Absents</dt>
                <dd class="font-semibold tabular-nums">{{ nbAbsents }}</dd>
              </div>
            </dl>

            <div
              class="flex h-24 flex-col justify-between rounded border border-slate-200 p-2.5 text-right"
            >
              <span class="text-[10px] uppercase tracking-wide text-slate-500"
                >Visa et signature du formateur</span
              >
              <span class="text-[10px] font-semibold text-slate-800">{{
                formateur
              }}</span>
            </div>
          </div>
        </article>
      </div>

      <!-- ═════════ Actions ═════════ -->
      <footer
        class="flex items-center justify-between gap-2 border-t border-border px-4 py-3"
      >
        <button
          type="button"
          @click="emit('close')"
          class="inline-flex h-8 items-center gap-2 rounded-[4px] border border-border bg-card px-3 text-xs font-medium text-foreground transition hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-dark/40"
        >
          Fermer
          <kbd
            class="rounded bg-muted px-1 font-mono text-[10px] text-muted-foreground"
            >Échap</kbd
          >
        </button>

        <button
          type="button"
          @click="printSheet"
          class="inline-flex h-8 items-center gap-1.5 rounded-[4px] bg-brand-dark px-3 text-xs font-medium text-white transition hover:bg-brand focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-dark/40 focus-visible:ring-offset-1"
        >
          <Printer class="h-3.5 w-3.5" />
          Imprimer la feuille A4
          <kbd
            class="rounded bg-white/20 px-1 font-mono text-[10px] text-white/90"
            >Entrée</kbd
          >
        </button>
      </footer>
    </div>
  </div>
</template>

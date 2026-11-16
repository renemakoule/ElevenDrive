<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from "vue";
import { useGuichetStore } from "../../stores/guichetStore";
import { useAuthStore } from "../../stores/authStore";
import { api } from "../../db/api";
import { realtimeSocket } from "../../services/realtimeSocketService";
import CodeAttendanceSheetModal from "./CodeAttendanceSheetModal.vue";
import { Check, Printer, Save, Users, Radio } from "lucide-vue-next";

const store = useGuichetStore();
const authStore = useAuthStore();
const isWsConnected = realtimeSocket.isConnected;
let unsubscribeWs: (() => void) | null = null;

/* ───────────────────────── Paramètres de séance ───────────────────────── */
const SEUIL_ADMISSIBLE = 32; // note minimale du test blanc (/40)
const SEUIL_EXCELLENT = 35;
const HEURES_CREDITEES = 2; // heures créditées par présence validée

const themeSelectionne = ref("Règles de Priorité & Intersections");
const THEMES_CODE = [
  "Signalisation Routière & Panneaux",
  "Règles de Priorité & Intersections",
  "Croisements, Dépassements & Voies Rapides",
  "Arrêt, Stationnement & Pannes",
  "Visibilité, Éclairage & Intempéries",
  "Vitesse, Distances de Sécurité & Freinage",
  "Alcool, Vigilance & Sécurité Routière",
  "Éco-conduite & Notions Mécaniques",
];

const moniteurs = ref<
  Array<{ id: string; nomComplet: string; role: string; telephone?: string }>
>([]);

const formateursList = computed(() => {
  const list = moniteurs.value
    .filter((m) => {
      const r = (m.role || "").toUpperCase();
      return r === "MONITEUR" || r === "DIRECTEUR";
    })
    .map((m) => (m.nomComplet || "").trim())
    .filter(Boolean);

  if (
    authStore.currentUser?.nom &&
    (authStore.currentRole === "moniteur" ||
      authStore.currentRole === "manager")
  ) {
    const me =
      `${authStore.currentUser.prenom || ""} ${authStore.currentUser.nom}`.trim();
    if (me && !list.includes(me)) {
      list.unshift(me);
    }
  }
  return Array.from(new Set(list));
});

const formateur = ref("");
const heureDebut = ref("14:00");
const heureFin = ref("16:00");
const isSheetModalOpen = ref(false);

interface CodeStudentItem {
  id: string;
  nom: string;
  numeroDossier: string;
  heuresActuelles: number;
  heuresMax: number;
  present: boolean;
  note: number;
}

const codeStudents = ref<CodeStudentItem[]>([]);

const syncStudentsFromStore = () => {
  if (store.searchResults && store.searchResults.length > 0) {
    codeStudents.value = store.searchResults.map((c) => ({
      id: c.id,
      nom: `${c.nom} ${c.prenoms}`.trim(),
      numeroDossier: c.numero_dossier || "DOS-ATTENTE",
      heuresActuelles: 0,
      heuresMax: 20,
      present: false,
      note: 0,
    }));
  } else {
    codeStudents.value = [];
  }
};

onMounted(async () => {
  try {
    const list = await api.getMoniteurs();
    moniteurs.value = Array.isArray(list) ? list : [];
    if (!formateur.value && formateursList.value.length > 0) {
      formateur.value = formateursList.value[0];
    }
  } catch {
    moniteurs.value = [];
  }
  if (!store.searchResults || store.searchResults.length === 0) {
    await store.search("");
  }
  syncStudentsFromStore();
  unsubscribeWs = realtimeSocket.subscribe(async (msg) => {
    if (
      msg.type === "CANDIDAT_DOSSIER_UPDATED" ||
      msg.type === "COURS_CODE_UPDATED"
    ) {
      await store.search("");
      syncStudentsFromStore();
    }
  });
});

onUnmounted(() => {
  if (unsubscribeWs) {
    unsubscribeWs();
    unsubscribeWs = null;
  }
});

watch(
  () => store.searchResults,
  () => {
    syncStudentsFromStore();
  },
  { deep: true },
);

/* ───────────────────────── Calculs & affichage ───────────────────────── */
const nbPresents = computed(
  () => codeStudents.value.filter((s) => s.present).length,
);

const pct = (value: number, max: number) =>
  max > 0 ? Math.min(100, Math.round((value / max) * 100)) : 0;
const isAssidu = (s: { heuresActuelles: number; heuresMax: number }) =>
  s.heuresActuelles >= s.heuresMax * 0.8;

const mention = (note: number) => {
  if (note >= SEUIL_EXCELLENT)
    return { label: "Excellent", dot: "bg-brand", text: "text-brand-dark" };
  if (note >= SEUIL_ADMISSIBLE)
    return { label: "Admissible", dot: "bg-sky-500", text: "text-sky-700" };
  return { label: "Ajourné", dot: "bg-rose-500", text: "text-rose-700" };
};

/* ───────────────────────── Actions ───────────────────────── */
const togglePresent = (student: any) => {
  student.present = !student.present;
  if (student.present && student.note === 0) {
    student.note = SEUIL_ADMISSIBLE;
  }
};

const clampNote = (student: any) => {
  const n = Number(student.note);
  student.note = Number.isFinite(n)
    ? Math.max(0, Math.min(40, Math.round(n)))
    : 0;
};

const saveAttendance = async () => {
  if (!formateur.value) {
    store.showToast("Veuillez sélectionner un formateur responsable.", "error");
    return;
  }
  const presencesPayload = codeStudents.value.map((s) => ({
    candidat_id: s.id,
    est_present: s.present,
    score_test: s.present ? s.note : 0,
  }));

  await store.enregistrerSeanceCode({
    theme_traite: themeSelectionne.value,
    formateur_nom: formateur.value,
    heure_debut: heureDebut.value,
    heure_fin: heureFin.value,
    presences: presencesPayload,
  });

  // Mise à jour des compteurs locaux d'heures
  codeStudents.value.forEach((s) => {
    if (s.present) {
      s.heuresActuelles += HEURES_CREDITEES;
    }
  });
};

const fieldClass =
  "h-8 w-full rounded-[4px] border border-border bg-card px-2.5 text-xs text-foreground transition focus:border-brand focus:outline-none focus:ring-2 focus:ring-brand/20";
</script>

<template>
  <div class="flex flex-1 flex-col gap-4 overflow-y-auto">
    <!-- ═════════ En-tête ═════════ -->
    <header class="flex items-start justify-between gap-3">
      <div class="min-w-0">
        <h1
          class="flex items-center gap-2 text-base font-semibold text-foreground"
        >
          Salle de code

          <kbd
            class="rounded bg-muted px-1 font-mono text-[10px] font-normal text-muted-foreground"
            >F5</kbd
          >
        </h1>
        <p class="text-xs text-muted-foreground">
          Chaque présence validée crédite +{{ HEURES_CREDITEES }} h de théorie
          sur le dossier de l'apprenant.
        </p>
      </div>

      <div class="shrink-0 text-right">
        <p class="text-xs text-muted-foreground">Présents en salle</p>
        <p class="text-lg font-semibold tabular-nums text-foreground">
          {{ nbPresents
          }}<span class="text-muted-foreground">
            / {{ codeStudents.length }}</span
          >
        </p>
      </div>
    </header>

    <!-- ═════════ Paramètres de la séance ═════════ -->
    <section
      class="grid grid-cols-1 gap-3 rounded-[4px] border border-border bg-card p-4 md:grid-cols-[1fr_1fr_auto]"
    >
      <div>
        <label for="theme" class="mb-1 block text-xs text-muted-foreground"
          >Thème traité</label
        >
        <select id="theme" v-model="themeSelectionne" :class="fieldClass">
          <option v-for="t in THEMES_CODE" :key="t" :value="t">{{ t }}</option>
        </select>
      </div>

      <div>
        <label for="formateur" class="mb-1 block text-xs text-muted-foreground"
          >Formateur responsable</label
        >
        <select id="formateur" v-model="formateur" :class="fieldClass">
          <option value="" disabled>
            {{
              formateursList.length === 0
                ? "Aucun moniteur / formateur enregistré"
                : "Sélectionner un formateur..."
            }}
          </option>
          <option v-for="f in formateursList" :key="f" :value="f">
            {{ f }}
          </option>
        </select>
      </div>

      <div class="flex items-end gap-2">
        <div class="w-24">
          <label
            for="heure-debut"
            class="mb-1 block text-xs text-muted-foreground"
            >Début</label
          >
          <input
            id="heure-debut"
            type="time"
            v-model="heureDebut"
            :class="[fieldClass, 'text-center tabular-nums']"
          />
        </div>
        <div class="w-24">
          <label
            for="heure-fin"
            class="mb-1 block text-xs text-muted-foreground"
            >Fin</label
          >
          <input
            id="heure-fin"
            type="time"
            v-model="heureFin"
            :class="[fieldClass, 'text-center tabular-nums']"
          />
        </div>
      </div>
    </section>

    <!-- ═════════ Pointage ═════════ -->
    <section
      class="flex flex-1 flex-col overflow-hidden rounded-[4px] border border-border bg-card"
    >
      <div class="flex items-center justify-between gap-3 px-4 py-3">
        <h2 class="text-sm font-semibold text-foreground">
          Apprenants inscrits
        </h2>
        <span class="text-[11px] text-muted-foreground"
          >Test blanc : admissible dès {{ SEUIL_ADMISSIBLE }} / 40</span
        >
      </div>

      <ul
        class="flex-1 divide-y divide-border overflow-y-auto border-t border-border"
      >
        <li
          v-if="codeStudents.length === 0"
          class="flex flex-col items-center justify-center py-12 px-4 text-center"
        >
          <Users class="w-8 h-8 text-muted-foreground/50 mb-2" />
          <p class="text-xs font-medium text-foreground">
            Aucun apprenant enregistré dans l'auto-école
          </p>
          <p class="text-[11px] text-muted-foreground mt-1 max-w-sm">
            Enregistrez des candidats via le guichet (touche F1 Inscription)
            pour pouvoir effectuer leur émargement et noter leurs tests de code.
          </p>
        </li>

        <li
          v-for="s in codeStudents"
          :key="s.id"
          class="flex items-center gap-3 px-4 py-2.5 transition"
          :class="s.present ? 'bg-brand-light/40' : 'hover:bg-muted/60'"
        >
          <!-- Zone de pointage (clic = présent / absent) -->
          <button
            type="button"
            role="checkbox"
            :aria-checked="s.present"
            @click="togglePresent(s)"
            class="flex min-w-0 flex-1 items-center gap-3 rounded-[4px] text-left focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40"
          >
            <span
              class="flex h-4 w-4 shrink-0 items-center justify-center rounded border transition"
              :class="
                s.present
                  ? 'border-brand-dark bg-brand-dark text-white'
                  : 'border-border bg-card text-transparent'
              "
            >
              <Check class="h-3 w-3" />
            </span>

            <span class="block min-w-0 flex-1">
              <span class="flex items-center gap-2">
                <span class="truncate text-xs font-medium text-foreground">{{
                  s.nom
                }}</span>
                <span
                  class="shrink-0 font-mono text-[10px] text-muted-foreground"
                  >{{ s.numeroDossier }}</span
                >
              </span>

              <span class="mt-1.5 flex items-center gap-2">
                <span
                  class="block h-1 w-24 overflow-hidden rounded-full bg-muted"
                >
                  <span
                    class="block h-full rounded-full bg-brand transition-all duration-500"
                    :style="{
                      width: pct(s.heuresActuelles, s.heuresMax) + '%',
                    }"
                  />
                </span>
                <span class="text-[11px] tabular-nums text-muted-foreground"
                  >{{ s.heuresActuelles }} / {{ s.heuresMax }} h</span
                >
                <span
                  v-if="isAssidu(s)"
                  class="text-[11px] font-medium text-brand-dark"
                  >· Assidu</span
                >
              </span>
            </span>
          </button>

          <!-- Mention & score -->
          <div class="flex shrink-0 items-center gap-3">
            <span
              v-if="s.present"
              class="flex items-center gap-1.5 text-[11px] font-medium"
              :class="mention(s.note).text"
            >
              <span
                class="h-1.5 w-1.5 rounded-full"
                :class="mention(s.note).dot"
              />
              {{ mention(s.note).label }}
            </span>
            <span v-else class="text-[11px] text-muted-foreground">Absent</span>

            <div class="relative w-16">
              <label :for="`score-${s.id}`" class="sr-only"
                >Score de {{ s.nom }}</label
              >
              <input
                :id="`score-${s.id}`"
                type="number"
                v-model.number="s.note"
                min="0"
                max="40"
                :disabled="!s.present"
                @blur="clampNote(s)"
                class="h-7 w-full rounded-[4px] border border-border bg-card pl-1 pr-6 text-right text-xs font-semibold tabular-nums text-foreground transition focus:border-brand focus:outline-none focus:ring-2 focus:ring-brand/20 disabled:bg-muted disabled:text-muted-foreground"
              />
              <span
                class="pointer-events-none absolute right-1.5 top-1/2 -translate-y-1/2 text-[10px] text-muted-foreground"
                >/40</span
              >
            </div>
          </div>
        </li>
      </ul>

      <!-- Actions -->
      <footer
        class="flex items-center justify-between gap-2 border-t border-border bg-muted px-4 py-3"
      >
        <button
          type="button"
          @click="isSheetModalOpen = true"
          class="inline-flex h-8 items-center gap-1.5 rounded-[4px] border border-border bg-card px-3 text-xs font-medium text-foreground transition hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40"
        >
          <Printer class="h-3.5 w-3.5 text-muted-foreground" />
          Feuille d'émargement A4
        </button>

        <button
          type="button"
          @click="saveAttendance"
          :disabled="store.loading"
          class="inline-flex h-8 items-center gap-1.5 rounded-[4px] bg-brand-dark px-3 text-xs font-medium text-white transition hover:bg-brand focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand/40 focus-visible:ring-offset-1 disabled:cursor-not-allowed disabled:opacity-60"
        >
          <Save class="h-3.5 w-3.5" />
          Valider la séance (+{{ HEURES_CREDITEES }} h)
        </button>
      </footer>
    </section>

    <!-- ═════════ Modale feuille A4 ═════════ -->
    <CodeAttendanceSheetModal
      v-if="isSheetModalOpen"
      :theme="themeSelectionne"
      :formateur="formateur"
      :heure-debut="heureDebut"
      :heure-fin="heureFin"
      :students="codeStudents"
      @close="isSheetModalOpen = false"
    />
  </div>
</template>

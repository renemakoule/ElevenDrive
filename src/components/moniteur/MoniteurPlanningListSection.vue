<script setup lang="ts">
import type { SeanceConduite } from "../../types";
import {
  Phone,
  MessageCircle,
  ChevronRight,
  Calendar,
  Check,
} from "lucide-vue-next";

const props = defineProps<{
  viewMode: "jour" | "semaine";
  sessions: SeanceConduite[];
  selectedDate: string;
  todayIso: string;
  dateAfficheeLongue: string;
  defaultKilometrage: number;
  weekDays: Array<{
    date: string;
    dayName: string;
    dayNumber: number;
    isToday: boolean;
    isSelected: boolean;
    count: number;
  }>;
  weekSessionsMap: Record<string, SeanceConduite[]>;
}>();

const emit = defineEmits<{
  openSession: [session: SeanceConduite];
  quickValidate: [session: SeanceConduite];
  resetToday: [];
  selectDayFromWeek: [iso: string];
}>();

const statusMap = {
  EFFECTUE: { label: "Terminé", dot: "bg-brand", text: "text-brand-dark" },
  EN_COURS: {
    label: "En cours",
    dot: "bg-sky-500 animate-pulse",
    text: "text-sky-600 dark:text-sky-400",
  },
  PROGRAMME: {
    label: "À venir",
    dot: "bg-muted-foreground/40",
    text: "text-muted-foreground",
  },
} as const;

const statusOf = (s: SeanceConduite) =>
  statusMap[s.statut as keyof typeof statusMap] ?? statusMap.PROGRAMME;

const pct = (value: number, max: number) =>
  max > 0 ? Math.min(100, Math.round((value / max) * 100)) : 0;

const num = (n: number) => n.toLocaleString("fr-FR");

const whatsappLink = (tel: string) =>
  `https://wa.me/${tel.replace(/[^0-9]/g, "")}`;

const trajet = (s: SeanceConduite) => {
  const debut = s.km_debut ?? props.defaultKilometrage ?? 0;
  const fin = s.km_fin ?? debut;
  return { debut, fin, distance: Math.max(0, fin - debut) };
};
</script>

<template>
  <!-- ═════════ Planning Réel (Vue Journée) ═════════ -->
  <section
    v-if="viewMode === 'jour'"
    class="overflow-hidden rounded-[4px] border border-border bg-card lg:col-span-8"
  >
    <div class="flex items-center justify-between px-4 py-3">
      <div class="flex items-center gap-2">
        <h2 class="text-sm font-semibold text-foreground capitalize">
          {{ dateAfficheeLongue }}
        </h2>
        <span
          v-if="selectedDate === todayIso"
          class="rounded-full border border-brand-border bg-brand-light px-2 py-0.5 text-[10px] font-semibold text-primary"
        >
          Aujourd'hui
        </span>
      </div>
      <span class="text-xs text-muted-foreground">
        {{ sessions.length }} leçon(s) planifiée(s)
      </span>
    </div>

    <ul
      v-if="sessions.length"
      class="divide-y divide-border border-t border-border"
    >
      <li
        v-for="session in sessions"
        :key="session.id"
        class="flex flex-col gap-3 px-4 py-3.5 transition sm:flex-row sm:items-center sm:gap-4"
        :class="
          session.statut === 'EN_COURS' ? 'bg-sky-500/5' : 'hover:bg-muted/50'
        "
      >
        <!-- Horaire -->
        <div class="flex shrink-0 items-baseline gap-2 sm:block sm:w-14">
          <div class="text-sm font-semibold tabular-nums text-foreground">
            {{ session.heure_debut }}
          </div>
          <div class="text-[11px] tabular-nums text-muted-foreground">
            {{ session.heure_fin }}
          </div>
        </div>

        <!-- Élève & progression réelle -->
        <div class="min-w-0 flex-1 space-y-2">
          <div>
            <div class="flex items-center gap-2">
              <h3 class="truncate text-sm font-semibold text-foreground">
                {{ session.candidat_nom }}
              </h3>
              <span
                class="flex shrink-0 items-center gap-1.5 text-[11px] font-medium"
                :class="statusOf(session).text"
              >
                <span
                  class="h-1.5 w-1.5 rounded-full"
                  :class="statusOf(session).dot"
                />
                {{ statusOf(session).label }}
              </span>
            </div>
            <p class="text-xs text-muted-foreground">
              {{ session.formule_nom }}
            </p>
          </div>

          <div class="flex flex-wrap items-center gap-x-4 gap-y-1.5">
            <div
              class="flex min-w-[170px] max-w-[240px] flex-1 items-center gap-2"
            >
              <div
                class="h-1 flex-1 overflow-hidden rounded-full bg-muted"
                role="progressbar"
                :aria-valuenow="
                  pct(
                    session.heures_deja_effectuees,
                    session.heures_forfait_total,
                  )
                "
                aria-valuemin="0"
                aria-valuemax="100"
                :aria-label="`Progression de ${session.candidat_nom}`"
              >
                <div
                  class="h-full rounded-full bg-brand transition-all duration-500"
                  :style="{
                    width:
                      pct(
                        session.heures_deja_effectuees,
                        session.heures_forfait_total,
                      ) + '%',
                  }"
                />
              </div>
              <span
                class="shrink-0 text-[11px] tabular-nums text-muted-foreground"
              >
                {{ session.heures_deja_effectuees }} /
                {{ session.heures_forfait_total }} h
              </span>
            </div>

            <span
              v-if="
                session.statut !== 'PROGRAMME' &&
                session.km_debut != null &&
                session.km_fin != null
              "
              class="text-[11px] tabular-nums text-muted-foreground"
            >
              {{ num(trajet(session).debut) }} →
              {{ num(trajet(session).fin) }} km
              <span class="font-medium text-foreground"
                >· +{{ trajet(session).distance }} km</span
              >
            </span>
          </div>
        </div>

        <!-- Actions tactiles -->
        <div class="flex flex-wrap sm:flex-nowrap shrink-0 items-center gap-1.5 pt-1 sm:pt-0">
          <a
            v-if="session.telephone_whatsapp"
            :href="`tel:${session.telephone_whatsapp}`"
            class="flex h-10 w-10 sm:h-8 sm:w-8 items-center justify-center rounded-[3px] border border-border bg-card text-muted-foreground transition hover:bg-muted hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-dark/40"
            aria-label="Appeler l'élève"
          >
            <Phone class="h-4 w-4 sm:h-3.5 sm:w-3.5" />
          </a>
          <a
            v-if="session.telephone_whatsapp"
            :href="whatsappLink(session.telephone_whatsapp)"
            target="_blank"
            rel="noopener noreferrer"
            class="flex h-10 w-10 sm:h-8 sm:w-8 items-center justify-center rounded-[3px] border border-border bg-card text-primary transition hover:bg-brand-light focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-dark/40"
            aria-label="Écrire sur WhatsApp"
          >
            <MessageCircle class="h-4 w-4 sm:h-3.5 sm:w-3.5" />
          </a>
          <!-- Bouton direct 1-clic pour valider la présence si leçon non terminée -->
          <button
            v-if="session.statut !== 'EFFECTUE'"
            type="button"
            @click.stop="emit('quickValidate', session)"
            class="flex-1 sm:flex-initial inline-flex h-10 sm:h-8 items-center justify-center gap-1.5 rounded-[3px] px-3 sm:px-2.5 text-xs font-semibold bg-emerald-600 hover:bg-emerald-700 text-white transition active:scale-95"
            title="Cocher l'élève présent et valider la leçon en 1 clic"
          >
            <Check class="h-3.5 w-3.5 stroke-[3]" />
            <span class="hidden sm:inline">Présent (1 clic)</span>
            <span class="sm:hidden">Présent</span>
          </button>

          <button
            type="button"
            @click="emit('openSession', session)"
            class="flex-1 sm:flex-initial inline-flex h-10 sm:h-8 items-center justify-center gap-1 rounded-[3px] px-3 text-xs font-medium transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-dark/40"
            :class="
              session.statut === 'EFFECTUE'
                ? 'border border-border bg-card text-foreground hover:bg-muted'
                : 'bg-brand-dark text-white hover:bg-brand'
            "
          >
            {{
              session.statut === "EFFECTUE"
                ? "Voir le livret"
                : "Émarger et valider"
            }}
            <ChevronRight class="h-3.5 w-3.5" />
          </button>
        </div>
      </li>
    </ul>

    <div
      v-else
      class="p-10 text-center text-xs text-muted-foreground border-t border-border space-y-3"
    >
      <div
        class="mx-auto flex h-10 w-10 items-center justify-center rounded-full bg-muted text-muted-foreground"
      >
        <Calendar class="h-5 w-5" />
      </div>
      <p class="font-medium text-foreground">
        Aucune leçon de conduite planifiée pour cette date
      </p>
      <p class="text-[11px] max-w-xs mx-auto">
        Sélectionnez un autre jour dans le bandeau de la semaine ci-dessus ou
        basculez en vue semaine.
      </p>
      <button
        v-if="selectedDate !== todayIso"
        type="button"
        @click="emit('resetToday')"
        class="inline-flex items-center gap-1.5 rounded-[4px] border border-brand-border bg-brand-light px-3 py-1.5 text-xs font-semibold text-brand-dark hover:bg-brand-light/80 transition"
      >
        Revenir à aujourd'hui
      </button>
    </div>
  </section>

  <!-- ═════════ VUE SEMAINE COMPLÈTE (Vue Hebdomadaire) ═════════ -->
  <section
    v-else
    class="overflow-hidden rounded-[4px] border border-border bg-card lg:col-span-8 p-4 space-y-3"
  >
    <div class="flex items-center justify-between pb-3 border-b border-border">
      <div>
        <h2 class="text-sm font-semibold text-foreground">
          Feuille de route hebdomadaire
        </h2>
        <p class="text-xs text-muted-foreground">
          Aperçu global de vos leçons pour les 7 jours de la semaine
        </p>
      </div>
      <span
        class="rounded-md border border-brand-border bg-brand-light px-2.5 py-1 text-xs font-semibold text-brand-dark"
      >
        Semaine du {{ weekDays[0]?.dayNumber }} au
        {{ weekDays[6]?.dayNumber }}
      </span>
    </div>

    <div class="grid gap-2.5 sm:grid-cols-2">
      <div
        v-for="day in weekDays"
        :key="day.date"
        class="rounded-[4px] border p-3.5 transition flex flex-col justify-between"
        :class="[
          day.isSelected
            ? 'border-brand ring-1 ring-brand bg-brand-light/20'
            : 'border-border bg-muted/20 hover:bg-muted/50',
        ]"
      >
        <div class="flex items-start justify-between">
          <div>
            <span
              class="text-[11px] font-bold text-muted-foreground uppercase"
              >{{ day.dayName }}</span
            >
            <div class="flex items-baseline gap-2 mt-0.5">
              <h3 class="text-base font-bold text-foreground">
                {{ day.dayNumber }}
              </h3>
              <span
                v-if="day.isToday"
                class="text-[10px] font-semibold text-brand bg-brand-light px-1.5 py-0.2 rounded"
                >Aujourd'hui</span
              >
            </div>
          </div>
          <span
            class="inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-medium"
            :class="
              day.count > 0
                ? 'bg-sky-500/10 text-sky-700 dark:text-sky-300 border border-sky-500/20'
                : 'bg-muted text-muted-foreground border border-border'
            "
          >
            {{ day.count }} leçon{{ day.count > 1 ? "s" : "" }}
          </span>
        </div>

        <!-- Aperçu de la première leçon si existante -->
        <div
          v-if="
            weekSessionsMap[day.date] && weekSessionsMap[day.date].length > 0
          "
          class="my-2.5 space-y-1 text-xs border-t border-border/60 pt-2"
        >
          <div
            v-for="s in weekSessionsMap[day.date].slice(0, 2)"
            :key="s.id"
            class="flex items-center justify-between text-muted-foreground"
          >
            <span class="truncate font-medium text-foreground max-w-[140px]">{{
              s.candidat_nom
            }}</span>
            <span class="font-mono text-[10px]">{{ s.heure_debut }}</span>
          </div>
          <span
            v-if="weekSessionsMap[day.date].length > 2"
            class="text-[10px] text-muted-foreground italic block"
          >
            +{{ weekSessionsMap[day.date].length - 2 }} autre(s) élève(s)...
          </span>
        </div>
        <div v-else class="my-2 text-[11px] text-muted-foreground italic">
          Aucune leçon planifiée
        </div>

        <button
          type="button"
          @click="emit('selectDayFromWeek', day.date)"
          class="w-full mt-2 inline-flex items-center justify-center gap-1 rounded-[4px] border border-border bg-background py-1.5 text-xs font-medium text-foreground hover:bg-muted transition"
        >
          <span>Consulter la journée</span>
          <ChevronRight class="h-3 w-3 text-muted-foreground" />
        </button>
      </div>
    </div>
  </section>
</template>

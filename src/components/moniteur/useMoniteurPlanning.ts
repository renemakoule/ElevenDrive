import { ref, reactive, computed, onMounted, onUnmounted } from "vue";
import { useGuichetStore } from "../../stores/guichetStore";
import { useCoursCodeStore } from "../../stores/coursCodeStore";
import { useSupportsCoursStore } from "../../stores/supportsCoursStore";
import { useEpreuvesStore } from "../../stores/epreuvesStore";
import { useQuizFormsStore } from "../../stores/quizFormsStore";
import { api } from "../../db/api";
import { realtimeSocket } from "../../services/realtimeSocketService";
import { speechAudioService } from "../../services/speechAudioService";
import type { SeanceConduite } from "../../types";

export function useMoniteurPlanning(getMoniteurId: () => string) {
  const store = useGuichetStore();
  const coursStore = useCoursCodeStore();
  const supportsStore = useSupportsCoursStore();
  const epreuvesStore = useEpreuvesStore();
  const quizStore = useQuizFormsStore();

  const isWsConnected = realtimeSocket.isConnected;
  const lastRealtimeUpdate = ref<string>("");
  let unsubscribeWs: (() => void) | null = null;

  const rawSessions = ref<SeanceConduite[]>([]);
  const isLoading = ref(false);

  const todayIso = new Date().toISOString().slice(0, 10);
  const selectedDate = ref<string>(todayIso);
  const viewMode = ref<"jour" | "semaine">("jour");
  const weekSessionsMap = reactive<Record<string, SeanceConduite[]>>({});

  const weekDays = computed(() => {
    const [y, m, d] = selectedDate.value.split("-").map(Number);
    const curr = new Date(y, m - 1, d);
    const day = curr.getDay();
    const diff = curr.getDate() - day + (day === 0 ? -6 : 1);
    const monday = new Date(y, m - 1, diff);

    const list = [];
    for (let i = 0; i < 7; i++) {
      const target = new Date(monday);
      target.setDate(monday.getDate() + i);
      const iso = target.toISOString().slice(0, 10);
      const dayName = new Intl.DateTimeFormat("fr-FR", {
        weekday: "short",
      }).format(target);
      const listSessions = weekSessionsMap[iso] || [];
      list.push({
        date: iso,
        dayName: dayName.replace(".", "").toUpperCase(),
        dayNumber: target.getDate(),
        isToday: iso === todayIso,
        isSelected: iso === selectedDate.value,
        count: listSessions.length,
      });
    }
    return list;
  });

  const dateAfficheeLongue = computed(() => {
    const [y, m, d] = selectedDate.value.split("-").map(Number);
    const target = new Date(y, m - 1, d);
    return new Intl.DateTimeFormat("fr-FR", {
      weekday: "long",
      day: "numeric",
      month: "long",
      year: "numeric",
    }).format(target);
  });

  const loadWeekSessions = async (force = false) => {
    try {
      const days = weekDays.value;
      for (const d of days) {
        if (force || !weekSessionsMap[d.date]) {
          api
            .getPlanningConduite(d.date, getMoniteurId())
            .then((res) => {
              weekSessionsMap[d.date] = res;
            })
            .catch(() => {});
        }
      }
    } catch {
      // Ignorer erreurs silencieuses
    }
  };

  const loadData = async () => {
    isLoading.value = true;
    try {
      const [sessionsData] = await Promise.all([
        api.getPlanningConduite(selectedDate.value, getMoniteurId()),
        store.loadVehicules(),
      ]);
      rawSessions.value = sessionsData;
      weekSessionsMap[selectedDate.value] = sessionsData;
      loadWeekSessions();
    } catch (err: any) {
      store.showToast("Erreur chargement planning: " + err.message, "error");
    } finally {
      isLoading.value = false;
    }
  };

  const selectDate = async (iso: string) => {
    selectedDate.value = iso;
    await loadData();
  };

  const goToPrevDay = async () => {
    const [y, m, d] = selectedDate.value.split("-").map(Number);
    selectedDate.value = new Date(y, m - 1, d - 1).toISOString().slice(0, 10);
    await loadData();
  };

  const goToNextDay = async () => {
    const [y, m, d] = selectedDate.value.split("-").map(Number);
    selectedDate.value = new Date(y, m - 1, d + 1).toISOString().slice(0, 10);
    await loadData();
  };

  const resetToToday = async () => {
    selectedDate.value = todayIso;
    await loadData();
  };

  const handleSelectDayFromWeek = async (iso: string) => {
    await selectDate(iso);
    viewMode.value = "jour";
  };

  onMounted(async () => {
    realtimeSocket.connect();
    await Promise.all([
      loadData(),
      coursStore.loadSeances(),
      supportsStore.loadSupports(),
      epreuvesStore.loadEpreuves(),
      quizStore.loadFormulaires(),
      quizStore.loadSoumissions(),
    ]);

    unsubscribeWs = realtimeSocket.subscribe(async (msg) => {
      lastRealtimeUpdate.value = new Date().toLocaleTimeString("fr-FR");
      if (
        msg.type === "PLANNING_CONDUITE_UPDATED" ||
        msg.type === "LIVRET_EMARGEMENT_UPDATED"
      ) {
        if (msg.data?.date_seance && typeof msg.data.date_seance === "string") {
          const hasOnlyGenerated = rawSessions.value.every((s) =>
            s.id.startsWith("plan-"),
          );
          if (hasOnlyGenerated && msg.data.date_seance !== selectedDate.value) {
            selectedDate.value = msg.data.date_seance.slice(0, 10);
          }
        }
        Object.keys(weekSessionsMap).forEach((k) => delete weekSessionsMap[k]);
        await loadData();
        loadWeekSessions(true);
      } else if (msg.type === "FLOTTE_VEHICULES_UPDATED") {
        await store.loadVehicules();
      } else if (msg.type === "COURS_CODE_UPDATED") {
        await coursStore.loadSeances();
      } else if (msg.type === "SUPPORTS_COURS_UPDATED") {
        await supportsStore.loadSupports();
      } else if (msg.type === "EPREUVES_UPDATED") {
        await epreuvesStore.loadEpreuves();
      } else if (
        msg.type === "QUIZ_FORMS_UPDATED" ||
        msg.type === "QUIZ_SUBMISSION_RECEIVED"
      ) {
        await Promise.all([
          quizStore.loadFormulaires(),
          quizStore.loadSoumissions(),
        ]);
      } else if (msg.type === "LIVRET_EMARGEMENT_LIVE" && msg.data?.sessionId) {
        const target = rawSessions.value.find(
          (s) => s.id === msg.data.sessionId,
        );
        if (target) {
          if (Array.isArray(msg.data.competences))
            target.competences = msg.data.competences;
          if (typeof msg.data.kmDebut === "number")
            target.km_debut = msg.data.kmDebut;
          if (typeof msg.data.kmFin === "number")
            target.km_fin = msg.data.kmFin;
          if (typeof msg.data.commentaire === "string")
            target.commentaire_moniteur = msg.data.commentaire;
        }
      }
    });
  });

  onUnmounted(() => {
    speechAudioService.arreter();
    if (unsubscribeWs) {
      unsubscribeWs();
      unsubscribeWs = null;
    }
  });

  return {
    isWsConnected,
    rawSessions,
    isLoading,
    todayIso,
    selectedDate,
    viewMode,
    weekSessionsMap,
    weekDays,
    dateAfficheeLongue,
    loadData,
    selectDate,
    goToPrevDay,
    goToNextDay,
    resetToToday,
    handleSelectDayFromWeek,
  };
}


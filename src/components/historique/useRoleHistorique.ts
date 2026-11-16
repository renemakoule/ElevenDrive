import { ref, computed, onMounted, onUnmounted } from "vue";
import { useAuthStore } from "../../stores/authStore";
import { useGuichetStore } from "../../stores/guichetStore";
import { useCoursCodeStore } from "../../stores/coursCodeStore";
import { useQuizFormsStore } from "../../stores/quizFormsStore";
import { realtimeSocket } from "../../services/realtimeSocketService";
import { api } from "../../db/api";
import type { SeanceConduite } from "../../types";
import { fcfa, type HistoriqueItem } from "./types";

export function useRoleHistorique() {
  const authStore = useAuthStore();
  const guichetStore = useGuichetStore();
  const coursStore = useCoursCodeStore();
  const quizStore = useQuizFormsStore();

  const isWsConnected = realtimeSocket.isConnected;
  const lastSyncAt = ref<string>("");
  const isRefreshing = ref(false);
  const drivingSessions = ref<SeanceConduite[]>([]);
  const liveEvents = ref<HistoriqueItem[]>([]);
  let unsubscribeWs: (() => void) | null = null;

  const currentRole = computed(() =>
    (authStore.currentRole || "secretaire").toLowerCase(),
  );

  const roleTitle = computed(() => {
    switch (currentRole.value) {
      case "secretaire":
      case "caissier":
        return "Guichet & Caisse";
      case "manager":
      case "directeur":
        return "Directeur & Gérance";
      case "moniteur":
        return "Moniteur & Enseignement";
      case "eleve":
      case "candidat":
        return "Espace Apprenant";
      default:
        return "Utilisateur";
    }
  });

  const refreshAllData = async () => {
    isRefreshing.value = true;
    try {
      const promises: Promise<any>[] = [
        guichetStore.loadStats(),
        coursStore.loadSeances(),
        quizStore.loadAll(),
        guichetStore.loadVehicules(),
      ];

      if (
        currentRole.value === "manager" ||
        currentRole.value === "directeur"
      ) {
        promises.push(guichetStore.loadGerantDashboard());
      }

      if (
        !guichetStore.searchResults ||
        guichetStore.searchResults.length === 0
      ) {
        promises.push(guichetStore.search(""));
      }

      try {
        const pl = await api.getPlanningConduite("ALL");
        drivingSessions.value = pl || [];
      } catch {
        // ignore
      }

      await Promise.allSettled(promises);
      lastSyncAt.value = new Date().toLocaleTimeString("fr-FR");
    } finally {
      isRefreshing.value = false;
    }
  };

  const baseHistoryItems = computed<HistoriqueItem[]>(() => {
    const items: HistoriqueItem[] = [];
    const currentUserId = authStore.currentUser?.id;
    const currentCandDossierId = guichetStore.activeDossierDetail?.dossier?.id;
    const currentCandName = guichetStore.activeDossierDetail
      ? `${guichetStore.activeDossierDetail.candidat.nom} ${guichetStore.activeDossierDetail.candidat.prenoms}`.trim()
      : `${authStore.currentUser?.nom || ""} ${authStore.currentUser?.prenom || ""}`.trim();

    if (
      guichetStore.gerantDashboardData?.derniers_paiements &&
      currentRole.value !== "eleve"
    ) {
      guichetStore.gerantDashboardData.derniers_paiements.forEach((p) => {
        const d = new Date(p.date_paiement);
        items.push({
          id: `gerant-paiement-${p.id}`,
          timestamp: d.getTime() || Date.now(),
          dateStr: d.toLocaleDateString("fr-FR"),
          heureStr: d.toLocaleTimeString("fr-FR", {
            hour: "2-digit",
            minute: "2-digit",
          }),
          type: "PAIEMENT",
          titre: `Paiement guichet : ${fcfa(p.montant)} FCFA`,
          description: `Règlement ${p.mode_reglement} pour ${p.candidat_nom} • Reçu n° ${p.numero_recu} (${p.formule_nom})`,
          acteur: "Guichet",
          badgeLabel: p.mode_reglement,
          badgeClass:
            "bg-emerald-500/15 text-emerald-700 dark:text-emerald-400 border-emerald-500/20",
          montant: p.montant,
        });
      });
    }

    if (guichetStore.activeDossierDetail?.paiements) {
      guichetStore.activeDossierDetail.paiements.forEach((p) => {
        const d = new Date(p.date_paiement);
        items.push({
          id: `paiement-${p.id}`,
          timestamp: d.getTime() || Date.now(),
          dateStr: d.toLocaleDateString("fr-FR"),
          heureStr: d.toLocaleTimeString("fr-FR", {
            hour: "2-digit",
            minute: "2-digit",
          }),
          type: "PAIEMENT",
          titre: `Paiement reçu : ${fcfa(p.montant)} FCFA`,
          description: `Règlement ${p.mode_reglement} pour ${p.candidat_nom_complet || currentCandName} • Reçu n° ${p.reference_recu || p.id.slice(0, 8)}`,
          acteur: "Guichet",
          badgeLabel: p.mode_reglement,
          badgeClass:
            "bg-emerald-500/15 text-emerald-700 dark:text-emerald-400 border-emerald-500/20",
          montant: p.montant,
        });
      });
    }

    drivingSessions.value.forEach((sc) => {
      if (
        currentRole.value === "eleve" &&
        currentCandDossierId &&
        sc.dossier_id !== currentCandDossierId &&
        !sc.candidat_nom?.toLowerCase().includes(currentCandName.toLowerCase())
      ) {
        return;
      }
      if (
        currentRole.value === "moniteur" &&
        currentUserId &&
        sc.moniteur_id &&
        sc.moniteur_id !== currentUserId
      ) {
        return;
      }

      const d = new Date(sc.date_seance || Date.now());
      const isValidee = sc.statut === "EFFECTUE";
      const compCount = Array.isArray(sc.competences)
        ? sc.competences.length
        : 0;

      items.push({
        id: `conduite-${sc.id}`,
        timestamp: d.getTime() || Date.now(),
        dateStr: d.toLocaleDateString("fr-FR"),
        heureStr: `${sc.heure_debut} - ${sc.heure_fin}`,
        type: "CONDUITE",
        titre: `Leçon pratique : ${sc.candidat_nom}`,
        description: `${sc.vehicule_immatriculation ? "Véhicule " + sc.vehicule_immatriculation + " • " : ""}Moniteur : ${sc.moniteur_nom || "Non assigné"}${compCount > 0 ? " • " + compCount + " compétences évaluées" : ""}`,
        acteur: sc.moniteur_nom || "Moniteur",
        badgeLabel: isValidee ? "Validée (Émargée)" : sc.statut,
        badgeClass: isValidee
          ? "bg-blue-500/15 text-blue-700 dark:text-blue-400 border-blue-500/20"
          : "bg-amber-500/15 text-amber-700 dark:text-amber-400 border-amber-500/20",
      });
    });

    coursStore.seances.forEach((s) => {
      const d = new Date(s.date_seance || Date.now());
      const isTerminee = s.statut_seance === "TERMINEE";
      const isOnline = s.type_seance === "EN_LIGNE";

      items.push({
        id: `cours-${s.id}`,
        timestamp: d.getTime() || Date.now(),
        dateStr: d.toLocaleDateString("fr-FR"),
        heureStr: `${s.heure_debut} - ${s.heure_fin}`,
        type: "CODE",
        titre: `${isOnline ? "Classe Virtuelle Visio" : "Cours Code en Salle"} : ${s.theme_traite}`,
        description: `Formateur : ${s.formateur_nom} • Inscrits : ${s.nombre_inscrits || 0} • Présents : ${s.nombre_presents || 0}`,
        acteur: s.formateur_nom,
        badgeLabel: isTerminee ? "Terminé (+2h code)" : s.statut_seance,
        badgeClass: isTerminee
          ? "bg-indigo-500/15 text-indigo-700 dark:text-indigo-400 border-indigo-500/20"
          : "bg-violet-500/15 text-violet-700 dark:text-violet-400 border-violet-500/20",
      });
    });

    quizStore.soumissions.forEach((sub) => {
      if (
        currentRole.value === "eleve" &&
        currentCandName &&
        !sub.candidat_nom.toLowerCase().includes(currentCandName.toLowerCase())
      ) {
        return;
      }

      const d = new Date(sub.submitted_at || Date.now());
      const reussite = sub.pourcentage >= 75;

      items.push({
        id: `quiz-sub-${sub.id}`,
        timestamp: d.getTime() || Date.now(),
        dateStr: d.toLocaleDateString("fr-FR"),
        heureStr: d.toLocaleTimeString("fr-FR", {
          hour: "2-digit",
          minute: "2-digit",
        }),
        type: "QUIZ",
        titre: `Évaluation Quiz : ${sub.quiz_titre}`,
        description: `Apprenant : ${sub.candidat_nom} • Score : ${sub.score_obtenu}/${sub.total_points} (${sub.pourcentage}%)`,
        acteur: sub.candidat_nom,
        badgeLabel: reussite ? "Validé (≥75%)" : "À réviser",
        badgeClass: reussite
          ? "bg-emerald-500/15 text-emerald-700 dark:text-emerald-400 border-emerald-500/20"
          : "bg-rose-500/15 text-rose-700 dark:text-rose-400 border-rose-500/20",
        score: sub.pourcentage,
      });
    });

    if (currentRole.value === "manager" || currentRole.value === "directeur") {
      const derogations =
        (guichetStore.gerantDashboardData as any)?.derogations_accordees || [];
      derogations.forEach((derog: any, idx: number) => {
        const d = new Date(derog.date_derogation || Date.now());
        items.push({
          id: `derogation-${derog.id || idx}`,
          timestamp: d.getTime() || Date.now(),
          dateStr: d.toLocaleDateString("fr-FR"),
          heureStr: d.toLocaleTimeString("fr-FR", {
            hour: "2-digit",
            minute: "2-digit",
          }),
          type: "DEROGATION",
          titre: `Dérogation PIN Manager accordée`,
          description: `Bénéficiaire : ${derog.candidat_nom || "Candidat"} • Raison : ${derog.motif || "Autorisation spéciale de passage"}`,
          acteur: derog.manager_nom || "Directeur",
          badgeLabel: "Code PIN Vérifié",
          badgeClass:
            "bg-amber-500/15 text-amber-700 dark:text-amber-400 border-amber-500/20",
        });
      });

      guichetStore.vehicules.forEach((veh) => {
        if (veh.statut !== "DISPONIBLE") {
          const d = new Date(veh.created_at || Date.now());
          items.push({
            id: `flotte-${veh.id}`,
            timestamp: d.getTime() || Date.now(),
            dateStr: d.toLocaleDateString("fr-FR"),
            heureStr: "Statut véhicule",
            type: "FLOTTE",
            titre: `Alerte Flotte : ${veh.immatriculation} (${veh.marque_modele})`,
            description: `Statut technique actuel : ${veh.statut} • Kilométrage compteur : ${veh.kilometrage_actuel || 0} km`,
            acteur: "Parc Automobile",
            badgeLabel: veh.statut,
            badgeClass:
              "bg-orange-500/15 text-orange-700 dark:text-orange-400 border-orange-500/20",
          });
        }
      });
    }

    if (
      currentRole.value === "secretaire" ||
      currentRole.value === "manager" ||
      currentRole.value === "caissier"
    ) {
      if (guichetStore.stats.total_especes_theorique > 0) {
        items.push({
          id: "caisse-session-active",
          timestamp: Date.now() - 3600000,
          dateStr: "Aujourd'hui",
          heureStr: "En cours",
          type: "CAISSE",
          titre: "Session de Caisse Active (Guichet)",
          description: `Tiroir théorique : ${fcfa(guichetStore.stats.total_especes_theorique)} FCFA • Fond de caisse initial : ${fcfa(guichetStore.stats.fond_de_caisse)} FCFA`,
          acteur: authStore.currentUser?.nom || "Caissier",
          badgeLabel: "Session Ouverte",
          badgeClass:
            "bg-emerald-500/15 text-emerald-700 dark:text-emerald-400 border-emerald-500/20",
        });
      }
    }

    return items;
  });

  const allHistoryItems = computed<HistoriqueItem[]>(() => {
    const merged = [...liveEvents.value, ...baseHistoryItems.value];
    const map = new Map<string, HistoriqueItem>();
    merged.forEach((item) => {
      if (!map.has(item.id)) {
        map.set(item.id, item);
      }
    });
    return Array.from(map.values()).sort((a, b) => b.timestamp - a.timestamp);
  });

  const kpiTotalActions = computed(() => allHistoryItems.value.length);
  const kpiTotalMontant = computed(() =>
    allHistoryItems.value
      .filter((i) => i.type === "PAIEMENT" && i.montant)
      .reduce((acc, curr) => acc + (curr.montant || 0), 0),
  );
  const kpiTotalConduite = computed(
    () => allHistoryItems.value.filter((i) => i.type === "CONDUITE").length,
  );
  const kpiTotalCode = computed(
    () => allHistoryItems.value.filter((i) => i.type === "CODE").length,
  );

  onMounted(async () => {
    realtimeSocket.connect();
    await refreshAllData();

    unsubscribeWs = realtimeSocket.subscribe(async (msg) => {
      lastSyncAt.value = new Date().toLocaleTimeString("fr-FR");
      const now = Date.now();
      const dStr = new Date().toLocaleDateString("fr-FR");
      const hStr = new Date().toLocaleTimeString("fr-FR", {
        hour: "2-digit",
        minute: "2-digit",
      });

      switch (msg.type) {
        case "PAIEMENT_ENREGISTRE": {
          const montant = msg.data?.montant || 0;
          const candNom = msg.data?.candidat_nom || "Candidat";
          const mode = msg.data?.mode_reglement || "Espèces";
          const recuRef = msg.data?.reference_recu || "REC-LIVE";

          liveEvents.value.unshift({
            id: `live-paiement-${now}-${Math.random()}`,
            timestamp: now,
            dateStr: dStr,
            heureStr: hStr,
            type: "PAIEMENT",
            titre: `Nouveau paiement : ${fcfa(montant)} FCFA`,
            description: `Règlement ${mode} pour ${candNom} • Reçu n° ${recuRef}`,
            acteur: "Guichet",
            badgeLabel: "Reçu en direct",
            badgeClass:
              "bg-emerald-500 text-white font-bold animate-pulse shadow-sm",
            montant,
            isRealtimeLive: true,
          });
          await guichetStore.loadStats();
          break;
        }
        case "LIVRET_EMARGEMENT_UPDATED":
        case "PLANNING_CONDUITE_UPDATED": {
          const cNom = msg.data?.candidat_nom || "Candidat";
          liveEvents.value.unshift({
            id: `live-conduite-${now}-${Math.random()}`,
            timestamp: now,
            dateStr: dStr,
            heureStr: hStr,
            type: "CONDUITE",
            titre: `Émargement conduite validé : ${cNom}`,
            description: `Validation tactile effectuée • Compétences et carnet de bord synchronisés`,
            acteur: "Moniteur",
            badgeLabel: "Émargé en direct",
            badgeClass:
              "bg-blue-500 text-white font-bold animate-pulse shadow-sm",
            isRealtimeLive: true,
          });
          try {
            const pl = await api.getPlanningConduite("ALL");
            drivingSessions.value = pl || [];
          } catch {}
          break;
        }
        case "COURS_CODE_UPDATED": {
          liveEvents.value.unshift({
            id: `live-code-${now}-${Math.random()}`,
            timestamp: now,
            dateStr: dStr,
            heureStr: hStr,
            type: "CODE",
            titre: `Séance de Code Théorique actualisée`,
            description: `Présences et validations d'heures de code synchronisées`,
            acteur: "Pédagogie",
            badgeLabel: "Actualisé en direct",
            badgeClass:
              "bg-indigo-500 text-white font-bold animate-pulse shadow-sm",
            isRealtimeLive: true,
          });
          await coursStore.loadSeances();
          break;
        }
        case "CLASSE_VIRTUELLE_JOIN": {
          const participant = msg.data?.nom || "Élève";
          liveEvents.value.unshift({
            id: `live-visio-${now}-${Math.random()}`,
            timestamp: now,
            dateStr: dStr,
            heureStr: hStr,
            type: "CODE",
            titre: `Connexion classe virtuelle Jitsi`,
            description: `${participant} a rejoint le salon de visioconférence en direct`,
            acteur: participant,
            badgeLabel: "Visio active",
            badgeClass:
              "bg-purple-500 text-white font-bold animate-pulse shadow-sm",
            isRealtimeLive: true,
          });
          break;
        }
        case "QUIZ_RESULTAT_UPDATED":
        case "QUIZ_SUBMISSION_RECEIVED": {
          const cNom = msg.data?.candidat_nom || "Apprenant";
          const score = msg.data?.pourcentage ?? msg.data?.score ?? 0;
          liveEvents.value.unshift({
            id: `live-quiz-${now}-${Math.random()}`,
            timestamp: now,
            dateStr: dStr,
            heureStr: hStr,
            type: "QUIZ",
            titre: `Résultat Quiz de Code reçu : ${cNom}`,
            description: `Score enregistré : ${score}% (${score >= 75 ? "Admis" : "À retravailler"})`,
            acteur: cNom,
            badgeLabel: `Score ${score}%`,
            badgeClass:
              score >= 75
                ? "bg-emerald-600 text-white font-bold shadow-sm"
                : "bg-rose-600 text-white font-bold shadow-sm",
            score,
            isRealtimeLive: true,
          });
          await quizStore.loadSoumissions();
          break;
        }
        case "CAISSE_SESSION_UPDATED": {
          liveEvents.value.unshift({
            id: `live-caisse-${now}-${Math.random()}`,
            timestamp: now,
            dateStr: dStr,
            heureStr: hStr,
            type: "CAISSE",
            titre: "Opération de Caisse / Pointage Z",
            description: `Mise à jour des fonds et pointage journalier`,
            acteur: "Guichet",
            badgeLabel: "Caisse Z",
            badgeClass: "bg-brand-dark text-white font-bold shadow-sm",
            isRealtimeLive: true,
          });
          await guichetStore.loadStats();
          break;
        }
        case "FLOTTE_VEHICULES_UPDATED": {
          liveEvents.value.unshift({
            id: `live-flotte-${now}-${Math.random()}`,
            timestamp: now,
            dateStr: dStr,
            heureStr: hStr,
            type: "FLOTTE",
            titre: "Flotte Automobile Mise à Jour",
            description: `Statuts de disponibilité ou fiches véhicules actualisés`,
            acteur: "Parc",
            badgeLabel: "Flotte",
            badgeClass: "bg-amber-600 text-white font-bold shadow-sm",
            isRealtimeLive: true,
          });
          await guichetStore.loadVehicules();
          break;
        }
        case "DOSSIER_ARCHIVE_CREATED": {
          const cNom = msg.data?.candidat_nom || "Candidat";
          const dosNum = msg.data?.numero_dossier || "DOS";
          liveEvents.value.unshift({
            id: `live-archive-${now}-${Math.random()}`,
            timestamp: now,
            dateStr: dStr,
            heureStr: hStr,
            type: "EXAMEN",
            titre: `Dossier scellé aux archives : ${cNom}`,
            description: `Dossier N° ${dosNum} transféré au Centre d'Archivage avec pièces et reçus`,
            acteur: "Archivage",
            badgeLabel: "Archivé",
            badgeClass: "bg-slate-700 text-white font-bold shadow-sm",
            isRealtimeLive: true,
          });
          break;
        }
      }
    });
  });

  onUnmounted(() => {
    if (unsubscribeWs) {
      unsubscribeWs();
      unsubscribeWs = null;
    }
  });

  return {
    isWsConnected,
    lastSyncAt,
    isRefreshing,
    currentRole,
    roleTitle,
    quizCount: computed(() => quizStore.soumissions.length),
    allHistoryItems,
    kpiTotalActions,
    kpiTotalMontant,
    kpiTotalConduite,
    kpiTotalCode,
    refreshAllData,
  };
}

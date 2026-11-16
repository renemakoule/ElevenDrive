import { defineStore } from "pinia";
import { api } from "../db/api";
import { neonSyncEngine } from "../services/neonSyncService";
import { realtimeSocket } from "../services/realtimeSocketService";
import { useAuthStore } from "./authStore";
import {
  execCloturerCaisse,
  execEnregistrerSeanceCode,
  execLoadGerantDashboard,
  execLoadSessionCaisse,
  execReprintReceipt,
  execUnlockConduiteWithPin,
} from "./guichet/caisseActions";
import {
  execAjouterVehicule,
  execCreerFormule,
  execLoadFormules,
  execLoadVehicules,
  execModifierPrixFormule,
  execModifierStatutVehicule,
  execSupprimerFormule,
  execViderCatalogueFormules,
} from "./guichet/flotteFormulesActions";
import { execImportCandidatesBatch } from "./guichet/importActions";
import type {
  Candidat,
  DossierDetailComplet,
  FormulePermis,
  TranchePaiementConfig,
  GuichetStats,
  Vehicule,
  CreerVehiculeInput,
  TableauGerantData,
  StatutVehicule,
  ImporterCandidatsBatchInput,
} from "../types";

export const useGuichetStore = defineStore("guichet", {
  state: () => ({
    searchQuery: "",
    searchResults: [] as Candidat[],
    selectedCandidatId: null as string | null,
    activeDossierDetail: null as DossierDetailComplet | null,
    formules: [] as FormulePermis[],
    stats: {
      total_candidats: 0,
      fond_de_caisse: 0,
      total_especes_theorique: 0,
      evenements_sync_en_attente: 0,
      mode_connectivite: "LOCAL_OFFLINE_READY",
    } as GuichetStats,

    activeView: "guichet" as
      | "guichet"
      | "code"
      | "caisse_z"
      | "gerant"
      | "moniteur"
      | "candidat"
      | "examens",
    isNewCandidateModalOpen: false,
    isImportCandidatesModalOpen: false,
    isReceiptModalOpen: false,
    isCashCloseModalOpen: false,
    isRapportZModalOpen: false,
    isManagerOverrideModalOpen: false,
    sessionCaisseDetail: null as any,
    derogationsAccorde: [] as string[],
    gerantDashboardData: null as TableauGerantData | null,
    vehicules: [] as Vehicule[],
    flotteLoading: false,

    lastReceipt: null as {
      numero_recu: string;
      code_securite_hash: string;
      candidat_nom: string;
      numero_dossier: string;
      formule_nom: string;
      montant_verse: number;
      solde_restant: number;
      mode_reglement: string;
      date_paiement: string;
      caissier_nom: string;
      est_duplicata?: boolean;
      nombre_impressions?: number;
    } | null,

    toastMessage: "" as string,
    toastType: "info" as "info" | "success" | "error",
    loading: false,
  }),

  actions: {
    async init() {
      const authStore = useAuthStore();
      if (!authStore.isAuthenticated || this.loading) {
        return;
      }
      this.loading = true;
      try {
        const [statsData, formulesData] = await Promise.all([
          api.getGuichetStats(),
          api.getFormules(),
        ]);
        this.stats = statsData || {
          total_candidats: 0,
          fond_de_caisse: 0,
          total_especes_theorique: 0,
          evenements_sync_en_attente: 0,
          mode_connectivite: "LOCAL_OFFLINE_READY",
        };
        this.formules = Array.isArray(formulesData) ? formulesData : [];
        neonSyncEngine.setPendingCount(
          this.stats?.evenements_sync_en_attente || 0,
        );

        await this.search("");
        if (this.searchResults && this.searchResults.length > 0) {
          await this.selectCandidat(this.searchResults[0].id);
        }

        // Synchronisation bidirectionnelle initiale non-bloquante (Push + Pull Cloud -> SQLite)
        if (neonSyncEngine.getState().isOnline) {
          neonSyncEngine.triggerAutoSync().catch(() => {});
        }

        // Pull automatique en temps réel quand un autre poste émet une mutation via WebSocket
        let wsSyncTimer: ReturnType<typeof setTimeout> | null = null;
        realtimeSocket.subscribe((msg) => {
          if (msg.type === "PONG" || !neonSyncEngine.getState().isOnline) return;
          if (wsSyncTimer) clearTimeout(wsSyncTimer);
          wsSyncTimer = setTimeout(() => {
            if (!neonSyncEngine.getState().isSyncing) {
              neonSyncEngine.triggerAutoSync().catch(() => {});
            }
          }, 1500);
        });
      } catch (err: any) {
        this.showToast(
          "Erreur initialisation guichet: " + err.message,
          "error",
        );
      } finally {
        this.loading = false;
      }
    },

    async search(query: string) {
      this.searchQuery = query;
      try {
        const results = await api.searchCandidats(query);
        this.searchResults = Array.isArray(results) ? results : [];
      } catch (err: any) {
        this.searchResults = [];
        this.showToast("Erreur recherche: " + err.message, "error");
      }
    },

    async loadStats() {
      try {
        const statsData = await api.getGuichetStats();
        this.stats = statsData;
        neonSyncEngine.setPendingCount(this.stats.evenements_sync_en_attente);
      } catch {
        // Ignorer en cas de micro-coupure réseau
      }
    },

    async selectCandidat(id: string) {
      this.selectedCandidatId = id;
      this.loading = true;
      try {
        this.activeDossierDetail = await api.getCandidatDossier(id);
      } catch (err: any) {
        this.showToast("Erreur chargement dossier: " + err.message, "error");
      } finally {
        this.loading = false;
      }
    },

    async processPayment(params: {
      montant: number;
      mode_reglement: string;
      reference_transaction?: string;
    }) {
      if (!this.activeDossierDetail) return;
      this.loading = true;
      try {
        const res = await api.enregistrerPaiement({
          dossier_id: this.activeDossierDetail.dossier.id,
          montant: params.montant,
          mode_reglement: params.mode_reglement,
          reference_transaction: params.reference_transaction,
        });

        await this.selectCandidat(this.activeDossierDetail.candidat.id);
        this.stats = await api.getGuichetStats();
        neonSyncEngine.setPendingCount(this.stats.evenements_sync_en_attente);
        this.loadSessionCaisse().catch(() => {});

        const authStore = useAuthStore();
        const caissierNom = authStore.currentUser?.nom
          ? `${authStore.currentUser.nom} ${authStore.currentUser.prenom} (${authStore.currentUser.roleLabel || "Guichet"})`
          : "Agent Guichet";

        this.lastReceipt = {
          numero_recu: res.numero_recu,
          code_securite_hash:
            res.code_securite_hash ||
            (res.numero_recu ? `SECHASH-${res.numero_recu}` : "SECHASH-000"),
          candidat_nom: `${this.activeDossierDetail.candidat.nom} ${this.activeDossierDetail.candidat.prenoms}`,
          numero_dossier: res.numero_dossier,
          formule_nom: this.activeDossierDetail.formule.intitule,
          montant_verse: res.montant_verse,
          solde_restant: res.nouveau_solde_restant,
          mode_reglement: params.mode_reglement,
          date_paiement: res.date_paiement,
          caissier_nom: caissierNom,
        };

        this.isReceiptModalOpen = true;
        this.showToast(
          `Encaissement de ${params.montant.toLocaleString("fr-FR")} FCFA validé ! Ticket généré.`,
          "success",
        );
      } catch (err: any) {
        this.showToast(err.message || "Erreur encaissement", "error");
      } finally {
        this.loading = false;
      }
    },

    async createCandidate(payload: {
      nom: string;
      prenoms: string;
      genre: "M" | "F";
      date_naissance: string;
      lieu_naissance: string;
      telephone_whatsapp: string;
      type_piece_id: string;
      numero_piece_id: string;
      formule_id: string;
      acompte_initial: number;
      mode_reglement: string;
      photo_url?: string;
      visite_medicale_statut?:
        | "VALIDE_SANS_RESERVE"
        | "VALIDE_AVEC_VERRES"
        | "EN_ATTENTE";
      mode_formation_theorique?: "PRESENTIEL" | "EN_LIGNE" | "HYBRIDE";
    }) {
      this.loading = true;
      try {
        const res = await api.creerCandidatAvecDossier(payload);
        await this.search("");
        await this.selectCandidat(res.candidat.id);
        this.isNewCandidateModalOpen = false;
        this.stats = await api.getGuichetStats();
        neonSyncEngine.setPendingCount(this.stats.evenements_sync_en_attente);
        this.loadSessionCaisse().catch(() => {});

        if (res.dossier_detail.paiements.length > 0) {
          const firstPay = res.dossier_detail.paiements[0];
          const authStore = useAuthStore();
          const caissierNom = authStore.currentUser?.nom
            ? `${authStore.currentUser.nom} ${authStore.currentUser.prenom} (${authStore.currentUser.roleLabel || "Guichet"})`
            : "Agent Guichet";

          this.lastReceipt = {
            numero_recu: firstPay.reference_recu,
            code_securite_hash: firstPay.reference_recu
              ? `SECHASH-${firstPay.reference_recu}`
              : "SECHASH-INIT",
            candidat_nom: `${res.candidat.nom} ${res.candidat.prenoms}`,
            numero_dossier: res.dossier_detail.dossier.numero_dossier,
            formule_nom: res.dossier_detail.formule.intitule,
            montant_verse: firstPay.montant,
            solde_restant: res.dossier_detail.dossier.solde_restant,
            mode_reglement: firstPay.mode_reglement,
            date_paiement: firstPay.date_paiement,
            caissier_nom: caissierNom,
          };
          this.isReceiptModalOpen = true;
        }

        this.showToast(
          `Candidat ${res.candidat.nom} inscrit avec succès ! Dossier: ${res.dossier_detail.dossier.numero_dossier}`,
          "success",
        );
      } catch (err: any) {
        this.showToast(err.message || "Erreur création candidat", "error");
      } finally {
        this.loading = false;
      }
    },

    async updateCandidatePhoto(candidatId: string, photoUrl: string) {
      this.loading = true;
      try {
        await api.updateCandidatPhoto(candidatId, photoUrl);
        if (
          this.activeDossierDetail &&
          this.activeDossierDetail.candidat.id === candidatId
        ) {
          this.activeDossierDetail.candidat.photo_url = photoUrl;
        }
        this.showToast(
          "Photo 4x4 d'identité mise à jour avec succès !",
          "success",
        );
      } catch (err: any) {
        this.showToast("Erreur mise à jour photo: " + err.message, "error");
      } finally {
        this.loading = false;
      }
    },

    loadSessionCaisse() {
      return execLoadSessionCaisse(this);
    },

    cloturerCaisse(payload: {
      total_especes_compte: number;
      observations?: string;
    }) {
      return execCloturerCaisse(this, payload);
    },

    reprintReceipt(paiementId: string) {
      return execReprintReceipt(this, paiementId);
    },

    unlockConduiteWithPin(pin: string, dossierId: string, motif?: string) {
      return execUnlockConduiteWithPin(this, pin, dossierId, motif);
    },

    loadGerantDashboard() {
      return execLoadGerantDashboard(this);
    },

    enregistrerSeanceCode(payload: {
      theme_traite: string;
      formateur_nom: string;
      heure_debut: string;
      heure_fin: string;
      presences: Array<{
        candidat_id: string;
        est_present: boolean;
        score_test: number;
      }>;
    }) {
      return execEnregistrerSeanceCode(this, payload);
    },

    showToast(message: string, type: "info" | "success" | "error" = "info") {
      this.toastMessage = message;
      this.toastType = type;
      setTimeout(() => {
        if (this.toastMessage === message) {
          this.toastMessage = "";
        }
      }, 4000);
    },

    triggerToast(message: string, type: "info" | "success" | "error" = "info") {
      this.showToast(message, type);
    },

    loadFormules() {
      return execLoadFormules(this);
    },

    creerFormule(input: {
      categorie_permis: string;
      intitule: string;
      description?: string;
      heures_theoriques: number;
      heures_pratiques: number;
      prix_total: number;
      acompte_minimum: number;
      tranches_paiement?: TranchePaiementConfig[];
    }) {
      return execCreerFormule(this, input);
    },

    modifierPrixFormule(
      formuleId: string,
      nouveauPrix: number,
      acompteMinimum?: number,
      motif?: string,
      tranchesPaiement?: TranchePaiementConfig[],
    ) {
      return execModifierPrixFormule(
        this,
        formuleId,
        nouveauPrix,
        acompteMinimum,
        motif,
        tranchesPaiement,
      );
    },

    supprimerFormule(formuleId: string) {
      return execSupprimerFormule(this, formuleId);
    },

    viderCatalogueFormules() {
      return execViderCatalogueFormules(this);
    },

    loadFlotte() {
      return this.loadVehicules();
    },

    loadVehicules() {
      return execLoadVehicules(this);
    },

    ajouterVehicule(input: CreerVehiculeInput) {
      return execAjouterVehicule(this, input);
    },

    modifierStatutVehicule(vehiculeId: string, statut: StatutVehicule) {
      return execModifierStatutVehicule(this, vehiculeId, statut);
    },

    importCandidatesBatch(payload: ImporterCandidatsBatchInput) {
      return execImportCandidatesBatch(this, payload);
    },
  },
});

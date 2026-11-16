import { defineStore } from "pinia";
import type {
  DossierArchiveItem,
  DocumentArchive,
  MotifArchivage,
  TypeDocumentImportant,
  Paiement,
} from "../types";
import { useAuthStore } from "./authStore";
import { useGuichetStore } from "./guichetStore";
import { realtimeSocket } from "../services/realtimeSocketService";
import { API_BASE_URL, R2_PUBLIC_URL } from "../config/env";

const STORAGE_KEY_ARCHIVES = "autoschool_dossiers_archives";
const STORAGE_KEY_DOCS = "autoschool_documents_importants";

const getAuthHeaders = (): Record<string, string> => {
  const headers: Record<string, string> = { "Content-Type": "application/json" };
  if (typeof localStorage !== "undefined") {
    const raw =
      localStorage.getItem("autoschool_active_session") ||
      localStorage.getItem("autoschool_web_session") ||
      localStorage.getItem("elevendrive_active_session");
    if (raw) {
      try {
        const sess = JSON.parse(raw);
        if (sess.token) headers["Authorization"] = `Bearer ${sess.token}`;
      } catch {}
    }
  }
  return headers;
};

export const useArchivageStore = defineStore("archivage", {
  state: () => ({
    archives: [] as DossierArchiveItem[],
    documentsImportants: [] as DocumentArchive[],
    loading: false,
    selectedArchiveId: "" as string,
    filterMotif: "ALL" as string,
    filterAnnee: "ALL" as string,
    searchQuery: "" as string,
    cloudR2Statut: {
      isConfigured: true,
      bucketName: (import.meta.env?.VITE_R2_BUCKET_NAME as string) || "",
      publicURL: R2_PUBLIC_URL,
      provider: "Cloudflare R2",
    },
  }),

  getters: {
    selectedArchive(state): DossierArchiveItem | null {
      if (!state.selectedArchiveId) return null;
      return (
        state.archives.find((a) => a.id === state.selectedArchiveId) || null
      );
    },

    filteredArchives(state): DossierArchiveItem[] {
      return state.archives.filter((item) => {
        // Filtre motif
        if (
          state.filterMotif !== "ALL" &&
          item.motif_archivage !== state.filterMotif
        ) {
          return false;
        }
        // Filtre année
        if (state.filterAnnee !== "ALL") {
          const year = new Date(item.date_archivage).getFullYear().toString();
          if (year !== state.filterAnnee) return false;
        }
        // Recherche
        if (state.searchQuery.trim()) {
          const q = state.searchQuery.toLowerCase().trim();
          const matchNom = `${item.candidat_nom} ${item.candidat_prenoms}`
            .toLowerCase()
            .includes(q);
          const matchNum = item.numero_dossier.toLowerCase().includes(q);
          const matchBoite = (item.reference_boite_physique || "")
            .toLowerCase()
            .includes(q);
          const matchFormule = item.formule_intitule.toLowerCase().includes(q);
          return matchNom || matchNum || matchBoite || matchFormule;
        }
        return true;
      });
    },

    anneesDisponibles(state): string[] {
      const years = new Set<string>();
      state.archives.forEach((a) => {
        const d = new Date(a.date_archivage);
        if (!isNaN(d.getFullYear())) {
          years.add(d.getFullYear().toString());
        }
      });
      return Array.from(years).sort((a, b) => b.localeCompare(a));
    },

    totalArchivesCount(state): number {
      return state.archives.length;
    },

    totalDocumentsArchived(state): number {
      return state.archives.reduce(
        (acc, a) => acc + (a.documents_count || 0),
        0,
      );
    },

    totalFondsArchives(state): number {
      return state.archives.reduce((acc, a) => acc + (a.total_paye || 0), 0);
    },
  },

  actions: {
    /**
     * Initialisation et chargement des archives
     */
    async loadArchives() {
      this.loading = true;
      try {
        // Tentative de récupération depuis le backend API
        const authStore = useAuthStore();
        const schoolId =
          authStore.currentSchool?.id ||
          authStore.currentUser?.auto_ecole_id ||
          "";

        // Récupération locale dans le cache
        const rawLocal = localStorage.getItem(STORAGE_KEY_ARCHIVES);
        let localItems: DossierArchiveItem[] = [];
        if (rawLocal) {
          try {
            localItems = JSON.parse(rawLocal);
          } catch {
            localItems = [];
          }
        }

        const rawDocs = localStorage.getItem(STORAGE_KEY_DOCS);
        if (rawDocs) {
          try {
            this.documentsImportants = JSON.parse(rawDocs);
          } catch {
            this.documentsImportants = [];
          }
        }

        // Si l'école est connectée, filtrer par auto-école
        if (schoolId) {
          this.archives = localItems.filter(
            (a) => !a.auto_ecole_id || a.auto_ecole_id === schoolId,
          );
        } else {
          this.archives = localItems;
        }

        // Vérification de la connectivité Cloudflare R2
        try {
          const res = await fetch(`${API_BASE_URL}/archivage/cloud/statut`, {
            headers: getAuthHeaders(),
          });
          if (res.ok) {
            const data = await res.json();
            if (data.success && data.data) {
              this.cloudR2Statut = {
                isConfigured: data.data.is_configured,
                bucketName: data.data.bucket_name || "",
                publicURL: data.data.public_url || R2_PUBLIC_URL,
                provider: "Cloudflare R2",
              };
            }
          }
        } catch {
          // Mode hors-ligne / fallback local
        }
      } finally {
        this.loading = false;
      }
    },

    /**
     * Sauvegarde locale
     */
    persistLocal() {
      try {
        localStorage.setItem(
          STORAGE_KEY_ARCHIVES,
          JSON.stringify(this.archives),
        );
        localStorage.setItem(
          STORAGE_KEY_DOCS,
          JSON.stringify(this.documentsImportants),
        );
      } catch {
        // ignore
      }
    },

    /**
     * Archive un dossier de candidat complet avec ses paiements et documents rattachés
     */
    async archiverDossier(payload: {
      candidatId: string;
      dossierId: string;
      motif: MotifArchivage;
      referenceBoite?: string;
      observations?: string;
      documentsFinaux?: {
        type: TypeDocumentImportant;
        titre: string;
        nomFichier: string;
        base64: string;
      }[];
    }) {
      const guichetStore = useGuichetStore();
      const authStore = useAuthStore();

      // Récupérer le dossier détaillé
      let dossierDetail = guichetStore.activeDossierDetail;
      if (!dossierDetail || dossierDetail.candidat.id !== payload.candidatId) {
        await guichetStore.selectCandidat(payload.candidatId);
        dossierDetail = guichetStore.activeDossierDetail;
      }

      if (!dossierDetail) {
        throw new Error(
          "Impossible de récupérer les détails du dossier à archiver.",
        );
      }

      const nowIso = new Date().toISOString();
      const archiveId = `ARCH-${Date.now()}-${Math.random().toString(36).slice(2, 7)}`;

      // Récupérer les documents existants liés à ce candidat
      const candDocs = this.documentsImportants.filter(
        (d) =>
          d.candidat_id === payload.candidatId ||
          d.dossier_id === payload.dossierId,
      );

      // Ajouter les documents finaux éventuels
      const addedDocs: DocumentArchive[] = (payload.documentsFinaux || []).map(
        (doc, idx) => ({
          id: `DOC-ARCH-${Date.now()}-${idx}`,
          dossier_id: payload.dossierId,
          candidat_id: payload.candidatId,
          type_document: doc.type,
          titre: doc.titre,
          nom_fichier: doc.nomFichier,
          fichier_url_ou_base64: doc.base64,
          date_televersement: nowIso,
          televerse_par_nom:
            `${authStore.currentUser?.nom || ""} ${authStore.currentUser?.prenom || ""}`.trim() ||
            "Agent",
          date_archivage: nowIso,
          est_archive: true,
        }),
      );

      // Si le candidat a une photo d'identité, l'ajouter comme document important
      if (dossierDetail.candidat.photo_url) {
        const photoExiste = candDocs.some(
          (d) => d.type_document === "PHOTO_IDENTITE",
        );
        if (!photoExiste) {
          addedDocs.push({
            id: `DOC-PHOTO-${Date.now()}`,
            dossier_id: payload.dossierId,
            candidat_id: payload.candidatId,
            type_document: "PHOTO_IDENTITE",
            titre: "Photo d'identité officielle (Webcam/Upload)",
            nom_fichier: "photo_identite.jpg",
            fichier_url_ou_base64: dossierDetail.candidat.photo_url,
            date_televersement: nowIso,
            televerse_par_nom: "Système Guichet",
            date_archivage: nowIso,
            est_archive: true,
          });
        }
      }

      const allDocsArchived = [...candDocs, ...addedDocs].map((d) => ({
        ...d,
        est_archive: true,
        date_archivage: d.date_archivage || nowIso,
      }));

      const dejaPaye =
        dossierDetail.dossier.montant_total -
        dossierDetail.dossier.solde_restant;

      const archiveItem: DossierArchiveItem = {
        id: archiveId,
        auto_ecole_id:
          authStore.currentSchool?.id ||
          authStore.currentUser?.auto_ecole_id ||
          dossierDetail.candidat.auto_ecole_id,
        dossier_id: payload.dossierId,
        candidat_id: payload.candidatId,
        numero_dossier: dossierDetail.dossier.numero_dossier,
        candidat_nom: dossierDetail.candidat.nom,
        candidat_prenoms: dossierDetail.candidat.prenoms,
        telephone: dossierDetail.candidat.telephone_whatsapp,
        categorie_permis:
          dossierDetail.formule.categorie_permis ||
          dossierDetail.candidat.categorie_permis ||
          "B",
        formule_intitule: dossierDetail.formule.intitule,
        montant_total: dossierDetail.dossier.montant_total,
        total_paye: Math.max(0, dejaPaye),
        solde_restant: dossierDetail.dossier.solde_restant,
        heures_code_effectuees: dossierDetail.dossier.heures_code_effectuees,
        heures_conduite_effectuees:
          dossierDetail.dossier.heures_conduite_effectuees,
        statut_dossier_origine: dossierDetail.dossier.statut_dossier,
        date_inscription: dossierDetail.dossier.date_inscription,
        date_archivage: nowIso,
        archive_par_id: authStore.currentUser?.id || "",
        archive_par_nom:
          `${authStore.currentUser?.nom || ""} ${authStore.currentUser?.prenom || ""}`.trim() ||
          "Directeur / Secrétaire",
        archive_par_role: authStore.currentRole || "Staff",
        motif_archivage: payload.motif,
        reference_boite_physique:
          payload.referenceBoite ||
          `BOITE-${new Date().getFullYear()}-${dossierDetail.dossier.numero_dossier}`,
        observations: payload.observations || "",
        documents_count: allDocsArchived.length,
        paiements_count: (dossierDetail.paiements || []).length,
        documents: allDocsArchived,
        paiements: dossierDetail.paiements || [],
      };

      // Mettre à jour le statut du dossier d'origine
      dossierDetail.dossier.statut_dossier = "ARCHIVE";

      // Ajouter en tête de liste
      this.archives.unshift(archiveItem);

      // Mettre à jour la collection des documents
      addedDocs.forEach((doc) => {
        this.documentsImportants.unshift(doc);
      });

      this.persistLocal();

      // Scellement asynchrone dans le Cloud Cloudflare R2 (Bucket : maryam-logistic-bucket)
      try {
        fetch(`${API_BASE_URL}/archivage/cloud/dossier`, {
          method: "POST",
          headers: getAuthHeaders(),
          body: JSON.stringify(archiveItem),
        })
          .then((res) => res.json())
          .then((cloudRes) => {
            if (cloudRes.success && cloudRes.data?.public_url) {
              archiveItem.observations =
                (archiveItem.observations
                  ? archiveItem.observations + " • "
                  : "") + `Cloudflare R2 Scellé : ${cloudRes.data.public_url}`;
              this.persistLocal();
            }
          })
          .catch(() => {});
      } catch {}

      // Diffusion temps réel WebSocket
      realtimeSocket.sendEvent("DOSSIER_ARCHIVE_CREATED", {
        archive_id: archiveId,
        candidat_nom: `${dossierDetail.candidat.nom} ${dossierDetail.candidat.prenoms}`,
        numero_dossier: dossierDetail.dossier.numero_dossier,
        motif: payload.motif,
        date_archivage: nowIso,
      });

      guichetStore.showToast(
        `Le dossier ${dossierDetail.dossier.numero_dossier} (${dossierDetail.candidat.nom}) a été transféré au Centre d'Archivage avec succès.`,
        "success",
      );

      return archiveItem;
    },

    /**
     * Restauration / Désarchivage d'un dossier
     */
    async desarchiverDossier(archiveId: string) {
      const guichetStore = useGuichetStore();
      const item = this.archives.find((a) => a.id === archiveId);
      if (!item) {
        throw new Error("Dossier archivé introuvable");
      }

      // Retirer des archives
      this.archives = this.archives.filter((a) => a.id !== archiveId);
      this.persistLocal();

      // Recharger le guichet
      await guichetStore.search(item.numero_dossier);
      if (guichetStore.selectedCandidatId) {
        await guichetStore.selectCandidat(guichetStore.selectedCandidatId);
        if (guichetStore.activeDossierDetail) {
          guichetStore.activeDossierDetail.dossier.statut_dossier = "EN_COURS";
        }
      }

      realtimeSocket.sendEvent("DOSSIER_DESARCHIVE", {
        archive_id: archiveId,
        numero_dossier: item.numero_dossier,
      });

      guichetStore.showToast(
        `Le dossier ${item.numero_dossier} a été restauré avec succès dans les dossiers actifs.`,
        "success",
      );
    },

    /**
     * Suppression définitive certifiée (réservée au Directeur après expiration de rétention légale)
     */
    async supprimerArchive(archiveId: string) {
      const guichetStore = useGuichetStore();
      const authStore = useAuthStore();
      if (authStore.currentRole !== "manager") {
        throw new Error(
          "Seul le Directeur est habilité à purger définitivement une archive.",
        );
      }

      const item = this.archives.find((a) => a.id === archiveId);
      this.archives = this.archives.filter((a) => a.id !== archiveId);
      this.persistLocal();

      guichetStore.showToast(
        `Archive ${item?.numero_dossier || archiveId} définitivement purgée du stockage.`,
        "info",
      );
    },

    /**
     * Téléversement d'un document important dans le stockage des documents
     */
    async televerserDocumentImportant(payload: {
      candidatId: string;
      dossierId: string;
      type: TypeDocumentImportant;
      titre: string;
      nomFichier: string;
      base64: string;
      mimeType?: string;
      taille?: number;
      notes?: string;
    }) {
      const authStore = useAuthStore();
      const guichetStore = useGuichetStore();

      const doc: DocumentArchive = {
        id: `DOC-IMP-${Date.now()}-${Math.random().toString(36).slice(2, 6)}`,
        dossier_id: payload.dossierId,
        candidat_id: payload.candidatId,
        type_document: payload.type,
        titre: payload.titre,
        nom_fichier: payload.nomFichier,
        taille_octets: payload.taille,
        fichier_url_ou_base64: payload.base64,
        mime_type: payload.mimeType || "application/pdf",
        date_televersement: new Date().toISOString(),
        televerse_par_nom:
          `${authStore.currentUser?.nom || ""} ${authStore.currentUser?.prenom || ""}`.trim() ||
          "Guichet",
        est_archive: false,
        notes: payload.notes,
      };

      this.documentsImportants.unshift(doc);
      this.persistLocal();

      // Téléversement asynchrone Cloudflare R2 (Bucket : maryam-logistic-bucket)
      try {
        fetch(`${API_BASE_URL}/archivage/cloud/upload`, {
          method: "POST",
          headers: getAuthHeaders(),
          body: JSON.stringify({
            candidat_id: payload.candidatId,
            dossier_id: payload.dossierId,
            type_document: payload.type,
            titre: payload.titre,
            nom_fichier: payload.nomFichier,
            fichier_base64: payload.base64,
            mime_type: payload.mimeType || "application/pdf",
          }),
        })
          .then((res) => res.json())
          .then((cloudRes) => {
            if (cloudRes.success && cloudRes.data?.public_url) {
              doc.fichier_url_ou_base64 = cloudRes.data.public_url;
              doc.notes =
                (doc.notes ? doc.notes + " • " : "") +
                `Cloudflare R2 : ${cloudRes.data.bucket}/${cloudRes.data.object_key}`;
              this.persistLocal();
            }
          })
          .catch(() => {});
      } catch {}

      guichetStore.showToast(
        `Document "${payload.titre}" téléversé et sécurisé dans les documents importants.`,
        "success",
      );

      return doc;
    },

    /**
     * Recherche les documents importants d'un candidat spécifique
     */
    getDocumentsForCandidat(candidatId: string): DocumentArchive[] {
      return this.documentsImportants.filter(
        (d) => d.candidat_id === candidatId,
      );
    },
  },
});

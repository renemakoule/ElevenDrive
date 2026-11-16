import { defineStore } from 'pinia';
import {
  offlineStorageService,
  formatBytes,
  type OfflineItem,
} from '../services/offlineStorageService';
import type { EpreuvePdf, SupportCours, TutorielVideo } from '../types';
import { useGuichetStore } from './guichetStore';
import { neonSyncEngine } from '../services/neonSyncService';

export const useOfflineDownloadsStore = defineStore('offlineDownloads', {
  state: () => ({
    items: [] as OfflineItem[],
    activeFilter: 'ALL' as 'ALL' | 'SUPPORT' | 'EPREUVE' | 'TUTORIEL',
    searchQuery: '',
    isLoading: false,
    isOnline: typeof navigator !== 'undefined' ? navigator.onLine : true,
    activeItemToPreview: null as OfflineItem | null,
  }),

  getters: {
    downloadedIds(state): Set<string> {
      return new Set(state.items.map((i) => i.original_id));
    },

    filteredItems(state): OfflineItem[] {
      return state.items.filter((item) => {
        if (state.activeFilter !== 'ALL' && item.type !== state.activeFilter) {
          return false;
        }
        if (state.searchQuery.trim()) {
          const q = state.searchQuery.toLowerCase().trim();
          const matchTitle = item.titre.toLowerCase().includes(q);
          const matchDesc = item.description?.toLowerCase().includes(q) || false;
          const matchCat = item.categorie_permis.toLowerCase().includes(q);
          if (!matchTitle && !matchDesc && !matchCat) {
            return false;
          }
        }
        return true;
      });
    },

    stats(state) {
      const total = state.items.length;
      const supportsCount = state.items.filter((i) => i.type === 'SUPPORT').length;
      const epreuvesCount = state.items.filter((i) => i.type === 'EPREUVE').length;
      const tutorielsCount = state.items.filter((i) => i.type === 'TUTORIEL').length;
      const totalBytes = state.items.reduce((acc, curr) => acc + (curr.taille_octets || 0), 0);

      return {
        total,
        supportsCount,
        epreuvesCount,
        tutorielsCount,
        totalBytes,
        totalSizeFormatted: formatBytes(totalBytes),
      };
    },
  },

  actions: {
    initListeners() {
      if (typeof window !== 'undefined') {
        window.addEventListener('online', () => {
          this.isOnline = true;
        });
        window.addEventListener('offline', () => {
          this.isOnline = false;
        });
        neonSyncEngine.subscribe((state) => {
          this.isOnline = state.isOnline;
        });
      }
    },

    async loadDownloads() {
      this.isLoading = true;
      try {
        this.items = await offlineStorageService.getAllItems();
      } catch (err) {
        console.error('Erreur chargement téléchargements hors-ligne:', err);
      } finally {
        this.isLoading = false;
      }
    },

    isDownloaded(originalId: string): boolean {
      return this.items.some((i) => i.original_id === originalId);
    },

    /**
     * Télécharge et enregistre une épreuve PDF hors-ligne
     */
    async downloadEpreuve(epreuve: EpreuvePdf) {
      const guichetStore = useGuichetStore();
      const existing = this.items.find((i) => i.original_id === epreuve.id);
      if (existing) {
        guichetStore.showToast(`"${epreuve.titre}" est déjà disponible hors-ligne.`, 'info');
        return;
      }

      const octets = epreuve.fichier_data ? Math.round(epreuve.fichier_data.length * 0.75) : 500000;
      const item: OfflineItem = {
        id: `epreuve_${epreuve.id}`,
        original_id: epreuve.id,
        type: 'EPREUVE',
        titre: epreuve.titre,
        categorie_permis: epreuve.categorie_permis,
        description: epreuve.description,
        date_telechargement: new Date().toISOString(),
        taille_octets: octets,
        taille_formatee: formatBytes(octets),
        fichier_nom: epreuve.fichier_nom || `${epreuve.titre}.pdf`,
        mime_type: 'application/pdf',
        blob_or_data: epreuve.fichier_data,
      };

      try {
        await offlineStorageService.saveItem(item);
        this.items.unshift(item);
        guichetStore.showToast(
          `Épreuve "${epreuve.titre}" enregistrée hors-ligne (${item.taille_formatee})`,
          'success'
        );
      } catch (err: any) {
        guichetStore.showToast(`Erreur enregistrement hors-ligne : ${err.message}`, 'error');
      }
    },

    /**
     * Télécharge et enregistre un support de cours théorique
     */
    async downloadSupport(support: SupportCours) {
      const guichetStore = useGuichetStore();
      const existing = this.items.find((i) => i.original_id === support.id);
      if (existing) {
        guichetStore.showToast(`"${support.titre}" est déjà disponible hors-ligne.`, 'info');
        return;
      }

      let octets = 250000;
      let blobOrData = support.fichier_data || '';

      if (support.type_support === 'PDF' && support.fichier_data) {
        octets = Math.round(support.fichier_data.length * 0.75);
      } else if (support.type_support === 'LIEN_WEB' && support.url_externe) {
        blobOrData = support.url_externe;
        octets = 15000; // métadonnées du cours web
      }

      const item: OfflineItem = {
        id: `support_${support.id}`,
        original_id: support.id,
        type: 'SUPPORT',
        titre: support.titre,
        categorie_permis: support.categorie_permis,
        description: support.description,
        date_telechargement: new Date().toISOString(),
        taille_octets: octets,
        taille_formatee: formatBytes(octets),
        fichier_nom: support.fichier_nom || `${support.titre}.pdf`,
        mime_type: support.type_support === 'PDF' ? 'application/pdf' : 'text/html',
        blob_or_data: blobOrData,
        url_externe: support.url_externe,
      };

      try {
        await offlineStorageService.saveItem(item);
        this.items.unshift(item);
        guichetStore.showToast(
          `Support "${support.titre}" enregistré hors-ligne (${item.taille_formatee})`,
          'success'
        );
      } catch (err: any) {
        guichetStore.showToast(`Erreur enregistrement hors-ligne : ${err.message}`, 'error');
      }
    },

    /**
     * Télécharge un tutoriel vidéo (MP4 direct en blob ou guide pédagogique hors-ligne)
     */
    async downloadTutoriel(tutoriel: TutorielVideo) {
      const guichetStore = useGuichetStore();
      const existing = this.items.find((i) => i.original_id === tutoriel.id);
      if (existing) {
        guichetStore.showToast(`"${tutoriel.titre}" est déjà disponible hors-ligne.`, 'info');
        return;
      }

      let blobOrData = tutoriel.url_video;
      let octets = 120000;
      let mimeType = 'video/mp4';

      // Si c'est un fichier direct MP4, tenter de le mettre en cache blob
      const urlLower = tutoriel.url_video.toLowerCase();
      if (urlLower.endsWith('.mp4') || urlLower.endsWith('.webm')) {
        try {
          const resp = await fetch(tutoriel.url_video);
          if (resp.ok) {
            const blob = await resp.blob();
            octets = blob.size;
            // Convertir blob en DataURL pour stockage pérenne
            blobOrData = await new Promise((res, rej) => {
              const reader = new FileReader();
              reader.onloadend = () => res(reader.result as string);
              reader.onerror = rej;
              reader.readAsDataURL(blob);
            });
          }
        } catch {
          // Fallback sur URL directe
          blobOrData = tutoriel.url_video;
        }
      }

      const item: OfflineItem = {
        id: `tutoriel_${tutoriel.id}`,
        original_id: tutoriel.id,
        type: 'TUTORIEL',
        titre: tutoriel.titre,
        categorie_permis: tutoriel.categorie_permis,
        description: tutoriel.description,
        date_telechargement: new Date().toISOString(),
        taille_octets: octets,
        taille_formatee: formatBytes(octets),
        mime_type: mimeType,
        blob_or_data: blobOrData,
        duree_minutes: tutoriel.duree_minutes,
        miniature_url: tutoriel.miniature_url,
        plateforme: tutoriel.plateforme,
        url_externe: tutoriel.url_video,
      };

      try {
        await offlineStorageService.saveItem(item);
        this.items.unshift(item);
        guichetStore.showToast(
          `Tutoriel "${tutoriel.titre}" enregistré hors-ligne (${item.taille_formatee})`,
          'success'
        );
      } catch (err: any) {
        guichetStore.showToast(`Erreur enregistrement hors-ligne : ${err.message}`, 'error');
      }
    },

    /**
     * Supprime un élément du stockage hors-ligne
     */
    async removeDownload(id: string) {
      const guichetStore = useGuichetStore();
      try {
        await offlineStorageService.deleteItem(id);
        this.items = this.items.filter((i) => i.id !== id);
        guichetStore.showToast('Élément supprimé du stockage hors-ligne.', 'info');
      } catch (err: any) {
        guichetStore.showToast(`Erreur suppression : ${err.message}`, 'error');
      }
    },

    /**
     * Supprime tous les éléments hors-ligne
     */
    async clearAllDownloads() {
      const guichetStore = useGuichetStore();
      try {
        await offlineStorageService.clearAll();
        this.items = [];
        guichetStore.showToast('Stockage hors-ligne vidé avec succès.', 'info');
      } catch (err: any) {
        guichetStore.showToast(`Erreur : ${err.message}`, 'error');
      }
    },

    openPreview(item: OfflineItem) {
      this.activeItemToPreview = item;
    },

    closePreview() {
      this.activeItemToPreview = null;
    },
  },
});

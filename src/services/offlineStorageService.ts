/**
 * Service de stockage hors-ligne IndexedDB
 * Permet de conserver les cours, épreuves et vidéos pour une lecture 100% hors-connexion
 * (Exactement comme l'espace téléchargement de YouTube / Netflix)
 */

export interface OfflineItem {
  id: string; // ex: 'epreuve_123', 'support_456', 'tutoriel_789'
  original_id: string;
  type: 'EPREUVE' | 'SUPPORT' | 'TUTORIEL';
  titre: string;
  categorie_permis: string; // 'B', 'A', 'TOUTES', etc.
  description?: string;
  date_telechargement: string; // ISO string
  taille_octets: number;
  taille_formatee: string; // ex: "2.4 Mo"
  fichier_nom?: string;
  mime_type: string;
  blob_or_data: string; // Data URL Base64 ou Blob URL
  // Spécifique tutoriel vidéo
  duree_minutes?: number;
  miniature_url?: string;
  plateforme?: string;
  url_externe?: string;
}

const DB_NAME = 'ElevendriveOfflineStore';
const DB_VERSION = 1;
const STORE_NAME = 'offline_downloads';

function openDB(): Promise<IDBDatabase> {
  return new Promise((resolve, reject) => {
    if (typeof window === 'undefined' || !window.indexedDB) {
      reject(new Error("IndexedDB n'est pas supporté par cet environnement."));
      return;
    }

    const request = indexedDB.open(DB_NAME, DB_VERSION);

    request.onupgradeneeded = (event) => {
      const db = (event.target as IDBOpenDBRequest).result;
      if (!db.objectStoreNames.contains(STORE_NAME)) {
        const store = db.createObjectStore(STORE_NAME, { keyPath: 'id' });
        store.createIndex('type', 'type', { unique: false });
        store.createIndex('categorie_permis', 'categorie_permis', { unique: false });
        store.createIndex('date_telechargement', 'date_telechargement', { unique: false });
      }
    };

    request.onsuccess = () => resolve(request.result);
    request.onerror = () => reject(request.error);
  });
}

export const offlineStorageService = {
  /**
   * Enregistre un document ou vidéo dans l'espace hors-ligne
   */
  async saveItem(item: OfflineItem): Promise<void> {
    const db = await openDB();
    return new Promise((resolve, reject) => {
      const tx = db.transaction(STORE_NAME, 'readwrite');
      const store = tx.objectStore(STORE_NAME);
      const req = store.put(item);

      req.onsuccess = () => resolve();
      req.onerror = () => reject(req.error);
    });
  },

  /**
   * Récupère un élément par son ID
   */
  async getItem(id: string): Promise<OfflineItem | undefined> {
    const db = await openDB();
    return new Promise((resolve, reject) => {
      const tx = db.transaction(STORE_NAME, 'readonly');
      const store = tx.objectStore(STORE_NAME);
      const req = store.get(id);

      req.onsuccess = () => resolve(req.result);
      req.onerror = () => reject(req.error);
    });
  },

  /**
   * Récupère tous les éléments téléchargés hors-ligne
   */
  async getAllItems(): Promise<OfflineItem[]> {
    const db = await openDB();
    return new Promise((resolve, reject) => {
      const tx = db.transaction(STORE_NAME, 'readonly');
      const store = tx.objectStore(STORE_NAME);
      const req = store.getAll();

      req.onsuccess = () => {
        const items: OfflineItem[] = req.result || [];
        // Trier par date décroissante (les plus récents d'abord)
        items.sort(
          (a, b) =>
            new Date(b.date_telechargement).getTime() -
            new Date(a.date_telechargement).getTime()
        );
        resolve(items);
      };
      req.onerror = () => reject(req.error);
    });
  },

  /**
   * Supprime un élément du stockage hors-ligne
   */
  async deleteItem(id: string): Promise<void> {
    const db = await openDB();
    return new Promise((resolve, reject) => {
      const tx = db.transaction(STORE_NAME, 'readwrite');
      const store = tx.objectStore(STORE_NAME);
      const req = store.delete(id);

      req.onsuccess = () => resolve();
      req.onerror = () => reject(req.error);
    });
  },

  /**
   * Supprime l'intégralité du stockage hors-ligne
   */
  async clearAll(): Promise<void> {
    const db = await openDB();
    return new Promise((resolve, reject) => {
      const tx = db.transaction(STORE_NAME, 'readwrite');
      const store = tx.objectStore(STORE_NAME);
      const req = store.clear();

      req.onsuccess = () => resolve();
      req.onerror = () => reject(req.error);
    });
  },

  /**
   * Calcule le stockage utilisé par les téléchargements
   */
  async getStorageUsage(): Promise<{ usedBytes: number; usedFormatted: string; count: number }> {
    const items = await this.getAllItems();
    const usedBytes = items.reduce((acc, curr) => acc + (curr.taille_octets || 0), 0);
    return {
      usedBytes,
      usedFormatted: formatBytes(usedBytes),
      count: items.length,
    };
  },
};

/**
 * Utilitaire de formatage de taille octets en Ko / Mo
 */
export function formatBytes(bytes: number, decimals = 1): string {
  if (!bytes || bytes === 0) return '0 Ko';
  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['Octets', 'Ko', 'Mo', 'Go'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i] || 'Mo'}`;
}

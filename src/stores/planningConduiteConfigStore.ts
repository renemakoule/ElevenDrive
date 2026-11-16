import { defineStore } from 'pinia';

export interface CreneauHoraireConfig {
  id: string;
  label: string;
  debut: string; // ex: "08:00"
  fin: string;   // ex: "10:00"
  actif: boolean;
  icon?: string;
}

export interface JourPassageConfig {
  id: number; // 1 = Lundi, 2 = Mardi, 3 = Mercredi, 4 = Jeudi, 5 = Vendredi, 6 = Samedi, 0 = Dimanche
  code: string;
  label: string;
  actif: boolean;
}

const STORAGE_KEY_CRENEAUX = 'autoschool_config_creneaux_conduite';
const STORAGE_KEY_JOURS = 'autoschool_config_jours_conduite';

const DEFAULT_CRENEAUX: CreneauHoraireConfig[] = [
  { id: 'cr-1', label: 'Matinée 1', debut: '08:00', fin: '10:00', actif: true },
  { id: 'cr-2', label: 'Matinée 2', debut: '10:30', fin: '12:30', actif: true },
  { id: 'cr-3', label: 'Après-midi 1', debut: '14:00', fin: '16:00', actif: true },
  { id: 'cr-4', label: 'Fin de journée', debut: '16:30', fin: '18:30', actif: true },
];

const DEFAULT_JOURS: JourPassageConfig[] = [
  { id: 1, code: 'LUN', label: 'Lundi', actif: true },
  { id: 2, code: 'MAR', label: 'Mardi', actif: true },
  { id: 3, code: 'MER', label: 'Mercredi', actif: true },
  { id: 4, code: 'JEU', label: 'Jeudi', actif: true },
  { id: 5, code: 'VEN', label: 'Vendredi', actif: true },
  { id: 6, code: 'SAM', label: 'Samedi', actif: true },
  { id: 0, code: 'DIM', label: 'Dimanche', actif: false },
];

export const usePlanningConduiteConfigStore = defineStore('planningConduiteConfig', {
  state: () => ({
    creneaux: [] as CreneauHoraireConfig[],
    joursPassage: [] as JourPassageConfig[],
    isConfigModalOpen: false,
  }),

  getters: {
    creneauxActifs(state): CreneauHoraireConfig[] {
      return state.creneaux.filter((c) => c.actif);
    },

    joursActifs(state): JourPassageConfig[] {
      return state.joursPassage.filter((j) => j.actif);
    },

    joursActifsLabels(): string {
      const actifs = this.joursActifs;
      if (actifs.length === 7) return 'Tous les jours (7j/7)';
      if (actifs.length === 6 && !actifs.some((j) => j.id === 0)) return 'Du Lundi au Samedi';
      return actifs.map((j) => j.label).join(', ');
    },
  },

  actions: {
    init() {
      // Charger créneaux
      if (typeof localStorage !== 'undefined') {
        const savedCreneaux = localStorage.getItem(STORAGE_KEY_CRENEAUX);
        if (savedCreneaux) {
          try {
            this.creneaux = JSON.parse(savedCreneaux);
          } catch {
            this.creneaux = [...DEFAULT_CRENEAUX];
          }
        } else {
          this.creneaux = [...DEFAULT_CRENEAUX];
        }

        // Charger jours de passage
        const savedJours = localStorage.getItem(STORAGE_KEY_JOURS);
        if (savedJours) {
          try {
            this.joursPassage = JSON.parse(savedJours);
          } catch {
            this.joursPassage = [...DEFAULT_JOURS];
          }
        } else {
          this.joursPassage = [...DEFAULT_JOURS];
        }
      } else {
        this.creneaux = [...DEFAULT_CRENEAUX];
        this.joursPassage = [...DEFAULT_JOURS];
      }
    },

    save() {
      if (typeof localStorage !== 'undefined') {
        localStorage.setItem(STORAGE_KEY_CRENEAUX, JSON.stringify(this.creneaux));
        localStorage.setItem(STORAGE_KEY_JOURS, JSON.stringify(this.joursPassage));
      }
    },

    ajouterCreneau(creneau: { label: string; debut: string; fin: string }) {
      const id = `cr-${Date.now()}`;
      this.creneaux.push({
        id,
        label: creneau.label.trim() || `Créneau ${creneau.debut}`,
        debut: creneau.debut,
        fin: creneau.fin,
        actif: true,
      });
      this.creneaux.sort((a, b) => a.debut.localeCompare(b.debut));
      this.save();
    },

    modifierCreneau(id: string, updates: Partial<CreneauHoraireConfig>) {
      const idx = this.creneaux.findIndex((c) => c.id === id);
      if (idx !== -1) {
        this.creneaux[idx] = { ...this.creneaux[idx], ...updates };
        this.creneaux.sort((a, b) => a.debut.localeCompare(b.debut));
        this.save();
      }
    },

    supprimerCreneau(id: string) {
      this.creneaux = this.creneaux.filter((c) => c.id !== id);
      this.save();
    },

    toggleJour(id: number) {
      const jour = this.joursPassage.find((j) => j.id === id);
      if (jour) {
        jour.actif = !jour.actif;
        this.save();
      }
    },

    setJourActif(id: number, actif: boolean) {
      const jour = this.joursPassage.find((j) => j.id === id);
      if (jour) {
        jour.actif = actif;
        this.save();
      }
    },

    resetToDefaults() {
      this.creneaux = JSON.parse(JSON.stringify(DEFAULT_CRENEAUX));
      this.joursPassage = JSON.parse(JSON.stringify(DEFAULT_JOURS));
      this.save();
    },

    /**
     * Vérifie si une date donnée (YYYY-MM-DD) correspond à un jour de passage coché
     */
    isJourDePassage(dateStr: string): boolean {
      if (!dateStr) return true;
      try {
        const [year, month, day] = dateStr.split('-').map(Number);
        const dateObj = new Date(year, month - 1, day);
        const dayOfWeek = dateObj.getDay(); // 0 = Dimanche, 1 = Lundi, etc.
        const jour = this.joursPassage.find((j) => j.id === dayOfWeek);
        return jour ? jour.actif : true;
      } catch {
        return true;
      }
    },

    getNomJour(dateStr: string): string {
      if (!dateStr) return '';
      try {
        const [year, month, day] = dateStr.split('-').map(Number);
        const dateObj = new Date(year, month - 1, day);
        const dayOfWeek = dateObj.getDay();
        const jour = this.joursPassage.find((j) => j.id === dayOfWeek);
        return jour ? jour.label : '';
      } catch {
        return '';
      }
    },
  },
});

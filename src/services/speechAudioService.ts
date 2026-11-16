import { ref, shallowRef } from 'vue';

export interface AudioSpeechState {
  isPlaying: boolean;
  isPaused: boolean;
  currentRate: number;
  currentSegmentIndex: number;
  totalSegments: number;
  currentSegmentText: string;
  voiceName: string;
}

class SpeechAudioService {
  private synth: SpeechSynthesis | null = null;
  private voices: SpeechSynthesisVoice[] = [];
  private currentUtterance: SpeechSynthesisUtterance | null = null;
  private heartbeatTimer: number | null = null;
  private segmentsQueue: string[] = [];
  private onSegmentChangeCallback: ((index: number, text: string) => void) | null = null;
  private onEndCallback: (() => void) | null = null;

  // États réactifs Vue 3
  public isPlaying = ref(false);
  public isPaused = ref(false);
  public currentRate = ref(1.0);
  public currentSegmentIndex = ref(0);
  public totalSegments = ref(0);
  public currentSegmentText = ref('');
  public selectedVoice = shallowRef<SpeechSynthesisVoice | null>(null);
  public availableVoices = shallowRef<SpeechSynthesisVoice[]>([]);
  public isSupported = ref(false);

  public currentLanguageCode = ref('fr');

  constructor() {
    if (typeof window !== 'undefined' && 'speechSynthesis' in window) {
      this.synth = window.speechSynthesis;
      this.isSupported.value = true;
      this.initVoices();
    }
  }

  private initVoices() {
    if (!this.synth) return;

    const load = () => {
      if (!this.synth) return;
      const allVoices = this.synth.getVoices();
      if (allVoices && allVoices.length > 0) {
        this.availableVoices.value = allVoices;
        const activeLang =
          (typeof localStorage !== 'undefined' &&
            (localStorage.getItem('elevendrive_lang') || localStorage.getItem('autoschool_lang'))) ||
          'fr';
        this.selectVoiceForLanguage(activeLang);
      }
    };

    load();
    if (typeof this.synth.onvoiceschanged !== 'undefined') {
      this.synth.onvoiceschanged = load;
    }
  }

  public selectVoiceForLanguage(langCode: string) {
    const list = this.availableVoices.value;
    const prefix = (langCode || 'fr').toLowerCase().slice(0, 2);
    this.currentLanguageCode.value = prefix;

    if (!list || list.length === 0) return;

    // 1. Voix locale pour la langue cible
    const bestLocal = list.find(
      (v) => v.lang.toLowerCase().startsWith(prefix) && v.localService
    );
    // 2. N'importe quelle voix pour cette langue
    const anyVoice = list.find((v) => v.lang.toLowerCase().startsWith(prefix));
    // 3. Fallback voix par défaut
    const defaultVoice = list.find((v) => v.default) || list[0];

    this.selectedVoice.value = bestLocal || anyVoice || defaultVoice || null;
  }

  private selectBestFrenchVoice() {
    this.selectVoiceForLanguage('fr');
  }

  /**
   * Découpe un texte brut en phrases et segments audibles et digestes
   */
  public decouperTexteEnSegments(texte: string): string[] {
    if (!texte || !texte.trim()) return [];

    // Nettoyage des retours chariots et espaces multiples
    const propre = texte
      .replace(/\r\n/g, '\n')
      .replace(/\t/g, ' ')
      .replace(/[ ]{2,}/g, ' ')
      .trim();

    // Découpage par ponctuations fortes (. ! ? :) ou doubles retours chariots
    const morceaux = propre
      .split(/(?<=[.!?:\n])\s+/)
      .map((s) => s.trim())
      .filter((s) => s.length > 1);

    const segmentsFinaux: string[] = [];

    // Si une phrase est trop longue (> 180 caractères), on la sous-découpe par virgules ou points-virgules
    for (const morceau of morceaux) {
      if (morceau.length <= 180) {
        segmentsFinaux.push(morceau);
      } else {
        const sousMorceaux = morceau
          .split(/(?<=[,;])\s+/)
          .map((s) => s.trim())
          .filter((s) => s.length > 0);

        if (sousMorceaux.length > 1) {
          segmentsFinaux.push(...sousMorceaux);
        } else {
          // Découpage brutal de sécurité si pas de ponctuation
          let reste = morceau;
          while (reste.length > 180) {
            let coupe = reste.lastIndexOf(' ', 180);
            if (coupe === -1) coupe = 180;
            segmentsFinaux.push(reste.slice(0, coupe).trim());
            reste = reste.slice(coupe).trim();
          }
          if (reste.length > 0) segmentsFinaux.push(reste);
        }
      }
    }

    return segmentsFinaux;
  }

  /**
   * Lance la lecture séquentielle d'une liste de segments avec suivi de progression
   */
  /**
   * Lance la lecture séquentielle d'une liste de segments avec suivi de progression
   */
  public lireSegments(
    segments: string[],
    indexDepart = 0,
    options?: {
      lang?: string;
      onSegmentChange?: (index: number, text: string) => void;
      onEnd?: () => void;
    }
  ) {
    if (!this.synth || !this.isSupported.value) return;

    this.arreter();

    if (options?.lang) {
      this.selectVoiceForLanguage(options.lang);
    }

    const segmentsFiltres = segments.filter((s) => s && s.trim().length > 0);
    if (segmentsFiltres.length === 0) return;

    this.segmentsQueue = segmentsFiltres;
    this.totalSegments.value = segmentsFiltres.length;
    this.currentSegmentIndex.value = Math.max(0, Math.min(indexDepart, segmentsFiltres.length - 1));
    this.onSegmentChangeCallback = options?.onSegmentChange || null;
    this.onEndCallback = options?.onEnd || null;

    this.isPlaying.value = true;
    this.isPaused.value = false;

    this.startHeartbeat();
    this.jouerSegmentCourant();
  }

  /**
   * Raccourci pour lire un texte complet en le découpant automatiquement
   */
  public lireTexte(
    texte: string,
    options?: {
      lang?: string;
      indexDepart?: number;
      onSegmentChange?: (index: number, text: string) => void;
      onEnd?: () => void;
    }
  ) {
    const segments = this.decouperTexteEnSegments(texte);
    if (segments.length === 0) return;
    this.lireSegments(segments, options?.indexDepart || 0, options);
  }

  private jouerSegmentCourant() {
    if (!this.synth || !this.isPlaying.value) return;

    const idx = this.currentSegmentIndex.value;
    if (idx < 0 || idx >= this.segmentsQueue.length) {
      this.terminerLecture();
      return;
    }

    const texte = this.segmentsQueue[idx];
    this.currentSegmentText.value = texte;

    if (this.onSegmentChangeCallback) {
      this.onSegmentChangeCallback(idx, texte);
    }

    // Annuler la diction précédente immédiate
    this.synth.cancel();

    const utterance = new SpeechSynthesisUtterance(texte);
    const defaultLocales: Record<string, string> = {
      fr: 'fr-FR',
      en: 'en-US',
      es: 'es-ES',
      ar: 'ar-SA',
    };
    const langCode = this.currentLanguageCode.value || 'fr';

    if (this.selectedVoice.value) {
      utterance.voice = this.selectedVoice.value;
      utterance.lang = this.selectedVoice.value.lang;
    } else {
      utterance.lang = defaultLocales[langCode] || 'fr-FR';
    }

    utterance.rate = this.currentRate.value;
    utterance.pitch = 1.0;

    utterance.onend = () => {
      if (!this.isPlaying.value || this.isPaused.value) return;

      if (this.currentSegmentIndex.value < this.segmentsQueue.length - 1) {
        this.currentSegmentIndex.value++;
        // Petite pause naturelle de 200ms entre deux phrases
        setTimeout(() => {
          this.jouerSegmentCourant();
        }, 180);
      } else {
        this.terminerLecture();
      }
    };

    utterance.onerror = (e) => {
      // Ignorer les erreurs "canceled" normales lors d'un stop/saut
      if (e.error !== 'canceled' && e.error !== 'interrupted') {
        console.warn('SpeechAudioService utterance error:', e);
      }
    };

    this.currentUtterance = utterance;
    this.synth.speak(utterance);
  }

  /**
   * Lit une courte phrase ponctuelle (ex: pour le quiz ou un feedback rapide)
   */
  public parlerCourt(texte: string, callback?: () => void) {
    if (!this.synth || !this.isSupported.value) {
      if (callback) callback();
      return;
    }

    this.synth.cancel();

    const utterance = new SpeechSynthesisUtterance(texte);
    if (this.selectedVoice.value) {
      utterance.voice = this.selectedVoice.value;
      utterance.lang = this.selectedVoice.value.lang;
    } else {
      utterance.lang = 'fr-FR';
    }

    utterance.rate = this.currentRate.value;
    utterance.onend = () => {
      if (callback) callback();
    };
    utterance.onerror = () => {
      if (callback) callback();
    };

    this.synth.speak(utterance);
  }

  public pause() {
    if (!this.synth || !this.isPlaying.value) return;
    this.synth.pause();
    this.isPaused.value = true;
  }

  public reprendre() {
    if (!this.synth || !this.isPlaying.value) return;
    this.synth.resume();
    this.isPaused.value = false;
  }

  public basculerPlayPause() {
    if (this.isPaused.value) {
      this.reprendre();
    } else if (this.isPlaying.value) {
      this.pause();
    }
  }

  public suivant() {
    if (this.currentSegmentIndex.value < this.segmentsQueue.length - 1) {
      this.currentSegmentIndex.value++;
      this.jouerSegmentCourant();
    }
  }

  public precedent() {
    if (this.currentSegmentIndex.value > 0) {
      this.currentSegmentIndex.value--;
      this.jouerSegmentCourant();
    }
  }

  public allerAuSegment(index: number) {
    if (index >= 0 && index < this.segmentsQueue.length) {
      this.currentSegmentIndex.value = index;
      this.jouerSegmentCourant();
    }
  }

  public setVitesse(vitesse: number) {
    this.currentRate.value = Math.max(0.5, Math.min(2.0, vitesse));
    // Si la lecture est en cours, relancer le segment avec le nouveau débit
    if (this.isPlaying.value && !this.isPaused.value) {
      this.jouerSegmentCourant();
    }
  }

  public arreter() {
    this.stopHeartbeat();
    if (this.synth) {
      this.synth.cancel();
    }
    this.isPlaying.value = false;
    this.isPaused.value = false;
    this.currentSegmentText.value = '';
    this.currentUtterance = null;
  }

  private terminerLecture() {
    this.arreter();
    if (this.onEndCallback) {
      this.onEndCallback();
    }
  }

  // Heartbeat pour éviter le gel des 15 secondes sur les moteurs Chromium / Tauri Webview
  private startHeartbeat() {
    this.stopHeartbeat();
    this.heartbeatTimer = window.setInterval(() => {
      if (this.synth && this.synth.speaking && !this.synth.paused) {
        this.synth.pause();
        this.synth.resume();
      }
    }, 12000);
  }

  private stopHeartbeat() {
    if (this.heartbeatTimer !== null) {
      clearInterval(this.heartbeatTimer);
      this.heartbeatTimer = null;
    }
  }
}

export const speechAudioService = new SpeechAudioService();

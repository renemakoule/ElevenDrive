/**
 * Moteur de notifications push natives (Bureau & Navigateur) et WebSocket temps réel
 */

import { getWsApiBaseUrl } from "../config/env";

type NotificationCallback = (event: {
  type: string;
  data: any;
  auto_ecole_id?: string;
}) => void;

const isTauriEnv = (): boolean => {
  return (
    typeof window !== "undefined" &&
    ("__TAURI_INTERNALS__" in window || "__TAURI__" in window)
  );
};

class NativeNotificationService {
  private ws: WebSocket | null = null;
  private reconnectTimer: any = null;
  private reconnectAttempts = 0;
  private currentToken: string | null = null;
  private listeners: Set<NotificationCallback> = new Set();
  private audioCtx: AudioContext | null = null;
  private soundEnabled = true;

  constructor() {
    // Vérifier les préférences sonores dans le stockage local
    try {
      const stored = localStorage.getItem("autoschool_sound_enabled");
      if (stored !== null) {
        this.soundEnabled = stored === "true";
      }
    } catch {
      // ignore
    }

    // Débloquer le contexte audio dès le premier clic ou interaction utilisateur (politique autoplay navigateur)
    this.initAudioUnlock();
  }

  private initAudioUnlock() {
    if (typeof window === "undefined") return;
    const unlock = () => {
      try {
        const AudioCtx =
          window.AudioContext || (window as any).webkitAudioContext;
        if (AudioCtx && !this.audioCtx) {
          this.audioCtx = new AudioCtx();
        }
        if (this.audioCtx && this.audioCtx.state === "suspended") {
          this.audioCtx.resume();
        }
      } catch {
        // ignore
      }
      window.removeEventListener("pointerdown", unlock);
      window.removeEventListener("keydown", unlock);
    };
    window.addEventListener("pointerdown", unlock, { once: true });
    window.addEventListener("keydown", unlock, { once: true });
  }

  /**
   * Vérifie et demande l'autorisation des notifications push natives du système d'exploitation
   */
  public async requestPermission(): Promise<boolean> {
    if (typeof window === "undefined") {
      return false;
    }

    // Dans l'environnement Desktop Tauri (Windows / macOS)
    if (isTauriEnv()) {
      try {
        localStorage.setItem("autoschool_desktop_push_enabled", "true");
      } catch {
        // ignore
      }
      this.playChime("success");

      // Tenter d'enregistrer la permission Web si supportée par la WebView2
      if ("Notification" in window && Notification.permission !== "denied") {
        try {
          await Notification.requestPermission();
        } catch {
          // ignore
        }
      }
      return true;
    }

    // Dans le navigateur Web standard
    if (!("Notification" in window)) {
      return false;
    }

    if (Notification.permission === "granted") {
      return true;
    }

    if (Notification.permission !== "denied") {
      try {
        const permission = await Notification.requestPermission();
        return permission === "granted";
      } catch {
        return false;
      }
    }

    return false;
  }

  public isPermissionGranted(): boolean {
    if (typeof window === "undefined") return false;

    // Sur le Desktop Tauri : vérification de la préférence enregistrée
    if (isTauriEnv()) {
      try {
        const stored = localStorage.getItem("autoschool_desktop_push_enabled");
        return stored === "true";
      } catch {
        return false;
      }
    }

    // Sur le Web standard
    if (!("Notification" in window)) return false;
    return Notification.permission === "granted";
  }

  public setSoundEnabled(enabled: boolean) {
    this.soundEnabled = enabled;
    try {
      localStorage.setItem("autoschool_sound_enabled", String(enabled));
    } catch {
      // ignore
    }
  }

  public isSoundEnabled(): boolean {
    return this.soundEnabled;
  }

  /**
   * Synthétise un son discret de notification moderne via l'API Web Audio (aucun fichier MP3 requis)
   */
  public playChime(type: "info" | "warning" | "critical" | "success" = "info") {
    if (!this.soundEnabled || typeof window === "undefined") return;

    try {
      const AudioCtx =
        window.AudioContext || (window as any).webkitAudioContext;
      if (!AudioCtx) return;

      if (!this.audioCtx) {
        this.audioCtx = new AudioCtx();
      }

      if (this.audioCtx.state === "suspended") {
        this.audioCtx.resume();
      }

      const now = this.audioCtx.currentTime;
      const osc = this.audioCtx.createOscillator();
      const gain = this.audioCtx.createGain();

      osc.connect(gain);
      gain.connect(this.audioCtx.destination);

      if (type === "critical") {
        // Double bip d'urgence
        osc.type = "sine";
        osc.frequency.setValueAtTime(659.25, now); // Mi5
        osc.frequency.setValueAtTime(880, now + 0.12); // La5
        gain.gain.setValueAtTime(0.2, now);
        gain.gain.exponentialRampToValueAtTime(0.001, now + 0.35);
        osc.start(now);
        osc.stop(now + 0.35);
      } else {
        // Carillon mélodique doux
        osc.type = "sine";
        osc.frequency.setValueAtTime(523.25, now); // Do5
        osc.frequency.exponentialRampToValueAtTime(783.99, now + 0.12); // Sol5
        gain.gain.setValueAtTime(0.12, now);
        gain.gain.exponentialRampToValueAtTime(0.001, now + 0.3);
        osc.start(now);
        osc.stop(now + 0.3);
      }
    } catch {
      // Ignorer si les restrictions du navigateur bloquent l'audio sans interaction préalable
    }
  }

  /**
   * Émet une notification push native OS (Windows / macOS / Navigateur)
   */
  public async showNotification(
    titre: string,
    options: {
      body: string;
      type?: "info" | "warning" | "critical" | "success";
      icon?: string;
      tag?: string;
      onClick?: () => void;
    },
  ) {
    this.playChime(options.type);

    const hasPermission = await this.requestPermission();
    if (!hasPermission) return;

    try {
      const notif = new Notification(titre, {
        body: options.body,
        icon: options.icon || "/app-icon.png",
        tag: options.tag,
      });

      notif.onclick = () => {
        window.focus();
        if (options.onClick) {
          options.onClick();
        }
        notif.close();
      };
    } catch {
      // Fallback gracieux si l'API Notification est restreinte
    }
  }

  /**
   * Résout l'URL WebSocket du backend en tenant compte des variables d'environnement,
   * de Tauri ou des ports Vite de développement.
   */
  private getWsUrl(token: string): string {
    const base = getWsApiBaseUrl();
    return base ? `${base}?token=${encodeURIComponent(token)}` : "";
  }

  /**
   * Établit la connexion WebSocket temps réel authentifiée
   */
  public connect(token: string) {
    if (typeof window === "undefined" || !token) return;
    if (
      this.ws &&
      (this.ws.readyState === WebSocket.OPEN ||
        this.ws.readyState === WebSocket.CONNECTING)
    ) {
      return;
    }

    this.currentToken = token;
    try {
      const wsUrl = this.getWsUrl(token);
      if (!wsUrl) return;
      this.ws = new WebSocket(wsUrl);

      this.ws.onopen = () => {
        this.reconnectAttempts = 0;
        console.log("⚡ [NotificationWS] Connecté en temps réel avec succès");
      };

      this.ws.onmessage = (event) => {
        try {
          const parsed = JSON.parse(event.data);
          this.handleIncomingMessage(parsed);
        } catch {
          // ignore non-json frames
        }
      };

      this.ws.onclose = () => {
        this.scheduleReconnect(token);
      };

      this.ws.onerror = () => {
        if (this.ws) {
          try {
            this.ws.close();
          } catch {}
        }
      };
    } catch {
      this.scheduleReconnect(token);
    }
  }

  public disconnect() {
    this.currentToken = null;
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
    if (this.ws) {
      try {
        this.ws.close();
      } catch {}
      this.ws = null;
    }
  }

  private scheduleReconnect(token: string) {
    if (this.reconnectTimer || !this.currentToken) return;
    this.reconnectAttempts++;
    const delay = Math.min(1000 * Math.pow(1.5, this.reconnectAttempts), 15000);
    this.reconnectTimer = setTimeout(() => {
      this.reconnectTimer = null;
      if (this.currentToken) {
        this.connect(this.currentToken);
      }
    }, delay);
  }

  public subscribe(cb: NotificationCallback): () => void {
    this.listeners.add(cb);
    return () => {
      this.listeners.delete(cb);
    };
  }

  private handleIncomingMessage(msg: {
    type: string;
    data: any;
    auto_ecole_id?: string;
  }) {
    // 1. Déclencher instantanément le carillon audio pour toute nouvelle notification
    if (msg.type === "NOTIFICATION_NEW" && msg.data) {
      const notif = msg.data;
      const type = (notif.type as any) || "info";
      this.playChime(type);

      // Notification Push Native Système d'exploitation
      this.showNotification(notif.titre || "ElevenDrive Notification", {
        body: notif.message || "",
        type: type,
        tag: notif.id,
      });
    } else if (msg.type === "SCHOOL_REGISTERED" && msg.data) {
      const school = msg.data;
      this.playChime("success");
      this.showNotification("🏢 Nouvelle Auto-École Enregistrée !", {
        body: `${school.nom_commercial} (${school.ville || "Actif"}, ${school.pays || "Réseau"}) • Directeur: ${school.directeur}`,
        type: "success",
        tag: `school-${school.id}`,
      });
    }

    // 2. Distribuer à tous les écouteurs Pinia / Vue
    this.listeners.forEach((listener) => {
      try {
        listener(msg);
      } catch (err) {
        console.error("Erreur listener notification:", err);
      }
    });
  }
}

export const nativeNotificationEngine = new NativeNotificationService();

import { ref } from "vue";
import { authBackend } from "./authBackendService";
import { getWsApiBaseUrl } from "../config/env";

export interface RealtimeWSMessage<T = any> {
  type: string;
  auto_ecole_id?: string;
  sender_id?: string;
  timestamp?: string;
  data?: T;
}

type RealtimeListener = (message: RealtimeWSMessage) => void;

class RealtimeSocketService {
  private socket: WebSocket | null = null;
  private listeners = new Set<RealtimeListener>();
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  private pingTimer: ReturnType<typeof setInterval> | null = null;
  private manualClose = false;

  public isConnected = ref(false);
  public lastEventType = ref<string>("");
  public lastEventAt = ref<string>("");

  private getWsBaseUrl(): string {
    const rawApi =
      (import.meta.env?.VITE_API_URL as string) ||
      "http://localhost:8080/api/auth";
    const rootUrl = rawApi
      .replace(/\/api\/auth\/?$/, "")
      .replace(/\/api\/?$/, "");
    const wsRoot = rootUrl
      .replace(/^http:\/\//i, "ws://")
      .replace(/^https:\/\//i, "wss://");
    return `${wsRoot}/api/ws`;
    return getWsApiBaseUrl();
  }

  public connect(): void {
    if (typeof window === "undefined") return;
    if (
      this.socket &&
      (this.socket.readyState === WebSocket.OPEN ||
        this.socket.readyState === WebSocket.CONNECTING)
    ) {
      return;
    }

    let token = "";
    let autoEcoleId = "";
    let userId = "";
    let role = "";

    try {
      const session =
        typeof authBackend.getActiveSession === "function"
          ? authBackend.getActiveSession()
          : null;
      if (session) {
        token = session.token || (session as any).access_token || "";
        autoEcoleId =
          session.auto_ecole?.id ||
          (session as any).school?.id ||
          session.user?.auto_ecole_id ||
          "";
        userId = session.user?.id || "";
        role = (session.user?.role || "").toUpperCase();
      }

      if (!autoEcoleId && typeof localStorage !== "undefined") {
        const webRaw = localStorage.getItem("autoschool_web_session");
        if (webRaw) {
          const webSess = JSON.parse(webRaw);
          autoEcoleId =
            webSess?.school?.id ||
            webSess?.auto_ecole?.id ||
            webSess?.user?.auto_ecole_id ||
            "";
          userId = userId || webSess?.user?.id || "";
          role =
            role || (webSess?.role || webSess?.user?.role || "").toUpperCase();
        }
      }
    } catch (err) {
      console.warn("[RealtimeWS] Lecture session impossible:", err);
      return;
    }

    if (!token || !autoEcoleId) {
      return;
    }

    this.manualClose = false;
    const params = new URLSearchParams({
      token,
      auto_ecole_id: autoEcoleId,
      user_id: userId,
      role,
    });

    const wsUrl = `${this.getWsBaseUrl()}?${params.toString()}`;

    try {
      this.socket = new WebSocket(wsUrl);

      this.socket.onopen = () => {
        this.isConnected.value = true;
        this.startPing();
      };

      this.socket.onmessage = (event: MessageEvent) => {
        try {
          const parsed: RealtimeWSMessage = JSON.parse(event.data);
          if (!parsed || !parsed.type) return;

          if (parsed.type !== "PONG") {
            this.lastEventType.value = parsed.type;
            this.lastEventAt.value =
              parsed.timestamp || new Date().toISOString();
          }

          this.listeners.forEach((listener) => {
            try {
              listener(parsed);
            } catch (err) {
              console.warn("[RealtimeWS] Erreur listener:", err);
            }
          });
        } catch {
          // Ignorer trames non-JSON
        }
      };

      this.socket.onclose = () => {
        this.isConnected.value = false;
        this.stopPing();
        this.socket = null;
        if (!this.manualClose) {
          this.scheduleReconnect();
        }
      };

      this.socket.onerror = () => {
        this.isConnected.value = false;
      };
    } catch {
      this.scheduleReconnect();
    }
  }

  public disconnect(): void {
    this.manualClose = true;
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
    this.stopPing();
    if (this.socket) {
      try {
        this.socket.close();
      } catch {
        // ignore
      }
      this.socket = null;
    }
    this.isConnected.value = false;
  }

  public subscribe(listener: RealtimeListener): () => void {
    this.listeners.add(listener);
    if (!this.isConnected.value) {
      this.connect();
    }
    return () => {
      this.listeners.delete(listener);
    };
  }

  public sendEvent(type: string, data: any): void {
    if (!this.socket || this.socket.readyState !== WebSocket.OPEN) {
      this.connect();
      return;
    }
    try {
      this.socket.send(
        JSON.stringify({
          type,
          data,
          timestamp: new Date().toISOString(),
        }),
      );
    } catch (err) {
      console.warn("[RealtimeWS] Impossible d’envoyer l’événement:", err);
    }
  }

  private scheduleReconnect(): void {
    if (this.reconnectTimer || this.manualClose) return;
    this.reconnectTimer = setTimeout(() => {
      this.reconnectTimer = null;
      this.connect();
    }, 2500);
  }

  private startPing(): void {
    this.stopPing();
    this.pingTimer = setInterval(() => {
      if (this.socket && this.socket.readyState === WebSocket.OPEN) {
        try {
          this.socket.send(JSON.stringify({ type: "PING" }));
        } catch {
          // ignore
        }
      }
    }, 25000);
  }

  private stopPing(): void {
    if (this.pingTimer) {
      clearInterval(this.pingTimer);
      this.pingTimer = null;
    }
  }
}

export const realtimeSocket = new RealtimeSocketService();

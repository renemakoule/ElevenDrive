import type {
  CreateSupportTicketPayload,
  SupportTicket,
  NetworkAlert,
  NetworkAnnouncement,
  CreateFeedbackPayload,
  SchoolFeedback,
} from "../types";
import { authBackend } from "./authBackendService";
import { API_BASE_URL } from "../config/env";

const getHeaders = () => {
  const session = authBackend.getActiveSession();
  const token = session?.token;
  return {
    "Content-Type": "application/json",
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
  };
};

export const supportService = {
  /**
   * Soumet un ticket d'assistance technique ou de signalement vers PostgreSQL
   */
  async createTicket(
    payload: CreateSupportTicketPayload,
  ): Promise<SupportTicket> {
    const url = `${API_BASE_URL}/support/ticket`;
    const res = await fetch(url, {
      method: "POST",
      headers: getHeaders(),
      body: JSON.stringify(payload),
    });
    const json = await res.json();
    if (!res.ok || !json.success) {
      throw new Error(
        json.error || "Impossible d'envoyer votre demande de support",
      );
    }
    return json.data;
  },

  /**
   * Récupère les alertes actives pour l'auto-école connectée
   */
  async getActiveAlerts(): Promise<NetworkAlert[]> {
    const url = `${API_BASE_URL}/communication/alerts`;
    const res = await fetch(url, { headers: getHeaders() });
    const json = await res.json();
    if (!res.ok || !json.success) {
      throw new Error(json.error || "Impossible de récupérer les alertes");
    }
    return json.data || [];
  },

  /**
   * Récupère les annonces réseau pour l'auto-école connectée
   */
  async getAnnouncements(): Promise<NetworkAnnouncement[]> {
    const url = `${API_BASE_URL}/communication/announcements`;
    const res = await fetch(url, { headers: getHeaders() });
    const json = await res.json();
    if (!res.ok || !json.success) {
      throw new Error(
        json.error || "Impossible de récupérer les annonces réseau",
      );
    }
    return json.data || [];
  },

  /**
   * Permet à l'auto-école de soumettre une notation ou un avis
   */
  async submitFeedback(
    payload: CreateFeedbackPayload,
  ): Promise<SchoolFeedback> {
    const url = `${API_BASE_URL}/communication/feedbacks`;
    const res = await fetch(url, {
      method: "POST",
      headers: getHeaders(),
      body: JSON.stringify(payload),
    });
    const json = await res.json();
    if (!res.ok || !json.success) {
      throw new Error(json.error || "Impossible de transmettre votre avis");
    }
    return json.data;
  },
};

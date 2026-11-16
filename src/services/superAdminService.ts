import type {
  SuperAdminStats,
  SchoolItemSummary,
  SuperAdminGlobalUser,
  SuperAdminGlobalCandidat,
  UpdateSchoolAdminPayload,
  UpdateSchoolLicencePayload,
  AuditLogEntry,
  NetworkAlert,
  CreateAlertPayload,
  NetworkAnnouncement,
  CreateAnnouncementPayload,
  SchoolFeedback,
  CreateFeedbackPayload,
  SupportTicket,
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

export const superAdminApi = {
  async getDashboardData(
    search: string = "",
  ): Promise<{ stats: SuperAdminStats; schools: SchoolItemSummary[] }> {
    const url = `${API_BASE_URL}/super-admin/dashboard${search ? `?q=${encodeURIComponent(search)}` : ""}`;
    const res = await fetch(url, { headers: getHeaders() });
    const json = await res.json();
    if (!res.ok || !json.success) {
      throw new Error(
        json.error || "Impossible de charger les données console Super-Admin",
      );
    }
    return json.data;
  },

  async getUsersData(search: string = ""): Promise<SuperAdminGlobalUser[]> {
    const url = `${API_BASE_URL}/super-admin/users${search ? `?q=${encodeURIComponent(search)}` : ""}`;
    const res = await fetch(url, { headers: getHeaders() });
    const json = await res.json();
    if (!res.ok || !json.success) {
      throw new Error(
        json.error || "Impossible de charger la liste des utilisateurs",
      );
    }
    return json.data?.users || [];
  },

  async getCandidatsData(
    search: string = "",
  ): Promise<SuperAdminGlobalCandidat[]> {
    const url = `${API_BASE_URL}/super-admin/candidats${search ? `?q=${encodeURIComponent(search)}` : ""}`;
    const res = await fetch(url, { headers: getHeaders() });
    const json = await res.json();
    if (!res.ok || !json.success) {
      throw new Error(
        json.error || "Impossible de charger la liste des candidats",
      );
    }
    return json.data?.candidats || [];
  },

  async toggleSchoolStatus(
    schoolId: string,
    estActif: boolean,
  ): Promise<boolean> {
    const url = `${API_BASE_URL}/super-admin/schools/status`;
    const res = await fetch(url, {
      method: "POST",
      headers: getHeaders(),
      body: JSON.stringify({ school_id: schoolId, est_actif: estActif }),
    });
    const json = await res.json();
    if (!res.ok || !json.success) {
      throw new Error(json.error || "Échec modification statut auto-école");
    }
    return true;
  },

  async updateSchool(payload: UpdateSchoolAdminPayload): Promise<boolean> {
    const url = `${API_BASE_URL}/super-admin/schools/update`;
    const res = await fetch(url, {
      method: "POST",
      headers: getHeaders(),
      body: JSON.stringify(payload),
    });
    const json = await res.json();
    if (!res.ok || !json.success) {
      throw new Error(json.error || "Échec de la mise à jour de l'auto-école");
    }
    return true;
  },

  async updateLicence(payload: UpdateSchoolLicencePayload): Promise<boolean> {
    const url = `${API_BASE_URL}/super-admin/schools/licence`;
    const res = await fetch(url, {
      method: "POST",
      headers: getHeaders(),
      body: JSON.stringify(payload),
    });
    const json = await res.json();
    if (!res.ok || !json.success) {
      throw new Error(json.error || "Échec de la configuration de licence");
    }
    return true;
  },

  async updateModules(
    schoolId: string,
    modulesActifs: string[],
  ): Promise<boolean> {
    const url = `${API_BASE_URL}/super-admin/schools/modules`;
    const res = await fetch(url, {
      method: "POST",
      headers: getHeaders(),
      body: JSON.stringify({
        school_id: schoolId,
        modules_actifs: modulesActifs,
      }),
    });
    const json = await res.json();
    if (!res.ok || !json.success) {
      throw new Error(json.error || "Échec de l'attribution des modules");
    }
    return true;
  },

  async deleteSchool(schoolId: string, confirmName: string): Promise<boolean> {
    const url = `${API_BASE_URL}/super-admin/schools/delete`;
    const res = await fetch(url, {
      method: "POST",
      headers: getHeaders(),
      body: JSON.stringify({ school_id: schoolId, confirm_name: confirmName }),
    });
    const json = await res.json();
    if (!res.ok || !json.success) {
      throw new Error(json.error || "Échec de la suppression de l'auto-école");
    }
    return true;
  },

  async getAuditLogs(): Promise<AuditLogEntry[]> {
    const url = `${API_BASE_URL}/super-admin/audit`;
    const res = await fetch(url, { headers: getHeaders() });
    const json = await res.json();
    if (!res.ok || !json.success) {
      throw new Error(json.error || "Impossible de charger les logs d'audit");
    }
    return json.data?.logs || [];
  },

  // --- Alertes Réseau ---
  async getAlerts(): Promise<NetworkAlert[]> {
    const url = `${API_BASE_URL}/super-admin/alerts`;
    const res = await fetch(url, { headers: getHeaders() });
    const json = await res.json();
    if (!res.ok || !json.success) {
      throw new Error(json.error || "Impossible de charger les alertes réseau");
    }
    return json.data || [];
  },

  async createAlert(payload: CreateAlertPayload): Promise<NetworkAlert> {
    const url = `${API_BASE_URL}/super-admin/alerts`;
    const res = await fetch(url, {
      method: "POST",
      headers: getHeaders(),
      body: JSON.stringify(payload),
    });
    const json = await res.json();
    if (!res.ok || !json.success) {
      throw new Error(json.error || "Échec de la diffusion de l'alerte");
    }
    return json.data;
  },

  async toggleAlert(alertId: string): Promise<boolean> {
    const url = `${API_BASE_URL}/super-admin/alerts/toggle`;
    const res = await fetch(url, {
      method: "POST",
      headers: getHeaders(),
      body: JSON.stringify({ id: alertId }),
    });
    const json = await res.json();
    if (!res.ok || !json.success) {
      throw new Error(
        json.error || "Échec de modification du statut de l'alerte",
      );
    }
    return json.data?.est_active;
  },

  async deleteAlert(alertId: string): Promise<boolean> {
    const url = `${API_BASE_URL}/super-admin/alerts?id=${encodeURIComponent(alertId)}`;
    const res = await fetch(url, {
      method: "DELETE",
      headers: getHeaders(),
    });
    const json = await res.json();
    if (!res.ok || !json.success) {
      throw new Error(json.error || "Échec de la suppression de l'alerte");
    }
    return true;
  },

  // --- Annonces Réseau & Événements ---
  async getAnnouncements(): Promise<NetworkAnnouncement[]> {
    const url = `${API_BASE_URL}/super-admin/announcements`;
    const res = await fetch(url, { headers: getHeaders() });
    const json = await res.json();
    if (!res.ok || !json.success) {
      throw new Error(json.error || "Impossible de charger les annonces");
    }
    return json.data || [];
  },

  async createAnnouncement(
    payload: CreateAnnouncementPayload,
  ): Promise<NetworkAnnouncement> {
    const url = `${API_BASE_URL}/super-admin/announcements`;
    const res = await fetch(url, {
      method: "POST",
      headers: getHeaders(),
      body: JSON.stringify(payload),
    });
    const json = await res.json();
    if (!res.ok || !json.success) {
      throw new Error(json.error || "Échec de publication de l'annonce");
    }
    return json.data;
  },

  async deleteAnnouncement(announcementId: string): Promise<boolean> {
    const url = `${API_BASE_URL}/super-admin/announcements?id=${encodeURIComponent(announcementId)}`;
    const res = await fetch(url, {
      method: "DELETE",
      headers: getHeaders(),
    });
    const json = await res.json();
    if (!res.ok || !json.success) {
      throw new Error(json.error || "Échec de la suppression de l'annonce");
    }
    return true;
  },

  // --- Avis & Retours d'expérience ---
  async getFeedbacks(starsFilter?: number): Promise<SchoolFeedback[]> {
    const query =
      starsFilter && starsFilter >= 1 ? `?etoiles=${starsFilter}` : "";
    const url = `${API_BASE_URL}/super-admin/feedbacks${query}`;
    const res = await fetch(url, { headers: getHeaders() });
    const json = await res.json();
    if (!res.ok || !json.success) {
      throw new Error(
        json.error || "Impossible de charger les avis des auto-écoles",
      );
    }
    return json.data || [];
  },

  async createFeedback(
    payload: CreateFeedbackPayload,
  ): Promise<SchoolFeedback> {
    const url = `${API_BASE_URL}/super-admin/feedbacks`;
    const res = await fetch(url, {
      method: "POST",
      headers: getHeaders(),
      body: JSON.stringify(payload),
    });
    const json = await res.json();
    if (!res.ok || !json.success) {
      throw new Error(json.error || "Échec de l'enregistrement de l'avis");
    }
    return json.data;
  },

  async replyFeedback(
    feedbackId: string,
    reponseAdmin: string,
    statut?: string,
  ): Promise<boolean> {
    const url = `${API_BASE_URL}/super-admin/feedbacks/reply`;
    const res = await fetch(url, {
      method: "POST",
      headers: getHeaders(),
      body: JSON.stringify({
        feedback_id: feedbackId,
        reponse_admin: reponseAdmin,
        statut,
      }),
    });
    const json = await res.json();
    if (!res.ok || !json.success) {
      throw new Error(json.error || "Échec de l'envoi de la réponse");
    }
    return true;
  },

  async updateFeedbackStatus(
    feedbackId: string,
    statut: "nouveau" | "en_cours" | "resolu",
  ): Promise<boolean> {
    const url = `${API_BASE_URL}/super-admin/feedbacks/status`;
    const res = await fetch(url, {
      method: "PATCH",
      headers: getHeaders(),
      body: JSON.stringify({ feedback_id: feedbackId, statut }),
    });
    const json = await res.json();
    if (!res.ok || !json.success) {
      throw new Error(json.error || "Échec de mise à jour du statut de l'avis");
    }
    return true;
  },

  async deleteFeedback(feedbackId: string): Promise<boolean> {
    const url = `${API_BASE_URL}/super-admin/feedbacks?id=${encodeURIComponent(feedbackId)}`;
    const res = await fetch(url, {
      method: "DELETE",
      headers: getHeaders(),
    });
    const json = await res.json();
    if (!res.ok || !json.success) {
      throw new Error(json.error || "Échec de suppression de l'avis");
    }
    return true;
  },

  // --- Tickets de Support Technique ---
  async getSupportTickets(): Promise<SupportTicket[]> {
    const url = `${API_BASE_URL}/super-admin/support-tickets`;
    const res = await fetch(url, { headers: getHeaders() });
    const json = await res.json();
    if (!res.ok || !json.success) {
      throw new Error(
        json.error || "Impossible de charger les tickets de support",
      );
    }
    return json.data || [];
  },
};

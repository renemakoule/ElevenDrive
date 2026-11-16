import type { SystemNotificationItem, SystemAuditLogItem } from "../types";
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

export const notificationApi = {
  async getNotifications(
    unreadOnly = false,
    limit = 50,
  ): Promise<{
    notifications: SystemNotificationItem[];
    unread_count: number;
  }> {
    const url = `${API_BASE_URL}/notifications?unread=${unreadOnly}&limit=${limit}`;
    const res = await fetch(url, { headers: getHeaders() });
    const json = await res.json();
    if (!res.ok || !json.success) {
      throw new Error(json.error || "Impossible de charger les notifications");
    }
    return json.data || { notifications: [], unread_count: 0 };
  },

  async markAsRead(id: string): Promise<boolean> {
    const url = `${API_BASE_URL}/notifications/read`;
    const res = await fetch(url, {
      method: "POST",
      headers: getHeaders(),
      body: JSON.stringify({ id }),
    });
    const json = await res.json();
    return Boolean(res.ok && json.success);
  },

  async markAllAsRead(): Promise<boolean> {
    const url = `${API_BASE_URL}/notifications/read-all`;
    const res = await fetch(url, {
      method: "POST",
      headers: getHeaders(),
    });
    const json = await res.json();
    return Boolean(res.ok && json.success);
  },

  async getAuditLogs(limit = 100): Promise<SystemAuditLogItem[]> {
    const url = `${API_BASE_URL}/notifications/audit?limit=${limit}`;
    const res = await fetch(url, { headers: getHeaders() });
    const json = await res.json();
    if (!res.ok || !json.success) {
      throw new Error(json.error || "Impossible de charger les logs d'audit");
    }
    return json.data?.logs || [];
  },
};

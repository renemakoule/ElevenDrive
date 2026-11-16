// ==============================================================================
// AutoSchool - Client CRUD Agences / Succursales
// ==============================================================================

import type { Agence } from "../../types";
import {
  fetchAuthWithFallback,
  getActiveSessionFromStorage,
} from "./authSessionHelpers";

/** Récupérer la liste des succursales (agences) de l'auto-école */
export async function fetchSchoolAgences(): Promise<Agence[]> {
  const session = getActiveSessionFromStorage();
  const token = session?.token;
  const localAgencesKey = `elevendrive_agences_${session?.auto_ecole?.id || "default"}`;

  if (token) {
    try {
      const response = await fetchAuthWithFallback("/agences", {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (response.ok) {
        const result = await response.json();
        if (
          result.success &&
          Array.isArray(result.data) &&
          result.data.length > 0
        ) {
          if (typeof localStorage !== "undefined") {
            localStorage.setItem(localAgencesKey, JSON.stringify(result.data));
          }
          return result.data;
        }
      }
    } catch {
      // Ignorer
    }
  }

  if (typeof localStorage !== "undefined") {
    const cached = localStorage.getItem(localAgencesKey);
    if (cached) {
      try {
        const parsed = JSON.parse(cached);
        if (Array.isArray(parsed) && parsed.length > 0) {
          return parsed;
        }
      } catch {
        // ignore
      }
    }
  }

  return session?.agence ? [session.agence] : [];
}

/** Créer une nouvelle agence pour l'auto-école dans PostgreSQL */
export async function createSchoolAgence(payload: {
  nom: string;
  code_agence?: string;
  adresse?: string;
  ville?: string;
  commune?: string;
  telephone?: string;
  est_siege?: boolean;
}): Promise<Agence> {
  const session = getActiveSessionFromStorage();
  const token = session?.token;

  let response: Response | null = null;
  try {
    response = await fetchAuthWithFallback("/agences", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
      },
      body: JSON.stringify(payload),
    });
  } catch (err) {
    console.warn("Échec réseau création agence:", err);
  }

  if (response) {
    const result = await response.json();
    if (!response.ok || !result.success) {
      throw new Error(
        result.error || "Erreur lors de la création de l'agence.",
      );
    }
    return result.data as Agence;
  }

  return {
    id: "ag-" + Date.now().toString(36),
    auto_ecole_id: session?.auto_ecole?.id || "default",
    nom: payload.nom,
    code_agence:
      payload.code_agence || `AG-${Math.floor(10 + Math.random() * 90)}`,
    adresse: payload.adresse || "",
    ville: payload.ville || "Douala",
    commune: payload.commune || "",
    telephone: payload.telephone || "",
    est_siege: payload.est_siege || false,
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  };
}

/** Modifier une agence existante dans PostgreSQL */
export async function updateSchoolAgence(
  id: string,
  updates: Partial<Agence>,
): Promise<Agence> {
  const session = getActiveSessionFromStorage();
  const token = session?.token;

  let response: Response | null = null;
  try {
    response = await fetchAuthWithFallback(`/agences/${encodeURIComponent(id)}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
      },
      body: JSON.stringify(updates),
    });
  } catch (err) {
    console.warn("Échec réseau modification agence:", err);
  }

  if (response) {
    const result = await response.json();
    if (!response.ok || !result.success) {
      throw new Error(
        result.error || "Erreur lors de la modification de l'agence.",
      );
    }
    return result.data as Agence;
  }

  return { id, ...updates } as Agence;
}

/** Supprimer une succursale dans PostgreSQL */
export async function deleteSchoolAgence(id: string): Promise<void> {
  const session = getActiveSessionFromStorage();
  const token = session?.token;

  let response: Response | null = null;
  try {
    response = await fetchAuthWithFallback(`/agences/${encodeURIComponent(id)}`, {
      method: "DELETE",
      headers: {
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
      },
    });
  } catch (err) {
    console.warn("Échec réseau suppression agence:", err);
  }

  if (response) {
    const result = await response.json();
    if (!response.ok || !result.success) {
      throw new Error(
        result.error || "Erreur lors de la suppression de l'agence.",
      );
    }
  }
}


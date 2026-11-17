// ==============================================================================
// AutoSchool - Opérations Login, Inscription École & Super-Admin
// ==============================================================================

import type {
  AuthSessionData,
  LoginCredentials,
  RegisterSchoolPayload,
  UserProfile,
} from "../../types";
import {
  API_BASE_URL,
  STORAGE_KEY_REFRESH_TOKEN,
  authenticateOfflineFromStorage,
  buildSessionFromAuthData,
  generateSchoolSlug,
  getMachineFingerprint,
  saveSessionToStorage,
} from "./authSessionHelpers";

/** 1. Inscription d'une nouvelle Auto-École via l'API Go */
export async function registerNewSchool(
  payload: RegisterSchoolPayload,
): Promise<AuthSessionData> {
  const slug = generateSchoolSlug(payload.nom_commercial);

  let cleanDevise = "XOF";
  const rawDevise = (payload.devise || "").toUpperCase();
  if (rawDevise.includes("XAF")) cleanDevise = "XAF";
  else if (rawDevise.includes("GNF")) cleanDevise = "GNF";
  else if (rawDevise.includes("EUR")) cleanDevise = "EUR";
  else if (rawDevise.includes("USD")) cleanDevise = "USD";
  else if (rawDevise.length === 3) cleanDevise = rawDevise;

  const telDirecteur = (payload.directeur_telephone || "").trim();
  const telEcole = (payload.telephone_ecole || telDirecteur).trim();

  const body = {
    nom_commercial: payload.nom_commercial.trim(),
    code_unique: slug,
    numero_agrement: payload.numero_agrement.trim().toUpperCase(),
    telephone_ecole: telEcole || telDirecteur,
    email_ecole: payload.email_ecole || payload.directeur_email || null,
    pays: payload.pays || "Cameroun",
    devise: cleanDevise,
    nom_agence: `Siège · ${payload.commune || payload.ville || "Principal"}`,
    adresse_agence:
      payload.adresse || `${payload.commune || ""}, ${payload.ville || ""}`,
    telephone_agence: telEcole || telDirecteur,
    ville_agence: payload.ville || "Douala",
    commune_agence: payload.commune || "Makepe",
    nom_directeur: payload.directeur_nom.trim(),
    prenoms_directeur: payload.directeur_prenoms.trim(),
    telephone_directeur: telDirecteur || telEcole,
    email_directeur: payload.directeur_email || null,
    mot_de_passe: payload.directeur_mot_de_passe,
    code_pin_manager: payload.code_pin_manager || null,
  };

  let response: Response;
  try {
    response = await fetch(`${API_BASE_URL}/register-school`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });
  } catch {
    throw new Error(
      `Impossible de contacter le serveur d'authentification (${API_BASE_URL}). Vérifiez que le backend Go est démarré ou votre connexion réseau.`,
    );
  }

  const result = await response.json();
  if (!response.ok || !result.success) {
    throw new Error(
      result.error || "Erreur lors de l'enregistrement de l'auto-école.",
    );
  }

  const data = result.data;
  const baseSession = buildSessionFromAuthData(
    data,
    "Agence Siège",
    payload.code_pin_manager || data.user.code_pin_manager || undefined,
  );

  const fallbackCity = payload.ville || baseSession.auto_ecole.ville || "Douala";
  const cityPrefix =
    fallbackCity
      .normalize("NFD")
      .replace(/[\u0300-\u036f]/g, "")
      .replace(/[^a-zA-Z0-9]/g, "")
      .slice(0, 3)
      .toUpperCase() || "HQ";

  const session: AuthSessionData = {
    ...baseSession,
    user: {
      ...baseSession.user,
      initial: (data.user.prenoms?.[0] || "D").toUpperCase(),
    },
    agence: data.agence || {
      id: "ag-primary",
      auto_ecole_id: baseSession.auto_ecole.id,
      code_agence: `AG-${cityPrefix}-01`,
      nom: `Siège · ${payload.commune || fallbackCity}`,
      adresse: payload.adresse || "Siège social",
      telephone: baseSession.auto_ecole.telephone,
      ville: fallbackCity,
      commune: payload.commune || "Makepe",
      est_siege: true,
    },
  };

  saveSessionToStorage(session);
  return session;
}

/** 2. Connexion d'un utilisateur via l'API Go */
export async function loginWithCredentials(
  credentials: LoginCredentials,
  schoolSlug?: string,
): Promise<AuthSessionData> {
  const identifiant = credentials.identifiant?.trim();
  const motDePasse = credentials.mot_de_passe;

  if (!identifiant) {
    throw new Error(
      "Veuillez renseigner votre numéro de téléphone ou adresse email.",
    );
  }
  if (!motDePasse) throw new Error("Veuillez renseigner votre mot de passe.");

  let response: Response;
  let networkFailed = false;
  try {
    response = await fetch(`${API_BASE_URL}/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        identifiant,
        mot_de_passe: motDePasse,
        school_slug: schoolSlug || undefined,
        machine_fingerprint: getMachineFingerprint(),
        machine_name:
          typeof navigator !== "undefined"
            ? navigator.userAgent.includes("Windows")
              ? "PC Windows"
              : "Poste Desktop"
            : "Poste Desktop",
        os:
          typeof navigator !== "undefined" &&
          navigator.userAgent.includes("Windows")
            ? "Windows"
            : "Linux",
      }),
    });
  } catch {
    networkFailed = true;
  }

  if (networkFailed || (response! && !response.ok && response.status >= 500)) {
    const offlineSession = authenticateOfflineFromStorage(identifiant);
    if (offlineSession) {
      saveSessionToStorage(offlineSession);
      return offlineSession;
    }
    if (networkFailed) {
      throw new Error(
        `Impossible de contacter le serveur Cloud (${API_BASE_URL}). ` +
        `Pour une première connexion sur ce poste, l'application doit joindre le serveur pour synchroniser votre compte. ` +
        `Vérifiez votre connexion Internet et assurez-vous que le serveur est en ligne.`,
      );
    }
  }

  const result = await response!.json();
  if (!response!.ok || !result.success) {
    if (
      result.error &&
      (result.error.includes("dial tcp") ||
        result.error.includes("no such host") ||
        result.error.includes("Cluster PostgreSQL Neon non joignable") ||
        result.error.includes("recherche globale utilisateur"))
    ) {
      const offlineSession = authenticateOfflineFromStorage(identifiant);
      if (offlineSession) {
        saveSessionToStorage(offlineSession);
        return offlineSession;
      }
      throw new Error(
        `Base cloud Neon temporairement injoignable (hors-ligne) et aucun profil local trouvé pour "${identifiant}".`,
      );
    }

    throw new Error(result.error || "Identifiant ou mot de passe incorrect.");
  }

  const session = buildSessionFromAuthData(result.data, "Agence principale");
  saveSessionToStorage(session);
  return session;
}

/** Connexion spécifique Super-Admin Console */
export async function loginAsSuperAdmin(
  identifiant: string,
  motDePasse: string,
): Promise<AuthSessionData> {
  const cleanId = identifiant.trim();
  if (!cleanId || !motDePasse) {
    throw new Error(
      "Identifiant et mot de passe super-administrateur requis.",
    );
  }

  const backendBase = API_BASE_URL.replace(/\/auth\/?$/, "");
  let response: Response;
  try {
    response = await fetch(`${backendBase}/super-admin/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        identifiant: cleanId,
        mot_de_passe: motDePasse,
      }),
    });
  } catch {
    throw new Error(
      `Serveur non joignable (${backendBase}). Assurez-vous que le backend Go est démarré.`,
    );
  }

  const result = await response.json();
  if (!response.ok || !result.success) {
    throw new Error(
      result.error || "Accès refusé : identifiants super-admin invalides.",
    );
  }

  const data = result.data;
  const userProfile: UserProfile = {
    id: data.user.id,
    nom: data.user.nom,
    prenom: data.user.prenoms,
    role: "super_admin",
    roleLabel: "Super-Admin",
    initial: "SA",
    agence: "ElevenDrive Cloud Console",
    auto_ecole_id: data.school.id,
    telephone: data.user.telephone,
    email: data.user.email || undefined,
  };

  const session: AuthSessionData = {
    user: userProfile,
    auto_ecole: data.school,
    agence: data.agence,
    token: data.access_token,
  };

  if (data.refresh_token && typeof localStorage !== "undefined") {
    localStorage.setItem(STORAGE_KEY_REFRESH_TOKEN, data.refresh_token);
  }

  saveSessionToStorage(session);
  return session;
}


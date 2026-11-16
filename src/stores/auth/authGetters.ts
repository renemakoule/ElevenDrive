// ==============================================================================
// AutoSchool - Profils par défaut, Permissions, Licences & Modules (Auth Getters)
// ==============================================================================

import type { AppPlatform, AutoEcole, UserProfile, UserRole } from "../../types";

export const DEFAULT_ACTIVE_MODULES: string[] = [
  "guichet",
  "caisse",
  "conduite",
  "examens",
  "flotte",
  "archivage_cloud",
];

export const USER_PROFILES: Record<UserRole, UserProfile> = {
  secretaire: {
    id: "",
    nom: "Secrétaire",
    prenom: "Guichet",
    role: "secretaire",
    roleLabel: "Secrétaire",
    initial: "S",
    agence: "",
    telephone: "",
    email: "",
  },
  manager: {
    id: "",
    nom: "Directeur",
    prenom: "Gérant",
    role: "manager",
    roleLabel: "Directeur",
    initial: "D",
    agence: "",
    telephone: "",
    email: "",
  },
  moniteur: {
    id: "",
    nom: "Moniteur",
    prenom: "Instructeur",
    role: "moniteur",
    roleLabel: "Moniteur",
    initial: "M",
    agence: "",
    telephone: "",
    email: "",
  },
  eleve: {
    id: "",
    nom: "Candidat",
    prenom: "Élève",
    role: "eleve",
    roleLabel: "Candidat",
    initial: "C",
    agence: "",
    dossier: "",
    telephone: "",
    email: "",
  },
  super_admin: {
    id: "super-admin-root",
    nom: "SUPER",
    prenom: "Administrateur",
    role: "super_admin",
    roleLabel: "Super-Admin",
    initial: "SA",
    agence: "ElevenDrive Platform Root",
    telephone: "",
    email: "admin@elevendrive.com",
  },
};

export function computeIsRoleAllowed(
  role: UserRole,
  platform: AppPlatform,
): boolean {
  if (role === "super_admin" || role === "manager" || role === "secretaire") {
    return true;
  }
  if (role === "moniteur" || role === "eleve") return platform === "pwa";
  return false;
}

export function computeRolePermissions(role: UserRole) {
  return {
    isDirecteur: role === "manager",
    isSecretaire: role === "secretaire",
    isMoniteur: role === "moniteur",
    isEleve: role === "eleve",
    canManageSchool: role === "manager",
    canManageAgences: role === "manager",
    canViewAgences: role === "manager" || role === "secretaire",
    canConfigurePin: role === "manager",
    canChangePassword: role !== "eleve",
    roleTitle:
      role === "manager"
        ? "Directeur / Gérant"
        : role === "secretaire"
          ? "Secrétaire Guichet"
          : role === "moniteur"
            ? "Moniteur de Conduite"
            : "Candidat / Élève Permis",
  };
}

export function computeIsLicenceExpired(
  role: UserRole,
  school: AutoEcole | null,
): boolean {
  if (role === "super_admin" || !school) return false;
  if (school.est_actif === false) return true;
  if (school.date_expiration_licence) {
    return new Date(school.date_expiration_licence).getTime() < Date.now();
  }
  return false;
}

export function computeTrialDaysRemaining(school: AutoEcole | null): number {
  if (!school?.date_expiration_licence) return 30;
  const diff = new Date(school.date_expiration_licence).getTime() - Date.now();
  return Math.max(0, Math.ceil(diff / (1000 * 60 * 60 * 24)));
}

export function parseActiveModules(
  role: UserRole,
  school: AutoEcole | null,
): string[] {
  if (role === "super_admin" || !school?.modules_actifs) {
    return [...DEFAULT_ACTIVE_MODULES];
  }
  try {
    const parsed = JSON.parse(school.modules_actifs);
    if (Array.isArray(parsed)) return parsed;
  } catch {}
  return [...DEFAULT_ACTIVE_MODULES];
}

export function computeIsModuleEnabled(
  role: UserRole,
  school: AutoEcole | null,
  moduleKey: string,
): boolean {
  if (role === "super_admin") return true;

  let required = moduleKey;
  if (moduleKey === "guichet") required = "guichet";
  else if (moduleKey === "caisse" || moduleKey === "finance") required = "caisse";
  else if (
    moduleKey === "conduite-planning" ||
    moduleKey === "cours-planning" ||
    moduleKey === "supports-cours" ||
    moduleKey === "tutoriels" ||
    moduleKey === "epreuves"
  ) {
    required = "conduite";
  } else if (
    moduleKey === "examens" ||
    moduleKey === "code" ||
    moduleKey === "stats"
  ) {
    required = "examens";
  } else if (moduleKey === "flotte") {
    required = "flotte";
  } else if (moduleKey === "archivage") {
    required = "archivage_cloud";
  } else {
    return true;
  }

  return parseActiveModules(role, school).includes(required);
}

export function computeWelcomeSubtitle(
  role: UserRole,
  schoolName?: string,
): string {
  const nom = schoolName || "elevendrive";
  switch (role) {
    case "secretaire":
      return `${nom} · Guichet 🌿`;
    case "manager":
      return `${nom} · Direction & Pilotage 🌿`;
    case "moniteur":
      return `${nom} · Livret mobile de conduite 🌿`;
    case "eleve":
      return `${nom} · Apprenant permis 🌿`;
    case "super_admin":
      return "Console Centrale d'Administration Root ⚡";
    default:
      return `${nom} · Espace Pro 🌿`;
  }
}


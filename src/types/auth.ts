export type UserRole =
  | "secretaire"
  | "manager"
  | "moniteur"
  | "eleve"
  | "super_admin";
export type AppPlatform = "desktop" | "pwa";

export interface UserProfile {
  id: string;
  nom: string;
  prenom: string;
  role: UserRole;
  roleLabel: string;
  initial: string;
  agence: string;
  auto_ecole_id?: string;
  agence_id?: string;
  telephone?: string;
  email?: string;
  dossier?: string;
  photo_url?: string;
  code_pin_manager?: string;
}

export interface SchoolAccessLinks {
  portal_url: string;
  moniteur_link: string;
  eleve_link: string;
  inscription_link: string;
}

export interface AutoEcole {
  id: string;
  code_unique: string; // Ex: 'la-reference'
  nom_commercial: string;
  numero_agrement?: string;
  telephone: string;
  email?: string;
  devise: string;
  pays: string;
  ville: string;
  logo_url?: string;
  est_actif: boolean;
  max_machines?: number;
  machines_actives?: number;
  date_expiration_licence?: string;
  type_licence?: string;
  modules_actifs?: string;
  links?: SchoolAccessLinks;
  created_at: string;
}

export interface Agence {
  id: string;
  auto_ecole_id: string;
  code_agence: string;
  nom: string;
  adresse: string;
  telephone: string;
  ville: string;
  commune: string;
  est_siege: boolean;
  created_at?: string;
  updated_at?: string;
}

export interface UserTeamMember {
  id: string;
  auto_ecole_id: string;
  agence_id?: string;
  agence_nom: string;
  nom: string;
  prenoms: string;
  telephone: string;
  email?: string;
  role: "DIRECTEUR" | "SECRETAIRE" | "CAISSIER" | "MONITEUR" | "ELEVE";
  statut: "ACTIF" | "SUSPENDU" | "ARCHIVE";
  created_at: string;
}

export interface RegisterUserPayload {
  nom: string;
  prenoms: string;
  telephone: string;
  email?: string;
  mot_de_passe: string;
  role: "DIRECTEUR" | "SECRETAIRE" | "CAISSIER" | "MONITEUR" | "ELEVE";
  agence_id?: string;
  code_pin_manager?: string;
}

export interface RegisterSchoolPayload {
  nom_commercial: string;
  numero_agrement: string;
  pays: string;
  ville: string;
  commune: string;
  adresse: string;
  telephone_ecole: string;
  email_ecole?: string;
  devise: string;
  directeur_nom: string;
  directeur_prenoms: string;
  directeur_telephone: string;
  directeur_email?: string;
  directeur_mot_de_passe: string;
  code_pin_manager?: string;
}

export interface LoginCredentials {
  identifiant: string; // Téléphone ou Email
  mot_de_passe: string;
}

export interface AuthSessionData {
  user: UserProfile;
  auto_ecole: AutoEcole;
  agence: Agence;
  token: string;
}

export interface OAuthResult {
  status: "authenticated" | "needs_registration";
  session?: AuthSessionData;
  provider: "google" | "facebook" | "linkedin" | "apple";
  profile?: {
    email?: string;
    nom?: string;
    prenoms?: string;
  };
  message?: string;
}


export interface SchoolFeedbackItem {
  id: string;
  autoEcoleNom: string;
  autoEcoleId?: string;
  auteur: string;
  role: string;
  etoiles: number;
  categorie:
    | "ergonomie"
    | "guichet"
    | "pwa_mobile"
    | "examens"
    | "support"
    | "general";
  commentaire: string;
  date: string;
  statut: "nouveau" | "en_cours" | "resolu";
  reponseAdmin?: string;
}


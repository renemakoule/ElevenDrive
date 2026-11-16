export interface NetworkAnnouncementItem {
  id: string;
  titre: string;
  contenu: string;
  typeFormat: "text_only" | "text_image" | "image_only";
  imageUrl?: string;
  categorie: "evenement" | "mise_a_jour" | "formation" | "general";
  cibleType: "all" | "single";
  cibleSchoolId?: string;
  cibleSchoolName?: string;
  dateEvenement?: string;
  lienAction?: string;
  datePublication: string;
}


export interface ProvisionStep {
  id: string;
  category: "file" | "data" | "read" | "stat";
  title: string;
  detail: string;
  targetPercent: number;
}

export const PROVISION_STEPS: ProvisionStep[] = [
  {
    id: "step-1",
    category: "file",
    title: "Génération de l'espace d'isolation tenant",
    detail: "Création du schéma PostgreSQL dédié & politiques RLS de sécurité",
    targetPercent: 8,
  },
  {
    id: "step-2",
    category: "data",
    title: "Initialisation de l'Agence Siège & Agences secondaires",
    detail: "Génération des codes succursales (AG-HQ-01) et géolocalisation",
    targetPercent: 18,
  },
  {
    id: "step-3",
    category: "file",
    title: "Déploiement de la base locale SQLite & synchronisation Neon",
    detail: "Miroir local-first pour tolérance absolue aux coupures réseau",
    targetPercent: 28,
  },
  {
    id: "step-4",
    category: "data",
    title: "Injection du référentiel ministériel de permis",
    detail:
      "Catégories A (Motos), B (Véhicules légers), C (Poids lourds), D (Transport)",
    targetPercent: 40,
  },
  {
    id: "step-5",
    category: "read",
    title: "Vérification de l'intégrité des formules d'apprentissage",
    detail:
      "Validation des tranches de paiement, seuils d'acompte et barèmes d'heures",
    targetPercent: 52,
  },
  {
    id: "step-6",
    category: "file",
    title: "Attribution du stockage cloud haute sécurité",
    detail:
      "Configuration des buckets d'archivage Cloudflare R2 pour CNI et photos",
    targetPercent: 65,
  },
  {
    id: "step-7",
    category: "data",
    title: "Configuration des terminaux de caisse & reçus thermiques",
    detail: "Modèles d'impression thermique 80mm/58mm et journal de caisse Z",
    targetPercent: 78,
  },
  {
    id: "step-8",
    category: "read",
    title: "Compilation de la banque officielle de questions de Code",
    detail:
      "Génération du pool d'entraînement et simulateur d'examen 40 questions",
    targetPercent: 88,
  },
  {
    id: "step-9",
    category: "stat",
    title: "Validation Super-Admin & Attribution de la Licence Gratuite",
    detail: "Attribution automatique de 30 jours d'essai complets tout-inclus",
    targetPercent: 100,
  },
];

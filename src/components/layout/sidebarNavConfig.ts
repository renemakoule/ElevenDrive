import {
  LayoutGrid,
  Receipt,
  CheckSquare,
  Award,
  FileText,
  TrendingUp,
  Car,
  ShieldAlert,
  Calendar,
  FileSignature,
  Brain,
  CalendarPlus,
  Tag,
  Users,
  BookOpen,
  Video,
  Film,
  History,
  Archive,
  DownloadCloud,
} from "lucide-vue-next";

export interface NavItem {
  id: string;
  label: string;
  icon: any;
  badge?: string;
  shortcut?: string;
}

export interface NavGroup {
  id: string;
  title: string;
  items: NavItem[];
}

export function getRoleNavGroups(role: string | null | undefined): NavGroup[] {
  switch (role) {
    case "secretaire":
      return [
        {
          id: "acc_caisse",
          title: "Accueil & Guichet",
          items: [
            { id: "overview", label: "Vue d'ensemble", icon: LayoutGrid },
            {
              id: "guichet",
              label: "Guichet & Encaissement",
              icon: Receipt,
              shortcut: "F2",
            },
            {
              id: "caisse",
              label: "Point Caisse (Rapport Z)",
              icon: FileText,
              shortcut: "F4",
            },
          ],
        },
        {
          id: "plannings",
          title: "Plannings & Formation",
          items: [
            {
              id: "conduite-planning",
              label: "Planning Conduite Pratique",
              icon: Calendar,
            },
            {
              id: "cours-planning",
              label: "Planning Cours (Visio/Salle)",
              icon: Video,
            },
            {
              id: "supports-cours",
              label: "Supports de Cours",
              icon: FileText,
            },
            { id: "tutoriels", label: "Tutoriels Vidéo", icon: Film },
            { id: "epreuves", label: "Banque d'Épreuves PDF", icon: BookOpen },
          ],
        },
        {
          id: "examens",
          title: "Pédagogie & Examens",
          items: [
            { id: "tarifs", label: "Tarifs & Permis", icon: Tag },
            {
              id: "code",
              label: "Présences Code en Salle",
              icon: CheckSquare,
              shortcut: "F5",
            },
            {
              id: "examens",
              label: "Examens d'État & Bordereaux",
              icon: Award,
              shortcut: "F9",
            },
          ],
        },
        {
          id: "registres",
          title: "Registres & Rétention",
          items: [
            {
              id: "historique",
              label: "Historique des Opérations",
              icon: History,
            },
            { id: "archivage", label: "Centre d'Archivage", icon: Archive },
          ],
        },
      ];

    case "manager":
      return [
        {
          id: "direction",
          title: "Pilotage & Accueil",
          items: [
            { id: "overview", label: "Vue d'ensemble", icon: LayoutGrid },
          ],
        },
        {
          id: "pedagogie",
          title: "Planning & Pédagogie",
          items: [
            {
              id: "conduite-planning",
              label: "Planning Conduite Pratique",
              icon: Calendar,
            },
            {
              id: "cours-planning",
              label: "Planning Cours (Visio/Salle)",
              icon: Video,
            },
            {
              id: "supports-cours",
              label: "Supports de Cours",
              icon: FileText,
            },
            { id: "tutoriels", label: "Tutoriels Vidéo", icon: Film },
            { id: "epreuves", label: "Banque d'Épreuves PDF", icon: BookOpen },
          ],
        },
        {
          id: "exploitation",
          title: "Exploitation & Flotte",
          items: [
            { id: "equipe", label: "Équipe & Collaborateurs", icon: Users },
            { id: "tarifs", label: "Catalogue Tarifs Permis", icon: Tag },
            { id: "flotte", label: "Flotte Automobile & Alertes", icon: Car },
          ],
        },
        {
          id: "finance_audit",
          title: "Finance, Examens & Audit",
          items: [
            {
              id: "finance",
              label: "Audit & Journal Financier",
              icon: TrendingUp,
            },
            { id: "stats", label: "Taux de Réussite & Examens", icon: Award },
            {
              id: "derogations",
              label: "Dérogations Manager (PIN)",
              icon: ShieldAlert,
            },
            {
              id: "historique",
              label: "Historique & Journal Audit",
              icon: History,
            },
            {
              id: "archivage",
              label: "Centre d'Archivage & Rétention",
              icon: Archive,
            },
          ],
        },
      ];

    case "moniteur":
      return [
        {
          id: "planning",
          title: "Planning & Sessions",
          items: [
            { id: "overview", label: "Planning du Jour", icon: Calendar },
            {
              id: "conduite-planning",
              label: "Planning Conduite Pratique",
              icon: Calendar,
            },
            {
              id: "cours-planning",
              label: "Planning Cours & Visio",
              icon: Video,
            },
          ],
        },
        {
          id: "evaluations",
          title: "Pédagogie & Évaluations",
          items: [
            {
              id: "supports-cours",
              label: "Supports de Cours & Audio",
              icon: FileText,
            },
            { id: "epreuves", label: "Banque d'Épreuves PDF", icon: BookOpen },
            { id: "quiz-forms", label: "Quiz & Réponses Élèves", icon: Brain },
          ],
        },
        {
          id: "suivi",
          title: "Livret, Flotte & Émargements",
          items: [
            {
              id: "livret",
              label: "Fiche & Émargement Tactile",
              icon: FileSignature,
            },
            { id: "vehicule", label: "Véhicule Assigné", icon: Car },
            {
              id: "historique",
              label: "Historique & Émargements",
              icon: History,
            },
          ],
        },
      ];

    case "eleve":
      return [
        {
          id: "parcours",
          title: "Mon Parcours",
          items: [
            { id: "overview", label: "Mon Espace Apprenant", icon: LayoutGrid },
            {
              id: "reservation",
              label: "Réserver une Conduite",
              icon: CalendarPlus,
            },
            { id: "recus", label: "Mes Reçus & Solde", icon: Receipt },
          ],
        },
        {
          id: "cours_quiz",
          title: "Cours & Entraînement",
          items: [
            {
              id: "cours-online",
              label: "Cours en Ligne (Visio)",
              icon: Video,
            },
            { id: "quiz", label: "Séries Quiz Code 24/7", icon: Brain },
            {
              id: "supports-cours",
              label: "Mes Supports de Cours",
              icon: FileText,
            },
            { id: "tutoriels", label: "Mes Tutoriels Vidéo", icon: Film },
            {
              id: "epreuves",
              label: "Mes Épreuves & Sujets PDF",
              icon: BookOpen,
            },
            {
              id: "telechargements",
              label: "Téléchargements Hors-Ligne",
              icon: DownloadCloud,
            },
          ],
        },
        {
          id: "activite",
          title: "Activité & Historique",
          items: [
            {
              id: "historique",
              label: "Mon Historique & Activités",
              icon: History,
            },
          ],
        },
      ];

    default:
      return [];
  }
}


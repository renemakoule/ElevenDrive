import {
  LayoutGrid,
  Calendar,
  Car,
  TrendingUp,
  Receipt,
  CheckSquare,
  FileSignature,
  Brain,
  Video,
  CalendarPlus,
} from "lucide-vue-next";

export interface MobileTabItem {
  id: string;
  label: string;
  icon: any;
}

export function getMobilePrimaryTabs(
  role: string | null | undefined,
): MobileTabItem[] {
  switch (role) {
    case "manager":
      return [
        { id: "overview", label: "Pilotage", icon: LayoutGrid },
        { id: "conduite-planning", label: "Conduite", icon: Calendar },
        { id: "flotte", label: "Flotte", icon: Car },
        { id: "finance", label: "Finance", icon: TrendingUp },
      ];

    case "secretaire":
      return [
        { id: "overview", label: "Accueil", icon: LayoutGrid },
        { id: "guichet", label: "Dossiers", icon: Receipt },
        { id: "code", label: "Code", icon: CheckSquare },
        { id: "conduite-planning", label: "Conduite", icon: Calendar },
      ];

    case "moniteur":
      return [
        { id: "overview", label: "Jour", icon: Calendar },
        { id: "conduite-planning", label: "Planning", icon: Car },
        { id: "livret", label: "Émarger", icon: FileSignature },
        { id: "quiz-forms", label: "Quiz", icon: Brain },
      ];

    case "eleve":
      return [
        { id: "overview", label: "Accueil", icon: LayoutGrid },
        { id: "cours-online", label: "Cours", icon: Video },
        { id: "quiz", label: "Quiz", icon: Brain },
        { id: "reservation", label: "Conduite", icon: CalendarPlus },
      ];

    default:
      return [];
  }
}

export interface NetworkAlertItem {
  id: string;
  titre: string;
  message: string;
  severite: "info" | "warning" | "critical";
  cibleType: "all" | "single";
  cibleSchoolId?: string;
  cibleSchoolName?: string;
  bloquante: boolean;
  dateEnvoi: string;
  estActive: boolean;
}


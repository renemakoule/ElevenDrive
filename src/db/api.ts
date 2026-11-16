// ==============================================================================
// AutoSchool - Passerelle de Données (Dual-Engine : Tauri SQLite & Go Cloud API)
// ==============================================================================

import { guichetApi } from "./guichetApi";
import { candidatsImportApi } from "./candidatsImportApi";
import { flotteConduiteApi } from "./flotteConduiteApi";
import { examensApi } from "./examensApi";
import { pedagogieContentApi } from "./pedagogieContentApi";
import { syncApi } from "./syncApi";

export const api = {
  ...guichetApi,
  ...candidatsImportApi,
  ...flotteConduiteApi,
  ...examensApi,
  ...pedagogieContentApi,
  ...syncApi,
};

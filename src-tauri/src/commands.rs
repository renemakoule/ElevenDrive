// ==============================================================================
// Tauri IPC Commands - Façade modulaire
// ==============================================================================

pub mod system_sync;
pub mod sync_pull;
pub mod sync_pull_extra;
pub mod candidats;
pub mod candidats_import;
pub mod caisse;
pub mod gerant_derogations;
pub mod formules_flotte_examens;
pub mod pedagogie_seances;
pub mod pedagogie_medias;

pub use system_sync::*;
pub use sync_pull::*;
pub use candidats::*;
pub use candidats_import::*;
pub use caisse::*;
pub use gerant_derogations::*;
pub use formules_flotte_examens::*;
pub use pedagogie_seances::*;
pub use pedagogie_medias::*;

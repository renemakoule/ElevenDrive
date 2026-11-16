pub mod models;
pub mod db;
pub mod commands;
pub mod escpos;
pub mod session_guard;

use db::Database;
use std::path::PathBuf;
use tauri::Manager;

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_process::init())
        .plugin(tauri_plugin_updater::Builder::new().build())
        .setup(|app| {
            // Localisation du fichier de base de données SQLite local
            let app_data_dir = app.path().app_data_dir().unwrap_or_else(|_| PathBuf::from("."));
            std::fs::create_dir_all(&app_data_dir).ok();
            let db_path = app_data_dir.join("autoecole_local.db");

            let database = Database::init(db_path)
                .expect("Impossible d'initialiser la base de données SQLite locale");

            app.manage(database);
            app.manage(session_guard::LocalSessionState::default());
            Ok(())
        })
        .invoke_handler(tauri::generate_handler![
            commands::get_guichet_stats,
            commands::search_candidats,
            commands::get_candidat_dossier,
            commands::get_formules,
            commands::enregistrer_paiement,
            commands::update_candidat_photo,
            commands::get_session_caisse_detail,
            commands::reimprimer_recu,
            commands::ouvrir_session_caisse,
            commands::verifier_code_manager_derogation,
            commands::enregistrer_seance_code_avec_presences,
            commands::get_impayes_et_alertes_gerant,
            commands::get_evenements_en_attente,
            commands::marquer_evenements_synchronises,
            commands::synchroniser_evenements_en_attente,
            commands::appliquer_snapshot_cloud,
            commands::get_epreuves,
            commands::creer_epreuve,
            commands::supprimer_epreuve,
            commands::incrementer_telechargement_epreuve,
            commands::get_seances_code,
            commands::creer_seance_code,
            commands::reprogrammer_seance_code,
            commands::get_moniteurs,
            commands::get_tutoriels_videos,
            commands::creer_tutoriel_video,
            commands::supprimer_tutoriel_video,
            commands::incrementer_vue_tutoriel,
            commands::get_supports_cours,
            commands::creer_support_cours,
            commands::supprimer_support_cours,
            commands::incrementer_telechargement_support,
            commands::creer_formule,
            commands::modifier_prix_formule,
            commands::supprimer_formule,
            commands::supprimer_toutes_formules,
            commands::creer_candidat_avec_dossier,
            commands::importer_candidats_batch,
            commands::cloturer_session_caisse,
            commands::accorder_derogation,
            commands::get_vehicules,
            commands::creer_vehicule,
            commands::update_vehicule_statut,
            commands::forcer_candidat_examen,
            commands::window_minimize,
            commands::window_toggle_maximize,
            commands::window_close,
            commands::window_start_drag,
            commands::ouvrir_url_externe,
            commands::sauvegarder_pdf_telechargements,
            session_guard::set_local_session_context,
            session_guard::clear_local_session_context,
        ])
        .run(tauri::generate_context!())
        .expect("Erreur lors de l'exécution de l'application Tauri");
}

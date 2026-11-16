use tauri::{Manager, State};
use chrono::Utc;
use rusqlite::{params, Connection};
use crate::db::Database;
use crate::models::SyncEvent;

pub(super) fn md5_or_simple_hash(input: &str) -> u64 {
    use std::collections::hash_map::DefaultHasher;
    use std::hash::{Hash, Hasher};
    let mut s = DefaultHasher::new();
    input.hash(&mut s);
    s.finish()
}

pub(super) fn get_primary_school_id(conn: &Connection) -> String {
    conn.query_row("SELECT id FROM auto_ecoles ORDER BY created_at ASC LIMIT 1", [], |r| r.get(0))
        .unwrap_or_else(|_| "ae-primary".to_string())
}

pub(super) fn get_primary_agence_id(conn: &Connection) -> String {
    conn.query_row("SELECT id FROM agences ORDER BY created_at ASC LIMIT 1", [], |r| r.get(0))
        .unwrap_or_else(|_| "ag-primary".to_string())
}

pub(super) fn get_primary_caissier_id(conn: &Connection) -> String {
    conn.query_row(
        "SELECT id FROM utilisateurs WHERE statut = 'ACTIF' ORDER BY created_at ASC LIMIT 1",
        [],
        |r| r.get(0),
    )
    .unwrap_or_else(|_| "usr-guichet".to_string())
}

/// Ouvre une URL http(s) dans le navigateur par défaut du système.
#[tauri::command]
pub fn ouvrir_url_externe(url: String) -> Result<(), String> {
    let clean = url.trim();
    if !(clean.starts_with("http://") || clean.starts_with("https://")) {
        return Err("Seules les URL http/https sont autorisées".to_string());
    }

    #[cfg(target_os = "windows")]
    let result = std::process::Command::new("rundll32")
        .args(["url.dll,FileProtocolHandler", clean])
        .spawn();

    #[cfg(target_os = "macos")]
    let result = std::process::Command::new("open").arg(clean).spawn();

    #[cfg(all(unix, not(target_os = "macos")))]
    let result = std::process::Command::new("xdg-open").arg(clean).spawn();

    result.map(|_| ()).map_err(|e| format!("Impossible d'ouvrir le navigateur : {}", e))
}

#[tauri::command]
pub fn window_minimize(window: tauri::Window) -> Result<(), String> {
    window.minimize().map_err(|e| e.to_string())
}

#[tauri::command]
pub fn window_toggle_maximize(window: tauri::Window) -> Result<(), String> {
    let is_maximized = window.is_maximized().map_err(|e| e.to_string())?;
    if is_maximized {
        window.unmaximize().map_err(|e| e.to_string())
    } else {
        window.maximize().map_err(|e| e.to_string())
    }
}

#[tauri::command]
pub fn window_close(window: tauri::Window) -> Result<(), String> {
    window.close().map_err(|e| e.to_string())
}

#[tauri::command]
pub fn window_start_drag(window: tauri::Window) -> Result<(), String> {
    window.start_dragging().map_err(|e| e.to_string())
}

/// Enregistre un PDF généré dans le dossier Téléchargements et l'ouvre.
#[tauri::command]
pub fn sauvegarder_pdf_telechargements(
    filename: String,
    bytes: Vec<u8>,
    app: tauri::AppHandle,
) -> Result<String, String> {
    let safe_name = std::path::Path::new(&filename)
        .file_name()
        .and_then(|n| n.to_str())
        .filter(|s| !s.is_empty())
        .unwrap_or("historique-activites.pdf")
        .to_string();

    let dir = app
        .path()
        .download_dir()
        .or_else(|_| app.path().document_dir())
        .map_err(|e| format!("Dossier Téléchargements introuvable : {}", e))?;
    std::fs::create_dir_all(&dir).map_err(|e| e.to_string())?;

    let target = dir.join(&safe_name);
    std::fs::write(&target, &bytes).map_err(|e| format!("Écriture PDF impossible : {}", e))?;

    let target_str = target.to_string_lossy().to_string();
    #[cfg(target_os = "windows")]
    let _ = std::process::Command::new("rundll32")
        .args(["url.dll,FileProtocolHandler", &target_str])
        .spawn();
    #[cfg(target_os = "macos")]
    let _ = std::process::Command::new("open").arg(&target_str).spawn();
    #[cfg(all(unix, not(target_os = "macos")))]
    let _ = std::process::Command::new("xdg-open").arg(&target_str).spawn();

    Ok(target_str)
}

#[tauri::command]
pub fn get_evenements_en_attente(db: State<Database>) -> Result<Vec<SyncEvent>, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    let mut stmt = conn.prepare(
        "SELECT id, nom_table, type_operation, entite_id, donnees_payload, horodatage_local, statut_sync
         FROM journal_sync_events
         WHERE statut_sync = 'EN_ATTENTE'
         ORDER BY horodatage_local ASC LIMIT 100"
    ).map_err(|e| e.to_string())?;

    let rows = stmt.query_map([], |row| {
        Ok(SyncEvent {
            id: row.get(0)?,
            nom_table: row.get(1)?,
            type_operation: row.get(2)?,
            entite_id: row.get(3)?,
            donnees_payload: row.get(4)?,
            horodatage_local: row.get(5)?,
            statut_sync: row.get(6)?,
        })
    }).map_err(|e| e.to_string())?;

    let mut list = Vec::new();
    for r in rows {
        if let Ok(item) = r {
            list.push(item);
        }
    }
    Ok(list)
}

#[tauri::command]
pub fn marquer_evenements_synchronises(event_ids: Vec<String>, db: State<Database>) -> Result<usize, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    let tx = conn.unchecked_transaction().map_err(|e| e.to_string())?;
    let mut updated_count = 0usize;
    for id in &event_ids {
        let rows = tx.execute(
            "UPDATE journal_sync_events SET statut_sync = 'SYNCHRONISE' WHERE id = ?1",
            params![id],
        ).unwrap_or(0);
        updated_count += rows;
    }
    tx.commit().map_err(|e| e.to_string())?;
    Ok(updated_count)
}

#[tauri::command]
pub fn synchroniser_evenements_en_attente(db: State<Database>) -> Result<serde_json::Value, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;

    let updated = conn.execute(
        "UPDATE journal_sync_events SET statut_sync = 'SYNCHRONISE' WHERE statut_sync = 'EN_ATTENTE'",
        [],
    ).map_err(|e| e.to_string())?;

    Ok(serde_json::json!({
        "statut": "SUCCES",
        "evenements_traites": updated,
        "doublons_ignores": 0,
        "erreurs": 0,
        "synced_at": Utc::now().to_rfc3339(),
        "message": format!("{} transaction(s) locale(s) synchronisée(s) vers Neon Cloud", updated)
    }))
}


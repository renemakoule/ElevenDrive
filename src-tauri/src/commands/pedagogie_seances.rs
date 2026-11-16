use tauri::State;
use uuid::Uuid;
use chrono::Utc;
use rusqlite::params;
use crate::db::Database;
use crate::models::{CreerSeanceCodePayload, SeanceCode};
use super::system_sync::{get_primary_agence_id, get_primary_school_id};

#[tauri::command]
pub fn enregistrer_seance_code_avec_presences(
    theme_traite: String,
    formateur_nom: String,
    heure_debut: String,
    heure_fin: String,
    presences: Vec<serde_json::Value>,
    db: State<Database>,
) -> Result<serde_json::Value, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;

    let seance_id = Uuid::new_v4().to_string();
    let today_date = Utc::now().format("%Y-%m-%d").to_string();
    let now_str = Utc::now().to_rfc3339();
    let auto_ecole_id = get_primary_school_id(&conn);
    let agence_id = get_primary_agence_id(&conn);

    conn.execute(
        "INSERT INTO seances_code (id, auto_ecole_id, agence_id, date_seance, heure_debut, heure_fin, theme_traite, formateur_nom, statut)
         VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, 'TERMINEE')",
        params![seance_id, auto_ecole_id, agence_id, today_date, heure_debut, heure_fin, theme_traite, formateur_nom],
    ).map_err(|e| e.to_string())?;

    let mut count_presents = 0;

    for item in &presences {
        let candidat_id = item["candidat_id"].as_str().unwrap_or_default();
        let est_present = item["est_present"].as_bool().unwrap_or(false);
        let score_test = item["score_test"].as_i64().unwrap_or(0);

        if candidat_id.is_empty() {
            continue;
        }

        let dossier_id: Option<String> = conn
            .query_row(
                "SELECT id FROM dossiers_inscription WHERE candidat_id = ?1 ORDER BY date_inscription DESC LIMIT 1",
                params![candidat_id],
                |r| r.get(0),
            )
            .ok();

        let presence_id = Uuid::new_v4().to_string();
        conn.execute(
            "INSERT INTO presences_code (id, seance_id, candidat_id, dossier_id, heure_pointage, score_test, est_present)
             VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7)",
            params![
                presence_id,
                seance_id,
                candidat_id,
                dossier_id.clone().unwrap_or_default(),
                now_str,
                score_test,
                if est_present { 1 } else { 0 }
            ],
        ).map_err(|e| e.to_string())?;

        if est_present {
            count_presents += 1;
            if let Some(ref d_id) = dossier_id {
                conn.execute(
                    "UPDATE dossiers_inscription SET heures_code_effectuees = heures_code_effectuees + 2 WHERE id = ?1",
                    params![d_id],
                ).ok();
            }
        }
    }

    let sync_event_id = Uuid::new_v4().to_string();
    let payload = serde_json::json!({
        "seance_id": seance_id,
        "date_seance": today_date,
        "heure_debut": heure_debut,
        "heure_fin": heure_fin,
        "theme_traite": theme_traite,
        "formateur_nom": formateur_nom,
        "presents": count_presents,
        "statut": "TERMINEE"
    }).to_string();

    conn.execute(
        "INSERT INTO journal_sync_events (id, agence_id, nom_table, type_operation, entite_id, donnees_payload, horodatage_local, statut_sync)
         VALUES (?1, ?2, 'seances_code', 'INSERT', ?3, ?4, ?5, 'EN_ATTENTE')",
        params![sync_event_id, agence_id, seance_id, payload, now_str],
    ).ok();

    Ok(serde_json::json!({
        "seance_id": seance_id,
        "total_presents": count_presents,
        "date_seance": today_date,
        "statut": "SUCCES"
    }))
}

#[tauri::command]
pub fn get_seances_code(
    filtre_type: Option<String>,
    db: State<Database>,
) -> Result<Vec<SeanceCode>, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;

    let _ = conn.execute(
        "UPDATE seances_code
         SET statut_seance = 'TERMINEE', statut = 'TERMINEE'
         WHERE COALESCE(statut_seance, 'PLANIFIEE') NOT IN ('TERMINEE', 'ANNULEE')
           AND (
             date_seance < date('now', 'localtime')
             OR (date_seance = date('now', 'localtime') AND COALESCE(NULLIF(heure_fin, ''), '23:59') <= strftime('%H:%M', 'now', 'localtime'))
           )",
        [],
    );

    let mut query = "SELECT id, auto_ecole_id, agence_id, date_seance, heure_debut, heure_fin, 
                            theme_traite, formateur_nom, moniteur_id, 
                            COALESCE(type_seance, 'PRESENTIEL'), lien_visio, mot_de_passe_visio, 
                            COALESCE(statut_seance, 'PLANIFIEE')
                     FROM seances_code".to_string();

    let has_filter = match &filtre_type {
        Some(t) => !t.trim().is_empty() && t != "ALL",
        None => false,
    };

    if has_filter {
        query.push_str(" WHERE type_seance = ?1");
    }
    query.push_str(" ORDER BY date_seance DESC, heure_debut ASC LIMIT 50");

    let mut stmt = conn.prepare(&query).map_err(|e| e.to_string())?;

    let map_row = |row: &rusqlite::Row| {
        Ok(SeanceCode {
            id: row.get(0)?,
            auto_ecole_id: row.get(1)?,
            agence_id: row.get(2)?,
            date_seance: row.get(3)?,
            heure_debut: row.get(4)?,
            heure_fin: row.get(5)?,
            theme_traite: row.get(6)?,
            formateur_nom: row.get(7)?,
            moniteur_id: row.get(8)?,
            type_seance: row.get(9)?,
            lien_visio: row.get(10)?,
            mot_de_passe_visio: row.get(11)?,
            statut_seance: row.get(12)?,
        })
    };

    if has_filter {
        let t = filtre_type.unwrap();
        let rows = stmt.query_map(params![t], map_row).map_err(|e| e.to_string())?;
        let mut list = Vec::new();
        for r in rows {
            if let Ok(item) = r { list.push(item); }
        }
        Ok(list)
    } else {
        let rows = stmt.query_map([], map_row).map_err(|e| e.to_string())?;
        let mut list = Vec::new();
        for r in rows {
            if let Ok(item) = r { list.push(item); }
        }
        Ok(list)
    }
}

#[tauri::command]
pub fn creer_seance_code(
    payload: CreerSeanceCodePayload,
    db: State<Database>,
) -> Result<SeanceCode, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    let seance_id = Uuid::new_v4().to_string();
    let auto_ecole_id = get_primary_school_id(&conn);
    let agence_id = get_primary_agence_id(&conn);

    conn.execute(
        "INSERT INTO seances_code (
            id, auto_ecole_id, agence_id, date_seance, heure_debut, heure_fin, 
            theme_traite, formateur_nom, moniteur_id, type_seance, lien_visio, 
            mot_de_passe_visio, statut_seance, statut
        ) VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9, ?10, ?11, ?12, 'PLANIFIEE', 'PLANIFIEE')",
        params![
            seance_id,
            auto_ecole_id,
            agence_id,
            payload.date_seance,
            payload.heure_debut,
            payload.heure_fin,
            payload.theme_traite,
            payload.formateur_nom,
            payload.moniteur_id,
            payload.type_seance,
            payload.lien_visio,
            payload.mot_de_passe_visio
        ],
    ).map_err(|e| e.to_string())?;

    let sync_id = Uuid::new_v4().to_string();
    let now_str = Utc::now().to_rfc3339();
    let sync_payload = serde_json::json!({
        "seance_id": seance_id,
        "date_seance": &payload.date_seance,
        "heure_debut": &payload.heure_debut,
        "heure_fin": &payload.heure_fin,
        "theme_traite": &payload.theme_traite,
        "formateur_nom": &payload.formateur_nom,
        "moniteur_id": &payload.moniteur_id,
        "type_seance": &payload.type_seance,
        "lien_visio": &payload.lien_visio,
        "mot_de_passe_visio": &payload.mot_de_passe_visio,
        "statut": "PLANIFIEE"
    }).to_string();
    let _ = conn.execute(
        "INSERT INTO journal_sync_events (id, agence_id, nom_table, type_operation, entite_id, donnees_payload, horodatage_local, statut_sync)
         VALUES (?1, ?2, 'seances_code', 'INSERT', ?3, ?4, ?5, 'EN_ATTENTE')",
        params![sync_id, &agence_id, &seance_id, sync_payload, now_str],
    );

    Ok(SeanceCode {
        id: seance_id,
        auto_ecole_id,
        agence_id,
        date_seance: payload.date_seance,
        heure_debut: payload.heure_debut,
        heure_fin: payload.heure_fin,
        theme_traite: payload.theme_traite,
        formateur_nom: payload.formateur_nom,
        moniteur_id: payload.moniteur_id,
        type_seance: payload.type_seance,
        lien_visio: payload.lien_visio,
        mot_de_passe_visio: payload.mot_de_passe_visio,
        statut_seance: "PLANIFIEE".to_string(),
    })
}

#[tauri::command]
pub fn reprogrammer_seance_code(
    id: String,
    date_seance: String,
    heure_debut: String,
    heure_fin: String,
    db: State<Database>,
    session: State<crate::session_guard::LocalSessionState>,
) -> Result<SeanceCode, String> {
    crate::session_guard::require_local_roles(&session, &["DIRECTEUR", "SECRETAIRE"])?;
    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    let updated = conn.execute(
        "UPDATE seances_code
         SET date_seance = ?1, heure_debut = ?2, heure_fin = ?3, statut_seance = 'PLANIFIEE', statut = 'PLANIFIEE'
         WHERE id = ?4",
        params![&date_seance, &heure_debut, &heure_fin, &id],
    ).map_err(|e| e.to_string())?;

    if updated == 0 {
        return Err("Séance introuvable en base locale".to_string());
    }

    let agence_id = get_primary_agence_id(&conn);
    let sync_id = Uuid::new_v4().to_string();
    let now_str = Utc::now().to_rfc3339();
    let sync_payload = serde_json::json!({
        "seance_id": &id,
        "date_seance": &date_seance,
        "heure_debut": &heure_debut,
        "heure_fin": &heure_fin
    }).to_string();
    let _ = conn.execute(
        "INSERT INTO journal_sync_events (id, agence_id, nom_table, type_operation, entite_id, donnees_payload, horodatage_local, statut_sync)
         VALUES (?1, ?2, 'seances_code', 'UPDATE', ?3, ?4, ?5, 'EN_ATTENTE')",
        params![sync_id, &agence_id, &id, sync_payload, now_str],
    );

    conn.query_row(
        "SELECT id, auto_ecole_id, agence_id, date_seance, heure_debut, heure_fin,
                theme_traite, formateur_nom, moniteur_id,
                COALESCE(type_seance, 'PRESENTIEL'), lien_visio, mot_de_passe_visio,
                COALESCE(statut_seance, 'PLANIFIEE')
         FROM seances_code WHERE id = ?1",
        params![&id],
        |row| {
            Ok(SeanceCode {
                id: row.get(0)?,
                auto_ecole_id: row.get(1)?,
                agence_id: row.get(2)?,
                date_seance: row.get(3)?,
                heure_debut: row.get(4)?,
                heure_fin: row.get(5)?,
                theme_traite: row.get(6)?,
                formateur_nom: row.get(7)?,
                moniteur_id: row.get(8)?,
                type_seance: row.get(9)?,
                lien_visio: row.get(10)?,
                mot_de_passe_visio: row.get(11)?,
                statut_seance: row.get(12)?,
            })
        },
    ).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn get_moniteurs(db: State<Database>) -> Result<Vec<serde_json::Value>, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    let mut stmt = conn.prepare(
        "SELECT id, nom, prenoms, telephone, role FROM utilisateurs WHERE role = 'MONITEUR' OR role = 'DIRECTEUR' ORDER BY nom ASC"
    ).map_err(|e| e.to_string())?;

    let rows = stmt.query_map([], |row| {
        Ok(serde_json::json!({
            "id": row.get::<_, String>(0)?,
            "nom": row.get::<_, String>(1)?,
            "prenoms": row.get::<_, String>(2)?,
            "telephone": row.get::<_, String>(3)?,
            "role": row.get::<_, String>(4)?,
        }))
    }).map_err(|e| e.to_string())?;

    let mut list = Vec::new();
    for r in rows {
        if let Ok(item) = r { list.push(item); }
    }
    Ok(list)
}


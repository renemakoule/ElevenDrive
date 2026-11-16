use tauri::State;
use uuid::Uuid;
use chrono::Utc;
use rusqlite::params;
use crate::db::Database;
use crate::models::*;

#[tauri::command]
pub fn get_epreuves(categorie: Option<String>, db: State<Database>) -> Result<Vec<EpreuvePdf>, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    let filter_active = matches!(&categorie, Some(cat) if !cat.trim().is_empty() && cat != "ALL");

    let mut query = "SELECT id, titre, categorie_permis, description, annee, session, duree_minutes, 
                            nombre_questions, avec_corrige, fichier_nom, fichier_taille, fichier_data, 
                            ajoute_par, telechargements, created_at 
                     FROM epreuves_pdf".to_string();

    if filter_active {
        query.push_str(" WHERE categorie_permis = ?1");
    }
    query.push_str(" ORDER BY annee DESC, created_at DESC");

    let mut stmt = conn.prepare(&query).map_err(|e| e.to_string())?;
    let map_row = |row: &rusqlite::Row| {
        Ok(EpreuvePdf {
            id: row.get(0)?,
            titre: row.get(1)?,
            categorie_permis: row.get(2)?,
            description: row.get(3)?,
            annee: row.get(4)?,
            session: row.get(5)?,
            duree_minutes: row.get(6)?,
            nombre_questions: row.get(7)?,
            avec_corrige: row.get::<_, i32>(8)? == 1,
            fichier_nom: row.get(9)?,
            fichier_taille: row.get(10)?,
            fichier_data: row.get(11)?,
            ajoute_par: row.get(12)?,
            telechargements_count: row.get(13)?,
            created_at: row.get(14)?,
        })
    };

    let mut list = Vec::new();
    if filter_active {
        let rows = stmt.query_map(params![categorie.unwrap()], map_row).map_err(|e| e.to_string())?;
        for r in rows { if let Ok(item) = r { list.push(item); } }
    } else {
        let rows = stmt.query_map([], map_row).map_err(|e| e.to_string())?;
        for r in rows { if let Ok(item) = r { list.push(item); } }
    }
    Ok(list)
}

#[tauri::command]
pub fn creer_epreuve(payload: CreerEpreuvePayload, db: State<Database>) -> Result<EpreuvePdf, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    let new_id = Uuid::new_v4().to_string();
    let created_at = Utc::now().to_rfc3339();

    conn.execute(
        "INSERT INTO epreuves_pdf (
            id, titre, categorie_permis, description, annee, session, duree_minutes, 
            nombre_questions, avec_corrige, fichier_nom, fichier_taille, fichier_data, 
            ajoute_par, telechargements, created_at
        ) VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9, ?10, ?11, ?12, ?13, 0, ?14)",
        params![
            new_id, payload.titre, payload.categorie_permis, payload.description,
            payload.annee, payload.session, payload.duree_minutes, payload.nombre_questions,
            if payload.avec_corrige { 1 } else { 0 },
            payload.fichier_nom, payload.fichier_taille, payload.fichier_data,
            payload.ajoute_par, created_at
        ],
    ).map_err(|e| e.to_string())?;

    Ok(EpreuvePdf {
        id: new_id,
        titre: payload.titre,
        categorie_permis: payload.categorie_permis,
        description: payload.description,
        annee: payload.annee,
        session: payload.session,
        duree_minutes: payload.duree_minutes,
        nombre_questions: payload.nombre_questions,
        avec_corrige: payload.avec_corrige,
        fichier_nom: payload.fichier_nom,
        fichier_taille: payload.fichier_taille,
        fichier_data: payload.fichier_data,
        ajoute_par: payload.ajoute_par,
        telechargements_count: 0,
        created_at,
    })
}

#[tauri::command]
pub fn supprimer_epreuve(id: String, db: State<Database>) -> Result<bool, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    let rows_affected = conn.execute("DELETE FROM epreuves_pdf WHERE id = ?1", params![id])
        .map_err(|e| e.to_string())?;
    Ok(rows_affected > 0)
}

#[tauri::command]
pub fn incrementer_telechargement_epreuve(id: String, db: State<Database>) -> Result<bool, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    let rows_affected = conn.execute(
        "UPDATE epreuves_pdf SET telechargements = telechargements + 1 WHERE id = ?1",
        params![id],
    ).map_err(|e| e.to_string())?;
    Ok(rows_affected > 0)
}

#[tauri::command]
pub fn get_tutoriels_videos(
    categorie: Option<String>,
    db: State<Database>,
) -> Result<Vec<TutorielVideo>, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    let filter_active = categorie.is_some() && categorie.as_deref() != Some("ALL") && categorie.as_deref() != Some("");

    let mut query = "SELECT id, titre, categorie_permis, description, url_video, plateforme, duree_minutes, miniature_url, vues_count, ajoute_par, created_at FROM tutoriels_video".to_string();
    if filter_active {
        query.push_str(" WHERE categorie_permis = ?1 OR categorie_permis = 'TOUTES'");
    }
    query.push_str(" ORDER BY created_at DESC");

    let mut stmt = conn.prepare(&query).map_err(|e| e.to_string())?;
    let map_row = |row: &rusqlite::Row| {
        Ok(TutorielVideo {
            id: row.get(0)?,
            titre: row.get(1)?,
            categorie_permis: row.get(2)?,
            description: row.get(3)?,
            url_video: row.get(4)?,
            plateforme: row.get(5)?,
            duree_minutes: row.get(6)?,
            miniature_url: row.get(7)?,
            vues_count: row.get(8)?,
            ajoute_par: row.get(9)?,
            created_at: row.get(10)?,
        })
    };

    let mut list = Vec::new();
    if filter_active {
        let rows = stmt.query_map(params![categorie.unwrap()], map_row).map_err(|e| e.to_string())?;
        for r in rows { if let Ok(item) = r { list.push(item); } }
    } else {
        let rows = stmt.query_map([], map_row).map_err(|e| e.to_string())?;
        for r in rows { if let Ok(item) = r { list.push(item); } }
    }
    Ok(list)
}

#[tauri::command]
pub fn creer_tutoriel_video(
    payload: CreerTutorielVideoPayload,
    db: State<Database>,
) -> Result<TutorielVideo, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    let new_id = Uuid::new_v4().to_string();
    let created_at = Utc::now().to_rfc3339();
    let plateforme = payload.plateforme.unwrap_or_else(|| "YOUTUBE".to_string());
    let duree_minutes = payload.duree_minutes.unwrap_or(5);

    conn.execute(
        "INSERT INTO tutoriels_video (
            id, titre, categorie_permis, description, url_video, plateforme, 
            duree_minutes, miniature_url, vues_count, ajoute_par, created_at
        ) VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, 0, ?9, ?10)",
        params![
            &new_id, &payload.titre, &payload.categorie_permis, &payload.description,
            &payload.url_video, &plateforme, duree_minutes, &payload.miniature_url,
            &payload.ajoute_par, &created_at
        ],
    ).map_err(|e| e.to_string())?;

    Ok(TutorielVideo {
        id: new_id,
        titre: payload.titre,
        categorie_permis: payload.categorie_permis,
        description: payload.description,
        url_video: payload.url_video,
        plateforme: Some(plateforme),
        duree_minutes: Some(duree_minutes),
        miniature_url: payload.miniature_url,
        vues_count: Some(0),
        ajoute_par: payload.ajoute_par,
        created_at: Some(created_at),
    })
}

#[tauri::command]
pub fn supprimer_tutoriel_video(id: String, db: State<Database>) -> Result<bool, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    conn.execute("DELETE FROM tutoriels_video WHERE id = ?1", params![id]).map_err(|e| e.to_string())?;
    Ok(true)
}

#[tauri::command]
pub fn incrementer_vue_tutoriel(id: String, db: State<Database>) -> Result<bool, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    conn.execute("UPDATE tutoriels_video SET vues_count = COALESCE(vues_count, 0) + 1 WHERE id = ?1", params![id]).map_err(|e| e.to_string())?;
    Ok(true)
}

#[tauri::command]
pub fn get_supports_cours(
    categorie: Option<String>,
    db: State<Database>,
) -> Result<Vec<SupportCours>, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    let filter_active = categorie.is_some() && categorie.as_deref() != Some("ALL") && categorie.as_deref() != Some("");

    let mut query = "SELECT id, titre, categorie_permis, description, type_support, fichier_nom, fichier_taille, fichier_data, url_externe, ajoute_par, telechargements, created_at FROM supports_cours".to_string();
    if filter_active {
        query.push_str(" WHERE categorie_permis = ?1 OR categorie_permis = 'TOUTES'");
    }
    query.push_str(" ORDER BY created_at DESC");

    let mut stmt = conn.prepare(&query).map_err(|e| e.to_string())?;
    let map_row = |row: &rusqlite::Row| {
        Ok(SupportCours {
            id: row.get(0)?,
            titre: row.get(1)?,
            categorie_permis: row.get(2)?,
            description: row.get(3)?,
            type_support: row.get(4)?,
            fichier_nom: row.get(5)?,
            fichier_taille: row.get(6)?,
            fichier_data: row.get(7)?,
            url_externe: row.get(8)?,
            ajoute_par: row.get(9)?,
            telechargements_count: row.get(10)?,
            created_at: row.get(11)?,
        })
    };

    let mut list = Vec::new();
    if filter_active {
        let rows = stmt.query_map(params![categorie.unwrap()], map_row).map_err(|e| e.to_string())?;
        for r in rows { if let Ok(item) = r { list.push(item); } }
    } else {
        let rows = stmt.query_map([], map_row).map_err(|e| e.to_string())?;
        for r in rows { if let Ok(item) = r { list.push(item); } }
    }
    Ok(list)
}

#[tauri::command]
pub fn creer_support_cours(
    payload: CreerSupportCoursPayload,
    db: State<Database>,
) -> Result<SupportCours, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    let new_id = Uuid::new_v4().to_string();
    let created_at = Utc::now().to_rfc3339();

    conn.execute(
        "INSERT INTO supports_cours (
            id, titre, categorie_permis, description, type_support, fichier_nom, 
            fichier_taille, fichier_data, url_externe, ajoute_par, telechargements, created_at
        ) VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9, ?10, 0, ?11)",
        params![
            &new_id, &payload.titre, &payload.categorie_permis, &payload.description,
            &payload.type_support, &payload.fichier_nom, &payload.fichier_taille,
            &payload.fichier_data, &payload.url_externe, &payload.ajoute_par, &created_at
        ],
    ).map_err(|e| e.to_string())?;

    Ok(SupportCours {
        id: new_id,
        titre: payload.titre,
        categorie_permis: payload.categorie_permis,
        description: payload.description,
        type_support: payload.type_support,
        fichier_nom: payload.fichier_nom,
        fichier_taille: payload.fichier_taille,
        fichier_data: payload.fichier_data,
        url_externe: payload.url_externe,
        ajoute_par: payload.ajoute_par,
        telechargements_count: Some(0),
        created_at: Some(created_at),
    })
}

#[tauri::command]
pub fn supprimer_support_cours(id: String, db: State<Database>) -> Result<bool, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    conn.execute("DELETE FROM supports_cours WHERE id = ?1", params![id]).map_err(|e| e.to_string())?;
    Ok(true)
}

#[tauri::command]
pub fn incrementer_telechargement_support(id: String, db: State<Database>) -> Result<bool, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    conn.execute("UPDATE supports_cours SET telechargements = COALESCE(telechargements, 0) + 1 WHERE id = ?1", params![id]).map_err(|e| e.to_string())?;
    Ok(true)
}


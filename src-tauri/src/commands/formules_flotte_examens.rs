use tauri::State;
use uuid::Uuid;
use chrono::Utc;
use rusqlite::params;
use crate::db::Database;
use crate::models::*;
use super::system_sync::{get_primary_agence_id, get_primary_school_id};

#[tauri::command]
pub fn get_formules(db: State<Database>) -> Result<Vec<FormulePermis>, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    let mut stmt = conn
        .prepare("SELECT id, categorie_permis, intitule, heures_theoriques, heures_pratiques, prix_total, acompte_minimum, est_actif, tranches_paiement FROM formules_permis WHERE est_actif = 1")
        .map_err(|e| e.to_string())?;

    let rows = stmt.query_map([], |row| {
        let tranches_raw: Option<String> = row.get(8).unwrap_or(None);
        let tranches_val = tranches_raw.and_then(|s| serde_json::from_str(&s).ok());
        Ok(FormulePermis {
            id: row.get(0)?,
            categorie_permis: row.get(1)?,
            intitule: row.get(2)?,
            heures_theoriques: row.get(3)?,
            heures_pratiques: row.get(4)?,
            prix_total: row.get(5)?,
            acompte_minimum: row.get(6)?,
            est_actif: row.get::<_, i32>(7)? == 1,
            tranches_paiement: tranches_val,
        })
    }).map_err(|e| e.to_string())?;

    let mut list = Vec::new();
    for r in rows {
        if let Ok(f) = r {
            list.push(f);
        }
    }
    Ok(list)
}

#[tauri::command]
pub fn creer_formule(payload: CreerFormulePayload, db: State<Database>) -> Result<FormulePermis, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    let new_id = Uuid::new_v4().to_string();
    let auto_ecole_id = get_primary_school_id(&conn);
    let _ = conn.execute(
        "INSERT OR IGNORE INTO auto_ecoles (id, code_unique, nom_commercial, numero_agrement, telephone)
         VALUES (?1, ?1, 'Auto-École Locale', 'LOCAL', '00000000')",
        params![&auto_ecole_id],
    );
    let tranches_str = payload.tranches_paiement.as_ref().map(|v| v.to_string());

    conn.execute(
        "INSERT INTO formules_permis (id, auto_ecole_id, categorie_permis, intitule, description, heures_theoriques, heures_pratiques, prix_total, acompte_minimum, est_actif, tranches_paiement)
         VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9, 1, ?10)",
        params![
            &new_id,
            &auto_ecole_id,
            &payload.categorie_permis,
            &payload.intitule,
            &payload.description,
            payload.heures_theoriques,
            payload.heures_pratiques,
            payload.prix_total,
            payload.acompte_minimum,
            &tranches_str
        ],
    ).map_err(|e| e.to_string())?;

    let agence_id = get_primary_agence_id(&conn);
    let sync_id = Uuid::new_v4().to_string();
    let sync_payload = serde_json::json!({
        "id": new_id,
        "categorie_permis": &payload.categorie_permis,
        "intitule": &payload.intitule,
        "description": &payload.description,
        "heures_theoriques": payload.heures_theoriques,
        "heures_pratiques": payload.heures_pratiques,
        "prix_total": payload.prix_total,
        "acompte_minimum": payload.acompte_minimum
    }).to_string();
    let now_sync = chrono::Utc::now().to_rfc3339();
    let _ = conn.execute(
        "INSERT INTO journal_sync_events (id, agence_id, nom_table, type_operation, entite_id, donnees_payload, horodatage_local, statut_sync)
         VALUES (?1, ?2, 'formules_permis', 'INSERT', ?3, ?4, ?5, 'EN_ATTENTE')",
        params![sync_id, &agence_id, &new_id, sync_payload, now_sync],
    );

    Ok(FormulePermis {
        id: new_id,
        categorie_permis: payload.categorie_permis,
        intitule: payload.intitule,
        heures_theoriques: payload.heures_theoriques,
        heures_pratiques: payload.heures_pratiques,
        prix_total: payload.prix_total,
        acompte_minimum: payload.acompte_minimum,
        est_actif: true,
        tranches_paiement: payload.tranches_paiement,
    })
}

#[tauri::command]
pub fn modifier_prix_formule(
    payload: ModifierPrixFormulePayload,
    db: State<Database>,
    session: State<crate::session_guard::LocalSessionState>,
) -> Result<FormulePermis, String> {
    crate::session_guard::require_local_roles(&session, &["DIRECTEUR"])?;
    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    let tranches_str = payload.tranches_paiement.as_ref().map(|v| v.to_string());

    if let Some(ref tr) = tranches_str {
        if let Some(acompte) = payload.acompte_minimum {
            conn.execute(
                "UPDATE formules_permis SET prix_total = ?1, acompte_minimum = ?2, tranches_paiement = ?3 WHERE id = ?4",
                params![payload.nouveau_prix, acompte, tr, &payload.formule_id],
            ).map_err(|e| e.to_string())?;
        } else {
            conn.execute(
                "UPDATE formules_permis SET prix_total = ?1, tranches_paiement = ?2 WHERE id = ?3",
                params![payload.nouveau_prix, tr, &payload.formule_id],
            ).map_err(|e| e.to_string())?;
        }
    } else if let Some(acompte) = payload.acompte_minimum {
        conn.execute(
            "UPDATE formules_permis SET prix_total = ?1, acompte_minimum = ?2 WHERE id = ?3",
            params![payload.nouveau_prix, acompte, &payload.formule_id],
        ).map_err(|e| e.to_string())?;
    } else {
        conn.execute(
            "UPDATE formules_permis SET prix_total = ?1 WHERE id = ?2",
            params![payload.nouveau_prix, &payload.formule_id],
        ).map_err(|e| e.to_string())?;
    }

    let agence_id = get_primary_agence_id(&conn);
    let sync_id = Uuid::new_v4().to_string();
    let now_sync = chrono::Utc::now().to_rfc3339();
    let sync_payload = serde_json::json!({
        "id": &payload.formule_id,
        "prix_total": payload.nouveau_prix,
        "acompte_minimum": payload.acompte_minimum
    }).to_string();
    let _ = conn.execute(
        "INSERT INTO journal_sync_events (id, agence_id, nom_table, type_operation, entite_id, donnees_payload, horodatage_local, statut_sync)
         VALUES (?1, ?2, 'formules_permis', 'UPDATE', ?3, ?4, ?5, 'EN_ATTENTE')",
        params![sync_id, &agence_id, &payload.formule_id, sync_payload, now_sync],
    );

    conn.query_row(
        "SELECT id, categorie_permis, intitule, heures_theoriques, heures_pratiques, prix_total, acompte_minimum, est_actif, tranches_paiement
         FROM formules_permis WHERE id = ?1",
        params![&payload.formule_id],
        |row| {
            let tranches_raw: Option<String> = row.get(8).unwrap_or(None);
            let tranches_val = tranches_raw.and_then(|s| serde_json::from_str(&s).ok());
            Ok(FormulePermis {
                id: row.get(0)?,
                categorie_permis: row.get(1)?,
                intitule: row.get(2)?,
                heures_theoriques: row.get(3)?,
                heures_pratiques: row.get(4)?,
                prix_total: row.get(5)?,
                acompte_minimum: row.get(6)?,
                est_actif: row.get::<_, i32>(7)? == 1,
                tranches_paiement: tranches_val,
            })
        },
    ).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn supprimer_formule(
    formule_id: String,
    db: State<Database>,
    session: State<crate::session_guard::LocalSessionState>,
) -> Result<bool, String> {
    crate::session_guard::require_local_roles(&session, &["DIRECTEUR"])?;
    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    let dossiers_count: i64 = conn
        .query_row(
            "SELECT count(*) FROM dossiers_inscription WHERE formule_id = ?1",
            params![&formule_id],
            |r| r.get(0),
        )
        .unwrap_or(0);

    if dossiers_count > 0 {
        conn.execute(
            "UPDATE formules_permis SET est_actif = 0 WHERE id = ?1",
            params![&formule_id],
        ).map_err(|e| e.to_string())?;
    } else {
        conn.execute(
            "DELETE FROM formules_permis WHERE id = ?1",
            params![&formule_id],
        ).map_err(|e| e.to_string())?;
    }

    let agence_id = get_primary_agence_id(&conn);
    let sync_id = Uuid::new_v4().to_string();
    let now_sync = chrono::Utc::now().to_rfc3339();
    let (op, sync_payload) = if dossiers_count > 0 {
        ("UPDATE", serde_json::json!({ "id": &formule_id, "est_actif": false }).to_string())
    } else {
        ("DELETE", serde_json::json!({ "id": &formule_id }).to_string())
    };
    let _ = conn.execute(
        "INSERT INTO journal_sync_events (id, agence_id, nom_table, type_operation, entite_id, donnees_payload, horodatage_local, statut_sync)
         VALUES (?1, ?2, 'formules_permis', ?3, ?4, ?5, ?6, 'EN_ATTENTE')",
        params![sync_id, &agence_id, op, &formule_id, sync_payload, now_sync],
    );
    Ok(true)
}

#[tauri::command]
pub fn supprimer_toutes_formules(
    db: State<Database>,
    session: State<crate::session_guard::LocalSessionState>,
) -> Result<usize, String> {
    crate::session_guard::require_local_roles(&session, &["DIRECTEUR"])?;
    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    conn.execute(
        "UPDATE formules_permis SET est_actif = 0 WHERE id IN (SELECT DISTINCT formule_id FROM dossiers_inscription)",
        [],
    ).map_err(|e| e.to_string())?;

    let count = conn.execute(
        "DELETE FROM formules_permis WHERE id NOT IN (SELECT DISTINCT formule_id FROM dossiers_inscription)",
        [],
    ).map_err(|e| e.to_string())?;

    Ok(count)
}

#[tauri::command]
pub fn get_vehicules(db: State<Database>) -> Result<Vec<Vehicule>, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    let mut stmt = conn.prepare(
        "SELECT id, agence_id, immatriculation, marque_modele, type_boite,
                kilometrage_actuel, date_expiration_assurance, date_visite_technique,
                statut, created_at, updated_at
         FROM vehicules ORDER BY created_at DESC"
    ).map_err(|e| e.to_string())?;

    let rows = stmt.query_map([], |row| {
        Ok(Vehicule {
            id: row.get(0)?,
            agence_id: row.get(1)?,
            immatriculation: row.get(2)?,
            marque_modele: row.get(3)?,
            type_boite: row.get(4)?,
            kilometrage_actuel: row.get(5)?,
            date_expiration_assurance: row.get(6)?,
            date_visite_technique: row.get(7)?,
            statut: row.get(8)?,
            created_at: row.get(9)?,
            updated_at: row.get(10)?,
        })
    }).map_err(|e| e.to_string())?;

    let mut list = Vec::new();
    for r in rows {
        if let Ok(v) = r {
            list.push(v);
        }
    }
    Ok(list)
}

#[tauri::command]
pub fn creer_vehicule(payload: CreerVehiculePayload, db: State<Database>) -> Result<Vehicule, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    let new_id = Uuid::new_v4().to_string();
    let now_str = Utc::now().to_rfc3339();

    let agence_id = match payload.agence_id {
        Some(ref a) if !a.trim().is_empty() => a.clone(),
        _ => get_primary_agence_id(&conn),
    };

    let statut = payload.statut.unwrap_or_else(|| "DISPONIBLE".to_string());
    let immat = payload.immatriculation.trim().to_uppercase();

    conn.execute(
        "INSERT INTO vehicules (
            id, agence_id, immatriculation, marque_modele, type_boite,
            kilometrage_actuel, date_expiration_assurance, date_visite_technique,
            statut, created_at, updated_at
        ) VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9, ?10, ?10)",
        params![
            &new_id,
            &agence_id,
            &immat,
            payload.marque_modele.trim(),
            &payload.type_boite,
            payload.kilometrage_actuel,
            &payload.date_expiration_assurance,
            &payload.date_visite_technique,
            &statut,
            &now_str
        ],
    ).map_err(|e| e.to_string())?;

    let sync_id = Uuid::new_v4().to_string();
    let sync_payload = serde_json::json!({
        "id": new_id,
        "immatriculation": &immat,
        "marque_modele": payload.marque_modele.trim(),
        "type_boite": &payload.type_boite,
        "kilometrage_actuel": payload.kilometrage_actuel,
        "date_expiration_assurance": &payload.date_expiration_assurance,
        "date_visite_technique": &payload.date_visite_technique,
        "statut": &statut
    }).to_string();
    let _ = conn.execute(
        "INSERT INTO journal_sync_events (id, agence_id, nom_table, type_operation, entite_id, donnees_payload, horodatage_local, statut_sync)
         VALUES (?1, ?2, 'vehicules', 'INSERT', ?3, ?4, ?5, 'EN_ATTENTE')",
        params![sync_id, &agence_id, &new_id, sync_payload, &now_str],
    );

    Ok(Vehicule {
        id: new_id,
        agence_id,
        immatriculation: immat,
        marque_modele: payload.marque_modele.trim().to_string(),
        type_boite: payload.type_boite,
        kilometrage_actuel: payload.kilometrage_actuel,
        date_expiration_assurance: payload.date_expiration_assurance,
        date_visite_technique: payload.date_visite_technique,
        statut,
        created_at: now_str.clone(),
        updated_at: now_str,
    })
}

#[tauri::command]
pub fn update_vehicule_statut(vehicule_id: String, statut: String, db: State<Database>) -> Result<serde_json::Value, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    let now_str = Utc::now().to_rfc3339();
    conn.execute(
        "UPDATE vehicules SET statut = ?1, updated_at = ?2 WHERE id = ?3",
        params![&statut, &now_str, &vehicule_id],
    ).map_err(|e| e.to_string())?;

    let agence_id = get_primary_agence_id(&conn);
    let sync_id = Uuid::new_v4().to_string();
    let sync_payload = serde_json::json!({
        "id": &vehicule_id,
        "statut": &statut
    }).to_string();
    let _ = conn.execute(
        "INSERT INTO journal_sync_events (id, agence_id, nom_table, type_operation, entite_id, donnees_payload, horodatage_local, statut_sync)
         VALUES (?1, ?2, 'vehicules', 'UPDATE', ?3, ?4, ?5, 'EN_ATTENTE')",
        params![sync_id, &agence_id, &vehicule_id, sync_payload, &now_str],
    );

    Ok(serde_json::json!({
        "success": true,
        "vehicule_id": vehicule_id,
        "statut": statut
    }))
}

#[tauri::command]
pub fn forcer_candidat_examen(
    payload: ForcerCandidatExamenPayload,
    db: State<Database>,
) -> Result<serde_json::Value, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;

    let candidat_id: String = conn
        .query_row(
            "SELECT candidat_id FROM dossiers_inscription WHERE id = ?1",
            params![&payload.dossier_id],
            |r| r.get(0),
        )
        .map_err(|_| "Dossier d'inscription introuvable en base locale".to_string())?;

    let count: i64 = conn
        .query_row(
            "SELECT count(*) FROM candidats_examen WHERE session_examen_id = ?1",
            params![&payload.session_examen_id],
            |r| r.get(0),
        )
        .unwrap_or(0);
    let numero_table = format!("TABLE-{:03}", count + 1);

    let new_id = Uuid::new_v4().to_string();

    conn.execute(
        "INSERT INTO candidats_examen (
            id, session_examen_id, dossier_id, candidat_id,
            numero_table, resultat_code, resultat_conduite, resultat_final
        ) VALUES (?1, ?2, ?3, ?4, ?5, 'EN_ATTENTE', 'EN_ATTENTE', 'EN_ATTENTE')",
        params![
            &new_id,
            &payload.session_examen_id,
            &payload.dossier_id,
            &candidat_id,
            &numero_table
        ],
    ).map_err(|e| e.to_string())?;

    let agence_id = get_primary_agence_id(&conn);
    let sync_id = Uuid::new_v4().to_string();
    let now_sync = chrono::Utc::now().to_rfc3339();
    let sync_payload = serde_json::json!({
        "id": &new_id,
        "session_examen_id": &payload.session_examen_id,
        "dossier_id": &payload.dossier_id,
        "candidat_id": &candidat_id,
        "numero_table": &numero_table
    }).to_string();
    let _ = conn.execute(
        "INSERT INTO journal_sync_events (id, agence_id, nom_table, type_operation, entite_id, donnees_payload, horodatage_local, statut_sync)
         VALUES (?1, ?2, 'candidats_examen', 'INSERT', ?3, ?4, ?5, 'EN_ATTENTE')",
        params![sync_id, &agence_id, &new_id, sync_payload, now_sync],
    );

    Ok(serde_json::json!({
        "success": true,
        "candidat_examen_id": new_id,
        "session_examen_id": payload.session_examen_id,
        "dossier_id": payload.dossier_id,
        "numero_table": numero_table,
        "est_derogation": true
    }))
}


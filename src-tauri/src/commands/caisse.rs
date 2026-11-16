use tauri::State;
use uuid::Uuid;
use chrono::Utc;
use rusqlite::params;
use crate::db::Database;
use crate::models::CloturerSessionCaissePayload;
use super::system_sync::{get_primary_agence_id, get_primary_caissier_id, md5_or_simple_hash};

#[tauri::command]
pub fn get_guichet_stats(db: State<Database>) -> Result<serde_json::Value, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;

    let total_candidats: i64 = conn
        .query_row("SELECT count(*) FROM candidats", [], |r| r.get(0))
        .unwrap_or(0);

    let active_session: Option<(f64, f64)> = conn
        .query_row(
            "SELECT fond_de_caisse, total_especes_theorique FROM sessions_caisse WHERE statut = 'OUVERTE' ORDER BY date_ouverture DESC LIMIT 1",
            [],
            |r| Ok((r.get(0)?, r.get(1)?)),
        )
        .ok();

    let pending_sync: i64 = conn
        .query_row(
            "SELECT count(*) FROM journal_sync_events WHERE statut_sync = 'EN_ATTENTE'",
            [],
            |r| r.get(0),
        )
        .unwrap_or(0);

    let (fond, total_theorique) = match active_session {
        Some(pair) => pair,
        None => {
            let esp_jour: f64 = conn
                .query_row(
                    "SELECT COALESCE(SUM(montant), 0.0)
                     FROM paiements
                     WHERE mode_reglement = 'ESPECES'
                       AND est_annule = 0
                       AND COALESCE(reference_transaction, '') <> 'REPRISE_HISTORIQUE'
                       AND substr(date_paiement, 1, 10) = date('now')",
                    [],
                    |r| r.get(0),
                )
                .unwrap_or(0.0);
            (0.0, esp_jour)
        }
    };

    Ok(serde_json::json!({
        "total_candidats": total_candidats,
        "fond_de_caisse": fond,
        "total_especes_theorique": total_theorique,
        "evenements_sync_en_attente": pending_sync,
        "mode_connectivite": "LOCAL_OFFLINE_READY"
    }))
}

#[tauri::command]
pub fn enregistrer_paiement(
    dossier_id: String,
    montant: f64,
    mode_reglement: String,
    reference_transaction: Option<String>,
    db: State<Database>,
) -> Result<serde_json::Value, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;

    if montant <= 0.0 {
        return Err("Le montant doit être supérieur à 0 FCFA".to_string());
    }

    let (_candidat_id, solde_actuel, _montant_total, num_dossier): (String, f64, f64, String) = conn
        .query_row(
            "SELECT candidat_id, solde_restant, montant_total, numero_dossier FROM dossiers_inscription WHERE id = ?1",
            params![dossier_id],
            |r| Ok((r.get(0)?, r.get(1)?, r.get(2)?, r.get(3)?)),
        )
        .map_err(|_| "Dossier introuvable".to_string())?;

    if montant > solde_actuel {
        return Err(format!("Le montant ({} FCFA) dépasse le solde restant ({} FCFA)", montant, solde_actuel));
    }

    let (session_id, agence_id, caissier_id): (String, String, String) = conn
        .query_row(
            "SELECT id, agence_id, caissier_id FROM sessions_caisse WHERE statut = 'OUVERTE' ORDER BY date_ouverture DESC LIMIT 1",
            [],
            |r| Ok((r.get(0)?, r.get(1)?, r.get(2)?)),
        )
        .map_err(|_| "Aucune session de caisse active. Veuillez ouvrir la caisse du jour.".to_string())?;

    let auto_ecole_id: String = conn
        .query_row("SELECT auto_ecole_id FROM agences WHERE id = ?1", params![agence_id], |r| r.get(0))
        .map_err(|e| e.to_string())?;

    let paiement_id = Uuid::new_v4().to_string();
    let recu_id = Uuid::new_v4().to_string();
    let now_str = Utc::now().to_rfc3339();

    let count_pay: i64 = conn.query_row("SELECT count(*) FROM paiements", [], |r| r.get(0)).unwrap_or(0);
    let num_recu = format!("REC-2026-10-{:04}", count_pay + 1);
    let sec_hash = format!("{:x}", md5_or_simple_hash(&format!("{}-{}-{}", paiement_id, montant, num_recu)));

    let nouveau_solde = (solde_actuel - montant).max(0.0);

    let tx = conn.unchecked_transaction().map_err(|e| e.to_string())?;

    tx.execute(
        "INSERT INTO paiements (id, auto_ecole_id, agence_id, dossier_id, session_caisse_id, montant, mode_reglement, reference_transaction, date_paiement, encaisse_par_id)
         VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9, ?10)",
        params![paiement_id, auto_ecole_id, agence_id, dossier_id, session_id, montant, mode_reglement, reference_transaction, now_str, caissier_id]
    ).map_err(|e| e.to_string())?;

    tx.execute(
        "INSERT INTO recus_caisse (id, paiement_id, numero_recu, code_securite_hash)
         VALUES (?1, ?2, ?3, ?4)",
        params![recu_id, paiement_id, num_recu, sec_hash]
    ).map_err(|e| e.to_string())?;

    tx.execute(
        "UPDATE dossiers_inscription SET solde_restant = ?1 WHERE id = ?2",
        params![nouveau_solde, dossier_id]
    ).map_err(|e| e.to_string())?;

    if mode_reglement == "ESPECES" {
        tx.execute(
            "UPDATE sessions_caisse SET total_especes_theorique = total_especes_theorique + ?1 WHERE id = ?2",
            params![montant, session_id]
        ).map_err(|e| e.to_string())?;
    }

    let sync_event_id = Uuid::new_v4().to_string();
    let payload = serde_json::json!({
        "paiement_id": paiement_id,
        "dossier_id": dossier_id,
        "montant": montant,
        "mode_reglement": mode_reglement,
        "reference_recu": num_recu,
        "reference_transaction": reference_transaction,
        "horodatage": now_str
    }).to_string();

    tx.execute(
        "INSERT INTO journal_sync_events (id, agence_id, nom_table, type_operation, entite_id, donnees_payload, horodatage_local, statut_sync)
         VALUES (?1, ?2, 'paiements', 'INSERT', ?3, ?4, ?5, 'EN_ATTENTE')",
        params![sync_event_id, agence_id, paiement_id, payload, now_str]
    ).map_err(|e| e.to_string())?;

    tx.commit().map_err(|e| e.to_string())?;

    Ok(serde_json::json!({
        "paiement_id": paiement_id,
        "numero_recu": num_recu,
        "code_securite_hash": sec_hash,
        "montant_verse": montant,
        "nouveau_solde_restant": nouveau_solde,
        "date_paiement": now_str,
        "numero_dossier": num_dossier
    }))
}

#[tauri::command]
pub fn get_session_caisse_detail(db: State<Database>) -> Result<serde_json::Value, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;

    let session_row = conn.query_row(
        "SELECT s.id, s.agence_id, s.caissier_id, COALESCE(u.nom || ' ' || u.prenoms, 'Agent Guichet'), s.date_ouverture, s.fond_de_caisse, s.total_especes_theorique, s.statut
         FROM sessions_caisse s
         LEFT JOIN utilisateurs u ON u.id = s.caissier_id
         WHERE s.statut = 'OUVERTE'
         ORDER BY s.date_ouverture DESC LIMIT 1",
        [],
        |r| {
            Ok((
                r.get::<_, String>(0)?,
                r.get::<_, String>(1)?,
                r.get::<_, String>(2)?,
                r.get::<_, String>(3)?,
                r.get::<_, String>(4)?,
                r.get::<_, f64>(5)?,
                r.get::<_, f64>(6)?,
                r.get::<_, String>(7)?,
            ))
        },
    );

    let (session_id, agence_id, caissier_id, caissier_nom, date_ouv, fond, _total_esp_theorique, statut, has_open_session) = match session_row {
        Ok((id, ag, cid, cnom, dov, fnd, esp, st)) => (id, ag, cid, cnom, dov, fnd, esp, st, true),
        Err(_) => (
            String::new(),
            String::new(),
            String::new(),
            "Session du jour".to_string(),
            Utc::now().format("%Y-%m-%d").to_string(),
            0.0,
            0.0,
            "OUVERTE".to_string(),
            false,
        ),
    };

    let sql_query = if has_open_session {
        "SELECT p.id, p.dossier_id, p.montant, p.mode_reglement, p.reference_transaction,
                p.date_paiement, r.numero_recu, r.nombre_impressions, c.nom || ' ' || c.prenoms, d.numero_dossier
         FROM paiements p
         LEFT JOIN recus_caisse r ON r.paiement_id = p.id
         JOIN dossiers_inscription d ON d.id = p.dossier_id
         JOIN candidats c ON c.id = d.candidat_id
         WHERE p.session_caisse_id = ?1
           AND p.est_annule = 0
           AND COALESCE(p.reference_transaction, '') <> 'REPRISE_HISTORIQUE'
         ORDER BY p.date_paiement DESC"
    } else {
        "SELECT p.id, p.dossier_id, p.montant, p.mode_reglement, p.reference_transaction,
                p.date_paiement, r.numero_recu, r.nombre_impressions, c.nom || ' ' || c.prenoms, d.numero_dossier
         FROM paiements p
         LEFT JOIN recus_caisse r ON r.paiement_id = p.id
         JOIN dossiers_inscription d ON d.id = p.dossier_id
         JOIN candidats c ON c.id = d.candidat_id
         WHERE substr(p.date_paiement, 1, 10) = date('now')
           AND p.est_annule = 0
           AND COALESCE(p.reference_transaction, '') <> 'REPRISE_HISTORIQUE'
         ORDER BY p.date_paiement DESC"
    };

    let mut stmt = conn.prepare(sql_query).map_err(|e| e.to_string())?;

    let mut total_especes = 0.0;
    let mut total_wave = 0.0;
    let mut total_orange = 0.0;
    let mut total_mtn = 0.0;
    let mut total_moov = 0.0;
    let mut total_virement = 0.0;

    let map_row = |row: &rusqlite::Row| {
        let montant: f64 = row.get(2)?;
        let mode: String = row.get(3)?;
        match mode.as_str() {
            "ESPECES" => total_especes += montant,
            "WAVE" => total_wave += montant,
            "ORANGE_MONEY" => total_orange += montant,
            "MTN_MOMO" => total_mtn += montant,
            "MOOV" => total_moov += montant,
            _ => total_virement += montant,
        };

        Ok(serde_json::json!({
            "id": row.get::<_, String>(0)?,
            "dossier_id": row.get::<_, String>(1)?,
            "montant": montant,
            "mode_reglement": mode,
            "reference_transaction": row.get::<_, Option<String>>(4)?,
            "date_paiement": row.get::<_, String>(5)?,
            "numero_recu": row.get::<_, Option<String>>(6)?.unwrap_or_default(),
            "nombre_impressions": row.get::<_, Option<i32>>(7)?.unwrap_or(1),
            "candidat_nom": row.get::<_, String>(8)?,
            "numero_dossier": row.get::<_, String>(9)?,
        }))
    };

    let mut paiements = Vec::new();
    if has_open_session {
        let pay_rows = stmt.query_map(params![&session_id], map_row).map_err(|e| e.to_string())?;
        for p in pay_rows.flatten() {
            paiements.push(p);
        }
    } else {
        let pay_rows = stmt.query_map([], map_row).map_err(|e| e.to_string())?;
        for p in pay_rows.flatten() {
            paiements.push(p);
        }
    }

    let total_encaisse_global = total_especes + total_wave + total_orange + total_mtn + total_moov + total_virement;

    Ok(serde_json::json!({
        "session_ouverte": has_open_session,
        "id": session_id,
        "agence_id": agence_id,
        "caissier_id": caissier_id,
        "caissier_nom": caissier_nom,
        "date_ouverture": date_ouv,
        "fond_de_caisse": fond,
        "total_especes_theorique": total_especes,
        "total_tiroir_attendu": fond + total_especes,
        "statut": statut,
        "ventilation": {
            "especes": total_especes,
            "wave": total_wave,
            "orange_money": total_orange,
            "mtn_momo": total_mtn,
            "moov": total_moov,
            "virement": total_virement,
            "total_encaisse": total_encaisse_global
        },
        "nombre_transactions": paiements.len(),
        "paiements": paiements
    }))
}

#[tauri::command]
pub fn reimprimer_recu(paiement_id: String, db: State<Database>) -> Result<serde_json::Value, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;

    conn.execute(
        "UPDATE recus_caisse SET nombre_impressions = nombre_impressions + 1 WHERE paiement_id = ?1",
        params![paiement_id],
    ).map_err(|e| e.to_string())?;

    let row = conn.query_row(
        "SELECT r.numero_recu, r.code_securite_hash, r.nombre_impressions, p.montant, p.mode_reglement, 
                p.date_paiement, c.nom || ' ' || c.prenoms, d.numero_dossier, f.intitule, d.solde_restant,
                COALESCE(u.nom || ' ' || u.prenoms, 'Agent Guichet')
         FROM recus_caisse r
         JOIN paiements p ON p.id = r.paiement_id
         JOIN dossiers_inscription d ON d.id = p.dossier_id
         JOIN candidats c ON c.id = d.candidat_id
         JOIN formules_permis f ON f.id = d.formule_id
         LEFT JOIN utilisateurs u ON u.id = p.encaisse_par_id
         WHERE r.paiement_id = ?1",
        params![paiement_id],
        |r| {
            Ok(serde_json::json!({
                "numero_recu": r.get::<_, String>(0)?,
                "code_securite_hash": r.get::<_, String>(1)?,
                "nombre_impressions": r.get::<_, i32>(2)?,
                "montant_verse": r.get::<_, f64>(3)?,
                "mode_reglement": r.get::<_, String>(4)?,
                "date_paiement": r.get::<_, String>(5)?,
                "candidat_nom": r.get::<_, String>(6)?,
                "numero_dossier": r.get::<_, String>(7)?,
                "formule_nom": r.get::<_, String>(8)?,
                "solde_restant": r.get::<_, f64>(9)?,
                "est_duplicata": r.get::<_, i32>(2)? > 1,
                "caissier_nom": r.get::<_, String>(10)?
            }))
        },
    ).map_err(|e| e.to_string())?;

    Ok(row)
}

#[tauri::command]
pub fn ouvrir_session_caisse(fond_de_caisse: f64, db: State<Database>) -> Result<String, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;

    let session_id = Uuid::new_v4().to_string();
    let now_str = Utc::now().to_rfc3339();
    let agence_id = get_primary_agence_id(&conn);
    let caissier_id = get_primary_caissier_id(&conn);

    conn.execute(
        "INSERT INTO sessions_caisse (id, agence_id, caissier_id, date_ouverture, fond_de_caisse, total_especes_theorique, statut)
         VALUES (?1, ?2, ?3, ?4, ?5, 0.0, 'OUVERTE')",
        params![session_id, agence_id, caissier_id, now_str, fond_de_caisse],
    ).map_err(|e| e.to_string())?;

    let sync_id = Uuid::new_v4().to_string();
    let sync_payload = serde_json::json!({
        "id": session_id,
        "caissier_id": caissier_id,
        "date_ouverture": now_str,
        "fond_de_caisse": fond_de_caisse
    }).to_string();
    let _ = conn.execute(
        "INSERT INTO journal_sync_events (id, agence_id, nom_table, type_operation, entite_id, donnees_payload, horodatage_local, statut_sync)
         VALUES (?1, ?2, 'sessions_caisse', 'INSERT', ?3, ?4, ?5, 'EN_ATTENTE')",
        params![sync_id, agence_id, session_id, sync_payload, now_str],
    );

    Ok(session_id)
}

#[tauri::command]
pub fn cloturer_session_caisse(
    payload: CloturerSessionCaissePayload,
    db: State<Database>,
    session: State<crate::session_guard::LocalSessionState>,
) -> Result<serde_json::Value, String> {
    crate::session_guard::require_local_roles(&session, &["DIRECTEUR", "SECRETAIRE", "CAISSIER"])?;
    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    let now_str = Utc::now().to_rfc3339();
    let agence_id = get_primary_agence_id(&conn);

    let open_session: Option<(String, String)> = conn
        .query_row(
            "SELECT id, agence_id FROM sessions_caisse WHERE statut = 'OUVERTE' ORDER BY date_ouverture DESC LIMIT 1",
            [],
            |r| Ok((r.get(0)?, r.get(1)?)),
        )
        .ok();

    let updated = conn.execute(
        "UPDATE sessions_caisse
         SET total_especes_compte = ?1,
             ecart_caisse = ?1 - (fond_de_caisse + total_especes_theorique),
             statut = 'CLOTUREE',
             date_fermeture = ?2,
             observations = ?3
         WHERE statut = 'OUVERTE'",
        params![payload.total_especes_compte, &now_str, &payload.observations],
    ).map_err(|e| e.to_string())?;

    let (sync_sess_id, sync_ag_id, op) = if updated == 0 {
        let session_id = Uuid::new_v4().to_string();
        let caissier_id = get_primary_caissier_id(&conn);

        conn.execute(
            "INSERT INTO sessions_caisse (id, agence_id, caissier_id, date_ouverture, date_fermeture, fond_de_caisse, total_especes_theorique, total_especes_compte, ecart_caisse, statut, observations)
             VALUES (?1, ?2, ?3, ?4, ?4, 0.0, 0.0, ?5, 0.0, 'CLOTUREE', ?6)",
            params![session_id, &agence_id, caissier_id, &now_str, payload.total_especes_compte, &payload.observations],
        ).map_err(|e| e.to_string())?;
        (session_id, agence_id, "INSERT")
    } else {
        open_session
            .map(|(id, ag)| (id, ag, "UPDATE"))
            .unwrap_or_else(|| (Uuid::new_v4().to_string(), agence_id, "UPDATE"))
    };

    let sync_id = Uuid::new_v4().to_string();
    let sync_payload = serde_json::json!({
        "id": sync_sess_id,
        "total_especes_compte": payload.total_especes_compte,
        "date_fermeture": now_str,
        "observations": payload.observations
    }).to_string();
    let _ = conn.execute(
        "INSERT INTO journal_sync_events (id, agence_id, nom_table, type_operation, entite_id, donnees_payload, horodatage_local, statut_sync)
         VALUES (?1, ?2, 'sessions_caisse', ?3, ?4, ?5, ?6, 'EN_ATTENTE')",
        params![sync_id, sync_ag_id, op, sync_sess_id, sync_payload, now_str],
    );

    Ok(serde_json::json!({
        "message": "Session de caisse clôturée avec succès"
    }))
}

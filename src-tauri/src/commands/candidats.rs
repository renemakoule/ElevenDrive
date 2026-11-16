use tauri::State;
use uuid::Uuid;
use chrono::Utc;
use rusqlite::params;
use crate::db::Database;
use crate::models::*;
use super::system_sync::{
    get_primary_agence_id, get_primary_caissier_id, get_primary_school_id, md5_or_simple_hash,
};

fn map_candidat_row(row: &rusqlite::Row) -> rusqlite::Result<Candidat> {
    Ok(Candidat {
        id: row.get(0)?,
        auto_ecole_id: row.get(1)?,
        agence_id: row.get(2)?,
        nom: row.get(3)?,
        prenoms: row.get(4)?,
        genre: row.get(5)?,
        date_naissance: row.get(6)?,
        lieu_naissance: row.get(7)?,
        nationalite: row.get(8)?,
        profession: row.get(9)?,
        quartier_residence: row.get(10)?,
        telephone_whatsapp: row.get(11)?,
        telephone_urgence: row.get(12)?,
        groupe_sanguin: row.get(13)?,
        photo_url: row.get(14)?,
        type_piece_id: row.get(15)?,
        numero_piece_id: row.get(16)?,
        date_expiration_piece: row.get(17)?,
        visite_medicale_statut: row.get(18)?,
        date_visite_medicale: row.get(19)?,
        created_at: row.get(20)?,
    })
}

#[tauri::command]
pub fn search_candidats(query: String, db: State<Database>) -> Result<Vec<Candidat>, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    let search_pattern = format!("%{}%", query.trim());

    let mut stmt = conn
        .prepare(
            "SELECT c.id, c.auto_ecole_id, c.agence_id, c.nom, c.prenoms, c.genre, c.date_naissance, 
                    c.lieu_naissance, c.nationalite, c.profession, c.quartier_residence, 
                    c.telephone_whatsapp, c.telephone_urgence, c.groupe_sanguin, c.photo_url, 
                    c.type_piece_id, c.numero_piece_id, c.date_expiration_piece, 
                    c.visite_medicale_statut, c.date_visite_medicale, c.created_at
             FROM candidats c
             LEFT JOIN dossiers_inscription d ON d.candidat_id = c.id
             WHERE c.nom LIKE ?1 
                OR c.prenoms LIKE ?1 
                OR c.telephone_whatsapp LIKE ?1 
                OR c.numero_piece_id LIKE ?1 
                OR d.numero_dossier LIKE ?1
             ORDER BY c.created_at DESC LIMIT 25",
        )
        .map_err(|e| e.to_string())?;

    let rows = stmt.query_map(params![search_pattern], map_candidat_row).map_err(|e| e.to_string())?;
    let mut list = Vec::new();
    for r in rows {
        if let Ok(c) = r {
            list.push(c);
        }
    }
    Ok(list)
}

#[tauri::command]
pub fn get_candidat_dossier(candidat_id: String, db: State<Database>) -> Result<Option<DossierDetailComplet>, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;

    let cand_res = conn.query_row(
        "SELECT id, auto_ecole_id, agence_id, nom, prenoms, genre, date_naissance, 
                lieu_naissance, nationalite, profession, quartier_residence, 
                telephone_whatsapp, telephone_urgence, groupe_sanguin, photo_url, 
                type_piece_id, numero_piece_id, date_expiration_piece, 
                visite_medicale_statut, date_visite_medicale, created_at
         FROM candidats WHERE id = ?1",
        params![candidat_id],
        map_candidat_row,
    );

    let candidat = match cand_res {
        Ok(c) => c,
        Err(_) => return Ok(None),
    };

    let dossier_res = conn.query_row(
        "SELECT id, candidat_id, formule_id, numero_dossier, date_inscription, 
                montant_total, remise_accordee, solde_restant, statut_dossier, 
                heures_code_effectuees, heures_conduite_effectuees, COALESCE(mode_formation_theorique, 'PRESENTIEL')
         FROM dossiers_inscription WHERE candidat_id = ?1 ORDER BY date_inscription DESC LIMIT 1",
        params![candidat.id],
        |row| {
            Ok(DossierInscription {
                id: row.get(0)?,
                candidat_id: row.get(1)?,
                formule_id: row.get(2)?,
                numero_dossier: row.get(3)?,
                date_inscription: row.get(4)?,
                montant_total: row.get(5)?,
                remise_accordee: row.get(6)?,
                solde_restant: row.get(7)?,
                statut_dossier: row.get(8)?,
                heures_code_effectuees: row.get(9)?,
                heures_conduite_effectuees: row.get(10)?,
                mode_formation_theorique: row.get(11)?,
            })
        },
    );

    let dossier = match dossier_res {
        Ok(d) => d,
        Err(_) => return Ok(None),
    };

    let formule = conn.query_row(
        "SELECT id, categorie_permis, intitule, heures_theoriques, heures_pratiques, prix_total, acompte_minimum, est_actif
         FROM formules_permis WHERE id = ?1",
        params![dossier.formule_id],
        |row| {
            Ok(FormulePermis {
                id: row.get(0)?,
                categorie_permis: row.get(1)?,
                intitule: row.get(2)?,
                heures_theoriques: row.get(3)?,
                heures_pratiques: row.get(4)?,
                prix_total: row.get(5)?,
                acompte_minimum: row.get(6)?,
                est_actif: row.get::<_, i32>(7)? == 1,
                tranches_paiement: None,
            })
        },
    ).map_err(|e| e.to_string())?;

    let mut pay_stmt = conn.prepare(
        "SELECT p.id, p.dossier_id, p.session_caisse_id, p.montant, p.mode_reglement, 
                p.reference_transaction, p.date_paiement, p.encaisse_par_id, r.numero_recu
         FROM paiements p
         LEFT JOIN recus_caisse r ON r.paiement_id = p.id
         WHERE p.dossier_id = ?1 AND p.est_annule = 0
         ORDER BY p.date_paiement ASC",
    ).map_err(|e| e.to_string())?;

    let cand_full_name = format!("{} {}", candidat.nom, candidat.prenoms);
    let mut cumul_paye = 0.0;
    let pay_rows = pay_stmt.query_map(params![dossier.id], |row| {
        let m: f64 = row.get(3)?;
        cumul_paye += m;
        let num_recu: Option<String> = row.get(8)?;
        Ok(Paiement {
            id: row.get(0)?,
            dossier_id: row.get(1)?,
            session_caisse_id: row.get(2)?,
            montant: m,
            mode_reglement: row.get(4)?,
            reference_transaction: row.get(5)?,
            date_paiement: row.get(6)?,
            encaisse_par_id: row.get(7)?,
            reference_recu: num_recu.unwrap_or_else(|| "SANS_RECU".to_string()),
            solde_apres_paiement: dossier.montant_total - cumul_paye,
            candidat_nom_complet: cand_full_name.clone(),
        })
    }).map_err(|e| e.to_string())?;

    let mut paiements = Vec::new();
    for p in pay_rows {
        if let Ok(item) = p {
            paiements.push(item);
        }
    }

    let code_pct = if formule.heures_theoriques > 0 {
        ((dossier.heures_code_effectuees as f64) / (formule.heures_theoriques as f64) * 100.0).min(100.0)
    } else {
        100.0
    };

    let conduite_pct = if formule.heures_pratiques > 0 {
        ((dossier.heures_conduite_effectuees as f64) / (formule.heures_pratiques as f64) * 100.0).min(100.0)
    } else {
        100.0
    };

    Ok(Some(DossierDetailComplet {
        candidat,
        dossier,
        formule,
        paiements,
        progression_code_pct: code_pct,
        progression_conduite_pct: conduite_pct,
    }))
}

#[tauri::command]
pub fn update_candidat_photo(
    candidat_id: String,
    photo_url: String,
    db: State<Database>,
) -> Result<bool, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    conn.execute(
        "UPDATE candidats SET photo_url = ?1 WHERE id = ?2",
        params![photo_url, candidat_id],
    )
    .map_err(|e| e.to_string())?;

    let sync_event_id = Uuid::new_v4().to_string();
    let now_str = Utc::now().to_rfc3339();
    let payload = serde_json::json!({
        "candidat_id": candidat_id,
        "action": "PHOTO_4X4_UPDATED"
    }).to_string();

    let agence_id: String = conn
        .query_row("SELECT agence_id FROM candidats WHERE id = ?1", params![candidat_id], |r| r.get(0))
        .unwrap_or_else(|_| "ag-cocody-01".to_string());

    conn.execute(
        "INSERT INTO journal_sync_events (id, agence_id, nom_table, type_operation, entite_id, donnees_payload, horodatage_local, statut_sync)
         VALUES (?1, ?2, 'candidats', 'UPDATE', ?3, ?4, ?5, 'EN_ATTENTE')",
        params![sync_event_id, agence_id, candidat_id, payload, now_str],
    ).ok();

    Ok(true)
}

#[tauri::command]
pub fn creer_candidat_avec_dossier(
    payload: CreerCandidatPayload,
    db: State<Database>,
) -> Result<serde_json::Value, String> {
    let conn = db.conn.lock().map_err(|e| e.to_string())?;

    let formule = conn.query_row(
        "SELECT id, categorie_permis, intitule, heures_theoriques, heures_pratiques, prix_total, acompte_minimum, est_actif
         FROM formules_permis WHERE id = ?1",
        params![&payload.formule_id],
        |row| {
            Ok(FormulePermis {
                id: row.get(0)?,
                categorie_permis: row.get(1)?,
                intitule: row.get(2)?,
                heures_theoriques: row.get(3)?,
                heures_pratiques: row.get(4)?,
                prix_total: row.get(5)?,
                acompte_minimum: row.get(6)?,
                est_actif: row.get::<_, i32>(7)? == 1,
                tranches_paiement: None,
            })
        },
    ).map_err(|_| "Formule de permis introuvable dans la base locale".to_string())?;

    let auto_ecole_id = get_primary_school_id(&conn);
    let agence_id = get_primary_agence_id(&conn);
    let caissier_id = get_primary_caissier_id(&conn);

    let cand_id = Uuid::new_v4().to_string();
    let dos_id = Uuid::new_v4().to_string();
    let now_str = Utc::now().to_rfc3339();
    let today_str = Utc::now().format("%Y-%m-%d").to_string();

    let count_dossiers: i64 = conn
        .query_row("SELECT count(*) FROM dossiers_inscription", [], |r| r.get(0))
        .unwrap_or(0);
    let num_dossier = format!("DOS-{}-{:04}", Utc::now().format("%Y"), count_dossiers + 1);

    let visite_statut = payload
        .visite_medicale_statut
        .clone()
        .unwrap_or_else(|| "VALIDE_SANS_RESERVE".to_string());
    let mode_formation = payload
        .mode_formation_theorique
        .clone()
        .unwrap_or_else(|| "PRESENTIEL".to_string());
    let solde_restant = (formule.prix_total - payload.acompte_initial).max(0.0);

    let tx = conn.unchecked_transaction().map_err(|e| e.to_string())?;

    tx.execute(
        "INSERT INTO candidats (
            id, auto_ecole_id, agence_id, nom, prenoms, genre, date_naissance,
            lieu_naissance, nationalite, telephone_whatsapp, photo_url,
            type_piece_id, numero_piece_id, visite_medicale_statut, created_at
        ) VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, 'Ivoirienne', ?9, ?10, ?11, ?12, ?13, ?14)",
        params![
            &cand_id,
            &auto_ecole_id,
            &agence_id,
            payload.nom.trim().to_uppercase(),
            payload.prenoms.trim(),
            &payload.genre,
            &payload.date_naissance,
            &payload.lieu_naissance,
            payload.telephone_whatsapp.trim(),
            &payload.photo_url,
            &payload.type_piece_id,
            payload.numero_piece_id.trim(),
            &visite_statut,
            &now_str
        ],
    ).map_err(|e| e.to_string())?;

    tx.execute(
        "INSERT INTO dossiers_inscription (
            id, auto_ecole_id, agence_id, candidat_id, formule_id, numero_dossier,
            date_inscription, montant_total, remise_accordee, solde_restant,
            statut_dossier, heures_code_effectuees, heures_conduite_effectuees,
            derogation_conduite_accordee, mode_formation_theorique
        ) VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, 0.0, ?9, 'EN_COURS', 0, 0, 0, ?10)",
        params![
            &dos_id,
            &auto_ecole_id,
            &agence_id,
            &cand_id,
            &formule.id,
            &num_dossier,
            &today_str,
            formule.prix_total,
            solde_restant,
            &mode_formation
        ],
    ).map_err(|e| e.to_string())?;

    let mut paiements = Vec::new();
    if payload.acompte_initial > 0.0 {
        let session_id: String = tx
            .query_row(
                "SELECT id FROM sessions_caisse WHERE statut = 'OUVERTE' ORDER BY date_ouverture DESC LIMIT 1",
                [],
                |r| r.get(0),
            )
            .unwrap_or_else(|_| {
                let s_id = Uuid::new_v4().to_string();
                let _ = tx.execute(
                    "INSERT INTO sessions_caisse (id, agence_id, caissier_id, date_ouverture, fond_de_caisse, total_especes_theorique, statut)
                     VALUES (?1, ?2, ?3, ?4, 50000.0, 0.0, 'OUVERTE')",
                    params![&s_id, &agence_id, &caissier_id, &now_str],
                );
                s_id
            });

        let pay_id = Uuid::new_v4().to_string();
        let recu_id = Uuid::new_v4().to_string();
        let count_pay: i64 = tx.query_row("SELECT count(*) FROM paiements", [], |r| r.get(0)).unwrap_or(0);
        let num_recu = format!("REC-{}-{:04}", Utc::now().format("%Y%m"), count_pay + 1);
        let sec_hash = format!("{:x}", md5_or_simple_hash(&format!("{}-{}-{}", pay_id, payload.acompte_initial, num_recu)));

        tx.execute(
            "INSERT INTO paiements (
                id, auto_ecole_id, agence_id, dossier_id, session_caisse_id,
                montant, mode_reglement, date_paiement, encaisse_par_id
            ) VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9)",
            params![
                &pay_id,
                &auto_ecole_id,
                &agence_id,
                &dos_id,
                &session_id,
                payload.acompte_initial,
                &payload.mode_reglement,
                &now_str,
                &caissier_id
            ],
        ).map_err(|e| e.to_string())?;

        tx.execute(
            "INSERT INTO recus_caisse (id, paiement_id, numero_recu, code_securite_hash)
             VALUES (?1, ?2, ?3, ?4)",
            params![&recu_id, &pay_id, &num_recu, &sec_hash],
        ).map_err(|e| e.to_string())?;

        if payload.mode_reglement == "ESPECES" {
            tx.execute(
                "UPDATE sessions_caisse SET total_especes_theorique = total_especes_theorique + ?1 WHERE id = ?2",
                params![payload.acompte_initial, &session_id],
            ).map_err(|e| e.to_string())?;
        }

        paiements.push(Paiement {
            id: pay_id,
            dossier_id: dos_id.clone(),
            session_caisse_id: session_id,
            montant: payload.acompte_initial,
            mode_reglement: payload.mode_reglement.clone(),
            reference_transaction: None,
            date_paiement: now_str.clone(),
            encaisse_par_id: caissier_id,
            reference_recu: num_recu,
            solde_apres_paiement: solde_restant,
            candidat_nom_complet: format!("{} {}", payload.nom.trim().to_uppercase(), payload.prenoms.trim()),
        });
    }

    let sync_event_id = Uuid::new_v4().to_string();
    let sync_payload = serde_json::json!({
        "candidat_id": cand_id,
        "dossier_id": dos_id,
        "formule_id": formule.id,
        "numero_dossier": num_dossier,
        "nom": &payload.nom,
        "prenoms": &payload.prenoms,
        "genre": &payload.genre,
        "date_naissance": &payload.date_naissance,
        "lieu_naissance": &payload.lieu_naissance,
        "telephone_whatsapp": &payload.telephone_whatsapp,
        "type_piece_id": &payload.type_piece_id,
        "numero_piece_id": &payload.numero_piece_id,
        "visite_medicale_statut": &visite_statut,
        "montant_total": formule.prix_total,
        "solde_restant": solde_restant,
        "mode_formation_theorique": &mode_formation
    }).to_string();

    tx.execute(
        "INSERT INTO journal_sync_events (id, agence_id, nom_table, type_operation, entite_id, donnees_payload, horodatage_local, statut_sync)
         VALUES (?1, ?2, 'candidats', 'INSERT', ?3, ?4, ?5, 'EN_ATTENTE')",
        params![sync_event_id, &agence_id, &cand_id, sync_payload, &now_str],
    ).map_err(|e| e.to_string())?;

    tx.commit().map_err(|e| e.to_string())?;

    let candidat = Candidat {
        id: cand_id.clone(),
        auto_ecole_id,
        agence_id,
        nom: payload.nom.trim().to_uppercase(),
        prenoms: payload.prenoms.trim().to_string(),
        genre: payload.genre,
        date_naissance: payload.date_naissance,
        lieu_naissance: payload.lieu_naissance,
        nationalite: "Ivoirienne".to_string(),
        profession: None,
        quartier_residence: None,
        telephone_whatsapp: payload.telephone_whatsapp.trim().to_string(),
        telephone_urgence: None,
        groupe_sanguin: None,
        photo_url: payload.photo_url,
        type_piece_id: payload.type_piece_id,
        numero_piece_id: payload.numero_piece_id.trim().to_string(),
        date_expiration_piece: None,
        visite_medicale_statut: visite_statut,
        date_visite_medicale: None,
        created_at: now_str,
    };

    let dossier = DossierInscription {
        id: dos_id,
        candidat_id: cand_id,
        formule_id: formule.id.clone(),
        numero_dossier: num_dossier,
        date_inscription: today_str,
        montant_total: formule.prix_total,
        remise_accordee: 0.0,
        solde_restant,
        statut_dossier: "EN_COURS".to_string(),
        heures_code_effectuees: 0,
        heures_conduite_effectuees: 0,
        mode_formation_theorique: mode_formation,
    };

    let dossier_detail = DossierDetailComplet {
        candidat: candidat.clone(),
        dossier,
        formule,
        paiements,
        progression_code_pct: 0.0,
        progression_conduite_pct: 0.0,
    };

    Ok(serde_json::json!({
        "candidat": candidat,
        "dossier_detail": dossier_detail
    }))
}

use std::collections::{HashMap, HashSet};
use chrono::Utc;
use rusqlite::params;
use serde::{Deserialize, Serialize};
use tauri::State;
use uuid::Uuid;

use crate::db::Database;
use crate::models::CreerCandidatPayload;
use super::system_sync::{
    get_primary_agence_id, get_primary_caissier_id, get_primary_school_id, md5_or_simple_hash,
};

#[derive(Debug, Serialize, Deserialize)]
pub struct ImporterCandidatsBatchPayload {
    pub candidats: Vec<CreerCandidatPayload>,
    pub ignorer_doublons: bool,
    pub impacter_caisse_jour: bool,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ImportErreurLigne {
    pub index: usize,
    pub nom: String,
    pub prenoms: String,
    pub numero_piece_id: String,
    pub raison: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ImportDossierCree {
    pub candidat_id: String,
    pub dossier_id: String,
    pub numero_dossier: String,
    pub nom_complet: String,
    pub solde_restant: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ImporterCandidatsBatchResult {
    pub total_recus: usize,
    pub total_crees: usize,
    pub total_ignores: usize,
    pub total_erreurs: usize,
    pub dossiers_crees: Vec<ImportDossierCree>,
    pub erreurs: Vec<ImportErreurLigne>,
}

#[tauri::command]
pub fn importer_candidats_batch(
    payload: ImporterCandidatsBatchPayload,
    db: State<Database>,
) -> Result<ImporterCandidatsBatchResult, String> {
    if payload.candidats.is_empty() {
        return Err("Aucun candidat à importer".to_string());
    }
    if payload.candidats.len() > 1000 {
        return Err("Limite dépassée : maximum 1000 candidats par import".to_string());
    }

    let conn = db.conn.lock().map_err(|e| e.to_string())?;
    let auto_ecole_id = get_primary_school_id(&conn);
    let agence_id = get_primary_agence_id(&conn);
    let caissier_id = get_primary_caissier_id(&conn);

    // Pré-charger les prix de formules
    let mut formules_prix: HashMap<String, f64> = HashMap::new();
    {
        let mut stmt = conn
            .prepare("SELECT id, prix_total FROM formules_permis")
            .map_err(|e| e.to_string())?;
        let rows = stmt
            .query_map([], |r| Ok((r.get::<_, String>(0)?, r.get::<_, f64>(1)?)))
            .map_err(|e| e.to_string())?;
        for row in rows.flatten() {
            formules_prix.insert(row.0, row.1);
        }
    }

    // Pré-charger les numéros de pièces existants
    let mut existing_pieces: HashSet<String> = HashSet::new();
    {
        let mut stmt = conn
            .prepare("SELECT LOWER(TRIM(numero_piece_id)) FROM candidats WHERE numero_piece_id IS NOT NULL")
            .map_err(|e| e.to_string())?;
        let rows = stmt
            .query_map([], |r| r.get::<_, String>(0))
            .map_err(|e| e.to_string())?;
        for p in rows.flatten() {
            if !p.is_empty() {
                existing_pieces.insert(p);
            }
        }
    }

    let mut count_dossiers: i64 = conn
        .query_row("SELECT count(*) FROM dossiers_inscription", [], |r| r.get(0))
        .unwrap_or(0);
    let mut count_pay: i64 = conn
        .query_row("SELECT count(*) FROM paiements", [], |r| r.get(0))
        .unwrap_or(0);

    let now_str = Utc::now().to_rfc3339();
    let today_str = Utc::now().format("%Y-%m-%d").to_string();
    let year_str = Utc::now().format("%Y").to_string();
    let ym_str = Utc::now().format("%Y%m").to_string();

    let mut result = ImporterCandidatsBatchResult {
        total_recus: payload.candidats.len(),
        total_crees: 0,
        total_ignores: 0,
        total_erreurs: 0,
        dossiers_crees: Vec::new(),
        erreurs: Vec::new(),
    };

    let mut tx = conn.unchecked_transaction().map_err(|e| e.to_string())?;

    // Préparer une session caisse si des paiements impactent la caisse du jour
    let needs_caisse = payload.impacter_caisse_jour
        && payload.candidats.iter().any(|c| c.acompte_initial > 0.0);
    let session_caisse_id: Option<String> = if needs_caisse {
        let s_id = tx
            .query_row(
                "SELECT id FROM sessions_caisse WHERE statut = 'OUVERTE' ORDER BY date_ouverture DESC LIMIT 1",
                [],
                |r| r.get::<_, String>(0),
            )
            .unwrap_or_else(|_| {
                let new_s_id = Uuid::new_v4().to_string();
                let _ = tx.execute(
                    "INSERT INTO sessions_caisse (id, agence_id, caissier_id, date_ouverture, fond_de_caisse, total_especes_theorique, statut)
                     VALUES (?1, ?2, ?3, ?4, 50000.0, 0.0, 'OUVERTE')",
                    params![&new_s_id, &agence_id, &caissier_id, &now_str],
                );
                new_s_id
            });
        Some(s_id)
    } else {
        None
    };

    for (idx, item) in payload.candidats.iter().enumerate() {
        let nom = item.nom.trim().to_uppercase();
        let prenoms = item.prenoms.trim().to_string();
        let tel = item.telephone_whatsapp.trim().to_string();
        let piece_raw = item.numero_piece_id.trim().to_string();
        let piece_norm = piece_raw.to_lowercase();

        if nom.is_empty() || prenoms.is_empty() || tel.is_empty() || piece_norm.is_empty() || item.formule_id.trim().is_empty() {
            result.total_erreurs += 1;
            result.erreurs.push(ImportErreurLigne {
                index: idx + 1,
                nom,
                prenoms,
                numero_piece_id: piece_raw,
                raison: "Champs obligatoires manquants".to_string(),
            });
            continue;
        }

        let Some(&prix_total) = formules_prix.get(&item.formule_id) else {
            result.total_erreurs += 1;
            result.erreurs.push(ImportErreurLigne {
                index: idx + 1,
                nom,
                prenoms,
                numero_piece_id: piece_raw,
                raison: "Formule de permis introuvable".to_string(),
            });
            continue;
        };

        if existing_pieces.contains(&piece_norm) {
            if payload.ignorer_doublons {
                result.total_ignores += 1;
                result.erreurs.push(ImportErreurLigne {
                    index: idx + 1,
                    nom,
                    prenoms,
                    numero_piece_id: piece_raw,
                    raison: "Doublon ignoré : numéro de pièce déjà existant".to_string(),
                });
            } else {
                result.total_erreurs += 1;
                result.erreurs.push(ImportErreurLigne {
                    index: idx + 1,
                    nom,
                    prenoms,
                    numero_piece_id: piece_raw,
                    raison: "Doublon : numéro de pièce déjà existant".to_string(),
                });
            }
            continue;
        }

        let sp = tx.savepoint().map_err(|e| e.to_string())?;
        let cand_id = Uuid::new_v4().to_string();
        let dos_id = Uuid::new_v4().to_string();
        count_dossiers += 1;
        let num_dossier = format!("DOS-{}-{:04}", year_str, count_dossiers);

        let genre = if item.genre.trim().eq_ignore_ascii_case("F") { "F" } else { "M" };
        let type_piece = if item.type_piece_id.trim().is_empty() { "CNI" } else { item.type_piece_id.trim() };
        let lieu = if item.lieu_naissance.trim().is_empty() { "Abidjan" } else { item.lieu_naissance.trim() };
        let visite_statut = item.visite_medicale_statut.clone().unwrap_or_else(|| "VALIDE_SANS_RESERVE".to_string());
        let mode_formation = item.mode_formation_theorique.clone().unwrap_or_else(|| "PRESENTIEL".to_string());

        let acompte = item.acompte_initial.max(0.0).min(prix_total);
        let solde_restant = (prix_total - acompte).max(0.0);

        let insert_cand = sp.execute(
            "INSERT INTO candidats (
                id, auto_ecole_id, agence_id, nom, prenoms, genre, date_naissance,
                lieu_naissance, nationalite, telephone_whatsapp, photo_url,
                type_piece_id, numero_piece_id, visite_medicale_statut, created_at
            ) VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, 'Ivoirienne', ?9, ?10, ?11, ?12, ?13, ?14)",
            params![
                &cand_id, &auto_ecole_id, &agence_id, &nom, &prenoms, genre,
                &item.date_naissance, lieu, &tel, &item.photo_url, type_piece,
                &piece_raw, &visite_statut, &now_str
            ],
        );
        if let Err(e) = insert_cand {
            count_dossiers -= 1;
            result.total_erreurs += 1;
            result.erreurs.push(ImportErreurLigne {
                index: idx + 1, nom, prenoms, numero_piece_id: piece_raw,
                raison: format!("Erreur insertion candidat : {}", e),
            });
            continue;
        }

        let insert_dos = sp.execute(
            "INSERT INTO dossiers_inscription (
                id, auto_ecole_id, agence_id, candidat_id, formule_id, numero_dossier,
                date_inscription, montant_total, remise_accordee, solde_restant,
                statut_dossier, heures_code_effectuees, heures_conduite_effectuees,
                derogation_conduite_accordee, mode_formation_theorique
            ) VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, 0.0, ?9, 'EN_COURS', 0, 0, 0, ?10)",
            params![
                &dos_id, &auto_ecole_id, &agence_id, &cand_id, &item.formule_id,
                &num_dossier, &today_str, prix_total, solde_restant, &mode_formation
            ],
        );
        if let Err(e) = insert_dos {
            count_dossiers -= 1;
            result.total_erreurs += 1;
            result.erreurs.push(ImportErreurLigne {
                index: idx + 1, nom, prenoms, numero_piece_id: piece_raw,
                raison: format!("Erreur insertion dossier : {}", e),
            });
            continue;
        }

        if acompte > 0.0 {
            count_pay += 1;
            let pay_id = Uuid::new_v4().to_string();
            let recu_id = Uuid::new_v4().to_string();
            let num_recu = format!("REC-{}-{:05}", ym_str, count_pay);
            let sec_hash = format!("{:x}", md5_or_simple_hash(&format!("{}-{}-{}", pay_id, acompte, num_recu)));
            let mode_reg = if item.mode_reglement.trim().is_empty() { "ESPECES".to_string() } else { item.mode_reglement.trim().to_uppercase() };
            let ref_tx: Option<&str> = if payload.impacter_caisse_jour { None } else { Some("REPRISE_HISTORIQUE") };

            // Si reprise historique sans session active, on rattache sans modifier total_especes_theorique
            let final_sess = if let Some(ref s) = session_caisse_id {
                s.clone()
            } else {
                sp.query_row(
                    "SELECT id FROM sessions_caisse ORDER BY date_ouverture DESC LIMIT 1",
                    [],
                    |r| r.get::<_, String>(0),
                )
                .unwrap_or_else(|_| {
                    let s_id = Uuid::new_v4().to_string();
                    let _ = sp.execute(
                        "INSERT INTO sessions_caisse (id, agence_id, caissier_id, date_ouverture, fond_de_caisse, total_especes_theorique, statut)
                         VALUES (?1, ?2, ?3, ?4, 0.0, 0.0, 'CLOTUREE')",
                        params![&s_id, &agence_id, &caissier_id, &now_str],
                    );
                    s_id
                })
            };

            if let Err(e) = sp.execute(
                "INSERT INTO paiements (
                    id, auto_ecole_id, agence_id, dossier_id, session_caisse_id,
                    montant, mode_reglement, reference_transaction, date_paiement, encaisse_par_id
                ) VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9, ?10)",
                params![
                    &pay_id, &auto_ecole_id, &agence_id, &dos_id, &final_sess,
                    acompte, &mode_reg, ref_tx, &now_str, &caissier_id
                ],
            ) {
                count_dossiers -= 1;
                count_pay -= 1;
                result.total_erreurs += 1;
                result.erreurs.push(ImportErreurLigne {
                    index: idx + 1, nom, prenoms, numero_piece_id: piece_raw,
                    raison: format!("Erreur paiement : {}", e),
                });
                continue;
            }

            let _ = sp.execute(
                "INSERT INTO recus_caisse (id, paiement_id, numero_recu, code_securite_hash) VALUES (?1, ?2, ?3, ?4)",
                params![&recu_id, &pay_id, &num_recu, &sec_hash],
            );

            if payload.impacter_caisse_jour && mode_reg == "ESPECES" {
                let _ = sp.execute(
                    "UPDATE sessions_caisse SET total_especes_theorique = total_especes_theorique + ?1 WHERE id = ?2",
                    params![acompte, &final_sess],
                );
            }
        }

        let sync_id = Uuid::new_v4().to_string();
        let sync_payload = serde_json::json!({
            "candidat_id": cand_id,
            "dossier_id": dos_id,
            "formule_id": item.formule_id,
            "numero_dossier": num_dossier,
            "nom": nom,
            "prenoms": prenoms,
            "genre": genre,
            "date_naissance": item.date_naissance,
            "lieu_naissance": lieu,
            "telephone_whatsapp": tel,
            "type_piece_id": type_piece,
            "numero_piece_id": piece_raw,
            "visite_medicale_statut": visite_statut,
            "montant_total": prix_total,
            "solde_restant": solde_restant,
            "mode_formation_theorique": mode_formation
        }).to_string();
        let _ = sp.execute(
            "INSERT INTO journal_sync_events (id, agence_id, nom_table, type_operation, entite_id, donnees_payload, horodatage_local, statut_sync)
             VALUES (?1, ?2, 'candidats', 'INSERT', ?3, ?4, ?5, 'EN_ATTENTE')",
            params![sync_id, &agence_id, &cand_id, sync_payload, &now_str],
        );

        sp.commit().map_err(|e| e.to_string())?;
        existing_pieces.insert(piece_norm);
        result.total_crees += 1;
        result.dossiers_crees.push(ImportDossierCree {
            candidat_id: cand_id,
            dossier_id: dos_id,
            numero_dossier: num_dossier,
            nom_complet: format!("{} {}", nom, prenoms),
            solde_restant,
        });
    }

    tx.commit().map_err(|e| e.to_string())?;
    Ok(result)
}

use crate::db::Database;
use rusqlite::params;
use serde_json::{json, Value};
use std::collections::HashSet;
use tauri::State;

fn str_val<'a>(v: &'a Value, key: &str) -> &'a str {
    v.get(key).and_then(|x| x.as_str()).unwrap_or("")
}

fn opt_str<'a>(v: &'a Value, key: &str) -> Option<&'a str> {
    v.get(key)
        .and_then(|x| x.as_str())
        .filter(|s| !s.is_empty())
}

fn f64_val(v: &Value, key: &str) -> f64 {
    v.get(key).and_then(|x| x.as_f64()).unwrap_or(0.0)
}

fn i64_val(v: &Value, key: &str, default: i64) -> i64 {
    v.get(key).and_then(|x| x.as_i64()).unwrap_or(default)
}

fn bool_int(v: &Value, key: &str, default: i64) -> i64 {
    match v.get(key) {
        Some(Value::Bool(b)) => if *b { 1 } else { 0 },
        Some(Value::Number(n)) => n.as_i64().unwrap_or(default),
        _ => default,
    }
}

/// Applique un instantané Cloud Neon dans SQLite local sans écraser les entités
/// ayant des mutations locales `'EN_ATTENTE'` dans `journal_sync_events`.
#[tauri::command]
pub fn appliquer_snapshot_cloud(snapshot: Value, db: State<Database>) -> Result<Value, String> {
    let mut conn = db.conn.lock().map_err(|e| e.to_string())?;

    // Récupérer les entite_id en attente de Push pour ne pas les écraser
    let mut pending_ids: HashSet<String> = HashSet::new();
    {
        let mut stmt = conn
            .prepare("SELECT entite_id FROM journal_sync_events WHERE statut_sync = 'EN_ATTENTE'")
            .map_err(|e| e.to_string())?;
        let rows = stmt
            .query_map([], |r| r.get::<_, String>(0))
            .map_err(|e| e.to_string())?;
        for id in rows.flatten() {
            pending_ids.insert(id);
        }
    }

    let tx = conn.transaction().map_err(|e| e.to_string())?;
    let mut ignores: i64 = 0;

    // 1. Auto-école (parent des agences, formules, candidats, etc.)
    if let Some(ae) = snapshot.get("auto_ecole").filter(|v| v.is_object()) {
        let ae_id = str_val(ae, "id");
        if !ae_id.is_empty() {
            tx.execute(
                "INSERT INTO auto_ecoles (id, code_unique, nom_commercial, numero_agrement, telephone, email, devise)
                 VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7)
                 ON CONFLICT(id) DO UPDATE SET
                   nom_commercial = excluded.nom_commercial,
                   numero_agrement = excluded.numero_agrement,
                   telephone = excluded.telephone,
                   email = excluded.email,
                   devise = excluded.devise",
                params![
                    ae_id,
                    str_val(ae, "code_unique"),
                    str_val(ae, "nom_commercial"),
                    str_val(ae, "numero_agrement"),
                    str_val(ae, "telephone"),
                    opt_str(ae, "email"),
                    opt_str(ae, "devise").unwrap_or("XOF"),
                ],
            )
            .map_err(|e| e.to_string())?;
        }
    }

    // 2. Agences
    if let Some(agences) = snapshot.get("agences").and_then(|v| v.as_array()) {
        for ag in agences {
            let id = str_val(ag, "id");
            let ae_id = str_val(ag, "auto_ecole_id");
            if id.is_empty() || ae_id.is_empty() {
                continue;
            }
            tx.execute(
                "INSERT INTO agences (id, auto_ecole_id, nom, adresse, telephone)
                 VALUES (?1, ?2, ?3, ?4, ?5)
                 ON CONFLICT(id) DO UPDATE SET
                   nom = excluded.nom,
                   adresse = excluded.adresse,
                   telephone = excluded.telephone",
                params![
                    id,
                    ae_id,
                    str_val(ag, "nom"),
                    str_val(ag, "adresse"),
                    str_val(ag, "telephone"),
                ],
            )
            .map_err(|e| e.to_string())?;
        }
    }

    // 3. Formules
    let mut formules_ok: i64 = 0;
    if let Some(formules) = snapshot.get("formules").and_then(|v| v.as_array()) {
        for f in formules {
            let id = str_val(f, "id");
            if id.is_empty() {
                continue;
            }
            if pending_ids.contains(id) {
                ignores += 1;
                continue;
            }
            tx.execute(
                "INSERT INTO formules_permis (id, auto_ecole_id, categorie_permis, intitule, description, heures_theoriques, heures_pratiques, prix_total, acompte_minimum, est_actif)
                 VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9, ?10)
                 ON CONFLICT(id) DO UPDATE SET
                   categorie_permis = excluded.categorie_permis,
                   intitule = excluded.intitule,
                   description = excluded.description,
                   heures_theoriques = excluded.heures_theoriques,
                   heures_pratiques = excluded.heures_pratiques,
                   prix_total = excluded.prix_total,
                   acompte_minimum = excluded.acompte_minimum,
                   est_actif = excluded.est_actif",
                params![
                    id,
                    str_val(f, "auto_ecole_id"),
                    str_val(f, "categorie_permis"),
                    str_val(f, "intitule"),
                    opt_str(f, "description"),
                    i64_val(f, "heures_theoriques", 20),
                    i64_val(f, "heures_pratiques", 20),
                    f64_val(f, "prix_total"),
                    f64_val(f, "acompte_minimum"),
                    bool_int(f, "est_actif", 1),
                ],
            )
            .map_err(|e| e.to_string())?;
            formules_ok += 1;
        }
    }

    // 4. Candidats
    let mut candidats_ok: i64 = 0;
    if let Some(candidats) = snapshot.get("candidats").and_then(|v| v.as_array()) {
        for c in candidats {
            let id = str_val(c, "id");
            if id.is_empty() {
                continue;
            }
            if pending_ids.contains(id) {
                ignores += 1;
                continue;
            }
            tx.execute(
                "INSERT INTO candidats (id, auto_ecole_id, agence_id, nom, prenoms, genre, date_naissance, lieu_naissance, nationalite, telephone_whatsapp, photo_url, type_piece_id, numero_piece_id, visite_medicale_statut)
                 VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9, ?10, ?11, ?12, ?13, ?14)
                 ON CONFLICT(id) DO UPDATE SET
                   nom = excluded.nom,
                   prenoms = excluded.prenoms,
                   genre = excluded.genre,
                   date_naissance = excluded.date_naissance,
                   lieu_naissance = excluded.lieu_naissance,
                   nationalite = excluded.nationalite,
                   telephone_whatsapp = excluded.telephone_whatsapp,
                   photo_url = excluded.photo_url,
                   type_piece_id = excluded.type_piece_id,
                   numero_piece_id = excluded.numero_piece_id,
                   visite_medicale_statut = excluded.visite_medicale_statut",
                params![
                    id,
                    str_val(c, "auto_ecole_id"),
                    str_val(c, "agence_id"),
                    str_val(c, "nom"),
                    str_val(c, "prenoms"),
                    str_val(c, "genre"),
                    str_val(c, "date_naissance"),
                    str_val(c, "lieu_naissance"),
                    str_val(c, "nationalite"),
                    str_val(c, "telephone_whatsapp"),
                    opt_str(c, "photo_url"),
                    str_val(c, "type_piece_id"),
                    str_val(c, "numero_piece_id"),
                    str_val(c, "visite_medicale_statut"),
                ],
            )
            .map_err(|e| e.to_string())?;
            candidats_ok += 1;
        }
    }

    // 5. Dossiers d'inscription
    let mut dossiers_ok: i64 = 0;
    if let Some(dossiers) = snapshot.get("dossiers").and_then(|v| v.as_array()) {
        for d in dossiers {
            let id = str_val(d, "id");
            let cand_id = str_val(d, "candidat_id");
            if id.is_empty() || cand_id.is_empty() {
                continue;
            }
            if pending_ids.contains(id) || pending_ids.contains(cand_id) {
                ignores += 1;
                continue;
            }
            let res = tx.execute(
                "INSERT INTO dossiers_inscription (id, auto_ecole_id, agence_id, candidat_id, formule_id, numero_dossier, date_inscription, montant_total, remise_accordee, solde_restant, statut_dossier, heures_code_effectuees, heures_conduite_effectuees, derogation_conduite_accordee)
                 VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9, ?10, ?11, ?12, ?13, ?14)
                 ON CONFLICT(id) DO UPDATE SET
                   formule_id = excluded.formule_id,
                   montant_total = excluded.montant_total,
                   remise_accordee = excluded.remise_accordee,
                   solde_restant = excluded.solde_restant,
                   statut_dossier = excluded.statut_dossier,
                   heures_code_effectuees = excluded.heures_code_effectuees,
                   heures_conduite_effectuees = excluded.heures_conduite_effectuees,
                   derogation_conduite_accordee = excluded.derogation_conduite_accordee",
                params![
                    id,
                    str_val(d, "auto_ecole_id"),
                    str_val(d, "agence_id"),
                    cand_id,
                    str_val(d, "formule_id"),
                    str_val(d, "numero_dossier"),
                    str_val(d, "date_inscription"),
                    f64_val(d, "montant_total"),
                    f64_val(d, "remise_accordee"),
                    f64_val(d, "solde_restant"),
                    str_val(d, "statut_dossier"),
                    i64_val(d, "heures_code_effectuees", 0),
                    i64_val(d, "heures_conduite_effectuees", 0),
                    bool_int(d, "derogation_conduite_accordee", 0),
                ],
            );
            if res.is_ok() {
                dossiers_ok += 1;
            }
        }
    }

    // 6. Paiements (assurer l'existence de l'utilisateur et de la session caisse pour PRAGMA foreign_keys)
    let mut paiements_ok: i64 = 0;
    if let Some(paiements) = snapshot.get("paiements").and_then(|v| v.as_array()) {
        for p in paiements {
            let id = str_val(p, "id");
            let dos_id = str_val(p, "dossier_id");
            if id.is_empty() || dos_id.is_empty() {
                continue;
            }
            if pending_ids.contains(id) {
                ignores += 1;
                continue;
            }
            let ae_id = str_val(p, "auto_ecole_id");
            let ag_id = str_val(p, "agence_id");
            let enc_id = str_val(p, "encaisse_par_id");
            let sess_id = str_val(p, "session_caisse_id");
            let date_p = str_val(p, "date_paiement");

            if !enc_id.is_empty() && !ae_id.is_empty() && !ag_id.is_empty() {
                let _ = tx.execute(
                    "INSERT OR IGNORE INTO utilisateurs (id, auto_ecole_id, agence_id, nom, prenoms, telephone, role)
                     VALUES (?1, ?2, ?3, 'Cloud', 'Sync', ?1, 'SECRETAIRE')",
                    params![enc_id, ae_id, ag_id],
                );
            }
            if !sess_id.is_empty() && !ag_id.is_empty() && !enc_id.is_empty() {
                let _ = tx.execute(
                    "INSERT OR IGNORE INTO sessions_caisse (id, agence_id, caissier_id, date_ouverture, statut)
                     VALUES (?1, ?2, ?3, ?4, 'FERMEE')",
                    params![sess_id, ag_id, enc_id, date_p],
                );
            }

            let res = tx.execute(
                "INSERT INTO paiements (id, auto_ecole_id, agence_id, dossier_id, session_caisse_id, montant, mode_reglement, reference_transaction, date_paiement, encaisse_par_id, est_annule)
                 VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9, ?10, ?11)
                 ON CONFLICT(id) DO UPDATE SET
                   montant = excluded.montant,
                   mode_reglement = excluded.mode_reglement,
                   est_annule = excluded.est_annule",
                params![
                    id,
                    ae_id,
                    ag_id,
                    dos_id,
                    sess_id,
                    f64_val(p, "montant"),
                    str_val(p, "mode_reglement"),
                    opt_str(p, "reference_transaction"),
                    date_p,
                    enc_id,
                    bool_int(p, "est_annule", 0),
                ],
            );
            if res.is_ok() {
                paiements_ok += 1;
            }
        }
    }

    // 7. Véhicules
    let mut vehicules_ok: i64 = 0;
    if let Some(vehicules) = snapshot.get("vehicules").and_then(|v| v.as_array()) {
        for v in vehicules {
            let id = str_val(v, "id");
            if id.is_empty() {
                continue;
            }
            if pending_ids.contains(id) {
                ignores += 1;
                continue;
            }
            let res = tx.execute(
                "INSERT INTO vehicules (id, agence_id, immatriculation, marque_modele, type_boite, kilometrage_actuel, date_expiration_assurance, date_visite_technique, statut)
                 VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9)
                 ON CONFLICT(id) DO UPDATE SET
                   marque_modele = excluded.marque_modele,
                   type_boite = excluded.type_boite,
                   kilometrage_actuel = excluded.kilometrage_actuel,
                   date_expiration_assurance = excluded.date_expiration_assurance,
                   date_visite_technique = excluded.date_visite_technique,
                   statut = excluded.statut",
                params![
                    id,
                    str_val(v, "agence_id"),
                    str_val(v, "immatriculation"),
                    str_val(v, "marque_modele"),
                    str_val(v, "type_boite"),
                    i64_val(v, "kilometrage_actuel", 0),
                    str_val(v, "date_expiration_assurance"),
                    str_val(v, "date_visite_technique"),
                    str_val(v, "statut"),
                ],
            );
            if res.is_ok() {
                vehicules_ok += 1;
            }
        }
    }

    let extra = super::sync_pull_extra::appliquer_tables_secondaires(&tx, &snapshot, &pending_ids);
    ignores += extra.ignores;

    tx.commit().map_err(|e| e.to_string())?;

    Ok(json!({
        "formules_appliquees": formules_ok,
        "candidats_appliques": candidats_ok,
        "dossiers_appliques": dossiers_ok,
        "paiements_appliques": paiements_ok,
        "vehicules_appliques": vehicules_ok,
        "utilisateurs_appliques": extra.utilisateurs,
        "sessions_caisse_appliquees": extra.sessions_caisse,
        "seances_code_appliquees": extra.seances_code,
        "sessions_examen_appliquees": extra.sessions_examen,
        "candidats_examen_appliquees": extra.candidats_examen,
        "ignores_en_attente": ignores,
    }))
}


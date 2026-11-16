use rusqlite::{params, Transaction};
use serde_json::Value;
use std::collections::HashSet;

fn s<'a>(v: &'a Value, key: &str) -> &'a str {
    v.get(key).and_then(|x| x.as_str()).unwrap_or("")
}

fn opt<'a>(v: &'a Value, key: &str) -> Option<&'a str> {
    v.get(key).and_then(|x| x.as_str()).filter(|x| !x.is_empty())
}

fn f64_v(v: &Value, key: &str) -> f64 {
    v.get(key).and_then(|x| x.as_f64()).unwrap_or(0.0)
}

fn opt_f64(v: &Value, key: &str) -> Option<f64> {
    v.get(key).and_then(|x| x.as_f64())
}

pub struct ExtraPullCounts {
    pub utilisateurs: i64,
    pub sessions_caisse: i64,
    pub seances_code: i64,
    pub sessions_examen: i64,
    pub candidats_examen: i64,
    pub ignores: i64,
}

pub fn appliquer_tables_secondaires(
    tx: &Transaction,
    snapshot: &Value,
    pending_ids: &HashSet<String>,
) -> ExtraPullCounts {
    let mut counts = ExtraPullCounts {
        utilisateurs: 0,
        sessions_caisse: 0,
        seances_code: 0,
        sessions_examen: 0,
        candidats_examen: 0,
        ignores: 0,
    };

    if let Some(users) = snapshot.get("utilisateurs").and_then(|v| v.as_array()) {
        for u in users {
            let id = s(u, "id");
            if id.is_empty() {
                continue;
            }
            if pending_ids.contains(id) {
                counts.ignores += 1;
                continue;
            }
            let res = tx.execute(
                "INSERT INTO utilisateurs (id, auto_ecole_id, agence_id, nom, prenoms, telephone, role, code_pin_manager, statut)
                 VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9)
                 ON CONFLICT(id) DO UPDATE SET
                   nom = excluded.nom,
                   prenoms = excluded.prenoms,
                   telephone = excluded.telephone,
                   role = excluded.role,
                   code_pin_manager = excluded.code_pin_manager,
                   statut = excluded.statut",
                params![
                    id,
                    s(u, "auto_ecole_id"),
                    s(u, "agence_id"),
                    s(u, "nom"),
                    s(u, "prenoms"),
                    s(u, "telephone"),
                    s(u, "role"),
                    opt(u, "code_pin_manager"),
                    s(u, "statut"),
                ],
            );
            if res.is_ok() {
                counts.utilisateurs += 1;
            }
        }
    }

    if let Some(sessions) = snapshot.get("sessions_caisse").and_then(|v| v.as_array()) {
        for sc in sessions {
            let id = s(sc, "id");
            if id.is_empty() {
                continue;
            }
            if pending_ids.contains(id) {
                counts.ignores += 1;
                continue;
            }
            let res = tx.execute(
                "INSERT INTO sessions_caisse (id, agence_id, caissier_id, date_ouverture, date_fermeture, fond_de_caisse, total_especes_theorique, total_especes_compte, ecart_caisse, statut, observations)
                 VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9, ?10, ?11)
                 ON CONFLICT(id) DO UPDATE SET
                   date_fermeture = excluded.date_fermeture,
                   total_especes_theorique = excluded.total_especes_theorique,
                   total_especes_compte = excluded.total_especes_compte,
                   ecart_caisse = excluded.ecart_caisse,
                   statut = excluded.statut,
                   observations = excluded.observations",
                params![
                    id,
                    s(sc, "agence_id"),
                    s(sc, "caissier_id"),
                    s(sc, "date_ouverture"),
                    opt(sc, "date_fermeture"),
                    f64_v(sc, "fond_de_caisse"),
                    f64_v(sc, "total_especes_theorique"),
                    opt_f64(sc, "total_especes_compte"),
                    opt_f64(sc, "ecart_caisse"),
                    s(sc, "statut"),
                    opt(sc, "observations"),
                ],
            );
            if res.is_ok() {
                counts.sessions_caisse += 1;
            }
        }
    }

    if let Some(seances) = snapshot.get("seances_code").and_then(|v| v.as_array()) {
        for sc in seances {
            let id = s(sc, "id");
            if id.is_empty() {
                continue;
            }
            if pending_ids.contains(id) {
                counts.ignores += 1;
                continue;
            }
            let res = tx.execute(
                "INSERT INTO seances_code (id, auto_ecole_id, agence_id, date_seance, heure_debut, heure_fin, theme_traite, formateur_nom, moniteur_id, type_seance, lien_visio, mot_de_passe_visio, statut_seance, statut)
                 VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9, ?10, ?11, ?12, ?13, ?13)
                 ON CONFLICT(id) DO UPDATE SET
                   date_seance = excluded.date_seance,
                   heure_debut = excluded.heure_debut,
                   heure_fin = excluded.heure_fin,
                   theme_traite = excluded.theme_traite,
                   formateur_nom = excluded.formateur_nom,
                   lien_visio = excluded.lien_visio,
                   statut_seance = excluded.statut_seance,
                   statut = excluded.statut",
                params![
                    id,
                    s(sc, "auto_ecole_id"),
                    s(sc, "agence_id"),
                    s(sc, "date_seance"),
                    s(sc, "heure_debut"),
                    s(sc, "heure_fin"),
                    s(sc, "theme_traite"),
                    s(sc, "formateur_nom"),
                    opt(sc, "moniteur_id"),
                    s(sc, "type_seance"),
                    opt(sc, "lien_visio"),
                    opt(sc, "mot_de_passe_visio"),
                    s(sc, "statut_seance"),
                ],
            );
            if res.is_ok() {
                counts.seances_code += 1;
            }
        }
    }

    if let Some(ex_sessions) = snapshot.get("sessions_examen").and_then(|v| v.as_array()) {
        for se in ex_sessions {
            let id = s(se, "id");
            if id.is_empty() {
                continue;
            }
            if pending_ids.contains(id) {
                counts.ignores += 1;
                continue;
            }
            let res = tx.execute(
                "INSERT INTO sessions_examen (id, auto_ecole_id, agence_id, code_session, date_examen, centre_examen, categorie_permis, statut_session, inspecteur_nom)
                 VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9)
                 ON CONFLICT(id) DO UPDATE SET
                   date_examen = excluded.date_examen,
                   centre_examen = excluded.centre_examen,
                   statut_session = excluded.statut_session,
                   inspecteur_nom = excluded.inspecteur_nom",
                params![
                    id,
                    s(se, "auto_ecole_id"),
                    s(se, "agence_id"),
                    s(se, "code_session"),
                    s(se, "date_examen"),
                    s(se, "centre_examen"),
                    s(se, "categorie_permis"),
                    s(se, "statut_session"),
                    opt(se, "inspecteur_nom"),
                ],
            );
            if res.is_ok() {
                counts.sessions_examen += 1;
            }
        }
    }

    if let Some(ce_list) = snapshot.get("candidats_examen").and_then(|v| v.as_array()) {
        for ce in ce_list {
            let id = s(ce, "id");
            if id.is_empty() {
                continue;
            }
            if pending_ids.contains(id) {
                counts.ignores += 1;
                continue;
            }
            let res = tx.execute(
                "INSERT INTO candidats_examen (id, session_examen_id, dossier_id, candidat_id, numero_table, resultat_code, resultat_conduite, resultat_final)
                 VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8)
                 ON CONFLICT(id) DO UPDATE SET
                   numero_table = excluded.numero_table,
                   resultat_code = excluded.resultat_code,
                   resultat_conduite = excluded.resultat_conduite,
                   resultat_final = excluded.resultat_final",
                params![
                    id,
                    s(ce, "session_examen_id"),
                    s(ce, "dossier_id"),
                    s(ce, "candidat_id"),
                    opt(ce, "numero_table"),
                    s(ce, "resultat_code"),
                    s(ce, "resultat_conduite"),
                    s(ce, "resultat_final"),
                ],
            );
            if res.is_ok() {
                counts.candidats_examen += 1;
            }
        }
    }

    counts
}


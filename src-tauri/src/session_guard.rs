use std::sync::Mutex;
use tauri::State;

#[derive(Clone, Debug)]
pub struct ActiveLocalSession {
    pub user_id: String,
    pub auto_ecole_id: String,
    pub role: String,
}

#[derive(Default)]
pub struct LocalSessionState {
    pub inner: Mutex<Option<ActiveLocalSession>>,
}

/// Normalise un libellé de rôle frontend ou backend vers le code canonique.
fn normalize_role(role: &str) -> String {
    match role.trim().to_uppercase().as_str() {
        "MANAGER" | "GERANT" | "DIRECTEUR" => "DIRECTEUR".to_string(),
        "SECRETAIRE" => "SECRETAIRE".to_string(),
        "CAISSIER" => "CAISSIER".to_string(),
        "MONITEUR" => "MONITEUR".to_string(),
        "ELEVE" | "CANDIDAT" => "CANDIDAT".to_string(),
        "SUPER_ADMIN" => "SUPER_ADMIN".to_string(),
        other => other.to_string(),
    }
}

/// Vérifie qu'une session locale est active et autorisée pour l'un des rôles fournis.
pub fn require_local_roles(
    state: &State<LocalSessionState>,
    allowed: &[&str],
) -> Result<ActiveLocalSession, String> {
    let guard = state.inner.lock().map_err(|e| e.to_string())?;
    let sess = guard
        .as_ref()
        .ok_or_else(|| "Session locale requise : veuillez vous connecter.".to_string())?;

    let user_role = normalize_role(&sess.role);
    if allowed.iter().any(|r| normalize_role(r) == user_role) {
        return Ok(sess.clone());
    }

    Err(format!(
        "Accès refusé côté local : privilèges insuffisants pour le rôle {}.",
        user_role
    ))
}

#[tauri::command]
pub fn set_local_session_context(
    user_id: String,
    auto_ecole_id: String,
    role: String,
    state: State<LocalSessionState>,
) -> Result<(), String> {
    let mut guard = state.inner.lock().map_err(|e| e.to_string())?;
    *guard = Some(ActiveLocalSession {
        user_id,
        auto_ecole_id,
        role: normalize_role(&role),
    });
    Ok(())
}

#[tauri::command]
pub fn clear_local_session_context(state: State<LocalSessionState>) -> Result<(), String> {
    let mut guard = state.inner.lock().map_err(|e| e.to_string())?;
    *guard = None;
    Ok(())
}


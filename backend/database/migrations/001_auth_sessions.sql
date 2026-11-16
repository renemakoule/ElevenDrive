-- ==============================================================================
-- AutoSchool - Migration Authentification & Sessions (Go Backend)
-- Version : 1.0.0
-- ==============================================================================

-- Table des sessions utilisateurs (gestion des refresh tokens & traçabilité)
CREATE TABLE IF NOT EXISTS auth_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES utilisateurs(id) ON DELETE CASCADE,
    auto_ecole_id UUID NOT NULL REFERENCES auto_ecoles(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    ip_address VARCHAR(45),
    user_agent TEXT,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Index pour recherches rapides et validation des sessions
CREATE INDEX IF NOT EXISTS idx_auth_sessions_user_id ON auth_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_auth_sessions_token_hash ON auth_sessions(token_hash);
CREATE INDEX IF NOT EXISTS idx_auth_sessions_expires_at ON auth_sessions(expires_at);

-- Trigger pour la mise à jour automatique de updated_at
CREATE TRIGGER trg_auth_sessions_updated_at
BEFORE UPDATE ON auth_sessions
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Activation de Row Level Security (RLS) sur les sessions
ALTER TABLE auth_sessions ENABLE ROW LEVEL SECURITY;

CREATE POLICY auth_sessions_isolation_policy ON auth_sessions
    FOR ALL USING (
        auto_ecole_id = NULLIF(current_setting('app.current_auto_ecole_id', true), '')::UUID 
        OR current_setting('app.is_admin', true) = 'true'
    );

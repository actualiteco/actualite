CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    user_id VARCHAR NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token_hash VARCHAR(64) NOT NULL UNIQUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,

    user_agent TEXT,
    ip_address INET,

    CONSTRAINT sessions_expires_after_creation
        CHECK (expires_at > created_at),

    CONSTRAINT sessions_revoked_after_creation
        CHECK (revoked_at IS NULL OR revoked_at >= created_at)
);

CREATE INDEX idx_sessions_user_id
    ON sessions(user_id);

CREATE INDEX idx_sessions_expires_at
    ON sessions(expires_at);

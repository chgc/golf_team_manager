CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    player_id TEXT,
    display_name TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('manager', 'player')),
    auth_provider TEXT NOT NULL CHECK (auth_provider IN ('dev_stub', 'line')),
    provider_subject TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (player_id) REFERENCES players (id) ON DELETE SET NULL,
    UNIQUE (auth_provider, provider_subject)
);

CREATE INDEX IF NOT EXISTS idx_users_role ON users (role);
CREATE INDEX IF NOT EXISTS idx_users_player_id ON users (player_id);

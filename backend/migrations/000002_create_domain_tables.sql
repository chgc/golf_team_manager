CREATE TABLE IF NOT EXISTS players (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    handicap REAL NOT NULL CHECK (handicap >= 0 AND handicap <= 54),
    phone TEXT,
    email TEXT,
    status TEXT NOT NULL CHECK (status IN ('active', 'inactive')),
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_players_status ON players (status);
CREATE INDEX IF NOT EXISTS idx_players_name ON players (name);

CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    session_date TEXT NOT NULL,
    course_name TEXT NOT NULL,
    course_address TEXT,
    max_players INTEGER NOT NULL CHECK (max_players > 0),
    registration_deadline TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('open', 'closed', 'confirmed', 'completed', 'cancelled')),
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (julianday(registration_deadline) <= julianday(session_date))
);

CREATE INDEX IF NOT EXISTS idx_sessions_status ON sessions (status);
CREATE INDEX IF NOT EXISTS idx_sessions_date ON sessions (session_date);

CREATE TABLE IF NOT EXISTS registrations (
    id TEXT PRIMARY KEY,
    player_id TEXT NOT NULL,
    session_id TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('confirmed', 'cancelled')),
    registered_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (player_id) REFERENCES players (id) ON DELETE RESTRICT,
    FOREIGN KEY (session_id) REFERENCES sessions (id) ON DELETE RESTRICT,
    UNIQUE (player_id, session_id)
);

CREATE INDEX IF NOT EXISTS idx_registrations_session_status ON registrations (session_id, status);
CREATE INDEX IF NOT EXISTS idx_registrations_player_status ON registrations (player_id, status);

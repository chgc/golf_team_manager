CREATE UNIQUE INDEX IF NOT EXISTS idx_users_player_id_unique
ON users (player_id)
WHERE player_id IS NOT NULL;

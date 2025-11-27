CREATE TABLE IF NOT EXISTS conversation_participants (
    id VARCHAR(26) PRIMARY KEY,
    conversation_id VARCHAR(26) NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    user_id VARCHAR(26) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
    left_at TIMESTAMPTZ DEFAULT NULL,
    UNIQUE(conversation_id, user_id)
);
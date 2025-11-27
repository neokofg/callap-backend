CREATE TABLE IF NOT EXISTS messages (
    id VARCHAR(26) PRIMARY KEY,
    conversation_id VARCHAR(26) NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender_id VARCHAR(26) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    message_type VARCHAR(20) DEFAULT 'text' CHECK (message_type IN ('text', 'image', 'file', 'system')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
    is_read BOOLEAN DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_messages_conversation_created ON messages (conversation_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_messages_sender_created ON messages (sender_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_messages_is_read ON messages (is_read);
CREATE TABLE IF NOT EXISTS friends (
    id VARCHAR(26) PRIMARY KEY,
    user_id VARCHAR(26) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    friend_id VARCHAR(26) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
        CHECK (status IN ('pending', 'accepted', 'rejected', 'blocked')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
    UNIQUE(user_id, friend_id)
);

CREATE INDEX IF NOT EXISTS idx_friends_user_status ON friends (user_id, status);
CREATE INDEX IF NOT EXISTS idx_friends_friend_status ON friends (friend_id, status);
CREATE INDEX IF NOT EXISTS idx_friends_created_at ON friends (created_at DESC);
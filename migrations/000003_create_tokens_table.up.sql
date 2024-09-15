CREATE TABLE IF NOT EXISTS tokens (
    hash bytea PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users ON DELETE CASCADE,
    expiry timestamp(0) with time zone NOT NULL,
    scope TEXT NOT NULL,
    created_at timestamp(8) WITH TIME ZONE NOT NULL DEFAULT NOW()
);
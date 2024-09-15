CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    activated BOOL NOT NULL,
    created_at timestamp(8) WITH TIME ZONE NOT NULL DEFAULT NOW()
)
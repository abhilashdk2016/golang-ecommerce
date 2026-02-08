CREATE TABLE refresh_tokens(
    id serial PRIMARY KEY,
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token varchar(500) UNIQUE NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);

CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);

CREATE INDEX idx_refresh_tokens_deleted_at ON refresh_tokens(deleted_at);


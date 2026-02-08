CREATE TABLE carts(
    id serial PRIMARY KEY,
    user_id integer UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);

CREATE INDEX idx_carts_user_id ON carts(user_id);

CREATE INDEX idx_carts_deleted_at ON carts(deleted_at);


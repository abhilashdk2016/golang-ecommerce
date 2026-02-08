CREATE TABLE categories(
    id serial PRIMARY KEY,
    name varchar(255) NOT NULL,
    description text,
    is_active boolean DEFAULT TRUE,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);

CREATE INDEX idx_categories_is_active ON categories(is_active);

CREATE INDEX idx_categories_deleted_at ON categories(deleted_at);


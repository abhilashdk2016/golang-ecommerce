CREATE TABLE products(
    id serial PRIMARY KEY,
    category_id integer NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    name varchar(255) NOT NULL,
    description text,
    price DECIMAL(10, 2) NOT NULL,
    stock integer DEFAULT 0,
    sku varchar(100) UNIQUE NOT NULL,
    is_active boolean DEFAULT TRUE,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);

CREATE INDEX idx_products_category_id ON products(category_id);

CREATE INDEX idx_products_sku ON products(sku);

CREATE INDEX idx_products_is_active ON products(is_active);

CREATE INDEX idx_products_deleted_at ON products(deleted_at);


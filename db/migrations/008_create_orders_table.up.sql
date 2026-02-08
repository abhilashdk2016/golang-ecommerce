CREATE TYPE order_status AS ENUM(
    'pending',
    'confirmed',
    'shipped',
    'delivered',
    'cancelled'
);

CREATE TABLE orders(
    id serial PRIMARY KEY,
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status order_status DEFAULT 'pending',
    total_amount DECIMAL(10, 2) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);

CREATE INDEX idx_orders_user_id ON orders(user_id);

CREATE INDEX idx_orders_status ON orders(status);

CREATE INDEX idx_orders_deleted_at ON orders(deleted_at);


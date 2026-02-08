CREATE TYPE user_role AS ENUM(
    'customer',
    'admin'
);

CREATE TABLE users(
    id serial PRIMARY KEY,
    email varchar(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name varchar(100) NOT NULL,
    last_name varchar(100) NOT NULL,
    phone varchar(20),
    is_active boolean DEFAULT TRUE,
    ROLE user_role DEFAULT 'customer',
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);

CREATE INDEX idx_users_email ON users(email);

CREATE INDEX idx_users_deleted_at ON users(deleted_at);


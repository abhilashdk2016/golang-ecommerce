CREATE TABLE product_images(
    id serial PRIMARY KEY,
    product_id integer NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    url varchar(500) NOT NULL,
    alt_text varchar(255),
    is_primary boolean DEFAULT FALSE,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);

CREATE INDEX idx_product_images_product_id ON product_images(product_id);

CREATE INDEX idx_product_images_is_primary ON product_images(is_primary);

CREATE INDEX idx_product_images_deleted_at ON product_images(deleted_at);


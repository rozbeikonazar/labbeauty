CREATE TABLE IF NOT EXISTS services (
    id bigserial PRIMARY KEY,
    time smallint,
    description TEXT NOT NULL,
    price integer NOT NULL, 
    category_id bigint REFERENCES categories (id) ON DELETE CASCADE,
    subcategory_id bigint REFERENCES subcategories (id) ON DELETE CASCADE
);

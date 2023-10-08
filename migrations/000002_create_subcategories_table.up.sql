CREATE TABLE  IF NOT EXISTS subcategories (
    id bigserial PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);
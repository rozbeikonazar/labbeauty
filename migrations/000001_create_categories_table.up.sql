CREATE TABLE IF NOT EXISTS categories (
    id bigserial PRIMARY KEY,
    title TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL, 
    photo_url TEXT NOT NULL UNIQUE 
);
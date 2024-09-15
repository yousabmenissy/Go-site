CREATE TABLE IF NOT EXISTS books (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    release_date DATE NOT NULL,
    publisher TEXT NOT NULL,
    country TEXT NOT NULL,
    isbn TEXT NOT NULL,
    pages INT NOT NULL,
    language TEXT NOT NULL,
    description TEXT NOT NULL,
    genres TEXT [] NOT NULL,
    cover_image_url TEXT NOT NULL,
    price NUMERIC(10, 2) NOT NULL
);
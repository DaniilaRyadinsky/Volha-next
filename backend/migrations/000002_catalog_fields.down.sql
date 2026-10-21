CREATE TABLE articles (
    id VARCHAR(100) PRIMARY KEY,
    title VARCHAR(300) NOT NULL,
    img VARCHAR(255) NOT NULL,
    text TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE products
    DROP COLUMN pricesale,
    DROP COLUMN is_featured;

ALTER TABLE slides
    DROP COLUMN badge,
    DROP COLUMN description,
    DROP COLUMN title;

ALTER TABLE categories
    DROP COLUMN description;

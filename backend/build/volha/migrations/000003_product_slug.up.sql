ALTER TABLE products
    ADD COLUMN slug VARCHAR(100);

UPDATE products
SET slug = article;

ALTER TABLE products
    ALTER COLUMN slug SET NOT NULL,
    ADD CONSTRAINT products_slug_key UNIQUE (slug);

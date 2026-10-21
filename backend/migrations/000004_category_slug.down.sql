ALTER TABLE categories
    RENAME CONSTRAINT categories_slug_key TO categories_uri_key;

ALTER TABLE categories
    RENAME COLUMN slug TO uri;

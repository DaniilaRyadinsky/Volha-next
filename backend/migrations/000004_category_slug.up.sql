ALTER TABLE categories
    RENAME COLUMN uri TO slug;

ALTER TABLE categories
    RENAME CONSTRAINT categories_uri_key TO categories_slug_key;

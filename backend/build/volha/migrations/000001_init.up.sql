CREATE TABLE brands (
    id VARCHAR(100) PRIMARY KEY,
    title VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE categories (
    id VARCHAR(100) PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    uri VARCHAR(100) NOT NULL UNIQUE,
    img VARCHAR(255) NOT NULL DEFAULT ''
);

CREATE TABLE countries (
    id VARCHAR(100) PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    friendly VARCHAR(100) NOT NULL DEFAULT ''
);

CREATE TABLE materials (
    id VARCHAR(100) PRIMARY KEY,
    title VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE colors (
    id VARCHAR(100) PRIMARY KEY,
    title VARCHAR(100) NOT NULL UNIQUE,
    hex VARCHAR(7) NOT NULL
);

CREATE TABLE products (
    id VARCHAR(100) PRIMARY KEY,
    title VARCHAR(100) NOT NULL UNIQUE,
    article VARCHAR(100) NOT NULL UNIQUE,
    brand_id VARCHAR(100) NOT NULL REFERENCES brands(id),
    category_id VARCHAR(100) NOT NULL REFERENCES categories(id),
    country_id VARCHAR(100) NOT NULL REFERENCES countries(id),
    width INT NOT NULL CHECK (width > 0),
    height INT NOT NULL CHECK (height > 0),
    depth INT NOT NULL CHECK (depth > 0),
    photos VARCHAR(255)[] NOT NULL DEFAULT '{}',
    price INT NOT NULL CHECK (price > 0),
    description TEXT NOT NULL DEFAULT ''
);

CREATE TABLE product_materials (
    product_id VARCHAR(100) NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    material_id VARCHAR(100) NOT NULL REFERENCES materials(id),
    PRIMARY KEY (product_id, material_id)
);

CREATE TABLE product_colors (
    product_id VARCHAR(100) NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    color_id VARCHAR(100) NOT NULL REFERENCES colors(id),
    PRIMARY KEY (product_id, color_id)
);

CREATE TABLE product_similarities (
    product_id VARCHAR(100) NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    similar_product_id VARCHAR(100) NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, similar_product_id),
    CHECK (product_id <> similar_product_id)
);

CREATE TABLE product_color_photos (
    product_id VARCHAR(100) NOT NULL,
    color_id VARCHAR(100) NOT NULL,
    photos VARCHAR(255)[] NOT NULL DEFAULT '{}',
    PRIMARY KEY (product_id, color_id),
    FOREIGN KEY (product_id, color_id)
        REFERENCES product_colors(product_id, color_id)
        ON DELETE CASCADE
);

CREATE TABLE slides (
    id VARCHAR(100) PRIMARY KEY,
    link VARCHAR(255) NOT NULL,
    img VARCHAR(255) NOT NULL,
    img762 VARCHAR(255) NOT NULL
);

CREATE TABLE articles (
    id VARCHAR(100) PRIMARY KEY,
    title VARCHAR(300) NOT NULL,
    img VARCHAR(255) NOT NULL,
    text TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_products_brand_id ON products(brand_id);
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_country_id ON products(country_id);
CREATE INDEX idx_products_price ON products(price);
CREATE INDEX idx_product_materials_material_product ON product_materials(material_id, product_id);
CREATE INDEX idx_product_colors_color_product ON product_colors(color_id, product_id);
CREATE INDEX idx_product_similarities_similar_product ON product_similarities(similar_product_id, product_id);

CREATE TABLE articles (
 id VARCHAR(100) PRIMARY KEY,
 slug VARCHAR(100) NOT NULL UNIQUE,
 title VARCHAR(300) NOT NULL,
 img VARCHAR(255)[] NOT NULL DEFAULT '{}',
 pdf VARCHAR(255) NOT NULL,
 created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_articles_created_at ON articles(created_at DESC, id);

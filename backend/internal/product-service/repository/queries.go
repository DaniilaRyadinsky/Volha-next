package repository

const productSelectQuery = `
SELECT
	p.id,
	p.title,
	p.article,
	p.slug,
	p.width,
	p.height,
	p.depth,
	p.photos,
	p.price,
	p.description,
	p.is_featured,
	p.pricesale,

	b.id      AS brand_id,
	b.title   AS brand_title,

	cat.id    AS category_id,
	cat.slug  AS category_slug,
	cat.title AS category_title,
	cat.description AS category_description,
	(SELECT COUNT(*)::int FROM products cp WHERE cp.category_id = cat.id) AS category_count,
	cat.img   AS category_img,

	co.id       AS country_id,
	co.title    AS country_title,
	co.friendly AS country_friendly,

	COALESCE(
		(
			SELECT json_agg(
				json_build_object(
					'id', m.id,
					'title', m.title
				)
				ORDER BY m.title, m.id
			)
			FROM product_materials pm
			JOIN materials m ON m.id = pm.material_id
			WHERE pm.product_id = p.id
		),
		'[]'::json
	) AS materials,

	COALESCE(
		(
			SELECT json_agg(
				json_build_object(
					'id', c.id,
					'title', c.title,
					'hex', c.hex
				)
				ORDER BY c.title, c.id
			)
			FROM product_colors pc
			JOIN colors c ON c.id = pc.color_id
			WHERE pc.product_id = p.id
		),
		'[]'::json
	) AS colors,

	COALESCE(
		(
			SELECT json_agg(
				json_build_object(
					'id', sp.id,
					'title', sp.title,
					'article', sp.article,
					'slug', sp.slug,
					'width', sp.width,
					'height', sp.height,
					'depth', sp.depth,
					'photos', sp.photos,
					'price', sp.price,
					'description', sp.description,
					'isFeatured', sp.is_featured,
					'pricesale', sp.pricesale
				)
				ORDER BY sp.id
			)
			FROM product_similarities ps
			JOIN products sp ON sp.id = ps.similar_product_id
			WHERE ps.product_id = p.id
		),
		'[]'::json
	) AS similar_products
FROM products p
JOIN brands     b   ON b.id   = p.brand_id
JOIN categories cat ON cat.id = p.category_id
JOIN countries  co  ON co.id  = p.country_id
`

const productListSelectQuery = `
SELECT
	p.id,
	p.title,
	p.article,
	p.slug,
	p.width,
	p.height,
	p.depth,
	p.photos,
	p.price,
	p.description,
	p.is_featured,
	p.pricesale,

	b.id      AS brand_id,
	b.title   AS brand_title,

	cat.id    AS category_id,
	cat.slug  AS category_slug,
	cat.title AS category_title,
	cat.description AS category_description,
	(SELECT COUNT(*)::int FROM products cp WHERE cp.category_id = cat.id) AS category_count,
	cat.img   AS category_img,

	co.id       AS country_id,
	co.title    AS country_title,
	co.friendly AS country_friendly,

	COALESCE(
		(
			SELECT json_agg(
				json_build_object(
					'id', m.id,
					'title', m.title
				)
				ORDER BY m.title, m.id
			)
			FROM product_materials pm
			JOIN materials m ON m.id = pm.material_id
			WHERE pm.product_id = p.id
		),
		'[]'::json
	) AS materials,

	COALESCE(
		(
			SELECT json_agg(
				json_build_object(
					'id', c.id,
					'title', c.title,
					'hex', c.hex
				)
				ORDER BY c.title, c.id
			)
			FROM product_colors pc
			JOIN colors c ON c.id = pc.color_id
			WHERE pc.product_id = p.id
		),
		'[]'::json
	) AS colors,

	COALESCE(
		(
			SELECT json_agg(
				json_build_object(
					'id', sp.id,
					'title', sp.title,
					'article', sp.article,
					'slug', sp.slug,
					'width', sp.width,
					'height', sp.height,
					'depth', sp.depth,
					'photos', sp.photos,
					'price', sp.price,
					'description', sp.description,
					'isFeatured', sp.is_featured,
					'pricesale', sp.pricesale
				)
				ORDER BY sp.id
			)
			FROM product_similarities ps
			JOIN products sp ON sp.id = ps.similar_product_id
			WHERE ps.product_id = p.id
		),
		'[]'::json
	) AS similar_products
FROM products p
JOIN brands     b   ON b.id   = p.brand_id
JOIN categories cat ON cat.id = p.category_id
JOIN countries  co  ON co.id  = p.country_id
`

const productCountQuery = `
SELECT COUNT(*)
FROM products p
`

const productCatalogSelectQuery = `
SELECT
	p.id,
	p.slug,
	cat.title,
	p.title,
	p.price,
	COALESCE(p.photos[1], ''),
	p.pricesale,
	p.is_featured
FROM products p
JOIN categories cat ON cat.id = p.category_id
`

const createProductQuery = `
WITH
	ins_product AS (
		INSERT INTO products (
			id, title, article, brand_id, category_id, country_id,
			width, height, depth, photos, price, description, is_featured, pricesale, slug
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
		RETURNING id
	),
	ins_materials AS (
		INSERT INTO product_materials (product_id, material_id)
		SELECT ins_product.id, material_id
		FROM ins_product
		CROSS JOIN unnest($16::varchar[]) AS t(material_id)
	),
	ins_colors AS (
		INSERT INTO product_colors (product_id, color_id)
		SELECT ins_product.id, color_id
		FROM ins_product
		CROSS JOIN unnest($17::varchar[]) AS t(color_id)
	),
	ins_similarities AS (
		INSERT INTO product_similarities (product_id, similar_product_id)
		SELECT ins_product.id, similar_product_id
		FROM ins_product
		CROSS JOIN unnest($18::varchar[]) AS t(similar_product_id)
	)
SELECT COUNT(*) FROM ins_product;
`

const updateProductQuery = `
WITH
	upd_product AS (
		UPDATE products SET
			title       = $2,
			article     = $3,
			brand_id    = $4,
			category_id = $5,
			country_id  = $6,
			width       = $7,
			height      = $8,
			depth       = $9,
			photos      = $10,
			price       = $11,
			description = $12,
			is_featured = $13,
			pricesale = $14,
			slug = $15
		WHERE id = $1
		RETURNING id
	),
	ins_materials AS (
		INSERT INTO product_materials (product_id, material_id)
		SELECT $1, m_id
		FROM unnest($16::varchar[]) AS t(m_id)
		WHERE EXISTS (SELECT 1 FROM upd_product)
		ON CONFLICT (product_id, material_id) DO NOTHING
	),
	del_materials AS (
		DELETE FROM product_materials pm
		WHERE pm.product_id = $1
		  AND NOT (pm.material_id = ANY($16::varchar[]))
		  AND EXISTS (SELECT 1 FROM upd_product)
	),
	ins_colors AS (
		INSERT INTO product_colors (product_id, color_id)
		SELECT $1, c_id
		FROM unnest($17::varchar[]) AS t(c_id)
		WHERE EXISTS (SELECT 1 FROM upd_product)
		ON CONFLICT (product_id, color_id) DO NOTHING
	),
	del_colors AS (
		DELETE FROM product_colors pc
		WHERE pc.product_id = $1
		  AND NOT (pc.color_id = ANY($17::varchar[]))
		  AND EXISTS (SELECT 1 FROM upd_product)
	),
	ins_similarities AS (
		INSERT INTO product_similarities (product_id, similar_product_id)
		SELECT $1, s_id
		FROM unnest($18::varchar[]) AS t(s_id)
		WHERE EXISTS (SELECT 1 FROM upd_product)
		ON CONFLICT (product_id, similar_product_id) DO NOTHING
	),
	del_similarities AS (
		DELETE FROM product_similarities ps
		WHERE ps.product_id = $1
		  AND NOT (ps.similar_product_id = ANY($18::varchar[]))
		  AND EXISTS (SELECT 1 FROM upd_product)
	)
SELECT COUNT(*) FROM upd_product;
`

const getDicByCatQuery = `
SELECT * FROM (
SELECT DISTINCT 'brand' as type, b.id, b.title, '' as extra1, '' as extra2, '' as extra3, '' as extra4
FROM brands b
JOIN products p ON b.id = p.brand_id
WHERE p.category_id = $1

UNION ALL

SELECT 'category', c.id, c.title, c.slug, c.img, c.description, (SELECT COUNT(*)::text FROM products p WHERE p.category_id = c.id) FROM categories c

UNION ALL

SELECT DISTINCT 'country' as type, c.id, c.title, c.friendly, '', '', ''
FROM countries c
JOIN products p ON c.id = p.country_id
WHERE p.category_id = $1

UNION ALL

SELECT DISTINCT 'material' as type, m.id, m.title, '', '', '', ''
FROM materials m
JOIN product_materials pm ON m.id = pm.material_id
JOIN products p ON pm.product_id = p.id
WHERE p.category_id = $1

UNION ALL

SELECT DISTINCT 'color' as type, c.id, c.title, c.hex, '', '', ''
FROM colors c
JOIN product_colors pc ON c.id = pc.color_id
JOIN products p ON pc.product_id = p.id
WHERE p.category_id = $1

UNION ALL

SELECT
	'stats',
	'',
	COALESCE(MIN(price)::text, '0'),
	COALESCE(MAX(price)::text, '0'),
	COALESCE(MIN(width)::text, '0') || ',' || COALESCE(MAX(width)::text, '0') || ',' ||
	COALESCE(MIN(height)::text, '0') || ',' || COALESCE(MAX(height)::text, '0') || ',' ||
	COALESCE(MIN(depth)::text, '0') || ',' || COALESCE(MAX(depth)::text, '0'),
	'',
	''
FROM products
WHERE category_id = $1
) AS dictionaries
ORDER BY
	CASE type
		WHEN 'brand' THEN 1
		WHEN 'category' THEN 2
		WHEN 'country' THEN 3
		WHEN 'material' THEN 4
		WHEN 'color' THEN 5
		WHEN 'stats' THEN 6
		ELSE 7
	END,
	title,
	id;
`

const getDicQuery = `
WITH dicts AS (
	SELECT 'brand' as type, id, title, '' as extra1, '' as extra2, '' as extra3, '' as extra4 FROM brands
	UNION ALL
	SELECT 'category', c.id, c.title, c.slug, c.img, c.description, (SELECT COUNT(*)::text FROM products p WHERE p.category_id = c.id) FROM categories c
	UNION ALL
	SELECT 'country', id, title, friendly, '', '', '' FROM countries
	UNION ALL
	SELECT 'material', id, title, '', '', '', '' FROM materials
	UNION ALL
	SELECT 'color', id, title, hex, '', '', '' FROM colors
),
stats AS (
	SELECT
		COALESCE(MIN(price)::text, '0') AS min_price,
		COALESCE(MAX(price)::text, '0') AS max_price,
		COALESCE(MIN(width)::text, '0') AS min_width,
		COALESCE(MAX(width)::text, '0') AS max_width,
		COALESCE(MIN(height)::text, '0') AS min_height,
		COALESCE(MAX(height)::text, '0') AS max_height,
		COALESCE(MIN(depth)::text, '0') AS min_depth,
		COALESCE(MAX(depth)::text, '0') AS max_depth
	FROM products
)
SELECT * FROM (
SELECT * FROM dicts
UNION ALL
SELECT 'stats', '', min_price, max_price, min_width || ',' || max_width || ',' || min_height || ',' || max_height || ',' || min_depth || ',' || max_depth, '', '' FROM stats
) AS dictionaries
ORDER BY
	CASE type
		WHEN 'brand' THEN 1
		WHEN 'category' THEN 2
		WHEN 'country' THEN 3
		WHEN 'material' THEN 4
		WHEN 'color' THEN 5
		WHEN 'stats' THEN 6
		ELSE 7
	END,
	title,
	id;
`

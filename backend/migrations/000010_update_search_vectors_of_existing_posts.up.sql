UPDATE posts
SET search_vector = setweight(to_tsvector('pg_catalog.simple', coalesce(content, '')), 'A')
WHERE search_vector IS NULL OR search_vector = ''::tsvector;

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_users_username_trgm ON users USING GIN (username gin_trgm_ops);
CREATE INDEX idx_users_display_name_trgm ON users USING GIN (display_name gin_trgm_ops);

CREATE TRIGGER posts_search_vector_update
    BEFORE INSERT OR UPDATE OF content ON posts
    FOR EACH ROW EXECUTE FUNCTION tsvector_update_trigger(search_vector, 'pg_catalog.simple', content);

UPDATE posts
SET search_vector = to_tsvector('pg_catalog.simple', content);

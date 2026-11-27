DROP TRIGGER IF EXISTS posts_search_vector_update ON posts;
DROP INDEX IF EXISTS idx_users_display_name_trgm;
DROP INDEX IF EXISTS idx_users_username_trgm;

CREATE OR REPLACE FUNCTION posts_on_soft_delete()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF NEW.deleted_at IS NOT NULL AND OLD.deleted_at IS NULL THEN
        DELETE FROM post_hashtags WHERE post_id = NEW.id;
    END IF;
    RETURN NEW;
END $$;

CREATE TRIGGER trg_posts_soft_delete
    AFTER UPDATE OF deleted_at ON posts
    FOR EACH ROW
    EXECUTE FUNCTION posts_on_soft_delete();

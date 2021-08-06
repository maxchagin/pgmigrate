CREATE OR REPLACE FUNCTION article_change()
  RETURNS trigger AS
$$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        UPDATE core.articles SET is_changed = true WHERE id = OLD.article_id;
        RETURN OLD;
    ELSE
        UPDATE core.articles SET is_changed = true WHERE id = NEW.article_id;
        RETURN NEW;
    END IF;
END;
$$
LANGUAGE plpgsql;


-- To add/delete tag bindings to article
CREATE TRIGGER articles_tags_change
  BEFORE INSERT OR DELETE
  ON articles_tags
  FOR EACH ROW
  EXECUTE PROCEDURE article_change();
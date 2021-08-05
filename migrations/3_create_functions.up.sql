--
-- Get list tags associated with an article (by article_id)
--
CREATE OR REPLACE FUNCTION tags(idx integer)
    RETURNS jsonb
    LANGUAGE plpgsql
AS
$function$
BEGIN
    RETURN
        (
            SELECT jsonb_strip_nulls(jsonb_agg(t))
            FROM (
                SELECT *
                FROM tags as tg
                        INNER JOIN articles_tags AS ats ON (tg.id = ats.tag_id)
                WHERE ats.article_id = idx
                ORDER BY ats.id ASC
            ) t
        );
END;
$function$
;

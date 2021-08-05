BEGIN;
--
-- Table articles
--
CREATE TABLE IF NOT EXISTS articles
(
    "id"              SERIAL  NOT NULL PRIMARY KEY,
    "title"           text,
    "url"             text,
    "text"            text,
    "is_changed"      boolean DEFAULT false,
    UNIQUE (id),
    UNIQUE (url)
);


--
-- Table tags
--
CREATE TABLE IF NOT EXISTS tags
(
    "id"                    SERIAL  NOT NULL PRIMARY KEY,
    "name"                  text    NOT NULL UNIQUE CHECK (name <> ''),
    UNIQUE (id)
);

--
-- Table relationship article and tag
--
CREATE TABLE IF NOT EXISTS articles_tags
(
    "id"         SERIAL  NOT NULL PRIMARY KEY,
    "article_id" integer NOT NULL references articles on DELETE CASCADE,
    "tag_id"     integer NOT NULL references tags on DELETE CASCADE,
    UNIQUE (id),
    UNIQUE (article_id, tag_id)
);

COMMIT;
-- INIT ALL
CREATE EXTENSION IF NOT EXISTS citext;
CREATE SCHEMA tpdb;

CREATE UNLOGGED TABLE "Service"
(
    id SERIAL PRIMARY KEY,
    forums INT DEFAULT 0,
    posts INT DEFAULT 0,
    threads INT DEFAULT 0,
    users INT DEFAULT 0
);
INSERT INTO "Service" DEFAULT VALUES;

CREATE UNLOGGED TABLE tpdb."User"
(
    id SERIAL PRIMARY KEY,
    nickname CITEXT COLLATE "C" UNIQUE NOT NULL,
    fullname TEXT NOT NULL,
    about TEXT,
    email CITEXT UNIQUE NOT NULL
);

CREATE UNLOGGED TABLE tpdb."Forum"
(
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    "user" CITEXT REFERENCES tpdb."User"(nickname) NOT NULL,
    slug CITEXT UNIQUE NOT NULL,
    posts INT DEFAULT 0,
    threads INT DEFAULT 0
);

CREATE UNLOGGED TABLE tpdb."Thread"
(
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    author CITEXT REFERENCES tpdb."User"(nickname) NOT NULL,
    forum CITEXT REFERENCES tpdb."Forum"(slug) NOT NULL,
    message TEXT NOT NULL,
    votes INT DEFAULT 0,
    slug CITEXT,
    created TIMESTAMP WITH TIME ZONE
);

CREATE UNLOGGED TABLE tpdb."Post"
(
    id SERIAL PRIMARY KEY,
    parent INT DEFAULT 0,
    author CITEXT REFERENCES tpdb."User"(nickname) NOT NULL,
    message TEXT NOT NULL,
    isEdited bool NOT NULL DEFAULT FALSE,
    forum CITEXT REFERENCES tpdb."Forum"(slug) NOT NULL,
    thread INT REFERENCES tpdb."Thread"(id) NOT NULL,
    created  TIMESTAMP WITH TIME ZONE,
--     path INT[] DEFAULT [0]
    path INT[] NOT NULL
);

CREATE UNLOGGED TABLE tpdb."Vote"
(
    id SERIAL PRIMARY KEY,
    threadid INT REFERENCES tpdb."Thread"(id) NOT NULL,
--     threadslug CITEXT REFERENCES tpdb."Thread"(Slug) NOT NULL,
    "user" CITEXT REFERENCES tpdb."User"(nickname) NOT NULL,
    vote INT NOT NULL,
    UNIQUE (threadid, "user")
);


-- SERVICE USERS
CREATE FUNCTION service_user_inc() RETURNS TRIGGER AS $user_inc$
BEGIN
    UPDATE "Service" SET users = users + 1 WHERE id = 1;
    RETURN NULL;
END;
$user_inc$ LANGUAGE plpgsql;
CREATE TRIGGER user_inc AFTER INSERT ON tpdb."User" EXECUTE PROCEDURE service_user_inc();

-- SERVICE THREADS
CREATE FUNCTION service_thread_inc() RETURNS TRIGGER AS $thread_inc$
BEGIN
    UPDATE "Service" SET threads = threads + 1 WHERE id = 1;
    RETURN NULL;
END;
$thread_inc$ LANGUAGE plpgsql;
CREATE TRIGGER thread_inc AFTER INSERT ON tpdb."Thread" EXECUTE PROCEDURE service_thread_inc();

-- SERVICE POSTS
CREATE FUNCTION service_post_inc() RETURNS TRIGGER AS $post_inc$
DECLARE
    record_count integer;
BEGIN
    SELECT COUNT(*) FROM newtbl INTO record_count;
    UPDATE "Service" SET posts = posts + record_count WHERE id = 1;
    RETURN NULL;
END;
$post_inc$ LANGUAGE plpgsql;
CREATE TRIGGER post_inc AFTER INSERT ON tpdb."Post"
    REFERENCING NEW TABLE as newtbl
    FOR EACH STATEMENT EXECUTE PROCEDURE service_post_inc();

-- SERVICE FORUMS
CREATE FUNCTION service_forum_inc() RETURNS TRIGGER AS $forum_inc$
BEGIN
    UPDATE "Service" SET forums = forums + 1 WHERE id = 1;
    RETURN NULL;
END;
$forum_inc$ LANGUAGE plpgsql;
CREATE TRIGGER forum_inc AFTER INSERT ON tpdb."Forum" EXECUTE PROCEDURE service_forum_inc();

-- POST INSERT
CREATE FUNCTION post_insert() RETURNS TRIGGER AS $post_insert$
DECLARE
    prevpath INT[];
BEGIN
    --     PATH
    IF new.parent != 0 THEN
        SELECT path FROM tpdb."Post" WHERE id = new.parent INTO prevpath;
        new.path := array_append(prevpath, new.id);
    ELSE
        new.path[1] := new.id;
    END IF;
    RETURN new;
END
$post_insert$ LANGUAGE plpgsql;
CREATE TRIGGER post_insert BEFORE INSERT ON tpdb."Post" FOR EACH ROW EXECUTE PROCEDURE post_insert();

-- POST INSERT FORUM
CREATE FUNCTION post_insert_forum() RETURNS TRIGGER AS $post_insert_forum$
BEGIN
    UPDATE tpdb."Forum" SET posts = posts + 1 WHERE slug = new.forum;
    RETURN new;
END
$post_insert_forum$ LANGUAGE plpgsql;
CREATE TRIGGER post_insert_forum AFTER INSERT ON tpdb."Post" FOR EACH ROW EXECUTE PROCEDURE post_insert_forum();

-- THREAD INSERT FORUM
CREATE FUNCTION thread_insert_forum() RETURNS TRIGGER AS $thread_insert_forum$
BEGIN
    UPDATE tpdb."Forum" SET threads = threads + 1 WHERE slug = new.forum;
    RETURN new;
END
$thread_insert_forum$ LANGUAGE plpgsql;
CREATE TRIGGER thread_insert_forum AFTER INSERT ON tpdb."Thread" FOR EACH ROW EXECUTE PROCEDURE thread_insert_forum();

-- FORUM INSTERT
CREATE FUNCTION forum_insert() RETURNS TRIGGER AS $forum_insert$
DECLARE
    nick text;
BEGIN
    SELECT nickname FROM tpdb."User" WHERE nickname = new."user" INTO nick;
    IF nick != '' THEN
        new."user" := nick;
    END IF;
    RETURN new;
END
$forum_insert$ LANGUAGE plpgsql;
CREATE TRIGGER forum_insert BEFORE INSERT ON tpdb."Forum" FOR EACH ROW EXECUTE PROCEDURE forum_insert();

-- THREAD INSTERT
CREATE FUNCTION thread_insert() RETURNS TRIGGER AS $thread_insert$
DECLARE
    nick text;
    forum text;
BEGIN
    SELECT nickname FROM tpdb."User" WHERE nickname = new."author" INTO nick;
    IF nick != '' THEN
        new."author" := nick;
    END IF;
    SELECT slug FROM tpdb."Forum" WHERE slug = new."forum" INTO forum;
    IF forum != '' THEN
        new."forum" := forum;
    END IF;
--     RAISE NOTICE 'inserted: "%", nickname: "%"', new, nick;
    RETURN new;
END
$thread_insert$ LANGUAGE plpgsql;
CREATE TRIGGER thread_insert BEFORE INSERT ON tpdb."Thread" FOR EACH ROW EXECUTE PROCEDURE thread_insert();

-- POST UPDATE
CREATE FUNCTION post_update() RETURNS TRIGGER AS $post_update$
BEGIN
    IF old.message != new.message THEN
        new.isedited := true;
    END IF;
    RETURN new;
END
$post_update$ LANGUAGE plpgsql;
CREATE TRIGGER post_update BEFORE UPDATE ON tpdb."Post" FOR EACH ROW EXECUTE PROCEDURE post_update();

CREATE INDEX IF NOT EXISTS idx_nickname_user ON tpdb."User" USING hash(nickname);
CREATE INDEX IF NOT EXISTS idx_nickname_user ON tpdb."User" USING btree(nickname);
CREATE INDEX IF NOT EXISTS idx_email_user ON tpdb."User" USING hash(email);

CREATE INDEX IF NOT EXISTS idx_slug_forum ON tpdb."Forum" USING hash(slug);

CREATE UNIQUE INDEX IF NOT EXISTS idx_slug_thread ON tpdb."Thread" USING btree(slug) WHERE slug <> '';
CREATE INDEX IF NOT EXISTS idx_forum_thread ON tpdb."Thread" USING hash(forum);
CREATE INDEX IF NOT EXISTS idx_created_thread ON tpdb."Thread" USING btree(created);

CREATE INDEX IF NOT EXISTS idx_thread_post ON tpdb."Post" USING btree(thread);
CREATE INDEX IF NOT EXISTS idx_created_post ON tpdb."Post" USING btree(created);
CREATE INDEX IF NOT EXISTS idx_path_post ON tpdb."Post" USING btree(path);
CREATE INDEX IF NOT EXISTS idx_gin_path_post ON tpdb."Post" USING gin(path);

CREATE INDEX IF NOT EXISTS idx_user_vote ON tpdb."Vote" USING hash("user");
CREATE INDEX IF NOT EXISTS idx_threadid_vote ON tpdb."Vote" USING hash(threadid);



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

CREATE UNLOGGED TABLE User
(
    id SERIAL PRIMARY KEY,
    nickname CITEXT COLLATE "C" UNIQUE NOT NULL,
    fullname TEXT NOT NULL,
    about TEXT,
    email CITEXT UNIQUE NOT NULL
);

CREATE UNLOGGED TABLE Forum
(
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    "user" CITEXT REFERENCES User(nickname) NOT NULL,
    slug CITEXT UNIQUE NOT NULL,
    posts INT DEFAULT 0,
    threads INT DEFAULT 0
);

CREATE UNLOGGED TABLE Thread
(
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    author CITEXT REFERENCES User(nickname) NOT NULL,
    forum CITEXT REFERENCES Forum(slug) NOT NULL,
    message TEXT NOT NULL,
    votes INT DEFAULT 0,
    slug CITEXT,
    created TIMESTAMP WITH TIME ZONE
);

CREATE UNLOGGED TABLE Post
(
    id SERIAL PRIMARY KEY,
    parent INT DEFAULT 0,
    author CITEXT REFERENCES User(nickname) NOT NULL,
    message TEXT NOT NULL,
    isEdited bool NOT NULL DEFAULT FALSE,
    forum CITEXT REFERENCES Forum(slug) NOT NULL,
    thread INT REFERENCES Thread(id) NOT NULL,
    created  TIMESTAMP WITH TIME ZONE,
--     path INT[] DEFAULT [0]
    path INT[] NOT NULL
);

CREATE UNLOGGED TABLE Vote
(
    id SERIAL PRIMARY KEY,
    threadid INT REFERENCES Thread(id) NOT NULL,
--     threadslug CITEXT REFERENCES Thread(Slug) NOT NULL,
    "user" CITEXT REFERENCES User(nickname) NOT NULL,
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
CREATE TRIGGER user_inc AFTER INSERT ON "User"EXECUTE PROCEDURE service_user_inc();

-- SERVICE THREADS
CREATE FUNCTION service_thread_inc() RETURNS TRIGGER AS $thread_inc$
BEGIN
    UPDATE "Service" SET threads = threads + 1 WHERE id = 1;
    RETURN NULL;
END;
$thread_inc$ LANGUAGE plpgsql;
CREATE TRIGGER thread_inc AFTER INSERT ON Thread EXECUTE PROCEDURE service_thread_inc();

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
CREATE TRIGGER post_inc AFTER INSERT ON Post
    REFERENCING NEW TABLE as newtbl
    FOR EACH STATEMENT EXECUTE PROCEDURE service_post_inc();

-- SERVICE FORUMS
CREATE FUNCTION service_forum_inc() RETURNS TRIGGER AS $forum_inc$
BEGIN
    UPDATE "Service" SET forums = forums + 1 WHERE id = 1;
    RETURN NULL;
END;
$forum_inc$ LANGUAGE plpgsql;
CREATE TRIGGER forum_inc AFTER INSERT ON Forum EXECUTE PROCEDURE service_forum_inc();

-- POST INSERT
CREATE FUNCTION post_insert() RETURNS TRIGGER AS $post_insert$
DECLARE
    prevpath INT[];
BEGIN
--     PATH
    IF new.parent != 0 THEN
        SELECT path FROM Post WHERE id = new.parent INTO prevpath;
        new.path := array_append(prevpath, new.id);
    ELSE
        new.path[1] := new.id;
    END IF;
    RETURN new;
END
$post_insert$ LANGUAGE plpgsql;
CREATE TRIGGER post_insert BEFORE INSERT ON Post FOR EACH ROW EXECUTE PROCEDURE post_insert();

-- POST INSERT FORUM
CREATE FUNCTION post_insert_forum() RETURNS TRIGGER AS $post_insert_forum$
BEGIN
    UPDATE Forum SET posts = posts + 1 WHERE slug = new.forum;
    RETURN new;
END
$post_insert_forum$ LANGUAGE plpgsql;
CREATE TRIGGER post_insert_forum AFTER INSERT ON Post FOR EACH ROW EXECUTE PROCEDURE post_insert_forum();

-- THREAD INSERT FORUM
CREATE FUNCTION thread_insert_forum() RETURNS TRIGGER AS $thread_insert_forum$
BEGIN
    UPDATE Forum SET threads = threads + 1 WHERE slug = new.forum;
    RETURN new;
END
$thread_insert_forum$ LANGUAGE plpgsql;
CREATE TRIGGER thread_insert_forum AFTER INSERT ON Thread FOR EACH ROW EXECUTE PROCEDURE thread_insert_forum();

-- FORUM INSTERT
CREATE FUNCTION forum_insert() RETURNS TRIGGER AS $forum_insert$
DECLARE
    nick text;
BEGIN
    SELECT nickname FROM "User"WHERE nickname = new."user" INTO nick;
    IF nick != '' THEN
        new."user" := nick;
    END IF;
    RETURN new;
END
$forum_insert$ LANGUAGE plpgsql;
CREATE TRIGGER forum_insert BEFORE INSERT ON Forum FOR EACH ROW EXECUTE PROCEDURE forum_insert();

-- THREAD INSTERT
CREATE FUNCTION thread_insert() RETURNS TRIGGER AS $thread_insert$
DECLARE
    nick text;
    forum text;
BEGIN
    SELECT nickname FROM "User"WHERE nickname = new."author" INTO nick;
    IF nick != '' THEN
        new."author" := nick;
    END IF;
    SELECT slug FROM Forum WHERE slug = new."forum" INTO forum;
    IF forum != '' THEN
        new."forum" := forum;
    END IF;
--     RAISE NOTICE 'inserted: "%", nickname: "%"', new, nick;
    RETURN new;
END
$thread_insert$ LANGUAGE plpgsql;
CREATE TRIGGER thread_insert BEFORE INSERT ON Thread FOR EACH ROW EXECUTE PROCEDURE thread_insert();

-- POST UPDATE
CREATE FUNCTION post_update() RETURNS TRIGGER AS $post_update$
BEGIN
    IF old.message != new.message THEN
        new.isedited := true;
    END IF;
    RETURN new;
END
$post_update$ LANGUAGE plpgsql;
CREATE TRIGGER post_update BEFORE UPDATE ON Post FOR EACH ROW EXECUTE PROCEDURE post_update();

CREATE INDEX IF NOT EXISTS idx_nickname_user ON "User"USING hash(nickname);
CREATE INDEX IF NOT EXISTS idx_nickname_user ON "User"USING btree(nickname);
CREATE INDEX IF NOT EXISTS idx_email_user ON "User"USING hash(email);

CREATE INDEX IF NOT EXISTS idx_slug_forum ON Forum USING hash(slug);

CREATE UNIQUE INDEX IF NOT EXISTS idx_slug_thread ON Thread USING btree(slug) WHERE slug <> '';
CREATE INDEX IF NOT EXISTS idx_forum_thread ON Thread USING hash(forum);
CREATE INDEX IF NOT EXISTS idx_created_thread ON Thread USING btree(created);

CREATE INDEX IF NOT EXISTS idx_thread_post ON Post USING btree(thread);
CREATE INDEX IF NOT EXISTS idx_created_post ON Post USING btree(created);
CREATE INDEX IF NOT EXISTS idx_path_post ON Post USING btree(path);
CREATE INDEX IF NOT EXISTS idx_gin_path_post ON Post USING gin(path);

CREATE INDEX IF NOT EXISTS idx_user_vote ON Vote USING hash("user");
CREATE INDEX IF NOT EXISTS idx_threadid_vote ON Vote USING hash(threadid);

-- -- VOTE INSTERT
-- CREATE FUNCTION vote_insert() RETURNS TRIGGER AS $vote_insert$
--     DECLARE
-- BEGIN
--     IF new.threadid IS NULL THEN
--         IF new.vote == 1 THEN
--             UPDATE Thread SET votes = votes + 1 WHERE slug = new.threadslug RETURNING id INTO new.threadid;
--         ELSE
--             UPDATE Thread SET votes = votes - 1 WHERE slug = new.threadslug RETURNING id INTO new.threadid;
--         END IF;
--     ELSE
--         IF new.vote == 1 THEN
--             UPDATE Thread SET votes = votes + 1 WHERE id = new.threadid RETURNING slug INTO new.threadslug;
--         ELSE
--             UPDATE Thread SET votes = votes - 1 WHERE id = new.threadid RETURNING slug INTO new.threadslug;
--         END IF;
--     END IF;
--     RETURN new;
-- END
-- $vote_insert$ LANGUAGE plpgsql;
-- CREATE TRIGGER vote_insert AFTER INSERT ON Vote FOR EACH ROW EXECUTE PROCEDURE vote_insert();

-- -- VOTE UPDATE
-- CREATE FUNCTION vote_update() RETURNS TRIGGER AS $vote_update$
-- BEGIN
--     IF old.vote != new.vote THEN
--         IF new.vote = 1 THEN
--             UPDATE Thread SET votes = votes + 2 WHERE id = new.threadid;
--         ELSE
--             UPDATE Thread SET votes = votes - 2 WHERE id = new.threadid;
--         END IF;
--     END IF;
--     RETURN new;
-- END
-- $vote_update$ LANGUAGE plpgsql;
-- CREATE TRIGGER vote_update AFTER UPDATE ON Vote FOR EACH ROW EXECUTE PROCEDURE vote_update();
